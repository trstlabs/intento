package autoibctx

import (
	//"fmt"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/trstlabs/trst/x/auto-ibc-tx/keeper"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

// BeginBlocker called every block, processes auto execution
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {

	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	logger := k.Logger(ctx)

	autoTxs := k.GetAutoTxsForBlock(ctx)

	cacheCtx, writeCache := ctx.CacheContext()
	for _, autoTx := range autoTxs {
		// attempt to self-send interchain account transaction

		err := k.SendAutoTx(cacheCtx, autoTx)
		if err == nil {
			logger.Info(
				"auto_tx",
				"owner", autoTx.Owner.String(),
			)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeAutoTx,
					sdk.NewAttribute(types.AttributeKeyAutoTxOwner, autoTx.Owner.String()),
				),
			)

		}

		isRecurring := autoTx.ExecTime.Before(autoTx.EndTime)
		isLastExec := false
		var gasUsed uint64
		if isRecurring {
			isLastExec = autoTx.ExecTime.Add(autoTx.Interval).After(autoTx.EndTime)
			if len(autoTx.AutoTxHistory) != 0 {
				prevEntry := autoTx.AutoTxHistory[len(autoTx.AutoTxHistory)-1].ActualExecTime
				gasUsed = uint64(autoTx.ExecTime.Sub(prevEntry).Milliseconds())

			} else {
				gasUsed = uint64(autoTx.ExecTime.Sub(autoTx.StartTime).Milliseconds())
			}

		} else {
			isLastExec = true
			gasUsed = uint64(autoTx.ExecTime.Sub(autoTx.StartTime).Milliseconds())
		}

		//deducts execution fees and distributes SDK-native coins from autoTx balance
		fee, err := k.DistributeCoins(ctx, autoTx, gasUsed, isRecurring, isLastExec, req.Header.ProposerAddress)
		if err != nil {
			logger.Info("auto execution", "error", err.Error())
		} else {
			fmt.Printf("write to cache\n")

			// if the autoTx execution is recurring and successful, we add a new entry to the queue with current entry time + interval
			if isRecurring {
				fmt.Printf("auto-executed recurring: %v \n", autoTx.TxID)
				k.RemoveFromAutoTxQueue(ctx, autoTx.TxID, autoTx.ExecTime)
				fmt.Printf("exec Time %+v \n", autoTx.ExecTime)
				nextExecTime := autoTx.ExecTime.Add(autoTx.Interval)
				fmt.Printf("exec Time new %+v \n", nextExecTime)
				historyEntry := types.AutoTxHistoryEntry{ScheduledExecTime: autoTx.ExecTime, ActualExecTime: time.Now(), ExecFee: fee}
				autoTx.AutoTxHistory = append(autoTx.AutoTxHistory, &historyEntry)

				if nextExecTime.Before(autoTx.EndTime) {
					k.InsertAutoTxQueue(ctx, autoTx.TxID, nextExecTime)
					autoTx.ExecTime = nextExecTime
					k.SetAutoTxInfo(ctx, autoTx.TxID, &autoTx)
					writeCache()
					continue
				}

				k.SetAutoTxInfo(ctx, autoTx.TxID, &autoTx)
			}

			writeCache()
		}

		k.RemoveFromAutoTxQueue(ctx, autoTx.TxID, autoTx.ExecTime)
		//_ = k.Delete(ctx, autoTx.Address)
		logger.Info(
			"expired",
			"autoTx", autoTx.TxID,
		)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeAutoTxExpired,
				sdk.NewAttribute(types.AttributeKeyAutoTxOwner, autoTx.Owner.String()),
			),
		)

	}

}
