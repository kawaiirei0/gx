package interfaces

import "time"

// ConfigStore 配置存储接口
type ConfigStore interface {
	// Load 加载配置
	Load() (*Config, error)

	// Save 保存配置
	Save(config *Config) error

	// EnsureConfigDir 确保配置目录存在
	EnsureConfigDir() error
}

// Config 应用配置
type Config struct {
	ActiveVersion   string            `json:"active_version"`    // 当前激活版本
	InstallPath     string            `json:"install_path"`      // 安装根目录
	Versions        map[string]string `json:"versions"`          // 版本到路径的映射
	LastUpdateCheck time.Time         `json:"last_update_check"` // 上次检查更新时间
}
