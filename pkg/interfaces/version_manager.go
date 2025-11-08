package interfaces

import "time"

// VersionManager 管理 Go 版本的安装、切换和检测
type VersionManager interface {
	// DetectInstalled 检测系统中已安装的 Go 版本
	DetectInstalled() ([]GoVersion, error)

	// GetActive 获取当前激活的版本
	GetActive() (*GoVersion, error)

	// Install 安装指定版本
	Install(version string, progress ProgressCallback) error

	// SwitchTo 切换到指定版本
	SwitchTo(version string) error

	// ListAvailable 获取可用的远程版本列表
	ListAvailable() ([]string, error)

	// GetLatest 获取最新稳定版本
	GetLatest() (string, error)

	// Uninstall 卸载指定版本
	Uninstall(version string) error
}

// GoVersion 表示一个 Go 版本的信息
type GoVersion struct {
	Version     string    `json:"version"`      // 例如: "1.21.5"
	Path        string    `json:"path"`         // 安装路径
	IsActive    bool      `json:"is_active"`    // 是否为当前激活版本
	InstallDate time.Time `json:"install_date"` // 安装日期
}

// ProgressCallback 下载进度回调函数
type ProgressCallback func(downloaded int64, total int64)
