package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
	"github.com/SayHeyD/sops-age-manager/test"
)

func mockKeygenInCmd(t *testing.T) func() {
	return key.SetKeygenRunner(
		func(file string) (string, error) {
			return "/usr/bin/age-keygen", nil
		},
		func(keygenPath, outputPath string) ([]byte, error) {
			if err := os.WriteFile(outputPath, []byte(testKeyContent), 0600); err != nil {
				return nil, err
			}
			return []byte("Public key: " + testKeyPublicKey + "\n"), nil
		},
	)
}

func TestCreateKeyCommandDefault(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	configPath := filepath.Join(testDir.Path, "config.yaml")
	t.Setenv("SOPS_AGE_MANAGER_CONFIG_DIR", configPath)

	cfg := config.NewConfig("", "", testDir.Path)
	if err := cfg.Write(); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	restoreKeygen := mockKeygenInCmd(t)
	defer restoreKeygen()

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	createForce = false
	createCopyPublic = false
	createSetEncryption = false
	createSetDecryption = false
	createSetBoth = false

	createKey("new-key")

	_ = w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "Created key \"new-key\"") {
		t.Errorf("expected created key message, got %q", output)
	}
	if !strings.Contains(output, "Public key: "+testKeyPublicKey) {
		t.Errorf("expected public key in output, got %q", output)
	}

	// Verify key file was created on disk
	createdFilePath := filepath.Join(testDir.Path, "new-key.txt")
	if _, err := os.Stat(createdFilePath); os.IsNotExist(err) {
		t.Fatalf("expected key file %q to exist", createdFilePath)
	}
}

func TestCreateKeyCommandWithCopyFlag(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	configPath := filepath.Join(testDir.Path, "config.yaml")
	t.Setenv("SOPS_AGE_MANAGER_CONFIG_DIR", configPath)

	cfg := config.NewConfig("", "", testDir.Path)
	if err := cfg.Write(); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	restoreKeygen := mockKeygenInCmd(t)
	defer restoreKeygen()

	var copiedText string
	origClipboard := clipboardWriteAllFunc
	defer func() { clipboardWriteAllFunc = origClipboard }()

	clipboardWriteAllFunc = func(text string) error {
		copiedText = text
		return nil
	}

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	createForce = false
	createCopyPublic = true
	createSetEncryption = false
	createSetDecryption = false
	createSetBoth = false
	defer func() { createCopyPublic = false }()

	createKey("copy-key")

	_ = w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if copiedText != testKeyPublicKey {
		t.Errorf("expected clipboard text %q, got %q", testKeyPublicKey, copiedText)
	}
	if !strings.Contains(output, "Copied public key of \"copy-key\" to clipboard") {
		t.Errorf("expected clipboard confirmation in output, got %q", output)
	}
}

func TestCreateKeyCommandWithUseFlag(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	configPath := filepath.Join(testDir.Path, "config.yaml")
	t.Setenv("SOPS_AGE_MANAGER_CONFIG_DIR", configPath)

	cfg := config.NewConfig("", "", testDir.Path)
	if err := cfg.Write(); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	restoreKeygen := mockKeygenInCmd(t)
	defer restoreKeygen()

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	createForce = false
	createCopyPublic = false
	createSetEncryption = false
	createSetDecryption = false
	createSetBoth = true
	defer func() { createSetBoth = false }()

	createKey("active-key")

	_ = w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "Set \"active-key\" as active decryption key") {
		t.Errorf("expected decryption message in output, got %q", output)
	}
	if !strings.Contains(output, "Set \"active-key\" as active encryption key") {
		t.Errorf("expected encryption message in output, got %q", output)
	}

	// Verify config
	updatedCfg, err := config.NewConfigFromFile()
	if err != nil {
		t.Fatalf("failed to read updated config: %v", err)
	}
	if updatedCfg.EncryptionKeyName != "active-key" || updatedCfg.DecryptionKeyName != "active-key" {
		t.Errorf("expected both keys set to active-key, got enc=%q dec=%q", updatedCfg.EncryptionKeyName, updatedCfg.DecryptionKeyName)
	}
}

func TestCreateKeyCommandWithEncryptionAndDecryptionFlags(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	configPath := filepath.Join(testDir.Path, "config.yaml")
	t.Setenv("SOPS_AGE_MANAGER_CONFIG_DIR", configPath)

	cfg := config.NewConfig("", "", testDir.Path)
	if err := cfg.Write(); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	restoreKeygen := mockKeygenInCmd(t)
	defer restoreKeygen()

	// 1. Encryption only
	createForce = false
	createCopyPublic = false
	createSetEncryption = true
	createSetDecryption = false
	createSetBoth = false

	createKey("enc-only")
	createSetEncryption = false

	updatedCfg, err := config.NewConfigFromFile()
	if err != nil {
		t.Fatalf("failed to read updated config: %v", err)
	}
	if updatedCfg.EncryptionKeyName != "enc-only" || updatedCfg.DecryptionKeyName != "" {
		t.Errorf("expected enc-only set for encryption only, got enc=%q dec=%q", updatedCfg.EncryptionKeyName, updatedCfg.DecryptionKeyName)
	}

	// 2. Decryption only
	createSetDecryption = true
	createKey("dec-only")
	createSetDecryption = false

	updatedCfg, err = config.NewConfigFromFile()
	if err != nil {
		t.Fatalf("failed to read updated config: %v", err)
	}
	if updatedCfg.DecryptionKeyName != "dec-only" || updatedCfg.EncryptionKeyName != "enc-only" {
		t.Errorf("expected dec-only set for decryption only, got enc=%q dec=%q", updatedCfg.EncryptionKeyName, updatedCfg.DecryptionKeyName)
	}
}

func TestCreateKeyCommandExecutionViaRoot(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	configPath := filepath.Join(testDir.Path, "config.yaml")
	t.Setenv("SOPS_AGE_MANAGER_CONFIG_DIR", configPath)

	cfg := config.NewConfig("", "", testDir.Path)
	if err := cfg.Write(); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	restoreKeygen := mockKeygenInCmd(t)
	defer restoreKeygen()

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	RootCmd.SetArgs([]string{"key", "create", "root-created-key"})
	err = RootCmd.Execute()

	_ = w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if err != nil {
		t.Fatalf("unexpected error executing RootCmd key create: %v", err)
	}

	if !strings.Contains(output, "Created key \"root-created-key\"") {
		t.Errorf("expected created key message in output, got %q", output)
	}
}
