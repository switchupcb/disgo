package wrapper

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/switchupcb/disgo/wrapper/internal/socket"
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

		logPayload(logSession(Logger.Info(), s.ID), payload.Op, payload.Data).Msg("received payload")

		if err = s.onPayload(bot, *payload); err != nil {
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
		s.Lock()
		atomic.AddInt32(&s.manager.pulses, 1)
		s.Unlock()

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
		s.reconnect("reconnecting session due to Opcode 7 Reconnect")

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
