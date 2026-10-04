//go:build darwin

package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

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

func AcquireInstanceLock(lockFilePath string) (*InstanceLock, error) {
	dir := filepath.Dir(lockFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("could not create lock directory: %w", err)
	}

	file, err := os.OpenFile(lockFilePath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("could not open lock file: %w", err)
	}

	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("another instance is already running")
	}

	_ = file.Truncate(0)
	_, _ = file.Seek(0, 0)
	_, _ = fmt.Fprintf(file, "%d\n", os.Getpid())

	return &InstanceLock{file: file}, nil
}

func (l *InstanceLock) Release() {
	if l == nil || l.file == nil {
		return
	}
	_ = syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
	_ = l.file.Close()
	l.file = nil
}
