package main

import (
	"fmt"
	"log"

	"github.com/yourusername/gx/internal/config"
	"github.com/yourusername/gx/internal/downloader"
	"github.com/yourusername/gx/internal/environment"
	"github.com/yourusername/gx/internal/installer"
	"github.com/yourusername/gx/internal/platform"
	"github.com/yourusername/gx/internal/version"
)

func main() {
	fmt.Println("=== gx Remote Version Query Demo ===\n")

	// 初始化依赖组件
	platformAdapter := platform.NewAdapter()
	
	configStore, err := config.NewStore()
	if err != nil {
		log.Fatalf("Failed to create config store: %v", err)
	}
	
	envManager := environment.NewManager(platformAdapter)
	downloaderInstance := downloader.NewDownloader()
	installerInstance := installer.NewInstaller(platformAdapter)

	// 创建版本管理器
	versionManager := version.NewManager(
		configStore,
		platformAdapter,
		envManager,
		downloaderInstance,
		installerInstance,
	)

	// 测试 GetLatest - 获取最新稳定版本
	fmt.Println("Fetching latest stable version...")
	latest, err := versionManager.GetLatest()
	if err != nil {
		log.Fatalf("Failed to get latest version: %v", err)
	}
	fmt.Printf("✓ Latest stable version: %s\n\n", latest)

	// 测试 ListAvailable - 获取所有可用版本
	fmt.Println("Fetching all available versions...")
	versions, err := versionManager.ListAvailable()
	if err != nil {
		log.Fatalf("Failed to list available versions: %v", err)
	}

	fmt.Printf("✓ Found %d available versions\n\n", len(versions))
	
	// 显示前 10 个版本
	fmt.Println("First 10 available versions:")
	displayCount := 10
	if len(versions) < displayCount {
		displayCount = len(versions)
	}
	
	for i := 0; i < displayCount; i++ {
		fmt.Printf("  %d. %s\n", i+1, versions[i])
	}

	if len(versions) > displayCount {
		fmt.Printf("  ... and %d more versions\n", len(versions)-displayCount)
	}

	fmt.Println("\n=== Demo completed successfully ===")
}
