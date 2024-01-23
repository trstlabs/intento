package keeper_test

import (
	//"fmt"
	"fmt"

	"time"

	sdkmath "cosmossdk.io/math"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/trst/x/auto-ibc-tx/keeper"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

// BeginBlocker called every block, processes auto execution
func FakeBeginBlocker(ctx sdk.Context, k keeper.Keeper, fakeProposer sdk.ConsAddress) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	available := k.GetRelayerRewardsAvailability(ctx)
	if !available {
		k.SetRelayerRewardsAvailability(ctx, true)
	}

	autoTxs := k.GetAutoTxsForBlock(ctx)

	for _, autoTx := range autoTxs {
		autoTx = k.GetAutoTxInfo(ctx, autoTx.TxID)
		if !k.AllowedToExecute(ctx, autoTx) {
			addAutoTxHistory(&autoTx, ctx.BlockTime(), sdk.Coin{}, false, nil, types.ErrAutoTxConditions)
			autoTx.ExecTime = autoTx.ExecTime.Add(autoTx.Interval)
			k.SetAutoTxInfo(ctx, &autoTx)
		}
		isRecurring := autoTx.ExecTime.Before(autoTx.EndTime)

		flexFee := calculateTimeBasedFlexFee(autoTx)
		fee, err := k.DistributeCoins(ctx, autoTx, flexFee, isRecurring, fakeProposer)

		k.RemoveFromAutoTxQueue(ctx, autoTx)
		if err != nil {
			fmt.Printf("err FakeBeginBlocker DistributeCoins: %v \n", err)
			errorString := fmt.Sprintf(types.ErrAutoTxDistribution, err.Error())
			addAutoTxHistory(&autoTx, ctx.BlockTime(), fee, false, nil, errorString)
		} else {
			err, executedLocally, msgResponses := k.SendAutoTx(ctx, &autoTx)
			if err != nil {
				addAutoTxHistory(&autoTx, ctx.BlockTime(), fee, executedLocally, msgResponses, fmt.Sprintf(types.ErrAutoTxMsgHandling, err.Error()))
			} else {
				addAutoTxHistory(&autoTx, ctx.BlockTime(), fee, executedLocally, msgResponses)
			}

			shouldRecur := isRecurring && (autoTx.ExecTime.Add(autoTx.Interval).Before(autoTx.EndTime) || autoTx.ExecTime.Add(autoTx.Interval) == autoTx.EndTime)
			allowedToRecur := (!autoTx.Configuration.StopOnSuccess && !autoTx.Configuration.StopOnFailure) || autoTx.Configuration.StopOnSuccess && err != nil || autoTx.Configuration.StopOnFailure && err == nil
			fmt.Printf("%v %v\n", shouldRecur, allowedToRecur)
			if shouldRecur && allowedToRecur {

				autoTx.ExecTime = autoTx.ExecTime.Add(autoTx.Interval)
				k.InsertAutoTxQueue(ctx, autoTx.TxID, autoTx.ExecTime)
			}

		}
		k.SetAutoTxInfo(ctx, &autoTx)
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

func calculateTimeBasedFlexFee(autoTx types.AutoTxInfo) sdkmath.Int {
	if len(autoTx.AutoTxHistory) != 0 {
		prevEntry := autoTx.AutoTxHistory[len(autoTx.AutoTxHistory)-1].ActualExecTime
		period := (autoTx.ExecTime.Sub(prevEntry))
		return sdk.NewInt(int64(period.Milliseconds()))
	}

	period := autoTx.ExecTime.Sub(autoTx.StartTime)
	return sdk.NewInt(int64(period.Milliseconds()))
}
