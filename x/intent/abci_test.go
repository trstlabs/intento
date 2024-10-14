package intent

import (
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
	configuration := types.ExecutionConfiguration{SaveResponses: true}
	action, sendToAddr := createTestTriggerAction(ctx, configuration, keepers)
	err := action.ValidateBasic()
	require.NoError(t, err)
	k := keepers.IntentKeeper

	k.SetActionInfo(ctx, &action)
	k.InsertActionQueue(ctx, action.ID, action.ExecTime)

	ctx2 := createNextExecutionContext(ctx, action.ExecTime)
	// test that action was added to the queue
	queue := k.GetActionsForBlock(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].ID)

	k.HandleAction(ctx2, k.Logger(ctx2), action, ctx2.BlockTime(), nil)
	action = k.GetActionInfo(ctx2, action.ID)
	ctx3 := createNextExecutionContext(ctx2, action.ExecTime)

	//queue in BeginBocker
	queue = k.GetActionsForBlock(ctx3)
	actionHistory := k.MustGetActionHistory(ctx3, queue[0].ID)
	// test that action history was updated
	require.Equal(t, ctx3.BlockHeader().Time, queue[0].ExecTime)
	require.Equal(t, 1, len(actionHistory))
	require.Equal(t, ctx2.BlockHeader().Time, actionHistory[0].ScheduledExecTime)
	require.Equal(t, ctx2.BlockHeader().Time, actionHistory[0].ActualExecTime)
	require.NotNil(t, ctx3.BlockHeader().Time, actionHistory[0].MsgResponses[0].Value)

	require.Equal(t, keepers.BankKeeper.GetAllBalances(ctx3, sendToAddr)[0].Amount, sdk.NewInt(100))

}

func TestBeginBlockerStopOnSuccess(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	configuration := types.ExecutionConfiguration{StopOnSuccess: true}
	action, _ := createTestTriggerAction(ctx, configuration, keepers)
	err := action.ValidateBasic()
	require.NoError(t, err)
	k := keepers.IntentKeeper

	k.SetActionInfo(ctx, &action)
	k.InsertActionQueue(ctx, action.ID, action.ExecTime)

	ctx2 := createNextExecutionContext(ctx, action.ExecTime)
	// test that action was added to the queue
	queue := k.GetActionsForBlock(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].ID)
	// BeginBlocker logic
	k.HandleAction(ctx2, k.Logger(ctx2), action, ctx.BlockTime(), nil)
	action = k.GetActionInfo(ctx2, action.ID)
	ctx3 := createNextExecutionContext(ctx2, action.ExecTime.Add(time.Hour))
	action = k.GetActionInfo(ctx3, action.ID)
	require.True(t, action.ExecTime.Before(ctx3.BlockTime()))

}

func TestBeginBlockerStopOnFailure(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	configuration := types.ExecutionConfiguration{StopOnFailure: true}
	action, _ := createBadAction(ctx, configuration, keepers)
	err := action.ValidateBasic()
	require.NoError(t, err)
	k := keepers.IntentKeeper

	k.SetActionInfo(ctx, &action)
	k.InsertActionQueue(ctx, action.ID, action.ExecTime)

	ctx2 := createNextExecutionContext(ctx, action.ExecTime)
	// test that action was added to the queue
	queue := k.GetActionsForBlock(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].ID)

	k.HandleAction(ctx2, k.Logger(ctx2), action, ctx.BlockTime(), nil)
	action = k.GetActionInfo(ctx2, action.ID)
	ctx3 := createNextExecutionContext(ctx2, action.ExecTime.Add(time.Hour))
	action = k.GetActionInfo(ctx3, action.ID)

	require.True(t, action.ExecTime.Before(ctx3.BlockTime()))

}

func TestErrorIsSavedToActionInfo(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	configuration := types.ExecutionConfiguration{StopOnFailure: true}
	action, emptyBalanceAcc := createTestTriggerAction(ctx, configuration, keepers)

	err := action.ValidateBasic()
	require.NoError(t, err)
	k := keepers.IntentKeeper

	k.SetActionInfo(ctx, &action)
	k.InsertActionQueue(ctx, action.ID, action.ExecTime)

	ctx2 := createNextExecutionContext(ctx, action.ExecTime)
	// test that action was added to the queue
	queue := k.GetActionsForBlock(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].ID)
	err = sendTokens(ctx, keepers, action.Owner, emptyBalanceAcc, sdk.NewInt64Coin("stake", 3_000_000_000_000))
	require.NoError(t, err)
	err = sendTokens(ctx, keepers, action.FeeAddress, emptyBalanceAcc, sdk.NewInt64Coin("stake", 3_000_000_000_000))
	require.NoError(t, err)
	k.HandleAction(ctx2, k.Logger(ctx2), action, ctx.BlockTime(), nil)

	action = k.GetActionInfo(ctx2, action.ID)
	ctx3 := createNextExecutionContext(ctx2, action.ExecTime.Add(time.Hour))
	action = k.GetActionInfo(ctx3, action.ID)
	actionHistory := k.MustGetActionHistory(ctx3, queue[0].ID)

	require.True(t, action.ExecTime.Before(ctx3.BlockTime()))
	require.NotNil(t, actionHistory[0].Errors)
	require.Contains(t, actionHistory[0].Errors[0], "balance too low")

}

func TestOwnerMustBeSignerForLocalAction(t *testing.T) {
	ctx, keepers, cdc := createTestContext(t)

	actionOwnerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	feeAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	toSendAcc, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
	require.Equal(t, keepers.BankKeeper.GetAllBalances(ctx, actionOwnerAddr)[0].Amount, sdk.NewInt(3_000_000_000_000))
	localMsg := &banktypes.MsgSend{
		FromAddress: toSendAcc.String(),
		ToAddress:   actionOwnerAddr.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
	}
	anys, _ := types.PackTxMsgAnys([]sdk.Msg{localMsg})

	action := types.ActionInfo{
		ID:         123,
		Owner:      actionOwnerAddr.String(),
		FeeAddress: feeAddr.String(),
		Msgs:       anys,
	}
	k := keepers.IntentKeeper

	err := action.GetTxMsgs(cdc)[0].ValidateBasic()
	require.NoError(t, err)

	fee, err := k.DistributeCoins(ctx, action, feeAddr, types.Denom, ctx.BlockHeader().ProposerAddress)

	require.NoError(t, err)
	executedLocally, _, err := k.TriggerAction(ctx, &action)
	require.Contains(t, err.Error(), "owner doesn't have permission to send this message: unauthorized")
	require.False(t, executedLocally)

	require.Equal(t, keepers.BankKeeper.GetAllBalances(ctx, feeAddr)[0].Amount, sdk.NewInt(3_000_000_000_000).Sub(fee.Amount))
}

func createTestContext(t *testing.T) (sdk.Context, keeper.TestKeepers, codec.Codec) {
	ctx, keepers, cdc := keeper.CreateTestInput(t, false)

	types.Denom = "stake"
	keepers.IntentKeeper.SetParams(ctx, types.Params{
		ActionFundsCommission: 2,
		ActionConstantFee:     1_000_000,                 // 1trst
		ActionFlexFeeMul:      3,                         //
		MaxActionDuration:     time.Hour * 24 * 366 * 10, // a little over 10 years
		MinActionDuration:     time.Second * 60,
		MinActionInterval:     time.Second * 20,
		GasFeeCoins:           sdk.NewCoins(sdk.NewCoin(types.Denom, sdk.OneInt())),
	})
	return ctx, keepers, cdc
}

func createTestTriggerAction(ctx sdk.Context, configuration types.ExecutionConfiguration, keepers keeper.TestKeepers) (types.ActionInfo, sdk.AccAddress) {
	actionOwnerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	fundedFeeAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	emptyBalanceAcc, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
	startTime := ctx.BlockHeader().Time
	execTime := ctx.BlockHeader().Time.Add(time.Hour)
	endTime := ctx.BlockHeader().Time.Add(time.Hour * 2)
	localMsg := &banktypes.MsgSend{
		FromAddress: actionOwnerAddr.String(),
		ToAddress:   emptyBalanceAcc.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
	}
	anys, _ := types.PackTxMsgAnys([]sdk.Msg{localMsg})

	action := types.ActionInfo{
		ID:            123,
		Owner:         actionOwnerAddr.String(),
		FeeAddress:    fundedFeeAddr.String(),
		ExecTime:      execTime,
		EndTime:       endTime,
		Interval:      time.Hour,
		StartTime:     startTime,
		Msgs:          anys,
		Configuration: &configuration,
		ICAConfig:     &types.ICAConfig{},
	}
	return action, emptyBalanceAcc
}

func sendTokens(ctx sdk.Context, keepers keeper.TestKeepers, from string, toAddr sdk.AccAddress, amount sdk.Coin) error {
	fromAddr, _ := sdk.AccAddressFromBech32(from)
	err := keepers.BankKeeper.SendCoins(ctx, fromAddr, toAddr, sdk.NewCoins(amount))

	return err
}

func createBadAction(ctx sdk.Context, configuration types.ExecutionConfiguration, keepers keeper.TestKeepers) (types.ActionInfo, sdk.AccAddress) {
	actionOwnerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
	feeAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	toSendAcc, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
	startTime := ctx.BlockHeader().Time
	execTime := ctx.BlockHeader().Time.Add(time.Hour)
	endTime := ctx.BlockHeader().Time.Add(time.Hour * 2)
	localMsg := &banktypes.MsgSend{
		FromAddress: actionOwnerAddr.String(),
		ToAddress:   toSendAcc.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
	}
	anys, _ := types.PackTxMsgAnys([]sdk.Msg{localMsg})

	action := types.ActionInfo{
		ID:            123,
		Owner:         actionOwnerAddr.String(),
		FeeAddress:    feeAddr.String(),
		ExecTime:      execTime,
		EndTime:       endTime,
		Interval:      time.Hour,
		StartTime:     startTime,
		Msgs:          anys,
		Configuration: &configuration,
	}
	return action, toSendAcc
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
	allowedToExecuteFunc      func(ctx sdk.Context, action types.ActionInfo) bool
	TriggerActionFunc         func(ctx sdk.Context, action types.ActionInfo) error
	DistributeCoinsFunc       func(ctx sdk.Context, action types.ActionInfo, flexFee uint64, isRecurring bool, isLastExec bool, proposer sdk.AccAddress) (uint64, error)
	RemoveFromActionQueueFunc func(ctx sdk.Context, actions ...types.ActionInfo)
	AddToActionQueueFunc      func(ctx sdk.Context, action types.ActionInfo)
	SetActionInfoFunc         func(ctx sdk.Context, id string, action *types.ActionInfo)
}
