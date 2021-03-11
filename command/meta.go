package command

import (
	"flag"
	"strings"

	"github.com/jkevlin/vault-cli/pkg/config"
	"github.com/jkevlin/vault-cli/pkg/configservice"
	"github.com/jkevlin/vault-cli/pkg/secretservice"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// FlagSetFlags is an enum to define what flags are present in the
// default FlagSet returned by Meta.FlagSet.
type FlagSetFlags uint

// Meta contains the meta-options and functionality that nearly every
// Nomad command inherits.
type Meta struct {
	Ui cli.Ui

	Config        *config.Config
	ConfigService configservice.ConfigService

	SecretService secretservice.SecretService

	ConfigPath string

	// // These are set by the command line flags.
	// context is the context name to use for this command
	context string
	// Whether to not-colorize output
	noColor bool

	// namespace to send API requests
	namespace    string
	name         string
	outputFormat string
}

// FlagSet returns a FlagSet with the common flags that every
// command implements. The exact behavior of FlagSet can be configured
// using the flags as the second parameter, for example to disable
// server settings on the commands that don't talk to a server.
func (m *Meta) FlagSet(n string) *flag.FlagSet {
	f := flag.NewFlagSet(n, flag.ContinueOnError)

	f.StringVar(&m.context, "c", "local", "")
	f.StringVar(&m.context, "context", "local", "")
	f.StringVar(&m.namespace, "n", "", "")
	f.StringVar(&m.namespace, "namespace", "", "")
	f.StringVar(&m.outputFormat, "o", "", "")
	f.StringVar(&m.outputFormat, "output", "", "")

	f.SetOutput(&uiErrorWriter{ui: m.Ui})

	return f
}

// AutocompleteFlags returns a set of flag completions for the given flag set.
func (m *Meta) AutocompleteFlags() complete.Flags {
	return complete.Flags{
		"-c":             complete.PredictAnything,
		"-context":       complete.PredictAnything,
		"-kubecliconfig": complete.PredictAnything,
		"-n":             complete.PredictAnything,
		"-namespace":     complete.PredictAnything,
		"-no-color":      complete.PredictNothing,
	}
}

// generalOptionsUsage returns the help string for the global options.
func generalOptionsUsage() string {

	helpText := `
  -context=<contextname>
    The name of the context to use for this run of the command
    Alias: -c
  -kubecliconfig=<kubecliconfig>
    The location if the cli config yaml file.

  -namespace=<namespace>
    The target namespace for queries and actions bound to a namespace.
    Overrides the VAULT_CLI_NAMESPACE environment variable if set.
    If set to '*', job and alloc subcommands query all namespaces authorized
    to user.
    Defaults to the "default" namespace.

  -no-color
    Disables colored command output. Alternatively, VAULT_CLI_NO_COLOR may be
    set.

  -output=<json|yaml|text>
    Alias: -o
`
	return strings.TrimSpace(helpText)
}

// funcVar is a type of flag that accepts a function that is the string given
// by the user.
type funcVar func(s string) error

func (f funcVar) Set(s string) error { return f(s) }
func (f funcVar) String() string     { return "" }
func (f funcVar) IsBoolFlag() bool   { return false }
