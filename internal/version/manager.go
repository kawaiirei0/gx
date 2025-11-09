package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/kawaiirei0/gx/internal/logger"
	"github.com/kawaiirei0/gx/pkg/constants"
	"github.com/kawaiirei0/gx/pkg/errors"
	"github.com/kawaiirei0/gx/pkg/interfaces"
)

// manager 实现 VersionManager 接口
type manager struct {
	configStore interfaces.ConfigStore
	platform    interfaces.PlatformAdapter
	envManager  interfaces.EnvironmentManager
	downloader  interfaces.Downloader
	installer   interfaces.Installer
}

// NewManager 创建新的版本管理器
func NewManager(configStore interfaces.ConfigStore, platform interfaces.PlatformAdapter, envManager interfaces.EnvironmentManager, downloader interfaces.Downloader, installer interfaces.Installer) interfaces.VersionManager {
	return &manager{
		configStore: configStore,
		platform:    platform,
		envManager:  envManager,
		downloader:  downloader,
		installer:   installer,
	}
}

// DetectInstalled 检测系统中已安装的 Go 版本
func (m *manager) DetectInstalled() ([]interfaces.GoVersion, error) {
	logger.Info("Detecting installed Go versions")
	
	cfg, err := m.configStore.Load()
	if err != nil {
		logger.Error("Failed to load config: %v", err)
		return nil, errors.ErrStorageFailed.WithCause(err).WithMessage("failed to load config")
	}

	var versions []interfaces.GoVersion

	// 扫描 gx 管理的版本目录
	versionsDir := cfg.InstallPath
	if _, err := os.Stat(versionsDir); err == nil {
		gxVersions, err := m.scanGxVersions(versionsDir, cfg.ActiveVersion)
		if err == nil {
			versions = append(versions, gxVersions...)
		}
	}

	// 扫描系统环境变量中的 Go 版本
	systemVersion, err := m.detectSystemGoVersion()
	if err == nil && systemVersion != nil {
		// 检查是否已经在列表中（避免重复）
		found := false
		for _, v := range versions {
			if v.Path == systemVersion.Path {
				found = true
				break
			}
		}
		if !found {
			versions = append(versions, *systemVersion)
		}
	}

	logger.Info("Detected %d installed Go versions", len(versions))
	return versions, nil
}

// GetActive 获取当前激活的版本
func (m *manager) GetActive() (*interfaces.GoVersion, error) {
	cfg, err := m.configStore.Load()
	if err != nil {
		logger.Error("Failed to load config: %v", err)
		return nil, errors.ErrStorageFailed.WithCause(err).WithMessage("failed to load config")
	}

	// 如果配置中有激活版本，返回该版本
	if cfg.ActiveVersion != "" {
		versionPath, ok := cfg.Versions[cfg.ActiveVersion]
		if !ok {
			logger.Warn("Active version %s not found in config", cfg.ActiveVersion)
			return nil, errors.ErrVersionNotFound.WithMessage("active version not found in config")
		}

		logger.Info("Active version: %s", cfg.ActiveVersion)
		return &interfaces.GoVersion{
			Version:  cfg.ActiveVersion,
			Path:     versionPath,
			IsActive: true,
		}, nil
	}

	// 否则尝试检测系统中的 Go 版本
	logger.Info("No active version in config, detecting system Go version")
	systemVersion, err := m.detectSystemGoVersion()
	if err != nil {
		logger.Warn("No active version found")
		return nil, errors.ErrVersionNotFound.WithMessage("no active version found")
	}

	logger.Info("System Go version detected: %s", systemVersion.Version)
	return systemVersion, nil
}

// scanGxVersions 扫描 gx 管理的版本目录
func (m *manager) scanGxVersions(versionsDir string, activeVersion string) ([]interfaces.GoVersion, error) {
	entries, err := os.ReadDir(versionsDir)
	if err != nil {
		return nil, err
	}

	var versions []interfaces.GoVersion
	versionRegex := regexp.MustCompile(`^go(\d+\.\d+(?:\.\d+)?)$`)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// 检查目录名是否符合 Go 版本格式
		dirName := entry.Name()
		matches := versionRegex.FindStringSubmatch(dirName)
		if len(matches) < 2 {
			continue
		}

		// 使用完整的目录名作为版本号（包含 "go" 前缀）
		// 例如：目录名 "go1.25.4" -> 版本号 "go1.25.4"
		fullVersion := dirName
		versionPath := filepath.Join(versionsDir, dirName)

		// 验证该目录是否包含有效的 Go 安装
		if !m.isValidGoInstallation(versionPath) {
			continue
		}

		// 获取安装日期
		info, err := entry.Info()
		var installDate time.Time
		if err == nil {
			installDate = info.ModTime()
		}

		versions = append(versions, interfaces.GoVersion{
			Version:     fullVersion,
			Path:        versionPath,
			IsActive:    fullVersion == activeVersion,
			InstallDate: installDate,
		})
	}

	return versions, nil
}

// detectSystemGoVersion 检测系统环境变量中的 Go 版本
func (m *manager) detectSystemGoVersion() (*interfaces.GoVersion, error) {
	// 尝试执行 go version 命令
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// 解析版本号，格式如: "go version go1.21.5 windows/amd64"
	versionStr := string(output)
	versionRegex := regexp.MustCompile(`go version (go\d+\.\d+(?:\.\d+)?)`)
	matches := versionRegex.FindStringSubmatch(versionStr)
	if len(matches) < 2 {
		return nil, errors.ErrInvalidVersion.WithMessage("failed to parse go version output")
	}

	// 使用完整版本号（包含 "go" 前缀）
	fullVersion := matches[1]

	// 获取 GOROOT 路径
	goroot := os.Getenv(constants.EnvGoRoot)
	if goroot == "" {
		// 如果 GOROOT 未设置，尝试通过 go env 获取
		cmd := exec.Command("go", "env", "GOROOT")
		output, err := cmd.Output()
		if err == nil {
			goroot = strings.TrimSpace(string(output))
		}
	}

	return &interfaces.GoVersion{
		Version:  fullVersion,
		Path:     goroot,
		IsActive: true,
	}, nil
}

// isValidGoInstallation 验证目录是否包含有效的 Go 安装
func (m *manager) isValidGoInstallation(path string) bool {
	// 检查 bin 目录是否存在
	binDir := filepath.Join(path, "bin")
	if _, err := os.Stat(binDir); os.IsNotExist(err) {
		return false
	}

	// 检查 go 可执行文件是否存在
	goExe := "go"
	if m.platform.GetOS() == constants.OSWindows {
		goExe = "go.exe"
	}

	goPath := filepath.Join(binDir, goExe)
	if _, err := os.Stat(goPath); os.IsNotExist(err) {
		return false
	}

	return true
}

// Install 安装指定版本
func (m *manager) Install(version string, progress interfaces.ProgressCallback) error {
	// 规范化版本号
	normalizedVersion := version
	if !strings.HasPrefix(version, "go") {
		normalizedVersion = "go" + version
	}

	logger.Info("Starting installation of Go version %s", normalizedVersion)

	// 创建恢复管理器
	recovery := errors.NewRecoveryManager()
	defer func() {
		// 确保在函数退出时执行清理
		if err := recovery.Cleanup(); err != nil {
			logger.Warn("Cleanup failed: %v", err)
		}
	}()

	// 加载配置
	cfg, err := m.configStore.Load()
	if err != nil {
		logger.Error("Failed to load config: %v", err)
		return errors.ErrStorageFailed.WithCause(err).WithMessage("failed to load config")
	}

	// 备份配置文件
	configPath := filepath.Join(cfg.InstallPath, "..", constants.ConfigFileName)
	backupPath, err := errors.BackupFile(configPath)
	if err == nil {
		logger.Debug("Config backed up to: %s", backupPath)
		recovery.AddRollback(func() error {
			return errors.RestoreFile(backupPath, configPath)
		})
		recovery.AddCleanup(func() error {
			return errors.SafeRemoveFile(backupPath)
		})
	}

	// 检查版本是否已安装
	if _, ok := cfg.Versions[normalizedVersion]; ok {
		logger.Warn("Version %s is already installed", normalizedVersion)
		return errors.ErrVersionAlreadyInstalled.WithMessage("version " + normalizedVersion + " is already installed")
	}

	// 确保安装目录存在
	if err := os.MkdirAll(cfg.InstallPath, 0755); err != nil {
		logger.Error("Failed to create install directory: %v", err)
		return errors.ErrInstallFailed.WithCause(err).WithMessage("failed to create install directory")
	}

	// 构建下载文件路径
	archiveExt := constants.ArchiveExtTarGz
	if m.platform.GetOS() == constants.OSWindows {
		archiveExt = constants.ArchiveExtZip
	}
	
	archiveFilename := normalizedVersion + "." + m.platform.GetOS() + "-" + m.platform.GetArch() + archiveExt
	archivePath := filepath.Join(cfg.InstallPath, archiveFilename)

	// 注册下载文件的清理
	errors.EnsureFileCleanup(recovery, archivePath)

	logger.Info("Downloading %s to %s", normalizedVersion, archivePath)
	// 下载安装包
	if err := m.downloader.Download(normalizedVersion, archivePath, progress); err != nil {
		logger.Error("Download failed: %v", err)
		// 执行回滚
		if rollbackErr := recovery.Rollback(); rollbackErr != nil {
			logger.Error("Rollback failed: %v", rollbackErr)
		}
		return err
	}
	logger.Info("Download completed successfully")

	// 构建安装目标路径
	versionPath := filepath.Join(cfg.InstallPath, normalizedVersion)

	// 注册版本目录的清理（如果安装失败）
	errors.EnsureDirectoryCleanup(recovery, versionPath)

	logger.Info("Installing %s to %s", normalizedVersion, versionPath)
	// 安装（解压）
	if err := m.installer.Install(archivePath, normalizedVersion, versionPath); err != nil {
		logger.Error("Installation failed: %v", err)
		// 执行回滚和清理
		if rollbackErr := recovery.CleanupAndRollback(); rollbackErr != nil {
			logger.Error("Recovery failed: %v", rollbackErr)
		}
		return err
	}
	logger.Info("Installation completed successfully")

	// 更新配置
	cfg.Versions[normalizedVersion] = versionPath
	if err := m.configStore.Save(cfg); err != nil {
		logger.Error("Failed to save config after installation: %v", err)
		// 配置保存失败，尝试清理已安装的版本
		if cleanupErr := recovery.CleanupAndRollback(); cleanupErr != nil {
			logger.Error("Failed to cleanup after config save failure: %v", cleanupErr)
		}
		return errors.ErrStorageFailed.WithCause(err).WithMessage("installation succeeded but failed to save config")
	}

	// 安装成功，清除清理函数（不需要清理）
	recovery.Clear()

	logger.Info("Go version %s installed successfully", normalizedVersion)
	return nil
}

// SwitchTo 切换到指定版本
func (m *manager) SwitchTo(version string) error {
	startTime := time.Now()
	logger.Info("Switching to Go version %s", version)

	// 加载配置
	cfg, err := m.configStore.Load()
	if err != nil {
		logger.Error("Failed to load config: %v", err)
		return errors.ErrStorageFailed.WithCause(err).WithMessage("failed to load config")
	}

	// 规范化版本号
	normalizedVersion := version
	if !strings.HasPrefix(version, "go") {
		normalizedVersion = "go" + version
	}

	// 检查版本是否已安装
	versionPath, ok := cfg.Versions[normalizedVersion]
	if !ok {
		logger.Error("Version %s is not installed", normalizedVersion)
		// 提供更友好的错误消息
		versionDisplay := strings.TrimPrefix(normalizedVersion, "go")
		return errors.ErrVersionNotInstalled.
			WithMessage(fmt.Sprintf("Go %s is not installed. Install it first using: gx install %s", versionDisplay, versionDisplay))
	}

	// 验证版本目录是否存在且有效
	if !m.isValidGoInstallation(versionPath) {
		logger.Error("Version %s installation is invalid or corrupted", normalizedVersion)
		versionDisplay := strings.TrimPrefix(normalizedVersion, "go")
		return errors.ErrVersionNotFound.
			WithMessage(fmt.Sprintf("Go %s installation is invalid or corrupted. Try reinstalling: gx uninstall %s && gx install %s", versionDisplay, versionDisplay, versionDisplay))
	}

	// 更新环境变量
	if err := m.envManager.SetGoRoot(versionPath); err != nil {
		logger.Error("Failed to set GOROOT: %v", err)
		return errors.ErrEnvironmentSetupFailed.WithCause(err).WithMessage("failed to set GOROOT")
	}

	if err := m.envManager.UpdatePath(versionPath); err != nil {
		logger.Error("Failed to update PATH: %v", err)
		return errors.ErrEnvironmentSetupFailed.WithCause(err).WithMessage("failed to update PATH")
	}

	// 更新配置中的激活版本
	cfg.ActiveVersion = normalizedVersion
	if err := m.configStore.Save(cfg); err != nil {
		logger.Error("Failed to save config: %v", err)
		return errors.ErrStorageFailed.WithCause(err).WithMessage("failed to save config")
	}

	// 验证切换是否成功
	activeVersion, err := m.GetActive()
	if err != nil {
		logger.Error("Failed to verify version switch: %v", err)
		return errors.ErrEnvironmentSetupFailed.WithCause(err).WithMessage("failed to verify version switch")
	}

	if activeVersion.Version != normalizedVersion {
		logger.Error("Version switch verification failed: expected %s, got %s", normalizedVersion, activeVersion.Version)
		return errors.ErrEnvironmentSetupFailed.WithMessage("version switch verification failed")
	}

	// 检查切换时间是否在 300ms 内
	elapsed := time.Since(startTime)
	if elapsed > 300*time.Millisecond {
		logger.Warn("Version switch took %v (target: 300ms)", elapsed)
	}

	logger.Info("Successfully switched to Go version %s (took %v)", normalizedVersion, elapsed)
	return nil
}

// ListAvailable 获取可用的远程版本列表
func (m *manager) ListAvailable() ([]string, error) {
	logger.Info("Fetching available Go versions from remote")
	// 通过 downloader 获取版本信息
	// 我们需要创建一个辅助方法来获取版本列表
	versions, err := m.fetchRemoteVersions()
	if err != nil {
		logger.Error("Failed to fetch remote versions: %v", err)
		return nil, err
	}

	// 提取版本号列表
	var versionList []string
	for _, v := range versions {
		versionList = append(versionList, v.Version)
	}

	logger.Info("Found %d available Go versions", len(versionList))
	return versionList, nil
}

// GetLatest 获取最新稳定版本
func (m *manager) GetLatest() (string, error) {
	logger.Info("Fetching latest stable Go version")
	versions, err := m.fetchRemoteVersions()
	if err != nil {
		logger.Error("Failed to fetch remote versions: %v", err)
		return "", err
	}

	// 查找最新的稳定版本
	// Go 官方 API 返回的版本列表通常按发布时间排序，最新的在前面
	for _, v := range versions {
		if v.Stable {
			logger.Info("Latest stable version: %s", v.Version)
			return v.Version, nil
		}
	}

	// 如果没有找到稳定版本，返回第一个版本
	if len(versions) > 0 {
		logger.Warn("No stable version found, returning first version: %s", versions[0].Version)
		return versions[0].Version, nil
	}

	logger.Error("No versions available")
	return "", errors.ErrVersionNotFound.WithMessage("no versions available")
}

// fetchRemoteVersions 从远程获取版本列表
// 这是一个辅助方法，通过 downloader 的 GetDownloadURL 来触发版本信息获取
func (m *manager) fetchRemoteVersions() ([]interfaces.RemoteVersion, error) {
	// 我们需要直接访问 Go 官方 API
	// 由于 downloader 已经有 fetchVersions 方法，我们需要创建一个新的接口方法
	// 或者直接在这里实现 HTTP 请求
	
	// 为了避免重复代码，我们直接实现 HTTP 请求
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(constants.GoVersionsAPIURL)
	if err != nil {
		return nil, errors.ErrNetworkError.WithCause(err).WithMessage("failed to fetch version list")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.ErrNetworkError.WithMessage("unexpected status code: " + resp.Status)
	}

	var versions []interfaces.RemoteVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, errors.ErrNetworkError.WithCause(err).WithMessage("failed to parse version list")
	}

	return versions, nil
}

// Uninstall 卸载指定版本
func (m *manager) Uninstall(version string) error {
	logger.Info("Uninstalling Go version %s", version)
	
	// 加载配置
	cfg, err := m.configStore.Load()
	if err != nil {
		logger.Error("Failed to load config: %v", err)
		return errors.ErrStorageFailed.WithCause(err).WithMessage("failed to load config")
	}

	// 检查版本是否已安装
	versionPath, ok := cfg.Versions[version]
	if !ok {
		logger.Warn("Version %s is not installed", version)
		return errors.ErrVersionNotInstalled.WithMessage("version " + version + " is not installed")
	}

	// 安全检查：不能卸载当前激活的版本
	if cfg.ActiveVersion == version {
		logger.Error("Cannot uninstall currently active version %s", version)
		return errors.ErrUninstallFailed.WithMessage("cannot uninstall the currently active version")
	}

	// 删除版本目录
	logger.Info("Removing version directory: %s", versionPath)
	if err := os.RemoveAll(versionPath); err != nil {
		logger.Error("Failed to remove version directory: %v", err)
		return errors.ErrUninstallFailed.WithCause(err).WithMessage("failed to remove version directory")
	}

	// 从配置中移除版本记录
	delete(cfg.Versions, version)
	if err := m.configStore.Save(cfg); err != nil {
		logger.Error("Failed to save config after uninstall: %v", err)
		return errors.ErrStorageFailed.WithCause(err).WithMessage("failed to save config after uninstall")
	}

	logger.Info("Successfully uninstalled Go version %s", version)
	return nil
}
