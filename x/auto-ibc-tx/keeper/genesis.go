package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"

	// authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	// "github.com/trstlabs/trst/x/AutoIbcTx/internal/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// InitGenesis sets supply information for genesis.
//
// CONTRACT: all types of accounts must have been already initialized/created
// InitGenesis initializes the trst module state
func (k Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState, msgHandler sdk.Handler) error {
	// NOTE: since the AutoIbcTx module is a module account, the auth module should
	// take care of importing the amount into the account except for the
	// genesis block
	fmt.Print("auto-ibc-tx InitGenesis... \n")
	if k.GetAutoIbcTxModuleBalance(ctx).IsZero() {
		err := k.InitializeAutoIbcTxModule(ctx, sdk.NewCoin(types.Denom, sdk.ZeroInt()))
		if err != nil {
			//commented out for tests
			panic(err)
		}
	}
	var maxTxID uint64
	for i, autoTxInfo := range gs.AutoTxInfos {
		err := k.importAutoTxInfo(ctx, autoTxInfo.TxID, autoTxInfo)
		if err != nil {
			return sdkerrors.Wrapf(err, "autoTxInfo %d with id: %d", i, autoTxInfo.TxID)
		}
		if autoTxInfo.TxID > maxTxID {
			maxTxID = autoTxInfo.TxID
		}
	}

	for i, seq := range gs.Sequences {
		err := k.importAutoIncrementID(ctx, seq.IDKey, seq.Value)
		if err != nil {
			return sdkerrors.Wrapf(err, "sequence number %d", i)
		}
	}

	// sanity check seq values
	if k.peekAutoIncrementID(ctx, types.KeyLastTxID) <= maxTxID {
		return sdkerrors.Wrapf(types.ErrInvalid, "seq %s must be greater %d ", string(types.KeyLastTxID), maxTxID)
	}

	fmt.Print("setting params...\n")
	k.SetParams(ctx, types.DefaultParams())

	return nil
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) *types.GenesisState {
	//var genState types.GenesisState
	genState := *types.DefaultGenesis()

	//genState.Params = keeper.GetParams(ctx)
	genState.Params = keeper.GetParams(ctx)

	keeper.IterateAutoTxInfos(ctx, func(txID uint64, info types.AutoTxInfo) bool {
		genState.AutoTxInfos = append(genState.AutoTxInfos, info)
		return false
	})

	for _, k := range [][]byte{types.KeyLastTxID, types.KeyLastTxAddrID} {
		genState.Sequences = append(genState.Sequences, types.Sequence{
			IDKey: k,
			Value: keeper.peekAutoIncrementID(ctx, k),
		})
	}

	return &genState
}

// GetAutoIbcTxModuleAccount returns the module account.
func (k Keeper) GetAutoIbcTxModuleAccount(ctx sdk.Context) (ModuleName authtypes.ModuleAccountI) {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// GetAutoIbcTxModuleBalance returns the module account balance
func (k Keeper) GetAutoIbcTxModuleBalance(ctx sdk.Context) sdk.Coin {
	return k.bankKeeper.GetBalance(ctx, k.GetAutoIbcTxModuleAccount(ctx).GetAddress(), "utrst")
}

// InitializeAutoIbcTxModule sets up the module account from genesis
func (k Keeper) InitializeAutoIbcTxModule(ctx sdk.Context, funds sdk.Coin) error {
	return k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(funds))
}
