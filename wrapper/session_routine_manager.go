package wrapper

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/switchupcb/websocket"
	"golang.org/x/sync/errgroup"
)

// manager represents a manager of a Session's goroutines.
type manager struct {
	// state represents the state of the Session's connection to Discord.
	state string

	// stateMutex is used to protect the Session's manager state from data races
	// by providing transactional functionality.
	stateMutex sync.RWMutex

	// signals represents a channel of signals.
	signals chan uint8

	// coroner represents a goroutine group to track the manager routine.
	coroner errgroup.Group

	// routines represents a goroutine counter that ensures all of the Session's goroutines
	// are spawned prior to returning from connect().
	routines sync.WaitGroup

	// cancel represents the cancellation signal for a Session's Context.
	cancel context.CancelFunc

	// actionError represents the error this manager detects upon a connection action (e.g., connecting, disconnecting).
	actionError chan error

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

	s.Context, s.manager.cancel = context.WithCancel(context.Background())
	s.manager.Group, s.Context = errgroup.WithContext(s.Context)
	s.manager.signals = make(chan uint8)
	s.manager.actionError = make(chan error, 1)

	// spawn the manager goroutine.
	s.manager.coroner.Go(func() error {
		if err := s.manage(bot); err != nil {
			return fmt.Errorf("manager: %w", err)
		}

		return nil
	})
}

// Session States represent the state of the Session's connection to Discord.
const (
	SessionStateNew = ""

	SessionStateConnecting          = "connecting (before websocket connection)"
	SessionStateConnectingWebsocket = "connecting (with websocket connection)"
	SessionStateConnected           = "connected"

	SessionStateDisconnecting          = "disconnecting (purposefully)"
	SessionStateDisconnectingError     = "disconnecting (due to an error)"
	SessionStateDisconnectingReconnect = "disconnecting (while reconnecting)"

	SessionStateDisconnectedFinal     = "disconnected (after connection)"
	SessionStateDisconnectedError     = "disconnected (due to an error)"
	SessionStateDisconnectedReconnect = "disconnected (while reconnecting)"

	SessionStateReconnecting = "reconnecting"
)

// State returns the state of the Session's connection to Discord.
func (s *Session) State() string {
	s.manager.stateMutex.RLock()
	defer s.manager.stateMutex.RUnlock()

	return s.manager.state
}

// setState sets the state of a Session.
func (s *Session) setState(state string) {
	s.manager.stateMutex.Lock()
	s.manager.state = state
	s.manager.stateMutex.Unlock()
}

// canReconnect returns whether the Session's fields are in a valid state to reconnect.
func (s *Session) canReconnect() bool {
	return s.ID != "" && s.Endpoint != "" && atomic.LoadInt64(&s.Seq) != 0
}

// Session Signals represent manager signals to perform actions to the Session.
const (
	sessionSignalConnect    = 1
	sessionSignalDisconnect = 2
	sessionSignalReconnect  = 3
)

// manage manages a Session's goroutines.
func (s *Session) manage(bot *Client) error {
	// spawn the coroner once the manager routine is alive.
	go s.coroner()

	defer func() {
		if s.State() != SessionStateDisconnectedReconnect {
			s.Unlock()
		}

		// wait until the previous connection's manager goroutines are closed.
		_ = s.manager.Wait()

		s.logClose("manager")
	}()

	var managedErr error

	for {
		select {
		case <-s.Context.Done():
			if s.State() == SessionStateDisconnectedReconnect {
				break
			}

			// wait until the previous connection's manager goroutines are closed.
			err := s.manager.Wait()
			if err != nil {
				closeErr := new(websocket.CloseError)

				if errors.As(err, closeErr) {
					if vErr := s.validateGatewayCloseError(closeErr); vErr == nil {
						// reconnect from a state where
						s.setState(SessionStateDisconnectedReconnect)

						// manager routines must be reset
						s.Context, s.manager.cancel = context.WithCancel(context.Background()) //nolint:fatcontext
						s.manager.Group, s.Context = errgroup.WithContext(s.Context)

						go func() {
							// send a connection signal.
							s.manager.signals <- sessionSignalConnect

							// read the s.manager.actionError send from a successful connection.
							e := <-s.manager.actionError
							LogSession(Logger.Info(), s.ID).Str(LogCtxClient, bot.ApplicationID).Msgf("captured result from close event reconnect: %q", e)
						}()

						s.Lock()

						break // to reconnect from the connect case logic.
					} // vErr == nil
				} // errors.As

				// TODO: Use errors.As: https://github.com/coder/websocket/issues/519
				if strings.Contains(err.Error(), "failed to close WebSocket: received header with unexpected rsv bits set") {
					err = nil
				}
			} // err != nil

			s.Lock()

			return err

		case signal := <-s.manager.signals:
			switch signal {
			case sessionSignalConnect:
				if s.State() != SessionStateDisconnectedReconnect {
					s.Lock()
				} else {
					s.Unlock()

					// wait until the previous connection's manager goroutines are closed.
					_ = s.manager.Wait()

					s.Lock()
				}

				LogSession(Logger.Info(), s.ID).Str(LogCtxClient, bot.ApplicationID).Msg("connecting session")

				if err := s.connect(bot); err != nil {
					managedErr = ErrorSession{SessionID: s.ID, State: s.State(), Type: ErrorSessionTypeGateway, Err: err}

					switch s.State() {
					case SessionStateConnectingWebsocket:
						go func() {
							// send a disconnection signal.
							s.manager.signals <- sessionSignalDisconnect
						}()

					// case SessionStateNew, SessionStateConnecting...
					default:
						return managedErr
					}

					break // to handle the error in the disconnect case logic.
				}

				s.setState(SessionStateConnected)
				s.manager.actionError <- nil
				s.Unlock()

			case sessionSignalDisconnect:
				if managedErr == nil && s.State() != SessionStateReconnecting {
					s.Lock()
				}

				// update the session's state and client close event code.
				code := FlagClientCloseEventCodeNormal

				switch {
				case managedErr != nil:
					s.setState(SessionStateDisconnectingError)
				case s.State() == SessionStateReconnecting:
					s.setState(SessionStateDisconnectingReconnect)
					code = FlagClientCloseEventCodeReconnect
				default:
					s.setState(SessionStateDisconnecting)
				}

				LogSession(Logger.Info(), s.ID).Msgf("%q session with code %d", s.State(), code)

				// disconnect the session.
				if err := s.disconnect(code); err != nil {
					managedErr = ErrorSession{
						SessionID: s.ID,
						State:     s.State(),
						Type:      ErrorSessionTypeGateway,
						Err: ErrorSessionDisconnect{
							Action: managedErr,
							Err:    err,
						},
					}

					// TODO: Use errors.As: https://github.com/coder/websocket/issues/519
					if strings.Contains(err.Error(), "failed to close WebSocket: received header with unexpected rsv bits set") {
						managedErr = nil
					}

					if s.State() != SessionStateDisconnectingReconnect {
						return managedErr
					}

					// validate error when reconnecting
					closeErr := new(websocket.CloseError)
					if errors.As(managedErr, closeErr) {
						if managedErr = s.validateGatewayCloseError(closeErr); managedErr != nil {
							return managedErr
						}
					}
				} // disconnect

				// update the session's state.
				switch {
				case s.State() == SessionStateDisconnectingError:
					s.setState(SessionStateDisconnectedError)

				case s.State() == SessionStateDisconnectingReconnect:
					s.setState(SessionStateDisconnectedReconnect)

				default:
					s.setState(SessionStateDisconnectedFinal)
				}

				LogSession(Logger.Info(), s.ID).Msgf("%q session with code %d", s.State(), code)

				if s.State() == SessionStateDisconnectedReconnect {
					// allow Discord to close the session.
					<-time.After(time.Second)

					go func() {
						// send a connection signal.
						s.manager.signals <- sessionSignalConnect
					}()

					break
				}

				// Destroy the manager when the bot isn't reconnecting.
				if managedErr != nil {
					return managedErr
				}

				return nil

			case sessionSignalReconnect:
				s.Lock()
				s.setState(SessionStateReconnecting)

				go func() {
					// send a disconnection signal.
					s.manager.signals <- sessionSignalDisconnect
				}()
			}
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
		LogSession(Logger.Info(), s.ID).
			Msgf("received unknown Gateway Close Event Code %d with reason %q",
				closeErr.Code, closeErr.Reason,
			)

		return closeErr
	}
}
