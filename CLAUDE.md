# CLAUDE.md

## Project Overview

This is a Bitrise Step that wraps the Bitrise CLI tool to install dependencies defined in tool version files. It supports multiple tool version file formats including `.tool-versions`, `.ruby-version`, `.node-version`, `.python-version`, `.java-version`, `.go-version`, `.terraform-version`, `.kubectl-version`, and `bitrise.yml` files with a `tools:` section.

## Development Commands

### Testing
```bash
# Run all tests
go test ./...

# Run tests in a specific package
go test ./installer

# Run a specific test
go test -run TestName ./installer

# Run tests with verbose output
go test -v ./...
```

### Building
```bash
# Build the step binary
go build -o bitrise-step-dependency-installer .
```

### Running the Step Locally
```bash
# Test the step using Bitrise CLI (requires .bitrise.secrets.yml)
bitrise run test
```

## Architecture

### Core Components

**main.go**: Entry point that reads environment variables, parses inputs, and orchestrates the installation process through the installer package.

**installer package**: Contains the core installation logic with three main types:

1. **ToolInstaller interface**: Defines the contract for all installer implementations
   - `Install() Result`: Executes the installation
   - `Description() string`: Returns a description of what's being installed

2. **RegularToolInstaller**: Handles standard tool version files (`.tool-versions`, `.ruby-version`, etc.)
   - Calls `bitrise tools setup --config <file>`

3. **BitriseToolInstaller**: Handles `bitrise.yml` files with `tools:` sections
   - Can target either:
     - Global `tools:` block (no workflow specified)
     - Workflow-specific `tools:` blocks (workflow ID specified)
   - Multiple workflow IDs create separate installer instances

### Key Design Patterns

**Strategy Pattern**: The `ToolInstaller` interface allows different installation strategies (regular vs bitrise) to be handled uniformly.

**Result Types**: Two result types are used:
- `Result`: Internal type with lowercase fields (`output []byte`, `err error`)
- `InstallerResult`: Public type with exported fields for returning to main

**File Type Detection**: File paths are examined by suffix to determine installer type:
- `bitrise.yml` suffix → `BitriseToolInstaller`
- All others → `RegularToolInstaller`

### Installation Flow

1. Parse environment variables in main.go
2. Split file paths and workflow IDs by newline
3. For each file path, create appropriate installer(s):
   - Regular files: one `RegularToolInstaller`
   - bitrise.yml: one or more `BitriseToolInstaller` instances (based on workflow list)
4. Execute all installers sequentially
5. Collect results and exit with error code if any installation failed

### Important Constraints

- **Only one bitrise.yml**: Multiple bitrise.yml files in the input list will cause an error
- **Workflow filtering**: The `bitrise_workflow_id_list` input only affects bitrise.yml files, not regular tool version files
- **Empty workflow list**: When no workflows are specified for bitrise.yml, it targets the global `tools:` block

## Module Path

The module path is `github.com/bitrise-io/bitrise-step-dependency-installer`. All imports should use this path.

## Environment Variables (Step Inputs)

- `tool_version_file_list`: Newline-separated list of tool version file paths (required)
- `bitrise_workflow_id_list`: Newline-separated list of workflow IDs for bitrise.yml files (optional)
- `verbose`: Enable verbose output ("true"/"false")
- `output_format`: Format for bitrise CLI output ("json", "plaintext", "bash")

## Known Issues

The test files currently reference old struct field names (`workflowList` instead of `workflow`) which will cause test failures. The installers use a single `workflow` string field, not a `workflowList` slice.
