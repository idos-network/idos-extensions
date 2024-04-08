package near

import (
	"context"
	"fmt"
	"testing"

	"github.com/eteu-technologies/near-api-go/pkg/client"
	"github.com/eteu-technologies/near-api-go/pkg/client/block"
	"github.com/eteu-technologies/near-api-go/pkg/types"
)

type stubClient struct {
	height        uint64
	grantsForResp []byte
}

func (sc *stubClient) BlockDetails(ctx context.Context, block block.BlockCharacteristic) (resp client.BlockView, err error) {
	return client.BlockView{
		Author: types.AccountID("whateverpoolv42.near"),
		Header: client.BlockHeaderView{
			Height: sc.height, // this is all we need for (*Backend).Height
		},
		Chunks: []client.ChunkHeaderView{},
	}, nil
}

func (sc *stubClient) ContractViewCallFunction(ctx context.Context, accountID, methodName, argsBase64 string, block block.BlockCharacteristic) (res client.CallResult, err error) {
	return client.CallResult{
		Result: sc.grantsForResp,
	}, nil
}

func TestBackend(t *testing.T) {
	wantHeight := uint64(123413241234)
	cl := &stubClient{
		height:        wantHeight,
		grantsForResp: []byte{}, // len zero to start (no grants)
	}
	be := &Backend{cl}

	height, err := be.Height(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if height != wantHeight {
		t.Errorf("wanted height %v got %v", wantHeight, height)
	}

	grantee := "jchappelow.testnet"
	registry := "grants.jchappelow.testnet"
	dataID := "blah"

	cl.grantsForResp = []byte(fmt.Sprintf(`[
		{
			"data_id": "%s",
			"grantee": "%s",
			"locked_until": 2690839560,
			"owner": "jchappelow.testnet"
		}
	]`, dataID, grantee))

	grants, err := be.GrantsFor(context.Background(), registry, grantee, dataID)
	if err != nil {
		t.Fatalf("HasGrant unexpectedly failed: %v", err)
	}
	if len(grants) == 0 {
		t.Error("expected to have grants, but did not")
	}

	cl.grantsForResp = []byte(`null`) // nil
	grants, err = be.GrantsFor(context.Background(), registry, grantee, dataID)
	if err != nil {
		t.Fatalf("HasGrant unexpectedly failed: %v", err)
	}
	if len(grants) > 0 {
		t.Error("expected to not have grants, but did")
	}

	cl.grantsForResp = []byte(`[]`)
	grants, err = be.GrantsFor(context.Background(), registry, grantee, dataID)
	if err != nil {
		t.Fatalf("HasGrant unexpectedly failed: %v", err)
	}
	if len(grants) > 0 {
		t.Error("expected to not have grants, but did")
	}

	// invalid account name
	entries, _ := be.GrantsFor(context.Background(), "not-a-registry", grantee, dataID)
	if len(entries) != 0 {
		t.Fatalf("HasGrant accepted a bogus registry acct name")
	}

	// valid implicit acct name
	hex64Chars := "beefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef"
	_, err = be.GrantsFor(context.Background(), hex64Chars, grantee, dataID)
	if err != nil {
		t.Fatalf("HasGrant unexpectedly failed: %v", err)
	}
}

func Test_isNearAcct(t *testing.T) {
	tests := []struct {
		name string
		acct string
		want bool
	}{
		{
			"ok .near",
			"blah.near",
			true,
		},
		{
			"ok 64 hex chars",
			"beefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef",
			true,
		},
		{
			"ok .testnet",
			"blah.testnet",
			true,
		},
		{
			"not ok",
			"wrong.eth",
			false,
		},
		{
			"not ok empty",
			"",
			false,
		},
		{
			"not ok non-hex 64 chars",
			"ppefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeefbeef",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNearAcct(tt.acct); got != tt.want {
				t.Errorf("isNearAcct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isValidPublicKey(t *testing.T) {
	tests := []struct {
		publicKey string
		want      bool
	}{
		{
			"",
			false,
		},
		{
			"blah.near",
			false,
		},
		{
			"derp:herp",
			false,
		},
		{
			"ed25519:INVALID",
			false,
		},
		{
			"ed25519:7dLLbzqc6kgGAC6smmJUUh9xqxH9habnLhptauAymmUJ",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.publicKey, func(t *testing.T) {
			if res, _ := isValidPublicKey(tt.publicKey); res != tt.want {
				t.Errorf("isNearAcct() = %v, want %v", res, tt.want)
			}
		})
	}
}
