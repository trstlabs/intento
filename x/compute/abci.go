package compute

import (
	//"fmt"
	"fmt"
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

	incentiveList, contracts := k.GetContractAddressesForBlock(ctx)
	//var rewardCoins sdk.Coins
	if len(incentiveList) > 0 {
		k.SetIncentiveCoins(ctx, incentiveList)

	}
	var gasUsed uint64
	cacheCtx, writeCache := ctx.CacheContext()
	for _, contract := range contracts {
		if contract.Address.Equals(contract.Address) && contract.ContractInfo.AutoMsg != nil {
			// attempt to self-execute
			// AutoMessage may mutate state thus we can use a cached context. If one of
			// the handlers fails, no state mutation is written and the error
			// message is logged.

			gas, err := k.SelfExecute(cacheCtx, contract.Address, contract.ContractInfo.AutoMsg, contract.ContractInfo.CallbackSig)
			if err == nil {
				logger.Info(
					"auto_msg",
					"gas", gas,
				)

				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EventTypeAutoMsgContract,
						sdk.NewAttribute(types.AttributeKeyContractAddr, contract.Address.String()),
					),
				)

			}
			gasUsed = gas

		}

		logger.Info(
			"executed",
			"contract", contract.Address.String(),
		)

		isRecurring := contract.ContractInfo.ExecTime.Before(contract.ContractInfo.EndTime)
		//deducts execution fees and distributes SDK-native coins from contract balance
		err := k.DistributeCoins(ctx, contract, gasUsed, isRecurring)
		if err != nil {
			//fmt.Printf("couldnt deduct fee %s\n", err)
			logger.Info(
				"contract payout creator",
				"err", err.Error(),
			)

		} else {
			//fmt.Printf("write to cache\n")
			writeCache()
			// if the contract execution is recurring and successful, we add a new entry to the queue with current entry time + interval
			if isRecurring {
				fmt.Printf("self-executed recurring :%s \n", contract.Address.String())
				k.RemoveFromContractQueue(ctx, contract.Address.String(), contract.ContractInfo.ExecTime)
				fmt.Printf("exec Time %+v \n", contract.ContractInfo.ExecTime)
				nextExecTime := contract.ContractInfo.ExecTime.Add(contract.ContractInfo.Interval)
				fmt.Printf("exec Time2 %+v \n", nextExecTime)
				if nextExecTime.Before(contract.ContractInfo.EndTime) {
					k.InsertContractQueue(ctx, contract.Address.String(), nextExecTime)
					contract.ContractInfo.ExecTime = nextExecTime
					k.SetContractInfo(ctx, contract)

					continue
				}

			}
		}
		//fmt.Printf("executed \n")
		k.RemoveFromContractQueue(ctx, contract.Address.String(), contract.ContractInfo.ExecTime)
		_ = k.Delete(ctx, contract.Address)
		logger.Info(
			"expired",
			"contract", contract.Address.String(),
		)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeContractExpired,
				sdk.NewAttribute(types.AttributeKeyContractAddr, contract.Address.String()),
			),
		)
	}

	return []abci.ValidatorUpdate{}
}
