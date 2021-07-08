package compute

import (
	//"fmt"
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
	//fmt.Printf("ABCI ENDBLOCK COMPUTE")
	// delete inactive items from store and its deposits
	k.IterateContractQueue(ctx, ctx.BlockHeader().Time, func(addr sdk.AccAddress) bool {
		//	fmt.Printf(" ENDBLOCK QUEue")
		logger.Info(
			"contract was expired",
			"contract", addr.String(),
		)
		//fmt.Printf("Deleting")
		err := k.CallLastMsg(ctx, addr)

		if err != nil {
			logger.Info(
				"Calling last msg unsuccesful", err,
			)
			err = k.ContractPayout(ctx, addr)
			logger.Info(
				"Contract creator payout", err,
			)
		}
		k.RemoveFromContractQueue(ctx, addr.String(), ctx.BlockHeader().Time)
		err = k.Delete(ctx, addr)
		logger.Info(
			"Deleted", err,
		)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeContractExpired,
				sdk.NewAttribute(types.AttributeKeyContractAddr, addr.String()),
			),
		)
		return false

		/*if err != nil {
			//fmt.Printf("contract.ContractInfo.CodeID")
			//fmt.Printf(".AttributeKeyContractAddr  is:  %s ", addr.String())
			return false
		}
		*/

	})

	return []abci.ValidatorUpdate{}
}
