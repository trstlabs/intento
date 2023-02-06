package keeper

import (
	"context"
	"fmt"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"

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
		return nil, status.Errorf(codes.NotFound, "no account found for connectionId: %s and portID: %s",
			req.ConnectionId, portID)
	}

	return types.NewQueryInterchainAccountResponse(ica), nil
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
	msg, err := icatypes.DeserializeCosmosTx(k.cdc, autoTxInfo.Data)
	if err != nil {
		return nil, err
	}
	makeReadableMsgData(&autoTxInfo, msg)

	return &types.QueryAutoTxResponse{
		AutoTxInfo: autoTxInfo,
	}, nil
}

// AutoTxs implements the Query/AutoTxs gRPC method
func (k Keeper) AutoTxs(c context.Context, req *types.QueryAutoTxsRequest) (*types.QueryAutoTxsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	autoTxs := make([]types.AutoTxInfo, 0)

	k.IterateAutoTxInfos(ctx, func(id uint64, info types.AutoTxInfo) bool {
		msg, _ := icatypes.DeserializeCosmosTx(k.cdc, info.Data)
		makeReadableMsgData(&info, msg)
		autoTxs = append(autoTxs, info)
		return false
	})

	return &types.QueryAutoTxsResponse{
		AutoTxInfos: autoTxs,
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
			msg, err := icatypes.DeserializeCosmosTx(k.cdc, autoTxInfo.Data)
			if err != nil {
				return false, err
			}
			makeReadableMsgData(&autoTxInfo, msg)
			autoTxInfo.Data = []byte(sdk.MsgTypeURL(msg[0]) + "'value': {" + msg[0].String() + "}")
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

func makeReadableMsgData(info *types.AutoTxInfo, msg []sdk.Msg) {
	info.Data = []byte(sdk.MsgTypeURL(msg[0]) + "," + msg[0].String())
	fmt.Printf(string(info.Data))
}
