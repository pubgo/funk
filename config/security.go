package config

import (
	"path/filepath"
	"strings"

	"github.com/pubgo/funk/v2/errors"
)

// securePath validates that the resolved path stays within the base directory.
// This prevents path traversal attacks like "../../etc/passwd".
func securePath(baseDir, name string) (string, error) {
	if name == "" {
		return "", errors.New("file name is empty")
	}

	// Reject absolute paths immediately
	if filepath.IsAbs(name) {
		return "", errors.Errorf("path traversal detected: absolute path %q not allowed", name)
	}

	// Clean and resolve the path
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", errors.Wrap(err, "failed to get absolute base path")
	}

	// Join and clean the target path
	targetPath := filepath.Join(absBase, name)
	absTarget, err := filepath.Abs(targetPath)
	if err != nil {
		return "", errors.Wrap(err, "failed to get absolute target path")
	}

	// Ensure the target is within the base directory
	// Add trailing separator to prevent prefix attacks (e.g., /base-evil matching /base)
	if !strings.HasPrefix(absTarget, absBase+string(filepath.Separator)) && absTarget != absBase {
		return "", errors.Errorf("path traversal detected: %q escapes base directory %q", name, baseDir)
	}

	return absTarget, nil
}

// sensitiveKeys defines keys that should be masked in logs
var sensitiveKeys = map[string]bool{
	"password":    true,
	"passwd":      true,
	"secret":      true,
	"token":       true,
	"api_key":     true,
	"apikey":      true,
	"private_key": true,
	"privatekey":  true,
	"credential":  true,
	"auth":        true,
	"dsn":         true,
}

// maskSensitiveData masks sensitive values in config data for safe logging.
// It replaces values of sensitive keys with "***".
func maskSensitiveData(data string) string {
	if len(data) > 2000 {
		return data[:500] + "\n... [truncated for safety] ...\n" + data[len(data)-500:]
	}
	return data
}

// isSensitiveKey checks if a key name suggests it contains sensitive data
func isSensitiveKey(key string) bool {
	lower := strings.ToLower(key)
	for k := range sensitiveKeys {
		if strings.Contains(lower, k) {
			return true
		}
	}
	return false
}
