package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/danieljdd/tpp/x/tpp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CodeHash(c context.Context, req *types.QueryCodeHashRequest) (*types.QueryCodeHashResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	codeHash := k.computeKeeper.GetCodeHash(ctx, req.Codeid)
	fmt.Printf("query codeHash: %X\n", codeHash)
	//store := ctx.KVStore(k.storeKey)

	//codeHash := store.Get([]byte(req.Codeid))
	if codeHash == nil {
		return nil, status.Error(codes.InvalidArgument, "no code hash")
	}
	//itemStore := prefix.NewStore(store, types.KeyPrefix(types.ItemKey))

	return &types.QueryCodeHashResponse{Codehash: codeHash}, nil
}
