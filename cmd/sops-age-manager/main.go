package main

import (
	"flag"
	"fmt"
	"github.com/SayHeyD/sops-age-manage/internal/key"
)

func main() {
	systemTrayIconFlagSet := false
	systemTrayIconFlagName := "tray-icon"

	systemTrayIcon := flag.Bool(systemTrayIconFlagName, false, "Add a tray icon for managing keys with a UI")
	activeKey := flag.String("key", "", "The key to select for sops")

	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		if f.Name == systemTrayIconFlagName {
			systemTrayIconFlagSet = true
		}
	})

	if systemTrayIconFlagSet {
		fmt.Println(*systemTrayIcon)
	}

	if *activeKey != "" {
		fmt.Println(*activeKey)
	}

	key.GetAvailableKeys()
}
