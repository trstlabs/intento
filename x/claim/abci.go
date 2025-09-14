package claim

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/claim/keeper"
)

// EndBlocker called every block, process vesting queue + airdrop expiry
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}
	// Process vesting queue
	k.IterateVestingQueue(ctx, ctx.BlockHeader().Time,
		func(recipientAddr sdk.AccAddress, action int32, period int32, endTime time.Time) bool {
			claimRecord, err := k.GetClaimRecord(ctx, recipientAddr)
			if err != nil {
				panic("Failed to get claim record")
			}

			claimRecord.Status[action].VestingPeriodsCompleted[period] = true

			if err := k.SetClaimRecord(ctx, claimRecord); err != nil {
				panic("Failed to set claim record")
			}

			k.RemoveEntryFromVestingQueue(ctx, recipientAddr.String(), endTime, byte(action), byte(period))

			return false // keep iterating
		},
	)

	// End Airdrop if time passed
	timeElapsed := ctx.BlockTime().Sub(params.AirdropStartTime)
	if timeElapsed > params.DurationUntilDecay+params.DurationOfDecay {
		if err := k.EndAirdrop(ctx); err != nil {
			panic(err)
		}
	}
}
