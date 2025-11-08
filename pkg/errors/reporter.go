package errors

import (
	"fmt"
	"strings"
)

// ErrorReporter 提供用户友好的错误报告
type ErrorReporter struct {
	verbose bool
}

// NewErrorReporter 创建新的错误报告器
func NewErrorReporter(verbose bool) *ErrorReporter {
	return &ErrorReporter{
		verbose: verbose,
	}
}

// Report 生成用户友好的错误报告
func (r *ErrorReporter) Report(err error) string {
	if err == nil {
		return ""
	}

	var builder strings.Builder

	// 获取自定义错误
	if customErr, ok := err.(*Error); ok {
		builder.WriteString(fmt.Sprintf("Error: %s\n", customErr.Message))

		// 如果是详细模式，显示错误代码
		if r.verbose {
			builder.WriteString(fmt.Sprintf("Code: %s\n", customErr.Code))
		}

		// 显示上下文信息
		if len(customErr.Context) > 0 && r.verbose {
			builder.WriteString("Context:\n")
			for key, value := range customErr.Context {
				builder.WriteString(fmt.Sprintf("  %s: %v\n", key, value))
			}
		}

		// 显示建议的解决方案
		if suggestion := r.getSuggestion(customErr); suggestion != "" {
			builder.WriteString(fmt.Sprintf("\nSuggestion: %s\n", suggestion))
		}

		// 如果是详细模式，显示错误链
		if r.verbose && customErr.Cause != nil {
			builder.WriteString(fmt.Sprintf("\nCaused by: %v\n", customErr.Cause))
			
			// 显示完整的错误链
			rootCause := GetRootCause(err)
			if rootCause != customErr.Cause {
				builder.WriteString(fmt.Sprintf("Root cause: %v\n", rootCause))
			}
		}
	} else {
		// 标准错误
		builder.WriteString(fmt.Sprintf("Error: %v\n", err))
	}

	return builder.String()
}

// getSuggestion 根据错误类型返回建议的解决方案
func (r *ErrorReporter) getSuggestion(err *Error) string {
	switch err.Code {
	case "VERSION_NOT_FOUND":
		return "Use 'gx list' to see installed versions or 'gx install <version>' to install a new version."
	
	case "VERSION_ALREADY_INSTALLED":
		return "The version is already installed. Use 'gx use <version>' to switch to it."
	
	case "VERSION_NOT_INSTALLED":
		return "Install the version first using 'gx install <version>'."
	
	case "DOWNLOAD_FAILED":
		return "Check your internet connection and try again. You can also check if the version exists using 'gx list --remote'."
	
	case "INSTALL_FAILED":
		return "Ensure you have sufficient disk space and permissions. Check the logs for more details."
	
	case "CHECKSUM_MISMATCH":
		return "The downloaded file may be corrupted. Try downloading again."
	
	case "NETWORK_ERROR":
		return "Check your internet connection and proxy settings."
	
	case "PERMISSION_DENIED":
		return "Try running the command with administrator/sudo privileges."
	
	case "ENVIRONMENT_SETUP_FAILED":
		return "You may need to restart your terminal or run 'source ~/.bashrc' (Linux/macOS) or restart your command prompt (Windows)."
	
	case "PLATFORM_NOT_SUPPORTED":
		return "Check the list of supported platforms using 'gx cross-build --help'."
	
	case "CONFIG_CORRUPTED":
		return "Your configuration file may be corrupted. Try removing ~/.gx/config.json and running the command again."
	
	case "DISK_SPACE_INSUFFICIENT":
		return "Free up some disk space and try again."
	
	case "ACTIVE_VERSION_IN_USE":
		return "Switch to a different version first using 'gx use <version>' before uninstalling."
	
	default:
		return "Check the logs for more details or run with --verbose flag for more information."
	}
}

// ReportWithRecovery 报告错误并提供恢复建议
func (r *ErrorReporter) ReportWithRecovery(err error) string {
	report := r.Report(err)
	
	// 检查错误是否可恢复
	if customErr, ok := err.(*Error); ok && customErr.IsRecoverable() {
		report += "\nThis error may be recoverable. Try the suggested solution above.\n"
	}
	
	return report
}

// FormatErrorChain 格式化错误链
func FormatErrorChain(err error) string {
	if err == nil {
		return ""
	}

	var builder strings.Builder
	builder.WriteString(err.Error())

	current := err
	depth := 0
	maxDepth := 10 // 防止无限循环

	for depth < maxDepth {
		unwrapped := Unwrap(current)
		if unwrapped == nil {
			break
		}
		
		builder.WriteString("\n  → ")
		builder.WriteString(unwrapped.Error())
		
		current = unwrapped
		depth++
	}

	return builder.String()
}

// Unwrap 解包错误
func Unwrap(err error) error {
	type unwrapper interface {
		Unwrap() error
	}

	if u, ok := err.(unwrapper); ok {
		return u.Unwrap()
	}
	return nil
}

// IsRecoverableError 检查错误是否可恢复
func IsRecoverableError(err error) bool {
	if customErr, ok := err.(*Error); ok {
		return customErr.IsRecoverable()
	}
	return false
}

// GetErrorContext 获取错误的所有上下文信息
func GetErrorContext(err error) map[string]interface{} {
	if customErr, ok := err.(*Error); ok {
		return customErr.Context
	}
	return nil
}
