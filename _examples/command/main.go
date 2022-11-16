package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	disgo "github.com/switchupcb/disgo/wrapper"
)

// Environment Variables
var (
	// token represents the bot's token.
	token = os.Getenv("TOKEN")

	// appid represents the bot's ApplicationID.
	//
	// Use Developer Mode to find it, or call GetCurrentUser (request) in your program
	// and set it programmatically.
	appid = os.Getenv("APPID")
)

var (
	// This program uses a sync.WaitGroup to prevent an immediate exit after the session is connected.
	wg sync.WaitGroup
)

func main() {
	// enable the logger for the API Wrapper.
	// zerolog.SetGlobalLevel(zerolog.DebugLevel)

	log.Println("Program is started.")

	// create a new Bot Client.
	bot := &disgo.Client{
		ApplicationID:  appid,                 // REQUIRED (for this example).
		Authentication: disgo.BotToken(token), // or BearerToken("TOKEN")
		Config:         disgo.DefaultConfig(),
		Handlers:       new(disgo.Handlers),
		Sessions:       []*disgo.Session{disgo.NewSession()},
	}

	log.Println("Creating an application command...")

	// Create a Create Global Application Command request.
	request := disgo.CreateGlobalApplicationCommand{
		Name:        "main",
		Description: "A basic command.",
	}

	// Register the new command by sending the request to Discord using the bot.
	//
	// returns a disgo.ApplicationCommand
	newCommand, err := request.Send(bot)
	if err != nil {
		log.Printf("failure sending command to Discord: %v", err)

		return
	}

	log.Println("Adding an event handler.")

	// Add an event handler to the bot.
	bot.Handle(disgo.FlagGatewayEventNameInteractionCreate, func(i *disgo.InteractionCreate) {
		log.Printf("main called by %s.", i.Interaction.User.Username)

		// see func declaration below.
		if err := onInteraction(bot, i.Interaction, newCommand); err != nil {
			log.Println(err)
		}

		// This call unblocks the main goroutine of the program from wg.Wait() [Line 100].
		wg.Done()
	})

	// add a tick to the WaitGroup counter.
	wg.Add(1)

	log.Println("Connecting to the Discord Gateway...")

	// Connect the session to the Discord Gateway (WebSocket Connection).
	if err := bot.Sessions[0].Connect(bot); err != nil {
		log.Printf("can't open websocket session to Discord Gateway: %v", err)

		return
	}

	log.Println("Successfully connected to the Discord Gateway. Waiting for an interaction...")

	// wg.Wait() unblocks once the event handler (defined above) calls wg.Done(),
	// which removes a tick from the wait group counter (such that the counter = 0).
	//
	// Alternatively, end the program using a SIGINT call via `Ctrl + C` from the terminal.
	//
	// The following code is equivalent to tools.InterceptSignal(tools.Signals, bot.Sessions...)
	interceptSIGINT(bot.Sessions[0])
	wg.Wait()

	log.Printf("Program executed successfully.")
}

// onInteraction deletes the Global Application Command, then disconnects the bot.
//
// In this example, onInteraction is called when a user sends a `/main` interaction to the bot.
func onInteraction(bot *disgo.Client, interaction *disgo.Interaction, command *disgo.ApplicationCommand) error {
	log.Println("Creating a response to the interaction...")

	// send an interaction response to reply to the user.
	requestCreateInteractionResponse := &disgo.CreateInteractionResponse{
		InteractionID: interaction.ID,

		// Interaction tokens are valid for 15 minutes,
		// but an initial response to an interaction must be sent within 3 seconds (of receiving it),
		// otherwise the token is invalidated.
		InteractionToken: interaction.Token,

		// https://discord.com/developers/docs/interactions/receiving-and-responding#responding-to-an-interaction
		InteractionResponse: &disgo.InteractionResponse{
			Type: disgo.FlagInteractionCallbackTypeCHANNEL_MESSAGE_WITH_SOURCE,

			// Any of the following objects can be used.
			//
			// Messages
			// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-messages
			//
			// Autocomplete
			// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-autocomplete
			//
			// Modal
			// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-response-object-modal
			Data: &disgo.Messages{
				Content: disgo.Pointer("Hello!"),
			},
		},
	}

	if err := requestCreateInteractionResponse.Send(bot); err != nil {
		return fmt.Errorf("error sending interaction response: %w", err)
	}

	log.Println("Deleting the application command...")

	// The following code is not necessarily required, but useful for the cleanup of this program.
	//
	// delete the Global Application Command.
	requestDeleteGlobalApplicationCommand := &disgo.DeleteGlobalApplicationCommand{CommandID: command.ID}
	if err := requestDeleteGlobalApplicationCommand.Send(bot); err != nil {
		return fmt.Errorf("error deleting Global Application Command: %w", err)
	}

	log.Println("Disconnecting from the Discord Gateway...")

	// Disconnect the session from the Discord Gateway (WebSocket Connection).
	if err := bot.Sessions[0].Disconnect(); err != nil {
		return fmt.Errorf("error closing connection to Discord Gateway: %w", err)
	}

	return nil
}

// interceptSIGINT intercepts the SIGINT signal for a graceful end of the program.
func interceptSIGINT(session *disgo.Session) {
	// create an buffered channel (reason in goroutine below).
	signalChannel := make(chan os.Signal, 1)

	// set the syscalls that signalChannel is sent.
	// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
	signal.Notify(signalChannel,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGKILL,
		syscall.SIGHUP,
	)

	// spawn a goroutines.
	// https://go.dev/tour/concurrency/1
	go func() {
		// block this goroutine until a signal is received.
		// https://go.dev/tour/concurrency/3
		<-signalChannel

		log.Println("Exiting program due to signal...")

		// Disconnect the session from the Discord Gateway (WebSocket Connection).
		if err := session.Disconnect(); err != nil {
			log.Printf("error closing connection to Discord Gateway: %v", err)

			os.Exit(0)
		}

		log.Println("Program exited successfully.")

		os.Exit(0)
	}()
}
