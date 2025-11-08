# Build System Overview

This document provides an overview of the gx build system implementation.

## Components

### 1. Makefile
**Location:** `Makefile`

A comprehensive GNU Make-based build system with the following targets:

- `make build` - Build for current platform
- `make build-all` - Build for all supported platforms
- `make release` - Create release packages with checksums
- `make install` - Install to GOPATH/bin
- `make clean` - Remove build artifacts
- `make test` - Run tests with coverage
- `make version` - Show version information
- `make deps` - Download dependencies
- `make fmt` - Format code
- `make lint` - Run linter

**Features:**
- Automatic version detection from git tags
- Parallel build support
- Cross-platform compilation
- Embedded version information via ldflags
- Binary size optimization (-s -w flags)

### 2. PowerShell Build Script
**Location:** `build.ps1`

Windows-native build script for users without Make:

```powershell
.\build.ps1 build              # Build current platform
.\build.ps1 build-all          # Build all platforms
.\build.ps1 release -Version v1.0.0  # Create release
.\build.ps1 install            # Install locally
.\build.ps1 clean              # Clean artifacts
.\build.ps1 test               # Run tests
.\build.ps1 version            # Show version
```

**Features:**
- Colored console output
- Progress indication
- Error handling
- Native Windows support
- Automatic checksum generation

### 3. Bash Build Script
**Location:** `build.sh`

Unix-like build script for Linux/macOS:

```bash
./build.sh build               # Build current platform
./build.sh build-all           # Build all platforms
VERSION=v1.0.0 ./build.sh release  # Create release
./build.sh install             # Install locally
./build.sh clean               # Clean artifacts
./build.sh test                # Run tests
./build.sh version             # Show version
```

**Features:**
- POSIX-compliant
- Colored output
- Environment variable support
- Portable across Unix systems

### 4. Release Helper Script
**Location:** `scripts/release.sh`

Automated release workflow script:

```bash
./scripts/release.sh v1.0.0
```

**Workflow:**
1. Validates version format
2. Checks for uncommitted changes
3. Verifies current branch
4. Runs tests
5. Updates changelog
6. Creates git tag
7. Builds release packages
8. Provides next steps

### 5. GitHub Actions Workflows

#### CI Workflow
**Location:** `.github/workflows/ci.yml`

Continuous integration workflow that:
- Runs tests on Windows, Linux, macOS
- Tests with Go 1.21, 1.22, 1.23
- Generates coverage reports
- Builds for all platforms
- Runs linter

Triggered on:
- Push to main/develop branches
- Pull requests to main/develop

#### Release Workflow
**Location:** `.github/workflows/release.yml`

Automated release workflow that:
- Builds all platform binaries
- Creates release packages
- Generates checksums
- Creates GitHub release
- Uploads artifacts

Triggered on:
- Git tags matching `v*`
- Manual workflow dispatch

## Version Management

### Version Information

Version information is embedded in binaries using Go ldflags:

```go
// cmd/gx/main.go
var (
    Version   = "dev"
    Commit    = "unknown"
    BuildDate = "unknown"
)
```

Set during build:
```bash
-ldflags "-X main.Version=v1.0.0 -X main.Commit=abc1234 -X main.BuildDate=2024-01-15T10:30:00Z"
```

### Version Detection

Automatic version detection from git:
```bash
git describe --tags --always --dirty
```

Manual override:
```bash
VERSION=v1.0.0 make build
```

## Supported Platforms

The build system supports cross-compilation for:

| OS      | Architecture | Binary Extension |
|---------|-------------|------------------|
| Windows | amd64       | .exe            |
| Windows | 386         | .exe            |
| Linux   | amd64       | (none)          |
| Linux   | 386         | (none)          |
| Linux   | arm64       | (none)          |
| macOS   | amd64       | (none)          |
| macOS   | arm64       | (none)          |

## Build Artifacts

### Directory Structure

```
gx/
├── build/                    # Build output
│   ├── windows_amd64/
│   │   └── gx.exe
│   ├── linux_amd64/
│   │   └── gx
│   └── ...
├── dist/                     # Release packages
│   ├── gx-v1.0.0-windows-amd64.zip
│   ├── gx-v1.0.0-linux-amd64.tar.gz
│   ├── checksums.txt
│   └── ...
└── coverage.out             # Test coverage
```

### Release Packages

- **Windows:** ZIP archives containing .exe
- **Linux/macOS:** tar.gz archives containing binary
- **Checksums:** SHA256 checksums for all packages

## Build Flags

### Standard Flags

```bash
-ldflags "-X main.Version=... -X main.Commit=... -X main.BuildDate=... -s -w"
```

- `-X main.Version`: Set version string
- `-X main.Commit`: Set git commit hash
- `-X main.BuildDate`: Set build timestamp
- `-s`: Strip symbol table
- `-w`: Strip DWARF debug info

### Cross-Compilation

```bash
GOOS=linux GOARCH=amd64 go build ...
```

## Documentation

- **BUILD.md** - Comprehensive build documentation
- **CHANGELOG.md** - Version history and changes
- **.github/RELEASE_PROCESS.md** - Release workflow guide
- **README.md** - Updated with build instructions

## Testing

### Local Testing

```bash
# Run tests
make test

# With coverage report
make test-coverage
```

### CI Testing

Automated testing on:
- Multiple OS (Windows, Linux, macOS)
- Multiple Go versions (1.21, 1.22, 1.23)
- Race condition detection
- Coverage reporting

## Usage Examples

### Development Build

```bash
# Quick build for testing
go build -o gx ./cmd/gx

# Or with build script
.\build.ps1 build
```

### Production Build

```bash
# Build with version info
VERSION=v1.0.0 make build

# Build all platforms
make build-all
```

### Release

```bash
# Automated release
./scripts/release.sh v1.0.0

# Manual release
VERSION=v1.0.0 make release
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

## Maintenance

### Adding New Platform

1. Add to PLATFORMS in Makefile:
   ```makefile
   PLATFORMS := \
       existing/platforms \
       newos/newarch
   ```

2. Add to $Platforms in build.ps1:
   ```powershell
   $Platforms = @(
       @{OS="newos"; Arch="newarch"}
   )
   ```

3. Add to PLATFORMS in build.sh:
   ```bash
   declare -a PLATFORMS=(
       "newos/newarch"
   )
   ```

### Updating Go Version

Update in:
- `.github/workflows/ci.yml` - Go version matrix
- `.github/workflows/release.yml` - Go setup version
- `go.mod` - Go version requirement

## Troubleshooting

### Common Issues

1. **"make: command not found"**
   - Use platform-specific scripts (build.ps1 or build.sh)

2. **"permission denied" on scripts**
   - Linux/macOS: `chmod +x build.sh scripts/release.sh`

3. **Version shows as "dev"**
   - Not in git repository
   - No git tags exist
   - Git not installed
   - Solution: Manually set VERSION environment variable

4. **Cross-compilation fails**
   - Ensure CGO is disabled (default for gx)
   - Check GOOS/GOARCH combination is valid

## Best Practices

1. **Always run tests before release**
   ```bash
   make test
   ```

2. **Use semantic versioning**
   - v1.0.0 for stable releases
   - v1.0.0-beta.1 for pre-releases

3. **Update CHANGELOG.md before release**

4. **Verify release packages**
   ```bash
   # Check checksums
   cat dist/checksums.txt
   
   # Test binary
   ./gx --version
   ```

5. **Tag releases in git**
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   ```

## Future Enhancements

Potential improvements:
- [ ] Docker build support
- [ ] Homebrew formula generation
- [ ] Chocolatey package generation
- [ ] Snap package support
- [ ] Code signing for binaries
- [ ] Notarization for macOS binaries
