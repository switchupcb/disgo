package wrapper

// SessionManager manages sessions.
type SessionManager struct {
	// All contains a pointer to every session that is being managed.
	All []*Session

	// Gateway represents a map of Discord Gateway (TCP WebSocket Connections) session IDs to Sessions.
	// map[ID]Session
	Gateway map[string]*Session

	// Voice represents a map of Discord Voice (UDP WebSocket Connection) session IDs to Sessions.
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
