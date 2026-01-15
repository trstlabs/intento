package keeper

import (
	"encoding/base64"
	"testing"

	"cosmossdk.io/math"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	osmosistwapv1beta1 "github.com/trstlabs/intento/x/intent/msg_registry/osmosis/twap/v1beta1"
	"github.com/trstlabs/intento/x/intent/types"
)

// Test data
var (
	testDenom   = "stake"
	testAmount1 = math.NewInt(500)
	testAmount2 = math.NewInt(1000)
)

func TestDifferenceModeFeedbackLoopSdkCoin(t *testing.T) {
	// Setup test environment
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin(testDenom, 3_000_000)))
	types.Denom = testDenom

	valAddr, ctx := delegateTokens(t, ctx, keeper, delAddr)

	flow := createBaseflow(delAddr, flowAddr)

	msgDelegate := newFakeMsgDelegate(delAddr, valAddr)
	msgDelegate.Amount = sdk.NewCoin(testDenom, testAmount1)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})
	msgWithdrawDelegatorRewardResp := distrtypes.MsgWithdrawDelegatorRewardResponse{Amount: sdk.NewCoins(sdk.NewCoin(testDenom, testAmount1))}
	msgWithdrawDelegatorRewardRespAny, _ := cdctypes.NewAnyWithValue(&msgWithdrawDelegatorRewardResp)
	msgWithdrawDelegatorRewardResp2 := distrtypes.MsgWithdrawDelegatorRewardResponse{Amount: sdk.NewCoins(sdk.NewCoin(testDenom, testAmount2))}
	msgWithdrawDelegatorRewardRespAny2, _ := cdctypes.NewAnyWithValue(&msgWithdrawDelegatorRewardResp2)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: []*cdctypes.Any{msgWithdrawDelegatorRewardRespAny}, Executed: true})
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: []*cdctypes.Any{msgWithdrawDelegatorRewardRespAny2}, Executed: true})
	msgDelegate.Amount = sdk.NewCoin(testDenom, testAmount2)
	flow.Conditions = &types.ExecutionConditions{
		FeedbackLoops: []*types.FeedbackLoop{
			{
				ResponseIndex:  0,
				ResponseKey:    "Amount",
				MsgsIndex:      0,
				MsgKey:         "Amount",
				ValueType:      "sdk.Coin",
				DifferenceMode: true,
			},
		},
	}

	// Run feedback loops
	err := keeper.RunFeedbackLoops(ctx, flow.ID, &flow.Msgs, flow.Conditions)
	require.NoError(t, err)

	// Verify the amount was updated with the difference
	var updatedMsg stakingtypes.MsgDelegate
	err = updatedMsg.Unmarshal(flow.Msgs[0].Value)
	require.NoError(t, err)

	// Calculate expected difference: new_value - old_value
	expectedDifference := testAmount2.Sub(testAmount1)
	require.True(t, updatedMsg.Amount.Amount.Equal(expectedDifference),
		"expected difference %s, got %s", expectedDifference, updatedMsg.Amount.Amount)
}

func TestDifferenceModeFeedbackLoopMathInt(t *testing.T) {
	// Setup test environment
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin(testDenom, 3_000_000)))
	types.Denom = testDenom

	valAddr, ctx := delegateTokens(t, ctx, keeper, delAddr)

	flow := createBaseflow(delAddr, flowAddr)

	// First execution with amount1
	msgDelegate := newFakeMsgDelegate(delAddr, valAddr)
	msgDelegate.Amount = sdk.NewCoin(testDenom, testAmount1)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})

	customResponse := &distrtypes.MsgWithdrawDelegatorRewardResponse{
		Amount: sdk.NewCoins(sdk.NewInt64Coin(testDenom, 1000)),
	}
	customResponseAny, _ := cdctypes.NewAnyWithValue(customResponse)

	// Store first response
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{
		MsgResponses: []*cdctypes.Any{customResponseAny},
		Executed:     true,
	})

	// Second execution with amount2
	msgDelegate.Amount = sdk.NewCoin(testDenom, testAmount2)
	flow.Conditions = &types.ExecutionConditions{
		FeedbackLoops: []*types.FeedbackLoop{
			{
				ResponseIndex:  0,
				ResponseKey:    "Amount.[0].Amount",
				MsgsIndex:      0,
				MsgKey:         "Amount.Amount",
				ValueType:      "math.Int",
				DifferenceMode: true,
			},
		},
	}

	// Store second response with different amount
	customResponse2 := &distrtypes.MsgWithdrawDelegatorRewardResponse{
		Amount: sdk.NewCoins(sdk.NewInt64Coin(testDenom, 1500)),
	}
	customResponseAny2, _ := cdctypes.NewAnyWithValue(customResponse2)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{
		MsgResponses: []*cdctypes.Any{customResponseAny2},
		Executed:     true,
	})

	// Run feedback loops
	err := keeper.RunFeedbackLoops(ctx, flow.ID, &flow.Msgs, flow.Conditions)
	require.NoError(t, err)

	// Verify the amount was updated with the difference (1500 - 1000 = 500)
	var updatedMsg stakingtypes.MsgDelegate
	err = updatedMsg.Unmarshal(flow.Msgs[0].Value)
	require.NoError(t, err)

	expectedDifference := math.NewInt(500) // 1500 - 1000
	require.True(t, updatedMsg.Amount.Amount.Equal(expectedDifference),
		"expected difference %s, got %s", expectedDifference, updatedMsg.Amount.Amount)
}

func TestDifferenceModeFeedbackLoopNoPreviousValue(t *testing.T) {
	// Setup test environment
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin(testDenom, 3_000_000)))
	types.Denom = testDenom

	valAddr, ctx := delegateTokens(t, ctx, keeper, delAddr)

	flow := createBaseflow(delAddr, flowAddr)

	// No previous execution history
	msgDelegate := newFakeMsgDelegate(delAddr, valAddr)
	msgDelegate.Amount = sdk.NewCoin(testDenom, testAmount1)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})

	// Try to run feedback loop with difference mode but no previous value
	flow.Conditions = &types.ExecutionConditions{
		FeedbackLoops: []*types.FeedbackLoop{
			{
				ResponseIndex:  0,
				ResponseKey:    "Amount",
				MsgsIndex:      0,
				MsgKey:         "Amount",
				ValueType:      "sdk.Coin",
				DifferenceMode: true,
			},
		},
	}

	// Should return an error about missing previous value
	err := keeper.RunFeedbackLoops(ctx, flow.ID, &flow.Msgs, flow.Conditions)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no execution history available for flow 1")
}

func TestDifferenceModeFeedbackLoopUnsupportedType(t *testing.T) {
	// Setup test environment
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin(testDenom, 3_000_000)))
	types.Denom = testDenom

	valAddr, ctx := delegateTokens(t, ctx, keeper, delAddr)

	flow := createBaseflow(delAddr, flowAddr)

	// First execution
	msgDelegate := newFakeMsgDelegate(delAddr, valAddr)
	msgDelegate.Amount = sdk.NewCoin(testDenom, testAmount1)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})

	// Store first response with an unsupported type (string)
	customMsg := &distrtypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: valAddr.String(),
	}
	// customMsgAny is not used in this test
	_ = customMsg
	customResponse := &distrtypes.MsgWithdrawDelegatorRewardResponse{
		Amount: sdk.NewCoins(sdk.NewInt64Coin(testDenom, 1000)),
	}
	customResponseAny, _ := cdctypes.NewAnyWithValue(customResponse)

	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{
		MsgResponses: []*cdctypes.Any{customResponseAny},
		Executed:     true,
	})

	// Second execution with difference mode on an unsupported type
	flow.Conditions = &types.ExecutionConditions{
		FeedbackLoops: []*types.FeedbackLoop{
			{
				ResponseIndex:  0,
				ResponseKey:    "DelegatorAddress", // This will be a string, which is not supported for difference mode
				MsgsIndex:      0,
				MsgKey:         "Amount",
				ValueType:      "string",
				DifferenceMode: true,
			},
		},
	}

	// Store second response
	customResponse2 := &distrtypes.MsgWithdrawDelegatorRewardResponse{
		Amount: sdk.NewCoins(sdk.NewInt64Coin(testDenom, 1500)),
	}
	customResponseAny2, _ := cdctypes.NewAnyWithValue(customResponse2)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{
		MsgResponses: []*cdctypes.Any{customResponseAny2},
		Executed:     true,
	})

	// Should return an error about unsupported type
	err := keeper.RunFeedbackLoops(ctx, flow.ID, &flow.Msgs, flow.Conditions)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found in interface")
}

func TestDifferenceModeFeedbackLoopNegativeDifference(t *testing.T) {
	// Setup test environment
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin(testDenom, 3_000_000)))
	types.Denom = testDenom

	valAddr, ctx := delegateTokens(t, ctx, keeper, delAddr)

	flow := createBaseflow(delAddr, flowAddr)

	// First execution with higher amount
	msgDelegate := newFakeMsgDelegate(delAddr, valAddr)
	msgDelegate.Amount = sdk.NewCoin(testDenom, testAmount2) // 1000
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})

	// Store first response
	customResponse1 := &distrtypes.MsgWithdrawDelegatorRewardResponse{
		Amount: sdk.NewCoins(sdk.NewInt64Coin(testDenom, 2000)),
	}
	customResponseAny1, _ := cdctypes.NewAnyWithValue(customResponse1)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{
		MsgResponses: []*cdctypes.Any{customResponseAny1},
		Executed:     true,
	})

	// Second execution with lower amount
	msgDelegate.Amount = sdk.NewCoin(testDenom, testAmount1) // 500
	flow.Conditions = &types.ExecutionConditions{
		FeedbackLoops: []*types.FeedbackLoop{
			{
				ResponseIndex:  0,
				ResponseKey:    "Amount",
				MsgsIndex:      0,
				MsgKey:         "Amount",
				ValueType:      "sdk.Coin",
				DifferenceMode: true,
			},
		},
	}

	// Store second response with lower amount (1500 - 2000 = -500, but we expect absolute value 500)
	customResponse2 := &distrtypes.MsgWithdrawDelegatorRewardResponse{
		Amount: sdk.NewCoins(sdk.NewInt64Coin(testDenom, 1500)),
	}
	customResponseAny2, _ := cdctypes.NewAnyWithValue(customResponse2)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{
		MsgResponses: []*cdctypes.Any{customResponseAny2},
		Executed:     true,
	})

	// Run feedback loops
	err := keeper.RunFeedbackLoops(ctx, flow.ID, &flow.Msgs, flow.Conditions)
	require.NoError(t, err)

	// Verify the amount was updated with the absolute difference (|1500 - 2000| = 500)
	var updatedMsg stakingtypes.MsgDelegate
	err = updatedMsg.Unmarshal(flow.Msgs[0].Value)
	require.NoError(t, err)

	expectedDifference := math.NewInt(500) // |1500 - 2000|
	require.True(t, updatedMsg.Amount.Amount.Equal(expectedDifference),
		"expected absolute difference %s, got %s", expectedDifference, updatedMsg.Amount.Amount)
}

func TestDifferenceModeWithTwapRecord(t *testing.T) {
	// Setup test environment
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin(testDenom, 3_000_000)))
	types.Denom = testDenom

	flow := createBaseflow(delAddr, flowAddr)

	// First execution with initial TWAP record
	twapB641 := "CMIYEkRpYmMvNDk4QTA3NTFDNzk4QTBEOUEzODlBQTM2OTExMjNEQURBNTdEQUE0RkUxNjVENUM3NTg5NDUwNUI4NzZCQTZFNBpEaWJjL0JFMDcyQzAzREE1NDRDRjI4MjQ5OTQxOEU3QkM2NEQzODYxNDg3OUIzRUU5NUY5QUQ5MUY2QzM3MjY3RDQ4MzYg17jpFCoLCI6u+cUGEOqJ9AYyDzU2ODkyNjgxMzUwMDA3MDoWMTc1NzY5NTMyNDM3Mzg3MTQ5MTczOUIYMTU4NTY4MzE5NDQwNzI3Nzg1OTQ2NjM5Sh41Nzg4MDE4MzIzNjcxNjAzNzYyMzIxNjY1ODI5OTdSHS0zMjA0MDE2MzUxNDcxMzQ0Njk0NDk3MDEzMDg2WgsIkonnxQYQ1orCBQ=="
	decoded1, err := base64.StdEncoding.DecodeString(twapB641)
	require.NoError(t, err)

	// Store first response
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{
		QueryResponses: []string{base64.StdEncoding.EncodeToString(decoded1)},
		Executed:       true,
	})

	// // Second execution with updated TWAP record
	twapB642 := "CMIYEkRpYmMvNDk4QTA3NTFDNzk4QTBEOUEzODlBQTM2OTExMjNEQURBNTdEQUE0RkUxNjVENUM3NTg5NDUwNUI4NzZCQTZFNBpEaWJjL0JFMDcyQzAzREE1NDRDRjI4MjQ5OTQxOEU3QkM2NEQzODYxNDg3OUIzRUU5NUY5QUQ5MUY2QzM3MjY3RDQ4MzYg17jpFCoLCI6u+cUGEOqJ9AYyDzU2ODkyNjgxMzUwMDA3MDoWMTc1NzY5NTMyNDM3Mzg3MTQ5MTczOUIYMTU4NTY4MzE5NDQwNzI3Nzg1OTQ2NjM5Sh41Nzg4MDE4MzIzNjcxNjAzNzYyMzIxNjY1ODI5OTdSHS0zMjA0MDE2MzUxNDcxMzQ0Njk0NDk3MDEzMDg2WgsIkonnxQYQ1orCBQ=="
	decoded2, err := base64.StdEncoding.DecodeString(twapB642)
	require.NoError(t, err)

	// Set up feedback loop to track P1LastSpotPrice from TWAP record
	flow.Conditions = &types.ExecutionConditions{
		Comparisons: []*types.Comparison{
			{
				ResponseIndex:  0,
				ResponseKey:    "",
				ValueType:      "osmosistwapv1beta1.TwapRecord.P1LastSpotPrice",
				DifferenceMode: true,
				ICQConfig: &types.ICQConfig{
					Response: []byte(decoded2),
				},
				Operand:  "0",
				Operator: types.ComparisonOperator_EQUAL,
			},
		},
	}

	// Run feedback loops
	isTrue, err := keeper.CompareResponseValue(ctx, flow.ID, flow.Msgs, *flow.Conditions.Comparisons[0])
	require.NoError(t, err)
	require.True(t, isTrue)

}

func TestDifferenceModeComparisonMathInt(t *testing.T) {
	// Setup test environment
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin(testDenom, 3_000_000)))
	types.Denom = testDenom

	valAddr, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flow := createBaseflow(delAddr, flowAddr)

	// First execution with initial amount
	msgDelegate := newFakeMsgDelegate(delAddr, valAddr)
	msgDelegate.Amount = sdk.NewCoin(testDenom, testAmount1)
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})

	// Store first response with initial balance
	initialBalance := math.NewInt(1000)
	customResponse := &distrtypes.MsgWithdrawDelegatorRewardResponse{
		Amount: sdk.NewCoins(sdk.NewCoin(testDenom, initialBalance)),
	}
	customResponseAny, _ := cdctypes.NewAnyWithValue(customResponse)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{
		MsgResponses: []*cdctypes.Any{customResponseAny},
		Executed:     true,
	})

	// Set up comparison with difference mode and 0 operand
	comparison := &types.Comparison{
		ResponseIndex:  0,
		ResponseKey:    "Amount.[0].Amount",
		Operand:        "0",
		Operator:       types.ComparisonOperator_EQUAL,
		ValueType:      "math.Int",
		DifferenceMode: true,
	}

	// Store second response with new balance (1500)
	newBalance := math.NewInt(1500)
	customResponse2 := &distrtypes.MsgWithdrawDelegatorRewardResponse{
		Amount: sdk.NewCoins(sdk.NewCoin(testDenom, newBalance)),
	}
	customResponseAny2, _ := cdctypes.NewAnyWithValue(customResponse2)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{
		MsgResponses: []*cdctypes.Any{customResponseAny2},
		Executed:     true,
	})

	// Get the latest response from history
	history, err := keeper.GetFlowHistory(ctx, flow.ID)
	require.NoError(t, err, "should get flow history without error")
	require.True(t, len(history) >= 2, "should have at least 2 history entries")
	latestResponse := history[len(history)-1].MsgResponses

	// The difference is 1500 - 1000 = 500
	// The comparison is: difference (500) == 0
	isTrue, err := keeper.CompareResponseValue(ctx, flow.ID, latestResponse, *comparison)
	require.NoError(t, err, "CompareResponseValue should not return an error")
	require.False(t, isTrue, "Comparison should be false as 500 is not equal to 0")

	// Add another response with the same balance (1500)
	// This should result in a difference of 0
	sameBalanceResponse := &distrtypes.MsgWithdrawDelegatorRewardResponse{
		Amount: sdk.NewCoins(sdk.NewCoin(testDenom, newBalance)),
	}
	sameBalanceResponseAny, _ := cdctypes.NewAnyWithValue(sameBalanceResponse)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{
		MsgResponses: []*cdctypes.Any{sameBalanceResponseAny},
		Executed:     true,
	})

	// Get the updated history
	history, err = keeper.GetFlowHistory(ctx, flow.ID)
	require.NoError(t, err, "should get flow history without error")
	require.True(t, len(history) >= 3, "should have at least 3 history entries")
	updatedResponse := history[len(history)-1].MsgResponses

	// The difference is 1500 - 1500 = 0
	// The comparison is: difference (0) == 0
	isTrue, err = keeper.CompareResponseValue(ctx, flow.ID, updatedResponse, *comparison)
	require.NoError(t, err, "CompareResponseValue should not return an error")
	require.True(t, isTrue, "Comparison should be true as 0 is equal to 0")
}

func TestDifferenceModeWithTwapRecordSpotPrice(t *testing.T) {
	// Setup test environment
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin(testDenom, 3_000_000)))
	types.Denom = testDenom

	valAddr, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flow := createBaseflow(delAddr, flowAddr)

	twapRecord := osmosistwapv1beta1.TwapRecord{
		P1LastSpotPrice:             math.LegacyNewDecWithPrec(499, 3),
		P0LastSpotPrice:             math.LegacyNewDecWithPrec(2345, 3),
		PoolId:                      1,
		GeometricTwapAccumulator:    math.LegacyNewDecWithPrec(100, 3),
		P1ArithmeticTwapAccumulator: math.LegacyNewDecWithPrec(654, 3),
		P0ArithmeticTwapAccumulator: math.LegacyNewDecWithPrec(857, 3),
	}
	twapRecordAny, err := cdctypes.NewAnyWithValue(&twapRecord)

	require.NoError(t, err)

	// Set up feedback loop to track P1LastSpotPrice from TWAP record
	flow.Conditions = &types.ExecutionConditions{
		Comparisons: []*types.Comparison{
			{
				ResponseIndex:  0,
				ResponseKey:    "",
				Operand:        "osmosistwapv1beta1.TwapRecord.GeometricTwapAccumulator",
				ValueType:      "osmosistwapv1beta1.TwapRecord",
				Operator:       types.ComparisonOperator_SMALLER_THAN,
				DifferenceMode: true,
				ICQConfig: &types.ICQConfig{
					Response: twapRecordAny.Value,
				},
			},
		},
	}

	// Store earlier response for geometric twap accumulator
	twapRecord.GeometricTwapAccumulator = math.LegacyNewDecWithPrec(200, 3)
	twapRecordAny, _ = cdctypes.NewAnyWithValue(&twapRecord)
	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{
		QueryResponses: []string{base64.StdEncoding.EncodeToString(twapRecordAny.Value)},
		Executed:       true,
	})

	// Create a message to update
	msgDelegate := newFakeMsgDelegate(delAddr, valAddr)
	msgDelegate.Amount = sdk.NewCoin(testDenom, math.NewInt(1000))
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})

	// Run feedback loops
	isTrue, err := keeper.CompareResponseValue(ctx, flow.ID, flow.Msgs, *flow.Conditions.Comparisons[0])
	require.NoError(t, err)
	require.True(t, isTrue)

}

func TestTwapAccelerationPositiveTrend(t *testing.T) {
	// This test verifies TWAP acceleration detection:
	// When Δaccumulator₂ > Δaccumulator₁, it indicates a positive price trend
	// (price × time)_now > (price × time)_previous

	// Setup test environment
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))
	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin(testDenom, 3_000_000)))
	types.Denom = testDenom

	valAddr, ctx := delegateTokens(t, ctx, keeper, delAddr)
	flow := createBaseflow(delAddr, flowAddr)

	// Interval 0: Initial accumulator value = 1000
	// This establishes the baseline
	twapRecord0 := osmosistwapv1beta1.TwapRecord{
		PoolId:                      1,
		P1ArithmeticTwapAccumulator: math.LegacyNewDecWithPrec(1000, 0), // 1000
		P0ArithmeticTwapAccumulator: math.LegacyNewDecWithPrec(500, 0),
		P1LastSpotPrice:             math.LegacyNewDecWithPrec(100, 0),
		P0LastSpotPrice:             math.LegacyNewDecWithPrec(50, 0),
		GeometricTwapAccumulator:    math.LegacyNewDecWithPrec(750, 0),
	}
	twapRecordAny0, err := cdctypes.NewAnyWithValue(&twapRecord0)
	require.NoError(t, err)

	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{
		QueryResponses: []string{base64.StdEncoding.EncodeToString(twapRecordAny0.Value)},
		Executed:       true,
	})

	// Interval 1: Accumulator = 1300
	// Δaccumulator₁ = 1300 - 1000 = 300
	twapRecord1 := osmosistwapv1beta1.TwapRecord{
		PoolId:                      1,
		P1ArithmeticTwapAccumulator: math.LegacyNewDecWithPrec(1300, 0), // 1300
		P0ArithmeticTwapAccumulator: math.LegacyNewDecWithPrec(650, 0),
		P1LastSpotPrice:             math.LegacyNewDecWithPrec(130, 0),
		P0LastSpotPrice:             math.LegacyNewDecWithPrec(65, 0),
		GeometricTwapAccumulator:    math.LegacyNewDecWithPrec(975, 0),
	}
	twapRecordAny1, err := cdctypes.NewAnyWithValue(&twapRecord1)
	require.NoError(t, err)

	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{
		QueryResponses: []string{base64.StdEncoding.EncodeToString(twapRecordAny1.Value)},
		Executed:       true,
	})

	// Interval 2: Accumulator = 1700
	// Δaccumulator₂ = 1700 - 1300 = 400
	// Since Δaccumulator₂ (400) > Δaccumulator₁ (300), this is acceleration (positive trend)
	twapRecord2 := osmosistwapv1beta1.TwapRecord{
		PoolId:                      1,
		P1ArithmeticTwapAccumulator: math.LegacyNewDecWithPrec(1700, 0), // 1700
		P0ArithmeticTwapAccumulator: math.LegacyNewDecWithPrec(850, 0),
		P1LastSpotPrice:             math.LegacyNewDecWithPrec(170, 0),
		P0LastSpotPrice:             math.LegacyNewDecWithPrec(85, 0),
		GeometricTwapAccumulator:    math.LegacyNewDecWithPrec(1275, 0),
	}
	twapRecordAny2, err := cdctypes.NewAnyWithValue(&twapRecord2)
	require.NoError(t, err)

	// Set up comparison with difference_mode = true
	// This will compute the delta between intervals and compare
	flow.Conditions = &types.ExecutionConditions{
		Comparisons: []*types.Comparison{
			{
				ResponseIndex:  0,
				ResponseKey:    "",
				ValueType:      "osmosistwapv1beta1.TwapRecord.P1ArithmeticTwapAccumulator",
				DifferenceMode: true,
				ICQConfig: &types.ICQConfig{
					Response: twapRecordAny2.Value,
				},
				Operand:  "delta", // Dynamic check: compares Δ₂ (curr) vs Δ₁ (prev)
				Operator: types.ComparisonOperator_LARGER_THAN,
			},
		},
	}

	// Create a message to update
	msgDelegate := newFakeMsgDelegate(delAddr, valAddr)
	msgDelegate.Amount = sdk.NewCoin(testDenom, math.NewInt(1000))
	flow.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgDelegate})

	// Run comparison
	// Expected: Δaccumulator₂ (400) > 300 → true (positive trend detected)
	isTrue, err := keeper.CompareResponseValue(ctx, flow.ID, flow.Msgs, *flow.Conditions.Comparisons[0])
	require.NoError(t, err)
	require.True(t, isTrue, "Expected positive trend: Δaccumulator₂ (400) > Δaccumulator₁ (300)")

	// Test negative case: Interval 3 with deceleration
	// Accumulator = 2000
	// Δaccumulator₃ = 2000 - 1700 = 300
	// Since Δaccumulator₃ (300) is NOT > Δaccumulator₂ (400), this is deceleration
	twapRecord3 := osmosistwapv1beta1.TwapRecord{
		PoolId:                      1,
		P1ArithmeticTwapAccumulator: math.LegacyNewDecWithPrec(2000, 0), // 2000
		P0ArithmeticTwapAccumulator: math.LegacyNewDecWithPrec(1000, 0),
		P1LastSpotPrice:             math.LegacyNewDecWithPrec(200, 0),
		P0LastSpotPrice:             math.LegacyNewDecWithPrec(100, 0),
		GeometricTwapAccumulator:    math.LegacyNewDecWithPrec(1500, 0),
	}
	twapRecordAny3, err := cdctypes.NewAnyWithValue(&twapRecord3)
	require.NoError(t, err)

	keeper.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{
		QueryResponses: []string{base64.StdEncoding.EncodeToString(twapRecordAny3.Value)},
		Executed:       true,
	})

	// Update comparison to check against previous delta (400)
	flow.Conditions.Comparisons[0].ICQConfig.Response = twapRecordAny3.Value
	flow.Conditions.Comparisons[0].Operand = "delta" // Dynamic check: compares Δ₃ (curr) vs Δ₂ (prev)

	// Run comparison again
	// Expected: Δaccumulator₃ (300) > 400 → false (deceleration, not a positive trend)
	isTrue, err = keeper.CompareResponseValue(ctx, flow.ID, flow.Msgs, *flow.Conditions.Comparisons[0])
	require.NoError(t, err)
	require.False(t, isTrue, "Expected no positive trend: Δaccumulator₃ (300) is not > Δaccumulator₂ (400)")
}
