package wrapper

import "fmt"

// Send Error Messages.
const (
	ErrSendMarshal = "an error occurred while marshalling a %v:\n%w"
	ErrSendRequest = "an error occurred while sending %v:\n%w"
	ErrQueryString = "an error occurred creating a URL Query String for %v:\n%w"
	ErrRedirect    = "an error occurred redirecting from %v due to a missing Location HTTP header"
)

// Status Code Error Messages.
const (
	ErrStatusCodeKnown   = "Status Code %d: %v"
	ErrStatusCodeUnknown = "Status Code %d: Unknown status code error from Discord"
)

// StatusCodeError handles a Discord API HTTP Status Code and returns the relevant error message.
func StatusCodeError(status int) error {
	if msg, ok := HTTPResponseCodes[status]; ok {
		return fmt.Errorf(ErrStatusCodeKnown, status, msg)
	}

	return fmt.Errorf(ErrStatusCodeUnknown, status)
}

// JSON Error Code Messages.
const (
	ErrJSONErrorCodeKnown = "JSON Error Code %d: %v"
	ErrJSONErrorUnknown   = "JSON Error Code %d: Unknown JSON Error Code from Discord"
)

// JSONCodeError handles a Discord API JSON Error Code and returns the relevant error message.
func JSONCodeError(status int) error {
	if msg, ok := JSONErrorCodes[status]; ok {
		return fmt.Errorf(ErrJSONErrorCodeKnown, status, msg)
	}

	return fmt.Errorf(ErrJSONErrorUnknown, status)
}
