# gx 快速开始指南

5 分钟快速上手 gx Go 版本管理工具。

## 第一步：安装 gx

### Windows

```powershell
git clone https://github.com/kawaiirei0/gx.git
cd gx
.\install.ps1
```

### Linux/macOS

```bash
git clone https://github.com/kawaiirei0/gx.git
cd gx
chmod +x install.sh
./install.sh
```

安装完成后，**重启终端**。

## 第二步：验证安装

```bash
gx --version
```

应该看到类似输出：
```
gx version 1.0.0 (commit: abc123, built: 2024-01-15)
```

## 第三步：安装 Go 版本

### 安装最新版本

```bash
gx install
```

### 安装特定版本

```bash
gx install 1.21.5
```

### 交互式选择版本

```bash
gx install -i
```

## 第四步：切换 Go 版本

```bash
gx use 1.21.5
```

或交互式选择：

```bash
gx use -i
```

## 第五步：验证 Go 版本

```bash
gx current
go version
```

## 常用命令速查

```bash
# 列出已安装的版本
gx list

# 列出可用的远程版本
gx list --remote

# 更新到最新版本
gx update

# 卸载版本
gx uninstall 1.20.0

# 运行 Go 程序
gx run main.go

# 构建项目
gx build

# 跨平台构建
gx cross-build --os linux --arch amd64 -o myapp
```

## 下一步

- 查看 [完整命令参考](COMMANDS.md)
- 了解 [架构设计](ARCHITECTURE.md)
- 阅读 [详细安装指南](INSTALLATION.md)

## 需要帮助？

- 查看 [故障排除](INSTALLATION.md#故障排除)
- 提交 [Issue](https://github.com/kawaiirei0/gx/issues)
- 参与 [讨论](https://github.com/kawaiirei0/gx/discussions)

---

**提示：** 使用 `gx --help` 查看所有可用命令。
