package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/kawaiirei0/gx/internal/ui"
)

var (
	uninstallForce bool
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall <version>",
	Short: "Uninstall a specific Go version",
	Long: `Uninstall a specific Go version managed by gx.
Cannot uninstall the currently active version.

Example:
  gx uninstall 1.21.5
  gx uninstall go1.21.5
  gx uninstall 1.21.5 --force`,
	Args: cobra.ExactArgs(1),
	RunE: runUninstall,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().BoolVarP(&uninstallForce, "force", "f", false, "skip confirmation prompt")
}

func runUninstall(cmd *cobra.Command, args []string) error {
	ctx, err := NewAppContext()
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	messenger := ui.NewMessenger(os.Stdout)
	prompter := ui.NewPrompter(os.Stdin, os.Stdout)
	errorFormatter := ui.NewErrorFormatter(os.Stderr)

	version := args[0]
	// 规范化版本号
	if !strings.HasPrefix(version, "go") {
		version = "go" + version
	}

	// 确认卸载（除非使用 --force）
	if !uninstallForce {
		confirmed, err := prompter.Confirm(
			fmt.Sprintf("Are you sure you want to uninstall Go %s?", strings.TrimPrefix(version, "go")),
			false,
		)
		if err != nil {
			return err
		}
		if !confirmed {
			messenger.Info("Uninstallation cancelled")
			return nil
		}
	}

	messenger.Info(fmt.Sprintf("Uninstalling Go %s...", strings.TrimPrefix(version, "go")))

	err = ctx.VersionManager.Uninstall(version)
	if err != nil {
		errorFormatter.Format(err)
		return err
	}

	messenger.Success(fmt.Sprintf("Go %s uninstalled successfully", strings.TrimPrefix(version, "go")))

	return nil
}
