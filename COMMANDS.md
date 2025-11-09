# gx 命令参考

本文档提供 gx 所有命令的详细参考。

## 目录

- [全局选项](#全局选项)
- [版本管理命令](#版本管理命令)
  - [install](#install)
  - [list](#list)
  - [use](#use)
  - [current](#current)
  - [update](#update)
  - [uninstall](#uninstall)
- [CLI 包装命令](#cli-包装命令)
  - [run](#run)
  - [build](#build)
  - [test](#test)
- [跨平台构建命令](#跨平台构建命令)
  - [cross-build](#cross-build)
- [使用技巧](#使用技巧)

## 全局选项

所有命令都支持以下全局选项：

### `-v, --verbose`

启用详细输出模式，显示更多调试信息。

```bash
gx install 1.21.5 --verbose
gx list -v
```

### `--config <file>`

指定配置文件路径。默认使用 `$HOME/.gx/config.json`。

```bash
gx --config /path/to/config.json install 1.21.5
```

### `--version`

显示 gx 的版本信息。

```bash
gx --version
# 输出: gx version 1.0.0 (commit: abc123, built: 2024-01-15)
```

### `-h, --help`

显示帮助信息。

```bash
gx --help
gx install --help
```

## 安装和配置命令

### init-install

将 gx 安装到系统 PATH，使其可以在任何位置使用。

#### 语法

```bash
gx init-install [flags]
```

#### 选项

- `-f, --force` - 强制重新安装，即使已经安装

#### 示例

```bash
# 首次安装
./build/gx init-install

# 强制重新安装
./build/gx init-install --force
```

#### 行为

**Windows:**
1. 将 gx.exe 复制到 `%LOCALAPPDATA%\gx\bin`
2. 使用 PowerShell 将该目录添加到用户 PATH 环境变量
3. 提示重启终端以使更改生效

**Linux/macOS:**
1. 提供两个安装选项：
   - `/usr/local/bin` - 系统级安装（需要 sudo）
   - `~/.local/bin` - 用户级安装（无需 sudo）
2. 复制 gx 到选定目录
3. 设置可执行权限
4. 如果安装到用户目录，自动添加到 shell 配置文件（.bashrc, .zshrc 等）
5. 提示重新加载 shell 配置或重启终端

#### 安装位置

| 操作系统 | 默认位置 | 需要权限 |
|---------|---------|---------|
| Windows | `%LOCALAPPDATA%\gx\bin` | 否 |
| Linux (系统级) | `/usr/local/bin` | sudo |
| Linux (用户级) | `~/.local/bin` | 否 |
| macOS (系统级) | `/usr/local/bin` | sudo |
| macOS (用户级) | `~/.local/bin` | 否 |

#### 验证安装

安装完成后，重启终端并运行：

```bash
# 检查 gx 是否在 PATH 中
which gx        # Linux/macOS
where.exe gx    # Windows

# 验证版本
gx --version
```

#### 故障排除

**问题：安装后仍然找不到 gx 命令**

解决方案：
- Windows: 重启命令提示符或 PowerShell
- Linux/macOS: 运行 `source ~/.bashrc` 或 `source ~/.zshrc`，或重启终端

**问题：权限被拒绝**

解决方案：
- Linux/macOS: 选择用户级安装（`~/.local/bin`）而不是系统级
- 或者使用 `sudo` 运行安装命令

**问题：PATH 没有自动更新**

解决方案：
手动添加到 shell 配置文件：

```bash
# bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc

# zsh
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

---

## 版本管理命令

### install

安装指定版本的 Go。

#### 语法

```bash
gx install [version] [flags]
```

#### 参数

- `version` (可选) - 要安装的 Go 版本号，可以带或不带 "go" 前缀
  - 示例：`1.21.5` 或 `go1.21.5`
  - 如果不指定，将安装最新稳定版本

#### 选项

- `-i, --interactive` - 交互式选择要安装的版本

#### 示例

```bash
# 安装特定版本
gx install 1.21.5
gx install go1.21.5

# 安装最新版本
gx install

# 交互式选择版本
gx install -i
gx install --interactive
```

#### 行为

1. 如果不指定版本，会查询最新稳定版本并提示确认
2. 如果使用 `-i` 标志，会显示可用版本列表供选择
3. 下载过程中显示进度条
4. 下载完成后验证 SHA256 校验和
5. 解压并安装到 `~/.gx/versions/` 目录
6. 安装完成后提示如何切换到新版本

#### 注意事项

- 需要网络连接
- 下载大小约 100-150 MB
- 如果版本已安装，会提示错误
- 安装过程中可以按 Ctrl+C 取消

---

### list

列出已安装的 Go 版本或可用的远程版本。

#### 语法

```bash
gx list [flags]
```

#### 选项

- `-r, --remote` - 列出可用的远程版本而不是已安装版本
- `-v, --verbose` - 显示详细信息（路径、安装日期等）

#### 示例

```bash
# 列出已安装的版本（简单格式）
gx list
# 输出:
#   1.20.12
# ✓ 1.21.5 (active)
#   1.22.0

# 列出已安装的版本（详细格式）
gx list --verbose
# 输出表格，包含状态、版本、路径、安装日期

# 列出可用的远程版本
gx list --remote
gx list -r
# 输出前 30 个可用版本
```

#### 输出格式

**简单格式：**
```
  1.20.12
✓ 1.21.5 (active)
  1.22.0
```

**详细格式（--verbose）：**
```
Status | Version | Path                          | Installed
-------|---------|-------------------------------|------------
✓      | 1.21.5  | /home/user/.gx/versions/...  | 2024-01-15
       | 1.20.12 | /home/user/.gx/versions/...  | 2024-01-10
```

**远程版本（--remote）：**
```
Available Go Versions

  1.22.0         1.21.6         1.21.5
  1.21.4         1.21.3         1.21.2
  ...

... and 50 more versions

To install a version, run:
  gx install <version>
```

---

### use

切换到指定的 Go 版本。

#### 语法

```bash
gx use [version] [flags]
```

#### 参数

- `version` (可选) - 要切换到的 Go 版本号
  - 如果不指定，会显示交互式选择界面

#### 选项

- `-i, --interactive` - 交互式选择版本（即使指定了版本也会显示选择界面）

#### 示例

```bash
# 切换到特定版本
gx use 1.21.5
gx use go1.21.5

# 交互式选择版本
gx use
gx use -i
```

#### 行为

1. 验证指定版本是否已安装
2. 更新配置文件中的 `active_version`
3. 设置 `GOROOT` 环境变量
4. 更新 `PATH` 环境变量
5. 持久化环境变量更改
6. 显示成功消息和后续步骤提示

#### 环境变量更新

**Windows:**
- 通过注册表更新用户环境变量
- 需要重启终端或命令提示符

**Linux/macOS:**
- 修改 shell 配置文件（`.bashrc`, `.zshrc` 等）
- 需要重新加载配置文件或重启终端

```bash
# Linux/macOS 重新加载配置
source ~/.bashrc   # bash
source ~/.zshrc    # zsh
```

#### 注意事项

- 切换操作通常在 300ms 内完成
- 无法切换到未安装的版本
- 切换后需要重启终端或重新加载配置

---

### current

显示当前激活的 Go 版本。

#### 语法

```bash
gx current [flags]
```

#### 选项

- `-v, --verbose` - 显示详细信息（安装路径）

#### 示例

```bash
# 显示当前版本
gx current
# 输出: Current Go version: 1.21.5

# 显示详细信息
gx current --verbose
# 输出:
# Current Go version: 1.21.5
# Installation path: /home/user/.gx/versions/go1.21.5
```

#### 行为

1. 从配置文件读取当前激活版本
2. 验证版本是否仍然存在
3. 显示版本信息

#### 错误情况

如果没有激活的版本：
```
Error: no active version set

To set a version, run:
  gx use <version>
```

---

### update

更新到最新的 Go 版本。

#### 语法

```bash
gx update [flags]
```

#### 选项

- `-s, --switch` - 安装后自动切换到新版本

#### 示例

```bash
# 检查并安装最新版本
gx update

# 安装后自动切换
gx update --switch
gx update -s
```

#### 行为

1. 查询 Go 官方 API 获取最新稳定版本
2. 检查该版本是否已安装
3. 如果已安装且是当前版本，显示提示并退出
4. 如果已安装但不是当前版本，询问是否切换
5. 如果未安装，下载并安装
6. 安装完成后询问是否切换（除非使用了 `--switch` 标志）

#### 示例输出

```bash
$ gx update
Checking for the latest Go version...
Latest version: 1.22.0
Installing Go 1.22.0...
Downloading ████████████████████ 100%
Go 1.22.0 installed successfully

Switch to Go 1.22.0 now? [Y/n]: y
Switching to 1.22.0...
Now using Go 1.22.0
```

---

### uninstall

卸载指定的 Go 版本。

#### 语法

```bash
gx uninstall <version> [flags]
```

#### 参数

- `version` (必需) - 要卸载的 Go 版本号

#### 选项

- `-f, --force` - 跳过确认提示，强制卸载

#### 示例

```bash
# 卸载特定版本（会提示确认）
gx uninstall 1.20.12
gx uninstall go1.20.12

# 强制卸载（跳过确认）
gx uninstall 1.20.12 --force
gx uninstall 1.20.12 -f
```

#### 行为

1. 验证版本是否已安装
2. 检查是否为当前激活版本（不能卸载激活版本）
3. 提示确认（除非使用 `--force`）
4. 删除版本目录
5. 更新配置文件

#### 确认提示

```bash
$ gx uninstall 1.20.12
Are you sure you want to uninstall Go 1.20.12? [y/N]: y
Uninstalling Go 1.20.12...
Go 1.20.12 uninstalled successfully
```

#### 错误情况

尝试卸载当前激活版本：
```
Error: cannot uninstall active version

Switch to another version first:
  gx use <version>
```

---

## CLI 包装命令

这些命令是对 Go 原生命令的包装，使用当前激活的 Go 版本执行。

### run

编译并运行 Go 程序。

#### 语法

```bash
gx run [flags] <file.go> [arguments...]
```

#### 参数

- `file.go` (必需) - 要运行的 Go 源文件
- `arguments...` (可选) - 传递给程序的参数

#### 标志

所有 `go run` 支持的标志都可以使用。

#### 示例

```bash
# 运行单个文件
gx run main.go

# 传递参数给程序
gx run main.go arg1 arg2 arg3

# 使用 race detector
gx run -race main.go

# 使用构建标志
gx run -ldflags="-s -w" main.go

# 运行多个文件
gx run main.go utils.go
```

#### 行为

1. 获取当前激活版本的 `go` 可执行文件路径
2. 构建完整命令：`go run [flags] <file.go> [arguments...]`
3. 执行命令，透传标准输入/输出/错误
4. 保留原始退出码

#### 等价命令

```bash
gx run main.go
# 等价于
/path/to/active/go/bin/go run main.go
```

---

### build

编译 Go 包和依赖。

#### 语法

```bash
gx build [flags] [packages]
```

#### 参数

- `packages` (可选) - 要构建的包，默认为当前目录

#### 标志

所有 `go build` 支持的标志都可以使用。

#### 示例

```bash
# 构建当前目录
gx build

# 指定输出文件
gx build -o myapp

# 构建特定文件
gx build main.go

# 使用 ldflags
gx build -ldflags="-s -w" .

# 构建特定包
gx build ./cmd/myapp

# 构建多个包
gx build ./...
```

#### 常用标志

- `-o <file>` - 指定输出文件名
- `-ldflags <flags>` - 链接器标志
  - `-s` - 去除符号表
  - `-w` - 去除 DWARF 调试信息
  - `-X` - 设置变量值
- `-tags <tags>` - 构建标签
- `-race` - 启用 race detector
- `-v` - 显示正在编译的包

#### 示例：设置版本信息

```bash
gx build -ldflags="-X main.version=1.0.0 -X main.commit=abc123" -o myapp
```

---

### test

运行测试。

#### 语法

```bash
gx test [flags] [packages]
```

#### 参数

- `packages` (可选) - 要测试的包，默认为当前目录

#### 标志

所有 `go test` 支持的标志都可以使用。

#### 示例

```bash
# 测试当前包
gx test

# 测试所有包
gx test ./...

# 详细输出
gx test -v

# 显示覆盖率
gx test -cover

# 生成覆盖率报告
gx test -coverprofile=coverage.out ./...

# 运行特定测试
gx test -run TestMyFunction

# 使用 race detector
gx test -race ./...

# 并行测试
gx test -parallel 4 ./...
```

#### 常用标志

- `-v` - 详细输出
- `-cover` - 显示覆盖率
- `-coverprofile <file>` - 生成覆盖率报告
- `-run <regexp>` - 只运行匹配的测试
- `-race` - 启用 race detector
- `-parallel <n>` - 并行运行测试
- `-timeout <duration>` - 测试超时时间
- `-short` - 运行短测试

#### 查看覆盖率报告

```bash
# 生成覆盖率报告
gx test -coverprofile=coverage.out ./...

# 查看覆盖率
go tool cover -func=coverage.out

# 生成 HTML 报告
go tool cover -html=coverage.out
```

---

## 跨平台构建命令

### cross-build

为不同的操作系统和架构编译 Go 程序。

#### 语法

```bash
gx cross-build [source] [flags]
```

#### 参数

- `source` (可选) - 源代码路径，默认为当前目录

#### 选项

- `--os <os>` (必需*) - 目标操作系统
  - 可选值：`windows`, `linux`, `darwin`
- `--arch <arch>` (必需*) - 目标架构
  - 可选值：`amd64`, `arm64`, `386`
- `-o, --output <path>` - 输出文件路径
- `--ldflags <flags>` - 链接器标志
- `--flags <flags>` - 额外的构建标志
- `--list-platforms` - 列出支持的平台

*注：使用 `--list-platforms` 时不需要 `--os` 和 `--arch`

#### 示例

```bash
# 为 Linux amd64 构建
gx cross-build --os linux --arch amd64 -o myapp

# 为 Windows 构建
gx cross-build --os windows --arch amd64 -o myapp.exe

# 为 macOS ARM64 构建
gx cross-build --os darwin --arch arm64 -o myapp

# 指定源代码路径
gx cross-build ./cmd/myapp --os linux --arch amd64 -o dist/myapp

# 使用 ldflags
gx cross-build --os linux --arch amd64 --ldflags="-s -w" -o myapp

# 使用额外的构建标志
gx cross-build --os linux --arch amd64 --flags="-tags=prod" -o myapp

# 列出支持的平台
gx cross-build --list-platforms
```

#### 支持的平台

```
Supported Platforms

  windows:
    • amd64
    • arm64
    • 386

  linux:
    • amd64
    • arm64
    • 386

  darwin:
    • amd64
    • arm64
```

#### 批量构建示例

```bash
# 为多个平台构建
gx cross-build --os linux --arch amd64 -o dist/myapp-linux-amd64
gx cross-build --os linux --arch arm64 -o dist/myapp-linux-arm64
gx cross-build --os windows --arch amd64 -o dist/myapp-windows-amd64.exe
gx cross-build --os darwin --arch amd64 -o dist/myapp-darwin-amd64
gx cross-build --os darwin --arch arm64 -o dist/myapp-darwin-arm64
```

#### 构建脚本示例

**Linux/macOS (build-all.sh):**
```bash
#!/bin/bash

VERSION="1.0.0"
LDFLAGS="-s -w -X main.version=$VERSION"

platforms=(
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "darwin/amd64"
    "darwin/arm64"
)

for platform in "${platforms[@]}"; do
    os="${platform%/*}"
    arch="${platform#*/}"
    output="dist/myapp-$os-$arch"
    
    if [ "$os" = "windows" ]; then
        output="$output.exe"
    fi
    
    echo "Building for $os/$arch..."
    gx cross-build --os "$os" --arch "$arch" --ldflags="$LDFLAGS" -o "$output"
done

echo "Build complete!"
```

**Windows (build-all.ps1):**
```powershell
$VERSION = "1.0.0"
$LDFLAGS = "-s -w -X main.version=$VERSION"

$platforms = @(
    @{os="linux"; arch="amd64"},
    @{os="linux"; arch="arm64"},
    @{os="windows"; arch="amd64"},
    @{os="darwin"; arch="amd64"},
    @{os="darwin"; arch="arm64"}
)

foreach ($platform in $platforms) {
    $os = $platform.os
    $arch = $platform.arch
    $output = "dist/myapp-$os-$arch"
    
    if ($os -eq "windows") {
        $output = "$output.exe"
    }
    
    Write-Host "Building for $os/$arch..."
    gx cross-build --os $os --arch $arch --ldflags=$LDFLAGS -o $output
}

Write-Host "Build complete!"
```

---

## 使用技巧

### 1. 快速切换版本

使用交互式模式快速切换：

```bash
gx use -i
```

### 2. 检查版本信息

```bash
# 查看 gx 管理的当前版本
gx current

# 查看实际的 go 版本
go version

# 查看 GOROOT
go env GOROOT
```

### 3. 清理旧版本

```bash
# 列出所有版本
gx list

# 卸载不需要的版本
gx uninstall 1.19.5
gx uninstall 1.20.0
```

### 4. 自动化脚本

在 CI/CD 或脚本中使用 gx：

```bash
#!/bin/bash

# 安装特定版本
gx install 1.21.5

# 切换版本
gx use 1.21.5

# 构建项目
gx build -o myapp

# 运行测试
gx test ./...
```

### 5. 项目特定版本

在项目中创建脚本来设置正确的 Go 版本：

```bash
#!/bin/bash
# setup-go.sh

REQUIRED_VERSION="1.21.5"

# 检查版本是否已安装
if ! gx list | grep -q "$REQUIRED_VERSION"; then
    echo "Installing Go $REQUIRED_VERSION..."
    gx install "$REQUIRED_VERSION"
fi

# 切换到所需版本
echo "Switching to Go $REQUIRED_VERSION..."
gx use "$REQUIRED_VERSION"

echo "Go version setup complete!"
go version
```

### 6. 环境变量验证

验证环境变量是否正确设置：

```bash
# 检查 GOROOT
echo $GOROOT

# 检查 PATH
echo $PATH | grep go

# 验证 go 命令路径
which go
```

### 7. 日志查看

如果遇到问题，查看日志：

```bash
# Linux/macOS
cat ~/.gx/logs/gx.log

# Windows
type %USERPROFILE%\.gx\logs\gx.log
```

### 8. 配置备份

备份配置文件：

```bash
# Linux/macOS
cp ~/.gx/config.json ~/.gx/config.json.backup

# Windows
copy %USERPROFILE%\.gx\config.json %USERPROFILE%\.gx\config.json.backup
```

---

## 常见问题

### Q: 切换版本后 `go version` 仍显示旧版本？

**A:** 需要重启终端或重新加载 shell 配置：

```bash
# bash
source ~/.bashrc

# zsh
source ~/.zshrc

# 或者重启终端
```

### Q: 如何完全卸载 gx？

**A:** 删除 gx 目录和配置：

```bash
# 删除 gx 数据目录
rm -rf ~/.gx

# 删除 gx 可执行文件
rm /usr/local/bin/gx  # 或你安装的位置

# 清理环境变量（手动编辑 shell 配置文件）
```

### Q: 下载速度慢怎么办？

**A:** 可以考虑：
1. 使用代理
2. 等待未来版本支持镜像源配置

### Q: 如何在多个项目中使用不同的 Go 版本？

**A:** 在每个项目中切换版本：

```bash
cd project-a
gx use 1.20.12

cd ../project-b
gx use 1.21.5
```

或者创建项目特定的设置脚本。

---

**文档版本：** 1.0  
**最后更新：** 2024-01-15
