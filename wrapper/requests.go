package wrapper

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	json "github.com/goccy/go-json"
	"github.com/gorilla/schema"
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
)

var (
	// https://discord.com/developers/docs/topics/rate-limits#global-rate-limit
	GlobalRateLimit = &Bucket{
		Limit:     FlagGlobalRequestRateLimit,
		Remaining: FlagGlobalRequestRateLimit,
	}

	// second is a temporary variable representing a second accounted for estimated latency.
	second = time.Second
)

var (
	// qsEncoder is used to create URL Query Strings from objects.
	qsEncoder = schema.NewEncoder()
)

// init runs at the start of the program.
func init() {
	// use `url` tags for the URL Query String encoder and decoder.
	qsEncoder.SetAliasTag("url")
}

// EndpointQueryString returns a URL Query String from a given object.
func EndpointQueryString(dst any) (string, error) {
	params := url.Values{}
	err := qsEncoder.Encode(dst, params)
	if err != nil {
		return "", err //nolint:wrapcheck
	}

	return params.Encode(), nil
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

RATELIMIT:
	// check the global rate limit prior to sending the request.
	fmt.Println("queued request")
	bot.Config.GlobalRateLimit.mu.Lock()
	fmt.Println("unlocked a request in the queue.")

	// when priority is set to true, it indicates that
	// a pending request is waiting to be sent first.
	if bot.Config.GlobalRateLimit.Priority > 0 {
		// unlock and re-lock this request's mutex to give the pending request priority
		// by sending this request (and subsequent requests) after that pending request.
		bot.Config.GlobalRateLimit.mu.Unlock()

		goto RATELIMIT
	}

	// reset the global rate limit bucket when the current time has passed its expiry.
	if !bot.Config.GlobalRateLimit.Expiry.IsZero() && time.Now().After(bot.Config.GlobalRateLimit.Expiry) {
		fmt.Println("reset expired bucket.")
		reset(bot.Config.GlobalRateLimit, second)
	}

	// when no requests remain in the global rate limit bucket,
	// wait until the bucket resets to send a request.
	bot.Config.GlobalRateLimit.muRem.Lock()
	if bot.Config.GlobalRateLimit.Remaining <= 0 {
		bot.Config.GlobalRateLimit.muRem.Unlock()

		fmt.Println("0 remaining... queued.")
		// expiry is zero when the bot has NEVER received a response from Discord.
		if bot.Config.GlobalRateLimit.Expiry.IsZero() {
			fmt.Println("first queued request")
			// a rate limit bucket with a limit of 0 indicates that
			// the bot has NEVER sent a response to Discord.
			if bot.Config.GlobalRateLimit.Limit == 0 {
				bot.Config.GlobalRateLimit.mu.Unlock()
				return fmt.Errorf("cannot use a global rate limit bucket with a Limit of zero")
			}

			// wait until Discord receives the first response of a request.
			fmt.Println("waiting for marked expiry")
			for {
				if !bot.Config.GlobalRateLimit.Expiry.IsZero() {
					fmt.Println("marked expiry")
					<-time.After(time.Until(bot.Config.GlobalRateLimit.Expiry))

					bot.Config.GlobalRateLimit.muRem.Lock()
					fmt.Println(time.Now(), "finished wait", bot.Config.GlobalRateLimit.Remaining)
					reset(bot.Config.GlobalRateLimit, second)
					fmt.Println("reset bucket with 0 remaining.")
					bot.Config.GlobalRateLimit.muRem.Unlock()

					fmt.Println("sending first queued request")
					break
				}
			}
		} else {
			fmt.Println("other queued request")
			// wait until the current bucket expires, then reset it.
			<-time.After(time.Until(bot.Config.GlobalRateLimit.Expiry))
			bot.Config.GlobalRateLimit.muRem.Lock()
			reset(bot.Config.GlobalRateLimit, second)
			bot.Config.GlobalRateLimit.muRem.Unlock()
		}

		// if the request was blocked due to regular bucket expiry, send it.
		if bot.Config.GlobalRateLimit.Priority == 0 {
			bot.Config.GlobalRateLimit.muRem.Lock()
			bot.Config.GlobalRateLimit.Remaining--
			bot.Config.GlobalRateLimit.muRem.Unlock()
			bot.Config.GlobalRateLimit.mu.Unlock()

			goto SEND
		}

		// if the request was blocked due to a 429 rate limit, re-queue it.
		bot.Config.GlobalRateLimit.muRem.Lock()
		bot.Config.GlobalRateLimit.Remaining--
		bot.Config.GlobalRateLimit.muRem.Unlock()
		bot.Config.GlobalRateLimit.mu.Unlock()

		goto RATELIMIT
	}

	bot.Config.GlobalRateLimit.Remaining--
	bot.Config.GlobalRateLimit.muRem.Unlock()
	bot.Config.GlobalRateLimit.mu.Unlock()

	goto SEND

PRIORITY:
	// ensure that the first request to read a 429 header is sent first,
	// and followed by subsequent requests.
	bot.Config.GlobalRateLimit.mu.Lock()

	// ensure that the first request is sent after "Retry-After" seconds,
	// if there were no subsequent requests queued.
	<-time.After(time.Until(bot.Config.GlobalRateLimit.Expiry))
	bot.Config.GlobalRateLimit.muRem.Lock()
	reset(bot.Config.GlobalRateLimit, second)
	bot.Config.GlobalRateLimit.Remaining--
	bot.Config.GlobalRateLimit.muRem.Unlock()
	bot.Config.GlobalRateLimit.Priority--
	bot.Config.GlobalRateLimit.mu.Unlock()

SEND:
	// send the request.
	fmt.Println(time.Now(), "sent")
	if err := bot.Config.Client.DoTimeout(request, response, bot.Config.Timeout); err != nil {
		return fmt.Errorf("%w", err)
	}
	fmt.Println("received", time.Now())

	// Set the expiry when the bot has received its FIRST response from Discord.
	//
	// Expiry will ONLY be set (multiple times) within the FIRST BUCKET of responses,
	// which may result in the FIRST BUCKET expiring late. However, this is fine since
	// subsequent buckets will be consistent with the exact rate limit.
	if bot.Config.GlobalRateLimit.Expiry.IsZero() {
		bot.Config.GlobalRateLimit.Expiry = time.Now().Add(second)
		fmt.Println("set expiry to ", bot.Config.GlobalRateLimit.Expiry)
	}

	// if a request was received AFTER the rate limit bucket expired,
	// it occurred in the next rate limit bucket.
	if time.Now().After(bot.Config.GlobalRateLimit.Expiry) {
		fmt.Println("received after")
		bot.Config.GlobalRateLimit.muRem.Lock()
		bot.Config.GlobalRateLimit.Remaining--
		bot.Config.GlobalRateLimit.muRem.Unlock()
	}

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
		timeAt429 := time.Now()
		fmt.Println("429 PARSED HEADER", header, "PARSED RETRY AFTER", retryafter)

		retry := retries < bot.Config.Retries
		retries++

		// expire the current rate limit immediately.
		if header.Global {
			bot.Config.GlobalRateLimit.mu.Lock()
			bot.Config.GlobalRateLimit.muRem.Lock()
			bot.Config.GlobalRateLimit.Remaining = 0
			bot.Config.GlobalRateLimit.muRem.Unlock()

			// set the global rate limit expiry such that the next request
			// will be sent after "Retry-After" seconds.
			bot.Config.GlobalRateLimit.Expiry = timeAt429.Add(time.Second * time.Duration(retryafter))

			if retry {
				// set priority to indicate that this request has priority when the bucket resets.
				bot.Config.GlobalRateLimit.Priority++
				bot.Config.GlobalRateLimit.mu.Unlock()

				goto PRIORITY
			}

			bot.Config.GlobalRateLimit.mu.Unlock()
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

func reset(bucket *Bucket, t time.Duration) {
	fmt.Println("old expiry of reset bucket at ", bucket.Expiry)
	bucket.Expiry = bucket.Expiry.Add(t)
	fmt.Println("new expiry of reset bucket at ", bucket.Expiry)
	if bucket.Remaining > 0 {
		bucket.Remaining = bucket.Limit
	} else {
		bucket.Remaining += bucket.Limit
	}
}
