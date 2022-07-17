package wrapper

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	"nhooyr.io/websocket"
)

// signal represents a manager Context Signal.
type signal string

// manager Context Signals.
const (
	// keySignal represents the Context key for a manager's signals.
	keySignal = signal("signal")

	// keyReason represents the Context key for a disconnection reason.
	keyReason = signal("reason")

	// signalDisconnect indicates that a Disconnection was called manually.
	signalDisconnect = 1

	// signalReconnect signals the manager to reconnect upon a successful disconnection.
	signalReconnect = 2
)

// manager represents a manager of a Session's goroutines.
type manager struct {
	// cancel represents the cancellation signal for a Session's Context.
	cancel context.CancelFunc

	// routines represents a goroutine counter that ensures all of the Session's goroutines
	// are spawned prior to returning from connect().
	routines sync.WaitGroup

	// err represents the error that this manager detected upon the closing of a Session's goroutines.e
	err chan error

	// errgroup ensures that all of the Session's goroutines are closed prior to returning
	// from Disconnect().
	//
	// IMPLEMENTATION
	// A Session's Context is cancelled to indicate a disconnection:
	// 1. Context is canceled (via function call or error).
	// 2. Goroutines read s.Context.Done() and close accordingly.
	// 3. errgroup.Wait() is called to block until all goroutines are closed.
	// 4. errgroup.Wait() result is returned once all goroutines are closed.
	//
	// As a result of 3, disconnection must NEVER occur on a Session's goroutine.
	// Otherwise, errorgroup.Wait() blocks the goroutine it's waiting on to be closed.
	//
	// ERRGROUP
	// errgroup manages a Session's goroutines: listen, heartbeat, pulse, respond.
	//
	// Upon connection, an (unmanaged) manager goroutine is used to monitor errgroup.Wait().
	//
	// When a disconnection is called purposefully, s.Conn and s.Context is closed.
	// This results in the eventual closing of a Session's goroutines.
	// When errgroup.Wait() returns nil, it indicates a successful disconnection.
	// Otherwise, a DisconnectErr will be returned.
	//
	// When an error occurs in a Session's goroutines, errgroup cancels the Session's context.
	// This results in the eventual closing of a Session's goroutines.
	// When errgroup.Wait() returns err (origin error), the state of the disconnection is managed
	// (since s.Conn may or may not need closing).
	// When managing the state of disconnection is successful, the manager routine returns err.
	// Otherwise, a DisconnectErr (which includes err) will be returned.
	//
	// The above indicates that disconnect() manages the STATE of disconnection, rather than performing
	// the ACTION of disconnection.
	//
	// This implementation allows a caller of disconnect() to use its return value to await disconnection.
	// For example, a channel can be used to receive the value that the manager routine sends.
	// Disconnect() is modified in this way to allow the end-user (developer) to only return from Disconnect()
	// when disconnection is completed.
	*errgroup.Group
}

// manage manages a Session's goroutines.
func (s *Session) manage(bot *Client) {
	s.manager.routines.Done()

	s.Lock()
	s.manager.err = make(chan error, 1)
	s.Unlock()

	// wait until all of a Session's goroutines are closed.
	if err := s.manager.Wait(); err != nil {
		if reason := s.Context.Value(keyReason); reason != nil {
			log.Println(reason)
		}

		if signal := s.Context.Value(keySignal); signal == signalDisconnect || signal == signalReconnect {
			log.Printf("Session %q purposely disconnected ungracefully", s.ID)
		}

		// when an error caused goroutines to close, manage the state of disconnection.
		disconnectErr := new(DisconnectError)
		closeErr := new(websocket.CloseError)
		switch {
		// when an error occurs from a purposeful disconnection.
		case errors.As(err, disconnectErr):
			s.manager.err <- err

		// when an error occurs from a WebSocket Close Error.
		case errors.As(err, closeErr):
			s.manager.err <- DisconnectError{
				SessionID: s.ID,
				Err:       err,
				Action:    nil,
			}

		default:
			if cErr := s.Conn.Close(websocket.StatusCode(FlagClientCloseEventCodeAway), ""); cErr != nil {
				s.manager.err <- DisconnectError{
					SessionID: s.ID,
					Err:       cErr,
					Action:    err,
				}

				return
			}

			s.manager.err <- err
		}

		return
	}

	// When a disconnection is called purposefully and is successful.
	if reason := s.Context.Value(keyReason); reason != nil {
		log.Println(reason)
	}

	signal := s.Context.Value(keySignal)
	switch signal {
	case signalDisconnect:
		log.Printf("successfully disconnected Session %q", s.ID)

	case signalReconnect:
		log.Printf("successfully disconnected Session %q (while reconnecting)", s.ID)

		// allow Discord to close the session.
		<-time.After(time.Second)

		s.Lock()

		// connect to the Discord Gateway again.
		if err := s.connect(bot); err != nil {
			s.manager.err <- fmt.Errorf("an error occurred while reconnecting Session %q: %w", s.ID, err)

			return
		}

		s.Unlock()

	default:
	}

	s.manager.err <- nil
}
