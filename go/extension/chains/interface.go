package chains

import (
	"context"
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
	FindGrants(ctx context.Context, registry, owner, grantee, dataId string) ([]*Grant, error)
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
