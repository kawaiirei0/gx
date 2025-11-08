//go:build linux || darwin

package platform

import (
	"os"
	"syscall"
)

// IsExecutable 检查文件是否可执行
func (a *adapter) IsExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// 检查文件模式是否包含可执行位
	mode := info.Mode()
	return mode&0111 != 0 // 检查任何可执行位（用户、组、其他）
}

// MakeExecutable 设置文件为可执行
func (a *adapter) MakeExecutable(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// 获取当前权限并添加可执行位
	currentMode := info.Mode()
	newMode := currentMode | 0111 // 添加所有用户的可执行权限

	// 设置新的文件权限
	return os.Chmod(path, newMode)
}

// getFileOwner 获取文件所有者（Unix 特定）
func getFileOwner(path string) (uid, gid int, err error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, 0, err
	}

	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, 0, nil
	}

	return int(stat.Uid), int(stat.Gid), nil
}
