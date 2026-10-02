package key

import (
	"path/filepath"
	"testing"

	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/test"
)

const wantedKeyName = "test_key"

const wantedFilePath = "/var/someDir/tmp/test_key.txt"

// ageKeyFileContent is NOT a real age key pair
const ageKeyFileContent = `# created: 2023-01-19T18:37:24+01:00
# public key: age1z9zvlcr2j3gt7mc9flmvyxm264v5aqyq0u2l46rlkg2c2fdzytgx7xl3qm
AGE-SECRET-KEY-HHS36XWKCVDKEKJ2M7WKQN3MFYUGIP4WWM7DT1CFANZUT5LT3K8ZRFZFGV3`

// ageKeyFileContentWithNewLine is NOT a real age key pair
const ageKeyFileContentWithNewLine = `# created: 2023-01-19T18:37:24+01:00
# public key: age1z9zvlcr2j3gt7mc9flmvyxm264v5aqyq0u2l46rlkg2c2fdzytgx7xl3qm
AGE-SECRET-KEY-HHS36XWKCVDKEKJ2M7WKQN3MFYUGIP4WWM7DT1CFANZUT5LT3K8ZRFZFGV3
`

// ageKeyPublicKey is NOT a real age public key
const ageKeyPublicKey = "age1z9zvlcr2j3gt7mc9flmvyxm264v5aqyq0u2l46rlkg2c2fdzytgx7xl3qm"

// ageKeyPrivateKey is NOT a real age private key
const ageKeyPrivateKey = "AGE-SECRET-KEY-HHS36XWKCVDKEKJ2M7WKQN3MFYUGIP4WWM7DT1CFANZUT5LT3K8ZRFZFGV3"

func TestNewKeyFunctionCreatesKeyWithCorrectName(t *testing.T) {
	t.Parallel()
	key := NewKey(wantedKeyName, wantedFilePath, ageKeyFileContent)

	if key.Name != wantedKeyName {
		t.Fatalf("Wanted name \"%s\" doesn't match with name on generated key: \"%s\"", wantedKeyName, key.Name)
	}
}

func TestNewKeyFunctionCreatesKeyWithCorrectPrivateKey(t *testing.T) {
	t.Parallel()
	key := NewKey(wantedKeyName, wantedFilePath, ageKeyFileContent)

	if key.PrivateKey != ageKeyPrivateKey {
		t.Fatalf("Wanted private key \"%s\" doesn't match with private key on generated key: \"%s\"", ageKeyPrivateKey, key.PrivateKey)
	}
}

func TestNewKeyFunctionCreatesKeyWithCorrectPublicKey(t *testing.T) {
	t.Parallel()
	key := NewKey(wantedKeyName, wantedFilePath, ageKeyFileContent)

	if key.PublicKey != ageKeyPublicKey {
		t.Fatalf("Wanted public key \"%s\" doesn't match with public key on generated key: \"%s\"", ageKeyPublicKey, key.PublicKey)
	}
}

func TestNewKeyFunctionCreatesKeyWithCorrectFilePath(t *testing.T) {
	t.Parallel()
	key := NewKey(wantedKeyName, wantedFilePath, ageKeyFileContent)

	if key.FileName != wantedFilePath {
		t.Fatalf("Wanted file path \"%s\" doesn't match with file path on generated key: \"%s\"", wantedFilePath, key.FileName)
	}
}

func TestNewKeyFunctionCreatesKeyWithCorrectPrivateKeyWithoutNewline(t *testing.T) {
	t.Parallel()
	key := NewKey(wantedKeyName, wantedFilePath, ageKeyFileContentWithNewLine)

	if key.PrivateKey != ageKeyPrivateKey {
		t.Fatalf("Wanted private key \"%s\" doesn't match with private key on generated key: \"%s\"", ageKeyPrivateKey, key.PrivateKey)
	}
}

func TestKeySetActiveFunctions(t *testing.T) {
	testDir := test.GenerateNewUniqueTestDir(t)
	defer testDir.CleanTestDir(t)

	configPath := filepath.Join(testDir.Path, "config.yaml")
	t.Setenv("SOPS_AGE_MANAGER_CONFIG_DIR", configPath)

	cfg := config.NewConfig("initial-enc", "initial-dec", testDir.Path)
	if err := cfg.Write(); err != nil {
		t.Fatalf("could not write initial config: %v", err)
	}

	targetKey := &Key{Name: "target-key"}

	// SetActiveEncryption
	targetKey.SetActiveEncryption()
	loadedCfg, err := config.NewConfigFromFile()
	if err != nil {
		t.Fatalf("could not read config: %v", err)
	}
	if loadedCfg.EncryptionKeyName != "target-key" || loadedCfg.DecryptionKeyName != "initial-dec" {
		t.Errorf("unexpected keys after SetActiveEncryption: enc=%s, dec=%s",
			loadedCfg.EncryptionKeyName, loadedCfg.DecryptionKeyName)
	}

	// SetActiveDecryption
	targetKey.SetActiveDecryption()
	loadedCfg, err = config.NewConfigFromFile()
	if err != nil {
		t.Fatalf("could not read config: %v", err)
	}
	if loadedCfg.EncryptionKeyName != "target-key" || loadedCfg.DecryptionKeyName != "target-key" {
		t.Errorf("unexpected keys after SetActiveDecryption: enc=%s, dec=%s",
			loadedCfg.EncryptionKeyName, loadedCfg.DecryptionKeyName)
	}

	// Reset and test SetActiveBoth
	cfg.EncryptionKeyName = "initial-enc"
	cfg.DecryptionKeyName = "initial-dec"
	if err := cfg.Write(); err != nil {
		t.Fatalf("could not reset config: %v", err)
	}

	targetKey.SetActiveBoth()
	loadedCfg, err = config.NewConfigFromFile()
	if err != nil {
		t.Fatalf("could not read config: %v", err)
	}
	if loadedCfg.EncryptionKeyName != "target-key" || loadedCfg.DecryptionKeyName != "target-key" {
		t.Errorf("unexpected keys after SetActiveBoth: enc=%s, dec=%s",
			loadedCfg.EncryptionKeyName, loadedCfg.DecryptionKeyName)
	}
}
