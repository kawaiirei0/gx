package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yourusername/gx/internal/config"
	"github.com/yourusername/gx/internal/crossbuilder"
	"github.com/yourusername/gx/internal/platform"
	"github.com/yourusername/gx/internal/version"
	"github.com/yourusername/gx/pkg/interfaces"
)

func main() {
	fmt.Println("=== Cross Builder Demo ===\n")

	// 创建平台适配器
	platformAdapter := platform.NewAdapter()
	fmt.Printf("Current platform: %s/%s\n\n", platformAdapter.GetOS(), platformAdapter.GetArch())

	// 创建配置存储
	configStore, err := config.NewStore()
	if err != nil {
		log.Fatalf("Failed to create config store: %v\n", err)
	}

	// 创建版本管理器
	versionManager := version.NewManager(configStore, platformAdapter, nil, nil, nil)

	// 创建跨平台构建器
	builder := crossbuilder.NewCrossBuilder(versionManager, platformAdapter)

	// 显示支持的平台
	fmt.Println("Supported platforms:")
	platforms := builder.GetSupportedPlatforms()
	for _, p := range platforms {
		fmt.Printf("  - %s/%s\n", p.OS, p.Arch)
	}
	fmt.Println()

	// 验证平台组合
	fmt.Println("Validating platform combinations:")
	testPlatforms := []struct {
		os   string
		arch string
	}{
		{"linux", "amd64"},
		{"windows", "amd64"},
		{"darwin", "arm64"},
		{"freebsd", "amd64"}, // 不支持的平台
	}

	for _, tp := range testPlatforms {
		err := builder.ValidatePlatform(tp.os, tp.arch)
		if err != nil {
			fmt.Printf("  ✗ %s/%s: %v\n", tp.os, tp.arch, err)
		} else {
			fmt.Printf("  ✓ %s/%s: supported\n", tp.os, tp.arch)
		}
	}
	fmt.Println()

	// 如果提供了命令行参数，执行实际构建
	if len(os.Args) > 1 {
		targetOS := "linux"
		targetArch := "amd64"
		sourcePath := "."
		outputPath := "demo_output"

		if len(os.Args) > 1 {
			targetOS = os.Args[1]
		}
		if len(os.Args) > 2 {
			targetArch = os.Args[2]
		}
		if len(os.Args) > 3 {
			sourcePath = os.Args[3]
		}
		if len(os.Args) > 4 {
			outputPath = os.Args[4]
		}

		fmt.Printf("Building for %s/%s...\n", targetOS, targetArch)
		config := interfaces.BuildConfig{
			SourcePath: sourcePath,
			OutputPath: outputPath,
			TargetOS:   targetOS,
			TargetArch: targetArch,
			BuildFlags: []string{"-v"},
			LDFlags:    "-s -w",
		}

		if err := builder.Build(config); err != nil {
			log.Fatalf("Build failed: %v\n", err)
		}
	} else {
		fmt.Println("Usage: crossbuilder_demo [target_os] [target_arch] [source_path] [output_path]")
		fmt.Println("Example: crossbuilder_demo linux amd64 . myapp")
	}
}
