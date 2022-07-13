package wrapper

import (
	"log"
	"os"
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

	s := &Session{mu: new(sessionMutex)}
	bot.Sessions = []*Session{s}

	// heartbeats represents the amount of heartbeat intervals to test.
	const heartbeats = 1

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
	t.Fatalf("TODO: slow thread would cause this case to be potentially impossible.")

	// once the Ready event has been received, wait to send a Heartbeat.
	for {
		select {
		case <-s.Connected:
			t.Fatalf("disconnected while waiting to send heartbeat")
		default:
			s.mu.heartbeat.Lock()

			// a Heartbeat is sent when s.heartbeat.send is full.
			if len(s.heartbeat.send) == 1 {
				s.mu.heartbeat.Unlock()

				goto ACK
			}

			s.mu.heartbeat.Unlock()
		}
	}

ACK:
	log.Println("ack")

	// once Heartbeat has been sent, wait to receive the respective HeartbeatACK.
	for {
		select {
		case <-s.Connected:
			t.Fatalf("disconnected while waiting for HeartbeatACK")
		default:
			s.mu.heartbeat.Lock()

			// a respective HeartbeatACK should be sent to s.heartbeat.ack.
			if len(s.heartbeat.ack) == 1 {
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
