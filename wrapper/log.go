package wrapper

import (
	"encoding/json"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// init is called at the start of the application.
func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
}

var (
	// Logger represents the Disgo Logger used to log information.
	Logger = zerolog.New(os.Stdout).Level(zerolog.Disabled)
)

// Logger Contexts
var (
	// logCtxClient represents the log key for a Client (Bot) Application ID.
	logCtxClient = "client"

	// logCtxCorrelation represents the log key for a Correlation ID.
	logCtxCorrelation = "xid"

	// logCtxRequest represents the log key for a Request ID.
	logCtxRequest = "request"

	// logCtxRoute represents the log key for a Route ID.
	logCtxRoute = "route"

	// logCtxResource represents the log key for a Resource ID.
	logCtxResource = "resource"

	// logCtxEndpoint represents the log key for an Endpoint.
	logCtxEndpoint = "endpoint"

	// logCtxBucket represents the log key for a Rate Limit Bucket ID.
	logCtxBucket = "bucket"

	// logCtxReset represents the log key for a Discord Bucket reset time.
	logCtxReset = "reset"

	// logCtxResponse represents the log key for an HTTP Request Response.
	logCtxResponse = "response"

	// logCtxResponseHeader represents the log key for an HTTP Request Response header.
	logCtxResponseHeader = "header"

	// logCtxResponseBody represents the log key for an HTTP Request Response body.
	logCtxResponseBody = "body"

	// logCtxSession represents the log key for a Discord Session ID.
	logCtxSession = "session"

	// logCtxPayload represents the log key for a Discord Gateway Payload.
	logCtxPayload = "payload"

	// logCtxPayloadOpcode represents the log key for a Discord Gateway Payload opcode.
	logCtxPayloadOpcode = "opcode"

	// logCtxPayloadData represents the log key for Discord Gateway Payload data.
	logCtxPayloadData = "data"

	// logCtxEvent represents the log key for a Discord Gateway Event.
	logCtxEvent = "event"

	// logCtxCommand represents the log key for a Discord Gateway command.
	logCtxCommand = "command"

	// logCtxCommandOpcode represents the log key for a Discord Gateway command opcode.
	logCtxCommandOpcode = "opcode"

	// logCtxCommandName represents the log key for a Discord Gateway command name.
	logCtxCommandName = "name"
)

// logRequest logs a request.
func logRequest(log *zerolog.Event, clientid, xid, routeid, resourceid, endpoint string) *zerolog.Event {
	return log.Timestamp().
		Str(logCtxClient, clientid).
		Dict(logCtxRequest, zerolog.Dict().
			Str(logCtxCorrelation, xid).
			Str(logCtxRoute, routeid).
			Str(logCtxResource, resourceid).
			Str(logCtxEndpoint, endpoint),
		)
}

// logResponse logs a response (typically using a logRequest).
func logResponse(log *zerolog.Event, header, body string) *zerolog.Event {
	return log.Dict(logCtxResponse, zerolog.Dict().
		Str(logCtxResponseHeader, header).
		Str(logCtxResponseBody, body),
	)
}

// logEventHandler logs an event handler action.
func logEventHandler(log *zerolog.Event, clientid, event string) *zerolog.Event {
	return log.Timestamp().
		Str(logCtxClient, clientid).
		Str(logCtxEvent, event)
}

// logSession logs a session.
func logSession(log *zerolog.Event, sessionid string) *zerolog.Event {
	return log.Timestamp().
		Str(logCtxSession, sessionid)
}

// logPayload logs a Discord Gateway Payload (typically using a logSession).
func logPayload(log *zerolog.Event, op int, data json.RawMessage) *zerolog.Event {
	return log.Dict(logCtxPayload, zerolog.Dict().
		Int(logCtxPayloadOpcode, op).
		Bytes(logCtxPayloadData, data),
	)
}

// logCommand logs a Gateway Command (typically using a logSession).
func logCommand(log *zerolog.Event, clientid string, op int, command string) *zerolog.Event {
	return log.Str(logCtxClient, clientid).
		Dict(logCtxCommand, zerolog.Dict().
			Int(logCtxCommandOpcode, op).
			Str(logCtxCommandName, command),
		)
}
