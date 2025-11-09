package wrapper

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kawaiirei0/gx/internal/logger"
	"github.com/kawaiirei0/gx/pkg/constants"
	"github.com/kawaiirei0/gx/pkg/errors"
	"github.com/kawaiirei0/gx/pkg/interfaces"
)

// cliWrapper 实现 CLIWrapper 接口
type cliWrapper struct {
	versionManager interfaces.VersionManager
	platform       interfaces.PlatformAdapter
}

// NewCLIWrapper 创建新的 CLI 包装器
func NewCLIWrapper(versionManager interfaces.VersionManager, platform interfaces.PlatformAdapter) interfaces.CLIWrapper {
	return &cliWrapper{
		versionManager: versionManager,
		platform:       platform,
	}
}

// Execute 执行 Go 命令
func (w *cliWrapper) Execute(command string, args []string) error {
	logger.Info("Executing Go command: %s %v", command, args)
	
	// 获取 Go 可执行文件路径
	goExe, err := w.GetGoExecutable()
	if err != nil {
		logger.Error("Failed to get Go executable: %v", err)
		return errors.Wrap(err, "OPERATION_FAILED", "failed to get Go executable").
			WithContext("command", command).
			WithContext("args", args)
	}
	logger.Debug("Using Go executable: %s", goExe)

	// 构建完整的命令参数
	// 第一个参数是 Go 子命令（如 run, build, test）
	cmdArgs := append([]string{command}, args...)

	// 创建命令
	cmd := exec.Command(goExe, cmdArgs...)

	// 透传标准输入、输出、错误流
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	if err := cmd.Run(); err != nil {
		// 保留原始退出码
		if exitErr, ok := err.(*exec.ExitError); ok {
			// 命令执行失败，但这是预期的行为（如编译错误）
			// 退出码已经通过 Stderr 传递给用户
			// 我们返回一个包装的错误，但保留退出码信息
			return &ExitError{
				ExitCode: exitErr.ExitCode(),
				Err:      exitErr,
				Command:  command,
				Args:     args,
			}
		}
		// 其他类型的错误（如命令未找到）
		return errors.ErrOperationFailed.
			WithCause(err).
			WithMessage("failed to execute go command").
			WithContext("command", command).
			WithContext("args", args).
			WithContext("go_executable", goExe)
	}

	return nil
}

// GetGoExecutable 获取当前使用的 Go 可执行文件路径
func (w *cliWrapper) GetGoExecutable() (string, error) {
	// 获取当前激活的版本
	activeVersion, err := w.versionManager.GetActive()
	if err != nil {
		return "", errors.ErrVersionNotFound.WithCause(err).WithMessage("no active Go version found")
	}

	// 验证版本路径是否存在
	if activeVersion.Path == "" {
		return "", errors.ErrVersionNotFound.WithMessage("active version path is empty")
	}

	// 构建 Go 可执行文件路径
	goExe := "go"
	if w.platform.GetOS() == constants.OSWindows {
		goExe = "go.exe"
	}

	goPath := filepath.Join(activeVersion.Path, "bin", goExe)

	// 验证文件是否存在
	if _, err := os.Stat(goPath); os.IsNotExist(err) {
		return "", errors.ErrNotFound.WithMessage(fmt.Sprintf("go executable not found at %s", goPath))
	}

	// 验证文件是否可执行（Unix 系统）
	if w.platform.GetOS() != constants.OSWindows {
		if !w.platform.IsExecutable(goPath) {
			return "", errors.ErrPermissionDenied.WithMessage(fmt.Sprintf("go executable at %s is not executable", goPath))
		}
	}

	return goPath, nil
}

// ExitError 包装退出错误，保留退出码
type ExitError struct {
	ExitCode int
	Err      error
	Command  string
	Args     []string
}

// Error 实现 error 接口
func (e *ExitError) Error() string {
	if e.Command != "" {
		return fmt.Sprintf("command 'go %s' exited with code %d", e.Command, e.ExitCode)
	}
	return fmt.Sprintf("command exited with code %d", e.ExitCode)
}

// Unwrap 返回原始错误
func (e *ExitError) Unwrap() error {
	return e.Err
}

// GetExitCode 获取退出码
func (e *ExitError) GetExitCode() int {
	return e.ExitCode
}

// IsExitError 检查错误是否为退出错误
func IsExitError(err error) (*ExitError, bool) {
	if exitErr, ok := err.(*ExitError); ok {
		return exitErr, true
	}
	return nil, false
}
