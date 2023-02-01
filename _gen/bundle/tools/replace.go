package tools

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	replacedDir = "_gen/bundle/tools/replaced"
	boundary    = "\n---\n"
)

var (
	bundleWarning   = `// Code generated by golang.org/x/tools/cmd/bundle. DO NOT EDIT.`
	generateWarning = `//go:generate bundle -o disgo.go -dst . -pkg disgo -prefix "" ./wrapper`
	bundlerWarning  = []byte(`// Code generated by github.com/switchupcb/disgo/_gen/tools/bundle. DO NOT EDIT.` + "\n")
)

// Replace replaces fields that have comments removed after being field aligned.
func Replace(filepath string) error {
	// replaced represents a map of fieldaligned structs (with comments removed)
	// to fieldaligned structs (with comments added back).
	replaced, err := initialize(replacedDir, boundary)
	if err != nil {
		return fmt.Errorf("initialize: %w", err)
	}

	// load the bundle's content.
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("error reading generated %v file: %w", filepath, err)
	}

	// strip the bundle's lines (in order to perform a comparison).
	content := strip(data)

	// replace the fieldaligned structs (without comments).
	for fieldaligned, fixed := range replaced {
		content = strings.Replace(content, fieldaligned, fixed, 1)
	}

	content = strings.Replace(content, bundleWarning, "", 1)
	content = strings.Replace(content, generateWarning, "", 1)

	// gofmt
	contentdata := []byte(content)
	fmtdata, err := format.Source(contentdata)
	if err != nil {
		if err = os.WriteFile(filepath, contentdata, filemodewrite); err != nil {
			return fmt.Errorf("error writing file %w", err)
		}

		return fmt.Errorf("error formatting generated code: %w", err)
	}

	var output []byte
	output = append(output, bundlerWarning...)
	output = append(output, fmtdata...)

	if err = os.WriteFile(filepath, output, filemodewrite); err != nil {
		return fmt.Errorf("error writing file %w", err)
	}

	return nil
}

// initialize initializes the replaced map by reading the files in the given directory.
//
// A boundary is used to determine the key and field.
func initialize(dir, boundary string) (map[string]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	output := make(map[string]string, len(files))
	for _, file := range files {
		filepath := path.Join(dir, file.Name())
		data, err := os.ReadFile(filepath)
		if err != nil {
			return nil, fmt.Errorf("error reading replaced %v file: %w", filepath, err)
		}

		// strip the file's lines (in order to perform a comparison).
		content := strip(data)

		fields := strings.Split(content, boundary)
		output[fields[0]] = fields[1]
	}

	return output, nil
}

var (
	newline = byte('\n')
)

// strip strips the lines of the given data.
func strip(data []byte) string {
	var stripped strings.Builder

	var line []byte
	for _, b := range data {
		if b == newline {
			stripped.WriteString(strings.TrimSpace(string(line)))
			stripped.WriteByte(newline)

			line = []byte{}

			continue
		}

		line = append(line, b)
	}

	stripped.WriteString(strings.TrimSpace(string(line)))

	return stripped.String()
}