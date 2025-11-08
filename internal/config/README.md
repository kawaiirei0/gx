# Config Store

配置存储模块，负责管理 gx 工具的配置文件。

## 功能

- ✅ 配置文件的加载和保存
- ✅ 配置目录的自动初始化 (~/.gx/)
- ✅ 配置数据结构的 JSON 序列化
- ✅ 默认配置的自动生成

## 使用示例

```go
package main

import (
    "github.com/yourusername/gx/internal/config"
)

func main() {
    // 创建配置存储
    store, err := config.NewStore()
    if err != nil {
        panic(err)
    }

    // 加载配置
    cfg, err := store.Load()
    if err != nil {
        panic(err)
    }

    // 修改配置
    cfg.ActiveVersion = "1.21.5"
    cfg.Versions["1.21.5"] = "/path/to/go1.21.5"

    // 保存配置
    if err := store.Save(cfg); err != nil {
        panic(err)
    }
}
```

## 配置文件结构

配置文件位于 `~/.gx/config.json`，格式如下：

```json
{
  "active_version": "1.21.5",
  "install_path": "/home/user/.gx/versions",
  "versions": {
    "1.21.5": "/home/user/.gx/versions/go1.21.5",
    "1.22.0": "/home/user/.gx/versions/go1.22.0"
  },
  "last_update_check": "2025-11-08T17:31:49+08:00"
}
```

## 接口

### Store

```go
type Store interface {
    // Load 加载配置
    Load() (*Config, error)

    // Save 保存配置
    Save(config *Config) error

    // EnsureConfigDir 确保配置目录存在
    EnsureConfigDir() error
}
```

### Config

```go
type Config struct {
    ActiveVersion   string            // 当前激活版本
    InstallPath     string            // 安装根目录
    Versions        map[string]string // 版本到路径的映射
    LastUpdateCheck time.Time         // 上次检查更新时间
}
```

## 测试

运行测试：

```bash
go test -v ./internal/config/
```

运行示例：

```bash
go run examples/config_demo.go
```
