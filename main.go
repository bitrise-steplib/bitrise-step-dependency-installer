package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bitrise-io/bitrise-step-dependency-installer/installer"
)

func main() {
	toolVersionFile := os.Getenv("tool_version_file")
	workflow := os.Getenv("bitrise_workflow_id")
	verboseMode := os.Getenv("verbose") == "true"

	if toolVersionFile == "" {
		fmt.Println("No tool version file provided in 'tool_version_file' environment variable.")
		os.Exit(1)
	} else if len(strings.Split(toolVersionFile, " ")) > 1 {
		fmt.Println("Multiple tool version files are not supported.")
		os.Exit(1)
	} else if len(strings.Split(workflow, " ")) > 1 {
		fmt.Println("Multiple workflows are not supported.")
		os.Exit(1)
	}

	output, err := installer.Install(toolVersionFile, workflow, verboseMode, exec.Command)

	if err != nil {
		fmt.Printf("Installation output: %s\n", string(output))
		fmt.Printf("Error during installation: %s\n", err)
		os.Exit(1)
	} else if verboseMode {
		fmt.Printf("Installation output: %s\n", string(output))
	}

	os.Exit(0)
}
