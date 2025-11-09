package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/kawaiirei0/gx/internal/logger"
	"github.com/kawaiirei0/gx/pkg/constants"
)

var (
	// 全局标志
	verbose bool
	config  string
	
	// 版本信息（由 main 包设置）
	appVersion   = "dev"
	appCommit    = "unknown"
	appBuildDate = "unknown"
)

// rootCmd 代表基础命令
var rootCmd = &cobra.Command{
	Use:   constants.AppName,
	Short: "Go version manager and development tool",
	Long: `gx is a cross-platform Go version manager and development tool.
It simplifies Go environment installation, version switching, and cross-platform compilation.`,
	Version: constants.AppVersion,
}

// SetVersionInfo 设置版本信息
func SetVersionInfo(v, c, bd string) {
	appVersion = v
	appCommit = c
	appBuildDate = bd
	// 更新 rootCmd 的版本信息
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", appVersion, appCommit, appBuildDate)
}

// Execute 执行根命令
func Execute() {
	// 初始化日志记录器
	if err := logger.Init(); err != nil {
		// 日志初始化失败不影响程序运行，只打印警告
		fmt.Fprintf(os.Stderr, "Warning: failed to initialize logger: %v\n", err)
	}
	defer logger.Close()

	logger.Info("gx %s started", appVersion)
	
	if err := rootCmd.Execute(); err != nil {
		logger.Error("Command execution failed: %v", err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	
	logger.Info("gx command completed successfully")
}

func init() {
	// 全局标志
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&config, "config", "", "config file (default is $HOME/.gx/config.json)")
	
	// 设置 PersistentPreRun 来处理 verbose 标志
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if verbose {
			logger.SetLevel(logger.LevelDebug)
			logger.Debug("Verbose mode enabled")
		}
	}
}
