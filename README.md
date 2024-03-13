# idOS Kwil extensions

This is a Kwil extension that exposes idOS smart contract functions (external view/pure) on various blockchains to Kwil.

Developers can call those functions in Kuneiform.

### Kuneiform Interface

The extension has two methods:

- `get_block_height() returns (int)`: Returns the most recent block height for the chain it is listening to.
- `has_grants(grantee_address, credential_id) returns (int)`: Returns whether or not a wallet has grants allowing it to access the credential id.  Returns a 1 if they have grants, or a 0 if they do not.

_For convenience, Kwil v0.5 added support for raising errors in an extension as a way of canceling action execution.  This extension can update the `has_grants` method to make this more consumable within Kuneiform._

### Supported Blockchains

This extension currently supports:

- Ethereum mainnet
- NEAR mainnet

### Configuration

The IDOS extension requires two configuration variables: `ETH_RPC_URL` and `NEAR_RPC_URL`.  They need to be provided at startup.

```bash
export ETH_RPC_URL=xxx
export NEAR_RPC_URL=xxx
```
By default ETH_RPC_URL is `http://127.0.0.1:8545/`. If you want to connect th chain run in your host machine, i.e. hardhat running locally in development environment, you should use `http://host.docker.internal:8545/` instead to access host machine port from the container.


### Initialization

When initializing the IDOS extension (in a Kuneiform schema), the database deployer needs to provide a contract address.  All extensions must also provide an alias for which the newly initialized extension can be referenced.  For example:

```SQL
database idos_db;

use idos {
    registry_address: "0xcd321d2211918124F24971604fCfA9f22370e46f",
    chain: "eth"
} as grants_eth;
```

In the above example, we are initializing the extension `idos` to the contract address `0xcd321d2211918124F24971604fCfA9f22370e46f` on NEAR, and aliasing the initialized extension as `grants_eth`.  The extension can now be used within an action block:

```SQL
action get_credential ($id) public {
    $can_access = grants_eth.has_grants(@caller, $id);
    SELECT CASE WHEN $can_access != 1
    THEN ERROR('caller does not have access') END;

    SELECT content
    FROM credentials
    WHERE id = $id;
}
```

The above query will prevent a caller from executing the SELECT query if they do not have grant(s) for the credential id.  This example can be found in the [test schema](./schemas/idos.kf).

The extensions should work for any Solidity smart contract that fulfills the following interface:

```TS
interface IKwilWhitelist {
    struct Grant {
      address owner;
      address grantee;
      string dataId;
      uint256 lockedUntil;
    }

    function grants_for(address grantee, string dataId)
        public view
        returns(Grant[] memory)
}
```

The extension can also be initialized to several chains at the same time:

```SQL
database idos_db;

use idos {
    registry_address: "0xcd321d2211918124F24971604fCfA9f22370e46f",
    chain: "eth"
} as grants_eth;

use idos {
    registry_address: "0xcd321d2211918124F24971604fCfA9f22370e46f",
    chain: "near"
} as grants_near;
```

## Run
If non-existing in `docker network list`, create a network to allow other containers to connect to extension via bridge mode: `docker network create kwil-dev`

```bash
# run with go
go run go/main.go

# or docker
make docker
# configure using docker-compose.yml
docker compose up -d
```

## Building

The [makefile](<./makefile>) contains functionalities for building Docker images.

### Build Local Architecture

For local tesing, you only need to build for your local architecture.  To do this, run:

```bash
make docker
```

### Build Multi-Architecture

When building for public usage, it's important that you build for all architectures you plan to support. For more info on this, see [this helpful blog](<https://www.thorsten-hans.com/how-to-build-multi-arch-docker-images-with-ease/>).

To build for linux amd64 and linux arm64, run:

```bash
make docker-multi-arch
```

To build multi-arch and push it to Dockerhub, run:

```bash
make docker-multi-arch PUSH=1
```
