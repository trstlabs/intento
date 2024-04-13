package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/cosmos/gogoproto/proto"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
)

var Denom = "uinto"
var ParseICAValue = "ICA_ADDR"

// GetTxMsgs unpacks sdk messages from any messages
func (actionInfo ActionInfo) GetTxMsgs(unpacker types.AnyUnpacker) (sdkMsgs []sdk.Msg) {

	for _, message := range actionInfo.Msgs {
		var sdkMsg sdk.Msg
		err := unpacker.UnpackAny(message, &sdkMsg)
		if err != nil {
			return nil
		}
		sdkMsgs = append(sdkMsgs, sdkMsg)
	}
	return sdkMsgs
}

// GetTxMsgs unpacks sdk messages from any messages
func GetTransferMsg(cdc codec.Codec, anyTransfer *types.Any) (transferMsg ibctransfertypes.MsgTransfer, err error) {
	//		// where Y is a field on MyStruct that implements UnpackInterfacesMessage itself
	//		err = s.Y.UnpackInterfaces(unpacker)
	//		if err != nil {
	//			return nil
	//		}
	//		return nil
	//	 }
	if err := proto.Unmarshal(anyTransfer.Value, &transferMsg); err != nil {
		return ibctransfertypes.MsgTransfer{}, err
	}

	// var sdkMsg sdk.Msg
	// err = cdc.UnpackAny(anyTransfer, &sdkMsg)
	// if err != nil {
	// 	return ibctransfertypes.MsgTransfer{}, err
	// }
	// sdkMsg.ProtoMessage()
	return transferMsg, nil
}
