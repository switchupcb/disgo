package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/switchupcb/disgo/tools"
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
	request := &disgo.CreateGlobalApplicationCommand{
		Name:        "followup",
		Description: disgo.Pointer("Showcase multiple types of interaction responses."),
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
		log.Printf("followup called by %s.", i.Interaction.User.Username)

		// see func declaration below.
		if err := onInteraction(bot, i.Interaction); err != nil {
			log.Println(err)
		}
	})

	log.Println("Connecting to the Discord Gateway...")

	// Connect the session to the Discord Gateway (WebSocket Connection).
	if err := bot.Sessions[0].Connect(bot); err != nil {
		log.Printf("can't open websocket session to Discord Gateway: %v", err)

		return
	}

	log.Println("Successfully connected to the Discord Gateway. Waiting for an interaction...")

	// end the program using a SIGINT call via `Ctrl + C` from the terminal.
	if err := tools.InterceptSignal(tools.Signals, bot.Sessions...); err != nil {
		log.Printf("error exiting program: %v", err)
	}

	log.Println("Deleting the application command...")

	// The following code is not necessarily required, but useful for the cleanup of this program.
	//
	// delete the Global Application Command.
	requestDeleteGlobalApplicationCommand := &disgo.DeleteGlobalApplicationCommand{CommandID: newCommand.ID}
	if err := requestDeleteGlobalApplicationCommand.Send(bot); err != nil {
		log.Printf("error deleting Global Application Command: %v", err)

		return
	}

	log.Printf("Program executed successfully.")
}

// onInteraction responds to the user in multiple ways.
//
// In this example, onInteraction is called when a user sends a `/followup` interaction to the bot.
func onInteraction(bot *disgo.Client, interaction *disgo.Interaction) error {
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
				Content: disgo.Pointer("This is a regular response. Wait..."),
			},
		},
	}

	if err := requestCreateInteractionResponse.Send(bot); err != nil {
		return fmt.Errorf("error sending original interaction response: %w", err)
	}

	log.Printf("Sent original interaction response.")

	// Wait 2 seconds, then edit the original response.
	timer := time.NewTimer(time.Second * 2)
	<-timer.C

	log.Println("Editing the original response to the interaction...")

	requestEditInteractionResponse := &disgo.EditOriginalInteractionResponse{
		ApplicationID:    bot.ApplicationID,
		InteractionToken: interaction.Token,
		Content:          disgo.Pointer2("This response is edited now."),
	}

	if _, err := requestEditInteractionResponse.Send(bot); err != nil {
		return fmt.Errorf("error sending edit interaction response: %w", err)
	}

	log.Printf("Edited original interaction response.")

	// Wait 2 seconds, then send a followup message.
	timer.Reset(time.Duration(time.Second * 2))
	<-timer.C

	log.Println("Sending a followup message to the interaction...")

	requestCreateFollowupMessage := &disgo.CreateFollowupMessage{
		ApplicationID:    bot.ApplicationID,
		InteractionToken: interaction.Token,
		Content:          disgo.Pointer("This is a followup message. Wait..."),
	}

	followupMessage, err := requestCreateFollowupMessage.Send(bot)
	if err != nil {
		return fmt.Errorf("error sending followup message interaction response: %w", err)
	}

	log.Println("Sent a followup message to the interaction.")

	// Wait 2 seconds, then edit the followup message.
	timer.Reset(time.Duration(time.Second * 2))
	<-timer.C

	log.Println("Editing the followup message to the interaction...")

	requestEditFollowupMessage := &disgo.EditFollowupMessage{
		ApplicationID:    bot.ApplicationID,
		InteractionToken: interaction.Token,
		MessageID:        followupMessage.ID,
		Content:          disgo.Pointer2("This followup message is edited now."),
	}

	if _, err := requestEditFollowupMessage.Send(bot); err != nil {
		log.Printf("Edited followup message to interaction response.")
	}

	log.Println("Edited the followup message to the interaction.")

	return nil
}
