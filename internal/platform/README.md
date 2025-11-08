# Platform Adapter

平台适配器提供了跨平台抽象层，使 gx 能够在 Windows、Linux 和 macOS 上无缝运行。

## 功能特性

### 1. 平台检测
- 自动检测操作系统类型（Windows、Linux、macOS）
- 自动检测系统架构（amd64、arm64、386）

### 2. 路径处理
- 跨平台路径规范化
- 自动处理不同操作系统的路径分隔符
- 路径连接和清理

### 3. 文件权限管理
- Unix 系统（Linux/macOS）：基于文件模式位的权限检查和设置
- Windows 系统：基于文件扩展名的可执行性检查

### 4. 工具函数
- 获取用户主目录
- 获取配置和安装目录
- 目录创建和检查
- 文件存在性检查
- 平台支持验证

## 使用示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/gx/internal/platform"
)

func main() {
    // 创建平台适配器
    adapter := platform.NewAdapter()
    
    // 获取平台信息
    fmt.Printf("OS: %s\n", adapter.GetOS())
    fmt.Printf("Arch: %s\n", adapter.GetArch())
    
    // 路径处理
    path := adapter.JoinPath("home", "user", ".gx")
    normalized := adapter.NormalizePath(path)
    fmt.Printf("Path: %s\n", normalized)
    
    // 获取主目录
    home, _ := adapter.GetHomeDir()
    fmt.Printf("Home: %s\n", home)
    
    // 文件权限
    if err := adapter.MakeExecutable("/path/to/file"); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## 架构设计

### 接口定义
`pkg/interfaces/platform.go` 定义了 `PlatformAdapter` 接口，确保所有平台实现提供一致的 API。

### 核心实现
`internal/platform/adapter.go` 提供了平台适配器的核心实现，使用 Go 标准库的 `runtime` 和 `path/filepath` 包。

### 平台特定实现
- `permissions_unix.go`：Unix 系统（Linux/macOS）的文件权限实现
  - 使用构建标签 `//go:build linux || darwin`
  - 基于文件模式位（0111）检查和设置可执行权限
  
- `permissions_windows.go`：Windows 系统的文件权限实现
  - 使用构建标签 `//go:build windows`
  - 基于文件扩展名（.exe、.bat、.cmd、.com）判断可执行性

### 工具函数
`internal/platform/utils.go` 提供了额外的平台相关工具函数：
- 获取可执行文件和压缩包扩展名
- 配置和安装目录管理
- 文件系统操作辅助函数
- 平台支持验证

## 支持的平台

| 操作系统 | 架构 | 状态 |
|---------|------|------|
| Windows | amd64 | ✅ 支持 |
| Windows | 386 | ✅ 支持 |
| Linux | amd64 | ✅ 支持 |
| Linux | arm64 | ✅ 支持 |
| Linux | 386 | ✅ 支持 |
| macOS | amd64 | ✅ 支持 |
| macOS | arm64 | ✅ 支持 |

## 测试

运行平台适配器测试：

```bash
go test -v ./internal/platform/...
```

测试覆盖：
- 平台检测
- 路径处理和规范化
- 文件权限检查和设置
- 工具函数
- 跨平台兼容性

## 设计决策

### 1. 构建标签
使用 Go 的构建标签（build tags）实现平台特定代码，确保：
- 编译时只包含目标平台的代码
- 避免运行时平台检查的开销
- 代码清晰分离，易于维护

### 2. 文件权限抽象
不同操作系统的文件权限模型差异很大：
- Unix：基于用户/组/其他的读/写/执行权限位
- Windows：基于文件扩展名和 ACL

我们的实现提供了统一的接口，隐藏了这些差异。

### 3. 路径处理
使用 Go 标准库的 `path/filepath` 包：
- 自动处理不同操作系统的路径分隔符
- 提供跨平台的路径操作
- 避免手动字符串拼接导致的错误

## 未来改进

- [ ] 支持更多平台（FreeBSD、OpenBSD 等）
- [ ] 增强 Windows ACL 权限处理
- [ ] 添加符号链接支持
- [ ] 提供更详细的文件系统元数据访问
