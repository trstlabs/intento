package alloc

import (
	"context"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/x/alloc/keeper"
	"github.com/trstlabs/intento/x/alloc/types"
)

// BeginBlocker to distribute specific rewards on every begin block
func BeginBlocker(ctx context.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	if err := k.DistributeInflation(sdk.UnwrapSDKContext(ctx)); err != nil {
		panic(fmt.Sprintf("Error distribute inflation: %s", err.Error()))
	}

}
