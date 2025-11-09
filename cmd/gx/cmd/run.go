package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/kawaiirei0/gx/internal/wrapper"
)

var runCmd = &cobra.Command{
	Use:   "run [flags] <file.go> [arguments...]",
	Short: "Run a Go program",
	Long: `Compile and run a Go program using the active Go version.
This is a wrapper around 'go run' command.

Example:
  gx run main.go
  gx run main.go arg1 arg2
  gx run -race main.go`,
	DisableFlagParsing: true, // 禁用标志解析，让所有参数传递给 go run
	RunE:               runRun,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runRun(cmd *cobra.Command, args []string) error {
	ctx, err := NewAppContext()
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	if len(args) == 0 {
		return fmt.Errorf("no Go file specified")
	}

	// 执行 go run 命令
	err = ctx.CLIWrapper.Execute("run", args)
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
