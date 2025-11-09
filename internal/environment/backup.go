package environment

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/kawaiirei0/gx/pkg/constants"
	"github.com/kawaiirei0/gx/pkg/errors"
)

// EnvBackup 环境变量备份结构
type EnvBackup struct {
	Timestamp time.Time         `json:"timestamp"`
	Variables map[string]string `json:"variables"`
}

// backupEnvironment 备份当前环境变量
func backupEnvironment(backupPath string) error {
	// 创建备份目录
	backupDir := filepath.Dir(backupPath)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return errors.ErrOperationFailed.WithCause(err).WithMessage("failed to create backup directory")
	}

	// 收集需要备份的环境变量
	backup := EnvBackup{
		Timestamp: time.Now(),
		Variables: make(map[string]string),
	}

	// 备份关键环境变量
	envVars := []string{
		constants.EnvGoRoot,
		constants.EnvGoPath,
		constants.EnvPath,
	}

	for _, key := range envVars {
		if value := os.Getenv(key); value != "" {
			backup.Variables[key] = value
		}
	}

	// 序列化为 JSON
	data, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		return errors.ErrOperationFailed.WithCause(err).WithMessage("failed to serialize backup")
	}

	// 写入备份文件
	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return errors.ErrOperationFailed.WithCause(err).WithMessage("failed to write backup file")
	}

	return nil
}

// restoreEnvironment 恢复环境变量
func restoreEnvironment(backupPath string) error {
	// 检查备份文件是否存在
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return errors.ErrNotFound.WithMessage("backup file not found")
	}

	// 读取备份文件
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return errors.ErrOperationFailed.WithCause(err).WithMessage("failed to read backup file")
	}

	// 解析 JSON
	var backup EnvBackup
	if err := json.Unmarshal(data, &backup); err != nil {
		return errors.ErrOperationFailed.WithCause(err).WithMessage("failed to parse backup file")
	}

	// 恢复环境变量
	for key, value := range backup.Variables {
		if err := os.Setenv(key, value); err != nil {
			return errors.ErrOperationFailed.WithCause(err).WithMessage("failed to restore environment variable: " + key)
		}
	}

	return nil
}

// listBackups 列出所有备份
func listBackups(backupDir string) ([]EnvBackup, error) {
	// 检查备份目录是否存在
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return []EnvBackup{}, nil
	}

	// 读取目录中的所有备份文件
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, errors.ErrOperationFailed.WithCause(err).WithMessage("failed to read backup directory")
	}

	var backups []EnvBackup
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// 只处理 JSON 文件
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		backupPath := filepath.Join(backupDir, entry.Name())
		data, err := os.ReadFile(backupPath)
		if err != nil {
			continue
		}

		var backup EnvBackup
		if err := json.Unmarshal(data, &backup); err != nil {
			continue
		}

		backups = append(backups, backup)
	}

	return backups, nil
}

// deleteBackup 删除备份文件
func deleteBackup(backupPath string) error {
	if err := os.Remove(backupPath); err != nil {
		return errors.ErrOperationFailed.WithCause(err).WithMessage("failed to delete backup file")
	}
	return nil
}
