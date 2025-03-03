package keeper

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
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
	err = keepers.IntentKeeper.CreateFlow(ctx, owner, label, msgs, duration, interval, startTime, feeFunds, configuration, types.HostedICAConfig{}, "", "", types.ExecutionConditions{})
	require.NoError(t, err)

	// Verify that the flow was created correctly
	flow := keepers.IntentKeeper.GetFlowInfo(ctx, 1)

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
	err = keepers.IntentKeeper.CreateFlow(ctx, owner, label, msgs, duration, interval, startTime, feeFunds, configuration, types.HostedICAConfig{}, "", "", types.ExecutionConditions{})
	require.NoError(t, err)

	// Verify that the flow was created correctly
	flow := keepers.IntentKeeper.GetFlowInfo(ctx, 1)

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
	err = keepers.IntentKeeper.CreateFlow(ctx, owner, label, msgs, duration, interval, startTime, feeFunds, configuration, types.HostedICAConfig{}, "", "", types.ExecutionConditions{})
	require.NoError(t, err)
	// Call the CreateFlow function
	err = keepers.IntentKeeper.CreateFlow(ctx, owner, label, msgs, duration, interval, startTime, feeFunds, configuration, types.HostedICAConfig{}, "", "", types.ExecutionConditions{})
	require.NoError(t, err)
	flows := keepers.IntentKeeper.GetFlowsForBlock(ctx.WithBlockTime(startTime.Add(interval)))
	require.Equal(t, len(flows), 2)

}
