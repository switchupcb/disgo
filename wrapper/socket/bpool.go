package socket

import (
	"bytes"
	"sync"
)

var (
	// bpool represents a synchronized bytes.Buffer pool.
	bpool sync.Pool
)

// get gets a buffer from the pool.
func get() *bytes.Buffer {
	if b := bpool.Get(); b != nil {
		return b.(*bytes.Buffer)
	}

	return new(bytes.Buffer)
}

// put puts a buffer into the pool.
func put(b *bytes.Buffer) {
	b.Reset()
	bpool.Put(b)
}
