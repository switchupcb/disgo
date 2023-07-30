package main

import (
	"flag"
	"log"
	"os"

	"github.com/switchupcb/disgo"
)

// Environment Variables.
var (
	// token represents the bot's token.
	token = os.Getenv("TOKEN")
)

// Command Line Flags.
var (
	channelID = flag.String("c", "", "Set the channel (ID) the message components will be sent to using -c.")
)

func main() {
	// enable the logger for the API Wrapper.
	// zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// parse the command line flags.
	flag.Parse()

	// ensure that the program has the necessary data to succeed.
	if token == "" {
		log.Println("The bot's token must be set, but is currently empty.")

		return
	}

	if *channelID == "" {
		log.Println("The channel to send the message to is not set.")
		flag.Usage()

		return
	}

	// create a new Bot Client with the information required to send a request.
	bot := &disgo.Client{
		Authentication: disgo.BotToken(token), // or BearerToken("TOKEN")
		Config:         disgo.DefaultConfig(),
	}

	// ensure that the bot has access to the channel.
	//
	// This is useful for the validation of this program, but unnecessary.
	getChannelRequest := disgo.GetChannel{ChannelID: *channelID}
	_, err := getChannelRequest.Send(bot)
	if err != nil {
		log.Printf("error occurred getting channel %q: %v", *channelID, err)

		return
	}

	// send an Action Row (containing a Button) to the channel.
	createMessageRequestActionRowButton := &disgo.CreateMessage{ //nolint:exhaustruct
		ChannelID: *channelID,
		Content:   disgo.Pointer("This is an Action Row (containing a Button)."),
		Components: []disgo.Component{
			&disgo.ActionRow{
				Type: disgo.FlagComponentTypeActionRow,
				Components: []disgo.Component{
					&disgo.Button{
						Type:     disgo.FlagComponentTypeButton,
						Style:    disgo.FlagButtonStyleRED,
						Label:    disgo.Pointer("Button Label."),
						Emoji:    nil,
						CustomID: disgo.Pointer("example-button"),
						URL:      nil,
						Disabled: nil,
					},
				},
			},
		},
	}

	message, err := createMessageRequestActionRowButton.Send(bot)
	if err != nil {
		log.Printf("error occurred sending a message to channel %q: %v", *channelID, err)

		return
	}

	log.Printf("Successfully sent message with ID %q", message.ID)

	// send an Action Row (containing a Select Menu) to the channel.
	createMessageRequestActionRowSelectMenu := &disgo.CreateMessage{ //nolint:exhaustruct
		ChannelID: *channelID,
		Content:   disgo.Pointer("This is an Action Row (containing a String Select Menu)."),
		Components: []disgo.Component{
			&disgo.ActionRow{
				Type: disgo.FlagComponentTypeActionRow,
				Components: []disgo.Component{
					&disgo.SelectMenu{
						Type:        disgo.FlagComponentTypeStringSelect,
						CustomID:    "example-select-menu",
						Placeholder: disgo.Pointer("Select an option."),
						Options: []disgo.SelectMenuOption{
							{
								Label:       "Yes",
								Value:       "yes",
								Description: disgo.Pointer("Yessir."),
								Emoji:       nil,
								Default:     nil,
							},
							{
								Label:       "No",
								Value:       "no",
								Description: disgo.Pointer("Nope!"),
								Emoji:       nil,
								Default:     nil,
							},
						},
						MinValues:    nil,
						MaxValues:    nil,
						Disabled:     nil,
						ChannelTypes: nil,
					},
				},
			},
		},
	}

	message, err = createMessageRequestActionRowSelectMenu.Send(bot)
	if err != nil {
		log.Printf("error occurred sending a message to channel %q: %v", *channelID, err)

		return
	}

	log.Printf("Successfully sent message with ID %q", message.ID)
}
