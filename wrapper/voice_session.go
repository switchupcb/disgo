package wrapper

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	json "github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
	"github.com/switchupcb/disgo/wrapper/socket"
	"github.com/switchupcb/websocket"
	"golang.org/x/sync/errgroup"
)

const (
	voiceWebSocketConnectionURLProtocol = "wss://"
	voiceEndpointParams                 = "?v=" + VersionDiscordVoiceGateway + "&encoding=json"
)

// VoiceSession represents a Discord Voice WebSocket Session.
type VoiceSession struct {
	// ID represents the Discord Gateway Websocket Session ID of the Voice Session.
	ID string

	// Nonce represents the heartbeat integer nonce.
	//
	// https://discord.com/developers/docs/topics/voice-connections#heartbeating
	Nonce int64

	// VoiceServerInfo represents Voice Server Update information for a session
	// connected to a voice channel.
	//
	// https://discord.com/developers/docs/topics/gateway-events#voice-server-update
	VoiceServerInfo *VoiceServerUpdate

	// Context carries request-scoped data for the Discord Voice Session.
	//
	// Context is also used as a signal for the Voice Session's goroutines.
	Context context.Context

	// Conn represents a WebSocket Connection to the Discord Voice Server.
	Conn *websocket.Conn

	// heartbeat contains the fields required to implement the heartbeat mechanism.
	heartbeat *voice_heartbeat

	// manager represents a manager of a Voice Session's goroutines.
	manager *voice_manager

	// RWMutex is used to protect the Session's variables from data races
	// by providing transactional functionality.
	sync.RWMutex
}

// isConnected returns whether the session is connected.
func (s *VoiceSession) isConnected() bool {
	if s.Context == nil {
		return false
	}

	select {
	case <-s.Context.Done():
		return false
	default:
		return true
	}
}

// canReconnect determines whether the session is in a valid state to reconnect.
func (s *VoiceSession) canReconnect() bool {
	return s.ID != "" && *s.VoiceServerInfo.Endpoint != "" && atomic.LoadInt64(&s.Nonce) != 0
}

// connect connects a session to a WebSocket Connection.
func (s *VoiceSession) connect(bot *Client, vc *VoiceChannelConnection) error {
	LogSession(Logger.Info(), s.ID).Str(LogCtxClient, bot.ApplicationID).Msg("connecting voice session")

	if s.isConnected() {
		return fmt.Errorf("voice session %q is already connected", s.ID)
	}

	var err error

	// connect to the Discord Voice Server Websocket.
	s.manager = new(voice_manager)
	s.Context, s.manager.cancel = context.WithCancel(context.Background())
	if s.Conn, _, err = websocket.Dial(
		s.Context,
		voiceWebSocketConnectionURLProtocol+*s.VoiceServerInfo.Endpoint+voiceEndpointParams,
		nil,
	); err != nil {
		return fmt.Errorf("error connecting to the Discord Voice Server: %w", err)
	}

	// handle the incoming Hello event upon connecting to the Voice Server.
	hello := new(VoiceHello)
	if err := readEventVoice(s, hello); err != nil {
		err = fmt.Errorf("error reading initial VoiceHello event: %w", err)
		sessionErr := ErrorSession{SessionID: s.ID, Err: err}
		if disconnectErr := s.disconnect(FlagClientCloseEventCodeNormal); disconnectErr != nil {
			sessionErr.Err = ErrorDisconnect{
				Action:     err,
				Err:        disconnectErr,
				Connection: ErrConnectionSessionVoice,
			}
		}

		return sessionErr
	}

	for _, handler := range vc.Handlers.VoiceHello {
		go handler(hello)
	}

	// begin sending heartbeat payloads every heartbeat_interval ms.
	ms := time.Millisecond * time.Duration(hello.HeartbeatInterval)
	s.heartbeat = &voice_heartbeat{
		interval: ms,
		ticker:   time.NewTicker(ms),
		send:     make(chan VoiceHeartbeat),

		// add a HeartbeatACK to the HeartbeatACK channel to prevent
		// the length of the HeartbeatACK channel from being 0 immediately,
		// which results in an attempt to reconnect.
		acks: 1,
	}

	// create a goroutine group for the Session.
	s.manager.Group, s.manager.signal = errgroup.WithContext(s.Context)
	s.manager.err = make(chan error, 1)

	// spawn the heartbeat pulse goroutine.
	s.manager.routines.Add(1)
	atomic.AddInt32(&s.manager.pulses, 1)
	s.manager.Go(func() error {
		s.pulse()
		return nil
	})

	// spawn the heartbeat beat goroutine.
	s.manager.routines.Add(1)
	s.manager.Go(func() error {
		if err := s.beat(); err != nil {
			return ErrorSession{
				SessionID: s.ID,
				Err:       fmt.Errorf("heartbeat: %w", err),
			}
		}

		return nil
	})

	// send the initial Identify or Resumed packet.
	if err := s.initial(bot, vc); err != nil {
		sessionErr := ErrorSession{SessionID: s.ID, Err: err}
		if disconnectErr := s.disconnect(FlagClientCloseEventCodeNormal); disconnectErr != nil {
			sessionErr.Err = ErrorDisconnect{
				Action:     err,
				Err:        disconnectErr,
				Connection: ErrConnectionSessionVoice,
			}
		}

		return sessionErr
	}

	// spawn the event listener listen goroutine.
	s.manager.routines.Add(1)
	s.manager.Go(func() error {
		if err := s.listen(vc); err != nil {
			return ErrorSession{
				SessionID: s.ID,
				Err:       fmt.Errorf("listen: %w", err),
			}
		}

		return nil
	})

	// spawn the manager goroutine.
	s.manager.routines.Add(1)
	go s.manage()

	// ensure that the Session's goroutines are spawned.
	s.manager.routines.Wait()

	return nil
}

// initial sends the initial Identify or Resume packet required to connect to the Voice Server,
// then handles the incoming Ready or Resumed packet that indicates a successful connection.
func (s *VoiceSession) initial(bot *Client, vc *VoiceChannelConnection) error {
	if !s.canReconnect() {
		// send an Opcode 0 Identify to the Discord Voice Server.
		identify := VoiceIdentify{
			ServerID:  s.VoiceServerInfo.GuildID,
			UserID:    bot.ApplicationID,
			SessionID: s.ID,
			Token:     s.VoiceServerInfo.Token,
		}

		if err := identify.SendEvent(s); err != nil {
			return err
		}

	} else {
		// send an Opcode 7 Resume to the Discord Voice Server to reconnect the session.
		resume := VoiceResume{
			ServerID:  s.VoiceServerInfo.GuildID,
			SessionID: s.ID,
			Token:     bot.Authentication.Token,
		}

		if err := resume.SendEvent(s); err != nil {
			return err
		}
	}

	// handle the incoming Ready or Resumed event.
	payload := getVoicePayload()
	defer putVoicePayload(payload)
	if err := socket.Read(s.Context, s.Conn, payload); err != nil {
		return fmt.Errorf("error reading initial voice payload: %w", err)
	}

	LogPayload(LogSession(Logger.Info(), s.ID), payload.Op, payload.Data).Msg("received initial voice payload")

	switch payload.Op {
	// When a connection is successful, the Discord Voice Server will respond with a Ready payload.
	case FlagVoiceOpcodeReadyServer:
		ready := new(VoiceReady)
		if err := json.Unmarshal(payload.Data, ready); err != nil {
			return fmt.Errorf("error reading ready event: %w", err)
		}

		LogSession(Logger.Info(), s.ID).Msg("received VoiceReady event")

		for _, handler := range vc.Handlers.VoiceReady {
			go handler(ready)
		}

		// Establish a Voice Connection (UDP).
		// https://discord.com/developers/docs/topics/voice-connections#establishing-a-voice-udp-connection
		if err := vc.connectUDP(ready); err != nil {
			return fmt.Errorf("error connecting to UDP Voice Server: %w", err)
		}

	// When a reconnection is successful, the Discord Voice Server will respond
	// with a Resumed payload.
	case FlagVoiceOpcodeResumed:
		LogSession(Logger.Info(), s.ID).Msg("received VoiceResumed event")

		for _, handler := range vc.Handlers.VoiceResumed {
			go handler(&VoiceResumed{})
		}

		// TODO: RESUME: ConnectUDP (?)

		// When a reconnection is unsuccessful, the Discord Voice Server will close
		// with an appropriate close event code.
	default:
		return fmt.Errorf("voice session %q received payload %d during connection which is unexpected", s.ID, payload.Op)
	}

	return nil
}

// disconnect disconnects a session from a WebSocket Connection using the given status code.
func (s *VoiceSession) disconnect(code int) error {
	id := s.ID
	LogSession(Logger.Info(), id).Msgf("disconnecting voice session with code %d", FlagClientCloseEventCodeNormal)

	s.manager.signal = context.WithValue(s.manager.signal, keySignal, signalDisconnect)

	// cancel the context to kill the goroutines of the Voice Session.
	defer s.manager.cancel()

	if err := s.Conn.Close(websocket.StatusCode(code), ""); err != nil {
		return fmt.Errorf("%w", err)
	}

	putVoiceSession(s)

	LogSession(Logger.Info(), id).Msgf("disconnected voice session with code %d", FlagClientCloseEventCodeNormal)

	return nil
}

// readEventVoice is a helper function for reading events from the Voice WebSocket Session.
func readEventVoice(s *VoiceSession, dst any) error {
	payload := new(VoicePayload)
	if err := socket.Read(s.Context, s.Conn, payload); err != nil {
		return fmt.Errorf("readEvent: %w", err)
	}

	if err := json.Unmarshal(payload.Data, dst); err != nil {
		return fmt.Errorf("readEvent: %w", err)
	}

	return nil
}

// writeEventVoice is a helper function for writing voice events to the WebSocket Session.
func writeEventVoice(s *VoiceSession, op int, name string, dst any) error {
	LogCommandVoice(log.Trace(), op, name).Msg("sending voice server command")

	// write the event to the WebSocket Connection.
	event, err := json.Marshal(dst)
	if err != nil {
		return fmt.Errorf("writeEvent: %w", err)
	}

	if err = socket.Write(s.Context, s.Conn, websocket.MessageText,
		VoicePayload{
			Op:   op,
			Data: event,
		}); err != nil {
		return fmt.Errorf("writeEvent: %w", err)
	}

	LogCommandVoice(log.Trace(), op, name).Msg("sending voice server command")

	return nil
}
