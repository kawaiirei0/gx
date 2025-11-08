package interfaces

// Installer 负责安装和卸载 Go 版本
type Installer interface {
	// Install 安装指定版本到目标路径
	Install(archivePath string, version string, destPath string) error

	// Uninstall 卸载指定版本
	Uninstall(version string, installPath string) error

	// Verify 验证安装是否成功
	Verify(installPath string, version string) error
}
