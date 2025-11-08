# Environment Manager

环境变量管理器，负责跨平台的环境变量设置和持久化。

## 功能

- 设置和获取 GOROOT、GOPATH 环境变量
- 更新 PATH 环境变量，添加 Go bin 目录
- 跨平台支持（Windows、Linux、macOS）
- 环境变量备份和恢复
- 持久化环境变量配置

## 平台特定实现

### Windows

- 使用 `setx` 命令设置用户环境变量
- 环境变量写入注册表，重启后仍然有效
- 支持通过 PowerShell 设置系统级环境变量（需要管理员权限）

### Linux/macOS

- 自动检测使用的 shell（bash、zsh）
- 更新相应的 shell 配置文件（.bashrc、.zshrc、.profile 等）
- 使用标记注释管理 gx 设置的环境变量
- 支持多个配置文件同时更新

## 使用示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/gx/internal/environment"
    "github.com/yourusername/gx/internal/platform"
)

func main() {
    // 创建平台适配器
    platformAdapter := platform.NewAdapter()
    
    // 创建环境管理器
    envManager := environment.NewManager(platformAdapter)
    
    // 设置 GOROOT
    goRoot := "/home/user/.gx/versions/go1.21.5"
    if err := envManager.SetGoRoot(goRoot); err != nil {
        fmt.Printf("Failed to set GOROOT: %v\n", err)
        return
    }
    
    // 更新 PATH
    if err := envManager.UpdatePath(goRoot); err != nil {
        fmt.Printf("Failed to update PATH: %v\n", err)
        return
    }
    
    // 获取当前 GOROOT
    currentGoRoot, err := envManager.GetGoRoot()
    if err != nil {
        fmt.Printf("Failed to get GOROOT: %v\n", err)
        return
    }
    fmt.Printf("Current GOROOT: %s\n", currentGoRoot)
    
    // 备份环境变量
    if err := envManager.Backup(); err != nil {
        fmt.Printf("Failed to backup environment: %v\n", err)
        return
    }
    
    fmt.Println("Environment configured successfully!")
}
```

## 环境变量持久化

### Windows

环境变量通过 `setx` 命令写入用户注册表：
- 位置：`HKEY_CURRENT_USER\Environment`
- 立即生效于新启动的进程
- 当前进程需要重新读取环境变量

### Linux/macOS

环境变量写入 shell 配置文件，格式如下：

```bash
# gx managed GOROOT
export GOROOT="/home/user/.gx/versions/go1.21.5"

# gx managed PATH
export PATH="/home/user/.gx/versions/go1.21.5/bin:$PATH"
```

- 使用标记注释 `# gx managed <VAR>` 标识 gx 管理的环境变量
- 支持更新和删除已有的 gx 环境变量设置
- 新的 shell 会话会自动加载更新后的环境变量

## 备份和恢复

环境变量备份保存在 `~/.gx/env_backup.json`：

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "variables": {
    "GOROOT": "/usr/local/go",
    "GOPATH": "/home/user/go",
    "PATH": "/usr/local/bin:/usr/bin:/bin"
  }
}
```

可以使用 `Restore()` 方法恢复到备份的状态。

## 注意事项

1. **Windows 权限**：设置系统级环境变量需要管理员权限
2. **Shell 重载**：Linux/macOS 上修改配置文件后，需要重新加载 shell 或启动新的 shell 会话
3. **PATH 更新**：更新 PATH 时会自动移除旧的 gx 管理的 Go bin 路径
4. **备份策略**：建议在重要操作前进行备份

## 错误处理

所有方法都返回详细的错误信息，包括：
- `ErrInvalidInput`：输入参数无效
- `ErrOperationFailed`：操作失败（包含底层错误）
- `ErrNotFound`：环境变量或备份文件不存在
- `ErrPlatformNotSupported`：不支持的平台
