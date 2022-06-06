package wrapper

import (
	"fmt"
	"strconv"

	"github.com/valyala/fasthttp"
)

// Global rate limits apply to the total number of requests a bot or user makes, independent of any per-route limits. You can read more on global rate limits below.
// Per-route rate limits exist for many individual endpoints, and may include the HTTP method (GET, POST, PUT, or DELETE). In some cases, per-route limits will be shared across a set of similar endpoints, indicated in the X-RateLimit-Bucket header. It's recommended to use this header as a unique identifier for a rate limit, which will allow you to group shared limits as you encounter them.
// Routes for controlling emojis do not follow the normal rate limit conventions. These routes are specifically limited on a per-guild basis to prevent abuse. This means that the quota returned by our APIs may be inaccurate, and you may encounter 429s.

// HTTP Header Rate Limit variables.
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
)

// peekRateLimitHeader peeks the Rate Limit Header values in an HTTP Header.
func peekRateLimitHeader(r *fasthttp.Response) (*RateLimitHeader, error) {
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
