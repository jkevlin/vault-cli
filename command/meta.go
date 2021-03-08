package command

import (
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
	// flagAddress string

	// // Whether to not-colorize output
	// noColor bool

	// // The region to send API requests
	// region string

	// // namespace to send API requests
	// namespace string

	// // token is used for ACLs to access privileged information
	// token string

	// caCert        string
	// caPath        string
	// clientCert    string
	// clientKey     string
	// tlsServerName string
	// insecure      bool
}

// AutocompleteFlags returns a set of flag completions for the given flag set.
func (m *Meta) AutocompleteFlags(fs FlagSetFlags) complete.Flags {
	if fs&FlagSetClient == 0 {
		return nil
	}

	return complete.Flags{
		"-address": complete.PredictAnything,
		"-region":  complete.PredictAnything,
		//"-namespace":       NamespacePredictor(m.Client, nil),
		"-no-color":        complete.PredictNothing,
		"-ca-cert":         complete.PredictFiles("*"),
		"-ca-path":         complete.PredictDirs("*"),
		"-client-cert":     complete.PredictFiles("*"),
		"-client-key":      complete.PredictFiles("*"),
		"-insecure":        complete.PredictNothing,
		"-tls-server-name": complete.PredictNothing,
		"-tls-skip-verify": complete.PredictNothing,
		"-token":           complete.PredictAnything,
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
  -address=<addr>
    The address of the Nomad server.
    Overrides the NOMAD_ADDR environment variable if set.
    Default = http://127.0.0.1:4646
  -region=<region>
    The region of the Nomad servers to forward commands to.
    Overrides the NOMAD_REGION environment variable if set.
    Defaults to the Agent's local region.
`

	namespaceText := `
  -namespace=<namespace>
    The target namespace for queries and actions bound to a namespace.
    Overrides the NOMAD_NAMESPACE environment variable if set.
    If set to '*', job and alloc subcommands query all namespaces authorized
    to user.
    Defaults to the "default" namespace.
`

	// note: that although very few commands use color explicitly, all of them
	// return red-colored text on error so we don't want to make this
	// configurable
	remainingText := `
  -no-color
    Disables colored command output. Alternatively, NOMAD_CLI_NO_COLOR may be
    set.
  -ca-cert=<path>
    Path to a PEM encoded CA cert file to use to verify the
    Nomad server SSL certificate.  Overrides the NOMAD_CACERT
    environment variable if set.
  -ca-path=<path>
    Path to a directory of PEM encoded CA cert files to verify
    the Nomad server SSL certificate. If both -ca-cert and
    -ca-path are specified, -ca-cert is used. Overrides the
    NOMAD_CAPATH environment variable if set.
  -client-cert=<path>
    Path to a PEM encoded client certificate for TLS authentication
    to the Nomad server. Must also specify -client-key. Overrides
    the NOMAD_CLIENT_CERT environment variable if set.
  -client-key=<path>
    Path to an unencrypted PEM encoded private key matching the
    client certificate from -client-cert. Overrides the
    NOMAD_CLIENT_KEY environment variable if set.
  -tls-server-name=<value>
    The server name to use as the SNI host when connecting via
    TLS. Overrides the NOMAD_TLS_SERVER_NAME environment variable if set.
  -tls-skip-verify
    Do not verify TLS certificate. This is highly not recommended. Verification
    will also be skipped if NOMAD_SKIP_VERIFY is set.
  -token
    The SecretID of an ACL token to use to authenticate API requests with.
    Overrides the NOMAD_TOKEN environment variable if set.
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
