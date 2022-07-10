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
	Limit     int16
	Remaining int16
	Pending   int16
	Priority  int32
	Date      time.Time
	Expiry    time.Time
	muQueue   sync.Mutex
	muAtomic  sync.Mutex
}

// Reset resets a Discord API Rate Limit Bucket and sets its expiry.
func (b *Bucket) Reset(expiry time.Time) {
	b.Expiry = expiry

	// Remaining = Limit - (Pending + Priority Requests)
	b.Remaining = b.Limit - (b.Pending + int16(atomic.LoadInt32(&b.Priority)))
	fmt.Println("reset expired bucket at", time.Now(),
		"\n\tto", b.Expiry, "remain", b.Remaining, "pending", b.Pending, "priority", atomic.LoadInt32(&b.Priority))
}

// Use uses the given amount of tokens for a Discord API Rate Limit Bucket.
func (b *Bucket) Use(amount int16) {
	b.muAtomic.Lock()
	defer b.muAtomic.Unlock()
	b.Remaining -= amount
	b.Pending += amount
	fmt.Println("sent", time.Now(), b.Remaining, b.Pending, atomic.LoadInt32(&b.Priority))
}

// Confirm confirms the usage of a given amount of tokens for a Discord API Rate Limit Bucket,
// using the bucket's current expiry and given (Discord Header) Date time.
func (b *Bucket) Confirm(amount int16, date time.Time, temp time.Time) {
	b.muAtomic.Lock()
	defer b.muAtomic.Unlock()
	fmt.Println("received", date,
		"\n\tfrom", temp,
		"\n\tat", time.Now(), b.Remaining, b.Pending, atomic.LoadInt32(&b.Priority))
	b.Pending -= amount

	switch {
	// Date is zero when a request has never been sent to Discord.
	//
	// set the Date of the current Bucket to the date of the current Discord Bucket.
	case b.Date.IsZero():
		b.Date = date
		b.Expiry = time.Now().Add(time.Second)

	// Date is EQUAL to the Discord Bucket's Date when the request applies to the current Bucket.
	case b.Date.Equal(date):

	// Date occurs AFTER a Discord Bucket's Date when the request applied to a previous Bucket.
	case b.Date.After(date):
		b.Remaining += amount
		fmt.Println("\taccount prior\n\tDate", b.Date, "Discord Date", date,
			"\n\t\tnow", time.Now(), b.Remaining, b.Pending, atomic.LoadInt32(&b.Priority),
			"\n\t\texp", b.Expiry)

	// Date occurs BEFORE a Discord Bucket's Date when the request applies to the next Bucket.
	//
	// set the current Bucket to the next Bucket using Date (and reset the new Bucket).
	case b.Date.Before(date):
		b.Date = date

		// The EXACT reset period of Discord's Global Rate Limit Bucket will always occur
		// BEFORE the current Bucket resets (due to this implementation).
		//
		// reset the current Bucket with an expiry that occurs a minimum of one second
		// AFTER the Discord Global Rate Limit Bucket was reset.
		//
		// This results in a Bucket's expiry that is eventually consistent with
		// Discord's Bucket expiry over time (and determined between requests).
		b.Reset(time.Now().Add(time.Second))
		b.Remaining -= amount
		fmt.Println("\taccount next\n\tDate", b.Date, "Discord Date", date,
			"\n\t\tnow", time.Now(), b.Remaining, b.Pending, atomic.LoadInt32(&b.Priority),
			"\n\t\texp", b.Expiry)
	}
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
