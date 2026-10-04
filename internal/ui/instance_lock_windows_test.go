//go:build windows

package ui

import (
	"testing"
)

func TestAcquireInstanceLockWindows(t *testing.T) {
	_, err := AcquireInstanceLock("dummy_path")
	if err == nil {
		t.Fatalf("expected error on Windows before implementation, got nil")
	}
}

func TestInstanceLockNilAndDoubleReleaseWindows(t *testing.T) {
	var nilLock *InstanceLock
	nilLock.Release()

	lock := &InstanceLock{}
	lock.Release()
	lock.Release()
}
