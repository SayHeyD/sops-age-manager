package cmd

import (
	"fmt"
	"strings"

	"github.com/SayHeyD/sops-age-manager/pkg/config"
	"github.com/SayHeyD/sops-age-manager/pkg/key"
	"github.com/spf13/cobra"
)

var (
	completionCommand = &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script for the specified shell",
		Long: `Generate shell completion scripts for SAM (sops-age-manager).

To load completions in your current shell session:

Bash:
  $ source <(sam completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ sam completion bash > /etc/bash_completion.d/sam
  # macOS:
  $ sam completion bash > $(brew --prefix)/etc/bash_completion.d/sam

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ sam completion zsh > "${fpath[1]}/_sam"

  # You will need to start a new shell for this setup to take effect.

Fish:
  $ sam completion fish | source

  # To load completions for each session, execute once:
  $ sam completion fish > ~/.config/fish/completions/sam.fish

PowerShell:
  PS> sam completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, add to your $PROFILE:
  PS> Add-Content -Path $PROFILE -Value "sam completion powershell | Out-String | Invoke-Expression"
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch strings.ToLower(args[0]) {
			case "bash":
				return cmd.Root().GenBashCompletionV2(cmd.OutOrStdout(), true)
			case "zsh":
				return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
			case "fish":
				return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
			default:
				return fmt.Errorf("unsupported shell type %q, supported shells are: bash, zsh, fish, powershell", args[0])
			}
		},
	}
)

func completeKeyNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	cfg, err := config.NewConfigFromFile()
	var keyDir string
	if err == nil && cfg != nil {
		keyDir = cfg.KeyDir
	}

	keys, err := key.FindAvailableKeys(keyDir)
	if err != nil || len(keys) == 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var names []string
	for _, k := range keys {
		names = append(names, k.Name)
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
