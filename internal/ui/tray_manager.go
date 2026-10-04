//go:build darwin || windows

package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop" //nolint:typecheck
	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
)

type TrayManager struct {
	desk    desktop.App
	entries []*keyEntry
	menu    *fyne.Menu
	logo    []byte
	watcher *ConfigWatcher
}

func NewTrayManager(desk desktop.App, keys []*key.Key, logo []byte) *TrayManager {
	tm := &TrayManager{
		desk: desk,
		logo: logo,
	}

	tm.entries = make([]*keyEntry, len(keys))
	for i, k := range keys {
		tm.entries[i] = newKeyEntry(k, tm.handleModeSelection)
	}

	var menuItems []*fyne.MenuItem
	if len(tm.entries) == 0 {
		noKeysItem := fyne.NewMenuItem("No keys found", func() {})
		noKeysItem.Disabled = true
		menuItems = []*fyne.MenuItem{noKeysItem}
	} else {
		menuItems = make([]*fyne.MenuItem, len(tm.entries))
		for i, entry := range tm.entries {
			menuItems[i] = entry.menuItem
		}
	}

	keySubMenu := fyne.NewMenuItem("Keys", func() {})
	keySubMenu.ChildMenu = fyne.NewMenu("Key menu", menuItems...)

	clearActiveKeysItem := fyne.NewMenuItem("Clear Active Keys", func() {
		key.ClearActiveKeys()
	})

	openConfigDirItem := fyne.NewMenuItem("Open Config Directory", func() {
		if err := OpenConfigDirectory(); err != nil {
			log.Printf("could not open config directory: %v", err)
		}
	})

	openKeyDirItem := fyne.NewMenuItem("Open Key Directory", func() {
		var keyDir string
		if cfg, err := config.NewConfigFromFile(); err == nil && cfg != nil {
			keyDir = cfg.KeyDir
		}
		if err := OpenKeyDirectory(keyDir); err != nil {
			log.Printf("could not open key directory: %v", err)
		}
	})

	tm.menu = fyne.NewMenu("SAM",
		keySubMenu,
		clearActiveKeysItem,
		fyne.NewMenuItemSeparator(),
		openConfigDirItem,
		openKeyDirItem,
	)

	if tm.desk != nil {
		tm.desk.SetSystemTrayMenu(tm.menu)
		if len(tm.logo) > 0 {
			tm.desk.SetSystemTrayIcon(fyne.NewStaticResource("Logo.png", tm.logo))
		}
	}

	return tm
}

func (tm *TrayManager) handleModeSelection(k *key.Key, mode KeyMode) {
	switch mode {
	case ModeBoth:
		k.SetActiveBoth()
	case ModeEncryption:
		k.SetActiveEncryption()
	case ModeDecryption:
		k.SetActiveDecryption()
	}
}

func (tm *TrayManager) Refresh() {
	appConfig, err := config.NewConfigFromFile()
	if err != nil {
		log.Printf("could not load configuration: %v", err)
		return
	}

	for _, entry := range tm.entries {
		entry.updateChecked(appConfig.EncryptionKeyName, appConfig.DecryptionKeyName)
	}

	if tm.menu != nil {
		tm.menu.Refresh()
	}
}

func (tm *TrayManager) StartWatching(configPath string) error {
	if tm.watcher != nil {
		tm.watcher.Stop()
	}

	watcher, err := NewConfigWatcher(configPath, tm.Refresh)
	if err != nil {
		return err
	}

	tm.watcher = watcher
	return nil
}

func (tm *TrayManager) StopWatching() {
	if tm.watcher != nil {
		tm.watcher.Stop()
		tm.watcher = nil
	}
}
