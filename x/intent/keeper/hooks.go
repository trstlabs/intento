package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

type IntentHooks interface {

	// AfterActionAuthz is called after an action on behalf using authz
	AfterActionAuthz(ctx sdk.Context, sender sdk.AccAddress)
	// AfterActionWasm is called after a MsgExecuteContract or MsgInstantiateContract
	AfterActionWasm(ctx sdk.Context, sender sdk.AccAddress)
}

var _ IntentHooks = MultiIntentHooks{}

// combine multiple module hooks, all hook functions are run in array sequence
type MultiIntentHooks []IntentHooks

// Creates hooks for the module Module
func NewMultiIntentHooks(hooks ...IntentHooks) MultiIntentHooks {
	return hooks
}

func (h MultiIntentHooks) AfterActionAuthz(ctx sdk.Context, sender sdk.AccAddress) {
	for i := range h {
		h[i].AfterActionAuthz(ctx, sender)
	}
}

func (h MultiIntentHooks) AfterActionWasm(ctx sdk.Context, sender sdk.AccAddress) {
	for i := range h {
		h[i].AfterActionWasm(ctx, sender)
	}
}
