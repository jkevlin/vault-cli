package command

import (
	"flag"
	"strings"

	"github.com/jkevlin/vault-cli/pkg/config"
	"github.com/jkevlin/vault-cli/pkg/secretservice"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// FlagSetFlags is an enum to define what flags are present in the
// default FlagSet returned by Meta.FlagSet.
type FlagSetFlags uint

const (
	FlagSetNone    FlagSetFlags = 0
	FlagSetClient  FlagSetFlags = 1 << iota
	FlagSetDefault              = FlagSetClient
)

// Meta contains the meta-options and functionality that nearly every
// Nomad command inherits.
type Meta struct {
	Ui cli.Ui

	Config *config.Config

	SecretService *secretservice.SecretService

	// // These are set by the command line flags.
	// context is the context name to use for this command
	context string
	// Whether to not-colorize output
	noColor bool

	// namespace to send API requests
	namespace string
}

// FlagSet returns a FlagSet with the common flags that every
// command implements. The exact behavior of FlagSet can be configured
// using the flags as the second parameter, for example to disable
// server settings on the commands that don't talk to a server.
func (m *Meta) FlagSet(n string, fs FlagSetFlags) *flag.FlagSet {
	f := flag.NewFlagSet(n, flag.ContinueOnError)

	// FlagSetClient is used to enable the settings for specifying
	// client connectivity options.
	if fs&FlagSetClient != 0 {
		f.StringVar(&m.namespace, "namespace", "", "")
		f.StringVar(&m.context, "context", "", "")

	}

	f.SetOutput(&uiErrorWriter{ui: m.Ui})

	return f
}

// AutocompleteFlags returns a set of flag completions for the given flag set.
func (m *Meta) AutocompleteFlags(fs FlagSetFlags) complete.Flags {
	if fs&FlagSetClient == 0 {
		return nil
	}

	return complete.Flags{
		"-kubecliconfig": complete.PredictAnything,
		"-namespace":     complete.PredictAnything,
		//"-namespace":       NamespacePredictor(m.Client, nil),
		"-no-color": complete.PredictNothing,
		// "-ca-cert":         complete.PredictFiles("*"),
		// "-ca-path":         complete.PredictDirs("*"),
		// "-client-cert":     complete.PredictFiles("*"),
		// "-client-key":      complete.PredictFiles("*"),
		// "-insecure":        complete.PredictNothing,
		// "-tls-server-name": complete.PredictNothing,
		// "-tls-skip-verify": complete.PredictNothing,
		"-c": complete.PredictAnything,
	}
}

type usageOptsFlags uint8

const (
	usageOptsDefault     usageOptsFlags = 0
	usageOptsNoNamespace                = 1 << iota
)

// generalOptionsUsage returns the help string for the global options.
func generalOptionsUsage(usageOpts usageOptsFlags) string {

	helpText := `
  -kubecliconfig=<kubecliconfig>
    The location if the cli config yaml file.
`

	namespaceText := `
  -namespace=<namespace>
    The target namespace for queries and actions bound to a namespace.
    Overrides the VAULT_CLI_NAMESPACE environment variable if set.
    If set to '*', job and alloc subcommands query all namespaces authorized
    to user.
    Defaults to the "default" namespace.
`

	// note: that although very few commands use color explicitly, all of them
	// return red-colored text on error so we don't want to make this
	// configurable
	remainingText := `
  -no-color
    Disables colored command output. Alternatively, VAULT_CLI_NO_COLOR may be
    set.
`

	if usageOpts&usageOptsNoNamespace == 0 {
		helpText = helpText + namespaceText
	}

	helpText = helpText + remainingText
	return strings.TrimSpace(helpText)
}

// funcVar is a type of flag that accepts a function that is the string given
// by the user.
type funcVar func(s string) error

func (f funcVar) Set(s string) error { return f(s) }
func (f funcVar) String() string     { return "" }
func (f funcVar) IsBoolFlag() bool   { return false }
