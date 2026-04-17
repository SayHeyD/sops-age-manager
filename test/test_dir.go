package test

import (
	"os"
	"testing"
)

type Dir struct {
	Path string
}

// GenerateNewUniqueTestDir creates and returns an empty directory for testing whatever you'd like
func GenerateNewUniqueTestDir(t *testing.T) *Dir {
	t.Helper()
	testDir := getTestBaseDir(t)

	return &Dir{
		Path: testDir,
	}
}

// CleanTestDir removes the test directory and all contents of it.
func (d *Dir) CleanTestDir(t *testing.T) {
	t.Helper()
	err := os.RemoveAll(d.Path)
	if err != nil {
		t.Fatalf("Could not delete test directory \"%s\": %v", d.Path, err)
	}
}
