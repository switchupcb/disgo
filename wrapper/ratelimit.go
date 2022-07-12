package wrapper

import (
	"strconv"
	"sync"
	"time"
)

// Routes for controlling emojis do not follow the normal rate limit conventions. These routes are specifically limited on a per-guild basis to prevent abuse. This means that the quota returned by our APIs may be inaccurate, and you may encounter 429s.

var (
	// DefaultRouteBucket represents the default rate limit Bucket for a route.
	//
	// The default rate limit Bucket is used to control the flow of the
	// "first request for any given route".
	//
	// This is necessary since Discord uses dynamic per-route rate limits.
	// As a result, a route's actual rate limit Bucket can NOT be discovered
	// until a request is sent (using the respective route).
	//
	// Use the default Bucket's Limit field-value to control
	// many requests of a given route can be sent (per second)
	// BEFORE the actual rate limit Bucket of the route is known.
	DefaultRouteBucket *Bucket = &Bucket{ //nolint:exhaustruct
		Limit: 1,
	}

	nilRouteBucket = "NIL"
)

// Bucket represents a Discord API Rate Limit Bucket.
type Bucket struct {
	// ID represents the Bucket ID.
	//
	// ID is only applicable to route rate limit Buckets.
	ID string

	// Limit represents the amount of requests a Bucket can send per reset.
	Limit int16

	// Remaining represents the amount of requests a Bucket can send until the next reset.
	Remaining int16

	// Pending represents the amount of requests that are sent and awaiting a response.
	Pending int16

	// Priority represents the amount of requests that have priority over
	// other requests in the Bucket.
	Priority int32

	// Date represents the time at which Discord received the last request of the Bucket.
	//
	// Date is only applicable to global rate limit Buckets.
	Date time.Time

	// Expiry represents the time at which the Bucket will reset (or become outdated).
	Expiry time.Time
}

// Reset resets a Discord API Rate Limit Bucket and sets its expiry.
func (b *Bucket) Reset(expiry time.Time) {
	b.Expiry = expiry

	// Remaining = Limit - Pending
	b.Remaining = b.Limit - b.Pending
}

// Use uses the given amount of tokens for a Discord API Rate Limit Bucket.
func (b *Bucket) Use(amount int16) {
	b.Remaining -= amount
	b.Pending += amount
}

// ConfirmDate confirms the usage of a given amount of tokens for a Discord API Rate Limit Bucket,
// using the bucket's current expiry and given (Discord Header) Date time.
//
// Used for the Global Rate Limit Bucket.
func (b *Bucket) ConfirmDate(amount int16, date time.Time) {
	b.Pending -= amount

	switch {
	// Date is zero when a request has never been sent to Discord.
	//
	// set the Date of the current Bucket to the date of the current Discord Bucket.
	case b.Date.IsZero():
		b.Date = date

		// The EXACT reset period of Discord's Global Rate Limit Bucket will always occur
		// BEFORE the current Bucket resets (due to this implementation).
		//
		// reset the current Bucket with an expiry that occurs [0, 1) seconds
		// AFTER the Discord Global Rate Limit Bucket will be reset.
		//
		// This results in a Bucket's expiry that is eventually consistent with
		// Discord's Bucket expiry over time (once determined between requests).
		b.Expiry = time.Now().Add(time.Second)

	// Date is EQUAL to the Discord Bucket's Date when the request applies to the current Bucket.
	case b.Date.Equal(date):

	// Date occurs BEFORE a Discord Bucket's Date when the request applies to the next Bucket.
	//
	// update the current Bucket to the next Bucket using Date.
	case b.Date.Before(date):
		b.Date = date

		// align the current Bucket's expiry to Discord's Bucket expiry.
		b.Expiry = time.Now().Add(time.Second)

	// Date occurs AFTER a Discord Bucket's Date when the request applied to a previous Bucket.
	case b.Date.After(date):
		b.Remaining += amount
	}
}

// ConfirmHeader confirms the usage of a given amount of tokens for a Discord API Rate Limit Bucket,
// using a give route ID and respective Discord Rate Limit Header.
//
// Used for Route Rate Limits.
func (b *Bucket) ConfirmHeader(amount int16, routeid uint16, header RateLimitHeader) {
	if b.Pending > 0 {
		b.Pending -= amount
	}

	b.ID = header.Bucket
	b.Limit = int16(header.Limit)
	b.Remaining = int16(header.Remaining) - b.Pending
	b.Expiry = time.Unix(header.Reset, 0)
}

// RateLimiter represents an interface for rate limits.
//
// RateLimiter is an interface which allows developers to use multi-application architectures,
// which run multiple applications on separate processes or servers.
type RateLimiter interface {
	// SetBucketHash maps a Route ID to a rate limit Bucket ID (Discord Hash).
	//
	// ID 0 is used as a Global Rate Limit Bucket or nil.
	SetBucketHash(uint16, string)

	// GetBucketHash gets a rate limit Bucket ID (Discord Hash) using a Route ID.
	GetBucketHash(uint16) string

	// SetBucketFromHash maps a Bucket ID to a rate limit Bucket.
	SetBucketFromHash(string, *Bucket)

	// GetBucketFromHash gets a rate limit Bucket using the given Bucket ID.
	GetBucketFromHash(string) *Bucket

	// SetBucket maps a Route ID to a rate limit Bucket.
	//
	// ID 0 is used as a Global Rate Limit Bucket or nil.
	SetBucket(uint16, *Bucket)

	// GetBucket gets a rate limit Bucket using the given Route ID.
	GetBucket(uint16) *Bucket

	// Lock locks the rate limiter.
	//
	// If the lock is already in use, the calling goroutine blocks until the rate limiter is available.
	//
	// This prevents multiple requests from being PROCESSED at once, which prevents race conditions.
	// In other words, a single request is PROCESSED from a rate limiter when Lock is implemented and called.
	//
	// This does NOT prevent multiple requests from being SENT at a time.
	Lock()

	// Unlock unlocks the rate limiter.
	//
	// If the rate limiter holds multiple locks, unlocking will unblock another goroutine,
	// which allows another request to be processed.
	Unlock()

	// StartTx starts a transaction with the rate limiter.
	//
	// If a transaction is already started, the calling goroutine blocks until the rate limiter is available.
	//
	// This prevents the transaction (of rate limit Bucket reads and writes) from concurrent manipulation.
	StartTx()

	// EndTx ends a transaction with the rate limiter.
	//
	// If the rate limiter holds multiple transactions, ending one will unblock another goroutine,
	// which allows another transaction to start.
	EndTx()
}

// RateLimit provides concurrency-safe rate limit functionality by implementing the RateLimiter interface.
type RateLimit struct {
	// ids represents a map of Route IDs to Bucket IDs (map[routeID]BucketID).
	ids map[uint16]string

	// buckets represents a map of Bucket IDs to rate limit Bucket (map[BucketID]*Bucket).
	buckets map[string]*Bucket

	// muQueue represents a mutex used to process a single request a time.
	muQueue sync.Mutex

	// muTx represents a mutex used to access multiple rate limit Buckets as a transaction.
	muTx sync.Mutex
}

func (r *RateLimit) SetBucketHash(id uint16, bucketid string) {
	// Default Route Bucket IDs are artificial IDs represented by
	// a string-formatted number (routeID).
	//
	// when a route ID is re-routed from a Default Bucket, remove the Default Bucket.
	if s := strconv.FormatUint(uint64(id), base10); s == r.GetBucketHash(id) {
		delete(r.buckets, s)
	}

	r.ids[id] = bucketid
}

func (r *RateLimit) GetBucketHash(id uint16) string {
	return r.ids[id]
}

func (r *RateLimit) SetBucketFromHash(bucketid string, bucket *Bucket) {
	r.buckets[bucketid] = bucket
}

func (r *RateLimit) GetBucketFromHash(bucketid string) *Bucket {
	return r.buckets[bucketid]
}

func (r *RateLimit) SetBucket(id uint16, bucket *Bucket) {
	r.buckets[r.ids[id]] = bucket
}

func (r *RateLimit) GetBucket(id uint16) *Bucket {
	// ID 0 is used as a Global Rate Limit Bucket or nil.
	if id != 0 {
		switch r.ids[id] {
		// when a non-global route is initialized and nil (see explanation below).
		case nilRouteBucket:
			return nil

		// This rate limiter implementation points ID 0 (which is reserved for a
		// Global Rate Limiter) to Bucket ID "".
		//
		// As a result (of the Default Bucket), non-0 IDs must be handled accordingly.
		case "":
			// when a non-global route is uninitialized, set it to the default bucket.
			//
			// While GetBucket can be called multiple times BEFORE a request is sent,
			// routeBucket.GetBucket is only called when the global rate limit has
			// been validated. As a result, the bucket that is allocated from this
			// call will ALWAYS be immediately used.
			//
			// The Bucket ID (Hash) is an artificial value while the route ID
			// is pointed to a Default Bucket. This results in subsequent calls to the
			// route's Default Bucket to return the same bucket.
			//
			// When a default bucket is exhausted, it will never expire.
			// As a result, the Remaining field-value will remain at 0 until
			// the pending request (and its respective Bucket) is confirmed.
			s := strconv.FormatUint(uint64(id), base10)
			r.SetBucketHash(id, s)
			if DefaultRouteBucket == nil {
				return nil
			}

			b := &Bucket{Remaining: DefaultRouteBucket.Limit} //nolint:exhaustruct
			r.SetBucketFromHash(s, b)

			return b
		}
	}

	return r.buckets[r.ids[id]]
}

func (r *RateLimit) Lock() {
	r.muQueue.Lock()
}

func (r *RateLimit) Unlock() {
	r.muQueue.Unlock()
}

func (r *RateLimit) StartTx() {
	r.muTx.Lock()
}

func (r *RateLimit) EndTx() {
	r.muTx.Unlock()
}
