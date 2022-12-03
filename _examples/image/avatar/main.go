package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/switchupcb/disgo"
	"github.com/switchupcb/disgo/tools"
)

// Environment Variables
var (
	// token represents the bot's token.
	token = os.Getenv("TOKEN")
)

// Command Line Flags
var (
	location = flag.String("i", "", "Set the location (filepath or URL) of the avatar image using -i.")
	remove   = flag.Bool("r", false, "Use -r to remove the avatar image after successfully setting it.")
)

func main() {
	// parse the command line flags.
	flag.Parse()

	// ensure that the program has the necessary data to succeed.
	if token == "" {
		fmt.Println("The bot's token must be set, but is currently empty.")

		return
	}

	if *location == "" {
		fmt.Println("The avatar's location must be set, but is currently empty.")
		flag.Usage()

		return
	}

	// Determine if the provided location is a filepath or URL,
	// by checking whether the file exists.
	var isFile bool
	if _, err := os.Stat(*location); err == nil {
		isFile = true
	}

	// In order to upload the file, you must determine the Image's Content Type
	// and Data URI scheme.
	//
	// https://discord.com/developers/docs/reference#image-data
	var image []byte
	var err error
	switch isFile {
	case true:
		// when the location is a filepath, load the avatar image locally.
		image, err = ioutil.ReadFile(*location)
		if err != nil {
			fmt.Printf("an error occurred reading the avatar image: %v", err)

			return
		}

	case false:
		// when the location is a URL, fetch the avatar image from the internet.
		response, err := http.Get(*location)
		if err != nil {
			fmt.Printf("an error occurred fetching the avatar image: %v", err)

			return
		}

		// read the HTTP Response Body ([]bytes) to determine the
		// Image's Content Type and Data URI scheme.
		defer response.Body.Close()

		image, err = ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("an error occurred reading the fetched avatar image: %v", err)

			return
		}
	}

	// format the Image Data into a Data URI.
	//
	// The following code is equivalent to tools.ImageDataURI(image).
	contentType := http.DetectContentType(image)
	encodedImage := base64.StdEncoding.EncodeToString(image)
	data := tools.DataURI(contentType, encodedImage)

	// create a new Bot Client with the information required to send a request.
	bot := &disgo.Client{
		Authentication: disgo.BotToken(token), // or BearerToken("TOKEN")
		Config:         disgo.DefaultConfig(),
	}

	// create a ModifyCurrentUser request.
	request := disgo.ModifyCurrentUser{
		Username: nil,
		Avatar:   &data,
	}

	// update the bot's avatar by sending the request to Discord using the bot.
	user, err := request.Send(bot)
	if err != nil {
		fmt.Printf("an error occurred updating the bot's avatar: %v", err)
		fmt.Printf("content type: %s", contentType)

		return
	}

	fmt.Printf("Successfully updated the avatar of %s#%s.\n", user.Username, user.Discriminator)
	fmt.Println("Find it at " + disgo.CDNEndpointUserAvatar(user.ID, *user.Avatar))

	if *remove {
		fmt.Println("Removing avatar...")
		request := disgo.ModifyCurrentUser{
			Username: nil,
			Avatar:   disgo.Pointer(""),
		}

		user, err := request.Send(bot)
		if err != nil {
			fmt.Printf("an error occurred removing the bot's avatar: %v", err)

			return
		}

		fmt.Printf("Successfully removed the avatar of %s#%s.\n", user.Username, user.Discriminator)
	}
}
