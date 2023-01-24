package cmd

import (
	"github.com/spf13/cobra"
)

var (
	configCommands = &cobra.Command{
		Use:   "config",
		Short: "Commands regarding the application configuration",
		Long:  `Commands to show the application config or the path to it`,
	}
)

func init() {
	configCommands.AddCommand(configDumpCommands)
}
