package inventory

import (
	"fmt"
	"io/ioutil"

	"github.com/mitchellh/go-homedir"
)

// GetHomeDir returns the home dir
func GetHomeDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	return home, nil
}

// ExpandHomePath translate home dir
func ExpandHomePath(path string) string {
	if path != "" && path[:1] == "~" {
		home, err := GetHomeDir()
		if err != nil {
			return ""
		}
		return home + path[1:]
	}
	return path
}

// ReadFile expands home dir
func ReadFile(filename string) ([]byte, error) {
	fn := ExpandHomePath(filename)
	return ioutil.ReadFile(fn)
}
