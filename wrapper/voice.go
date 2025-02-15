package wrapper

import (
	"errors"
	"fmt"
	"net"
)

// VoiceChannelConnection represents a Discord Voice Channel Connection.
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
type VoiceChannelConnection struct {
	// State represents the Voice Status of the bot.
	State GatewayVoiceStateUpdate

	// Session represents the Discord Gateway WebSocket Session of the VoiceChannelConnection.
	GatewaySession *Session

	// VoiceSession represents the Discord Voice WebSocket Session of the VoiceChannelConnection.
	VoiceSession *VoiceSession

	// Connection represents the Voice UDP connection of the VoiceChannelConnection.
	Connection *net.UDPConn

	// Handlers represents a VoiceChannelConnection's Voice Session event handlers.
	Handlers *VoiceHandlers
}

// addDefaultVoiceStateUpdate adds a default event handler for the VoiceStateUpdate event to the bot.
func addDefaultVoiceStateUpdate(bot *Client) error {
	return bot.Handle(FlagGatewayEventNameVoiceStateUpdate, func(v *VoiceStateUpdate) {
		// check whether this Voice State Update describes the bot.
		if v.UserID != bot.ApplicationID {
			return
		}

		// update the Voice State for the Voice State Update Voice Channel's Connection.
		vc := bot.Sessions.GetVoiceChannelConnection(v.SessionID, *v.GuildID)
		if vc == nil {
			return
		}

		vc.VoiceSession.Lock()
		vc.VoiceSession.ID = v.SessionID
		vc.State.ChannelID = v.ChannelID
		vc.State.SelfMute = v.SelfMute
		vc.State.SelfDeaf = v.SelfDeaf
		vc.VoiceSession.Unlock()
	})
}

// addDefaultHandlerVoiceServerUpdate adds a default event handler for the VoiceServerUpdate event to the bot.
func addDefaultHandlerVoiceServerUpdate(bot *Client) error {
	return bot.Handle(FlagGatewayEventNameVoiceServerUpdate, func(v *VoiceServerUpdate) {
		vc := bot.Sessions.GetVoiceChannelConnection(SessionManagerVoiceKeyUnknownSession, v.GuildID)
		if vc == nil {
			return
		}

		vc.VoiceSession.Lock()

		// check that the provided GuildID matches the incoming Voice Server Update GuildID.
		if vc.State.GuildID == v.GuildID {
			vc.VoiceSession.VoiceServerInfo = v
			vc.VoiceSession.VoiceServerInfo.Endpoint = v.Endpoint
		}

		vc.VoiceSession.Unlock()
	})
}

// addDefaultHandlerSessionDescription adds a default event handler for the SessionDescription event to the Voice Session.
func addDefaultHandlerSessionDescription(vc *VoiceChannelConnection) error {
	return vc.Handle(FlagVoiceOpcodeNameSessionDescription, func(sd *SessionDescription) {
		// TODO: Encryption and Decryption in connectUDP()
		// https://discord.com/developers/docs/topics/voice-connections#transport-encryption-and-sending-voice
	})
}

// VoiceConnection connects the bot to a Discord Voice Channel using the Discord Gateway.
func (vc *VoiceChannelConnection) Connect(bot *Client) error {
	if bot.ApplicationID == "" {
		return errors.New("ConnectVoice: Client must have an ApplicationID to connect to voice channel." +
			"Set `bot.ApplicationID` before connecting to a voice channel.") //lint:ignore ST1005 format help message.
	}

	// check that the user (developer) has provided a ChannelID.
	if vc.State.ChannelID == nil || *vc.State.ChannelID == "" {
		return errors.New("ConnectVoice: Voice ChannelID must be non-nil and non-empty to connect to voice channel")
	}

	if vc.GatewaySession == nil || !vc.GatewaySession.isConnected() {
		return errors.New("ConnectVoice: Session must be connected to the Discord Gateway to connect to voice channel")
	}

	if !bot.Config.Gateway.IntentSet[FlagIntentGUILD_VOICE_STATES] {
		return errors.New("ConnectVoice: Session must be connected to the Discord Gateway with the GUILD_VOICE_STATES intent. " +
			"Use `bot.Config.Gateway.EnableIntent(FlagIntentGUILD_VOICE_STATES)` before connecting the Gateway Session to the Discord Gateway.") //lint:ignore ST1005 format help message.
	}

	vc.VoiceSession = newVoiceSession()

	if vc.Handlers == nil {
		vc.Handlers = new(VoiceHandlers)
	}

	if len(bot.Handlers.VoiceStateUpdate) == 0 {
		if err := addDefaultVoiceStateUpdate(bot); err != nil {
			return fmt.Errorf("ConnectVoice: %w", err)
		}
	}

	if len(bot.Handlers.VoiceServerUpdate) == 0 {
		if err := addDefaultHandlerVoiceServerUpdate(bot); err != nil {
			return fmt.Errorf("ConnectVoice: %w", err)
		}
	}

	if len(vc.Handlers.SessionDescription) == 0 {
		if err := addDefaultHandlerSessionDescription(vc); err != nil {
			return fmt.Errorf("ConnectVoice: %w", err)
		}
	}

	// Store the Voice Connection into the bot's Voice Session Manager.
	bot.Sessions.StoreVoiceChannelConnection(vc.GatewaySession.ID, vc.State.GuildID, vc)

	// Send an Opcode 4 Gateway Voice State Update to the Discord Gateway.
	//
	// According to Discord, the response events of this send event should never be cached.
	//
	// Disclaimer. The bot will not receive response events when the voice channel is full,
	// unless the bot has the MANAGE_CHANNELS permission.
	if err := vc.State.SendEvent(bot, vc.GatewaySession); err != nil {
		return fmt.Errorf("voice: %w", err)
	}

VOICESERVERUPDATE:
	// Wait for the Voice State Update and Voice Server Update events.
	for {
		vc.VoiceSession.RLock()

		if vc.VoiceSession.VoiceServerInfo != nil && vc.VoiceSession.VoiceServerInfo.Endpoint != nil {
			vc.VoiceSession.RUnlock()

			break
		}

		select {
		case <-vc.GatewaySession.Context.Done():
			vc.VoiceSession.RUnlock()

			return <-vc.GatewaySession.manager.err
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
	if err := vc.VoiceSession.connect(bot, vc); err != nil {
		return fmt.Errorf("voice: %w", err)
	}

	return nil
}
