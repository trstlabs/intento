package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/danieljdd/tpp/x/tpp/types"
)

/*
// GetDepositParams returns the current DepositParams from the global param store
func (k Keeper) GetActiveParams(ctx sdk.Context) types.ActiveParams {
	var activeParams types.ActiveParams

	k.paramSpace.Get(ctx, types.ParamStoreKeyActiveParams, &activeParams)
	return activeParams
}
*/
/**/
// ParamKeyTable for tpp module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&types.Params{})
}

// GetParams returns the total set of curating parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the curating parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
