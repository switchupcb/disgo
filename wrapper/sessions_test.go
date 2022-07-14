package wrapper

import (
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

// TestConnect tests Connect(), Disconnect(), heartbeat(), listen(), and onPayload()
// in order to ensure that WebSocket functionality works.
func TestConnect(t *testing.T) {
	// timeout is used to prevent this test from lasting longer than expected.
	timeout := time.NewTimer(time.Second * 8)

	bot := &Client{
		Authentication: BotToken(os.Getenv("TOKEN")),
		Config:         DefaultConfig(),
		Handlers:       new(Handlers),
		Sessions:       []*Session{new(Session)},
	}

	s := bot.Sessions[0]

	// a channel is used to wait for the Ready event.
	wait := make(chan int)

	// a Ready event is sent upon a successful connection.
	bot.Handle(FlagGatewayEventNameReady, func(*Ready) {
		// once the Ready event has been received, wait to send a Heartbeat.
		//
		// NOTE: use jitter functionality to queue the first Heartbeat faster.
		s.mu.Lock()

		s.heartbeat.interval = time.Second
		s.heartbeat.ticker.Reset(s.heartbeat.interval)

		s.mu.Unlock()

		wait <- 0
	})

	// connect to the Discord Gateway (WebSocket Session).
	if err := bot.Connect(s); err != nil {
		t.Fatalf("%v", err)
	}

	// connecting to a connected session should result in an error.
	if err := bot.Connect(s); err == nil {
		t.Fatalf("expected error while connecting to already connected session")
	}

	// wait until the Ready event is received.
	select {
	case <-timeout.C:
		t.Fatalf("test took longer than expected while receiving the Ready event")
	case <-s.Context.Done():
		t.Fatalf("disconnected before handling Ready event")
	case <-wait:
		break
	}

	for {
		select {
		case <-timeout.C:
			t.Fatalf("test took longer than expected while sending heartbeat")
		case <-s.Context.Done():
			t.Fatalf("disconnected while waiting to send heartbeat")
		default:
			s.mu.Lock()

			// a Heartbeat is sent when the amount of ACKs since the
			// last Heartbeat was sent is 0.
			if atomic.LoadUint32(&s.heartbeat.acks) == 0 {
				s.mu.Unlock()

				goto ACK
			}

			s.mu.Unlock()
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
			s.mu.Lock()

			// a respective HeartbeatACK should be sent to s.heartbeat.ack.
			if atomic.LoadUint32(&s.heartbeat.acks) == 1 {
				s.mu.Unlock()

				goto DISCONNECT
			}

			s.mu.Unlock()
		}
	}

DISCONNECT:
	// disconnect from the Discord Gateway (WebSocket Connection).
	if err := bot.Sessions[0].Disconnect(); err != nil {
		t.Fatalf("%v", err)
	}

	// disconnecting from a disconnected session should result in an error.
	if err := bot.Sessions[0].Disconnect(); err == nil {
		t.Fatalf("expected error while disconnecting from already disconnected session")
	}
}

// TestReconnect tests Connect(), Disconnect(), heartbeat(), listen(), and onPayload()
// in order to ensure that WebSocket reconnection functionality works.
func TestReconnect(t *testing.T) {
	// timeout is used to prevent this test from lasting longer than expected.
	timeout := time.NewTimer(time.Second * 8)

	bot := &Client{
		Authentication: BotToken(os.Getenv("TOKEN")),
		Config:         DefaultConfig(),
		Handlers:       &Handlers{},
		Sessions:       []*Session{new(Session)},
	}

	s := bot.Sessions[0]

	// a channel is used to wait for the resumed event which signals a reconnection.
	wait := make(chan int)

	// a Ready event is sent upon a successful connection.
	bot.Handle(FlagGatewayEventNameReady, func(*Ready) {
		// manipulate s.heartbeat.ack to test reconnect functionality.
		s.mu.Lock()

		// simulate a Heartbeat not receiving a respective HeartbeatACK.
		//
		// NOTE: use jitter functionality to queue the first Heartbeat faster.
		s.heartbeat.acks = 0
		s.heartbeat.interval = time.Millisecond
		s.heartbeat.ticker.Reset(s.heartbeat.interval)

		s.mu.Unlock()

		wait <- 0
	})

	// a Resumed event is sent upon a successful reconnection.
	bot.Handle(FlagGatewayEventNameResumed, func(*Resumed) {
		wait <- 0
	})

	// connect to the Discord Gateway (WebSocket Session).
	if err := bot.Connect(s); err != nil {
		t.Fatalf("%v", err)
	}

	// wait until the Ready event is received.
	for {
		s.mu.Lock()

		select {
		case <-timeout.C:
			t.Fatalf("test took longer than expected while receiving the Ready event")
		case <-s.Context.Done():
			t.Fatalf("disconnected before handling Ready event")
		case <-wait:
			s.mu.Unlock()

			goto AGAIN
		default:
			s.mu.Unlock()
		}
	}

AGAIN:
	// wait until another Ready or Resumed event is handled by the bot.
	select {
	case <-timeout.C:
		t.Fatalf("test took too long to reconnect")
	case <-wait:
		break
	}

	// disconnect from the Discord Gateway (WebSocket Connection).
	if err := bot.Sessions[0].Disconnect(); err != nil {
		log.Println(err)
	}
}
