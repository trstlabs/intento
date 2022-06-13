package claim

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/trst/x/claim/keeper"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {

	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}

	k.IterateVestingQueue(ctx, ctx.BlockHeader().Time, func(period int32) bool { return true })
	// End Airdrop
	timeElapsed := ctx.BlockTime().Sub(params.AirdropStartTime)
	if timeElapsed > params.DurationUntilDecay+params.DurationOfDecay {
		// airdrop time passed
		err := k.EndAirdrop(ctx)
		if err != nil {
			panic(err)
		}
	}
}
