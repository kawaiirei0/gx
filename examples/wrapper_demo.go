package main

import (
	"fmt"
	"os"

	"github.com/kawaiirei0/gx/internal/config"
	"github.com/kawaiirei0/gx/internal/downloader"
	"github.com/kawaiirei0/gx/internal/environment"
	"github.com/kawaiirei0/gx/internal/installer"
	"github.com/kawaiirei0/gx/internal/platform"
	"github.com/kawaiirei0/gx/internal/version"
	"github.com/kawaiirei0/gx/internal/wrapper"
)

func main() {
	// 创建平台适配器
	platformAdapter := platform.NewAdapter()
	fmt.Printf("Platform: %s/%s\n", platformAdapter.GetOS(), platformAdapter.GetArch())

	// 创建配置存储
	configStore, err := config.NewStore()
	if err != nil {
		fmt.Printf("Error creating config store: %v\n", err)
		os.Exit(1)
	}

	// 创建环境管理器
	envManager := environment.NewManager(platformAdapter)

	// 创建下载器
	downloaderInstance := downloader.NewDownloader()

	// 创建安装器
	installerInstance := installer.NewInstaller(platformAdapter)

	// 创建版本管理器
	versionManager := version.NewManager(
		configStore,
		platformAdapter,
		envManager,
		downloaderInstance,
		installerInstance,
	)

	// 创建 CLI Wrapper
	cliWrapper := wrapper.NewCLIWrapper(versionManager, platformAdapter)

	// 获取 Go 可执行文件路径
	fmt.Println("\n=== Getting Go Executable Path ===")
	goPath, err := cliWrapper.GetGoExecutable()
	if err != nil {
		fmt.Printf("Error getting go executable: %v\n", err)
		fmt.Println("Note: Make sure you have a Go version installed and activated using gx")
		os.Exit(1)
	}
	fmt.Printf("Go executable: %s\n", goPath)

	// 执行 go version 命令
	fmt.Println("\n=== Executing: go version ===")
	err = cliWrapper.Execute("version", []string{})
	if err != nil {
		if exitErr, ok := err.(*wrapper.ExitError); ok {
			fmt.Printf("Command exited with code: %d\n", exitErr.GetExitCode())
			os.Exit(exitErr.GetExitCode())
		} else {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

	// 执行 go env 命令
	fmt.Println("\n=== Executing: go env GOROOT GOPATH ===")
	err = cliWrapper.Execute("env", []string{"GOROOT", "GOPATH"})
	if err != nil {
		if exitErr, ok := err.(*wrapper.ExitError); ok {
			fmt.Printf("Command exited with code: %d\n", exitErr.GetExitCode())
			os.Exit(exitErr.GetExitCode())
		} else {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("\n=== CLI Wrapper Demo Complete ===")
}
