package wrapper

import (
	"fmt"
	"log"
	"strconv"
	"time"

	json "github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
)

// Conversion Constants.
const (
	base10 = 10
	bit64  = 64
)

// HTTP Header Variables.
var (
	// contentTypeURL represents an HTTP Header indicating a payload with an encoded URL Query String.
	contentTypeURL = []byte("application/x-www-form-urlencoded")

	// contentTypeJSON represents an HTTP Header that indicates a payload with a JSON body.
	contentTypeJSON = []byte("application/json")

	// contentTypeMulti represents an HTTP Header that indicates a payload with a multi-part (file) body.
	contentTypeMulti = []byte("multipart/form-data")

	// headerLocation represents a byte representation of "Location" for HTTP Header functionality.
	headerLocation = []byte("Location")

	// headerAuthorizationKey represents the key for an "Authorization" HTTP Header.
	headerAuthorizationKey = "Authorization"
)

// HTTP Header Rate Limit Variables.
var (
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

	// headerRetryAfter represents a byte representation of "Retry-After" for HTTP Header functionality.
	headerRateLimitRetryAfter = []byte(FlagRateLimitHeaderRetryAfter)

	// headerDate represents a byte representation of "Date" for HTTP Header functionality.
	headerDate = []byte(FlagRateLimitHeaderDate)

	// msPerSecond represents the amount of milliseconds in a second.
	msPerSecond float64 = 1000
)

// Custom Rate Limit Variables
var (
	// Interaction endpoints are not bound to the bot's Global Rate Limit.
	// https://discord.com/developers/docs/interactions/receiving-and-responding#endpoints
	interactionRouteIDs = map[uint16]bool{
		18: true, 19: true, 20: true, 21: true, 22: true, 23: true, 24: true, 25: true,
	}
)

// peekDate peeks an HTTP Header for the Date.
func peekDate(r *fasthttp.Response) (time.Time, error) {
	date, err := time.Parse(time.RFC1123, string(r.Header.PeekBytes(headerDate)))
	if err != nil {
		return time.Time{}, fmt.Errorf("error occurred parsing the \"Date\" HTTP Header: %w", err)
	}

	return date, nil
}

// peekHeaderRateLimit peeks an HTTP Header for Rate Limit Header values.
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

// peekHeader429 peeks an HTTP Header with a 429 Status Code for the Rate Limit Header "Retry-After".
func peekHeader429(r *fasthttp.Response) (int64, error) {
	retryafter, err := strconv.ParseInt(string(r.Header.PeekBytes(headerRateLimitRetryAfter)), base10, bit64)
	if err != nil {
		return 0, fmt.Errorf(ErrRateLimit, string(headerRateLimitRetryAfter), err)
	}

	return retryafter, nil
}

// SendRequest sends a fasthttp.Request using the given route ID, HTTP method, URI, content type and body,
// then parses the response into dst.
func SendRequest(bot *Client, routeid uint16, method, uri string, content, body []byte, dst any) error { //nolint:gocyclo,maintidx
	retries := 0
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	request.Header.SetMethod(method)
	request.Header.SetContentTypeBytes(content)
	request.Header.Set(headerAuthorizationKey, bot.Authentication.Header)
	request.SetRequestURI(uri)
	request.SetBodyRaw(body)
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)

	// Interaction endpoints are not bound to the bot's Global Rate Limit.
	if interactionRouteIDs[routeid] {
		goto SEND
	}

RATELIMIT:
	// a single request or response is PROCESSED at any point in time.
	bot.Config.Request.RateLimiter.Lock()

	// check Global and Route Rate Limit Buckets prior to sending the current request.
	for {
		bot.Config.Request.RateLimiter.StartTx()
		globalBucket := bot.Config.Request.RateLimiter.GetBucket(0)

		// stop waiting when the Global Rate Limit Bucket is NOT empty.
		if isNotEmpty(globalBucket) {
			routeBucket := bot.Config.Request.RateLimiter.GetBucket(routeid)

			if isNotEmpty(routeBucket) {
				bot.Config.Request.RateLimiter.EndTx()

				goto USE
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

USE:
	bot.Config.Request.RateLimiter.StartTx()

	// func() is required to allow a jump over a variable declaration (from goto SEND).
	func() {
		globalBucket := bot.Config.Request.RateLimiter.GetBucket(0)
		if globalBucket != nil {
			globalBucket.Use(1)
		}

		if routeBucket := bot.Config.Request.RateLimiter.GetBucket(routeid); routeBucket != nil {
			routeBucket.Use(1)
		}
	}()

	bot.Config.Request.RateLimiter.EndTx()
	bot.Config.Request.RateLimiter.Unlock()

SEND:
	// send the request.
	if err := bot.Config.Request.Client.DoTimeout(request, response, bot.Config.Request.Timeout); err != nil {
		return fmt.Errorf("%w", err)
	}

	var header RateLimitHeader

	// confirm the response with the rate limiter.
	//
	// Interaction endpoints are not bound to the bot's Global Rate Limit.
	if !interactionRouteIDs[routeid] { //nolint:nestif
		log.Println("\n" + time.Now().String() + "\n" + response.Header.String())

		// parse the Rate Limit Header for per-route rate limit functionality.
		header = peekHeaderRateLimit(response)

		// parse the Date header for global rate limit functionality.
		date, err := peekDate(response)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		bot.Config.Request.RateLimiter.StartTx()

		// confirm the global rate limit Bucket.
		globalBucket := bot.Config.Request.RateLimiter.GetBucket(0)
		if globalBucket != nil {
			globalBucket.ConfirmDate(1, date)
		}

		// confirm the route rate limit Bucket (if applicable).
		routeBucket := bot.Config.Request.RateLimiter.GetBucket(routeid)
		switch {
		// when there is no Discord Bucket, remove the route's mapping to a rate limit Bucket.
		case header.Bucket == "":
			if routeBucket != nil {
				bot.Config.Request.RateLimiter.SetBucketHash(routeid, nilRouteBucket)
			}

		// when the route's Bucket ID does NOT match the Discord Bucket, update it.
		case routeBucket.ID != header.Bucket:
			var pending int16
			if routeBucket != nil {
				pending = routeBucket.Pending
			}

			// update the route ID mapping to a rate limit Bucket ID.
			bot.Config.Request.RateLimiter.SetBucketHash(routeid, header.Bucket)

			// map the Bucket ID to the updated Rate Limit Bucket.
			if bucket := bot.Config.Request.RateLimiter.GetBucketFromHash(header.Bucket); bucket != nil {
				routeBucket = bucket
			} else {
				routeBucket = getBucket()
				bot.Config.Request.RateLimiter.SetBucketFromHash(header.Bucket, routeBucket)
			}

			routeBucket.Pending += pending
			routeBucket.ID = header.Bucket
		}

		if routeBucket != nil {
			routeBucket.ConfirmHeader(1, routeid, header)
		}

		if response.StatusCode() != fasthttp.StatusTooManyRequests {
			bot.Config.Request.RateLimiter.EndTx()
		}
	}

	// follow redirects.
	if fasthttp.StatusCodeIsRedirect(response.StatusCode()) {
		location := response.Header.PeekBytes(headerLocation)
		if len(location) == 0 {
			return fmt.Errorf(ErrRedirect, uri)
		}

		request.URI().UpdateBytes(location)

		goto USE
	}

	// handle the response.
	switch response.StatusCode() {
	case fasthttp.StatusOK:
		// parse the response data.
		if err := json.Unmarshal(response.Body(), dst); err != nil {
			return fmt.Errorf("%w", err)
		}

		return nil

	// process the rate limit.
	case fasthttp.StatusTooManyRequests:
		log.Println(string(response.Body()))

		// parse the rate limit response data for `retry_after`.
		var data RateLimitResponse
		if err := json.Unmarshal(response.Body(), &data); err != nil {
			return fmt.Errorf("%w", err)
		}

		// determine the reset time.
		var reset time.Time
		if data.RetryAfter == 0 {
			// when the 429 is from Discord, use the `retry_after` value (ms).
			reset = time.Now().Add(time.Millisecond * time.Duration(data.RetryAfter*msPerSecond))
		} else {
			// when the 429 is a Cloudflare Ban, use the `"Retry-After"` value (s).
			retryafter, err := peekHeader429(response)
			if err != nil {
				return fmt.Errorf("%w", err)
			}
			reset = time.Now().Add(time.Second * time.Duration(retryafter))
		}

		log.Printf("429 Reset Time: %v", reset)

		retry := retries < bot.Config.Request.Retries
		retries++

		switch header.Global {
		case true:
			// when the current time is BEFORE the reset time,
			// all requests must wait until the 429 expires.
			if time.Now().Before(reset) {
				if globalBucket := bot.Config.Request.RateLimiter.GetBucket(0); globalBucket != nil {
					globalBucket.Remaining = 0
					globalBucket.Expiry = reset.Add(time.Millisecond)
				}
			}

			bot.Config.Request.RateLimiter.EndTx()

		case false:
			// do NOT block other requests while waiting for a Route Rate Limit.
			bot.Config.Request.RateLimiter.EndTx()

			// when the current time is BEFORE the reset time,
			// requests with the same Rate Limit Bucket must wait until the 429 expires.
			if time.Now().Before(reset) {
				if routeBucket := bot.Config.Request.RateLimiter.GetBucket(routeid); routeBucket != nil {
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
