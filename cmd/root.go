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
	"runtime"
)

var (
	showVersion bool

	appVersion string

	RootCmd = &cobra.Command{
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

func Execute(version string) {
	appVersion = version

	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("executing rootCmd: %v", err)
	}
}

func init() {
	cobra.OnInitialize()

	RootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "Shows the current version of sam")

	RootCmd.AddCommand(keyCommands)
	RootCmd.AddCommand(configCommands)
}

func executeSops(args []string) {
	if showVersion {
		fmt.Printf("sam version: %s (%s)\n", appVersion, runtime.Version())
		return
	}

	appConfig, err := config.NewConfigFromFile()
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
		log.Printf("Could not find decryption key \"%s\"", appConfig.DecryptionKeyName)
	}

	var passThroughOut bytes.Buffer
	var passThroughErr bytes.Buffer

	sopsCmd := exec.Command(args[0], args[1:]...)

	sopsCmd.Stdout = &passThroughOut
	sopsCmd.Stderr = &passThroughErr

	err = sopsCmd.Run()
	if err != nil {
		fmt.Printf("sops error: %v: %s", err, passThroughErr.String())
		return
	}

	fmt.Print(passThroughOut.String())
}
