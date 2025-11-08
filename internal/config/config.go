package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/yourusername/gx/pkg/constants"
	"github.com/yourusername/gx/pkg/errors"
	"github.com/yourusername/gx/pkg/interfaces"
)

// fileStore 基于文件的配置存储实现
type fileStore struct {
	configPath string
}

// NewStore 创建新的配置存储
func NewStore() (interfaces.ConfigStore, error) {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return nil, errors.ErrStorageFailed.WithCause(err)
	}

	return &fileStore{
		configPath: configPath,
	}, nil
}

// Load 加载配置文件
func (s *fileStore) Load() (*interfaces.Config, error) {
	// 如果配置文件不存在，返回默认配置
	if _, err := os.Stat(s.configPath); os.IsNotExist(err) {
		return s.getDefaultConfig()
	}

	// 读取配置文件
	data, err := os.ReadFile(s.configPath)
	if err != nil {
		return nil, errors.ErrStorageFailed.
			WithCause(err).
			WithMessage("failed to read config file").
			WithContext("config_path", s.configPath)
	}

	// 解析 JSON
	var config interfaces.Config
	if err := json.Unmarshal(data, &config); err != nil {
		// 配置文件损坏，尝试恢复
		backupPath := s.configPath + ".backup"
		if _, backupErr := os.Stat(backupPath); backupErr == nil {
			// 备份文件存在，尝试从备份恢复
			backupData, backupReadErr := os.ReadFile(backupPath)
			if backupReadErr == nil {
				var backupConfig interfaces.Config
				if backupParseErr := json.Unmarshal(backupData, &backupConfig); backupParseErr == nil {
					// 备份文件有效，使用备份配置
					config = backupConfig
					// 尝试恢复配置文件
					os.WriteFile(s.configPath, backupData, 0644)
					
					if config.Versions == nil {
						config.Versions = make(map[string]string)
					}
					return &config, nil
				}
			}
		}
		
		// 无法恢复，返回错误
		return nil, errors.ErrConfigCorrupted.
			WithCause(err).
			WithMessage("failed to parse config file and no valid backup found").
			WithContext("config_path", s.configPath).
			AsRecoverable()
	}

	// 如果 Versions 为 nil，初始化为空 map
	if config.Versions == nil {
		config.Versions = make(map[string]string)
	}

	return &config, nil
}

// Save 保存配置到文件
func (s *fileStore) Save(config *interfaces.Config) error {
	// 确保配置目录存在
	if err := s.EnsureConfigDir(); err != nil {
		return err
	}

	// 序列化为 JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return errors.ErrStorageFailed.
			WithCause(err).
			WithMessage("failed to serialize config").
			WithContext("config_path", s.configPath)
	}

	// 如果配置文件已存在，先备份
	if _, err := os.Stat(s.configPath); err == nil {
		backupPath, backupErr := errors.BackupFile(s.configPath)
		if backupErr == nil {
			// 备份成功，保存成功后清理备份
			defer os.Remove(backupPath)
		} else {
			// 备份失败，记录警告但继续保存
			// 这是优雅降级的一个例子
		}
	}

	// 写入临时文件
	tmpPath := s.configPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return errors.ErrStorageFailed.
			WithCause(err).
			WithMessage("failed to write config file").
			WithContext("config_path", s.configPath).
			WithContext("temp_path", tmpPath)
	}

	// 原子性地替换配置文件
	if err := os.Rename(tmpPath, s.configPath); err != nil {
		// 清理临时文件
		os.Remove(tmpPath)
		return errors.ErrStorageFailed.
			WithCause(err).
			WithMessage("failed to replace config file").
			WithContext("config_path", s.configPath)
	}

	return nil
}

// EnsureConfigDir 确保配置目录存在
func (s *fileStore) EnsureConfigDir() error {
	configDir := filepath.Dir(s.configPath)

	// 检查目录是否存在
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		// 创建目录，权限 0755
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return errors.ErrStorageFailed.WithCause(err).WithMessage("failed to create config directory")
		}
	}

	return nil
}

// getDefaultConfig 获取默认配置
func (s *fileStore) getDefaultConfig() (*interfaces.Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.ErrStorageFailed.WithCause(err).WithMessage("failed to get home directory")
	}

	installPath := filepath.Join(homeDir, constants.DefaultInstallDir)

	return &interfaces.Config{
		ActiveVersion:   "",
		InstallPath:     installPath,
		Versions:        make(map[string]string),
		LastUpdateCheck: time.Time{},
	}, nil
}

// GetDefaultConfig 获取默认配置（保留向后兼容）
func GetDefaultConfig() (*interfaces.Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	installPath := filepath.Join(homeDir, constants.DefaultInstallDir)

	return &interfaces.Config{
		ActiveVersion:   "",
		InstallPath:     installPath,
		Versions:        make(map[string]string),
		LastUpdateCheck: time.Time{},
	}, nil
}

// GetConfigDir 获取配置目录路径
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, constants.ConfigDir), nil
}

// GetConfigFilePath 获取配置文件路径
func GetConfigFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, constants.ConfigFileName), nil
}
