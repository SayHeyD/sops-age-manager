package cmd

import (
	"fmt"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
	"github.com/spf13/cobra"
)

var (
	listKeyWithPrivateKeys = false
	listKeyWithFileName    = false

	listKeyCommand = &cobra.Command{
		Use:   "list",
		Short: "Lists all available sops keys",
		Long: `List all available sops keys. By default only name and public key are listed.
Filepath and private key can be displayed optionally`,
		Run: func(cmd *cobra.Command, args []string) {
			listKeys()
		},
	}
)

func init() {
	listKeyCommand.Flags().BoolVarP(&listKeyWithPrivateKeys, "private-key", "p", false, "displays the private-key")
	listKeyCommand.Flags().BoolVarP(&listKeyWithFileName, "file", "f", false, "displays the filepath of the key")
}

func listKeys() {
	keys := key.GetAvailableKeys("")

	for _, ageKey := range keys {
		fmt.Println()

		fmt.Printf("%s\n", ageKey.Name)

		fmt.Println("║")

		if listKeyWithFileName {
			fmt.Printf("╠═ File path: %s\n", ageKey.FileName)
		}

		if listKeyWithPrivateKeys {
			fmt.Printf("╠═ Private Key: %s\n", ageKey.PrivateKey)
		}

		fmt.Printf("╚═ Public Key: %s\n", ageKey.PublicKey)
	}
}
