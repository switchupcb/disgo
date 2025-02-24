package wrapper

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/switchupcb/websocket"
)

// Session represents a Discord Gateway WebSocket Session.
type Session struct {
	// ID represents the session ID of the Session.
	ID string

	// Seq represents the last sequence number received by the client.
	//
	// https://discord.com/developers/docs/topics/gateway#heartbeat
	Seq int64

	// Endpoint represents the endpoint that is used to reconnect to the Gateway.
	Endpoint string

	// Shard represents the [shard_id, num_shards] for the Session.
	//
	// https://discord.com/developers/docs/topics/gateway#sharding
	Shard *[2]int

	// Context carries request-scoped data for the Discord Gateway Connection.
	//
	// Context is also used as a signal for the Session's goroutines.
	Context context.Context

	// Cancel represents the cancellation signal for a Session Context.
	cancel context.CancelFunc

	// Conn represents a WebSocket Connection to the Discord Gateway.
	Conn *websocket.Conn

	// state represents the state of the Session's connection to Discord.
	state string

	// stateMutex is used to protect the Session's manager state from data races
	// by providing transactional functionality.
	stateMutex sync.RWMutex

	// heartbeat contains the fields required to implement the heartbeat mechanism.
	heartbeat *heartbeat

	// manager represents a manager of a Session's goroutines.
	manager *manager

	// client_manager represents the *Client Session Manager of the Session.
	client_manager *SessionManager

	// RateLimiter represents an object that provides rate limit functionality.
	RateLimiter RateLimiter

	// RWMutex is used to protect the Session's variables from data races
	// by providing transactional functionality.
	sync.RWMutex
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
	s.stateMutex.RLock()
	defer s.stateMutex.RUnlock()

	return s.state
}

// setState sets the state of a Session.
func (s *Session) setState(state string) {
	s.stateMutex.Lock()
	s.state = state
	s.stateMutex.Unlock()
}

// canReconnect returns whether the Session's fields are in a valid state to reconnect.
func (s *Session) canReconnect() bool {
	return s.ID != "" && s.Endpoint != "" && atomic.LoadInt64(&s.Seq) != 0
}

// Connect connects a session to the Discord Gateway (WebSocket Connection).
func (s *Session) Connect(bot *Client) error {
	if bot == nil {
		return errors.New("cannot connect session using a nil Client")
	}

	if bot.Sessions == nil {
		bot.Sessions = NewSessionManager()
	}

	if bot.Handlers == nil {
		bot.Handlers = new(Handlers)
	}

	s.Lock()
	s.client_manager = bot.Sessions

	if s.manager != nil && s.State() == SessionStateConnected {
		s.Unlock()

		return fmt.Errorf("session %q is already connected", s.ID)
	}

	s.spawnManager(bot)

	s.manager.signals <- sessionSignalConnect
	s.Unlock()

	// wait until the Session has connected
	for {
		select {
		// Context is cancelled during connection when the manager returns an error
		// or disconnects from another goroutine call.
		case <-s.manager.context.Done():
			return s.manager.coroner.Wait() //nolint:wrapcheck
		default:
			break
		}

		// Session is SessionStateConnected after connection.
		//
		// proof: Calling Connect() during connection cannot happen while the manager exists.
		if s.State() == SessionStateConnected {
			break
		}
	}

	return nil
}

// Disconnect disconnects a session from the Discord Gateway.
func (s *Session) Disconnect() error {
	s.Lock()
	if s.manager == nil || s.State() != SessionStateConnected {
		s.Unlock()

		return errors.New("cannot disconnect session that isn't connected")
	}

	s.manager.signals <- sessionSignalDisconnect
	s.Unlock()

	// Session is disconnected from a Disconnect() call when the coroner shuts down.
	if err := s.manager.coroner.Wait(); err != nil {
		return err //nolint:wrapcheck
	}

	return nil
}

// Reconnect reconnects an already connected session to the Discord Gateway
// by disconnecting the session, then connecting again.
func (s *Session) Reconnect() error {
	s.Lock()
	if s.manager == nil || s.State() != SessionStateConnected {
		s.Unlock()

		return errors.New("cannot reconnect session that isn't connected")
	}

	s.manager.signals <- sessionSignalReconnect
	s.Unlock()

	// wait until the manager has received the sessionSignalReconnect
	// or changed state to another signal.
	for {
		if s.State() != SessionStateConnected {
			break
		}
	}

	// wait until the Session has reconnected
	// or has experienced an error during reconnection.
	for {
		select {
		// Context is cancelled during reconnection when the manager returns an error
		// or disconnects from another goroutine call.
		case <-s.manager.context.Done():
			return s.manager.coroner.Wait() //nolint:wrapcheck
		default:
			break
		}

		// Session is SessionStateConnected after reconnection.
		//
		// proof: Calling Connect() during reconnection cannot happen while the manager exists.
		if s.State() == SessionStateConnected {
			break
		}
	}

	return nil
}
