package interfaces

// CLIWrapper 包装和转发 Go 原生命令
type CLIWrapper interface {
	// Execute 执行 Go 命令
	// command: Go 命令名称（如 "run", "build", "test"）
	// args: 命令参数
	Execute(command string, args []string) error

	// GetGoExecutable 获取当前使用的 Go 可执行文件路径
	GetGoExecutable() (string, error)
}
