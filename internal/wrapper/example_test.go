package wrapper_test

import (
	"fmt"
	"os"

	"github.com/yourusername/gx/internal/config"
	"github.com/yourusername/gx/internal/downloader"
	"github.com/yourusername/gx/internal/environment"
	"github.com/yourusername/gx/internal/installer"
	"github.com/yourusername/gx/internal/platform"
	"github.com/yourusername/gx/internal/version"
	"github.com/yourusername/gx/internal/wrapper"
)

// Example_executeGoVersion 演示如何执行 go version 命令
func Example_executeGoVersion() {
	// 创建依赖
	platformAdapter := platform.NewAdapter()
	configStore, err := config.NewStore()
	if err != nil {
		fmt.Printf("Error creating config store: %v\n", err)
		return
	}

	envManager := environment.NewManager(platformAdapter)
	downloaderInstance := downloader.NewDownloader()
	installerInstance := installer.NewInstaller(platformAdapter)

	versionManager := version.NewManager(
		configStore,
		platformAdapter,
		envManager,
		downloaderInstance,
		installerInstance,
	)

	// 创建 CLI Wrapper
	cliWrapper := wrapper.NewCLIWrapper(versionManager, platformAdapter)

	// 执行 go version 命令
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
}

// Example_executeGoRun 演示如何执行 go run 命令
func Example_executeGoRun() {
	// 创建依赖
	platformAdapter := platform.NewAdapter()
	configStore, err := config.NewStore()
	if err != nil {
		fmt.Printf("Error creating config store: %v\n", err)
		return
	}

	envManager := environment.NewManager(platformAdapter)
	downloaderInstance := downloader.NewDownloader()
	installerInstance := installer.NewInstaller(platformAdapter)

	versionManager := version.NewManager(
		configStore,
		platformAdapter,
		envManager,
		downloaderInstance,
		installerInstance,
	)

	// 创建 CLI Wrapper
	cliWrapper := wrapper.NewCLIWrapper(versionManager, platformAdapter)

	// 执行 go run main.go
	err = cliWrapper.Execute("run", []string{"main.go"})
	if err != nil {
		if exitErr, ok := err.(*wrapper.ExitError); ok {
			fmt.Printf("Command exited with code: %d\n", exitErr.GetExitCode())
			os.Exit(exitErr.GetExitCode())
		} else {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
}

// Example_getGoExecutable 演示如何获取 Go 可执行文件路径
func Example_getGoExecutable() {
	// 创建依赖
	platformAdapter := platform.NewAdapter()
	configStore, err := config.NewStore()
	if err != nil {
		fmt.Printf("Error creating config store: %v\n", err)
		return
	}

	envManager := environment.NewManager(platformAdapter)
	downloaderInstance := downloader.NewDownloader()
	installerInstance := installer.NewInstaller(platformAdapter)

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
	goPath, err := cliWrapper.GetGoExecutable()
	if err != nil {
		fmt.Printf("Error getting go executable: %v\n", err)
		return
	}

	fmt.Printf("Go executable path: %s\n", goPath)
}

// Example_executeWithArguments 演示如何执行带参数的 Go 命令
func Example_executeWithArguments() {
	// 创建依赖
	platformAdapter := platform.NewAdapter()
	configStore, err := config.NewStore()
	if err != nil {
		fmt.Printf("Error creating config store: %v\n", err)
		return
	}

	envManager := environment.NewManager(platformAdapter)
	downloaderInstance := downloader.NewDownloader()
	installerInstance := installer.NewInstaller(platformAdapter)

	versionManager := version.NewManager(
		configStore,
		platformAdapter,
		envManager,
		downloaderInstance,
		installerInstance,
	)

	// 创建 CLI Wrapper
	cliWrapper := wrapper.NewCLIWrapper(versionManager, platformAdapter)

	// 执行 go build -o output.exe main.go
	err = cliWrapper.Execute("build", []string{"-o", "output.exe", "main.go"})
	if err != nil {
		if exitErr, ok := err.(*wrapper.ExitError); ok {
			fmt.Printf("Command exited with code: %d\n", exitErr.GetExitCode())
			os.Exit(exitErr.GetExitCode())
		} else {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("Build successful!")
}
