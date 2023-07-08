package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	. "github.com/switchupcb/disgo"
	"golang.org/x/sync/errgroup"
)

// TestRequestGlobalRateLimit tests the global rate limit mechanism (with the Default Bucket mechanism disabled)
// for HTTP requests.
func TestRequestGlobalRateLimit(t *testing.T) {
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

	// ensure that the next test starts with a full bucket.
	time.After(time.Second * 1)
}

// TestRequestRouteRateLimit tests the per-route rate limit mechanism (with the Default Bucket mechanism enabled)
// for HTTP requests.
func TestRequestRouteRateLimit(t *testing.T) {
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

	// ensure that the next test starts with a full bucket.
	time.After(time.Second * 1)
}

// TestGatewayIdentifyRateLimit tests the Identify rate limit mechanism for the Discord Gateway.
func TestGatewayIdentifyRateLimit(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	bot := &Client{
		Authentication: BotToken(os.Getenv("TOKEN")),
		Config:         DefaultConfig(),
		Handlers:       new(Handlers),
		Sessions:       NewSessionManager(),
	}

	// a counter is used to count the amount of Ready events.
	readyCount := 0

	// a Ready event is sent upon a successful connection.
	if err := bot.Handle(FlagGatewayEventNameReady, func(*Ready) {
		readyCount++
	}); err != nil {
		t.Fatalf("%v", err)
	}

	s1 := NewSession()
	s2 := NewSession()

	// call Connect at the same time.
	eg := new(errgroup.Group)

	eg.Go(func() error {
		// connect to the Discord Gateway (WebSocket Session).
		if err := s1.Connect(bot); err != nil {
			return fmt.Errorf("s1: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		// connect to the Discord Gateway (WebSocket Session).
		if err := s2.Connect(bot); err != nil {
			return fmt.Errorf("s2: %w", err)
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		t.Fatalf("%v", err)
	}

	// disconnect from the Discord Gateway (WebSocket Connection).
	if err := s1.Disconnect(); err != nil {
		t.Fatalf("s1: %v", err)
	}

	// disconnect from the Discord Gateway (WebSocket Connection).
	if err := s2.Disconnect(); err != nil {
		t.Fatalf("s2: %v", err)
	}

	if readyCount != 2 {
		t.Fatalf("expect to receive 2 Ready events but got %d", readyCount)
	}

	// allow Discord to close each session.
	<-time.After(time.Second * 5)
}
