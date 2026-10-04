//go:build !darwin && !windows
// +build !darwin,!windows

package cmd

import (
	"log"

	"github.com/SayHeyD/sops-age-manager/pkg/config"
)

func launchUI(appConfig *config.Config, appLogo []byte) {
	log.Fatal("UI is currently only supported on macOS and Windows")
}
