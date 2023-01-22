package main

import (
	"flag"
	"fmt"
	"github.com/SayHeyD/sops-age-manage/internal/key"
	"log"
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

	if *activeKey == "" {
		log.Fatal("No key name provided")
	}

	keys := key.GetAvailableKeys()

	for _, ageKey := range keys {
		if ageKey.Name == *activeKey {
			ageKey.SetActive()
			fmt.Println(fmt.Sprintf("Set \"%s\" as active key", ageKey.Name))
			return
		}
	}

	log.Fatal(fmt.Sprintf("No key with name \"%s\" found", *activeKey))
}
