package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/gx/internal/wrapper"
)

var buildCmd = &cobra.Command{
	Use:   "build [flags] [packages]",
	Short: "Compile packages and dependencies",
	Long: `Compile packages and dependencies using the active Go version.
This is a wrapper around 'go build' command.

Example:
  gx build
  gx build main.go
  gx build -o myapp
  gx build -ldflags="-s -w" .`,
	DisableFlagParsing: true, // 禁用标志解析，让所有参数传递给 go build
	RunE:               runBuild,
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

func runBuild(cmd *cobra.Command, args []string) error {
	ctx, err := NewAppContext()
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	// 执行 go build 命令
	err = ctx.CLIWrapper.Execute("build", args)
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
