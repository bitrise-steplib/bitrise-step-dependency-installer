package installer

import (
	"fmt"
	"testing"
)

// mockInstaller is a test double for the ToolInstaller interface
type mockInstaller struct {
	description      string
	output           []byte
	err              error
	installCallCount int
}

func (m *mockInstaller) Install() ([]byte, error) {
	m.installCallCount++
	return m.output, m.err
}

func (m *mockInstaller) Description() string {
	return m.description
}

// mockFactory is a test double for InstallerFactory
type mockFactory struct {
	bitriseInstallers []*mockInstaller
	regularInstallers []*mockInstaller
	bitriseCallCount  int
	regularCallCount  int
}

func (f *mockFactory) CreateBitriseInstaller(filePath, workflow, outputFormat string, verboseMode bool) ToolInstaller {
	if f.bitriseCallCount >= len(f.bitriseInstallers) {
		return &mockInstaller{description: fmt.Sprintf("mock bitrise installer %d", f.bitriseCallCount)}
	}
	installer := f.bitriseInstallers[f.bitriseCallCount]
	f.bitriseCallCount++
	return installer
}

func (f *mockFactory) CreateRegularInstaller(filePath, outputFormat string, verboseMode bool) ToolInstaller {
	if f.regularCallCount >= len(f.regularInstallers) {
		return &mockInstaller{description: fmt.Sprintf("mock regular installer %d", f.regularCallCount)}
	}
	installer := f.regularInstallers[f.regularCallCount]
	f.regularCallCount++
	return installer
}

func TestInstallerManager_InstallAll_EmptyFileList(t *testing.T) {
	factory := &mockFactory{}
	manager := NewInstallerManager(factory)

	results, err := manager.InstallAll([]string{}, []string{}, "plaintext", false)

	if err != nil {
		t.Errorf("expected no error for empty file list, got: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected 0 results for empty file list, got: %d", len(results))
	}
}

func TestInstallerManager_InstallAll_MultipleBitriseYMLFiles(t *testing.T) {
	fileList := []string{
		"path/to/bitrise.yml",
		"another/bitrise.yml",
	}

	factory := &mockFactory{}
	manager := NewInstallerManager(factory)

	results, err := manager.InstallAll(fileList, []string{}, "plaintext", false)

	if err == nil {
		t.Error("expected error for multiple bitrise.yml files, got nil")
	}

	expectedErrMsg := "multiple bitrise.yml files detected in the input list, only one is allowed"
	if err.Error() != expectedErrMsg {
		t.Errorf("expected error message %q, got %q", expectedErrMsg, err.Error())
	}

	if results != nil {
		t.Errorf("expected nil results on error, got: %v", results)
	}
}

func TestInstallerManager_InstallAll_InvokesCorrectInstallers(t *testing.T) {
	fileList := []string{
		".tool-versions",
		"bitrise.yml",
		".ruby-version",
	}

	factory := &mockFactory{
		regularInstallers: []*mockInstaller{
			{description: "regular1", output: []byte("output1"), err: nil},
			{description: "regular2", output: []byte("output2"), err: nil},
		},
		bitriseInstallers: []*mockInstaller{
			{description: "bitrise1", output: []byte("output3"), err: nil},
		},
	}

	manager := NewInstallerManager(factory)
	results, err := manager.InstallAll(fileList, []string{}, "json", false)

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if factory.regularCallCount != 2 {
		t.Errorf("expected 2 regular installer creations, got: %d", factory.regularCallCount)
	}

	if factory.bitriseCallCount != 1 {
		t.Errorf("expected 1 bitrise installer creation, got: %d", factory.bitriseCallCount)
	}

	// Verify Install() was called on each
	if factory.regularInstallers[0].installCallCount != 1 {
		t.Errorf("expected Install() called once on first regular installer, got: %d", factory.regularInstallers[0].installCallCount)
	}

	if factory.bitriseInstallers[0].installCallCount != 1 {
		t.Errorf("expected Install() called once on bitrise installer, got: %d", factory.bitriseInstallers[0].installCallCount)
	}

	// Verify results
	if len(results) != 3 {
		t.Errorf("expected 3 results, got: %d", len(results))
	}
}

func TestInstallerManager_InstallAll_WithWorkflows(t *testing.T) {
	fileList := []string{"bitrise.yml"}
	workflows := []string{"build", "test"}

	factory := &mockFactory{
		bitriseInstallers: []*mockInstaller{
			{description: "bitrise-build", output: []byte("build output"), err: nil},
			{description: "bitrise-test", output: []byte("test output"), err: nil},
		},
	}

	manager := NewInstallerManager(factory)
	results, err := manager.InstallAll(fileList, workflows, "plaintext", false)

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	// Should create 2 bitrise installers (one per workflow)
	if factory.bitriseCallCount != 2 {
		t.Errorf("expected 2 bitrise installer creations, got: %d", factory.bitriseCallCount)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 results, got: %d", len(results))
	}
}

func TestInstallerManager_InstallAll_HandlesErrors(t *testing.T) {
	fileList := []string{".tool-versions", ".ruby-version"}

	factory := &mockFactory{
		regularInstallers: []*mockInstaller{
			{description: "installer1", output: []byte("success"), err: nil},
			{description: "installer2", output: nil, err: fmt.Errorf("installation failed")},
		},
	}

	manager := NewInstallerManager(factory)
	results, err := manager.InstallAll(fileList, []string{}, "plaintext", false)

	// InstallAll should not return error, but collect errors in results
	if err != nil {
		t.Errorf("expected no error from InstallAll, got: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got: %d", len(results))
	}

	// First should succeed
	if results[0].Err != nil {
		t.Errorf("expected first result to have no error, got: %v", results[0].Err)
	}

	// Second should have error
	if results[1].Err == nil {
		t.Error("expected second result to have error, got nil")
	}
}
