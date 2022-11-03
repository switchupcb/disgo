package tools

import (
	"fmt"
	"go/format"
	"strings"
)

// Clean cleans the generated code.
func Clean(code []byte) ([]byte, error) {
	content := removeApplicationID(string(code))
	content = fixPointerFunc(content)

	// gofmt
	contentdata := []byte(content)
	fmtdata, err := format.Source(contentdata)
	if err != nil {
		return contentdata, fmt.Errorf("an error occurred while formatting the generated code.\n%w", err)
	}

	return fmtdata, nil
}

// removeApplicationID removes all ApplicationID fields from dasgo.
func removeApplicationID(content string) string {
	var file strings.Builder

	lines := strings.Split(content, "\n")
	for _, line := range lines {
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

		if keep {
			file.WriteString(line + "\n")
		}
	}

	return file.String()
}

// fixPointerFunc fixes the generated pointer function in dasgo.
func fixPointerFunc(content string) string {
	var file strings.Builder

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if "func Pointer(v T) *T {" == line {
			file.WriteString("func Pointer[T any](v T) *T {" + "\n")
		} else {
			file.WriteString(line + "\n")
		}
	}

	return file.String()
}
