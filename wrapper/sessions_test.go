package wrapper

import (
	"os"
	"sync/atomic"
	"testing"

	_ "github.com/joho/godotenv/autoload"
)

// TestConnect tests Connect(), Disconnect(), heartbeat(), listen(), and onPayload()
// in order to ensure that WebSocket functionality works.
func TestConnect(t *testing.T) {
	t.Parallel()

	bot := &Client{
		Authentication: BotToken(os.Getenv("TOKEN")),
		Config:         DefaultConfig(),
		Handlers:       &Handlers{},
	}

	s := NewSession()
	bot.Sessions = []*Session{s}

	// connect to the Discord Gateway (WebSocket Session).
	if err := bot.Connect(bot.Sessions[0]); err != nil {
		t.Fatalf("%v", err)
	}

	// wait until the Ready event is received.
	for {
		select {
		case <-s.Connected:
			t.Fatalf("disconnected before handling Ready event")
		default:
			s.mu.resource.RLock()

			// a Ready event handler sets s.ID upon receiving a Ready event.
			if s.ID != "" {
				s.mu.resource.RUnlock()

				goto HEARTBEAT
			}

			s.mu.resource.RUnlock()
		}
	}

HEARTBEAT:
	// once the Ready event has been received, wait to send a Heartbeat.
	for {
		select {
		case <-s.Connected:
			t.Fatalf("disconnected while waiting to send heartbeat")
		default:
			s.mu.heartbeat.Lock()

			// a Heartbeat is sent when the amount of ACKs since the
			// last Heartbeat was sent is 0.
			if atomic.LoadUint32(&s.heartbeat.acks) == 0 {
				s.mu.heartbeat.Unlock()

				goto ACK
			}

			s.mu.heartbeat.Unlock()
		}
	}

ACK:
	// once Heartbeat has been sent, wait to receive the respective HeartbeatACK.
	for {
		select {
		case <-s.Connected:
			t.Fatalf("disconnected while waiting for HeartbeatACK")
		default:
			s.mu.heartbeat.Lock()

			// a respective HeartbeatACK should be sent to s.heartbeat.ack.
			if atomic.LoadUint32(&s.heartbeat.acks) == 1 {
				s.mu.heartbeat.Unlock()

				goto DISCONNECT
			}

			s.mu.heartbeat.Unlock()
		}
	}

DISCONNECT:
	// disconnect from the Discord Gateway (WebSocket Connection).
	s.mu.connect.Lock()
	defer s.mu.connect.Unlock()

	if err := bot.Sessions[0].Disconnect(FlagClientCloseEventCodeNormal); err != nil {
		t.Fatalf("%v", err)
	}
}

// TestReconnect tests Connect(), Disconnect(), heartbeat(), listen(), and onPayload()
// in order to ensure that WebSocket reconnection functionality works.
func TestReconnect(t *testing.T) {
	t.Parallel()

	bot := &Client{
		Authentication: BotToken(os.Getenv("TOKEN")),
		Config:         DefaultConfig(),
		Handlers:       &Handlers{},
	}

	s := NewSession()
	bot.Sessions = []*Session{s}

	// connect to the Discord Gateway (WebSocket Session).
	if err := bot.Connect(bot.Sessions[0]); err != nil {
		t.Fatalf("%v", err)
	}

	// wait until the Ready event is received.
	for {
		select {
		case <-s.Connected:
			t.Fatalf("disconnected before handling Ready event")
		default:
			s.mu.resource.RLock()

			// a Ready event handler sets s.ID upon receiving a Ready event.
			if s.ID != "" {
				s.mu.resource.RUnlock()

				goto HEARTBEAT
			}

			s.mu.resource.RUnlock()
		}
	}

HEARTBEAT:
	// once the Ready event has been received, wait to send a Heartbeat.
	for {
		select {
		case <-s.Connected:
			t.Fatalf("disconnected while waiting to send heartbeat")
		default:
			s.mu.heartbeat.Lock()

			// a Heartbeat is sent when the amount of ACKs since the
			// last Heartbeat was sent is 0.
			if atomic.LoadUint32(&s.heartbeat.acks) == 0 {
				s.mu.heartbeat.Unlock()

				goto ACK
			}

			s.mu.heartbeat.Unlock()
		}
	}

ACK:
	// once Heartbeat has been sent, wait to receive the respective HeartbeatACK.
	for {
		select {
		case <-s.Connected:
			t.Fatalf("disconnected while waiting for HeartbeatACK")
		default:
			s.mu.heartbeat.Lock()

			// a respective HeartbeatACK should be sent to s.heartbeat.ack.
			if atomic.LoadUint32(&s.heartbeat.acks) == 1 {
				s.mu.heartbeat.Unlock()

				goto DISCONNECT
			}

			s.mu.heartbeat.Unlock()
		}
	}

DISCONNECT:
	// disconnect from the Discord Gateway (WebSocket Connection).
	s.mu.connect.Lock()
	defer s.mu.connect.Unlock()

	if err := bot.Sessions[0].Disconnect(FlagClientCloseEventCodeNormal); err != nil {
		t.Fatalf("%v", err)
	}
}

// TODO: manipulate s.heartbeat.ack on HEARTBEAT to test reconnect functionality.
