# Downloader Package

## Overview

The downloader package provides functionality to download Go installation packages from the official Go website. It includes features for:

- Fetching available Go versions from the official API
- Generating download URLs for specific versions and platforms
- Downloading files with progress tracking
- Verifying downloaded files using SHA256 checksums

## Usage

### Creating a Downloader

```go
import "github.com/kawaiirei0/gx/internal/downloader"

dl := downloader.NewDownloader()
```

### Getting Download URL

```go
url, err := dl.GetDownloadURL("1.21.5", "linux", "amd64")
if err != nil {
    // Handle error
}
// url: https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
```

### Downloading a Go Version

```go
progress := func(downloaded, total int64) {
    percent := float64(downloaded) / float64(total) * 100
    fmt.Printf("\rDownloading: %.2f%%", percent)
}

err := dl.Download("1.21.5", "/tmp/go1.21.5.tar.gz", progress)
if err != nil {
    // Handle error
}
```

## Features

### Automatic Version Normalization

The downloader automatically adds the "go" prefix to version numbers if not present:

```go
// Both work the same
dl.GetDownloadURL("1.21.5", "linux", "amd64")
dl.GetDownloadURL("go1.21.5", "linux", "amd64")
```

### SHA256 Verification

All downloaded files are automatically verified against the official SHA256 checksums from the Go API. If verification fails, the download is rejected and the temporary file is cleaned up.

### Progress Tracking

The downloader supports progress callbacks that receive the number of bytes downloaded and the total file size:

```go
type ProgressCallback func(downloaded int64, total int64)
```

### Atomic Downloads

Files are downloaded to a temporary location first, verified, and only then moved to the final destination. This ensures that partial or corrupted downloads don't leave invalid files.

## Error Handling

The downloader uses custom error types from the `pkg/errors` package:

- `ErrNetworkError`: Network-related errors (connection failures, timeouts)
- `ErrVersionNotFound`: Requested version or platform not available
- `ErrDownloadFailed`: Download process failed
- `ErrChecksumMismatch`: SHA256 verification failed

## Implementation Details

### API Integration

The downloader queries the official Go versions API at `https://go.dev/dl/?mode=json` to get:
- Available versions
- Download URLs
- File sizes
- SHA256 checksums

### Platform Support

Supports all platforms available from the official Go downloads:
- Windows (amd64, 386)
- Linux (amd64, arm64, 386)
- macOS/Darwin (amd64, arm64)

### HTTP Client Configuration

- Default timeout: 30 minutes (suitable for large downloads on slow connections)
- Follows redirects automatically
- Supports HTTP/2 when available
