package installer

import (
	"fmt"
	"os/exec"
	"strings"
)

type ExecCommandFunc func(name string, args ...string) *exec.Cmd

const (
	BitriseYMLFileName = "bitrise.yml"
)

func Install(toolVersionFile string, workflow string, verboseMode bool, commandExecutor ExecCommandFunc) ([]byte, error) {

	if strings.HasSuffix(toolVersionFile, BitriseYMLFileName) {
		return installBitriseYML(toolVersionFile, workflow, verboseMode, commandExecutor)
	} else {
		return installRegularToolFile(toolVersionFile, verboseMode, commandExecutor)
	}
}

func installRegularToolFile(toolVersionFile string, verboseMode bool, commandExecutor ExecCommandFunc) ([]byte, error) {
	if verboseMode {
		fmt.Printf("Installing tools from %s\n", toolVersionFile)
	}

	command := commandExecutor("bitrise", "tools", "setup", "--config", toolVersionFile)
	return runCommand(verboseMode, command)
}

func installBitriseYML(toolVersionFile string, workflow string, verboseMode bool, commandExecutor ExecCommandFunc) ([]byte, error) {
	// Validate bitrise.yml before installing
	if err := validateBitriseYML(toolVersionFile, verboseMode, commandExecutor); err != nil {
		return nil, err
	}

	cmdArgs := []string{"tools", "setup", "--config", toolVersionFile}
	if workflow == "" {
		if verboseMode {
			fmt.Printf("Installing tools from %s global tools block\n", toolVersionFile)
		}
	} else {
		cmdArgs = append(cmdArgs, "--workflow", workflow)
		if verboseMode {
			fmt.Printf("Installing tools from %s workflow %s\n", toolVersionFile, workflow)
		}
	}

	command := commandExecutor("bitrise", cmdArgs...)
	return runCommand(verboseMode, command)
}

func validateBitriseYML(toolVersionFile string, verboseMode bool, commandExecutor ExecCommandFunc) error {
	if verboseMode {
		fmt.Printf("Validating %s\n", toolVersionFile)
	}

	command := commandExecutor("bitrise", "validate", "--config", toolVersionFile)
	output, err := command.CombinedOutput()

	if err != nil {
		return fmt.Errorf("bitrise.yml validation failed: %s\n%s", err, string(output))
	}

	if verboseMode {
		fmt.Printf("Validation passed for %s\n", toolVersionFile)
	}

	return nil
}

func runCommand(verboseMode bool, cmd *exec.Cmd) ([]byte, error) {
	if verboseMode {
		fmt.Printf("Running command: %s %s\n", cmd.Path, strings.Join(cmd.Args[1:], " "))
	}
	return cmd.CombinedOutput()
}
