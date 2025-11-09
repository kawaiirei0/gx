package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/kawaiirei0/gx/internal/logger"
	"github.com/kawaiirei0/gx/internal/ui"
)

var (
	doctorFix bool
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check and fix gx configuration issues",
	Long: `Diagnose and optionally fix common gx configuration problems.
This command will:
  - Check if configured versions actually exist
  - Verify active version is valid
  - Clean up invalid entries

Example:
  gx doctor           # check only
  gx doctor --fix     # check and fix`,
	RunE: runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	doctorCmd.Flags().BoolVarP(&doctorFix, "fix", "f", false, "automatically fix issues")
}

func runDoctor(cmd *cobra.Command, args []string) error {
	logger.Info("Doctor command started")

	ctx, err := NewAppContext()
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	messenger := ui.NewMessenger(os.Stdout)
	prompter := ui.NewPrompter(os.Stdin, os.Stdout)
	errorFormatter := ui.NewErrorFormatter(os.Stderr)

	messenger.Section("gx Configuration Doctor")
	fmt.Println()

	// 加载配置
	cfg, err := ctx.ConfigStore.Load()
	if err != nil {
		errorFormatter.Format(fmt.Errorf("failed to load config: %w", err))
		return err
	}

	messenger.Info(fmt.Sprintf("Install path: %s", cfg.InstallPath))
	fmt.Println()

	// 检查问题
	var issues []string
	invalidVersions := make(map[string]string)

	// 1. 检查版本目录是否存在
	messenger.Info("Checking configured versions...")
	for version, path := range cfg.Versions {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			issue := fmt.Sprintf("Version %s: path does not exist (%s)", version, path)
			issues = append(issues, issue)
			invalidVersions[version] = path
			messenger.Warning(fmt.Sprintf("  ✗ %s", version))
		} else {
			messenger.Info(fmt.Sprintf("  ✓ %s", version))
		}
	}

	// 2. 检查激活版本
	fmt.Println()
	messenger.Info("Checking active version...")
	if cfg.ActiveVersion != "" {
		if _, ok := cfg.Versions[cfg.ActiveVersion]; !ok {
			issue := fmt.Sprintf("Active version %s is not in versions list", cfg.ActiveVersion)
			issues = append(issues, issue)
			messenger.Warning(fmt.Sprintf("  ✗ Active version %s not found in configuration", cfg.ActiveVersion))
		} else if _, exists := invalidVersions[cfg.ActiveVersion]; exists {
			issue := fmt.Sprintf("Active version %s points to non-existent path", cfg.ActiveVersion)
			issues = append(issues, issue)
			messenger.Warning(fmt.Sprintf("  ✗ Active version %s path does not exist", cfg.ActiveVersion))
		} else {
			messenger.Info(fmt.Sprintf("  ✓ Active version: %s", cfg.ActiveVersion))
		}
	} else {
		messenger.Info("  No active version set")
	}

	// 显示结果
	fmt.Println()
	if len(issues) == 0 {
		messenger.Success("No issues found!")
		return nil
	}

	messenger.Warning(fmt.Sprintf("Found %d issue(s):", len(issues)))
	for i, issue := range issues {
		fmt.Printf("  %d. %s\n", i+1, issue)
	}

	// 修复问题
	if len(invalidVersions) > 0 {
		fmt.Println()
		shouldFix := doctorFix
		if !doctorFix {
			confirmed, err := prompter.Confirm("Do you want to fix these issues?", true)
			if err != nil {
				return err
			}
			shouldFix = confirmed
		}

		if shouldFix {
			fmt.Println()
			messenger.Info("Fixing issues...")

			// 删除无效的版本记录
			for version := range invalidVersions {
				delete(cfg.Versions, version)
				messenger.Info(fmt.Sprintf("  Removed invalid version: %s", version))
			}

			// 如果激活版本无效，清除它
			if cfg.ActiveVersion != "" {
				if _, exists := invalidVersions[cfg.ActiveVersion]; exists {
					cfg.ActiveVersion = ""
					messenger.Info("  Cleared invalid active version")
				}
			}

			// 保存配置
			if err := ctx.ConfigStore.Save(cfg); err != nil {
				errorFormatter.Format(fmt.Errorf("failed to save config: %w", err))
				return err
			}

			fmt.Println()
			messenger.Success("Issues fixed successfully!")
			fmt.Println()
			messenger.Info("You can now:")
			fmt.Println("  - Install Go versions: gx install <version>")
			fmt.Println("  - List versions: gx list")
		} else {
			fmt.Println()
			messenger.Info("No changes made. Run 'gx doctor --fix' to fix issues automatically.")
		}
	}

	logger.Info("Doctor command completed")
	return nil
}
