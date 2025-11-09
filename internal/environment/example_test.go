package environment_test

import (
	"fmt"
	"os"

	"github.com/kawaiirei0/gx/internal/environment"
	"github.com/kawaiirei0/gx/internal/platform"
)

// ExampleNewManager 演示如何创建环境管理器
func ExampleNewManager() {
	// 创建平台适配器
	platformAdapter := platform.NewAdapter()

	// 创建环境管理器
	envManager := environment.NewManager(platformAdapter)

	fmt.Printf("Environment manager created for platform: %s\n", platformAdapter.GetOS())
	// Output: Environment manager created for platform: windows
}

// ExampleManager_SetGoRoot 演示如何设置 GOROOT
func ExampleManager_SetGoRoot() {
	platformAdapter := platform.NewAdapter()
	envManager := environment.NewManager(platformAdapter)

	// 设置 GOROOT（示例路径）
	goRoot := "/home/user/.gx/versions/go1.21.5"
	
	// 注意：这个示例不会实际执行，因为路径可能不存在
	if err := envManager.SetGoRoot(goRoot); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("GOROOT set successfully")
}

// ExampleManager_UpdatePath 演示如何更新 PATH
func ExampleManager_UpdatePath() {
	platformAdapter := platform.NewAdapter()
	envManager := environment.NewManager(platformAdapter)

	goRoot := "/home/user/.gx/versions/go1.21.5"

	// 更新 PATH
	if err := envManager.UpdatePath(goRoot); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("PATH updated successfully")
}

// ExampleManager_GetGoRoot 演示如何获取当前 GOROOT
func ExampleManager_GetGoRoot() {
	platformAdapter := platform.NewAdapter()
	envManager := environment.NewManager(platformAdapter)

	// 先设置一个测试值
	os.Setenv("GOROOT", "/usr/local/go")

	// 获取 GOROOT
	goRoot, err := envManager.GetGoRoot()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Current GOROOT: %s\n", goRoot)
	// Output: Current GOROOT: /usr/local/go
}

// ExampleManager_Backup 演示如何备份环境变量
func ExampleManager_Backup() {
	platformAdapter := platform.NewAdapter()
	envManager := environment.NewManager(platformAdapter)

	// 设置一些测试环境变量
	os.Setenv("GOROOT", "/usr/local/go")
	os.Setenv("GOPATH", "/home/user/go")

	// 备份环境变量
	if err := envManager.Backup(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Environment variables backed up successfully")
}

// ExampleManager_Restore 演示如何恢复环境变量
func ExampleManager_Restore() {
	platformAdapter := platform.NewAdapter()
	envManager := environment.NewManager(platformAdapter)

	// 恢复环境变量
	if err := envManager.Restore(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Environment variables restored successfully")
}
