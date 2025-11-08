package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/yourusername/gx/pkg/errors"
)

// ErrorFormatter 错误格式化器
type ErrorFormatter struct {
	writer io.Writer
}

// NewErrorFormatter 创建新的错误格式化器
func NewErrorFormatter(writer io.Writer) *ErrorFormatter {
	return &ErrorFormatter{
		writer: writer,
	}
}

// Format 格式化并显示错误
func (ef *ErrorFormatter) Format(err error) {
	if err == nil {
		return
	}

	// 检查是否是自定义错误类型
	if gxErr, ok := err.(*errors.Error); ok {
		ef.formatGxError(gxErr)
		return
	}

	// 普通错误
	fmt.Fprintf(ef.writer, "✗ Error: %s\n", err.Error())
}

// formatGxError 格式化 gx 自定义错误
func (ef *ErrorFormatter) formatGxError(err *errors.Error) {
	fmt.Fprintln(ef.writer)
	fmt.Fprintf(ef.writer, "✗ %s\n", err.Message)

	// 显示错误代码
	if err.Code != "" {
		fmt.Fprintf(ef.writer, "  Error Code: %s\n", err.Code)
	}

	// 显示原因
	if err.Cause != nil {
		fmt.Fprintf(ef.writer, "  Cause: %s\n", err.Cause.Error())
	}

	// 根据错误类型提供建议
	ef.provideSuggestions(err)

	fmt.Fprintln(ef.writer)
}

// provideSuggestions 根据错误类型提供解决建议
func (ef *ErrorFormatter) provideSuggestions(err *errors.Error) {
	suggestions := ef.getSuggestions(err)
	if len(suggestions) == 0 {
		return
	}

	fmt.Fprintln(ef.writer)
	fmt.Fprintln(ef.writer, "Suggestions:")
	for _, suggestion := range suggestions {
		fmt.Fprintf(ef.writer, "  • %s\n", suggestion)
	}
}

// getSuggestions 获取错误建议
func (ef *ErrorFormatter) getSuggestions(err *errors.Error) []string {
	var suggestions []string

	switch {
	case strings.Contains(err.Code, "VERSION_NOT_FOUND"):
		suggestions = append(suggestions,
			"Check if the version number is correct",
			"Run 'gx list --remote' to see available versions",
			"Make sure you have internet connection",
		)

	case strings.Contains(err.Code, "VERSION_NOT_INSTALLED"):
		suggestions = append(suggestions,
			"Install the version first using 'gx install <version>'",
			"Run 'gx list' to see installed versions",
		)

	case strings.Contains(err.Code, "VERSION_ALREADY_INSTALLED"):
		suggestions = append(suggestions,
			"The version is already installed",
			"Use 'gx use <version>' to switch to it",
			"Use 'gx uninstall <version>' to remove it first if you want to reinstall",
		)

	case strings.Contains(err.Code, "NETWORK_ERROR"):
		suggestions = append(suggestions,
			"Check your internet connection",
			"Verify that https://go.dev is accessible",
			"Try again later if the service is temporarily unavailable",
		)

	case strings.Contains(err.Code, "DOWNLOAD_FAILED"):
		suggestions = append(suggestions,
			"Check your internet connection",
			"Ensure you have enough disk space",
			"Try downloading again",
		)

	case strings.Contains(err.Code, "CHECKSUM_MISMATCH"):
		suggestions = append(suggestions,
			"The downloaded file may be corrupted",
			"Try downloading again",
			"Check your internet connection stability",
		)

	case strings.Contains(err.Code, "INSTALL_FAILED"):
		suggestions = append(suggestions,
			"Ensure you have write permissions to the installation directory",
			"Check if you have enough disk space",
			"Try running with administrator/sudo privileges if needed",
		)

	case strings.Contains(err.Code, "UNINSTALL_FAILED"):
		suggestions = append(suggestions,
			"Make sure the version is not currently in use",
			"Close any programs using this Go version",
			"Try running with administrator/sudo privileges if needed",
		)

	case strings.Contains(err.Code, "ENVIRONMENT_SETUP_FAILED"):
		suggestions = append(suggestions,
			"Check if you have permissions to modify environment variables",
			"On Windows, try running as administrator",
			"On Linux/macOS, check your shell configuration files",
			"You may need to restart your terminal for changes to take effect",
		)

	case strings.Contains(err.Code, "STORAGE_FAILED"):
		suggestions = append(suggestions,
			"Check if the configuration directory (~/.gx) is writable",
			"Ensure you have enough disk space",
			"Verify file permissions",
		)

	case strings.Contains(err.Code, "PLATFORM_NOT_SUPPORTED"):
		suggestions = append(suggestions,
			"Your platform may not be supported by this Go version",
			"Check the Go release notes for platform support",
			"Try a different version",
		)

	default:
		suggestions = append(suggestions,
			"Run with --verbose flag for more details",
			"Check the logs at ~/.gx/logs/gx.log",
		)
	}

	return suggestions
}

// FormatHelp 格式化帮助信息
func (ef *ErrorFormatter) FormatHelp(command string, description string, examples []string) {
	fmt.Fprintln(ef.writer)
	fmt.Fprintf(ef.writer, "Command: %s\n", command)
	fmt.Fprintf(ef.writer, "%s\n", description)

	if len(examples) > 0 {
		fmt.Fprintln(ef.writer)
		fmt.Fprintln(ef.writer, "Examples:")
		for _, example := range examples {
			fmt.Fprintf(ef.writer, "  %s\n", example)
		}
	}

	fmt.Fprintln(ef.writer)
}
