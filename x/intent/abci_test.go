package intent

import (
	"fmt"
	"testing"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	keeper "github.com/trstlabs/intento/x/intent/keeper"
	"github.com/trstlabs/intento/x/intent/types"
)

func TestBeginBlocker(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	configuration := types.ExecutionConfiguration{SaveMsgResponses: true}
	autoTx, sendToAddr := createTestSendAction(ctx, configuration, keepers)
	err := autoTx.ValidateBasic()
	require.NoError(t, err)
	k := keepers.IntentKeeper

	k.SetActionInfo(ctx, &autoTx)
	k.InsertActionQueue(ctx, autoTx.ID, autoTx.ExecTime)

	ctx2 := createNextExecutionContext(ctx, autoTx.ExecTime)
	// test that autoTx was added to the queue
	queue := k.GetActionsForBlock(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].ID)

	fakeActionExec(k, ctx2, autoTx)
	autoTx = k.GetActionInfo(ctx2, autoTx.ID)
	ctx3 := createNextExecutionContext(ctx2, autoTx.ExecTime)

	//queue in BeginBocker
	queue = k.GetActionsForBlock(ctx3)
	autoTxHistory := k.GetActionHistory(ctx3, queue[0].ID)
	// test that autoTx history was updated
	require.Equal(t, ctx3.BlockHeader().Time, queue[0].ExecTime)
	require.Equal(t, 1, len(autoTxHistory.History))
	require.Equal(t, ctx2.BlockHeader().Time, autoTxHistory.History[0].ScheduledExecTime)
	require.Equal(t, ctx2.BlockHeader().Time, autoTxHistory.History[0].ActualExecTime)
	require.NotNil(t, ctx3.BlockHeader().Time, autoTxHistory.History[0].MsgResponses[0].Value)
	require.Equal(t, keepers.BankKeeper.GetAllBalances(ctx3, sendToAddr)[0].Amount, sdk.NewInt(100))

}

func TestBeginBlockerStopOnSuccess(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	configuration := types.ExecutionConfiguration{StopOnSuccess: true}
	autoTx, _ := createTestSendAction(ctx, configuration, keepers)
	err := autoTx.ValidateBasic()
	require.NoError(t, err)
	k := keepers.IntentKeeper

	k.SetActionInfo(ctx, &autoTx)
	k.InsertActionQueue(ctx, autoTx.ID, autoTx.ExecTime)

	ctx2 := createNextExecutionContext(ctx, autoTx.ExecTime)
	// test that autoTx was added to the queue
	queue := k.GetActionsForBlock(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].ID)
	// BeginBlocker logic
	// check dependent txs
	// setting new ExecTime and adding a new entry into the queue based on interval
	//fmt.Printf("action will recur: %v \n", autoTx.ID)

	fakeActionExec(k, ctx2, autoTx)
	autoTx = k.GetActionInfo(ctx2, autoTx.ID)
	ctx3 := createNextExecutionContext(ctx2, autoTx.ExecTime.Add(time.Hour))
	autoTx = k.GetActionInfo(ctx3, autoTx.ID)
	fmt.Printf("%v\n", autoTx.ExecTime)
	require.True(t, autoTx.ExecTime.Before(ctx3.BlockTime()))

}

func TestBeginBlockerStopOnFailure(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	configuration := types.ExecutionConfiguration{StopOnFailure: true}
	autoTx, _ := createBadAction(ctx, configuration, keepers)
	err := autoTx.ValidateBasic()
	require.NoError(t, err)
	k := keepers.IntentKeeper

	k.SetActionInfo(ctx, &autoTx)
	k.InsertActionQueue(ctx, autoTx.ID, autoTx.ExecTime)

	ctx2 := createNextExecutionContext(ctx, autoTx.ExecTime)
	// test that autoTx was added to the queue
	queue := k.GetActionsForBlock(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].ID)

	// BeginBlocker logic
	// check dependent txs
	// setting new ExecTime and adding a new entry into the queue based on interval
	//fmt.Printf("action will recur: %v \n", autoTx.ID)
	fakeActionExec(k, ctx2, autoTx)
	autoTx = k.GetActionInfo(ctx2, autoTx.ID)
	ctx3 := createNextExecutionContext(ctx2, autoTx.ExecTime.Add(time.Hour))
	autoTx = k.GetActionInfo(ctx3, autoTx.ID)

	require.True(t, autoTx.ExecTime.Before(ctx3.BlockTime()))

}

func TestErrorDistributionIsSaved(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	configuration := types.ExecutionConfiguration{StopOnFailure: true}
	autoTx, emptyBalanceAcc := createTestSendAction(ctx, configuration, keepers)

	err := autoTx.ValidateBasic()
	require.NoError(t, err)
	k := keepers.IntentKeeper

	k.SetActionInfo(ctx, &autoTx)
	k.InsertActionQueue(ctx, autoTx.ID, autoTx.ExecTime)

	ctx2 := createNextExecutionContext(ctx, autoTx.ExecTime)
	// test that autoTx was added to the queue
	queue := k.GetActionsForBlock(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].ID)
	err = sendTokens(ctx, keepers, autoTx.Owner, emptyBalanceAcc, sdk.NewInt64Coin("stake", 3_000_000_000_000))
	require.NoError(t, err)
	err = sendTokens(ctx, keepers, autoTx.FeeAddress, emptyBalanceAcc, sdk.NewInt64Coin("stake", 3_000_000_000_000))
	require.NoError(t, err)
	fakeActionExec(k, ctx2, autoTx)

	autoTx = k.GetActionInfo(ctx2, autoTx.ID)
	ctx3 := createNextExecutionContext(ctx2, autoTx.ExecTime.Add(time.Hour))
	autoTx = k.GetActionInfo(ctx3, autoTx.ID)
	autoTxHistory := k.GetActionHistory(ctx3, queue[0].ID)

	require.True(t, autoTx.ExecTime.Before(ctx3.BlockTime()))
	require.NotNil(t, autoTxHistory.History[0].Errors)
	require.Contains(t, autoTxHistory.History[0].Errors[0], "distribution")

}

func fakeActionExec(k keeper.Keeper, ctx sdk.Context, autoTx types.ActionInfo) {
	timeOfBlock := ctx.BlockHeader().Time
	autoTx = k.GetActionInfo(ctx, autoTx.ID)
	autoTxHistory, _ := k.TryGetActionHistory(ctx, autoTx.ID)
	if !k.AllowedToExecute(ctx, autoTx) {
		k.AddActionHistory(ctx, &autoTxHistory, &autoTx, timeOfBlock, sdk.Coin{}, false, nil, types.ErrActionConditions)
		autoTx.ExecTime = autoTx.ExecTime.Add(autoTx.Interval)
		k.SetActionInfo(ctx, &autoTx)
	}
	isRecurring := autoTx.ExecTime.Before(autoTx.EndTime)

	flexFee := calculateTimeBasedFlexFee(autoTx, autoTxHistory)
	fee, err := k.DistributeCoins(ctx, autoTx, flexFee, isRecurring, ctx.BlockHeader().ProposerAddress)

	k.RemoveFromActionQueue(ctx, autoTx)
	if err != nil {
		errorString := fmt.Sprintf(types.ErrActionDistribution, err.Error())
		k.AddActionHistory(ctx, &autoTxHistory, &autoTx, timeOfBlock, fee, false, nil, errorString)
	} else {
		err, executedLocally, msgResponses := k.SendAction(ctx, &autoTx)
		if err != nil {
			k.AddActionHistory(ctx, &autoTxHistory, &autoTx, ctx.BlockTime(), fee, executedLocally, msgResponses, fmt.Sprintf(types.ErrActionMsgHandling, err.Error()))
		} else {
			k.AddActionHistory(ctx, &autoTxHistory, &autoTx, ctx.BlockTime(), fee, executedLocally, msgResponses)
		}

		shouldRecur := isRecurring && (autoTx.ExecTime.Add(autoTx.Interval).Before(autoTx.EndTime) || autoTx.ExecTime.Add(autoTx.Interval) == autoTx.EndTime)
		allowedToRecur := (!autoTx.Configuration.StopOnSuccess && !autoTx.Configuration.StopOnFailure) || autoTx.Configuration.StopOnSuccess && err != nil || autoTx.Configuration.StopOnFailure && err == nil
		fmt.Printf("%v %v\n", shouldRecur, allowedToRecur)
		if shouldRecur && allowedToRecur {

			autoTx.ExecTime = autoTx.ExecTime.Add(autoTx.Interval)
			k.InsertActionQueue(ctx, autoTx.ID, autoTx.ExecTime)
		}

	}
	k.SetActionInfo(ctx, &autoTx)

}

func TestOwnerMustBeSignerForLocalAction(t *testing.T) {
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

	autoTx := types.ActionInfo{
		ID:         123,
		Owner:      autoTxOwnerAddr.String(),
		FeeAddress: feeAddr.String(),
		Msgs:       anys,
	}
	k := keepers.IntentKeeper

	err := autoTx.GetTxMsgs(cdc)[0].ValidateBasic()
	require.NoError(t, err)

	autoTxHistory, _ := k.TryGetActionHistory(ctx, autoTx.ID)
	flexFee := calculateTimeBasedFlexFee(autoTx, autoTxHistory)

	fee, err := k.DistributeCoins(ctx, autoTx, flexFee, true, ctx.BlockHeader().ProposerAddress)

	require.NoError(t, err)
	err, executedLocally, _ := k.SendAction(ctx, &autoTx)
	require.Contains(t, err.Error(), "owner doesn't have permission to send this message: unauthorized")
	require.False(t, executedLocally)

	require.Equal(t, keepers.BankKeeper.GetAllBalances(ctx, feeAddr)[0].Amount, sdk.NewInt(3_000_000_000_000).Sub(fee.Amount))
}

func createTestContext(t *testing.T) (sdk.Context, keeper.TestKeepers, codec.Codec) {
	ctx, keepers, cdc := keeper.CreateTestInput(t, false)

	types.Denom = "stake"
	keepers.IntentKeeper.SetParams(ctx, types.Params{
		ActionFundsCommission:      2,
		ActionConstantFee:          1_000_000,                 // 1trst
		ActionFlexFeeMul:           3,                         // 3*calculated time-based flexFee
		RecurringActionConstantFee: 1_000_000,                 // 1trst
		MaxActionDuration:          time.Hour * 24 * 366 * 10, // a little over 10 years
		MinActionDuration:          time.Second * 60,
		MinActionInterval:          time.Second * 20,
	})
	return ctx, keepers, cdc
}

func createTestSendAction(ctx sdk.Context, configuration types.ExecutionConfiguration, keepers keeper.TestKeepers) (types.ActionInfo, sdk.AccAddress) {
	autoTxOwnerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	fundedFeeAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	emptyBalanceAcc, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
	startTime := ctx.BlockHeader().Time
	execTime := ctx.BlockHeader().Time.Add(time.Hour)
	endTime := ctx.BlockHeader().Time.Add(time.Hour * 2)
	localMsg := &banktypes.MsgSend{
		FromAddress: autoTxOwnerAddr.String(),
		ToAddress:   emptyBalanceAcc.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
	}
	anys, _ := types.PackTxMsgAnys([]sdk.Msg{localMsg})

	autoTx := types.ActionInfo{
		ID:            123,
		Owner:         autoTxOwnerAddr.String(),
		FeeAddress:    fundedFeeAddr.String(),
		ExecTime:      execTime,
		EndTime:       endTime,
		Interval:      time.Hour,
		StartTime:     startTime,
		Msgs:          anys,
		Configuration: &configuration,
		ICAConfig:     &types.ICAConfig{},
	}
	return autoTx, emptyBalanceAcc
}

func sendTokens(ctx sdk.Context, keepers keeper.TestKeepers, from string, toAddr sdk.AccAddress, amount sdk.Coin) error {
	fromAddr, _ := sdk.AccAddressFromBech32(from)
	err := keepers.BankKeeper.SendCoins(ctx, fromAddr, toAddr, sdk.NewCoins(amount))

	return err
}

func createBadAction(ctx sdk.Context, configuration types.ExecutionConfiguration, keepers keeper.TestKeepers) (types.ActionInfo, sdk.AccAddress) {
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

	autoTx := types.ActionInfo{
		ID:            123,
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
	AllowedToExecuteFunc      func(ctx sdk.Context, autoTx types.ActionInfo) bool
	SendActionFunc            func(ctx sdk.Context, autoTx types.ActionInfo) error
	DistributeCoinsFunc       func(ctx sdk.Context, autoTx types.ActionInfo, flexFee uint64, isRecurring bool, isLastExec bool, proposer sdk.AccAddress) (uint64, error)
	RemoveFromActionQueueFunc func(ctx sdk.Context, autoTxs ...types.ActionInfo)
	AddToActionQueueFunc      func(ctx sdk.Context, autoTx types.ActionInfo)
	SetActionInfoFunc         func(ctx sdk.Context, id string, autoTx *types.ActionInfo)
}

func createTestSendActions(ctx sdk.Context, count int, keepers keeper.TestKeepers) []types.ActionInfo {
	autoTxs := make([]types.ActionInfo, count)
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
		autoTxs[i] = types.ActionInfo{
			ID:            uint64(i),
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
