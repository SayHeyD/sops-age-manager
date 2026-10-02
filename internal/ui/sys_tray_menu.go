//go:build darwin

package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop" //nolint:typecheck
	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
)

func CreateSysTrayMenu(a fyne.App, keys []*key.Key, appConfig *config.Config, logo []byte) {
	desk, ok := a.(desktop.App)
	if !ok {
		return
	}

	tm := NewTrayManager(desk, keys, logo)
	tm.Refresh()
}
