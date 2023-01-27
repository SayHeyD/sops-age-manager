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

	var wantedEncryptionKey *key.Key
	var wantedDecryptionKey *key.Key

	keys := key.GetAvailableKeys("")

	for _, foundKey := range keys {
		if appConfig.EncryptionKeyName == foundKey.Name {
			wantedEncryptionKey = foundKey
		}

		if appConfig.DecryptionKeyName == foundKey.Name {
			wantedDecryptionKey = foundKey
		}
	}

	if wantedEncryptionKey == nil {
		log.Printf("Could not find encryption key \"%s\"", appConfig.EncryptionKeyName)
	} else {
		err = os.Setenv("SOPS_AGE_KEY", wantedEncryptionKey.PrivateKey)
		if err != nil {
			log.Fatalf("could not set env variable: %v", err)
		}
	}

	if wantedDecryptionKey == nil {
		log.Printf("Could not find encryption key \"%s\"", appConfig.DecryptionKeyName)
	} else {
		args = append([]string{"--age", wantedDecryptionKey.PublicKey}, args...)
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
