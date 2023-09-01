# idOS Kwil extensions

This is a Kwil extension that exposes idOS smart contract functions (external view/pure) on various blockchains to Kwil.
Developers can call those functions in Kuneiform.

## configuration

This extension will read configuration from system environment variables.

```bash
export ETH_RPC_URL=xxx
export NEAR_RPC_URL=xxx
```

## run

```bash
# run with go
go run main.go

# or docker
make docker
# put your configuration in .env file
docker compose up -d
```
