package keeper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/trstlabs/trst/x/compute/internal/types"
)

const (
	QueryListContractByCode      = types.QueryListContractByCode
	QueryGetContract             = types.QueryGetContract
	QueryGetContractPublicState  = types.QueryGetContractPublicState
	QueryGetContractPrivateState = types.QueryGetContractPrivateState
	QueryGetCode                 = types.QueryGetCode
	QueryListCode                = types.QueryListCode
	QueryContractAddress         = types.QueryContractAddress
	QueryContractKey             = types.QueryContractKey
	QueryContractHash            = types.QueryContractHash
	QueryMasterCertificate       = types.QueryMasterCertificate
)

// NewLegacyQuerier creates a new querier
func NewLegacyQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		var (
			rsp interface{}
			err error
			bz  []byte
		)
		switch path[0] {
		case QueryGetContract:
			addr, err := sdk.AccAddressFromBech32(path[1])
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
			}
			rsp, err = queryContractInfo(ctx, addr, keeper)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
			}
		case QueryGetContractPublicState:
			addr, err := sdk.AccAddressFromBech32(path[1])
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
			}
			rsp, err = queryContractPublicState(ctx, addr, keeper)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
			}
		case QueryListContractByCode:
			codeID, err := strconv.ParseUint(path[1], 10, 64)
			if err != nil {
				return nil, sdkerrors.Wrapf(types.ErrInvalid, "code id: %s", err.Error())
			}
			rsp, err = queryContractListByCode(ctx, codeID, keeper)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
			}
		case QueryGetContractPrivateState:

			return queryPrivateContractState(ctx, path[1], req, keeper)

		case QueryGetCode:
			codeID, err := strconv.ParseUint(path[1], 10, 64)
			if err != nil {
				return nil, sdkerrors.Wrapf(types.ErrInvalid, "code id: %s", err.Error())
			}
			rsp, err = QueryCode(ctx, codeID, keeper)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
			}
		case QueryListCode:
			rsp, err = queryCodeList(ctx, keeper)
		/*
			case QueryContractHistory:
				contractAddr, err := sdk.AccAddressFromBech32(path[1])
				if err != nil {
					return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
				}
				rsp, err = queryContractHistory(ctx, contractAddr, keeper)
		*/
		case QueryContractAddress:
			bz, err = queryContractAddress(ctx, path[1], keeper)
		case QueryContractKey:
			addr, err := sdk.AccAddressFromBech32(path[1])
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
			}
			bz, err = queryContractKey(ctx, addr, keeper)
		case QueryContractHash:
			addr, err := sdk.AccAddressFromBech32(path[1])
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
			}
			bz, err = queryContractHash(ctx, addr, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("unknown data query endpoint %v", path[0]))
		}
		if err != nil {
			return nil, err
		}

		if bz != nil {
			return bz, nil
		}

		if rsp == nil || reflect.ValueOf(rsp).IsNil() {
			return nil, nil
		}

		//bz, err = keeper.legacyAmino.MarshalJSON(rsp)
		bz, err = json.MarshalIndent(rsp, "", "  ")
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return bz, nil
	}
}

func queryPrivateContractState(ctx sdk.Context, bech string, req abci.RequestQuery, keeper Keeper) (json.RawMessage, error) {
	contractAddr, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, bech)
	}

	/*
		var resultData []types.Model
			switch queryMethod {
			case QueryMethodContractStateAll:
				// this returns a serialized json object (which internally encoded binary fields properly)
				for iter := keeper.GetContractState(ctx, contractAddr); iter.Valid(); iter.Next() {
					resultData = append(resultData, types.Model{
						Key:   iter.Key(),
						Value: iter.Value(),
					})
				}
				if resultData == nil {
					resultData = make([]types.Model, 0)
				}
			case QueryMethodContractStateRaw:
				// this returns the raw data from the state, base64-encoded
				return keeper.QueryRaw(ctx, contractAddr, data), nil

			case QueryMethodContractStateSmart:
	*/

	// we enforce a subjective gas limit on all queries to avoid infinite loops
	ctx = ctx.WithGasMeter(sdk.NewGasMeter(keeper.queryGasLimit))
	// this returns raw bytes (must be base64-encoded)
	return keeper.QueryPrivate(ctx, contractAddr, req.Data, false)

	/*
			default:
				return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, queryMethod)
			}


		bz, err := json.Marshal(resultData)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return bz, nil
	*/

}
