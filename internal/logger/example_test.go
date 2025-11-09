package logger_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kawaiirei0/gx/internal/logger"
)

// ExampleLogger 演示日志记录器的基本使用
func ExampleLogger() {
	// 初始化日志记录器
	if err := logger.Init(); err != nil {
		fmt.Printf("Warning: failed to initialize logger: %v\n", err)
	}
	defer logger.Close()

	// 记录不同级别的日志
	logger.Info("Application started")
	logger.Debug("Debug information: %s", "some debug data")
	logger.Warn("This is a warning message")
	logger.Error("An error occurred: %v", fmt.Errorf("example error"))

	// 输出日志文件位置
	logPath, _ := logger.GetLogFilePath()
	fmt.Printf("Logs are written to: %s\n", logPath)
}

// ExampleLogger_withCustomLevel 演示设置日志级别
func ExampleLogger_withCustomLevel() {
	if err := logger.Init(); err != nil {
		fmt.Printf("Warning: failed to initialize logger: %v\n", err)
	}
	defer logger.Close()

	// 设置为 DEBUG 级别，记录所有日志
	logger.SetLevel(logger.LevelDebug)
	logger.Debug("This debug message will be logged")

	// 设置为 ERROR 级别，只记录错误
	logger.SetLevel(logger.LevelError)
	logger.Info("This info message will NOT be logged")
	logger.Error("This error message will be logged")
}

// ExampleNewLogger 演示创建独立的日志记录器实例
func ExampleNewLogger() {
	// 创建自定义日志记录器
	customLogger, err := logger.NewLogger()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}
	defer customLogger.Close()

	// 使用自定义日志记录器
	customLogger.Info("Custom logger message")
	customLogger.SetLevel(logger.LevelWarn)
	customLogger.Warn("Warning from custom logger")
}

// ExampleEnsureLogDir 演示确保日志目录存在
func ExampleEnsureLogDir() {
	// 确保日志目录存在
	if err := logger.EnsureLogDir(); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
		return
	}

	// 获取日志目录路径
	logDir, err := logger.GetLogDir()
	if err != nil {
		fmt.Printf("Failed to get log directory: %v\n", err)
		return
	}

	// 检查目录是否存在
	if info, err := os.Stat(logDir); err == nil && info.IsDir() {
		fmt.Printf("Log directory exists: %s\n", logDir)
	}
}

// ExampleGetLogFilePath 演示获取日志文件路径
func ExampleGetLogFilePath() {
	logPath, err := logger.GetLogFilePath()
	if err != nil {
		fmt.Printf("Failed to get log file path: %v\n", err)
		return
	}

	fmt.Printf("Log file path: %s\n", filepath.Base(logPath))
	// Output: Log file path: gx.log
}
