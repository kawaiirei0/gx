package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/gx/internal/logger"
	"github.com/yourusername/gx/internal/ui"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate-config",
	Short: "Migrate configuration to new format",
	Long: `Migrate configuration file to use consistent version number format.
This command will update version numbers to include the 'go' prefix.

Example:
  gx migrate-config`,
	RunE: runMigrate,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func runMigrate(cmd *cobra.Command, args []string) error {
	logger.Info("Config migration started")

	ctx, err := NewAppContext()
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	messenger := ui.NewMessenger(os.Stdout)
	errorFormatter := ui.NewErrorFormatter(os.Stderr)

	messenger.Section("Configuration Migration")
	fmt.Println()

	// 加载配置
	cfg, err := ctx.ConfigStore.Load()
	if err != nil {
		errorFormatter.Format(fmt.Errorf("failed to load config: %w", err))
		return err
	}

	messenger.Info("Checking configuration format...")
	fmt.Println()

	// 检查是否需要迁移
	needsMigration := false
	migratedVersions := make(map[string]string)
	var newActiveVersion string

	// 迁移 versions 映射
	for version, path := range cfg.Versions {
		normalizedVersion := version
		if !strings.HasPrefix(version, "go") {
			normalizedVersion = "go" + version
			needsMigration = true
			messenger.Info(fmt.Sprintf("  %s → %s", version, normalizedVersion))
		}
		migratedVersions[normalizedVersion] = path
	}

	// 迁移 active_version
	if cfg.ActiveVersion != "" && !strings.HasPrefix(cfg.ActiveVersion, "go") {
		newActiveVersion = "go" + cfg.ActiveVersion
		needsMigration = true
		fmt.Println()
		messenger.Info(fmt.Sprintf("Active version: %s → %s", cfg.ActiveVersion, newActiveVersion))
	} else {
		newActiveVersion = cfg.ActiveVersion
	}

	if !needsMigration {
		fmt.Println()
		messenger.Success("Configuration is already in the correct format")
		return nil
	}

	fmt.Println()
	messenger.Warning("Configuration needs migration")
	fmt.Println()

	// 更新配置
	cfg.Versions = migratedVersions
	cfg.ActiveVersion = newActiveVersion

	// 保存配置
	messenger.Info("Saving migrated configuration...")
	if err := ctx.ConfigStore.Save(cfg); err != nil {
		errorFormatter.Format(fmt.Errorf("failed to save config: %w", err))
		return err
	}

	fmt.Println()
	messenger.Success("Configuration migrated successfully!")
	fmt.Println()
	messenger.Info("You can now use gx commands normally:")
	fmt.Println("  gx list")
	fmt.Println("  gx current")
	fmt.Println("  gx use <version>")

	logger.Info("Config migration completed successfully")
	return nil
}
