package wrapper

import (
	"sync"
	"time"
)

// Per-route rate limits exist for many individual endpoints, and may include the HTTP method (GET, POST, PUT, or DELETE). In some cases, per-route limits will be shared across a set of similar endpoints, indicated in the X-RateLimit-Bucket header. It's recommended to use this header as a unique identifier for a rate limit, which will allow you to group shared limits as you encounter them.
// Routes for controlling emojis do not follow the normal rate limit conventions. These routes are specifically limited on a per-guild basis to prevent abuse. This means that the quota returned by our APIs may be inaccurate, and you may encounter 429s.
// Interaction endpoints are not bound to the bot's Global Rate Limit.

// RateLimiter represents an interface for rate limits.
//
// RateLimiter is an interface which allows developers to use multi-application architectures,
// which run multiple applications on separate processes or servers.
type RateLimiter interface {
	// GetBucket gets a rate limiter using the given bucket ID.
	GetBucket(string) *Bucket

	// SetBucket maps a bucket ID to a rate limit bucket.
	SetBucket(string, *Bucket)
}

// Bucket represents a Discord API Rate Limit Bucket.
type Bucket struct {
	Limit     int
	Remaining int
	Priority  int
	Expiry    time.Time
	mu        sync.Mutex
}

// RateLimit provides concurrency-safe rate limit functionality by implementing the RateLimiter interface.
type RateLimit struct {
	// buckets represents a synchronized map of buckets to rate limiters (map[string]*Bucket).
	buckets sync.Map
}

func (r *RateLimit) GetBucket(id string) *Bucket {
	if v, ok := r.buckets.Load(id); ok {
		return v.(*Bucket) //nolint:forcetypeassert
	}

	return nil
}

func (r *RateLimit) SetBucket(id string, bucket *Bucket) {
	r.buckets.Store(id, bucket)
}
