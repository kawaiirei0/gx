package environment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kawaiirei0/gx/internal/logger"
	"github.com/kawaiirei0/gx/pkg/constants"
	"github.com/kawaiirei0/gx/pkg/errors"
	"github.com/kawaiirei0/gx/pkg/interfaces"
)

// manager 实现 EnvironmentManager 接口
type manager struct {
	platform interfaces.PlatformAdapter
	backupPath string
}

// NewManager 创建新的环境变量管理器
func NewManager(platform interfaces.PlatformAdapter) interfaces.EnvironmentManager {
	homeDir, _ := platform.GetHomeDir()
	backupPath := filepath.Join(homeDir, constants.ConfigDir, "env_backup.json")
	
	return &manager{
		platform: platform,
		backupPath: backupPath,
	}
}

// SetGoRoot 设置 GOROOT 环境变量
func (m *manager) SetGoRoot(path string) error {
	logger.Info("Setting GOROOT to: %s", path)
	
	if path == "" {
		logger.Error("GOROOT path cannot be empty")
		return errors.ErrInvalidInput.WithMessage("GOROOT path cannot be empty")
	}

	// 验证路径存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.Error("GOROOT path does not exist: %s", path)
		return errors.ErrInvalidInput.
			WithMessage(fmt.Sprintf("GOROOT path does not exist: %s", path)).
			WithContext("path", path)
	}

	// 规范化路径
	normalizedPath := m.platform.NormalizePath(path)
	logger.Debug("Normalized GOROOT path: %s", normalizedPath)

	// 备份当前环境变量
	oldGoRoot := os.Getenv(constants.EnvGoRoot)

	// 根据平台设置环境变量
	err := m.setEnvPersistent(constants.EnvGoRoot, normalizedPath)
	if err != nil {
		logger.Error("Failed to set GOROOT: %v", err)
		return errors.Wrap(err, "ENVIRONMENT_SETUP_FAILED", "failed to set GOROOT").
			WithContext("old_value", oldGoRoot).
			WithContext("new_value", normalizedPath)
	}
	
	logger.Info("GOROOT set successfully")
	return nil
}

// SetGoPath 设置 GOPATH 环境变量
func (m *manager) SetGoPath(path string) error {
	if path == "" {
		return errors.ErrInvalidInput.WithMessage("GOPATH cannot be empty")
	}

	// 规范化路径
	normalizedPath := m.platform.NormalizePath(path)

	// 根据平台设置环境变量
	return m.setEnvPersistent(constants.EnvGoPath, normalizedPath)
}

// UpdatePath 更新 PATH 环境变量，添加 Go bin 目录
func (m *manager) UpdatePath(goRoot string) error {
	if goRoot == "" {
		return errors.ErrInvalidInput.WithMessage("GOROOT cannot be empty")
	}

	// 构建 Go bin 目录路径
	goBinPath := filepath.Join(goRoot, "bin")

	// 验证 bin 目录存在
	if _, err := os.Stat(goBinPath); os.IsNotExist(err) {
		return errors.ErrInvalidInput.
			WithMessage(fmt.Sprintf("Go bin directory does not exist: %s", goBinPath)).
			WithContext("go_root", goRoot).
			WithContext("bin_path", goBinPath)
	}

	// 获取当前 PATH
	currentPath := os.Getenv(constants.EnvPath)

	// 检查是否已经在 PATH 中
	pathSeparator := m.getPathSeparator()
	paths := strings.Split(currentPath, pathSeparator)
	
	// 移除旧的 Go bin 路径（可能有多个）
	var newPaths []string
	for _, p := range paths {
		normalizedP := m.platform.NormalizePath(p)
		// 跳过包含 .gx/versions 的路径
		if !strings.Contains(normalizedP, constants.DefaultInstallDir) {
			newPaths = append(newPaths, p)
		}
	}

	// 将新的 Go bin 路径添加到开头
	newPaths = append([]string{goBinPath}, newPaths...)
	newPath := strings.Join(newPaths, pathSeparator)

	// 持久化 PATH 更新
	if err := m.setEnvPersistent(constants.EnvPath, newPath); err != nil {
		return errors.Wrap(err, "ENVIRONMENT_SETUP_FAILED", "failed to update PATH").
			WithContext("old_path", currentPath).
			WithContext("new_path", newPath).
			WithContext("go_bin", goBinPath)
	}
	
	return nil
}

// GetGoRoot 获取当前 GOROOT
func (m *manager) GetGoRoot() (string, error) {
	goRoot := os.Getenv(constants.EnvGoRoot)
	if goRoot == "" {
		return "", errors.ErrNotFound.WithMessage("GOROOT not set")
	}
	return goRoot, nil
}

// GetGoPath 获取当前 GOPATH
func (m *manager) GetGoPath() (string, error) {
	goPath := os.Getenv(constants.EnvGoPath)
	if goPath == "" {
		// GOPATH 有默认值
		homeDir, err := m.platform.GetHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homeDir, "go"), nil
	}
	return goPath, nil
}

// Backup 备份当前环境变量配置
func (m *manager) Backup() error {
	return backupEnvironment(m.backupPath)
}

// Restore 恢复环境变量配置
func (m *manager) Restore() error {
	return restoreEnvironment(m.backupPath)
}

// getPathSeparator 获取 PATH 分隔符
func (m *manager) getPathSeparator() string {
	if m.platform.GetOS() == constants.OSWindows {
		return ";"
	}
	return ":"
}

// setEnvPersistent 持久化设置环境变量（平台特定实现）
func (m *manager) setEnvPersistent(key, value string) error {
	// 首先设置当前进程的环境变量
	if err := os.Setenv(key, value); err != nil {
		return errors.ErrOperationFailed.WithCause(err).WithMessage(fmt.Sprintf("failed to set %s", key))
	}

	// 根据平台持久化环境变量
	switch m.platform.GetOS() {
	case constants.OSWindows:
		return m.setEnvWindows(key, value)
	case constants.OSLinux, constants.OSDarwin:
		return m.setEnvUnix(key, value)
	default:
		return errors.ErrPlatformNotSupported.WithMessage(fmt.Sprintf("unsupported platform: %s", m.platform.GetOS()))
	}
}
