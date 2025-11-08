package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/gx/internal/ui"
	"github.com/yourusername/gx/pkg/constants"
	"github.com/yourusername/gx/pkg/interfaces"
)

var (
	targetOS     string
	targetArch   string
	outputPath   string
	ldflags      string
	buildFlags   []string
	listPlatforms bool
)

var crossBuildCmd = &cobra.Command{
	Use:   "cross-build [source]",
	Short: "Cross-compile Go programs for different platforms",
	Long: `Cross-compile Go programs for different operating systems and architectures.
Supports building for Windows, Linux, and macOS from any platform.

Example:
  gx cross-build --os linux --arch amd64 -o myapp
  gx cross-build --os windows --arch amd64 -o myapp.exe .
  gx cross-build --list-platforms`,
	RunE: runCrossBuild,
}

func init() {
	rootCmd.AddCommand(crossBuildCmd)
	
	crossBuildCmd.Flags().StringVar(&targetOS, "os", "", "target operating system (windows, linux, darwin)")
	crossBuildCmd.Flags().StringVar(&targetArch, "arch", "", "target architecture (amd64, arm64, 386)")
	crossBuildCmd.Flags().StringVarP(&outputPath, "output", "o", "", "output file path")
	crossBuildCmd.Flags().StringVar(&ldflags, "ldflags", "", "linker flags")
	crossBuildCmd.Flags().StringSliceVar(&buildFlags, "flags", []string{}, "additional build flags")
	crossBuildCmd.Flags().BoolVar(&listPlatforms, "list-platforms", false, "list supported platforms")
}

func runCrossBuild(cmd *cobra.Command, args []string) error {
	ctx, err := NewAppContext()
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	messenger := ui.NewMessenger(os.Stdout)
	errorFormatter := ui.NewErrorFormatter(os.Stderr)

	// 如果请求列出平台，显示支持的平台列表
	if listPlatforms {
		return listSupportedPlatforms(ctx)
	}

	// 验证必需参数
	if targetOS == "" {
		errorFormatter.Format(fmt.Errorf("--os flag is required"))
		fmt.Println()
		messenger.Info("Use --list-platforms to see supported platforms")
		return fmt.Errorf("missing required flag")
	}

	if targetArch == "" {
		errorFormatter.Format(fmt.Errorf("--arch flag is required"))
		fmt.Println()
		messenger.Info("Use --list-platforms to see supported platforms")
		return fmt.Errorf("missing required flag")
	}

	// 获取源代码路径
	sourcePath := "."
	if len(args) > 0 {
		sourcePath = args[0]
	}

	// 构建配置
	buildConfig := interfaces.BuildConfig{
		SourcePath: sourcePath,
		OutputPath: outputPath,
		TargetOS:   targetOS,
		TargetArch: targetArch,
		LDFlags:    ldflags,
		BuildFlags: buildFlags,
	}

	// 执行跨平台构建
	messenger.Info(fmt.Sprintf("Building for %s/%s...", targetOS, targetArch))
	err = ctx.CrossBuilder.Build(buildConfig)
	if err != nil {
		errorFormatter.Format(err)
		return err
	}

	messenger.Success("Build completed successfully")
	if outputPath != "" {
		messenger.Info(fmt.Sprintf("Output: %s", outputPath))
	}

	return nil
}

func listSupportedPlatforms(ctx *AppContext) error {
	messenger := ui.NewMessenger(os.Stdout)
	platforms := ctx.CrossBuilder.GetSupportedPlatforms()

	messenger.Section("Supported Platforms")
	fmt.Println()

	// 按操作系统分组显示
	platformsByOS := make(map[string][]string)
	for _, p := range platforms {
		platformsByOS[p.OS] = append(platformsByOS[p.OS], p.Arch)
	}

	// 按顺序显示
	osOrder := []string{constants.OSWindows, constants.OSLinux, constants.OSDarwin}
	for _, os := range osOrder {
		if archs, ok := platformsByOS[os]; ok {
			fmt.Printf("  %s:\n", os)
			for _, arch := range archs {
				fmt.Printf("    • %s\n", arch)
			}
			fmt.Println()
		}
	}

	messenger.Section("Example Usage")
	examples := []string{
		"gx cross-build --os linux --arch amd64 -o myapp",
		"gx cross-build --os windows --arch amd64 -o myapp.exe",
		"gx cross-build --os darwin --arch arm64 -o myapp",
	}
	messenger.List(examples, "")

	return nil
}
