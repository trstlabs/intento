package keeper

import (
	"testing"

	math "cosmossdk.io/math"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
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
