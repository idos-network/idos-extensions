package near

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/eteu-technologies/near-api-go/pkg/client"
	"github.com/eteu-technologies/near-api-go/pkg/client/block"
	"github.com/idos-network/idos-extensions/extension/chains"
)

// isNearAcct checks if the string is a valid near account name. This is either
// an "implicit" account, which is 64 hexadecimal characters, or a named account
// that is suffixed by the network name. For example. fractal.near or
// app0.testnet.
func isNearAcct(acct string) bool {
	if strings.HasSuffix(acct, ".near") || strings.HasSuffix(acct, ".testnet") ||
		strings.HasSuffix(acct, ".shardnet") || strings.HasSuffix(acct, ".guildnet") {
		return true
	}
	if len(acct) != 64 {
		return false
	}
	_, err := hex.DecodeString(acct)
	return err == nil
}

// chainClient is implemented by the near-api-go/pkg/client.Client type, but is
// defined as an interface so it may be stubbed out for testing.
type chainClient interface {
	BlockDetails(ctx context.Context, block block.BlockCharacteristic) (resp client.BlockView, err error)
	ContractViewCallFunction(ctx context.Context, accountID, methodName, argsBase64 string, block block.BlockCharacteristic) (res client.CallResult, err error)
}

// Backend is the NEAR implementation of a chains.ChainBackend for the Fractal
// extension.
type Backend struct {
	cl chainClient
}

// New creates a new NEAR backend using the provided RPC URL.
func New(url string) (*Backend, error) {
	cl, err := client.NewClient(url)
	if err != nil {
		return nil, err
	}
	return &Backend{&cl}, nil
}

func (nb *Backend) Height(ctx context.Context) (uint64, error) {
	res, err := nb.cl.BlockDetails(ctx, block.FinalityFinal())
	if err != nil {
		return 0, err
	}
	return res.Header.Height, nil
}

// The following args structure may need to be updated when the NEAR registry
// contract is implemented.  For now we are assuming:
//  1. method name is `grants_for`
//  2. arg1 is `grantee` of type `AccountId` (see near_sdk), represented as
//     a string in the marshalled request.
//  3. arg2 is `dataId`, a string
//  4. the return is a list of grant objects
type grantArgs struct {
	Grantee string `json:"grantee"`
	DataID  string `json:"data_id"`
}

type grantResp struct {
	Owner       string `json:"owner"`
	LockedUntil uint64 `json:"locked_until"` // uint256

	grantArgs // no idea why
}

func base64CallArgs(thing any) (string, error) {
	b, err := json.Marshal(thing)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (nb *Backend) GrantsFor(ctx context.Context, registry, acct, resource string) ([]*chains.Grant, error) {
	if !isNearAcct(registry) {
		return nil, errors.New("invalid NEAR account name of registry contract")
	}
	if !isNearAcct(acct) {
		return nil, errors.New("invalid NEAR account name of grantee")
	}

	base64Args, err := base64CallArgs(grantArgs{
		Grantee: acct,
		DataID:  resource,
	})
	if err != nil {
		return nil, err
	}
	res, err := nb.cl.ContractViewCallFunction(ctx, registry, `grants_for`, base64Args,
		block.FinalityFinal())
	if err != nil {
		return nil, err
	}

	var grantList []grantResp
	if err = json.Unmarshal(res.Result, &grantList); err != nil {
		return nil, fmt.Errorf("unmarshal failed (%w) - res: %v",
			err, string(res.Result))
	}

	grants := make([]*chains.Grant, len(grantList))
	for i := range grantList {
		gIn := &grantList[i]
		grants[i] = &chains.Grant{
			Owner:       gIn.Owner,
			LockedUntil: gIn.LockedUntil,
			Grantee:     gIn.Grantee,
			DataID:      gIn.DataID,
		}
	}

	return grants, nil
}

type driver struct{}

func (d *driver) New(url string) (chains.ChainBackend, error) {
	return New(url)
}

func (d *driver) Chain() string {
	return "near"
}

func init() {
	chains.RegisterDriver(&driver{})
}
