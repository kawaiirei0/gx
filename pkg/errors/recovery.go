package errors

import (
	"fmt"
	"os"
	"path/filepath"
)

// RecoveryManager 管理错误恢复和清理操作
type RecoveryManager struct {
	cleanupFuncs []CleanupFunc
	rollbackFuncs []RollbackFunc
}

// CleanupFunc 清理函数类型
type CleanupFunc func() error

// RollbackFunc 回滚函数类型
type RollbackFunc func() error

// NewRecoveryManager 创建新的恢复管理器
func NewRecoveryManager() *RecoveryManager {
	return &RecoveryManager{
		cleanupFuncs: make([]CleanupFunc, 0),
		rollbackFuncs: make([]RollbackFunc, 0),
	}
}

// AddCleanup 添加清理函数
func (rm *RecoveryManager) AddCleanup(fn CleanupFunc) {
	rm.cleanupFuncs = append(rm.cleanupFuncs, fn)
}

// AddRollback 添加回滚函数
func (rm *RecoveryManager) AddRollback(fn RollbackFunc) {
	rm.rollbackFuncs = append(rm.rollbackFuncs, fn)
}

// Cleanup 执行所有清理函数
func (rm *RecoveryManager) Cleanup() error {
	var errs []error
	
	// 按照添加的逆序执行清理（后进先出）
	for i := len(rm.cleanupFuncs) - 1; i >= 0; i-- {
		if err := rm.cleanupFuncs[i](); err != nil {
			errs = append(errs, err)
		}
	}
	
	if len(errs) > 0 {
		return ErrCleanupFailed.WithMessage(fmt.Sprintf("%d cleanup operations failed", len(errs)))
	}
	
	return nil
}

// Rollback 执行所有回滚函数
func (rm *RecoveryManager) Rollback() error {
	var errs []error
	
	// 按照添加的逆序执行回滚（后进先出）
	for i := len(rm.rollbackFuncs) - 1; i >= 0; i-- {
		if err := rm.rollbackFuncs[i](); err != nil {
			errs = append(errs, err)
		}
	}
	
	if len(errs) > 0 {
		return ErrRecoveryFailed.WithMessage(fmt.Sprintf("%d rollback operations failed", len(errs)))
	}
	
	return nil
}

// CleanupAndRollback 执行清理和回滚
func (rm *RecoveryManager) CleanupAndRollback() error {
	cleanupErr := rm.Cleanup()
	rollbackErr := rm.Rollback()
	
	if cleanupErr != nil && rollbackErr != nil {
		return ErrPartialFailure.WithMessage("both cleanup and rollback failed").
			WithContext("cleanup_error", cleanupErr).
			WithContext("rollback_error", rollbackErr)
	}
	
	if cleanupErr != nil {
		return cleanupErr
	}
	
	return rollbackErr
}

// Clear 清空所有注册的函数
func (rm *RecoveryManager) Clear() {
	rm.cleanupFuncs = make([]CleanupFunc, 0)
	rm.rollbackFuncs = make([]RollbackFunc, 0)
}

// SafeRemoveAll 安全地删除目录，带错误处理
func SafeRemoveAll(path string) error {
	if path == "" {
		return ErrInvalidInput.WithMessage("path cannot be empty")
	}
	
	// 验证路径存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 路径不存在，认为清理成功
		return nil
	}
	
	// 删除目录
	if err := os.RemoveAll(path); err != nil {
		return ErrCleanupFailed.WithCause(err).
			WithMessage(fmt.Sprintf("failed to remove directory: %s", path)).
			WithContext("path", path)
	}
	
	return nil
}

// SafeRemoveFile 安全地删除文件，带错误处理
func SafeRemoveFile(path string) error {
	if path == "" {
		return ErrInvalidInput.WithMessage("path cannot be empty")
	}
	
	// 验证文件存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 文件不存在，认为清理成功
		return nil
	}
	
	// 删除文件
	if err := os.Remove(path); err != nil {
		return ErrCleanupFailed.WithCause(err).
			WithMessage(fmt.Sprintf("failed to remove file: %s", path)).
			WithContext("path", path)
	}
	
	return nil
}

// EnsureDirectoryCleanup 确保目录在操作失败时被清理
func EnsureDirectoryCleanup(rm *RecoveryManager, dirPath string) {
	rm.AddCleanup(func() error {
		return SafeRemoveAll(dirPath)
	})
}

// EnsureFileCleanup 确保文件在操作失败时被清理
func EnsureFileCleanup(rm *RecoveryManager, filePath string) {
	rm.AddCleanup(func() error {
		return SafeRemoveFile(filePath)
	})
}

// BackupFile 备份文件
func BackupFile(srcPath string) (backupPath string, err error) {
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return "", ErrNotFound.WithMessage(fmt.Sprintf("source file not found: %s", srcPath))
	}
	
	// 创建备份文件路径
	backupPath = srcPath + ".backup"
	
	// 读取源文件
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return "", ErrOperationFailed.WithCause(err).WithMessage("failed to read source file")
	}
	
	// 写入备份文件
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return "", ErrOperationFailed.WithCause(err).WithMessage("failed to write backup file")
	}
	
	return backupPath, nil
}

// RestoreFile 从备份恢复文件
func RestoreFile(backupPath, targetPath string) error {
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return ErrNotFound.WithMessage(fmt.Sprintf("backup file not found: %s", backupPath))
	}
	
	// 读取备份文件
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return ErrRecoveryFailed.WithCause(err).WithMessage("failed to read backup file")
	}
	
	// 写入目标文件
	if err := os.WriteFile(targetPath, data, 0644); err != nil {
		return ErrRecoveryFailed.WithCause(err).WithMessage("failed to restore file")
	}
	
	return nil
}

// CreateTempDir 创建临时目录，并注册清理函数
func CreateTempDir(rm *RecoveryManager, prefix string) (string, error) {
	tmpDir, err := os.MkdirTemp("", prefix)
	if err != nil {
		return "", ErrOperationFailed.WithCause(err).WithMessage("failed to create temp directory")
	}
	
	// 注册清理函数
	if rm != nil {
		EnsureDirectoryCleanup(rm, tmpDir)
	}
	
	return tmpDir, nil
}

// MoveWithRollback 移动文件/目录，支持回滚
func MoveWithRollback(rm *RecoveryManager, src, dst string) error {
	// 检查源是否存在
	srcInfo, err := os.Stat(src)
	if err != nil {
		return ErrNotFound.WithCause(err).WithMessage(fmt.Sprintf("source not found: %s", src))
	}
	
	// 如果目标已存在，先备份
	var dstBackup string
	if _, err := os.Stat(dst); err == nil {
		dstBackup = dst + ".rollback"
		if err := os.Rename(dst, dstBackup); err != nil {
			return ErrOperationFailed.WithCause(err).WithMessage("failed to backup destination")
		}
		
		// 注册回滚函数：恢复目标
		if rm != nil {
			rm.AddRollback(func() error {
				if _, err := os.Stat(dstBackup); err == nil {
					return os.Rename(dstBackup, dst)
				}
				return nil
			})
		}
	}
	
	// 确保目标目录存在
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return ErrOperationFailed.WithCause(err).WithMessage("failed to create destination directory")
	}
	
	// 移动文件/目录
	if err := os.Rename(src, dst); err != nil {
		// 如果移动失败，尝试恢复目标备份
		if dstBackup != "" {
			os.Rename(dstBackup, dst)
		}
		return ErrOperationFailed.WithCause(err).WithMessage("failed to move")
	}
	
	// 注册回滚函数：移回源位置
	if rm != nil {
		rm.AddRollback(func() error {
			if _, err := os.Stat(dst); err == nil {
				return os.Rename(dst, src)
			}
			return nil
		})
	}
	
	// 清理目标备份
	if dstBackup != "" {
		if rm != nil {
			rm.AddCleanup(func() error {
				return SafeRemoveAll(dstBackup)
			})
		} else {
			os.RemoveAll(dstBackup)
		}
	}
	
	// 如果是目录，验证移动成功
	if srcInfo.IsDir() {
		if _, err := os.Stat(dst); err != nil {
			return ErrOperationFailed.WithMessage("move verification failed")
		}
	}
	
	return nil
}
