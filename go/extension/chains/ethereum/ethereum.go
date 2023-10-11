package ethereum

import (
	"context"
	"fmt"

	"github.com/idos-network/idos-extensions/extension/chains"
	"github.com/idos-network/idos-extensions/extension/chains/ethereum/registry"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Backend is the Ethereum implementation of a chains.ChainBackend for the
// Fractal extension.
type Backend struct {
	cl *ethclient.Client
}

// New creates a new Ethereum backend using the provided RPC URL.
func New(url string) (*Backend, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("dial rpc failed: %w", err)
	}
	return &Backend{client}, nil
}

func (b *Backend) Height(ctx context.Context) (uint64, error) {
	return b.cl.BlockNumber(ctx)
}

func (b *Backend) GrantsFor(ctx context.Context, contract, addr, resource string) ([]*chains.Grant, error) {
	if !common.IsHexAddress(contract) {
		return nil, nil //errors.New("invalid contract address")
	}
	if !common.IsHexAddress(addr) {
		return nil, nil //errors.New("invalid grantee address")
	}

	reg, err := registry.NewRegistryCaller(common.HexToAddress(contract), b.cl)
	if err != nil {
		return nil, fmt.Errorf("create registry failed: %w", err)
	}
	grantList, err := reg.GrantsFor(&bind.CallOpts{
		Context: ctx,
	}, common.HexToAddress(addr), resource)
	if err != nil {
		return nil, fmt.Errorf("get grants for failed: %w", err)
	}

	grants := make([]*chains.Grant, len(grantList))
	for i := range grantList {
		gIn := &grantList[i]
		grants[i] = &chains.Grant{
			Owner:       gIn.Owner.Hex(),
			LockedUntil: gIn.LockedUntil.Uint64(),
			Grantee:     gIn.Grantee.Hex(),
			DataID:      gIn.DataId,
		}
	}

	return grants, nil
}

type driver struct{}

func (d *driver) New(url string) (chains.ChainBackend, error) {
	return New(url)
}

func (d *driver) Chain() string {
	return "eth"
}

func init() {
	chains.RegisterDriver(&driver{})
}
