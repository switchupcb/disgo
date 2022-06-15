package wrapper

import (
	"errors"
	"fmt"
)

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

/// Gateway Errors
const (
	errRemoveInvalidEventHandler = "cannot remove event handler for %s since there is no event handler at index %d"
)

var (
	errOpcodeReconnect = errors.New("Opcode Reconnect")
)

// ErrorDisconnect represents a WebSocket disconnection error that occurs when
// a WebSocket attempt to gracefully disconnect fails.
type ErrorDisconnect struct {
	// SessionID represents the ID of the Session that is disconnecting.
	SessionID string

	// Err represents the error that occurred while disconnecting.
	Err error

	// Action represents the error that prompted the disconnection (if applicable).
	Action error
}

func (e ErrorDisconnect) Error() string {
	return fmt.Sprintf("an error occurred disconnecting the session %s from the Discord Gateway\n"+
		"\tDisconnect(): %v\n"+
		"\treason: %v\n",
		e.SessionID, e.Err, e.Action)
}

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
	return fmt.Sprintf("an error occurred while %s a %v Event\n\t%v", e.Action, e.Event, e.Err)
}

const (
	ErrorEventActionUnmarshal = "unmarshalling"
	ErrorEventActionMarshal   = "marshalling"
	ErrorEventActionRead      = "reading"
	ErrorEventActionWrite     = "writing"
)

// ErrorEventHandler represents an Event Handler error that occurs when
// an attempt to {action} an event handler fails.
type ErrorEventHandler struct {
	// Event represents the name of the event handler involved in this error.
	Event string

	// Err represents the error that occurred while performing the action.
	Err error

	// Action represents the action that prompted the error.
	//
	// ErrorEventHandlerAction can be four values:
	// ErrorEventHandlerAdd:       an error occurred while adding an event handler.
	// ErrorEventHandlerRemove:    an error occurred while removing an event handler.
	Action string
}

func (e ErrorEventHandler) Error() string {
	return fmt.Sprintf("an error occurred while %s a %s event handler\n\t%v", e.Action, e.Event, e.Err)
}

const (
	ErrorEventHandlerAdd    = "adding"
	ErrorEventHandlerRemove = "removing"
)
