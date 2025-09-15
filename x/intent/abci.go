package intent

import (
	"time"

	storetypes "cosmossdk.io/store/types"
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

	flows := k.GetFlowsForBlockAndPruneQueue(ctx)
	if len(flows) == 0 {
		return
	}

	ctx = ctx.WithGasMeter(storetypes.NewGasMeter(types.MaxGasTotal))

	for _, flow := range flows {
		// Estimate a gas cost before actually running it
		estimatedGas := uint64(100_000)
		if ctx.GasMeter().Limit()-ctx.GasMeter().GasConsumed() < estimatedGas {
			logger.Info("Skipping remaining flows due to block gas limit")
			break
		}

		// Check if ICQConfig is present and submit an interchain query if applicable
		if (len(flow.Conditions.FeedbackLoops) > 0 && flow.Conditions.FeedbackLoops[0].ICQConfig != nil) ||
			(len(flow.Conditions.Comparisons) > 0 && flow.Conditions.Comparisons[0].ICQConfig != nil) {
			k.SubmitInterchainQueries(ctx, flow, logger)
			// If the query is submitted, we skip handling this flow for now
			continue
		}
		k.HandleFlow(ctx, logger, flow, ctx.BlockHeader().Time)
	}
}

// ensureRelayerRewardsAvailable checks if relayer rewards are available and sets them if not
func ensureRelayerRewardsAvailable(ctx sdk.Context, k keeper.Keeper) {
	if !k.GetRelayerRewardsAvailability(ctx) {
		k.SetRelayerRewardsAvailability(ctx, true)
	}
}
