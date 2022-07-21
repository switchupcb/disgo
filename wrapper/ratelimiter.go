package wrapper

import (
	"time"
)

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
	b.Expiry = time.Unix(int64(header.Reset)+1, 0)
}
