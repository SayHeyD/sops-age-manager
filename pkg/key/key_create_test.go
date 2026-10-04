package key

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/SayHeyD/sops-age-manager/test"
)

const mockValidKeyFile = `# created: 2026-10-05T01:00:00Z
# public key: age1z9zvlcr2j3gt7mc9flmvyxm264v5aqyq0u2l46rlkg2c2fdzytgx7xl3qm
AGE-SECRET-KEY-1HHS36XWKCVDKEKJ2M7WKQN3MFYUGIP4WWM7DT1CFANZUT5LT3K8ZRFZFGV3
`

func mockKeygenRunner(t *testing.T, content string) func(string, string) ([]byte, error) {
	return func(keygenPath, outputPath string) ([]byte, error) {
		if err := os.WriteFile(outputPath, []byte(content), 0600); err != nil {
			return nil, err
		}
		return []byte("Public key: age1z9zvlcr2j3gt7mc9flmvyxm264v5aqyq0u2l46rlkg2c2fdzytgx7xl3qm\n"), nil
	}
}

func TestCreateKeySuccess(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	restore := SetKeygenRunner(
		func(file string) (string, error) {
			return "/usr/bin/age-keygen", nil
		},
		mockKeygenRunner(t, mockValidKeyFile),
	)
	defer restore()

	k, err := CreateKey("test-key", testDir.Path, false)
	if err != nil {
		t.Fatalf("unexpected error creating key: %v", err)
	}

	if k.Name != "test-key" {
		t.Errorf("expected key name %q, got %q", "test-key", k.Name)
	}

	expectedPath := filepath.Join(testDir.Path, "test-key.txt")
	if k.FileName != expectedPath {
		t.Errorf("expected file path %q, got %q", expectedPath, k.FileName)
	}

	if k.PublicKey != "age1z9zvlcr2j3gt7mc9flmvyxm264v5aqyq0u2l46rlkg2c2fdzytgx7xl3qm" {
		t.Errorf("unexpected public key: %q", k.PublicKey)
	}

	if k.PrivateKey != "AGE-SECRET-KEY-1HHS36XWKCVDKEKJ2M7WKQN3MFYUGIP4WWM7DT1CFANZUT5LT3K8ZRFZFGV3" {
		t.Errorf("unexpected private key: %q", k.PrivateKey)
	}

	// Verify file permissions
	info, err := os.Stat(expectedPath)
	if err != nil {
		t.Fatalf("failed to stat created key file: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected mode 0600, got %v", info.Mode().Perm())
	}
}

func TestCreateKeyNestedDirectory(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	restore := SetKeygenRunner(
		func(file string) (string, error) {
			return "/usr/bin/age-keygen", nil
		},
		mockKeygenRunner(t, mockValidKeyFile),
	)
	defer restore()

	k, err := CreateKey("nested/sub/cluster-key.txt", testDir.Path, false)
	if err != nil {
		t.Fatalf("unexpected error creating nested key: %v", err)
	}

	expectedName := filepath.Join("nested", "sub", "cluster-key")
	if k.Name != "nested/sub/cluster-key" && k.Name != expectedName {
		t.Errorf("expected key name nested/sub/cluster-key, got %q", k.Name)
	}

	expectedPath := filepath.Join(testDir.Path, "nested", "sub", "cluster-key.txt")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("expected file %q to exist", expectedPath)
	}
}

func TestCreateKeyAlreadyExistsWithoutForce(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	existingPath := filepath.Join(testDir.Path, "existing.txt")
	if err := os.WriteFile(existingPath, []byte(mockValidKeyFile), 0600); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	restore := SetKeygenRunner(
		func(file string) (string, error) {
			return "/usr/bin/age-keygen", nil
		},
		mockKeygenRunner(t, mockValidKeyFile),
	)
	defer restore()

	_, err := CreateKey("existing", testDir.Path, false)
	if err == nil {
		t.Fatal("expected error when key file already exists without force, got nil")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestCreateKeyAlreadyExistsWithForce(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	existingPath := filepath.Join(testDir.Path, "existing.txt")
	if err := os.WriteFile(existingPath, []byte("old content"), 0600); err != nil {
		t.Fatalf("failed to create existing file: %v", err)
	}

	restore := SetKeygenRunner(
		func(file string) (string, error) {
			return "/usr/bin/age-keygen", nil
		},
		mockKeygenRunner(t, mockValidKeyFile),
	)
	defer restore()

	k, err := CreateKey("existing", testDir.Path, true)
	if err != nil {
		t.Fatalf("unexpected error when overwriting with force: %v", err)
	}
	if k == nil || k.PublicKey == "" {
		t.Fatalf("expected valid key returned, got %v", k)
	}
}

func TestCreateKeyEmptyName(t *testing.T) {
	_, err := CreateKey("   ", "", false)
	if err == nil {
		t.Fatal("expected error for empty key name, got nil")
	}
}

func TestCreateKeyMissingAgeKeygen(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	restore := SetKeygenRunner(
		func(file string) (string, error) {
			return "", exec.ErrNotFound
		},
		nil,
	)
	defer restore()

	_, err := CreateKey("my-key", testDir.Path, false)
	if err == nil {
		t.Fatal("expected error when age-keygen not found, got nil")
	}
	if !strings.Contains(err.Error(), "age-keygen not found in PATH") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestCreateKeyCommandFails(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	restore := SetKeygenRunner(
		func(file string) (string, error) {
			return "/usr/bin/age-keygen", nil
		},
		func(keygenPath, outputPath string) ([]byte, error) {
			return []byte("permission denied"), errors.New("exit status 1")
		},
	)
	defer restore()

	_, err := CreateKey("my-key", testDir.Path, false)
	if err == nil {
		t.Fatal("expected error when keygen fails, got nil")
	}
	if !strings.Contains(err.Error(), "failed to generate age key") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestCreateKeyInvalidKeyOutput(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	restore := SetKeygenRunner(
		func(file string) (string, error) {
			return "/usr/bin/age-keygen", nil
		},
		mockKeygenRunner(t, "invalid content with no keys"),
	)
	defer restore()

	_, err := CreateKey("my-key", testDir.Path, false)
	if err == nil {
		t.Fatal("expected error when generated key is invalid, got nil")
	}
	if !strings.Contains(err.Error(), "does not contain valid age keys") {
		t.Errorf("unexpected error message: %v", err)
	}
}
