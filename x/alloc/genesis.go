package alloc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/alloc/keeper"
	"github.com/trstlabs/intento/x/alloc/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {

	k.SetParams(ctx, genState.Params)

}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params: k.GetParams(ctx),
	}
}
