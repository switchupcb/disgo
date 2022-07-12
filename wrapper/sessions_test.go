package wrapper

import (
	"os"
	"testing"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func TestDataRace(t *testing.T) {
	bot := &Client{
		// Set the Authentication Header using BotToken() or BearerToken().
		Authentication: BotToken(os.Getenv("TOKEN")),
		Authorization:  &Authorization{},
		Config:         DefaultConfig(),
		Handlers:       &Handlers{},
	}

	s := &Session{
		mu: new(sessionMutex),
	}

	bot.Sessions = []*Session{s}
	err := bot.Connect(bot.Sessions[0])
	if err != nil {
		t.Fatalf("%v", err)
	}

	const heartbeats = 1
	for {

		if s.Connected == nil {
			t.Fatalf("disconnected before calling disconnect")
		}

		s.mu.resource.Lock()
		if s.ID != "" {
			s.mu.resource.Unlock()
			// wait for a heartbeat to test heartbeat functionality.
			stop := time.Now().Add(heartbeats * s.heartbeat.interval)
			for time.Now().Before(stop) {
				if s.Connected == nil {
					t.Fatalf("Session disconnected while waiting for heartbeat")
				}
			}

			err = bot.Sessions[0].Disconnect(1000)
			if err != nil {
				t.Fatalf("%v", err)
			}

			return
		}
		s.mu.resource.Unlock()

	}

}
