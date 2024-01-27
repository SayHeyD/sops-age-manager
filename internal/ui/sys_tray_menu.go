package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop" //nolint:typecheck
	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
	"github.com/atotto/clipboard"
	"log"
)

func CreateSysTrayMenu(a fyne.App, keys []*key.Key, config *config.Config) {
	var menuItems []*fyne.MenuItem

	for _, ageKey := range keys {

		selectedKey := *ageKey

		ageKeyEncryptionDecryptionMenuEntry := fyne.NewMenuItem("Encryption and decryption", func() {
			selectedKey.SetActiveDecryption()
			selectedKey.SetActiveEncryption()
		})

		ageKeyEncryptionMenuEntry := fyne.NewMenuItem("Encryption", func() {
			selectedKey.SetActiveEncryption()
		})

		ageKeyDecryptionMenuEntry := fyne.NewMenuItem("Decryption", func() {
			selectedKey.SetActiveDecryption()
		})

		if config.DecryptionKeyName == selectedKey.Name && config.EncryptionKeyName == selectedKey.Name {
			ageKeyEncryptionDecryptionMenuEntry.Checked = true
		}

		if config.EncryptionKeyName == selectedKey.Name && config.DecryptionKeyName != selectedKey.Name {
			ageKeyEncryptionDecryptionMenuEntry.Checked = true
		}

		if config.EncryptionKeyName != selectedKey.Name && config.DecryptionKeyName == selectedKey.Name {
			ageKeyEncryptionDecryptionMenuEntry.Checked = true
		}

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

	if desk, ok := a.(desktop.App); ok {
		keyMenu := fyne.NewMenuItem("Keys", func() {})
		keyMenu.ChildMenu = fyne.NewMenu("Key menu", menuItems...)

		m := fyne.NewMenu("SAM", keyMenu)
		desk.SetSystemTrayMenu(m)
	}
}
