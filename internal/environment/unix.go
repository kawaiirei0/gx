//go:build linux || darwin

package environment

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yourusername/gx/pkg/errors"
)

// setEnvUnix 在 Unix 系统（Linux/macOS）上持久化设置环境变量
func (m *manager) setEnvUnix(key, value string) error {
	homeDir, err := m.platform.GetHomeDir()
	if err != nil {
		return errors.ErrOperationFailed.WithCause(err).WithMessage("failed to get home directory")
	}

	// 对于 PATH，我们需要特殊处理：只保存 Go bin 路径
	valueToSave := value
	if key == "PATH" {
		// 从完整的 PATH 中提取 Go bin 路径（第一个路径）
		pathSeparator := ":"
		paths := strings.Split(value, pathSeparator)
		if len(paths) > 0 {
			valueToSave = paths[0] // 只保存 Go bin 路径
		}
	}

	// 检测使用的 shell
	shell := os.Getenv("SHELL")
	var rcFiles []string

	if strings.Contains(shell, "zsh") {
		rcFiles = []string{
			filepath.Join(homeDir, ".zshrc"),
			filepath.Join(homeDir, ".zprofile"),
		}
	} else if strings.Contains(shell, "bash") {
		rcFiles = []string{
			filepath.Join(homeDir, ".bashrc"),
			filepath.Join(homeDir, ".bash_profile"),
			filepath.Join(homeDir, ".profile"),
		}
	} else {
		// 默认尝试常见的配置文件
		rcFiles = []string{
			filepath.Join(homeDir, ".profile"),
			filepath.Join(homeDir, ".bashrc"),
		}
	}

	// 尝试更新每个存在的配置文件
	updated := false
	for _, rcFile := range rcFiles {
		if _, err := os.Stat(rcFile); err == nil {
			if err := m.updateShellRC(rcFile, key, valueToSave); err != nil {
				// 继续尝试其他文件
				continue
			}
			updated = true
		}
	}

	// 如果没有找到任何配置文件，创建 .profile
	if !updated {
		profilePath := filepath.Join(homeDir, ".profile")
		if err := m.updateShellRC(profilePath, key, valueToSave); err != nil {
			return err
		}
	}

	return nil
}

// updateShellRC 更新 shell 配置文件中的环境变量
func (m *manager) updateShellRC(rcFile, key, value string) error {
	// 读取现有内容
	var lines []string
	existingFile := false

	if _, err := os.Stat(rcFile); err == nil {
		existingFile = true
		file, err := os.Open(rcFile)
		if err != nil {
			return errors.ErrOperationFailed.WithCause(err).WithMessage(fmt.Sprintf("failed to open %s", rcFile))
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			return errors.ErrOperationFailed.WithCause(err).WithMessage(fmt.Sprintf("failed to read %s", rcFile))
		}
	}

	// 查找并更新或添加环境变量
	marker := fmt.Sprintf("# gx managed %s", key)
	exportLine := m.buildExportLine(key, value)
	
	found := false
	var newLines []string

	for i, line := range lines {
		// 跳过旧的 gx 管理的该环境变量
		if strings.Contains(line, marker) {
			// 跳过标记行和下一行（export 行）
			if i+1 < len(lines) {
				continue
			}
		}
		
		// 检查是否是旧的 export 行
		if strings.HasPrefix(strings.TrimSpace(line), fmt.Sprintf("export %s=", key)) && 
		   i > 0 && strings.Contains(lines[i-1], marker) {
			// 跳过，因为已经在上面处理了标记行
			continue
		}

		newLines = append(newLines, line)
	}

	// 添加新的环境变量设置
	if !found {
		// 添加空行（如果文件不为空且最后一行不是空行）
		if len(newLines) > 0 && newLines[len(newLines)-1] != "" {
			newLines = append(newLines, "")
		}
		
		newLines = append(newLines, marker)
		newLines = append(newLines, exportLine)
	}

	// 写回文件
	content := strings.Join(newLines, "\n") + "\n"
	
	// 如果是新文件，使用 0644 权限；否则保持原有权限
	perm := os.FileMode(0644)
	if existingFile {
		if info, err := os.Stat(rcFile); err == nil {
			perm = info.Mode()
		}
	}

	if err := os.WriteFile(rcFile, []byte(content), perm); err != nil {
		return errors.ErrOperationFailed.WithCause(err).WithMessage(fmt.Sprintf("failed to write %s", rcFile))
	}

	return nil
}

// buildExportLine 构建 export 语句
func (m *manager) buildExportLine(key, value string) string {
	if key == "PATH" {
		// 对于 PATH，value 是 Go bin 路径，需要追加到现有 PATH
		return fmt.Sprintf(`export %s="%s:$%s"`, key, value, key)
	}
	return fmt.Sprintf(`export %s="%s"`, key, value)
}

// removeFromShellRC 从 shell 配置文件中移除环境变量
func (m *manager) removeFromShellRC(rcFile, key string) error {
	// 读取现有内容
	if _, err := os.Stat(rcFile); os.IsNotExist(err) {
		return nil // 文件不存在，无需删除
	}

	file, err := os.Open(rcFile)
	if err != nil {
		return errors.ErrOperationFailed.WithCause(err).WithMessage(fmt.Sprintf("failed to open %s", rcFile))
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return errors.ErrOperationFailed.WithCause(err).WithMessage(fmt.Sprintf("failed to read %s", rcFile))
	}

	// 过滤掉 gx 管理的环境变量
	marker := fmt.Sprintf("# gx managed %s", key)
	var newLines []string
	skip := false

	for _, line := range lines {
		if strings.Contains(line, marker) {
			skip = true
			continue
		}
		
		if skip && strings.HasPrefix(strings.TrimSpace(line), fmt.Sprintf("export %s=", key)) {
			skip = false
			continue
		}
		
		newLines = append(newLines, line)
	}

	// 写回文件
	content := strings.Join(newLines, "\n") + "\n"
	
	info, err := os.Stat(rcFile)
	if err != nil {
		return errors.ErrOperationFailed.WithCause(err).WithMessage(fmt.Sprintf("failed to stat %s", rcFile))
	}

	if err := os.WriteFile(rcFile, []byte(content), info.Mode()); err != nil {
		return errors.ErrOperationFailed.WithCause(err).WithMessage(fmt.Sprintf("failed to write %s", rcFile))
	}

	return nil
}
