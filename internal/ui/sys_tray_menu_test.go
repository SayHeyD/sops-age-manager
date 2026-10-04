//go:build darwin || windows

package ui

import (
	"testing"

	"fyne.io/fyne/v2"
	fyneTest "fyne.io/fyne/v2/test"
	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
)

type mockFyneDesktopApp struct {
	fyne.App
	menu          *fyne.Menu
	icon          fyne.Resource
	menuCallCount int
	iconCallCount int
}

func (m *mockFyneDesktopApp) SetSystemTrayMenu(menu *fyne.Menu) {
	m.menu = menu
	m.menuCallCount++
}

func (m *mockFyneDesktopApp) SetSystemTrayIcon(icon fyne.Resource) {
	m.icon = icon
	m.iconCallCount++
}

func TestCreateSysTrayMenuWithDesktopApp(t *testing.T) {
	cleanup := setupTestConfig(t, "key-1", "key-2")
	defer cleanup()

	baseApp := fyneTest.NewApp()
	deskApp := &mockFyneDesktopApp{App: baseApp}

	keys := []*key.Key{
		{Name: "key-1", PublicKey: "pub1", PrivateKey: "priv1"},
	}
	cfg := config.NewConfig("key-1", "key-2", "")
	logo := []byte("logo-bytes")

	tm := CreateSysTrayMenu(deskApp, keys, cfg, logo)
	if tm != nil {
		defer tm.StopWatching()
	}

	if deskApp.menu == nil {
		t.Fatal("expected menu to be created and set")
	}

	if deskApp.menu.Label != "SAM" {
		t.Errorf("expected menu label 'SAM', got '%s'", deskApp.menu.Label)
	}

	if deskApp.icon == nil || deskApp.icon.Name() != "Logo.png" {
		t.Errorf("expected icon 'Logo.png', got %v", deskApp.icon)
	}

	if deskApp.menuCallCount != 1 {
		t.Errorf("expected SetSystemTrayMenu to be called exactly 1 time, got %d", deskApp.menuCallCount)
	}
	if deskApp.iconCallCount != 1 {
		t.Errorf("expected SetSystemTrayIcon to be called exactly 1 time, got %d", deskApp.iconCallCount)
	}
}

func TestCreateSysTrayMenuNoKeys(t *testing.T) {
	cleanup := setupTestConfig(t, "", "")
	defer cleanup()

	baseApp := fyneTest.NewApp()
	deskApp := &mockFyneDesktopApp{App: baseApp}

	cfg := config.NewConfig("", "", "")
	logo := []byte("logo-bytes")

	tm := CreateSysTrayMenu(deskApp, nil, cfg, logo)
	if tm != nil {
		defer tm.StopWatching()
	}

	if deskApp.menu == nil {
		t.Fatal("expected menu to be created and set")
	}

	if len(deskApp.menu.Items) != 1 {
		t.Fatalf("expected 1 root menu item, got %d", len(deskApp.menu.Items))
	}

	keysItem := deskApp.menu.Items[0]
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
}

func TestCreateSysTrayMenuWithNonDesktopApp(t *testing.T) {
	plainApp := fyneTest.NewApp()
	keys := []*key.Key{
		{Name: "key-1", PublicKey: "pub1", PrivateKey: "priv1"},
	}
	cfg := config.NewConfig("key-1", "key-2", "")
	logo := []byte("logo-bytes")

	// Should safely return early without panic
	CreateSysTrayMenu(plainApp, keys, cfg, logo)
}
