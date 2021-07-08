package keeper

import (
	"fmt"
	"os"
	"path/filepath"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/danieljdd/tpp/x/tpp/types"
)

// InitGenesis initializes the TPP module state
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

	k.InitializeContract(ctx)

	for _, elem := range state.ItemList {

		k.SetItem(ctx, *elem)
		k.InsertListedItemQueue(ctx, elem.Id, *elem, elem.Endtime)
		//if (elem.Buyer != "") {k.SetBuyer(ctx, elem.Id, elem.Buyer)}

	}

	for _, elem := range state.EstimatorList {

		k.SetEstimator(ctx, *elem)
	}

	k.SetParams(ctx, types.DefaultParams())
	// this line is used by starport scaffolding # genesis/module/init
	// Set all the estimator

	// Set estimator count
	k.SetEstimatorCount(ctx, int64(len(state.EstimatorList)))

	// Set buyer count
	//k.SetBuyerCount(ctx, int64(len(state.BuyerList)))

	// Set item count
	k.SetItemCount(ctx, int64(len(state.ItemList)))

}

// ExportGenesis exports the TPP module state
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesis()

	itemList := k.GetAllItem(ctx)
	for _, elem := range itemList {
		elem := elem
		genesis.ItemList = append(genesis.ItemList, &elem)
	}

	estimatorList := k.GetAllEstimator(ctx)
	for _, elem := range estimatorList {
		elem := elem
		genesis.EstimatorList = append(genesis.EstimatorList, &elem)
	}

	/*	// Get all buyer
		buyerList := k.GetAllBuyer(ctx)
		for _, elem := range buyerList {
			elem := elem
			genesis.BuyerList = append(genesis.BuyerList, &elem)
		}
	*/

	return genesis
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

// InitializeTPPModule sets up the module account from genesis
func (k Keeper) InitializeContract(ctx sdk.Context) error {
	addr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	// ensure reward pool module account is set
	if addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	userHomeDir, _ := os.UserHomeDir()

	wasm, err := os.ReadFile(filepath.Join(userHomeDir, "tpp", "contract.wasm.gz"))
	if err != nil {
		panic(err)
	}
	var codeID uint64
	var hash string

	codeID, err = k.computeKeeper.Create(ctx, addr, wasm, "", "")
	if err != nil {
		panic(err)
	}

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(fmt.Sprint(codeID)), []byte(hash))

	return nil
}
