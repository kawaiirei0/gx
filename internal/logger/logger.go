package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/yourusername/gx/pkg/constants"
)

// LogLevel 日志级别
type LogLevel int

const (
	// LevelDebug 调试级别
	LevelDebug LogLevel = iota
	// LevelInfo 信息级别
	LevelInfo
	// LevelWarn 警告级别
	LevelWarn
	// LevelError 错误级别
	LevelError
)

// String 返回日志级别的字符串表示
func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger 日志记录器接口
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	SetLevel(level LogLevel)
	Close() error
}

// fileLogger 基于文件的日志记录器实现
type fileLogger struct {
	mu       sync.Mutex
	file     *io.Writer
	level    LogLevel
	logPath  string
	disabled bool
}

var (
	// defaultLogger 默认日志记录器实例
	defaultLogger Logger
	once          sync.Once
)

// Init 初始化日志记录器
func Init() error {
	var err error
	once.Do(func() {
		defaultLogger, err = NewLogger()
	})
	return err
}

// NewLogger 创建新的日志记录器
func NewLogger() (Logger, error) {
	logPath, err := GetLogFilePath()
	if err != nil {
		// 如果无法获取日志路径，创建一个禁用的日志记录器
		return &fileLogger{
			disabled: true,
			level:    LevelInfo,
		}, nil
	}

	// 确保日志目录存在
	if err := EnsureLogDir(); err != nil {
		// 如果无法创建日志目录，创建一个禁用的日志记录器
		return &fileLogger{
			disabled: true,
			level:    LevelInfo,
		}, nil
	}

	// 打开日志文件（追加模式）
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		// 如果无法打开日志文件，创建一个禁用的日志记录器
		return &fileLogger{
			disabled: true,
			level:    LevelInfo,
		}, nil
	}

	writer := io.Writer(file)
	return &fileLogger{
		file:     &writer,
		level:    LevelInfo,
		logPath:  logPath,
		disabled: false,
	}, nil
}

// log 写入日志
func (l *fileLogger) log(level LogLevel, format string, args ...interface{}) {
	if l.disabled {
		return
	}

	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level.String(), message)

	if l.file != nil {
		(*l.file).Write([]byte(logLine))
	}
}

// Debug 记录调试级别日志
func (l *fileLogger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, format, args...)
}

// Info 记录信息级别日志
func (l *fileLogger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, format, args...)
}

// Warn 记录警告级别日志
func (l *fileLogger) Warn(format string, args ...interface{}) {
	l.log(LevelWarn, format, args...)
}

// Error 记录错误级别日志
func (l *fileLogger) Error(format string, args ...interface{}) {
	l.log(LevelError, format, args...)
}

// SetLevel 设置日志级别
func (l *fileLogger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// Close 关闭日志记录器
func (l *fileLogger) Close() error {
	if l.disabled || l.file == nil {
		return nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if closer, ok := (*l.file).(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// GetLogDir 获取日志目录路径
func GetLogDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, constants.ConfigDir, "logs"), nil
}

// GetLogFilePath 获取日志文件路径
func GetLogFilePath() (string, error) {
	logDir, err := GetLogDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(logDir, "gx.log"), nil
}

// EnsureLogDir 确保日志目录存在
func EnsureLogDir() error {
	logDir, err := GetLogDir()
	if err != nil {
		return err
	}

	// 检查目录是否存在
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		// 创建目录，权限 0755
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}
	}

	return nil
}

// 全局日志函数，使用默认日志记录器

// Debug 记录调试级别日志
func Debug(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debug(format, args...)
	}
}

// Info 记录信息级别日志
func Info(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(format, args...)
	}
}

// Warn 记录警告级别日志
func Warn(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warn(format, args...)
	}
}

// Error 记录错误级别日志
func Error(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(format, args...)
	}
}

// SetLevel 设置默认日志记录器的日志级别
func SetLevel(level LogLevel) {
	if defaultLogger != nil {
		defaultLogger.SetLevel(level)
	}
}

// Close 关闭默认日志记录器
func Close() error {
	if defaultLogger != nil {
		return defaultLogger.Close()
	}
	return nil
}

// GetDefaultLogger 获取默认日志记录器
func GetDefaultLogger() Logger {
	if defaultLogger == nil {
		Init()
	}
	return defaultLogger
}
