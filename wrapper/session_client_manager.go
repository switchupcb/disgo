package wrapper

import "sync"

// SessionManager manages sessions.
type SessionManager struct {
	// Gateway represents a map of Discord Gateway (TCP WebSocket Connections) session IDs to Sessions.
	// map[ID]Session (map[string]*Session)
	Gateway *sync.Map

	// Voice represents a map of Discord Voice (UDP WebSocket Connection) session IDs to Sessions.
	// map[ID]Session (map[string]*Session)
	Voice *sync.Map
}

// NewSessionManager creates a new SessionManager.
func NewSessionManager() *SessionManager {
	return &SessionManager{
		Gateway: new(sync.Map),
		Voice:   new(sync.Map),
	}
}
