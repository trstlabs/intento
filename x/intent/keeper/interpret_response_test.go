package keeper

import (
	"encoding/base64"
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
	flow := createBaseflow(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flow)
	require.NoError(t, err)
	require.Equal(t, int64(-1), executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	flow.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(1000)))
	flow.Conditions.FeedbackLoops = []*types.FeedbackLoop{{ResponseIndex: 0, ResponseKey: "Amount", MsgsIndex: 0, MsgKey: "Amount", ValueType: "sdk.Coin"}}
	err = keeper.RunFeedbackLoops(ctx, flow.ID, &flow.Msgs, flow.Conditions)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(flow.Msgs[0], &msgDelegate)
	require.NoError(t, err)
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(101)))

	executedLocally, _, err = keeper.TriggerFlow(ctx, &flow)
	require.NoError(t, err)
	require.Equal(t, int64(-1), executedLocally)
}

func TestParseCoinFromMsgExec(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))

	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flow := createBaseflow(delAddr, flowAddr)
	// Wrap MsgWithdrawDelegatorReward in MsgExec
	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	anyReward, err := cdctypes.NewAnyWithValue(msgWithdrawDelegatorReward)
	require.NoError(t, err)

	msgExec := &authztypes.MsgExec{
		Grantee: delAddr.String(),
		Msgs:    []*cdctypes.Any{anyReward},
	}
	flow.Msgs, err = types.PackTxMsgAnys([]sdk.Msg{msgExec})
	require.NoError(t, err)
	err = flow.ValidateBasic()
	require.NoError(t, err)
	msgWithdrawDelegatorRewardResp := distrtypes.MsgWithdrawDelegatorRewardResponse{Amount: sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1000)))}
	msgWithdrawDelegatorRewardRespAny, err := cdctypes.NewAnyWithValue(&msgWithdrawDelegatorRewardResp)
	require.NoError(t, err)

	msgExecResp := authztypes.MsgExecResponse{Results: [][]byte{msgWithdrawDelegatorRewardRespAny.Value}}
	msgExecRespAny, err := cdctypes.NewAnyWithValue(&msgExecResp)
	require.NoError(t, err)

	msgResponses, _, err := keeper.HandleDeepResponses(ctx, []*cdctypes.Any{msgExecRespAny}, sdk.AccAddress{}, flow, 0)
	require.NoError(t, err)
	// require.Equal(t, int64(-1), executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	flow.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(1000)))
	flow.Conditions.FeedbackLoops = []*types.FeedbackLoop{{
		ResponseIndex: 0,
		ResponseKey:   "Amount",
		MsgsIndex:     0,
		MsgKey:        "Amount",
		ValueType:     "sdk.Coin",
	}}

	err = keeper.RunFeedbackLoops(ctx, flow.ID, &flow.Msgs, flow.Conditions)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(flow.Msgs[0], &msgDelegate)
	require.NoError(t, err)
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(1000)))
}

func TestParseInnerString(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flow := createBaseflow(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flow)
	require.NoError(t, err)
	require.Equal(t, int64(-1), executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	msgDelegate.Amount.Denom = "test"
	flow.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount.Denom, "test")
	flow.Conditions.FeedbackLoops = []*types.FeedbackLoop{{ResponseIndex: 0, ResponseKey: "Amount.[0].Denom", MsgsIndex: 0, MsgKey: "Amount.Denom", ValueType: "string"}}
	err = keeper.RunFeedbackLoops(ctx, flow.ID, &flow.Msgs, flow.Conditions)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(flow.Msgs[0], &msgDelegate)
	require.NoError(t, err)
	require.Equal(t, msgDelegate.Amount.Denom, "stake")

	executedLocally, _, err = keeper.TriggerFlow(ctx, &flow)
	require.NoError(t, err)
	require.Equal(t, int64(-1), executedLocally)
}

func TestParseInnerStringFail(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flow := createBaseflow(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flow)
	require.NoError(t, err)
	require.Equal(t, int64(-1), executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	msgDelegate.Amount.Denom = "test"
	flow.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount.Denom, "test")
	flow.Conditions.FeedbackLoops = []*types.FeedbackLoop{{ResponseIndex: 0, ResponseKey: "Amount.[0].Amount", MsgsIndex: 0, MsgKey: "Amount.Denom", ValueType: "string"}}
	err = keeper.RunFeedbackLoops(ctx, flow.ID, &flow.Msgs, flow.Conditions)
	require.Error(t, err)

}

func TestParseInnerInt(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flow := createBaseflow(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flow)
	require.NoError(t, err)
	require.Equal(t, int64(-1), executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	flow.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(1000)))
	flow.Conditions.FeedbackLoops = []*types.FeedbackLoop{{ResponseIndex: 0, ResponseKey: "Amount.[0].Amount", MsgsIndex: 0, MsgKey: "Amount.Amount", ValueType: "math.Int"}}
	err = keeper.RunFeedbackLoops(ctx, flow.ID, &flow.Msgs, flow.Conditions)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(flow.Msgs[0], &msgDelegate)
	require.NoError(t, err)
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(101)))

	executedLocally, _, err = keeper.TriggerFlow(ctx, &flow)
	require.NoError(t, err)
	require.Equal(t, int64(-1), executedLocally)
}

func TestCompareInnerIntTrue(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flow := createBaseflow(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flow)
	require.NoError(t, err)
	require.Equal(t, int64(-1), executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	flow.Conditions = &types.ExecutionConditions{}
	flow.Conditions.Comparisons = []*types.Comparison{{ResponseIndex: 0, ResponseKey: "Amount.[0].Amount", ValueType: "math.Int", Operator: 0, Operand: "101"}}
	boolean, err := keeper.CompareResponseValue(ctx, flow.ID, msgResponses, *flow.Conditions.Comparisons[0])
	require.NoError(t, err)

	require.True(t, boolean)
}

func TestCompareCoinTrue(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flow := createBaseflow(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flow)
	require.NoError(t, err)
	require.Equal(t, int64(-1), executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})
	flow.Conditions = &types.ExecutionConditions{}
	flow.Conditions.Comparisons = []*types.Comparison{{ResponseIndex: 0, ResponseKey: "Amount.[0]", ValueType: "sdk.Coin", Operator: 0, Operand: "101stake"}}
	boolean, err := keeper.CompareResponseValue(ctx, flow.ID, msgResponses, *flow.Conditions.Comparisons[0])
	require.NoError(t, err)

	require.True(t, boolean)
}

func TestCompareCoinLargerThanTrue(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flow := createBaseflow(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flow)
	require.NoError(t, err)
	require.Equal(t, int64(-1), executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	flow.Conditions = &types.ExecutionConditions{}
	comparison := types.Comparison{ResponseIndex: 0, ResponseKey: "Amount.[0]", ValueType: "sdk.Coin", Operator: 4, Operand: "11stake"}
	flow.Conditions.Comparisons = []*types.Comparison{&comparison}

	boolean, err := keeper.CompareResponseValue(ctx, flow.ID, msgResponses, *flow.Conditions.Comparisons[0])
	require.NoError(t, err)
	require.True(t, boolean)

	boolean, err = keeper.allowedToExecute(ctx, flow)
	require.NoError(t, err)
	require.True(t, boolean)
}

func TestCompareIntFalse(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flow := createBaseflow(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flow)
	require.NoError(t, err)
	require.Equal(t, int64(-1), executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	flow.Conditions = &types.ExecutionConditions{}
	flow.Conditions.Comparisons = []*types.Comparison{{ResponseIndex: 0, ResponseKey: "Amount.[0].Amount", ValueType: "math.Int", Operator: 0, Operand: "100000000000"}}
	boolean, err := keeper.CompareResponseValue(ctx, flow.ID, msgResponses, *flow.Conditions.Comparisons[0])
	require.NoError(t, err)

	require.False(t, boolean)
}

func TestCompareDenomString(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flow := createBaseflow(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flow)
	require.NoError(t, err)
	require.Equal(t, int64(-1), executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	flow.Conditions = &types.ExecutionConditions{}
	flow.Conditions.Comparisons = []*types.Comparison{{ResponseIndex: 0, ResponseKey: "Amount.[0].Denom", ValueType: "string", Operator: 0, Operand: "stake"}}
	boolean, err := keeper.CompareResponseValue(ctx, flow.ID, msgResponses, *flow.Conditions.Comparisons[0])
	require.NoError(t, err)

	require.True(t, boolean)
}

func TestInvalidResponseIndex(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flow := createBaseflow(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flow)
	require.NoError(t, err)
	require.Equal(t, int64(-1), executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	flow.Conditions = &types.ExecutionConditions{}
	flow.Conditions.Comparisons = []*types.Comparison{{ResponseIndex: 1, ResponseKey: "Amount.[0].Denom", ValueType: "string", Operator: 0, Operand: "sta"}}
	_, err = keeper.CompareResponseValue(ctx, flow.ID, msgResponses, *flow.Conditions.Comparisons[0])
	require.Error(t, err)

	require.Contains(t, err.Error(), "number of responses")
}

func TestCompareDenomStringContains(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flow := createBaseflow(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flow)
	require.NoError(t, err)
	require.Equal(t, int64(-1), executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	flow.Conditions = &types.ExecutionConditions{}
	flow.Conditions.Comparisons = []*types.Comparison{{ResponseIndex: 0, ResponseKey: "Amount.[0].Denom", ValueType: "string", Operator: 1, Operand: "sta"}}
	boolean, err := keeper.CompareResponseValue(ctx, flow.ID, msgResponses, *flow.Conditions.Comparisons[0])
	require.NoError(t, err)

	require.True(t, boolean)
}

func TestCompareArrayCoinsContainsTrue(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flow := createBaseflow(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flow)
	require.NoError(t, err)
	require.Equal(t, int64(-1), executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	flow.Conditions = &types.ExecutionConditions{}
	flow.Conditions.Comparisons = []*types.Comparison{{ResponseIndex: 0, ResponseKey: "Amount", ValueType: "sdk.Coins", Operator: 1, Operand: "101stake"}}
	boolean, err := keeper.CompareResponseValue(ctx, flow.ID, msgResponses, *flow.Conditions.Comparisons[0])
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
	flow := createBaseflow(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	flow.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(1000)))
	flow.Conditions.FeedbackLoops = []*types.FeedbackLoop{{ResponseIndex: 0, ResponseKey: "", MsgsIndex: 0, MsgKey: "Amount.Amount", ValueType: "math.Int"}}
	queryCallback, err := math.NewInt(39999999999).Marshal()
	require.NoError(t, err)
	flow.Conditions.FeedbackLoops[0].ICQConfig = &types.ICQConfig{Response: queryCallback}
	err = keeper.RunFeedbackLoops(ctx, flow.ID, &flow.Msgs, flow.Conditions)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(flow.Msgs[0], &msgDelegate)
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
	flow := createBaseflow(delAddr, flowAddr)

	// Create three test messages
	msg1 := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	msg2 := newFakeMsgDelegate(delAddr, val)
	msg3 := newFakeMsgDelegate(delAddr, val)

	// Set initial amounts
	msg2.Amount = sdk.NewCoin("stake", math.NewInt(1000))
	msg3.Amount = sdk.NewCoin("stake", math.NewInt(2000))

	// Pack all messages
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msg1, msg2, msg3})

	// Set up feedback loops that will cause potential duplicates
	flow.Conditions = &types.ExecutionConditions{
		FeedbackLoops: []*types.FeedbackLoop{
			{
				// First feedback loop after first message
				ResponseIndex: 0,
				ResponseKey:   "Amount.[0].Amount",
				MsgsIndex:     1,
				MsgKey:        "Amount.Amount",
				ValueType:     "math.Int",
			},
			{
				// Second feedback loop after first message (same point as first)
				ResponseIndex: 0,
				ResponseKey:   "Amount.[0].Amount",
				MsgsIndex:     2,
				MsgKey:        "Amount.Amount",
				ValueType:     "math.Int",
			},
		},
	}

	// Execute the flow
	executedLocally, responses, err := keeper.TriggerFlow(ctx, &flow)
	require.NoError(t, err)
	require.Equal(t, int64(-1), executedLocally)

	// Verify we have the expected number of responses (3 messages)
	require.Len(t, responses, 3, "expected 3 message responses")

	// Create a flow history entry with the actual responses
	historyEntry := &types.FlowHistoryEntry{
		ScheduledExecTime: flow.ExecTime,
		ActualExecTime:    flow.ExecTime,
		ExecFee:           sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(0))),
		Executed:          true,
		TimedOut:          false,
		MsgResponses:      responses, // Use the actual responses from TriggerFlow
		Errors:            []string{},
	}
	keeper.SetFlowHistoryEntry(ctx, flow.ID, historyEntry)
	history, _ := keeper.GetFlowHistory(ctx, flow.ID)

	// Manually run the feedback loops to ensure they're processed
	err = keeper.RunFeedbackLoops(ctx, flow.ID, &flow.Msgs, flow.Conditions)
	require.NoError(t, err)

	// Verify we have exactly one history entry
	require.Len(t, history, 1, "expected exactly one history entry")

	// Verify the history entry has the correct number of message responses
	require.Len(t, history[0].MsgResponses, 3, "expected exactly 3 message responses in history")
	require.Empty(t, history[0].Errors, "expected no errors in history entry")

	// Verify the messages were modified by the feedback loops
	// Unpack the messages from the flow to check if they were modified

	var modifiedMsg2, modifiedMsg3 *stakingtypes.MsgDelegate
	err = keeper.cdc.UnpackAny(flow.Msgs[1], &modifiedMsg2)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(flow.Msgs[2], &modifiedMsg3)
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
	flow := createBaseflow(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	flow.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(1000)))
	flow.Conditions.FeedbackLoops = []*types.FeedbackLoop{{ResponseIndex: 0, ResponseKey: "", MsgsIndex: 0, MsgKey: "Amount", ValueType: "sdk.Coin"}}
	coin := sdk.NewCoin("stake", math.NewInt(39999999999))
	queryCallback, err := coin.Marshal()
	require.NoError(t, err)
	flow.Conditions.FeedbackLoops[0].ICQConfig = &types.ICQConfig{Response: queryCallback}
	err = keeper.RunFeedbackLoops(ctx, flow.ID, &flow.Msgs, flow.Conditions)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(flow.Msgs[0], &msgDelegate)
	require.NoError(t, err)
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(39999999999)))
}

func TestParseStringICQ(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flow := createBaseflow(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	flow.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", math.NewInt(1000)))
	flow.Conditions.FeedbackLoops = []*types.FeedbackLoop{{ResponseIndex: 0, ResponseKey: "", MsgsIndex: 0, MsgKey: "Amount.Denom", ValueType: "string", ICQConfig: &types.ICQConfig{ConnectionId: "connection-0", ChainId: "test"}}}
	queryCallback := []byte("fuzz")
	flow.Conditions.FeedbackLoops[0].ICQConfig = &types.ICQConfig{Response: queryCallback}
	err := keeper.RunFeedbackLoops(ctx, flow.ID, &flow.Msgs, flow.Conditions)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(flow.Msgs[0], &msgDelegate)
	require.NoError(t, err)
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("fuzz", math.NewInt(1000)))
}

func TestCompareCoinTrueICQ(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flow := createBaseflow(delAddr, flowAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerFlow(ctx, &flow)
	require.NoError(t, err)
	require.Equal(t, int64(-1), executedLocally)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: msgResponses})

	flow.Conditions = &types.ExecutionConditions{}
	flow.Conditions.Comparisons = []*types.Comparison{{ResponseIndex: 0, ResponseKey: "", ValueType: "sdk.Coin", Operator: 4, Operand: "101stake"}}

	coin := sdk.NewCoin("stake", math.NewInt(39999999999))
	queryCallback, err := coin.Marshal()
	require.NoError(t, err)
	flow.Conditions.Comparisons[0].ICQConfig = &types.ICQConfig{Response: queryCallback}
	boolean, err := keeper.CompareResponseValue(ctx, flow.ID, msgResponses, *flow.Conditions.Comparisons[0])
	require.NoError(t, err)

	require.True(t, boolean)
}

func TestCompareFromWasmResponse(t *testing.T) {
	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))

	b64 := "eyJuYW1lIjoiQVRDIC8gVVNEQyIsInRyZWFzdXJ5Ijoib3NtbzF4Zjh0Nzg5bGUyeHB5NzllZjQ1dnBlZjVodDVkOWd1eG01MGFhZSIsInVybCI6Imh0dHBzOi8vd3d3LmFzdG9uaWMuaW8vIiwiZGlzdF9pbmRleCI6IjQzMDA4NTE3Mjk5Mjg3NTIuODUwMjMwODc1MTI0MjY0MDI1IiwibGFzdF91cGRhdGVkIjoiMTc0NzMxNzYwODQxNDQxMjk3NSIsIm91dF9kZW5vbSI6ImliYy84NTZDNkYwMUU0OEVDOUE4OUVGRjM1QjJCNzVDM0JDMENERjUxQzIzQUIwOTRDQ0M2QUUzMUI2RDMwQzJBNkVEIiwib3V0X3N1cHBseSI6IjUwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwIiwib3V0X3JlbWFpbmluZyI6IjAiLCJpbl9kZW5vbSI6ImliYy80OThBMDc1MUM3OThBMEQ5QTM4OUFBMzY5MTEyM0RBREE1N0RBQTRGRTE2NUQ1Qzc1ODk0NTA1Qjg3NkJBNkU0IiwiaW5fc3VwcGx5IjoiMCIsInNwZW50X2luIjoiMTE2Mzk3NTA0MzIiLCJzaGFyZXMiOiI3NDQxNTY2MyIsInN0YXJ0X3RpbWUiOiIxNzQ3MzE0MDAwMDAwMDAwMDAwIiwiZW5kX3RpbWUiOiIxNzQ3MzE3NjAwMDAwMDAwMDAwIiwiY3VycmVudF9zdHJlYW1lZF9wcmljZSI6IjAuMDAwMDAwMDAwMDAwMDAwMjU2Iiwic3RhdHVzIjoiZmluYWxpemVkIiwicGF1c2VfZGF0ZSI6bnVsbCwic3RyZWFtX2NyZWF0aW9uX2Rlbm9tIjoidW9zbW8iLCJzdHJlYW1fY3JlYXRpb25fZmVlIjoiMzAwMDAwMDAwIiwic3RyZWFtX2V4aXRfZmVlX3BlcmNlbnQiOiIwLjA0MiJ9"
	queryKey := "current_streamed_price"
	compareValue := "0.000000000000000256"

	decoded, err := base64.StdEncoding.DecodeString(b64)
	require.NoError(t, err)

	comparison := types.Comparison{
		ResponseIndex: 0,
		ResponseKey:   queryKey,
		Operand:       compareValue,
		Operator:      types.ComparisonOperator_EQUAL,
		ValueType:     "math.Dec",
		ICQConfig:     &types.ICQConfig{Response: decoded},
	}

	ok, err := keeper.CompareResponseValue(ctx, 1, nil, comparison)
	require.NoError(t, err)
	require.True(t, ok, "comparison should succeed for equality")

	// Change operator to NOT_EQUAL and check
	comparison.Operator = types.ComparisonOperator_NOT_EQUAL
	ok, err = keeper.CompareResponseValue(ctx, 1, nil, comparison)
	require.NoError(t, err)
	require.False(t, ok, "comparison should fail for not equal")

	compareValue = "0.000000000000000300"
	comparison.Operand = compareValue
	comparison.Operator = types.ComparisonOperator_SMALLER_THAN
	ok, err = keeper.CompareResponseValue(ctx, 1, nil, comparison)
	require.NoError(t, err, "error comparing value: %v", err)
	require.True(t, ok, "comparison should succeed for smaller than")

	compareValue = "0.10"
	ok, err = keeper.CompareResponseValue(ctx, 1, nil, comparison)
	require.NoError(t, err, "error comparing value: %v", err)
	require.True(t, ok, "comparison should succeed for smaller than")

	comparison.Operator = types.ComparisonOperator_LARGER_THAN
	ok, err = keeper.CompareResponseValue(ctx, 1, nil, comparison)
	require.NoError(t, err, "error comparing value: %v", err)
	require.False(t, ok, "comparison should not succeed for larger than")

	compareValue = "aSDasfa"
	comparison.Operand = compareValue
	comparison.Operator = types.ComparisonOperator_SMALLER_THAN
	ok, err = keeper.CompareResponseValue(ctx, 1, nil, comparison)
	require.Error(t, err, "error comparing value: %v", err)
	require.False(t, ok, "comparison should not succeed for smaller than")

}

func TestFeedbackLoopFromWasmResponse(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	b64 := "eyJtaW5fc3RyZWFtX3NlY29uZHMiOiIxMjAiLCJtaW5fc2Vjb25kc191bnRpbF9zdGFydF90aW1lIjoiMTIwIiwiYWNjZXB0ZWRfaW5fZGVub20iOiJmYWN0b3J5L29zbW8xbno3cWRwN2VnMzBzcjk1OXd2cnduOWo5MzcwaDR4dDZ0dG0waDMvdXNzb3NtbyIsInN0cmVhbV9jcmVhdGlvbl9kZW5vbSI6InVvc21vIiwic3RyZWFtX2NyZWF0aW9uX2ZlZSI6IjEwMDAwMDAwIiwiZXhpdF9mZWVfcGVyY2VudCI6IjAuMSIsImZlZV9jb2xsZWN0b3IiOiJvc21vMW56N3FkcDdlZzMwc3I5NTl3dnJ3bjlqOTM3MGg0eHQ2dHRtMGgzIiwicHJvdG9jb2xfYWRtaW4iOiJvc21vMW56N3FkcDdlZzMwc3I5NTl3dnJ3bjlqOTM3MGg0eHQ2dHRtMGgzIn0="
	responseKey := "stream_creation_fee"

	decoded, err := base64.StdEncoding.DecodeString(b64)
	require.NoError(t, err)

	conditions := &types.ExecutionConditions{
		FeedbackLoops: []*types.FeedbackLoop{{
			MsgsIndex:     0,
			ResponseIndex: 0,
			ResponseKey:   responseKey,
			ValueType:     "math.Int",
			ICQConfig:     &types.ICQConfig{Response: decoded},
			MsgKey:        "Amount.Amount", // Example, adjust as needed
		}},
	}

	val, ctx := delegateTokens(t, ctx, keeper, delAddr)

	msgDelegate := newFakeMsgDelegate(sdk.AccAddress("test"), val)
	msgDelegate.Amount = sdk.NewCoin("stake", math.NewInt(1000))
	msgs, err := types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	require.NoError(t, err)

	err = keeper.RunFeedbackLoops(ctx, 1, &msgs, conditions)
	require.NoError(t, err)

	err = keeper.cdc.UnpackAny(msgs[0], &msgDelegate)
	require.NoError(t, err)

	expectedAmount := sdk.NewCoin("stake", math.NewInt(10000000))
	require.Equal(t, expectedAmount, msgDelegate.Amount, "amount should be updated by feedback loop")
}

func TestTwapResponse(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	b64 := "CMIYEkRpYmMvNDk4QTA3NTFDNzk4QTBEOUEzODlBQTM2OTExMjNEQURBNTdEQUE0RkUxNjVENUM3NTg5NDUwNUI4NzZCQTZFNBpEaWJjL0JFMDcyQzAzREE1NDRDRjI4MjQ5OTQxOEU3QkM2NEQzODYxNDg3OUIzRUU5NUY5QUQ5MUU2QzM3MjY3RDQ4MzYg17jpFCoLCI6u+cUGEOqJ9AYyDzU2ODkyNjgxMzUwMDA3MDoWMTc1NzY5NTMyNDM3Mzg3MTQ5MTczOUIYMTU4NTY4MzE5NDQwNzI3Nzg1OTQ2NjM5Sh41Nzg4MDE4MzIzNjcxNjAzNzYyMzIxNjY1ODI5OTdSHS0zMjA0MDE2MzUxNDcxMzQ0Njk0NDk3MDEzMDg2WgsIkonnxQYQ1orCBQ=="
	decoded, err := base64.StdEncoding.DecodeString(b64)
	require.NoError(t, err)

	comparison := types.Comparison{
		ResponseIndex: 0,
		ResponseKey:   "",
		Operand:       "0.000000000000000100",
		Operator:      types.ComparisonOperator_LARGER_THAN,
		ValueType:     "osmosistwapv1beta1.TwapRecord",
		ICQConfig:     &types.ICQConfig{Response: decoded},
	}

	ok, err := keeper.CompareResponseValue(ctx, 1, nil, comparison)
	require.NoError(t, err)
	require.True(t, ok, "comparison should succeed for larger than")

	comparison.ValueType = "osmosistwapv1beta1.TwapRecord.P1LastSpotPrice"
	ok, err = keeper.CompareResponseValue(ctx, 1, nil, comparison)
	require.NoError(t, err)
	require.True(t, ok, "comparison should succeed for larger than")

	comparison.ValueType = "osmosistwapv1beta1.TwapRecord.dffadsf"
	_, err = keeper.CompareResponseValue(ctx, 1, nil, comparison)
	require.Error(t, err)

	comparison.ValueType = "osmosistwapv1beta1.TwapRecord.LastErrorTime"
	_, err = keeper.CompareResponseValue(ctx, 1, nil, comparison)
	require.Error(t, err)

	feedbackLoop := types.FeedbackLoop{
		ResponseIndex: 0,
		ResponseKey:   "",
		ValueType:     "osmosistwapv1beta1.TwapRecord",
		ICQConfig:     &types.ICQConfig{Response: decoded},
		MsgKey:        "Amount.Amount",
	}
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	msgDelegate := newFakeMsgDelegate(sdk.AccAddress("test"), val)
	msgs, err := types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	require.NoError(t, err)

	err = keeper.RunFeedbackLoops(ctx, 1, &msgs, &types.ExecutionConditions{FeedbackLoops: []*types.FeedbackLoop{&feedbackLoop}})
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot assign math.Dec to math.Int")
}

func TestBalanceResponse(t *testing.T) {
	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	b64 := "NTE1MjM2NjU2OQ=="
	decoded, err := base64.StdEncoding.DecodeString(b64)
	require.NoError(t, err)

	comparison := types.Comparison{
		ResponseIndex: 0,
		ResponseKey:   "balance",
		Operand:       "100",
		Operator:      types.ComparisonOperator_LARGER_THAN,
		ValueType:     "math.Int",
		ICQConfig:     &types.ICQConfig{Response: decoded},
	}

	ok, err := keeper.CompareResponseValue(ctx, 1, nil, comparison)
	require.NoError(t, err)
	require.True(t, ok, "comparison should succeed for larger than")

}

// func TestCompareStringFromWasmResponse(t *testing.T) {
// 	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))

// 	b64 := "eyJtaW5fc3RyZWFtX3NlY29uZHMiOiIxMjAiLCJtaW5fc2Vjb25kc191bnRpbF9zdGFydF90aW1lIjoiMTIwIiwiYWNjZXB0ZWRfaW5fZGVub20iOiJmYWN0b3J5L29zbW8xbno3cWRwN2VnMzBzcjk1OXd2cnduOWo5MzcwaDR4dDZ0dG0waDMvdXNzb3NtbyIsInN0cmVhbV9jcmVhdGlvbl9kZW5vbSI6InVvc21vIiwic3RyZWFtX2NyZWF0aW9uX2ZlZSI6IjEwMDAwMDAwIiwiZXhpdF9mZWVfcGVyY2VudCI6IjAuMSIsImZlZV9jb2xsZWN0b3IiOiJvc21vMW56N3FkcDdlZzMwc3I5NTl3dnJ3bjlqOTM3MGg0eHQ2dHRtMGgzIiwicHJvdG9jb2xfYWRtaW4iOiJvc21vMW56N3FkcDdlZzMwc3I5NTl3dnJ3bjlqOTM3MGg0eHQ2dHRtMGgzIn0="
// 	queryKey := "fee_collector"
// 	compareValue := "osmo1nz7qdp7eg30sr959wvrwn9j9370h4xt6ttm0h3"

// 	decoded, err := base64.StdEncoding.DecodeString(b64)
// 	require.NoError(t, err)

// 	comparison := types.Comparison{
// 		ResponseIndex: 0,
// 		ResponseKey:   queryKey,
// 		Operand:       compareValue,
// 		Operator:      types.ComparisonOperator_EQUAL,
// 		ValueType:     "string",
// 		ICQConfig:     &types.ICQConfig{Response: decoded},
// 	}

// 	ok, err := keeper.CompareResponseValue(ctx, 1, nil, comparison)
// 	require.NoError(t, err)
// 	require.True(t, ok, "comparison should succeed for equality on string field")

// 	// Change operator to NOT_EQUAL and check
// 	comparison.Operator = types.ComparisonOperator_NOT_EQUAL
// 	ok, err = keeper.CompareResponseValue(ctx, 1, nil, comparison)
// 	require.NoError(t, err)
// 	require.False(t, ok, "comparison should fail for not equal on string field")

// 	compareValue = "osmo1somethingelse"
// 	comparison.Operand = compareValue
// 	comparison.Operator = types.ComparisonOperator_EQUAL
// 	ok, err = keeper.CompareResponseValue(ctx, 1, nil, comparison)
// 	require.NoError(t, err)
// 	require.False(t, ok, "comparison should fail for wrong string value")
// }
