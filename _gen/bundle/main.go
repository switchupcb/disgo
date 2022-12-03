package main

import (
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	exeDir     = "_gen/bundle"
	bundlePath = "disgo.go"
)

func main() {
	if err := check(); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	// disgo generation
	if err := generate(); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}

// check checks that the current working directory is `disgo/_gen/bundle`.
func check() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting the current working directory.\n%w", err)
	}

	cwdDir := filepath.Dir(cwd)
	base := filepath.Base(cwdDir) + "/" + filepath.Base(cwd)
	if base != exeDir && filepath.Base(filepath.Dir(cwdDir)) != "disgo" {
		return fmt.Errorf("This executable must be run from disgo/" + exeDir)
	}

	return nil
}

// generate generates the disgo bundle.
func generate() error {
	if err := os.Chdir("../../"); err != nil {
		return fmt.Errorf("chdir: %w", err)
	}

	bundle := exec.Command("bundle",
		"-o", bundlePath,
		"-dst", ".",
		"-pkg", "disgo",
		"-prefix", "\"\"",
		"./wrapper",
	)

	// will return format error on success.
	std, _ := bundle.CombinedOutput()
	fmt.Printf("WARNING: %v", string(std))

	if err := imports(); err != nil {
		return fmt.Errorf("imports: %w", err)
	}

	return nil
}

const (
	filemodewrite = 0644
)

var (
	skip = map[string]bool{
		`"encoding/json"`: true,
	}
)

// imports fixes the imports of the bundler.
func imports() error {
	data, err := os.ReadFile(bundlePath)
	if err != nil {
		return fmt.Errorf("error reading generated %v file: %v", bundlePath, err)
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

	if err = os.WriteFile(bundlePath, fmtdata, filemodewrite); err != nil {
		return fmt.Errorf("error while writing file", err)
	}

	return nil
}
