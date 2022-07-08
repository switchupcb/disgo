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
	headerDate       = []byte("Date")
	headerDateString = "Date:"
)

var (
	// https://discord.com/developers/docs/topics/rate-limits#global-rate-limit
	GlobalRateLimit = &Bucket{
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

// peekHeader429 peeks an HTTP Header with a 429 Status Code for the Rate Limit Header "RetryAfter".
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
		return time.Time{}, fmt.Errorf("error occurred parsing the \"Date\" HTTP Header: %v", err)
	}

	return date, nil
}

// SendRequest sends a fasthttp.Request using the given URI, HTTP method, content type and body,
// then parses the response into dst.
func SendRequest(bot *Client, method, uri string, content, body []byte, dst any) error {
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

	fmt.Println("queued request")

RATELIMIT:
	// check the global rate limit prior to sending the request.
	bot.Config.GlobalRateLimit.mu.Lock()
	fmt.Println("unlocked a request in the queue.")

	// reset the global rate limit bucket when the current time has passed its expiry.
	bot.Config.GlobalRateLimit.muExpiry.RLock()
	if time.Now().After(bot.Config.GlobalRateLimit.Expiry) || time.Now().Equal(bot.Config.GlobalRateLimit.Expiry) {
		fmt.Println("expiry passed")
		bot.Config.GlobalRateLimit.muExpiry.RUnlock()
		bot.Config.GlobalRateLimit.Reset(time.Now().Truncate(time.Second).Add(time.Second))
	} else {
		bot.Config.GlobalRateLimit.muExpiry.RUnlock()
	}

	// when there are more or equal pending requests compared to the limit.
	if atomic.LoadInt32(&bot.Config.GlobalRateLimit.Pending) >= GlobalRateLimit.Limit {
		fmt.Println("requeue request from pending")
		bot.Config.GlobalRateLimit.mu.Unlock()

		goto RATELIMIT
	}

	/*
		// when priority is greater than 0, another request will be sent first.
		if atomic.LoadInt32(&bot.Config.GlobalRateLimit.Priority) > 0 {
			fmt.Println("requeue request from priority")

			// unlock and re-lock this request's mutex to give the pending request priority
			// by sending this request (and subsequent requests) after that pending request.
			bot.Config.GlobalRateLimit.mu.Unlock()

			goto RATELIMIT
		}
	*/

	// when no requests remain in the global rate limit bucket,
	// wait until the bucket resets to send a request.
	if atomic.LoadInt32(&bot.Config.GlobalRateLimit.Remaining) <= 0 {
		fmt.Println("no requests remain")

		// wait until the current bucket expires, then reset it.
		bot.Config.GlobalRateLimit.muExpiry.RLock()
		expiry := bot.Config.GlobalRateLimit.Expiry
		bot.Config.GlobalRateLimit.muExpiry.RUnlock()

		fmt.Println("waiting", expiry, time.Now())
		<-time.After(time.Until(expiry))
		fmt.Println("finished waiting.", time.Now())
		bot.Config.GlobalRateLimit.mu.Unlock()

		goto RATELIMIT
	}

	bot.Config.GlobalRateLimit.Use(1)
	bot.Config.GlobalRateLimit.mu.Unlock()

	goto SEND

PRIORITY:
	// ensure that requests (prompted by 429s) are sent first,
	// and followed by subsequent requests.
	bot.Config.GlobalRateLimit.mu.Lock()
	bot.Config.GlobalRateLimit.Use(1)
	bot.Config.GlobalRateLimit.mu.Unlock()

SEND:
	// send the request.
	sent := time.Now()
	if err := bot.Config.Client.DoTimeout(request, response, bot.Config.Timeout); err != nil {
		return fmt.Errorf("%w", err)
	}

	date, err := peekDate(response)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	bot.Config.GlobalRateLimit.mu.Lock()
	bot.Config.GlobalRateLimit.Confirm(1, date, sent)
	bot.Config.GlobalRateLimit.mu.Unlock()

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
		retryafter, err := peekHeader429(response)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		header := peekHeaderRateLimit(response)

		var data RateLimitResponse
		if err := json.Unmarshal(response.Body(), &data); err != nil {
			return fmt.Errorf("%w", err)
		}

		fmt.Println("\t429", header, "Retry-After", retryafter, "retry_after", data.RetryAfter)
		fmt.Println("\n", response.Header.String(), string(response.Body()))

		// add json unmarshal for precision
		retry := retries < bot.Config.Retries
		retries++

		// expire the current rate limit immediately.
		timeAt429 := time.Now()
		if header.Global {
			bot.Config.GlobalRateLimit.mu.Lock()
			atomic.StoreInt32(&bot.Config.GlobalRateLimit.Remaining, 0)

			// set the global rate limit expiry such that the next request
			// will be sent after "Retry-After" seconds.
			bot.Config.GlobalRateLimit.muExpiry.Lock()
			bot.Config.GlobalRateLimit.Expiry = timeAt429.Add(time.Second * time.Duration(retryafter))
			bot.Config.GlobalRateLimit.muExpiry.Unlock()

			if retry {
				// set priority to indicate that this request has priority when the bucket resets.
				// bot.Config.GlobalRateLimit.Priority++
				bot.Config.GlobalRateLimit.mu.Unlock()

				goto PRIORITY
			}

			defer bot.Config.GlobalRateLimit.mu.Unlock()
		}

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
