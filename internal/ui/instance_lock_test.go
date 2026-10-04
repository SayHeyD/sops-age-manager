package ui

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/test"
)

func TestAcquireInstanceLock(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	lockPath := filepath.Join(testDir.Path, "sam.lock")

	lock1, err := AcquireInstanceLock(lockPath)
	if err != nil {
		t.Fatalf("expected lock1 to succeed, got error: %v", err)
	}
	if lock1 == nil || lock1.file == nil {
		t.Fatalf("expected non-nil lock1 with open file")
	}
	defer func() {
		if lock1 != nil {
			lock1.Release()
		}
	}()

	// Verify PID was written to the lock file
	buf := make([]byte, 64)
	n, err := lock1.file.ReadAt(buf, 0)
	if err != nil && n == 0 {
		t.Fatalf("failed to read lock file: %v", err)
	}
	pidStr := strings.TrimSpace(string(buf[:n]))
	pid, err := strconv.Atoi(pidStr)
	if err != nil || pid != os.Getpid() {
		t.Errorf("expected PID %d in lock file, got %q", os.Getpid(), pidStr)
	}

	// Attempting to acquire a second lock on the same path must fail
	lock2, err := AcquireInstanceLock(lockPath)
	if err == nil {
		if lock2 != nil {
			lock2.Release()
		}
		t.Fatalf("expected lock2 to fail while lock1 is held, but succeeded")
	}

	// Releasing lock1 allows lock3 to be acquired
	lock1.Release()
	lock1 = nil

	lock3, err := AcquireInstanceLock(lockPath)
	if err != nil {
		t.Fatalf("expected lock3 to succeed after lock1 was released, got error: %v", err)
	}
	lock3.Release()
}

func TestInstanceLockNilAndDoubleRelease(t *testing.T) {
	var nilLock *InstanceLock
	// Should not panic on nil
	nilLock.Release()

	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	lockPath := filepath.Join(testDir.Path, "sam.lock")
	lock, err := AcquireInstanceLock(lockPath)
	if err != nil {
		t.Fatalf("failed to acquire lock: %v", err)
	}

	// First release
	lock.Release()
	// Second release should be a no-op and not panic
	lock.Release()
}

func TestGetLockFilePath(t *testing.T) {
	// Custom config path
	customPath := filepath.Join(string(filepath.Separator)+"custom", "dir", "config.yaml")
	customCfg := &config.Config{
		Path: customPath,
	}
	expected := filepath.Join(string(filepath.Separator)+"custom", "dir", "sam.lock")
	if got := getLockFilePath(customCfg); got != expected {
		t.Errorf("expected %s, got %s", expected, got)
	}

	// Nil config fallback
	fallbackPath := getLockFilePath(nil)
	if fallbackPath == "" {
		t.Errorf("expected non-empty fallback lock file path")
	}
}
