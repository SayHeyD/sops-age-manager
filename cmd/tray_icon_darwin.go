//go:build darwin
// +build darwin

package cmd

import (
	"github.com/SayHeyD/sops-age-manager/internal/ui"
	"github.com/SayHeyD/sops-age-manager/pkg/config"
)

func launchUI(appConfig *config.Config, appLogo []byte) {
	ui.Init(appConfig, appLogo)
}
