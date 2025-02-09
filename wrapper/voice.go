package wrapper

import (
	"fmt"
	"net"
	"sync"
)

// VoiceConnection represents a Discord Voice Channel Connection.
//
// A Discord Voice Channel Connection is composed of three connections:
//
//  1. Gateway WebSocket Session (TCP): Used to connect to the Voice Websocket Session and
//     receive information about who is in the voice channel.
//
//  2. Voice WebSocket Session (TCP): Used to connect to the Voice UDP Connection and
//     receive information about who is speaking in the voice channel.
//
//  3. Voice Connection (UDP): Used to send and receive audio from Discord.
type VoiceConnection struct {
	// State represents the Voice Status of the bot.
	State GatewayVoiceStateUpdate

	// Session represents the Discord Gateway WebSocket Session of the VoiceConnection.
	Session *Session

	// VoiceSession represents the Discord Voice WebSocket Session of the VoiceConnection.
	VoiceSession *VoiceSession

	// Connection represents the Voice UDP connection of the VoiceConnection.
	Connection *net.UDPConn

	// Handlers represents a Voice Connection's Voice Session event handlers.
	Handlers *VoiceHandlers
}

// protectedHandlerVoiceStateUpdate represents a protected event handler for the VoiceStateUpdate event.
func protectedHandlerVoiceStateUpdate(bot *Client, v *VoiceStateUpdate) {
	// sessionID to map[guildID]*VoiceConnection.
	if v, ok := bot.Sessions.Voice.Load(v.SessionID); ok {
		gvcMap := v.(*sync.Map)

		// map[guildID]*VoiceConnection.
		v2, _ := gvcMap.Load(v.GuildID)
		vc := v2.(*VoiceConnection)
		if vc.Session.ID == v.SessionID {
			// wait <- 0
		}

		// TODO: concurrency checks, can this be empty with a valid voicestateupdate?
	} else {
		// Session disconnected, reconnected with new id, voice connection map wasn't reset?
		//vc.client_manager.Voice.Store(vc.Session.ID, new(sync.Map))

		// goto SESSIONMANAGER
	}
}

// protectedHandlerVoiceServerUpdate represents a protected event handler for the VoiceServerUpdate event.
func protectedHandlerVoiceServerUpdate(bot *Client, v *VoiceServerUpdate) {
	// check that the provided GuildID matches the incoming Voice Server Update GuildID.
	if vc.State.GuildID == v.GuildID {
		vc.VoiceSession.Lock()
		vc.VoiceSession.VoiceServerInfo = v
		vc.VoiceSession.Unlock()

		// A null endpoint means that the voice server is reallocating.
		if v.Endpoint == nil {
			// disconnect from the current voice server.
			vc.VoiceSession.disconnect(FlagClientCloseEventCodeNormal)

			// TODO: wait until a new voice server is allocated.
		}
	}
}

// VoiceConnection connects the bot to a Discord Voice Channel using the Discord Gateway.
func (vc *VoiceConnection) Connect(bot *Client) error {
	if vc.Session == nil || !vc.Session.isConnected() {
		return fmt.Errorf("ConnectVoice: Session must be connected to the Discord Gateway to connect to voice channel")
	}

	if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_VOICE_STATES] {
		return fmt.Errorf("ConnectVoice: Session must be connected to the Discord Gateway with the GUILD_VOICE_STATES intent. " +
			"Use `bot.Config.Gateway.EnableIntent(FlagIntentGUILD_VOICE_STATES)` before connecting the given session to the Discord Gateway to enable it.") //lint:ignore ST1005 format help message.
	}

	// check that the user (developer) has provided a ChannelID.
	if vc.State.ChannelID == nil || *vc.State.ChannelID == "" {
		return fmt.Errorf("ConnectVoice: Voice ChannelID must be non-nil and non-empty to connect to voice channel")
	}

	if vc.Handlers == nil {
		vc.Handlers = new(VoiceHandlers)
	}

	vc.VoiceSession = newVoiceSession()

SESSIONMANAGER:
	// Store the Voice Connection into the bot's Session Manager.
	//
	// sessionID to map[guildID]*VoiceConnection.
	if v, ok := vc.client_manager.Voice.Load(vc.Session.ID); ok {
		gvcMap := v.(*sync.Map)

		// map[guildID]*VoiceConnection.
		gvcMap.Store(vc.State.GuildID, vc)
	} else {
		vc.client_manager.Voice.Store(vc.Session.ID, new(sync.Map))

		goto SESSIONMANAGER
	}
	// TODO: SessionManager voice connection reset on gateway session disconnect?

	// a channel is used to wait for the each event.
	wait := make(chan int)

	// Send an Opcode 4 Gateway Voice State Update to the Discord Gateway.
	//
	// According to Discord, the response events of this send event should never be cached.
	//
	// Disclaimer. The bot will not receive response events when the voice channel is full,
	// unless the bot has the MANAGE_CHANNELS permission.
	vc.State.SendEvent(bot, vc.Session)

	// Wait for the Voice State Update and Voice Server Update events.
	select {
	case <-vc.Session.Context.Done():
		return <-vc.Session.manager.err
	case <-wait:
		break
	}

VOICESERVERUPDATE:
	for {
		vc.VoiceSession.RLock()

		if vc.VoiceSession.VoiceServerInfo != nil && vc.VoiceSession.VoiceServerInfo.Endpoint != nil {
			vc.VoiceSession.RUnlock()

			break
		}

		select {
		case <-vc.Session.Context.Done():
			vc.VoiceSession.RUnlock()

			return <-vc.Session.manager.err
		default:
			vc.VoiceSession.RUnlock()
			//lint:ignore SA4011 break into for loop.
			break
		}
	}

	// A null endpoint means that the voice server is reallocating.
	if vc.VoiceSession.VoiceServerInfo.Endpoint == nil {
		goto VOICESERVERUPDATE
	}

	// Establish a Voice WebSocket Connection (TCP).
	// https://discord.com/developers/docs/topics/voice-connections#establishing-a-voice-websocket-connection
	// TODO: Spawn on other routine for VoiceConnection manager between UDP and VoiceSession?
	if err := vc.VoiceSession.connect(bot, vc); err != nil {
		return fmt.Errorf("voice: %w", err)
	}

	Logger.Printf("works")

	return nil
}
