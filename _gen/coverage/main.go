package main

import (
	"fmt"
	"maps"
)

func main() {
	// unreferenced represents a map of endpoints that are not referenced in the maps "unused" or "endpointsGraph".
	unreferenced := filterEndpoints(endpoints, maps.Keys(endpointGraph), maps.Keys(unused))
	fmt.Printf("%d endpoints unreferenced.\n", len(unreferenced))
	for endpoint := range unreferenced {
		fmt.Printf("%v\n", endpoint)
	}
	fmt.Println()

	// unrequired represents a map of endpoints that are referenced in the maps "unused" or "endpointsGraph",
	// but not in the map "endpoints" (containing all endpoints).
	unrequired := filterEndpoints(combineEndpoints(maps.Keys(endpointGraph), maps.Keys(unused)), maps.Keys(endpoints))
	fmt.Printf("%d endpoints unrequired.\n", len(unrequired))
	for endpoint := range unrequired {
		fmt.Printf("%v\n", endpoint)
	}
	fmt.Println()

	// unused endpoints are not used in the coverage test.
	fmt.Printf("%d endpoints marked unused on purpose.\n\n", len(unused))

	// orderOfRequests represents the order requests should be called in the coverage test.
	orderOfRequests := filterOutput(unused, findOrder(endpointGraph))
	fmt.Printf("Here is the coverage test order of requests.\n\n")
	for i, endpoint := range orderOfRequests {
		fmt.Printf("%d. %v\n", i, endpoint)
	}
	fmt.Println()

	// uncalledRequests represents the requests that should be called in the coverage test, but aren't.
	uncalledRequests, err := checkCoverageTest(orderOfRequests)
	fmt.Printf("Here are requests which are not present in the current coverage test, but should be.\n\n")
	if err != nil {
		fmt.Printf("%q", err)

		return
	}

	for i, endpoint := range uncalledRequests {
		fmt.Printf("%d. %v\n", i, endpoint)
	}
}
