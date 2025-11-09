package installer

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/kawaiirei0/gx/pkg/constants"
	"github.com/kawaiirei0/gx/pkg/errors"
	"github.com/kawaiirei0/gx/pkg/interfaces"
)

// goInstaller 实现 Installer 接口
type goInstaller struct {
	platform interfaces.PlatformAdapter
}

// NewInstaller 创建新的安装器
func NewInstaller(platform interfaces.PlatformAdapter) interfaces.Installer {
	return &goInstaller{
		platform: platform,
	}
}

// Install 安装指定版本到目标路径
func (i *goInstaller) Install(archivePath string, version string, destPath string) error {
	// 创建恢复管理器
	recovery := errors.NewRecoveryManager()
	
	// 确保目标目录存在
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return errors.ErrInstallFailed.
			WithCause(err).
			WithMessage("failed to create destination directory").
			WithContext("dest_path", destPath)
	}

	// 注册清理函数：如果安装失败，删除目标目录
	errors.EnsureDirectoryCleanup(recovery, destPath)

	// 根据文件扩展名选择解压方法
	var extractErr error
	if strings.HasSuffix(archivePath, constants.ArchiveExtZip) {
		extractErr = i.extractZip(archivePath, destPath)
		if extractErr != nil {
			extractErr = errors.ErrInstallFailed.
				WithCause(extractErr).
				WithMessage("failed to extract zip archive").
				WithContext("archive_path", archivePath).
				WithContext("dest_path", destPath)
		}
	} else if strings.HasSuffix(archivePath, constants.ArchiveExtTarGz) {
		extractErr = i.extractTarGz(archivePath, destPath)
		if extractErr != nil {
			extractErr = errors.ErrInstallFailed.
				WithCause(extractErr).
				WithMessage("failed to extract tar.gz archive").
				WithContext("archive_path", archivePath).
				WithContext("dest_path", destPath)
		}
	} else {
		extractErr = errors.ErrInstallFailed.
			WithMessage("unsupported archive format").
			WithContext("archive_path", archivePath)
	}

	if extractErr != nil {
		// 执行清理
		recovery.Cleanup()
		return extractErr
	}

	// 验证安装
	if err := i.Verify(destPath, version); err != nil {
		// 验证失败，清理安装目录
		recovery.Cleanup()
		return errors.Wrap(err, "INSTALL_FAILED", "installation verification failed").
			WithContext("dest_path", destPath).
			WithContext("version", version)
	}

	// 安装成功，清除清理函数（不需要清理）
	recovery.Clear()
	return nil
}

// extractZip 解压 ZIP 文件
func (i *goInstaller) extractZip(archivePath string, destPath string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		if err := i.extractZipFile(file, destPath); err != nil {
			return err
		}
	}

	return nil
}

// extractZipFile 解压单个 ZIP 文件条目
func (i *goInstaller) extractZipFile(file *zip.File, destPath string) error {
	// 构建目标路径
	// Go 压缩包的结构是 go/bin/go, go/src/... 等
	// 我们需要去掉顶层的 "go" 目录
	targetPath := filepath.Join(destPath, i.stripTopDir(file.Name))

	// 如果是目录
	if file.FileInfo().IsDir() {
		return os.MkdirAll(targetPath, file.Mode())
	}

	// 确保父目录存在
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	// 创建文件
	destFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer destFile.Close()

	// 打开源文件
	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 复制内容
	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}

	return nil
}

// extractTarGz 解压 tar.gz 文件
func (i *goInstaller) extractTarGz(archivePath string, destPath string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := i.extractTarFile(header, tarReader, destPath); err != nil {
			return err
		}
	}

	return nil
}

// extractTarFile 解压单个 tar 文件条目
func (i *goInstaller) extractTarFile(header *tar.Header, reader io.Reader, destPath string) error {
	// 构建目标路径，去掉顶层 "go" 目录
	targetPath := filepath.Join(destPath, i.stripTopDir(header.Name))

	switch header.Typeflag {
	case tar.TypeDir:
		return os.MkdirAll(targetPath, os.FileMode(header.Mode))

	case tar.TypeReg:
		// 确保父目录存在
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		// 创建文件
		file, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
		if err != nil {
			return err
		}
		defer file.Close()

		// 复制内容
		if _, err := io.Copy(file, reader); err != nil {
			return err
		}

		return nil

	case tar.TypeSymlink:
		// 处理符号链接
		return os.Symlink(header.Linkname, targetPath)

	default:
		// 忽略其他类型
		return nil
	}
}

// stripTopDir 去掉路径中的顶层目录
// 例如: "go/bin/go" -> "bin/go"
func (i *goInstaller) stripTopDir(path string) string {
	parts := strings.Split(filepath.ToSlash(path), "/")
	if len(parts) > 1 && parts[0] == "go" {
		return filepath.Join(parts[1:]...)
	}
	return path
}

// Verify 验证安装是否成功
func (i *goInstaller) Verify(installPath string, version string) error {
	// 检查 bin 目录是否存在
	binDir := filepath.Join(installPath, "bin")
	if _, err := os.Stat(binDir); os.IsNotExist(err) {
		return errors.ErrInstallFailed.WithMessage("bin directory not found")
	}

	// 检查 go 可执行文件是否存在
	goExe := "go"
	if runtime.GOOS == constants.OSWindows {
		goExe = "go.exe"
	}

	goPath := filepath.Join(binDir, goExe)
	if _, err := os.Stat(goPath); os.IsNotExist(err) {
		return errors.ErrInstallFailed.WithMessage("go executable not found")
	}

	// 确保 go 可执行文件有执行权限（Unix 系统）
	if runtime.GOOS != constants.OSWindows {
		if err := i.platform.MakeExecutable(goPath); err != nil {
			return errors.ErrInstallFailed.WithCause(err).WithMessage("failed to set executable permission")
		}
	}

	// 验证版本号
	cmd := exec.Command(goPath, "version")
	output, err := cmd.Output()
	if err != nil {
		return errors.ErrInstallFailed.WithCause(err).WithMessage("failed to execute go version")
	}

	// 解析版本号
	versionStr := string(output)
	versionRegex := regexp.MustCompile(`go version go(\d+\.\d+(?:\.\d+)?)`)
	matches := versionRegex.FindStringSubmatch(versionStr)
	if len(matches) < 2 {
		return errors.ErrInstallFailed.WithMessage("failed to parse go version output")
	}

	installedVersion := matches[1]
	expectedVersion := strings.TrimPrefix(version, "go")

	if installedVersion != expectedVersion {
		return errors.ErrInstallFailed.WithMessage(fmt.Sprintf("version mismatch: expected %s, got %s", expectedVersion, installedVersion))
	}

	return nil
}

// Uninstall 卸载指定版本
func (i *goInstaller) Uninstall(version string, installPath string) error {
	// 检查路径是否存在
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		return errors.ErrVersionNotInstalled.WithMessage("installation path does not exist")
	}

	// 删除整个安装目录
	if err := os.RemoveAll(installPath); err != nil {
		return errors.ErrUninstallFailed.WithCause(err).WithMessage("failed to remove installation directory")
	}

	return nil
}
