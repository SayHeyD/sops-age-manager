package cmd

import (
	"fmt"
	"github.com/SayHeyD/sops-age-manager/pkg/config"
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
		Run: func(cmd *cobra.Command, args []string) {
			executeSops(args)
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("executing rootCmd: %v", err)
	}
}

func init() {
	cobra.OnInitialize()

	rootCmd.AddCommand(keyCommands)
	rootCmd.AddCommand(configCommands)
}

func executeSops(args []string) {
	appConfig, err := config.NewConfigFromFile("")
	if err != nil {
		log.Fatalf("execute sops: %v", err)
	}

	fmt.Printf("Args: %s\nKeyName: %s\n", args, appConfig.KeyName)
	//out, err := exec.Command("sops", args...).Output()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(string(out))
}
