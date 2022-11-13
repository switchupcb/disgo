package wrapper

import (
	"fmt"
)

// Send Request Error Messages.
const (
	ErrSendMarshal = "error while marshalling a %v:\n%w"
	ErrSendRequest = "error while sending %v:\n%w"
	ErrQueryString = "error creating a URL Query String for %v:\n%w"
	ErrMultipart   = "error creating multipart form: %w"
	ErrRedirect    = "error redirecting from %v due to a missing Location HTTP header"
	ErrRateLimit   = "error converting the HTTP rate limit header %v\n\t%w"
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

// Event Handler Error Messages.
const (
	errHandleNotRemoved   = "event handler for %s was not added"
	errRemoveInvalidIndex = "event handler for %s cannot be removed since there is no event handler at index %d"
)

// ErrorEvent represents a WebSocket error that occurs when
// a WebSocket attempt to {action} an event fails.
type ErrorEvent struct {
	// Event represents the name of the event involved in this error.
	Event string

	// Err represents the error that occurred while performing the action.
	Err error

	// Action represents the action that prompted the error.
	//
	// ErrorEventAction's can be one of four values:
	// ErrorEventActionUnmarshal: an error occurred while unmarshalling the Event from a JSON.
	// ErrorEventActionMarshal:   an error occurred while marshalling the Event to a JSON.
	// ErrorEventActionRead:      an error occurred while reading the Event from a Websocket Connection.
	// ErrorEventActionWrite:     an error occurred while writing the Event to a Websocket Connection.
	Action string
}

func (e ErrorEvent) Error() string {
	return fmt.Sprintf("error while %s a %v Event\n\t%v", e.Action, e.Event, e.Err)
}

const (
	ErrorEventActionUnmarshal = "unmarshalling"
	ErrorEventActionMarshal   = "marshalling"
	ErrorEventActionRead      = "reading"
	ErrorEventActionWrite     = "writing"
)

// DisconnectError represents a WebSocket disconnection error that occurs when
// a WebSocket attempt to gracefully disconnect fails.
type DisconnectError struct {
	// SessionID represents the ID of the Session that is disconnecting.
	SessionID string

	// Err represents the error that occurred while disconnecting.
	Err error

	// Action represents the error that prompted the disconnection (if applicable).
	Action error
}

func (e DisconnectError) Error() string {
	return fmt.Sprintf("error disconnecting the session %q from the Discord Gateway\n"+
		"\tDisconnect(): %v\n"+
		"\treason: %v\n",
		e.SessionID, e.Err, e.Action)
}
