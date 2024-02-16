package keeper

import (
	"context"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"

	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

// InterchainAccountFromAddress implements the Query/InterchainAccountFromAddress gRPC method
func (k Keeper) InterchainAccountFromAddress(goCtx context.Context, req *types.QueryInterchainAccountFromAddressRequest) (*types.QueryInterchainAccountFromAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	portID, err := icatypes.NewControllerPortID(req.Owner)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "could not find account: %s", err)
	}

	ica, found := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, req.ConnectionId, portID)
	if !found {
		return nil, status.Errorf(codes.NotFound, "no account found for connectionID: %s and portID: %s",
			req.ConnectionId, portID)
	}

	return types.NewQueryInterchainAccountResponse(ica), nil
}

// Params returns params of the mint module.
func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

// AutoTx implements the Query/AutoTxgRPC method
func (k Keeper) AutoTx(c context.Context, req *types.QueryAutoTxRequest) (*types.QueryAutoTxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	txID, err := strconv.ParseUint(req.Id, 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	autoTxInfo, err := k.TryGetAutoTxInfo(ctx, txID)
	if err != nil {
		return nil, err
	}
	// for msg := range autoTxInfo.Msgs{
	// 	makeReadableMsgData(&autoTxInfo, msg)
	// }

	return &types.QueryAutoTxResponse{
		AutoTxInfo: autoTxInfo,
	}, nil
}

// AutoTxHistory implements the Query/AutoTxHistory method
func (k Keeper) AutoTxHistory(c context.Context, req *types.QueryAutoTxHistoryRequest) (*types.QueryAutoTxHistoryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	txID, err := strconv.ParseUint(req.Id, 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	autoTxHistory, err := k.TryGetAutoTxHistory(ctx, txID)
	if err != nil {
		return nil, err
	}

	return &types.QueryAutoTxHistoryResponse{
		History: autoTxHistory.History,
	}, nil
}

// AutoTxs implements the Query/AutoTxs gRPC method
func (k Keeper) AutoTxs(c context.Context, req *types.QueryAutoTxsRequest) (*types.QueryAutoTxsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	autoTxs := make([]types.AutoTxInfo, 0)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.AutoTxKeyPrefix)
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(_ []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			var c types.AutoTxInfo
			k.cdc.MustUnmarshal(value, &c)
			autoTxs = append(autoTxs, c)

		}
		return true, nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryAutoTxsResponse{
		AutoTxInfos: autoTxs,
		Pagination:  pageRes,
	}, nil
}

// AutoTxsForOwner implements the Query/AutoTxsForOwner gRPC method
func (k Keeper) AutoTxsForOwner(c context.Context, req *types.QueryAutoTxsForOwnerRequest) (*types.QueryAutoTxsForOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	autoTxs := make([]types.AutoTxInfo, 0)

	ownerAddress, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, err
	}
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetAutoTxsByOwnerPrefix(ownerAddress))
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(key []byte, _ []byte, accumulate bool) (bool, error) {
		if accumulate {
			autoTxID := types.GetIDFromBytes(key /* [types.TimeTimeLen:] */)
			autoTxInfo := k.GetAutoTxInfo(ctx, autoTxID)
			// msg, err := icatypes.DeserializeCosmosTx(k.cdc, autoTxInfo.Data)
			// if err != nil {
			// 	return false, err
			// }
			// makeReadableMsgData(&autoTxInfo, msg)
			autoTxs = append(autoTxs, autoTxInfo)

		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryAutoTxsForOwnerResponse{
		AutoTxInfos: autoTxs,
		Pagination:  pageRes,
	}, nil
}

// AutoTxsForOwner implements the Query/AutoTxsForOwner gRPC method
func (k Keeper) AutoTxIbcTxUsage(c context.Context, req *types.QueryAutoTxIbcUsageRequest) (*types.QueryAutoTxIbcUsageResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	autoTxIbcUsage := make([]types.AutoTxIbcUsage, 0)

	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.AutoTxIbcUsageKeyPrefix)
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(_ []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			var c types.AutoTxIbcUsage
			k.cdc.MustUnmarshal(value, &c)
			autoTxIbcUsage = append(autoTxIbcUsage, c)

		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryAutoTxIbcUsageResponse{
		AutoTxIbcUsage: autoTxIbcUsage,
		Pagination:     pageRes,
	}, nil
}

// func makeReadableMsgData(info *types.AutoTxInfo, msg []sdk.Msg) {
// 	info.Data = []byte(sdk.MsgTypeURL(msg[0]) + "," + msg[0].String())
// 	//fmt.Printf(string(info.Data))
// }
