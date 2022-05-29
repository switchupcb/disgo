package wrapper

import (
	"fmt"
	"net/url"
	"time"

	json "github.com/goccy/go-json"
	"github.com/gorilla/schema"
	"github.com/valyala/fasthttp"
)

// timeout is a temporary variable that represents the amount of time
// a request will wait for a response.
// TODO: refactor with timeout usage.
var timeout time.Duration

// Send Error Messages.
const (
	ErrSendMarshal = "an error occurred while marshalling a %v:\n%w"
	ErrSendRequest = "an error occurred while sending %v:\n%w"
	ErrQueryString = "an error occurred creating a URL Query String for %v:\n%w"
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

var (
	// contentTypeURL represents an HTTP header indicating a payload with an encoded URL Query String.
	contentTypeURL = []byte("application/x-www-form-urlencoded")

	// contentTypeJSON represents an HTTP header that indicates a payload with a JSON body.
	contentTypeJSON = []byte("application/json")

	// contentTypeMulti represents an HTTP header that indicates a payload with a multi-part (file) body.
	contentTypeMulti = []byte("multipart/form-data")
)

// SendRequest sends a fasthttp.Request using the given URI, HTTP method, content type and body,
// then parses the response into dst.
func SendRequest(client *fasthttp.Client, method, uri string, content, body []byte, dst any) error {
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	request.SetRequestURI(uri)
	request.Header.SetMethod(method)
	request.Header.SetContentTypeBytes(content)
	request.SetBodyRaw(body)

	// receive the response from the request.
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)
	if err := client.DoTimeout(request, response, timeout); err != nil {
		return fmt.Errorf("%w", err)
	}

	// unmarshal the JSON response into dst.
	if response.StatusCode() == fasthttp.StatusOK {
		err := json.Unmarshal(response.Body(), dst)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		return nil
	}

	return StatusCodeError(response.StatusCode())
}

var (
	// qsEncoder is used to create URL Query Strings from objects.
	qsEncoder = schema.NewEncoder()

	// qsDecoder is used to parse URL Query Strings into objects.
	qsDecoder = schema.NewDecoder()
)

// init runs at the start of the program.
func init() {

	// use `url` tags for the URL Query String encoder and decoder.
	qsEncoder.SetAliasTag("url")
	qsDecoder.SetAliasTag("url")
}

// EndpointQueryString return a URL Query String from a given object.
func EndpointQueryString(dst any) (string, error) {
	params := url.Values{}
	err := qsEncoder.Encode(dst, params)
	if err != nil {
		return "", err
	}

	return params.Encode(), nil
}
