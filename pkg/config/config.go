package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// SaveConfig writes yaml to a file
func (c *Config) SaveConfig(filename string) error {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, bytes, 0644)
}

func (c *Config) String() string {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// // GetHomeDir returns the home dir
// func (c *Config) GetHomeDir() (string, error) {
// 	home, err := homedir.Dir()
// 	if err != nil {
// 		return "", fmt.Errorf(err.Error())
// 	}
// 	return home, nil
// }

// // ExpandHomePath translate home dir
// func (c *Config) ExpandHomePath(path string) string {
// 	if path != "" && path[:1] == "~" {
// 		home, err := c.GetHomeDir()
// 		if err != nil {
// 			return ""
// 		}
// 		return home + path[1:]
// 	}
// 	return path
// }

// // ReadFile expands home dir
// func (c *Config) ReadFile(filename string) ([]byte, error) {
// 	fn := c.ExpandHomePath(filename)
// 	return ioutil.ReadFile(fn)
// }
