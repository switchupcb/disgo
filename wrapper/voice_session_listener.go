package wrapper

import (
	"sync/atomic"

	"github.com/switchupcb/disgo/wrapper/socket"
)

// listen listens to the connection for payloads from the Discord Voice Server.
func (s *VoiceSession) listen(vc *VoiceConnection) error {
	s.manager.routines.Done()

	var err error

	for {
		payload := getVoicePayload()
		if err = socket.Read(s.Context, s.Conn, payload); err != nil {
			break
		}

		LogPayload(LogSession(Logger.Info(), s.ID), payload.Op, payload.Data).Msg("received voice payload")

		if err = s.onPayload(vc, *payload); err != nil {
			break
		}
	}

	s.Lock()
	defer s.Unlock()
	defer s.logClose("listen")

	select {
	case <-s.Context.Done():
		return nil

	default:
		return err
	}
}

// onPayload handles an Discord Voice Server Payload.
func (s *VoiceSession) onPayload(vc *VoiceConnection, payload VoicePayload) error {
	defer putVoicePayload(&payload)

	// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-opcodes
	switch payload.Op {
	case FlagVoiceOpcodeSpeaking:
		go vc.handle(FlagVoiceOpcodeNameSpeaking, payload.Data)

	// handle the successful acknowledgement of the client's last heartbeat.
	case FlagVoiceOpcodeHeartbeatACK:
		s.Lock()
		atomic.AddUint32(&s.heartbeat.acks, 1)
		s.Unlock()

	case FlagVoiceOpcodeClientDisconnect:
		go vc.handle(FlagVoiceOpcodeNameClientDisconnect, payload.Data)
	}

	return nil
}
