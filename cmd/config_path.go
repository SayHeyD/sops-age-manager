package cmd

import (
	"fmt"
	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/spf13/cobra"
	"log"
)

var (
	configPathCommands = &cobra.Command{
		Use:   "path",
		Short: "Returns the path of the current config",
		Long: `Prints out the path of the currently loaded application configuration.
Stored in the SOPS_AGE_MANAGER_CONFIG_DIR environment variable. If the variable is empty the
default path of "$HOME/.sops-age-manager/config.yaml" will be used.`,
		Run: func(cmd *cobra.Command, args []string) {
			configPath()
		},
	}
)

func configPath() {
	appConfig, err := config.NewConfigFromFile()
	if err != nil {
		log.Fatalf("print config path: %v", err)
	}
	fmt.Printf("Current application config path: \"%s\"", appConfig.Path)
}
