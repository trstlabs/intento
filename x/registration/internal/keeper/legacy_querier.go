package keeper

import (
	"encoding/hex"
	"encoding/json"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/danieljdd/trst/x/registration/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryEncryptedSeed     = types.QueryEncryptedSeed
	QueryMasterCertificate = types.QueryMasterCertificate
)

// controls error output on querier - set true when testing/debugging
const debug = false

// NewQuerier creates a new querier
func NewLegacyQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		var (
			rsp interface{}
			err error
		)
		switch path[0] {
		case QueryEncryptedSeed:
			pubKey, err := hex.DecodeString(path[1])
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
			}
			rsp, err = queryEncryptedSeed(ctx, pubKey, keeper)
		case QueryMasterCertificate:
			rsp, err = queryMasterKey(ctx, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown data query endpoint")
		}
		if err != nil {
			return nil, err
		}
		if rsp == nil || reflect.ValueOf(rsp).IsNil() {
			return nil, nil
		}
		bz, err := json.MarshalIndent(rsp, "", "  ")
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return bz, nil
	}
}
