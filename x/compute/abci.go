package compute

import (
	//"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/danieljdd/trst/x/compute/internal/keeper"
	"github.com/danieljdd/trst/x/compute/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {

	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	logger := k.Logger(ctx)
	//fmt.Printf("ABCI ENDBLOCK COMPUTE")
	// delete inactive items from store and its deposits
	k.IterateContractQueue(ctx, ctx.BlockHeader().Time, func(item types.ContractInfoWithAddress) bool {

		logger.Info(
			"contract was expired",
			"contract", item.Address.String(),
		)

		//err := k.CallLastMsg(ctx, item.Address)
		if item.LastMsg != nil {
			res, err := k.Execute(ctx, item.Address, item.Address, item.LastMsg, sdk.NewCoins(sdk.NewCoin("utrst", sdk.ZeroInt())), nil)
			if err != nil {
				_ = k.ContractPayout(ctx, item.Address)
				logger.Info(
					"Error lastMsg, creator payout", item.Address,
				)
			}
			k.SetContractResult(ctx, item.Address, res)
		}

		k.RemoveFromContractQueue(ctx, item.Address.String(), item.ContractInfo.EndTime)
		_ = k.Delete(ctx, item.Address)
		logger.Info(
			"Deleted contract", item.Address,
		)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeContractExpired,
				sdk.NewAttribute(types.AttributeKeyContractAddr, item.Address.String()),
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
