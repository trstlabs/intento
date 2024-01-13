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
		fmt.Println("FAKE BEGIN BLOCKER")
		logger.Debug("autotx execution", "id", autoTx.TxID)

		isRecurring := autoTx.ExecTime.Before(autoTx.EndTime)

		flexFee := calculateTimeBasedFlexFee(autoTx, isRecurring)
		fee, err := k.DistributeCoins(ctx, autoTx, flexFee, isRecurring, fakeProposer)

		k.RemoveFromAutoTxQueue(ctx, autoTx)
		if err != nil {
			fmt.Println("auto_tx", "distribution err", err.Error())
			addAutoTxHistory(&autoTx, timeOfBlock, fee, false, nil, err)
		} else {
			err, executedLocally, msgResponses := k.SendAutoTx(ctx, &autoTx)
			addAutoTxHistory(&autoTx, timeOfBlock, fee, executedLocally, msgResponses, err)
			if err != nil {
				fmt.Printf("execution error: %v \n", err.Error())
			}
			// setting new ExecTime and adding a new entry into the queue based on interval
			willRecur := isRecurring && (autoTx.ExecTime.Add(autoTx.Interval).Before(autoTx.EndTime) || autoTx.ExecTime.Add(autoTx.Interval) == autoTx.EndTime)
			if willRecur {
				fmt.Printf("auto-tx will recur: %v \n", autoTx.TxID)
				autoTx.ExecTime = autoTx.ExecTime.Add(autoTx.Interval)
				k.InsertAutoTxQueue(ctx, autoTx.TxID, autoTx.ExecTime)
			}

			k.SetAutoTxInfo(ctx, &autoTx)
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeAutoTx,
				sdk.NewAttribute(types.AttributeKeyAutoTxID, fmt.Sprint(autoTx.TxID)),
				sdk.NewAttribute(types.AttributeKeyAutoTxOwner, autoTx.Owner),
			),
		)
	}
}

func addAutoTxHistory(autoTx *types.AutoTxInfo, actualExecTime time.Time, execFee sdk.Coin, executedLocally bool, msgResponses []*cdctypes.Any, err ...error) {
	historyEntry := types.AutoTxHistoryEntry{
		ScheduledExecTime: autoTx.ExecTime,
		ActualExecTime:    actualExecTime,
		ExecFee:           execFee,
	}
	if len(err) == 1 && err[0] != nil {
		historyEntry.Errors = append(historyEntry.Errors, err[0].Error())
	}
	if executedLocally {
		historyEntry.Executed = true
		historyEntry.MsgResponses = msgResponses
	}
	autoTx.AutoTxHistory = append(autoTx.AutoTxHistory, &historyEntry)
}

func calculateTimeBasedFlexFee(autoTx types.AutoTxInfo, isRecurring bool) sdkmath.Int {
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
