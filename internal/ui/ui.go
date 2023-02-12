package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
	"github.com/atotto/clipboard"
	"log"
)

func Init(config *config.Config) {
	a := app.New()

	keys := key.GetAvailableKeys("")

	var menuItems []*fyne.MenuItem

	for _, ageKey := range keys {

		ageKeyEncryptionDecryptionMenuEntry := fyne.NewMenuItem("Encryption and decryption", func() {
			fmt.Println(ageKey.Name)
			ageKey.SetActiveDecryption()
			ageKey.SetActiveEncryption()
		})

		ageKeyEncryptionMenuEntry := fyne.NewMenuItem("Encryption", func() {
			ageKey.SetActiveEncryption()
		})

		ageKeyDecryptionMenuEntry := fyne.NewMenuItem("Decryption", func() {
			ageKey.SetActiveDecryption()
		})

		fmt.Println(config.DecryptionKeyName)

		if config.DecryptionKeyName == ageKey.Name && config.EncryptionKeyName == ageKey.Name {
			ageKeyEncryptionDecryptionMenuEntry.Checked = true
		}

		if config.EncryptionKeyName == ageKey.Name && config.DecryptionKeyName != ageKey.Name {
			ageKeyEncryptionDecryptionMenuEntry.Checked = true
		}

		if config.EncryptionKeyName != ageKey.Name && config.DecryptionKeyName == ageKey.Name {
			ageKeyEncryptionDecryptionMenuEntry.Checked = true
		}

		ageKeyMenu := fyne.NewMenuItem(ageKey.Name, func() {})
		ageKeyMenu.ChildMenu = fyne.NewMenu("key options for "+ageKey.Name,
			ageKeyEncryptionDecryptionMenuEntry,
			ageKeyEncryptionMenuEntry,
			ageKeyDecryptionMenuEntry,
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Copy key name", func() {
				err := clipboard.WriteAll(ageKey.Name)
				if err != nil {
					log.Fatalf("could not set clipboard text: %v", err)
				}
			}),
			fyne.NewMenuItem("Copy public key", func() {
				err := clipboard.WriteAll(ageKey.PublicKey)
				if err != nil {
					log.Fatalf("could not set clipboard text: %v", err)
				}
			}),
			fyne.NewMenuItem("Copy private key", func() {
				err := clipboard.WriteAll(ageKey.PrivateKey)
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

	a.Run()
}
