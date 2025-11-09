package main

import (
	"fmt"

	"github.com/kawaiirei0/gx/internal/config"
	"github.com/kawaiirei0/gx/internal/downloader"
	"github.com/kawaiirei0/gx/internal/installer"
	"github.com/kawaiirei0/gx/internal/platform"
)

func main() {
	fmt.Println("=== gx Downloader Demo ===\n")

	// 创建必要的组件
	platformAdapter := platform.NewAdapter()
	configStore, err := config.NewStore()
	if err != nil {
		fmt.Printf("Error creating config store: %v\n", err)
		return
	}

	dl := downloader.NewDownloader()
	_ = installer.NewInstaller(platformAdapter) // Create but don't use in demo

	// 演示 1: 获取下载 URL
	fmt.Println("1. Getting download URL for Go 1.21.5...")
	url, err := dl.GetDownloadURL("1.21.5", platformAdapter.GetOS(), platformAdapter.GetArch())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("   URL: %s\n\n", url)
	}

	// 演示 2: 列出已安装的版本
	fmt.Println("2. Detecting installed Go versions...")
	
	// 注意: 这里需要 EnvironmentManager，但为了演示简化，我们只展示概念
	// 在实际使用中，需要完整的依赖注入
	fmt.Println("   (Requires EnvironmentManager - see full implementation)\n")

	// 演示 3: 下载和安装流程（模拟）
	fmt.Println("3. Installation process overview:")
	fmt.Println("   a) Download archive with progress tracking")
	fmt.Println("   b) Verify SHA256 checksum")
	fmt.Println("   c) Extract to installation directory")
	fmt.Println("   d) Verify installation (check go executable)")
	fmt.Println("   e) Update configuration\n")

	// 演示 4: 进度回调示例
	fmt.Println("4. Progress callback example:")
	progressCallback := func(downloaded, total int64) {
		if total > 0 {
			percent := float64(downloaded) / float64(total) * 100
			fmt.Printf("\r   Downloading: %.2f%% (%d/%d bytes)", percent, downloaded, total)
		}
	}
	fmt.Printf("   Callback function: %T\n\n", progressCallback)

	// 演示 5: 验证功能
	fmt.Println("5. Installation verification:")
	fmt.Println("   - Checks bin/ directory exists")
	fmt.Println("   - Verifies go executable is present")
	fmt.Println("   - Sets executable permissions (Unix)")
	fmt.Println("   - Runs 'go version' to validate")
	fmt.Println("   - Compares version output with expected\n")

	// 演示 6: 完整的版本管理器使用（需要所有依赖）
	fmt.Println("6. Complete Version Manager usage:")
	fmt.Println("   To use the full version manager, initialize with:")
	fmt.Println("   - ConfigStore: manages configuration")
	fmt.Println("   - PlatformAdapter: handles OS-specific operations")
	fmt.Println("   - EnvironmentManager: manages environment variables")
	fmt.Println("   - Downloader: downloads Go distributions")
	fmt.Println("   - Installer: extracts and verifies installations")
	fmt.Println()
	fmt.Println("   Example:")
	fmt.Println("   vm := version.NewManager(configStore, platform, envMgr, downloader, installer)")
	fmt.Println("   err := vm.Install(\"1.21.5\", progressCallback)")

	// 显示平台信息
	fmt.Println("\n=== Platform Information ===")
	fmt.Printf("OS: %s\n", platformAdapter.GetOS())
	fmt.Printf("Architecture: %s\n", platformAdapter.GetArch())
	fmt.Printf("Path Separator: %s\n", platformAdapter.PathSeparator())

	// 显示配置信息
	cfg, err := configStore.Load()
	if err != nil {
		fmt.Printf("\nConfig: Error loading - %v\n", err)
	} else {
		fmt.Printf("\nConfig Path: %s\n", cfg.InstallPath)
		fmt.Printf("Active Version: %s\n", cfg.ActiveVersion)
		fmt.Printf("Installed Versions: %d\n", len(cfg.Versions))
	}

	fmt.Println("\n=== Demo Complete ===")
}

// demonstrateFullInstallation shows how to use the complete version manager
// Note: Requires EnvironmentManager implementation to run
// Uncomment when EnvironmentManager is implemented
/*
func demonstrateFullInstallation() {
	// Create all necessary components
	platformAdapter := platform.NewAdapter()
	configStore, _ := config.NewStore()
	dl := downloader.NewDownloader()
	inst := installer.NewInstaller(platformAdapter)
	
	// Note: envManager needs to be implemented
	// envManager := environment.NewManager(platformAdapter)
	
	// vm := version.NewManager(configStore, platformAdapter, envManager, dl, inst)

	// Define progress callback
	progress := func(downloaded, total int64) {
		if total > 0 {
			percent := float64(downloaded) / float64(total) * 100
			fmt.Printf("\rDownloading: %.2f%%", percent)
		}
	}

	// Install Go 1.21.5
	// err := vm.Install("1.21.5", progress)
	// if err != nil {
	//     fmt.Printf("Installation failed: %v\n", err)
	//     return
	// }

	fmt.Println("\nInstallation completed successfully!")
}
*/
