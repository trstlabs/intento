package types

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/cosmos/gogoproto/proto"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

var Denom = "uinto"
var ParseICAValue = "ICA_ADDR"
var MaxGas uint64 = 1_000_000
var MaxGasTotal uint64 = 500_000_000

// GetTxMsgs unpacks sdk messages from any messages
func (flow Flow) GetTxMsgs(unpacker types.AnyUnpacker) (sdkMsgs []sdk.Msg) {

	for _, message := range flow.Msgs {
		var sdkMsg sdk.Msg
		err := unpacker.UnpackAny(message, &sdkMsg)
		if err != nil {
			return nil
		}
		sdkMsgs = append(sdkMsgs, sdkMsg)
	}
	return sdkMsgs
}

func (flow Flow) IsInterchain() bool {
	if flow.SelfHostedICA.ConnectionID != "" {
		return true
	}
	if flow.TrustlessAgent.AgentAddress != "" {
		return true
	}
	return false
}

func (flow Flow) IsLocal() bool {
	return !flow.IsInterchain()
}

func (flow Flow) IsICQ() bool {
	if flow.Conditions.FeedbackLoops != nil {
		for _, loop := range flow.Conditions.FeedbackLoops {
			if loop.ICQConfig != nil {
				return true
			}
		}
	}
	if flow.Conditions.Comparisons != nil {
		for _, comparison := range flow.Conditions.Comparisons {
			if comparison.ICQConfig != nil {
				return true
			}
		}
	}
	return false
}

var GasFeeCoinsSupported sdk.Coins = sdk.Coins{sdk.NewCoin(Denom, math.NewInt(10))}

// GetTxMsgs unpacks sdk messages from any messages
func GetTransferMsg(cdc codec.Codec, anyTransfer *types.Any) (transferMsg ibctransfertypes.MsgTransfer, err error) {

	if err := proto.Unmarshal(anyTransfer.Value, &transferMsg); err != nil {
		return ibctransfertypes.MsgTransfer{}, err
	}

	if err := transferMsg.ValidateBasic(); err != nil {
		return ibctransfertypes.MsgTransfer{}, err
	}

	return transferMsg, nil
}
