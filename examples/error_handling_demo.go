package main

import (
	"fmt"
	"os"

	"github.com/yourusername/gx/pkg/errors"
)

func main() {
	fmt.Println("=== Error Handling Demo ===\n")

	// Demo 1: Basic error with context
	fmt.Println("1. Basic Error with Context:")
	err1 := errors.ErrVersionNotFound.
		WithMessage("version go1.21.0 not found").
		WithContext("requested_version", "go1.21.0").
		WithContext("available_versions", []string{"go1.20.0", "go1.19.0"})
	
	reporter := errors.NewErrorReporter(true)
	fmt.Println(reporter.Report(err1))
	fmt.Println()

	// Demo 2: Error wrapping
	fmt.Println("2. Error Wrapping:")
	baseErr := fmt.Errorf("network timeout")
	err2 := errors.Wrap(baseErr, "DOWNLOAD_FAILED", "failed to download Go archive").
		WithContext("url", "https://go.dev/dl/go1.21.0.tar.gz").
		WithContext("attempt", 3)
	
	fmt.Println(reporter.Report(err2))
	fmt.Println()

	// Demo 3: Recoverable error
	fmt.Println("3. Recoverable Error:")
	err3 := errors.ErrConfigCorrupted.
		WithMessage("config file is corrupted but backup exists").
		WithContext("config_path", "~/.gx/config.json").
		WithContext("backup_path", "~/.gx/config.json.backup").
		AsRecoverable()
	
	fmt.Println(reporter.ReportWithRecovery(err3))
	fmt.Println()

	// Demo 4: Recovery Manager
	fmt.Println("4. Recovery Manager Demo:")
	demoRecoveryManager()
	fmt.Println()

	// Demo 5: Error chain
	fmt.Println("5. Error Chain:")
	err5 := errors.ErrInstallFailed.
		WithCause(errors.ErrDownloadFailed.WithCause(baseErr)).
		WithMessage("installation failed due to download error")
	
	fmt.Println("Error chain:")
	fmt.Println(errors.FormatErrorChain(err5))
	fmt.Println()

	// Demo 6: Error type checking
	fmt.Println("6. Error Type Checking:")
	if errors.IsType(err1, errors.ErrVersionNotFound) {
		fmt.Println("✓ Error is of type VERSION_NOT_FOUND")
	}
	
	if errors.IsRecoverableError(err3) {
		fmt.Println("✓ Error is recoverable")
	}
	fmt.Println()
}

func demoRecoveryManager() {
	recovery := errors.NewRecoveryManager()
	
	// Create a temp directory
	tmpDir, err := errors.CreateTempDir(recovery, "demo-")
	if err != nil {
		fmt.Printf("Failed to create temp dir: %v\n", err)
		return
	}
	
	fmt.Printf("Created temp directory: %s\n", tmpDir)
	
	// Create a temp file
	tmpFile := tmpDir + "/test.txt"
	os.WriteFile(tmpFile, []byte("test content"), 0644)
	fmt.Printf("Created temp file: %s\n", tmpFile)
	
	// Register cleanup
	errors.EnsureFileCleanup(recovery, tmpFile)
	
	// Simulate operation
	fmt.Println("Simulating operation...")
	
	// Cleanup
	fmt.Println("Executing cleanup...")
	if err := recovery.Cleanup(); err != nil {
		fmt.Printf("Cleanup failed: %v\n", err)
	} else {
		fmt.Println("✓ Cleanup successful")
	}
	
	// Verify cleanup
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		fmt.Println("✓ Temp directory removed")
	}
}
