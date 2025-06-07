package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	"github.com/trstlabs/intento/x/intent/msg_registry"
	"github.com/trstlabs/intento/x/intent/types"
)

func handleMsgData(msgData *sdk.MsgData) (proto.Message, int, error) {
	entry, ok := msg_registry.MsgRegistry[msgData.MsgType]
	if !ok {
		return nil, -1, nil // unknown msg type
	}

	msgResponse := entry.NewResponse()
	if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
		return nil, -1, errorsmod.Wrapf(types.ErrJSONUnmarshal, "cannot unmarshal response message: %s", err.Error())
	}

	return msgResponse, entry.RewardType, nil
}
