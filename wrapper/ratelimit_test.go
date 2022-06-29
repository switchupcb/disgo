package wrapper

import (
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
)

// TestGlobalRateLimit tests the global rate limit mechanism.
func TestGlobalRateLimit(t *testing.T) {
	// setup the bot.
	bot := &Client{
		Authentication: BotToken(os.Getenv("TOKEN")),
		Config:         DefaultConfig(),
	}
	bot.Config.Retries = 0

	// prepare the request.
	request := new(GetCurrentBotApplicationInformation)
	requests := 101

	// prepare the test tracking variables.
	errs := make(chan error)
	responses := make(chan *Application, requests)

	// send 51 requests concurrently.
	for i := 1; i <= requests; i++ {
		go func(id int) {
			t.Log("Spawned request goroutine", id)
			app, err := request.Send(bot)
			if err != nil {
				errs <- err
			}

			t.Log("Request", id, ":", app)
			responses <- app
		}(i)
	}

	// wait until all requests are sent.
	for {
		select {
		case err := <-errs:
			t.Fatalf("%v", err)
		default:
			if len(responses) == cap(responses) {
				return
			}
		}
	}
}
