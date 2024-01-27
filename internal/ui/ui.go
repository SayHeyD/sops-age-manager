package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop" //nolint:typecheck
	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
	"github.com/atotto/clipboard"
	"log"
	"time"
)

/** Hide macos dock application when app is started
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

int
SetActivationPolicy(void) {
    [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
    return 0;
}
*/
import "C"

func setActivationPolicy() {
	C.SetActivationPolicy()
}

func Init(config *config.Config) {
	a := app.New()

	keys := key.GetAvailableKeys("")

	var menuItems []*fyne.MenuItem

	for _, ageKey := range keys {

		menuKey := *ageKey

		ageKeyEncryptionDecryptionMenuEntry := fyne.NewMenuItem("Encryption and decryption", func() {
			menuKey.SetActiveDecryption()
			menuKey.SetActiveEncryption()
		})

		ageKeyEncryptionMenuEntry := fyne.NewMenuItem("Encryption", func() {
			menuKey.SetActiveEncryption()
		})

		ageKeyDecryptionMenuEntry := fyne.NewMenuItem("Decryption", func() {
			menuKey.SetActiveDecryption()
		})

		if config.DecryptionKeyName == menuKey.Name && config.EncryptionKeyName == menuKey.Name {
			ageKeyEncryptionDecryptionMenuEntry.Checked = true
		}

		if config.EncryptionKeyName == menuKey.Name && config.DecryptionKeyName != menuKey.Name {
			ageKeyEncryptionDecryptionMenuEntry.Checked = true
		}

		if config.EncryptionKeyName != menuKey.Name && config.DecryptionKeyName == menuKey.Name {
			ageKeyEncryptionDecryptionMenuEntry.Checked = true
		}

		ageKeyMenu := fyne.NewMenuItem(menuKey.Name, func() {})
		ageKeyMenu.ChildMenu = fyne.NewMenu("key options for "+menuKey.Name,
			ageKeyEncryptionDecryptionMenuEntry,
			ageKeyEncryptionMenuEntry,
			ageKeyDecryptionMenuEntry,
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Copy key name", func() {
				err := clipboard.WriteAll(menuKey.Name)
				if err != nil {
					log.Fatalf("could not set clipboard text: %v", err)
				}
			}),
			fyne.NewMenuItem("Copy public key", func() {
				err := clipboard.WriteAll(menuKey.PublicKey)
				if err != nil {
					log.Fatalf("could not set clipboard text: %v", err)
				}
			}),
			fyne.NewMenuItem("Copy private key", func() {
				err := clipboard.WriteAll(menuKey.PrivateKey)
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

	a.Lifecycle().SetOnStarted(func() {
		go func() {
			time.Sleep(200 * time.Millisecond)
			setActivationPolicy()
		}()
	})

	a.Run()
}
