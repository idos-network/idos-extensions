package chains

import (
	"context"
	"time"
)

func LockedGrantsFor(gc GrantChecker, ctx context.Context, registry, addr, resource string) ([]*Grant, error) {
	now := uint64(time.Now().Unix())

	allGrants, err := gc.GrantsFor(ctx, registry, addr, resource)
	if err != nil {
		return nil, err
	}
	if len(allGrants) < 1 {
		return allGrants, nil
	}

	result := make([]*Grant, 0, len(allGrants))
	for _, g := range allGrants {
		if g.LockedUntil > now {
			result = append(result, g)
		}
	}

	return result, nil
}
