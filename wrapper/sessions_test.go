package wrapper

import (
	"testing"
)

const key = "OTg4NjA5NTMwMzc0MDA0NzU3.GreAn5.KayNltB_tollJDOj19t69sW51e9C6C_M5ES5mY"

func TestDataRace(t *testing.T) {
	bot := &Client{
		// Set the Authentication Header using BotToken() or BearerToken().
		Authentication: BotToken(key),
		Authorization:  &Authorization{},
		Config:         DefaultConfig(),
		Handlers:       &Handlers{},
	}

	err := bot.Connect(&Session{mu: new(sessionMutex)})
	if err != nil {
		t.Errorf("%v", err)
	}
}
