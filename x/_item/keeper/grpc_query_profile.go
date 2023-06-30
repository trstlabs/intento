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

func (k Keeper) AllProfiles(c context.Context, req *types.QueryAllProfilesRequest) (*types.QueryAllProfilesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var profiles []*types.Profile
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	profileStore := prefix.NewStore(store, types.KeyPrefix(types.ProfileKey))

	pageRes, err := query.Paginate(profileStore, req.Pagination, func(key []byte, value []byte) error {
		var profile types.Profile
		if err := k.cdc.Unmarshal(value, &profile); err != nil {
			return err
		}

		profiles = append(profiles, &profile)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllProfilesResponse{Profile: profiles, Pagination: pageRes}, nil
}

func (k Keeper) Profile(c context.Context, req *types.QueryGetProfileRequest) (*types.QueryGetProfileResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var profile types.Profile
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProfileKey))
	k.cdc.MustUnmarshal(store.Get(append(types.KeyPrefix(types.ProfileKey), []byte(req.Owner)...)), &profile)

	return &types.QueryGetProfileResponse{Profile: &profile}, nil
}
