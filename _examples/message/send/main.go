package main

import (
	"flag"
	"log"
	"os"

	"github.com/switchupcb/disgo"
)

// Environment Variables
var (
	// token represents the bot's token.
	token = os.Getenv("TOKEN")
)

// Command Line Flags
var (
	channelID = flag.String("c", "", "Set the channel (ID) the message will be sent to using -c.")
	msg       = flag.String("m", "", "Set the text content of the message using -m.")
	location  = flag.String("f", "", "Set the location (filepath or URL) of the file using -f.")
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

	if *msg == "" && *location == "" {
		log.Println("The message has no content to send. Set the message, file location, or both.")
		flag.Usage()

		return
	}

	var files []*disgo.File
	if *location != "" {
		file, err := getFile(*location)
		if err != nil {
			log.Printf("an error occurred getting the file: %v", err)

			return
		}

		// The following code is equivalent to `files = []*disgo.File{file}`.
		files = make([]*disgo.File, 0)
		files = append(files, file)
	}

	// create a new Bot Client with the information required to send a request.
	bot := &disgo.Client{
		Authentication: disgo.BotToken(token), // or BearerToken("TOKEN")
		Config:         disgo.DefaultConfig(),
	}

	// ensure that the bot has access to the channel.
	//
	// This is useful for the validation of this program, but not necessary.
	getChannelRequest := disgo.GetChannel{ChannelID: *channelID}
	_, err := getChannelRequest.Send(bot)
	if err != nil {
		log.Printf("error occurred getting channel %q: %v", *channelID, err)

		return
	}

	// send a message in the channel.
	//
	// Explicitly defining every field of a Create Message struct is not necessary.
	// In any case, a description of each field is found in the Discord API Documentation.
	// https://discord.com/developers/docs/resources/channel#create-message-jsonform-params
	createMessageRequest := disgo.CreateMessage{
		ChannelID:        *channelID,
		Content:          msg,
		Nonce:            nil,
		TTS:              nil,
		Embeds:           nil,
		AllowedMentions:  nil,
		MessageReference: nil,
		Components:       nil,
		StickerIDS:       nil,
		Files:            files,
		Attachments:      nil,
		Flags:            nil,
	}

	message, err := createMessageRequest.Send(bot)
	if err != nil {
		log.Printf("error occurred sending a message to channel %q: %v", *channelID, err)

		return
	}

	log.Printf("Successfully sent message with ID %q", message.ID)
}
