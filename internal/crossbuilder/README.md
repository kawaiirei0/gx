# Cross Builder

跨平台构建器，用于在一个操作系统上编译出其他操作系统的可执行文件。

## 功能

- 支持跨平台编译（Windows、Linux、macOS）
- 自动设置 GOOS 和 GOARCH 环境变量
- 验证平台组合是否支持
- 自动添加平台特定的文件扩展名（Windows 的 .exe）
- 支持自定义构建标志和链接器标志

## 支持的平台组合

- windows/amd64
- windows/386
- linux/amd64
- linux/arm64
- linux/386
- darwin/amd64
- darwin/arm64

## 使用示例

```go
package main

import (
    "github.com/kawaiirei0/gx/internal/crossbuilder"
    "github.com/kawaiirei0/gx/pkg/interfaces"
)

func main() {
    // 创建跨平台构建器
    builder := crossbuilder.NewCrossBuilder(versionManager, platform)
    
    // 配置构建
    config := interfaces.BuildConfig{
        SourcePath: ".",
        OutputPath: "myapp",
        TargetOS:   "linux",
        TargetArch: "amd64",
        BuildFlags: []string{"-v"},
        LDFlags:    "-s -w",
    }
    
    // 执行构建
    if err := builder.Build(config); err != nil {
        log.Fatal(err)
    }
}
```

## 实现细节

### 环境变量设置

构建器会自动设置以下环境变量：
- `GOOS`: 目标操作系统
- `GOARCH`: 目标架构

### 输出文件处理

- 对于 Windows 目标平台，自动添加 `.exe` 扩展名
- 对于非 Windows 平台，自动移除 `.exe` 扩展名（如果存在）

### 错误处理

- 验证源代码路径是否存在
- 验证平台组合是否支持
- 验证 Go 可执行文件是否存在
- 捕获并包装构建错误

## 相关需求

- Requirement 7.1: 设置 GOOS 和 GOARCH 环境变量
- Requirement 7.2: 调用 go build 并传递正确的平台参数
- Requirement 7.3: 支持指定的平台组合
- Requirement 7.4: 在指定输出目录生成目标平台的可执行文件
- Requirement 7.5: 显示目标平台信息和输出文件路径
