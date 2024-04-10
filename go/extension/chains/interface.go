package chains

import (
	"context"
	"fmt"
)

// Grant models a Fractal grant for certain data to a grantee.
type Grant struct {
	Owner       string `json:"owner"`
	LockedUntil uint64 `json:"locked_until"` // uint256
	Grantee     string `json:"grantee"`
	DataID      string `json:"data_id"`
}

// GrantChecker specifies the method required of a type to retrieve grants.
type GrantChecker interface {
	GrantsFor(ctx context.Context, registry, addr, resource string) ([]*Grant, error)
	// NOTE: HasGrant just returns a boolean at present, but presumably we'll
	// want to return the entire slice of grants, which includes lock times.
}

// ChainBackend must be implemented for a block chain backend to provide height
// and account grant data to the extension.
type ChainBackend interface {
	GrantChecker
	Height(context.Context) (uint64, error)
	IsValidPublicKey(public_key string) bool
	// TODO: Close(), for ws connections
}

// NewChainBackend creates a new ChainBackend, which is an instance of one of
// the implementations in the various sub-packages.
func NewChainBackend(chain, url string) (ChainBackend, error) {
	maker, ok := drivers[chain]
	if !ok {
		return nil, fmt.Errorf("unrecognized chain %v", chain)
	}
	return maker.New(url)
}

// Driver is a type that creates a new chain backend using a specified URL. The
// Driver system is used so that packages implementing a ChainBackend may
// register with the chains package, allowing NewChainBackend to create the
// backend by name. An alternate approach is to define a constructor in the
// consumer (extension package) imports the sub-packages directly, calling them
// directly e.g. near.New(...) vs chains.NewChainBackend("near", ...).
type Driver interface {
	New(url string) (ChainBackend, error)
	Chain() string
}

var (
	drivers = map[string]Driver{}
)

// RegisterDriver should be called by packages that implement ChainBackend so
// that NewChainBackend may instantiate them by name.
func RegisterDriver(d Driver) {
	drivers[d.Chain()] = d
}
