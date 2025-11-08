# 贡献指南

感谢你对 gx 项目的关注！我们欢迎各种形式的贡献，包括但不限于：

- 🐛 报告 Bug
- 💡 提出新功能建议
- 📝 改进文档
- 🔧 提交代码修复或新功能
- 🧪 编写测试用例
- 🌍 翻译文档

## 目录

- [行为准则](#行为准则)
- [如何贡献](#如何贡献)
- [开发环境设置](#开发环境设置)
- [代码规范](#代码规范)
- [提交规范](#提交规范)
- [Pull Request 流程](#pull-request-流程)
- [测试指南](#测试指南)
- [文档贡献](#文档贡献)

## 行为准则

本项目采用 [Contributor Covenant](https://www.contributor-covenant.org/) 行为准则。参与本项目即表示你同意遵守其条款。

### 我们的承诺

为了营造一个开放和友好的环境，我们承诺：

- 使用友好和包容的语言
- 尊重不同的观点和经验
- 优雅地接受建设性批评
- 关注对社区最有利的事情
- 对其他社区成员表示同理心

## 如何贡献

### 报告 Bug

如果你发现了 Bug，请通过 [GitHub Issues](https://github.com/yourusername/gx/issues) 报告。

**好的 Bug 报告应该包含：**

1. **清晰的标题** - 简洁描述问题
2. **复现步骤** - 详细的步骤说明
3. **预期行为** - 你期望发生什么
4. **实际行为** - 实际发生了什么
5. **环境信息** - 操作系统、Go 版本、gx 版本等
6. **日志输出** - 相关的错误日志（`~/.gx/logs/gx.log`）
7. **截图** - 如果适用

**Bug 报告模板：**

```markdown
## Bug 描述
简要描述问题

## 复现步骤
1. 执行命令 `gx install 1.21.5`
2. 然后执行 `gx use 1.21.5`
3. 看到错误...

## 预期行为
应该成功切换到 Go 1.21.5

## 实际行为
显示错误：version not found

## 环境信息
- OS: Windows 11
- gx 版本: v1.0.0
- Go 版本: 1.20.5

## 日志输出
```
[ERROR] Failed to switch version: ...
```

## 截图
（如果适用）
```

### 提出功能建议

我们欢迎新功能建议！请通过 [GitHub Discussions](https://github.com/yourusername/gx/discussions) 或 Issues 提出。

**好的功能建议应该包含：**

1. **问题描述** - 你想解决什么问题？
2. **建议方案** - 你建议如何实现？
3. **替代方案** - 是否考虑过其他方案？
4. **使用场景** - 谁会使用这个功能？
5. **示例** - 提供使用示例

## 开发环境设置

### 前置要求

- Go 1.19 或更高版本
- Git
- Make (Linux/macOS) 或 PowerShell (Windows)

### 克隆仓库

```bash
git clone https://github.com/yourusername/gx.git
cd gx
```

### 安装依赖

```bash
go mod download
```

### 构建项目

```bash
# Linux/macOS
make build

# Windows
.\build.ps1 build
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test -cover ./...

# 运行特定包的测试
go test ./internal/version/...

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 运行示例

```bash
# 运行示例程序
go run examples/config_demo.go
go run examples/downloader_demo.go
```

## 代码规范

### Go 代码风格

我们遵循 Go 官方代码风格指南：

1. **格式化：** 使用 `gofmt` 或 `goimports` 格式化代码
2. **命名：** 遵循 Go 命名约定
   - 包名：小写，单个单词
   - 导出标识符：大写开头（PascalCase）
   - 未导出标识符：小写开头（camelCase）
3. **注释：** 为导出的函数、类型和常量添加文档注释
4. **错误处理：** 不要忽略错误，适当处理或返回

### 代码组织

```go
// 包注释
package mypackage

import (
    // 标准库
    "fmt"
    "os"
    
    // 第三方库
    "github.com/spf13/cobra"
    
    // 本项目包
    "github.com/yourusername/gx/pkg/interfaces"
)

// 常量
const (
    DefaultTimeout = 30
)

// 变量
var (
    ErrNotFound = errors.New("not found")
)

// 类型定义
type MyStruct struct {
    Field1 string
    Field2 int
}

// 构造函数
func NewMyStruct() *MyStruct {
    return &MyStruct{}
}

// 方法
func (m *MyStruct) Method() error {
    // 实现
    return nil
}

// 辅助函数
func helperFunction() {
    // 实现
}
```

### 注释规范

**包注释：**
```go
// Package version provides Go version management functionality.
// It handles version detection, installation, and switching.
package version
```

**函数注释：**
```go
// Install downloads and installs the specified Go version.
// It returns an error if the version is already installed or if the download fails.
//
// Example:
//   err := vm.Install("go1.21.5", progressCallback)
func (vm *VersionManager) Install(version string, progress ProgressCallback) error {
    // 实现
}
```

**复杂逻辑注释：**
```go
// 检查版本是否已安装
// 1. 扫描版本目录
// 2. 验证 go 可执行文件存在
// 3. 执行 go version 确认版本号
if vm.IsInstalled(version) {
    return ErrVersionExists
}
```

### 错误处理

**好的错误处理：**
```go
// ✅ 正确：包装错误，添加上下文
file, err := os.Open(path)
if err != nil {
    return fmt.Errorf("failed to open config file %s: %w", path, err)
}
defer file.Close()

// ✅ 正确：检查特定错误类型
if errors.Is(err, ErrVersionNotFound) {
    // 处理版本未找到的情况
}

// ✅ 正确：使用 errors.As 提取错误
var exitErr *ExitError
if errors.As(err, &exitErr) {
    os.Exit(exitErr.GetExitCode())
}
```

**不好的错误处理：**
```go
// ❌ 错误：忽略错误
file, _ := os.Open(path)

// ❌ 错误：丢失错误上下文
if err != nil {
    return err
}

// ❌ 错误：使用 panic 处理可恢复错误
if err != nil {
    panic(err)
}
```

### 测试规范

**测试文件命名：** `*_test.go`

**测试函数命名：** `TestXxx` 或 `TestXxx_Scenario`

**示例测试：**
```go
func TestVersionManager_Install(t *testing.T) {
    // 准备
    vm := NewVersionManager(mockConfig)
    version := "go1.21.5"
    
    // 执行
    err := vm.Install(version, nil)
    
    // 断言
    if err != nil {
        t.Errorf("Install() error = %v, want nil", err)
    }
    
    // 验证
    installed, err := vm.IsInstalled(version)
    if !installed {
        t.Error("Version should be installed")
    }
}

func TestVersionManager_Install_AlreadyInstalled(t *testing.T) {
    // 测试已安装版本的情况
    vm := NewVersionManager(mockConfig)
    version := "go1.21.5"
    
    // 第一次安装
    _ = vm.Install(version, nil)
    
    // 第二次安装应该返回错误
    err := vm.Install(version, nil)
    if !errors.Is(err, ErrVersionExists) {
        t.Errorf("Install() error = %v, want ErrVersionExists", err)
    }
}
```

**表驱动测试：**
```go
func TestNormalizeVersion(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"with go prefix", "go1.21.5", "go1.21.5", false},
        {"without go prefix", "1.21.5", "go1.21.5", false},
        {"invalid format", "invalid", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := NormalizeVersion(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("NormalizeVersion() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("NormalizeVersion() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## 提交规范

我们使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范。

### 提交消息格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type 类型

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响代码运行）
- `refactor`: 重构（既不是新功能也不是 Bug 修复）
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动
- `ci`: CI 配置文件和脚本的变动

### Scope 范围

- `version`: 版本管理
- `download`: 下载功能
- `install`: 安装功能
- `cli`: CLI 命令
- `wrapper`: CLI 包装器
- `crossbuild`: 跨平台构建
- `config`: 配置管理
- `ui`: 用户界面
- `docs`: 文档

### 示例

```bash
# 新功能
feat(version): add support for beta versions

# Bug 修复
fix(download): handle network timeout correctly

# 文档更新
docs(readme): update installation instructions

# 重构
refactor(config): simplify config loading logic

# 性能优化
perf(version): cache version list to improve performance
```

## Pull Request 流程

### 1. Fork 仓库

点击 GitHub 页面右上角的 "Fork" 按钮。

### 2. 创建分支

```bash
# 克隆你的 fork
git clone https://github.com/your-username/gx.git
cd gx

# 添加上游仓库
git remote add upstream https://github.com/yourusername/gx.git

# 创建特性分支
git checkout -b feature/my-feature
```

### 3. 开发和测试

```bash
# 进行修改
# ...

# 运行测试
go test ./...

# 运行 linter
go vet ./...
golangci-lint run

# 格式化代码
gofmt -s -w .
```

### 4. 提交更改

```bash
# 添加更改
git add .

# 提交（遵循提交规范）
git commit -m "feat(version): add support for beta versions"
```

### 5. 同步上游

```bash
# 获取上游更新
git fetch upstream

# 合并到你的分支
git rebase upstream/main
```

### 6. 推送到 Fork

```bash
git push origin feature/my-feature
```

### 7. 创建 Pull Request

1. 访问你的 Fork 页面
2. 点击 "New Pull Request"
3. 填写 PR 描述

**PR 描述模板：**

```markdown
## 变更说明
简要描述这个 PR 做了什么

## 变更类型
- [ ] Bug 修复
- [ ] 新功能
- [ ] 重构
- [ ] 文档更新
- [ ] 性能优化
- [ ] 其他

## 相关 Issue
Closes #123

## 测试
描述你如何测试这些更改

## 检查清单
- [ ] 代码遵循项目风格指南
- [ ] 已添加/更新测试
- [ ] 所有测试通过
- [ ] 已更新文档
- [ ] 提交消息遵循规范
```

### 8. 代码审查

- 响应审查意见
- 进行必要的修改
- 推送更新

### 9. 合并

PR 被批准后，维护者会合并你的代码。

## 测试指南

### 单元测试

**编写单元测试的原则：**

1. **独立性：** 测试之间不应相互依赖
2. **可重复：** 测试结果应该一致
3. **快速：** 单元测试应该快速执行
4. **清晰：** 测试意图应该明确

**使用 Mock：**

```go
// 定义 mock 接口
type MockDownloader struct {
    DownloadFunc func(version, dest string, progress ProgressCallback) error
}

func (m *MockDownloader) Download(version, dest string, progress ProgressCallback) error {
    if m.DownloadFunc != nil {
        return m.DownloadFunc(version, dest, progress)
    }
    return nil
}

// 在测试中使用
func TestInstall_WithMock(t *testing.T) {
    mockDownloader := &MockDownloader{
        DownloadFunc: func(version, dest string, progress ProgressCallback) error {
            // 模拟下载成功
            return nil
        },
    }
    
    vm := NewVersionManager(mockDownloader)
    err := vm.Install("go1.21.5", nil)
    
    if err != nil {
        t.Errorf("Install() failed: %v", err)
    }
}
```

### 集成测试

**使用临时目录：**

```go
func TestIntegration_InstallAndSwitch(t *testing.T) {
    // 创建临时目录
    tmpDir, err := os.MkdirTemp("", "gx-test-*")
    if err != nil {
        t.Fatal(err)
    }
    defer os.RemoveAll(tmpDir)
    
    // 使用临时目录进行测试
    config := &Config{
        InstallPath: tmpDir,
    }
    
    vm := NewVersionManager(config)
    
    // 测试完整流程
    // ...
}
```

### 测试覆盖率

```bash
# 生成覆盖率报告
go test -coverprofile=coverage.out ./...

# 查看覆盖率
go tool cover -func=coverage.out

# 生成 HTML 报告
go tool cover -html=coverage.out -o coverage.html
```

**覆盖率目标：**
- 整体覆盖率：> 80%
- 核心包覆盖率：> 90%

## 文档贡献

### 文档类型

1. **README.md** - 项目概览和快速开始
2. **ARCHITECTURE.md** - 架构设计文档
3. **CONTRIBUTING.md** - 贡献指南（本文档）
4. **BUILD.md** - 构建说明
5. **代码注释** - 函数和类型的文档注释

### 文档风格

- 使用清晰、简洁的语言
- 提供代码示例
- 使用适当的格式（标题、列表、代码块）
- 保持文档更新

### 文档审查

文档更改也需要通过 PR 流程，确保：
- 拼写和语法正确
- 格式一致
- 示例代码可运行
- 链接有效

## 发布流程

（仅限维护者）

1. 更新 CHANGELOG.md
2. 更新版本号
3. 创建 Git tag
4. 构建发布包
5. 发布到 GitHub Releases

## 获取帮助

如果你有任何问题：

- 📖 查看 [文档](README.md)
- 💬 在 [Discussions](https://github.com/yourusername/gx/discussions) 提问
- 🐛 在 [Issues](https://github.com/yourusername/gx/issues) 报告问题

## 致谢

感谢所有贡献者！你们的贡献让 gx 变得更好。

---

再次感谢你的贡献！🎉
