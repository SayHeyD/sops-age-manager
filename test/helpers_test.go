package test

import (
	"os"
	"testing"
)

func TestGetTestBaseDirReturnsANonEmptyString(t *testing.T) {
	dirPath := getTestBaseDir(t)

	if dirPath == "" {
		t.Fatalf("GenerateNewUniqueTestDir() does not return an object")
	}

	err := os.RemoveAll(dirPath)
	if err != nil {
		t.Fatalf("Could not delete test directory \"%s\": %v", dirPath, err)
	}
}

func TestGetTestBaseDirReturnsCreatesADirectory(t *testing.T) {
	dirPath := getTestBaseDir(t)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Fatalf("Directory was not created: %s", dirPath)
	}

	err := os.RemoveAll(dirPath)
	if err != nil {
		t.Fatalf("Could not delete test directory \"%s\": %v", dirPath, err)
	}
}
