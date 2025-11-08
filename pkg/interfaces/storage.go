package interfaces

// Storage 管理版本信息的持久化存储
type Storage interface {
	// SaveVersion 保存版本信息
	SaveVersion(version *GoVersion) error

	// GetVersion 获取指定版本信息
	GetVersion(version string) (*GoVersion, error)

	// GetAllVersions 获取所有已安装版本
	GetAllVersions() ([]GoVersion, error)

	// DeleteVersion 删除版本信息
	DeleteVersion(version string) error

	// SetActiveVersion 设置当前激活版本
	SetActiveVersion(version string) error

	// GetActiveVersion 获取当前激活版本
	GetActiveVersion() (string, error)
}
