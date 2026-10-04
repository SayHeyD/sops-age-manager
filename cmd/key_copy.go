package cmd

import (
	"fmt"
	"log"

	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

var (
	copyPublicKey  bool
	copyPrivateKey bool
	copyKeyName    bool

	clipboardWriteAllFunc = clipboard.WriteAll

	copyKeyCommand = &cobra.Command{
		Use:   "copy",
		Short: "Copy key attributes (public key, private key, or name) to clipboard",
		Long: `Copies the recipient public key (age1...), private key, or key name 
of the specified key to the system clipboard. By default, the public key is copied.`,
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completeKeyNames,
		Run: func(cmd *cobra.Command, args []string) {
			copyKey(args[0])
		},
	}
)

func init() {
	usage := `Usage:
  sam key copy <KEY_NAME> [flags]

Arguments:
  KEY_NAME     key to copy attributes from, required

Flags:
  -p, --public    copies the public key (default)
      --private   copies the private key
  -n, --name      copies the key name
  -h, --help      help for copy
`

	copyKeyCommand.PersistentFlags().BoolVarP(&copyPublicKey, "public", "p", false, "copies the public key (recipient)")
	copyKeyCommand.PersistentFlags().BoolVar(&copyPrivateKey, "private", false, "copies the private key")
	copyKeyCommand.PersistentFlags().BoolVarP(&copyKeyName, "name", "n", false, "copies the key name")

	copyKeyCommand.SetUsageTemplate(usage)
}

func copyKey(keyName string) {
	var keyDir string
	if appConfig, err := config.NewConfigFromFile(); err == nil {
		keyDir = appConfig.KeyDir
	}

	keys := key.GetAvailableKeys(keyDir)

	for _, ageKey := range keys {
		if ageKey.Name == keyName {
			contentToCopy := ageKey.PublicKey
			attributeName := "public key"

			if copyPrivateKey {
				contentToCopy = ageKey.PrivateKey
				attributeName = "private key"
			} else if copyKeyName {
				contentToCopy = ageKey.Name
				attributeName = "key name"
			}

			if err := clipboardWriteAllFunc(contentToCopy); err != nil {
				log.Fatalf("could not copy to clipboard: %v", err)
			}

			fmt.Printf("Copied %s of \"%s\" to clipboard\n", attributeName, ageKey.Name)
			return
		}
	}

	log.Fatalf("No key with name \"%s\" found", keyName)
}
