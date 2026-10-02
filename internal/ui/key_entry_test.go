//go:build darwin

package ui

import (
	"testing"

	"github.com/SayHeyD/sops-age-manager/pkg/key"
)

func TestKeyEntryUpdateChecked(t *testing.T) {
	testKey := &key.Key{
		Name:       "test-key",
		PublicKey:  "age1publickeytest",
		PrivateKey: "AGE-SECRET-KEY-test",
	}

	var selectedMode KeyMode
	var selectedKey *key.Key
	entry := newKeyEntry(testKey, func(k *key.Key, mode KeyMode) {
		selectedKey = k
		selectedMode = mode
	})

	if entry.menuItem.Label != "test-key" {
		t.Fatalf("expected menu label 'test-key', got '%s'", entry.menuItem.Label)
	}

	if entry.menuItem.ChildMenu == nil {
		t.Fatal("expected child menu to be non-nil")
	}

	// 1. Neither encryption nor decryption
	entry.updateChecked("other-key-1", "other-key-2")
	if entry.bothItem.Checked || entry.encryptionItem.Checked || entry.decryptionItem.Checked {
		t.Errorf("expected no checked items when key is not active, got both=%v, enc=%v, dec=%v",
			entry.bothItem.Checked, entry.encryptionItem.Checked, entry.decryptionItem.Checked)
	}

	// 2. Encryption only
	entry.updateChecked("test-key", "other-key-2")
	if !entry.encryptionItem.Checked || entry.bothItem.Checked || entry.decryptionItem.Checked {
		t.Errorf("expected only encryption checked, got both=%v, enc=%v, dec=%v",
			entry.bothItem.Checked, entry.encryptionItem.Checked, entry.decryptionItem.Checked)
	}

	// 3. Decryption only
	entry.updateChecked("other-key-1", "test-key")
	if !entry.decryptionItem.Checked || entry.bothItem.Checked || entry.encryptionItem.Checked {
		t.Errorf("expected only decryption checked, got both=%v, enc=%v, dec=%v",
			entry.bothItem.Checked, entry.encryptionItem.Checked, entry.decryptionItem.Checked)
	}

	// 4. Both encryption and decryption
	entry.updateChecked("test-key", "test-key")
	if !entry.bothItem.Checked || entry.encryptionItem.Checked || entry.decryptionItem.Checked {
		t.Errorf("expected only bothItem checked, got both=%v, enc=%v, dec=%v",
			entry.bothItem.Checked, entry.encryptionItem.Checked, entry.decryptionItem.Checked)
	}

	// Test action callbacks
	entry.bothItem.Action()
	if selectedKey != testKey || selectedMode != ModeBoth {
		t.Errorf("expected action to trigger ModeBoth on testKey, got key=%v, mode=%v", selectedKey, selectedMode)
	}

	entry.encryptionItem.Action()
	if selectedKey != testKey || selectedMode != ModeEncryption {
		t.Errorf("expected action to trigger ModeEncryption on testKey, got key=%v, mode=%v", selectedKey, selectedMode)
	}

	entry.decryptionItem.Action()
	if selectedKey != testKey || selectedMode != ModeDecryption {
		t.Errorf("expected action to trigger ModeDecryption on testKey, got key=%v, mode=%v", selectedKey, selectedMode)
	}
}
