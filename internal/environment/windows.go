//go:build windows

package environment

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/yourusername/gx/pkg/errors"
)

// setEnvWindows 在 Windows 上持久化设置环境变量
func (m *manager) setEnvWindows(key, value string) error {
	// 使用 setx 命令设置用户环境变量
	// setx 会将环境变量写入注册表，使其持久化
	cmd := exec.Command("setx", key, value)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.ErrOperationFailed.
			WithCause(err).
			WithMessage(fmt.Sprintf("failed to set %s on Windows: %s", key, string(output)))
	}

	return nil
}

// setEnvWindowsSystem 在 Windows 上设置系统级环境变量（需要管理员权限）
func (m *manager) setEnvWindowsSystem(key, value string) error {
	// 使用 PowerShell 设置系统环境变量
	psScript := fmt.Sprintf(`[Environment]::SetEnvironmentVariable('%s', '%s', 'Machine')`, key, value)
	
	cmd := exec.Command("powershell", "-Command", psScript)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.ErrOperationFailed.
			WithCause(err).
			WithMessage(fmt.Sprintf("failed to set system %s on Windows: %s", key, string(output)))
	}

	return nil
}

// getEnvWindows 从 Windows 注册表获取环境变量
func (m *manager) getEnvWindows(key string) (string, error) {
	// 使用 PowerShell 读取用户环境变量
	psScript := fmt.Sprintf(`[Environment]::GetEnvironmentVariable('%s', 'User')`, key)
	
	cmd := exec.Command("powershell", "-Command", psScript)
	
	output, err := cmd.Output()
	if err != nil {
		return "", errors.ErrOperationFailed.
			WithCause(err).
			WithMessage(fmt.Sprintf("failed to get %s from Windows registry", key))
	}

	return strings.TrimSpace(string(output)), nil
}

// refreshEnvWindows 刷新 Windows 环境变量（通知系统环境变量已更改）
func (m *manager) refreshEnvWindows() error {
	// 广播 WM_SETTINGCHANGE 消息通知系统环境变量已更改
	// 这需要使用 Windows API，这里使用 PowerShell 作为替代方案
	psScript := `
		$HWND_BROADCAST = [IntPtr] 0xffff
		$WM_SETTINGCHANGE = 0x1a
		if (-not ("Win32.NativeMethods" -as [Type])) {
			Add-Type -Namespace Win32 -Name NativeMethods -MemberDefinition @"
				[DllImport("user32.dll", SetLastError = true, CharSet = CharSet.Auto)]
				public static extern IntPtr SendMessageTimeout(
					IntPtr hWnd, uint Msg, UIntPtr wParam, string lParam,
					uint fuFlags, uint uTimeout, out UIntPtr lpdwResult);
"@
		}
		$result = [UIntPtr]::Zero
		[Win32.NativeMethods]::SendMessageTimeout($HWND_BROADCAST, $WM_SETTINGCHANGE, [UIntPtr]::Zero, "Environment", 2, 5000, [ref] $result)
	`
	
	cmd := exec.Command("powershell", "-Command", psScript)
	
	if err := cmd.Run(); err != nil {
		// 刷新失败不是致命错误，只记录警告
		return nil
	}

	return nil
}
