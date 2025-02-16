package wrapper

import "strings"

type hash func(routeid string, parameters ...string) (string, string)

// Hash returns a hashing function which hashes a request using its routeID.
//
// n represents the degree of resources used to hash the request.
//
//		Per-Route: n = 0
//		Per-Resource: n = 1
//	 ...
func Hash(n int) hash {
	return func(r string, p ...string) (string, string) {
		return r, strings.Join(p[:n], "")
	}
}

var (
	// HashPerRoute hashes a request that uses a per-route rate limit algorithm.
	HashPerRoute = Hash(0)
)
