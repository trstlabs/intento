package trst

import (
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/danieljdd/trst/x/item/keeper"
	"github.com/danieljdd/trst/x/item/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {

	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	logger := k.Logger(ctx)

	// delete inactive items from store and its deposits
	k.IterateListedItemsQueue(ctx, ctx.BlockHeader().Time, func(item types.Item) bool {
		logger.Info(
			"Item was expired",
			"item", item.Id,
			//"title", item.GetTitle(),
		)
		for _, estimator := range item.Estimatorlist {
			key := append(types.Uint64ToByte(item.Id), []byte(estimator)...)
			if item.Highestestimator == estimator && item.Transferable || item.Lowestestimator == estimator && !item.Transferable {
				k.DeleteEstimationWithoutDeposit(ctx, key)
			} else if item.Bestestimator == estimator {
				k.DeleteEstimationWithReward(ctx, key)
			} else {
				k.DeleteEstimation(ctx, key)
			}
		}
		errDelete := k.DeleteItemContract(ctx, item.Contract)
		if errDelete != nil {
			panic("error deleting item contract")
		}
		k.RemoveFromListedItemQueue(ctx, item.Id, item.Endtime)
		k.RemoveFromItemSeller(ctx, item.Id, item.Seller)
		k.DeleteItem(ctx, item.Id)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeItemExpired,
				sdk.NewAttribute(types.AttributeKeyItemID, strconv.FormatUint(item.Id, 10)),
			),
		)

		return false
	})

	return []abci.ValidatorUpdate{}
}
