package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var Denom = "utrst"

// GetTxMsgs fetches cached any messages
func (autoTxInfo AutoTxInfo) GetTxMsgs() []sdk.Msg {
	var sdkMsgs []sdk.Msg
	for _, message := range autoTxInfo.Msgs {
		sdkMsg, ok := message.GetCachedValue().(sdk.Msg)
		if !ok {
			return nil
		}
		sdkMsgs = append(sdkMsgs, sdkMsg)
	}

	return sdkMsgs
}
