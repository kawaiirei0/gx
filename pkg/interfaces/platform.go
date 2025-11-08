package interfaces

// PlatformAdapter 提供跨平台抽象层
type PlatformAdapter interface {
	// GetOS 获取操作系统类型
	GetOS() string

	// GetArch 获取架构类型
	GetArch() string

	// PathSeparator 获取路径分隔符
	PathSeparator() string

	// NormalizePath 规范化路径
	NormalizePath(path string) string

	// IsExecutable 检查文件是否可执行
	IsExecutable(path string) bool

	// MakeExecutable 设置文件为可执行
	MakeExecutable(path string) error

	// GetHomeDir 获取用户主目录
	GetHomeDir() (string, error)

	// JoinPath 连接路径
	JoinPath(elem ...string) string
}
