package cmd

import (
	"fmt"
	"log"

	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
	"github.com/spf13/cobra"
)

var (
	createForce         bool
	createCopyPublic    bool
	createSetEncryption bool
	createSetDecryption bool
	createSetBoth       bool

	createKeyCommand = &cobra.Command{
		Use:     "create <KEY_NAME>",
		Aliases: []string{"generate", "new"},
		Short:   "Create a new age key using age-keygen",
		Long: `Creates a new age key file with the specified name in the configured key directory 
using the pre-installed 'age-keygen' utility. 

Key files are saved with a '.txt' extension (e.g. <KEY_NAME>.txt) and standard 0600 file permissions. 
Nested key names like 'k8s/production' are placed in the appropriate subdirectories.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			createKey(args[0])
		},
	}
)

func init() {
	usage := `Usage:
  sam key create <KEY_NAME> [flags]

Aliases:
  create, generate, new

Arguments:
  KEY_NAME     name of the key to create, required

Flags:
  -f, --force        overwrite existing key file if already exists
  -c, --copy         copy the public key (recipient) to clipboard
  -e, --encryption   set newly created key as active encryption key
  -d, --decryption   set newly created key as active decryption key
  -u, --use          set newly created key as active key for both encryption and decryption
  -h, --help         help for create
`

	createKeyCommand.PersistentFlags().BoolVarP(&createForce, "force", "f", false, "overwrite existing key file if already exists")
	createKeyCommand.PersistentFlags().BoolVarP(&createCopyPublic, "copy", "c", false, "copy public key recipient to clipboard")
	createKeyCommand.PersistentFlags().BoolVarP(&createSetEncryption, "encryption", "e", false, "set newly created key as active encryption key")
	createKeyCommand.PersistentFlags().BoolVarP(&createSetDecryption, "decryption", "d", false, "set newly created key as active decryption key")
	createKeyCommand.PersistentFlags().BoolVarP(&createSetBoth, "use", "u", false, "set newly created key as active for both encryption and decryption")

	createKeyCommand.SetUsageTemplate(usage)
}

func createKey(keyName string) {
	var keyDir string
	if appConfig, err := config.NewConfigFromFile(); err == nil && appConfig != nil {
		keyDir = appConfig.KeyDir
	}

	createdKey, err := key.CreateKey(keyName, keyDir, createForce)
	if err != nil {
		log.Fatalf("failed to create key: %v", err)
	}

	fmt.Printf("Created key \"%s\" (%s)\n", createdKey.Name, createdKey.FileName)
	fmt.Printf("Public key: %s\n", createdKey.PublicKey)

	if createCopyPublic {
		if err := clipboardWriteAllFunc(createdKey.PublicKey); err != nil {
			log.Fatalf("could not copy public key to clipboard: %v", err)
		}
		fmt.Printf("Copied public key of \"%s\" to clipboard\n", createdKey.Name)
	}

	if createSetBoth {
		createdKey.SetActiveBoth()
		fmt.Printf("Set \"%s\" as active decryption key\n", createdKey.Name)
		fmt.Printf("Set \"%s\" as active encryption key\n", createdKey.Name)
	} else {
		if createSetDecryption {
			createdKey.SetActiveDecryption()
			fmt.Printf("Set \"%s\" as active decryption key\n", createdKey.Name)
		}
		if createSetEncryption {
			createdKey.SetActiveEncryption()
			fmt.Printf("Set \"%s\" as active encryption key\n", createdKey.Name)
		}
	}
}
