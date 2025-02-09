package wrapper

import "sync"

// SessionManager manages sessions.
type SessionManager struct {
	// Gateway represents a map of Discord Gateway Session IDs to Sessions.
	//
	// map[SessionID]Session (map[string]*Session)
	Gateway *sync.Map

	// Voice represents a map of Discord Gateway Session IDs to
	// a set of Discord Guild IDs to Discord Voice Connections.
	//
	// map[SessionID]map[GuildID]VoiceConnection (map[string]map[string]*VoiceConnection)
	Voice *sync.Map
}

// NewSessionManager creates a new SessionManager.
func NewSessionManager() *SessionManager {
	return &SessionManager{
		Gateway: new(sync.Map),
		Voice:   new(sync.Map),
	}
}
