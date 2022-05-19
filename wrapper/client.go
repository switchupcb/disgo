package wrapper

import "github.com/valyala/fasthttp"

// Client represents a Discord Application.
type Client struct {
	ApplicationID string
	client        *fasthttp.Client
}
