package platform

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/yourusername/gx/pkg/constants"
)

// GetExecutableExtension 获取当前平台的可执行文件扩展名
func GetExecutableExtension() string {
	if runtime.GOOS == constants.OSWindows {
		return ".exe"
	}
	return ""
}

// GetArchiveExtension 获取当前平台的压缩包扩展名
func GetArchiveExtension() string {
	if runtime.GOOS == constants.OSWindows {
		return constants.ArchiveExtZip
	}
	return constants.ArchiveExtTarGz
}

// GetConfigDir 获取配置目录的完整路径
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, constants.ConfigDir), nil
}

// GetInstallDir 获取安装目录的完整路径
func GetInstallDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, constants.DefaultInstallDir), nil
}

// EnsureDir 确保目录存在，如果不存在则创建
func EnsureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// IsDirectory 检查路径是否为目录
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetPlatformString 获取平台字符串（格式：os/arch）
func GetPlatformString() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}

// IsSupportedPlatform 检查是否为支持的平台
func IsSupportedPlatform(os, arch string) bool {
	supportedPlatforms := map[string][]string{
		constants.OSWindows: {constants.ArchAMD64, constants.Arch386},
		constants.OSLinux:   {constants.ArchAMD64, constants.ArchARM64, constants.Arch386},
		constants.OSDarwin:  {constants.ArchAMD64, constants.ArchARM64},
	}

	archs, ok := supportedPlatforms[os]
	if !ok {
		return false
	}

	for _, a := range archs {
		if a == arch {
			return true
		}
	}
	return false
}
