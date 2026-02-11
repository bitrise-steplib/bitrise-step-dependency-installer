package installer

import (
	"log"
	"os/exec"
	"strings"
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
func TestInstaller_Install_CommandArgs(t *testing.T) {

	tests := []struct {
		name     string
		filePath string
		workflow string
		expected string
	}{
		{
			name:     "tool-versions file",
			filePath: ".tool-versions",
			workflow: "",
			expected: "tools setup --config .tool-versions",
		},
		{
			name:     "ruby-version file",
			filePath: ".ruby-version",
			workflow: "",
			expected: "tools setup --config .ruby-version",
		},
		{
			name:     "bitrise.yml file with workflow",
			filePath: "bitrise.yml",
			workflow: "test",
			expected: "tools setup --config bitrise.yml --workflow test",
		},
		{
			name:     "bitrise.yml file with no workflow",
			filePath: "bitrise.yml",
			workflow: "",
			expected: "tools setup --config bitrise.yml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("Running test: %s\n", tt.name)
			capture := &captureCommandArgsRegular{}

			_, _ = Install(tt.filePath, tt.workflow, false, capture.mockExec)

			if strings.Trim(strings.Join(capture.capturedArgs, " "), " ") != tt.expected {
				t.Errorf("unexpected args:\ngot:  %v\nwant: %v", capture.capturedArgs, tt.expected)
			}
		})
	}
}
