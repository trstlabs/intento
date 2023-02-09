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

	available := k.GetRelayerRewardsAvailability(ctx)
	if !available {
		k.SetRelayerRewardsAvailability(ctx, true)
	}

	logger := k.Logger(ctx)

	autoTxs := k.GetAutoTxsForBlock(ctx)
	logger.Debug("auto_ibc-txs", "txs", len(autoTxs))

	for _, autoTx := range autoTxs {
		// check dependent txs
		if !k.AllowedToExecute(ctx, &autoTx) {
			updateAutoTxHistory(&autoTx, types.ErrAutoTxContinue)
			k.SetAutoTxInfo(ctx, &autoTx)
			continue
		}

		logger.Debug("auto_tx", "owner", autoTx.Owner)

		isRecurring := autoTx.ExecTime.Before(autoTx.EndTime)

		flexFee := calculateFlexFee(autoTx, isRecurring)
		if fee, err := k.DistributeCoins(ctx, autoTx, flexFee, isRecurring, req.Header.ProposerAddress); err != nil {
			logger.Error("auto_tx", "distribution err", err.Error())
			addAutoTxHistory(&autoTx, ctx.BlockHeader().Time, fee, err)

		} else {
			if err := k.SendAutoTx(ctx, autoTx); err != nil {
				logger.Error("auto_tx", "err", err)
				addAutoTxHistory(&autoTx, ctx.BlockHeader().Time, sdk.Coin{}, err)
			}
			addAutoTxHistory(&autoTx, ctx.BlockHeader().Time, fee)
			willRecur := autoTx.ExecTime.Before(autoTx.EndTime) && autoTx.ExecTime.Add(autoTx.Interval).Before(autoTx.EndTime)
			k.RemoveFromAutoTxQueue(ctx, autoTx)
			if willRecur {
				fmt.Printf("auto-tx will recur: %v \n", autoTx.TxID)
				// adding next execTime and a new entry into the queue based on interval
				autoTx.ExecTime = autoTx.ExecTime.Add(autoTx.Interval)
				k.InsertAutoTxQueue(ctx, autoTx.TxID, autoTx.ExecTime)
			}
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeAutoTx,
				sdk.NewAttribute(types.AttributeKeyAutoTxOwner, autoTx.Owner),
			),
		)

		k.SetAutoTxInfo(ctx, &autoTx)

	}
}

func addAutoTxHistory(autoTx *types.AutoTxInfo, actualExecTime time.Time, execFee sdk.Coin, err ...error) {
	historyEntry := types.AutoTxHistoryEntry{
		ScheduledExecTime: autoTx.ExecTime,
		ActualExecTime:    actualExecTime,
		ExecFee:           execFee,
	}
	if len(err) > 0 {
		historyEntry.Error = err[0].Error()
	}
	autoTx.AutoTxHistory = append(autoTx.AutoTxHistory, &historyEntry)

}

func updateAutoTxHistory(autoTx *types.AutoTxInfo, err error) {
	//autoTx.AutoTxHistory[len(autoTx.AutoTxHistory)-1].Retries = autoTx.AutoTxHistory[len(autoTx.AutoTxHistory)-1].Retries + 1
	autoTx.AutoTxHistory[len(autoTx.AutoTxHistory)-1].Error = err.Error()
}

func calculateFlexFee(autoTx types.AutoTxInfo, isRecurring bool) sdk.Int {
	if len(autoTx.AutoTxHistory) != 0 {
		prevEntry := autoTx.AutoTxHistory[len(autoTx.AutoTxHistory)-1].ActualExecTime
		return sdk.NewInt((autoTx.ExecTime.Sub(prevEntry)).Milliseconds())
	}
	return sdk.NewInt((autoTx.ExecTime.Sub(autoTx.StartTime)).Milliseconds())

}
