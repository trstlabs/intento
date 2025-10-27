package intent

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
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

	k.SetFlow(ctx, &flow)
	k.InsertFlowQueue(ctx, flow.ID, flow.ExecTime)

	ctx2 := createNextExecutionContext(ctx, flow.ExecTime)
	// test that flow was added to the queue
	queue := k.GetFlowsForBlockAndPruneQueue(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].ID)

	k.HandleFlow(ctx2, k.Logger(ctx2), flow, ctx2.BlockTime())
	flow = k.GetFlow(ctx2, flow.ID)
	ctx3 := createNextExecutionContext(ctx2, flow.ExecTime)

	//queue in BeginBocker
	queue = k.GetFlowsForBlockAndPruneQueue(ctx3)
	flowHistory := k.MustGetFlowHistory(ctx3, queue[0].ID)
	// test that flow history was updated
	require.Equal(t, ctx3.BlockHeader().Time, queue[0].ExecTime)
	require.Equal(t, 1, len(flowHistory))
	require.Equal(t, ctx2.BlockHeader().Time, flowHistory[0].ScheduledExecTime)
	require.Equal(t, ctx2.BlockHeader().Time, flowHistory[0].ActualExecTime)
	require.NotNil(t, ctx3.BlockHeader().Time, flowHistory[0].MsgResponses[0].Value)

	require.Equal(t, keepers.BankKeeper.GetAllBalances(ctx3, sendToAddr)[0].Amount, math.NewInt(100))

}

func TestBeginBlockerTransfer(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	configuration := types.ExecutionConfiguration{SaveResponses: true}
	flow, _ := createInvalidTestFlowTransfer(ctx, configuration, keepers)
	err := flow.ValidateBasic()
	require.NoError(t, err)
	k := keepers.IntentKeeper

	k.SetFlow(ctx, &flow)
	k.InsertFlowQueue(ctx, flow.ID, flow.ExecTime)

	ctx2 := createNextExecutionContext(ctx, flow.ExecTime)
	// test that flow was added to the queue
	queue := k.GetFlowsForBlockAndPruneQueue(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].ID)

	k.HandleFlow(ctx2, k.Logger(ctx2), flow, ctx2.BlockTime())
	flow = k.GetFlow(ctx2, flow.ID)
	ctx3 := createNextExecutionContext(ctx2, flow.ExecTime)

	// queue in BeginBlocker window (at nextExec)
	queue = k.GetFlowsForBlockAndPruneQueue(ctx3)
	require.Equal(t, 1, len(queue))
	flowHistory := k.MustGetFlowHistory(ctx3, queue[0].ID)
	// test that flow history was updated
	require.Equal(t, ctx3.BlockHeader().Time, queue[0].ExecTime)
	require.Equal(t, 1, len(flowHistory))
	require.Equal(t, ctx2.BlockHeader().Time, flowHistory[0].ScheduledExecTime)
	require.Equal(t, ctx2.BlockHeader().Time, flowHistory[0].ActualExecTime)
	require.NotNil(t, flowHistory[0].Errors)
	require.Equal(t, flowHistory[0].Errors[0], "msg handling error: could not execute local flow: 10010: invalid coins")
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
		k.SetFlow(ctx, &flow)
		k.InsertFlowQueue(ctx, flow.ID, flow.ExecTime)
	}
	ctx2 := createNextExecutionContext(ctx, flow.ExecTime)
	queue := k.GetFlowsForBlockAndPruneQueue(ctx2)
	BeginBlocker(ctx2, k)
	queue2 := k.GetFlowsForBlockAndPruneQueue(ctx2)
	ctx3 := createNextExecutionContext(ctx2, ctx2.BlockTime().Add(time.Minute))
	queue3 := k.GetFlowsForBlockAndPruneQueue(ctx3)
	require.NotEqual(t, len(queue), len(queue2))
	require.Equal(t, len(queue2), len(queue3)) //should be the same in later block if not executed

	BeginBlocker(ctx3, k)
	queue4 := k.GetFlowsForBlockAndPruneQueue(ctx3)
	require.Equal(t, len(queue4), 0)
}

func TestBeginBlockerStopOnSuccess(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	configuration := types.ExecutionConfiguration{StopOnSuccess: true}
	flow, _ := createTestFlow(ctx, configuration, keepers)
	err := flow.ValidateBasic()
	require.NoError(t, err)
	k := keepers.IntentKeeper

	k.SetFlow(ctx, &flow)
	k.InsertFlowQueue(ctx, flow.ID, flow.ExecTime)

	ctx2 := createNextExecutionContext(ctx, flow.ExecTime)
	// test that flow was added to the queue
	queue := k.GetFlowsForBlockAndPruneQueue(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].ID)
	// BeginBlocker logic
	k.HandleFlow(ctx2, k.Logger(ctx2), flow, ctx.BlockTime())
	flow = k.GetFlow(ctx2, flow.ID)
	ctx3 := createNextExecutionContext(ctx2, flow.ExecTime.Add(time.Hour))
	flow = k.GetFlow(ctx3, flow.ID)
	require.True(t, flow.ExecTime.Before(ctx3.BlockTime()))

}

func TestBeginBlockerStopOnFailure(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	configuration := types.ExecutionConfiguration{StopOnFailure: true}
	flow, _ := createBadFlow(ctx, configuration, keepers)
	err := flow.ValidateBasic()
	require.NoError(t, err)
	k := keepers.IntentKeeper

	k.SetFlow(ctx, &flow)
	k.InsertFlowQueue(ctx, flow.ID, flow.ExecTime)

	ctx2 := createNextExecutionContext(ctx, flow.ExecTime)
	// test that flow was added to the queue
	queue := k.GetFlowsForBlockAndPruneQueue(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].ID)

	k.HandleFlow(ctx2, k.Logger(ctx2), flow, ctx.BlockTime())
	flow = k.GetFlow(ctx2, flow.ID)
	ctx3 := createNextExecutionContext(ctx2, flow.ExecTime.Add(time.Hour))
	flow = k.GetFlow(ctx3, flow.ID)
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

	k.SetFlow(ctx, &flow)
	k.InsertFlowQueue(ctx, flow.ID, flow.ExecTime)

	ctx2 := createNextExecutionContext(ctx, flow.ExecTime)
	// test that flow was added to the queue
	queue := k.GetFlowsForBlockAndPruneQueue(ctx2)
	require.Equal(t, 1, len(queue))
	require.Equal(t, uint64(123), queue[0].ID)
	k.HandleFlow(ctx2, k.Logger(ctx2), queue[0], ctx.BlockTime())
	flow = k.GetFlow(ctx2, flow.ID)
	ctx3 := createNextExecutionContext(ctx2, flow.ExecTime.Add(time.Hour))
	flow = k.GetFlow(ctx3, flow.ID)
	flowHistory := k.MustGetFlowHistory(ctx3, flow.ID)
	require.NotNil(t, flowHistory[0].Errors)
	//exec time is not updated
	require.True(t, flow.ExecTime.Before(ctx3.BlockTime()))
	//flow also not in queue
	iterator := k.FlowQueueIterator(ctx, ctx3.BlockTime())
	require.Nil(t, iterator.Key())

}

func TestErrorIsSavedToflow(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	configuration := types.ExecutionConfiguration{StopOnFailure: true}
	flow, emptyBalanceAcc := createTestFlow(ctx, configuration, keepers)

	err := flow.ValidateBasic()
	require.NoError(t, err)
	k := keepers.IntentKeeper

	k.SetFlow(ctx, &flow)
	k.InsertFlowQueue(ctx, flow.ID, flow.ExecTime)

	ctx2 := createNextExecutionContext(ctx, flow.ExecTime)
	// test that flow was added to the queue
	flows := k.GetFlowsForBlockAndPruneQueue(ctx2)
	require.Equal(t, 1, len(flows))
	require.Equal(t, uint64(123), flows[0].ID)
	err = sendTokens(ctx, keepers, flow.Owner, emptyBalanceAcc, sdk.NewInt64Coin("stake", 3_000_000_000_000))
	require.NoError(t, err)
	err = sendTokens(ctx, keepers, flow.FeeAddress, emptyBalanceAcc, sdk.NewInt64Coin("stake", 3_000_000_000_000))
	require.NoError(t, err)
	k.HandleFlow(ctx2, k.Logger(ctx2), flow, ctx.BlockTime())

	flow = k.GetFlow(ctx2, flow.ID)
	ctx3 := createNextExecutionContext(ctx2, flow.ExecTime.Add(time.Hour))
	flow = k.GetFlow(ctx3, flow.ID)
	flowHistory := k.MustGetFlowHistory(ctx3, flows[0].ID)

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

	flow := types.Flow{
		ID:         123,
		Owner:      flowOwnerAddr.String(),
		FeeAddress: feeAddr.String(),
		Msgs:       anys,
	}
	k := keepers.IntentKeeper

	executedLocally, _, err := k.TriggerFlow(ctx, &flow)
	require.Contains(t, err.Error(), "owner doesn't have permission to send this message: unauthorized")
	require.NotEqual(t, int64(-1), executedLocally)

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

	flow := types.Flow{
		ID:         123,
		Owner:      feeAddr.String(),
		FeeAddress: flowOwnerAddr.String(),
		Msgs:       anys,
	}
	k := keepers.IntentKeeper

	flow.Configuration = &types.ExecutionConfiguration{}
	k.HandleFlow(ctx, k.Logger(ctx), flow, ctx.BlockHeader().Time)
	history, err := k.GetFlowHistory(ctx, flow.ID)
	require.Nil(t, err)

	require.Contains(t, history[len(history)-1].Errors[0], "owner doesn't have permission to send this message: unauthorized")
}

func TestEvaluateComparison_AutocompoundErrFlowConditions_AllowsExecution(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	k := keepers.IntentKeeper

	// Create a simple flow
	configuration := types.ExecutionConfiguration{SaveResponses: true}
	flow, _ := createTestFlow(ctx, configuration, keepers)
	k.SetFlow(ctx, &flow)

	// Record a latest history entry with ErrFlowConditions for the same flow
	errEntry := &types.FlowHistoryEntry{Errors: []string{types.ErrFlowConditions}}
	k.SetFlowHistoryEntry(ctx, flow.ID, errEntry)

	// Set a dummy comparison (will be skipped by autocompound shortcut)
	flow.Conditions = &types.ExecutionConditions{}
	flow.Conditions.Comparisons = []*types.Comparison{{
		ResponseIndex: 0,
		ResponseKey:   "",
		ValueType:     "math.Int",
		Operator:      types.ComparisonOperator_EQUAL,
		Operand:       "0",
	}}

	k.HandleFlow(ctx, k.Logger(ctx), flow, ctx.BlockHeader().Time)

	history, err := k.GetFlowHistory(ctx, flow.ID)
	require.NoError(t, err)
	require.NotEmpty(t, history)
	// Should execute despite previous ErrFlowConditions
	require.True(t, history[len(history)-1].Executed)
}

func TestEvaluateComparison_SkipWhenNoResponsesAndNoICQ_AllowsExecution(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	k := keepers.IntentKeeper

	configuration := types.ExecutionConfiguration{}
	flow, _ := createTestFlow(ctx, configuration, keepers)
	k.SetFlow(ctx, &flow)

	// Add history entries with nil/empty MsgResponses
	k.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: nil})
	k.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: nil})

	// Comparison without ICQ should be skipped and not block execution
	flow.Conditions = &types.ExecutionConditions{}
	flow.Conditions.Comparisons = []*types.Comparison{{
		ResponseIndex: 0,
		ResponseKey:   "Amount.[0].Amount",
		ValueType:     "math.Int",
		Operator:      types.ComparisonOperator_EQUAL,
		Operand:       "101",
	}}

	k.HandleFlow(ctx, k.Logger(ctx), flow, ctx.BlockHeader().Time)

	history, err := k.GetFlowHistory(ctx, flow.ID)
	require.NoError(t, err)
	require.NotEmpty(t, history)
	// Should execute even though there were no msg responses and no ICQ
	require.True(t, history[len(history)-1].Executed)
}

func TestEvaluateComparison_UsesLatestMsgResponsesFromOtherFlow(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	k := keepers.IntentKeeper

	configuration := types.ExecutionConfiguration{SaveResponses: true}
	flow, _ := createTestFlow(ctx, configuration, keepers)
	k.SetFlow(ctx, &flow)

	// Prepare another flow ID with a history entry that has MsgResponses
	otherFlowID := uint64(1111)
	resp := distributiontypes.MsgWithdrawDelegatorRewardResponse{Amount: sdk.NewCoins(sdk.NewInt64Coin("stake", 101))}
	respAny, err := cdctypes.NewAnyWithValue(&resp)
	require.NoError(t, err)
	k.SetFlowHistoryEntry(ctx, otherFlowID, &types.FlowHistoryEntry{MsgResponses: []*cdctypes.Any{respAny}})

	// Comparison targets the other flow's responses: Amount contains 101stake
	flow.Conditions = &types.ExecutionConditions{}
	flow.Conditions.UseAndForComparisons = true
	flow.Conditions.Comparisons = []*types.Comparison{{
		FlowID:        otherFlowID,
		ResponseIndex: 0,
		ResponseKey:   "Amount",
		ValueType:     "sdk.Coins",
		Operator:      types.ComparisonOperator_CONTAINS,
		Operand:       "101stake",
	}}

	k.HandleFlow(ctx, k.Logger(ctx), flow, ctx.BlockHeader().Time)

	history, err := k.GetFlowHistory(ctx, flow.ID)
	require.NoError(t, err)
	require.NotEmpty(t, history)
	// Should execute because comparison evaluated using latest msg responses of other flow
	require.True(t, history[len(history)-1].Executed)
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

func createTestFlow(ctx sdk.Context, configuration types.ExecutionConfiguration, keepers keeper.TestKeepers) (types.Flow, sdk.AccAddress) {
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

	flow := types.Flow{
		ID:            123,
		Owner:         flowOwnerAddr.String(),
		FeeAddress:    fundedFeeAddr.String(),
		ExecTime:      execTime,
		EndTime:       endTime,
		Interval:      time.Hour,
		StartTime:     startTime,
		Msgs:          anys,
		Configuration: &configuration,
		SelfHostedICA: &types.ICAConfig{},
		Conditions:    &types.ExecutionConditions{},
	}
	return flow, emptyBalanceAcc
}

func createInvalidTestFlowTransfer(ctx sdk.Context, configuration types.ExecutionConfiguration, keepers keeper.TestKeepers) (types.Flow, sdk.AccAddress) {
	flowOwnerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	fundedFeeAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	emptyBalanceAcc, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))
	startTime := ctx.BlockHeader().Time
	execTime := ctx.BlockHeader().Time.Add(time.Hour)
	endTime := ctx.BlockHeader().Time.Add(time.Hour * 2)
	localMsg := &ibctransfertypes.MsgTransfer{
		SourcePort:    "transfer",
		SourceChannel: "channel-0",
		Token: sdk.Coin{
			Amount: math.NewInt(100),
			Denom:  "10",
		},
		Sender:           flowOwnerAddr.String(),
		Receiver:         "",
		TimeoutHeight:    clienttypes.Height{RevisionNumber: 0, RevisionHeight: 0},
		TimeoutTimestamp: 0,
		Memo:             "hello",
	}

	anys, _ := types.PackTxMsgAnys([]sdk.Msg{localMsg})

	flow := types.Flow{
		ID:            123,
		Owner:         flowOwnerAddr.String(),
		FeeAddress:    fundedFeeAddr.String(),
		ExecTime:      execTime,
		EndTime:       endTime,
		Interval:      time.Hour,
		StartTime:     startTime,
		Msgs:          anys,
		Configuration: &configuration,
		SelfHostedICA: &types.ICAConfig{},
		Conditions:    &types.ExecutionConditions{},
	}
	return flow, emptyBalanceAcc
}

func createNoFeeFlow(ctx sdk.Context, configuration types.ExecutionConfiguration, keepers keeper.TestKeepers) (types.Flow, sdk.AccAddress) {
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

	flow := types.Flow{
		ID:            123,
		Owner:         flowOwnerAddr.String(),
		FeeAddress:    fundedFeeAddr.String(),
		ExecTime:      execTime,
		EndTime:       endTime,
		Interval:      time.Hour,
		StartTime:     startTime,
		Msgs:          anys,
		Configuration: &configuration,
		SelfHostedICA: &types.ICAConfig{},
		Conditions:    &types.ExecutionConditions{},
	}
	return flow, emptyBalanceAcc
}

func sendTokens(ctx sdk.Context, keepers keeper.TestKeepers, from string, toAddr sdk.AccAddress, amount sdk.Coin) error {
	fromAddr, _ := sdk.AccAddressFromBech32(from)
	err := keepers.BankKeeper.SendCoins(ctx, fromAddr, toAddr, sdk.NewCoins(amount))

	return err
}

func createBadFlow(ctx sdk.Context, configuration types.ExecutionConfiguration, keepers keeper.TestKeepers) (types.Flow, sdk.AccAddress) {
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

	flow := types.Flow{
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

func TestHandleFlow_WithFalseFeedbackLoopIndex(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	k := keepers.IntentKeeper

	// Create valid funded accounts for all addresses
	ownerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 10_000_000)))
	addr1, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 10_000_000)))
	addr2, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 10_000_000)))
	addr3, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 10_000_000)))

	// Create a simple flow with 2 messages using valid addresses
	msg1 := &banktypes.MsgSend{
		FromAddress: addr1.String(),
		ToAddress:   addr2.String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("stake", 10)),
	}
	anyMsg1, err := types.PackTxMsgAnys([]sdk.Msg{msg1})
	require.NoError(t, err)
	msg2 := &banktypes.MsgSend{
		FromAddress: addr2.String(),
		ToAddress:   addr3.String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("stake", 20)),
	}
	anyMsg2, err := types.PackTxMsgAnys([]sdk.Msg{msg2})
	require.NoError(t, err)

	flow := types.Flow{
		ID:         999,
		Owner:      ownerAddr.String(),
		FeeAddress: ownerAddr.String(),
		ExecTime:   ctx.BlockHeader().Time,
		Interval:   time.Hour,
		EndTime:    ctx.BlockHeader().Time.Add(time.Hour * 2),
		Conditions: &types.ExecutionConditions{
			FeedbackLoops: []*types.FeedbackLoop{
				{
					// Intentionally invalid index (out of bounds)
					MsgsIndex: 5,
					FlowID:    0, //999
				},
			},
		},
		Configuration: &types.ExecutionConfiguration{
			SaveResponses: true,
		},
	}
	flow.Msgs = append(anyMsg1, anyMsg2...) // Add both messages to the flow

	require.NoError(t, flow.ValidateBasic())

	k.SetFlow(ctx, &flow)
	k.InsertFlowQueue(ctx, flow.ID, flow.ExecTime)

	ctx2 := createNextExecutionContext(ctx, flow.ExecTime)

	// Run HandleFlow, which internally calls handleRunFeedbackLoops
	// We expect this to handle the invalid MsgsIndex and not panic
	// but possibly to log or return an error related to slicing.
	// Let's capture any panic for test safety.
	var panicked bool
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		k.HandleFlow(ctx2, k.Logger(ctx2), flow, ctx2.BlockTime())
	}()

	require.False(t, panicked, "Expected no panic despite invalid MsgsIndex")

	// Now check flow history for errors related to feedback loop execution
	flowHistory := k.MustGetFlowHistory(ctx2, flow.ID)
	require.NotEmpty(t, flowHistory)

	for _, h := range flowHistory {
		for _, errStr := range h.Errors {
			if errStr != "" {
				require.Contains(t, errStr, "invalid feedback loop MsgsIndex: 5 (only 2 messages in flow 999)")

				t.Logf("Found flow history error: %s", errStr)
			}
		}
	}

	require.Equal(t, len(flowHistory), 1, "Expected one flow history entry for flow 999")
	require.Equal(t, len(flowHistory[0].MsgResponses), 0, "Expected zero message responses in flow history due to invalid feedback loop")

	ctx3 := createNextExecutionContext(ctx, flow.ExecTime.Add(time.Hour))

	// Also check flow queue state after execution
	queue := k.GetFlowsForBlockAndPruneQueue(ctx3)
	require.NotNil(t, queue)
	t.Logf("Flows queued after HandleFlow: %d", len(queue))
}

func TestHandleFlow_WithGoodFeedbackLoopIndex(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	k := keepers.IntentKeeper

	// Create valid funded accounts for all addresses
	ownerAddr, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 10_000_000)))
	addr2, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 10_000_000)))

	// Create a simple flow with 2 messages using valid addresses
	msg1 := &banktypes.MsgSend{
		FromAddress: ownerAddr.String(),
		ToAddress:   addr2.String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("stake", 10)),
	}
	anyMsg1, err := types.PackTxMsgAnys([]sdk.Msg{msg1})
	require.NoError(t, err)
	msg2 := &banktypes.MsgSend{
		FromAddress: ownerAddr.String(),
		ToAddress:   addr2.String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("stake", 20)),
	}
	anyMsg2, err := types.PackTxMsgAnys([]sdk.Msg{msg2})
	require.NoError(t, err)

	flow := types.Flow{
		ID:         999,
		Owner:      ownerAddr.String(),
		FeeAddress: ownerAddr.String(),
		ExecTime:   ctx.BlockHeader().Time,
		Interval:   time.Hour,
		EndTime:    ctx.BlockHeader().Time.Add(time.Hour * 2),
		Conditions: &types.ExecutionConditions{
			FeedbackLoops: []*types.FeedbackLoop{
				{
					MsgsIndex:     1,
					FlowID:        11,
					MsgKey:        "Amount",
					ValueType:     "sdk.Coin",
					ResponseIndex: 0, // Use the first response from the flow history entry
					ResponseKey:   "Amount.[0].Amount",
				},
			},
		},
		Configuration: &types.ExecutionConfiguration{
			SaveResponses: true,
		},
	}
	flow.Msgs = append(anyMsg1, anyMsg2...) // Add both messages to the flow

	require.NoError(t, flow.ValidateBasic())

	k.SetFlow(ctx, &flow)
	k.InsertFlowQueue(ctx, flow.ID, flow.ExecTime)

	resp := distributiontypes.MsgWithdrawDelegatorRewardResponse{Amount: sdk.NewCoins(sdk.NewInt64Coin("stake", 100))}
	respAny, err := types.PackTxMsgAnys([]sdk.Msg{&resp})
	require.NoError(t, err)
	k.SetFlowHistoryEntry(ctx, 11, &types.FlowHistoryEntry{MsgResponses: respAny})
	flow11 := flow
	flow11.ID = 11
	k.SetFlow(ctx, &flow11)

	ctx2 := createNextExecutionContext(ctx, flow.ExecTime)

	// Run HandleFlow, which internally calls handleRunFeedbackLoops
	// We expect this to handle the invalid MsgsIndex and not panic
	// but possibly to log or return an error related to slicing.
	// Let's capture any panic for test safety.
	var panicked bool
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		k.HandleFlow(ctx2, k.Logger(ctx2), flow, ctx2.BlockTime())
	}()

	require.False(t, panicked, "Expected no panic despite invalid MsgsIndex")

	// Now check flow history for errors related to feedback loop execution
	flowHistory := k.MustGetFlowHistory(ctx2, flow.ID)
	require.NotEmpty(t, flowHistory)

	for _, h := range flowHistory {
		for _, errStr := range h.Errors {
			if errStr != "" {
				require.Contains(t, errStr, "invalid feedback loop MsgsIndex: 5 (only 2 messages in flow 999)")

				t.Logf("Found flow history error: %s", errStr)
			}
		}
	}

	require.Equal(t, len(flowHistory), 1, "Expected one flow history entry for flow 999")
	require.Equal(t, len(flowHistory[0].Errors), 0, "Expected no errors in flow history")
	require.Equal(t, len(flowHistory[0].MsgResponses), 2, "Expected 2 message responses")

	flowUpdated := k.GetFlow(ctx2, flow.ID)
	require.NotEqual(t, flow.Msgs[1], flowUpdated.Msgs[1], "Expected flow messages to be updated after HandleFlow execution")

	ctx3 := createNextExecutionContext(ctx, flow.ExecTime.Add(time.Hour))

	// Also check flow queue state after execution
	queue := k.GetFlowsForBlockAndPruneQueue(ctx3)
	require.NotNil(t, queue)
	t.Logf("Flows queued after HandleFlow: %d", len(queue))
}

func TestHandleLocalFlow_SchedulesPreUpgrade(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	k := keepers.IntentKeeper

	configuration := types.ExecutionConfiguration{SaveResponses: true}
	flow, _ := createTestFlow(ctx, configuration, keepers)

	require.NoError(t, flow.ValidateBasic())

	k.SetFlow(ctx, &flow)
	k.InsertFlowQueue(ctx, flow.ID, flow.ExecTime)

	prevExec := flow.ExecTime
	nextExec := prevExec.Add(flow.Interval)

	flowHistory := types.FlowHistoryEntry{
		ScheduledExecTime: prevExec,
		ActualExecTime:    prevExec,
		MsgResponses:      []*cdctypes.Any{{TypeUrl: "", Value: nil}},
		Errors:            nil,
	}
	k.SetFlowHistoryEntry(ctx, flow.ID, &flowHistory)
	// Build context at prevExec with intento-1 chain to trigger scheduling logic
	header2 := tmproto.Header{
		Height:          3_849_000, // below upgrade
		Time:            flow.ExecTime,
		ChainID:         "intento-1",
		ProposerAddress: ctx.BlockHeader().ProposerAddress,
	}
	ctx2 := sdk.NewContext(ctx.MultiStore(), header2, false, ctx.Logger())
	queue := k.GetFlowsForBlockAndPruneQueue(ctx2)
	require.Equal(t, 1, len(queue))

	k.HandleFlow(ctx2, k.Logger(ctx2), flow, ctx2.BlockTime())
	flow = k.GetFlow(ctx2, flow.ID)

	ctx3 := createNextExecutionContext(ctx2, nextExec)
	queue2 := k.GetFlowsForBlockAndPruneQueue(ctx3)
	require.Equal(t, 1, len(queue2))
}

func TestHandleLocalFlow_SchedulesPostUpgrade(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	k := keepers.IntentKeeper

	configuration := types.ExecutionConfiguration{SaveResponses: true}
	flow, _ := createTestFlow(ctx, configuration, keepers)
	require.NoError(t, flow.ValidateBasic())

	k.SetFlow(ctx, &flow)
	k.InsertFlowQueue(ctx, flow.ID, flow.ExecTime)

	prevExec := flow.ExecTime
	nextExec := prevExec.Add(flow.Interval)

	flowHistory := types.FlowHistoryEntry{
		ScheduledExecTime: prevExec,
		ActualExecTime:    prevExec,
		MsgResponses:      []*cdctypes.Any{{TypeUrl: "", Value: nil}},
		Errors:            nil,
	}
	k.SetFlowHistoryEntry(ctx, flow.ID, &flowHistory)
	header := tmproto.Header{
		Height:          3_900_000,
		Time:            prevExec,
		ChainID:         "intento-1",
		ProposerAddress: ctx.BlockHeader().ProposerAddress,
	}
	ctx2 := sdk.NewContext(ctx.MultiStore(), header, false, ctx.Logger())

	queue := k.GetFlowsForBlockAndPruneQueue(ctx2)
	require.Equal(t, 1, len(queue))

	k.HandleFlow(ctx2, k.Logger(ctx2), flow, ctx2.BlockTime())
	flow = k.GetFlow(ctx2, flow.ID)

	headerNext := tmproto.Header{
		Height:          header.Height + 1,
		Time:            nextExec,
		ChainID:         "intento-1",
		ProposerAddress: ctx.BlockHeader().ProposerAddress,
	}
	ctx3 := sdk.NewContext(ctx.MultiStore(), headerNext, false, ctx.Logger())
	queue2 := k.GetFlowsForBlockAndPruneQueue(ctx3)
	require.Equal(t, 1, len(queue2))
	require.Equal(t, flow.ID, queue2[0].ID)
}
