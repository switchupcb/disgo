package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

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

const (
	// set cleanup to true to remove the created Application Command upon program termination.
	cleanup = false
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
	//
	// Use subcommand groups and subcommands to create a /calculate command (unusable),
	// which provides a subcommand group "add" (with usable commands "int", "string"),
	// and subcommand "subtract" (usable).
	//
	// COMMAND /calculate
	//  GROUP "add"
	//    SUBCOMMAND "int"     (adds two integers)
	//      OPTION 1: The first integer.
	//      OPTION 2: The second integer.
	//
	//  SUBCOMMAND "string"    (adds two strings)
	//    OPTION 1: The first string.
	//    OPTION 2: The second string.
	//
	//  SUBCOMMAND "subtract"  (adds two doubles)
	//    OPTION 1: The first double.
	//    OPTION 2: The second double.
	//
	request := &disgo.CreateGlobalApplicationCommand{
		Name:        "calculate",
		Description: "Calculate an operation using two values.",

		// create the `/calculate <group>` subcommand groups.
		Options: []*disgo.ApplicationCommandOption{

			// subcommand group `/calculate add`
			{
				Name:        "add",
				Description: "Add two values together.",
				Type:        disgo.FlagApplicationCommandOptionTypeSUB_COMMAND_GROUP,

				Options: []*disgo.ApplicationCommandOption{
					// create the `/calculate add int` subcommand.
					{
						Name:        "int",
						Description: "Add two integers together.",
						Type:        disgo.FlagApplicationCommandOptionTypeSUB_COMMAND,

						// create the OPTIONS for the `/calculate add int` subcommand.
						Options: []*disgo.ApplicationCommandOption{
							{
								Name:        "addend",
								Description: "The first integer to add.",
								Type:        disgo.FlagApplicationCommandOptionTypeINTEGER,
							},
							{
								Name:        "summand",
								Description: "The second integer to add.",
								Type:        disgo.FlagApplicationCommandOptionTypeINTEGER,
							},
						},
					},

					// create the `/calculate add string` subcommand.
					{
						Name:        "string",
						Description: "Concatenate two strings.",
						Type:        disgo.FlagApplicationCommandOptionTypeSUB_COMMAND,

						// create the OPTIONS for the `/calculate add string` subcommand.
						Options: []*disgo.ApplicationCommandOption{
							{
								Name:        "addend",
								Description: "The first string to concatenate.",
								Type:        disgo.FlagApplicationCommandOptionTypeSTRING,
							},
							{
								Name:        "summand",
								Description: "The second string to concatenate.",
								Type:        disgo.FlagApplicationCommandOptionTypeSTRING,
							},
						},
					},
				},
			},

			// subcommand `/calculate subtract`
			{
				Name:        "subtract",
				Description: "Subtract one value from another.",
				Type:        disgo.FlagApplicationCommandOptionTypeSUB_COMMAND,

				// create the OPTIONS for the `/calculate subtract` subcommand.
				Options: []*disgo.ApplicationCommandOption{
					{
						Name:        "minuend",
						Description: "The number to subtract from.",
						Type:        disgo.FlagApplicationCommandOptionTypeNUMBER,
						Required:    disgo.Pointer(true),
					},
					{
						Name:        "subtrahend",
						Description: "The number to subtract.",
						Type:        disgo.FlagApplicationCommandOptionTypeNUMBER,
						Required:    disgo.Pointer(true),

						// Explicitly defining every field of an option struct is not necessary.
						// In any case, a description of each field is found in the Discord API Documentation.
						// https://discord.com/developers/docs/interactions/application-commands#application-command-object

						NameLocalizations:        nil,
						DescriptionLocalizations: nil,
						Choices:                  nil,
						Options:                  nil,
						ChannelTypes:             nil,
						MinValue:                 nil,
						MaxValue:                 nil,
						Autocomplete:             nil,
					},
				},
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
	bot.Handle(disgo.FlagGatewayEventNameInteractionCreate, func(i *disgo.InteractionCreate) {
		log.Printf("calculate called by %s.", i.Interaction.User.Username)

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

	log.Println("Successfully connected to the Discord Gateway. Waiting for interactions...")

	// End the program using a SIGINT call via `Ctrl + C` from the terminal.
	//
	// a blocking call (<-signalChannel) is made to prevent the main thread from returning.
	//
	// The following code is equivalent to tools.InterceptSignal(tools.Signals, bot.Sessions...)
	interceptSIGINT(bot, newCommand)

	log.Printf("Program executed successfully.")
}

// onInteraction calculates an operation based on the user's input, then responds to the interaction.
//
// In this example, onInteraction is called when a user sends a `/calculate` interaction to the bot.
func onInteraction(bot *disgo.Client, interaction *disgo.Interaction) error {
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
	// Note: It's recommended to specify the amount explicitly when possible (i.e 5).
	// Otherwise, the tools package provides an easy way to determine the amount of options.
	optionMap := tools.OptionsToMap(nil, options, tools.NumOptions(options))

	// This example stores the calculation in the following variable.
	var calculation string

	// Perform the given calculation based on the user's input.
	//
	// Use the option name to select it from the map.
	if _, ok := optionMap["add"]; ok {
		// add two integers (in a safe manner).
		if group, ok := optionMap["int"]; ok && len(group.Options) != 0 {
			var sum int64

			for _, option := range group.Options {
				value, err := option.Value.Int64()
				if err != nil {
					return fmt.Errorf("An option expected to be an integer was not one: %w", err)
				}

				sum += value
			}

			calculation = strconv.FormatInt(sum, 10)
		}

		// concatenate two strings (in a safe manner).
		if group, ok := optionMap["string"]; ok {
			for _, option := range group.Options {
				calculation += option.Value.String()
			}
		}
	}

	// subtract two doubles.
	if _, ok := optionMap["subtract"]; ok {
		// The options for this command are required, but server-side validation is provided.
		minuend, ok := optionMap["minuend"]
		if !ok {
			return fmt.Errorf("The minuend was not provided.")
		}

		subtrahend, ok := optionMap["subtrahend"]
		if !ok {
			return fmt.Errorf("The subtrahend was not provided.")
		}

		minuendFloat, err := minuend.Value.Float64()
		if err != nil {
			return fmt.Errorf("An option expected to be a double (float) was not one: %w", err)
		}

		subtrahendFloat, err := subtrahend.Value.Float64()
		if err != nil {
			return fmt.Errorf("An option expected to be a double (float) was not one: %w", err)
		}

		difference := minuendFloat - subtrahendFloat
		calculation = fmt.Sprintf("%.3f", difference)
	}

	// calculation is empty when no subcommand options were provided.
	if calculation == "" {
		calculation = "There was nothing to do."
	}

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
				Content: disgo.Pointer(calculation), // or &calculation
			},
		},
	}

	if err := requestCreateInteractionResponse.Send(bot); err != nil {
		return fmt.Errorf("error sending interaction response: %w", err)
	}

	log.Println("Sent a response to the interaction.")

	return nil
}

// interceptSIGINT intercepts the SIGINT signal for a graceful end of the program.
func interceptSIGINT(bot *disgo.Client, command *disgo.ApplicationCommand) {
	// create an buffered channel (reason in goroutine below).
	signalChannel := make(chan os.Signal, 1)

	// set the syscalls that signalChannel is sent.
	// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
	signal.Notify(signalChannel,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	)

	// block this goroutine until a signal is received.
	// https://go.dev/tour/concurrency/3
	<-signalChannel

	log.Println("Exiting program due to signal...")

	if cleanup {
		log.Println("Deleting the application command...")

		// delete the Global Application Command.
		requestDeleteGlobalApplicationCommand := &disgo.DeleteGlobalApplicationCommand{CommandID: command.ID}
		if err := requestDeleteGlobalApplicationCommand.Send(bot); err != nil {
			log.Printf("error deleting Global Application Command: %v", err)
		}
	}

	log.Println("Disconnecting from the Discord Gateway...")

	// Disconnect the session from the Discord Gateway (WebSocket Connection).
	if err := bot.Sessions[0].Disconnect(); err != nil {
		log.Printf("error closing connection to Discord Gateway: %v", err)

		os.Exit(1)
	}
}
