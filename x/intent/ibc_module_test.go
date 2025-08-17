package intent

import (
	"encoding/base64"
	"fmt"
	"testing"

	"cosmossdk.io/math"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/stretchr/testify/require"
	keeper "github.com/trstlabs/intento/x/intent/keeper"
	elysamm "github.com/trstlabs/intento/x/intent/msg_registry/elys/amm"
	elyscommitment "github.com/trstlabs/intento/x/intent/msg_registry/elys/commitment"
	elysestaking "github.com/trstlabs/intento/x/intent/msg_registry/elys/estaking"
	"github.com/trstlabs/intento/x/intent/types"
)

func TestHandleWithdrawElysStakingRewardsPacket(t *testing.T) {
	// Setup
	ctx, keepers, _ := createTestContext(t)
	k := keepers.IntentKeeper

	flow, _ := createTestFlow(ctx, types.ExecutionConfiguration{}, keepers)

	eMsg := elysestaking.MsgWithdrawElysStakingRewards{}
	any, _ := cdctypes.NewAnyWithValue(&eMsg)
	msg := authztypes.MsgExec{Msgs: []*cdctypes.Any{any}}
	anys, _ := types.PackTxMsgAnys([]sdk.Msg{&msg})
	flow.Msgs = anys
	k.SetFlow(ctx, &flow)
	relayer, _ := keeper.CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))

	k.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: nil})
	k.SetTmpFlowID(ctx, flow.ID, "icacontroller-into12m09f4a8jeam4ysm7udq6449qf49grklr2c50xs3hzkuryh0znmqyql2u9", "channel-1", 0)
	// Create a mock packet
	packetData := channeltypes.Packet{
		Data:               []byte{},
		SourcePort:         "icacontroller-into12m09f4a8jeam4ysm7udq6449qf49grklr2c50xs3hzkuryh0znmqyql2u9",
		SourceChannel:      "channel-1",
		DestinationPort:    "icahost",
		DestinationChannel: "channel-98",
		Sequence:           0,
	}
	// Create a mock acknowledgement
	ackBase64 := "eyJyZXN1bHQiOiJFcElCQ2lVdlkyOXpiVzl6TG1GMWRHaDZMbll4WW1WMFlURXVUWE5uUlhobFkxSmxjM0J2Ym5ObEVta0tad3BKQ2tScFltTXZSakE0TWtJMk5VTTRPRVUwUWpaRU5VVkdNVVJDTWpRelEwUkJNVVF6TXpGRU1EQXlOelU1UlRrek9FRXdSalZEUkROR1JrUkROVVExTTBJelJUTTBPUklCTXdvTENnVjFaV1JsYmhJQ09Ea0tEUW9HZFdWa1pXNWlFZ014T0RnPSJ9"
	ackBytes, err := base64.StdEncoding.DecodeString(ackBase64)
	require.NoError(t, err)

	IBCModule := NewIBCModule(keepers.IntentKeeper)
	// Call the OnAcknowledgementPacket function
	err = IBCModule.OnAcknowledgementPacket(
		ctx,
		packetData,
		ackBytes,
		relayer,
	)

	// Verify no error occurred
	require.NoError(t, err)
	flowHistory, err := keepers.IntentKeeper.GetFlowHistory(ctx, flow.ID)
	require.NoError(t, err)

	require.Nil(t, flowHistory[len(flowHistory)-1].Errors)

}

func TestFeedbackLoopWithdrawSwapStake(t *testing.T) {
	ctx, keepers, _ := createTestContext(t)
	k := keepers.IntentKeeper

	// Create a basic flow with just one withdraw->swap->stake sequence
	flow, _ := createTestFlow(ctx, types.ExecutionConfiguration{}, keepers)
	flow.Configuration.SaveResponses = true
	// Create messages
	withdrawMsg := &elysestaking.MsgWithdrawElysStakingRewards{}
	withdrawAny, _ := cdctypes.NewAnyWithValue(withdrawMsg)

	swapMsg := &elysamm.MsgSwapExactAmountIn{
		TokenIn: sdk.Coin{Denom: "tokenA", Amount: math.NewInt(0)},
	}
	swapAny, _ := cdctypes.NewAnyWithValue(swapMsg)

	stakeMsg := &elyscommitment.MsgStake{}
	stakeAny, _ := cdctypes.NewAnyWithValue(stakeMsg)

	// Wrap withdraw in MsgExec
	exec := &authztypes.MsgExec{Msgs: []*cdctypes.Any{withdrawAny}}
	execAny, _ := cdctypes.NewAnyWithValue(exec)

	flow.Msgs = []*cdctypes.Any{execAny, swapAny, stakeAny}

	// Setup feedback loops
	flow.Conditions = &types.ExecutionConditions{
		FeedbackLoops: []*types.FeedbackLoop{
			// Update swap amount from withdraw response
			{
				ResponseIndex: 0,
				ResponseKey:   "Amount.[0].Amount",
				MsgsIndex:     1,
				MsgKey:        "TokenIn.Amount",
				ValueType:     "sdk.Int",
			},
			// Update stake amount from swap response
			{
				ResponseIndex: 1,
				ResponseKey:   "TokenOutAmount",
				MsgsIndex:     2,
				MsgKey:        "Amount",
				ValueType:     "sdk.Int",
			},
		},
	}

	k.SetFlow(ctx, &flow)

	// Simulate withdraw response
	withdrawResp := &elysestaking.MsgWithdrawElysStakingRewardsResponse{
		Amount: sdk.Coins{{Amount: math.NewInt(100)}},
	}
	withdrawRespAny, _ := cdctypes.NewAnyWithValue(withdrawResp)

	historyEntry := &types.FlowHistoryEntry{
		MsgResponses: []*cdctypes.Any{},
		ExecFee:      sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100))),
	}
	k.SetCurrentFlowHistoryEntry(ctx, flow.ID, historyEntry)

	// Set up tmp flow ID
	portID := "icacontroller-test"
	channelID := "channel-1"
	seq := uint64(1)
	k.SetTmpFlowID(ctx, flow.ID, portID, channelID, seq)
	fmt.Print("fllow]n", flow.Conditions)
	// Process response
	err := k.HandleResponseAndSetFlowResult(ctx, portID, channelID, make(sdk.AccAddress, 20), seq, []*cdctypes.Any{withdrawRespAny})
	require.NoError(t, err)

	// Verify swap message was updated
	updatedFlow := k.GetFlow(ctx, flow.ID)
	var updatedSwapMsg elysamm.MsgSwapExactAmountIn
	updatedSwapMsg.Unmarshal(updatedFlow.Msgs[1].Value)

	require.NotEqual(t, updatedFlow.Msgs[1].Value, swapAny.Value)
	require.Equal(t, math.NewInt(100), updatedSwapMsg.TokenIn.Amount)
}

func TestFeedbackLoopWithdrawSwapStakeTwice(t *testing.T) {
	// 	Processes two complete withdraw→swap→stake sequences
	// Verifies each feedback loop updates the correct message
	// Maintains proper response indexing across multiple sequences
	// Checks the final state after all operations

	ctx, keepers, _ := createTestContext(t)
	k := keepers.IntentKeeper

	// Create a flow with two withdraw->swap->stake sequences
	flow, _ := createTestFlow(ctx, types.ExecutionConfiguration{}, keepers)
	flow.Configuration.SaveResponses = true

	// Create messages for first sequence
	withdrawMsg1 := &elysestaking.MsgWithdrawElysStakingRewards{}
	withdrawAny1, _ := cdctypes.NewAnyWithValue(withdrawMsg1)
	swapMsg1 := &elysamm.MsgSwapExactAmountIn{
		TokenIn: sdk.Coin{Denom: "tokenA", Amount: math.NewInt(0)},
	}
	swapAny1, _ := cdctypes.NewAnyWithValue(swapMsg1)
	stakeMsg1 := &elyscommitment.MsgStake{}
	stakeAny1, _ := cdctypes.NewAnyWithValue(stakeMsg1)

	// Create messages for second sequence
	withdrawMsg2 := &elysestaking.MsgWithdrawElysStakingRewards{}
	withdrawAny2, _ := cdctypes.NewAnyWithValue(withdrawMsg2)
	swapMsg2 := &elysamm.MsgSwapExactAmountIn{
		TokenIn: sdk.Coin{Denom: "tokenA", Amount: math.NewInt(0)},
	}
	swapAny2, _ := cdctypes.NewAnyWithValue(swapMsg2)
	stakeMsg2 := &elyscommitment.MsgStake{}
	stakeAny2, _ := cdctypes.NewAnyWithValue(stakeMsg2)

	// Wrap withdraws in MsgExec
	exec1 := &authztypes.MsgExec{Msgs: []*cdctypes.Any{withdrawAny1}}
	execAny1, _ := cdctypes.NewAnyWithValue(exec1)
	exec2 := &authztypes.MsgExec{Msgs: []*cdctypes.Any{withdrawAny2}}
	execAny2, _ := cdctypes.NewAnyWithValue(exec2)

	// Set up flow with both sequences
	flow.Msgs = []*cdctypes.Any{
		execAny1, swapAny1, stakeAny1, // First sequence
		execAny2, swapAny2, stakeAny2, // Second sequence
	}

	// Setup feedback loops for both sequences
	flow.Conditions = &types.ExecutionConditions{
		FeedbackLoops: []*types.FeedbackLoop{
			// First sequence
			{
				ResponseIndex: 0, // First withdraw response
				ResponseKey:   "Amount.[0].Amount",
				MsgsIndex:     1, // First swap
				MsgKey:        "TokenIn.Amount",
				ValueType:     "sdk.Int",
			},
			{
				ResponseIndex: 1, // First swap response
				ResponseKey:   "TokenOutAmount",
				MsgsIndex:     2, // First stake
				MsgKey:        "Amount",
				ValueType:     "sdk.Int",
			},
			// Second sequence
			{
				ResponseIndex: 2, // Second withdraw response
				ResponseKey:   "Amount.[0].Amount",
				MsgsIndex:     4, // Second swap
				MsgKey:        "TokenIn.Amount",
				ValueType:     "sdk.Int",
			},
			{
				ResponseIndex: 3, // Second swap response
				ResponseKey:   "TokenOutAmount",
				MsgsIndex:     5, // Second stake
				MsgKey:        "Amount",
				ValueType:     "sdk.Int",
			},
		},
	}

	k.SetFlow(ctx, &flow)

	// Simulate first withdraw response (100 tokens)
	withdrawResp1 := &elysestaking.MsgWithdrawElysStakingRewardsResponse{
		Amount: sdk.Coins{{Amount: math.NewInt(100)}},
	}
	withdrawRespAny1, _ := cdctypes.NewAnyWithValue(withdrawResp1)

	// Set up tmp flow ID for first withdraw
	portID := "icacontroller-test"
	channelID := "channel-1"
	seq1 := uint64(1)
	k.SetTmpFlowID(ctx, flow.ID, portID, channelID, seq1)

	historyEntry := &types.FlowHistoryEntry{
		MsgResponses: []*cdctypes.Any{},
		ExecFee:      sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100))),
	}
	k.SetCurrentFlowHistoryEntry(ctx, flow.ID, historyEntry)

	// Process first withdraw response
	err := k.HandleResponseAndSetFlowResult(ctx, portID, channelID, make(sdk.AccAddress, 20), seq1, []*cdctypes.Any{withdrawRespAny1})
	require.NoError(t, err)

	// Verify first swap message was updated
	updatedFlow := k.GetFlow(ctx, flow.ID)

	var updatedSwapMsg1 elysamm.MsgSwapExactAmountIn
	updatedSwapMsg1.Unmarshal(updatedFlow.Msgs[1].Value)
	require.Equal(t, math.NewInt(100), updatedSwapMsg1.TokenIn.Amount)

	// Simulate first swap response (80 tokens out)
	swapResp1 := &elysamm.MsgSwapExactAmountInResponse{
		TokenOutAmount: math.NewInt(80),
	}
	swapRespAny1, _ := cdctypes.NewAnyWithValue(swapResp1)

	// Set up tmp flow ID for first swap
	seq2 := uint64(2)
	k.SetTmpFlowID(ctx, flow.ID, portID, channelID, seq2)

	// Process first swap response
	err = k.HandleResponseAndSetFlowResult(ctx, portID, channelID, make(sdk.AccAddress, 20), seq2, []*cdctypes.Any{swapRespAny1})
	require.NoError(t, err)

	// Verify first stake message was updated
	updatedFlow = k.GetFlow(ctx, flow.ID)
	var updatedStakeMsg1 elyscommitment.MsgStake
	require.NoError(t, updatedStakeMsg1.Unmarshal(updatedFlow.Msgs[2].Value))
	require.Equal(t, math.NewInt(80), updatedStakeMsg1.Amount)

	// Simulate second withdraw response (200 tokens)
	withdrawResp2 := &elysestaking.MsgWithdrawElysStakingRewardsResponse{
		Amount: sdk.Coins{{Amount: math.NewInt(200)}},
	}
	withdrawRespAny2, _ := cdctypes.NewAnyWithValue(withdrawResp2)

	// Set up tmp flow ID for second withdraw
	seq3 := uint64(3)
	k.SetTmpFlowID(ctx, flow.ID, portID, channelID, seq3)

	// Process second withdraw response
	err = k.HandleResponseAndSetFlowResult(ctx, portID, channelID, make(sdk.AccAddress, 20), seq3, []*cdctypes.Any{withdrawRespAny2})
	require.NoError(t, err)

	// Verify second swap message was updated
	updatedFlow = k.GetFlow(ctx, flow.ID)
	var updatedSwapMsg2 elysamm.MsgSwapExactAmountIn
	require.NoError(t, updatedSwapMsg2.Unmarshal(updatedFlow.Msgs[4].Value))
	require.Equal(t, math.NewInt(200), updatedSwapMsg2.TokenIn.Amount)

	// Simulate second swap response (180 tokens out)
	swapResp2 := &elysamm.MsgSwapExactAmountInResponse{
		TokenOutAmount: math.NewInt(180),
	}
	swapRespAny2, _ := cdctypes.NewAnyWithValue(swapResp2)

	// Set up tmp flow ID for second swap
	seq4 := uint64(4)
	k.SetTmpFlowID(ctx, flow.ID, portID, channelID, seq4)

	// Process second swap response
	err = k.HandleResponseAndSetFlowResult(ctx, portID, channelID, make(sdk.AccAddress, 20), seq4, []*cdctypes.Any{swapRespAny2})
	require.NoError(t, err)

	// Verify second stake message was updated
	updatedFlow = k.GetFlow(ctx, flow.ID)
	var updatedStakeMsg2 elyscommitment.MsgStake
	require.NoError(t, updatedStakeMsg2.Unmarshal(updatedFlow.Msgs[5].Value))
	require.Equal(t, math.NewInt(180), updatedStakeMsg2.Amount)
}
