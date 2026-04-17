package test

import (
	"os"
	"testing"
)

func TestIfGenerateNewUniqueTestDirDoesNotReturnNil(t *testing.T) {
	dir := GenerateNewUniqueTestDir(t)

	if dir == nil {
		t.Fatalf("GenerateNewUniqueTestDir() does not return an object")
	}

	err := os.RemoveAll(dir.Path)
	if err != nil {
		t.Fatalf("Could not delete test directory \"%s\": %v", dir.Path, err)
	}
}

func TestIfGenerateNewUniqueTestDirCreatesATestDirectory(t *testing.T) {
	dir := GenerateNewUniqueTestDir(t)

	if _, err := os.Stat(dir.Path); os.IsNotExist(err) {
		t.Fatalf("Test directory was not created on hard drive: \"%s\"", dir.Path)
	}

	err := os.RemoveAll(dir.Path)
	if err != nil {
		t.Fatalf("Could not delete test directory \"%s\": %v", dir.Path, err)
	}
}

func TestIfCleanTestDirDeletesTheTestDirectory(t *testing.T) {
	dir := GenerateNewUniqueTestDir(t)

	dir.CleanTestDir(t)
	if _, err := os.Stat(dir.Path); os.IsNotExist(err) {
		return
	}

	t.Fatalf("Directory still exists: \"%s\"", dir.Path)
}
