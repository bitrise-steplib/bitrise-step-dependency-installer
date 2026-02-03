package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/bitrise-step-dependency-installer/installer"
)

func main() {
	toolVersionFileListStr := os.Getenv("tool_version_file_list")

	if toolVersionFileListStr == "" {
		fmt.Println("No tool version files provided in 'tool_version_file_list' environment variable.")
		os.Exit(1)
	}

	workflowListStr := os.Getenv("bitrise_workflow_id_list")
	verboseMode := os.Getenv("verbose") == "true"
	outputFormat := strings.ToLower(os.Getenv("output_format"))

	if outputFormat == "" {
		outputFormat = "plaintext"
	}
	if outputFormat != "plaintext" && outputFormat != "json" && outputFormat != "bash" {
		fmt.Printf("Invalid output format: %s. Supported formats are: plaintext, json, bash.\n", outputFormat)
		os.Exit(1)
	}

	toolVersionFileList := strings.Split(toolVersionFileListStr, "\n")
	validWorkflows := validWorkflows(strings.Split(workflowListStr, "\n"))
	installer := installer.NewInstallerManager(&installer.DefaultInstallerFactory{})
	results, err := installer.InstallAll(toolVersionFileList, validWorkflows, outputFormat, verboseMode)

	if err != nil {
		fmt.Printf("Error during installation: %s\n", err)
		os.Exit(1)
	}

	wasError := false
	for _, result := range results {
		if result.Err != nil {
			fmt.Printf("Error during installation for %s: %s\n", result.ToolInstallerDescription, result.Err)
			wasError = true
		} else if verboseMode {
			fmt.Printf("Installation output for %s: %s\n", result.ToolInstallerDescription, string(result.Output))
		}
	}

	if wasError {
		os.Exit(1)
	}
	os.Exit(0)
}

func validWorkflows(workflowList []string) []string {
	validList := []string{}
	for _, workflow := range workflowList {
		trimmed := strings.TrimSpace(workflow)
		if trimmed != "" {
			validList = append(validList, trimmed)
		}
	}
	return validList
}
