package trst

import (
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/trstlabs/trst/x/item/keeper"
	"github.com/trstlabs/trst/x/item/types"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {

	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	logger := k.Logger(ctx)

	// delete inactive items from store and its deposits
	k.IterateListedItemsByEndTime(ctx, ctx.BlockHeader().Time, func(item types.Item) bool {
		logger.Info(
			"Item was expired",
			"item", item.Id,
			//"title", item.GetTitle(),
		)

		k.RemoveFromListedItemQueue(ctx, item.Id, item.ListingDuration.EndTime)
		k.RemoveFromSellerItems(ctx, item.Id, item.Transfer.Seller)
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
