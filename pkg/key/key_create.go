package key

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	lookPath            = exec.LookPath
	keygenCommandRunner = func(keygenPath string, outputPath string) ([]byte, error) {
		cmd := exec.Command(keygenPath, "-o", outputPath)
		return cmd.CombinedOutput()
	}
)

// SetKeygenRunner sets custom lookPath and command runner functions for testing and returns a restore function.
func SetKeygenRunner(customLookPath func(string) (string, error), customRunner func(string, string) ([]byte, error)) func() {
	origLookPath := lookPath
	origRunner := keygenCommandRunner
	if customLookPath != nil {
		lookPath = customLookPath
	}
	if customRunner != nil {
		keygenCommandRunner = customRunner
	}
	return func() {
		lookPath = origLookPath
		keygenCommandRunner = origRunner
	}
}

// CreateKey creates a new age key with the specified name in the configured key directory.
// It invokes the pre-installed 'age-keygen' binary.
// If force is false and the key file already exists, an error is returned.
func CreateKey(keyName string, keyDirPath string, force bool) (*Key, error) {
	cleanName := strings.TrimSpace(keyName)
	if cleanName == "" {
		return nil, fmt.Errorf("key name cannot be empty")
	}

	cleanName = strings.TrimSuffix(cleanName, ".txt")

	keyDirPath, err := GetKeyDirPath(keyDirPath)
	if err != nil {
		return nil, err
	}

	targetFilePath := filepath.Join(keyDirPath, cleanName+".txt")

	if !force {
		if _, err := os.Stat(targetFilePath); err == nil {
			return nil, fmt.Errorf("key file already exists at %q: use --force to overwrite", targetFilePath)
		}
	}

	targetDir := filepath.Dir(targetFilePath)
	if err := os.MkdirAll(targetDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create directory %q: %w", targetDir, err)
	}

	keygenPath, err := lookPath("age-keygen")
	if err != nil {
		return nil, fmt.Errorf("age-keygen not found in PATH: please install age (e.g. brew install age, apt install age, winget install FiloSottile.age)")
	}

	output, err := keygenCommandRunner(keygenPath, targetFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to generate age key: %s (%w)", strings.TrimSpace(string(output)), err)
	}

	if err := os.Chmod(targetFilePath, 0600); err != nil {
		return nil, fmt.Errorf("failed to set key file permissions: %w", err)
	}

	fileContent, err := os.ReadFile(targetFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read created key file: %w", err)
	}

	createdKey := NewKey(cleanName, targetFilePath, string(fileContent))
	if createdKey.PublicKey == "" || createdKey.PrivateKey == "" {
		return nil, fmt.Errorf("generated key file does not contain valid age keys")
	}

	return createdKey, nil
}
