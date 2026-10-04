//go:build windows

package ui

import (
	"fmt"
)

// TODO: Implement single instance locking for Windows using LockFileEx or a named mutex (CreateMutex).
func AcquireInstanceLock(lockFilePath string) (*InstanceLock, error) {
	return nil, fmt.Errorf("instance locking is not yet implemented on Windows")
}

func (l *InstanceLock) Release() {
	if l == nil || l.file == nil {
		return
	}
	_ = l.file.Close()
	l.file = nil
}
