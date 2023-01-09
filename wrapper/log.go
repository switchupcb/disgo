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
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

var (
	// Logger represents the Disgo Logger used to log information.
	Logger = zerolog.New(os.Stdout)
)

// Logger Contexts
const (
	// LogCtxClient represents the log key for a Client (Bot) Application ID.
	LogCtxClient = "client"

	// LogCtxCorrelation represents the log key for a Correlation ID.
	LogCtxCorrelation = "xid"

	// LogCtxRequest represents the log key for a Request ID.
	LogCtxRequest = "request"

	// LogCtxRoute represents the log key for a Route ID.
	LogCtxRoute = "route"

	// LogCtxResource represents the log key for a Resource ID.
	LogCtxResource = "resource"

	// LogCtxEndpoint represents the log key for an Endpoint.
	LogCtxEndpoint = "endpoint"

	// LogCtxRequestBody represents the log key for an HTTP Request Body.
	LogCtxRequestBody = "body"

	// LogCtxBucket represents the log key for a Rate Limit Bucket ID.
	LogCtxBucket = "bucket"

	// LogCtxReset represents the log key for a Discord Bucket reset time.
	LogCtxReset = "reset"

	// LogCtxResponse represents the log key for an HTTP Request Response.
	LogCtxResponse = "response"

	// LogCtxResponseHeader represents the log key for an HTTP Request Response header.
	LogCtxResponseHeader = "header"

	// LogCtxResponseBody represents the log key for an HTTP Request Response body.
	LogCtxResponseBody = "body"

	// LogCtxRequestRateLimitCode represents the log key for an HTTP Rate Limit Response code.
	LogCtxRequestRateLimitCode = "code"

	// LogCtxSession represents the log key for a Discord Session ID.
	LogCtxSession = "session"

	// LogCtxPayload represents the log key for a Discord Gateway Payload.
	LogCtxPayload = "payload"

	// LogCtxPayloadOpcode represents the log key for a Discord Gateway Payload opcode.
	LogCtxPayloadOpcode = "opcode"

	// LogCtxPayloadData represents the log key for Discord Gateway Payload data.
	LogCtxPayloadData = "data"

	// LogCtxEvent represents the log key for a Discord Gateway Event.
	LogCtxEvent = "event"

	// LogCtxCommand represents the log key for a Discord Gateway command.
	LogCtxCommand = "command"

	// LogCtxCommandOpcode represents the log key for a Discord Gateway command opcode.
	LogCtxCommandOpcode = "opcode"

	// LogCtxCommandName represents the log key for a Discord Gateway command name.
	LogCtxCommandName = "name"
)

// LogRequest logs a request.
func LogRequest(log *zerolog.Event, clientid, xid, routeid, resourceid, endpoint string) *zerolog.Event {
	return log.Timestamp().
		Str(LogCtxClient, clientid).
		Dict(LogCtxRequest, zerolog.Dict().
			Str(LogCtxCorrelation, xid).
			Str(LogCtxRoute, routeid).
			Str(LogCtxResource, resourceid).
			Str(LogCtxEndpoint, endpoint),
		)
}

// LogRequestBody logs a request with its body.
func LogRequestBody(log *zerolog.Event, clientid, xid, routeid, resourceid, endpoint, body string) *zerolog.Event {
	return log.Timestamp().
		Str(LogCtxClient, clientid).
		Dict(LogCtxRequest, zerolog.Dict().
			Str(LogCtxCorrelation, xid).
			Str(LogCtxRoute, routeid).
			Str(LogCtxResource, resourceid).
			Str(LogCtxEndpoint, endpoint).
			Str(LogCtxRequestBody, body),
		)
}

// LogResponse logs a response (typically using LogRequest).
func LogResponse(log *zerolog.Event, header, body string) *zerolog.Event {
	return log.Dict(LogCtxResponse, zerolog.Dict().
		Str(LogCtxResponseHeader, header).
		Str(LogCtxResponseBody, body),
	)
}

// LogEventHandler logs an event handler action.
func LogEventHandler(log *zerolog.Event, clientid, event string) *zerolog.Event {
	return log.Timestamp().
		Str(LogCtxClient, clientid).
		Str(LogCtxEvent, event)
}

// LogSession logs a session.
func LogSession(log *zerolog.Event, sessionid string) *zerolog.Event {
	return log.Timestamp().
		Str(LogCtxSession, sessionid)
}

// LogPayload logs a Discord Gateway Payload (typically using LogSession).
func LogPayload(log *zerolog.Event, op int, data json.RawMessage) *zerolog.Event {
	return log.Dict(LogCtxPayload, zerolog.Dict().
		Int(LogCtxPayloadOpcode, op).
		Bytes(LogCtxPayloadData, data),
	)
}

// LogCommand logs a Gateway Command (typically using a LogSession).
func LogCommand(log *zerolog.Event, clientid string, op int, command string) *zerolog.Event {
	return log.Str(LogCtxClient, clientid).
		Dict(LogCtxCommand, zerolog.Dict().
			Int(LogCtxCommandOpcode, op).
			Str(LogCtxCommandName, command),
		)
}
