package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

func (k Keeper) HandleRelayerReward(ctx sdk.Context, relayer sdk.AccAddress, rewardType int) {
	if !k.GetRelayerRewardsAvailability(ctx) {
		return
	}
	p := k.GetParams(ctx)

	fmt.Printf("p.RelayerRewards %v\n", p.RelayerRewards)
	incentiveCoin := sdk.NewCoin(types.Denom, sdk.NewInt(p.RelayerRewards[rewardType]))

	fmt.Printf("relayer %v\n", relayer.String())
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, relayer, sdk.NewCoins(incentiveCoin))
	if err != nil {
		//set incentives unavailable
		k.SetRelayerRewardsAvailability(ctx, false)
	}

}

func (k Keeper) SetRelayerRewardsAvailability(ctx sdk.Context, rewardsAvailable bool) {
	store := ctx.KVStore(k.storeKey)
	value := []byte("false")
	if rewardsAvailable {
		value = []byte("true")
	}
	store.Set([]byte(types.KeyRelayerRewardsAvailability), value)
}

// GetRelayerRewardsAvailability returns the rewards availability bool
func (k Keeper) GetRelayerRewardsAvailability(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	value := store.Get([]byte(types.KeyRelayerRewardsAvailability))
	return string(value) == "true"
}
