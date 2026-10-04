//go:build windows
// +build windows

package cmd

import (
	"unsafe"

	"github.com/SayHeyD/sops-age-manager/internal/ui"
	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"golang.org/x/sys/windows"
)

var (
	kernel32              = windows.NewLazySystemDLL("kernel32.dll")
	user32                = windows.NewLazySystemDLL("user32.dll")
	getConsoleWindow      = kernel32.NewProc("GetConsoleWindow")
	getConsoleProcessList = kernel32.NewProc("GetConsoleProcessList")
	showWindow            = user32.NewProc("ShowWindow")
	freeConsole           = kernel32.NewProc("FreeConsole")
)

func launchUI(appConfig *config.Config, appLogo []byte) {
	hideConsoleIfOwned()
	ui.Init(appConfig, appLogo)
}

func hideConsoleIfOwned() {
	hwnd, _, _ := getConsoleWindow.Call()
	if hwnd == 0 {
		return
	}

	var pids [2]uint32
	count, _, _ := getConsoleProcessList.Call(uintptr(unsafe.Pointer(&pids[0])), 2)
	// If only 1 process is attached to this console, it was created specifically for this process
	// (e.g. launched by double-clicking in Windows Explorer). In this case, hide and free the console.
	if count == 1 {
		showWindow.Call(hwnd, uintptr(windows.SW_HIDE))
		freeConsole.Call()
	}
}
