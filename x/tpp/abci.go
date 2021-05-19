package tpp

import (
	"fmt"
	"time"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/danieljdd/tpp/x/tpp/keeper"
	"github.com/danieljdd/tpp/x/tpp/types"
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
			key := item.Id + "-" + element
				k.DeleteEstimator(ctx, key)
			}
			k.DeleteItem(ctx, item.Id)
		// called when items become inactive
		//keeper.AfterItemFailedMinDeposit(ctx, proposal.ProposalId)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeItemExpired,
				sdk.NewAttribute(types.AttributeKeyItemID, fmt.Sprintf("%d", item.Id)),

			),
		)


		return false
	})

	return []abci.ValidatorUpdate{}
}