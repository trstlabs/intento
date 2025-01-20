package intent

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/intent/keeper"
	"github.com/trstlabs/intento/x/intent/types"
)

// BeginBlocker called every block, processes auto execution
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	ensureRelayerRewardsAvailable(ctx, k)

	logger := k.Logger(ctx)
	actions := k.GetActionsForBlock(ctx)

	timeOfBlock := ctx.BlockHeader().Time
	for _, action := range actions {
		// Check if ICQConfig is present and submit an interchain query if applicable
		if (action.Conditions.FeedbackLoops != nil && action.Conditions.FeedbackLoops[0].ICQConfig != nil) || (action.Conditions.Comparisons != nil && action.Conditions.Comparisons[0].ICQConfig != nil) {
			k.SubmitInterchainQueries(ctx, action, logger)
			k.RemoveFromActionQueue(ctx, action)
			// If the query is submitted, we skip handling this action for now
			continue

		}
		k.HandleAction(ctx, logger, action, timeOfBlock, nil)
	}
}

// ensureRelayerRewardsAvailable checks if relayer rewards are available and sets them if not
func ensureRelayerRewardsAvailable(ctx sdk.Context, k keeper.Keeper) {
	if !k.GetRelayerRewardsAvailability(ctx) {
		k.SetRelayerRewardsAvailability(ctx, true)
	}
}
