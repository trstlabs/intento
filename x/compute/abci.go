package compute

import (
	//"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/trstlabs/trst/x/compute/internal/keeper"
	"github.com/trstlabs/trst/x/compute/internal/types"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {

	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	logger := k.Logger(ctx)

	//addressList := k.GetAllContractAddresses

	// delete inactive contracts from store and its deposits
	k.IterateContractQueue(ctx, ctx.BlockHeader().Time, func(contract types.ContractInfoWithAddress) bool {

		logger.Info(
			"contract was expired",
			"contract", contract.Address.String(),
		)

		if contract.ContractInfo.AutoMsg != nil {

			res, err := k.Execute(ctx, contract.Address, contract.Address, contract.ContractInfo.AutoMsg, sdk.NewCoins(sdk.NewCoin("utrst", sdk.ZeroInt())), contract.ContractInfo.CallbackSig)
			if err != nil {
				logger.Info(
					"contract",
					"err", err.Error(),
				)

				k.SetContractResult(ctx, contract.Address, &sdk.Result{Log: err.Error()})
				logger.Info(
					"contract",
					"err", err.Error(),
				)
			} else {
				k.SetContractResult(ctx, contract.Address, res)

			}
			err = k.ContractPayoutCreator(ctx, contract.Address)
			if err != nil {
				logger.Info(
					"contract payout creator",
					"err", err.Error(),
				)
			}
			k.RemoveFromContractQueue(ctx, contract.Address.String(), contract.ContractInfo.EndTime)
			_ = k.Delete(ctx, contract.Address)
			logger.Info(
				"deleted",
				"contract", contract.Address.String(),
			)

		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeContractExpired,
				sdk.NewAttribute(types.AttributeKeyContractAddr, contract.Address.String()),
			),
		)

		return false

	})

	return []abci.ValidatorUpdate{}
}
