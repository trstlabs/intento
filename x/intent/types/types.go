package types

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	proto "github.com/cosmos/gogoproto/proto"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
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
func (actionInfo ActionInfo) ActionAuthzSignerOk(codec codec.Codec) bool {
	for _, message := range actionInfo.Msgs {
		var sdkMsg sdk.Msg
		err := codec.UnpackAny(message, &sdkMsg)
		if err != nil {
			return false
		}

		// fmt.Printf("signer: %v owner %v \n", sdkMsg.GetSigners()[0].String(), actionInfo.Owner)
		if ((message.TypeUrl) == sdk.MsgTypeURL(&authztypes.MsgExec{})) {
			var authzMsg authztypes.MsgExec
			if err := proto.Unmarshal(message.Value, &authzMsg); err != nil {
				return false
			}

			for _, message := range authzMsg.Msgs {
				signers, _, err := codec.GetMsgV1Signers(message)
				if err != nil {
					return false
				}
				for _, acct := range signers {
					if sdk.AccAddress(acct).String() != actionInfo.Owner {
						return false
					}
				}
				var sdkMsgAuthZ sdk.Msg
				err = codec.UnpackAny(message, &sdkMsgAuthZ)
				if err != nil {
					return false
				}
				m, ok := sdkMsgAuthZ.(sdk.HasValidateBasic)
				if !ok {
					continue
				}

				if err := m.ValidateBasic(); err != nil {
					return false
				}
			}
		}
	}
	return true
}

var GasFeeCoinsSupported sdk.Coins = sdk.Coins{sdk.NewCoin(Denom, math.NewInt(10))}

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
