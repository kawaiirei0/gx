//go:build windows

package platform

import (
	"os"
	"path/filepath"
	"strings"
)

// IsExecutable 检查文件是否可执行（Windows）
func (a *adapter) IsExecutable(path string) bool {
	// 在 Windows 上，检查文件扩展名
	ext := strings.ToLower(filepath.Ext(path))
	executableExts := []string{".exe", ".bat", ".cmd", ".com"}

	for _, execExt := range executableExts {
		if ext == execExt {
			// 同时检查文件是否存在
			if _, err := os.Stat(path); err == nil {
				return true
			}
		}
	}

	return false
}

// MakeExecutable 设置文件为可执行（Windows）
func (a *adapter) MakeExecutable(path string) error {
	// 在 Windows 上，可执行性由文件扩展名决定
	// 如果文件没有 .exe 扩展名，我们不需要做任何事情
	// 只需要确保文件存在且可读
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// 检查文件是否可读
	if info.Mode()&0400 == 0 {
		// 如果文件不可读，尝试添加读权限
		return os.Chmod(path, info.Mode()|0400)
	}

	return nil
}
