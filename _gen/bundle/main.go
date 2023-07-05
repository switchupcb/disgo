package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/switchupcb/disgo/_gen/bundle/tools"
)

var (
	showFieldaligned = flag.Bool("f", false, "Use -f to showcase the output from the fieldalignment tool.")
)

const (
	exeDir        = "_gen/bundle"
	bundlePath    = "disgo.go"
	pkg           = "package disgo"
	filemodewrite = 0644
)

func main() {
	if err := check(); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	flag.Parse()

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
		return fmt.Errorf("error getting the current working directory: %w", err)
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

	// clear the bundled file.
	bundle := `//go:generate bundle -o disgo.go -dst . -pkg disgo -prefix "" ./wrapper`
	cleared := strings.Join([]string{bundle, pkg}, "\n")
	if err := os.WriteFile(bundlePath, []byte(cleared), filemodewrite); err != nil {
		return fmt.Errorf("clear: %w", err)
	}

	bundlegen := exec.Command("go", "generate")
	std, err := bundlegen.CombinedOutput()
	if err != nil {
		return fmt.Errorf("bundle (go generate): %v", string(std))
	}

	// fix the imports of the bundle.
	if err := tools.Imports(bundlePath); err != nil {
		return fmt.Errorf("imports: %w", err)
	}

	// fieldalign the bundle (until there is nothing else to fieldalign).
	var fieldalignmentOutput []byte
	for i := 0; i < 5; i++ {
		fieldalignment := exec.Command("fieldalignment", "-fix", bundlePath)
		std, err = fieldalignment.CombinedOutput()
		if err != nil && err.Error() == "exit status 3" {
			fieldalignmentOutput = append(fieldalignmentOutput, std...)
			fmt.Printf("WARNING (fieldalignment): %v\n", err)
			fmt.Println("running fieldalignment again...\n")

			continue

		} else if err != nil {
			return fmt.Errorf("fieldalignment: %w", err)
		}

		fieldalignmentOutput = append(fieldalignmentOutput, std...)

		if i == 4 {
			return fmt.Errorf("fieldalignment: more calls to fieldalignment were required")
		}

		break
	}

	if *showFieldaligned {
		fmt.Println(string(fieldalignmentOutput))
	}

	// add removed comments to the bundle.
	if err := tools.Replace(bundlePath); err != nil {
		return fmt.Errorf("replace: %w", err)
	}

	return nil
}
