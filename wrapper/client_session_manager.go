package wrapper

import "sync"

const (
	// SessionManagerVoiceKeyUnknownSession represents an unknown Session ID key for SessionManager.Voice
	SessionManagerVoiceKeyUnknownSession = "UNKNOWN"
)

// SessionManager manages sessions.
type SessionManager struct {
	// Gateway represents a map of Discord Gateway Session IDs to Sessions.
	//
	// map[SessionID]Session (map[string]*Session)
	Gateway *sync.Map

	// Voice represents a map of Discord Gateway Session IDs
	// to a set (map) of Discord GuildIDs to Discord Voice Channel Connections.
	//
	// Use the SessionID Key `SessionManagerVoiceKeyUnknownSession` to find a Discord Voice Channel Connection
	// using a GuildID when the SessionID is unknown.
	//
	// map[SessionID]map[GuildID]*VoiceChannelConnection (map[string]map[string]*VoiceChannelConnection)
	Voice *sync.Map
}

// NewSessionManager creates a new SessionManager.
func NewSessionManager() *SessionManager {
	return &SessionManager{
		Gateway: new(sync.Map),
		Voice:   new(sync.Map),
	}
}

// RemoveGatewaySession removes a Gateway Session and its Voice Channel Connections from the Session Manager.
func (sm *SessionManager) RemoveGatewaySession(id string) {
	// remove the mapped Gateway Session ID.
	sm.Gateway.Delete(id)

	// remove the mapped Voice Channel Connections.
	//
	// v = map[GuildID]*VoiceChannelConnection
	if v, ok := sm.Voice.Load(id); ok {
		knownSessionIDMap := v.(*sync.Map)

		// remove the mapped Voice Channel Connections with an unknown Session ID.
		//
		// u = map[GuildID]*VoiceChannelConnection
		if u, ok := sm.Voice.Load(SessionManagerVoiceKeyUnknownSession); ok {
			unknownSessionIDMap := u.(*sync.Map)

			knownSessionIDMap.Range(func(key, value any) bool {
				// key = Guild ID
				guildID := key.(string)
				unknownSessionIDMap.Delete(guildID)

				return true
			})
		}

		knownSessionIDMap.Clear()
	}

	sm.Voice.Delete(id)
}

// StoreVoiceChannelConnection stores a Voice Channel Connection.
func (sm *SessionManager) StoreVoiceChannelConnection(sessionid string, guildid string, vc *VoiceChannelConnection) {
LOADMAP:
	// v = map[GuildID]*VoiceChannelConnection
	if v, ok := sm.Voice.Load(sessionid); ok {
		guildIDvoiceChannelConnectionMap := v.(*sync.Map)
		guildIDvoiceChannelConnectionMap.Store(vc.State.GuildID, vc)
	} else {
		// Store the Gateway Session ID into the bot's Session Manager.
		sm.Voice.Store(sessionid, new(sync.Map))

		goto LOADMAP
	}

	if sessionid != SessionManagerVoiceKeyUnknownSession {
		sm.StoreVoiceChannelConnection(SessionManagerVoiceKeyUnknownSession, guildid, vc)
	}
}

// GetVoiceChannelConnection gets a Voice Channel Connection using a given Guild ID.
func (sm *SessionManager) GetVoiceChannelConnection(sessionid string, guildid string) *VoiceChannelConnection {
	// v = map[GuildID]*VoiceChannelConnection
	if v, ok := sm.Voice.Load(sessionid); ok {
		guildIDvoiceChannelConnectionMap := v.(*sync.Map)

		// v2 = *VoiceChannelConnection
		if v2, ok := guildIDvoiceChannelConnectionMap.Load(guildid); ok {
			return v2.(*VoiceChannelConnection)
		}
	}

	return nil
}
