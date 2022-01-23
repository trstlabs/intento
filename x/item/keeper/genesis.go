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

	k.SetParams(ctx, types.DefaultParams())
	k.InitializeContract(ctx)
	// NOTE: since the item module is a module account, the auth module should
	// take care of importing the amount into the account except for the
	// genesis block
	if k.GetItemModuleBalance(ctx).IsZero() {
		err := k.InitializeItemModule(ctx, sdk.NewCoin("utrst", sdk.ZeroInt()))
		if err != nil {
			panic(err)
		}
	}

	for _, elem := range state.ItemList {

		k.SetItem(ctx, *elem)
		k.InsertListedItemQueue(ctx, elem.Id, *elem, elem.EndTime)
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

	itemList := k.GetAllItem(ctx)
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

// InitializeContract sets up the module account from genesis
func (k Keeper) InitializeContract(ctx sdk.Context) error {
	addr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	params := k.GetParams(ctx)
	// ensure reward pool module account is set
	if addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	userHomeDir, _ := os.UserHomeDir()

	wasm, err := os.ReadFile(filepath.Join(userHomeDir, "trst", "contract.wasm.gz"))
	if err != nil {
		panic(err)
	}
	var codeID uint64
	var hash string

	codeID, err = k.computeKeeper.Create(ctx, addr, wasm, "", "", params.MaxActivePeriod, "Estimation aggregation", "This code is used to enable indepedent pricing through aggregating estimations. Users send their estimations after which the creator can reveal the final price. The funds are automatically sent back at the end of the contract period. This code is used internally.")
	if err != nil {
		panic(err)
	}

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(fmt.Sprint(codeID)), []byte(hash))

	return nil
}
