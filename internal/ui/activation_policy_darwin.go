//go:build darwin

package ui

import (
	"time"

	"fyne.io/fyne/v2"
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

func setupAppLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		go func() {
			time.Sleep(200 * time.Millisecond)
			C.SetActivationPolicy()
		}()
	})
}
