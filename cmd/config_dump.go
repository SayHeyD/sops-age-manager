package cmd

import (
	"fmt"
	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/spf13/cobra"
	"log"
)

var (
	configDumpCommands = &cobra.Command{
		Use:   "dump",
		Short: "Return the application configuration",
		Long:  `Prints out the complete application configuration file into the console`,
		Run: func(cmd *cobra.Command, args []string) {
			dumpConfig()
		},
	}
)

func dumpConfig() {
	appConfig, err := config.NewConfigFromFile()
	if err != nil {
		log.Fatalf("dump config: %v", err)
	}

	rawConfig, err := appConfig.Raw()
	if err != nil {
		log.Fatalf("dump raw config: %v", err)
	}
	fmt.Println()
	fmt.Println(rawConfig)
}
