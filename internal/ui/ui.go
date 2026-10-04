//go:build darwin || windows

package ui

import (
	"log"

	"fyne.io/fyne/v2/app"
	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
)

func Init(config *config.Config, logo []byte) {
	lockFilePath := getLockFilePath(config)
	lock, err := AcquireInstanceLock(lockFilePath)
	if err != nil {
		log.Printf("SAM is already running: %v", err)
		return
	}
	defer lock.Release()

	a := app.New()

	var keyDir string
	if config != nil {
		keyDir = config.KeyDir
	}

	keys, err := key.FindAvailableKeys(keyDir)
	if err != nil {
		log.Printf("could not read key files: %v", err)
	}

	CreateSysTrayMenu(a, keys, config, logo)

	setupAppLifecycle(a)

	a.Run()
}
