package claim

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/claim/keeper"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {

	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}

	k.IterateVestingQueue(ctx, ctx.BlockHeader().Time, func(recipientAddr sdk.AccAddress, action int32, period int32, endTime time.Time) bool {
		claimRecord, err := k.GetClaimRecord(ctx, recipientAddr)
		if err != nil {
			panic("Failed to get claim record ")
		}
		//fmt.Printf("iteration for %s; current time %s; action: %v; period: %v\n", endTime, ctx.BlockHeader().Time, action, period)
		claimRecord.Status[action].VestingPeriodCompleted[period] = true
		err = k.SetClaimRecord(ctx, claimRecord)
		if err != nil {
			panic("Failed to set claim record ")

		}
		k.RemoveEntryFromVestingQueue(ctx, recipientAddr.String(), endTime)

		return true
	})
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
