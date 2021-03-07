package command

import (
	"fmt"
	"os"

	"github.com/jkevlin/vault-cli/pkg/config"
	"github.com/jkevlin/vault-cli/pkg/configstore/configfile"
	homedir "github.com/mitchellh/go-homedir"
)

var (
	configPath string
	app        *config.Config
)

func AppInit(path string) error {
	if path == "" {
		if path = os.Getenv("VAULTCONFIG"); path == "" {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			configPath := home + "/.vaultcli"
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				err = os.Mkdir(configPath, 0755)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
			path = home + "/.vaultcli/config.yaml"
		}
	}
	store := configfile.NewConfigFileStore()
	_app, err := store.Read(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	app = _app
	return err
}
