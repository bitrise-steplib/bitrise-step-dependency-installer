package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Integration tests that run the step binary via "go run main.go".
// These tests require the bitrise CLI to be installed on the system.

// runStep runs the step with the given environment variables and returns stdout, stderr, and error
func runStep(t *testing.T, env map[string]string) (string, string, error) {
	t.Helper()

	cmd := exec.Command("go", "run", "main.go")

	// Set environment variables
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	cmd.Env = append(cmd.Env, "verbose=true")

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// TestIntegration_SingleToolVersionsFile tests installing from a single .tool-versions file
// with multiple tool definitions
func TestIntegration_SingleToolVersionsFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup: Create temp directory with .tool-versions file
	tmpDir := t.TempDir()
	toolVersionsPath := filepath.Join(tmpDir, ".tool-versions")

	toolVersionsContent := `nodejs 20.10.0
python 3.11.6
golang 1.21.5
`
	if err := os.WriteFile(toolVersionsPath, []byte(toolVersionsContent), 0644); err != nil {
		t.Fatalf("failed to write .tool-versions file: %v", err)
	}

	// Execute
	stdout, stderr, err := runStep(t, map[string]string{
		"tool_version_file_list": toolVersionsPath,
		"output_format":          "plaintext",
	})

	// Log output for debugging
	t.Logf("stdout: %s", stdout)
	t.Logf("stderr: %s", stderr)

	// Verify - step should complete without error
	if err != nil {
		t.Errorf("step failed: %v", err)
	}
}

// TestIntegration_ToolVersionsAndGoVersion tests installing from both a .tool-versions
// and a .go-version file
func TestIntegration_ToolVersionsAndGoVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup: Create temp directory with both files
	tmpDir := t.TempDir()
	toolVersionsPath := filepath.Join(tmpDir, ".tool-versions")
	goVersionPath := filepath.Join(tmpDir, ".go-version")

	toolVersionsContent := `nodejs 20.10.0
`
	goVersionContent := "1.21.5\n"

	if err := os.WriteFile(toolVersionsPath, []byte(toolVersionsContent), 0644); err != nil {
		t.Fatalf("failed to write .tool-versions file: %v", err)
	}
	if err := os.WriteFile(goVersionPath, []byte(goVersionContent), 0644); err != nil {
		t.Fatalf("failed to write .go-version file: %v", err)
	}

	// Execute with newline-separated file list
	fileList := toolVersionsPath + "\n" + goVersionPath
	stdout, stderr, err := runStep(t, map[string]string{
		"tool_version_file_list": fileList,
		"output_format":          "json",
	})

	// Log output for debugging
	t.Logf("stdout: %s", stdout)
	t.Logf("stderr: %s", stderr)

	// Verify
	if err != nil {
		t.Errorf("step failed: %v", err)
	}
}

// TestIntegration_BitriseYMLGlobalTools tests installing from a bitrise.yml with
// tools defined at workflow level, called without specifying a workflow
func TestIntegration_BitriseYMLGlobalTools(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup: Create temp directory with bitrise.yml
	tmpDir := t.TempDir()
	bitriseYMLPath := filepath.Join(tmpDir, "bitrise.yml")

	bitriseYMLContent := `format_version: "13"
default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git

workflows:
  primary:
    tools:
      go: "1.21.5"
      node: "20.10.0"
    steps:
      - script:
          inputs:
            - content: echo "Hello"
`
	if err := os.WriteFile(bitriseYMLPath, []byte(bitriseYMLContent), 0644); err != nil {
		t.Fatalf("failed to write bitrise.yml file: %v", err)
	}

	// Execute: No workflow specified
	stdout, stderr, err := runStep(t, map[string]string{
		"tool_version_file_list": bitriseYMLPath,
		"output_format":          "plaintext",
	})

	// Log output for debugging
	t.Logf("stdout: %s", stdout)
	t.Logf("stderr: %s", stderr)

	// Note: Without a workflow specified, bitrise CLI may fail or use default behavior
	// This test documents the actual behavior
	if err != nil {
		t.Logf("step returned error (expected when no workflow specified): %v", err) // TODO: Fix with BE-1960
	}
}

// TestIntegration_BitriseYMLSingleWorkflow tests installing from a bitrise.yml
// with a single workflow that defines a tools block
func TestIntegration_BitriseYMLSingleWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup: Create temp directory with bitrise.yml
	tmpDir := t.TempDir()
	bitriseYMLPath := filepath.Join(tmpDir, "bitrise.yml")

	bitriseYMLContent := `format_version: "13"
default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git

workflows:
  build:
    tools:
      go: "1.21.5"
      node: "20.10.0"
    steps:
      - script:
          inputs:
            - content: echo "Building..."
`
	if err := os.WriteFile(bitriseYMLPath, []byte(bitriseYMLContent), 0644); err != nil {
		t.Fatalf("failed to write bitrise.yml file: %v", err)
	}

	// Execute: Single workflow specified
	stdout, stderr, err := runStep(t, map[string]string{
		"tool_version_file_list":   bitriseYMLPath,
		"bitrise_workflow_id_list": "build",
		"output_format":            "bash",
	})

	// Log output for debugging
	t.Logf("stdout: %s", stdout)
	t.Logf("stderr: %s", stderr)

	// Verify
	if err != nil {
		t.Errorf("step failed: %v", err)
	}
}

// TestIntegration_BitriseYMLMultipleWorkflows tests installing from a bitrise.yml
// with multiple workflows that define tools blocks
func TestIntegration_BitriseYMLMultipleWorkflows(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup: Create temp directory with bitrise.yml
	tmpDir := t.TempDir()
	bitriseYMLPath := filepath.Join(tmpDir, "bitrise.yml")

	bitriseYMLContent := `format_version: "13"
default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git

workflows:
  build:
    tools:
      go: "1.21.5"
    steps:
      - script:
          inputs:
            - content: echo "Building..."

  test:
    tools:
      go: "1.21.5"
      node: "20.10.0"
    steps:
      - script:
          inputs:
            - content: echo "Testing..."

  deploy:
    tools:
      go: "1.21.5"
      python: "3.11.6"
    steps:
      - script:
          inputs:
            - content: echo "Deploying..."
`
	if err := os.WriteFile(bitriseYMLPath, []byte(bitriseYMLContent), 0644); err != nil {
		t.Fatalf("failed to write bitrise.yml file: %v", err)
	}

	// Execute: Multiple workflows specified (newline-separated)
	workflows := "build\ntest\ndeploy"
	stdout, stderr, err := runStep(t, map[string]string{
		"tool_version_file_list":   bitriseYMLPath,
		"bitrise_workflow_id_list": workflows,
		"output_format":            "json",
	})

	// Log output for debugging
	t.Logf("stdout: %s", stdout)
	t.Logf("stderr: %s", stderr)

	// Verify
	if err != nil {
		t.Errorf("step failed: %v", err)
	}
}

// TestIntegration_BitriseYMLGlobalAndGoVersion tests installing from both a bitrise.yml
// with workflow tools and a .go-version file
func TestIntegration_BitriseYMLGlobalAndGoVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup: Create temp directory with both files
	tmpDir := t.TempDir()
	bitriseYMLPath := filepath.Join(tmpDir, "bitrise.yml")
	goVersionPath := filepath.Join(tmpDir, ".go-version")

	bitriseYMLContent := `format_version: "13"
default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git

workflows:
  primary:
    tools:
      node: "20.10.0"
    steps:
      - script:
          inputs:
            - content: echo "Hello"
`
	goVersionContent := "1.21.5\n"

	if err := os.WriteFile(bitriseYMLPath, []byte(bitriseYMLContent), 0644); err != nil {
		t.Fatalf("failed to write bitrise.yml file: %v", err)
	}
	if err := os.WriteFile(goVersionPath, []byte(goVersionContent), 0644); err != nil {
		t.Fatalf("failed to write .go-version file: %v", err)
	}

	// Execute: bitrise.yml with workflow + .go-version
	fileList := bitriseYMLPath + "\n" + goVersionPath
	stdout, stderr, err := runStep(t, map[string]string{
		"tool_version_file_list":   fileList,
		"bitrise_workflow_id_list": "primary",
		"output_format":            "plaintext",
	})

	// Log output for debugging
	t.Logf("stdout: %s", stdout)
	t.Logf("stderr: %s", stderr)

	// Verify
	if err != nil {
		t.Errorf("step failed: %v", err)
	}
}
