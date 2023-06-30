package compute

import (
	//"fmt"
	"fmt"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/trst/x/compute/internal/keeper"
	"github.com/trstlabs/trst/x/compute/internal/types"
)

// BeginBlocker called every block, processes auto execution
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {

	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	logger := k.Logger(ctx)

	contracts := k.GetContractAddressesForBlock(ctx)
	/*//var rewardCoins sdk.Coins
	if len(incentiveList) > 0 {
		k.SetIncentiveCoins(ctx, incentiveList)

	}*/
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
					"contract", contract.Address.String(),
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

		isRecurring := contract.ContractInfo.ExecTime.Before(contract.ContractInfo.EndTime)
		isLastExec := false
		if isRecurring {
			isLastExec = contract.ContractInfo.ExecTime.Add(contract.ContractInfo.Interval).After(contract.ContractInfo.EndTime)
		} else {
			isLastExec = true
		}
		//deducts execution fees and distributes SDK-native coins from contract balance
		fee, err := k.DistributeCoins(ctx, contract, gasUsed, isRecurring, isLastExec, req.Header.ProposerAddress)
		if err != nil {
			logger.Info("auto execution", "error", err.Error())
		} else {
			fmt.Printf("write to cache\n")

			// if the contract execution is recurring and successful, we add a new entry to the queue with current entry time + interval
			if isRecurring {
				fmt.Printf("auto-executed recurring :%s \n", contract.Address.String())
				k.RemoveFromContractQueue(ctx, contract.Address.String(), contract.ContractInfo.ExecTime)
				fmt.Printf("exec Time %+v \n", contract.ContractInfo.ExecTime)
				nextExecTime := contract.ContractInfo.ExecTime.Add(contract.ContractInfo.Interval)
				fmt.Printf("exec Time new %+v \n", nextExecTime)
				historyEntry := types.ExecHistoryEntry{ScheduledExecTime: contract.ContractInfo.ExecTime, ActualExecTime: time.Now(), ExecFee: fee}
				contract.ContractInfo.ExecHistory = append(contract.ContractInfo.ExecHistory, &historyEntry)

				if nextExecTime.Before(contract.ContractInfo.EndTime) {
					k.InsertContractQueue(ctx, contract.Address.String(), nextExecTime)
					contract.ContractInfo.ExecTime = nextExecTime
					k.SetContractInfo(ctx, contract)
					writeCache()
					continue
				}

				k.SetContractInfo(ctx, contract)
			}

			writeCache()
		}

		k.RemoveFromContractQueue(ctx, contract.Address.String(), contract.ContractInfo.ExecTime)
		//_ = k.Delete(ctx, contract.Address)
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

}
