package intent

import (
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/intent/keeper"
	"github.com/trstlabs/intento/x/intent/types"
)

// BeginBlocker called every block, processes auto execution
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	ensureRelayerRewardsAvailable(ctx, k)

	logger := k.Logger(ctx)
	actions := k.GetActionsForBlock(ctx)

	timeOfBlock := ctx.BlockHeader().Time
	for _, action := range actions {
		// Check if ICQConfig is present and submit an interchain query if applicable
		if action.Conditions != nil && action.Conditions.ICQConfig != nil {
			k.SubmitInterchainQuery(ctx, action, logger)
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
