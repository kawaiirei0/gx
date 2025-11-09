package platform_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/kawaiirei0/gx/internal/platform"
)

// TestPlatformAdapter 测试平台适配器的基本功能
func TestPlatformAdapter(t *testing.T) {
	adapter := platform.NewAdapter()

	t.Run("GetOS", func(t *testing.T) {
		os := adapter.GetOS()
		if os == "" {
			t.Error("GetOS returned empty string")
		}
		t.Logf("Operating System: %s", os)
	})

	t.Run("GetArch", func(t *testing.T) {
		arch := adapter.GetArch()
		if arch == "" {
			t.Error("GetArch returned empty string")
		}
		t.Logf("Architecture: %s", arch)
	})

	t.Run("PathSeparator", func(t *testing.T) {
		sep := adapter.PathSeparator()
		if sep == "" {
			t.Error("PathSeparator returned empty string")
		}
		t.Logf("Path Separator: %q", sep)
	})

	t.Run("GetHomeDir", func(t *testing.T) {
		home, err := adapter.GetHomeDir()
		if err != nil {
			t.Errorf("GetHomeDir failed: %v", err)
		}
		if home == "" {
			t.Error("GetHomeDir returned empty string")
		}
		t.Logf("Home Directory: %s", home)
	})

	t.Run("JoinPath", func(t *testing.T) {
		path := adapter.JoinPath("home", "user", "test")
		expected := filepath.Join("home", "user", "test")
		if path != expected {
			t.Errorf("JoinPath failed: got %s, want %s", path, expected)
		}
		t.Logf("Joined Path: %s", path)
	})

	t.Run("NormalizePath", func(t *testing.T) {
		testPath := filepath.Join(".", "test", "..", "file.txt")
		normalized := adapter.NormalizePath(testPath)
		if normalized == "" {
			t.Error("NormalizePath returned empty string")
		}
		t.Logf("Original: %s, Normalized: %s", testPath, normalized)
	})
}

// TestFilePermissions 测试文件权限相关功能
func TestFilePermissions(t *testing.T) {
	adapter := platform.NewAdapter()

	// 创建临时测试文件
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_executable")

	// 创建测试文件
	f, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	f.Close()

	t.Run("IsExecutable_BeforeMake", func(t *testing.T) {
		isExec := adapter.IsExecutable(testFile)
		t.Logf("File is executable (before): %v", isExec)
	})

	t.Run("MakeExecutable", func(t *testing.T) {
		err := adapter.MakeExecutable(testFile)
		if err != nil {
			t.Errorf("MakeExecutable failed: %v", err)
		}
	})

	t.Run("IsExecutable_AfterMake", func(t *testing.T) {
		isExec := adapter.IsExecutable(testFile)
		t.Logf("File is executable (after): %v", isExec)
	})
}

// TestPlatformUtils 测试平台工具函数
func TestPlatformUtils(t *testing.T) {
	t.Run("GetExecutableExtension", func(t *testing.T) {
		ext := platform.GetExecutableExtension()
		t.Logf("Executable Extension: %q", ext)
	})

	t.Run("GetArchiveExtension", func(t *testing.T) {
		ext := platform.GetArchiveExtension()
		if ext == "" {
			t.Error("GetArchiveExtension returned empty string")
		}
		t.Logf("Archive Extension: %s", ext)
	})

	t.Run("GetPlatformString", func(t *testing.T) {
		platform := platform.GetPlatformString()
		if platform == "" {
			t.Error("GetPlatformString returned empty string")
		}
		t.Logf("Platform String: %s", platform)
	})

	t.Run("IsSupportedPlatform", func(t *testing.T) {
		testCases := []struct {
			os       string
			arch     string
			expected bool
		}{
			{"windows", "amd64", true},
			{"linux", "amd64", true},
			{"darwin", "arm64", true},
			{"freebsd", "amd64", false},
			{"linux", "mips", false},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("%s/%s", tc.os, tc.arch), func(t *testing.T) {
				result := platform.IsSupportedPlatform(tc.os, tc.arch)
				if result != tc.expected {
					t.Errorf("IsSupportedPlatform(%s, %s) = %v, want %v",
						tc.os, tc.arch, result, tc.expected)
				}
			})
		}
	})

	t.Run("FileExists", func(t *testing.T) {
		// Test with existing file
		tmpFile := filepath.Join(t.TempDir(), "test.txt")
		os.WriteFile(tmpFile, []byte("test"), 0644)

		if !platform.FileExists(tmpFile) {
			t.Error("FileExists returned false for existing file")
		}

		// Test with non-existing file
		if platform.FileExists(filepath.Join(t.TempDir(), "nonexistent.txt")) {
			t.Error("FileExists returned true for non-existing file")
		}
	})

	t.Run("IsDirectory", func(t *testing.T) {
		tmpDir := t.TempDir()

		if !platform.IsDirectory(tmpDir) {
			t.Error("IsDirectory returned false for directory")
		}

		tmpFile := filepath.Join(tmpDir, "file.txt")
		os.WriteFile(tmpFile, []byte("test"), 0644)

		if platform.IsDirectory(tmpFile) {
			t.Error("IsDirectory returned true for file")
		}
	})

	t.Run("EnsureDir", func(t *testing.T) {
		tmpDir := t.TempDir()
		testDir := filepath.Join(tmpDir, "test", "nested", "dir")

		err := platform.EnsureDir(testDir)
		if err != nil {
			t.Errorf("EnsureDir failed: %v", err)
		}

		if !platform.IsDirectory(testDir) {
			t.Error("EnsureDir did not create directory")
		}
	})
}
