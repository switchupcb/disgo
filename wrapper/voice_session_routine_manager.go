package wrapper

import (
	"context"
	"errors"
	"fmt"
	"sync"

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

	// when an error caused goroutines to close, manage the state of disconnection.
	if err != nil {
		disconnectErr := new(ErrorSessionDisconnect)
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
				s.manager.err <- ErrorSessionDisconnect{
					Action: err,
					Err:    cErr,
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
	code, ok := VoiceCloseEventCodes[int(closeErr.Code)]
	switch ok {
	// Voice Close Event Code is known.
	case true:
		LogSession(Logger.Info(), s.ID).
			Msgf("received Voice Close Event Code %d %s: %s",
				code.Code, code.Description, code.Explanation,
			)

		return closeErr

	// Voice Close Event Code is unknown.
	default:

		// when another goroutine calls disconnect(),
		// s.Conn.Close is called before s.cancel which will result in
		// a CloseError with the close code that Disgo uses to reconnect.
		if closeErr.Code == websocket.StatusCode(FlagClientCloseEventCodeReconnect) {
			return nil
		}

		LogSession(Logger.Info(), s.ID).
			Msgf("received unknown Voice Close Event Code %d with reason %q",
				closeErr.Code, closeErr.Reason,
			)

		return closeErr
	}
}
