# gx Quick Installation Script for Windows
# This script builds gx and runs init-install

$ErrorActionPreference = "Stop"

Write-Host "╔════════════════════════════════════════╗" -ForegroundColor Blue
Write-Host "║     gx - Go Version Manager            ║" -ForegroundColor Blue
Write-Host "║     Quick Installation Script          ║" -ForegroundColor Blue
Write-Host "╚════════════════════════════════════════╝" -ForegroundColor Blue
Write-Host ""

# Check if Go is installed
try {
    $goVersion = & go version
    Write-Host "✓ Found Go: $goVersion" -ForegroundColor Green
    Write-Host ""
} catch {
    Write-Host "Error: Go is not installed" -ForegroundColor Red
    Write-Host "Please install Go first: https://golang.org/dl/"
    exit 1
}

# Build gx
Write-Host "Building gx..." -ForegroundColor Yellow
if (Test-Path "build.ps1") {
    & .\build.ps1 build
} else {
    & go build -o build\gx.exe .\cmd\gx
}

if (-not (Test-Path "build\gx.exe")) {
    Write-Host "Error: Build failed" -ForegroundColor Red
    exit 1
}

Write-Host "✓ Build complete" -ForegroundColor Green
Write-Host ""

# Run init-install
Write-Host "Running installation..." -ForegroundColor Yellow
Write-Host ""
& .\build\gx.exe init-install

Write-Host ""
Write-Host "╔════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║     Installation Complete!             ║" -ForegroundColor Green
Write-Host "╚════════════════════════════════════════╝" -ForegroundColor Green
