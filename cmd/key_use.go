package cmd

import (
	"fmt"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
	"github.com/spf13/cobra"
	"log"
)

var (
	setEncryptionKey bool
	setDecryptionKey bool

	useKeyCommand = &cobra.Command{
		Use:   "use",
		Short: "Use the key with the given name",
		Long: `Uses the key given by name. The age key will be used for operations
performed with sops f.ex. decrypting and encrypting files. Decryption and encryption 
keys can be set independently. Not specifying any flags will set the key for both
decryption and encryption.`,
		Run: func(cmd *cobra.Command, args []string) {
			setActiveKey(args[0])
		},
	}
)

func init() {
	usage := `Usage:
  sam key use <KEY_NAME> [flags]

Arguments:
  KEY_NAME     key to use for sops commands, required

Flags:
  -e, --encryption   sets the key for encryption only
  -d, --decryption   sets the key for decryption only
  -h, --help         help for key
`

	useKeyCommand.PersistentFlags().BoolVarP(&setEncryptionKey, "encryption", "e", false, "")
	useKeyCommand.PersistentFlags().BoolVarP(&setDecryptionKey, "decryption", "d", false, "")

	useKeyCommand.SetUsageTemplate(usage)
}

func setActiveKey(keyName string) {
	keys := key.GetAvailableKeys("")

	for _, ageKey := range keys {
		if ageKey.Name == keyName {
			if setDecryptionKey || (!setDecryptionKey && !setEncryptionKey) {
				ageKey.SetActiveDecryption()
				fmt.Printf("Set \"%s\" as active decryption key\n", ageKey.Name)
			}

			if setEncryptionKey || (!setDecryptionKey && !setEncryptionKey) {
				ageKey.SetActiveEncryption()
				fmt.Printf("Set \"%s\" as active encryption key\n", ageKey.Name)
			}

			return
		}
	}

	log.Fatalf("No key with name \"%s\" found", keyName)
}
