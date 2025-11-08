package interfaces

// Downloader 负责下载 Go 安装包
type Downloader interface {
	// Download 下载指定版本的 Go 安装包
	Download(version string, destPath string, progress ProgressCallback) error

	// GetDownloadURL 获取下载 URL
	GetDownloadURL(version string, os string, arch string) (string, error)
}

// RemoteVersion 表示远程可用的 Go 版本信息
type RemoteVersion struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	Files   []File `json:"files"`
}

// File 表示一个可下载的文件信息
type File struct {
	Filename string `json:"filename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	SHA256   string `json:"sha256"`
	Size     int64  `json:"size"`
}
