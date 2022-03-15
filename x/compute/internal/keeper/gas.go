package keeper

import (

	//"encoding/json"

	//"log"

	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	wasm "github.com/trstlabs/trst/go-cosmwasm"

	"github.com/trstlabs/trst/x/compute/internal/types"
)

func gasForContract(ctx sdk.Context) uint64 {
	meter := ctx.GasMeter()
	remaining := (meter.Limit() - meter.GasConsumed()) * types.GasMultiplier
	if remaining > types.MaxGas {
		return types.MaxGas
	}
	return remaining
}

func consumeGas(ctx sdk.Context, gas uint64) {
	consumed := (gas / types.GasMultiplier) + 1
	ctx.GasMeter().ConsumeGas(consumed, "wasm contract")
	// throw OutOfGas error if we ran out (got exactly to zero due to better limit enforcing)
	if ctx.GasMeter().IsOutOfGas() {
		fmt.Printf("Out of Gas \n")
		panic(sdk.ErrorOutOfGas{})
	}
}

// MultipliedGasMeter wraps the GasMeter from context and multiplies all reads by out defined multiplier
type MultipiedGasMeter struct {
	originalMeter sdk.GasMeter
}

var _ wasm.GasMeter = MultipiedGasMeter{}

func (m MultipiedGasMeter) GasConsumed() sdk.Gas {
	return m.originalMeter.GasConsumed() * types.GasMultiplier
}

func gasMeter(ctx sdk.Context) MultipiedGasMeter {
	return MultipiedGasMeter{
		originalMeter: ctx.GasMeter(),
	}
}
