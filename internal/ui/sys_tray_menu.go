//go:build darwin || windows

package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop" //nolint:typecheck
	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
)

func CreateSysTrayMenu(a fyne.App, keys []*key.Key, appConfig *config.Config, logo []byte) *TrayManager {
	desk, ok := a.(desktop.App)
	if !ok {
		return nil
	}

	tm := NewTrayManager(desk, keys, logo)
	tm.Refresh()

	configPath := ""
	if appConfig != nil && appConfig.Path != "" {
		configPath = appConfig.Path
	} else {
		loadedCfg, err := config.NewConfigFromFile()
		if err == nil && loadedCfg != nil {
			configPath = loadedCfg.Path
		}
	}

	if configPath != "" {
		if err := tm.StartWatching(configPath); err != nil {
			log.Printf("could not start config watcher: %v", err)
		}
	}

	return tm
}
