package crossbuilder

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/yourusername/gx/internal/logger"
	"github.com/yourusername/gx/pkg/constants"
	"github.com/yourusername/gx/pkg/errors"
	"github.com/yourusername/gx/pkg/interfaces"
)

// crossBuilder 实现 CrossBuilder 接口
type crossBuilder struct {
	versionManager interfaces.VersionManager
	platform       interfaces.PlatformAdapter
}

// NewCrossBuilder 创建新的跨平台构建器
func NewCrossBuilder(versionManager interfaces.VersionManager, platform interfaces.PlatformAdapter) interfaces.CrossBuilder {
	return &crossBuilder{
		versionManager: versionManager,
		platform:       platform,
	}
}

// supportedPlatforms 定义支持的平台组合
var supportedPlatforms = []interfaces.PlatformInfo{
	{OS: constants.OSWindows, Arch: constants.ArchAMD64},
	{OS: constants.OSWindows, Arch: constants.Arch386},
	{OS: constants.OSLinux, Arch: constants.ArchAMD64},
	{OS: constants.OSLinux, Arch: constants.ArchARM64},
	{OS: constants.OSLinux, Arch: constants.Arch386},
	{OS: constants.OSDarwin, Arch: constants.ArchAMD64},
	{OS: constants.OSDarwin, Arch: constants.ArchARM64},
}

// Build 执行跨平台构建
func (cb *crossBuilder) Build(config interfaces.BuildConfig) error {
	logger.Info("Starting cross-platform build for %s/%s", config.TargetOS, config.TargetArch)
	
	// 验证平台组合
	if err := cb.ValidatePlatform(config.TargetOS, config.TargetArch); err != nil {
		logger.Error("Invalid platform: %v", err)
		return err
	}

	// 验证源代码路径
	if config.SourcePath == "" {
		config.SourcePath = "."
	}

	// 检查源代码路径是否存在
	if _, err := os.Stat(config.SourcePath); os.IsNotExist(err) {
		logger.Error("Source path not found: %s", config.SourcePath)
		return errors.ErrNotFound.
			WithMessage(fmt.Sprintf("source path not found: %s", config.SourcePath)).
			WithContext("source_path", config.SourcePath)
	}
	logger.Debug("Source path: %s", config.SourcePath)

	// 获取 Go 可执行文件路径
	goExe, err := cb.getGoExecutable()
	if err != nil {
		logger.Error("Failed to get Go executable: %v", err)
		return errors.Wrap(err, "OPERATION_FAILED", "failed to get Go executable for cross-build").
			WithContext("target_os", config.TargetOS).
			WithContext("target_arch", config.TargetArch)
	}

	// 处理输出路径，添加平台特定的文件扩展名
	outputPath := cb.normalizeOutputPath(config.OutputPath, config.TargetOS)

	// 如果指定了输出路径，确保输出目录存在
	if outputPath != "" {
		outputDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return errors.ErrOperationFailed.
				WithCause(err).
				WithMessage("failed to create output directory").
				WithContext("output_dir", outputDir)
		}
	}

	// 构建命令参数
	args := []string{"build"}

	// 添加输出路径
	if outputPath != "" {
		args = append(args, "-o", outputPath)
	}

	// 添加链接器标志
	if config.LDFlags != "" {
		args = append(args, "-ldflags", config.LDFlags)
	}

	// 添加额外的构建标志
	if len(config.BuildFlags) > 0 {
		args = append(args, config.BuildFlags...)
	}

	// 添加源代码路径
	args = append(args, config.SourcePath)

	// 创建命令
	cmd := exec.Command(goExe, args...)

	// 设置环境变量
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", config.TargetOS))
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOARCH=%s", config.TargetArch))

	// 透传标准输出和错误流
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行构建命令
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return errors.ErrOperationFailed.
				WithCause(exitErr).
				WithMessage(fmt.Sprintf("build failed for %s/%s", config.TargetOS, config.TargetArch)).
				WithContext("target_os", config.TargetOS).
				WithContext("target_arch", config.TargetArch).
				WithContext("output_path", outputPath).
				WithContext("exit_code", exitErr.ExitCode())
		}
		return errors.ErrOperationFailed.
			WithCause(err).
			WithMessage("failed to execute build command").
			WithContext("target_os", config.TargetOS).
			WithContext("target_arch", config.TargetArch)
	}

	// 验证输出文件是否存在
	if outputPath != "" {
		if _, err := os.Stat(outputPath); err != nil {
			return errors.ErrOperationFailed.
				WithMessage("build completed but output file not found").
				WithContext("output_path", outputPath).
				WithContext("target_os", config.TargetOS).
				WithContext("target_arch", config.TargetArch)
		}
	}

	// 显示构建成功信息
	fmt.Printf("Build successful for %s/%s\n", config.TargetOS, config.TargetArch)
	if outputPath != "" {
		absPath, _ := filepath.Abs(outputPath)
		fmt.Printf("Output: %s\n", absPath)
	}

	return nil
}

// ValidatePlatform 验证平台组合是否支持
func (cb *crossBuilder) ValidatePlatform(targetOS, targetArch string) error {
	if targetOS == "" {
		return errors.ErrInvalidInput.WithMessage("target OS cannot be empty")
	}

	if targetArch == "" {
		return errors.ErrInvalidInput.WithMessage("target architecture cannot be empty")
	}

	// 检查平台组合是否在支持列表中
	for _, platform := range supportedPlatforms {
		if platform.OS == targetOS && platform.Arch == targetArch {
			return nil
		}
	}

	return errors.ErrPlatformNotSupported.WithMessage(
		fmt.Sprintf("platform %s/%s is not supported", targetOS, targetArch),
	)
}

// GetSupportedPlatforms 获取支持的平台组合列表
func (cb *crossBuilder) GetSupportedPlatforms() []interfaces.PlatformInfo {
	// 返回副本以防止外部修改
	platforms := make([]interfaces.PlatformInfo, len(supportedPlatforms))
	copy(platforms, supportedPlatforms)
	return platforms
}

// getGoExecutable 获取当前使用的 Go 可执行文件路径
func (cb *crossBuilder) getGoExecutable() (string, error) {
	// 获取当前激活的版本
	activeVersion, err := cb.versionManager.GetActive()
	if err != nil {
		return "", errors.ErrVersionNotFound.WithCause(err).WithMessage("no active Go version found")
	}

	// 验证版本路径是否存在
	if activeVersion.Path == "" {
		return "", errors.ErrVersionNotFound.WithMessage("active version path is empty")
	}

	// 构建 Go 可执行文件路径
	goExe := "go"
	if cb.platform.GetOS() == constants.OSWindows {
		goExe = "go.exe"
	}

	goPath := filepath.Join(activeVersion.Path, "bin", goExe)

	// 验证文件是否存在
	if _, err := os.Stat(goPath); os.IsNotExist(err) {
		return "", errors.ErrNotFound.WithMessage(fmt.Sprintf("go executable not found at %s", goPath))
	}

	return goPath, nil
}

// normalizeOutputPath 规范化输出路径，添加平台特定的文件扩展名
func (cb *crossBuilder) normalizeOutputPath(outputPath, targetOS string) string {
	if outputPath == "" {
		return ""
	}

	// 如果目标平台是 Windows，确保输出文件有 .exe 扩展名
	if targetOS == constants.OSWindows {
		if !strings.HasSuffix(strings.ToLower(outputPath), ".exe") {
			outputPath += ".exe"
		}
	} else {
		// 对于非 Windows 平台，移除 .exe 扩展名（如果有）
		if strings.HasSuffix(strings.ToLower(outputPath), ".exe") {
			outputPath = outputPath[:len(outputPath)-4]
		}
	}

	return outputPath
}
