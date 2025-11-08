package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/gx/internal/wrapper"
)

var testCmd = &cobra.Command{
	Use:   "test [flags] [packages]",
	Short: "Test packages",
	Long: `Test packages using the active Go version.
This is a wrapper around 'go test' command.

Example:
  gx test
  gx test ./...
  gx test -v
  gx test -cover ./...`,
	DisableFlagParsing: true, // 禁用标志解析，让所有参数传递给 go test
	RunE:               runTest,
}

func init() {
	rootCmd.AddCommand(testCmd)
}

func runTest(cmd *cobra.Command, args []string) error {
	ctx, err := NewAppContext()
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	// 执行 go test 命令
	err = ctx.CLIWrapper.Execute("test", args)
	if err != nil {
		// 检查是否是退出错误
		if exitErr, ok := err.(*wrapper.ExitError); ok {
			// 保留原始退出码
			os.Exit(exitErr.GetExitCode())
		}
		return err
	}

	return nil
}
