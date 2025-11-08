package errors

import (
	"errors"
	"fmt"
)

// 错误类型定义
var (
	// ErrVersionNotFound 版本未找到
	ErrVersionNotFound = NewError("VERSION_NOT_FOUND", "version not found")

	// ErrVersionAlreadyInstalled 版本已安装
	ErrVersionAlreadyInstalled = NewError("VERSION_ALREADY_INSTALLED", "version already installed")

	// ErrVersionNotInstalled 版本未安装
	ErrVersionNotInstalled = NewError("VERSION_NOT_INSTALLED", "version not installed")

	// ErrDownloadFailed 下载失败
	ErrDownloadFailed = NewError("DOWNLOAD_FAILED", "download failed")

	// ErrInstallFailed 安装失败
	ErrInstallFailed = NewError("INSTALL_FAILED", "installation failed")

	// ErrUninstallFailed 卸载失败
	ErrUninstallFailed = NewError("UNINSTALL_FAILED", "uninstall failed")

	// ErrInvalidVersion 无效的版本号
	ErrInvalidVersion = NewError("INVALID_VERSION", "invalid version format")

	// ErrEnvironmentSetupFailed 环境变量设置失败
	ErrEnvironmentSetupFailed = NewError("ENVIRONMENT_SETUP_FAILED", "environment setup failed")

	// ErrStorageFailed 存储操作失败
	ErrStorageFailed = NewError("STORAGE_FAILED", "storage operation failed")

	// ErrNetworkError 网络错误
	ErrNetworkError = NewError("NETWORK_ERROR", "network error")

	// ErrPermissionDenied 权限不足
	ErrPermissionDenied = NewError("PERMISSION_DENIED", "permission denied")

	// ErrChecksumMismatch 校验和不匹配
	ErrChecksumMismatch = NewError("CHECKSUM_MISMATCH", "checksum verification failed")

	// ErrInvalidInput 无效的输入
	ErrInvalidInput = NewError("INVALID_INPUT", "invalid input")

	// ErrNotFound 资源未找到
	ErrNotFound = NewError("NOT_FOUND", "resource not found")

	// ErrOperationFailed 操作失败
	ErrOperationFailed = NewError("OPERATION_FAILED", "operation failed")

	// ErrPlatformNotSupported 平台不支持
	ErrPlatformNotSupported = NewError("PLATFORM_NOT_SUPPORTED", "platform not supported")

	// ErrConfigCorrupted 配置文件损坏
	ErrConfigCorrupted = NewError("CONFIG_CORRUPTED", "configuration file is corrupted")

	// ErrDiskSpaceInsufficient 磁盘空间不足
	ErrDiskSpaceInsufficient = NewError("DISK_SPACE_INSUFFICIENT", "insufficient disk space")

	// ErrArchiveCorrupted 压缩包损坏
	ErrArchiveCorrupted = NewError("ARCHIVE_CORRUPTED", "archive file is corrupted")

	// ErrTimeout 操作超时
	ErrTimeout = NewError("TIMEOUT", "operation timed out")

	// ErrCancelled 操作被取消
	ErrCancelled = NewError("CANCELLED", "operation was cancelled")

	// ErrActiveVersionInUse 当前激活版本正在使用
	ErrActiveVersionInUse = NewError("ACTIVE_VERSION_IN_USE", "cannot modify active version")

	// ErrCleanupFailed 清理失败
	ErrCleanupFailed = NewError("CLEANUP_FAILED", "cleanup operation failed")

	// ErrRecoveryFailed 恢复失败
	ErrRecoveryFailed = NewError("RECOVERY_FAILED", "recovery operation failed")

	// ErrPartialFailure 部分操作失败
	ErrPartialFailure = NewError("PARTIAL_FAILURE", "operation partially failed")
)

// Error 自定义错误类型
type Error struct {
	Code       string
	Message    string
	Cause      error
	Recoverable bool // 标记错误是否可恢复
	Context    map[string]interface{} // 错误上下文信息
}

// NewError 创建新的错误
func NewError(code, message string) *Error {
	return &Error{
		Code:       code,
		Message:    message,
		Recoverable: false,
		Context:    make(map[string]interface{}),
	}
}

// Error 实现 error 接口
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// WithCause 添加原因
func (e *Error) WithCause(cause error) *Error {
	return &Error{
		Code:       e.Code,
		Message:    e.Message,
		Cause:      cause,
		Recoverable: e.Recoverable,
		Context:    e.Context,
	}
}

// WithMessage 添加详细消息
func (e *Error) WithMessage(message string) *Error {
	return &Error{
		Code:       e.Code,
		Message:    fmt.Sprintf("%s: %s", e.Message, message),
		Cause:      e.Cause,
		Recoverable: e.Recoverable,
		Context:    e.Context,
	}
}

// WithContext 添加上下文信息
func (e *Error) WithContext(key string, value interface{}) *Error {
	newContext := make(map[string]interface{})
	for k, v := range e.Context {
		newContext[k] = v
	}
	newContext[key] = value
	
	return &Error{
		Code:       e.Code,
		Message:    e.Message,
		Cause:      e.Cause,
		Recoverable: e.Recoverable,
		Context:    newContext,
	}
}

// AsRecoverable 标记错误为可恢复
func (e *Error) AsRecoverable() *Error {
	return &Error{
		Code:       e.Code,
		Message:    e.Message,
		Cause:      e.Cause,
		Recoverable: true,
		Context:    e.Context,
	}
}

// IsRecoverable 检查错误是否可恢复
func (e *Error) IsRecoverable() bool {
	return e.Recoverable
}

// GetContext 获取上下文信息
func (e *Error) GetContext(key string) (interface{}, bool) {
	val, ok := e.Context[key]
	return val, ok
}

// Unwrap 返回原始错误
func (e *Error) Unwrap() error {
	return e.Cause
}

// Is 实现 errors.Is 接口
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// Wrap 包装标准错误为自定义错误
func Wrap(err error, code, message string) *Error {
	if err == nil {
		return nil
	}
	
	// 如果已经是自定义错误，保留原有信息
	if customErr, ok := err.(*Error); ok {
		return customErr.WithMessage(message)
	}
	
	return NewError(code, message).WithCause(err)
}

// IsType 检查错误是否为指定类型
func IsType(err error, target *Error) bool {
	return errors.Is(err, target)
}

// GetRootCause 获取错误链的根本原因
func GetRootCause(err error) error {
	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			return err
		}
		err = unwrapped
	}
}
