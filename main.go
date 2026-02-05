package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bitrise-io/bitrise-step-dependency-installer/installer"
)

func main() {
	toolVersionFile := os.Getenv("tool_version_file")

	if toolVersionFile == "" {
		fmt.Println("No tool version file provided in 'tool_version_file' environment variable.")
		os.Exit(1)
	}

	workflow := os.Getenv("bitrise_workflow_id")
	verboseMode := os.Getenv("verbose") == "true"

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
