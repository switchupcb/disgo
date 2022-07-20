package wrapper

import (
	"fmt"
	"strconv"
	"sync/atomic"
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
	headerRateLimitRetryAfter = []byte("Retry-After")

	// headerDate represents a byte representation of "Date" for HTTP Header functionality.
	headerDate = []byte("Date")

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
func SendRequest(bot *Client, routeid uint16, method, uri string, content, body []byte, dst any) error { //nolint:gocognit,cyclop,funlen,gocyclo
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
	// a single request is PROCESSED from a queue at any point in time.
	bot.Config.RateLimiter.Lock()

	// check global and route rate limit Buckets prior to sending the current request.
	for {
		bot.Config.RateLimiter.StartTx()
		globalBucket := bot.Config.RateLimiter.GetBucket(0)

		// when a Bucket contains a priority request, it will be sent before the current request.
		if containsPriority(globalBucket) {
			bot.Config.RateLimiter.EndTx()
			bot.Config.RateLimiter.Unlock()

			goto RATELIMIT
		}

		// stop waiting when the Global Rate Limit Bucket is NOT empty.
		if isNotEmpty(globalBucket) {
			routeBucket := bot.Config.RateLimiter.GetBucket(routeid)

			// TODO: A request should not be sent if a Global Rate Limit Bucket hit a 429.
			if containsPriority(routeBucket) {
				bot.Config.RateLimiter.EndTx()
				bot.Config.RateLimiter.Unlock()

				goto RATELIMIT
			}

			if isNotEmpty(routeBucket) {
				bot.Config.RateLimiter.EndTx()

				break
			}

			if isExpired(routeBucket) {
				routeBucket.Reset(time.Time{})
			}

			// TODO: Multiple requests would be blocked by a single route rate limit bucket.
			bot.Config.RateLimiter.EndTx()

			continue
		}

		// reset the global rate limit bucket when the current Bucket has passed its expiry.
		if isExpired(globalBucket) {
			globalBucket.Reset(time.Now().Add(time.Second))
		}

		bot.Config.RateLimiter.EndTx()
	}

USE:
	bot.Config.RateLimiter.StartTx()

	// func() is required to allow a jump over a variable declaration (from goto SEND).
	func() {
		globalBucket := bot.Config.RateLimiter.GetBucket(0)
		if globalBucket != nil {
			globalBucket.Use(1)
		}

		if routeBucket := bot.Config.RateLimiter.GetBucket(routeid); routeBucket != nil {
			routeBucket.Use(1)
		}
	}()

	bot.Config.RateLimiter.EndTx()
	bot.Config.RateLimiter.Unlock()

SEND:
	// send the request.
	if err := bot.Config.Client.DoTimeout(request, response, bot.Config.Timeout); err != nil {
		return fmt.Errorf("%w", err)
	}

	// parse the Rate Limit Header for per-route rate limit functionality.
	header := peekHeaderRateLimit(response)
	fmt.Println("\nTime Now", time.Now(), "\n"+response.Header.String())

	// confirm the response with the rate limiter.
	//
	// Interaction endpoints are not bound to the bot's Global Rate Limit.
	if !interactionRouteIDs[routeid] { // nolint:nestif
		// parse the Date header for global rate limit functionality.
		date, err := peekDate(response)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		bot.Config.RateLimiter.StartTx()

		// confirm the global rate limit Bucket.
		globalBucket := bot.Config.RateLimiter.GetBucket(0)
		if globalBucket != nil {
			globalBucket.ConfirmDate(1, date)
		}

		// confirm the route rate limit Bucket (if applicable).
		routeBucket := bot.Config.RateLimiter.GetBucket(routeid)
		if header.Bucket == "" {
			// when there is no Discord Bucket, remove the route's mapping to a rate limit Bucket.
			if routeBucket != nil {
				bot.Config.RateLimiter.SetBucketHash(routeid, nilRouteBucket)
			}
		} else {
			// ensure the route Bucket is up to date.
			if routeBucket == nil || routeBucket.ID != header.Bucket {
				// update the route ID mapping to a rate limit Bucket ID.
				//
				// TODO: Potential memory leak due to excessive storage of buckets
				// if IDs change between resets.
				bot.Config.RateLimiter.SetBucketHash(routeid, header.Bucket)

				// update the Bucket ID mapping to a rate limit Bucket.
				if bucket := bot.Config.RateLimiter.GetBucketFromHash(header.Bucket); bucket != nil {
					routeBucket = bucket
				} else {
					bot.Config.RateLimiter.SetBucketFromHash(header.Bucket, new(Bucket))
				}
			}

			routeBucket.ConfirmHeader(1, routeid, header)
		}

		bot.Config.RateLimiter.EndTx()
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
		// TODO: Handle priority.
		fmt.Println(response.Header.String(), response.Body())

		// prevent other requests in the queue from being sent by incrementing priority.
		if header.Global {
			if globalBucket := bot.Config.RateLimiter.GetBucket(0); globalBucket != nil {
				atomic.AddInt32(&globalBucket.Priority, 1)
			}
		} else {
			// TODO
		}

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

		// a single request is PROCESSED from a mutex queue at any point in time.
		retry := retries < bot.Config.Retries
		retries++

		bot.Config.RateLimiter.Lock()

		if header.Global {
			// when the current time is BEFORE the reset time,
			// all requests must wait until the 429 expires.
			if time.Now().Before(reset) {
				bot.Config.RateLimiter.StartTx()
				if globalBucket := bot.Config.RateLimiter.GetBucket(0); globalBucket != nil {
					globalBucket.Remaining = 0
					globalBucket.Expiry = reset.Add(time.Second + 1)
				}
				bot.Config.RateLimiter.EndTx()
			}
		} else {
			// TODO
		}

		if globalBucket := bot.Config.RateLimiter.GetBucket(0); globalBucket != nil {
			atomic.AddInt32(&globalBucket.Priority, -1)
		}

		if retry {
			// priorityWait contains the same functionality as the RATELIMIT `for` loop,
			// but without a priority counter check.
			priorityWait(bot.Config.RateLimiter, routeid)

			goto USE
		}

		bot.Config.RateLimiter.Unlock()

		return StatusCodeError(fasthttp.StatusTooManyRequests)

	// retry the request on a bad gateway server error.
	case fasthttp.StatusBadGateway:
		if retries < bot.Config.Retries {
			retries++

			goto RATELIMIT
		}

		return StatusCodeError(fasthttp.StatusBadGateway)

	default:
		return StatusCodeError(response.StatusCode())
	}
}

// priorityWait waits until a priority request is ready to be sent.
func priorityWait(r RateLimiter, routeid uint16) {
	for {
		r.StartTx()
		globalBucket := r.GetBucket(0)

		if isNotEmpty(globalBucket) {
			routeBucket := r.GetBucket(routeid)

			if isNotEmpty(routeBucket) {
				r.EndTx()

				break
			}

			if isExpired(routeBucket) {
				routeBucket.Reset(time.Time{})
			}

			r.EndTx()

			continue
		}

		if isExpired(globalBucket) {
			globalBucket.Reset(time.Now().Add(time.Second))
		}

		r.EndTx()
	}
}

// containsPriority determines whether a rate limit Bucket contains a priority request.
func containsPriority(b *Bucket) bool {
	// a rate limit bucket contains priority when
	// 1. the bucket exists AND
	// 2. there is one or more priority request token(s).
	return b != nil && atomic.LoadInt32(&b.Priority) > 0
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
