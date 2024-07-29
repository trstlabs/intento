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

func (k Keeper) ActionHistory(ctx context.Context, req *types.QueryActionHistoryRequest) (*types.QueryActionHistoryResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Convert id from the request
	id, err := strconv.ParseUint(req.Id, 10, 64)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	// Assuming ActionHistoryEntry items are stored with keys prefixed by id
	store := prefix.NewStore(sdkCtx.KVStore(k.storeKey), types.GetActionHistoryKey(id))

	// Paginate over the prefixed store
	var historyEntries []types.ActionHistoryEntry
	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var historyEntry types.ActionHistoryEntry
		if err := k.cdc.Unmarshal(value, &historyEntry); err != nil {
			return err
		}
		historyEntries = append(historyEntries, historyEntry)
		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "paginate: %v", err)
	}

	// Return paginated results
	return &types.QueryActionHistoryResponse{
		History:    historyEntries,
		Pagination: pageRes,
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
			actionID := types.GetIDFromBytes(key)
			actionInfo := k.GetActionInfo(ctx, actionID)

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

// HostedAccount implements the Query/HostedAccount gRPC method
func (k Keeper) HostedAccount(c context.Context, req *types.QueryHostedAccountRequest) (*types.QueryHostedAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	hosted, err := k.TryGetHostedAccount(ctx, req.Address)
	if err != nil {
		return nil, err
	}

	return &types.QueryHostedAccountResponse{
		HostedAccount: hosted,
	}, nil
}

// HostedAccounts implements the Query/HostedAccounts gRPC method
func (k Keeper) HostedAccounts(c context.Context, req *types.QueryHostedAccountsRequest) (*types.QueryHostedAccountsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	hostedAccounts := make([]types.HostedAccount, 0)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostedAccountKeyPrefix)
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(_ []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			var c types.HostedAccount
			k.cdc.MustUnmarshal(value, &c)
			hostedAccounts = append(hostedAccounts, c)

		}
		return true, nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryHostedAccountsResponse{
		HostedAccounts: hostedAccounts,
		Pagination:     pageRes,
	}, nil
}

// HostedAccountsByAdmin implements the Query/HostedAccountsByAdmin gRPC method
func (k Keeper) HostedAccountsByAdmin(c context.Context, req *types.QueryHostedAccountsByAdminRequest) (*types.QueryHostedAccountsByAdminResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	hostedAccounts := make([]types.HostedAccount, 0)

	admin, err := sdk.AccAddressFromBech32(req.Admin)
	if err != nil {
		return nil, err
	}
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetHostedAccountsByAdminPrefix(admin))
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(key []byte, _ []byte, accumulate bool) (bool, error) {
		if accumulate {
			hostedAccountAddress := string(key)
			actionInfo := k.GetHostedAccount(ctx, hostedAccountAddress)

			hostedAccounts = append(hostedAccounts, actionInfo)

		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryHostedAccountsByAdminResponse{
		HostedAccounts: hostedAccounts,
		Pagination:     pageRes,
	}, nil
}
