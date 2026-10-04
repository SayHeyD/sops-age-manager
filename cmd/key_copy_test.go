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

const (
	testKeyContent = `# created: 2023-01-19T18:37:24+01:00
# public key: age1z9zvlcr2j3gt7mc9flmvyxm264v5aqyq0u2l46rlkg2c2fdzytgx7xl3qm
AGE-SECRET-KEY-HHS36XWKCVDKEKJ2M7WKQN3MFYUGIP4WWM7DT1CFANZUT5LT3K8ZRFZFGV3`
	testKeyPublicKey  = "age1z9zvlcr2j3gt7mc9flmvyxm264v5aqyq0u2l46rlkg2c2fdzytgx7xl3qm"
	testKeyPrivateKey = "AGE-SECRET-KEY-HHS36XWKCVDKEKJ2M7WKQN3MFYUGIP4WWM7DT1CFANZUT5LT3K8ZRFZFGV3"
)

func setupTestKeyEnv(t *testing.T) (*test.Dir, func()) {
	t.Helper()
	testDir := test.GenerateNewUniqueTestDir(t)

	keyFilePath := filepath.Join(testDir.Path, "test-key.txt")
	if err := os.WriteFile(keyFilePath, []byte(testKeyContent), 0600); err != nil {
		t.Fatalf("failed to write test key file: %v", err)
	}

	configPath := filepath.Join(testDir.Path, "config.yaml")
	t.Setenv("SOPS_AGE_MANAGER_CONFIG_DIR", configPath)

	cfg := config.NewConfig("", "", testDir.Path)
	if err := cfg.Write(); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	return testDir, func() {
		testDir.CleanTestDir(t)
	}
}

func TestCopyKeyDefaultCopiesPublicKey(t *testing.T) {
	_, cleanup := setupTestKeyEnv(t)
	defer cleanup()

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

	copyPublicKey = false
	copyPrivateKey = false
	copyKeyName = false

	copyKey("test-key")

	_ = w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if copiedText != testKeyPublicKey {
		t.Errorf("expected copied text to be %q, got %q", testKeyPublicKey, copiedText)
	}

	if !strings.Contains(output, "Copied public key of \"test-key\" to clipboard") {
		t.Errorf("unexpected output: %q", output)
	}
}

func TestCopyKeyWithPrivateKeyFlag(t *testing.T) {
	_, cleanup := setupTestKeyEnv(t)
	defer cleanup()

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

	copyPublicKey = false
	copyPrivateKey = true
	copyKeyName = false
	defer func() { copyPrivateKey = false }()

	copyKey("test-key")

	_ = w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if copiedText != testKeyPrivateKey {
		t.Errorf("expected copied text to be %q, got %q", testKeyPrivateKey, copiedText)
	}

	if !strings.Contains(output, "Copied private key of \"test-key\" to clipboard") {
		t.Errorf("unexpected output: %q", output)
	}
}

func TestCopyKeyWithNameFlag(t *testing.T) {
	_, cleanup := setupTestKeyEnv(t)
	defer cleanup()

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

	copyPublicKey = false
	copyPrivateKey = false
	copyKeyName = true
	defer func() { copyKeyName = false }()

	copyKey("test-key")

	_ = w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if copiedText != "test-key" {
		t.Errorf("expected copied text to be %q, got %q", "test-key", copiedText)
	}

	if !strings.Contains(output, "Copied key name of \"test-key\" to clipboard") {
		t.Errorf("unexpected output: %q", output)
	}
}

func TestCopyKeyCommandExecution(t *testing.T) {
	_, cleanup := setupTestKeyEnv(t)
	defer cleanup()

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

	RootCmd.SetArgs([]string{"key", "copy", "test-key"})
	err = RootCmd.Execute()

	_ = w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	if err != nil {
		t.Fatalf("unexpected error executing RootCmd: %v", err)
	}

	if copiedText != testKeyPublicKey {
		t.Errorf("expected copied text to be %q, got %q", testKeyPublicKey, copiedText)
	}
}
