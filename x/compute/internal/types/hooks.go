package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type ComputeHooks interface {

	// AfterExitPool is called after ExitPool, ExitSwapShareAmountIn, and ExitSwapExternAmountOut
	AfterComputeInstantiate(ctx sdk.Context, sender sdk.AccAddress)
	// AfterSwap is called after SwapExactAmountIn and SwapExactAmountOut
	AfterComputeExecuted(ctx sdk.Context, sender sdk.AccAddress)
}

var _ ComputeHooks = MultiComputeHooks{}

// combine multiple module hooks, all hook functions are run in array sequence
type MultiComputeHooks []ComputeHooks

// Creates hooks for the module Module
func NewMultiComputeHooks(hooks ...ComputeHooks) MultiComputeHooks {
	return hooks
}

func (h MultiComputeHooks) AfterComputeInstantiate(ctx sdk.Context, sender sdk.AccAddress) {
	for i := range h {
		h[i].AfterComputeInstantiate(ctx, sender)
	}
}

func (h MultiComputeHooks) AfterComputeExecuted(ctx sdk.Context, sender sdk.AccAddress) {
	for i := range h {
		h[i].AfterComputeExecuted(ctx, sender)
	}
}
