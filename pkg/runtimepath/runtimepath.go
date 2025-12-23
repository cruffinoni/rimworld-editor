package runtimepath

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const defaultDirName = ".therapon"

var (
	baseDir string
	initErr error
	once    sync.Once
)

// BaseDir returns the directory used for all runtime state. This is always
// ~/.therapon within the current user's home directory. The directory is
// created on first use. Panics if it cannot be determined.
func BaseDir() string {
	dir, err := base()
	if err != nil {
		panic(err)
	}
	return dir
}

// LogsPath returns the path to the logs directory combined with any additional elements.
func LogsPath(elements ...string) string {
	dir := filepath.Join(BaseDir(), "logs")
	if len(elements) == 0 {
		return dir
	}
	return filepath.Join(append([]string{dir}, elements...)...)
}

// ConfigDir returns the runtime configuration directory, creating it if necessary.
func ConfigDir() string {
	dir := filepath.Join(BaseDir(), "config")
	_ = os.MkdirAll(dir, 0o755)
	return dir
}

func base() (string, error) {
	once.Do(func() {
		home, err := os.UserHomeDir()
		if err != nil || home == "" {
			initErr = fmt.Errorf("determine user home directory: %w", err)
			return
		}
		dir := filepath.Join(home, defaultDirName)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			initErr = fmt.Errorf("ensure therapon home %s: %w", dir, err)
			return
		}
		baseDir = dir
	})
	return baseDir, initErr
}
