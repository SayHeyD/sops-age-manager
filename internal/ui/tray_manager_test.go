//go:build darwin || windows

package ui

import (
	"os"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
	"github.com/SayHeyD/sops-age-manager/test"
)

type mockDesktopApp struct {
	menu          *fyne.Menu
	icon          fyne.Resource
	menuCallCount int
	iconCallCount int
}

func (m *mockDesktopApp) SetSystemTrayMenu(menu *fyne.Menu) {
	m.menu = menu
	m.menuCallCount++
}

func (m *mockDesktopApp) SetSystemTrayIcon(icon fyne.Resource) {
	m.icon = icon
	m.iconCallCount++
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

func TestNewTrayManagerNoKeys(t *testing.T) {
	desk := &mockDesktopApp{}
	logo := []byte("fake-png-data")

	tm := NewTrayManager(desk, nil, logo)
	if tm == nil {
		t.Fatal("expected NewTrayManager to return a non-nil TrayManager")
	}

	if len(tm.entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(tm.entries))
	}

	if desk.menu == nil {
		t.Fatal("expected SetSystemTrayMenu to be called")
	}

	if len(desk.menu.Items) != 1 {
		t.Fatalf("expected 1 root menu item, got %d", len(desk.menu.Items))
	}

	keysItem := desk.menu.Items[0]
	if keysItem.Label != "Keys" || keysItem.ChildMenu == nil {
		t.Fatalf("expected 'Keys' menu item with child menu, got label '%s'", keysItem.Label)
	}

	if len(keysItem.ChildMenu.Items) != 1 {
		t.Fatalf("expected 1 item in Keys child menu, got %d", len(keysItem.ChildMenu.Items))
	}

	noKeysItem := keysItem.ChildMenu.Items[0]
	if noKeysItem.Label != "No keys found" {
		t.Fatalf("expected 'No keys found' item, got '%s'", noKeysItem.Label)
	}

	if !noKeysItem.Disabled {
		t.Fatal("expected 'No keys found' item to be disabled")
	}

	// Calling Refresh when no keys are loaded should not panic
	tm.Refresh()
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

	if desk.menuCallCount != 1 {
		t.Errorf("expected SetSystemTrayMenu to be called exactly 1 time, got %d", desk.menuCallCount)
	}
	if desk.iconCallCount != 1 {
		t.Errorf("expected SetSystemTrayIcon to be called exactly 1 time, got %d", desk.iconCallCount)
	}

	// Calling Refresh again should update checkmarks in-place without re-calling SetSystemTrayMenu
	tm.Refresh()
	if desk.menuCallCount != 1 {
		t.Errorf("expected SetSystemTrayMenu count to remain 1 after second Refresh, got %d", desk.menuCallCount)
	}
	if desk.iconCallCount != 1 {
		t.Errorf("expected SetSystemTrayIcon count to remain 1 after second Refresh, got %d", desk.iconCallCount)
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

func TestTrayManagerWatchAndExternalConfigChange(t *testing.T) {
	cleanup := setupTestConfig(t, "key-1", "key-1")
	defer cleanup()

	cfg, err := config.NewConfigFromFile()
	if err != nil {
		t.Fatalf("could not read config: %v", err)
	}

	desk := &mockDesktopApp{}
	key1 := &key.Key{Name: "key-1", PublicKey: "pub1", PrivateKey: "priv1"}
	key2 := &key.Key{Name: "key-2", PublicKey: "pub2", PrivateKey: "priv2"}
	keys := []*key.Key{key1, key2}

	tm := NewTrayManager(desk, keys, nil)
	tm.Refresh()

	if err := tm.StartWatching(cfg.Path); err != nil {
		t.Fatalf("failed to start watching: %v", err)
	}
	defer tm.StopWatching()

	// Initially key-1 is both encryption and decryption
	if !tm.entries[0].bothItem.Checked || tm.entries[1].bothItem.Checked {
		t.Fatalf("unexpected initial state: entry0 both=%v, entry1 both=%v",
			tm.entries[0].bothItem.Checked, tm.entries[1].bothItem.Checked)
	}

	// External CLI updates config to key-2 for decryption
	cfg.DecryptionKeyName = "key-2"
	if err := cfg.Write(); err != nil {
		t.Fatalf("could not write updated config: %v", err)
	}

	// Wait for watcher to trigger refresh
	time.Sleep(200 * time.Millisecond)

	// Verify key-1 is now encryption only, and key-2 is decryption only
	if !tm.entries[0].encryptionItem.Checked || tm.entries[0].bothItem.Checked {
		t.Errorf("expected entry0 to be encryption only, got enc=%v, both=%v",
			tm.entries[0].encryptionItem.Checked, tm.entries[0].bothItem.Checked)
	}
	if !tm.entries[1].decryptionItem.Checked || tm.entries[1].bothItem.Checked {
		t.Errorf("expected entry1 to be decryption only, got dec=%v, both=%v",
			tm.entries[1].decryptionItem.Checked, tm.entries[1].bothItem.Checked)
	}

	// External CLI updates config to key-2 for both
	cfg.EncryptionKeyName = "key-2"
	if err := cfg.Write(); err != nil {
		t.Fatalf("could not write updated config: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	if !tm.entries[1].bothItem.Checked || tm.entries[0].bothItem.Checked || tm.entries[0].encryptionItem.Checked {
		t.Errorf("expected entry1 to be both, got entry0 enc=%v, entry1 both=%v",
			tm.entries[0].encryptionItem.Checked, tm.entries[1].bothItem.Checked)
	}

	// Stop watching idempotency
	tm.StopWatching()
	tm.StopWatching()
}

func TestTrayManagerHandleModeSelectionUpdatesViaWatcher(t *testing.T) {
	cleanup := setupTestConfig(t, "key-1", "key-1")
	defer cleanup()

	cfg, err := config.NewConfigFromFile()
	if err != nil {
		t.Fatalf("could not read config: %v", err)
	}

	desk := &mockDesktopApp{}
	key1 := &key.Key{Name: "key-1", PublicKey: "pub1", PrivateKey: "priv1"}
	key2 := &key.Key{Name: "key-2", PublicKey: "pub2", PrivateKey: "priv2"}
	keys := []*key.Key{key1, key2}

	tm := NewTrayManager(desk, keys, nil)
	tm.Refresh()

	if err := tm.StartWatching(cfg.Path); err != nil {
		t.Fatalf("failed to start watching: %v", err)
	}
	defer tm.StopWatching()

	// Initial check: key-1 is both
	if !tm.entries[0].bothItem.Checked {
		t.Fatalf("expected key-1 bothItem checked initially")
	}

	// Trigger UI click on key2 encryption
	tm.entries[1].encryptionItem.Action()

	// Wait for watcher to trigger refresh
	time.Sleep(200 * time.Millisecond)

	if !tm.entries[1].encryptionItem.Checked || tm.entries[1].bothItem.Checked {
		t.Errorf("expected key-2 to be encryption only via watcher, got enc=%v, both=%v",
			tm.entries[1].encryptionItem.Checked, tm.entries[1].bothItem.Checked)
	}
	if !tm.entries[0].decryptionItem.Checked || tm.entries[0].bothItem.Checked {
		t.Errorf("expected key-1 to be decryption only via watcher, got dec=%v, both=%v",
			tm.entries[0].decryptionItem.Checked, tm.entries[0].bothItem.Checked)
	}

	// Trigger UI click on key2 both
	tm.entries[1].bothItem.Action()
	time.Sleep(200 * time.Millisecond)

	if !tm.entries[1].bothItem.Checked || tm.entries[0].bothItem.Checked {
		t.Errorf("expected key-2 to be both via watcher, got key2 both=%v, key1 both=%v",
			tm.entries[1].bothItem.Checked, tm.entries[0].bothItem.Checked)
	}
}
