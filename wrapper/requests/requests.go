// Package requests provides Discord API Requests.
package requests

import (
	"fmt"
	"time"

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
