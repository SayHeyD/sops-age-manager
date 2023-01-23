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
}
