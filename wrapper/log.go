package wrapper

import (
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
	Logger = zerolog.New(os.Stdout).Level(zerolog.DebugLevel)
)

// Logger Contexts
var (
	// logCtxClient represents the log key for a Client (Bot) Application ID.
	logCtxClient = "client"

	// logCtxBucket represents the log key for a Rate Limit Bucket ID.
	logCtxBucket = "bucket"

	// logCtxReset represents the log key for a Discord Bucket reset time.
	logCtxReset = "reset"

	// logCtxRequest represents the log key for a Request ID.
	logCtxRequest = "request"

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

	// logCtxCommand represents the log key for a Discord Gateway command.
	logCtxCommand = "command"

	// logCtxCommandOpcode represents the log key for a Discord Gateway command opcode.
	logCtxCommandOpcode = "opcode"

	// logCtxCommandName represents the log key for a Discord Gateway command name.
	logCtxCommandName = "name"
)
