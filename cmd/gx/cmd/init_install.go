package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/gx/internal/logger"
	"github.com/yourusername/gx/internal/ui"
)

var (
	initInstallForce bool
)

var initInstallCmd = &cobra.Command{
	Use:   "init-install",
	Short: "Install gx to system PATH",
	Long: `Install gx to system PATH so it can be used from anywhere.
This command will:
  - Copy gx binary to a system location
  - Add gx to PATH environment variable
  - Make the changes persistent across sessions

Example:
  gx init-install
  gx init-install --force`,
	RunE: runInitInstall,
}

func init() {
	rootCmd.AddCommand(initInstallCmd)
	initInstallCmd.Flags().BoolVarP(&initInstallForce, "force", "f", false, "force reinstall even if already installed")
}

func runInitInstall(cmd *cobra.Command, args []string) error {
	logger.Info("Init-install command started")

	messenger := ui.NewMessenger(os.Stdout)
	prompter := ui.NewPrompter(os.Stdin, os.Stdout)
	errorFormatter := ui.NewErrorFormatter(os.Stderr)

	messenger.Section("gx Installation")
	fmt.Println()

	// 获取当前可执行文件路径
	exePath, err := os.Executable()
	if err != nil {
		errorFormatter.Format(fmt.Errorf("failed to get executable path: %w", err))
		return err
	}

	// 解析符号链接
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		errorFormatter.Format(fmt.Errorf("failed to resolve executable path: %w", err))
		return err
	}

	messenger.Info(fmt.Sprintf("Current executable: %s", exePath))
	fmt.Println()

	// 检查是否已经在系统路径中
	if !initInstallForce {
		if isInSystemPath(exePath) {
			messenger.Success("gx is already installed in system PATH")
			messenger.Info("Use --force to reinstall")
			return nil
		}
	}

	// 根据操作系统选择安装方法
	switch runtime.GOOS {
	case "windows":
		return installWindows(exePath, messenger, prompter, errorFormatter)
	case "linux", "darwin":
		return installUnix(exePath, messenger, prompter, errorFormatter)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func installWindows(exePath string, messenger *ui.Messenger, prompter *ui.Prompter, errorFormatter *ui.ErrorFormatter) error {
	messenger.Info("Installing gx on Windows...")
	fmt.Println()

	// 目标安装目录
	installDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "gx", "bin")
	targetPath := filepath.Join(installDir, "gx.exe")

	messenger.Info(fmt.Sprintf("Installation directory: %s", installDir))

	// 确认安装
	confirmed, err := prompter.Confirm("Proceed with installation?", true)
	if err != nil {
		return err
	}
	if !confirmed {
		messenger.Info("Installation cancelled")
		return nil
	}

	fmt.Println()

	// 创建安装目录
	messenger.Info("Creating installation directory...")
	if err := os.MkdirAll(installDir, 0755); err != nil {
		errorFormatter.Format(fmt.Errorf("failed to create directory: %w", err))
		return err
	}

	// 复制可执行文件
	messenger.Info("Copying gx executable...")
	if err := copyFile(exePath, targetPath); err != nil {
		errorFormatter.Format(fmt.Errorf("failed to copy executable: %w", err))
		return err
	}

	// 添加到 PATH
	messenger.Info("Adding to system PATH...")
	if err := addToWindowsPath(installDir); err != nil {
		errorFormatter.Format(fmt.Errorf("failed to add to PATH: %w", err))
		return err
	}

	fmt.Println()
	messenger.Success("gx installed successfully!")
	fmt.Println()
	messenger.Info("Installation complete. Please restart your terminal or command prompt")
	messenger.Info("for the changes to take effect.")
	fmt.Println()
	messenger.Info("You can now use 'gx' from anywhere:")
	fmt.Println("  gx --version")
	fmt.Println("  gx install 1.21.5")

	logger.Info("Init-install completed successfully on Windows")
	return nil
}

func installUnix(exePath string, messenger *ui.Messenger, prompter *ui.Prompter, errorFormatter *ui.ErrorFormatter) error {
	osName := "Linux"
	if runtime.GOOS == "darwin" {
		osName = "macOS"
	}

	messenger.Info(fmt.Sprintf("Installing gx on %s...", osName))
	fmt.Println()

	// 检查是否有 sudo 权限
	hasSudo := checkSudoAccess()

	// 选择安装位置
	var installDir string
	var needsSudo bool

	if hasSudo {
		// 提供选项：系统级或用户级
		options := []string{
			"/usr/local/bin (system-wide, requires sudo)",
			fmt.Sprintf("%s/.local/bin (user only, no sudo)", os.Getenv("HOME")),
		}

		selected, err := prompter.Select("Select installation location:", options)
		if err != nil {
			return err
		}

		if selected == 0 {
			installDir = "/usr/local/bin"
			needsSudo = true
		} else {
			installDir = filepath.Join(os.Getenv("HOME"), ".local", "bin")
			needsSudo = false
		}
	} else {
		// 只能安装到用户目录
		installDir = filepath.Join(os.Getenv("HOME"), ".local", "bin")
		needsSudo = false
		messenger.Info(fmt.Sprintf("Installing to user directory: %s", installDir))
		messenger.Info("(sudo not available for system-wide installation)")
	}

	targetPath := filepath.Join(installDir, "gx")

	fmt.Println()
	messenger.Info(fmt.Sprintf("Installation directory: %s", installDir))

	// 确认安装
	confirmed, err := prompter.Confirm("Proceed with installation?", true)
	if err != nil {
		return err
	}
	if !confirmed {
		messenger.Info("Installation cancelled")
		return nil
	}

	fmt.Println()

	// 创建安装目录
	messenger.Info("Creating installation directory...")
	if err := os.MkdirAll(installDir, 0755); err != nil {
		errorFormatter.Format(fmt.Errorf("failed to create directory: %w", err))
		return err
	}

	// 复制可执行文件
	messenger.Info("Copying gx executable...")
	if needsSudo {
		// 使用 sudo 复制
		if err := copyFileWithSudo(exePath, targetPath); err != nil {
			errorFormatter.Format(fmt.Errorf("failed to copy executable: %w", err))
			return err
		}
	} else {
		if err := copyFile(exePath, targetPath); err != nil {
			errorFormatter.Format(fmt.Errorf("failed to copy executable: %w", err))
			return err
		}
	}

	// 设置可执行权限
	messenger.Info("Setting executable permissions...")
	if needsSudo {
		if err := chmodWithSudo(targetPath, 0755); err != nil {
			errorFormatter.Format(fmt.Errorf("failed to set permissions: %w", err))
			return err
		}
	} else {
		if err := os.Chmod(targetPath, 0755); err != nil {
			errorFormatter.Format(fmt.Errorf("failed to set permissions: %w", err))
			return err
		}
	}

	// 添加到 PATH（如果需要）
	if !needsSudo && installDir == filepath.Join(os.Getenv("HOME"), ".local", "bin") {
		messenger.Info("Adding to PATH in shell configuration...")
		if err := addToUnixPath(installDir); err != nil {
			messenger.Warning(fmt.Sprintf("Failed to add to PATH automatically: %v", err))
			messenger.Info("You may need to add it manually to your shell configuration")
		}
	}

	fmt.Println()
	messenger.Success("gx installed successfully!")
	fmt.Println()

	if needsSudo {
		messenger.Info("Installation complete. You can now use 'gx' from anywhere:")
	} else {
		messenger.Info("Installation complete. Please restart your terminal or run:")
		fmt.Println("  source ~/.bashrc  (bash)")
		fmt.Println("  source ~/.zshrc   (zsh)")
		fmt.Println()
		messenger.Info("Then you can use 'gx' from anywhere:")
	}

	fmt.Println("  gx --version")
	fmt.Println("  gx install 1.21.5")

	logger.Info("Init-install completed successfully on %s", osName)
	return nil
}

// 辅助函数

func isInSystemPath(exePath string) bool {
	// 检查可执行文件是否已经在 PATH 中的某个目录
	exeDir := filepath.Dir(exePath)

	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		if dir == exeDir {
			return true
		}
	}

	return false
}

func copyFile(src, dst string) error {
	// 读取源文件
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// 写入目标文件
	return os.WriteFile(dst, data, 0755)
}

func copyFileWithSudo(src, dst string) error {
	// 使用 sudo cp 命令
	cmd := fmt.Sprintf("sudo cp '%s' '%s'", src, dst)
	return executeShellCommand(cmd)
}

func chmodWithSudo(path string, mode os.FileMode) error {
	// 使用 sudo chmod 命令
	cmd := fmt.Sprintf("sudo chmod %o '%s'", mode, path)
	return executeShellCommand(cmd)
}

func executeShellCommand(cmd string) error {
	shell := "/bin/sh"
	if runtime.GOOS == "darwin" {
		shell = "/bin/bash"
	}

	process := exec.Command(shell, "-c", cmd)
	process.Stdin = os.Stdin
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr

	return process.Run()
}

func checkSudoAccess() bool {
	// 检查是否可以使用 sudo
	cmd := exec.Command("sudo", "-n", "true")
	err := cmd.Run()
	return err == nil
}

func addToWindowsPath(dir string) error {
	// 使用 PowerShell 添加到用户 PATH
	script := fmt.Sprintf(`
$path = [Environment]::GetEnvironmentVariable('Path', 'User')
if ($path -notlike '*%s*') {
    $newPath = $path + ';%s'
    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
    Write-Host 'Added to PATH'
} else {
    Write-Host 'Already in PATH'
}
`, dir, dir)

	cmd := exec.Command("powershell", "-Command", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func addToUnixPath(dir string) error {
	home := os.Getenv("HOME")

	// 检测使用的 shell
	shells := []struct {
		name       string
		configFile string
	}{
		{"bash", filepath.Join(home, ".bashrc")},
		{"zsh", filepath.Join(home, ".zshrc")},
		{"fish", filepath.Join(home, ".config", "fish", "config.fish")},
	}

	pathLine := fmt.Sprintf("\n# Added by gx\nexport PATH=\"%s:$PATH\"\n", dir)
	if runtime.GOOS == "darwin" {
		// macOS 也检查 .bash_profile 和 .zprofile
		shells = append(shells,
			struct {
				name       string
				configFile string
			}{"bash", filepath.Join(home, ".bash_profile")},
			struct {
				name       string
				configFile string
			}{"zsh", filepath.Join(home, ".zprofile")},
		)
	}

	updated := false
	for _, shell := range shells {
		if _, err := os.Stat(shell.configFile); err == nil {
			// 检查是否已经添加
			content, err := os.ReadFile(shell.configFile)
			if err != nil {
				continue
			}

			if !strings.Contains(string(content), dir) {
				// 追加到文件
				f, err := os.OpenFile(shell.configFile, os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					continue
				}
				defer f.Close()

				if _, err := f.WriteString(pathLine); err != nil {
					continue
				}

				updated = true
			}
		}
	}

	if !updated {
		return fmt.Errorf("no shell configuration file found")
	}

	return nil
}


