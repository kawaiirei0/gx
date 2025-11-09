//go:build windows

package environment

import (
	"github.com/kawaiirei0/gx/pkg/errors"
)

// setEnvUnix is a stub for Windows builds
func (m *manager) setEnvUnix(key, value string) error {
	return errors.ErrPlatformNotSupported.WithMessage("Unix environment management not available on Windows")
}
