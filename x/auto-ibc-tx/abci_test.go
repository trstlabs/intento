package autoibctx

import (
	"fmt"
	"testing"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	keeper "github.com/trstlabs/trst/x/auto-ibc-tx/keeper"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

func TestBeginBlocker(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	configuration := types.ExecutionConfiguration{SaveMsgResponses: true}
	autoTx, sendToAddr := createTestSendAutoTx(ctx, configuration, keepers)
	err := autoTx.ValidateBasic()
	require.NoError(t, err)
	k := keepers.AutoIbcTxKeeper

	k.SetAutoTxInfo(ctx, &autoTx)
	k.InsertAutoTxQueue(ctx, autoTx.TxID, autoTx.ExecTime)

	ctx2 := createNextExecutionContext(ctx, autoTx.ExecTime)
	// test that autoTx was added to the queue
	queue := k.GetAutoTxsForBlock(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].TxID)

	fakeAutoTxExec(k, ctx2, autoTx)
	autoTx = k.GetAutoTxInfo(ctx2, autoTx.TxID)
	ctx3 := createNextExecutionContext(ctx2, autoTx.ExecTime)

	//queue in BeginBocker
	queue = k.GetAutoTxsForBlock(ctx3)

	// test that autoTx history was updated
	require.Equal(t, ctx3.BlockHeader().Time, queue[0].ExecTime)
	require.Equal(t, 1, len(queue[0].AutoTxHistory))
	require.Equal(t, ctx2.BlockHeader().Time, queue[0].AutoTxHistory[0].ScheduledExecTime)
	require.Equal(t, ctx2.BlockHeader().Time, queue[0].AutoTxHistory[0].ActualExecTime)
	require.NotNil(t, ctx3.BlockHeader().Time, queue[0].AutoTxHistory[0].MsgResponses[0].Value)
	require.Equal(t, keepers.BankKeeper.GetAllBalances(ctx3, sendToAddr)[0].Amount, sdk.NewInt(100))

}

func TestBeginBlockerStopOnSuccess(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	configuration := types.ExecutionConfiguration{StopOnSuccess: true}
	autoTx, _ := createTestSendAutoTx(ctx, configuration, keepers)
	err := autoTx.ValidateBasic()
	require.NoError(t, err)
	k := keepers.AutoIbcTxKeeper

	k.SetAutoTxInfo(ctx, &autoTx)
	k.InsertAutoTxQueue(ctx, autoTx.TxID, autoTx.ExecTime)

	ctx2 := createNextExecutionContext(ctx, autoTx.ExecTime)
	// test that autoTx was added to the queue
	queue := k.GetAutoTxsForBlock(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].TxID)
	// BeginBlocker logic
	// check dependent txs
	// setting new ExecTime and adding a new entry into the queue based on interval
	//fmt.Printf("auto-tx will recur: %v \n", autoTx.TxID)

	fakeAutoTxExec(k, ctx2, autoTx)
	autoTx = k.GetAutoTxInfo(ctx2, autoTx.TxID)
	ctx3 := createNextExecutionContext(ctx2, autoTx.ExecTime.Add(time.Hour))
	autoTx = k.GetAutoTxInfo(ctx3, autoTx.TxID)
	fmt.Printf("%v\n", autoTx.ExecTime)
	require.True(t, autoTx.ExecTime.Before(ctx3.BlockTime()))

}

func TestBeginBlockerStopOnFailure(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	configuration := types.ExecutionConfiguration{StopOnFailure: true}
	autoTx, _ := createBadAutoTx(ctx, configuration, keepers)
	err := autoTx.ValidateBasic()
	require.NoError(t, err)
	k := keepers.AutoIbcTxKeeper

	k.SetAutoTxInfo(ctx, &autoTx)
	k.InsertAutoTxQueue(ctx, autoTx.TxID, autoTx.ExecTime)

	ctx2 := createNextExecutionContext(ctx, autoTx.ExecTime)
	// test that autoTx was added to the queue
	queue := k.GetAutoTxsForBlock(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].TxID)

	// BeginBlocker logic
	// check dependent txs
	// setting new ExecTime and adding a new entry into the queue based on interval
	//fmt.Printf("auto-tx will recur: %v \n", autoTx.TxID)
	fakeAutoTxExec(k, ctx2, autoTx)
	autoTx = k.GetAutoTxInfo(ctx2, autoTx.TxID)
	ctx3 := createNextExecutionContext(ctx2, autoTx.ExecTime.Add(time.Hour))
	autoTx = k.GetAutoTxInfo(ctx3, autoTx.TxID)

	require.True(t, autoTx.ExecTime.Before(ctx3.BlockTime()))

}

func fakeAutoTxExec(k keeper.Keeper, ctx sdk.Context, autoTx types.AutoTxInfo) {
	autoTx = k.GetAutoTxInfo(ctx, autoTx.TxID)
	if !k.AllowedToExecute(ctx, autoTx) {
		addAutoTxHistory(&autoTx, ctx.BlockTime(), sdk.Coin{}, false, nil, types.ErrAutoTxConditions)
		autoTx.ExecTime = autoTx.ExecTime.Add(autoTx.Interval)
		k.SetAutoTxInfo(ctx, &autoTx)
	}
	isRecurring := autoTx.ExecTime.Before(autoTx.EndTime)

	flexFee := calculateTimeBasedFlexFee(autoTx)
	fee, err := k.DistributeCoins(ctx, autoTx, flexFee, isRecurring, ctx.BlockHeader().ProposerAddress)

	k.RemoveFromAutoTxQueue(ctx, autoTx)
	if err != nil {
		addAutoTxHistory(&autoTx, ctx.BlockTime(), fee, false, nil, err)
	} else {
		err, executedLocally, msgResponses := k.SendAutoTx(ctx, &autoTx)
		addAutoTxHistory(&autoTx, ctx.BlockTime(), fee, executedLocally, msgResponses, err)
		fmt.Printf("err: %v \n", err)
		shouldRecur := isRecurring && (autoTx.ExecTime.Add(autoTx.Interval).Before(autoTx.EndTime) || autoTx.ExecTime.Add(autoTx.Interval) == autoTx.EndTime)
		allowedToRecur := (!autoTx.Configuration.StopOnSuccess && !autoTx.Configuration.StopOnFailure) || autoTx.Configuration.StopOnSuccess && err != nil || autoTx.Configuration.StopOnFailure && err == nil
		fmt.Printf("%v %v\n", shouldRecur, allowedToRecur)
		if shouldRecur && allowedToRecur {

			autoTx.ExecTime = autoTx.ExecTime.Add(autoTx.Interval)
			k.InsertAutoTxQueue(ctx, autoTx.TxID, autoTx.ExecTime)
		}

		k.SetAutoTxInfo(ctx, &autoTx)
	}

}

func TestBeginBlockerStressTest(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	k := keepers.AutoIbcTxKeeper

	autoTxs := createTestSendAutoTxs(ctx, 100, keepers)

	for _, autoTx := range autoTxs {
		k.InsertAutoTxQueue(ctx, autoTx.TxID, autoTx.ExecTime)
		k.SetAutoTxInfo(ctx, &autoTx)
	}

	ctx2 := createNextExecutionContext(ctx, autoTxs[0].ExecTime)
	queue := k.GetAutoTxsForBlock(ctx2)

	// test that all autoTxs were added to the queue
	require.Equal(t, len(autoTxs), len(queue))

	// BeginBlocker logic
	for _, autoTx := range queue {
		fakeAutoTxExec(k, ctx2, autoTx)
	}

	// information for the next execution
	ctx3 := createNextExecutionContext(ctx, autoTxs[0].ExecTime.Add(autoTxs[0].Interval))
	queue = k.GetAutoTxsForBlock(ctx3)

	// test that autoTx history was updated for all entries
	for _, autoTx := range queue {
		require.Equal(t, 1, len(autoTx.AutoTxHistory))
	}
}

func TestOwnerMustBeSignerForLocalAutoTx(t *testing.T) {
	ctx, keepers, cdc := createTestContext(t)

	autoTxOwnerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	feeAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	toSendAcc, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
	require.Equal(t, keepers.BankKeeper.GetAllBalances(ctx, autoTxOwnerAddr)[0].Amount, sdk.NewInt(3_000_000_000_000))
	localMsg := &banktypes.MsgSend{
		FromAddress: toSendAcc.String(),
		ToAddress:   autoTxOwnerAddr.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
	}
	anys, _ := types.PackTxMsgAnys([]sdk.Msg{localMsg})

	autoTx := types.AutoTxInfo{
		TxID:       123,
		Owner:      autoTxOwnerAddr.String(),
		FeeAddress: feeAddr.String(),
		Msgs:       anys,
	}
	k := keepers.AutoIbcTxKeeper

	err := autoTx.GetTxMsgs(cdc)[0].ValidateBasic()
	require.NoError(t, err)

	feeBeforeFeeParams := calculateTimeBasedFlexFee(autoTx)

	fee, err := k.DistributeCoins(ctx, autoTx, feeBeforeFeeParams, true, ctx.BlockHeader().ProposerAddress)

	require.NoError(t, err)
	err, executedLocally, _ := k.SendAutoTx(ctx, &autoTx)
	require.Contains(t, err.Error(), "owner doesn't have permission to send this message: unauthorized")
	require.False(t, executedLocally)

	require.Equal(t, keepers.BankKeeper.GetAllBalances(ctx, feeAddr)[0].Amount, sdk.NewInt(3_000_000_000_000).Sub(fee.Amount))
}

func createTestContext(t *testing.T) (sdk.Context, keeper.TestKeepers, codec.Codec) {
	ctx, keepers, cdc := keeper.CreateTestInput(t, false)

	types.Denom = "stake"
	keepers.AutoIbcTxKeeper.SetParams(ctx, types.Params{
		AutoTxFundsCommission:      2,
		AutoTxConstantFee:          1_000_000,                 // 1trst
		AutoTxFlexFeeMul:           3,                         // 3*calculated time-based flexFee
		RecurringAutoTxConstantFee: 1_000_000,                 // 1trst
		MaxAutoTxDuration:          time.Hour * 24 * 366 * 10, // a little over 10 years
		MinAutoTxDuration:          time.Second * 60,
		MinAutoTxInterval:          time.Second * 20,
	})
	return ctx, keepers, cdc
}

func createTestSendAutoTx(ctx sdk.Context, configuration types.ExecutionConfiguration, keepers keeper.TestKeepers) (types.AutoTxInfo, sdk.AccAddress) {
	autoTxOwnerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	feeAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	toSendAcc, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
	startTime := ctx.BlockHeader().Time
	execTime := ctx.BlockHeader().Time.Add(time.Hour)
	endTime := ctx.BlockHeader().Time.Add(time.Hour * 2)
	localMsg := &banktypes.MsgSend{
		FromAddress: autoTxOwnerAddr.String(),
		ToAddress:   toSendAcc.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
	}
	anys, _ := types.PackTxMsgAnys([]sdk.Msg{localMsg})

	autoTx := types.AutoTxInfo{
		TxID:          123,
		Owner:         autoTxOwnerAddr.String(),
		FeeAddress:    feeAddr.String(),
		ExecTime:      execTime,
		EndTime:       endTime,
		Interval:      time.Hour,
		StartTime:     startTime,
		Msgs:          anys,
		Configuration: &configuration,
	}
	return autoTx, toSendAcc
}

func createBadAutoTx(ctx sdk.Context, configuration types.ExecutionConfiguration, keepers keeper.TestKeepers) (types.AutoTxInfo, sdk.AccAddress) {
	autoTxOwnerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
	feeAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	toSendAcc, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
	startTime := ctx.BlockHeader().Time
	execTime := ctx.BlockHeader().Time.Add(time.Hour)
	endTime := ctx.BlockHeader().Time.Add(time.Hour * 2)
	localMsg := &banktypes.MsgSend{
		FromAddress: autoTxOwnerAddr.String(),
		ToAddress:   toSendAcc.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
	}
	anys, _ := types.PackTxMsgAnys([]sdk.Msg{localMsg})

	autoTx := types.AutoTxInfo{
		TxID:          123,
		Owner:         autoTxOwnerAddr.String(),
		FeeAddress:    feeAddr.String(),
		ExecTime:      execTime,
		EndTime:       endTime,
		Interval:      time.Hour,
		StartTime:     startTime,
		Msgs:          anys,
		Configuration: &configuration,
	}
	return autoTx, toSendAcc
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

func createTestSendAutoTxs(ctx sdk.Context, count int, keepers keeper.TestKeepers) []types.AutoTxInfo {
	autoTxs := make([]types.AutoTxInfo, count)
	startTime := ctx.BlockHeader().Time
	execTime := ctx.BlockHeader().Time.Add(time.Hour)
	endTime := ctx.BlockHeader().Time.Add(time.Hour * 2)

	for i := 0; i < count; i++ {
		autoTxOwnerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
		feeAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
		toSendAcc, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
		localMsg := &banktypes.MsgSend{
			FromAddress: autoTxOwnerAddr.String(),
			ToAddress:   toSendAcc.String(),
			Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
		}
		anys, _ := types.PackTxMsgAnys([]sdk.Msg{localMsg})
		autoTxs[i] = types.AutoTxInfo{
			TxID:          uint64(i),
			Owner:         autoTxOwnerAddr.String(),
			FeeAddress:    feeAddr.String(),
			ExecTime:      execTime,
			EndTime:       endTime,
			Interval:      time.Hour,
			StartTime:     startTime,
			Msgs:          anys,
			Configuration: &types.ExecutionConfiguration{SaveMsgResponses: false},
		}
	}
	return autoTxs
}
