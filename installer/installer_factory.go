package installer

type InstallerFactory interface {
	CreateBitriseInstaller(filePath, workflow, outputFormat string, verboseMode bool) ToolInstaller
	CreateRegularInstaller(filePath, outputFormat string, verboseMode bool) ToolInstaller
}

type DefaultInstallerFactory struct{}

func (f *DefaultInstallerFactory) CreateBitriseInstaller(filePath, workflow, outputFormat string, verboseMode bool) ToolInstaller {
	return NewBitriseToolInstaller(filePath, workflow, outputFormat, verboseMode)
}

func (f *DefaultInstallerFactory) CreateRegularInstaller(filePath, outputFormat string, verboseMode bool) ToolInstaller {
	return NewRegularToolInstaller(filePath, outputFormat, verboseMode)
}
