package cmd

import (
	"bufio"
	"fmt"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	listKeyWithPrivateKeys     = false
	listKeyWithFileName        = false
	skipPrivateKeyConfirmation = false

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
	listKeyCommand.Flags().BoolVarP(&listKeyWithPrivateKeys, "private-key", "p", false,
		"displays the private-key")
	listKeyCommand.Flags().BoolVarP(&listKeyWithFileName, "file", "f", false,
		"displays the filepath of the key")
	listKeyCommand.Flags().BoolVarP(&skipPrivateKeyConfirmation, "yes", "y", false,
		"automatically confirm that you want private keys to be shown")
}

func listKeys() {
	keys := key.GetAvailableKeys("")

	if listKeyWithPrivateKeys {
		if !skipPrivateKeyConfirmation {
			fmt.Printf("Please confirm that you want to display the pricate keys (yes/y): ")

			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("could not retrieve user input: ", err)
				return
			}

			input = strings.TrimSuffix(input, "\n")

			normalizedUserInput := strings.ToLower(input)

			if normalizedUserInput != "yes" && normalizedUserInput != "y" {
				fmt.Println("Question not answered with \"yes\" or \"y\". Aborting ...")
				return
			}
		}
	}

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
