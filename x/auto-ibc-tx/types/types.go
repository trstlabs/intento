package types

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var Denom = "utrst"
var ParseICAValue = "ICA_ADDR"

// GetTxMsgs unpacks sdk messages from any messages
func (autoTxInfo AutoTxInfo) GetTxMsgs(unpacker types.AnyUnpacker) (sdkMsgs []sdk.Msg) {

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
