package main

import (
	"iter"
	"maps"
)

// filterEndpoints removes unused endpoints from an endpoint map.
func filterEndpoints(endpoints map[string]bool, unused ...iter.Seq[string]) map[string]bool {
	e := maps.Clone(endpoints)
	for _, unusedEndpoints := range unused {
		for unusedEndpoint := range unusedEndpoints {
			delete(e, unusedEndpoint)
		}
	}

	return e
}

// combineEndpoints combines endpoints into a map.
func combineEndpoints(endpointsIter ...iter.Seq[string]) map[string]bool {
	e := make(map[string]bool)
	for _, endpoints := range endpointsIter {
		for endpoint := range endpoints {
			e[endpoint] = true
		}
	}

	return e
}
