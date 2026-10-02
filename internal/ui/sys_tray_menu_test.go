//go:build darwin

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
	menu *fyne.Menu
	icon fyne.Resource
}

func (m *mockFyneDesktopApp) SetSystemTrayMenu(menu *fyne.Menu) {
	m.menu = menu
}

func (m *mockFyneDesktopApp) SetSystemTrayIcon(icon fyne.Resource) {
	m.icon = icon
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

	CreateSysTrayMenu(deskApp, keys, cfg, logo)

	if deskApp.menu == nil {
		t.Fatal("expected menu to be created and set")
	}

	if deskApp.menu.Label != "SAM" {
		t.Errorf("expected menu label 'SAM', got '%s'", deskApp.menu.Label)
	}

	if deskApp.icon == nil || deskApp.icon.Name() != "Logo.png" {
		t.Errorf("expected icon 'Logo.png', got %v", deskApp.icon)
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
