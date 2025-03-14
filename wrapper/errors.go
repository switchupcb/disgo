package wrapper

import (
	"fmt"
)

// Send Request Error Messages.
const (
	errSendMarshal = "marshalling an HTTP body: %w"
	errUnmarshal   = "error unmarshalling into %T: %w"
	errRateLimit   = "error converting the HTTP rate limit header %q: %w"
)

// ErrorRequest represents an HTTP Request error that occurs when an attempt to send a request fails.
type ErrorRequest struct {
	// Err represents the error that occurred while performing the action.
	Err error

	// ClientID represents the Application ID of the request sender.
	ClientID string

	// CorrelationID represents the ID used to correlate the request to other logs.
	CorrelationID string

	// RouteID represents the ID (hash) of the Disgo Route.
	RouteID string

	// ResourceID represents the ID (hash) of the resource for the route.
	ResourceID string

	// Endpoint represents the endpoint the request was sent to.
	Endpoint string
}

func (e ErrorRequest) Error() string {
	return fmt.Errorf("REQUEST ERROR: client %q: x %q: route: %q: resource: %q: endpoint: %q: error: %w",
		e.ClientID, e.CorrelationID, e.RouteID, e.ResourceID, e.Endpoint, e.Err).Error()
}

// ErrorStatusCode represents an HTTP Request error that occurs when an unexpected response is returned.
type ErrorStatusCode struct {
	// StatusCode represents the HTTP Status Code received from a response.
	StatusCode int
}

// Status Code Error Messages.
const (
	errStatusCodeKnown   = "status code %d: %v"
	errStatusCodeUnknown = "status code %d: unknown status code error from Discord"
)

func (e ErrorStatusCode) Error() string {
	return fmt.Sprintf("STATUS CODE ERROR: status code: %d: msg: %v", e.StatusCode, StatusCodeError(e.StatusCode))
}

// StatusCodeError returns the relevant message for a Discord API HTTP Status Code.
func StatusCodeError(status int) string {
	if msg, ok := HTTPResponseCodes[status]; ok {
		return fmt.Sprintf(errStatusCodeKnown, status, msg)
	}

	return fmt.Sprintf(errStatusCodeUnknown, status)
}

// JSON Error Code Messages.
const (
	errJSONErrorCodeKnown = "JSON Error Code %d: %v"
	errJSONErrorUnknown   = "JSON Error Code %d: Unknown JSON Error Code from Discord"
)

// JSONCodeError handles a Discord API JSON Error Code and returns the relevant error message.
func JSONCodeError(status int) error {
	if msg, ok := JSONErrorCodes[status]; ok {
		return fmt.Errorf(errJSONErrorCodeKnown, status, msg)
	}

	return fmt.Errorf(errJSONErrorUnknown, status)
}

// Event Handler Error Messages.
const (
	errHandleNotRemoved   = "event handler was not added"
	errRemoveInvalidIndex = "event handler cannot be removed since there is no event handler at index %d"
)

// ErrorEventHandler represents an Event Handler error that occurs when an attempt to
// add or remove an event handler fails.
type ErrorEventHandler struct {
	// Err represents the error that occurred while performing the action.
	Err error

	// ClientID represents the Application ID of the event handler owner.
	ClientID string

	// Event represents the event of the involved handler.
	Event string
}

func (e ErrorEventHandler) Error() string {
	return fmt.Errorf("EVENT HANDLER ERROR: client %q: event %q: error: %w",
		e.ClientID, e.Event, e.Err).Error()
}

const (
	ErrorEventActionUnmarshal = "unmarshalling"
	ErrorEventActionMarshal   = "marshalling"
	ErrorEventActionRead      = "reading"
	ErrorEventActionWrite     = "writing"
)

// ErrorEvent represents a WebSocket error that occurs when an attempt to {action} an event fails.
type ErrorEvent struct {
	// ClientID represents the Application ID of the event handler caller.
	ClientID string

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
	return fmt.Errorf("EVENT ERROR: client %q: event %q: action: %q, error: %w",
		e.ClientID, e.Event, e.Action, e.Err).Error()
}

// ErrorSession represents a WebSocket Session error that occurs during an active session.
type ErrorSession struct {
	// Err represents the error that occurred.
	Err error

	// SessionID represents the ID of the Session.
	SessionID string

	// State represents the state of the session.
	State string

	// Type represents the type of connection (e.g., Discord Gateway, Discord Voice).
	Type string
}

const (
	ErrorSessionTypeGateway = "Discord Gateway"
	ErrorSessionTypeVoice   = "Discord Voice"
)

func (e ErrorSession) Error() string {
	return fmt.Errorf("SESSION ERROR: %q session %q: state: %q error: %w", e.Type, e.SessionID, e.State, e.Err).Error()
}

// ErrorSessionDisconnect represents a disconnection error that occurs when
// an attempt to gracefully disconnect from a session fails.
type ErrorSessionDisconnect struct {
	// Action represents the error that prompted the disconnection (if applicable).
	Action error

	// Err represents the error that occurred while disconnecting.
	Err error
}

func (e ErrorSessionDisconnect) Error() string {
	return fmt.Errorf(
		"\tDisconnect(): %v\n"+
			"\treason: %w\n",
		e.Err, e.Action,
	).Error() //lint:ignore ST1005 readability
}
