package main

import "fmt"

// filterOutput removes unused endpoints from the endpoint output slice.
func filterOutput(unused map[string]bool, endpoints []string) []string {
	output := make([]string, 0, len(unused))

	for _, endpoint := range endpoints {
		if unused[endpoint] {
			continue
		}

		output = append(output, endpoint)
	}

	return output
}

// findOrder finds the optimal order of endpoints using a dependency graph.
func findOrder(endpoints map[string][]string) []string {
	numEndpoints := len(endpoints)

	// dependents represents a map of dependent endpoints to
	// the respective amount of dependencies (map[dependency]numDependencies).
	dependents := make(map[string]int, numEndpoints)

	// calculate the number of dependencies for each dependent.
	for endpoint, dependencies := range endpoints {
		dependents[endpoint] = len(dependencies)
	}

	// queue represents a first-in first-out data structure.
	//
	// queue can't have more entries than the number of endpoints,
	// so initialize a map of length 0 with capacity = numEndpoints.
	queue := make([]string, 0, numEndpoints)

	// fill the queue with dependent endpoints that have no dependencies (i.e `CreateGuild`).
	for endpoint, numDependencies := range dependents {
		if numDependencies == 0 {
			queue = append(queue, endpoint)
		}
	}

	// output represents the returned result.
	output := make([]string, 0, numEndpoints)

	// add dependent endpoints with no dependencies to the output list.
	for len(queue) > 0 {
		// select the first entry in the queue (of dependent endpoints with no dependencies).
		current := queue[0]

		// remove the first entry out the queue (of dependent endpoints with no dependencies).
		queue = queue[1:]

		// add the entry to the output.
		output = append(output, current)

		// decrease the amount of endpoints remaining.
		delete(endpoints, current)
		numEndpoints--

		// The operation above removes any amount of dependencies from the queue.
		//
		// add endpoints with no dependencies to the queue.
		for endpoint, dependencies := range endpoints {
			// when a dependency is added to the output,
			// remove that endpoint (current) from dependent (endpoint)s' dependencies.
			if contains(dependencies, current) {
				dependents[endpoint]--

				// when the endpoint is dependent on no other endpoints, add it to the queue.
				if dependents[endpoint] == 0 {
					queue = append(queue, endpoint)
				}
			}
		}
	}

	if numEndpoints != 0 {
		fmt.Printf("WARNING: dependency cycle occurred (i.e [a: b],[b: a]) or necessary endpoint is unused.\n\n")
		fmt.Println("Examine the following endpoints.")

		for endpoint, dependencies := range endpoints {
			fmt.Println(endpoint, "depends on", dependencies)
		}

		return []string{}
	}

	return output
}

// contains returns whether the slice s contains the string x.
func contains(s []string, x string) bool {
	for _, item := range s {
		if x == item {
			return true
		}
	}

	return false
}
