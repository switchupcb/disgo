package wrapper

import (
	"log"
	"strconv"
	"sync"
)

var (
	// DefaultRouteBucket represents the default rate limit Bucket for a route.
	//
	// The Default Rate Limit Bucket is used to control the flow of the
	// "first request for any given route".
	//
	// This is necessary since Discord uses dynamic per-route rate limits.
	// As a result, a route's actual rate limit Bucket can NOT be discovered
	// until a request is sent (using the respective route).
	//
	// Use the Default Bucket's Limit field-value to control how
	// many requests of a given route can be sent (per second)
	// BEFORE the actual rate limit Bucket of the route is known.
	DefaultRouteBucket *Bucket = &Bucket{ //nolint:exhaustruct
		Limit: 1,
	}

	nilRouteBucket = "NIL"
)

// RateLimit provides concurrency-safe rate limit functionality by implementing the RateLimiter interface.
type RateLimit struct {
	// ids represents a map of Route IDs to Bucket IDs (map[routeID]BucketID).
	ids map[uint16]string

	// buckets represents a map of Bucket IDs to rate limit Bucket (map[BucketID]*Bucket).
	buckets map[string]*Bucket

	// entries represents a map of Bucket IDs to a count of references to Routes (map[BucketID]count).
	//
	// Used to safely remove a Bucket once it's no longer in use.
	entries map[string]int

	// muQueue represents a mutex used to process a single request a time.
	muQueue sync.Mutex

	// muTx represents a mutex used to access multiple rate limit Buckets as a transaction.
	muTx sync.Mutex
}

func (r *RateLimit) SetBucketHash(routeid uint16, bucketid string) {
	currentBucketID := r.ids[routeid]

	// when the current Bucket ID is not the same as the new Bucket ID.
	if currentBucketID != bucketid {
		// update the entries for the current Bucket ID.
		if currentBucketID != "" {
			r.entries[currentBucketID]--

			// when the current Bucket ID is no longer referenced by a Route,
			// delete the respective Bucket (to allow Garbage Collection).
			if r.entries[currentBucketID] <= 0 {
				putBucket(r.buckets[currentBucketID])
				delete(r.entries, currentBucketID)
				delete(r.buckets, currentBucketID)

				log.Println("deleted bucket", currentBucketID)
			}
		}

		// set the Route ID to the new Bucket ID.
		r.ids[routeid] = bucketid

		// update the entries for the new Bucket ID.
		r.entries[bucketid]++

		log.Println("set route", routeid, "to bucket", bucketid)
	}
}

func (r *RateLimit) GetBucketHash(routeid uint16) string {
	return r.ids[routeid]
}

func (r *RateLimit) SetBucketFromHash(bucketid string, bucket *Bucket) {
	r.buckets[bucketid] = bucket

	log.Printf("set bucket %s to %p", bucketid, bucket)
}

func (r *RateLimit) GetBucketFromHash(bucketid string) *Bucket {
	return r.buckets[bucketid]
}

func (r *RateLimit) SetBucket(routeid uint16, bucket *Bucket) {
	r.buckets[r.ids[routeid]] = bucket
}

func (r *RateLimit) GetBucket(routeid uint16) *Bucket {
	// ID 0 is used as a Global Rate Limit Bucket or nil.
	if routeid != 0 {
		switch r.ids[routeid] {
		// when a non-global route is initialized and nil (see explanation below).
		case nilRouteBucket:
			return nil

		// This rate limiter implementation points ID 0 (which is reserved for a
		// Global Rate Limiter) to Bucket ID "".
		//
		// As a result (of the Default Bucket), non-0 Route IDs must be handled accordingly.
		case "":
			// when a non-global route is uninitialized, set it to the default bucket.
			//
			// While GetBucket can be called multiple times BEFORE a request is sent,
			// routeBucket.GetBucket is only called when the global rate limit has
			// been validated. As a result, the bucket that is allocated from this
			// call will ALWAYS be immediately used.
			//
			// The Bucket ID (Hash) is an artificial value while the Route ID
			// is pointed to a Default Bucket. This results in subsequent calls to the
			// Route's Default Bucket to return the same bucket.
			//
			// When a Default Bucket is exhausted, it will never expire.
			// As a result, the Remaining field-value will remain at 0 until
			// the pending request (and its respective Bucket) is confirmed.
			s := strconv.FormatUint(uint64(routeid), base10)
			r.SetBucketHash(routeid, s)
			if DefaultRouteBucket == nil {
				return nil
			}

			b := getBucket()
			b.Remaining = DefaultRouteBucket.Limit
			r.SetBucketFromHash(s, b)

			return b
		}
	}

	return r.buckets[r.ids[routeid]]
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
