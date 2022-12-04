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
		return content, fmt.Errorf("clean: error formatting generated code: %w", err)
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
				if linefields[j] != "`json:\"-\"`" && strings.Contains(linefields[j], "json") {
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

var (
	lines = map[string]string{
		"func Pointer(v T) *T {":                       "func Pointer[T any](v T) *T {",
		"func Pointer2(v T, null ...bool) **T {":       "func Pointer2[T any](v T, null ...bool) **T {",
		"func IsValue(p *T) bool {":                    "func IsValue[T any](p *T) bool {",
		"func IsValue2(dp **T) bool {":                 "func IsValue2[T any](dp **T) bool {",
		"func PointerCheck(dp **T) PointerIndicator {": "func PointerCheck[T any](dp **T) PointerIndicator {",
	}
)

// fixPointerFunc fixes the generated pointer functions in dasgo.
func fixPointerFunc(line string) string {
	if fixed, ok := lines[line]; ok {
		return fixed
	}

	return ""
}
