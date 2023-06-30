package autoibctx

import (
	//"fmt"
	"fmt"

	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
	logger.Debug("auto_txs", "amount", len(autoTxs))

	timeOfBlock := ctx.BlockHeader().Time
	for _, autoTx := range autoTxs {
		// check dependent txs
		if !k.AllowedToExecute(ctx, autoTx) {
			addAutoTxHistory(&autoTx, timeOfBlock, sdk.Coin{}, false, types.ErrAutoTxConditions)
			autoTx.ExecTime = autoTx.ExecTime.Add(autoTx.Interval)
			k.SetAutoTxInfo(ctx, &autoTx)
			continue
		}

		logger.Debug("auto_tx", "owner", autoTx.Owner)

		isRecurring := autoTx.ExecTime.Before(autoTx.EndTime)

		flexFee := calculateTimeBasedFlexFee(autoTx, isRecurring)
		fee, err := k.DistributeCoins(ctx, autoTx, flexFee, isRecurring, req.Header.ProposerAddress)
		if err != nil {
			logger.Error("auto_tx", "distribution err", err.Error())
			addAutoTxHistory(&autoTx, timeOfBlock, fee, false, err)
		} else {
			err, executedLocally := k.SendAutoTx(ctx, autoTx)
			addAutoTxHistory(&autoTx, timeOfBlock, fee, executedLocally, err)
		}

		k.RemoveFromAutoTxQueue(ctx, autoTx)
		// updagting ExecTime and adding a new entry into the queue based on interval
		willRecur := isRecurring && (autoTx.ExecTime.Add(autoTx.Interval).Before(autoTx.EndTime) || autoTx.ExecTime.Add(autoTx.Interval) == autoTx.EndTime)
		if willRecur {
			fmt.Printf("auto-tx will recur: %v \n", autoTx.TxID)
			autoTx.ExecTime = autoTx.ExecTime.Add(autoTx.Interval)
			k.InsertAutoTxQueue(ctx, autoTx.TxID, autoTx.ExecTime)
		}

		k.SetAutoTxInfo(ctx, &autoTx)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeAutoTx,
				sdk.NewAttribute(types.AttributeKeyAutoTxID, fmt.Sprint(autoTx.TxID)),
				sdk.NewAttribute(types.AttributeKeyAutoTxOwner, autoTx.Owner),
			),
		)
	}
}

func addAutoTxHistory(autoTx *types.AutoTxInfo, actualExecTime time.Time, execFee sdk.Coin, executedLocally bool, err ...error) {
	historyEntry := types.AutoTxHistoryEntry{
		ScheduledExecTime: autoTx.ExecTime,
		ActualExecTime:    actualExecTime,
		ExecFee:           execFee,
	}
	if len(err) == 1 && err[0] != nil {
		historyEntry.Error = err[0].Error()
	}
	if executedLocally {
		historyEntry.Executed = true
	}
	autoTx.AutoTxHistory = append(autoTx.AutoTxHistory, &historyEntry)
}

func calculateTimeBasedFlexFee(autoTx types.AutoTxInfo, isRecurring bool) sdk.Int {
	if len(autoTx.AutoTxHistory) != 0 {
		prevEntry := autoTx.AutoTxHistory[len(autoTx.AutoTxHistory)-1].ActualExecTime
		period := (autoTx.ExecTime.Sub(prevEntry))
		return sdk.NewInt(int64(period.Minutes()))
	}
	//return sdk.NewInt(int64((autoTx.ExecTime.Sub(autoTx.StartTime)).Minutes()))
	period := autoTx.ExecTime.Sub(autoTx.StartTime)
	if period.Seconds() <= 60 {
		//base fee so we do not have a zero fee
		return sdk.NewInt(1_000)
	}
	return sdk.NewInt(int64(period.Minutes()))
}
