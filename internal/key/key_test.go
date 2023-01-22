package key

import (
	"os"
	"testing"
)

const wantedKeyName = "test_key"

const wantedFilePath = "/var/someDir/tmp/test_key.txt"

// ageKeyFileContent is NOT a real age key pair
const ageKeyFileContent = `# created: 2023-01-19T18:37:24+01:00
# public key: age1z9zvlcr2j3gt7mc9flmvyxm264v5aqyq0u2l46rlkg2c2fdzytgx7xl3qm
AGE-SECRET-KEY-HHS36XWKCVDKEKJ2M7WKQN3MFYUGIP4WWM7DT1CFANZUT5LT3K8ZRFZFGV3`

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

func TestSetActiveSetsEnvVarCorrectly(t *testing.T) {
	t.Parallel()
	key := NewKey(wantedKeyName, wantedFilePath, ageKeyFileContent)
	key.SetActive()

	envVarValue := os.Getenv("SOPS_AGE_KEY")

	if envVarValue != ageKeyPrivateKey {
		t.Fatalf("Set active did not set the env \"SOPS_AGE_KEY\" to \"%s\". Value is: \"%s\"", ageKeyPrivateKey, key.PublicKey)
	}
}
