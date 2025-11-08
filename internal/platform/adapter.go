package platform

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/yourusername/gx/pkg/interfaces"
)

// adapter 实现 PlatformAdapter 接口
type adapter struct {
	os   string
	arch string
}

// NewAdapter 创建新的平台适配器
func NewAdapter() interfaces.PlatformAdapter {
	return &adapter{
		os:   runtime.GOOS,
		arch: runtime.GOARCH,
	}
}

// GetOS 获取操作系统类型
func (a *adapter) GetOS() string {
	return a.os
}

// GetArch 获取架构类型
func (a *adapter) GetArch() string {
	return a.arch
}

// PathSeparator 获取路径分隔符
func (a *adapter) PathSeparator() string {
	return string(os.PathSeparator)
}

// NormalizePath 规范化路径
func (a *adapter) NormalizePath(path string) string {
	// 使用 filepath.Clean 清理路径
	cleaned := filepath.Clean(path)
	
	// 转换为绝对路径（如果可能）
	if abs, err := filepath.Abs(cleaned); err == nil {
		return abs
	}
	
	return cleaned
}

// JoinPath 连接路径
func (a *adapter) JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}

// GetHomeDir 获取用户主目录
func (a *adapter) GetHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home, nil
}
