package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/gx/internal/logger"
	"github.com/yourusername/gx/internal/ui"
)

var (
	installInteractive bool
)

var installCmd = &cobra.Command{
	Use:   "install [version]",
	Short: "Install a specific Go version",
	Long: `Install a specific Go version from the official Go distribution.
If no version is specified, installs the latest stable version.

Example:
  gx install 1.21.5
  gx install        # installs latest version
  gx install -i     # interactive version selection`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVarP(&installInteractive, "interactive", "i", false, "interactive version selection")
}

func runInstall(cmd *cobra.Command, args []string) error {
	logger.Info("Install command started")
	
	ctx, err := NewAppContext()
	if err != nil {
		logger.Error("Failed to initialize context: %v", err)
		return fmt.Errorf("failed to initialize: %w", err)
	}

	messenger := ui.NewMessenger(os.Stdout)
	prompter := ui.NewPrompter(os.Stdin, os.Stdout)
	errorFormatter := ui.NewErrorFormatter(os.Stderr)

	var versionToInstall string

	// 交互式版本选择
	if installInteractive && len(args) == 0 {
		messenger.Info("Fetching available Go versions...")

		versions, err := ctx.VersionManager.ListAvailable()
		if err != nil {
			errorFormatter.Format(err)
			return err
		}

		// 只显示稳定版本（前20个）
		stableVersions := versions
		if len(stableVersions) > 20 {
			stableVersions = stableVersions[:20]
		}

		// 格式化版本显示
		displayVersions := make([]string, len(stableVersions))
		for i, v := range stableVersions {
			displayVersions[i] = strings.TrimPrefix(v, "go")
		}

		selected, err := prompter.SelectVersion(displayVersions, 10)
		if err != nil {
			return err
		}

		versionToInstall = selected
		if !strings.HasPrefix(versionToInstall, "go") {
			versionToInstall = "go" + versionToInstall
		}
	} else if len(args) == 0 {
		// 如果没有指定版本，获取最新版本
		messenger.Info("Fetching latest Go version...")
		latest, err := ctx.VersionManager.GetLatest()
		if err != nil {
			errorFormatter.Format(err)
			return err
		}
		versionToInstall = latest
		messenger.Info(fmt.Sprintf("Latest version: %s", strings.TrimPrefix(versionToInstall, "go")))

		// 确认安装
		confirmed, err := prompter.Confirm(fmt.Sprintf("Install Go %s?", strings.TrimPrefix(versionToInstall, "go")), true)
		if err != nil {
			return err
		}
		if !confirmed {
			messenger.Info("Installation cancelled")
			return nil
		}
	} else {
		versionToInstall = args[0]
		// 规范化版本号
		if !strings.HasPrefix(versionToInstall, "go") {
			versionToInstall = "go" + versionToInstall
		}
	}

	messenger.Info(fmt.Sprintf("Installing Go %s...", strings.TrimPrefix(versionToInstall, "go")))

	// 创建进度条
	var progressBar *ui.ProgressBar

	// 创建进度回调
	progressCallback := func(downloaded, total int64) {
		if progressBar == nil && total > 0 {
			progressBar = ui.NewProgressBar(os.Stdout, total, "Downloading")
		}
		if progressBar != nil {
			progressBar.Update(downloaded)
		}
	}

	// 执行安装
	err = ctx.VersionManager.Install(versionToInstall, progressCallback)
	if err != nil {
		if progressBar != nil {
			fmt.Println() // 换行
		}
		errorFormatter.Format(err)
		return err
	}

	// 完成进度条
	if progressBar != nil {
		progressBar.Finish()
	}

	messenger.Success(fmt.Sprintf("Go %s installed successfully", strings.TrimPrefix(versionToInstall, "go")))
	fmt.Println()
	messenger.Info("To use this version, run:")
	fmt.Printf("  gx use %s\n", strings.TrimPrefix(versionToInstall, "go"))

	logger.Info("Install command completed successfully for version %s", versionToInstall)
	return nil
}
