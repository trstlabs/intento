package keeper

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"
	elysamm "github.com/trstlabs/intento/x/intent/msg_registry/elys/amm"
	elyscommitment "github.com/trstlabs/intento/x/intent/msg_registry/elys/commitment"
	elysestaking "github.com/trstlabs/intento/x/intent/msg_registry/elys/estaking"
	"github.com/trstlabs/intento/x/intent/types"
)

func TestCreateFlow(t *testing.T) {
	// Create a mock context and keeper
	ctx, keepers, _ := CreateTestInput(t, false)
	types.Denom = sdk.DefaultBondDenom
	// Create a mock owner and fee funds
	owner, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	sendTo, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	feeFunds := sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100))

	// Create a mock label, port ID, and messages
	label := "test-label"

	localMsg := newFakeMsgSend(owner, sendTo)
	msgs, err := types.PackTxMsgAnys([]sdk.Msg{localMsg})
	require.NoError(t, err)

	duration := 10 * time.Minute
	interval := 1 * time.Minute
	startTime := time.Now().UTC()
	configuration := types.ExecutionConfiguration{SaveResponses: false}

	// Call the CreateFlow function
	err = keepers.IntentKeeper.CreateFlow(ctx, owner, label, msgs, duration, interval, startTime, feeFunds, configuration, types.TrustlessAgentConfig{}, "", "", types.ExecutionConditions{})
	require.NoError(t, err)

	// Verify that the flow was created correctly
	flow := keepers.IntentKeeper.Getflow(ctx, 1)

	require.Equal(t, uint64(1), flow.ID)
	require.Equal(t, owner.String(), flow.Owner)
	require.Equal(t, label, flow.Label)
	addr, _ := sdk.AccAddressFromBech32(flow.FeeAddress)
	require.Equal(t, feeFunds, keepers.BankKeeper.GetAllBalances(ctx, addr))
	require.Equal(t, interval, flow.Interval)
	require.Equal(t, startTime, flow.StartTime)
	require.Equal(t, configuration, *flow.Configuration)
}

func TestCreateFlowWithZeroFeeFundsWorks(t *testing.T) {
	// Create a mock context and keeper
	ctx, keepers, _ := CreateTestInput(t, false)
	types.Denom = sdk.DefaultBondDenom
	// Create a mock owner and fee funds
	owner := sdk.AccAddress("owner")
	feeFunds := sdk.Coins{}
	sendTo, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))

	// Create a mock label, port ID, and messages
	label := "test-label"
	localMsg := newFakeMsgSend(owner, sendTo)
	msgs, err := types.PackTxMsgAnys([]sdk.Msg{localMsg})
	require.NoError(t, err)

	duration := 10 * time.Minute
	interval := 1 * time.Minute
	startTime := time.Now().UTC()
	configuration := types.ExecutionConfiguration{SaveResponses: false}

	// Call the CreateFlow function
	err = keepers.IntentKeeper.CreateFlow(ctx, owner, label, msgs, duration, interval, startTime, feeFunds, configuration, types.TrustlessAgentConfig{}, "", "", types.ExecutionConditions{})
	require.NoError(t, err)

	// Verify that the flow was created correctly
	flow := keepers.IntentKeeper.Getflow(ctx, 1)

	require.Equal(t, uint64(1), flow.ID)
	require.Equal(t, owner.String(), flow.Owner)
	require.Equal(t, label, flow.Label)
	addr, _ := sdk.AccAddressFromBech32(flow.FeeAddress)
	require.Equal(t, sdk.Coins{}, keepers.BankKeeper.GetAllBalances(ctx, addr))
	require.Equal(t, interval, flow.Interval)
	require.Equal(t, startTime, flow.StartTime)
	require.Equal(t, configuration, *flow.Configuration)
}

func TestGetFlowsForBlock(t *testing.T) {
	// Create a mock context and keeper
	ctx, keepers, _ := CreateTestInput(t, false)
	types.Denom = sdk.DefaultBondDenom
	// Create a mock owner and fee funds
	owner, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	sendTo, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	feeFunds := sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100))

	// Create a mock label, port ID, and messages
	label := "test-label"

	localMsg := newFakeMsgSend(owner, sendTo)
	msgs, err := types.PackTxMsgAnys([]sdk.Msg{localMsg})
	require.NoError(t, err)

	duration := 10 * time.Minute
	interval := 1 * time.Minute
	startTime := time.Now().UTC()
	configuration := types.ExecutionConfiguration{SaveResponses: false}

	// Call the CreateFlow function
	err = keepers.IntentKeeper.CreateFlow(ctx, owner, label, msgs, duration, interval, startTime, feeFunds, configuration, types.TrustlessAgentConfig{}, "", "", types.ExecutionConditions{})
	require.NoError(t, err)
	// Call the CreateFlow function
	err = keepers.IntentKeeper.CreateFlow(ctx, owner, label, msgs, duration, interval, startTime, feeFunds, configuration, types.TrustlessAgentConfig{}, "", "", types.ExecutionConditions{})
	require.NoError(t, err)
	flows := keepers.IntentKeeper.GetFlowsForBlock(ctx.WithBlockTime(startTime.Add(interval)))
	require.Equal(t, len(flows), 2)

}

func TestIncrementalExecutionWithFeedbackLoops(t *testing.T) {
	ctx, keepers, _ := CreateTestInput(t, false)
	k := keepers.IntentKeeper
	cdc := keepers.IntentKeeper.cdc

	// Create a mock label, port ID, and messages
	duration := 10 * time.Minute
	interval := 1 * time.Minute
	startTime := time.Now().UTC()
	configuration := types.ExecutionConfiguration{SaveResponses: true, WalletFallback: true}
	owner, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	types.Denom = "stake"

	portID := "icacontroller-test"
	channelID := "channel-1"
	connectionID := "connection-1"

	// Call the CreateFlow function
	err := keepers.IntentKeeper.CreateFlow(ctx, owner, "label", []*cdctypes.Any{}, duration, interval, startTime, sdk.Coins{}, configuration, types.TrustlessAgentConfig{}, portID, connectionID, types.ExecutionConditions{})
	require.NoError(t, err)
	flow := keepers.IntentKeeper.Getflow(ctx, 1)
	require.NotNil(t, flow.FeeAddress)
	// Message creation helpers
	newWithdrawMsg := func() *cdctypes.Any {
		msg := &elysestaking.MsgWithdrawElysStakingRewards{}
		any, _ := cdctypes.NewAnyWithValue(msg)
		return any
	}
	newSwapMsg := func() *cdctypes.Any {
		msg := &elysamm.MsgSwapExactAmountIn{TokenIn: sdk.Coin{Denom: "tokenA", Amount: math.NewInt(0)}}
		any, _ := cdctypes.NewAnyWithValue(msg)
		return any
	}
	newStakeMsg := func() *cdctypes.Any {
		msg := &elyscommitment.MsgStake{}
		any, _ := cdctypes.NewAnyWithValue(msg)
		return any
	}

	// Build flow with two sequences
	flow.Msgs = []*cdctypes.Any{
		newWithdrawMsg(), // Index 0 - First withdraw (needs response)
		newSwapMsg(),     // Index 1 - First swap (updated by feedback, needs response)
		newStakeMsg(),    // Index 2 - First stake (updated by feedback)
		newWithdrawMsg(), // Index 3 - Second withdraw (needs response)
		newSwapMsg(),     // Index 4 - Second swap (updated by feedback, needs response)
		newStakeMsg(),    // Index 5 - Second stake (updated by feedback)
	}

	k.SetFlowHistoryEntry(ctx, flow.ID, &types.FlowHistoryEntry{MsgResponses: nil, ExecFee: sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(100)))})

	// Setup feedback loops
	flow.Conditions = &types.ExecutionConditions{
		FeedbackLoops: []*types.FeedbackLoop{
			// First sequence
			{
				ResponseIndex: 0, // Withdraw#0
				ResponseKey:   "Amount.[0].Amount",
				MsgsIndex:     1, // Swap#1
				MsgKey:        "TokenIn.Amount",
				ValueType:     "sdk.Int",
			},
			{
				ResponseIndex: 1, // Swap#1
				ResponseKey:   "TokenOutAmount",
				MsgsIndex:     2, // Stake#2
				MsgKey:        "Amount",
				ValueType:     "sdk.Int",
			},
			// Second sequence
			{
				ResponseIndex: 3, // Withdraw#3
				ResponseKey:   "Amount.[0].Amount",
				MsgsIndex:     4, // Swap#4
				MsgKey:        "TokenIn.Amount",
				ValueType:     "sdk.Int",
			},
			{
				ResponseIndex: 4, // Swap#4
				ResponseKey:   "TokenOutAmount",
				MsgsIndex:     5, // Stake#5
				MsgKey:        "Amount",
				ValueType:     "sdk.Int",
			},
		},
	}

	k.Setflow(ctx, &flow)

	// Test cases in execution order
	tests := []struct {
		name          string
		seq           uint64
		createResp    func() proto.Message
		expectedMsgs  []int // Expected message indices in history
		validateIndex int   // Which message to check was updated
	}{
		{
			name: "Withdraw response #0",
			seq:  1,
			createResp: func() proto.Message {
				return &elysestaking.MsgWithdrawElysStakingRewardsResponse{
					Amount: sdk.Coins{{Amount: math.NewInt(100)}},
				}
			},
			expectedMsgs:  []int{0}, // withdraw 0 is done
			validateIndex: 1,        // Check swap was updated
		},
		{
			name: "Swap response #1",
			seq:  2,
			createResp: func() proto.Message {
				return &elysamm.MsgSwapExactAmountInResponse{
					TokenOutAmount: math.NewInt(80),
				}
			},
			expectedMsgs:  []int{0, 1}, // swap 1 is now done
			validateIndex: 2,
		},
		{
			name: "Stake response #2",
			seq:  3,
			createResp: func() proto.Message {
				return &elyscommitment.MsgStakeResponse{}
			},
			expectedMsgs:  []int{0, 1, 2}, // stake 2 is now done
			validateIndex: 4,              // swap 4 was waiting on stake 2
		},
		{
			name: "Withdraw response #3",
			seq:  4,
			createResp: func() proto.Message {
				return &elysestaking.MsgWithdrawElysStakingRewardsResponse{

					Amount: sdk.Coins{{Amount: math.NewInt(200)}},
				}
			},
			expectedMsgs:  []int{0, 1, 2, 3}, // withdraw 3 done
			validateIndex: 5,                 // stake 5 waiting on this
		},
		{
			name: "Swap response #4",
			seq:  5,
			createResp: func() proto.Message {
				return &elysamm.MsgSwapExactAmountInResponse{}
			},
			expectedMsgs:  []int{0, 1, 2, 3, 4}, // swap 4 done
			validateIndex: -1,
		},
		{
			name: "Stake response #5",
			seq:  6,
			createResp: func() proto.Message {
				return &elyscommitment.MsgStakeResponse{}
			},
			expectedMsgs:  []int{0, 1, 2, 3, 4, 5}, // all done
			validateIndex: -1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set sequence mapping
			k.SetTmpFlowID(ctx, flow.ID, portID, channelID, tc.seq)

			// Create response
			respAny, err := cdctypes.NewAnyWithValue(tc.createResp())
			require.NoError(t, err)

			// Process response
			err = k.HandleResponseAndSetFlowResult(ctx, portID, channelID, make(sdk.AccAddress, 20), tc.seq, []*cdctypes.Any{respAny})
			require.NoError(t, err)

			// Verify execution history
			history := k.MustGetFlowHistory(ctx, flow.ID)
			require.Len(t, history[0].MsgResponses, len(tc.expectedMsgs), "unexpected number of executed messages")

			for _, msgIdx := range tc.expectedMsgs {
				require.Contains(t, history[0].MsgResponses[msgIdx].TypeUrl, flow.Msgs[msgIdx].TypeUrl)
			}

			// Verify target message was updated
			updatedFlow := k.Getflow(ctx, flow.ID)
			var updatedMsg proto.Message
			if tc.validateIndex >= 0 {
				var updatedMsg proto.Message
				require.NoError(t, cdc.UnpackAny(updatedFlow.Msgs[tc.validateIndex], &updatedMsg))
			}

			switch msg := updatedMsg.(type) {
			case *elysamm.MsgSwapExactAmountIn:
				expected, ok := tc.createResp().(*elysestaking.MsgWithdrawElysStakingRewardsResponse)
				require.True(t, ok, "expected WithdrawElysStakingRewardsResponse")
				require.Equal(t, expected.Amount[0].Amount, msg.TokenIn.Amount)

			case *elyscommitment.MsgStake:
				expected, ok := tc.createResp().(*elysamm.MsgSwapExactAmountInResponse)
				require.True(t, ok, "expected MsgSwapExactAmountInResponse")
				require.Equal(t, expected.TokenOutAmount, msg.Amount)
			}

		})
	}

	// Final verification
	finalHistory := k.MustGetFlowHistory(ctx, flow.ID)
	require.Len(t, finalHistory[0].MsgResponses, 6, "all messages should be executed")

}
