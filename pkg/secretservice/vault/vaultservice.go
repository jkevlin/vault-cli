package vault

import (
	"errors"

	"github.com/hashicorp/vault/api"
	"github.com/jkevlin/vault-cli/pkg/secretservice"
)

type vaultservice struct {
	Client *api.Client
}

// NewVaultService should return a pointer to a vaultservice client
func NewVaultService(c *api.Client) secretservice.SecretService {
	return &vaultservice{Client: c}
}

// Delete is to satisfy a lint error for this interface
func (vs *vaultservice) Delete(path string) (*api.Secret, error) {
	return vs.Client.Logical().Delete(path)
}

// List is to satisfy a lint error for this interface
func (vs *vaultservice) List(path string) (*api.Secret, error) {
	return vs.Client.Logical().List(path)
}

// Read is to satisfy a lint error for this interface
func (vs *vaultservice) Read(path string) (*api.Secret, error) {
	return vs.Client.Logical().Read(path)
}

func (vs *vaultservice) ReadWithData(path string, data map[string][]string) (*api.Secret, error) {
	return vs.Client.Logical().ReadWithData(path, data)
}

// Write is to satisfy a lint error for this interface
func (vs *vaultservice) Write(path string, data map[string]interface{}) (*api.Secret, error) {
	return vs.Client.Logical().Write(path, data)
}

// IsKVv2 check version
func (vs *vaultservice) IsKVv2(path string) (string, bool, error) {
	mountPath, version, err := kvPreflightVersionRequest(vs.Client, path)
	if err != nil {
		return "", false, err
	}

	return mountPath, version == 2, nil
}

// kvPreflightVersionRequest do a preflight call
func kvPreflightVersionRequest(client *api.Client, path string) (string, int, error) {
	// We don't want to use a wrapping call here so save any custom value and
	// reservice after
	currentWrappingLookupFunc := client.CurrentWrappingLookupFunc()
	client.SetWrappingLookupFunc(nil)
	defer client.SetWrappingLookupFunc(currentWrappingLookupFunc)
	currentOutputCurlString := client.OutputCurlString()
	client.SetOutputCurlString(false)
	defer client.SetOutputCurlString(currentOutputCurlString)

	r := client.NewRequest("GET", "/v1/sys/internal/ui/mounts/"+path)
	resp, err := client.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		// If we get a 404 we are using an older version of vault, default to
		// version 1
		if resp != nil && resp.StatusCode == 404 {
			return "", 1, nil
		}

		return "", 0, err
	}

	secret, err := api.ParseSecret(resp.Body)
	if err != nil {
		return "", 0, err
	}
	if secret == nil {
		return "", 0, errors.New("nil response from pre-flight request")
	}
	var mountPath string
	if mountPathRaw, ok := secret.Data["path"]; ok {
		mountPath = mountPathRaw.(string)
	}
	options := secret.Data["options"]
	if options == nil {
		return mountPath, 1, nil
	}
	versionRaw := options.(map[string]interface{})["version"]
	if versionRaw == nil {
		return mountPath, 1, nil
	}
	version := versionRaw.(string)
	switch version {
	case "", "1":
		return mountPath, 1, nil
	case "2":
		return mountPath, 2, nil
	}

	return mountPath, 1, nil
}
