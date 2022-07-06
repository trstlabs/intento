package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

type ComputeHooks interface {

	// AfterExitPool is called after ExitPool, ExitSwapShareAmountIn, and ExitSwapExternAmountOut
	AfterRecurringSend(ctx sdk.Context, sender sdk.AccAddress)
	// AfterAutoSwap is called after SwapExactAmountIn and SwapExactAmountOut
	AfterAutoSwap(ctx sdk.Context, sender sdk.AccAddress)
}

var _ ComputeHooks = MultiComputeHooks{}

// combine multiple module hooks, all hook functions are run in array sequence
type MultiComputeHooks []ComputeHooks

// Creates hooks for the module Module
func NewMultiComputeHooks(hooks ...ComputeHooks) MultiComputeHooks {
	return hooks
}

func (h MultiComputeHooks) AfterRecurringSend(ctx sdk.Context, sender sdk.AccAddress) {
	for i := range h {
		h[i].AfterRecurringSend(ctx, sender)
	}
}

func (h MultiComputeHooks) AfterAutoSwap(ctx sdk.Context, sender sdk.AccAddress) {
	for i := range h {
		h[i].AfterAutoSwap(ctx, sender)
	}
}
