//go:build windows

package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

func AcquireInstanceLock(lockFilePath string) (*InstanceLock, error) {
	dir := filepath.Dir(lockFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("could not create lock directory: %w", err)
	}

	file, err := os.OpenFile(lockFilePath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("could not open lock file: %w", err)
	}

	var overlapped windows.Overlapped
	err = windows.LockFileEx(
		windows.Handle(file.Fd()),
		windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY,
		0,
		1,
		0,
		&overlapped,
	)
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
	var overlapped windows.Overlapped
	_ = windows.UnlockFileEx(windows.Handle(l.file.Fd()), 0, 1, 0, &overlapped)
	_ = l.file.Close()
	l.file = nil
}
