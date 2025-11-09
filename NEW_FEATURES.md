# 🎉 gx 新功能说明

## 新增功能

### 1. 一键安装到系统 PATH

现在你可以使用 `gx init-install` 命令自动将 gx 添加到系统 PATH，无需手动配置！

#### 快速开始

**方式一：使用安装脚本（推荐）**

```bash
# Linux/macOS
chmod +x install.sh
./install.sh

# Windows (PowerShell)
.\install.ps1
```

**方式二：手动运行 init-install**

```bash
# 1. 构建 gx
go build -o build/gx ./cmd/gx

# 2. 运行 init-install
./build/gx init-install
```

#### 功能特性

✅ **跨平台支持**
- Windows: 自动安装到 `%LOCALAPPDATA%\gx\bin`
- Linux/macOS: 可选择系统级或用户级安装

✅ **自动配置 PATH**
- Windows: 通过 PowerShell 修改用户环境变量
- Linux/macOS: 自动更新 shell 配置文件

✅ **交互式安装**
- 显示安装位置
- 请求确认
- 提供清晰的后续步骤

✅ **智能检测**
- 检测是否已安装
- 支持强制重新安装（`--force`）

#### 使用示例

```bash
# 首次安装
$ ./build/gx init-install

╔════════════════════════════════════════╗
║     gx Installation                    ║
╚════════════════════════════════════════╝

Current executable: /path/to/gx/build/gx
Installing gx on Linux...

Installation directory: /home/user/.local/bin
Proceed with installation? [Y/n]: y

Creating installation directory...
Copying gx executable...
Setting executable permissions...
Adding to PATH in shell configuration...

✓ gx installed successfully!

Installation complete. Please restart your terminal or run:
  source ~/.bashrc  (bash)
  source ~/.zshrc   (zsh)

Then you can use 'gx' from anywhere:
  gx --version
  gx install 1.21.5
```

---

### 2. 版本检测问题修复

修复了版本号格式不一致导致的问题。

#### 修复的问题

❌ **修复前：**
```bash
$ gx list
✓ 1.24.5 (active)    # 不带 "go" 前缀

$ gx current
✓ Current Go version: 1.21.5    # 与 list 不一致

$ gx install 1.25.4
✗ version already installed    # 错误判断
```

✅ **修复后：**
```bash
$ gx list
✓ go1.24.5 (active)    # 统一格式

$ gx current
✓ Current Go version: go1.24.5    # 一致

$ gx install 1.25.4
ℹ Installing Go 1.25.4...    # 正确工作
```

#### 技术细节

- 统一使用完整版本号格式（带 "go" 前缀）
- 修复版本扫描逻辑
- 修复系统 Go 版本检测
- 确保配置文件一致性

---

## 完整使用流程

### 新用户

```bash
# 1. 克隆仓库
git clone https://github.com/yourusername/gx.git
cd gx

# 2. 一键安装
./install.sh        # Linux/macOS
.\install.ps1       # Windows

# 3. 重启终端

# 4. 开始使用
gx --version
gx install 1.21.5
gx use 1.21.5
```

### 现有用户

如果你已经在使用 gx，建议重新构建并运行 init-install：

```bash
# 1. 更新代码
git pull

# 2. 重新构建
go build -o build/gx ./cmd/gx

# 3. 运行 init-install（可选）
./build/gx init-install

# 4. 验证版本检测
./build/gx list
./build/gx current
```

---

## 文档更新

新增和更新的文档：

📚 **新增文档：**
- `INSTALLATION.md` - 详细安装指南
- `QUICKSTART.md` - 5分钟快速上手
- `docs/VERSION_FIX.md` - 版本检测修复说明
- `SUMMARY.md` - 功能实现总结

📝 **更新文档：**
- `README.md` - 添加快速安装说明
- `COMMANDS.md` - 添加 init-install 命令文档
- `docs/README.md` - 更新文档导航

---

## 常见问题

### Q: 安装后找不到 gx 命令？

**A:** 需要重启终端或重新加载配置：

```bash
# Linux/macOS
source ~/.bashrc  # bash
source ~/.zshrc   # zsh

# Windows
# 重启 PowerShell 或命令提示符
```

### Q: 版本号显示格式变了？

**A:** 这是正常的。新版本统一使用完整格式（带 "go" 前缀），例如 `go1.21.5` 而不是 `1.21.5`。这样更符合 Go 官方命名规范。

### Q: 如何卸载 gx？

**A:** 参考 `INSTALLATION.md` 中的卸载说明：

```bash
# Windows
Remove-Item -Recurse -Force "$env:LOCALAPPDATA\gx"

# Linux/macOS
rm ~/.local/bin/gx  # 或 sudo rm /usr/local/bin/gx
rm -rf ~/.gx
```

### Q: 可以自定义安装位置吗？

**A:** 可以。不使用 init-install，手动复制文件到你想要的位置，然后手动添加到 PATH。

---

## 反馈和建议

如果你遇到问题或有建议：

- 📝 [提交 Issue](https://github.com/yourusername/gx/issues)
- 💬 [参与讨论](https://github.com/yourusername/gx/discussions)
- 📖 [查看文档](docs/README.md)

---

## 致谢

感谢所有用户的反馈和建议，让 gx 变得更好！

---

**发布日期：** 2024-01-15  
**版本：** v1.0.0
