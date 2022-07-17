package wrapper

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/switchupcb/disgo/wrapper/internal/socket"
	"nhooyr.io/websocket"
)

// listen listens to the connection for payloads from the Discord Gateway.
func (s *Session) listen(bot *Client) error {
	s.manager.routines.Done()

	var err error

	for {
		payload := getPayload()
		if err = socket.Read(s.Context, s.Conn, payload); err != nil {
			break
		}

		log.Println("PAYLOAD", payload.Op, string(payload.Data))
		if err = s.onPayload(bot, *payload); err != nil {
			break
		}
	}

	s.Lock()
	defer s.Unlock()

	select {
	case <-s.Context.Done():
		return nil

	default:
		closeErr := new(websocket.CloseError)
		if errors.As(err, closeErr) {
			if gcErr := s.handleGatewayCloseError(bot, closeErr); gcErr == nil {
				return nil
			}
		}

		return err
	}
}

// onPayload handles an Discord Gateway Payload.
func (s *Session) onPayload(bot *Client, payload GatewayPayload) error {
	defer putPayload(&payload)

	// https://discord.com/developers/docs/topics/opcodes-and-status-codes#gateway-gateway-opcodes
	switch payload.Op {
	// run the bot's event handlers.
	case FlagGatewayOpcodeDispatch:
		atomic.StoreInt64(&s.Seq, *payload.SequenceNumber)
		go bot.handle(*payload.EventName, payload.Data)

	// send an Opcode 1 Heartbeat to the Discord Gateway.
	case FlagGatewayOpcodeHeartbeat:
		s.manager.Go(func() error {
			if err := s.respond(payload.Data); err != nil {
				return fmt.Errorf("respond: %w", err)
			}

			return nil
		})

	// handle the successful acknowledgement of the client's last heartbeat.
	case FlagGatewayOpcodeHeartbeatACK:
		s.Lock()
		atomic.AddUint32(&s.heartbeat.acks, 1)
		s.Unlock()

	// occurs when the Discord Gateway is shutting down the connection, while signalling the client to reconnect.
	case FlagGatewayOpcodeReconnect:
		s.manager.Go(func() error {
			log.Printf("reconnecting Session %q due to Opcode 7 Reconnect", s.ID)
			s.Context = context.WithValue(s.Context, keySignal, signalReconnect)

			return s.disconnect(FlagClientCloseEventCodeReconnect)
		})

		return nil

	// in the context of onPayload, an Invalid Session occurs when an active session is invalidated.
	case FlagGatewayOpcodeInvalidSession:
		// wait for Discord to close the session, then complete a fresh connect.
		<-time.NewTimer(invalidSessionWaitTime).C

		s.Lock()
		defer s.Unlock()

		if err := s.initial(bot, 0); err != nil {
			return err
		}
	}

	return nil
}

// handleGatewayCloseError handles a Discord Gateway WebSocket CloseError.
//
// returns the given closeErr if a disconnect is warranted.
func (s *Session) handleGatewayCloseError(bot *Client, closeErr *websocket.CloseError) error {
	code, ok := GatewayCloseEventCodes[int(closeErr.Code)]
	switch ok {
	// Gateway Close Event Code is known.
	case true:
		log.Printf(
			"Session %q received Gateway Close Event Code %d %s: %s",
			s.ID, code.Code, code.Description, code.Explanation,
		)

		if code.Reconnect {
			s.manager.Go(func() error {
				log.Printf("reconnecting Session %q due to Gateway Close Event Code %d", s.ID, code.Code)
				s.Context = context.WithValue(s.Context, keySignal, signalReconnect)

				return s.disconnect(FlagClientCloseEventCodeReconnect)
			})

			return nil
		}

		return closeErr

	// Gateway Close Event Code is unknown.
	default:

		// when another goroutine calls disconnect(),
		// s.Conn.Close is called before s.cancel which will result in
		// a CloseError with the close code that Disgo uses to reconnect.
		if closeErr.Code == websocket.StatusCode(FlagClientCloseEventCodeReconnect) {
			return nil
		}

		log.Printf(
			"Session %q received unknown Gateway Close Event Code %d with reason %q",
			s.ID, closeErr.Code, closeErr.Reason,
		)

		return closeErr
	}
}
