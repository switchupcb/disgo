package integration_test

import (
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	. "github.com/switchupcb/disgo"
)

// TestConnect tests Connect(), Disconnect(), heartbeat(), listen(), and onPayload()
// in order to ensure that WebSocket functionality works.
func TestConnect(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// timeout is used to prevent this test from lasting longer than expected.
	timeout := time.NewTimer(time.Second * 4)

	bot := &Client{
		Authentication: BotToken(os.Getenv("TOKEN")),
		Config:         DefaultConfig(),
		Handlers:       new(Handlers),
	}

	s := NewSession()

	// a channel is used to wait for the Ready event.
	wait := make(chan int)

	// a Ready event is sent upon a successful connection.
	bot.Handle(FlagGatewayEventNameReady, func(*Ready) {
		wait <- 0
	})

	// connect to the Discord Gateway (WebSocket Session).
	if err := s.Connect(bot); err != nil {
		t.Fatalf("%v", err)
	}

	// connecting to a connected session should result in an error.
	if err := s.Connect(bot); err == nil {
		t.Fatalf("expected error while connecting to already connected session")
	}

	// wait until a Heartbeat is sent.
	for {
		select {
		case <-timeout.C:
			t.Fatalf("test took longer than expected while sending heartbeat")
		case <-s.Context.Done():
			t.Fatalf("disconnected while waiting to send heartbeat")
		default:
			// a Heartbeat is sent when the amount of ACKs since the
			// last Heartbeat was sent is 0.
			if s.Monitor() == 0 {
				goto ACK
			}
		}
	}

ACK:
	// once Heartbeat has been sent, wait to receive the respective HeartbeatACK.
	for {
		select {
		case <-timeout.C:
			t.Fatalf("test took longer than expected while waiting for HeartbeatACK")
		case <-s.Context.Done():
			t.Fatalf("disconnected while waiting for HeartbeatACK")
		default:
			// a respective HeartbeatACK should be sent to s.heartbeat.ack.
			if s.Monitor() == 1 {
				goto DISCONNECT
			}
		}
	}

DISCONNECT:
	// wait until the Ready event is received.
	select {
	case <-timeout.C:
		t.Fatalf("test took longer than expected while receiving the Ready event")
	case <-s.Context.Done():
		t.Fatalf("disconnected before handling Ready event")
	case <-wait:
		break
	}

	// disconnect from the Discord Gateway (WebSocket Connection).
	if err := s.Disconnect(); err != nil {
		t.Fatalf("%v", err)
	}

	// disconnecting from a disconnected session should result in an error.
	if err := s.Disconnect(); err == nil {
		t.Fatalf("expected error while disconnecting from already disconnected session")
	}

	// allow Discord to close the session.
	<-time.After(time.Second * 5)
}

// TestReconnect tests Connect(), Disconnect(), heartbeat(), listen(), and onPayload()
// in order to ensure that WebSocket reconnection functionality works.
func TestReconnect(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// timeout is used to prevent this test from lasting longer than expected.
	timeout := time.NewTimer(time.Second * 8)

	bot := &Client{
		Authentication: BotToken(os.Getenv("TOKEN")),
		Config:         DefaultConfig(),
		Handlers:       &Handlers{},
	}

	s := NewSession()

	// a channel is used to wait for the resumed event which signals a reconnection.
	wait := make(chan int)

	// a Ready event is sent upon a successful connection.
	bot.Handle(FlagGatewayEventNameReady, func(*Ready) {
		wait <- 0
	})

	// a Resumed event is sent upon a successful reconnection.
	bot.Handle(FlagGatewayEventNameResumed, func(*Resumed) {
		wait <- 0
	})

	// connect to the Discord Gateway (WebSocket Session).
	if err := s.Connect(bot); err != nil {
		t.Fatalf("%v", err)
	}

	// wait until the Ready event is received.
	for {
		s.Lock()

		select {
		case <-timeout.C:
			t.Fatalf("test took longer than expected while receiving the Ready event")
		case <-s.Context.Done():
			t.Fatalf("disconnected before handling Ready event")
		case <-wait:
			s.Unlock()

			goto RECONNECT
		default:
			s.Unlock()
		}
	}

RECONNECT:
	// reconnect.
	if err := s.Reconnect(bot); err != nil {
		t.Fatalf("%v", err)
	}

	// wait until another Ready or Resumed event is handled by the bot.
	select {
	case <-timeout.C:
		t.Fatalf("test took too long to reconnect")
	case <-wait:
		break
	}

	// disconnect from the Discord Gateway (WebSocket Connection).
	if err := s.Disconnect(); err != nil {
		t.Fatalf("%v", err)
	}
}
