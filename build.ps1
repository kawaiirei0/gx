# gx Build Script for Windows
# PowerShell build script for users without Make

param(
    [Parameter(Position=0)]
    [string]$Target = "build",
    
    [Parameter()]
    [string]$Version = "",
    
    [Parameter()]
    [switch]$Help
)

# Configuration
$BinaryName = "gx"
$MainPackage = "./cmd/gx"
$BuildDir = "build"
$DistDir = "dist"

# Get version information
if ([string]::IsNullOrEmpty($Version)) {
    try {
        $Version = git describe --tags --always --dirty 2>$null
        if ([string]::IsNullOrEmpty($Version)) {
            $Version = "dev"
        }
    } catch {
        $Version = "dev"
    }
}

try {
    $Commit = git rev-parse --short HEAD 2>$null
    if ([string]::IsNullOrEmpty($Commit)) {
        $Commit = "unknown"
    }
} catch {
    $Commit = "unknown"
}

$BuildDate = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")

# Build flags
$LDFlags = "-X main.Version=$Version -X main.Commit=$Commit -X main.BuildDate=$BuildDate -s -w"

# Platform configurations
$Platforms = @(
    @{OS="windows"; Arch="amd64"},
    @{OS="windows"; Arch="386"},
    @{OS="linux"; Arch="amd64"},
    @{OS="linux"; Arch="386"},
    @{OS="linux"; Arch="arm64"},
    @{OS="darwin"; Arch="amd64"},
    @{OS="darwin"; Arch="arm64"}
)

function Show-Help {
    Write-Host "gx Build Script for Windows"
    Write-Host ""
    Write-Host "Usage: .\build.ps1 [target] [-Version <version>]"
    Write-Host ""
    Write-Host "Available targets:"
    Write-Host "  build          - Build for current platform (default)"
    Write-Host "  build-all      - Build for all supported platforms"
    Write-Host "  release        - Create release packages for all platforms"
    Write-Host "  install        - Install gx to GOPATH\bin"
    Write-Host "  clean          - Remove build artifacts"
    Write-Host "  test           - Run tests"
    Write-Host "  version        - Show version information"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\build.ps1 build"
    Write-Host "  .\build.ps1 build-all -Version v1.0.0"
    Write-Host "  .\build.ps1 release"
}

function Build-Current {
    Write-Host "Building $BinaryName for current platform..."
    
    if (-not (Test-Path $BuildDir)) {
        New-Item -ItemType Directory -Path $BuildDir | Out-Null
    }
    
    # Detect current platform
    $CurrentOS = if ($IsWindows -or $env:OS -match "Windows") { "windows" } 
                 elseif ($IsMacOS) { "darwin" } 
                 else { "linux" }
    $CurrentArch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
    
    $OutputPath = Join-Path $BuildDir "$BinaryName.exe"
    
    # Set GOOS and GOARCH explicitly for current platform
    $env:GOOS = $CurrentOS
    $env:GOARCH = $CurrentArch
    
    & go build -ldflags $LDFlags -o $OutputPath $MainPackage
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Build complete: $OutputPath" -ForegroundColor Green
    } else {
        Write-Host "Build failed!" -ForegroundColor Red
        exit 1
    }
}

function Build-All {
    Write-Host "Building for all platforms..."
    
    Clean-Build
    
    foreach ($Platform in $Platforms) {
        $OS = $Platform.OS
        $Arch = $Platform.Arch
        $PlatformDir = "${OS}_${Arch}"
        $OutputDir = Join-Path $BuildDir $PlatformDir
        
        Write-Host "Building for $OS/$Arch..."
        
        if (-not (Test-Path $OutputDir)) {
            New-Item -ItemType Directory -Path $OutputDir | Out-Null
        }
        
        $BinaryExt = if ($OS -eq "windows") { ".exe" } else { "" }
        $OutputPath = Join-Path $OutputDir "$BinaryName$BinaryExt"
        
        $env:GOOS = $OS
        $env:GOARCH = $Arch
        
        & go build -ldflags $LDFlags -o $OutputPath $MainPackage
        
        if ($LASTEXITCODE -ne 0) {
            Write-Host "Build failed for $OS/$Arch!" -ForegroundColor Red
            exit 1
        }
    }
    
    Write-Host "All builds complete!" -ForegroundColor Green
}

function Create-Release {
    Write-Host "Creating release packages..."
    
    Build-All
    
    if (-not (Test-Path $DistDir)) {
        New-Item -ItemType Directory -Path $DistDir | Out-Null
    }
    
    foreach ($Platform in $Platforms) {
        $OS = $Platform.OS
        $Arch = $Platform.Arch
        $PlatformDir = "${OS}_${Arch}"
        $ArchiveName = "$BinaryName-$Version-$OS-$Arch"
        
        Write-Host "Packaging $ArchiveName..."
        
        $SourceDir = Join-Path $BuildDir $PlatformDir
        
        if ($OS -eq "windows") {
            $ArchivePath = Join-Path $DistDir "$ArchiveName.zip"
            Compress-Archive -Path (Join-Path $SourceDir "$BinaryName.exe") -DestinationPath $ArchivePath -Force
        } else {
            $ArchivePath = Join-Path $DistDir "$ArchiveName.tar.gz"
            $BinaryPath = Join-Path $SourceDir $BinaryName
            
            # Use tar if available, otherwise skip non-Windows archives
            if (Get-Command tar -ErrorAction SilentlyContinue) {
                & tar -czf $ArchivePath -C $SourceDir $BinaryName
            } else {
                Write-Host "  Warning: tar not found, skipping $OS archive" -ForegroundColor Yellow
                continue
            }
        }
        
        if (Test-Path $ArchivePath) {
            $Size = (Get-Item $ArchivePath).Length
            Write-Host "  Created: $ArchivePath ($([math]::Round($Size/1MB, 2)) MB)" -ForegroundColor Green
        }
    }
    
    Write-Host "`nRelease packages created in $DistDir/" -ForegroundColor Green
    Get-ChildItem $DistDir | Format-Table Name, Length, LastWriteTime
    
    # Generate checksums
    Write-Host "`nGenerating checksums..."
    $ChecksumFile = Join-Path $DistDir "checksums.txt"
    Get-ChildItem $DistDir -File | ForEach-Object {
        $Hash = (Get-FileHash $_.FullName -Algorithm SHA256).Hash.ToLower()
        "$Hash  $($_.Name)" | Out-File -FilePath $ChecksumFile -Append -Encoding utf8
    }
    Write-Host "Checksums saved to $ChecksumFile" -ForegroundColor Green
}

function Install-Binary {
    Write-Host "Installing $BinaryName..."
    
    & go install -ldflags $LDFlags $MainPackage
    
    if ($LASTEXITCODE -eq 0) {
        $GoPath = & go env GOPATH
        Write-Host "Installed to $GoPath\bin\$BinaryName.exe" -ForegroundColor Green
    } else {
        Write-Host "Installation failed!" -ForegroundColor Red
        exit 1
    }
}

function Clean-Build {
    Write-Host "Cleaning build artifacts..."
    
    if (Test-Path $BuildDir) {
        Remove-Item -Recurse -Force $BuildDir
    }
    
    if (Test-Path $DistDir) {
        Remove-Item -Recurse -Force $DistDir
    }
    
    Write-Host "Clean complete" -ForegroundColor Green
}

function Run-Tests {
    Write-Host "Running tests..."
    
    & go test -v -race -coverprofile=coverage.out ./...
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Tests passed!" -ForegroundColor Green
    } else {
        Write-Host "Tests failed!" -ForegroundColor Red
        exit 1
    }
}

function Show-Version {
    Write-Host "Version: $Version"
    Write-Host "Commit: $Commit"
    Write-Host "Build Date: $BuildDate"
}

# Main execution
if ($Help) {
    Show-Help
    exit 0
}

switch ($Target.ToLower()) {
    "build" { Build-Current }
    "build-all" { Build-All }
    "release" { Create-Release }
    "install" { Install-Binary }
    "clean" { Clean-Build }
    "test" { Run-Tests }
    "version" { Show-Version }
    default {
        Write-Host "Unknown target: $Target" -ForegroundColor Red
        Write-Host ""
        Show-Help
        exit 1
    }
}
