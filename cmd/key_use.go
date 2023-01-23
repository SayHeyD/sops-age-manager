package cmd

import (
	"fmt"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
	"github.com/spf13/cobra"
	"log"
)

var (
	useKeyCommand = &cobra.Command{
		Use:   "use",
		Short: "Use the key with the given name",
		Long: `Uses the key given by name. The age key will be used for operations
performed with sops f.ex. decrypting and encrypting files.`,
		Run: func(cmd *cobra.Command, args []string) {
			setActiveKey(args[0])
		},
	}
)

func init() {
	usage := `Usage:
  sam key <KEY_NAME> [flags]

Arguments:
  KEY_NAME     key to use for sops commands, required

Flags:
  -h, --help   help for key
`

	useKeyCommand.SetUsageTemplate(usage)
}

func setActiveKey(keyName string) {
	keys := key.GetAvailableKeys("")

	for _, ageKey := range keys {
		if ageKey.Name == keyName {
			ageKey.SetActive()
			fmt.Println(fmt.Sprintf("Set \"%s\" as active key", ageKey.Name))
			return
		}
	}

	log.Fatal(fmt.Sprintf("No key with name \"%s\" found", keyName))
}
