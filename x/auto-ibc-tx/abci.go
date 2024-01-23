package autoibctx

import (
	"fmt"

	"time"

	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
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

	timeOfBlock := ctx.BlockHeader().Time
	for _, autoTx := range autoTxs {
		// check dependent txs
		if !k.AllowedToExecute(ctx, autoTx) {
			addAutoTxHistory(&autoTx, timeOfBlock, sdk.Coin{}, false, nil, types.ErrAutoTxConditions)
			autoTx.ExecTime = autoTx.ExecTime.Add(autoTx.Interval)
			k.SetAutoTxInfo(ctx, &autoTx)
			continue
		}

		logger.Debug("auto_tx execution", "id", autoTx.TxID)

		isRecurring := autoTx.ExecTime.Before(autoTx.EndTime)

		flexFee := calculateTimeBasedFlexFee(autoTx)
		fee, err := k.DistributeCoins(ctx, autoTx, flexFee, isRecurring, req.Header.ProposerAddress)

		k.RemoveFromAutoTxQueue(ctx, autoTx)
		if err != nil {
			errorString := fmt.Sprintf(types.ErrAutoTxDistribution, err.Error())
			addAutoTxHistory(&autoTx, timeOfBlock, fee, false, nil, errorString)
		} else {
			err, executedLocally, msgResponses := k.SendAutoTx(ctx, &autoTx)
			if err != nil {
				addAutoTxHistory(&autoTx, ctx.BlockTime(), fee, executedLocally, msgResponses, fmt.Sprintf(types.ErrAutoTxMsgHandling, err.Error()))
			} else {
				addAutoTxHistory(&autoTx, ctx.BlockTime(), fee, executedLocally, msgResponses)
			}

			// setting new ExecTime and adding a new entry into the queue based on interval
			shouldRecur := isRecurring && (autoTx.ExecTime.Add(autoTx.Interval).Before(autoTx.EndTime) || autoTx.ExecTime.Add(autoTx.Interval) == autoTx.EndTime)
			allowedToRecur := (!autoTx.Configuration.StopOnSuccess && !autoTx.Configuration.StopOnFailure) || autoTx.Configuration.StopOnSuccess && err != nil || autoTx.Configuration.StopOnFailure && err == nil

			if shouldRecur && allowedToRecur {
				//fmt.Printf("auto-tx will recur: %v \n", autoTx.TxID)
				autoTx.ExecTime = autoTx.ExecTime.Add(autoTx.Interval)
				k.InsertAutoTxQueue(ctx, autoTx.TxID, autoTx.ExecTime)
			}
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

func addAutoTxHistory(autoTx *types.AutoTxInfo, actualExecTime time.Time, execFee sdk.Coin, executedLocally bool, msgResponses []*cdctypes.Any, err ...string) {
	historyEntry := types.AutoTxHistoryEntry{
		ScheduledExecTime: autoTx.ExecTime,
		ActualExecTime:    actualExecTime,
		ExecFee:           execFee,
	}

	if executedLocally {
		historyEntry.Executed = true
		historyEntry.MsgResponses = msgResponses
	}

	if len(err) == 1 && err[0] != "" {
		historyEntry.Errors = append(historyEntry.Errors, err[0])
		historyEntry.Executed = false
	}
	autoTx.AutoTxHistory = append(autoTx.AutoTxHistory, &historyEntry)
}

// we may reimplement this as a configuration-based gas fee
func calculateTimeBasedFlexFee(autoTx types.AutoTxInfo) sdkmath.Int {
	if len(autoTx.AutoTxHistory) != 0 {
		prevEntry := autoTx.AutoTxHistory[len(autoTx.AutoTxHistory)-1].ActualExecTime
		period := (autoTx.ExecTime.Sub(prevEntry))
		return sdk.NewInt(int64(period.Milliseconds()))
	}

	period := autoTx.ExecTime.Sub(autoTx.StartTime)
	if period.Seconds() <= 60 {
		//base fee so we do not have a zero fee
		return sdk.NewInt(60_000)
	}
	return sdk.NewInt(int64(period.Milliseconds()))
}
