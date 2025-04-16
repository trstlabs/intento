package types

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/stretchr/testify/require"
)

func TestGetInvalidTransferMsgDoesNotPanic(t *testing.T) {
	// Set up the codec
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	ibctransfertypes.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)

	// Create a sample MsgTransfer
	msg := ibctransfertypes.MsgTransfer{
		SourcePort:       "transfer",
		SourceChannel:    "channel-6",
		Token:            sdk.Coin{Denom: "10", Amount: math.NewInt(400)},
		Sender:           "into18ajuj6drylfdvt4d37peexle47ljxqa5v8r6n8",
		Receiver:         "",
		TimeoutHeight:    clienttypes.Height{RevisionNumber: 0, RevisionHeight: 0},
		TimeoutTimestamp: 2526374086000000000,
		Memo:             "",
	}

	anyMsg, _ := PackTxMsgAnys([]sdk.Msg{&msg})

	// Call the function under test
	_, err := GetTransferMsg(cdc, anyMsg[0])
	require.Error(t, err)

}
