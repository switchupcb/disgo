package wrapper

import (
	"fmt"

	json "github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
)

// Client is ONLY defined here for the purpose of the proof of concept.
type Client struct {
	ApplicationID string
	client        *fasthttp.Client
	ctx           *fasthttp.RequestCtx
}

const (
	ErrStatusCodeKnown   = "Status Code %d: %v"
	ErrStatusCodeUnknown = "Status Code %d: Unknown JSON error from Discord"
)

// StatusCodeError handles a Discord API HTTP Status Code and returns the relevant error.
func StatusCodeError(status int) error {
	if msg, ok := JSONErrorCodes[status]; ok {
		return fmt.Errorf(ErrStatusCodeKnown, status, msg)
	}

	return fmt.Errorf(ErrStatusCodeUnknown, status)
}

// ParseResponse parses the response of a Discord API Request with a JSON Body into dst.
func ParseResponseJSON(ctx *fasthttp.RequestCtx, dst any) error {
	defer fasthttp.ReleaseResponse(&ctx.Response)

	if ctx.Response.StatusCode() == fasthttp.StatusOK {
		err := json.Unmarshal(ctx.Response.Body(), dst)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		return nil
	}

	return StatusCodeError(ctx.Response.StatusCode())
}
