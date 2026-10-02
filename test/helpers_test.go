package test

import (
	"os"
	"testing"
)

func TestGetTestBaseDirReturnsANonEmptyString(t *testing.T) {
	dirPath := getTestBaseDir(t)

	if dirPath == "" {
		t.Fatalf("getTestBaseDir() does not return a non-empty string")
	}
}

func TestGetTestBaseDirReturnsCreatesADirectory(t *testing.T) {
	dirPath := getTestBaseDir(t)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Fatalf("Directory was not created: %s", dirPath)
	}
}
