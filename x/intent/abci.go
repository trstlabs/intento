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
	flows := k.GetFlowsForBlock(ctx)

	timeOfBlock := ctx.BlockHeader().Time
	for _, flow := range flows {
		// Check if ICQConfig is present and submit an interchain query if applicable
		if (flow.Conditions.FeedbackLoops != nil && flow.Conditions.FeedbackLoops[0].ICQConfig != nil) || (flow.Conditions.Comparisons != nil && flow.Conditions.Comparisons[0].ICQConfig != nil) {
			k.SubmitInterchainQueries(ctx, flow, logger)
			k.RemoveFromFlowQueue(ctx, flow)
			// If the query is submitted, we skip handling this flow for now
			continue

		}
		k.HandleFlow(ctx, logger, flow, timeOfBlock, nil)
	}
}

// ensureRelayerRewardsAvailable checks if relayer rewards are available and sets them if not
func ensureRelayerRewardsAvailable(ctx sdk.Context, k keeper.Keeper) {
	if !k.GetRelayerRewardsAvailability(ctx) {
		k.SetRelayerRewardsAvailability(ctx, true)
	}
}
