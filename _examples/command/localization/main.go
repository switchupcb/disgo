package main

import (
	"fmt"
	"log"
	"os"

	"github.com/switchupcb/disgo"
	"github.com/switchupcb/disgo/tools"
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
		Name: "hello",

		// Be sure to adhere to Application Command Naming rules.
		// https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-naming
		NameLocalizations: &map[string]string{
			// Top 10. Locales by Population
			//
			// Discord doesn't support every locale.
			// https://discord.com/developers/docs/reference#locales
			disgo.FlagLocalesEnglishUS:    "hello",
			disgo.FlagLocalesEnglishUK:    "mate",
			disgo.FlagLocalesChineseChina: "你好",
			disgo.FlagLocalesHindi:        "नमस्ते",
			disgo.FlagLocalesSpanish:      "hola",
			disgo.FlagLocalesFrench:       "bonjour",
			// 6. Arabic
			// 7. Bengali
			disgo.FlagLocalesRussian:             "привет",
			disgo.FlagLocalesPortugueseBrazilian: "olá",
			// 10. Indonesian
			// 11. Pronouns
		},

		Description: disgo.Pointer("Say hello."),
		DescriptionLocalizations: &map[string]string{
			disgo.FlagLocalesEnglishUS:           "Say hello.",
			disgo.FlagLocalesEnglishUK:           "Say hello.",
			disgo.FlagLocalesChineseChina:        "问好。",
			disgo.FlagLocalesHindi:               "नमस्ते बोलो।",
			disgo.FlagLocalesSpanish:             "Di hola.",
			disgo.FlagLocalesFrench:              "Dis bonjour.",
			disgo.FlagLocalesRussian:             "Скажи привет.",
			disgo.FlagLocalesPortugueseBrazilian: "Diga olá.",
		},

		// Localization is also supported in Application Command Options.
		// https://discord.com/developers/docs/interactions/application-commands#localization
		Options: nil,
	}

	// Register the new command by sending the request to Discord using the bot.
	//
	// returns a disgo.ApplicationCommand
	newCommand, err := request.Send(bot)
	if err != nil {
		log.Printf("failure sending command to Discord: %v", err)

		return
	}

	// save the map defined in the returned command from Discord for later usage.
	if newCommand.NameLocalizations == nil {
		log.Println("error: returned command from Discord does not contain localizations.")

		return
	}

	locales := newCommand.NameLocalizations

	log.Println("Adding an event handler.")

	// Add an event handler to the bot.
	//
	// ensure that the event handler is added to the bot.
	if err := bot.Handle(disgo.FlagGatewayEventNameInteractionCreate, func(i *disgo.InteractionCreate) {
		log.Printf("hello called by %s.", i.Interaction.User.Username)

		// see func declaration below.
		if err := onInteraction(bot, i.Interaction, *locales); err != nil {
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

// onInteraction responds to the user based on their locale.
//
// In this example, onInteraction is called when a user sends a `/hello` interaction to the bot.
func onInteraction(bot *disgo.Client, interaction *disgo.Interaction, locales map[string]string) error {
	log.Println("Creating a response to the interaction...")

	// determine the response.
	var locale string
	if interaction.Locale != nil {
		locale = locales[*interaction.Locale]
	}

	if locale == "" {
		locale = "The current locale is not supported by this command."
	}

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
				Content: disgo.Pointer(locale), // or &locale
			},
		},
	}

	if err := requestCreateInteractionResponse.Send(bot); err != nil {
		return fmt.Errorf("error sending interaction response: %w", err)
	}

	log.Println("Sent a response to the interaction.")

	return nil
}
