package autoibctx

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	proto "github.com/gogo/protobuf/proto"
	msgregistry "github.com/trstlabs/trst/x/auto-ibc-tx/msg_registry"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

func handleMsgData(ctx sdk.Context, msgData *sdk.MsgData) (string, int, error) {
	fmt.Printf("handling data for typeurl: %v and data: %v\n ", msgData.MsgType, msgData.Data)

	switch msgData.MsgType {
	//authz
	case sdk.MsgTypeURL(&authztypes.MsgExec{}):
		msgResponse := &authztypes.MsgExecResponse{}
		if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
			return "", -1, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal authz exec response message: %s", err.Error())
		}
		return msgResponse.String(), types.KeyAutoTxIncentiveForAuthzTx, nil
		//sdk
	case sdk.MsgTypeURL(&banktypes.MsgSend{}):
		msgResponse := &banktypes.MsgSendResponse{}
		if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
			return "", -1, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal send response message: %s", err.Error())
		}
		return msgResponse.String(), types.KeyAutoTxIncentiveForSDKTx, nil
	case sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}):
		msgResponse := &stakingtypes.MsgDelegateResponse{}
		if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
			return "", -1, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal delegate response message: %s", err.Error())
		}
		return msgResponse.String(), types.KeyAutoTxIncentiveForSDKTx, nil
	case sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}):
		msgResponse := &stakingtypes.MsgUndelegateResponse{}
		if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
			return "", -1, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal undelegate response message: %s", err.Error())
		}
		return msgResponse.String(), types.KeyAutoTxIncentiveForSDKTx, nil
	case sdk.MsgTypeURL(&msgregistry.MsgExecuteContract{}):
		msgResponse := &msgregistry.MsgExecuteContractResponse{}
		if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
			return "", -1, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal msg execute response message: %s", err.Error())
		}
		return msgResponse.String(), types.KeyAutoTxIncentiveForWasmTx, nil
		//wasm
	case sdk.MsgTypeURL(&msgregistry.MsgInstantiateContract{}):
		msgResponse := &msgregistry.MsgInstantiateContractResponse{}
		if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
			return "", -1, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal msg MsgInstantiateContract response message: %s", err.Error())
		}
		return msgResponse.String(), types.KeyAutoTxIncentiveForWasmTx, nil
	case sdk.MsgTypeURL(&msgregistry.MsgSwapExactAmountIn{}):
		msgResponse := &msgregistry.MsgSwapExactAmountInResponse{}
		if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
			return "", -1, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal MsgSwapExactAmountIn response message: %s", err.Error())
		}
		return msgResponse.String(), types.KeyAutoTxIncentiveForOsmoTx, nil
	case sdk.MsgTypeURL(&msgregistry.MsgSwapExactAmountOut{}):
		msgResponse := &msgregistry.MsgSwapExactAmountOutResponse{}
		if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
			return "", -1, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal MsgSwapExactAmountOut response message: %s", err.Error())
		}
		return msgResponse.String(), types.KeyAutoTxIncentiveForOsmoTx, nil
	case sdk.MsgTypeURL(&msgregistry.MsgJoinPool{}):
		msgResponse := &msgregistry.MsgJoinPoolResponse{}
		if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
			return "", -1, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal MsgJoinPool response message: %s", err.Error())
		}
		return msgResponse.String(), types.KeyAutoTxIncentiveForOsmoTx, nil
	case sdk.MsgTypeURL(&msgregistry.MsgExitPool{}):
		msgResponse := &msgregistry.MsgExitPoolResponse{}
		if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
			return "", -1, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal MsgExitPool response message: %s", err.Error())
		}
		return msgResponse.String(), types.KeyAutoTxIncentiveForOsmoTx, nil

	// TODO: handle other messages
	default:

		return "", -1, nil
	}
}
