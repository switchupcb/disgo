package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tools "github.com/switchupcb/disgo/tools"
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
		Name:        "autocomplete",
		Description: disgo.Pointer("Learn about autocompletion."),

		Options: []*disgo.ApplicationCommandOption{
			{
				Name:        "freewill",
				Description: "Do you have it?",
				Type:        disgo.FlagApplicationCommandOptionTypeSTRING,
				Required:    disgo.Pointer(true),

				// The following choices are the ONLY valid choices for this option.
				Choices: []*disgo.ApplicationCommandOptionChoice{
					{
						Name:  "Yes",
						Value: "y",
					},
					{
						Name:  "No",
						Value: "n",
					},
				},
			},
			{
				Name:        "confirm",
				Description: "Confirm your answer.",
				Type:        disgo.FlagApplicationCommandOptionTypeSTRING,
				Required:    disgo.Pointer(true),

				// Autocomplete choices will be provided, but the user can still
				// input any value for this option.
				Autocomplete: disgo.Pointer(true),
			},
		},
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
	//
	// ensure that the event handler is added to the bot.
	if err := bot.Handle(disgo.FlagGatewayEventNameInteractionCreate, func(i *disgo.InteractionCreate) {
		log.Println("Received interaction.")

		// see func declaration below.
		if err := onInteraction(bot, i.Interaction); err != nil {
			log.Println(err)
		}
	}); err != nil {
		// when the Handle(eventname, function) parameters are not configured correctly.
		log.Printf("Failed to add event handler to bot: %v", err)

		os.Exit(1)
	}

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

// onInteraction sends the user autocomplete choices.
//
// In this example, onInteraction is called when a user sends a `/autocomplete` interaction to the bot.
func onInteraction(bot *disgo.Client, interaction *disgo.Interaction) error {
	// User input will be returned as partial data to this bot application.
	//
	// check whether autocomplete data or a command submission has been received.
	switch interaction.Type {
	case disgo.FlagInteractionTypeAPPLICATION_COMMAND_AUTOCOMPLETE:
		log.Println("The interaction contains autocompletion data. Sending choices...")

		// access the interaction's options in the order provided by the user.
		//
		// The following commands assist with type converting Interaction Data
		// into the respective data structs.
		//
		// 	ApplicationCommand()
		//	MessageComponent()
		// 	ModalSubmit()
		//
		options := interaction.ApplicationCommand().Options

		// Alternatively, convert the option slice into a map (in an efficient manner).
		//
		// Note: It's recommended to specify the amount explicitly when possible (i.e 3).
		// Otherwise, the tools package provides an easy way to determine the amount of options.
		optionMap := tools.OptionsToMap(nil, options, tools.NumOptions(options))

		// determine the choices that will be sent to the user depending on the incoming options.
		//
		// Discord will highlight the closest choice to the user's input.
		choices := []*disgo.ApplicationCommandOptionChoice{
			{
				Name:  "Yes",
				Value: "y",
			},
			{
				Name:  "No",
				Value: "n",
			},
		}

		// Note: In order to receive autocompletion data, one option must be non-nil,
		// but this does NOT imply that the non-nil option is NOT empty (i.e "").
		freewill := optionMap["freewill"]
		confirm := optionMap["confirm"]

		// When the user is completing their first option, send both choices.
		if (freewill != nil && confirm == nil) || (freewill == nil && confirm != nil) {
			log.Println("Sending both choices...")

			// When the user has completed both options, determine which option is focused
			// to provide the correct choice.
		} else if freewill != nil && freewill.Focused != nil && *freewill.Focused {
			if confirm.Value.String() == "Yes" {
				choices = choices[1:] // n
			} else if confirm.Value.String() == "No" {
				choices = choices[0:1] // y
			}

			log.Println("Sending choice opposite of confirm...")
		} else if confirm != nil && confirm.Focused != nil && *confirm.Focused {
			if freewill.Value.String() == "y" {
				choices = choices[1:] // n
			} else {
				choices = choices[0:1] // y
			}

			log.Println("Sending choice opposite of freewill...")
		} else {
			log.Println("Unknown choice state encountered. Sending both choices...")
		}

		requestAutocomplete := &disgo.CreateInteractionResponse{
			InteractionID:    interaction.ID,
			InteractionToken: interaction.Token,
			InteractionResponse: &disgo.InteractionResponse{
				Type: disgo.FlagInteractionCallbackTypeAPPLICATION_COMMAND_AUTOCOMPLETE_RESULT,
				Data: &disgo.Autocomplete{
					Choices: choices,
				},
			},
		}

		if err := requestAutocomplete.Send(bot); err != nil {
			return fmt.Errorf("failure sending autocompletion command to Discord: %v", err)
		}

		log.Println("Sent choices.")

	case disgo.FlagInteractionTypeAPPLICATION_COMMAND:
		log.Println("The interaction contains command data. Sending interaction response...")

		// determine the response.
		options := interaction.ApplicationCommand().Options
		optionMap := tools.OptionsToMap(nil, options, tools.NumOptions(options))
		freewill := optionMap["freewill"]
		confirm := optionMap["confirm"]

		response := "Hmmm. I guess you aren't sure..."
		if freewill.Value.String() == "y" && strings.ToLower(confirm.Value.String()) == "yes" {
			response = "Where there's a will there's a way."
		} else if freewill.Value.String() == "n" && strings.ToLower(confirm.Value.String()) == "no" {
			response = "Fate awaits you."
		}

		// When the interaction is submitted, send an interaction response to reply to the user.
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
					Content: disgo.Pointer(response), // &response
				},
			},
		}

		if err := requestCreateInteractionResponse.Send(bot); err != nil {
			return fmt.Errorf("error sending interaction response: %w", err)
		}

		log.Println("Sent interaction response.")
	}

	return nil
}
