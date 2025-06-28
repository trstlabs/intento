package keeper

import (
	"testing"

	math "cosmossdk.io/math"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/trstlabs/intento/x/intent/types"
)

func TestParseCoin(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flowInfo := createBaseFlowInfo(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flowInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flowInfo.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	flowInfo.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(1000)))
	flowInfo.Conditions.FeedbackLoops = []*types.FeedbackLoop{{ResponseIndex: 0, ResponseKey: "Amount", MsgsIndex: 0, MsgKey: "Amount", ValueType: "sdk.Coin"}}
	err = keeper.RunFeedbackLoops(ctx, flowInfo.ID, &flowInfo.Msgs, flowInfo.Conditions)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(flowInfo.Msgs[0], &msgDelegate)
	require.NoError(t, err)
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(101)))

	executedLocally, _, err = keeper.TriggerFlow(ctx, &flowInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
}

func TestParseCoinFromMsgExec(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))

	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flowInfo := createBaseFlowInfo(delAddr, flowAddr)

	// Wrap MsgWithdrawDelegatorReward in MsgExec
	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	anyReward, err := cdctypes.NewAnyWithValue(msgWithdrawDelegatorReward)
	require.NoError(t, err)

	msgExec := &authztypes.MsgExec{
		Grantee: delAddr.String(),
		Msgs:    []*cdctypes.Any{anyReward},
	}
	flowInfo.Msgs, err = types.PackTxMsgAnys([]sdk.Msg{msgExec})
	require.NoError(t, err)
	err = flowInfo.ValidateBasic()
	require.NoError(t, err)

	// executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flowInfo)
	// require.NoError(t, err)
	msgWithdrawDelegatorRewardResp := distrtypes.MsgWithdrawDelegatorRewardResponse{Amount: sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1000)))}
	msgWithdrawDelegatorRewardRespAny, err := cdctypes.NewAnyWithValue(&msgWithdrawDelegatorRewardResp)
	require.NoError(t, err)

	msgExecResp := authztypes.MsgExecResponse{Results: [][]byte{msgWithdrawDelegatorRewardRespAny.Value}}
	msgExecRespAny, err := cdctypes.NewAnyWithValue(&msgExecResp)
	require.NoError(t, err)

	msgResponses, _, err := keeper.HandleDeepResponses(ctx, []*cdctypes.Any{msgExecRespAny}, sdk.AccAddress{}, flowInfo, 0)
	require.NoError(t, err)
	// require.True(t, executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flowInfo.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	flowInfo.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(1000)))
	flowInfo.Conditions.FeedbackLoops = []*types.FeedbackLoop{{
		ResponseIndex: 0,
		ResponseKey:   "Amount",
		MsgsIndex:     0,
		MsgKey:        "Amount",
		ValueType:     "sdk.Coin",
	}}

	err = keeper.RunFeedbackLoops(ctx, flowInfo.ID, &flowInfo.Msgs, flowInfo.Conditions)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(flowInfo.Msgs[0], &msgDelegate)
	require.NoError(t, err)
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(1000)))
}

func TestParseInnerString(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flowInfo := createBaseFlowInfo(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flowInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flowInfo.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	msgDelegate.Amount.Denom = "test"
	flowInfo.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount.Denom, "test")
	flowInfo.Conditions.FeedbackLoops = []*types.FeedbackLoop{{ResponseIndex: 0, ResponseKey: "Amount.[0].Denom", MsgsIndex: 0, MsgKey: "Amount.Denom", ValueType: "string"}}
	err = keeper.RunFeedbackLoops(ctx, flowInfo.ID, &flowInfo.Msgs, flowInfo.Conditions)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(flowInfo.Msgs[0], &msgDelegate)
	require.NoError(t, err)
	require.Equal(t, msgDelegate.Amount.Denom, "stake")

	executedLocally, _, err = keeper.TriggerFlow(ctx, &flowInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
}

func TestParseInnerStringFail(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flowInfo := createBaseFlowInfo(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flowInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flowInfo.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	msgDelegate.Amount.Denom = "test"
	flowInfo.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount.Denom, "test")
	flowInfo.Conditions.FeedbackLoops = []*types.FeedbackLoop{{ResponseIndex: 0, ResponseKey: "Amount.[0].Amount", MsgsIndex: 0, MsgKey: "Amount.Denom", ValueType: "string"}}
	err = keeper.RunFeedbackLoops(ctx, flowInfo.ID, &flowInfo.Msgs, flowInfo.Conditions)
	require.Error(t, err)

}

func TestParseInnerInt(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flowInfo := createBaseFlowInfo(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flowInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flowInfo.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	flowInfo.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(1000)))
	flowInfo.Conditions.FeedbackLoops = []*types.FeedbackLoop{{ResponseIndex: 0, ResponseKey: "Amount.[0].Amount", MsgsIndex: 0, MsgKey: "Amount.Amount", ValueType: "sdk.Int"}}
	err = keeper.RunFeedbackLoops(ctx, flowInfo.ID, &flowInfo.Msgs, flowInfo.Conditions)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(flowInfo.Msgs[0], &msgDelegate)
	require.NoError(t, err)
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(101)))

	executedLocally, _, err = keeper.TriggerFlow(ctx, &flowInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
}

func TestCompareInnerIntTrue(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flowInfo := createBaseFlowInfo(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flowInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flowInfo.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	flowInfo.Conditions = &types.ExecutionConditions{}
	flowInfo.Conditions.Comparisons = []*types.Comparison{{ResponseIndex: 0, ResponseKey: "Amount.[0].Amount", ValueType: "sdk.Int", Operator: 0, Operand: "101"}}
	boolean, err := keeper.CompareResponseValue(ctx, flowInfo.ID, msgResponses, *flowInfo.Conditions.Comparisons[0])
	require.NoError(t, err)

	require.True(t, boolean)
}

func TestCompareCoinTrue(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flowInfo := createBaseFlowInfo(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flowInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flowInfo.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})
	flowInfo.Conditions = &types.ExecutionConditions{}
	flowInfo.Conditions.Comparisons = []*types.Comparison{{ResponseIndex: 0, ResponseKey: "Amount.[0]", ValueType: "sdk.Coin", Operator: 0, Operand: "101stake"}}
	boolean, err := keeper.CompareResponseValue(ctx, flowInfo.ID, msgResponses, *flowInfo.Conditions.Comparisons[0])
	require.NoError(t, err)

	require.True(t, boolean)
}

func TestCompareCoinLargerThanTrue(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flowInfo := createBaseFlowInfo(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flowInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flowInfo.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	flowInfo.Conditions = &types.ExecutionConditions{}
	comparison := types.Comparison{ResponseIndex: 0, ResponseKey: "Amount.[0]", ValueType: "sdk.Coin", Operator: 4, Operand: "11stake"}
	flowInfo.Conditions.Comparisons = []*types.Comparison{&comparison}

	boolean, err := keeper.CompareResponseValue(ctx, flowInfo.ID, msgResponses, *flowInfo.Conditions.Comparisons[0])
	require.NoError(t, err)
	require.True(t, boolean)

	boolean, err = keeper.allowedToExecute(ctx, flowInfo)
	require.NoError(t, err)
	require.True(t, boolean)
}

func TestCompareIntFalse(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flowInfo := createBaseFlowInfo(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flowInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flowInfo.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	flowInfo.Conditions = &types.ExecutionConditions{}
	flowInfo.Conditions.Comparisons = []*types.Comparison{{ResponseIndex: 0, ResponseKey: "Amount.[0].Amount", ValueType: "sdk.Int", Operator: 0, Operand: "100000000000"}}
	boolean, err := keeper.CompareResponseValue(ctx, flowInfo.ID, msgResponses, *flowInfo.Conditions.Comparisons[0])
	require.NoError(t, err)

	require.False(t, boolean)
}

func TestCompareDenomString(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flowInfo := createBaseFlowInfo(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flowInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flowInfo.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	flowInfo.Conditions = &types.ExecutionConditions{}
	flowInfo.Conditions.Comparisons = []*types.Comparison{{ResponseIndex: 0, ResponseKey: "Amount.[0].Denom", ValueType: "string", Operator: 0, Operand: "stake"}}
	boolean, err := keeper.CompareResponseValue(ctx, flowInfo.ID, msgResponses, *flowInfo.Conditions.Comparisons[0])
	require.NoError(t, err)

	require.True(t, boolean)
}

func TestInvalidResponseIndex(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flowInfo := createBaseFlowInfo(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flowInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flowInfo.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	flowInfo.Conditions = &types.ExecutionConditions{}
	flowInfo.Conditions.Comparisons = []*types.Comparison{{ResponseIndex: 1, ResponseKey: "Amount.[0].Denom", ValueType: "string", Operator: 0, Operand: "sta"}}
	_, err = keeper.CompareResponseValue(ctx, flowInfo.ID, msgResponses, *flowInfo.Conditions.Comparisons[0])
	require.Error(t, err)

	require.Contains(t, err.Error(), "number of responses")
}

func TestCompareDenomStringContains(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flowInfo := createBaseFlowInfo(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flowInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flowInfo.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	flowInfo.Conditions = &types.ExecutionConditions{}
	flowInfo.Conditions.Comparisons = []*types.Comparison{{ResponseIndex: 0, ResponseKey: "Amount.[0].Denom", ValueType: "string", Operator: 1, Operand: "sta"}}
	boolean, err := keeper.CompareResponseValue(ctx, flowInfo.ID, msgResponses, *flowInfo.Conditions.Comparisons[0])
	require.NoError(t, err)

	require.True(t, boolean)
}

func TestCompareArrayCoinsContainsTrue(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flowInfo := createBaseFlowInfo(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flowInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flowInfo.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	flowInfo.Conditions = &types.ExecutionConditions{}
	flowInfo.Conditions.Comparisons = []*types.Comparison{{ResponseIndex: 0, ResponseKey: "Amount", ValueType: "sdk.Coins", Operator: 1, Operand: "101stake"}}
	boolean, err := keeper.CompareResponseValue(ctx, flowInfo.ID, msgResponses, *flowInfo.Conditions.Comparisons[0])
	require.NoError(t, err)

	require.True(t, boolean)
}

func TestCompareArrayCoinsContainsFalse(t *testing.T) {
	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	fakeCoins, _ := sdk.ParseCoinsNormalized("1000abc,3000000000degf")
	msgResponse := distrtypes.MsgWithdrawDelegatorRewardResponse{Amount: fakeCoins}
	any, _ := cdctypes.NewAnyWithValue(&msgResponse)
	msgResponses := []*cdctypes.Any{any}
	responseComparison := []*types.Comparison{{ResponseIndex: 0, ResponseKey: "Amount", ValueType: "sdk.Coins", Operator: 0, Operand: "100aaa"}}
	boolean, err := keeper.CompareResponseValue(ctx, 1, msgResponses, *responseComparison[0])
	require.NoError(t, err)

	require.False(t, boolean)
}

func TestCompareArrayEquals(t *testing.T) {
	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	fakeCoins, _ := sdk.ParseCoinsNormalized("1000abc,3000000000degf")
	msgResponse := distrtypes.MsgWithdrawDelegatorRewardResponse{Amount: fakeCoins}
	any, _ := cdctypes.NewAnyWithValue(&msgResponse)
	msgResponses := []*cdctypes.Any{any}
	responseComparison := []*types.Comparison{{ResponseIndex: 0, ResponseKey: "Amount", ValueType: "sdk.Coins", Operator: 0, Operand: "1000abc,3000000000degf"}}
	boolean, err := keeper.CompareResponseValue(ctx, 1, msgResponses, *responseComparison[0])
	require.NoError(t, err)

	require.True(t, boolean)
}

func TestParseAmountICQ(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flowInfo := createBaseFlowInfo(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	flowInfo.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(1000)))
	flowInfo.Conditions.FeedbackLoops = []*types.FeedbackLoop{{ResponseIndex: 0, ResponseKey: "", MsgsIndex: 0, MsgKey: "Amount.Amount", ValueType: "sdk.Int"}}
	queryCallback, err := math.NewInt(39999999999).Marshal()
	require.NoError(t, err)
	flowInfo.Conditions.FeedbackLoops[0].ICQConfig = &types.ICQConfig{Response: queryCallback}
	err = keeper.RunFeedbackLoops(ctx, flowInfo.ID, &flowInfo.Msgs, flowInfo.Conditions)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(flowInfo.Msgs[0], &msgDelegate)
	require.NoError(t, err)
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(39999999999)))
}

func TestFeedbackLoopNoDuplicates(t *testing.T) {
	// Setup test environment
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)

	// Create a flow with three messages and two feedback loops
	flowInfo := createBaseFlowInfo(delAddr, flowAddr)

	// Create three test messages
	msg1 := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	msg2 := newFakeMsgDelegate(delAddr, val)
	msg3 := newFakeMsgDelegate(delAddr, val)

	// Set initial amounts
	msg2.Amount = sdk.NewCoin("stake", math.NewInt(1000))
	msg3.Amount = sdk.NewCoin("stake", math.NewInt(2000))

	// Pack all messages
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msg1, msg2, msg3})

	// Set up feedback loops that will cause potential duplicates
	flowInfo.Conditions = &types.ExecutionConditions{
		FeedbackLoops: []*types.FeedbackLoop{
			{
				// First feedback loop after first message
				ResponseIndex: 0,
				ResponseKey:   "Amount.[0].Amount",
				MsgsIndex:     1,
				MsgKey:        "Amount.Amount",
				ValueType:     "sdk.Int",
			},
			{
				// Second feedback loop after first message (same point as first)
				ResponseIndex: 0,
				ResponseKey:   "Amount.[0].Amount",
				MsgsIndex:     2,
				MsgKey:        "Amount.Amount",
				ValueType:     "sdk.Int",
			},
		},
	}

	// Execute the flow
	executedLocally, responses, err := keeper.TriggerFlow(ctx, &flowInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)

	// Verify we have the expected number of responses (3 messages)
	require.Len(t, responses, 3, "expected 3 message responses")

	// Create a flow history entry with the actual responses
	historyEntry := &types.FlowHistoryEntry{
		ScheduledExecTime: flowInfo.ExecTime,
		ActualExecTime:    flowInfo.ExecTime,
		ExecFee:           sdk.NewCoin("stake", math.NewInt(0)),
		Executed:          true,
		TimedOut:          false,
		MsgResponses:      responses, // Use the actual responses from TriggerFlow
		Errors:            []string{},
	}
	keeper.SetFlowHistoryEntry(ctx, flowInfo.ID, historyEntry)
	history, _ := keeper.GetFlowHistory(ctx, flowInfo.ID)

	// Manually run the feedback loops to ensure they're processed
	err = keeper.RunFeedbackLoops(ctx, flowInfo.ID, &flowInfo.Msgs, flowInfo.Conditions)
	require.NoError(t, err)

	// Verify we have exactly one history entry
	require.Len(t, history, 1, "expected exactly one history entry")

	// Verify the history entry has the correct number of message responses
	require.Len(t, history[0].MsgResponses, 3, "expected exactly 3 message responses in history")
	require.Empty(t, history[0].Errors, "expected no errors in history entry")

	// Verify the messages were modified by the feedback loops
	// Unpack the messages from the flow to check if they were modified

	var modifiedMsg2, modifiedMsg3 *stakingtypes.MsgDelegate
	err = keeper.cdc.UnpackAny(flowInfo.Msgs[1], &modifiedMsg2)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(flowInfo.Msgs[2], &modifiedMsg3)
	require.NoError(t, err)

	// Verify the amounts were modified (not equal to original)
	require.NotEqual(t, msg2.Amount, modifiedMsg2.Amount, "second message should be modified by feedback loop")
	require.NotEqual(t, msg3.Amount, modifiedMsg3.Amount, "third message should be modified by feedback loop")
}

func TestParseCoinICQ(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flowInfo := createBaseFlowInfo(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	flowInfo.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(1000)))
	flowInfo.Conditions.FeedbackLoops = []*types.FeedbackLoop{{ResponseIndex: 0, ResponseKey: "", MsgsIndex: 0, MsgKey: "Amount", ValueType: "sdk.Coin"}}
	coin := sdk.NewCoin("stake", math.NewInt(39999999999))
	queryCallback, err := coin.Marshal()
	require.NoError(t, err)
	flowInfo.Conditions.FeedbackLoops[0].ICQConfig = &types.ICQConfig{Response: queryCallback}
	err = keeper.RunFeedbackLoops(ctx, flowInfo.ID, &flowInfo.Msgs, flowInfo.Conditions)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(flowInfo.Msgs[0], &msgDelegate)
	require.NoError(t, err)
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(39999999999)))
}

func TestParseStringICQ(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flowInfo := createBaseFlowInfo(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	flowInfo.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(1000)))
	flowInfo.Conditions.FeedbackLoops = []*types.FeedbackLoop{{ResponseIndex: 0, ResponseKey: "", MsgsIndex: 0, MsgKey: "Amount.Denom", ValueType: "string", ICQConfig: &types.ICQConfig{ConnectionId: "connection-0", ChainId: "test"}}}
	queryCallback := []byte("fuzz")
	flowInfo.Conditions.FeedbackLoops[0].ICQConfig = &types.ICQConfig{Response: queryCallback}
	err := keeper.RunFeedbackLoops(ctx, flowInfo.ID, &flowInfo.Msgs, flowInfo.Conditions)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(flowInfo.Msgs[0], &msgDelegate)
	require.NoError(t, err)
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("fuzz", math.NewInt(1000)))
}

func TestCompareCoinTrueICQ(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flowInfo := createBaseFlowInfo(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flowInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flowInfo.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	flowInfo.Conditions = &types.ExecutionConditions{}
	flowInfo.Conditions.Comparisons = []*types.Comparison{{ResponseIndex: 0, ResponseKey: "", ValueType: "sdk.Coin", Operator: 4, Operand: "101stake"}}

	coin := sdk.NewCoin("stake", math.NewInt(39999999999))
	queryCallback, err := coin.Marshal()
	require.NoError(t, err)
	flowInfo.Conditions.Comparisons[0].ICQConfig = &types.ICQConfig{Response: queryCallback}
	boolean, err := keeper.CompareResponseValue(ctx, flowInfo.ID, msgResponses, *flowInfo.Conditions.Comparisons[0])
	require.NoError(t, err)

	require.True(t, boolean)
}
