package wrapper

import (
	"fmt"
	"strconv"
	"strings"
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
	// headerDate       = []byte("Date")
	headerDateString = "Date:"

	// msPerSecond represents the amount of milliseconds in a second.
	msPerSecond float64 = 1000
)

var (
	// https://discord.com/developers/docs/topics/rate-limits#global-rate-limit
	GlobalRateLimit = &Bucket{ //nolint:exhaustruct
		Limit:     FlagGlobalRequestRateLimit,
		Remaining: FlagGlobalRequestRateLimit,
	}
)

// peekHeaderRateLimit peeks an HTTP Header for Rate Limit Header values.
func peekHeaderRateLimit(r *fasthttp.Response) *RateLimitHeader {
	limit, _ := strconv.Atoi(string(r.Header.PeekBytes(headerRateLimit)))
	remaining, _ := strconv.Atoi(string(r.Header.PeekBytes(headerRateLimitRemaining)))
	reset, _ := strconv.ParseInt(string(r.Header.PeekBytes(headerRateLimitReset)), base10, bit64)
	resetafter, _ := strconv.ParseFloat(string(r.Header.PeekBytes(headerRateLimitResetAfter)), bit64)
	global, _ := strconv.ParseBool(string(r.Header.PeekBytes(headerRateLimitGlobal)))

	return &RateLimitHeader{
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

// peekDate peeks an HTTP Header for the Date.
func peekDate(r *fasthttp.Response) (time.Time, error) {
	// https://github.com/valyala/fasthttp/issues/1339
	var d string
	for _, line := range strings.Split(r.Header.String(), "\n") {
		if strings.Contains(line, headerDateString) {
			d = line[6 : len(line)-1]
			break
		}
	}

	date, err := time.Parse(time.RFC1123, d)
	if err != nil {
		return time.Time{}, fmt.Errorf("error occurred parsing the \"Date\" HTTP Header: %w", err)
	}

	return date, nil
}

// SendRequest sends a fasthttp.Request using the given URI, HTTP method, content type and body,
// then parses the response into dst.
func SendRequest(bot *Client, method, uri string, content, body []byte, dst any) error { //nolint:gocognit,cyclop,funlen
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

RATELIMIT:
	// a single request is PROCESSED from a queue at any point in time.
	bot.Config.RateLimiter.Lock()

	// check global and route rate limit Buckets prior to sending the current request.
	for {
		bot.Config.RateLimiter.StartTx()

		// when priority is greater than 0, another request will be sent first.
		if atomic.LoadInt32(&bot.Config.GlobalRateLimit.Priority) > 0 {
			bot.Config.RateLimiter.EndTx()
			bot.Config.RateLimiter.Unlock()

			goto RATELIMIT
		}

		// when no requests remain in the global rate limit bucket,
		// wait until the bucket resets to send a request.
		if isEmpty(bot.Config.GlobalRateLimit) {
			bot.Config.RateLimiter.EndTx()

			break
		}

		// reset the global rate limit bucket when the current Bucket has passed its expiry.
		if isExpired(bot.Config.GlobalRateLimit) {
			bot.Config.GlobalRateLimit.Reset(time.Now().Add(time.Second))
		}

		bot.Config.RateLimiter.EndTx()
	}

	goto SEND

PRIORITY:
	// priorityWait contains the same functionality as the above for loop,
	// but without a priority counter check.
	priorityWait(bot.Config.RateLimiter, bot.Config.GlobalRateLimit)

SEND:
	bot.Config.RateLimiter.StartTx()
	bot.Config.GlobalRateLimit.Use(1)
	bot.Config.RateLimiter.EndTx()
	bot.Config.RateLimiter.Unlock()

	// send the request.
	if err := bot.Config.Client.DoTimeout(request, response, bot.Config.Timeout); err != nil {
		return fmt.Errorf("%w", err)
	}

	date, err := peekDate(response)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	bot.Config.RateLimiter.StartTx()
	bot.Config.GlobalRateLimit.Confirm(1, date)
	bot.Config.RateLimiter.EndTx()

	// follow redirects.
	if fasthttp.StatusCodeIsRedirect(response.StatusCode()) {
		location := response.Header.PeekBytes(headerLocation)
		if len(location) == 0 {
			return fmt.Errorf(ErrRedirect, uri)
		}

		request.URI().UpdateBytes(location)

		goto SEND
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
		// prevent other requests in the queue from being sent by incrementing priority.
		atomic.AddInt32(&bot.Config.GlobalRateLimit.Priority, 1)

		// parse the rate limit header.
		_, err := peekHeader429(response)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		header := peekHeaderRateLimit(response)

		// parse the rate limit response data.
		var data RateLimitResponse
		if err := json.Unmarshal(response.Body(), &data); err != nil {
			return fmt.Errorf("%w", err)
		}

		reset := time.Now().Add(time.Millisecond * time.Duration(data.RetryAfter*msPerSecond))
		retry := retries < bot.Config.Retries
		retries++

		// a single request is PROCESSED from a mutex queue at any point in time.
		bot.Config.RateLimiter.Lock()
		if header.Global {
			// when the current time is BEFORE the reset time,
			// all requests must wait until the 429 expires.
			if time.Now().Before(reset) {
				bot.Config.RateLimiter.StartTx()
				bot.Config.GlobalRateLimit.Remaining = 0
				bot.Config.GlobalRateLimit.Expiry = reset.Add(time.Second + 1)
				bot.Config.RateLimiter.EndTx()
			}
		}

		atomic.AddInt32(&bot.Config.GlobalRateLimit.Priority, -1)

		if retry {
			goto PRIORITY
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
func priorityWait(r RateLimiter, b *Bucket) {
	for {
		r.StartTx()
		if isEmpty(b) {
			r.EndTx()

			break
		}

		if isExpired(b) {
			b.Reset(time.Now().Add(time.Second))
		}

		r.EndTx()
	}
}

// isExpired determines whether a rate limit Bucket is expired.
func isExpired(b *Bucket) bool {
	return !b.Expiry.IsZero() && time.Now().After(b.Expiry)
}

// isEmpty determines whether a rate limit Bucket is empty.
func isEmpty(b *Bucket) bool {
	return b.Remaining > 0
}
