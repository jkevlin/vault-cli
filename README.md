# vault-cli

vault-cli is a vault automation tool, used to configure a vault server
with all of the namespaces, endpoints, policies, roles auth endpoins, etc.

vault-cli stores its state in convienent yaml format.  This allows a company to
maintain configuration control over the contents of a vault server.

## Try it out

In first terminal window

```bash
vault server -dev -dev-root-token-id root -dev-listen-address 127.0.0.1:8200
```

In second terminal

- Clone https://github.com/jkevlin/vault-cli

```bash
cd vault-cli
go mod vendor
go build
./vault-cli put vaultnamespace -c=local -namespace=root "local-*"
```
