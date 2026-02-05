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
	workflowArg := ""
	if workflow == "" {
		if verboseMode {
			fmt.Printf("Installing tools from %s global tools block\n", toolVersionFile)
		}
	} else {
		workflowArg = "--workflow " + workflow
		if verboseMode {
			fmt.Printf("Installing tools from %s workflow %s\n", toolVersionFile, workflow)
		}
	}

	command := commandExecutor("bitrise", "tools", "setup", "--config", toolVersionFile, workflowArg)
	return runCommand(verboseMode, command)
}

func runCommand(verboseMode bool, cmd *exec.Cmd) ([]byte, error) {
	if verboseMode {
		fmt.Printf("Running command: %s %s\n", cmd.Path, strings.Join(cmd.Args, " "))
	}
	return cmd.CombinedOutput()
}
