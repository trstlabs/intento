package tpp

import (
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/danieljdd/tpp/x/tpp/keeper"
	"github.com/danieljdd/tpp/x/tpp/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {

	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	logger := k.Logger(ctx)

	// delete inactive items from store and its deposits
	k.IterateInactiveItemsQueue(ctx, ctx.BlockHeader().Time, func(item types.Item) bool {
		logger.Info(
			"Item was expired",
			"item", item.Id,
			"title", item.GetTitle(),
		)
		for _, element := range item.Estimatorlist {
			//apply this to each element
			key := append(types.Uint64ToByte(item.Id), []byte(element)...)

			k.DeleteEstimation(ctx, key)
		}
		k.DeleteItem(ctx, item.Id)
		k.RemoveFromItemSeller(ctx, item.Id, item.Seller)
		// called when items become inactive
		//keeper.AfterItemFailedMinDeposit(ctx, proposal.ProposalId)

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
