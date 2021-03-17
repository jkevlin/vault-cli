# vault-cli

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
