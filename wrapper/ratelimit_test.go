package wrapper

import (
	"context"
	"os"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

// TestGlobalRateLimit tests the global rate limit mechanism.
func TestGlobalRateLimit(t *testing.T) {
	// setup the bot.
	bot := &Client{
		Authentication: BotToken(os.Getenv("TOKEN")),
		Config:         DefaultConfig(),
	}
	bot.Config.Retries = 0
	DefaultRouteBucket.Limit = FlagGlobalRequestRateLimit

	// prepare the request.
	request := new(GetCurrentBotApplicationInformation)
	requests := 151

	// prepare the test tracking variables.
	eg, ctx := errgroup.WithContext(context.Background())

	// send 51 requests concurrently.
	for i := 1; i <= requests; i++ {
		select {
		case <-ctx.Done():
			t.Fatalf("%v", eg.Wait())
		default:
		}

		id := i
		eg.Go(func() error {
			t.Log("Spawned request goroutine", id)

			app, err := request.Send(bot)
			if err != nil {
				return err
			}

			t.Log("Request", id, ":", app)

			return nil
		})
	}

	// wait until all requests are sent and responses are received.
	if err := eg.Wait(); err != nil {
		t.Fatalf("%v", err)
	}

	// ensure that the next test starts from a full bucket.
	time.After(time.Second)
}

// TestRouteRateLimit tests the per-route rate limit mechanism.
func TestRouteRateLimit(t *testing.T) {

}
