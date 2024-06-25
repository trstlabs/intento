package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
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

// grantee should always be action owner
func (actionInfo ActionInfo) ActionAuthzSignerOk(unpacker types.AnyUnpacker) bool {
	for _, message := range actionInfo.Msgs {
		var sdkMsg sdk.Msg
		err := unpacker.UnpackAny(message, &sdkMsg)
		if err != nil {
			return false
		}

		// fmt.Printf("signer: %v owner %v \n", sdkMsg.GetSigners()[0].String(), actionInfo.Owner)
		if sdkMsg.GetSigners()[0].String() != actionInfo.Owner && ((message.TypeUrl) == sdk.MsgTypeURL(&authztypes.MsgExec{})) {
			var authzMsg authztypes.MsgExec
			if err := proto.Unmarshal(message.Value, &authzMsg); err != nil {
				return false
			}
			for _, message := range authzMsg.Msgs {
				var sdkMsgAuthZ sdk.Msg
				err := unpacker.UnpackAny(message, &sdkMsgAuthZ)
				if err != nil {
					return false
				}
				//fmt.Printf("signer3: %v \n", sdkMsgAuthZ.GetSigners()[0].String())
				if sdkMsgAuthZ.GetSigners()[0].String() != "" && sdkMsgAuthZ.GetSigners()[0].String() != actionInfo.Owner {
					fmt.Printf("false: %v %v \n", sdkMsgAuthZ.GetSigners()[0].String(), actionInfo.Owner)
					// fmt.Printf("sdkMsg: %v \n", sdkMsgAuthZ)
					return false
				}
			}
		}
	}
	return true
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
