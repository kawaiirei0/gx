# gx - Go Version Manager

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue.svg)](https://golang.org/)

一个跨平台的 Go 版本管理和开发工具，简化 Go 环境的安装、版本切换和跨平台编译。

## ✨ 特性

- 🚀 **快速版本切换** - 在多个 Go 版本之间快速切换（< 300ms）
- 📦 **简单安装** - 一键安装任意 Go 版本
- 🔄 **自动更新** - 轻松更新到最新的 Go 版本
- 🌍 **跨平台支持** - 支持 Windows、Linux 和 macOS
- 🛠️ **CLI 包装器** - 统一的命令行接口包装 Go 原生命令
- 🏗️ **跨平台构建** - 在任意平台上为其他平台编译程序
- 📊 **友好的用户界面** - 进度条、交互式选择和清晰的错误提示

## 📋 目录

- [安装](#-安装)
- [快速开始](#-快速开始)
- [命令参考](#-命令参考)
- [使用示例](#-使用示例)
- [配置](#-配置)
- [开发](#-开发)
- [贡献](#-贡献)
- [许可证](#-许可证)

## 🚀 安装

### 快速安装（推荐）

**前置要求：** Go 1.19 或更高版本

```bash
# 克隆仓库
git clone https://github.com/kawaiirei0/gx.git
cd gx

# 一键安装（自动添加到系统 PATH）
# Linux/macOS
chmod +x install.sh
./install.sh

# Windows (PowerShell)
.\install.ps1
```

这个脚本会：
1. 构建 gx
2. 自动将 gx 添加到系统 PATH
3. 让你可以在任何位置使用 `gx` 命令

### 手动安装

如果你想手动控制安装过程：

```bash
# 1. 构建 gx
# Windows (PowerShell)
.\build.ps1 build

# Linux/macOS
./build.sh build

# 2. 运行 init-install 命令
# Windows
.\build\gx.exe init-install

# Linux/macOS
./build/gx init-install
```

`init-install` 命令会：
- 将 gx 复制到系统目录
- 自动添加到 PATH 环境变量
- 在 Windows 上：安装到 `%LOCALAPPDATA%\gx\bin`
- 在 Linux/macOS 上：可选择安装到 `/usr/local/bin` 或 `~/.local/bin`

### 传统安装方式

使用 Go 的标准安装方式：

```bash
# 安装到 GOPATH/bin
go install ./cmd/gx

# 或使用构建脚本
# Windows
.\build.ps1 install

# Linux/macOS
make install
```

### 从发布版本安装

从 [Releases](https://github.com/kawaiirei0/gx/releases) 页面下载适合你操作系统的预编译二进制文件。

## 🎯 快速开始

```bash
# 1. 安装最新版本的 Go
gx install

# 2. 切换到已安装的版本
gx use 1.21.5

# 3. 验证当前版本
gx current

# 4. 列出所有已安装的版本
gx list

# 5. 使用 gx 运行 Go 程序
gx run main.go

# 6. 跨平台构建
gx cross-build --os linux --arch amd64 -o myapp
```

## 📚 文档

- **[README.md](README.md)** - 项目概览和快速开始（本文档）
- **[INSTALLATION.md](INSTALLATION.md)** - 详细的安装指南和故障排除
- **[COMMANDS.md](COMMANDS.md)** - 详细的命令参考文档
- **[ARCHITECTURE.md](ARCHITECTURE.md)** - 架构设计和技术细节
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - 贡献指南和开发规范
- **[BUILD.md](BUILD.md)** - 构建说明和发布流程
- **[CHANGELOG.md](CHANGELOG.md)** - 版本变更历史
- **[examples/README.md](examples/README.md)** - 示例程序说明

## 📖 命令参考

完整的命令参考请查看 [COMMANDS.md](COMMANDS.md)。

### 版本管理命令

#### `gx install [version]`

安装指定版本的 Go。如果不指定版本，将安装最新稳定版本。

```bash
# 安装特定版本
gx install 1.21.5

# 安装最新版本
gx install

# 交互式选择版本
gx install -i
```

**选项：**
- `-i, --interactive` - 交互式版本选择

#### `gx list`

列出所有已安装的 Go 版本。

```bash
# 简单列表
gx list

# 详细信息（包含路径和安装日期）
gx list -v

# 列出可用的远程版本
gx list --remote
```

**选项：**
- `-r, --remote` - 列出可用的远程版本
- `-v, --verbose` - 显示详细信息

#### `gx use [version]`

切换到指定的 Go 版本。

```bash
# 切换到特定版本
gx use 1.21.5

# 交互式选择版本
gx use -i
```

**选项：**
- `-i, --interactive` - 交互式版本选择

**注意：** 切换版本后，可能需要重启终端或重新加载 shell 配置文件。

#### `gx current`

显示当前激活的 Go 版本。

```bash
gx current

# 显示详细信息（包含安装路径）
gx current -v
```

#### `gx update`

更新到最新的 Go 版本。

```bash
# 检查并安装最新版本
gx update

# 安装后自动切换到新版本
gx update --switch
```

**选项：**
- `-s, --switch` - 安装后自动切换到新版本

#### `gx uninstall <version>`

卸载指定的 Go 版本。

```bash
# 卸载特定版本（会提示确认）
gx uninstall 1.21.5

# 强制卸载（跳过确认）
gx uninstall 1.21.5 --force
```

**选项：**
- `-f, --force` - 跳过确认提示

**注意：** 无法卸载当前激活的版本。

### CLI 包装命令

这些命令是对 Go 原生命令的包装，使用当前激活的 Go 版本执行。

#### `gx run [flags] <file.go> [arguments...]`

编译并运行 Go 程序。

```bash
# 运行单个文件
gx run main.go

# 传递参数给程序
gx run main.go arg1 arg2

# 使用 race detector
gx run -race main.go
```

#### `gx build [flags] [packages]`

编译 Go 包和依赖。

```bash
# 构建当前目录
gx build

# 指定输出文件
gx build -o myapp

# 使用 ldflags
gx build -ldflags="-s -w" .

# 构建特定文件
gx build main.go
```

#### `gx test [flags] [packages]`

运行测试。

```bash
# 测试当前包
gx test

# 测试所有包
gx test ./...

# 详细输出
gx test -v

# 显示覆盖率
gx test -cover ./...
```

### 跨平台构建命令

#### `gx cross-build [source]`

为不同的操作系统和架构编译 Go 程序。

```bash
# 为 Linux amd64 构建
gx cross-build --os linux --arch amd64 -o myapp

# 为 Windows 构建
gx cross-build --os windows --arch amd64 -o myapp.exe

# 为 macOS ARM64 构建
gx cross-build --os darwin --arch arm64 -o myapp

# 列出支持的平台
gx cross-build --list-platforms

# 使用额外的构建标志
gx cross-build --os linux --arch amd64 --ldflags="-s -w" -o myapp
```

**选项：**
- `--os <os>` - 目标操作系统（windows, linux, darwin）
- `--arch <arch>` - 目标架构（amd64, arm64, 386）
- `-o, --output <path>` - 输出文件路径
- `--ldflags <flags>` - 链接器标志
- `--flags <flags>` - 额外的构建标志
- `--list-platforms` - 列出支持的平台

**支持的平台：**
- Windows: amd64, 386
- Linux: amd64, arm64, 386
- macOS (darwin): amd64, arm64

### 全局选项

所有命令都支持以下全局选项：

- `-v, --verbose` - 详细输出
- `--config <file>` - 指定配置文件（默认：`$HOME/.gx/config.json`）
- `--version` - 显示版本信息
- `-h, --help` - 显示帮助信息

## 💡 使用示例

### 场景 1：设置新的开发环境

```bash
# 安装最新版本的 Go
gx install

# 切换到新安装的版本
gx use 1.21.5

# 验证安装
gx current
go version
```

### 场景 2：在不同项目间切换 Go 版本

```bash
# 项目 A 需要 Go 1.20
cd project-a
gx use 1.20.12

# 项目 B 需要 Go 1.21
cd ../project-b
gx use 1.21.5
```

### 场景 3：跨平台构建和部署

```bash
# 在 Windows 上为 Linux 服务器构建
gx cross-build --os linux --arch amd64 -o myserver

# 为多个平台构建
gx cross-build --os linux --arch amd64 -o dist/myapp-linux
gx cross-build --os windows --arch amd64 -o dist/myapp-windows.exe
gx cross-build --os darwin --arch arm64 -o dist/myapp-macos
```

### 场景 4：保持 Go 版本最新

```bash
# 检查并安装最新版本
gx update

# 自动切换到最新版本
gx update --switch

# 查看所有可用版本
gx list --remote
```

### 场景 5：清理旧版本

```bash
# 列出已安装的版本
gx list

# 卸载不再需要的版本
gx uninstall 1.19.5
gx uninstall 1.20.0
```

## ⚙️ 配置

gx 的配置和数据存储在用户主目录下的 `.gx` 文件夹中：

```
~/.gx/
├── config.json          # 配置文件
├── versions/            # 已安装的 Go 版本
│   ├── go1.21.5/
│   ├── go1.20.12/
│   └── ...
└── logs/                # 日志文件
    └── gx.log
```

### 配置文件格式

`~/.gx/config.json`:

```json
{
  "active_version": "go1.21.5",
  "install_path": "/home/user/.gx/versions",
  "versions": {
    "go1.21.5": "/home/user/.gx/versions/go1.21.5",
    "go1.20.12": "/home/user/.gx/versions/go1.20.12"
  },
  "last_update_check": "2024-01-15T10:30:00Z"
}
```

### 环境变量

gx 会自动管理以下环境变量：

- `GOROOT` - 指向当前激活的 Go 版本
- `PATH` - 包含当前 Go 版本的 bin 目录

**Windows:** 环境变量通过注册表持久化  
**Linux/macOS:** 环境变量写入 shell 配置文件（`.bashrc`, `.zshrc` 等）

## 🛠️ 开发

### 项目结构

```
gx/
├── cmd/
│   └── gx/                    # 主程序入口
│       ├── main.go
│       └── cmd/               # CLI 命令实现
├── internal/                  # 内部实现（不对外暴露）
│   ├── config/                # 配置管理
│   ├── crossbuilder/          # 跨平台构建
│   ├── downloader/            # 下载管理
│   ├── environment/           # 环境变量管理
│   ├── installer/             # 安装管理
│   ├── logger/                # 日志记录
│   ├── platform/              # 平台适配
│   ├── ui/                    # 用户界面组件
│   ├── utils/                 # 工具函数
│   ├── version/               # 版本管理
│   └── wrapper/               # CLI 包装器
├── pkg/                       # 公共包（可被外部引用）
│   ├── constants/             # 常量定义
│   ├── errors/                # 错误类型
│   └── interfaces/            # 核心接口定义
├── examples/                  # 示例程序
├── scripts/                   # 构建和发布脚本
├── go.mod
├── go.sum
├── README.md
├── BUILD.md                   # 构建文档
└── CHANGELOG.md               # 变更日志
```

### 构建

详细的构建说明请参阅 [BUILD.md](BUILD.md)。

**快速构建：**

```bash
# Windows
.\build.ps1 build

# Linux/macOS
make build
```

**运行测试：**

```bash
# 运行所有测试
go test ./...

# 运行测试并显示覆盖率
go test -cover ./...

# 运行特定包的测试
go test ./internal/version/...
```

### 核心接口

gx 采用接口驱动的设计，主要接口包括：

- **VersionManager** - 管理 Go 版本的安装、切换和检测
- **Downloader** - 负责下载 Go 安装包
- **Installer** - 负责安装和卸载 Go 版本
- **EnvironmentManager** - 管理系统环境变量
- **CLIWrapper** - 包装和转发 Go 原生命令
- **CrossBuilder** - 处理跨平台编译
- **PlatformAdapter** - 提供跨平台抽象层

详细的架构设计请参阅 [ARCHITECTURE.md](ARCHITECTURE.md)。

## 🤝 贡献

欢迎贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解如何参与项目。

### 贡献流程

1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交你的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 开启一个 Pull Request

### 代码规范

- 遵循 Go 官方代码风格指南
- 使用 `gofmt` 格式化代码
- 添加适当的注释和文档
- 为新功能编写测试
- 确保所有测试通过

## 📝 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 🙏 致谢

- [Go 官方团队](https://golang.org/) - 提供优秀的编程语言
- [Cobra](https://github.com/spf13/cobra) - 强大的 CLI 框架
- 所有贡献者和用户

## 📞 联系方式

- 问题反馈：[GitHub Issues](https://github.com/kawaiirei0/gx/issues)
- 功能建议：[GitHub Discussions](https://github.com/kawaiirei0/gx/discussions)

---

**注意：** 本项目仍在积极开发中。欢迎反馈和建议！
