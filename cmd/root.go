package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var (
	rootCmd = &cobra.Command{
		Use:   "sam",
		Short: "Sops-Age-Manager (SAM) is a tool for managing multiple age keys when using mozilla/sops",
		Long: `Sops-Age-Manager (SAM) is a tool for managing the age key used by sops.
This wrapper for sops should provide key selection by name, rather than
by using the private or public key.

GitHub: https://github.com/SayHeyD/sops-age-manager`,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("")
	}
}

func init() {
	cobra.OnInitialize()

	rootCmd.AddCommand(keyCommands)

	/*rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")

	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(initCmd)*/
}
