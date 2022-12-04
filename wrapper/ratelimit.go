package wrapper

import (
	"sync"
)

const (
	nilRouteBucket = "NIL"
)

var (
	// IgnoreGlobalRateLimitRouteIDs represents a set of Route IDs that do NOT adhere to the Global Rate Limit.
	//
	// Interaction endpoints are not bound to the bot's Global Rate Limit.
	// https://discord.com/developers/docs/interactions/receiving-and-responding#endpoints
	IgnoreGlobalRateLimitRouteIDs = map[string]bool{
		"18": true, "19": true, "20": true, "21": true, "22": true, "23": true, "24": true, "25": true,
	}
)

// RateLimit provides concurrency-safe rate limit functionality by implementing the RateLimiter interface.
type RateLimit struct {
	// ids represents a map of Route IDs to Bucket IDs (map[routeID]BucketID).
	ids map[string]string

	// buckets represents a map of Bucket IDs to rate limit Bucket (map[BucketID]*Bucket).
	buckets map[string]*Bucket

	// entries represents a map of Bucket IDs to a count of references to Routes (map[BucketID]count).
	//
	// Used to safely remove a Bucket once it's no longer in use.
	entries map[string]int

	// DefaultBucket represents a Default Rate Limit Bucket, which is used to control
	// the rate of the "first request(s) for any given route".
	//
	// This is necessary since Discord uses dynamic per-route token rate limits.
	// As a result, a route's actual Rate Limit Bucket can NOT be discovered until a request (for that route) is sent.
	//
	// A Default Rate Limit Bucket can be set at multiple levels.
	//	Route (RateLimit.DefaultBucket):    Used when a per-route request's bucket is NOT initialized (i.e route 16*).
	// 	Resource (Route Bucket Hash):       Used when a per-resource request's bucket is NOT initialized (i.e route 16, resource 32*).
	// 	Resource(s) (Resource Bucket Hash): Used when an nth degree per-resource request's bucket is NOT initialized (i.e route 16; resource 32; resource 7*).
	// 	So on and so forth...
	//	*A request is NOT initialized when it has never been set (to a bucket or nil).
	//
	// Set the DefaultBucket to `nil` to disable the Default Rate Limit Bucket mechanism.
	//
	// Use a Default Bucket's Limit field-value to control how many requests of
	// a given route can be sent (per second) BEFORE the actual Rate Limit Bucket of that route is known.
	DefaultBucket *Bucket

	// muQueue represents a mutex used to process a single request a time.
	muQueue sync.Mutex

	// muTx represents a mutex used to access multiple rate limit Buckets as a transaction.
	muTx sync.Mutex
}

func (r *RateLimit) SetBucketID(routeid string, bucketid string) {
	currentBucketID := r.ids[routeid]

	// when the current Bucket ID is not the same as the new Bucket ID.
	if currentBucketID != bucketid {
		// update the entries for the current Bucket ID.
		if currentBucketID != "" {
			r.entries[currentBucketID]--

			// when the current Bucket ID is no longer referenced by a Route,
			// delete the respective Bucket (and recycle it).
			if r.entries[currentBucketID] <= 0 {
				if currentBucket := r.buckets[currentBucketID]; currentBucket != nil {
					putBucket(currentBucket)
				}
				delete(r.entries, currentBucketID)
				delete(r.buckets, currentBucketID)

				Logger.Info().Timestamp().Str(LogCtxRequest, routeid).Str(LogCtxBucket, currentBucketID).Msg("deleted bucket")
			}
		}

		// set the Route ID to the new Bucket ID.
		r.ids[routeid] = bucketid

		// update the entries for the new Bucket ID.
		r.entries[bucketid]++

		Logger.Info().Timestamp().Str(LogCtxRequest, routeid).Str(LogCtxBucket, bucketid).Msg("set route to bucket")
	}
}

func (r *RateLimit) GetBucketID(routeid string) string {
	return r.ids[routeid]
}

func (r *RateLimit) SetBucketFromID(bucketid string, bucket *Bucket) {
	r.buckets[bucketid] = bucket

	Logger.Info().Timestamp().Str(LogCtxBucket, bucketid).Msgf("set bucket to object %p", bucket)
}

func (r *RateLimit) GetBucketFromID(bucketid string) *Bucket {
	return r.buckets[bucketid]
}

func (r *RateLimit) SetBucket(routeid string, bucket *Bucket) {
	r.buckets[r.ids[routeid]] = bucket
}

func (r *RateLimit) GetBucket(routeid string, resourceid string) *Bucket {
	requestid := routeid + resourceid

	// ID 0 is used as a Global Rate Limit Bucket (or nil).
	if routeid != GlobalRateLimitRouteID {
		switch r.ids[requestid] {
		// when a non-global route is initialized and (BucketID == "NIL"), NO rate limit applies.
		case nilRouteBucket:
			return nil

		// This rate limiter implementation points the Route ID 0 (which is reserved for a
		// Global Rate Limit) to Bucket ID "".
		//
		// As a result (of the Default Bucket mechanism), non-0 Route IDs must be handled accordingly.
		case "":
			// when a non-global route is uninitialized, set it to the Default Bucket.
			//
			// While GetBucket() can be called multiple times BEFORE a request is sent,
			// this case is only true the FIRST time a GetBucket() call is made (for that request),
			// As a result, the Bucket that is allocated from this call will ALWAYS be
			// immediately used.
			//
			// The Route's Bucket ID (Hash) is set to an artificial value while
			// the Route ID is pointed to a Default Bucket. This results in
			// subsequent calls to the Route's Default Bucket to return to
			// the same initialized bucket.
			//
			// When a Default Bucket is exhausted, it will never expire.
			// As a result, the Remaining field-value will remain at 0 until
			// the pending request (and its actual Bucket) is confirmed.
			//
			// requestID = routeid + resourceid
			// temporaryBucketID = requestID
			r.SetBucketID(requestid, requestid)

			// DefaultBucket (Per-Route) = RateLimit.DefaultBucket
			if "" == resourceid {
				if r.DefaultBucket == nil {
					return nil
				}

				b := getBucket()
				b.Limit = r.DefaultBucket.Limit
				b.Remaining = r.DefaultBucket.Limit
				r.SetBucketFromID(requestid, b)

				return b
			}

			// DefaultBucket (Per-Resource) = GetBucket(routeid, "")
			defaultBucket := r.GetBucket(routeid, "")
			if defaultBucket == nil {
				return nil
			}

			b := getBucket()
			b.Limit = defaultBucket.Limit
			b.Remaining = defaultBucket.Limit
			r.SetBucketFromID(requestid, b)

			return b
		}
	}

	return r.buckets[r.ids[requestid]]
}

func (r *RateLimit) SetDefaultBucket(bucket *Bucket) {
	r.DefaultBucket = bucket
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
