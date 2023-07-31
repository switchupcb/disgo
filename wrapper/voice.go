package wrapper

import (
	"fmt"

	"github.com/switchupcb/websocket"
)

const (
	voiceWebSocketConnectionURLProtocol = "wss://"
)

// ConnectVoice connects a session to a Discord Voice Channel.
func (s *Session) ConnectVoice(bot *Client, vc GatewayVoiceStateUpdate) error {
	// Set up the event handler for the Voice State Update and Voice Server Update events.
	//
	// According to Discord, the response events of this send event should never be cached.
	//
	// Disclaimer. The bot will not receive response events when the voice channel is full,
	// unless the bot has the MANAGE_CHANNELS permission.
	if bot.Handlers == nil {
		return fmt.Errorf(errNoHandlers) //lint:ignore ST1005 format help message.
	}

	// check that the user (developer) has provided a ChannelID.
	if vc.ChannelID == nil || *vc.ChannelID == "" {
		return fmt.Errorf("ConnectVoice: Voice ChannelID must be non-nil and non-empty to connect to voice channel")
	}

	// a channel is used to wait for the each event.
	wait := make(chan int)

	// Voice State Update event handler.
	if err := bot.Handle(FlagGatewayEventNameVoiceStateUpdate, func(v *VoiceStateUpdate) {
		if s.ID == v.SessionID {
			wait <- 0
		}
	}); err != nil {
		return fmt.Errorf("ConnectVoice: %w", err)
	}

	// TODO: defer removal of handler on error

	// Voice Server Update event handler.
	if err := bot.Handle(FlagGatewayEventNameVoiceServerUpdate, func(v *VoiceServerUpdate) {
		// check that the provided GuildID matches the incoming Voice Server Update GuildID.
		if vc.GuildID == v.GuildID {
			s.Lock()
			s.VoiceServerInfo = v
			s.Unlock()

			// TODO: A null endpoint means that the voice server is reallocating.
			// Disconnect from the current voice server and wait until a new voice server is allocated.
		}
	}); err != nil {
		return fmt.Errorf("ConnectVoice: %w", err)
	}

	// TODO: defer removal of handler on error

	// connect to the Gateway.
	// https://discord.com/developers/docs/topics/voice-connections#retrieving-voice-server-information
	if !s.isConnected() {
		if err := s.Connect(bot); err != nil {
			return fmt.Errorf("voice: %w", err)
		}
	}

	// Send an Opcode 4 Gateway Voice State Update to the Discord Gateway.
	vc.SendEvent(bot, s)

	// Wait for the Voice State Update and Voice Server Update events.
	select {
	case <-s.Context.Done():
		return <-s.manager.err
	case <-wait:
		break
	}

VOICESERVERUPDATE:
	for {
		s.RLock()

		if s.VoiceServerInfo != nil && s.VoiceServerInfo.Endpoint != nil {
			s.RUnlock()

			break
		}

		select {
		case <-s.Context.Done():
			s.RUnlock()

			return <-s.manager.err
		default:
			s.RUnlock()
			//lint:ignore SA4011 break into for loop.
			break
		}
	}

	// Establish a Voice WebSocket Connection (UDP).
	// https://discord.com/developers/docs/topics/voice-connections#establishing-a-voice-websocket-connection
	//
	// A null endpoint means that the voice server is reallocating.
	s.RLock()

	if s.VoiceServerInfo.Endpoint == nil {
		goto VOICESERVERUPDATE
	}

	var err error

	// connect to the Voice WebSocket Connection.
	// TODO: Dial to new VoiceSession.
	if s.Conn, _, err = websocket.Dial(
		s.Context,
		voiceWebSocketConnectionURLProtocol+*s.VoiceServerInfo.Endpoint+gatewayEndpointParams,
		nil); err != nil {
		return fmt.Errorf("error connecting to the Discord Voice Gateway: %w", err)
	}

	s.RUnlock()

	Logger.Printf("works")
	// Once connected to the voice WebSocket endpoint, we can send an Opcode 0 Identify payload with our server_id, user_id, session_id, and token:
	//
	// equivalent to initial() function

	// The voice server should respond with an Opcode 2 Ready payload
	//
	// ? https://discord.com/developers/docs/topics/voice-connections#establishing-a-voice-websocket-connection-example-voice-ready-payload

	// https://discord.com/developers/docs/topics/voice-connections#establishing-a-voice-udp-connection
	// https://discord.com/developers/docs/topics/voice-connections#ip-discovery

	// https://discord.com/developers/docs/topics/opcodes-and-status-codes#voice
	// opcode 1 send

	// https://discord.com/developers/docs/topics/voice-connections#establishing-a-voice-udp-connection-encryption-modes
	// 	opcode 4 receive

	// store into client manager sessions as voice session

	// Connection is established, create channel for external library to process voice
	// https://discord.com/developers/docs/topics/voice-connections#encrypting-and-sending-voice

	return nil
}
