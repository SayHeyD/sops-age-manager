//go:build !darwin

package ui

import (
	"fyne.io/fyne/v2"
)

func setupAppLifecycle(a fyne.App) {
	// No-op on non-macOS platforms
}
