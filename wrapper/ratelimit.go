package wrapper

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Per-route rate limits exist for many individual endpoints, and may include the HTTP method (GET, POST, PUT, or DELETE). In some cases, per-route limits will be shared across a set of similar endpoints, indicated in the X-RateLimit-Bucket header. It's recommended to use this header as a unique identifier for a rate limit, which will allow you to group shared limits as you encounter them.
// Routes for controlling emojis do not follow the normal rate limit conventions. These routes are specifically limited on a per-guild basis to prevent abuse. This means that the quota returned by our APIs may be inaccurate, and you may encounter 429s.

// RateLimiter represents an interface for rate limits.
//
// RateLimiter is an interface which allows developers to use multi-application architectures,
// which run multiple applications on separate processes or servers.
type RateLimiter interface {
	// GetBucket gets a rate limiter using the given bucket.
	GetBucket(string) *rate.Limiter

	// GetBucketExpiry gets the Epoch time at which the rate limit for the given bucket resets.
	GetBucketExpiry(string) time.Time

	// SetBucket sets a bucket with a rate limiter.
	SetBucket(string, *rate.Limiter)

	// SetBucketExpiry sets the Epoch time at which the rate limits for the given bucket resets.
	SetBucketExpiry(string, time.Time)
}

// RateLimit provides concurrency-safe rate limit functionality by implementing the RateLimiter interface.
type RateLimit struct {
	// buckets represents a synchronized map of buckets to rate limiters (map[string]*rate.Limiter).
	buckets *sync.Map

	// expiry represents a synchronized map of buckets to the time at which they expire (map[string]*rate.Limiter).
	expiry *sync.Map
}

func (r RateLimit) GetBucket(bucket string) *rate.Limiter {
	if v, ok := r.buckets.Load(bucket); ok {
		return v.(*rate.Limiter)
	}

	return nil
}

func (r RateLimit) GetBucketExpiry(bucket string) time.Time {
	if v, ok := r.expiry.Load(bucket); ok {
		return v.(time.Time)
	}

	return time.Now().Add(time.Hour)
}

func (r RateLimit) SetBucket(bucket string, ratelimit *rate.Limiter) {
	r.buckets.Store(bucket, ratelimit)
}

func (r RateLimit) SetBucketExpiry(bucket string, expiry time.Time) {
	r.expiry.Store(bucket, expiry)
}
