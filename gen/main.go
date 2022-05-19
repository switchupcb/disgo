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
	downloadFlag  = flag.Bool("d", false, "Downloads the latest copy of dasgo from Github.")
	skipEndpoints = flag.Bool("xe", false, "Skips the generation of endpoint functions.")
)

const (
	// redirect `>` is not guaranteed to work, so files must be written.
	filemodewrite = 0644

	downloadURL   = "https://github.com/switchupcb/dasgo/archive/main.zip"
	inputDownload = "input/dasgo.zip"

	outputEndpoints = "../wrapper/endpoints.go"
	outputDasgo     = "../wrapper/dasgo.go"
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
		os.Exit(2)
	}

	if !*skipEndpoints {
		inputEndpoints := filepath.Join(absfilepath, "endpoints.go")
		data, err := os.ReadFile(inputEndpoints)
		if err != nil {
			fmt.Printf("endpoint error: %v", err)
			os.Exit(3)
		}

		std, err := tools.Endpoints(data)
		if err != nil {
			fmt.Printf("endpoint error: %v", err)
			os.Exit(3)
		}

		err = os.WriteFile(outputEndpoints, std, filemodewrite)
		if err != nil {
			fmt.Printf("endpoint error: %v", err)
			os.Exit(3)
		}

		e := os.Remove(inputEndpoints)
		if e != nil {
			fmt.Printf("endpoint error: %v", err)
			os.Exit(3)
		}
	}

	if err = convert(absfilepath); err != nil {
		fmt.Printf("%v", err)
		os.Exit(4)
	}

	// disgo generation
	if err = generate(); err != nil {
		fmt.Printf("%v", err)
		os.Exit(5)
	}
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
		return fmt.Errorf("snowflake error: %v", err)
	}

	err = os.WriteFile(outputDasgo, std, filemodewrite)
	if err != nil {
		return fmt.Errorf("snowflake error: %v", err)
	}

	return nil
}

// generate generates Disgo code.
func generate() error {
	// requests
	if err := os.Chdir("../"); err != nil {
		return fmt.Errorf("chdir error (send): %v", err)
	}

	sendgen := exec.Command("copygen", "-yml", "wrapper/copygen/setup/setup.yml", "-xm")
	std, err := sendgen.CombinedOutput()
	if err != nil {
		return fmt.Errorf("copygen error (send): %v", string(std))
	}

	if err := os.Chdir("gen"); err != nil {
		return fmt.Errorf("chdir error (send): %v", err)
	}

	// events

	// clean
	data, err := os.ReadFile(outputDasgo)
	if err != nil {
		return fmt.Errorf("clean error: %v", err)
	}

	std, err = tools.Clean(data)
	if err != nil {
		return fmt.Errorf("clean error: %v", err)
	}

	err = os.WriteFile(outputDasgo, std, filemodewrite)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
