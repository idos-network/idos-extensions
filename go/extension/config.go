package extension

import "fmt"

type Config struct {
	EthRpcUrl      string `env:"ETH_RPC_URL" envDefault:"https://ethereum-sepolia-rpc.publicnode.com"`
	ArbitrumRpcUrl string `env:"ARBITRUM_RPC_URL" envDefault:"https://sepolia-rollup.arbitrum.io/rpc"`
	NearRpcUrl     string `env:"NEAR_RPC_URL" envDefault:"https://rpc.testnet.near.org"`
	Port           int    `env:"PORT" envDefault:"50055"`
}

func (c *Config) ListenAddr() string {
	return fmt.Sprintf(":%d", c.Port)
}

func (c *Config) ChainURLs() map[string]string {
	rpcURLs := make(map[string]string)
	if url := c.EthRpcUrl; url != "" {
		rpcURLs["eth"] = url
	}
	if url := c.ArbitrumRpcUrl; url != "" {
		rpcURLs["arbitrum"] = url
	}
	if url := c.NearRpcUrl; url != "" {
		rpcURLs["near"] = url
	}
	return rpcURLs
}
