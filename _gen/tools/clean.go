package tools

import (
	"fmt"
	"go/format"
	"strings"
)

// Clean cleans the generated code.
func Clean(code []byte) ([]byte, error) {
	var file strings.Builder

	// traverse every line in the dasgo.go file.
	lines := strings.Split(string(code), "\n")
	for _, line := range lines {
		if fixed := fixPointerFunc(line); fixed != "" {
			line = fixed
		}

		keep := removeApplicationID(line) && removePayloadJSON(line)
		if keep {
			file.WriteString(line + "\n")
		}
	}

	// gofmt
	content := []byte(file.String())
	fmtdata, err := format.Source(content)
	if err != nil {
		return content, fmt.Errorf("an error occurred while formatting the generated code.\n%w", err)
	}

	return fmtdata, nil
}

// removeApplicationID removes an ApplicationID field from dasgo.
func removeApplicationID(line string) bool {
	keep := true

	// determine if the line contains `ApplicationID` (with no other characters before/after).
	linefields := strings.Fields(line)
	for i, field := range linefields {
		// ApplicationID fields with json tags are kept, so check for a json tag.
		if field == "ApplicationID" {
			keep = false

			n := len(linefields)
			for j := i + 1; j < n; j++ {
				if strings.Contains(linefields[j], "json") {
					keep = true
				}
			}

			break
		}
	}

	return keep
}

// removePayloadJSON removes a PayloadJSON field from dasgo.
func removePayloadJSON(line string) bool {
	keep := true

	// determine if the line contains `PayloadJSON` (with no other characters before/after).
	linefields := strings.Fields(line)
	for _, field := range linefields {
		if field == "PayloadJSON" {
			keep = false
		}
	}

	return keep
}

// fixPointerFunc fixes the generated pointer function in dasgo.
func fixPointerFunc(line string) string {
	if "func Pointer(v T) *T {" == line {
		return "func Pointer[T any](v T) *T {"
	}

	return ""
}
