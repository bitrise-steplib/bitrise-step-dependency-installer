package installer

import (
	"fmt"
	"os/exec"
)

type BitriseToolInstaller struct {
	filePath     string
	workflow     string
	outputFormat string
	verboseMode  bool
	ExecCommand  ExecCommandFunc
}

func NewBitriseToolInstaller(filePath, workflow, outputFormat string, verboseMode bool) *BitriseToolInstaller {
	return &BitriseToolInstaller{
		filePath:     filePath,
		workflow:     workflow,
		outputFormat: outputFormat,
		verboseMode:  verboseMode,
		ExecCommand:  exec.Command,
	}
}

func (b *BitriseToolInstaller) Install() ([]byte, error) {
	if b.workflow == "" {
		if b.verboseMode {
			fmt.Printf("Installing tools from %s global tools block\n", b.filePath)
		}
		return b.ExecCommand("bitrise", "tools", "setup", "--config", b.filePath, "--output-format", b.outputFormat).CombinedOutput()
	}
	if b.verboseMode {
		fmt.Printf("Installing tools from %s workflow %s\n", b.filePath, b.workflow)
	}
	return b.ExecCommand("bitrise", "tools", "setup", "--config", b.filePath, "--workflow", b.workflow, "--output-format", b.outputFormat).CombinedOutput()
}

func (b *BitriseToolInstaller) Description() string {
	if b.workflow != "" {
		return fmt.Sprintf("BitriseToolInstaller for file: %s and workflow: %s", b.filePath, b.workflow)
	}
	return fmt.Sprintf("BitriseToolInstaller for file: %s (no workflow)", b.filePath)
}
