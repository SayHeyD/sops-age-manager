package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop" //nolint:typecheck
	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
	"github.com/atotto/clipboard"
	"log"
)

// TODO: restructure the whole file, this reads like spaghetti leftovers of a 3 year old

const (
	ConfigEncryption           = 0
	ConfigDecryption           = 1
	ConfigEncryptionDecryption = 2
)

func CreateSysTrayMenu(a fyne.App, keys []*key.Key, config *config.Config) {
	var menuItems []*fyne.MenuItem

	desk, ok := a.(desktop.App)

	for _, ageKey := range keys {

		selectedKey := *ageKey

		ageKeyEncryptionDecryptionMenuEntry := fyne.NewMenuItem("Encryption and decryption", func() {
			UpdateSysTrayMenu(desk, menuItems, selectedKey, ConfigEncryptionDecryption)
		})

		ageKeyEncryptionMenuEntry := fyne.NewMenuItem("Encryption", func() {
			UpdateSysTrayMenu(desk, menuItems, selectedKey, ConfigEncryption)
		})

		ageKeyDecryptionMenuEntry := fyne.NewMenuItem("Decryption", func() {
			UpdateSysTrayMenu(desk, menuItems, selectedKey, ConfigDecryption)
		})

		ageKeyMenu := fyne.NewMenuItem(selectedKey.Name, func() {})
		ageKeyMenu.ChildMenu = fyne.NewMenu("key options for "+selectedKey.Name,
			ageKeyEncryptionDecryptionMenuEntry,
			ageKeyEncryptionMenuEntry,
			ageKeyDecryptionMenuEntry,
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Copy key name", func() {
				err := clipboard.WriteAll(selectedKey.Name)
				if err != nil {
					log.Fatalf("could not set clipboard text: %v", err)
				}
			}),
			fyne.NewMenuItem("Copy public key", func() {
				err := clipboard.WriteAll(selectedKey.PublicKey)
				if err != nil {
					log.Fatalf("could not set clipboard text: %v", err)
				}
			}),
			fyne.NewMenuItem("Copy private key", func() {
				err := clipboard.WriteAll(selectedKey.PrivateKey)
				if err != nil {
					log.Fatalf("could not set clipboard text: %v", err)
				}
			}))

		menuItems = append(menuItems, ageKeyMenu)
	}

	if ok {
		setMenu(desk, menuItems)
	}
}

func setMenu(desk desktop.App, menuItems []*fyne.MenuItem) {
	keyMenu := fyne.NewMenuItem("Keys", func() {})
	keyMenu.ChildMenu = fyne.NewMenu("Key menu", menuItems...)

	m := fyne.NewMenu("SAM", keyMenu)
	desk.SetSystemTrayMenu(m)
}

func UpdateSysTrayMenu(desk desktop.App, menuItems []*fyne.MenuItem, key key.Key, setMode uint) {

	samConfig, err := config.NewConfigFromFile()
	if err != nil {
		log.Fatal(err)
	}

	currentEncryptionKey := samConfig.EncryptionKeyName
	currentDecryptionKey := samConfig.DecryptionKeyName

	for _, menuItem := range menuItems {

		childMenuItems := menuItem.ChildMenu.Items

		for _, childMenuItem := range childMenuItems {
			childMenuItem.Checked = false
		}

		if key.Name == menuItem.Label {

			if setMode == ConfigEncryptionDecryption {
				for _, childMenuItem := range childMenuItems {
					if childMenuItem.Label == "Encryption and decryption" {
						childMenuItem.Checked = true
						key.SetActiveEncryption()
						key.SetActiveDecryption()
					} else {
						childMenuItem.Checked = false
					}
				}
			} else if setMode == ConfigEncryption {
				for _, childMenuItem := range childMenuItems {
					if childMenuItem.Label == "Encryption" {
						childMenuItem.Checked = true
						key.SetActiveEncryption()
					} else {
						childMenuItem.Checked = false
					}
				}
			} else if setMode == ConfigDecryption {
				for _, childMenuItem := range childMenuItems {
					if childMenuItem.Label == "Decryption" {
						childMenuItem.Checked = true
						key.SetActiveDecryption()
					} else {
						childMenuItem.Checked = false
					}
				}
			}
		}

		if setMode == ConfigEncryption && currentEncryptionKey == menuItem.Label {
			for _, childMenuItem := range childMenuItems {
				if childMenuItem.Label == "Decryption" {
					childMenuItem.Checked = true
				}
			}
		}

		if setMode == ConfigDecryption && currentDecryptionKey == menuItem.Label {
			for _, childMenuItem := range childMenuItems {
				if childMenuItem.Label == "Encryption" {
					childMenuItem.Checked = true
				}
			}
		}
	}

	setMenu(desk, menuItems)
}
