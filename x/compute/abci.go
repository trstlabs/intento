package compute

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/danieljdd/tpp/x/compute/internal/keeper"
	"github.com/danieljdd/tpp/x/compute/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {

	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	logger := k.Logger(ctx)
	fmt.Printf("ABCI ENDBLOCK COMPUTE")
	// delete inactive items from store and its deposits
	/*k.IterateContractsQueue(ctx, ctx.BlockHeader().Time, func(contract types.ContractInfoWithAddress) bool {
		fmt.Printf(" ENDBLOCK QUEue")
		logger.Info(
			"contract was expired",
			"contract", contract.ContractInfo.CodeID,
			//"title", item.GetTitle(),
		)
		fmt.Printf("Deleting")
		err := k.Delete(ctx, contract.Address)
		if err != nil {
			fmt.Printf("contract.ContractInfo.CodeID")
			fmt.Printf("contract.ContractInfo.CodeID  is:  %s ", strconv.FormatUint(contract.ContractInfo.CodeID, 10))
			return false
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeContractExpired,
				sdk.NewAttribute(types.AttributeKeyCodeID, strconv.FormatUint(contract.ContractInfo.CodeID, 10)),
			),
		)
		return false
	})*/
	k.IterateContractQueue(ctx, ctx.BlockHeader().Time, func(addr sdk.AccAddress) bool {
		fmt.Printf(" ENDBLOCK QUEue")
		logger.Info(
			"contract was expired",
			"contract", addr.String(),
			//"title", item.GetTitle(),
		)
		fmt.Printf("Deleting")
		err := k.Delete(ctx, addr)
		if err != nil {
			//fmt.Printf("contract.ContractInfo.CodeID")
			fmt.Printf(".AttributeKeyContractAddr  is:  %s ", addr.String())
			return false
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeContractExpired,
				sdk.NewAttribute(types.AttributeKeyContractAddr, addr.String()),
			),
		)
		return false
	})

	return []abci.ValidatorUpdate{}
}
