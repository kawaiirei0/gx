# Installer Package

## Overview

The installer package handles the extraction and installation of Go distribution archives. It supports both ZIP (Windows) and tar.gz (Linux/macOS) formats and includes verification to ensure successful installation.

## Usage

### Creating an Installer

```go
import (
    "github.com/kawaiirei0/gx/internal/installer"
    "github.com/kawaiirei0/gx/internal/platform"
)

platformAdapter := platform.NewAdapter()
inst := installer.NewInstaller(platformAdapter)
```

### Installing a Go Version

```go
err := inst.Install(
    "/tmp/go1.21.5.tar.gz",  // Archive path
    "go1.21.5",               // Version
    "/home/user/.gx/versions/go1.21.5", // Destination
)
if err != nil {
    // Handle error
}
```

### Verifying an Installation

```go
err := inst.Verify("/home/user/.gx/versions/go1.21.5", "go1.21.5")
if err != nil {
    // Installation is invalid
}
```

### Uninstalling a Version

```go
err := inst.Uninstall("go1.21.5", "/home/user/.gx/versions/go1.21.5")
if err != nil {
    // Handle error
}
```

## Features

### Archive Format Support

- **ZIP**: Used for Windows distributions
- **tar.gz**: Used for Linux and macOS distributions

The installer automatically detects the format based on file extension.

### Directory Structure Normalization

Go distribution archives contain a top-level `go/` directory. The installer automatically strips this directory to create a clean installation:

```
Archive structure:        Installed structure:
go/                       (install path)/
├── bin/                  ├── bin/
│   └── go                │   └── go
├── src/                  ├── src/
└── pkg/                  └── pkg/
```

### Installation Verification

After extraction, the installer verifies:

1. **Directory structure**: Checks that `bin/` directory exists
2. **Executable presence**: Verifies `go` (or `go.exe` on Windows) exists
3. **Executable permissions**: Sets execute permissions on Unix systems
4. **Version validation**: Runs `go version` and verifies the output matches expected version

### Automatic Cleanup

If installation fails at any stage, the installer automatically cleans up partial installations to prevent corrupted state.

## Supported Archive Formats

### ZIP Format (Windows)

- Extracts all files and directories
- Preserves file permissions from archive
- Handles nested directory structures

### tar.gz Format (Linux/macOS)

- Supports regular files, directories, and symbolic links
- Preserves Unix file permissions
- Handles hard links and special files gracefully

## Error Handling

The installer uses custom error types:

- `ErrInstallFailed`: Installation process failed
- `ErrVersionNotInstalled`: Version not found during uninstall
- `ErrUninstallFailed`: Uninstall process failed

All errors include detailed context about what failed and why.

## Platform Compatibility

The installer works across:
- Windows (10+)
- Linux (all major distributions)
- macOS (10.15+)

Platform-specific operations (like setting executable permissions) are handled through the PlatformAdapter interface.

## Implementation Details

### Permission Handling

On Unix systems, the installer ensures the `go` executable has execute permissions. On Windows, this is handled automatically by the OS.

### Symbolic Links

The installer preserves symbolic links in tar.gz archives, which is important for Go's internal structure.

### Error Recovery

If extraction fails midway, the installer:
1. Stops the extraction process
2. Removes all partially extracted files
3. Returns a detailed error message

This ensures the installation directory is never left in a corrupted state.
