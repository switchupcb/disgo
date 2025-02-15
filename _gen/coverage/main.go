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

	fmt.Printf("%d endpoints marked unused on purpose.\n\n", len(unused))

	fmt.Printf("Here is the coverage test order of requests.\n")
	for i, endpoint := range filterOutput(unused, findOrder(endpointGraph)) {
		fmt.Printf("%d. %v\n", i, endpoint)
	}
}
