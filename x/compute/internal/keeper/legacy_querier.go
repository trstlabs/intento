package keeper

import (
	"encoding/json"
	"reflect"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/trstlabs/trst/x/compute/internal/types"
)

const (
	QueryListContractByCode = types.QueryListContractByCode
	QueryGetContract        = types.QueryGetContract
	QueryGetContractResult  = types.QueryGetContractResult
	QueryGetContractState   = types.QueryGetContractState
	QueryGetCode            = types.QueryGetCode
	QueryListCode           = types.QueryListCode
	QueryContractAddress    = types.QueryContractAddress
	QueryContractKey        = types.QueryContractKey
	QueryContractHash       = types.QueryContractHash
	QueryMasterCertificate  = types.QueryMasterCertificate
	//QueryContractHistory    = "contract-history"
)

const QueryMethodContractStateSmart = "smart"

/*
const (
	QueryMethodContractStateSmart = "smart"
	QueryMethodContractStateAll   = "all"
	QueryMethodContractStateRaw   = "raw"
)
*/

// NewLegacyQuerier creates a new querier
func NewLegacyQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		var (
			rsp interface{}
			err error
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
		case QueryGetContractResult:
			addr, err := sdk.AccAddressFromBech32(path[1])
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
			}
			rsp, err = queryContractResult(ctx, addr, keeper)
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
		case QueryGetContractState:
			//if len(path) < 3 {
			//	return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "path invalid, unknown data query endpoint")
			//	}
			return queryContractState(ctx, path[1], req, keeper)
			//rsp, err = queryContractState(ctx, path[1], path[2], req.Data, keeper)
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
			rsp, err = queryContractAddress(ctx, path[1], keeper)
		case QueryContractKey:
			addr, err := sdk.AccAddressFromBech32(path[1])
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
			}
			rsp, err = queryContractKey(ctx, addr, keeper)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
			}
		case QueryContractHash:
			addr, err := sdk.AccAddressFromBech32(path[1])
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
			}
			rsp, err = queryContractHash(ctx, addr, keeper)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
			}
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

func queryContractState(ctx sdk.Context, bech string, req abci.RequestQuery, keeper Keeper) (json.RawMessage, error) {
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
	return keeper.QuerySmart(ctx, contractAddr, req.Data, false)

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
