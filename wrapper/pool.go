package wrapper

import (
	"sync"
	"time"
)

// rlbpool represents a synchronized Rate Limit Bucket pool.
var rlbpool sync.Pool

// getBucket gets a Bucket from a pool.
func getBucket() *Bucket {
	if b := rlbpool.Get(); b != nil {
		return b.(*Bucket) //nolint:forcetypeassert
	}

	return new(Bucket)
}

// putBucket puts a Rate Limit Bucket into the pool.
func putBucket(b *Bucket) {
	b.ID = ""
	b.Limit = 0
	b.Remaining = 0
	b.Pending = 0
	b.Date = time.Time{}
	b.Expiry = time.Time{}

	rlbpool.Put(b)
}

// spool represents a synchronized Session pool.
var spool sync.Pool

// NewSession gets a Session from a pool.
func NewSession() *Session {
	if s := spool.Get(); s != nil {
		return s.(*Session) //nolint:forcetypeassert
	}

	return new(Session)
}

// putSession puts a Session into the pool.
func putSession(s *Session) {
	s.Lock()
	defer s.Unlock()

	// reset the Session.
	s.ID = ""
	s.Seq = 0
	s.Endpoint = ""
	s.Shard = nil
	s.Context = nil
	s.cancel = nil
	s.Conn = nil
	s.setState(SessionStateNew)
	s.heartbeat = nil
	s.manager = nil
	s.client_manager = nil
	s.RateLimiter = nil

	spool.Put(s)
}

// gpool represents a synchronized Gateway Payload pool.
var gpool sync.Pool

// getPayload gets a Gateway Payload from the pool.
func getPayload() *GatewayPayload {
	if g := gpool.Get(); g != nil {
		return g.(*GatewayPayload) //nolint:forcetypeassert
	}

	return new(GatewayPayload)
}

// putPayload puts a Gateway Payload into the pool.
func putPayload(g *GatewayPayload) {
	// reset the Gateway Payload.
	g.Op = 0
	g.Data = nil
	g.SequenceNumber = nil
	g.EventName = nil

	gpool.Put(g)
}

// vspool represents a synchronized Voice Session pool.
var vspool sync.Pool

// newVoiceSession gets a Voice Session from a pool.
func newVoiceSession() *VoiceSession {
	if vs := vspool.Get(); vs != nil {
		return vs.(*VoiceSession) //nolint:forcetypeassert
	}

	return new(VoiceSession)
}

// putVoiceSession puts a Voice Session into the pool.
func putVoiceSession(vs *VoiceSession) {
	vs.Lock()
	defer vs.Unlock()

	// reset the Session.
	vs.ID = ""
	vs.Nonce = 0
	vs.Context = nil
	vs.Conn = nil
	vs.heartbeat = nil
	vs.manager = nil

	vspool.Put(vs)
}

// vpool represents a synchronized Voice Payload pool.
var vpool sync.Pool

// getVoicePayload gets a Voice Payload from the pool.
func getVoicePayload() *VoicePayload {
	if v := vpool.Get(); v != nil {
		return v.(*VoicePayload) //nolint:forcetypeassert
	}

	return new(VoicePayload)
}

// putVoicePayload puts a Voice Payload into the pool.
func putVoicePayload(v *VoicePayload) {
	// reset the Voice Payload.
	v.Op = 0
	v.Data = nil

	vpool.Put(v)
}
