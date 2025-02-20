package wrapper

import (
	"fmt"

	json "github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
	"github.com/switchupcb/disgo/wrapper/socket"
	"github.com/switchupcb/websocket"
)

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
