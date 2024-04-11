package chains

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type stubGrantChecker struct {
	timelockDeltasFromNow *[]int64
}

func (sc *stubGrantChecker) GrantsFor(_ctx context.Context, registry, addr, resource string) ([]*Grant, error) {
	now := uint64(time.Now().Unix())

	result := make([]*Grant, 0, len(*sc.timelockDeltasFromNow))

	for _, delta := range *sc.timelockDeltasFromNow {
		result = append(result, &Grant{
			Owner:       "fake owner",
			LockedUntil: uint64(int64(now) + delta),
			Grantee:     "fake grantee",
			DataID:      "fake data_id",
		})
	}
	return result, nil
}

func Test_LockedGrantsFor(t *testing.T) {
	tests := []struct {
		timelockDeltasFromNow []int64
		wantedLen             int
	}{
		{
			[]int64{},
			0,
		},
		{
			[]int64{10},
			1,
		},
		{
			[]int64{10, 20},
			2,
		},
		{
			[]int64{-10},
			0,
		},
		{
			[]int64{-10, 10},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%+v", tt), func(t *testing.T) {
			mockGrantChecker := &stubGrantChecker{timelockDeltasFromNow: &tt.timelockDeltasFromNow}
			if res, _ := LockedGrantsFor(mockGrantChecker, nil, "irrelevant", "irrelevant", "irrelevant"); len(res) != tt.wantedLen {
				t.Errorf("LockedGrantsFor() = %+v, wanted len %v", res, tt.wantedLen)
			}
		})
	}

}
