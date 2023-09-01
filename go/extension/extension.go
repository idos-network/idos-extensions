package extension

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/idos-network/idos-extensions/extension/chains"

	"github.com/kwilteam/kwil-extensions/server"
	"github.com/kwilteam/kwil-extensions/types"
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
				"get_block_height": server.WithOutputsCheck(e.BlockHeight, 1),
				"has_grants":       server.WithInputsCheck(server.WithOutputsCheck(e.GrantsFor, 1), 2),
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

	// TODO: inspect the grants themselves? e.g. locktime not passed

	var exist uint8
	if len(grants) > 0 {
		exist = 1
	}
	return encodeScalarValues(exist)
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
