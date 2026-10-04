package cmd

import (
	"fmt"

	"github.com/SayHeyD/sops-age-manager/pkg/key"
	"github.com/spf13/cobra"
)

var (
	clearEncryptionKey bool
	clearDecryptionKey bool

	clearKeyCommand = &cobra.Command{
		Use:   "clear",
		Short: "Clears the active encryption and/or decryption keys",
		Long: `Clears the currently active age keys. When keys are cleared, SAM operates in pass-through 
mode, allowing sops to use its default key configuration without injecting SAM keys.`,
		Run: func(cmd *cobra.Command, args []string) {
			clearActiveKeys()
		},
	}
)

func init() {
	usage := `Usage:
  sam key clear [flags]

Flags:
  -e, --encryption   clears the active encryption key only
  -d, --decryption   clears the active decryption key only
  -h, --help         help for clear
`

	clearKeyCommand.PersistentFlags().BoolVarP(&clearEncryptionKey, "encryption", "e", false, "clears the active encryption key")
	clearKeyCommand.PersistentFlags().BoolVarP(&clearDecryptionKey, "decryption", "d", false, "clears the active decryption key")

	clearKeyCommand.SetUsageTemplate(usage)
}

func clearActiveKeys() {
	clearedDecryptionKeyTemplate := "Cleared active decryption key\n"
	clearedEncryptionKeyTemplate := "Cleared active encryption key\n"

	if !clearDecryptionKey && !clearEncryptionKey {
		key.ClearActiveKeys()
		fmt.Printf(clearedDecryptionKeyTemplate)
		fmt.Printf(clearedEncryptionKeyTemplate)
	} else {
		if clearDecryptionKey {
			key.ClearActiveDecryption()
			fmt.Printf(clearedDecryptionKeyTemplate)
		}

		if clearEncryptionKey {
			key.ClearActiveEncryption()
			fmt.Printf(clearedEncryptionKeyTemplate)
		}
	}
}
