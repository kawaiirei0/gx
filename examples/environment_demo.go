package main

import (
	"fmt"
	"os"

	"github.com/kawaiirei0/gx/internal/environment"
	"github.com/kawaiirei0/gx/internal/platform"
)

func main() {
	fmt.Println("=== gx Environment Manager Demo ===\n")

	// 创建平台适配器
	platformAdapter := platform.NewAdapter()
	fmt.Printf("Platform: %s/%s\n", platformAdapter.GetOS(), platformAdapter.GetArch())

	// 创建环境管理器
	envManager := environment.NewManager(platformAdapter)
	fmt.Println("Environment manager created\n")

	// 演示 1: 获取当前环境变量
	fmt.Println("--- Current Environment ---")
	if goRoot, err := envManager.GetGoRoot(); err == nil {
		fmt.Printf("GOROOT: %s\n", goRoot)
	} else {
		fmt.Printf("GOROOT: not set (%v)\n", err)
	}

	if goPath, err := envManager.GetGoPath(); err == nil {
		fmt.Printf("GOPATH: %s\n", goPath)
	} else {
		fmt.Printf("GOPATH: not set (%v)\n", err)
	}
	fmt.Println()

	// 演示 2: 备份当前环境变量
	fmt.Println("--- Backup Environment ---")
	if err := envManager.Backup(); err != nil {
		fmt.Printf("Failed to backup: %v\n", err)
	} else {
		fmt.Println("Environment backed up successfully")
	}
	fmt.Println()

	// 演示 3: 设置新的 GOROOT（示例）
	// 注意：这只是演示，实际使用时需要确保路径存在
	homeDir, _ := platformAdapter.GetHomeDir()
	exampleGoRoot := platformAdapter.JoinPath(homeDir, ".gx", "versions", "go1.21.5")
	
	fmt.Println("--- Set GOROOT (Demo) ---")
	fmt.Printf("Would set GOROOT to: %s\n", exampleGoRoot)
	fmt.Println("(Skipped in demo to avoid modifying your system)")
	fmt.Println()

	// 演示 4: 更新 PATH（示例）
	fmt.Println("--- Update PATH (Demo) ---")
	fmt.Printf("Would add to PATH: %s/bin\n", exampleGoRoot)
	fmt.Println("(Skipped in demo to avoid modifying your system)")
	fmt.Println()

	// 演示 5: 平台特定信息
	fmt.Println("--- Platform-Specific Details ---")
	switch platformAdapter.GetOS() {
	case "windows":
		fmt.Println("On Windows:")
		fmt.Println("  - Environment variables are set via 'setx' command")
		fmt.Println("  - Changes are written to user registry")
		fmt.Println("  - New processes will see the updated values")
	case "linux", "darwin":
		fmt.Println("On Unix (Linux/macOS):")
		fmt.Println("  - Environment variables are written to shell RC files")
		shell := os.Getenv("SHELL")
		if shell != "" {
			fmt.Printf("  - Detected shell: %s\n", shell)
		}
		fmt.Println("  - New shell sessions will load the updated values")
		fmt.Println("  - Run 'source ~/.bashrc' (or ~/.zshrc) to reload in current shell")
	}
	fmt.Println()

	fmt.Println("=== Demo Complete ===")
	fmt.Println("\nTo actually use the environment manager:")
	fmt.Println("1. Install a Go version using gx")
	fmt.Println("2. The environment manager will automatically configure GOROOT and PATH")
	fmt.Println("3. Use 'gx use <version>' to switch between installed versions")
}
