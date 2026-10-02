//go:build darwin

package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
	"github.com/atotto/clipboard"
)

type keyEntry struct {
	key            *key.Key
	menuItem       *fyne.MenuItem
	bothItem       *fyne.MenuItem
	encryptionItem *fyne.MenuItem
	decryptionItem *fyne.MenuItem
}

func newKeyEntry(k *key.Key, onSelectMode func(k *key.Key, mode KeyMode)) *keyEntry {
	entry := &keyEntry{key: k}

	entry.bothItem = fyne.NewMenuItem("Encryption and decryption", func() {
		onSelectMode(k, ModeBoth)
	})
	entry.encryptionItem = fyne.NewMenuItem("Encryption", func() {
		onSelectMode(k, ModeEncryption)
	})
	entry.decryptionItem = fyne.NewMenuItem("Decryption", func() {
		onSelectMode(k, ModeDecryption)
	})

	entry.menuItem = fyne.NewMenuItem(k.Name, func() {})
	entry.menuItem.ChildMenu = fyne.NewMenu(
		"key options for "+k.Name,
		entry.bothItem,
		entry.encryptionItem,
		entry.decryptionItem,
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Copy key name", copyToClipboard(k.Name)),
		fyne.NewMenuItem("Copy public key", copyToClipboard(k.PublicKey)),
		fyne.NewMenuItem("Copy private key", copyToClipboard(k.PrivateKey)),
	)

	return entry
}

func copyToClipboard(content string) func() {
	return func() {
		if err := clipboard.WriteAll(content); err != nil {
			log.Printf("could not copy to clipboard: %v", err)
		}
	}
}

func (e *keyEntry) updateChecked(activeEncKey, activeDecKey string) {
	isEnc := e.key.Name == activeEncKey
	isDec := e.key.Name == activeDecKey

	e.bothItem.Checked = isEnc && isDec
	e.encryptionItem.Checked = isEnc && !isDec
	e.decryptionItem.Checked = !isEnc && isDec
}
