//go:build darwin

package ui

import (
	"os"
	"testing"

	"fyne.io/fyne/v2"
	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
	"github.com/SayHeyD/sops-age-manager/test"
)

type mockDesktopApp struct {
	menu *fyne.Menu
	icon fyne.Resource
}

func (m *mockDesktopApp) SetSystemTrayMenu(menu *fyne.Menu) {
	m.menu = menu
}

func (m *mockDesktopApp) SetSystemTrayIcon(icon fyne.Resource) {
	m.icon = icon
}

func setupTestConfig(t *testing.T, encKey, decKey string) func() {
	t.Helper()
	testDir := test.GenerateNewUniqueTestDir(t)
	configPath := testDir.Path + string(os.PathSeparator) + "config.yaml"

	oldEnv := os.Getenv("SOPS_AGE_MANAGER_CONFIG_DIR")
	if err := os.Setenv("SOPS_AGE_MANAGER_CONFIG_DIR", configPath); err != nil {
		t.Fatalf("could not set config env: %v", err)
	}

	cfg := config.NewConfig(encKey, decKey, testDir.Path)
	if err := cfg.Write(); err != nil {
		t.Fatalf("could not write test config: %v", err)
	}

	return func() {
		testDir.CleanTestDir(t)
		if oldEnv != "" {
			_ = os.Setenv("SOPS_AGE_MANAGER_CONFIG_DIR", oldEnv)
		} else {
			_ = os.Unsetenv("SOPS_AGE_MANAGER_CONFIG_DIR")
		}
	}
}

func TestNewTrayManager(t *testing.T) {
	desk := &mockDesktopApp{}
	keys := []*key.Key{
		{Name: "key-1", PublicKey: "pub1", PrivateKey: "priv1"},
		{Name: "key-2", PublicKey: "pub2", PrivateKey: "priv2"},
	}
	logo := []byte("fake-png-data")

	tm := NewTrayManager(desk, keys, logo)
	if tm == nil {
		t.Fatal("expected NewTrayManager to return a non-nil TrayManager")
	}

	if len(tm.entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(tm.entries))
	}

	if string(tm.logo) != string(logo) {
		t.Fatalf("expected logo '%s', got '%s'", logo, tm.logo)
	}
}

func TestTrayManagerRefresh(t *testing.T) {
	cleanup := setupTestConfig(t, "key-1", "key-2")
	defer cleanup()

	desk := &mockDesktopApp{}
	keys := []*key.Key{
		{Name: "key-1", PublicKey: "pub1", PrivateKey: "priv1"},
		{Name: "key-2", PublicKey: "pub2", PrivateKey: "priv2"},
	}
	logo := []byte("test-logo")

	tm := NewTrayManager(desk, keys, logo)
	tm.Refresh()

	if desk.menu == nil {
		t.Fatal("expected SetSystemTrayMenu to be called")
	}

	if desk.menu.Label != "SAM" {
		t.Fatalf("expected root menu label 'SAM', got '%s'", desk.menu.Label)
	}

	if len(desk.menu.Items) != 1 {
		t.Fatalf("expected 1 root menu item, got %d", len(desk.menu.Items))
	}

	keysItem := desk.menu.Items[0]
	if keysItem.Label != "Keys" || keysItem.ChildMenu == nil {
		t.Fatalf("expected 'Keys' menu item with child menu, got label '%s'", keysItem.Label)
	}

	if len(keysItem.ChildMenu.Items) != 2 {
		t.Fatalf("expected 2 key menu items, got %d", len(keysItem.ChildMenu.Items))
	}

	// Verify checkmarks on key-1 (encryption only) and key-2 (decryption only)
	entry1 := tm.entries[0]
	if !entry1.encryptionItem.Checked || entry1.decryptionItem.Checked || entry1.bothItem.Checked {
		t.Errorf("entry1 checkmarks invalid: enc=%v, dec=%v, both=%v",
			entry1.encryptionItem.Checked, entry1.decryptionItem.Checked, entry1.bothItem.Checked)
	}

	entry2 := tm.entries[1]
	if entry2.encryptionItem.Checked || !entry2.decryptionItem.Checked || entry2.bothItem.Checked {
		t.Errorf("entry2 checkmarks invalid: enc=%v, dec=%v, both=%v",
			entry2.encryptionItem.Checked, entry2.decryptionItem.Checked, entry2.bothItem.Checked)
	}

	if desk.icon == nil || desk.icon.Name() != "Logo.png" {
		t.Errorf("expected tray icon 'Logo.png', got %v", desk.icon)
	}
}

func TestTrayManagerHandleModeSelection(t *testing.T) {
	cleanup := setupTestConfig(t, "key-1", "key-1")
	defer cleanup()

	desk := &mockDesktopApp{}
	testKey := &key.Key{Name: "test-target", PublicKey: "pub", PrivateKey: "priv"}
	keys := []*key.Key{testKey}

	tm := NewTrayManager(desk, keys, nil)

	// ModeEncryption
	tm.handleModeSelection(testKey, ModeEncryption)
	cfg, err := config.NewConfigFromFile()
	if err != nil {
		t.Fatalf("could not read config: %v", err)
	}
	if cfg.EncryptionKeyName != "test-target" {
		t.Errorf("expected encryption key 'test-target', got '%s'", cfg.EncryptionKeyName)
	}

	// ModeDecryption
	tm.handleModeSelection(testKey, ModeDecryption)
	cfg, err = config.NewConfigFromFile()
	if err != nil {
		t.Fatalf("could not read config: %v", err)
	}
	if cfg.DecryptionKeyName != "test-target" {
		t.Errorf("expected decryption key 'test-target', got '%s'", cfg.DecryptionKeyName)
	}

	// ModeBoth
	tm.handleModeSelection(testKey, ModeBoth)
	cfg, err = config.NewConfigFromFile()
	if err != nil {
		t.Fatalf("could not read config: %v", err)
	}
	if cfg.EncryptionKeyName != "test-target" || cfg.DecryptionKeyName != "test-target" {
		t.Errorf("expected both keys 'test-target', got enc='%s', dec='%s'",
			cfg.EncryptionKeyName, cfg.DecryptionKeyName)
	}
}
