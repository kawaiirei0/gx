# gx 安装指南

本文档提供 gx 的详细安装说明，包括多种安装方式和故障排除。

## 目录

- [快速安装](#快速安装)
- [手动安装](#手动安装)
- [安装位置](#安装位置)
- [验证安装](#验证安装)
- [卸载](#卸载)
- [故障排除](#故障排除)

## 快速安装

这是最简单的安装方式，适合大多数用户。

### Windows

1. 打开 PowerShell
2. 克隆仓库并运行安装脚本：

```powershell
git clone https://github.com/kawaiirei0/gx.git
cd gx
.\install.ps1
```

3. 按照提示完成安装
4. 重启 PowerShell
5. 验证安装：`gx --version`

### Linux

1. 打开终端
2. 克隆仓库并运行安装脚本：

```bash
git clone https://github.com/kawaiirei0/gx.git
cd gx
chmod +x install.sh
./install.sh
```

3. 选择安装位置（系统级或用户级）
4. 按照提示完成安装
5. 重新加载 shell 配置或重启终端：
   ```bash
   source ~/.bashrc  # bash
   source ~/.zshrc   # zsh
   ```
6. 验证安装：`gx --version`

### macOS

1. 打开终端
2. 克隆仓库并运行安装脚本：

```bash
git clone https://github.com/kawaiirei0/gx.git
cd gx
chmod +x install.sh
./install.sh
```

3. 选择安装位置（系统级或用户级）
4. 如果选择系统级安装，可能需要输入密码
5. 按照提示完成安装
6. 重新加载 shell 配置或重启终端：
   ```bash
   source ~/.zshrc   # zsh (macOS 默认)
   source ~/.bashrc  # bash
   ```
7. 验证安装：`gx --version`

## 手动安装

如果你想更精细地控制安装过程，可以使用手动安装方式。

### 步骤 1：构建 gx

**Windows:**
```powershell
.\build.ps1 build
```

**Linux/macOS:**
```bash
./build.sh build
# 或
make build
```

构建完成后，可执行文件位于 `build/` 目录。

### 步骤 2：运行 init-install

**Windows:**
```powershell
.\build\gx.exe init-install
```

**Linux/macOS:**
```bash
./build/gx init-install
```

### 步骤 3：按照提示完成安装

`init-install` 命令会：

1. 显示当前可执行文件的位置
2. 询问安装位置（Linux/macOS）
3. 请求确认
4. 复制文件到目标位置
5. 添加到 PATH 环境变量
6. 显示后续步骤

## 安装位置

### Windows

gx 会被安装到：
```
%LOCALAPPDATA%\gx\bin\gx.exe
```

通常是：
```
C:\Users\<kawaiirei0>\AppData\Local\gx\bin\gx.exe
```

这个目录会被自动添加到用户的 PATH 环境变量。

### Linux

你可以选择两个安装位置：

**系统级安装（需要 sudo）：**
```
/usr/local/bin/gx
```
- 所有用户都可以使用
- 需要管理员权限
- 已经在系统 PATH 中

**用户级安装（推荐）：**
```
~/.local/bin/gx
```
- 只有当前用户可以使用
- 不需要管理员权限
- 会自动添加到 shell 配置文件

### macOS

与 Linux 相同，可以选择：

**系统级安装（需要 sudo）：**
```
/usr/local/bin/gx
```

**用户级安装（推荐）：**
```
~/.local/bin/gx
```

## 验证安装

安装完成后，验证 gx 是否正确安装：

### 检查 PATH

**Windows:**
```powershell
where.exe gx
```

应该显示：
```
C:\Users\<kawaiirei0>\AppData\Local\gx\bin\gx.exe
```

**Linux/macOS:**
```bash
which gx
```

应该显示：
```
/usr/local/bin/gx
# 或
/home/<username>/.local/bin/gx
```

### 检查版本

```bash
gx --version
```

应该显示类似：
```
gx version 1.0.0 (commit: abc123, built: 2024-01-15)
```

### 测试基本命令

```bash
# 显示帮助
gx --help

# 列出已安装版本
gx list

# 列出可用版本
gx list --remote
```

## 卸载

### Windows

1. 删除安装目录：
   ```powershell
   Remove-Item -Recurse -Force "$env:LOCALAPPDATA\gx"
   ```

2. 从 PATH 中移除（可选）：
   ```powershell
   $path = [Environment]::GetEnvironmentVariable('Path', 'User')
   $newPath = ($path.Split(';') | Where-Object { $_ -notlike '*gx*' }) -join ';'
   [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
   ```

3. 删除配置和数据：
   ```powershell
   Remove-Item -Recurse -Force "$env:USERPROFILE\.gx"
   ```

### Linux/macOS

1. 删除可执行文件：
   ```bash
   # 系统级安装
   sudo rm /usr/local/bin/gx
   
   # 用户级安装
   rm ~/.local/bin/gx
   ```

2. 从 shell 配置中移除 PATH（如果是用户级安装）：
   ```bash
   # 编辑配置文件，删除包含 gx 的行
   nano ~/.bashrc  # 或 ~/.zshrc
   ```

3. 删除配置和数据：
   ```bash
   rm -rf ~/.gx
   ```

## 故障排除

### 问题 1：安装后找不到 gx 命令

**症状：**
```bash
gx: command not found
```

**解决方案：**

**Windows:**
1. 重启 PowerShell 或命令提示符
2. 如果仍然不行，检查 PATH：
   ```powershell
   $env:Path -split ';' | Select-String gx
   ```
3. 手动添加到 PATH（如果需要）：
   ```powershell
   $path = [Environment]::GetEnvironmentVariable('Path', 'User')
   $newPath = $path + ';' + "$env:LOCALAPPDATA\gx\bin"
   [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
   ```

**Linux/macOS:**
1. 重新加载 shell 配置：
   ```bash
   source ~/.bashrc  # bash
   source ~/.zshrc   # zsh
   ```
2. 或者重启终端
3. 检查 PATH：
   ```bash
   echo $PATH | grep gx
   ```
4. 手动添加到配置文件（如果需要）：
   ```bash
   echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
   source ~/.bashrc
   ```

### 问题 2：权限被拒绝

**症状：**
```bash
Permission denied
```

**解决方案：**

**Linux/macOS:**
1. 如果尝试系统级安装，使用 sudo：
   ```bash
   sudo ./build/gx init-install
   ```
2. 或者选择用户级安装（不需要 sudo）

### 问题 3：构建失败

**症状：**
```bash
Build failed
```

**解决方案：**

1. 检查 Go 版本：
   ```bash
   go version
   ```
   需要 Go 1.19 或更高版本

2. 更新依赖：
   ```bash
   go mod download
   go mod tidy
   ```

3. 清理并重新构建：
   ```bash
   # Windows
   .\build.ps1 clean
   .\build.ps1 build
   
   # Linux/macOS
   make clean
   make build
   ```

### 问题 4：Windows 上 PowerShell 执行策略错误

**症状：**
```
.\install.ps1 : File cannot be loaded because running scripts is disabled
```

**解决方案：**

临时允许脚本执行：
```powershell
Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process
.\install.ps1
```

或者永久更改（需要管理员权限）：
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### 问题 5：macOS 上"无法验证开发者"错误

**症状：**
```
"gx" cannot be opened because the developer cannot be verified
```

**解决方案：**

1. 打开系统偏好设置 > 安全性与隐私
2. 点击"仍要打开"
3. 或者在终端中运行：
   ```bash
   xattr -d com.apple.quarantine ~/.local/bin/gx
   ```

### 问题 6：Linux 上 ~/.local/bin 不在 PATH 中

**症状：**
安装到 `~/.local/bin` 后仍然找不到命令

**解决方案：**

手动添加到 PATH：

```bash
# bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc

# zsh
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# fish
echo 'set -gx PATH $HOME/.local/bin $PATH' >> ~/.config/fish/config.fish
source ~/.config/fish/config.fish
```

## 高级安装选项

### 自定义安装位置

如果你想安装到自定义位置，可以手动复制文件：

```bash
# 构建
./build.sh build

# 复制到自定义位置
cp build/gx /your/custom/path/

# 添加到 PATH
export PATH="/your/custom/path:$PATH"

# 永久添加（添加到 shell 配置文件）
echo 'export PATH="/your/custom/path:$PATH"' >> ~/.bashrc
```

### 从 GOPATH 安装

使用 Go 的标准安装方式：

```bash
go install ./cmd/gx
```

这会将 gx 安装到 `$GOPATH/bin`。确保 `$GOPATH/bin` 在你的 PATH 中：

```bash
export PATH="$GOPATH/bin:$PATH"
```

### 使用符号链接

在 Linux/macOS 上，你可以创建符号链接而不是复制文件：

```bash
# 构建
./build.sh build

# 创建符号链接
sudo ln -s $(pwd)/build/gx /usr/local/bin/gx
```

## 多版本管理

如果你需要同时保留多个版本的 gx：

```bash
# 安装到带版本号的目录
cp build/gx ~/.local/bin/gx-1.0.0

# 创建符号链接指向当前版本
ln -s ~/.local/bin/gx-1.0.0 ~/.local/bin/gx

# 切换版本时，只需更新符号链接
ln -sf ~/.local/bin/gx-1.1.0 ~/.local/bin/gx
```

## 获取帮助

如果遇到其他问题：

1. 查看日志文件：`~/.gx/logs/gx.log`
2. 提交 Issue：https://github.com/kawaiirei0/gx/issues
3. 查看文档：https://github.com/kawaiirei0/gx/tree/main/docs

---

**文档版本：** 1.0  
**最后更新：** 2024-01-15
