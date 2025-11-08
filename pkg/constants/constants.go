package constants

const (
	// AppName 应用名称
	AppName = "gx"

	// AppVersion 应用版本
	AppVersion = "0.1.0"

	// DefaultInstallDir 默认安装目录
	DefaultInstallDir = ".gx/versions"

	// ConfigFileName 配置文件名
	ConfigFileName = "config.json"

	// ConfigDir 配置目录
	ConfigDir = ".gx"

	// GoDownloadURL Go 官方下载地址
	GoDownloadURL = "https://go.dev/dl/"

	// GoVersionsAPIURL Go 版本列表 API
	GoVersionsAPIURL = "https://go.dev/dl/?mode=json"

	// MinGoVersion 最低支持的 Go 版本
	MinGoVersion = "1.16"
)

// 平台相关常量
const (
	// OSWindows Windows 操作系统
	OSWindows = "windows"

	// OSLinux Linux 操作系统
	OSLinux = "linux"

	// OSDarwin macOS 操作系统
	OSDarwin = "darwin"

	// ArchAMD64 AMD64 架构
	ArchAMD64 = "amd64"

	// ArchARM64 ARM64 架构
	ArchARM64 = "arm64"

	// Arch386 386 架构
	Arch386 = "386"
)

// 文件扩展名
const (
	// ArchiveExtTarGz tar.gz 压缩包
	ArchiveExtTarGz = ".tar.gz"

	// ArchiveExtZip zip 压缩包
	ArchiveExtZip = ".zip"
)

// 环境变量名称
const (
	// EnvGoRoot GOROOT 环境变量
	EnvGoRoot = "GOROOT"

	// EnvGoPath GOPATH 环境变量
	EnvGoPath = "GOPATH"

	// EnvPath PATH 环境变量
	EnvPath = "PATH"
)
