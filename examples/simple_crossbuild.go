package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/kawaiirei0/gx/pkg/constants"
)

func main() {
	fmt.Println("=== Simple Cross Build Test ===\n")

	// Test 1: Build for Linux AMD64
	fmt.Println("Test 1: Building for linux/amd64...")
	if err := crossBuild("linux", "amd64", "./examples/test_app", "test_linux_amd64"); err != nil {
		log.Printf("Failed: %v\n", err)
	} else {
		fmt.Println("✓ Success\n")
	}

	// Test 2: Build for Windows AMD64
	fmt.Println("Test 2: Building for windows/amd64...")
	if err := crossBuild("windows", "amd64", "./examples/test_app", "test_windows_amd64"); err != nil {
		log.Printf("Failed: %v\n", err)
	} else {
		fmt.Println("✓ Success\n")
	}

	// Test 3: Build for Darwin ARM64
	fmt.Println("Test 3: Building for darwin/arm64...")
	if err := crossBuild("darwin", "arm64", "./examples/test_app", "test_darwin_arm64"); err != nil {
		log.Printf("Failed: %v\n", err)
	} else {
		fmt.Println("✓ Success\n")
	}

	fmt.Println("All tests completed!")
}

func crossBuild(targetOS, targetArch, sourcePath, outputPath string) error {
	// Add platform-specific extension
	if targetOS == constants.OSWindows {
		outputPath += ".exe"
	}

	// Create build command
	cmd := exec.Command("go", "build", "-o", outputPath, sourcePath)

	// Set environment variables
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", targetOS))
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOARCH=%s", targetArch))

	// Set output streams
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute build
	if err := cmd.Run(); err != nil {
		return err
	}

	// Verify output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return fmt.Errorf("output file not created: %s", outputPath)
	}

	fmt.Printf("  Output: %s\n", outputPath)
	return nil
}
