package keeper

import (
	"fmt"
	"os"
	"path/filepath"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/trstlabs/trst/x/item/types"
)

// InitGenesis initializes the trst module state
func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) {
	//k.SetParams(ctx, state.Params)

	k.SetParams(ctx, state.Params)
	k.InitializeContracts(ctx)
	// NOTE: since the item module is a module account, the auth module should
	// take care of importing the amount into the account except for the
	// genesis block
	if k.GetItemModuleBalance(ctx).IsZero() {
		err := k.InitializeItemModule(ctx, sdk.NewCoin("utrst", sdk.ZeroInt()))
		if err != nil {
			panic(err)
		}
		k.InitializeItemIncentiveModule(ctx, sdk.NewCoin("utrst", sdk.ZeroInt()))

	}

	for _, elem := range state.ItemList {

		k.SetItem(ctx, *elem)
		k.InsertListedItemQueue(ctx, *elem)
		//if (elem.Buyer != "") {k.SetBuyer(ctx, elem.Id, elem.Buyer)}

	}

	for _, elem := range state.ProfileList {

		k.SetProfile(ctx, *elem, elem.Owner)
	}

	// this line is used by starport scaffolding # genesis/module/init
	// Set all the estimator

	// Set estimator count
	//k.SetEstimationInfoCount(ctx, int64(len(state.ProfileList)))

	// Set buyer count
	//k.SetBuyerCount(ctx, int64(len(state.BuyerList)))

	// Set item count
	k.SetItemCount(ctx, int64(len(state.ItemList)))

}

// ExportGenesis exports the trst module state
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesis()

	itemList := k.GetAllItems(ctx)
	for _, elem := range itemList {
		elem := elem
		genesis.ItemList = append(genesis.ItemList, &elem)
	}

	profileList := k.GetAllProfiles(ctx)
	for _, elem := range profileList {
		elem := elem
		genesis.ProfileList = append(genesis.ProfileList, &elem)
	}

	genesis.Params = k.GetParams(ctx)

	/*	// Get all buyer
		buyerList := k.GetAllBuyer(ctx)
		for _, elem := range buyerList {
			elem := elem
			genesis.BuyerList = append(genesis.BuyerList, &elem)
		}
	*/

	return genesis
}

// GetItemModuleAccount returns the module account.
func (k Keeper) GetItemModuleAccount(ctx sdk.Context) (ModuleName authtypes.ModuleAccountI) {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// GetItemModuleBalance returns the module account balance
func (k Keeper) GetItemModuleBalance(ctx sdk.Context) sdk.Coin {
	return k.bankKeeper.GetBalance(ctx, k.GetItemModuleAccount(ctx).GetAddress(), "utrst")
}

// InitializeItemModule sets up the module account from genesis
func (k Keeper) InitializeItemModule(ctx sdk.Context, funds sdk.Coin) error {
	return k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(funds))
}

// InitializeItemIncentiveModule creates the module account for item incentives.
func (k Keeper) InitializeItemIncentiveModule(ctx sdk.Context, amount sdk.Coin) {

	moduleAcc := authtypes.NewEmptyModuleAccount(
		types.ItemIncentivesModuleAcctName, authtypes.Minter)

	k.accountKeeper.SetModuleAccount(ctx, moduleAcc)

	err := k.bankKeeper.MintCoins(ctx, types.ItemIncentivesModuleAcctName, sdk.NewCoins(amount))
	if err != nil {
		panic(err)
	}
}

// InitializeContracts sets up the module contracts from genesis
func (k Keeper) InitializeContracts(ctx sdk.Context) error {
	addr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	params := k.GetParams(ctx)
	// ensure reward pool module account is set
	if addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	userHomeDir, _ := os.UserHomeDir()

	transferCode, err := os.ReadFile(filepath.Join(userHomeDir, "trst", "wasm_code", "transfer_contract.wasm.gz"))
	if err != nil {
		panic(err)
	}
	estimateOnlyCode, err := os.ReadFile(filepath.Join(userHomeDir, "trst", "wasm_code", "estimate_only_contract.wasm.gz"))
	if err != nil {
		panic(err)
	}
	//var codeID uint64
	//var hash string

	_, err = k.computeKeeper.Create(ctx, addr, transferCode, "", "", params.MaxActivePeriod, "Estimation aggregation", " for Trustless Transfer items. This code is used internally. This code is used to enable indepedent pricing through aggregating estimations. Users send their estimations after which the creator can reveal the final price. The funds are automatically redistributed at the end of the contract period.")
	if err != nil {
		panic(err)
	}
	_, err = k.computeKeeper.Create(ctx, addr, estimateOnlyCode, "", "", params.MaxActivePeriod, "Estimation aggregation", " This code is used internally. This code is used to enable indepedent pricing through aggregating estimations. Users send their estimations after which the creator can reveal the final price. The funds are automatically redistributed at the end of the contract period.")
	if err != nil {
		panic(err)
	}

	//store := ctx.KVStore(k.storeKey)
	//store.Set([]byte(fmt.Sprint(codeID)), []byte(hash))

	return nil
}
