package wrapper

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/switchupcb/websocket"
	"golang.org/x/sync/errgroup"
)

// voice_manager represents a manager of a Voice Session's goroutines.
type voice_manager struct {
	// routines represents a goroutine counter that ensures all of the Voice Session's goroutines
	// are spawned prior to returning from connect().
	routines sync.WaitGroup

	// cancel represents the cancellation signal for a Voice Session's Context.
	cancel context.CancelFunc

	// signal represents the Context Signal for a Voice Session upon disconnection.
	signal context.Context

	// err represents the error that this manager detected upon the closing of a Voice Session's goroutines.
	err chan error

	// pulses represents the amount of goroutines that can generate heartbeat pulses.
	//
	// pulses ensures that pulse goroutines always have a receiver channel for heartbeats
	// by preventing the heartbeat goroutine from closing before other pulse goroutines.
	pulses int32

	// errgroup ensures that all of the Voice Session's goroutines are closed prior to returning
	// from disconnect().
	//
	// IMPLEMENTATION
	// A Voice Session's Context is cancelled to indicate a disconnection:
	// 1. Context is canceled (via function call or error).
	// 2. Goroutines read s.Context.Done() and close accordingly.
	// 3. errgroup.Wait() is called to block until all goroutines are closed.
	// 4. errgroup.Wait() result is returned once all goroutines are closed.
	//
	// As a result of 3, disconnection must NEVER occur on a Voice Session's goroutine.
	// Otherwise, errorgroup.Wait() blocks the goroutine it's waiting on to be closed.
	// In other words, disconnection MUST occur on another goroutine.
	//
	// ERRGROUP
	// errgroup manages a Voice Session's goroutines: listen, heartbeat, pulse, respond.
	//
	// Upon connection, an (unmanaged) manager goroutine is used to monitor errgroup.Wait().
	//
	// When a disconnection is called purposefully, s.Conn and s.Context is closed.
	// This results in the eventual closing of a Voice Session's goroutines.
	// When errgroup.Wait() returns nil, it indicates a successful disconnection.
	// Otherwise, a DisconnectError will be returned.
	//
	// When an error occurs in a Voice Session's goroutines, errgroup cancels the Voice Session's context.
	// This results in the eventual closing of a Voice Session's goroutines.
	// When errgroup.Wait() returns err (origin error), the state of the disconnection is managed
	// (since s.Conn may or may not need closing).
	// When managing the state of disconnection is successful, the manager routine returns err.
	// Otherwise, a DisconnectError (which includes err) will be returned.
	//
	// The above indicates that the manager manages the STATE of disconnection, while disconnect()
	// performs the ACTION of disconnection.
	//
	// This implementation allows a caller of disconnect() to use its return value to await disconnection.
	// For example, a channel can be used to receive the value that the manager routine sends.
	// VoiceConnection.Disconnect() is modified in this way to allow the end-user (developer) to only return from Disconnect()
	// when disconnection is fully completed (with goroutines closed).
	*errgroup.Group
}

// decrementPulses safely decrements the pulses counter of a Voice Session manager.
func (s *VoiceSession) decrementPulses() {
	s.Lock()
	defer s.Unlock()

	atomic.AddInt32(&s.manager.pulses, -1)
}

// logClose safely logs the close of a Voice Session's goroutine.
func (s *VoiceSession) logClose(routine string) {
	LogSession(Logger.Info(), s.ID).Msgf("closed %s routine", routine)
}

// reconnect spawns a goroutine for reconnection which prompts the manager
// to reconnect upon a disconnection.
func (s *VoiceSession) reconnect(reason string) {
	s.manager.Go(func() error {
		s.Lock()
		defer s.logClose("reconnect")
		defer s.Unlock()

		LogSession(Logger.Info(), s.ID).Msg(reason)

		s.manager.signal = context.WithValue(s.manager.signal, keySignal, signalReconnect)
		if err := s.disconnect(FlagClientCloseEventCodeReconnect); err != nil {
			return fmt.Errorf("reconnect: %w", err)
		}

		return nil
	})
}

// manage manages a Voice Session's goroutines.
func (s *VoiceSession) manage() {
	s.manager.routines.Done()
	defer func() {
		s.Lock()
		s.logClose("manager")
		s.Unlock()
	}()

	// wait until all of a Voice Session's goroutines are closed.
	err := s.manager.Wait()
	s.Lock()
	defer s.Unlock()

	// log the reason for disconnection (if applicable).
	if reason := s.manager.signal.Value(keyReason); reason != nil {
		LogSession(Logger.Info(), s.ID).Msgf("%v", reason)
	}

	// when a signal is provided, it indicates that the disconnection was purposeful.
	signal := s.manager.signal.Value(keySignal)
	switch signal {
	case signalDisconnect:
		LogSession(Logger.Info(), s.ID).Msg("successfully disconnected")

		s.manager.err <- nil

		return

	case signalReconnect:
		LogSession(Logger.Info(), s.ID).Msg("successfully disconnected (while reconnecting)")

		// allow Discord to close the session.
		<-time.After(time.Second)

		s.manager.err <- nil

		return
	}

	// when an error caused goroutines to close, manage the state of disconnection.
	if err != nil {
		disconnectErr := new(ErrorDisconnect)
		closeErr := new(websocket.CloseError)
		switch {
		// when an error occurs from a purposeful disconnection.
		case errors.As(err, disconnectErr):
			s.manager.err <- err

		// when an error occurs from a WebSocket Close Error.
		case errors.As(err, closeErr):
			s.manager.err <- s.handleGatewayCloseError(closeErr)

		default:
			if cErr := s.Conn.Close(websocket.StatusCode(FlagClientCloseEventCodeAway), ""); cErr != nil {
				s.manager.err <- ErrorDisconnect{
					Action:     err,
					Err:        cErr,
					Connection: ErrConnectionSession,
				}

				return
			}

			s.manager.err <- err
		}

		return
	}

	s.manager.err <- nil
}

// handleGatewayCloseError handles a WebSocket CloseError.
func (s *VoiceSession) handleGatewayCloseError(closeErr *websocket.CloseError) error {
	code, ok := GatewayCloseEventCodes[int(closeErr.Code)]
	switch ok {
	// Gateway Close Event Code is known.
	case true:
		LogSession(Logger.Info(), s.ID).
			Msgf("received Gateway Close Event Code %d %s: %s",
				code.Code, code.Description, code.Explanation,
			)

		if code.Reconnect {
			s.reconnect(fmt.Sprintf("reconnecting due to Gateway Close Event Code %d", code.Code))

			return nil
		}

		return closeErr

	// Gateway Close Event Code is unknown.
	default:

		// when another goroutine calls disconnect(),
		// s.Conn.Close is called before s.cancel which will result in
		// a CloseError with the close code that Disgo uses to reconnect.
		if closeErr.Code == websocket.StatusCode(FlagClientCloseEventCodeReconnect) {
			return nil
		}

		LogSession(Logger.Info(), s.ID).
			Msgf("received unknown Gateway Close Event Code %d with reason %q",
				closeErr.Code, closeErr.Reason,
			)

		return closeErr
	}
}

// Wait blocks until the calling Voice Session has disconnected, then returns the reason
// (disgo.SignalReason) for disconnecting and the disconnection error (if it exists).
//
// If Wait() is called on a Voice Session that isn't connected, it will return immediately
// with code SignalNone.
//
// It's NOT recommended to modify a Voice Session after it has disconnected,
// since it will be cleared and placed into a memory pool shortly after.
func (s *VoiceSession) Wait() (int, error) {
	if !s.isConnected() {
		return SignalNone, nil
	}

	// NOTE: Wait() is equivalent to the s.manage() s.manager.Wait() handling logic,
	// but without the management of the disconnection state,
	// and without the usage of a channel that tells another goroutine to unblock.
	//
	// wait until all of a Session's goroutines are closed.
	err := s.manager.Wait()
	s.Lock()
	defer s.Unlock()

	// when a signal is provided, it indicates that the disconnection was purposeful.
	signal := s.manager.signal.Value(keySignal)
	switch signal {
	case signalDisconnect:
		return SignalDisconnect, nil

	case signalReconnect:
		return SignalReconnect, nil
	}

	// when an error caused goroutines to close.
	if err != nil {
		disconnectErr := new(ErrorDisconnect)
		closeErr := new(websocket.CloseError)
		switch {
		// when an error occurs from a purposeful disconnection.
		case errors.As(err, disconnectErr):
			if signal != nil {
				if signalValue, ok := signal.(int); ok {
					return signalValue, err //nolint:wrapcheck
				}
			}

			return SignalDisconnectError, err //nolint:wrapcheck

		// when an error occurs from a WebSocket Close Error.
		case errors.As(err, closeErr):
			return SignalError, s.handleGatewayCloseError(closeErr)
		}

		return SignalError, err //nolint:wrapcheck
	}

	return SignalUndefined, nil
}
