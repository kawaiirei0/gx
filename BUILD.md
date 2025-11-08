# Building gx

This document describes how to build and release gx from source.

## Prerequisites

- Go 1.21 or later
- Git (for version information)
- Make (optional, for Makefile usage)
- zip/tar (for creating release packages)

## Quick Start

### Build for Current Platform

**Using Make (Linux/macOS/Windows with Make):**
```bash
make build
```

**Using PowerShell (Windows):**
```powershell
.\build.ps1 build
```

**Using Bash (Linux/macOS):**
```bash
./build.sh build
```

**Using Go directly:**
```bash
go build -o build/gx ./cmd/gx
```

The binary will be created in the `build/` directory.

## Build Targets

### Build for All Platforms

Build binaries for all supported platforms:

```bash
# Using Make
make build-all

# Using PowerShell
.\build.ps1 build-all

# Using Bash
./build.sh build-all
```

Supported platforms:
- Windows: amd64, 386
- Linux: amd64, 386, arm64
- macOS (Darwin): amd64, arm64

### Create Release Packages

Build all platforms and create compressed release packages:

```bash
# Using Make
make release

# Using PowerShell
.\build.ps1 release

# Using Bash
./build.sh release
```

This will:
1. Build binaries for all platforms
2. Create `.zip` files for Windows builds
3. Create `.tar.gz` files for Linux/macOS builds
4. Generate SHA256 checksums in `dist/checksums.txt`

Release packages will be created in the `dist/` directory.

### Install Locally

Install gx to your `$GOPATH/bin`:

```bash
# Using Make
make install

# Using PowerShell
.\build.ps1 install

# Using Bash
./build.sh install
```

## Version Management

Version information is embedded in the binary during build using Git tags and ldflags.

### Automatic Version Detection

By default, the build scripts will:
1. Use `git describe --tags --always --dirty` to determine the version
2. Use the short commit hash
3. Use the current UTC timestamp as build date

### Manual Version Override

You can override the version:

**Make:**
```bash
make build VERSION=v1.0.0
```

**PowerShell:**
```powershell
.\build.ps1 build -Version v1.0.0
```

**Bash:**
```bash
VERSION=v1.0.0 ./build.sh build
```

### Version Information in Binary

The version information is displayed with:
```bash
gx --version
```

Output example:
```
gx version v1.0.0 (commit: abc1234, built: 2024-01-15T10:30:00Z)
```

## Build Flags

The build process uses the following ldflags:

```
-X main.Version=<version>
-X main.Commit=<commit>
-X main.BuildDate=<date>
-s -w  # Strip debug information and symbol table
```

These flags:
- Embed version information into the binary
- Reduce binary size by stripping debug symbols

## Development Builds

For development, you can build without version information:

```bash
go build -o gx ./cmd/gx
```

Or use the build scripts which will automatically use "dev" as the version.

## Testing

Run tests before building:

```bash
# Using Make
make test

# Using PowerShell
.\build.ps1 test

# Using Bash
./build.sh test
```

## Cleaning Build Artifacts

Remove all build artifacts:

```bash
# Using Make
make clean

# Using PowerShell
.\build.ps1 clean

# Using Bash
./build.sh clean
```

This removes:
- `build/` directory (compiled binaries)
- `dist/` directory (release packages)

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Build Release
        run: make release
      
      - name: Create GitHub Release
        uses: softprops/action-gh-release@v1
        with:
          files: dist/*
```

## Build Script Details

### Makefile

The Makefile provides the most comprehensive build experience with:
- Automatic platform detection
- Parallel builds support
- Dependency management
- Test coverage reports

### PowerShell Script (build.ps1)

The PowerShell script is designed for Windows users who don't have Make installed:
- Native Windows support
- Colored output
- Error handling
- Progress indication

### Bash Script (build.sh)

The Bash script provides a Make alternative for Linux/macOS:
- POSIX-compliant
- Colored output
- Portable across Unix-like systems

## Troubleshooting

### "command not found: make"

Use the platform-specific build scripts:
- Windows: `.\build.ps1`
- Linux/macOS: `./build.sh`

### "permission denied" on Linux/macOS

Make the script executable:
```bash
chmod +x build.sh
```

### Cross-compilation Issues

Ensure you have the necessary cross-compilation tools:
- For CGO-enabled builds, you may need cross-compilers
- gx doesn't use CGO by default, so cross-compilation should work out of the box

### Version Shows as "dev"

This happens when:
- Git is not installed
- You're not in a Git repository
- No tags exist in the repository

To fix, either:
- Install Git and ensure you're in a Git repo
- Manually specify the version: `VERSION=v1.0.0 make build`

## Release Checklist

Before creating a release:

1. [ ] Update version in git tag
2. [ ] Run tests: `make test`
3. [ ] Build all platforms: `make build-all`
4. [ ] Create release packages: `make release`
5. [ ] Verify checksums: `cat dist/checksums.txt`
6. [ ] Test binaries on target platforms
7. [ ] Create GitHub release with artifacts
8. [ ] Update documentation

## Additional Make Targets

```bash
make help           # Show all available targets
make deps           # Download dependencies
make fmt            # Format code
make lint           # Run linter (requires golangci-lint)
make test-coverage  # Generate HTML coverage report
make version        # Show version information
```
