package installer

import (
	"os/exec"
	"testing"
)

// captureCommandArgsRegular captures arguments for regular installer tests
type captureCommandArgsRegular struct {
	capturedCommand string
	capturedArgs    []string
}

func (c *captureCommandArgsRegular) mockExec(command string, args ...string) *exec.Cmd {
	c.capturedCommand = command
	c.capturedArgs = args
	// Return a cmd that will execute successfully
	cmd := exec.Command("true")
	return cmd
}

func TestRegularToolInstaller_Description(t *testing.T) {
	tests := []struct {
		name         string
		filePath     string
		expectedDesc string
	}{
		{
			name:         "tool-versions file",
			filePath:     ".tool-versions",
			expectedDesc: "RegularToolInstaller for file: .tool-versions",
		},
		{
			name:         "ruby-version file",
			filePath:     ".ruby-version",
			expectedDesc: "RegularToolInstaller for file: .ruby-version",
		},
		{
			name:         "python-version with path",
			filePath:     "/home/user/project/.python-version",
			expectedDesc: "RegularToolInstaller for file: /home/user/project/.python-version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installer := NewRegularToolInstaller(tt.filePath, "plaintext", false)

			actual := installer.Description()
			if actual != tt.expectedDesc {
				t.Errorf("expected description %q, got: %q", tt.expectedDesc, actual)
			}
		})
	}
}

func TestRegularToolInstaller_Install_CommandArgs(t *testing.T) {
	capture := &captureCommandArgsRegular{}

	installer := &RegularToolInstaller{
		filePath:     ".tool-versions",
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
	expected := []string{"tools", "setup", "--config", ".tool-versions", "--output-format", "plaintext"}
	if !sliceEqualRegular(capture.capturedArgs, expected) {
		t.Errorf("unexpected args:\ngot:  %v\nwant: %v", capture.capturedArgs, expected)
	}
}

func TestRegularToolInstaller_Install_AllOutputFormats(t *testing.T) {
	formats := []string{"plaintext", "json", "bash"}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			capture := &captureCommandArgsRegular{}

			installer := &RegularToolInstaller{
				filePath:     ".tool-versions",
				outputFormat: format,
				verboseMode:  false,
				ExecCommand:  capture.mockExec,
			}

			_, _ = installer.Install()

			// Verify --output-format flag
			foundFormat := false
			for i, arg := range capture.capturedArgs {
				if arg == "--output-format" && i+1 < len(capture.capturedArgs) {
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

// sliceEqualRegular compares two string slices for equality
func sliceEqualRegular(a, b []string) bool {
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
