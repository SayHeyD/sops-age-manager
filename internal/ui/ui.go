//go:build darwin

package ui

import (
	"log"
	"time"

	"github.com/SayHeyD/sops-age-manager/pkg/key"

	"fyne.io/fyne/v2/app"
	"github.com/SayHeyD/sops-age-manager/pkg/config"
)

/* // Hide macos dock application when app is started
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

int
SetActivationPolicy(void) {
    [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
    return 0;
}
*/
import "C" //nolint:typecheck

func setActivationPolicy() {
	C.SetActivationPolicy()
}

func Init(config *config.Config, logo []byte) {
	lockFilePath := getLockFilePath(config)
	lock, err := AcquireInstanceLock(lockFilePath)
	if err != nil {
		log.Printf("SAM is already running: %v", err)
		return
	}
	defer lock.Release()

	a := app.New()

	keys := key.GetAvailableKeys(config.KeyDir)

	CreateSysTrayMenu(a, keys, config, logo)

	a.Lifecycle().SetOnStarted(func() {
		go func() {
			time.Sleep(200 * time.Millisecond)
			setActivationPolicy()
		}()
	})

	a.Run()
}
