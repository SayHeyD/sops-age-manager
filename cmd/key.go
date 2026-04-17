package cmd

import (
	"github.com/spf13/cobra"
)

var (
	keyCommands = &cobra.Command{
		Use:   "key",
		Short: "Commands for managing the detected age keys",
		Long:  `Commands for choosing and listing the available age keys to use with sops.`,
	}
)

func init() {
	keyCommands.AddCommand(listKeyCommand)
	keyCommands.AddCommand(useKeyCommand)
}
