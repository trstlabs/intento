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

func (k Keeper) EstimatorAll(c context.Context, req *types.QueryAllEstimatorRequest) (*types.QueryAllEstimatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var estimators []*types.Estimator
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	estimatorStore := prefix.NewStore(store, types.KeyPrefix(types.EstimatorKey))

	pageRes, err := query.Paginate(estimatorStore, req.Pagination, func(key []byte, value []byte) error {
		var estimator types.Estimator
		if err := k.cdc.UnmarshalBinaryBare(value, &estimator); err != nil {
			return err
		}

		estimators = append(estimators, &estimator)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllEstimatorResponse{Estimator: estimators, Pagination: pageRes}, nil
}

func (k Keeper) Estimator(c context.Context, req *types.QueryGetEstimatorRequest) (*types.QueryGetEstimatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var estimator types.Estimator
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.KeyPrefix(types.EstimatorKey+req.Itemid)), &estimator)

	return &types.QueryGetEstimatorResponse{Estimator: &estimator}, nil
}
