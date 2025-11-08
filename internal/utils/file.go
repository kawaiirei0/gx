package utils

import (
	"os"
	"path/filepath"
)

// EnsureDir 确保目录存在，如果不存在则创建
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirExists 检查目录是否存在
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// RemoveDir 删除目录及其内容
func RemoveDir(path string) error {
	return os.RemoveAll(path)
}

// GetHomeDir 获取用户主目录
func GetHomeDir() (string, error) {
	return os.UserHomeDir()
}

// JoinPath 连接路径
func JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}
