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
		Sender:           "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs",
		Receiver:         "",
		TimeoutHeight:    clienttypes.Height{RevisionNumber: 0, RevisionHeight: 0},
		TimeoutTimestamp: 2526374086000000000,
		Memo:             "",
	}

	anyMsg, _ := PackTxMsgAnys([]sdk.Msg{&msg})

	// Call the function under test
	_, err := GetTransferMsg(cdc, anyMsg[0])

	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid coins")

}
func TestGetValidTransferMsgDoesNotError(t *testing.T) {
	// Set up the codec
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	ibctransfertypes.RegisterInterfaces(interfaceRegistry)

	cdc := codec.NewProtoCodec(interfaceRegistry)

	var codec codec.Codec = cdc

	// Create a sample MsgTransfer
	msg := ibctransfertypes.MsgTransfer{
		SourcePort:       "transfer",
		SourceChannel:    "channel-6",
		Token:            sdk.Coin{Denom: "uatom", Amount: math.NewInt(400)},
		Sender:           "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs",
		Receiver:         "into1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx",
		TimeoutHeight:    clienttypes.Height{RevisionNumber: 0, RevisionHeight: 0},
		TimeoutTimestamp: 2526374086000000000,
		Memo:             "",
	}

	anyMsg, _ := PackTxMsgAnys([]sdk.Msg{&msg})

	// Call the function under test
	msg, err := GetTransferMsg(codec, anyMsg[0])
	require.NoError(t, err)

	require.Equal(t, msg, ibctransfertypes.MsgTransfer{
		SourcePort:       "transfer",
		SourceChannel:    "channel-6",
		Token:            sdk.Coin{Denom: "uatom", Amount: math.NewInt(400)},
		Sender:           "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs",
		Receiver:         "into1ykql5ktedxkpjszj5trzu8f5dxajvgv95nuwjx",
		TimeoutHeight:    clienttypes.Height{RevisionNumber: 0, RevisionHeight: 0},
		TimeoutTimestamp: 2526374086000000000,
		Memo:             "",
	})
}
