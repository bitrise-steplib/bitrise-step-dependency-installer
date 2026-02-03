package installer

import (
	"fmt"
	"os/exec"
	"strings"
)

type ExecCommandFunc func(name string, args ...string) *exec.Cmd

type InstallerResult struct {
	ToolInstallerDescription string
	Output                   []byte
	Err                      error
}

type ToolInstaller interface {
	Install() ([]byte, error)
	Description() string
}

const (
	BitriseYMLFileName = "bitrise.yml"
)

type InstallerManager struct {
	installerFactory InstallerFactory
}

func NewInstallerManager(factory InstallerFactory) *InstallerManager {
	return &InstallerManager{
		installerFactory: factory,
	}
}

func (m *InstallerManager) InstallAll(toolVersionFileList []string, workflowList []string, outputFormat string, verboseMode bool) ([]InstallerResult, error) {
	var toolInstallers []ToolInstaller
	bitriseCount := 0

	for _, filePath := range toolVersionFileList {
		if filePath == "" {
			continue
		}

		switch {
		case strings.HasSuffix(filePath, BitriseYMLFileName):
			bitriseCount++
			if bitriseCount > 1 {
				return nil, fmt.Errorf("multiple bitrise.yml files detected in the input list, only one is allowed")
			}
			bitriseInstallers := createBitriseInstallers(m.installerFactory, filePath, workflowList, outputFormat, verboseMode)
			toolInstallers = append(toolInstallers, bitriseInstallers...)
		default:
			toolInstaller := m.installerFactory.CreateRegularInstaller(filePath, outputFormat, verboseMode)
			toolInstallers = append(toolInstallers, toolInstaller)
		}
	}

	var results []InstallerResult
	for _, installer := range toolInstallers {
		output, err := installer.Install()
		results = append(results, InstallerResult{
			ToolInstallerDescription: installer.Description(),
			Output:                   output,
			Err:                      err,
		})
	}
	return results, nil
}

func createBitriseInstallers(factory InstallerFactory, toolVersionFile string, workflowList []string, outputFormat string, verboseMode bool) []ToolInstaller {
	if len(workflowList) == 0 {
		return []ToolInstaller{factory.CreateBitriseInstaller(toolVersionFile, "", outputFormat, verboseMode)}
	}

	var installers []ToolInstaller
	for _, workflow := range workflowList {
		installer := factory.CreateBitriseInstaller(toolVersionFile, workflow, outputFormat, verboseMode)
		installers = append(installers, installer)
	}
	return installers
}
