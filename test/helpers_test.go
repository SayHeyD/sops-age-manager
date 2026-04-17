package test

import (
	"os"
	"testing"
)

func TestGetTestBaseDirReturnsANonEmptyString(t *testing.T) {
	dirPath := getTestBaseDir(t)
	defer os.RemoveAll(dirPath)

	if dirPath == "" {
		t.Fatalf("getTestBaseDir() does not return a string")
	}
}

func TestGetTestBaseDirReturnsCreatesADirectory(t *testing.T) {
	dirPath := getTestBaseDir(t)
	defer os.RemoveAll(dirPath)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Fatalf("Directory was not created: %s", dirPath)
	}
}
