package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/gx/internal/ui"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the current active Go version",
	Long: `Display the currently active Go version managed by gx.

Example:
  gx current`,
	RunE: runCurrent,
}

func init() {
	rootCmd.AddCommand(currentCmd)
}

func runCurrent(cmd *cobra.Command, args []string) error {
	ctx, err := NewAppContext()
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	messenger := ui.NewMessenger(os.Stdout)
	errorFormatter := ui.NewErrorFormatter(os.Stderr)

	activeVersion, err := ctx.VersionManager.GetActive()
	if err != nil {
		errorFormatter.Format(err)
		return err
	}

	messenger.Success(fmt.Sprintf("Current Go version: %s", strings.TrimPrefix(activeVersion.Version, "go")))

	if verbose {
		fmt.Println()
		messenger.Info(fmt.Sprintf("Installation path: %s", activeVersion.Path))
	}

	return nil
}
