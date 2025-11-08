package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/gx/internal/ui"
)

var (
	listRemote bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed Go versions",
	Long: `List all Go versions installed by gx.
Use --remote flag to list available versions from the official Go distribution.

Example:
  gx list
  gx list --remote`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&listRemote, "remote", "r", false, "list available remote versions")
}

func runList(cmd *cobra.Command, args []string) error {
	ctx, err := NewAppContext()
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	if listRemote {
		return listRemoteVersions(ctx)
	}

	return listInstalledVersions(ctx)
}

func listInstalledVersions(ctx *AppContext) error {
	messenger := ui.NewMessenger(os.Stdout)
	errorFormatter := ui.NewErrorFormatter(os.Stderr)

	versions, err := ctx.VersionManager.DetectInstalled()
	if err != nil {
		errorFormatter.Format(err)
		return err
	}

	if len(versions) == 0 {
		messenger.Warning("No Go versions installed by gx")
		fmt.Println()
		messenger.Info("To install a version, run:")
		fmt.Println("  gx install <version>")
		return nil
	}

	messenger.Section("Installed Go Versions")
	fmt.Println()

	// 按版本号排序
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Version < versions[j].Version
	})

	// 准备表格数据
	if verbose {
		headers := []string{"Status", "Version", "Path", "Installed"}
		rows := make([][]string, len(versions))

		for i, v := range versions {
			status := " "
			if v.IsActive {
				status = "✓"
			}

			installDate := ""
			if !v.InstallDate.IsZero() {
				installDate = v.InstallDate.Format("2006-01-02")
			}

			rows[i] = []string{
				status,
				strings.TrimPrefix(v.Version, "go"),
				v.Path,
				installDate,
			}
		}

		messenger.Table(headers, rows)
	} else {
		// 简单列表显示
		for _, v := range versions {
			marker := " "
			status := ""
			if v.IsActive {
				marker = "✓"
				status = " (active)"
			}
			fmt.Printf("%s %s%s\n", marker, strings.TrimPrefix(v.Version, "go"), status)
		}
	}

	fmt.Println()
	if !verbose {
		messenger.Info("Use --verbose flag for more details")
	}

	return nil
}

func listRemoteVersions(ctx *AppContext) error {
	messenger := ui.NewMessenger(os.Stdout)
	errorFormatter := ui.NewErrorFormatter(os.Stderr)

	spinner := ui.NewSpinner(os.Stdout, "Fetching available Go versions...")

	// 模拟加载动画（在实际网络请求期间）
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				spinner.Tick()
				// 短暂延迟
				for i := 0; i < 100000000; i++ {
				}
			}
		}
	}()

	versions, err := ctx.VersionManager.ListAvailable()
	done <- true
	spinner.Clear()

	if err != nil {
		errorFormatter.Format(err)
		return err
	}

	if len(versions) == 0 {
		messenger.Warning("No versions available")
		return nil
	}

	messenger.Section("Available Go Versions")
	fmt.Println()

	// 显示前 30 个版本
	maxDisplay := 30
	if len(versions) < maxDisplay {
		maxDisplay = len(versions)
	}

	// 分列显示
	columns := 3
	rows := (maxDisplay + columns - 1) / columns

	for row := 0; row < rows; row++ {
		for col := 0; col < columns; col++ {
			idx := row + col*rows
			if idx < maxDisplay {
				version := versions[idx]
				displayVersion := strings.TrimPrefix(version, "go")
				fmt.Printf("  %-15s", displayVersion)
			}
		}
		fmt.Println()
	}

	if len(versions) > maxDisplay {
		fmt.Println()
		messenger.Info(fmt.Sprintf("... and %d more versions", len(versions)-maxDisplay))
	}

	fmt.Println()
	messenger.Info("To install a version, run:")
	fmt.Println("  gx install <version>")

	return nil
}
