//go:build darwin || windows

package ui

import (
	"path/filepath"
	"testing"

	"github.com/SayHeyD/sops-age-manager/test"
)

func TestOpenDirectoryCreatesDirectoryIfMissing(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	targetDir := filepath.Join(testDir.Path, "nested", "deep", "dir")

	var openedPath string
	origOpenCommand := openCommandFunc
	defer func() { openCommandFunc = origOpenCommand }()

	openCommandFunc = func(path string) error {
		openedPath = path
		return nil
	}

	if err := OpenDirectory(targetDir); err != nil {
		t.Fatalf("unexpected error calling OpenDirectory: %v", err)
	}

	if openedPath != targetDir {
		t.Errorf("expected openedPath to be %q, got %q", targetDir, openedPath)
	}
}

func TestOpenConfigDirectory(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	configFilePath := filepath.Join(testDir.Path, "config.yaml")
	t.Setenv("SOPS_AGE_MANAGER_CONFIG_DIR", configFilePath)

	var openedPath string
	origOpenCommand := openCommandFunc
	defer func() { openCommandFunc = origOpenCommand }()

	openCommandFunc = func(path string) error {
		openedPath = path
		return nil
	}

	if err := OpenConfigDirectory(); err != nil {
		t.Fatalf("unexpected error calling OpenConfigDirectory: %v", err)
	}

	if openedPath != testDir.Path {
		t.Errorf("expected openedPath to be %q, got %q", testDir.Path, openedPath)
	}
}

func TestOpenKeyDirectory(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	var openedPath string
	origOpenCommand := openCommandFunc
	defer func() { openCommandFunc = origOpenCommand }()

	openCommandFunc = func(path string) error {
		openedPath = path
		return nil
	}

	if err := OpenKeyDirectory(testDir.Path); err != nil {
		t.Fatalf("unexpected error calling OpenKeyDirectory: %v", err)
	}

	if openedPath != testDir.Path {
		t.Errorf("expected openedPath to be %q, got %q", testDir.Path, openedPath)
	}
}
