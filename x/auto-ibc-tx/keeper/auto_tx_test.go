package keeper

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

func TestCreateAutoTx(t *testing.T) {
	// Create a mock context and keeper
	ctx, keepers := CreateTestInput(t, false)
	types.Denom = sdk.DefaultBondDenom
	// Create a mock owner and fee funds
	owner, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	sendTo, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))
	feeFunds := sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 100))

	// Create a mock label, port ID, and messages
	label := "test-label"
	portID := "test-port-id"

	localMsg := newFakeMsgSend(owner, sendTo)
	msgs, err := types.PackTxMsgAnys([]sdk.Msg{localMsg})
	require.NoError(t, err)

	// Create a mock connection ID, duration, interval, start time, and dependencies
	connectionID := "test-connection-id"
	duration := 10 * time.Minute
	interval := 1 * time.Minute
	startTime := time.Now().UTC()
	dependsOn := []uint64{1, 2, 3}

	// Call the CreateAutoTx function
	err = keepers.AutoIbcTxKeeper.CreateAutoTx(ctx, owner, label, portID, msgs, connectionID, duration, interval, startTime, feeFunds, dependsOn)
	require.NoError(t, err)

	// Verify that the auto transaction was created correctly
	autoTx := keepers.AutoIbcTxKeeper.GetAutoTxInfo(ctx, 1)

	require.Equal(t, uint64(1), autoTx.TxID)
	require.Equal(t, owner.String(), autoTx.Owner)
	require.Equal(t, label, autoTx.Label)
	addr, _ := sdk.AccAddressFromBech32(autoTx.FeeAddress)
	require.Equal(t, feeFunds, keepers.BankKeeper.GetAllBalances(ctx, addr))
	require.Equal(t, interval, autoTx.Interval)
	require.Equal(t, startTime, autoTx.StartTime)
	require.Equal(t, portID, autoTx.PortID)
	require.Equal(t, connectionID, autoTx.ConnectionID)
	require.Equal(t, dependsOn, autoTx.DependsOnTxIds)
}

func TestCreateAutoTxWithZeroFundsWorks(t *testing.T) {
	// Create a mock context and keeper
	ctx, keepers := CreateTestInput(t, false)
	types.Denom = sdk.DefaultBondDenom
	// Create a mock owner and fee funds
	owner := sdk.AccAddress("owner")
	feeFunds := sdk.Coins{}
	sendTo, _ := CreateFakeFundedAccount(ctx, keepers.AccountKeeper, keepers.BankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000_000)))

	// Create a mock label, port ID, and messages
	label := "test-label"
	portID := "test-port-id"

	localMsg := newFakeMsgSend(owner, sendTo)
	msgs, err := types.PackTxMsgAnys([]sdk.Msg{localMsg})
	require.NoError(t, err)

	// Create a mock connection ID, duration, interval, start time, and dependencies
	connectionID := "test-connection-id"
	duration := 10 * time.Minute
	interval := 1 * time.Minute
	startTime := time.Now().UTC()
	dependsOn := []uint64{1, 2, 3}

	// Call the CreateAutoTx function
	err = keepers.AutoIbcTxKeeper.CreateAutoTx(ctx, owner, label, portID, msgs, connectionID, duration, interval, startTime, feeFunds, dependsOn)
	require.NoError(t, err)

	// Verify that the auto transaction was created correctly
	autoTx := keepers.AutoIbcTxKeeper.GetAutoTxInfo(ctx, 1)

	require.Equal(t, uint64(1), autoTx.TxID)
	require.Equal(t, owner.String(), autoTx.Owner)
	require.Equal(t, label, autoTx.Label)
	addr, _ := sdk.AccAddressFromBech32(autoTx.FeeAddress)
	require.Equal(t, sdk.Coins{}, keepers.BankKeeper.GetAllBalances(ctx, addr))
	require.Equal(t, interval, autoTx.Interval)
	require.Equal(t, startTime, autoTx.StartTime)
	require.Equal(t, portID, autoTx.PortID)
	require.Equal(t, connectionID, autoTx.ConnectionID)
	require.Equal(t, dependsOn, autoTx.DependsOnTxIds)
}
