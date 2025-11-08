# CLI Wrapper

CLI Wrapper 组件负责包装和转发 Go 原生命令，提供统一的命令执行接口。

## 功能

- 执行 Go 命令（run, build, test 等）
- 透传标准输入、输出、错误流
- 保留原始命令的退出码
- 自动获取当前激活版本的 Go 可执行文件路径

## 接口

```go
type CLIWrapper interface {
    // Execute 执行 Go 命令
    Execute(command string, args []string) error
    
    // GetGoExecutable 获取当前使用的 Go 可执行文件路径
    GetGoExecutable() (string, error)
}
```

## 使用示例

```go
package main

import (
    "fmt"
    "github.com/yourusername/gx/internal/wrapper"
    "github.com/yourusername/gx/internal/version"
    "github.com/yourusername/gx/internal/platform"
)

func main() {
    // 创建依赖
    platformAdapter := platform.NewAdapter()
    versionManager := version.NewManager(...)
    
    // 创建 CLI Wrapper
    cliWrapper := wrapper.NewCLIWrapper(versionManager, platformAdapter)
    
    // 执行 go run 命令
    err := cliWrapper.Execute("run", []string{"main.go", "--arg1", "value1"})
    if err != nil {
        if exitErr, ok := err.(*wrapper.ExitError); ok {
            fmt.Printf("Command exited with code: %d\n", exitErr.GetExitCode())
        } else {
            fmt.Printf("Error: %v\n", err)
        }
    }
}
```

## 实现细节

### 命令执行流程

1. 通过 `GetGoExecutable()` 获取当前激活版本的 Go 可执行文件路径
2. 使用 `os/exec` 包创建命令
3. 将标准输入、输出、错误流连接到当前进程
4. 执行命令并保留退出码

### 退出码处理

CLI Wrapper 使用自定义的 `ExitError` 类型来保留命令的退出码：

```go
type ExitError struct {
    ExitCode int
    Err      error
}
```

这样可以让调用者区分：
- 命令执行失败（如编译错误）：返回 `ExitError`，包含退出码
- 系统错误（如命令未找到）：返回普通错误

### 路径验证

`GetGoExecutable()` 方法会进行以下验证：

1. 检查当前是否有激活的 Go 版本
2. 验证版本路径是否存在
3. 检查 Go 可执行文件是否存在
4. 在 Unix 系统上验证文件是否具有可执行权限

## 错误处理

可能返回的错误类型：

- `ErrVersionNotFound`: 未找到激活的 Go 版本
- `ErrNotFound`: Go 可执行文件不存在
- `ErrPermissionDenied`: 文件没有可执行权限
- `ErrOperationFailed`: 命令执行失败
- `ExitError`: 命令执行完成但返回非零退出码

## 支持的命令

CLI Wrapper 可以执行任何 Go 命令，常用的包括：

- `run`: 运行 Go 程序
- `build`: 编译 Go 程序
- `test`: 运行测试
- `mod`: 模块管理
- `get`: 下载依赖
- `install`: 安装包
- `fmt`: 格式化代码
- `vet`: 代码检查

所有参数和标志都会原样传递给 Go 命令。
