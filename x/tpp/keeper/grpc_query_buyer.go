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
/*
func (k Keeper) BuyerAll(c context.Context, req *types.QueryAllBuyerRequest) (*types.QueryAllBuyerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var buyers []*types.Buyer
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	buyerStore := prefix.NewStore(store, types.KeyPrefix(types.BuyerKey))

	pageRes, err := query.Paginate(buyerStore, req.Pagination, func(key []byte, value []byte) error {
		var buyer types.Buyer
		if err := k.cdc.UnmarshalBinaryBare(value, &buyer); err != nil {
			return err
		}

		buyers = append(buyers, &buyer)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllBuyerResponse{Buyer: buyers, Pagination: pageRes}, nil
}

func (k Keeper) Buyer(c context.Context, req *types.QueryGetBuyerRequest) (*types.QueryGetBuyerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var buyer types.Buyer
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BuyerKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.KeyPrefix(types.BuyerKey+req.Itemid)), &buyer)

	return &types.QueryGetBuyerResponse{Buyer: &buyer}, nil
}


*/

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

