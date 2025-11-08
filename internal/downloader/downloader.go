package downloader

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/yourusername/gx/internal/logger"
	"github.com/yourusername/gx/pkg/constants"
	"github.com/yourusername/gx/pkg/errors"
	"github.com/yourusername/gx/pkg/interfaces"
)

// httpDownloader 基于 HTTP 的下载器实现
type httpDownloader struct {
	client  *http.Client
	baseURL string
	apiURL  string
}

// NewDownloader 创建新的下载器
func NewDownloader() interfaces.Downloader {
	return &httpDownloader{
		client: &http.Client{
			Timeout: 30 * time.Minute, // 下载超时时间
		},
		baseURL: constants.GoDownloadURL,
		apiURL:  constants.GoVersionsAPIURL,
	}
}

// GetDownloadURL 获取指定版本和平台的下载 URL
func (d *httpDownloader) GetDownloadURL(version string, os string, arch string) (string, error) {
	// 规范化版本号（确保有 "go" 前缀）
	if !strings.HasPrefix(version, "go") {
		version = "go" + version
	}

	logger.Debug("Getting download URL for %s (%s/%s)", version, os, arch)
	
	// 查询远程版本信息
	versions, err := d.fetchVersions()
	if err != nil {
		logger.Error("Failed to fetch versions: %v", err)
		return "", errors.ErrNetworkError.WithCause(err).WithMessage("failed to fetch version list")
	}

	// 查找匹配的版本和文件
	for _, v := range versions {
		if v.Version != version {
			continue
		}

		// 查找匹配平台的文件
		for _, file := range v.Files {
			if file.OS == os && file.Arch == arch {
				return d.baseURL + file.Filename, nil
			}
		}
	}

	return "", errors.ErrVersionNotFound.WithMessage(fmt.Sprintf("version %s not found for %s/%s", version, os, arch))
}

// fetchVersions 从 Go 官方 API 获取版本列表
func (d *httpDownloader) fetchVersions() ([]interfaces.RemoteVersion, error) {
	resp, err := d.client.Get(d.apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var versions []interfaces.RemoteVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, err
	}

	return versions, nil
}

// Download 下载指定版本的 Go 安装包
func (d *httpDownloader) Download(version string, destPath string, progress interfaces.ProgressCallback) error {
	logger.Info("Starting download of Go version %s", version)
	
	// 创建恢复管理器
	recovery := errors.NewRecoveryManager()
	defer func() {
		if err := recovery.Cleanup(); err != nil {
			logger.Warn("Download cleanup failed: %v", err)
		}
	}()
	
	// 获取下载 URL
	url, err := d.GetDownloadURL(version, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		logger.Error("Failed to get download URL: %v", err)
		return err
	}
	logger.Info("Download URL: %s", url)

	// 获取文件信息（包括 SHA256）
	fileInfo, err := d.getFileInfo(version, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		logger.Error("Failed to get file info: %v", err)
		// 优雅降级：如果无法获取文件信息，继续下载但跳过校验
		logger.Warn("Continuing download without checksum verification")
		fileInfo = nil
	} else {
		logger.Info("Expected file size: %d bytes, SHA256: %s", fileInfo.Size, fileInfo.SHA256)
	}

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "gx-download-*")
	if err != nil {
		return errors.ErrDownloadFailed.WithCause(err).WithMessage("failed to create temp file")
	}
	tmpPath := tmpFile.Name()
	
	// 注册临时文件清理
	errors.EnsureFileCleanup(recovery, tmpPath)

	// 下载文件
	logger.Info("Downloading to temporary file: %s", tmpPath)
	expectedSize := int64(0)
	if fileInfo != nil {
		expectedSize = fileInfo.Size
	}
	
	if err := d.downloadFile(url, tmpFile, expectedSize, progress); err != nil {
		tmpFile.Close()
		logger.Error("Download failed: %v", err)
		return err
	}
	tmpFile.Close()
	logger.Info("Download completed")

	// 验证 SHA256（如果有文件信息）
	if fileInfo != nil && fileInfo.SHA256 != "" {
		logger.Info("Verifying checksum...")
		if err := d.verifyChecksum(tmpPath, fileInfo.SHA256); err != nil {
			logger.Error("Checksum verification failed: %v", err)
			return err
		}
		logger.Info("Checksum verified successfully")
	} else {
		logger.Warn("Skipping checksum verification (file info not available)")
	}

	// 确保目标目录存在
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		logger.Error("Failed to create destination directory: %v", err)
		return errors.ErrDownloadFailed.WithCause(err).WithMessage("failed to create destination directory")
	}

	// 移动文件到目标位置
	logger.Info("Moving file to destination: %s", destPath)
	if err := os.Rename(tmpPath, destPath); err != nil {
		// 如果 Rename 失败（可能跨文件系统），尝试复制
		logger.Warn("Rename failed, trying copy: %v", err)
		if err := d.copyFile(tmpPath, destPath); err != nil {
			logger.Error("Failed to copy file: %v", err)
			return errors.ErrDownloadFailed.WithCause(err).WithMessage("failed to move file to destination")
		}
		// 复制成功后删除临时文件
		os.Remove(tmpPath)
	}

	logger.Info("Download completed successfully: %s", destPath)
	return nil
}

// getFileInfo 获取指定版本和平台的文件信息
func (d *httpDownloader) getFileInfo(version string, os string, arch string) (*interfaces.File, error) {
	// 规范化版本号
	if !strings.HasPrefix(version, "go") {
		version = "go" + version
	}

	versions, err := d.fetchVersions()
	if err != nil {
		return nil, errors.ErrNetworkError.WithCause(err).WithMessage("failed to fetch version list")
	}

	for _, v := range versions {
		if v.Version != version {
			continue
		}

		for _, file := range v.Files {
			if file.OS == os && file.Arch == arch {
				return &file, nil
			}
		}
	}

	return nil, errors.ErrVersionNotFound.WithMessage(fmt.Sprintf("file info not found for %s/%s", os, arch))
}

// downloadFile 下载文件并显示进度
func (d *httpDownloader) downloadFile(url string, dest *os.File, expectedSize int64, progress interfaces.ProgressCallback) error {
	resp, err := d.client.Get(url)
	if err != nil {
		return errors.ErrDownloadFailed.WithCause(err).WithMessage("failed to start download")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.ErrDownloadFailed.WithMessage(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}

	// 使用响应头中的 Content-Length，如果没有则使用预期大小
	totalSize := resp.ContentLength
	if totalSize <= 0 {
		totalSize = expectedSize
	}

	// 创建进度读取器
	reader := &progressReader{
		reader:   resp.Body,
		total:    totalSize,
		callback: progress,
	}

	// 复制数据
	_, err = io.Copy(dest, reader)
	if err != nil {
		return errors.ErrDownloadFailed.WithCause(err).WithMessage("failed to write file")
	}

	return nil
}

// verifyChecksum 验证文件的 SHA256 校验和
func (d *httpDownloader) verifyChecksum(filePath string, expectedChecksum string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.ErrChecksumMismatch.WithCause(err).WithMessage("failed to open file for verification")
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return errors.ErrChecksumMismatch.WithCause(err).WithMessage("failed to calculate checksum")
	}

	actualChecksum := hex.EncodeToString(hash.Sum(nil))
	if actualChecksum != expectedChecksum {
		return errors.ErrChecksumMismatch.WithMessage(fmt.Sprintf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum))
	}

	return nil
}

// copyFile 复制文件（用于跨文件系统移动）
func (d *httpDownloader) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	return destFile.Sync()
}

// progressReader 包装 io.Reader 以提供进度回调
type progressReader struct {
	reader   io.Reader
	total    int64
	current  int64
	callback interfaces.ProgressCallback
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.current += int64(n)

	if pr.callback != nil {
		pr.callback(pr.current, pr.total)
	}

	return n, err
}
