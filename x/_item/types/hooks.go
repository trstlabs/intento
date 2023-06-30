package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type ItemHooks interface {
	// AfterPoolCreated is called after CreatePool
	//AfterItemCreated(ctx sdk.Context, sender sdk.AccAddress, poolId uint64)
	// AfterJoinPool is called after JoinPool, JoinSwapExternAmountIn, and JoinSwapShareAmountOut
	//AfterItemEstimated(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, enterCoins sdk.Coins, shareOutAmount sdk.Int)
	// AfterExitPool is called after ExitPool, ExitSwapShareAmountIn, and ExitSwapExternAmountOut
	AfterItemBought(ctx sdk.Context, sender sdk.AccAddress)
	// AfterSwap is called after SwapExactAmountIn and SwapExactAmountOut
	AfterItemTokenized(ctx sdk.Context, sender sdk.AccAddress)
}

var _ ItemHooks = MultiItemHooks{}

// combine multiple module hooks, all hook functions are run in array sequence
type MultiItemHooks []ItemHooks

// Creates hooks for the module Module
func NewMultiItemHooks(hooks ...ItemHooks) MultiItemHooks {
	return hooks
}

/*func (h MultiItemHooks) AfterItemCreated(ctx sdk.Context, sender sdk.AccAddress, poolId uint64) {
	for i := range h {
		h[i].AfterItemCreated(ctx, sender, poolId)
	}
}

func (h MultiItemHooks) AfterItemEstimated(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, enterCoins sdk.Coins, shareOutAmount sdk.Int) {
	for i := range h {
		h[i].AfterItemEstimated(ctx, sender, poolId, enterCoins, shareOutAmount)
	}
}
*/
func (h MultiItemHooks) AfterItemBought(ctx sdk.Context, sender sdk.AccAddress) {
	for i := range h {
		h[i].AfterItemBought(ctx, sender)
	}
}

func (h MultiItemHooks) AfterItemTokenized(ctx sdk.Context, sender sdk.AccAddress) {
	for i := range h {
		h[i].AfterItemTokenized(ctx, sender)
	}
}
