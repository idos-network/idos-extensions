package extension

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/idos-network/idos-extensions/extension/chains"

	// Register the chains with ChainBackend implementations.
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

// NOTE: we are using a driver system whereby the individual chain
// implementations register with the `chains` package, allowing us to
// instantiate a chain by-name via the chains.NewChainBackend constructor.
//
// A possibly more straight forward approach might be to define the following
// newChainBackend locally that use the constructors directly from the
// individual sub-packages packages. A mater of taste.

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
