package wrapper

import (
	"os"
	"testing"

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

	err := bot.Connect(&Session{mu: new(sessionMutex)})
	if err != nil {
		t.Errorf("%v", err)
	}
}
