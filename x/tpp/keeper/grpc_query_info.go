package keeper

import (
	"context"
	"encoding/hex"

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

	CodeHash := k.computeKeeper.GetCodeHash(ctx, req.Codeid)
	//	fmt.Printf("query codeHash: %X\n", CodeHash)
	//	fmt.Printf("query codeHash string: %s", CodeHash)
	//	fmt.Printf("query codeHash string string: %s", string(CodeHash))
	//store := ctx.KVStore(k.storeKey)

	//codeHash := store.Get([]byte(req.Codeid))
	if CodeHash == nil {
		return nil, status.Error(codes.InvalidArgument, "no code hash")
	}

	//itemStore := prefix.NewStore(store, types.KeyPrefix(types.ItemKey))

	return &types.QueryCodeHashResponse{Codehash: CodeHash, CodehashString: hex.EncodeToString(CodeHash)}, nil
}
