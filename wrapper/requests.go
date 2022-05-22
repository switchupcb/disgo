package wrapper

import (
	"fmt"
	"time"

	json "github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
)

// timeout is a temporary variable that represents the amount of time
// a request will wait for a response.
// TODO: refactor with timeout usage.
var timeout time.Duration

// TODO: Set HTTP METHOD (Get, Post, Put, Patch, Delete) [using Copygen Options]
const (
	TODO = fasthttp.MethodPost
)

// Status Code Error Messages.
const (
	ErrStatusCodeKnown   = "Status Code %d: %v"
	ErrStatusCodeUnknown = "Status Code %d: Unknown status code error from Discord"
)

// StatusCodeError handles a Discord API HTTP Status Code and returns the relevant error.
func StatusCodeError(status int) error {
	if msg, ok := JSONErrorCodes[status]; ok {
		return fmt.Errorf(ErrStatusCodeKnown, status, msg)
	}

	return fmt.Errorf(ErrStatusCodeUnknown, status)
}

// ContentTypeJSON represents an HTTP header that indicates a JSON body.
var ContentTypeJSON = []byte("application/json")

// SendRequest sends a fasthttp.Request with a JSON body using the given URI, method, and body,
// then parses the response into dst.
func SendRequest(dst any, client *fasthttp.Client, method, uri string, body []byte) error {
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	request.SetRequestURI(uri)
	request.Header.SetMethod(method)
	request.Header.SetContentTypeBytes(ContentTypeJSON)
	request.SetBodyRaw(body)

	// receive the response from the request.
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)
	if err := client.DoTimeout(request, response, timeout); err != nil {
		return fmt.Errorf("%w", err)
	}

	// unmarshal the response into dst.
	if response.StatusCode() == fasthttp.StatusOK {
		err := json.Unmarshal(response.Body(), dst)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		return nil
	}

	return StatusCodeError(response.StatusCode())
}

// Send Error Messages.
const (
	ErrSendMarshal = "an error occurred while marshalling a %v: \n%w"
	ErrSendRequest = "an error occurred while sending %v: \n%w"
)
