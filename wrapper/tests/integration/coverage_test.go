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

	command := new(ApplicationCommand)

	// set the bot's Application ID.
	requestGetCurrentBotApplicationInformation := &GetCurrentBotApplicationInformation{}
	app, err := requestGetCurrentBotApplicationInformation.Send(bot)
	if err != nil {
		t.Fatal(fmt.Errorf("GetCurrentBotApplicationInformation: %w", err))
	}

	bot.ApplicationID = app.ID
	if app.ID == "" {
		t.Fatal("GetCurrentBotApplicationInformation: expected non-null Application ID")
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

		var err error
		if command, err = request.Send(bot); err != nil {
			return fmt.Errorf("failure sending CreateGlobalApplicationCommand to Discord: %w", err)
		}

		if command.ID == "" {
			return fmt.Errorf("CreateGlobalApplicationCommand: expected non-null Global Application Command")
		}

		return nil
	})

	eg.Go(func() error {
		request := &ListVoiceRegions{}
		regions, err := request.Send(bot)
		if err != nil {
			return fmt.Errorf("ListVoiceRegions: %w", err)
		}

		if len(regions) == 0 {
			return fmt.Errorf("ListVoiceRegions: expected non-empty Voice Regions Array")
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
	// commands.
	eg.Go(func() error {
		return testCommands(bot, command)
	})

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

// testCommands tests all endpoints that are dependent on a global command.
func testCommands(bot *Client, command *ApplicationCommand) error {
	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		commands, err := new(GetGlobalApplicationCommands).Send(bot)
		if err != nil {
			return fmt.Errorf("GetGlobalApplicationCommands: %w", err)
		}

		if len(commands) == 0 {
			return fmt.Errorf("GetGlobalApplicationCommands: expected non-empty Global Application Command List")
		}

		return nil
	})

	eg.Go(func() error {
		getGlobalApplicationCommand := &GetGlobalApplicationCommand{CommandID: command.ID}
		got, err := getGlobalApplicationCommand.Send(bot)
		if err != nil {
			return fmt.Errorf("GetGlobalApplicationCommand: %w", err)
		}

		if got == nil {
			return fmt.Errorf("GetGlobalApplicationCommand: expected non-null Global Application Command")
		}

		return nil
	})

	eg.Go(func() error {
		editGlobalApplicationCommand := &EditGlobalApplicationCommand{
			CommandID:   command.ID,
			Name:        "notmain",
			Description: "This is not a main global command.",
		}

		editedCommand, err := editGlobalApplicationCommand.Send(bot)
		if err != nil {
			return fmt.Errorf("EditGlobalApplicationCommand: %w", err)
		}

		if editedCommand.ID == "" {
			return fmt.Errorf("EditGlobalApplicationCommand: expected non-null Global Application Command")
		}

		return nil
	})

	// wait until all requests that depend on the main command have been processed.
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", eg.Wait())
	default:
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("%w", eg.Wait())
	}

	deleteGlobalApplicationCommand := &DeleteGlobalApplicationCommand{
		CommandID: command.ID,
	}

	if err := deleteGlobalApplicationCommand.Send(bot); err != nil {
		return fmt.Errorf("DeleteGlobalApplicationCommand: %w", err)
	}

	return nil
}
