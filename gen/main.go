package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/switchupcb/disgo/gen/tools"
)

var (
	downloadFlag = flag.Bool("d", false, "Downloads the latest copy of dasgo from Github.")
)

const (
	// redirect `>` is not guaranteed to work, so files must be written.
	filemodewrite = 0644

	outputEndpoints = "../wrapper/endpoints.go"
	outputDasgo     = "../wrapper/dasgo.go"
	outputSend      = "../wrapper/send.go"
)

func main() {
	if err := check(); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	flag.Parse()

	// download the latest copy of dasgo from Github.
	if *downloadFlag {
		if err := download("https://github.com/switchupcb/dasgo/archive/main.zip", "dasgo.zip"); err != nil {
			fmt.Printf("%v", err)
			os.Exit(1)
		}
	}

	// dasgo generation
	absfilepath, err := filepath.Abs("dasgo-main/dasgo")
	if err != nil {
		fmt.Printf("an error occurred determining the unzipped dasgo source code filepath.\n%v", err)
		os.Exit(1)
	}

	if err = convert(absfilepath); err != nil {
		fmt.Printf("%v", err)
		os.Exit(2)
	}

	// disgo generation
	// generate()
}

// check checks that the current
func check() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting the current working directory.\n%w", err)
	}

	if filepath.Base(cwd) != "gen" && filepath.Base(filepath.Dir(cwd)) != "disgo" {
		return fmt.Errorf("This executable must be run from disgo/gen")
	}

	return nil
}

// download downloads and extracts (.zip) a file from a URL.
func download(url, output string) error {
	curl := exec.Command("curl", "-L", url, "-o", output)
	std, err := curl.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error downloading file from url.\n%v", string(std))
	}

	unzip := exec.Command("unzip", output)
	std, err = unzip.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error unzipping file.\n%v", string(std))
	}

	return nil
}

// convert converts Dasgo objects to the Disgo standard.
func convert(abspath string) error {
	// endpoints
	inputEndpoints := filepath.Join(abspath, "endpoints.go")
	std, err := tools.Endpoints(inputEndpoints)
	if err != nil {
		return fmt.Errorf("endpoint error: %v", err)
	}

	err = os.WriteFile(outputEndpoints, std, filemodewrite)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	e := os.Remove(inputEndpoints)
	if e != nil {
		return fmt.Errorf("%w", err)
	}

	// nstruct

	// xstruct
	xstruct := exec.Command("tools/xstruct", "-d", abspath+"/...", "-p", "disgo", "-g")
	std, err = xstruct.CombinedOutput()
	if err != nil {
		return fmt.Errorf("xstruct error: %v", string(std))
	}

	err = os.WriteFile(outputDasgo, std, filemodewrite)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// flagstd

	return nil
}
