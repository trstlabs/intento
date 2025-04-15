package intent

import (
	"testing"
	"time"

	"cosmossdk.io/math"
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
	flow, sendToAddr := createTestFlow(ctx, configuration, keepers)
	err := flow.ValidateBasic()
	require.NoError(t, err)
	k := keepers.IntentKeeper

	k.SetFlowInfo(ctx, &flow)
	k.InsertFlowQueue(ctx, flow.ID, flow.ExecTime)

	ctx2 := createNextExecutionContext(ctx, flow.ExecTime)
	// test that flow was added to the queue
	queue := k.GetFlowsForBlock(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].ID)

	k.HandleFlow(ctx2, k.Logger(ctx2), flow, ctx2.BlockTime(), nil)
	flow = k.GetFlowInfo(ctx2, flow.ID)
	ctx3 := createNextExecutionContext(ctx2, flow.ExecTime)

	//queue in BeginBocker
	queue = k.GetFlowsForBlock(ctx3)
	flowHistory := k.MustGetFlowHistory(ctx3, queue[0].ID)
	// test that flow history was updated
	require.Equal(t, ctx3.BlockHeader().Time, queue[0].ExecTime)
	require.Equal(t, 1, len(flowHistory))
	require.Equal(t, ctx2.BlockHeader().Time, flowHistory[0].ScheduledExecTime)
	require.Equal(t, ctx2.BlockHeader().Time, flowHistory[0].ActualExecTime)
	require.NotNil(t, ctx3.BlockHeader().Time, flowHistory[0].MsgResponses[0].Value)

	require.Equal(t, keepers.BankKeeper.GetAllBalances(ctx3, sendToAddr)[0].Amount, math.NewInt(100))

}

func TestBeginBlockerLoad(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	types.MaxGasTotal = 200_000
	configuration := types.ExecutionConfiguration{SaveResponses: true}
	k := keepers.IntentKeeper
	flow, _ := createTestFlow(ctx, configuration, keepers)
	for i := range 10 {
		flow, _ := createTestFlow(ctx, configuration, keepers)
		flow.ID = uint64(i)
		k.SetFlowInfo(ctx, &flow)
		k.InsertFlowQueue(ctx, flow.ID, flow.ExecTime)
	}
	ctx2 := createNextExecutionContext(ctx, flow.ExecTime)
	queue := k.GetFlowsForBlock(ctx2)
	BeginBlocker(ctx2, k)
	queue2 := k.GetFlowsForBlock(ctx2)
	ctx3 := createNextExecutionContext(ctx2, ctx2.BlockTime().Add(time.Minute))
	queue3 := k.GetFlowsForBlock(ctx3)
	require.NotEqual(t, len(queue), len(queue2))
	require.Equal(t, len(queue2), len(queue3)) //should be the same in later block if not executed

	BeginBlocker(ctx3, k)
	queue4 := k.GetFlowsForBlock(ctx3)
	require.NotEqual(t, len(queue4), 0)
}

func TestBeginBlockerStopOnSuccess(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	configuration := types.ExecutionConfiguration{StopOnSuccess: true}
	flow, _ := createTestFlow(ctx, configuration, keepers)
	err := flow.ValidateBasic()
	require.NoError(t, err)
	k := keepers.IntentKeeper

	k.SetFlowInfo(ctx, &flow)
	k.InsertFlowQueue(ctx, flow.ID, flow.ExecTime)

	ctx2 := createNextExecutionContext(ctx, flow.ExecTime)
	// test that flow was added to the queue
	queue := k.GetFlowsForBlock(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].ID)
	// BeginBlocker logic
	k.HandleFlow(ctx2, k.Logger(ctx2), flow, ctx.BlockTime(), nil)
	flow = k.GetFlowInfo(ctx2, flow.ID)
	ctx3 := createNextExecutionContext(ctx2, flow.ExecTime.Add(time.Hour))
	flow = k.GetFlowInfo(ctx3, flow.ID)
	require.True(t, flow.ExecTime.Before(ctx3.BlockTime()))

}

func TestBeginBlockerStopOnFailure(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	configuration := types.ExecutionConfiguration{StopOnFailure: true}
	flow, _ := createBadFlow(ctx, configuration, keepers)
	err := flow.ValidateBasic()
	require.NoError(t, err)
	k := keepers.IntentKeeper

	k.SetFlowInfo(ctx, &flow)
	k.InsertFlowQueue(ctx, flow.ID, flow.ExecTime)

	ctx2 := createNextExecutionContext(ctx, flow.ExecTime)
	// test that flow was added to the queue
	queue := k.GetFlowsForBlock(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].ID)

	k.HandleFlow(ctx2, k.Logger(ctx2), flow, ctx.BlockTime(), nil)
	flow = k.GetFlowInfo(ctx2, flow.ID)
	ctx3 := createNextExecutionContext(ctx2, flow.ExecTime.Add(time.Hour))
	flow = k.GetFlowInfo(ctx3, flow.ID)
	//exec time is not updated
	require.True(t, flow.ExecTime.Before(ctx3.BlockTime()))
	//flow also not in queue
	iterator := k.FlowQueueIterator(ctx, ctx3.BlockTime())
	require.Nil(t, iterator.Key())
}

func TestBeginBlockerAlwaysStopOnLowBalance(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	configuration := types.ExecutionConfiguration{StopOnFailure: false}
	flow, _ := createNoFeeFlow(ctx, configuration, keepers)
	err := flow.ValidateBasic()
	require.NoError(t, err)
	k := keepers.IntentKeeper

	k.SetFlowInfo(ctx, &flow)
	k.InsertFlowQueue(ctx, flow.ID, flow.ExecTime)

	ctx2 := createNextExecutionContext(ctx, flow.ExecTime)
	// test that flow was added to the queue
	queue := k.GetFlowsForBlock(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].ID)

	k.HandleFlow(ctx2, k.Logger(ctx2), flow, ctx.BlockTime(), nil)
	flow = k.GetFlowInfo(ctx2, flow.ID)
	ctx3 := createNextExecutionContext(ctx2, flow.ExecTime.Add(time.Hour))
	flow = k.GetFlowInfo(ctx3, flow.ID)
	//exec time is not updated
	require.True(t, flow.ExecTime.Before(ctx3.BlockTime()))
	//flow also not in queue
	iterator := k.FlowQueueIterator(ctx, ctx3.BlockTime())
	require.Nil(t, iterator.Key())

}

func TestErrorIsSavedToFlowInfo(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	configuration := types.ExecutionConfiguration{StopOnFailure: true}
	flow, emptyBalanceAcc := createTestFlow(ctx, configuration, keepers)

	err := flow.ValidateBasic()
	require.NoError(t, err)
	k := keepers.IntentKeeper

	k.SetFlowInfo(ctx, &flow)
	k.InsertFlowQueue(ctx, flow.ID, flow.ExecTime)

	ctx2 := createNextExecutionContext(ctx, flow.ExecTime)
	// test that flow was added to the queue
	queue := k.GetFlowsForBlock(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].ID)
	err = sendTokens(ctx, keepers, flow.Owner, emptyBalanceAcc, sdk.NewInt64Coin("stake", 3_000_000_000_000))
	require.NoError(t, err)
	err = sendTokens(ctx, keepers, flow.FeeAddress, emptyBalanceAcc, sdk.NewInt64Coin("stake", 3_000_000_000_000))
	require.NoError(t, err)
	k.HandleFlow(ctx2, k.Logger(ctx2), flow, ctx.BlockTime(), nil)

	flow = k.GetFlowInfo(ctx2, flow.ID)
	ctx3 := createNextExecutionContext(ctx2, flow.ExecTime.Add(time.Hour))
	flow = k.GetFlowInfo(ctx3, flow.ID)
	flowHistory := k.MustGetFlowHistory(ctx3, queue[0].ID)

	require.True(t, flow.ExecTime.Before(ctx3.BlockTime()))
	require.NotNil(t, flowHistory[0].Errors)
	require.Contains(t, flowHistory[0].Errors[0], "balance too low")

}

func TestOwnerMustBeSignerForLocalFlow(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)

	flowOwnerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	feeAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	toSendAcc, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
	require.Equal(t, keepers.BankKeeper.GetAllBalances(ctx, flowOwnerAddr)[0].Amount, math.NewInt(3_000_000_000_000))
	localMsg := &banktypes.MsgSend{
		FromAddress: toSendAcc.String(),
		ToAddress:   flowOwnerAddr.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))),
	}
	anys, _ := types.PackTxMsgAnys([]sdk.Msg{localMsg})

	flow := types.FlowInfo{
		ID:         123,
		Owner:      flowOwnerAddr.String(),
		FeeAddress: feeAddr.String(),
		Msgs:       anys,
	}
	k := keepers.IntentKeeper

	executedLocally, _, err := k.TriggerFlow(ctx, &flow)
	require.Contains(t, err.Error(), "owner doesn't have permission to send this message: unauthorized")
	require.False(t, executedLocally)

}

func TestNotExecutingWithWrongSigner(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)

	flowOwnerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	feeAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	toSendAcc, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
	require.Equal(t, keepers.BankKeeper.GetAllBalances(ctx, flowOwnerAddr)[0].Amount, math.NewInt(3_000_000_000_000))
	localMsg := &banktypes.MsgSend{
		FromAddress: flowOwnerAddr.String(),
		ToAddress:   toSendAcc.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))),
	}
	anys, _ := types.PackTxMsgAnys([]sdk.Msg{localMsg})

	flow := types.FlowInfo{
		ID:         123,
		Owner:      feeAddr.String(),
		FeeAddress: flowOwnerAddr.String(),
		Msgs:       anys,
	}
	k := keepers.IntentKeeper

	flow.Configuration = &types.ExecutionConfiguration{}
	k.HandleFlow(ctx, k.Logger(ctx), flow, ctx.BlockHeader().Time, nil)
	history, err := k.GetFlowHistory(ctx, flow.ID)
	require.Nil(t, err)

	require.Contains(t, history[len(history)-1].Errors[0], "owner doesn't have permission to send this message: unauthorized")
}

func TestAllowedToExecuteWithNoStopOnFailure(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)

	flowOwnerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	feeAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	toSendAcc, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
	require.Equal(t, keepers.BankKeeper.GetAllBalances(ctx, flowOwnerAddr)[0].Amount, math.NewInt(3_000_000_000_000))
	localMsg := &banktypes.MsgSend{
		FromAddress: flowOwnerAddr.String(),
		ToAddress:   toSendAcc.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))),
	}
	anys, _ := types.PackTxMsgAnys([]sdk.Msg{localMsg})

	flow := types.FlowInfo{
		ID:         123,
		Owner:      flowOwnerAddr.String(),
		FeeAddress: feeAddr.String(),
		Msgs:       anys,
	}
	k := keepers.IntentKeeper
	k.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: nil})
	k.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: nil})

	flow.Configuration = &types.ExecutionConfiguration{}
	flow.Conditions = &types.ExecutionConditions{}
	flow.Conditions.Comparisons = []*types.Comparison{{ResponseIndex: 0, ResponseKey: "Amount.[0].Amount", ValueType: "sdk.Int", Operator: 0, Operand: "101"}}
	k.HandleFlow(ctx, k.Logger(ctx), flow, ctx.BlockHeader().Time, nil)
	history, err := k.GetFlowHistory(ctx, flow.ID)
	require.Nil(t, err)
	require.True(t, history[len(history)-1].Executed)
}

func TestNotAllowedToExecuteWithStopOnFailure(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)

	flowOwnerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	feeAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	toSendAcc, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
	require.Equal(t, keepers.BankKeeper.GetAllBalances(ctx, flowOwnerAddr)[0].Amount, math.NewInt(3_000_000_000_000))
	localMsg := &banktypes.MsgSend{
		FromAddress: flowOwnerAddr.String(),
		ToAddress:   toSendAcc.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))),
	}
	anys, _ := types.PackTxMsgAnys([]sdk.Msg{localMsg})

	flow := types.FlowInfo{
		ID:         123,
		Owner:      flowOwnerAddr.String(),
		FeeAddress: feeAddr.String(),
		Msgs:       anys,
	}
	k := keepers.IntentKeeper
	k.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: nil})
	k.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: nil})

	flow.Configuration = &types.ExecutionConfiguration{StopOnFailure: true}
	flow.Conditions = &types.ExecutionConditions{}
	flow.Conditions.Comparisons = []*types.Comparison{{ResponseIndex: 0, ResponseKey: "Amount.[0].Amount", ValueType: "sdk.Int", Operator: 0, Operand: "101"}}
	k.HandleFlow(ctx, k.Logger(ctx), flow, ctx.BlockHeader().Time, nil)
	history, err := k.GetFlowHistory(ctx, flow.ID)
	require.Nil(t, err)
	require.False(t, history[len(history)-1].Executed)
}

func TestAllowedToExecuteWithNoStopOnFailureAndUseAndForComparisons(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)

	flowOwnerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	feeAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	toSendAcc, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
	require.Equal(t, keepers.BankKeeper.GetAllBalances(ctx, flowOwnerAddr)[0].Amount, math.NewInt(3_000_000_000_000))
	localMsg := &banktypes.MsgSend{
		FromAddress: flowOwnerAddr.String(),
		ToAddress:   toSendAcc.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))),
	}
	anys, _ := types.PackTxMsgAnys([]sdk.Msg{localMsg})

	flow := types.FlowInfo{
		ID:         123,
		Owner:      flowOwnerAddr.String(),
		FeeAddress: feeAddr.String(),
		Msgs:       anys,
	}
	k := keepers.IntentKeeper
	k.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: nil})
	k.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: nil})

	flow.Configuration = &types.ExecutionConfiguration{}
	flow.Conditions = &types.ExecutionConditions{UseAndForComparisons: true}
	flow.Conditions.Comparisons = []*types.Comparison{{ResponseIndex: 0, ResponseKey: "Amount.[0].Amount", ValueType: "sdk.Int", Operator: 0, Operand: "101"}}
	k.HandleFlow(ctx, k.Logger(ctx), flow, ctx.BlockHeader().Time, nil)
	history, err := k.GetFlowHistory(ctx, flow.ID)
	require.Nil(t, err)
	require.False(t, history[len(history)-1].Executed)
}

func createTestContext(t *testing.T) (sdk.Context, keeper.TestKeepers, codec.Codec) {
	ctx, keepers, cdc := keeper.CreateTestInput(t, false)

	types.Denom = "stake"
	keepers.IntentKeeper.SetParams(ctx, types.Params{
		FlowFundsCommission: 2,
		BurnFeePerMsg:       1_000_000,                 // 1trst
		FlowFlexFeeMul:      3,                         //
		MaxFlowDuration:     time.Hour * 24 * 366 * 10, // a little over 10 years
		MinFlowDuration:     time.Second * 60,
		MinFlowInterval:     time.Second * 20,
		GasFeeCoins:         sdk.NewCoins(sdk.NewCoin(types.Denom, math.OneInt())),
	})
	return ctx, keepers, cdc
}

func createTestFlow(ctx sdk.Context, configuration types.ExecutionConfiguration, keepers keeper.TestKeepers) (types.FlowInfo, sdk.AccAddress) {
	flowOwnerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	fundedFeeAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	emptyBalanceAcc, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
	startTime := ctx.BlockHeader().Time
	execTime := ctx.BlockHeader().Time.Add(time.Hour)
	endTime := ctx.BlockHeader().Time.Add(time.Hour * 2)
	localMsg := &banktypes.MsgSend{
		FromAddress: flowOwnerAddr.String(),
		ToAddress:   emptyBalanceAcc.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))),
	}
	anys, _ := types.PackTxMsgAnys([]sdk.Msg{localMsg})

	flow := types.FlowInfo{
		ID:            123,
		Owner:         flowOwnerAddr.String(),
		FeeAddress:    fundedFeeAddr.String(),
		ExecTime:      execTime,
		EndTime:       endTime,
		Interval:      time.Hour,
		StartTime:     startTime,
		Msgs:          anys,
		Configuration: &configuration,
		ICAConfig:     &types.ICAConfig{},
		Conditions:    &types.ExecutionConditions{},
	}
	return flow, emptyBalanceAcc
}

func createNoFeeFlow(ctx sdk.Context, configuration types.ExecutionConfiguration, keepers keeper.TestKeepers) (types.FlowInfo, sdk.AccAddress) {
	flowOwnerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	fundedFeeAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
	emptyBalanceAcc, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
	startTime := ctx.BlockHeader().Time
	execTime := ctx.BlockHeader().Time.Add(time.Hour)
	endTime := ctx.BlockHeader().Time.Add(time.Hour * 2)
	localMsg := &banktypes.MsgSend{
		FromAddress: flowOwnerAddr.String(),
		ToAddress:   emptyBalanceAcc.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))),
	}
	anys, _ := types.PackTxMsgAnys([]sdk.Msg{localMsg})

	flow := types.FlowInfo{
		ID:            123,
		Owner:         flowOwnerAddr.String(),
		FeeAddress:    fundedFeeAddr.String(),
		ExecTime:      execTime,
		EndTime:       endTime,
		Interval:      time.Hour,
		StartTime:     startTime,
		Msgs:          anys,
		Configuration: &configuration,
		ICAConfig:     &types.ICAConfig{},
		Conditions:    &types.ExecutionConditions{},
	}
	return flow, emptyBalanceAcc
}

func sendTokens(ctx sdk.Context, keepers keeper.TestKeepers, from string, toAddr sdk.AccAddress, amount sdk.Coin) error {
	fromAddr, _ := sdk.AccAddressFromBech32(from)
	err := keepers.BankKeeper.SendCoins(ctx, fromAddr, toAddr, sdk.NewCoins(amount))

	return err
}

func createBadFlow(ctx sdk.Context, configuration types.ExecutionConfiguration, keepers keeper.TestKeepers) (types.FlowInfo, sdk.AccAddress) {
	flowOwnerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
	feeAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	toSendAcc, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
	startTime := ctx.BlockHeader().Time
	execTime := ctx.BlockHeader().Time.Add(time.Hour)
	endTime := ctx.BlockHeader().Time.Add(time.Hour * 2)
	localMsg := &banktypes.MsgSend{
		FromAddress: flowOwnerAddr.String(),
		ToAddress:   toSendAcc.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))),
	}
	anys, _ := types.PackTxMsgAnys([]sdk.Msg{localMsg})

	flow := types.FlowInfo{
		ID:            123,
		Owner:         flowOwnerAddr.String(),
		FeeAddress:    feeAddr.String(),
		ExecTime:      execTime,
		EndTime:       endTime,
		Interval:      time.Hour,
		StartTime:     startTime,
		Msgs:          anys,
		Configuration: &configuration,
	}
	return flow, toSendAcc
}

func createNextExecutionContext(ctx sdk.Context, nextExecTime time.Time) sdk.Context {
	return sdk.NewContext(ctx.MultiStore(), tmproto.Header{
		Height:          ctx.BlockHeight() + 1111,
		Time:            nextExecTime,
		ChainID:         ctx.ChainID(),
		ProposerAddress: ctx.BlockHeader().ProposerAddress,
	}, false, ctx.Logger())
}
