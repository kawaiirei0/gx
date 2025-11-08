# gx 架构文档

本文档描述了 gx Go 版本管理工具的架构设计、核心组件和实现细节。

## 目录

- [设计理念](#设计理念)
- [架构概览](#架构概览)
- [核心组件](#核心组件)
- [数据流](#数据流)
- [错误处理](#错误处理)
- [平台适配](#平台适配)
- [性能考虑](#性能考虑)
- [安全性](#安全性)

## 设计理念

gx 的设计遵循以下核心理念：

1. **单一职责原则** - 每个组件只负责一个明确的功能
2. **接口驱动** - 通过接口定义组件边界，便于测试和扩展
3. **跨平台兼容** - 统一的 API，平台差异由适配层处理
4. **快速响应** - 最小化启动开销，优化常用操作
5. **用户友好** - 清晰的错误提示和进度反馈

## 架构概览

### 高层架构图

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI Layer                            │
│                    (Cobra Commands)                          │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────┴────────────────────────────────────┐
│                    Command Router                            │
└────┬──────────┬──────────┬──────────┬──────────────────────┘
     │          │          │          │
     ▼          ▼          ▼          ▼
┌─────────┐ ┌──────┐ ┌─────────┐ ┌──────────┐
│ Version │ │ CLI  │ │  Cross  │ │  Config  │
│ Manager │ │Wrapper│ │ Builder │ │  Store   │
└────┬────┘ └──┬───┘ └────┬────┘ └──────────┘
     │         │          │
     ▼         ▼          ▼
┌──────────────────────────────────┐
│      Platform Adapter            │
│  (OS/Arch Abstraction Layer)     │
└──────────────────────────────────┘
     │         │          │
     ▼         ▼          ▼
┌──────────┐ ┌──────┐ ┌──────┐
│Downloader│ │ Env  │ │ File │
│          │ │ Mgr  │ │System│
└──────────┘ └──────┘ └──────┘
```

### 分层架构

1. **表示层 (Presentation Layer)**
   - CLI 命令定义和参数解析
   - 用户交互（进度条、提示、错误显示）

2. **应用层 (Application Layer)**
   - 命令路由和协调
   - 业务逻辑编排

3. **领域层 (Domain Layer)**
   - 核心业务逻辑
   - 版本管理、下载、安装等

4. **基础设施层 (Infrastructure Layer)**
   - 平台适配
   - 文件系统操作
   - 网络通信

## 核心组件

### 1. Version Manager

**职责：** 管理 Go 版本的生命周期

**接口定义：**

```go
type VersionManager interface {
    // 检测系统中已安装的 Go 版本
    DetectInstalled() ([]GoVersion, error)
    
    // 获取当前激活的版本
    GetActive() (*GoVersion, error)
    
    // 安装指定版本
    Install(version string, progress ProgressCallback) error
    
    // 切换到指定版本
    SwitchTo(version string) error
    
    // 获取可用的远程版本列表
    ListAvailable() ([]string, error)
    
    // 获取最新稳定版本
    GetLatest() (string, error)
    
    // 卸载指定版本
    Uninstall(version string) error
}
```

**实现细节：**

- **版本存储：** `~/.gx/versions/go{version}/`
- **版本检测：** 扫描版本目录，验证 `go` 可执行文件
- **版本切换：** 更新配置文件和环境变量
- **并发安全：** 使用文件锁防止并发安装冲突

**关键算法：**

```go
// 版本切换流程
func (vm *versionManager) SwitchTo(version string) error {
    1. 验证版本是否已安装
    2. 获取版本安装路径
    3. 更新配置文件中的 active_version
    4. 调用 EnvironmentManager 更新环境变量
    5. 验证切换是否成功
    6. 记录日志
}
```

### 2. Downloader

**职责：** 从官方源下载 Go 安装包

**接口定义：**

```go
type Downloader interface {
    // 下载指定版本的 Go 安装包
    Download(version string, destPath string, progress ProgressCallback) error
    
    // 获取下载 URL
    GetDownloadURL(version string, os string, arch string) (string, error)
    
    // 验证下载文件的完整性
    VerifyChecksum(filePath string, expectedSHA256 string) error
}
```

**实现细节：**

- **下载源：** `https://go.dev/dl/`
- **进度跟踪：** 通过回调函数报告下载进度
- **完整性验证：** SHA256 校验和验证
- **错误恢复：** 支持断点续传（如果服务器支持）
- **超时处理：** 30 秒连接超时，5 分钟读取超时

**下载流程：**

```go
func (d *downloader) Download(version, destPath string, progress ProgressCallback) error {
    1. 获取下载 URL 和 SHA256
    2. 创建临时文件
    3. 发起 HTTP GET 请求
    4. 流式写入文件，报告进度
    5. 验证 SHA256 校验和
    6. 移动到目标路径
    7. 清理临时文件
}
```

### 3. Installer

**职责：** 安装和卸载 Go 版本

**接口定义：**

```go
type Installer interface {
    // 安装 Go 版本（解压和配置）
    Install(archivePath string, version string) error
    
    // 卸载 Go 版本
    Uninstall(version string) error
    
    // 验证安装是否成功
    Verify(version string) error
}
```

**实现细节：**

- **解压：** 支持 `.tar.gz` (Linux/macOS) 和 `.zip` (Windows)
- **权限设置：** 在 Unix 系统上设置可执行权限
- **安装验证：** 执行 `go version` 验证安装
- **清理：** 卸载时删除整个版本目录

### 4. Environment Manager

**职责：** 管理系统环境变量

**接口定义：**

```go
type EnvironmentManager interface {
    // 设置 GOROOT 环境变量
    SetGOROOT(path string) error
    
    // 更新 PATH 环境变量
    UpdatePATH(goPath string) error
    
    // 获取当前环境变量
    GetEnv(key string) string
    
    // 持久化环境变量更改
    PersistEnv() error
}
```

**平台特定实现：**

**Windows:**
```go
// 通过注册表修改用户环境变量
func (em *windowsEnvManager) SetGOROOT(path string) error {
    key := `Environment`
    registry.SetStringValue(HKEY_CURRENT_USER, key, "GOROOT", path)
}
```

**Linux/macOS:**
```go
// 修改 shell 配置文件
func (em *unixEnvManager) SetGOROOT(path string) error {
    shells := []string{".bashrc", ".zshrc", ".profile"}
    for _, shell := range shells {
        appendToFile(shell, "export GOROOT=" + path)
    }
}
```

### 5. CLI Wrapper

**职责：** 包装和转发 Go 原生命令

**接口定义：**

```go
type CLIWrapper interface {
    // 执行 Go 命令
    Execute(command string, args []string) error
    
    // 获取当前使用的 Go 可执行文件路径
    GetGoExecutable() (string, error)
}
```

**实现细节：**

- **命令执行：** 使用 `os/exec` 包
- **流透传：** 直接连接 stdin/stdout/stderr
- **退出码保留：** 保留原始命令的退出码
- **环境变量：** 继承当前进程的环境变量

**执行流程：**

```go
func (w *cliWrapper) Execute(command string, args []string) error {
    1. 获取当前激活版本的 go 可执行文件路径
    2. 构建完整命令：goPath + command + args
    3. 创建 exec.Cmd
    4. 连接 stdin/stdout/stderr
    5. 执行命令
    6. 等待完成并返回退出码
}
```

### 6. Cross Builder

**职责：** 处理跨平台编译

**接口定义：**

```go
type CrossBuilder interface {
    // 跨平台构建
    Build(config BuildConfig) error
    
    // 获取支持的平台列表
    GetSupportedPlatforms() []Platform
}

type BuildConfig struct {
    SourcePath   string
    OutputPath   string
    TargetOS     string
    TargetArch   string
    BuildFlags   []string
    LDFlags      string
}
```

**实现细节：**

- **环境变量设置：** 设置 `GOOS` 和 `GOARCH`
- **平台验证：** 验证目标平台组合是否支持
- **文件扩展名：** 自动添加平台特定扩展名（Windows 的 `.exe`）
- **构建标志：** 支持自定义构建标志和链接器标志

**支持的平台矩阵：**

| OS      | amd64 | arm64 | 386 |
|---------|-------|-------|-----|
| Windows | ✓     | ✓     | ✓   |
| Linux   | ✓     | ✓     | ✓   |
| Darwin  | ✓     | ✓     | ✗   |

### 7. Platform Adapter

**职责：** 提供跨平台抽象层

**接口定义：**

```go
type PlatformAdapter interface {
    // 获取操作系统类型
    GetOS() string
    
    // 获取架构类型
    GetArch() string
    
    // 获取路径分隔符
    PathSeparator() string
    
    // 规范化路径
    NormalizePath(path string) string
    
    // 检查文件是否可执行
    IsExecutable(path string) bool
    
    // 设置文件为可执行
    MakeExecutable(path string) error
}
```

**实现细节：**

- **路径处理：** 统一使用 `filepath` 包处理路径
- **权限管理：** Unix 系统使用 `chmod +x`，Windows 无需处理
- **环境检测：** 使用 `runtime.GOOS` 和 `runtime.GOARCH`

### 8. Config Store

**职责：** 持久化工具配置和状态

**接口定义：**

```go
type ConfigStore interface {
    // 加载配置
    Load() (*Config, error)
    
    // 保存配置
    Save(config *Config) error
    
    // 获取配置文件路径
    GetConfigPath() string
}

type Config struct {
    ActiveVersion    string            `json:"active_version"`
    InstallPath      string            `json:"install_path"`
    Versions         map[string]string `json:"versions"`
    LastUpdateCheck  time.Time         `json:"last_update_check"`
}
```

**实现细节：**

- **存储格式：** JSON
- **配置位置：** `~/.gx/config.json`
- **原子写入：** 先写入临时文件，再重命名
- **备份：** 保存前备份旧配置

## 数据流

### 安装流程

```
用户命令: gx install 1.21.5
    │
    ▼
┌─────────────────────┐
│  CLI Command        │
│  (install.go)       │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  Version Manager    │
│  - 验证版本格式      │
│  - 检查是否已安装    │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  Downloader         │
│  - 获取下载 URL      │
│  - 下载安装包        │
│  - 验证校验和        │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  Installer          │
│  - 解压安装包        │
│  - 设置权限          │
│  - 验证安装          │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  Config Store       │
│  - 更新配置文件      │
└─────────────────────┘
```

### 版本切换流程

```
用户命令: gx use 1.21.5
    │
    ▼
┌─────────────────────┐
│  CLI Command        │
│  (use.go)           │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  Version Manager    │
│  - 验证版本已安装    │
│  - 获取版本路径      │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  Config Store       │
│  - 更新 active_version│
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  Environment Manager│
│  - 设置 GOROOT       │
│  - 更新 PATH         │
│  - 持久化环境变量    │
└─────────────────────┘
```

### 命令包装流程

```
用户命令: gx run main.go
    │
    ▼
┌─────────────────────┐
│  CLI Command        │
│  (run.go)           │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  CLI Wrapper        │
│  - 获取 go 路径      │
│  - 构建命令          │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│  Process Executor   │
│  - 执行命令          │
│  - 透传 I/O          │
│  - 返回退出码        │
└─────────────────────┘
```

## 错误处理

### 错误类型层次

```go
// 基础错误类型
var (
    ErrVersionNotFound      = errors.New("version not found")
    ErrVersionExists        = errors.New("version already installed")
    ErrNoActiveVersion      = errors.New("no active version set")
    ErrDownloadFailed       = errors.New("download failed")
    ErrInvalidVersion       = errors.New("invalid version format")
    ErrPlatformNotSupported = errors.New("platform not supported")
    ErrChecksumMismatch     = errors.New("checksum verification failed")
    ErrInstallationFailed   = errors.New("installation failed")
)
```

### 错误处理策略

1. **错误包装：** 使用 `fmt.Errorf` 和 `%w` 保留错误链
2. **上下文信息：** 在每一层添加相关上下文
3. **用户友好：** 在 UI 层转换为易懂的错误消息
4. **日志记录：** 所有错误都记录到日志文件
5. **优雅降级：** 非关键错误不中断程序

**示例：**

```go
// 底层错误
err := os.Open(path)
if err != nil {
    return fmt.Errorf("failed to open file %s: %w", path, err)
}

// 中间层包装
err := downloader.Download(version, path, callback)
if err != nil {
    return fmt.Errorf("failed to download Go %s: %w", version, err)
}

// UI 层处理
if err != nil {
    errorFormatter.Format(err)  // 显示友好的错误消息
    logger.Error("Installation failed: %v", err)  // 记录详细错误
    return err
}
```

## 平台适配

### 平台差异处理

| 功能 | Windows | Linux/macOS |
|------|---------|-------------|
| 环境变量持久化 | 注册表 | Shell 配置文件 |
| 路径分隔符 | `\` | `/` |
| 可执行权限 | 不需要 | `chmod +x` |
| 安装包格式 | `.zip` | `.tar.gz` |
| 默认 Shell | PowerShell/CMD | bash/zsh |

### 平台检测

```go
func detectPlatform() (os, arch string) {
    os = runtime.GOOS
    arch = runtime.GOARCH
    
    // 规范化操作系统名称
    switch os {
    case "darwin":
        os = "darwin"  // macOS
    case "windows":
        os = "windows"
    case "linux":
        os = "linux"
    }
    
    return os, arch
}
```

## 性能考虑

### 启动时间优化

- **延迟加载：** 只在需要时加载配置和检测版本
- **缓存：** 缓存版本列表和配置信息
- **并发：** 使用 goroutine 并发执行独立操作

**目标：** 启动时间 < 100ms

### 版本切换优化

- **符号链接：** 使用符号链接而非复制文件（如果平台支持）
- **配置缓存：** 内存中缓存配置，减少文件 I/O
- **环境变量批量更新：** 一次性更新所有环境变量

**目标：** 切换时间 < 300ms

### 下载优化

- **流式处理：** 边下载边写入，不占用大量内存
- **进度更新节流：** 限制进度回调频率（每 100ms）
- **HTTP/2：** 使用 HTTP/2 协议提升下载速度

## 安全性

### 下载安全

1. **HTTPS：** 只从 HTTPS 源下载
2. **证书验证：** 验证服务器证书
3. **校验和验证：** 验证 SHA256 校验和
4. **官方源：** 只从 `go.dev` 下载

### 文件系统安全

1. **路径验证：** 防止路径遍历攻击
2. **权限控制：** 配置文件权限设置为 600
3. **原子操作：** 使用临时文件和重命名保证原子性

### 命令执行安全

1. **参数验证：** 验证所有用户输入
2. **命令注入防护：** 使用 `exec.Command` 而非 shell 执行
3. **环境隔离：** 不继承敏感环境变量

## 扩展性

### 插件系统（未来）

预留插件接口，支持：
- 自定义下载源
- 版本管理策略
- 构建流程定制

### 配置扩展

支持通过配置文件扩展：
- 镜像源配置
- 代理设置
- 自定义安装路径

## 测试策略

### 单元测试

- 每个组件独立测试
- 使用 mock 对象隔离依赖
- 覆盖率目标：80%+

### 集成测试

- 测试组件间交互
- 使用临时目录和测试服务器
- 覆盖主要用户场景

### 平台测试

- CI/CD 在三个平台上运行测试
- 测试平台特定功能
- 验证跨平台兼容性

---

**文档版本：** 1.0  
**最后更新：** 2024-01-15
