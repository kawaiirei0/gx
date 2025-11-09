package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/kawaiirei0/gx/internal/ui"
)

var (
	autoSwitch bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update to the latest Go version",
	Long: `Install the latest stable Go version.
Optionally switch to the new version automatically with --switch flag.

Example:
  gx update
  gx update --switch`,
	RunE: runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVarP(&autoSwitch, "switch", "s", false, "automatically switch to the new version after installation")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	ctx, err := NewAppContext()
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	messenger := ui.NewMessenger(os.Stdout)
	prompter := ui.NewPrompter(os.Stdin, os.Stdout)
	errorFormatter := ui.NewErrorFormatter(os.Stderr)

	messenger.Info("Checking for the latest Go version...")

	latest, err := ctx.VersionManager.GetLatest()
	if err != nil {
		errorFormatter.Format(err)
		return err
	}

	messenger.Info(fmt.Sprintf("Latest version: %s", strings.TrimPrefix(latest, "go")))

	// 检查是否已安装
	versions, err := ctx.VersionManager.DetectInstalled()
	if err != nil {
		errorFormatter.Format(err)
		return err
	}

	alreadyInstalled := false
	isActive := false
	for _, v := range versions {
		if v.Version == latest {
			alreadyInstalled = true
			isActive = v.IsActive
			break
		}
	}

	if alreadyInstalled && isActive {
		messenger.Success(fmt.Sprintf("You are already using the latest version (%s)", strings.TrimPrefix(latest, "go")))
		return nil
	}

	if alreadyInstalled {
		messenger.Success(fmt.Sprintf("Latest version (%s) is already installed", strings.TrimPrefix(latest, "go")))

		// 询问是否切换
		if !autoSwitch {
			confirmed, err := prompter.Confirm(
				fmt.Sprintf("Switch to Go %s now?", strings.TrimPrefix(latest, "go")),
				true,
			)
			if err != nil {
				return err
			}
			autoSwitch = confirmed
		}

		if autoSwitch {
			messenger.Info(fmt.Sprintf("Switching to %s...", strings.TrimPrefix(latest, "go")))
			err = ctx.VersionManager.SwitchTo(latest)
			if err != nil {
				errorFormatter.Format(err)
				return err
			}
			messenger.Success(fmt.Sprintf("Now using Go %s", strings.TrimPrefix(latest, "go")))
		} else {
			fmt.Println()
			messenger.Info("To use this version later, run:")
			fmt.Printf("  gx use %s\n", strings.TrimPrefix(latest, "go"))
		}
		return nil
	}

	messenger.Info(fmt.Sprintf("Installing Go %s...", strings.TrimPrefix(latest, "go")))

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
	err = ctx.VersionManager.Install(latest, progressCallback)
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

	messenger.Success(fmt.Sprintf("Go %s installed successfully", strings.TrimPrefix(latest, "go")))

	// 询问是否切换
	if !autoSwitch {
		fmt.Println()
		confirmed, err := prompter.Confirm(
			fmt.Sprintf("Switch to Go %s now?", strings.TrimPrefix(latest, "go")),
			true,
		)
		if err != nil {
			return err
		}
		autoSwitch = confirmed
	}

	// 如果设置了自动切换标志，切换到新版本
	if autoSwitch {
		messenger.Info(fmt.Sprintf("Switching to %s...", strings.TrimPrefix(latest, "go")))
		err = ctx.VersionManager.SwitchTo(latest)
		if err != nil {
			errorFormatter.Format(err)
			return err
		}
		messenger.Success(fmt.Sprintf("Now using Go %s", strings.TrimPrefix(latest, "go")))
	} else {
		fmt.Println()
		messenger.Info("To use this version, run:")
		fmt.Printf("  gx use %s\n", strings.TrimPrefix(latest, "go"))
	}

	return nil
}
