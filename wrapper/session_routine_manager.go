package wrapper

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/switchupcb/websocket"
	"golang.org/x/sync/errgroup"
)

// manager represents a manager of a Session's goroutines.
type manager struct {
	// signals represents a channel of signals.
	signals chan uint8

	// coroner represents a goroutine group to track the manager routine.
	coroner *errgroup.Group

	// context is used as a context for the manager routine.
	context context.Context

	// routines represents a goroutine counter that ensures all of the Session's goroutines
	// are spawned prior to returning from connect().
	routines sync.WaitGroup

	// pulses represents the amount of goroutines that can generate heartbeat pulses.
	//
	// pulses ensures that pulse goroutines always have a receiver channel for heartbeats
	// by preventing the heartbeat goroutine from closing before other pulse goroutines.
	pulses int32

	// errgroup ensures all of the Session's goroutines are closed prior to returning
	// from Disconnect().
	//
	// IMPLEMENTATION
	// A session is managed by multiple goroutine groups.
	//
	// 1. The "coroner" is an unmanaged routine which reports an error from the manager when the Session is no longer active.
	//   a. The "coroner" is responsible for returning an error to a session Disconnect() call.
	//   b. The "coroner" is not responsible for returning an error to a session Connect() or Reconnect() call,
	//     because the coroner only receives an error from a manager when the manager is dead.
	//
	// 2. The "manager" is a tracked routine (by the coroner) which manages the state of the connection to Discord.
	//   a. The "manager" is responsible for returning an error to a session Connect() or Reconnect() call,
	//     because the manager cannot shutdown after these calls are made (in comparison to a final Disconnect()).
	//   b. The manager manages a Session's goroutines: listen, heartbeat, pulse, respond.
	//
	// USING CONTEXT CANCELLATION (to deactivate the session):
	// 1. Context is cancelled (via function call or error in a goroutine).
	// 2. Goroutines read s.Context.Done() and close accordingly.
	// 3. errgroup.Wait() is called from a manager to block until all goroutines are closed.
	// 4. errgroup.Wait() is called from a coroner to return a result once manager goroutine is dead.
	//
	//
	// USING ERRGROUPS (to close goroutines).
	// s.Conn and s.Context is closed when a disconnection is called purposefully.
	//   - This results in the eventual closing of a Session's goroutines.
	//   - A successful disconnection has occurred when errgroup.Wait() returns nil.
	//   - Otherwise, an error is returned.
	//
	*errgroup.Group
}

// spawnManager spawns a tracked manager.
func (s *Session) spawnManager(bot *Client) {
	s.manager = new(manager)
	s.manager.coroner, s.manager.context = errgroup.WithContext(context.Background())
	s.manager.signals = make(chan uint8)

	// spawn the manager goroutine.
	s.manager.coroner.Go(func() error {
		if err := s.manage(bot); err != nil {
			return fmt.Errorf("manager: %w", err)
		}

		return nil
	})
}

// Session Signals represent manager signals to perform actions to the Session.
const (
	sessionSignalConnect    = 1
	sessionSignalDisconnect = 2
	sessionSignalReconnect  = 3
)

// manage manages a Session's goroutines.
func (s *Session) manage(bot *Client) error { //nolint:maintidx
	// spawn the coroner once the manager routine is alive.
	go s.coroner()

	// create a temporary context for a new session (which is reset upon connection).
	s.Context = context.Background()

	var managerErr error

	defer func() {
		// remove the session from the client.
		s.client_manager.RemoveGatewaySession(s.ID)

		s.logClose("manager")
	}()

	for {
		select {
		// <-s.Context.Done() when all managed routines are closing
		// due to reconnection (while awaiting a connection signal) or
		// due to an unexpected error in a managed routine.
		case <-s.Context.Done():
			LogSession(Logger.Info(), s.ID).Str(LogCtxClient, bot.ApplicationID).Msgf("received signal: <-s.Context.Done with state %q", s.State())

			// wait until the session's manager goroutines are closed (with s.Unlocked).
			//
			// proof: s.manager.Wait() returns instantly when SessionStateDisconnectedReconnect (with s.Locked).
			err := s.manager.Wait()

			// All session routines are closed when
			//
			// reconnecting (while waiting for another signal)
			if s.State() == SessionStateDisconnectedReconnect {
				break
			}

			// disconnecting (unexpectedly)
			if err != nil {
				// TODO: Use errors.As: https://github.com/coder/websocket/issues/519
				if strings.Contains(err.Error(), "failed to close WebSocket: received header with unexpected rsv bits set") {
					return nil
				}

				closeErr := new(websocket.CloseError)
				if errors.As(err, closeErr) {
					if vErr := s.validateGatewayCloseError(closeErr); vErr == nil {
						// reconnect from a state where
						s.setState(SessionStateDisconnectedReconnect)

						// send a connection signal.
						go func() {
							s.manager.signals <- sessionSignalConnect
						}()

						s.Lock()

						break // to reconnect from the connect case logic.
					} // vErr == nil
				} // errors.As
			} // err != nil

			return nil

		case signal := <-s.manager.signals:
			switch signal {
			case sessionSignalConnect:
				LogSession(Logger.Info(), s.ID).Str(LogCtxClient, bot.ApplicationID).Msgf("received signal: connect with state %q", s.State())

				switch s.State() {
				// SessionStateNew when Connect() on new session.
				case SessionStateNew:
					s.Lock()

					LogSession(Logger.Info(), s.ID).Str(LogCtxClient, bot.ApplicationID).Msg("connecting session")

					if err := s.connect(bot); err != nil {
						managerErr = ErrorSession{SessionID: s.ID, State: s.State(), Type: ErrorSessionTypeGateway, Err: err}

						// disconnect when error occurred after websocket connection
						if s.State() == SessionStateConnectingWebsocket {
							// send a disconnection signal.
							go func() {
								s.manager.signals <- sessionSignalDisconnect
							}()

							break // to handle error after disconnection
						}

						s.Unlock()

						return managerErr
					}

					s.setState(SessionStateConnected)
					s.Unlock()

				// SessionStateDisconnectedReconnect when reconnecting from disconnected session.
				case SessionStateDisconnectedReconnect:
					// s.Lock() called during reconnection signal.

					LogSession(Logger.Info(), s.ID).Str(LogCtxClient, bot.ApplicationID).Msg("reconnecting session")

					if err := s.connect(bot); err != nil {
						managerErr = ErrorSession{SessionID: s.ID, State: s.State(), Type: ErrorSessionTypeGateway, Err: err}

						// disconnect when error occurred after websocket connection
						if s.State() == SessionStateConnectingWebsocket {
							// send a disconnection signal.
							go func() {
								s.manager.signals <- sessionSignalDisconnect
							}()

							break // to handle error after disconnection
						}

						s.Unlock()

						return managerErr
					}

					LogSession(Logger.Info(), s.ID).Str(LogCtxClient, bot.ApplicationID).Msg("connected session")
					s.setState(SessionStateConnected)
					s.Unlock()

				default:
					return fmt.Errorf("unexpected state during session connection: %v", s.State())
				}

			case sessionSignalDisconnect:
				LogSession(Logger.Info(), s.ID).Str(LogCtxClient, bot.ApplicationID).Msgf("received signal: disconnect with state %q", s.State())

				// update the session's state and client close event code.
				var code int

				switch {
				case managerErr != nil:
					// s.Lock() called before error.

					s.setState(SessionStateDisconnectingError)
					code = FlagClientCloseEventCodeNormal

					LogSession(Logger.Info(), s.ID).Str(LogCtxClient, bot.ApplicationID).Msgf("%q session with code %d", s.State(), code)

				case s.State() == SessionStateReconnecting:
					// s.Lock() called during reconnection signal.

					s.setState(SessionStateDisconnectingReconnect)
					code = FlagClientCloseEventCodeReconnect

					LogSession(Logger.Info(), s.ID).Str(LogCtxClient, bot.ApplicationID).Msgf("%q session with code %d", s.State(), code)

				default:
					s.Lock()

					s.setState(SessionStateDisconnecting)
					code = FlagClientCloseEventCodeNormal

					LogSession(Logger.Info(), s.ID).Str(LogCtxClient, bot.ApplicationID).Msgf("%q session with code %d", s.State(), code)
				}

				// disconnect the session.
				if err := s.disconnect(code); err != nil {
					// validate the disconnection error.
					closeErr := new(websocket.CloseError)

					// TODO: Use errors.As: https://github.com/coder/websocket/issues/519
					if strings.Contains(err.Error(), "failed to close WebSocket: received header with unexpected rsv bits set") {
						err = nil
					} else if errors.As(err, closeErr) {
						err = s.validateGatewayCloseError(closeErr)
					}

					if managerErr != nil {
						s.Unlock()

						// wait until the session's manager goroutines are closed (with s.Unlocked).
						_ = s.manager.Wait()

						return ErrorSession{
							SessionID: s.ID,
							State:     s.State(),
							Type:      ErrorSessionTypeGateway,
							Err: ErrorSessionDisconnect{
								Action: managerErr,
								Err:    err,
							},
						}
					}

					if err != nil {
						s.Unlock()

						// wait until the session's manager goroutines are closed (with s.Unlocked).
						_ = s.manager.Wait()

						return ErrorSession{
							SessionID: s.ID,
							State:     s.State(),
							Type:      ErrorSessionTypeGateway,
							Err: ErrorSessionDisconnect{
								Action: nil,
								Err:    err,
							},
						}
					}
				} // disconnect

				// update the session's state.
				switch {
				case managerErr != nil:
					s.setState(SessionStateDisconnectedError)
				case s.State() == SessionStateDisconnectingReconnect:
					s.setState(SessionStateDisconnectedReconnect)
				default:
					s.setState(SessionStateDisconnectedFinal)
				}

				// wait until the session's manager goroutines are closed (with s.Unlocked).
				s.Unlock()
				_ = s.manager.Wait()

				s.Lock()
				LogSession(Logger.Info(), s.ID).Str(LogCtxClient, bot.ApplicationID).Msgf("%q session with code %d", s.State(), code)

				if s.State() == SessionStateDisconnectedReconnect {
					// allow Discord to close the session.
					<-time.After(time.Second)

					// send a connection signal.
					go func() {
						s.manager.signals <- sessionSignalConnect
					}()

					break
				}

				s.Unlock()

				return nil

			case sessionSignalReconnect:
				LogSession(Logger.Info(), s.ID).Str(LogCtxClient, bot.ApplicationID).Msgf("received signal: reconnect with state %q", s.State())

				s.Lock()
				s.setState(SessionStateReconnecting)

				// send a disconnection signal.
				go func() {
					s.manager.signals <- sessionSignalDisconnect
				}()
			} // switch signal
		} // select
	} // for
}

// logClose safely logs the close of a Session's goroutine.
func (s *Session) logClose(routine string) {
	LogSession(Logger.Info(), s.ID).Msgf("closed %s routine", routine)
}

// validateGatewayCloseError validates a WebSocket CloseError
// and returns whether to reconnect (when error == nil).
func (s *Session) validateGatewayCloseError(closeErr *websocket.CloseError) error {
	code, ok := GatewayCloseEventCodes[int(closeErr.Code)]

	switch ok {
	// Gateway Close Event Code is known.
	case true:
		LogSession(Logger.Info(), s.ID).
			Msgf("received Gateway Close Event Code %d %s: %s",
				code.Code, code.Description, code.Explanation,
			)

		if code.Reconnect {
			return nil
		}

		return closeErr

	// Gateway Close Event Code is unknown.
	default:
		// when another goroutine returns an error,
		// s.Conn.Close is called before s.cancel which will result in
		// a CloseError with the close code that Disgo uses to reconnect.
		if closeErr.Code == websocket.StatusCode(FlagClientCloseEventCodeNormal) ||
			closeErr.Code == websocket.StatusCode(FlagClientCloseEventCodeReconnect) {
			return nil
		}

		LogSession(Logger.Info(), s.ID).
			Msgf("received unknown Gateway Close Event Code %d with reason %q",
				closeErr.Code, closeErr.Reason,
			)

		return closeErr
	}
}
