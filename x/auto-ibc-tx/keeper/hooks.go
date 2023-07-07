package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

type AutoIbcTxHooks interface {

	// AfterAutoTxAuthz is called after an action on behalf using authz
	AfterAutoTxAuthz(ctx sdk.Context, sender sdk.AccAddress)
	// AfterAutoTxWasm is called after a MsgExecuteContract or MsgInstantiateContract
	AfterAutoTxWasm(ctx sdk.Context, sender sdk.AccAddress)
}

var _ AutoIbcTxHooks = MultiAutoIbcTxHooks{}

// combine multiple module hooks, all hook functions are run in array sequence
type MultiAutoIbcTxHooks []AutoIbcTxHooks

// Creates hooks for the module Module
func NewMultiAutoIbcTxHooks(hooks ...AutoIbcTxHooks) MultiAutoIbcTxHooks {
	return hooks
}

func (h MultiAutoIbcTxHooks) AfterAutoTxAuthz(ctx sdk.Context, sender sdk.AccAddress) {
	for i := range h {
		h[i].AfterAutoTxAuthz(ctx, sender)
	}
}

func (h MultiAutoIbcTxHooks) AfterAutoTxWasm(ctx sdk.Context, sender sdk.AccAddress) {
	for i := range h {
		h[i].AfterAutoTxWasm(ctx, sender)
	}
}
