package test

import (
	"os"
	"testing"
)

// getTestBaseDir returns the base path for creating test directories and
// creates the base directory if it does not exist
func getTestBaseDir(t *testing.T) string {
	t.Helper()

	dir, err := os.MkdirTemp("", "sops-age-manager-")
	if err != nil {
		t.Fatalf("Could not create testing directory: %v", err)
	}

	return dir
}
