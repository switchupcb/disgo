package http

import (
	"time"

	"github.com/valyala/fasthttp"
)

// HTTP Methods
const (
	GET    = fasthttp.MethodGet
	POST   = fasthttp.MethodPost
	PUT    = fasthttp.MethodPut
	PATCH  = fasthttp.MethodPatch
	DELETE = fasthttp.MethodDelete
)

// ContentTypeJSON represents an HTTP header that indicates a JSON body.
var ContentTypeJSON = []byte("application/json")

// timeout is a temporary variable that represents the amount of time
// a request will wait for a response.
// TODO: refactor with timeout usage.
var timeout time.Duration

// ReleaseResponse is a wrapper for the fasthttp.ReleaseResponse func.
func ReleaseResponse(ctx *fasthttp.RequestCtx) {
	fasthttp.ReleaseResponse(&ctx.Response)
}

// SendRequestJSON sends a fasthttp.Request with a JSON body using the given URI, method, and body.
// The resulting fasthttp.Response is NOT released unless an error occurs.
func SendRequestJSON(client *fasthttp.Client, ctx *fasthttp.RequestCtx, method, uri string, body []byte) error {
	// setup the request.
	ctx.Request = *fasthttp.AcquireRequest()
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Request.SetRequestURI(uri)
	ctx.Request.Header.SetMethod(method)
	ctx.Request.Header.SetContentTypeBytes(ContentTypeJSON)
	ctx.Request.SetBodyRaw(body)

	// receive the response from the request.
	ctx.Response = *fasthttp.AcquireResponse()
	err := client.DoTimeout(&ctx.Request, &ctx.Response, timeout)

	// release the request from the pool.
	fasthttp.ReleaseRequest(&ctx.Request)
	if err != nil {
		fasthttp.ReleaseResponse(&ctx.Response)
		return err
	}

	return nil
}
