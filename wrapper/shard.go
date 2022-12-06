package wrapper

import "time"

// ShardManager represents an interface for Shard Management.
//
// ShardManager is an interface which allows developers to use multi-application architectures,
// which run multiple applications on separate processes or servers.
type ShardManager interface {
	// GetLimit gets the limit of the ShardManager.
	GetLimit() *ShardLimit

	// SetLimit sets the limit of the ShardManager.
	SetLimit(ShardLimit)

	// Connect connects to the Discord Gateway using the Shard Manager.
	Connect(bot *Client)

	// Reconnect connects to the Discord Gateway using the Shard Manager.
	Reconnect(bot *Client)

	// Disconnect disconnects from the Discord Gateway using the Shard Manager.
	Disconnect(bot *Client)

	// ConnectSession connects a session of a bot to the Discord Gateway using the Shard Manager.
	ConnectSession(bot *Client, session *Session)

	// ReconnectSession reconnects a session to the Discord Gateway using the Shard Manager.
	ReconnectSession(bot *Client, session *Session)

	// DisconnectSession disconnects a session from the Discord Gateway using the Shard Manager.
	DisconnectSession(session *Session)

	// Identify determines how a Session identifies to the Discord Gateway.
	//
	// Called from the session.go initial function.
	Identify(bot *Client, session *Session)

	// Ready is called when a Session receives a ready event.
	//
	// Called from the session.go initial function.
	Ready(bot *Client, session *Session, event *Ready)
}

// ShardLimit contains information about sharding limits.
type ShardLimit struct {
	// MaxStarts represents the maximum amount of WebSocket Sessions a bot can start per day.
	//
	// This is equivalent to the maximum amount of Shards a bot can create per day.
	//
	// Discord represents this value from the "total" field of the SessionStartLimit object.
	// https://discord.com/developers/docs/topics/gateway#session-start-limit-object
	MaxStarts int

	// RemainingStarts represents the remaining number of "starts" that the bot is allowed
	// until the reset time.
	//
	// This is equivalent to the remaining number of Shards that the bot can create
	// for the rest of the day.
	//
	// Discord represents this value from the "remaining" field of the SessionStartLimit object.
	// https://discord.com/developers/docs/topics/gateway#session-start-limit-object
	RemainingStarts int

	// Reset represents the time at which the Session Start Rate Limit resets (daily).
	//
	// Discord represents this value from the "reset_after" field of the SessionStartLimit object.
	// https://discord.com/developers/docs/topics/gateway#session-start-limit-object
	Reset time.Time

	// MaxConcurrency represents the number of Identify SendEvents the bot can send every 5 seconds.
	MaxConcurrency int
}

// SessionManager manages sessions.
type SessionManager struct {
	// All contains a pointer to every session that is being managed.
	All []*Session

	// Gateway represents a map of Discord Gateway (TCP WebSocket Connections) Session IDs to Sessions.
	// map[ID]Session
	Gateway map[string]*Session

	// Voice represents a map of Discord Voice (UDP WebSocket Connection) Session IDs to Sessions.
	// map[ID]Session
	Voice map[string]*Session
}

// NewSessionManager creates a new SessionManager.
func NewSessionManager() *SessionManager {
	return &SessionManager{
		All:     []*Session{},
		Gateway: make(map[string]*Session),
		Voice:   make(map[string]*Session),
	}
}

// NewSession creates a managed Session and returns it.
func (sm *SessionManager) NewSession() *Session {
	session := NewSession()
	sm.All = append(sm.All, session)

	return session
}
