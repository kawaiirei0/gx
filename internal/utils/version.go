package utils

import (
	"regexp"
	"strings"

	"github.com/yourusername/gx/pkg/errors"
)

// ValidateVersion 验证版本号格式
func ValidateVersion(version string) error {
	// 支持格式: 1.21.5, go1.21.5, 1.21
	pattern := `^(go)?(\d+)\.(\d+)(\.(\d+))?$`
	matched, err := regexp.MatchString(pattern, version)
	if err != nil {
		return err
	}
	if !matched {
		return errors.ErrInvalidVersion.WithMessage(version)
	}
	return nil
}

// NormalizeVersion 标准化版本号格式
func NormalizeVersion(version string) string {
	// 移除 "go" 前缀
	version = strings.TrimPrefix(version, "go")
	return version
}

// AddGoPrefix 添加 "go" 前缀
func AddGoPrefix(version string) string {
	if strings.HasPrefix(version, "go") {
		return version
	}
	return "go" + version
}

// CompareVersions 比较两个版本号
// 返回: -1 (v1 < v2), 0 (v1 == v2), 1 (v1 > v2)
func CompareVersions(v1, v2 string) int {
	v1 = NormalizeVersion(v1)
	v2 = NormalizeVersion(v2)

	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var p1, p2 int
		if i < len(parts1) {
			p1 = parseIntOrZero(parts1[i])
		}
		if i < len(parts2) {
			p2 = parseIntOrZero(parts2[i])
		}

		if p1 < p2 {
			return -1
		}
		if p1 > p2 {
			return 1
		}
	}

	return 0
}

func parseIntOrZero(s string) int {
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		}
	}
	return result
}
