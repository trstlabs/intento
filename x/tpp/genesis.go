package tpp
/*
import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/danieljdd/tpp/x/tpp/keeper"
	"github.com/danieljdd/tpp/x/tpp/types"
)

// InitGenesis initializes the TPP module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetParams(ctx, types.DefaultParams())
	// this line is used by starport scaffolding # genesis/module/init
	// Set all the estimator
	for _, elem := range genState.EstimatorList {
		k.SetEstimator(ctx, *elem)
	}

	// Set estimator count
	k.SetEstimatorCount(ctx, int64(len(genState.EstimatorList)))

	// Set all the buyer
	for _, elem := range genState.BuyerList {
		k.SetBuyer(ctx, *elem)
	}

	// Set buyer count
	k.SetBuyerCount(ctx, int64(len(genState.BuyerList)))

	// Set all the item
	for _, elem := range genState.ItemList {
		k.SetItem(ctx, *elem)
	}

	// Set item count
	k.SetItemCount(ctx, int64(len(genState.ItemList)))

}

// ExportGenesis returns the TPP module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export
	// Get all estimator
	estimatorList := k.GetAllEstimator(ctx)
	for _, elem := range estimatorList {
		elem := elem
		genesis.EstimatorList = append(genesis.EstimatorList, &elem)
	}

	// Get all buyer
	buyerList := k.GetAllBuyer(ctx)
	for _, elem := range buyerList {
		elem := elem
		genesis.BuyerList = append(genesis.BuyerList, &elem)
	}

	// Get all item
	itemList := k.GetAllItem(ctx)
	for _, elem := range itemList {
		elem := elem
		genesis.ItemList = append(genesis.ItemList, &elem)
	}

	return genesis
}
*/