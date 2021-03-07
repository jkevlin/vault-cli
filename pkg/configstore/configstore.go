package configstore

import (
	"github.com/jkevlin/vault-cli/pkg/config"
)

// Store Interface for controller operations needed by task workers
//go:generate counterfeiter -o fakes/store.go --fake-name FakeStore . Store
type Store interface {
	Read(path string) (*config.Config, error)
	Write(path string, cfg *config.Config) error
}
