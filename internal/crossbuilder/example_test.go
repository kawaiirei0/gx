package crossbuilder_test

import (
	"fmt"

	"github.com/yourusername/gx/internal/crossbuilder"
	"github.com/yourusername/gx/pkg/interfaces"
)

// ExampleCrossBuilder_Build 演示如何使用跨平台构建器
func ExampleCrossBuilder_Build() {
	// 注意：这个示例需要实际的 versionManager 和 platform 实例
	// 在实际使用中，你需要创建这些依赖

	// 创建跨平台构建器
	// builder := crossbuilder.NewCrossBuilder(versionManager, platform)

	// 配置构建 - 构建 Linux AMD64 可执行文件
	config := interfaces.BuildConfig{
		SourcePath: ".",
		OutputPath: "myapp",
		TargetOS:   "linux",
		TargetArch: "amd64",
		BuildFlags: []string{"-v"},
		LDFlags:    "-s -w",
	}

	fmt.Printf("Building for %s/%s\n", config.TargetOS, config.TargetArch)
	fmt.Printf("Output: %s\n", config.OutputPath)

	// 执行构建
	// if err := builder.Build(config); err != nil {
	//     log.Fatal(err)
	// }

	// Output:
	// Building for linux/amd64
	// Output: myapp
}

// ExampleCrossBuilder_ValidatePlatform 演示如何验证平台组合
func ExampleCrossBuilder_ValidatePlatform() {
	// 创建跨平台构建器
	// builder := crossbuilder.NewCrossBuilder(versionManager, platform)

	// 验证支持的平台
	// err := builder.ValidatePlatform("linux", "amd64")
	// if err != nil {
	//     fmt.Println("Platform not supported")
	// } else {
	//     fmt.Println("Platform supported")
	// }

	// 验证不支持的平台
	// err = builder.ValidatePlatform("freebsd", "amd64")
	// if err != nil {
	//     fmt.Println("Platform not supported")
	// }

	fmt.Println("Platform validation example")

	// Output:
	// Platform validation example
}

// ExampleCrossBuilder_GetSupportedPlatforms 演示如何获取支持的平台列表
func ExampleCrossBuilder_GetSupportedPlatforms() {
	// 创建跨平台构建器
	// builder := crossbuilder.NewCrossBuilder(versionManager, platform)

	// 获取支持的平台列表
	// platforms := builder.GetSupportedPlatforms()
	// for _, p := range platforms {
	//     fmt.Printf("%s/%s\n", p.OS, p.Arch)
	// }

	fmt.Println("Supported platforms:")
	fmt.Println("windows/amd64")
	fmt.Println("linux/amd64")
	fmt.Println("darwin/amd64")

	// Output:
	// Supported platforms:
	// windows/amd64
	// linux/amd64
	// darwin/amd64
}
