package key

import (
	"os"
	"testing"
)

// ageKeyFileContent is NOT a real age key pair
const ageKeyFileContent = `# created: 2023-01-19T18:37:24+01:00
# public key: age1z9zvlcr2j3gt7mc9flmvyxm264v5aqyq0u2l46rlkg2c2fdzytgx7xl3qm
AGE-SECRET-KEY-HHS36XWKCVDKEKJ2M7WKQN3MFYUGIP4WWM7DT1CFANZUT5LT3K8ZRFZFGV3`

// ageKeyPublicKey is NOT a real age public key
const ageKeyPublicKey = "age1z9zvlcr2j3gt7mc9flmvyxm264v5aqyq0u2l46rlkg2c2fdzytgx7xl3qm"

// ageKeyPrivateKey is NOT a real age private key
const ageKeyPrivateKey = "AGE-SECRET-KEY-HHS36XWKCVDKEKJ2M7WKQN3MFYUGIP4WWM7DT1CFANZUT5LT3K8ZRFZFGV3"

func TestNewKeyFunctionCreatesKeyWithCorrectName(t *testing.T) {
	wantedKeyName := "test_key"
	key := NewKey(wantedKeyName, ageKeyFileContent)

	if key.Name != wantedKeyName {
		t.Fatalf("Wanted name \"%s\" doesn't match with name on generated key: \"%s\"", wantedKeyName, key.Name)
	}
}

func TestNewKeyFunctionCreatesKeyWithCorrectPrivateKey(t *testing.T) {
	key := NewKey("test_key", ageKeyFileContent)

	if key.PrivateKey != ageKeyPrivateKey {
		t.Fatalf("Wanted private key \"%s\" doesn't match with private key on generated key: \"%s\"", ageKeyPrivateKey, key.PrivateKey)
	}
}

func TestNewKeyFunctionCreatesKeyWithCorrectPublicKey(t *testing.T) {
	key := NewKey("test_key", ageKeyFileContent)

	if key.PublicKey != ageKeyPublicKey {
		t.Fatalf("Wanted public key \"%s\" doesn't match with public key on generated key: \"%s\"", ageKeyPublicKey, key.PublicKey)
	}
}

func TestSetActiveSetsEnvVarCorrectly(t *testing.T) {
	key := NewKey("test_key", ageKeyFileContent)
	key.SetActive()

	envVarValue := os.Getenv("SOPS_AGE_KEY")

	if envVarValue != ageKeyPrivateKey {
		t.Fatalf("Set active did not set the env \"SOPS_AGE_KEY\" to \"%s\". Value is: \"%s\"", ageKeyPrivateKey, key.PublicKey)
	}
}
