package wrapper

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

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
	Limit     int32
	Remaining int32
	Pending   int32
	Expiry    time.Time
	mu        sync.Mutex
	muExpiry  sync.RWMutex
}

// Reset resets a Discord API Rate Limit Bucket and sets its expiry.
func (b *Bucket) Reset(expiry time.Time) {
	b.muExpiry.Lock()
	b.Expiry = expiry
	b.muExpiry.Unlock()

	// Remaining = Limit - (Pending + Priority Requests)
	atomic.StoreInt32(&b.Remaining, b.Limit-(atomic.LoadInt32(&b.Pending)))
	fmt.Println("reset expired bucket at", time.Now(),
		"\n\tnow", time.Now(), "remain", atomic.LoadInt32(&b.Remaining), "pending", atomic.LoadInt32(&b.Pending),
		"\n\tto", b.Expiry, "remain", atomic.LoadInt32(&b.Remaining))
}

// Use uses the given amount of tokens for a Discord API Rate Limit Bucket.
func (b *Bucket) Use(amount int32) {
	atomic.AddInt32(&b.Remaining, -amount)
	atomic.AddInt32(&b.Pending, amount)
	fmt.Println("sent", time.Now(), atomic.LoadInt32(&b.Remaining), atomic.LoadInt32(&b.Pending))
}

// Confirm confirms the usage of a given amount of tokens for a Discord API Rate Limit Bucket,
// using the bucket's current expiry and given (Discord Header) Date time.
func (b *Bucket) Confirm(amount int32, date time.Time, temp time.Time) {
	fmt.Println("received", date,
		"\n\tfrom", temp,
		"\n\tat", time.Now(), atomic.LoadInt32(&b.Remaining), atomic.LoadInt32(&b.Pending))
	atomic.AddInt32(&b.Pending, -amount)

	// NOTE: The following edge cases are currently optimized for global rate limits.
	b.muExpiry.RLock()

	// when a request was received by Discord AFTER or ON the current rate limit bucket expiry,
	// account for it the NEXT rate limit bucket.
	//
	// i.e Request received by Discord at 4s, Confirmation at 4s, Bucket.Reset() at 4.01s.
	//
	// The received pending request (date = 4s) is deducted from the expiring bucket (with expiry = 4s).
	// As a result, the bucket resets without a pending request at 4.01s which results in one more
	// remaining request being allocated than allowed in the new bucket (with expiry = 5s).
	if date.Equal(b.Expiry) {
		atomic.AddInt32(&b.Remaining, -amount)
		fmt.Println("\taccount next", date,
			"\n\t\tnow", time.Now(), atomic.LoadInt32(&b.Remaining), atomic.LoadInt32(&b.Pending),
			"\n\t\texp", b.Expiry)
		b.muExpiry.RUnlock()
		return
	}

	// when a request was received by Discord BEFORE a rate limit bucket reset,
	// but not known at time of reset, account for its confirmation in hindsight.
	//
	// i.e Request received by Discord at 3.9s, Bucket.Reset() at 4s, Confirmation at 4.1s.
	//
	// The bucket resets with a pending request at 4s which results in one less remaining request.
	// However, this pending request was received (date = 3s) in a - now expired - bucket (with expiry = 4s)
	// so the new bucket (with expiry = 5s) can send one more request than it allocated.
	if date.Equal(b.Expiry.Add((-time.Second * 2))) {
		atomic.AddInt32(&b.Remaining, amount)
		fmt.Println("\taccount prior\n\tdate", date,
			"\n\t\tnow", time.Now(), atomic.LoadInt32(&b.Remaining), atomic.LoadInt32(&b.Pending),
			"\n\t\texp", b.Expiry)
		b.muExpiry.RUnlock()
		return
	}

	b.muExpiry.RUnlock()
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
