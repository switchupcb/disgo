package integration_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	. "github.com/switchupcb/disgo/wrapper"
	"golang.org/x/sync/errgroup"
)

// TestGlobalRateLimit tests the global rate limit mechanism (with the Default Bucket mechanism disabled).
func TestGlobalRateLimit(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// setup the bot.
	bot := &Client{
		Authentication: BotToken(os.Getenv("TOKEN")),
		Config:         DefaultConfig(),
	}
	bot.Config.Request.Retries = 0
	bot.Config.Request.RateLimiter.SetDefaultBucket(nil)

	// prepare the request.
	request := new(GetCurrentBotApplicationInformation)
	requests := 101

	// prepare the test tracking variables.
	eg, ctx := errgroup.WithContext(context.Background())

	// send the requests concurrently.
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
	time.After(time.Second * 2)
}

// TestRouteRateLimit tests the per-route rate limit mechanism (with the Default Bucket mechanism enabled).
func TestRouteRateLimit(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// setup the bot.
	bot := &Client{
		Authentication: BotToken(os.Getenv("TOKEN")),
		Config:         DefaultConfig(),
	}
	bot.Config.Request.Retries = 0
	bot.Config.Request.RateLimiter.SetDefaultBucket(
		&Bucket{Limit: 1}, //nolint:exhaustruct
	)

	// prepare the request.
	request := GetUser{UserID: os.Getenv("APPID")}
	requests := 31

	// prepare the test tracking variables.
	eg, ctx := errgroup.WithContext(context.Background())

	// send the requests concurrently.
	for i := 1; i <= requests; i++ {
		select {
		case <-ctx.Done():
			t.Fatalf("%v", eg.Wait())
		default:
		}

		id := i
		eg.Go(func() error {
			t.Log("Spawned request goroutine", id)

			user, err := request.Send(bot)
			if err != nil {
				return err
			}

			t.Log("Request", id, ":", user)

			return nil
		})
	}

	// wait until all requests are sent and responses are received.
	if err := eg.Wait(); err != nil {
		t.Fatalf("%v", err)
	}
}
