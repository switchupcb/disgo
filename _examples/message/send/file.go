package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/switchupcb/disgo"
)

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
			return nil, fmt.Errorf("an error occurred reading the file: %w", err)
		}

	case false:
		// when the location is a URL, fetch the file from the internet.
		response, err := http.Get(location) //nolint:gosec
		if err != nil {
			return nil, fmt.Errorf("an error occurred fetching the file: %w", err)
		}

		// read the HTTP Response Body ([]bytes) to determine the file's Content Type and Data.
		defer response.Body.Close()

		data, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("an error occurred reading the fetched file data: %w", err)
		}
	}

	return &disgo.File{
		Name:        path.Base(location),
		ContentType: http.DetectContentType(data),
		Data:        data,
	}, nil
}
