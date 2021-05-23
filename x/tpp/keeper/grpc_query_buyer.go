package keeper

import (
	"context"

	//	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	//	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/danieljdd/tpp/x/tpp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) BuyerItems(c context.Context, req *types.QueryBuyerItemsRequest) (*types.QueryBuyerItemsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	//var items []*types.Item
	ctx := sdk.UnwrapSDKContext(c)

	items := k.GetAllBuyerItems(ctx, req.Buyer)

	//store := ctx.KVStore(k.storeKey)
	//itemStore := prefix.NewStore(store, types.InactiveItemQueuePrefix)

	return &types.QueryBuyerItemsResponse{Item: items}, nil

}
