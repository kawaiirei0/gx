//go:build linux || darwin

package environment

import (
	"github.com/yourusername/gx/pkg/errors"
)

// setEnvWindows is a stub for Unix builds
func (m *manager) setEnvWindows(key, value string) error {
	return errors.ErrPlatformNotSupported.WithMessage("Windows environment management not available on Unix")
}
