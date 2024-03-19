package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/intent/types"

	// authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	// "github.com/trstlabs/intento/x/Intent/internal/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// InitGenesis sets supply information for genesis.
//
// CONTRACT: all types of accounts must have been already initialized/created
// InitGenesis initializes the intent module state
func (k Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState) error {

	k.SetParams(ctx, types.DefaultParams())
	// NOTE: since the Intent module is a module account, the auth module should
	// take care of importing the amount into the account except for the
	// genesis block
	//fmt.Print("intent InitGenesis... \n")
	if k.GetIntentModuleBalance(ctx).IsZero() {
		err := k.InitializeIntentModule(ctx, sdk.NewCoin(types.Denom, sdk.ZeroInt()))
		if err != nil {
			return err
		}
	}
	var maxID uint64
	for i, autoTxInfo := range gs.ActionInfos {
		err := k.importActionInfo(ctx, autoTxInfo.ID, autoTxInfo)
		if err != nil {
			return errorsmod.Wrapf(err, "autoTxInfo %d with id: %d", i, autoTxInfo.ID)
		}
		if autoTxInfo.ID > maxID {
			maxID = autoTxInfo.ID
		}
	}

	for i, seq := range gs.Sequences {
		err := k.importAutoIncrementID(ctx, seq.IDKey, seq.Value)
		if err != nil {
			return errorsmod.Wrapf(err, "sequence number %d", i)
		}
	}

	// sanity check seq values
	if k.peekAutoIncrementID(ctx, types.KeyLastID) <= maxID {
		return errorsmod.Wrapf(types.ErrInvalidGenesis, "seq %s must be greater %d ", string(types.KeyLastID), maxID)
	}

	return nil
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) *types.GenesisState {
	//var genState types.GenesisState
	genState := *types.DefaultGenesis()

	genState.Params = keeper.GetParams(ctx)

	keeper.IterateActionInfos(ctx, func(id uint64, info types.ActionInfo) bool {
		genState.ActionInfos = append(genState.ActionInfos, info)
		return false
	})

	for _, k := range [][]byte{types.KeyLastID, types.KeyLastTxAddrID} {
		genState.Sequences = append(genState.Sequences, types.Sequence{
			IDKey: k,
			Value: keeper.peekAutoIncrementID(ctx, k),
		})
	}

	return &genState
}

// GetIntentModuleAccount returns the module account.
func (k Keeper) GetIntentModuleAccount(ctx sdk.Context) (ModuleName authtypes.ModuleAccountI) {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// GetIntentModuleBalance returns the module account balance
func (k Keeper) GetIntentModuleBalance(ctx sdk.Context) sdk.Coin {
	return k.bankKeeper.GetBalance(ctx, k.GetIntentModuleAccount(ctx).GetAddress(), "uinto")
}

// InitializeIntentModule sets up the module account from genesis
func (k Keeper) InitializeIntentModule(ctx sdk.Context, funds sdk.Coin) error {
	return k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(funds))
}
