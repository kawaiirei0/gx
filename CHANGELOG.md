# Changelog

All notable changes to gx will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-01-15

### Added
- **init-install 命令**: 自动将 gx 安装到系统 PATH
  - 跨平台支持（Windows、Linux、macOS）
  - 自动配置环境变量
  - 交互式安装流程
  - 支持系统级和用户级安装（Linux/macOS）
- **快速安装脚本**:
  - `install.sh` for Linux/macOS
  - `install.ps1` for Windows
- **完整文档体系**:
  - `INSTALLATION.md` - 详细安装指南（15+ KB）
  - `QUICKSTART.md` - 5分钟快速上手指南
  - `NEW_FEATURES.md` - 新功能说明
  - `SUMMARY.md` - 功能实现总结
  - `docs/VERSION_FIX.md` - 版本检测修复说明
  - 更新 `COMMANDS.md` 添加 init-install 文档
  - 更新 `README.md` 添加快速安装说明

### Fixed
- **版本检测问题**: 统一版本号格式（带 "go" 前缀）
  - 修复 `gx list` 和 `gx current` 显示不一致的问题
  - 修复安装检查错误判断版本已安装的问题
  - 修复激活状态标记错误的问题
  - 修复 `scanGxVersions()` 方法的版本号提取逻辑
  - 修复 `detectSystemGoVersion()` 方法的版本号格式

### Changed
- 版本号显示格式统一为 `go1.x.x`（带 "go" 前缀）
- 改进错误提示信息的友好性

## [Unreleased]

### Added
- Initial implementation of gx Go version manager
- Version detection and management
- Go installation and switching
- CLI wrapper for Go commands
- Cross-platform build support
- Multi-platform support (Windows, Linux, macOS)

### Changed

### Deprecated

### Removed

### Fixed

### Security

## [0.1.0] - YYYY-MM-DD

### Added
- Initial release
- Core version management functionality
- Download and install Go versions
- Switch between installed versions
- Wrap common Go CLI commands (run, build, test)
- Cross-platform compilation support
- Configuration management
- Logging system
- Error handling and recovery

---

## Release Notes Template

When creating a new release, copy this template:

```markdown
## [X.Y.Z] - YYYY-MM-DD

### Added
- New features

### Changed
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Removed features

### Fixed
- Bug fixes

### Security
- Security fixes
```

## Version Guidelines

- **Major version (X.0.0)**: Incompatible API changes
- **Minor version (0.X.0)**: New functionality in a backward compatible manner
- **Patch version (0.0.X)**: Backward compatible bug fixes

## Links

[Unreleased]: https://github.com/kawaiirei0/gx/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/kawaiirei0/gx/releases/tag/v0.1.0
