package ui

import (
	"os"
	"path/filepath"

	"github.com/SayHeyD/sops-age-manager/pkg/config"
)

type InstanceLock struct {
	file *os.File
}

func getLockFilePath(cfg *config.Config) string {
	if cfg != nil && cfg.Path != "" {
		return filepath.Join(filepath.Dir(cfg.Path), "sam.lock")
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		return filepath.Join(homeDir, ".sops-age-manager", "sam.lock")
	}

	return filepath.Join(os.TempDir(), "sops-age-manager.lock")
}
