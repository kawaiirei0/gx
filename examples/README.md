# gx 示例程序

本目录包含 gx 各个组件的示例程序，帮助开发者理解如何使用 gx 的内部 API。

## 目录

- [配置管理](#配置管理)
- [下载器](#下载器)
- [环境管理](#环境管理)
- [错误处理](#错误处理)
- [跨平台构建](#跨平台构建)
- [CLI 包装器](#cli-包装器)
- [版本管理](#版本管理)

## 运行示例

所有示例都可以使用 `go run` 命令运行：

```bash
# 运行配置管理示例
go run examples/config_demo.go

# 运行下载器示例
go run examples/downloader_demo.go

# 运行所有示例
go run examples/*.go
```

或者使用 gx 运行：

```bash
gx run examples/config_demo.go
```

## 示例说明

### 配置管理

**文件：** `config_demo.go`

演示如何使用配置存储来管理 gx 的配置信息。

**功能：**
- 创建配置存储
- 加载配置文件
- 修改配置
- 保存配置
- 验证配置持久化

**运行：**
```bash
go run examples/config_demo.go
```

**输出示例：**
```
=== Config Store Demo ===

1. Ensuring config directory exists...
   ✓ Config directory ready

2. Loading configuration...
   ✓ Config loaded
   - Active Version: 
   - Install Path: /home/user/.gx/versions
   - Versions: 0 installed

3. Updating configuration...

4. Saving configuration...
   ✓ Config saved successfully

5. Reloading to verify...
   ✓ Config reloaded
   - Active Version: 1.21.5
   - Installed Versions:
     • 1.21.5 -> /home/user/.gx/versions/go1.21.5
     • 1.22.0 -> /home/user/.gx/versions/go1.22.0
   - Last Update Check: 2024-01-15T10:30:00Z

=== Demo Complete ===
```

**学习要点：**
- 配置文件位置：`~/.gx/config.json`
- 配置结构：`Config` 类型
- 原子保存：先写临时文件再重命名

---

### 下载器

**文件：** `downloader_demo.go`

演示如何使用下载器下载 Go 安装包。

**功能：**
- 获取下载 URL
- 下载文件并显示进度
- 验证 SHA256 校验和
- 处理下载错误

**运行：**
```bash
go run examples/downloader_demo.go
```

**输出示例：**
```
=== Downloader Demo ===

1. Getting download URL for Go 1.21.5...
   ✓ URL: https://go.dev/dl/go1.21.5.linux-amd64.tar.gz

2. Downloading Go 1.21.5...
   Progress: [████████████████████] 100% (150.5 MB / 150.5 MB)
   ✓ Download complete

3. Verifying checksum...
   ✓ Checksum verified

=== Demo Complete ===
```

**学习要点：**
- 下载 URL 格式
- 进度回调机制
- SHA256 校验和验证
- 错误处理和重试

---

### 环境管理

**文件：** `environment_demo.go`

演示如何管理系统环境变量。

**功能：**
- 设置 GOROOT
- 更新 PATH
- 持久化环境变量
- 平台特定处理

**运行：**
```bash
go run examples/environment_demo.go
```

**输出示例：**
```
=== Environment Manager Demo ===

1. Current environment:
   GOROOT: /usr/local/go
   PATH: /usr/local/bin:/usr/bin:/bin

2. Setting GOROOT to /home/user/.gx/versions/go1.21.5...
   ✓ GOROOT updated

3. Updating PATH...
   ✓ PATH updated

4. Persisting environment variables...
   ✓ Environment persisted

5. New environment:
   GOROOT: /home/user/.gx/versions/go1.21.5
   PATH: /home/user/.gx/versions/go1.21.5/bin:/usr/local/bin:/usr/bin:/bin

=== Demo Complete ===
```

**学习要点：**
- Windows：注册表修改
- Linux/macOS：Shell 配置文件修改
- 环境变量持久化策略

---

### 错误处理

**文件：** `error_handling_demo.go`

演示 gx 的错误处理机制。

**功能：**
- 自定义错误类型
- 错误包装和传播
- 错误检查和处理
- 用户友好的错误消息

**运行：**
```bash
go run examples/error_handling_demo.go
```

**输出示例：**
```
=== Error Handling Demo ===

1. Testing version not found error...
   Error: version not found: go1.99.0
   Error type: ErrVersionNotFound

2. Testing wrapped error...
   Error: failed to install version: failed to download: network timeout
   Unwrapped: network timeout

3. Testing error formatting...
   User message: Unable to install Go 1.99.0
   Suggestion: Please check the version number and try again

=== Demo Complete ===
```

**学习要点：**
- 错误类型定义
- 使用 `fmt.Errorf` 和 `%w` 包装错误
- 使用 `errors.Is` 和 `errors.As` 检查错误
- 用户友好的错误消息

---

### 跨平台构建

**文件：** `crossbuilder_demo.go`

演示如何使用跨平台构建器。

**功能：**
- 设置目标平台
- 执行跨平台构建
- 处理平台特定选项
- 验证构建输出

**运行：**
```bash
go run examples/crossbuilder_demo.go
```

**输出示例：**
```
=== Cross Builder Demo ===

1. Building for linux/amd64...
   ✓ Build successful: dist/myapp-linux-amd64

2. Building for windows/amd64...
   ✓ Build successful: dist/myapp-windows-amd64.exe

3. Building for darwin/arm64...
   ✓ Build successful: dist/myapp-darwin-arm64

=== Demo Complete ===
```

**学习要点：**
- GOOS 和 GOARCH 环境变量
- 支持的平台组合
- 平台特定文件扩展名
- 构建标志和链接器标志

---

### CLI 包装器

**文件：** `wrapper_demo.go`

演示如何使用 CLI 包装器执行 Go 命令。

**功能：**
- 获取 Go 可执行文件路径
- 执行 Go 命令
- 透传输入/输出
- 保留退出码

**运行：**
```bash
go run examples/wrapper_demo.go
```

**输出示例：**
```
=== CLI Wrapper Demo ===

1. Getting Go executable path...
   ✓ Go path: /home/user/.gx/versions/go1.21.5/bin/go

2. Executing 'go version'...
   go version go1.21.5 linux/amd64
   ✓ Command executed successfully

3. Executing 'go env GOROOT'...
   /home/user/.gx/versions/go1.21.5
   ✓ Command executed successfully

=== Demo Complete ===
```

**学习要点：**
- 命令执行机制
- 标准流透传
- 退出码处理
- 错误处理

---

### 版本管理

**文件：** `list_versions.go`, `list_remote_versions.go`

演示如何列出本地和远程的 Go 版本。

**功能：**
- 检测已安装版本
- 查询远程可用版本
- 获取最新版本
- 版本信息格式化

**运行：**
```bash
# 列出本地版本
go run examples/list_versions.go

# 列出远程版本
go run examples/list_remote_versions.go
```

**输出示例：**
```
=== List Versions Demo ===

Installed Go Versions:
  1.20.12
✓ 1.21.5 (active)
  1.22.0

Available Remote Versions:
  1.22.0  1.21.6  1.21.5  1.21.4
  1.21.3  1.21.2  1.21.1  1.21.0
  ...

Latest stable version: 1.22.0

=== Demo Complete ===
```

**学习要点：**
- 版本检测算法
- Go 官方 API 使用
- 版本号解析和比较
- 版本信息展示

---

## 测试应用

**目录：** `test_app/`

一个简单的测试应用，用于验证 gx 的各项功能。

**文件：**
- `main.go` - 主程序
- `utils.go` - 工具函数

**功能：**
- 显示 Go 版本信息
- 显示环境变量
- 测试跨平台构建

**运行：**
```bash
# 使用 gx 运行
gx run examples/test_app/main.go

# 使用 gx 构建
gx build -o test_app examples/test_app

# 跨平台构建
gx cross-build --os linux --arch amd64 -o test_app-linux examples/test_app
```

---

## 高级示例

### 简单跨平台构建

**文件：** `simple_crossbuild.go`

演示最简单的跨平台构建用法。

```bash
go run examples/simple_crossbuild.go
```

### 测试下载 URL

**文件：** `test_download_url.go`

测试不同平台和版本的下载 URL 生成。

```bash
go run examples/test_download_url.go
```

---

## 开发建议

### 1. 使用示例学习

示例程序是学习 gx 内部 API 的最佳方式：

1. 阅读示例代码
2. 运行示例程序
3. 修改示例代码进行实验
4. 参考示例编写自己的代码

### 2. 错误处理

所有示例都展示了正确的错误处理方式：

```go
result, err := someFunction()
if err != nil {
    log.Fatalf("Operation failed: %v", err)
}
```

### 3. 资源清理

注意资源清理，使用 `defer`：

```go
file, err := os.Open(path)
if err != nil {
    return err
}
defer file.Close()
```

### 4. 日志记录

使用 gx 的日志系统：

```go
import "github.com/yourusername/gx/internal/logger"

logger.Info("Operation started")
logger.Error("Operation failed: %v", err)
```

---

## 常见问题

### Q: 示例程序需要安装 gx 吗？

**A:** 不需要。示例程序直接使用 gx 的内部包，可以在开发环境中运行。

### Q: 可以修改示例程序吗？

**A:** 当然可以！示例程序就是用来学习和实验的。

### Q: 示例程序会修改我的系统吗？

**A:** 大多数示例程序只是演示，不会修改系统。但某些示例（如环境管理）可能会创建临时文件或目录。

### Q: 如何调试示例程序？

**A:** 使用 Go 的调试工具：

```bash
# 使用 delve
dlv debug examples/config_demo.go

# 使用 VS Code 或 GoLand 的调试功能
```

---

## 贡献示例

如果你有好的示例想要分享：

1. 创建新的示例文件
2. 添加清晰的注释
3. 更新本 README
4. 提交 Pull Request

示例程序应该：
- 简单易懂
- 专注于一个功能
- 包含完整的错误处理
- 有清晰的输出

---

**文档版本：** 1.0  
**最后更新：** 2024-01-15
