package installer

import (
	"os/exec"
	"testing"
)

// captureCommandArgs captures arguments passed to exec.Command for verification
type captureCommandArgs struct {
	capturedCommand string
	capturedArgs    []string
}

func (c *captureCommandArgs) mockExec(command string, args ...string) *exec.Cmd {
	c.capturedCommand = command
	c.capturedArgs = args
	// Return a cmd that will fail but won't actually execute
	cmd := exec.Command("true") // Use 'true' which exists on all Unix systems
	return cmd
}

func TestBitriseToolInstaller_Description_NoWorkflow(t *testing.T) {
	installer := NewBitriseToolInstaller("/path/to/bitrise.yml", "", "plaintext", false)

	expected := "BitriseToolInstaller for file: /path/to/bitrise.yml (no workflow)"
	actual := installer.Description()

	if actual != expected {
		t.Errorf("expected description %q, got: %q", expected, actual)
	}
}

func TestBitriseToolInstaller_Description_WithWorkflow(t *testing.T) {
	installer := NewBitriseToolInstaller("/path/to/bitrise.yml", "test", "plaintext", false)

	expected := "BitriseToolInstaller for file: /path/to/bitrise.yml and workflow: test"
	actual := installer.Description()

	if actual != expected {
		t.Errorf("expected description %q, got: %q", expected, actual)
	}
}

func TestBitriseToolInstaller_Install_CommandArgs_NoWorkflow(t *testing.T) {
	capture := &captureCommandArgs{}

	installer := &BitriseToolInstaller{
		filePath:     "/path/to/config.yml",
		workflow:     "",
		outputFormat: "json",
		verboseMode:  false,
		ExecCommand:  capture.mockExec,
	}

	_, _ = installer.Install()

	// Verify command
	if capture.capturedCommand != "bitrise" {
		t.Errorf("expected command 'bitrise', got: %q", capture.capturedCommand)
	}

	// Verify args
	expected := []string{"tools", "setup", "--config", "/path/to/config.yml", "--output-format", "json"}
	if !sliceEqual(capture.capturedArgs, expected) {
		t.Errorf("unexpected args:\ngot:  %v\nwant: %v", capture.capturedArgs, expected)
	}

	// Verify --workflow is NOT present
	for _, arg := range capture.capturedArgs {
		if arg == "--workflow" {
			t.Errorf("expected NO --workflow flag for empty workflow, got args: %v", capture.capturedArgs)
		}
	}
}

func TestBitriseToolInstaller_Install_CommandArgs_WithWorkflow(t *testing.T) {
	capture := &captureCommandArgs{}

	installer := &BitriseToolInstaller{
		filePath:     "bitrise.yml",
		workflow:     "build",
		outputFormat: "plaintext",
		verboseMode:  false,
		ExecCommand:  capture.mockExec,
	}

	_, _ = installer.Install()

	// Verify command
	if capture.capturedCommand != "bitrise" {
		t.Errorf("expected command 'bitrise', got: %q", capture.capturedCommand)
	}

	// Verify args
	expected := []string{"tools", "setup", "--config", "bitrise.yml", "--workflow", "build", "--output-format", "plaintext"}
	if !sliceEqual(capture.capturedArgs, expected) {
		t.Errorf("unexpected args:\ngot:  %v\nwant: %v", capture.capturedArgs, expected)
	}
}

func TestBitriseToolInstaller_Install_AllOutputFormats(t *testing.T) {
	formats := []string{"plaintext", "json", "bash"}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			capture := &captureCommandArgs{}

			installer := &BitriseToolInstaller{
				filePath:     "test.yml",
				workflow:     "",
				outputFormat: format,
				verboseMode:  false,
				ExecCommand:  capture.mockExec,
			}

			_, _ = installer.Install()

			// Verify --output-format flag
			foundFormat := false
			for i, arg := range capture.capturedArgs {
				if arg == "--output-format" {
					if capture.capturedArgs[i+1] == format {
						foundFormat = true
					}
				}
			}

			if !foundFormat {
				t.Errorf("expected --output-format %s in args: %v", format, capture.capturedArgs)
			}
		})
	}
}

// sliceEqual compares two string slices for equality
func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
