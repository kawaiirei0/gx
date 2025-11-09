# Logger Package

日志记录包，提供统一的日志记录功能。

## 功能特性

- 支持多个日志级别：DEBUG, INFO, WARN, ERROR
- 自动创建日志目录（~/.gx/logs/）
- 线程安全的日志写入
- 时间戳格式化
- 优雅降级（如果无法创建日志文件，不会影响程序运行）

## 使用方法

### 初始化日志记录器

```go
import "github.com/kawaiirei0/gx/internal/logger"

func main() {
    // 初始化日志记录器
    if err := logger.Init(); err != nil {
        // 日志初始化失败不会影响程序运行
        fmt.Println("Warning: failed to initialize logger:", err)
    }
    defer logger.Close()
    
    // 使用日志记录器
    logger.Info("Application started")
}
```

### 记录不同级别的日志

```go
// 调试信息
logger.Debug("Debug message: %s", debugInfo)

// 一般信息
logger.Info("User installed Go version %s", version)

// 警告信息
logger.Warn("Configuration file not found, using defaults")

// 错误信息
logger.Error("Failed to download file: %v", err)
```

### 设置日志级别

```go
// 只记录 INFO 及以上级别的日志
logger.SetLevel(logger.LevelInfo)

// 记录所有级别的日志（包括 DEBUG）
logger.SetLevel(logger.LevelDebug)

// 只记录错误
logger.SetLevel(logger.LevelError)
```

### 创建独立的日志记录器实例

```go
// 创建新的日志记录器实例
customLogger, err := logger.NewLogger()
if err != nil {
    // 处理错误
}
defer customLogger.Close()

customLogger.Info("Custom logger message")
```

## 日志格式

日志文件格式：
```
[2024-01-15 14:30:45] [INFO] Application started
[2024-01-15 14:30:46] [INFO] Detecting installed Go versions
[2024-01-15 14:30:47] [WARN] No active version found
[2024-01-15 14:30:50] [ERROR] Failed to download: network timeout
```

## 日志文件位置

- Windows: `C:\Users\<username>\.gx\logs\gx.log`
- Linux/macOS: `~/.gx/logs/gx.log`

## 注意事项

1. 日志记录器在无法创建日志文件时会自动禁用，不会影响程序正常运行
2. 日志文件以追加模式打开，不会覆盖已有日志
3. 日志记录是线程安全的，可以在并发环境中使用
4. 建议在程序退出时调用 `logger.Close()` 关闭日志文件
