package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/test"
)

func setupTestClearEnv(t *testing.T, initialEnc, initialDec string) (*test.Dir, func()) {
	t.Helper()
	testDir := test.GenerateNewUniqueTestDir(t)

	configPath := filepath.Join(testDir.Path, "config.yaml")
	t.Setenv("SOPS_AGE_MANAGER_CONFIG_DIR", configPath)

	cfg := config.NewConfig(initialEnc, initialDec, testDir.Path)
	if err := cfg.Write(); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	return testDir, func() {
		testDir.CleanTestDir(t)
	}
}

func TestClearActiveKeysDefaultClearsBoth(t *testing.T) {
	_, cleanup := setupTestClearEnv(t, "enc-key", "dec-key")
	defer cleanup()

	clearEncryptionKey = false
	clearDecryptionKey = false

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	clearActiveKeys()

	_ = w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	cfg, err := config.NewConfigFromFile()
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	if cfg.EncryptionKeyName != "" || cfg.DecryptionKeyName != "" {
		t.Errorf("expected both keys to be empty, got enc=%q, dec=%q", cfg.EncryptionKeyName, cfg.DecryptionKeyName)
	}

	if !strings.Contains(output, "Cleared active decryption key") || !strings.Contains(output, "Cleared active encryption key") {
		t.Errorf("unexpected output: %q", output)
	}
}

func TestClearActiveKeysEncryptionOnly(t *testing.T) {
	_, cleanup := setupTestClearEnv(t, "enc-key", "dec-key")
	defer cleanup()

	clearEncryptionKey = true
	clearDecryptionKey = false
	defer func() { clearEncryptionKey = false }()

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	clearActiveKeys()

	_ = w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	cfg, err := config.NewConfigFromFile()
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	if cfg.EncryptionKeyName != "" {
		t.Errorf("expected encryption key to be empty, got %q", cfg.EncryptionKeyName)
	}
	if cfg.DecryptionKeyName != "dec-key" {
		t.Errorf("expected decryption key to remain %q, got %q", "dec-key", cfg.DecryptionKeyName)
	}

	if !strings.Contains(output, "Cleared active encryption key") {
		t.Errorf("unexpected output: %q", output)
	}
	if strings.Contains(output, "Cleared active decryption key") {
		t.Errorf("output should not contain decryption clear message: %q", output)
	}
}

func TestClearActiveKeysDecryptionOnly(t *testing.T) {
	_, cleanup := setupTestClearEnv(t, "enc-key", "dec-key")
	defer cleanup()

	clearEncryptionKey = false
	clearDecryptionKey = true
	defer func() { clearDecryptionKey = false }()

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	clearActiveKeys()

	_ = w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	cfg, err := config.NewConfigFromFile()
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	if cfg.EncryptionKeyName != "enc-key" {
		t.Errorf("expected encryption key to remain %q, got %q", "enc-key", cfg.EncryptionKeyName)
	}
	if cfg.DecryptionKeyName != "" {
		t.Errorf("expected decryption key to be empty, got %q", cfg.DecryptionKeyName)
	}

	if !strings.Contains(output, "Cleared active decryption key") {
		t.Errorf("unexpected output: %q", output)
	}
	if strings.Contains(output, "Cleared active encryption key") {
		t.Errorf("output should not contain encryption clear message: %q", output)
	}
}

func TestClearKeyCommandExecution(t *testing.T) {
	_, cleanup := setupTestClearEnv(t, "enc-key", "dec-key")
	defer cleanup()

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	RootCmd.SetArgs([]string{"key", "clear"})
	err = RootCmd.Execute()

	_ = w.Close()
	os.Stdout = origStdout

	if err != nil {
		t.Fatalf("unexpected error executing RootCmd: %v", err)
	}

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	cfg, err := config.NewConfigFromFile()
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	if cfg.EncryptionKeyName != "" || cfg.DecryptionKeyName != "" {
		t.Errorf("expected keys to be cleared, got enc=%q, dec=%q", cfg.EncryptionKeyName, cfg.DecryptionKeyName)
	}
}
