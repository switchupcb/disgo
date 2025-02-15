package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	// integrationTestFilepath represents the integrationTestFilepath from ./coverage
	integrationTestFilepath = "../../wrapper/tests/integration/coverage_test.go"
)

// checkCoverageTest checks the integration coverage test for the given endpoints
// and outputs the missing requests.
func checkCoverageTest(endpoints []string) ([]string, error) {
	// Read the integration test file.
	fileBytes, err := os.ReadFile(integrationTestFilepath)
	if err != nil {
		return nil, fmt.Errorf("coverage test check: %w", err)
	}
	file := string(fileBytes)

	// Check the file for each request struct instantiation used to send an event (e.g., &Request, new(Request))
	uncalled := make([]string, 0, len(endpoints))
	for i := range endpoints {
		if !strings.Contains(file, "&"+endpoints[i]) && !strings.Contains(file, "new("+endpoints[i]+")") {
			uncalled = append(uncalled, endpoints[i])
		}
	}

	return uncalled, nil
}
