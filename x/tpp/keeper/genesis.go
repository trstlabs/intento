package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/danieljdd/tpp/x/tpp/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// InitGenesis initializes the curating module state
func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) {
	//k.SetParams(ctx, state.Params)

	// NOTE: since the reward pool is a module account, the auth module should
	// take care of importing the amount into the account except for the
	// genesis block
	if k.GetTPPModuleBalance(ctx).IsZero() {
		err := k.InitializeTPPModule(ctx, sdk.NewCoin("tpp", sdk.ZeroInt()))
		if err != nil {
			panic(err)
		}
	}

}



// GetTPPModuleAccount returns the module account.
func (k Keeper) GetTPPModuleAccount(ctx sdk.Context) (ModuleName authtypes.ModuleAccountI) {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// GetTPPModuleBalance returns the module account balance
func (k Keeper) GetTPPModuleBalance(ctx sdk.Context) sdk.Coin {
	return k.bankKeeper.GetBalance(ctx, k.GetTPPModuleAccount(ctx).GetAddress(), "tpp")
}

// InitializeTPPModule sets up the module account from genesis
func (k Keeper) InitializeTPPModule(ctx sdk.Context, funds sdk.Coin) error {
	return k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(funds))
}