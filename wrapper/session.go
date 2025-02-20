package wrapper

import (
	"context"
	"errors"
	"fmt"
	"sync"

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

	// Conn represents a WebSocket Connection to the Discord Gateway.
	Conn *websocket.Conn

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

	if err := <-s.manager.actionError; err != nil {
		return err
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

	if err := <-s.manager.actionError; err != nil {
		return err
	}

	// Reset the session.
	putSession(s)

	return nil
}

// Reconnect reconnects an already connected session to the Discord Gateway
// by disconnecting the session, then connecting again.
func (s *Session) Reconnect(bot *Client) error {
	s.Lock()
	if s.manager == nil || s.State() != SessionStateConnected {
		s.Unlock()

		return errors.New("cannot reconnect session that isn't connected")
	}

	s.manager.signals <- sessionSignalReconnect
	s.Unlock()

	if err := <-s.manager.actionError; err != nil {
		return fmt.Errorf("reconnect: %w", err)
	}

	return nil
}
