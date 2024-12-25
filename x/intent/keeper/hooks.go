package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

type IntentHooks interface {

	// AfterActionLocal is called after an action on behalf using authz
	AfterActionLocal(ctx sdk.Context, sender sdk.AccAddress)
	// AfterActionICA is called after a MsgExecuteContract or MsgInstantiateContract
	AfterActionICA(ctx sdk.Context, sender sdk.AccAddress)
}

var _ IntentHooks = MultiIntentHooks{}

// combine multiple module hooks, all hook functions are run in array sequence
type MultiIntentHooks []IntentHooks

// Creates hooks for the module Module
func NewMultiIntentHooks(hooks ...IntentHooks) MultiIntentHooks {
	return hooks
}

func (h MultiIntentHooks) AfterActionLocal(ctx sdk.Context, sender sdk.AccAddress) {
	for i := range h {
		h[i].AfterActionLocal(ctx, sender)
	}
}

func (h MultiIntentHooks) AfterActionICA(ctx sdk.Context, sender sdk.AccAddress) {
	for i := range h {
		h[i].AfterActionICA(ctx, sender)
	}
}
