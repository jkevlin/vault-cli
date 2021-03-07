package configfile

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/jkevlin/vault-cli/pkg/config"
	"github.com/jkevlin/vault-cli/pkg/configstore"
	"github.com/mitchellh/go-homedir"
)

type configfile struct {
}

// NewConfigFileStore should return a pointer to a configfile client
func NewConfigFileStore() configstore.Store {
	return &configfile{}
}

func (cf *configfile) Read(path string) (*config.Config, error) {
	fn := expandHomePath(path)
	bytes, err := ioutil.ReadFile(fn)
	if err != nil && strings.Contains(err.Error(), "no such file") {
		bytes = getDefaultConfig()
	} else if err != nil {
		return nil, err
	}
	c := config.Config{}
	err = yaml.Unmarshal([]byte(bytes), &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (cf *configfile) Write(path string, cfg *config.Config) error {
	return nil
}

// getHomeDir returns the home dir
func getHomeDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	return home, nil
}

// expandHomePath translate home dir
func expandHomePath(path string) string {
	if path != "" && path[:1] == "~" {
		home, err := getHomeDir()
		if err != nil {
			return ""
		}
		return home + path[1:]
	}
	return path
}

func getDefaultConfig() []byte {
	bytes := []byte(`apiVersion: v1
kind: Config
contexts:
- name: local
  context:
    cluster: local
    namespace: nextgen
    session:
      token: root
      lease-duration: 7200
      expires: 2582395696
      renewable: true
    user: localuser
clusters:
- name: local
  cluster:
    certificate-authority: ""
    insecure-skip-tls-verify: true
    server: http://127.0.0.1:8200
current-context: local
users:
- name: localuser
  user:
    client-certificate: ""
    client-key: ""
    password: ""
    username: ""
    roleID: ""
    secretID: ""
    ignore-namespace-on-auth: false
`)
	return bytes

}
