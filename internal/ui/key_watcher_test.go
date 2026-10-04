//go:build darwin || windows

package ui

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/SayHeyD/sops-age-manager/test"
)

func TestKeyWatcherFileEvents(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	var changeCount int32
	kw, err := NewKeyWatcher(testDir.Path, func() {
		atomic.AddInt32(&changeCount, 1)
	})
	if err != nil {
		t.Fatalf("failed to create KeyWatcher: %v", err)
	}
	defer kw.Stop()

	// 1. Create a key file
	key1Path := filepath.Join(testDir.Path, "key-1.txt")
	content := "# public key: age1samplekey1\nAGE-SECRET-KEY-1SAMPLEKEY1\n"
	if err := os.WriteFile(key1Path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write key file: %v", err)
	}

	// Wait for debounce
	time.Sleep(150 * time.Millisecond)
	if atomic.LoadInt32(&changeCount) < 1 {
		t.Errorf("expected at least 1 change event on file create, got %d", atomic.LoadInt32(&changeCount))
	}

	// Reset count
	atomic.StoreInt32(&changeCount, 0)

	// 2. Modify the key file
	if err := os.WriteFile(key1Path, []byte(content+"# update\n"), 0600); err != nil {
		t.Fatalf("failed to modify key file: %v", err)
	}

	time.Sleep(150 * time.Millisecond)
	if atomic.LoadInt32(&changeCount) < 1 {
		t.Errorf("expected at least 1 change event on file modify, got %d", atomic.LoadInt32(&changeCount))
	}

	// Reset count
	atomic.StoreInt32(&changeCount, 0)

	// 3. Remove the key file
	if err := os.Remove(key1Path); err != nil {
		t.Fatalf("failed to remove key file: %v", err)
	}

	time.Sleep(150 * time.Millisecond)
	if atomic.LoadInt32(&changeCount) < 1 {
		t.Errorf("expected at least 1 change event on file remove, got %d", atomic.LoadInt32(&changeCount))
	}
}

func TestKeyWatcherSubdirectoryFileEvents(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	var changeCount int32
	kw, err := NewKeyWatcher(testDir.Path, func() {
		atomic.AddInt32(&changeCount, 1)
	})
	if err != nil {
		t.Fatalf("failed to create KeyWatcher: %v", err)
	}
	defer kw.Stop()

	// Create subdirectory
	subDir := filepath.Join(testDir.Path, "subfolder")
	if err := os.MkdirAll(subDir, 0700); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	atomic.StoreInt32(&changeCount, 0)

	// Create key file in subdirectory
	subKeyPath := filepath.Join(subDir, "nested-key.txt")
	content := "# public key: age1samplekeynested\nAGE-SECRET-KEY-1SAMPLENESTED\n"
	if err := os.WriteFile(subKeyPath, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write nested key file: %v", err)
	}

	time.Sleep(150 * time.Millisecond)
	if atomic.LoadInt32(&changeCount) < 1 {
		t.Errorf("expected change event for nested file creation, got %d", atomic.LoadInt32(&changeCount))
	}
}

func TestKeyWatcherStopCleanly(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	var changeCount int32
	kw, err := NewKeyWatcher(testDir.Path, func() {
		atomic.AddInt32(&changeCount, 1)
	})
	if err != nil {
		t.Fatalf("failed to create KeyWatcher: %v", err)
	}

	kw.Stop()
	// Multiple stop calls should be safe
	kw.Stop()

	// Write file after stop
	keyPath := filepath.Join(testDir.Path, "key-after-stop.txt")
	_ = os.WriteFile(keyPath, []byte("dummy"), 0600)

	time.Sleep(100 * time.Millisecond)
	if atomic.LoadInt32(&changeCount) != 0 {
		t.Errorf("expected 0 change events after stop, got %d", atomic.LoadInt32(&changeCount))
	}
}
