package keeper

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/intent/types"
)

func (k Keeper) HandleRelayerReward(ctx sdk.Context, relayer sdk.AccAddress, rewardType int) {
	if !k.GetRelayerRewardsAvailability(ctx) {
		return
	}
	p, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("p.RelayerRewards %v\n", p.RelayerRewards)
	incentiveCoin := sdk.NewCoin(types.Denom, math.NewInt(p.RelayerRewards[rewardType]))

	// fmt.Printf("relayer %v\n", relayer.String())
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, relayer, sdk.NewCoins(incentiveCoin))
	if err != nil {
		//set incentives unavailable
		k.SetRelayerRewardsAvailability(ctx, false)
	}

}

func (k Keeper) SetRelayerRewardsAvailability(ctx sdk.Context, rewardsAvailable bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	value := []byte("false")
	if rewardsAvailable {
		value = []byte("true")
	}
	store.Set([]byte(types.KeyRelayerRewardsAvailability), value)
}

// GetRelayerRewardsAvailability returns the rewards availability
func (k Keeper) GetRelayerRewardsAvailability(ctx sdk.Context) bool {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	value := store.Get([]byte(types.KeyRelayerRewardsAvailability))
	return string(value) == "true"
}
