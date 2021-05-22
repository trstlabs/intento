package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/danieljdd/tpp/x/tpp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ItemAll(c context.Context, req *types.QueryAllItemRequest) (*types.QueryAllItemResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var items []*types.Item
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	itemStore := prefix.NewStore(store, types.KeyPrefix(types.ItemKey))

	pageRes, err := query.Paginate(itemStore, req.Pagination, func(key []byte, value []byte) error {
		var item types.Item
		if err := k.cdc.UnmarshalBinaryBare(value, &item); err != nil {
			return err
		}

		items = append(items, &item)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllItemResponse{Item: items, Pagination: pageRes}, nil
}

func (k Keeper) InactiveItemsAll(c context.Context, req *types.QueryAllInactiveItemsRequest) (*types.QueryAllInactiveItemsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	//var items []*types.Item
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	itemStore := prefix.NewStore(store, types.InactiveItemQueuePrefix)
	var items []string
	pageRes, err := query.Paginate(itemStore, req.Pagination, func(key []byte, value []byte) error {

		items = append(items, string(value))
		return nil

	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	//items := k.GetAllInactiveItems(ctx)

	return &types.QueryAllInactiveItemsResponse{Item: items, Pagination: pageRes}, nil
}

func (k Keeper) Item(c context.Context, req *types.QueryGetItemRequest) (*types.QueryGetItemResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var item types.Item
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ItemKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.KeyPrefix(types.ItemKey+req.Id)), &item)

	return &types.QueryGetItemResponse{Item: &item}, nil
}

func (k Keeper) SellerItems(c context.Context, req *types.QuerySellerItemsRequest) (*types.QuerySellerItemsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	//var items []*types.Item
	ctx := sdk.UnwrapSDKContext(c)

	items := k.GetAllSellerItems(ctx, req.Seller)

	//store := ctx.KVStore(k.storeKey)
	//itemStore := prefix.NewStore(store, types.InactiveItemQueuePrefix)

	return &types.QuerySellerItemsResponse{Item: items}, nil

}
