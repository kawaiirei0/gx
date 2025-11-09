package cmd

import (
	"github.com/kawaiirei0/gx/internal/crossbuilder"
	"github.com/kawaiirei0/gx/internal/downloader"
	"github.com/kawaiirei0/gx/internal/environment"
	"github.com/kawaiirei0/gx/internal/installer"
	"github.com/kawaiirei0/gx/internal/platform"
	"github.com/kawaiirei0/gx/internal/version"
	"github.com/kawaiirei0/gx/internal/wrapper"
	"github.com/kawaiirei0/gx/pkg/interfaces"
	
	configpkg "github.com/kawaiirei0/gx/internal/config"
)

// AppContext 包含应用程序的所有依赖
type AppContext struct {
	VersionManager interfaces.VersionManager
	CLIWrapper     interfaces.CLIWrapper
	CrossBuilder   interfaces.CrossBuilder
	ConfigStore    interfaces.ConfigStore
	Platform       interfaces.PlatformAdapter
	EnvManager     interfaces.EnvironmentManager
}

// NewAppContext 创建新的应用程序上下文
func NewAppContext() (*AppContext, error) {
	// 初始化平台适配器
	platformAdapter := platform.NewAdapter()

	// 初始化配置存储
	configStore, err := configpkg.NewStore()
	if err != nil {
		return nil, err
	}

	// 初始化环境管理器
	envManager := environment.NewManager(platformAdapter)

	// 初始化下载器
	downloaderInstance := downloader.NewDownloader()

	// 初始化安装器
	installerInstance := installer.NewInstaller(platformAdapter)

	// 初始化版本管理器
	versionManager := version.NewManager(
		configStore,
		platformAdapter,
		envManager,
		downloaderInstance,
		installerInstance,
	)

	// 初始化 CLI 包装器
	cliWrapper := wrapper.NewCLIWrapper(versionManager, platformAdapter)

	// 初始化跨平台构建器
	crossBuilderInstance := crossbuilder.NewCrossBuilder(versionManager, platformAdapter)

	return &AppContext{
		VersionManager: versionManager,
		CLIWrapper:     cliWrapper,
		CrossBuilder:   crossBuilderInstance,
		ConfigStore:    configStore,
		Platform:       platformAdapter,
		EnvManager:     envManager,
	}, nil
}
