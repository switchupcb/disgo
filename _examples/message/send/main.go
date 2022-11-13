package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	disgo "github.com/switchupcb/disgo/wrapper"
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
	// parse the command line flags.
	flag.Parse()

	// ensure that the program has the necessary data to succeed.
	if token == "" {
		fmt.Println("The bot's token must be set, but is currently empty.")

		return
	}

	if *channelID == "" {
		fmt.Println("The channel to send the message to is not set.")
		flag.Usage()

		return
	}

	if *msg == "" && *location == "" {
		fmt.Println("The message has no content to send. Set the message, file location, or both.")
		flag.Usage()

		return
	}

	var files []disgo.File
	if *location != "" {
		file, err := getFile(*location)
		if err != nil {
			fmt.Printf("an error occurred getting the file: %v", err)

			return
		}

		// The following code is equivalent to `files = []*disgo.File{file}`.
		files = make([]disgo.File, 0)
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
		fmt.Printf("error occurred getting channel %q: %v", *channelID, err)

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
		fmt.Printf("error occurred sending a message to channel %q: %v", *channelID, err)

		return
	}

	fmt.Printf("Successfully sent message with ID %q", message.ID)
}

// getFile returns a disgo.File for usage in a message.
func getFile(location string) (*disgo.File, error) {
	// Determine if the provided location is a filepath or URL,
	// by checking whether the file exists.
	var isFile bool
	if _, err := os.Stat(location); err == nil {
		isFile = true
	}

	// In order to upload a file, you should set the File's Name, Content Type, and Data accordingly.
	var data []byte
	var err error
	switch isFile {
	case true:
		// when the location is a filepath, load the file locally.
		data, err = ioutil.ReadFile(location)
		if err != nil {
			return nil, fmt.Errorf("an error occurred reading the file: %v", err)
		}

	case false:
		// when the location is a URL, fetch the file from the internet.
		response, err := http.Get(location)
		if err != nil {
			return nil, fmt.Errorf("an error occurred fetching the file: %v", err)
		}

		// read the HTTP Response Body ([]bytes) to determine the file's Content Type and Data.
		defer response.Body.Close()

		data, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("an error occurred reading the fetched file data: %v", err)
		}
	}

	return &disgo.File{
		Name:        path.Base(location),
		ContentType: http.DetectContentType(data),
		Data:        data,
	}, nil
}
