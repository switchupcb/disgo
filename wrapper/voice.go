package wrapper

import (
	"fmt"
)

// VoiceConnection represents a Discord Voice Channel Connection.
//
// A Discord Voice Channel Connection is composed of three connections:
//
//  1. Gateway WebSocket Session: Used to connect to the Voice Websocket Session and
//     receive information about who is in the voice channel.
//
//  2. Voice WebSocket Session: Used to connect to the Voice UDP Connection and
//     receive information about who is speaking in the voice channel.
//
//  3. Voice UDP Connection: Used to send and receive audio from Discord.
type VoiceConnection struct {
	// State represents the Voice Status of the bot.
	State GatewayVoiceStateUpdate

	// Session represents the Discord Gateway WebSocket Session of the VoiceConnection.
	Session *Session

	// VoiceSession represents the Discord Voice WebSocket Session of the VoiceConnection.
	VoiceSession *VoiceSession

	// Connection represents the Voice UDP connection of the VoiceConnection.
	Connection *UDPConnection

	// Handlers represents a Voice Connection's Voice Session event handlers.
	Handlers *VoiceHandlers

	// client_manager represents the *Client Session Manager of the VoiceConnection.
	client_manager *SessionManager
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

	// Set up the event handler for the Voice State Update and Voice Server Update events.
	//
	// According to Discord, the response events of this send event should never be cached.
	//
	// Disclaimer. The bot will not receive response events when the voice channel is full,
	// unless the bot has the MANAGE_CHANNELS permission.
	if bot.Handlers == nil {
		return fmt.Errorf(errNoHandlers) //lint:ignore ST1005 format help message.
	}

	if vc.Handlers == nil {
		vc.Handlers = new(VoiceHandlers)
	}

	vc.VoiceSession = newVoiceSession()
	vc.client_manager = bot.Sessions

	// a channel is used to wait for the each event.
	wait := make(chan int)

	// TODO: forward handlers to voice connection for simplicity?

	// Voice State Update event handler.
	if err := bot.Handle(FlagGatewayEventNameVoiceStateUpdate, func(v *VoiceStateUpdate) {
		if vc.Session.ID == v.SessionID {
			wait <- 0
		}
	}); err != nil {
		return fmt.Errorf("ConnectVoice: %w", err)
	}

	// TODO: defer removal of handler on error

	// Voice Server Update event handler.
	if err := bot.Handle(FlagGatewayEventNameVoiceServerUpdate, func(v *VoiceServerUpdate) {
		// check that the provided GuildID matches the incoming Voice Server Update GuildID.
		if vc.State.GuildID == v.GuildID {
			vc.VoiceSession.Lock()
			vc.VoiceSession.VoiceServerInfo = v
			vc.VoiceSession.Unlock()

			// TODO: A null endpoint means that the voice server is reallocating.
			// Disconnect from the current voice server and wait until a new voice server is allocated.
		}
	}); err != nil {
		return fmt.Errorf("ConnectVoice: %w", err)
	}

	// TODO: defer removal of handler on error

	// Send an Opcode 4 Gateway Voice State Update to the Discord Gateway.
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

	// Establish a Voice WebSocket Connection (UDP).
	// https://discord.com/developers/docs/topics/voice-connections#establishing-a-voice-websocket-connection
	//
	// A null endpoint means that the voice server is reallocating.
	if vc.VoiceSession.VoiceServerInfo.Endpoint == nil {
		goto VOICESERVERUPDATE
	}

	// TODO: Spawn on other routine for VoiceConnection manager between UDP and VoiceSession?
	// connect to the Discord Voice Server.
	if err := vc.VoiceSession.connect(bot, vc); err != nil {
		return fmt.Errorf("voice: %w", err)
	}

	Logger.Printf("works")

	// https://discord.com/developers/docs/topics/voice-connections#establishing-a-voice-udp-connection
	// https://discord.com/developers/docs/topics/voice-connections#ip-discovery
	//
	// 	https://discord.com/developers/docs/topics/opcodes-and-status-codes#voice
	// 	opcode 1 send
	//
	// 	https://discord.com/developers/docs/topics/voice-connections#establishing-a-voice-udp-connection-encryption-modes
	// 	opcode 4 receive

	// Connection is established, create channel for external library to process voice
	// https://discord.com/developers/docs/topics/voice-connections#encrypting-and-sending-voice

	// Store Voice Connection into the bot's Session Manager.
	if slice, ok := vc.client_manager.Voice.Load(vc.Session.ID); ok {
		vcs := slice.([]*VoiceConnection)
		vc.client_manager.Voice.Store(vc.Session.ID, append(vcs, vc))
	} else {
		vc.client_manager.Voice.Store(vc.Session.ID, []*VoiceConnection{vc})
	}

	return nil
}
