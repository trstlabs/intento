package keeper

import (
	"testing"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/stretchr/testify/require"
	"github.com/trstlabs/intento/x/intent/types"
)

func TestParseCoin(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))
	actionAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	actionInfo := createBaseActionInfo(delAddr, actionAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerAction(ctx, &actionInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetActionHistoryEntry(ctx, actionInfo.ID, &types.ActionHistoryEntry{MsgResponses: msgResponses})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	actionInfo.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", sdk.NewInt(1000)))
	actionInfo.Conditions.UseResponseValue = &types.UseResponseValue{ResponseIndex: 0, ResponseKey: "Amount", MsgsIndex: 0, MsgKey: "Amount", ValueType: "sdk.Coin"}
	err = keeper.UseResponseValue(ctx, actionInfo.ID, &actionInfo.Msgs, actionInfo.Conditions, nil)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(actionInfo.Msgs[0], &msgDelegate)
	require.NoError(t, err)
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", sdk.NewInt(101)))

	executedLocally, _, err = keeper.TriggerAction(ctx, &actionInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
}

func TestParseInnerString(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))
	actionAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	actionInfo := createBaseActionInfo(delAddr, actionAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerAction(ctx, &actionInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetActionHistoryEntry(ctx, actionInfo.ID, &types.ActionHistoryEntry{MsgResponses: msgResponses})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	msgDelegate.Amount.Denom = "test"
	actionInfo.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount.Denom, "test")
	actionInfo.Conditions.UseResponseValue = &types.UseResponseValue{ResponseIndex: 0, ResponseKey: "Amount.[0].Denom", MsgsIndex: 0, MsgKey: "Amount.Denom", ValueType: "string"}
	err = keeper.UseResponseValue(ctx, actionInfo.ID, &actionInfo.Msgs, actionInfo.Conditions, nil)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(actionInfo.Msgs[0], &msgDelegate)
	require.NoError(t, err)
	require.Equal(t, msgDelegate.Amount.Denom, "stake")

	executedLocally, _, err = keeper.TriggerAction(ctx, &actionInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
}

func TestParseInnerStringFail(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))
	actionAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	actionInfo := createBaseActionInfo(delAddr, actionAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerAction(ctx, &actionInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetActionHistoryEntry(ctx, actionInfo.ID, &types.ActionHistoryEntry{MsgResponses: msgResponses})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	msgDelegate.Amount.Denom = "test"
	actionInfo.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount.Denom, "test")
	actionInfo.Conditions.UseResponseValue = &types.UseResponseValue{ResponseIndex: 0, ResponseKey: "Amount.[0].Amount", MsgsIndex: 0, MsgKey: "Amount.Denom", ValueType: "string"}
	err = keeper.UseResponseValue(ctx, actionInfo.ID, &actionInfo.Msgs, actionInfo.Conditions, nil)
	require.Error(t, err)

}

func TestParseInnerInt(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))
	actionAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	actionInfo := createBaseActionInfo(delAddr, actionAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerAction(ctx, &actionInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetActionHistoryEntry(ctx, actionInfo.ID, &types.ActionHistoryEntry{MsgResponses: msgResponses})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	actionInfo.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", sdk.NewInt(1000)))
	actionInfo.Conditions.UseResponseValue = &types.UseResponseValue{ResponseIndex: 0, ResponseKey: "Amount.[0].Amount", MsgsIndex: 0, MsgKey: "Amount.Amount", ValueType: "sdk.Int"}
	err = keeper.UseResponseValue(ctx, actionInfo.ID, &actionInfo.Msgs, actionInfo.Conditions, nil)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(actionInfo.Msgs[0], &msgDelegate)
	require.NoError(t, err)
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", sdk.NewInt(101)))

	executedLocally, _, err = keeper.TriggerAction(ctx, &actionInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
}

func TestCompareInnerIntTrue(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))
	actionAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	actionInfo := createBaseActionInfo(delAddr, actionAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerAction(ctx, &actionInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetActionHistoryEntry(ctx, actionInfo.ID, &types.ActionHistoryEntry{MsgResponses: msgResponses})

	actionInfo.Conditions = &types.ExecutionConditions{}
	actionInfo.Conditions.ResponseComparison = &types.ResponseComparison{ResponseIndex: 0, ResponseKey: "Amount.[0].Amount", ValueType: "sdk.Int", ComparisonOperator: 0, ComparisonOperand: "101"}
	boolean, err := keeper.CompareResponseValue(ctx, actionInfo.ID, msgResponses, *actionInfo.Conditions.ResponseComparison, nil)
	require.NoError(t, err)

	require.True(t, boolean)
}

func TestCompareCoinTrue(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))
	actionAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	actionInfo := createBaseActionInfo(delAddr, actionAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerAction(ctx, &actionInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetActionHistoryEntry(ctx, actionInfo.ID, &types.ActionHistoryEntry{MsgResponses: msgResponses})

	actionInfo.Conditions = &types.ExecutionConditions{}
	actionInfo.Conditions.ResponseComparison = &types.ResponseComparison{ResponseIndex: 0, ResponseKey: "Amount.[0]", ValueType: "sdk.Coin", ComparisonOperator: 0, ComparisonOperand: "101stake"}
	boolean, err := keeper.CompareResponseValue(ctx, actionInfo.ID, msgResponses, *actionInfo.Conditions.ResponseComparison, nil)
	require.NoError(t, err)

	require.True(t, boolean)
}

func TestCompareIntFalse(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))
	actionAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	actionInfo := createBaseActionInfo(delAddr, actionAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerAction(ctx, &actionInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetActionHistoryEntry(ctx, actionInfo.ID, &types.ActionHistoryEntry{MsgResponses: msgResponses})

	actionInfo.Conditions = &types.ExecutionConditions{}
	actionInfo.Conditions.ResponseComparison = &types.ResponseComparison{ResponseIndex: 0, ResponseKey: "Amount.[0].Amount", ValueType: "sdk.Int", ComparisonOperator: 0, ComparisonOperand: "100000000000"}
	boolean, err := keeper.CompareResponseValue(ctx, actionInfo.ID, msgResponses, *actionInfo.Conditions.ResponseComparison, nil)
	require.NoError(t, err)

	require.False(t, boolean)
}

func TestCompareDenomString(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))
	actionAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	actionInfo := createBaseActionInfo(delAddr, actionAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerAction(ctx, &actionInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetActionHistoryEntry(ctx, actionInfo.ID, &types.ActionHistoryEntry{MsgResponses: msgResponses})

	actionInfo.Conditions = &types.ExecutionConditions{}
	actionInfo.Conditions.ResponseComparison = &types.ResponseComparison{ResponseIndex: 0, ResponseKey: "Amount.[0].Denom", ValueType: "string", ComparisonOperator: 0, ComparisonOperand: "stake"}
	boolean, err := keeper.CompareResponseValue(ctx, actionInfo.ID, msgResponses, *actionInfo.Conditions.ResponseComparison, nil)
	require.NoError(t, err)

	require.True(t, boolean)
}

func TestInvalidResponseIndex(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))
	actionAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	actionInfo := createBaseActionInfo(delAddr, actionAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerAction(ctx, &actionInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetActionHistoryEntry(ctx, actionInfo.ID, &types.ActionHistoryEntry{MsgResponses: msgResponses})

	actionInfo.Conditions = &types.ExecutionConditions{}
	actionInfo.Conditions.ResponseComparison = &types.ResponseComparison{ResponseIndex: 1, ResponseKey: "Amount.[0].Denom", ValueType: "string", ComparisonOperator: 0, ComparisonOperand: "sta"}
	_, err = keeper.CompareResponseValue(ctx, actionInfo.ID, msgResponses, *actionInfo.Conditions.ResponseComparison, nil)
	require.Error(t, err)

	require.Contains(t, err.Error(), "number of responses")
}

func TestCompareDenomStringContains(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))
	actionAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	actionInfo := createBaseActionInfo(delAddr, actionAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerAction(ctx, &actionInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetActionHistoryEntry(ctx, actionInfo.ID, &types.ActionHistoryEntry{MsgResponses: msgResponses})

	actionInfo.Conditions = &types.ExecutionConditions{}
	actionInfo.Conditions.ResponseComparison = &types.ResponseComparison{ResponseIndex: 0, ResponseKey: "Amount.[0].Denom", ValueType: "string", ComparisonOperator: 1, ComparisonOperand: "sta"}
	boolean, err := keeper.CompareResponseValue(ctx, actionInfo.ID, msgResponses, *actionInfo.Conditions.ResponseComparison, nil)
	require.NoError(t, err)

	require.True(t, boolean)
}

func TestCompareArrayCoinsContainsTrue(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))
	actionAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	actionInfo := createBaseActionInfo(delAddr, actionAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerAction(ctx, &actionInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetActionHistoryEntry(ctx, actionInfo.ID, &types.ActionHistoryEntry{MsgResponses: msgResponses})

	actionInfo.Conditions = &types.ExecutionConditions{}
	actionInfo.Conditions.ResponseComparison = &types.ResponseComparison{ResponseIndex: 0, ResponseKey: "Amount", ValueType: "sdk.Coins", ComparisonOperator: 1, ComparisonOperand: "101stake"}
	boolean, err := keeper.CompareResponseValue(ctx, actionInfo.ID, msgResponses, *actionInfo.Conditions.ResponseComparison, nil)
	require.NoError(t, err)

	require.True(t, boolean)
}

func TestCompareArrayCoinsContainsFalse(t *testing.T) {
	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))
	fakeCoins, _ := sdk.ParseCoinsNormalized("1000abc,3000000000degf")
	msgResponse := distrtypes.MsgWithdrawDelegatorRewardResponse{Amount: fakeCoins}
	any, _ := cdctypes.NewAnyWithValue(&msgResponse)
	msgResponses := []*cdctypes.Any{any}
	responseComparison := &types.ResponseComparison{ResponseIndex: 0, ResponseKey: "Amount", ValueType: "sdk.Coins", ComparisonOperator: 0, ComparisonOperand: "100aaa"}
	boolean, err := keeper.CompareResponseValue(ctx, 1, msgResponses, *responseComparison, nil)
	require.NoError(t, err)

	require.False(t, boolean)
}

func TestCompareArrayEquals(t *testing.T) {
	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))
	fakeCoins, _ := sdk.ParseCoinsNormalized("1000abc,3000000000degf")
	msgResponse := distrtypes.MsgWithdrawDelegatorRewardResponse{Amount: fakeCoins}
	any, _ := cdctypes.NewAnyWithValue(&msgResponse)
	msgResponses := []*cdctypes.Any{any}
	responseComparison := &types.ResponseComparison{ResponseIndex: 0, ResponseKey: "Amount", ValueType: "sdk.Coins", ComparisonOperator: 0, ComparisonOperand: "1000abc,3000000000degf"}
	boolean, err := keeper.CompareResponseValue(ctx, 1, msgResponses, *responseComparison, nil)
	require.NoError(t, err)

	require.True(t, boolean)
}

func TestParseAmountICQ(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))
	actionAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	actionInfo := createBaseActionInfo(delAddr, actionAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	actionInfo.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", sdk.NewInt(1000)))
	actionInfo.Conditions.UseResponseValue = &types.UseResponseValue{ResponseIndex: 0, ResponseKey: "", MsgsIndex: 0, MsgKey: "Amount.Amount", ValueType: "sdk.Int", FromICQ: true}
	queryCallback, err := sdk.NewInt(39999999999).Marshal()
	require.NoError(t, err)
	err = keeper.UseResponseValue(ctx, actionInfo.ID, &actionInfo.Msgs, actionInfo.Conditions, queryCallback)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(actionInfo.Msgs[0], &msgDelegate)
	require.NoError(t, err)
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", sdk.NewInt(39999999999)))
}

func TestParseCoinICQ(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))
	actionAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	actionInfo := createBaseActionInfo(delAddr, actionAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	actionInfo.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", sdk.NewInt(1000)))
	actionInfo.Conditions.UseResponseValue = &types.UseResponseValue{ResponseIndex: 0, ResponseKey: "", MsgsIndex: 0, MsgKey: "Amount", ValueType: "sdk.Coin", FromICQ: true}
	coin := sdk.NewCoin("stake", sdk.NewInt(39999999999))
	queryCallback, err := coin.Marshal()
	require.NoError(t, err)
	err = keeper.UseResponseValue(ctx, actionInfo.ID, &actionInfo.Msgs, actionInfo.Conditions, queryCallback)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(actionInfo.Msgs[0], &msgDelegate)
	require.NoError(t, err)
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", sdk.NewInt(39999999999)))
}

func TestParseStringICQ(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))
	actionAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	actionInfo := createBaseActionInfo(delAddr, actionAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})

	msgDelegate := newFakeMsgDelegate(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	actionInfo.Conditions = &types.ExecutionConditions{}
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", sdk.NewInt(1000)))
	actionInfo.Conditions.UseResponseValue = &types.UseResponseValue{ResponseIndex: 0, ResponseKey: "", MsgsIndex: 0, MsgKey: "Amount.Denom", ValueType: "string", FromICQ: true}
	text := []byte("fuzz")

	err := keeper.UseResponseValue(ctx, actionInfo.ID, &actionInfo.Msgs, actionInfo.Conditions, text)
	require.NoError(t, err)
	err = keeper.cdc.UnpackAny(actionInfo.Msgs[0], &msgDelegate)
	require.NoError(t, err)
	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("fuzz", sdk.NewInt(1000)))
}

func TestCompareCoinTrueICQ(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))
	actionAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
	actionInfo := createBaseActionInfo(delAddr, actionAddr)

	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})
	executedLocally, msgResponses, err := keeper.TriggerAction(ctx, &actionInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)
	keeper.SetActionHistoryEntry(ctx, actionInfo.ID, &types.ActionHistoryEntry{MsgResponses: msgResponses})

	actionInfo.Conditions = &types.ExecutionConditions{}
	actionInfo.Conditions.ResponseComparison = &types.ResponseComparison{ResponseIndex: 0, ResponseKey: "", ValueType: "sdk.Coin", ComparisonOperator: 4, ComparisonOperand: "101stake", FromICQ: true}

	coin := sdk.NewCoin("stake", sdk.NewInt(39999999999))
	queryCallback, err := coin.Marshal()
	require.NoError(t, err)
	boolean, err := keeper.CompareResponseValue(ctx, actionInfo.ID, msgResponses, *actionInfo.Conditions.ResponseComparison, queryCallback)
	require.NoError(t, err)

	require.True(t, boolean)
}

// func TestParseCoinAuthZ(t *testing.T) {
// 	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))
// 	actionAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
// 	types.Denom = "stake"
// 	val, ctx := delegateTokens(t, ctx, keeper, delAddr)
// 	actionInfo := createBaseActionInfo(delAddr, actionAddr)

// 	msgs := newFakeMsgAuthZWithdrawDelegatorReward(delAddr, delAddr, val)
// 	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgs})
// 	executedLocally, msgResponses, err := keeper.TriggerAction(ctx, &actionInfo)
// 	require.NoError(t, err)
// 	require.True(t, executedLocally)
// 	keeper.SetActionHistoryEntry(ctx, actionInfo.ID, &types.ActionHistoryEntry{MsgResponses: msgResponses})

// 	msgDelegate := newFakeMsgDelegate(delAddr, val)
// 	actionInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
// 	actionInfo.Conditions = &types.ExecutionConditions{}
// 	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", sdk.NewInt(1000)))
// 	actionInfo.Conditions.UseResponseValue = &types.UseResponseValue{ResponseIndex: 0, ResponseKey: "Amount", MsgsIndex: 0, MsgKey: "Amount", ValueType: "sdk.Coin"}
// 	err = keeper.UseResponseValue(ctx, actionInfo.ID, &actionInfo.Msgs, actionInfo.Conditions, nil)
// 	require.NoError(t, err)
// 	err = keeper.cdc.UnpackAny(actionInfo.Msgs[0], &msgDelegate)
// 	require.NoError(t, err)
// 	require.Equal(t, msgDelegate.Amount, sdk.NewCoin("stake", sdk.NewInt(101)))

// 	executedLocally, _, err = keeper.TriggerAction(ctx, &actionInfo)
// 	require.NoError(t, err)
// 	require.True(t, executedLocally)
// }

// func newFakeMsgAuthZWithdrawDelegatorReward(delegator, grantee sdk.AccAddress, validator stakingtypes.Validator) *authztypes.MsgExec {
// 	msgWithdrawDelegatorReward := &distrtypes.MsgWithdrawDelegatorReward{
// 		DelegatorAddress: delegator.String(),
// 		ValidatorAddress: validator.GetOperator().String(),
// 	}
// 	anys, _ := types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})

// 	msgExec := &authztypes.MsgExec{Grantee: grantee.String(), Msgs: anys}
// 	return msgExec
// }
