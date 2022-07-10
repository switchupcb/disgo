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
	bot.Config.GlobalRateLimit.muQueue.Lock()
	fmt.Println("processing a request in the queue.")

	for {
		bot.Config.GlobalRateLimit.muAtomic.Lock()
		// when no requests remain in the global rate limit bucket,
		// wait until the bucket resets to send a request.
		if bot.Config.GlobalRateLimit.Remaining > 0 {
			bot.Config.GlobalRateLimit.muAtomic.Unlock()

			break
		}

		// reset the global rate limit bucket when the current Bucket has passed its expiry.
		if !bot.Config.GlobalRateLimit.Expiry.IsZero() &&
			time.Now().After(bot.Config.GlobalRateLimit.Expiry) {
			fmt.Println("expiry passed")
			bot.Config.GlobalRateLimit.Reset(time.Now().Add(time.Second))
		}
		bot.Config.GlobalRateLimit.muAtomic.Unlock()
	}

PRIORITY:
	bot.Config.GlobalRateLimit.Use(1)
	bot.Config.GlobalRateLimit.muQueue.Unlock()

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

	bot.Config.GlobalRateLimit.Confirm(1, date, sent)

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
		// increment priority to prevent requests in the queue from being sent.
		atomic.AddInt32(&bot.Config.GlobalRateLimit.Priority, 1)

		// parse the rate limit header.
		retryafter, err := peekHeader429(response)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		header := peekHeaderRateLimit(response)

		// parse the rate limit response data.
		var data RateLimitResponse
		if err := json.Unmarshal(response.Body(), &data); err != nil {
			return fmt.Errorf("%w", err)
		}

		fmt.Println("\t429", header, "Retry-After", retryafter, "retry_after", data.RetryAfter,
			fmt.Sprintf("\n%s%s", response.Header.String(), string(response.Body())))

		retry := retries < bot.Config.Retries
		retries++

		// expire the current rate limit immediately.
		reset := time.Now().Add(time.Millisecond * time.Duration(data.RetryAfter*1000))
		if header.Global {
			bot.Config.GlobalRateLimit.muQueue.Lock()
			fmt.Println("unlocked 429 from queue at", time.Now(), "\n\treset", reset)
			if !time.Now().After(reset) {
				fmt.Println("waiting", data.RetryAfter*1000, "ms due to 429")
				<-time.After(time.Until(reset))
				fmt.Println("finished wait due to 429")
				bot.Config.GlobalRateLimit.Reset(time.Now().Add(time.Second))
			}
			atomic.AddInt32(&bot.Config.GlobalRateLimit.Priority, -1)

			if retry {
				goto PRIORITY
			}

			defer bot.Config.GlobalRateLimit.muQueue.Unlock()
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
