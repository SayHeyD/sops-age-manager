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

Use the base command with '--' after which you can execute what you want. 
The sops configuration will be applied automatically.

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
		for index, arg := range args {
			if arg == "sops" {
				argsUntilSops := make([]string, len(args[:index+1]))
				argsAfterSops := make([]string, len(args[index+1:]))

				copy(argsUntilSops, args[:index+1])
				copy(argsAfterSops, args[index+1:])

				firstArgHalf := append(argsUntilSops, "--age", wantedEncryptionKey.PublicKey)
				args = append(firstArgHalf, argsAfterSops...)
			}
		}

	}

	if wantedDecryptionKey == nil {
		log.Printf("Could not find decryption key \"%s\"", appConfig.DecryptionKeyName)
	} else {
		err = os.Setenv("SOPS_AGE_KEY", wantedDecryptionKey.PrivateKey)
		if err != nil {
			log.Fatalf("could not set env variable: %v", err)
		}
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
