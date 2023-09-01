//go:build online

package near_test

// NOTE: go test -v -tags online

import (
	"context"
	"testing"
	"time"

	"github.com/kwilteam/extension-fractal-demo/extension/chains/near"
)

func TestBackend_Online(t *testing.T) {
	be, err := near.New("https://rpc.testnet.near.org")
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	height, err := be.Height(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Height: %d", height)

	// jsonArgs := `{"account_id": "linear-protocol.near"}`
	// argsB64 := base64.StdEncoding.EncodeToString([]byte(jsonArgs))
	// callRes, err := be.cl.ContractViewCallFunction(ctx, "nodeasy.poolv1.near",
	// 	"get_account_staked_balance", argsB64, block.FinalityFinal())

	// Check the deployed (fake) registry contract at
	// https://explorer.testnet.near.org/accounts/grants.jchappelow.testnet
	//
	// NOTE that this deployment always returns the same response from
	// grants_from because I could not figure out how to ensure it was a
	// read-only (view) method. Despite no `&mut self`, read-only RPC's always
	// errored claiming "storage_write" was being used. Regardless, the method
	// has the correct inputs/outputs so we can test an actual testnet RPC.
	grants, err := be.GrantsFor(ctx, "grants.jchappelow.testnet",
		"dummy.testnet", "whatever_data")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("grants to \"dummy.testnet\" for data_id \"whaterver_data\" "+
		"on \"grants.jchappelow.testnet\"? ", len(grants) > 0)
}
