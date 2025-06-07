package intent

import (
	"encoding/base64"
	"testing"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/stretchr/testify/require"
	keeper "github.com/trstlabs/intento/x/intent/keeper"
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
	k.SetFlowInfo(ctx, &flow)
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
