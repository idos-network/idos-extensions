package extension

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/idos-network/idos-extensions/extension/chains"

	"github.com/kwilteam/kwil-extensions/server"
	"github.com/kwilteam/kwil-extensions/types"

	"github.com/mr-tron/base58"
)

const (
	metaKeyRegistryAddress = "registry_address"
	metaKeyChain           = "chain"
)

type metadata map[string]string

func (m metadata) RegistryAddress() string {
	return m[metaKeyRegistryAddress]
}

func (m metadata) Chain() string {
	return m[metaKeyChain]
}

var requiredMetadata = map[string]string{
	metaKeyRegistryAddress: "",
	metaKeyChain:           "eth",
}

type FractalExt struct {
	backends map[string]chains.ChainBackend
}

func NewFractalExt(rpcURLs map[string]string) (*FractalExt, error) {
	// Construct the chain backends (RPC clients) for each URL.
	backends := make(map[string]chains.ChainBackend, len(rpcURLs))
	for chain, url := range rpcURLs {
		be, err := startChainBackend(chain, url)
		if err != nil {
			return nil, err
		}
		backends[chain] = be
	}

	return &FractalExt{
		backends: backends,
	}, nil
}

func (e *FractalExt) BuildServer(logger *log.Logger) (*server.ExtensionServer, error) {
	return server.Builder().
		Named(e.Name()).
		WithInitializer(initialize).
		WithLoggerFunc(func(l string) {
			logger.Println(l)
		}).
		WithMethods(
			map[string]server.MethodFunc{
				"get_block_height":               server.WithOutputsCheck(e.BlockHeight, 1),
				"has_grants":                     server.WithInputsCheck(server.WithOutputsCheck(e.GrantsFor, 1), 2),
				"has_locked_grants":              server.WithInputsCheck(server.WithOutputsCheck(e.LockedGrantsFor, 1), 2),
				"implicit_address_to_public_key": server.WithInputsCheck(server.WithOutputsCheck(e.ImplicitAddressToPublicKey, 1), 1),
				"determine_wallet_type":          server.WithInputsCheck(server.WithOutputsCheck(e.DetermineWalletType, 1), 1),
				"is_valid_public_key":            server.WithInputsCheck(server.WithOutputsCheck(e.IsValidPublicKey, 1), 1),
			}).
		Build()
}

func (e *FractalExt) Name() string {
	return "idos"
}

// getMetadata returns the metadata from given context. An error will be
// returned if a required metadata is not provided. Default values should be
// provided during initialization, not execution.
func (e *FractalExt) getMetadata(ctx *types.ExecutionContext) (metadata, error) {
	for k := range requiredMetadata {
		if _, ok := ctx.Metadata[k]; !ok {
			return nil, fmt.Errorf("metadata %s is required", k)
		}
	}
	return metadata(ctx.Metadata), nil
}

func (e *FractalExt) chainBackend(ctx *types.ExecutionContext) (chains.ChainBackend, string, error) {
	m, err := e.getMetadata(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("get metadata failed: %w", err)
	}

	be, ok := e.backends[m.Chain()]
	if !ok {
		return nil, "", errors.New("unsupported")
	}
	regAddr := m.RegistryAddress()
	return be, regAddr, nil
}

// BlockHeight returns the current block height
func (e *FractalExt) BlockHeight(ctx *types.ExecutionContext, _ ...*types.ScalarValue) ([]*types.ScalarValue, error) {
	be, _, err := e.chainBackend(ctx)
	if err != nil {
		return nil, err
	}
	height, err := be.Height(ctx.Ctx)
	if err != nil {
		return nil, err
	}
	return encodeScalarValues(height)
}

// GrantsFor returns whether the given user has grants
// @return1 1 for true, 0 for false
func (e *FractalExt) GrantsFor(ctx *types.ExecutionContext, values ...*types.ScalarValue) ([]*types.ScalarValue, error) {
	be, regAddr, err := e.chainBackend(ctx)
	if err != nil {
		return nil, err
	}
	granteeAddress, err := values[0].String()
	if err != nil {
		return nil, fmt.Errorf("convert value to string failed: %w", err)
	}

	dataId, err := values[1].String()
	if err != nil {
		return nil, fmt.Errorf("convert value to string failed: %w", err)
	}

	grants, err := be.GrantsFor(ctx.Ctx, regAddr, granteeAddress, dataId)
	if err != nil {
		return nil, fmt.Errorf("failed to check grants: %w", err)
	}

	var exist uint8
	if grants != nil && len(grants) > 0 {
		exist = 1
	}
	return encodeScalarValues(exist)
}

// @return1 1 for true, 0 for false
func (e *FractalExt) LockedGrantsFor(ctx *types.ExecutionContext, values ...*types.ScalarValue) ([]*types.ScalarValue, error) {
	be, regAddr, err := e.chainBackend(ctx)
	if err != nil {
		return nil, err
	}
	ownerAddress, err := values[0].String()
	if err != nil {
		return nil, fmt.Errorf("convert value to string failed: %w", err)
	}

	dataId, err := values[1].String()
	if err != nil {
		return nil, fmt.Errorf("convert value to string failed: %w", err)
	}

	lockedGrants, err := chains.LockedGrantsFor(be, ctx.Ctx, regAddr, ownerAddress, dataId)
	if err != nil {
		return nil, fmt.Errorf("failed to check for locked grants: %w", err)
	}

	var result uint8
	if len(lockedGrants) > 0 {
		result = 1
	}
	return encodeScalarValues(result)
}

func (e *FractalExt) ImplicitAddressToPublicKey(ctx *types.ExecutionContext, values ...*types.ScalarValue) ([]*types.ScalarValue, error) {
	inputHex, err := values[0].String()
	if err != nil {
		return nil, fmt.Errorf("convert value to string failed: %w", err)
	}
	binaryString, _ := hex.DecodeString(inputHex)
	base58 := base58.Encode(binaryString)
	var public_key string
	if len(inputHex) != 64 || base58 == "" {
		public_key = ""
	} else {
		public_key = fmt.Sprintf("ed25519:%s", base58)
	}

	return encodeScalarValues(public_key)
}

func (e *FractalExt) IsValidPublicKey(ctx *types.ExecutionContext, values ...*types.ScalarValue) ([]*types.ScalarValue, error) {
	be, _, err := e.chainBackend(ctx)
	if err != nil {
		return nil, err
	}

	public_key, err := values[0].String()
	if err != nil {
		return nil, fmt.Errorf("convert value to string failed: %w", err)
	}

	var result uint8
	if be.IsValidPublicKey(public_key) {
		result = 1
	}
	return encodeScalarValues(result)
}

var EVM_ADDRESS_REGEX = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

// This has very dumb logic: eth address returns EVM type, and NEAR returns otherwise.
// TODO: make the logic more detailed and return error is the address is neither EVM no NEAR.
func (e *FractalExt) DetermineWalletType(ctx *types.ExecutionContext, values ...*types.ScalarValue) ([]*types.ScalarValue, error) {
	address, err := values[0].String()
	if err != nil {
		return nil, fmt.Errorf("convert value to string failed: %w", err)
	}

	var wallet_type string
	if EVM_ADDRESS_REGEX.MatchString(address) {
		wallet_type = "EVM"
	} else {
		wallet_type = "NEAR"
	}

	return encodeScalarValues(wallet_type)
}

// initialize checks that the meta data includes all required fields and applies
// any default values.
func initialize(ctx context.Context, metadata map[string]string) (map[string]string, error) {
	for k, def := range requiredMetadata {
		_, ok := metadata[k]
		if ok {
			continue
		}
		if def == "" {
			return nil, fmt.Errorf("metadata %s is required", k)
		}
		metadata[k] = def
	}

	return metadata, nil
}

func encodeScalarValues(values ...any) ([]*types.ScalarValue, error) {
	scalarValues := make([]*types.ScalarValue, len(values))
	for i, v := range values {
		scalarValue, err := types.NewScalarValue(v)
		if err != nil {
			return nil, fmt.Errorf("convert value to scalar failed: %w", err)
		}

		scalarValues[i] = scalarValue
	}

	return scalarValues, nil
}
