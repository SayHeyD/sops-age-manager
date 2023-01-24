package key

import (
	"github.com/SayHeyD/sops-age-manager/test"
	"os"
	"testing"
)

type TestKeyFiles struct {
	Name        string
	FileName    string
	FileContent string
	PrivateKey  string
	PublicKey   string
}

// getTestKeys returns age key values for testing. The returned keys are not real key pairs.
func getTestKeys() []*TestKeyFiles {
	return []*TestKeyFiles{
		{
			Name:     "first-key",
			FileName: "first-key.txt",
			FileContent: `# created: 2023-01-19T18:37:24+01:00
# public key: age1z9zvlcr2j3gt7mc9flmvyxm264v5aqyq0u2l46rlkg2c2fdzytgx7xl3qm
AGE-SECRET-KEY-HHS36XWKCVDKEKJ2M7WKQN3MFYUGIP4WWM7DT1CFANZUT5LT3K8ZRFZFGV3`,
			PrivateKey: "AGE-SECRET-KEY-HHS36XWKCVDKEKJ2M7WKQN3MFYUGIP4WWM7DT1CFANZUT5LT3K8ZRFZFGV3",
			PublicKey:  "age1z9zvlcr2j3gt7mc9flmvyxm264v5aqyq0u2l46rlkg2c2fdzytgx7xl3qm",
		},
		{
			Name:     "second-key",
			FileName: "second-key.txt",
			FileContent: `# created: 2023-01-19T18:37:24+01:00
# public key: agevyxm264v5aqyq0u21z9zvlcr2j3gt7mc9flml46rlkg2c2fdzytgx7xl3qm
AGE-SECRET-KEY-36XWKCVDKEKJ2M7WKQN3MFYUHHSGIP4WWM7DT1CFANZUT5LT3K8ZRFZFGV3`,
			PrivateKey: "AGE-SECRET-KEY-36XWKCVDKEKJ2M7WKQN3MFYUHHSGIP4WWM7DT1CFANZUT5LT3K8ZRFZFGV3",
			PublicKey:  "agevyxm264v5aqyq0u21z9zvlcr2j3gt7mc9flml46rlkg2c2fdzytgx7xl3qm",
		},
		{
			Name:     "third-key",
			FileName: "third-key.txt",
			FileContent: `# created: 2023-01-19T18:37:24+01:00
# public key: agel46rlkg2c2fdr2j3gt7mc9flmvyxmzytgx7xl3qm1z9zvlc264v5aqyq0u2
AGE-SECRET-KEY-HFANZUT5LT3K8ZRFZFHS36XWKCVDKEKJ2M7WKQN3MFYUGIP4WWM7DT1CGV3`,
			PrivateKey: "AGE-SECRET-KEY-HFANZUT5LT3K8ZRFZFHS36XWKCVDKEKJ2M7WKQN3MFYUGIP4WWM7DT1CGV3",
			PublicKey:  "agel46rlkg2c2fdr2j3gt7mc9flmvyxmzytgx7xl3qm1z9zvlc264v5aqyq0u2",
		},
	}
}

/*
prepareKeyTestDir generates a temporary directory, with a unique name, where key files will be written to
in order to execute tests in an isolated directory. The string returned is the absolute filepath of the directory.
*/
func prepareKeyTestDir(t *testing.T) *test.Dir {
	testDir := test.GenerateNewUniqueTestDir(t)

	keys := getTestKeys()

	for _, key := range keys {
		filePath := testDir.Path + string(os.PathSeparator) + key.FileName

		keyFile, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("Could not create file for testing \"%s\": %v", filePath, err)
		}

		if _, err := keyFile.WriteString(key.FileContent); err != nil {
			t.Fatalf("Could not write to file \"%s\": %v", filePath, err)
		}

		if err := keyFile.Close(); err != nil {
			t.Fatalf("Could not close file \"%s\": %v", filePath, err)
		}
	}

	return testDir
}

func TestGetAvailableKeysReturnsCorrectAmountOfKeys(t *testing.T) {
	t.Parallel()
	testDir := prepareKeyTestDir(t)
	defer testDir.CleanTestDir(t)

	keys := GetAvailableKeys(testDir.Path)

	lengthOfFetchedKeys := len(keys)
	lengthOfTestKeys := len(getTestKeys())

	if lengthOfFetchedKeys != lengthOfTestKeys {
		t.Fatalf(
			"Length of Keys fetched by \"GetAvailableKeys\" (%d) does not match length of \"getTestKeys()\" (%d)",
			lengthOfFetchedKeys, lengthOfTestKeys,
		)
	}
}

func TestGetAvailableKeysReturnsCorrectKeys(t *testing.T) {
	t.Parallel()
	testDir := prepareKeyTestDir(t)
	defer testDir.CleanTestDir(t)

	keys := GetAvailableKeys(testDir.Path)
	wantedKeys := getTestKeys()

	for i, key := range keys {
		wantedKey := wantedKeys[i]

		if wantedKey == nil {
			t.Fatalf("GetAvailableKeys() returned a key that should not exist: \"%v\"", key)
		}

		if key.Name != wantedKey.Name {
			t.Fatalf("key name \"%s\" does not match with expected key name \"%s\"", key.Name, wantedKey.Name)
		}

		if key.PrivateKey != wantedKey.PrivateKey {
			t.Fatalf("private key \"%s\" does not match with expected private key \"%s\"", key.PrivateKey, wantedKey.PrivateKey)
		}

		if key.PublicKey != wantedKey.PublicKey {
			t.Fatalf("public key \"%s\" does not match with expected public key \"%s\"", key.PublicKey, wantedKey.PublicKey)
		}
	}
}
