package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yourusername/gx/internal/config"
)

func main() {
	fmt.Println("=== Config Store Demo ===\n")

	// 创建配置存储
	store, err := config.NewStore()
	if err != nil {
		log.Fatalf("Failed to create config store: %v", err)
	}

	// 确保配置目录存在
	fmt.Println("1. Ensuring config directory exists...")
	if err := store.EnsureConfigDir(); err != nil {
		log.Fatalf("Failed to ensure config directory: %v", err)
	}
	fmt.Println("   ✓ Config directory ready")

	// 加载配置（如果不存在会返回默认配置）
	fmt.Println("\n2. Loading configuration...")
	cfg, err := store.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Printf("   ✓ Config loaded\n")
	fmt.Printf("   - Active Version: %s\n", cfg.ActiveVersion)
	fmt.Printf("   - Install Path: %s\n", cfg.InstallPath)
	fmt.Printf("   - Versions: %d installed\n", len(cfg.Versions))

	// 修改配置
	fmt.Println("\n3. Updating configuration...")
	cfg.ActiveVersion = "1.21.5"
	cfg.Versions["1.21.5"] = cfg.InstallPath + "/go1.21.5"
	cfg.Versions["1.22.0"] = cfg.InstallPath + "/go1.22.0"
	cfg.LastUpdateCheck = time.Now()

	// 保存配置
	fmt.Println("\n4. Saving configuration...")
	if err := store.Save(cfg); err != nil {
		log.Fatalf("Failed to save config: %v", err)
	}
	fmt.Println("   ✓ Config saved successfully")

	// 重新加载验证
	fmt.Println("\n5. Reloading to verify...")
	reloadedCfg, err := store.Load()
	if err != nil {
		log.Fatalf("Failed to reload config: %v", err)
	}
	fmt.Printf("   ✓ Config reloaded\n")
	fmt.Printf("   - Active Version: %s\n", reloadedCfg.ActiveVersion)
	fmt.Printf("   - Installed Versions:\n")
	for version, path := range reloadedCfg.Versions {
		fmt.Printf("     • %s -> %s\n", version, path)
	}
	fmt.Printf("   - Last Update Check: %s\n", reloadedCfg.LastUpdateCheck.Format(time.RFC3339))

	fmt.Println("\n=== Demo Complete ===")
}
