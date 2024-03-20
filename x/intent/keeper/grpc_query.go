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

	"github.com/trstlabs/intento/x/intent/types"
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

// Action implements the Query/ActiongRPC method
func (k Keeper) Action(c context.Context, req *types.QueryActionRequest) (*types.QueryActionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	id, err := strconv.ParseUint(req.Id, 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	actionInfo, err := k.TryGetActionInfo(ctx, id)
	if err != nil {
		return nil, err
	}
	// for msg := range actionInfo.Msgs{
	// 	makeReadableMsgData(&actionInfo, msg)
	// }

	return &types.QueryActionResponse{
		ActionInfo: actionInfo,
	}, nil
}

// ActionHistory implements the Query/ActionHistory method
func (k Keeper) ActionHistory(c context.Context, req *types.QueryActionHistoryRequest) (*types.QueryActionHistoryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	id, err := strconv.ParseUint(req.Id, 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	actionHistory, err := k.TryGetActionHistory(ctx, id)
	if err != nil {
		return nil, err
	}

	return &types.QueryActionHistoryResponse{
		History: actionHistory.History,
	}, nil
}

// Actions implements the Query/Actions gRPC method
func (k Keeper) Actions(c context.Context, req *types.QueryActionsRequest) (*types.QueryActionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	actions := make([]types.ActionInfo, 0)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ActionKeyPrefix)
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(_ []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			var c types.ActionInfo
			k.cdc.MustUnmarshal(value, &c)
			actions = append(actions, c)

		}
		return true, nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryActionsResponse{
		ActionInfos: actions,
		Pagination:  pageRes,
	}, nil
}

// ActionsForOwner implements the Query/ActionsForOwner gRPC method
func (k Keeper) ActionsForOwner(c context.Context, req *types.QueryActionsForOwnerRequest) (*types.QueryActionsForOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	actions := make([]types.ActionInfo, 0)

	ownerAddress, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, err
	}
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetActionsByOwnerPrefix(ownerAddress))
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(key []byte, _ []byte, accumulate bool) (bool, error) {
		if accumulate {
			autoID := types.GetIDFromBytes(key /* [types.TimeTimeLen:] */)
			actionInfo := k.GetActionInfo(ctx, autoID)
			// msg, err := icatypes.DeserializeCosmosTx(k.cdc, actionInfo.Data)
			// if err != nil {
			// 	return false, err
			// }
			// makeReadableMsgData(&actionInfo, msg)
			actions = append(actions, actionInfo)

		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryActionsForOwnerResponse{
		ActionInfos: actions,
		Pagination:  pageRes,
	}, nil
}

// ActionsForOwner implements the Query/ActionsForOwner gRPC method
func (k Keeper) ActionIbcTxUsage(c context.Context, req *types.QueryActionIbcUsageRequest) (*types.QueryActionIbcUsageResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	actionIbcUsage := make([]types.ActionIbcUsage, 0)

	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ActionIbcUsageKeyPrefix)
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(_ []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			var c types.ActionIbcUsage
			k.cdc.MustUnmarshal(value, &c)
			actionIbcUsage = append(actionIbcUsage, c)

		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryActionIbcUsageResponse{
		ActionIbcUsage: actionIbcUsage,
		Pagination:     pageRes,
	}, nil
}

// func makeReadableMsgData(info *types.ActionInfo, msg []sdk.Msg) {
// 	info.Data = []byte(sdk.MsgTypeURL(msg[0]) + "," + msg[0].String())
// 	//fmt.Printf(string(info.Data))
// }
