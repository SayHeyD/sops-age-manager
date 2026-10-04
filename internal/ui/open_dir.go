//go:build darwin || windows

package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
)

var openCommandFunc = func(dirPath string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", dirPath)
	case "windows":
		cmd = exec.Command("explorer", dirPath)
	default:
		cmd = exec.Command("xdg-open", dirPath)
	}
	return cmd.Start()
}

func OpenDirectory(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("could not create directory '%s': %w", dirPath, err)
		}
	}

	return openCommandFunc(dirPath)
}

func OpenConfigDirectory() error {
	dirPath, err := config.GetConfigDirPath()
	if err != nil {
		return err
	}
	return OpenDirectory(dirPath)
}

func OpenKeyDirectory(customKeyDir string) error {
	dirPath, err := key.GetKeyDirPath(customKeyDir)
	if err != nil {
		return err
	}
	return OpenDirectory(dirPath)
}
