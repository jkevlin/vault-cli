package command

import (
	"strings"

	"github.com/posener/complete"
)

type PutVaultNamespaceCommand struct {
	Meta
}

func (c *PutVaultNamespaceCommand) Help() string {
	helpText := `
Usage: nomad put vaultnamespace [options]
  Bootstrap is used to bootstrap the ACL system and get an initial token.
General Options:
  ` + generalOptionsUsage(usageOptsDefault|usageOptsNoNamespace) + `
`
	return strings.TrimSpace(helpText)
}

func (c *PutVaultNamespaceCommand) AutocompleteFlags() complete.Flags {
	return mergeAutocompleteFlags(c.Meta.AutocompleteFlags(FlagSetClient),
		complete.Flags{})
}

func (c *PutVaultNamespaceCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *PutVaultNamespaceCommand) Synopsis() string {
	return "Bootstrap the ACL system for initial token"
}

func (c *PutVaultNamespaceCommand) Name() string { return "acl bootstrap" }

func (c *PutVaultNamespaceCommand) Run(args []string) int {
	return 0
}
