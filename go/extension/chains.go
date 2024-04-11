package extension

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/idos-network/idos-extensions/extension/chains"
	"github.com/idos-network/idos-extensions/extension/chains/ethereum"
	"github.com/idos-network/idos-extensions/extension/chains/near"
)

func startChainBackend(chain, url string) (chains.ChainBackend, error) {
	be, err := newChainBackend(chain, url)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	height, err := be.Height(ctx)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("unable to get height of chain %v: %w", chain, err)
	}
	fmt.Printf("Started chain %v at height %v\n", chain, height)
	return be, nil
}

func newChainBackend(chain, url string) (chains.ChainBackend, error) {
	switch chain {
	case "eth":
		return ethereum.New(url)
	case "arbitrum":
		return ethereum.New(url)
	case "near":
		return near.New(url)
	default:
		return nil, errors.New("unsupported chain")
	}
}
