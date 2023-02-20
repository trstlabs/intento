package autoibctx

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	keeper "github.com/trstlabs/trst/x/auto-ibc-tx/keeper"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

func TestBeginBlocker(t *testing.T) {
	ctx, keepers := createTestContext(t)
	autoTx := createTestAutoTx(ctx, keepers)
	k := keepers.AutoIbcTxKeeper

	k.SetAutoTxInfo(ctx, &autoTx)
	k.InsertAutoTxQueue(ctx, autoTx.TxID, autoTx.ExecTime)

	ctx2 := createNextExecutionContext(ctx, autoTx.ExecTime)

	// test that autoTx was added to the queue
	queue := k.GetAutoTxsForBlock(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].TxID)

	// BeginBlocker
	isRecurring := autoTx.ExecTime.Before(autoTx.EndTime)
	flexFee := calculateFlexFee(autoTx, isRecurring)
	fee, err := k.DistributeCoins(ctx2, autoTx, flexFee, isRecurring, ctx2.BlockHeader().ProposerAddress)
	require.NoError(t, err)

	addAutoTxHistory(&autoTx, ctx2.BlockHeader().Time, fee)

	k.RemoveFromAutoTxQueue(ctx2, autoTx)
	autoTx.ExecTime = autoTx.ExecTime.Add(autoTx.Interval)
	k.InsertAutoTxQueue(ctx2, autoTx.TxID, autoTx.ExecTime)
	k.SetAutoTxInfo(ctx2, &autoTx)

	// in OnAcknowledgementPacket we check result and set Executed
	autoTx.AutoTxHistory[0].Executed = true

	// information for the next execution
	ctx3 := createNextExecutionContext(ctx2, autoTx.ExecTime)
	canExecute := k.AllowedToExecute(ctx, &autoTx)
	require.True(t, canExecute)

	//queue in BeginBocker
	queue = k.GetAutoTxsForBlock(ctx3)

	// test that autoTx history was updated
	require.Equal(t, 1, len(queue[0].AutoTxHistory))
	require.Equal(t, ctx2.BlockHeader().Time, queue[0].AutoTxHistory[0].ScheduledExecTime)
	require.Equal(t, ctx2.BlockHeader().Time, queue[0].AutoTxHistory[0].ActualExecTime)
	//require.Equal(t, uint64(0), queue[0].AutoTxHistory[0].Retries)
	require.Equal(t, ctx3.BlockHeader().Time, queue[0].ExecTime)
}

func TestBeginBlockerStressTest(t *testing.T) {
	ctx, keepers := createTestContext(t)
	//autoTx := createTestAutoTx(ctx, keepers)
	k := keepers.AutoIbcTxKeeper

	autoTxs := createTestAutoTxs(ctx, 100, keepers)

	for _, autoTx := range autoTxs {
		k.InsertAutoTxQueue(ctx, autoTx.TxID, autoTx.ExecTime)
		k.SetAutoTxInfo(ctx, &autoTx)
	}

	ctx2 := createNextExecutionContext(ctx, autoTxs[0].ExecTime)
	queue := k.GetAutoTxsForBlock(ctx2)

	// test that all autoTxs were added to the queue
	require.Equal(t, len(autoTxs), len(queue))

	// BeginBlocker
	for _, autoTx := range queue {
		isRecurring := autoTx.ExecTime.Before(autoTx.EndTime)
		flexFee := calculateFlexFee(autoTx, isRecurring)
		fee, err := k.DistributeCoins(ctx2, autoTx, flexFee, isRecurring, ctx2.BlockHeader().ProposerAddress)
		require.NoError(t, err)

		addAutoTxHistory(&autoTx, ctx2.BlockHeader().Time, fee)

		k.RemoveFromAutoTxQueue(ctx2, autoTx)
		k.InsertAutoTxQueue(ctx2, autoTx.TxID, autoTx.ExecTime.Add(autoTx.Interval))
		k.SetAutoTxInfo(ctx2, &autoTx)
	}

	// information for the next execution
	ctx3 := createNextExecutionContext(ctx, autoTxs[0].ExecTime.Add(autoTxs[0].Interval))
	queue = k.GetAutoTxsForBlock(ctx3)

	// test that autoTx history was updated for all entries
	for _, autoTx := range queue {
		require.Equal(t, 1, len(autoTx.AutoTxHistory))

	}
}

func TestBeginBlockerWithRetry(t *testing.T) {
	ctx, keepers := createTestContext(t)
	autoTx := createTestAutoTx(ctx, keepers)
	k := keepers.AutoIbcTxKeeper

	k.SetAutoTxInfo(ctx, &autoTx)
	k.InsertAutoTxQueue(ctx, autoTx.TxID, autoTx.ExecTime)

	ctx2 := createNextExecutionContext(ctx, autoTx.ExecTime)

	// test that autoTx was added to the queue
	queue := k.GetAutoTxsForBlock(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].TxID)

	// BeginBlocker
	isRecurring := autoTx.ExecTime.Before(autoTx.EndTime)
	flexFee := calculateFlexFee(autoTx, isRecurring)
	fee, err := k.DistributeCoins(ctx2, autoTx, flexFee, isRecurring, ctx2.BlockHeader().ProposerAddress)
	require.NoError(t, err)

	addAutoTxHistory(&autoTx, ctx2.BlockHeader().Time, fee)

	k.RemoveFromAutoTxQueue(ctx2, autoTx)
	autoTx.ExecTime = autoTx.ExecTime.Add(autoTx.Interval)
	k.InsertAutoTxQueue(ctx2, autoTx.TxID, autoTx.ExecTime)
	k.SetAutoTxInfo(ctx2, &autoTx)

	// information for the next execution
	ctx3 := createNextExecutionContext(ctx2, autoTx.ExecTime)
	queue = k.GetAutoTxsForBlock(ctx3)
	require.Equal(t, 1, len(queue[0].AutoTxHistory))
	require.Equal(t, ctx2.BlockHeader().Time, queue[0].AutoTxHistory[0].ScheduledExecTime)
	require.Equal(t, ctx3.BlockHeader().Time, queue[0].ExecTime)

	//We have no Executed from ibc_module.go so AllowedToExecute will reinsert the tx and update retry count
	canExecute := k.AllowedToExecute(ctx, &autoTx)
	require.True(t, canExecute)
	updateAutoTxHistory(&autoTx, types.ErrAutoTxContinue)
	k.SetAutoTxInfo(ctx, &autoTx)

	// information for the next execution
	ctx4 := createNextExecutionContext(ctx2, autoTx.ExecTime.Add(time.Second))
	queue = k.GetAutoTxsForBlock(ctx4)
	require.NotEmpty(t, queue)
	//require.Equal(t, uint64(1), queue[0].AutoTxHistory[0].Retries)

	// // test that autoTx history was updated

}

func createTestContext(t *testing.T) (sdk.Context, keeper.TestKeepers) {
	ctx, keepers := keeper.CreateTestInput(t, false)

	types.Denom = "stake"

	return ctx, keepers
}

func createTestAutoTx(ctx sdk.Context, keepers keeper.TestKeepers) types.AutoTxInfo {
	autoTxOwnerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	feeAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	startTime := ctx.BlockHeader().Time
	execTime := ctx.BlockHeader().Time.Add(time.Hour)
	endTime := ctx.BlockHeader().Time.Add(time.Hour * 2)
	autoTx := types.AutoTxInfo{
		TxID:       123,
		Owner:      autoTxOwnerAddr.String(),
		FeeAddress: feeAddr.String(),
		ExecTime:   execTime,
		EndTime:    endTime,
		Interval:   time.Hour,
		StartTime:  startTime,
		//MaxRetries: 2,
	}
	return autoTx
}

func createNextExecutionContext(ctx sdk.Context, nextExecTime time.Time) sdk.Context {
	return sdk.NewContext(ctx.MultiStore(), tmproto.Header{
		Height:          ctx.BlockHeight() + 1111,
		Time:            nextExecTime,
		ChainID:         ctx.ChainID(),
		ProposerAddress: ctx.BlockHeader().ProposerAddress,
	}, false, ctx.Logger())
}

type KeeperMock struct {
	AllowedToExecuteFunc      func(ctx sdk.Context, autoTx types.AutoTxInfo) bool
	SendAutoTxFunc            func(ctx sdk.Context, autoTx types.AutoTxInfo) error
	DistributeCoinsFunc       func(ctx sdk.Context, autoTx types.AutoTxInfo, flexFee uint64, isRecurring bool, isLastExec bool, proposer sdk.AccAddress) (uint64, error)
	RemoveFromAutoTxQueueFunc func(ctx sdk.Context, autoTxs ...types.AutoTxInfo)
	AddToAutoTxQueueFunc      func(ctx sdk.Context, autoTx types.AutoTxInfo)
	SetAutoTxInfoFunc         func(ctx sdk.Context, txID string, autoTx *types.AutoTxInfo)
}

func createTestAutoTxs(ctx sdk.Context, count int, keepers keeper.TestKeepers) []types.AutoTxInfo {
	autoTxs := make([]types.AutoTxInfo, count)
	startTime := ctx.BlockHeader().Time
	execTime := ctx.BlockHeader().Time.Add(time.Hour)
	endTime := ctx.BlockHeader().Time.Add(time.Hour * 2)
	for i := 0; i < count; i++ {
		autoTxOwnerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
		feeAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
		autoTxs[i] = types.AutoTxInfo{
			TxID:       uint64(i),
			Owner:      autoTxOwnerAddr.String(),
			FeeAddress: feeAddr.String(),
			ExecTime:   execTime,
			EndTime:    endTime,
			Interval:   time.Hour,
			StartTime:  startTime,
		}
	}
	return autoTxs
}
