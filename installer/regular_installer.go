package installer

import (
	"fmt"
	"os/exec"
)

type RegularToolInstaller struct {
	filePath     string
	outputFormat string
	verboseMode  bool
	ExecCommand  ExecCommandFunc // Exported for testing
}

// NewRegularToolInstaller creates a new RegularToolInstaller with default command execution
func NewRegularToolInstaller(filePath, outputFormat string, verboseMode bool) *RegularToolInstaller {
	return &RegularToolInstaller{
		filePath:     filePath,
		outputFormat: outputFormat,
		verboseMode:  verboseMode,
		ExecCommand:  exec.Command,
	}
}

func (r *RegularToolInstaller) Install() ([]byte, error) {
	if r.verboseMode {
		fmt.Printf("Installing tools from: %s\n", r.filePath)
	}

	return r.ExecCommand("bitrise", "tools", "setup", "--config", r.filePath, "--output-format", r.outputFormat).CombinedOutput()
}

func (r *RegularToolInstaller) Description() string {
	return fmt.Sprintf("RegularToolInstaller for file: %s", r.filePath)
}
