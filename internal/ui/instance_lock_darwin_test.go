//go:build darwin

package ui

import (
	"path/filepath"
	"testing"

	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/test"
)

func TestInitAbortsWhenAlreadyLocked(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	cfgPath := filepath.Join(testDir.Path, "config.yaml")
	t.Setenv("SOPS_AGE_MANAGER_CONFIG_DIR", cfgPath)

	cfg := config.NewConfig("key1", "key1", testDir.Path)
	cfg.Path = cfgPath
	if err := cfg.Write(); err != nil {
		t.Fatalf("could not write test config: %v", err)
	}

	lockPath := getLockFilePath(cfg)
	existingLock, err := AcquireInstanceLock(lockPath)
	if err != nil {
		t.Fatalf("failed to acquire test lock: %v", err)
	}
	defer existingLock.Release()

	// Init should detect existing lock and return immediately without blocking or panicking
	Init(cfg, nil)
}
