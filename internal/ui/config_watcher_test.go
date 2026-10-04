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

func TestConfigWatcherDetectsFileChanges(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	configFilePath := filepath.Join(testDir.Path, "config.yaml")
	if err := os.WriteFile(configFilePath, []byte("initial: true\n"), 0600); err != nil {
		t.Fatalf("could not write initial config: %v", err)
	}

	var changeCount int32
	cw, err := NewConfigWatcher(configFilePath, func() {
		atomic.AddInt32(&changeCount, 1)
	})
	if err != nil {
		t.Fatalf("failed to create config watcher: %v", err)
	}
	defer cw.Stop()

	// Give watcher a moment to initialize
	time.Sleep(50 * time.Millisecond)

	// Modify config file
	if err := os.WriteFile(configFilePath, []byte("encryptionKey: key-1\n"), 0600); err != nil {
		t.Fatalf("could not update config file: %v", err)
	}

	// Wait for debounce and callback
	time.Sleep(200 * time.Millisecond)

	if count := atomic.LoadInt32(&changeCount); count < 1 {
		t.Errorf("expected at least 1 change notification, got %d", count)
	}

	// Test Stop idempotency
	cw.Stop()
	cw.Stop()
}

func TestConfigWatcherInvalidPath(t *testing.T) {
	_, err := NewConfigWatcher("/nonexistent_directory/config.yaml", func() {})
	if err == nil {
		t.Error("expected error for non-existent directory")
	}
}
