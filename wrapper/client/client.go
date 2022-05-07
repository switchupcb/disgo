// Package client represents a Discord Application.
package client

import (
	"github.com/switchupcb/disgo/wrapper/resources"
	"github.com/valyala/fasthttp"
)

// Client is ONLY defined here for the purpose of the proof of concept.
type Client struct {
	ApplicationID resources.Snowflake
	client        *fasthttp.Client
	ctx           *fasthttp.RequestCtx
	Token         string
}
