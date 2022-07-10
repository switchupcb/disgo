package wrapper

import (
	"sync"
	"time"
)

// Routes for controlling emojis do not follow the normal rate limit conventions. These routes are specifically limited on a per-guild basis to prevent abuse. This means that the quota returned by our APIs may be inaccurate, and you may encounter 429s.
// Interaction endpoints are not bound to the bot's Global Rate Limit.

// Bucket represents a Discord API Rate Limit Bucket.
type Bucket struct {
	Limit     int16
	Remaining int16
	Pending   int16
	Priority  int32
	Date      time.Time
	Expiry    time.Time
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

// Confirm confirms the usage of a given amount of tokens for a Discord API Rate Limit Bucket,
// using the bucket's current expiry and given (Discord Header) Date time.
func (b *Bucket) Confirm(amount int16, date time.Time) {
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

// RateLimiter represents an interface for rate limits.
//
// RateLimiter is an interface which allows developers to use multi-application architectures,
// which run multiple applications on separate processes or servers.
type RateLimiter interface {
	// GetBucket gets a rate limit Bucket using the Route ID.
	GetBucket(uint16) *Bucket

	// SetBucket maps a rate limit Bucket to a given Route ID.
	SetBucket(uint16, *Bucket)

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
	// buckets represents a synchronized map of Route IDs to rate limit Buckets (map[uint16]*Bucket).
	buckets sync.Map

	// muQueue represents a mutex used to process a single request a time.
	muQueue sync.Mutex

	// muTx represents a mutex used to access multiple rate limit Buckets as a transaction.
	muTx sync.Mutex
}

func (r *RateLimit) GetBucket(id uint16) *Bucket {
	if v, ok := r.buckets.Load(id); ok {
		return v.(*Bucket) //nolint:forcetypeassert
	}

	return nil
}

func (r *RateLimit) SetBucket(id uint16, bucket *Bucket) {
	r.buckets.Store(id, bucket)
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
