package wrapper

import (
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

// TestConnect tests Connect(), Disconnect(), heartbeat(), listen(), and onPayload()
// in order to ensure that WebSocket functionality works.
func TestConnect(t *testing.T) {
	bot := &Client{
		Authentication: BotToken(os.Getenv("TOKEN")),
		Config:         DefaultConfig(),
		Handlers:       &Handlers{},
	}

	s := new(Session)
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
			s.mu.RLock()

			// a Ready event handler sets s.ID upon receiving a Ready event.
			if s.ID != "" {
				s.mu.RUnlock()

				goto HEARTBEAT
			}

			s.mu.RUnlock()
		}
	}

HEARTBEAT:
	// once the Ready event has been received, wait to send a Heartbeat.
	//
	// NOTE: use jitter functionality to queue the first Heartbeat faster.
	s.muHeartbeat.Lock()

	s.heartbeat.interval = time.Second
	s.heartbeat.ticker.Reset(s.heartbeat.interval)

	s.muHeartbeat.Unlock()

	for {
		select {
		case <-s.Connected:
			t.Fatalf("disconnected while waiting to send heartbeat")
		default:
			s.muHeartbeat.Lock()

			// a Heartbeat is sent when the amount of ACKs since the
			// last Heartbeat was sent is 0.
			if atomic.LoadUint32(&s.heartbeat.acks) == 0 {
				s.muHeartbeat.Unlock()

				goto ACK
			}

			s.muHeartbeat.Unlock()
		}
	}

ACK:
	// once Heartbeat has been sent, wait to receive the respective HeartbeatACK.
	for {
		select {
		case <-s.Connected:
			t.Fatalf("disconnected while waiting for HeartbeatACK")
		default:
			s.muHeartbeat.Lock()

			// a respective HeartbeatACK should be sent to s.heartbeat.ack.
			if atomic.LoadUint32(&s.heartbeat.acks) == 1 {
				s.muHeartbeat.Unlock()

				goto DISCONNECT
			}

			s.muHeartbeat.Unlock()
		}
	}

DISCONNECT:
	// disconnect from the Discord Gateway (WebSocket Connection).
	if err := bot.Sessions[0].Disconnect(FlagClientCloseEventCodeNormal); err != nil {
		t.Fatalf("%v", err)
	}
}

// TestReconnect tests Connect(), Disconnect(), heartbeat(), listen(), and onPayload()
// in order to ensure that WebSocket reconnection functionality works.
func TestReconnect(t *testing.T) {
	bot := &Client{
		Authentication: BotToken(os.Getenv("TOKEN")),
		Config:         DefaultConfig(),
		Handlers:       &Handlers{},
	}

	s := new(Session)
	bot.Sessions = []*Session{s}

	// a channel is used to wait for the resumed event which signals a reconnection.
	wait := make(chan int)

	// a Resumed event is sent upon a successful reconnection.
	bot.Handle(FlagGatewayEventNameResumed, func(*Resumed) {
		// disconnect from the Discord Gateway (WebSocket Connection).
		if err := bot.Sessions[0].Disconnect(FlagClientCloseEventCodeNormal); err != nil {
			log.Println(err)
			wait <- 1

			return
		}

		wait <- 0
	})

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
			s.mu.RLock()

			// a Ready event handler sets s.ID upon receiving a Ready event.
			if s.ID != "" {
				s.mu.RUnlock()

				goto MANIPULATE
			}

			s.mu.RUnlock()
		}
	}

MANIPULATE:
	// once the Ready event has been received, manipulate s.heartbeat.ack
	// to test reconnect functionality.
	s.muHeartbeat.Lock()

	// simulate a Heartbeat not receiving a respective HeartbeatACK.
	//
	// NOTE: use jitter functionality to queue the first Heartbeat faster.
	s.heartbeat.acks = 0
	s.heartbeat.interval = time.Millisecond
	s.heartbeat.ticker.Reset(s.heartbeat.interval)

	s.muHeartbeat.Unlock()

	// wait until a reconnect is triggered by waiting until s.Connected is closed
	// after a heartbeat interval triggers a heartbeat.
	<-time.After(time.Second)

	select {
	case <-s.Connected:
		break
	default:
		t.Fatalf("test did not disconnect after heartbeat ack manipulation")
	}

	// timeout is used to prevent this test from lasting longer than expected.
	timeout := time.NewTimer(time.Second * 5)

	// wait until a Resumed event is handled by the bot.
	for {
		select {
		case <-timeout.C:
			t.Fatalf("test took too long to reconnect")
		case exit := <-wait:
			switch exit {
			case 0:
				return
			case 1:
				t.Fatalf("an error occurred while disconnecting after reconnecting")
			}
		}
	}
}
