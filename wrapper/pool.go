package wrapper

import "sync"

var (
	// gpool represents a synchronized Gateway Payload pool.
	gpool sync.Pool
)

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
