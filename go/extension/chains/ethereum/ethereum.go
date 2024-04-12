package ethereum

import (
	"context"
	"fmt"
	"regexp"

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
		return nil, nil
	}
	if !common.IsHexAddress(addr) {
		return nil, nil
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

func (b *Backend) FindGrants(ctx context.Context, contract, owner, grantee, resource string) ([]*chains.Grant, error) {
	if !common.IsHexAddress(contract) {
		return nil, nil
	}
	if !common.IsHexAddress(owner) {
		return nil, nil
	}
	if !common.IsHexAddress(grantee) {
		return nil, nil
	}

	reg, err := registry.NewRegistryCaller(common.HexToAddress(contract), b.cl)
	if err != nil {
		return nil, fmt.Errorf("create registry failed: %w", err)
	}
	grantList, err := reg.FindGrants(&bind.CallOpts{Context: ctx}, common.HexToAddress(owner), common.HexToAddress(grantee), resource)
	if err != nil {
		return nil, fmt.Errorf("finding grants failed: %w", err)
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

var publicKeyRegex = regexp.MustCompile("^(:?0x)?[0-9a-fA-F]{130}$")

// Extracted this to be able to test without having to create an RPC connection.
func isValidPublicKey(publicKey string) bool {
	return publicKeyRegex.MatchString(publicKey)
}

func (b *Backend) IsValidPublicKey(publicKey string) bool {
	return isValidPublicKey(publicKey)
}
