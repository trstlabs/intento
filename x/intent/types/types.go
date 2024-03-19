package types

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
)

var Denom = "uinto"
var ParseICAValue = "ICA_ADDR"

// GetTxMsgs unpacks sdk messages from any messages
func (autoTxInfo ActionInfo) GetTxMsgs(unpacker types.AnyUnpacker) (sdkMsgs []sdk.Msg) {

	for _, message := range autoTxInfo.Msgs {
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
func GetTransferMsg(unpacker types.AnyUnpacker, any *types.Any) (transferMsg ibctransfertypes.MsgTransfer, err error) {

	err = unpacker.UnpackAny(any, &transferMsg)
	if err != nil {
		return ibctransfertypes.MsgTransfer{}, err
	}

	return transferMsg, nil
}
