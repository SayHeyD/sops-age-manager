package test

import (
	"os"
	"testing"
)

// getTestBaseDir returns the base path for creating test directories and
// creates the base directory if it does not exist
func getTestBaseDir(t *testing.T) string {
	t.Helper()
	tmpDir := os.TempDir()

	if tmpDir[len(tmpDir)-1:] != string(os.PathSeparator) {
		tmpDir += string(os.PathSeparator)
	}

	dir := tmpDir + "sops-age-manager"

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			t.Fatalf("Could not create testing directories: %v", err)
		}
	} else if err != nil {
		t.Fatalf("Could not check if testing directory exist: %v", err)
	}

	return dir
}
