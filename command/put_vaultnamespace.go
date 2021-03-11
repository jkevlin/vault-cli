package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/jkevlin/vault-cli/pkg/inventory"
	vaultapi "github.com/jkevlin/vault-go/api/v1"
	"github.com/mitchellh/go-homedir"
	"github.com/posener/complete"
	"gopkg.in/yaml.v2"
)

const (
	envVaultCLIConfigDir  = "VAULTCLICONFIG"
	configDefaultDir      = ".vaultcli"
	configDefaultFileName = "config.yaml"
)

type PutVaultNamespaceCommand struct {
	Meta
}

func (c *PutVaultNamespaceCommand) Help() string {
	helpText := `
Usage: vault-cli put vaultnamespace [options]

General Options:
  ` + generalOptionsUsage() + `
`
	return strings.TrimSpace(helpText)
}

func (c *PutVaultNamespaceCommand) AutocompleteFlags() complete.Flags {
	return mergeAutocompleteFlags(c.Meta.AutocompleteFlags(),
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
	var (
		dirname string
	)

	flagSet := c.Meta.FlagSet(c.Name())
	flagSet.Usage = func() { c.Ui.Output(c.Help()) }
	flagSet.StringVar(&dirname, "f", "hack/sample/vaultnamespace", "")
	if err := flagSet.Parse(args); err != nil {
		return 1
	}
	args = flagSet.Args()
	vaultnamespacefilespec := args[0]

	var path string
	if path = os.Getenv(envVaultCLIConfigDir); path == "" {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		configPath := home + "/" + configDefaultDir
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			err = os.Mkdir(configPath, 0755)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		path = configPath + "/" + configDefaultFileName
		c.ConfigPath = path
	}
	cfg, err := c.ConfigService.Read(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting config: %s\n", err.Error())
		return 1
	}
	c.Config = cfg

	secretsvc, err := c.Config.GetServiceFromContext(c.ConfigPath, "local", "nextgen")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting service from config: %s\n", err.Error())
		return 1
	}
	c.SecretService = secretsvc

	files, err := inventory.GetFiles(dirname, vaultnamespacefilespec)
	if err != nil {
		fmt.Printf("get files error: %s\n", err.Error())
		return 1
	}
	if len(files) == 0 {
		fmt.Printf("Vault Namespace (%s) not found in inventory", vaultnamespacefilespec)
		return 1
	}

	for _, f := range files {
		filename := dirname + "/" + f
		data, err := inventory.ReadFile(filename + ".yaml")
		if err != nil {
			fmt.Println("error reading file: ", err.Error())
			return 1
		}
		vaultNamespace := vaultapi.VaultNamespace{}
		err = yaml.Unmarshal(data, &vaultNamespace)
		if err != nil {
			fmt.Printf("unable to marshal vaultnamespace: %s\n", err.Error())
			return 1
		}

		if vaultNamespace.Spec.NamespaceBase != "" {
			c.SecretService.GetClient().SetNamespace(vaultNamespace.Spec.NamespaceBase)
		}

		secret, err := c.SecretService.Read(fmt.Sprintf("/sys/namespaces/%s", vaultNamespace.Spec.NamespaceName))
		if err == nil && secret != nil {
			fmt.Printf("Vault Namespace: %s.yaml exists\n", f)
			return 0
		}
		m := make(map[string]interface{})
		_, err = c.SecretService.Write(fmt.Sprintf("/sys/namespaces/%s", vaultNamespace.Spec.NamespaceName), m)
		if err != nil {
			fmt.Printf("Vault Namespace: %s.yaml %s\n", f, err)
			return 1
		}
		fmt.Printf("Vault Namespace: %s.yaml write, OK\n", f)
	}

	return 0
}
