package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/jkevlin/vault-cli/pkg/secretservice"
)

// SetSession Builds a session object
// expireSkewFactor in seconds
func (c *Config) SetSession(token string, duration *int64, renewable *bool, expireSkewFactor int64) *Session {
	if token != "" || duration != nil || renewable != nil {
		session := Session{
			Token:         token,
			LeaseDuration: duration,
			Renewable:     renewable,
		}
		// reduce the expire time by 30 minutes
		zero := int64(0)
		if duration != nil && *duration > zero {
			expire := (time.Now().UTC().Unix() + int64(*duration)) - expireSkewFactor
			session.Expires = &expire
		}
		return &session
	}
	return nil
}

// vaultSessionExpireSkewFactor the amount of time to subtract from Expire to account for clock skew
const vaultSessionExpireSkewFactor = int64(30 * 60)

// // VaultCertLogin will get a token from vault
// // CertLogin reqires TLS and a CA.crt
// func (c *Config) VaultCertLogin(namespace, url, cert, key, cacert string, insecureSkipVerify bool) (*api.Secret, error) {
// 	cert = c.ExpandHomePath(cert)
// 	key = c.ExpandHomePath(key)
// 	clientCertKey, err := tls.LoadX509KeyPair(cert, key)
// 	if err != nil {
// 		return nil, err
// 	}

// 	client := &http.Client{}
// 	if cacert != "" {
// 		caCert, err := c.ReadFile(cacert)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if len(caCert) == 0 {
// 			return nil, errors.New("could not read caCert")
// 		}
// 		caCertPool := x509.NewCertPool()
// 		caCertPool.AppendCertsFromPEM(caCert)

// 		client = &http.Client{
// 			Transport: &http.Transport{
// 				TLSClientConfig: &tls.Config{
// 					RootCAs:            caCertPool,
// 					Certificates:       []tls.Certificate{clientCertKey},
// 					InsecureSkipVerify: insecureSkipVerify,
// 				},
// 			},
// 		}
// 	} else {
// 		client = &http.Client{
// 			Transport: &http.Transport{
// 				TLSClientConfig: &tls.Config{
// 					Certificates:       []tls.Certificate{clientCertKey},
// 					InsecureSkipVerify: insecureSkipVerify,
// 				},
// 			},
// 		}
// 	}

// 	req, err := http.NewRequest("POST", url+"/v1/auth/cert/login", nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if namespace != "" {
// 		req.Header.Add("X-Vault-Namespace", namespace)
// 	}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	jdata, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	secret := api.Secret{}
// 	err = json.Unmarshal([]byte(jdata), &secret)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if secret.Auth == nil || secret.Auth.ClientToken == "" {
// 		return nil, fmt.Errorf("could not get token: body:%v", string(jdata))
// 	}
// 	return &secret, nil
// }

// // VaultUserPassLogin will get a token from vault
// func (c *Config) VaultUserPassLogin(namespace, authurl, username, password, cacert string, insecureSkipVerify bool) (*api.Secret, error) {
// 	client := &http.Client{}
// 	if cacert != "" {
// 		caCert, err := c.ReadFile(cacert)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		caCertPool := x509.NewCertPool()
// 		caCertPool.AppendCertsFromPEM(caCert)

// 		client.Transport = &http.Transport{
// 			TLSClientConfig: &tls.Config{
// 				RootCAs:            caCertPool,
// 				InsecureSkipVerify: insecureSkipVerify,
// 			},
// 		}
// 	}
// 	values := map[string]string{"password": password}

// 	jsonValue, _ := json.Marshal(values)
// 	nsPath := ""
// 	if namespace != "" {
// 		nsPath = namespace + "/"
// 	}
// 	req, err := http.NewRequest("POST", authurl+"/v1/"+nsPath+"auth/userpass/login/"+username, bytes.NewBuffer(jsonValue))
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	jdata, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	secret := api.Secret{}
// 	err = json.Unmarshal([]byte(jdata), &secret)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &secret, nil
// }

// // VaultAppRoleLogin will get a token from vault
// func (c *Config) VaultAppRoleLogin(namespace, authurl, roleID, secretID, cacert string, insecureSkipVerify bool) (*api.Secret, error) {
// 	client := &http.Client{}
// 	if cacert != "" {
// 		caCert, err := c.ReadFile(cacert)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		caCertPool := x509.NewCertPool()
// 		caCertPool.AppendCertsFromPEM(caCert)

// 		client.Transport = &http.Transport{
// 			TLSClientConfig: &tls.Config{
// 				RootCAs:            caCertPool,
// 				InsecureSkipVerify: insecureSkipVerify,
// 			},
// 		}
// 	}
// 	values := map[string]string{"role_id": roleID, "secret_id": secretID}

// 	jsonValue, _ := json.Marshal(values)
// 	nsPath := ""
// 	if namespace != "" {
// 		nsPath = namespace + "/"
// 	}
// 	req, err := http.NewRequest("POST", authurl+"/v1/"+nsPath+"auth/approle/login", bytes.NewBuffer(jsonValue))
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	jdata, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	secret := api.Secret{}
// 	err = json.Unmarshal([]byte(jdata), &secret)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &secret, nil
// }

// GetSession will return an existing session or create a new one and save the config
func (cfg *Config) GetSession(secretsvc secretservice.SecretService, configfile, contextName string, forceNewSession bool) (*Session, error) {
	for _, c := range cfg.Contexts {
		if c.Name == contextName {
			now := time.Now().UTC().Unix()
			if c.Session.Expires == nil {
				zero := int64(0)
				c.Session.Expires = &zero
			}

			if forceNewSession || c.Session.Token == "" || now > *c.Session.Expires {
				cluster := cfg.GetClusterByName(c.Cluster)
				user := cfg.GetUserByName(c.User)
				if cluster == nil || cluster.Server == "" {
					return nil, errors.New("cluster must have server address")
				}
				if user == nil {
					return nil, errors.New("user must have cert and key")
				}
				var err error
				var response *api.Secret
				ns := c.Namespace
				if user.IgnoreNamespaceOnAuth == true {
					ns = ""
				}
				if user.ClientCert != "" {
					response, err = secretsvc.CertLogin(ns, cluster.Server, "cert", user.ClientCert, user.ClientKey, cluster.CertAuth, cluster.InsecureSkipTLSVerify)
				} else if user.Username != "" {
					response, err = secretsvc.UserPassLogin(ns, cluster.Server, "userpass", user.Username, user.Password, cluster.CertAuth, cluster.InsecureSkipTLSVerify)
				} else if user.RoleID != "" {
					response, err = secretsvc.AppRoleLogin(ns, cluster.Server, "approle", user.RoleID, user.SecretID, cluster.CertAuth, cluster.InsecureSkipTLSVerify)

				} else {
					return nil, fmt.Errorf("GetSession login requires credentials")
				}
				if err != nil {
					return nil, err
				}

				if response.Auth != nil {
					duration := int64(response.Auth.LeaseDuration)
					session := cfg.SetSession(response.Auth.ClientToken, &duration, &response.Auth.Renewable, vaultSessionExpireSkewFactor)
					c.Session = *session
					cfg.SaveConfig(configfile)
				}
			}
			return &c.Session, nil
		}
	}

	return nil, errors.New("context not found")
}

// StopSession will return an existing session or create a new one and save the config
func (cfg *Config) StopSession(configfile, contextName string) error {

	for _, c := range cfg.Contexts {
		if c.Name == contextName {
			c.Session = Session{}
		}
	}
	return nil
}
