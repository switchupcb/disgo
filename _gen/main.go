package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/switchupcb/disgo/_gen/tools"
)

var (
	downloadFlag  = flag.Bool("d", false, "Downloads the latest copy of dasgo from Github.")
	skipEndpoints = flag.Bool("xe", false, "Skips the generation of endpoint functions.")
)

const (
	exeDir = "_gen"

	// redirect `>` is not guaranteed to work, so files must be written.
	filemodewrite = 0644

	downloadURL   = "https://github.com/switchupcb/dasgo/archive/main.zip"
	inputDownload = "input/dasgo.zip"

	outputEndpoints = "../wrapper/endpoints.go"
	outputDasgo     = "../wrapper/dasgo.go"

	endpointErr = "endpoint error: %w\nUse -xe to skip the generation of endpoint functions."
)

func main() {
	if err := check(); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	flag.Parse()

	// download the latest copy of dasgo from Github.
	if *downloadFlag {
		if err := download(downloadURL, inputDownload); err != nil {
			fmt.Printf("%v", err)
			os.Exit(1)
		}
	}

	// dasgo generation
	absfilepath, err := filepath.Abs("input/dasgo-main/dasgo")
	if err != nil {
		fmt.Printf("an error occurred determining the unzipped dasgo source code filepath.\n%v", err)
		os.Exit(1)
	}

	if !*skipEndpoints {
		if err = endpoints(absfilepath); err != nil {
			fmt.Printf("%v", err)
			os.Exit(1)
		}
	}

	if err = convert(absfilepath); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	// disgo generation
	if err = generate(); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}

// check checks that the current
func check() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting the current working directory.\n%w", err)
	}

	if filepath.Base(cwd) != exeDir && filepath.Base(filepath.Dir(cwd)) != "disgo" {
		return fmt.Errorf("This executable must be run from disgo/" + exeDir)
	}

	return nil
}

// download downloads and extracts (.zip) a file from a URL.
func download(url, output string) error {
	curl := exec.Command("curl", "-L", url, "-o", output, "--create-dirs")
	std, err := curl.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error downloading file from url.\n%v", string(std))
	}

	unzip := exec.Command("unzip", "-o", output, "-d", "input")
	std, err = unzip.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error unzipping file.\n%v", string(std))
	}

	return nil
}

// endpoints generates endpoint functions from Dasgo endpoints.
func endpoints(path string) error {
	inputEndpoints := filepath.Join(path, "endpoints.go")
	data, err := os.ReadFile(inputEndpoints)
	if err != nil {
		return fmt.Errorf(endpointErr, err)
	}

	std, err := tools.Endpoints(data)
	if err != nil {
		return fmt.Errorf(endpointErr, err)
	}

	err = os.WriteFile(outputEndpoints, std, filemodewrite)
	if err != nil {
		return fmt.Errorf(endpointErr, err)
	}

	err = os.Remove(inputEndpoints)
	if err != nil {
		return fmt.Errorf(endpointErr, err)
	}

	return nil
}

// convert converts Dasgo objects to the Disgo standard.
func convert(abspath string) error {
	dasgopath := abspath + "/..."

	// nstruct

	// xstruct
	xstruct := exec.Command("tools/xstruct", "-d", dasgopath, "-p", "wrapper", "-g")
	std, err := xstruct.CombinedOutput()
	if err != nil {
		return fmt.Errorf("xstruct error: %v", string(std))
	}

	// snowflake
	std, err = tools.Snowflake(std)
	if err != nil {
		return fmt.Errorf("snowflake error: %w", err)
	}

	err = os.WriteFile(outputDasgo, std, filemodewrite)
	if err != nil {
		return fmt.Errorf("snowflake error: %w", err)
	}

	return nil
}

// generate generates Disgo code.
func generate() error {
	// requests
	if err := os.Chdir("../"); err != nil {
		return fmt.Errorf("chdir error (send): %w", err)
	}

	// send
	sendgen := exec.Command("copygen", "-yml", "wrapper/copygen/requests/setup.yml", "-xm")
	std, err := sendgen.CombinedOutput()
	if err != nil {
		return fmt.Errorf("copygen error (send): %v", string(std))
	}

	// events
	handlegen := exec.Command("copygen", "-yml", "wrapper/copygen/events/setup.yml", "-xm")
	std, err = handlegen.CombinedOutput()
	if err != nil {
		return fmt.Errorf("copygen error (handle): %v", string(std))
	}

	// commands
	commandgen := exec.Command("copygen", "-yml", "wrapper/copygen/commands/setup.yml", "-xm")
	std, err = commandgen.CombinedOutput()
	if err != nil {
		return fmt.Errorf("copygen error (command): %v", string(std))
	}

	// reset
	if err := os.Chdir(exeDir); err != nil {
		return fmt.Errorf("chdir error (generate): %w", err)
	}

	// clean
	data, err := os.ReadFile(outputDasgo)
	if err != nil {
		return fmt.Errorf("clean error: %w", err)
	}

	std, err = tools.Clean(data)
	if err != nil {
		return fmt.Errorf("clean error: %w", err)
	}

	err = os.WriteFile(outputDasgo, std, filemodewrite)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
