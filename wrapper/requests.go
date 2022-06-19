package wrapper

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	json "github.com/goccy/go-json"
	"github.com/gorilla/schema"
	"github.com/valyala/fasthttp"
	"golang.org/x/time/rate"
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
	// TODO: requires clarification on whether Discord uses:
	// 	A. rolling rate limit which allows 50 burst then a request every 1/50th second.
	// 	B. bucket rate limit which allows one burst request every 1/50th second.
	GlobalRateLimit = rate.NewLimiter(rate.Limit(FlagGlobalRequestRateLimit), FlagGlobalRequestRateLimit)
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
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	request.Header.SetMethod(method)
	request.Header.SetContentTypeBytes(content)
	request.Header.Set(headerAuthorizationKey, bot.Authentication.Header)
	request.SetRequestURI(uri)
	request.SetBodyRaw(body)

	// receive the response from the request.
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)

	// send the request.
	retries := 0

	// check the ratelimit of the request.
	ratelimit := bot.Config.RateLimiter.GetBucket("bucket???")
	if ratelimit != nil {
		if err := ratelimit.Wait(context.Background()); err != nil {
			// when the burst size (which should be >= 1) is 0.
			return fmt.Errorf("%w", err)
		}
	}

	// TODO: Handle Global Rate Limit.

SEND:
	if err := bot.Config.Client.DoTimeout(request, response, bot.Config.Timeout); err != nil {
		return fmt.Errorf("%w", err)
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
		// process the ratelimit.
		if ratelimit == nil || time.Now().After(bot.Config.RateLimiter.GetBucketExpiry("bucket???")) {
			header, err := peekHeaderRateLimit(response)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			// create a new rate limiter that allows "Remaining" requests to be made
			// from now until the bucket expires.
			newRateLimit := rate.NewLimiter(
				rate.Limit((float64(time.Second)*header.ResetAfter)/float64(header.Remaining)),
				1,
			)

			bot.Config.RateLimiter.SetBucket(header.Bucket, newRateLimit)
			bot.Config.RateLimiter.SetBucketExpiry(header.Bucket, time.Unix(header.Reset, 0))
		}

		// parse the response data.
		if err := json.Unmarshal(response.Body(), dst); err != nil {
			return fmt.Errorf("%w", err)
		}

		return nil

	// process the ratelimit.
	case fasthttp.StatusTooManyRequests:
		header, err := peekHeaderRateLimit(response)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		retryafter, err := peekHeader429(response)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		// expire the current ratelimit immediately (unless it's global).
		if !header.Global {
			bot.Config.RateLimiter.SetBucket(header.Bucket, nil)
			bot.Config.RateLimiter.SetBucketExpiry(header.Bucket, time.Now())
		} else {
			// TODO: Handle Global Rate Limit Trigger.
		}

		// wait until the "Retry-After" header to retry the request.
		<-time.After(time.Second * time.Duration(retryafter))

		// send another request (which sets a valid rate limit bucket).
		if retries < bot.Config.Retries {
			retries++

			goto SEND
		}

		return StatusCodeError(fasthttp.StatusTooManyRequests)

	// retry the request on a bad gateway server error.
	case fasthttp.StatusBadGateway:
		if retries < bot.Config.Retries {
			retries++

			goto SEND
		}

		return StatusCodeError(fasthttp.StatusBadGateway)

	default:
		return StatusCodeError(response.StatusCode())
	}
}

// peekHeaderRateLimit peeks an HTTP Header for Rate Limit Header values.
func peekHeaderRateLimit(r *fasthttp.Response) (*RateLimitHeader, error) {
	// TODO: Ensure correct []byte to int, int64, float64, bool conversion.
	fmt.Println(r.Header.String())

	limit, err := strconv.Atoi(string(r.Header.PeekBytes(headerRateLimit)))
	if err != nil {
		return nil, fmt.Errorf("an error occurred converting a rate limit header to an int:\n%w", err)
	}

	remaining, err := strconv.Atoi(string(r.Header.PeekBytes(headerRateLimitRemaining)))
	if err != nil {
		return nil, fmt.Errorf("an error occurred converting a rate limit header to an int:\n%w", err)
	}

	reset, err := strconv.ParseInt(string(r.Header.PeekBytes(headerRateLimitReset)), base10, bit64)
	if err != nil {
		return nil, fmt.Errorf("an error occurred converting a rate limit header to an int64:\n%w", err)
	}

	resetafter, err := strconv.ParseFloat(string(r.Header.PeekBytes(headerRateLimitResetAfter)), bit64)
	if err != nil {
		return nil, fmt.Errorf("an error occurred converting a rate limit header to a float64:\n%w", err)
	}

	global, err := strconv.ParseBool(string(r.Header.PeekBytes(headerRateLimitGlobal)))
	if err != nil {
		return nil, fmt.Errorf("an error occurred converting a rate limit header to a bool:\n%w", err)
	}

	ratelimit := &RateLimitHeader{
		Limit:      limit,
		Remaining:  remaining,
		Reset:      reset,
		ResetAfter: resetafter,
		Bucket:     string(r.Header.PeekBytes(headerRateLimitBucket)),
		Global:     global,
		Scope:      string(r.Header.PeekBytes(headerRateLimitScope)),
	}

	return ratelimit, nil
}

// peekHeader429 peeks an HTTP Header with a 429 Status Code for the Rate Limit Header "RetryAfter".
func peekHeader429(r *fasthttp.Response) (int64, error) {
	retryafter, err := strconv.ParseInt(string(r.Header.PeekBytes(headerRateLimitRetryAfter)), base10, bit64)
	if err != nil {
		return 0, fmt.Errorf("an error occurred converting a rate limit header to an int64:\n%w", err)
	}

	return retryafter, nil
}
