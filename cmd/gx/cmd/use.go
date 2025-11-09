package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/kawaiirei0/gx/internal/ui"
)

var (
	useInteractive bool
)

var useCmd = &cobra.Command{
	Use:   "use [version]",
	Short: "Switch to a specific Go version",
	Long: `Switch to a specific Go version that has been installed.
This updates the environment variables to point to the selected version.

Example:
  gx use 1.21.5
  gx use go1.21.5
  gx use -i         # interactive version selection`,
	Args: cobra.MaximumNArgs(1),
	RunE: runUse,
}

func init() {
	rootCmd.AddCommand(useCmd)
	useCmd.Flags().BoolVarP(&useInteractive, "interactive", "i", false, "interactive version selection")
}

func runUse(cmd *cobra.Command, args []string) error {
	ctx, err := NewAppContext()
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	messenger := ui.NewMessenger(os.Stdout)
	prompter := ui.NewPrompter(os.Stdin, os.Stdout)
	errorFormatter := ui.NewErrorFormatter(os.Stderr)

	var version string

	// 交互式版本选择
	if (useInteractive || len(args) == 0) {
		versions, err := ctx.VersionManager.DetectInstalled()
		if err != nil {
			errorFormatter.Format(err)
			return err
		}

		if len(versions) == 0 {
			messenger.Warning("No Go versions installed")
			messenger.Info("Install a version first using:")
			fmt.Println("  gx install <version>")
			return nil
		}

		// 构建版本选项列表
		versionOptions := make([]string, len(versions))
		for i, v := range versions {
			marker := ""
			if v.IsActive {
				marker = " (current)"
			}
			versionOptions[i] = fmt.Sprintf("%s%s", v.Version, marker)
		}

		selected, err := prompter.Select("Select a Go version:", versionOptions)
		if err != nil {
			return err
		}

		version = versions[selected].Version
	} else {
		version = args[0]
		// 规范化版本号
		if !strings.HasPrefix(version, "go") {
			version = "go" + version
		}
	}

	messenger.Info(fmt.Sprintf("Switching to Go %s...", strings.TrimPrefix(version, "go")))

	err = ctx.VersionManager.SwitchTo(version)
	if err != nil {
		errorFormatter.Format(err)
		return err
	}

	messenger.Success(fmt.Sprintf("Now using Go %s", strings.TrimPrefix(version, "go")))
	fmt.Println()

	// 根据操作系统提供不同的提示
	if runtime.GOOS == "windows" {
		messenger.Info("Note: You may need to restart your terminal or command prompt")
		messenger.Info("for the environment changes to take effect.")
	} else {
		messenger.Info("Note: You may need to restart your terminal or run:")
		fmt.Println("  source ~/.bashrc  (bash)")
		fmt.Println("  source ~/.zshrc   (zsh)")
		messenger.Info("for the environment changes to take effect.")
	}

	return nil
}
