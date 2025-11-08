package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestConfigStore 测试配置存储的基本功能
func TestConfigStore(t *testing.T) {
	// 创建临时目录用于测试
	tempDir := t.TempDir()
	
	// 创建测试用的 store
	configPath := filepath.Join(tempDir, "config.json")
	store := &fileStore{
		configPath: configPath,
	}

	// 测试 1: 加载不存在的配置文件应返回默认配置
	t.Run("Load non-existent config returns default", func(t *testing.T) {
		config, err := store.Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if config == nil {
			t.Fatal("Load() returned nil config")
		}
		if config.ActiveVersion != "" {
			t.Errorf("Expected empty ActiveVersion, got %s", config.ActiveVersion)
		}
		if config.Versions == nil {
			t.Error("Expected initialized Versions map")
		}
	})

	// 测试 2: 保存配置
	t.Run("Save config", func(t *testing.T) {
		config := &Config{
			ActiveVersion:   "1.21.5",
			InstallPath:     "/home/user/.gx/versions",
			Versions:        map[string]string{"1.21.5": "/home/user/.gx/versions/go1.21.5"},
			LastUpdateCheck: time.Now(),
		}

		err := store.Save(config)
		if err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		// 验证文件是否存在
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Error("Config file was not created")
		}
	})

	// 测试 3: 加载已保存的配置
	t.Run("Load saved config", func(t *testing.T) {
		config, err := store.Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}

		if config.ActiveVersion != "1.21.5" {
			t.Errorf("Expected ActiveVersion '1.21.5', got '%s'", config.ActiveVersion)
		}

		if len(config.Versions) != 1 {
			t.Errorf("Expected 1 version, got %d", len(config.Versions))
		}

		if path, ok := config.Versions["1.21.5"]; !ok {
			t.Error("Expected version 1.21.5 in Versions map")
		} else if path != "/home/user/.gx/versions/go1.21.5" {
			t.Errorf("Expected path '/home/user/.gx/versions/go1.21.5', got '%s'", path)
		}
	})

	// 测试 4: 确保配置目录存在
	t.Run("EnsureConfigDir creates directory", func(t *testing.T) {
		newTempDir := t.TempDir()
		newConfigPath := filepath.Join(newTempDir, "subdir", "config.json")
		newStore := &fileStore{
			configPath: newConfigPath,
		}

		err := newStore.EnsureConfigDir()
		if err != nil {
			t.Fatalf("EnsureConfigDir() error = %v", err)
		}

		// 验证目录是否创建
		configDir := filepath.Dir(newConfigPath)
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			t.Error("Config directory was not created")
		}
	})
}

// TestNewStore 测试创建新的配置存储
func TestNewStore(t *testing.T) {
	store, err := NewStore()
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	if store == nil {
		t.Fatal("NewStore() returned nil")
	}
}
