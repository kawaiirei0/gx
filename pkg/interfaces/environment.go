package interfaces

// EnvironmentManager 管理系统环境变量
type EnvironmentManager interface {
	// SetGoRoot 设置 GOROOT 环境变量
	SetGoRoot(path string) error

	// SetGoPath 设置 GOPATH 环境变量
	SetGoPath(path string) error

	// UpdatePath 更新 PATH 环境变量
	UpdatePath(goRoot string) error

	// GetGoRoot 获取当前 GOROOT
	GetGoRoot() (string, error)

	// GetGoPath 获取当前 GOPATH
	GetGoPath() (string, error)

	// Backup 备份当前环境变量配置
	Backup() error

	// Restore 恢复环境变量配置
	Restore() error
}
