package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/trstlabs/trst/x/item/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ItemAll(c context.Context, req *types.QueryAllItemsRequest) (*types.QueryAllItemsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var items []*types.Item
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	itemStore := prefix.NewStore(store, types.KeyPrefix(types.ItemKey))

	pageRes, err := query.Paginate(itemStore, req.Pagination, func(key []byte, value []byte) error {
		var item types.Item
		if err := k.cdc.Unmarshal(value, &item); err != nil {
			return err
		}

		items = append(items, &item)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllItemsResponse{Item: items, Pagination: pageRes}, nil
}

func (k Keeper) ListedItemsAll(c context.Context, req *types.QueryAllListedItemsRequest) (*types.QueryAllListedItemsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var items []*types.Item
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	itemStore := prefix.NewStore(store, types.ListedItemQueuePrefix)

	pageRes, err := query.Paginate(itemStore, req.Pagination, func(key []byte, value []byte) error {
		var item types.Item
		if err := k.cdc.Unmarshal(value, &item); err != nil {
			return err
		}

		//k.cdc.MustUnmarshal(value, &item) //strconv.FormatUint(types.GetItemIDFromBytes(value), 10))
		items = append(items, &item)
		return nil

	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	//items := k.GetAllListedItems(ctx)

	return &types.QueryAllListedItemsResponse{Item: items, Pagination: pageRes}, nil
}

func (k Keeper) Item(c context.Context, req *types.QueryGetItemRequest) (*types.QueryGetItemResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var item types.Item
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ItemKey))
	k.cdc.MustUnmarshal(store.Get(append(types.KeyPrefix(types.ItemKey), types.Uint64ToByte(req.Id)...)), &item)

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
	//itemStore := prefix.NewStore(store, types.ListedItemQueuePrefix)

	return &types.QuerySellerItemsResponse{Item: items}, nil

}
