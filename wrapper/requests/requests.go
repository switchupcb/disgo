// Package requests provides Discord API Requests.
package requests

import (
	"fmt"
	"time"

	json "github.com/goccy/go-json"
	"github.com/switchupcb/disgo/wrapper/resources"
	"github.com/valyala/fasthttp"
)

// Client is ONLY defined here for the purpose of the proof of concept.
type Client struct {
	ApplicationID resources.Snowflake
	client        *fasthttp.Client
	ctx           *fasthttp.RequestCtx
	timeout       time.Duration
}

// HandleStatus handles a Discord API HTTP Status Code and returns the relevant error.
func HandleStatus(status int) error {
	if msg, ok := resources.JSONErrorCodes[status]; ok {
		return fmt.Errorf("Status Code %d: %v", status, msg)
	}

	return fmt.Errorf("Status Code %d: Unknown JSON error from Discord", status)
}

// ParseResponse parses the response of a Discord API Request with a JSON Body into dst.
func ParseResponseJSON(ctx *fasthttp.RequestCtx, dst any) error {
	if ctx.Response.StatusCode() == fasthttp.StatusOK {
		err := json.Unmarshal(ctx.Response.Body(), dst)
		if err != nil {
			fasthttp.ReleaseResponse(&ctx.Response)
			return err
		}
	} else {
		err := HandleStatus(ctx.Response.StatusCode())
		fasthttp.ReleaseResponse(&ctx.Response)
		return err
	}
	fasthttp.ReleaseResponse(&ctx.Response)

	return nil
}
