package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/rs/zerolog"
	. "github.com/switchupcb/disgo/wrapper"
)

// TestCoverage tests 100+ endpoints (requests) and respective events from the Discord API.
func TestCoverage(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	bot := &Client{
		Authentication: BotToken(os.Getenv("COVERAGE_TEST_TOKEN")),
		Config:         DefaultConfig(),
		Handlers:       new(Handlers),
		Sessions:       []*Session{NewSession()},
	}

	// set the bot's Application ID.
	requestGetCurrentBotApplicationInformation := &GetCurrentBotApplicationInformation{}
	app, err := requestGetCurrentBotApplicationInformation.Send(bot)
	if app.ID == "" {
		t.Fatal("GetCurrentBotApplicationInformation: expected non-null Application ID")
	}

	bot.ApplicationID = app.ID
	if err != nil {
		t.Fatal(fmt.Errorf("GetCurrentBotApplicationInformation: %w", err))
	}

	initializeEventHandlers(bot)

	// Connect the session to the Discord Gateway (WebSocket Connection).
	if err := bot.Sessions[0].Connect(bot); err != nil {
		t.Fatalf("can't open websocket session to Discord: %v", err)
	}

	eg, ctx := errgroup.WithContext(context.Background())

	// Call endpoints with no dependencies.
	//
	// create a global application command.
	eg.Go(func() error {
		request := CreateGlobalApplicationCommand{
			Name:        "main",
			Description: "A basic command",
		}

		newCommand, err := request.Send(bot)
		if err != nil {
			return fmt.Errorf("failure sending command to Discord: %v", err)
		}

		if newCommand.ID == "" {
			return fmt.Errorf("CreateGlobalApplicationCommand: expected non-null Global Application Command")
		}

		if err != nil {
			return fmt.Errorf("CreateGlobalApplicationCommand: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		request := &ListVoiceRegions{}
		regions, err := request.Send(bot)
		if len(regions) == 0 {
			return fmt.Errorf("ListVoiceRegions: expected non-empty Voice Regions Array")
		}

		if err != nil {
			return fmt.Errorf("ListVoiceRegions: %w", err)
		}

		return nil
	})

	// wait until all required requests have been processed.
	select {
	case <-ctx.Done():
		t.Fatalf("%v", eg.Wait())
	default:
	}

	if err := eg.Wait(); err != nil {
		t.Fatalf("%v", err)
	}

	// Call endpoints with one or more dependencies.
	//
	// var guild *Guild

	// Disconnect the session from the Discord Gateway (WebSocket Connection).
	if err := bot.Sessions[0].Disconnect(); err != nil {
		t.Fatalf("%v", err)
	}

	// allow Discord to close the session.
	<-time.After(time.Second * 5)
}

// initializeEventHandlers initializes the event handlers necessary for this test.
func initializeEventHandlers(bot *Client) {

}
