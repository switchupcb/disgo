package wrapper

import (
	"fmt"
	"strconv"
	"time"

	json "github.com/goccy/go-json"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
)

// Conversion Constants.
const (
	base10              = 10
	bit64               = 64
	msPerSecond float64 = 1000
)

// Content Types
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type
var (
	// ContentTypeURLQueryString is an HTTP Header Content Type that indicates
	// a payload with an encoded URL Query String.
	ContentTypeURLQueryString = []byte("application/x-www-form-urlencoded")

	// ContentTypeJSON is an HTTP Header Content Type that indicates a payload with a JSON body.
	ContentTypeJSON = []byte("application/json")

	// ContentTypeMultipartForm is an HTTP Header Content Type that indicates
	// a payload with multiple content types.
	ContentTypeMultipartForm = []byte("multipart/form-data")

	// ContentTypeJPEG is an HTTP Header Content Type that indicates a payload with a JPEG image.
	ContentTypeJPEG = []byte("image/jpeg")

	// ContentTypePNG is an HTTP Header Content Type that indicates a payload with a PNG image.
	ContentTypePNG = []byte("image/png")

	// ContentTypeWebP is an HTTP Header Content Type that indicates a payload with a WebP image.
	ContentTypeWebP = []byte("image/webp")

	// ContentTypeGIF is an HTTP Header Content Type that indicates a payload with a GIF animated image.
	ContentTypeGIF = []byte("image/gif")
)

// HTTP Header Variables.
const (
	// headerAuthorizationKey represents the key for an "Authorization" HTTP Header.
	headerAuthorizationKey = "Authorization"
)

// HTTP Header Rate Limit Variables.
var (
	// headerDate represents a byte representation of "Date" for HTTP Header functionality.
	headerDate = []byte(FlagRateLimitHeaderDate)

	// headerRetryAfter represents a byte representation of "Retry-After" for HTTP Header functionality.
	headerRateLimitRetryAfter = []byte(FlagRateLimitHeaderRetryAfter)

	// headerRateLimit represents a byte representation of "X-RateLimit-Limit" for HTTP Header functionality.
	headerRateLimit = []byte(FlagRateLimitHeaderLimit)

	// headerRateLimitRemaining represents a byte representation of "X-RateLimit-Remaining" for HTTP Header functionality.
	headerRateLimitRemaining = []byte(FlagRateLimitHeaderRemaining)

	// headerRateLimitReset represents a byte representation of "X-RateLimit-Reset" for HTTP Header functionality.
	headerRateLimitReset = []byte(FlagRateLimitHeaderReset)

	// headerRateLimitResetAfter represents a byte representation of "X-RateLimit-Reset-After" for HTTP Header functionality.
	headerRateLimitResetAfter = []byte(FlagRateLimitHeaderResetAfter)

	// headerRateLimitBucket represents a byte representation of "X-RateLimit-Bucket" for HTTP Header functionality.
	headerRateLimitBucket = []byte(FlagRateLimitHeaderBucket)

	// headerRateLimitGlobal represents a byte representation of "X-RateLimit-Global" for HTTP Header functionality.
	headerRateLimitGlobal = []byte(FlagRateLimitHeaderGlobal)

	// headerRateLimitScope represents a byte representation of "X-RateLimit-Scope" for HTTP Header functionality.
	headerRateLimitScope = []byte(FlagRateLimitHeaderScope)
)

// peekDate peeks an HTTP Header for the Date.
func peekDate(r *fasthttp.Response) (time.Time, error) {
	date, err := time.Parse(time.RFC1123, string(r.Header.PeekBytes(headerDate)))
	if err != nil {
		return time.Time{}, fmt.Errorf("error occurred parsing the \"Date\" HTTP Header: %w", err)
	}

	return date, nil
}

// peekHeaderRetryAfter peeks an HTTP Header for the Rate Limit Header "Retry-After".
func peekHeaderRetryAfter(r *fasthttp.Response) (float64, error) {
	retryafter, err := strconv.ParseFloat(string(r.Header.PeekBytes(headerRateLimitRetryAfter)), bit64)
	if err != nil {
		return 0, fmt.Errorf(errRateLimit, string(headerRateLimitRetryAfter), err)
	}

	return retryafter, nil
}

// peekHeaderRateLimit peeks an HTTP Header for Discord Rate Limit Header values.
func peekHeaderRateLimit(r *fasthttp.Response) RateLimitHeader {
	limit, _ := strconv.Atoi(string(r.Header.PeekBytes(headerRateLimit)))
	remaining, _ := strconv.Atoi(string(r.Header.PeekBytes(headerRateLimitRemaining)))
	reset, _ := strconv.ParseFloat(string(r.Header.PeekBytes(headerRateLimitReset)), bit64)
	resetafter, _ := strconv.ParseFloat(string(r.Header.PeekBytes(headerRateLimitResetAfter)), bit64)
	global, _ := strconv.ParseBool(string(r.Header.PeekBytes(headerRateLimitGlobal)))

	return RateLimitHeader{
		Limit:      limit,
		Remaining:  remaining,
		Reset:      reset,
		ResetAfter: resetafter,
		Bucket:     string(r.Header.PeekBytes(headerRateLimitBucket)),
		Global:     global,
		Scope:      string(r.Header.PeekBytes(headerRateLimitScope)),
	}
}

// SendRequest sends a fasthttp.Request using the given route ID, HTTP method, URI, content type and body,
// then parses the response into dst.
func SendRequest(bot *Client, xid, routeid, resourceid, method, uri string, content, body []byte, dst any) error { //nolint:gocyclo,maintidx
	retries := 0
	requestid := routeid + resourceid
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	request.Header.SetMethod(method)
	request.Header.SetContentTypeBytes(content)
	request.Header.Set(headerAuthorizationKey, bot.Authentication.Header)
	request.SetRequestURI(uri)
	request.SetBodyRaw(body)
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)

	// Certain endpoints are not bound to the bot's Global Rate Limit.
	if IgnoreGlobalRateLimitRouteIDs[requestid] {
		goto SEND
	}

RATELIMIT:
	// a single request or response is PROCESSED at any point in time.
	bot.Config.Request.RateLimiter.Lock()

	LogRequest(Logger.Trace(), bot.ApplicationID, xid, routeid, resourceid, uri).Msg("processing request")
	if Logger.GetLevel() == zerolog.TraceLevel {
		LogRequestBody(Logger.Trace(), bot.ApplicationID, xid, routeid, resourceid, uri, string(body)).Msg("")
	}

	// check Global and Route Rate Limit Buckets prior to sending the current request.
	for {
		bot.Config.Request.RateLimiter.StartTx()

		globalBucket := bot.Config.Request.RateLimiter.GetBucket(GlobalRateLimitRouteID, "")

		// stop waiting when the Global Rate Limit Bucket is NOT empty.
		if isNotEmpty(globalBucket) {
			routeBucket := bot.Config.Request.RateLimiter.GetBucket(routeid, resourceid)

			if isNotEmpty(routeBucket) {
				break
			}

			if isExpired(routeBucket) {
				// When a Route Bucket expires, its new expiry becomes unknown.
				// As a result, it will never reset (again) until a pending request's
				// response sets a new expiry.
				routeBucket.Reset(time.Time{})
			}

			var wait time.Time
			if routeBucket != nil {
				wait = routeBucket.Expiry
			}

			// do NOT block other requests due to a Route Rate Limit.
			bot.Config.Request.RateLimiter.EndTx()
			bot.Config.Request.RateLimiter.Unlock()

			// reduce CPU usage by blocking the current goroutine
			// until it's eligible for action.
			if routeBucket != nil {
				<-time.After(time.Until(wait))
			}

			goto RATELIMIT
		}

		// reset the Global Rate Limit Bucket when the current Bucket has passed its expiry.
		if isExpired(globalBucket) {
			globalBucket.Reset(time.Now().Add(time.Second))
		}

		bot.Config.Request.RateLimiter.EndTx()
	}

	if globalBucket := bot.Config.Request.RateLimiter.GetBucket(GlobalRateLimitRouteID, ""); globalBucket != nil {
		globalBucket.Use(1)
	}

	if routeBucket := bot.Config.Request.RateLimiter.GetBucket(routeid, resourceid); routeBucket != nil {
		routeBucket.Use(1)
	}

	bot.Config.Request.RateLimiter.EndTx()
	bot.Config.Request.RateLimiter.Unlock()

SEND:
	LogRequest(Logger.Trace(), bot.ApplicationID, xid, routeid, resourceid, uri).Msg("sending request")

	// send the request.
	if err := bot.Config.Request.Client.DoTimeout(request, response, bot.Config.Request.Timeout); err != nil {
		return fmt.Errorf("%w", err)
	}

	LogResponse(LogRequest(Logger.Info(), bot.ApplicationID, xid, routeid, resourceid, uri),
		response.Header.String(), string(response.Body()),
	).Msg("")

	var header RateLimitHeader

	// confirm the response with the rate limiter.
	//
	// Certain endpoints are not bound to the bot's Global Rate Limit.
	if !IgnoreGlobalRateLimitRouteIDs[requestid] {
		// parse the Rate Limit Header for per-route rate limit functionality.
		header = peekHeaderRateLimit(response)

		// parse the Date header for Global Rate Limit functionality.
		date, err := peekDate(response)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		bot.Config.Request.RateLimiter.StartTx()

		// confirm the Global Rate Limit Bucket.
		globalBucket := bot.Config.Request.RateLimiter.GetBucket(GlobalRateLimitRouteID, "")
		if globalBucket != nil {
			globalBucket.ConfirmDate(1, date)
		}

		// confirm the Route Rate Limit Bucket (if applicable).
		routeBucket := bot.Config.Request.RateLimiter.GetBucket(routeid, resourceid)
		switch {
		// when there is no Discord Bucket, remove the route's mapping to a rate limit Bucket.
		case header.Bucket == "":
			if routeBucket != nil {
				bot.Config.Request.RateLimiter.SetBucketID(requestid, nilRouteBucket)
				routeBucket = nil
			}

		// when the route's Bucket ID does NOT match the Discord Bucket, update it.
		case routeBucket == nil && header.Bucket != "" || routeBucket.ID != header.Bucket:
			var pending int16
			if routeBucket != nil {
				pending = routeBucket.Pending
			}

			// update the route ID mapping to a rate limit Bucket ID.
			bot.Config.Request.RateLimiter.SetBucketID(requestid, header.Bucket)

			// map the Bucket ID to the updated Rate Limit Bucket.
			if bucket := bot.Config.Request.RateLimiter.GetBucketFromID(header.Bucket); bucket != nil {
				routeBucket = bucket
			} else {
				routeBucket = getBucket()
				bot.Config.Request.RateLimiter.SetBucketFromID(header.Bucket, routeBucket)
			}

			routeBucket.Pending += pending
			routeBucket.ID = header.Bucket
		}

		if routeBucket != nil {
			routeBucket.ConfirmHeader(1, header)
		}

		if response.StatusCode() != fasthttp.StatusTooManyRequests {
			bot.Config.Request.RateLimiter.EndTx()
		}
	}

	// handle the response.
	switch response.StatusCode() {
	case fasthttp.StatusOK, fasthttp.StatusCreated:
		// parse the response data.
		if err := json.Unmarshal(response.Body(), dst); err != nil {
			return fmt.Errorf(errUnmarshal, dst, err)
		}

		return nil

	case fasthttp.StatusNoContent:
		return nil

	// process the rate limit.
	case fasthttp.StatusTooManyRequests:
		retry := retries < bot.Config.Request.Retries
		retries++

		switch header.Scope { //nolint:gocritic
		// Discord per-resource (shared) rate limit headers include the per-route (user) bucket.
		//
		// when a per-resource rate limit is encountered, send another request or return.
		case RateLimitScopeValueShared:
			bot.Config.Request.RateLimiter.EndTx()

			if retry || bot.Config.Request.RetryShared {
				goto RATELIMIT
			}

			return StatusCodeError(response.StatusCode())
		}

		// parse the rate limit response data for `retry_after`.
		var data RateLimitResponse
		if err := json.Unmarshal(response.Body(), &data); err != nil {
			return fmt.Errorf("%w", err)
		}

		// determine the reset time.
		var reset time.Time
		if data.RetryAfter == 0 {
			// when the 429 is from Discord, use the `retry_after` value (s).
			reset = time.Now().Add(time.Millisecond * time.Duration(data.RetryAfter*msPerSecond))
		} else {
			// when the 429 is from a Cloudflare Ban, use the `"Retry-After"` value (s).
			retryafter, err := peekHeaderRetryAfter(response)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			reset = time.Now().Add(time.Millisecond * time.Duration(retryafter*msPerSecond))
		}

		if data.Code == nil {
			LogRequest(Logger.Debug(), bot.ApplicationID, xid, routeid, resourceid, uri).
				Time(LogCtxReset, reset).Msg("")
		} else {
			LogRequest(Logger.Debug(), bot.ApplicationID, xid, routeid, resourceid, uri).
				Time(LogCtxReset, reset).
				Err(JSONCodeError(*data.Code)).Msg("")
		}

		switch header.Global {
		// when the global request rate limit is encountered.
		case true:
			// when the current time is BEFORE the reset time,
			// all requests must wait until the 429 expires.
			if time.Now().Before(reset) {
				if globalBucket := bot.Config.Request.RateLimiter.GetBucket(GlobalRateLimitRouteID, ""); globalBucket != nil {
					globalBucket.Remaining = 0
					globalBucket.Expiry = reset.Add(time.Millisecond)
				}
			}

			bot.Config.Request.RateLimiter.EndTx()

		// when a per-route (user) rate limit is encountered.
		case false:
			// do NOT block other requests while waiting for a Route Rate Limit.
			bot.Config.Request.RateLimiter.EndTx()

			// when the current time is BEFORE the reset time,
			// requests with the same Rate Limit Bucket must wait until the 429 expires.
			if time.Now().Before(reset) {
				if routeBucket := bot.Config.Request.RateLimiter.GetBucket(routeid, resourceid); routeBucket != nil {
					routeBucket.Remaining = 0
					routeBucket.Expiry = reset.Add(time.Millisecond)
				}
			}
		}

		if retry {
			goto RATELIMIT
		}

		return StatusCodeError(fasthttp.StatusTooManyRequests)

	// retry the request on a bad gateway server error.
	case fasthttp.StatusBadGateway:
		if retries < bot.Config.Request.Retries {
			retries++

			goto RATELIMIT
		}

		return StatusCodeError(fasthttp.StatusBadGateway)

	default:
		return StatusCodeError(response.StatusCode())
	}
}

// isExpired determines whether a rate limit Bucket is expired.
func isExpired(b *Bucket) bool {
	// a rate limit bucket is expired when
	// 1. the bucket exists AND
	// 2. the current time occurs after the non-zero expiry time.
	return b != nil && !b.Expiry.IsZero() && time.Now().After(b.Expiry)
}

// isNotEmpty determines whether a rate limit Bucket is NOT empty.
func isNotEmpty(b *Bucket) bool {
	// a rate limit bucket is NOT empty when
	// 1. the bucket does not exist OR
	// 2. there is one or more remaining request token(s).
	return b == nil || b.Remaining > 0
}
