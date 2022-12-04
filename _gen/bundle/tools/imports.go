package tools

import (
	"fmt"
	"go/format"
	"os"
	"strings"
)

const (
	filemodewrite = 0644
)

var (
	// skip represents a set of imports to remove.
	skip = map[string]bool{
		`"encoding/json"`: true,
	}
)

// Imports fixes the imports of the bundler.
func Imports(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("error reading generated %v file: %v", filepath, err)
	}

	var output []byte

	content := string(data)
	for _, line := range strings.Split(content, "\n") {
		if skip[strings.TrimSpace(line)] {
			continue
		}

		output = append(output, []byte(line+"\n")...)
	}

	// gofmt
	fmtdata, err := format.Source(output)
	if err != nil {
		return fmt.Errorf("error while formatting the generated code.\n%w", err)
	}

	if err = os.WriteFile(filepath, fmtdata, filemodewrite); err != nil {
		return fmt.Errorf("error while writing file", err)
	}

	return nil
}
