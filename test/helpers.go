package test

import (
	"github.com/google/uuid"
	"os"
	"testing"
)

// getTestBaseDir returns the base path for creating test directories and
// creates the base directory if it does not exist
func getTestBaseDir(t *testing.T) string {
	dir := os.TempDir() + "sops-age-manager"

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			t.Fatalf("Could not create testing directories: %v", err)
		}
	} else if err != nil {
		t.Fatalf("Could not check if testing directory exist: %v", err)
	}

	return dir
}

// GenerateNewUniqueTestDir creates and returns an empty directory for testing whatever you'd like
func GenerateNewUniqueTestDir(t *testing.T) string {
	testDir := getTestBaseDir(t) + string(os.PathSeparator) + uuid.NewString()

	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		if err := os.Mkdir(testDir, os.ModePerm); err != nil {
			t.Fatalf("Could not create testing directories: %v", err)
		}
	} else if err != nil {
		t.Fatalf("Could not check if testing directory exist: %v", err)
	}

	return testDir
}
