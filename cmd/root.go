package cmd

import (
	"bytes"
	"fmt"
	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
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

	var wantedKey *key.Key
	keys := key.GetAvailableKeys("")
	for _, foundKey := range keys {
		if appConfig.KeyName == foundKey.Name {
			wantedKey = foundKey
			break
		}
	}

	if wantedKey == nil {
		log.Fatalf("Could not find key \"%s\"", appConfig.KeyName)
	}

	args = append([]string{"--age", wantedKey.PublicKey}, args...)
	err = os.Setenv("SOPS_AGE_KEY", wantedKey.PrivateKey)
	if err != nil {
		log.Fatalf("could not set env variable: %v", err)
	}

	fmt.Printf("sops %s\n", args)

	var sopsOut bytes.Buffer
	var stderr bytes.Buffer

	sopsCmd := exec.Command("sops", args...)

	sopsCmd.Stdout = &sopsOut
	sopsCmd.Stderr = &stderr

	err = sopsCmd.Run()
	if err != nil {
		fmt.Printf("sops error: %v: %s", err, stderr.String())
		return
	}

	fmt.Print(sopsOut.String())
}
