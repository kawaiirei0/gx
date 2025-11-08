package interfaces

// CrossBuilder 处理跨平台编译逻辑
type CrossBuilder interface {
	// Build 执行跨平台构建
	Build(config BuildConfig) error

	// ValidatePlatform 验证平台组合是否支持
	ValidatePlatform(targetOS, targetArch string) error

	// GetSupportedPlatforms 获取支持的平台组合列表
	GetSupportedPlatforms() []PlatformInfo
}

// BuildConfig 构建配置
type BuildConfig struct {
	SourcePath string   // 源代码路径（目录或文件）
	OutputPath string   // 输出文件路径
	TargetOS   string   // 目标操作系统
	TargetArch string   // 目标架构
	BuildFlags []string // 额外的构建标志
	LDFlags    string   // 链接器标志
}

// PlatformInfo 平台信息
type PlatformInfo struct {
	OS   string // 操作系统
	Arch string // 架构
}
