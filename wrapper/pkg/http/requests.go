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
	defer fasthttp.ReleaseResponse(&ctx.Response)

	// release the request from the pool if something goes wrong.
	if err := client.DoTimeout(&ctx.Request, &ctx.Response, timeout); err != nil {
		fasthttp.ReleaseRequest(&ctx.Request)
		return err
	}

	return nil
}
