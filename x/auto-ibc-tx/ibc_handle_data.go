package autoibctx

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	proto "github.com/gogo/protobuf/proto"
	msgregistry "github.com/trstlabs/trst/x/auto-ibc-tx/msg_registry"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

func handleMsgData(ctx sdk.Context, msgData *sdk.MsgData) (string, int, error) {
	fmt.Printf("handling data for typeurl: %v and data: %v\n", msgData.MsgType, msgData.Data)

	var msgResponse proto.Message
	var rewardType int

	switch msgData.MsgType {

	// authz
	case sdk.MsgTypeURL(&authztypes.MsgExec{}):
		msgResponse = &authztypes.MsgExecResponse{}
		rewardType = types.KeyAutoTxIncentiveForAuthzTx

	// sdk
	case sdk.MsgTypeURL(&banktypes.MsgSend{}):
		msgResponse = &banktypes.MsgSendResponse{}
		rewardType = types.KeyAutoTxIncentiveForSDKTx

	case sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}):
		msgResponse = &stakingtypes.MsgDelegateResponse{}
		rewardType = types.KeyAutoTxIncentiveForSDKTx

	case sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}):
		msgResponse = &stakingtypes.MsgUndelegateResponse{}
		rewardType = types.KeyAutoTxIncentiveForSDKTx

	case sdk.MsgTypeURL(&distrtypes.MsgWithdrawDelegatorReward{}):
		msgResponse = &distrtypes.MsgWithdrawDelegatorRewardResponse{}
		rewardType = types.KeyAutoTxIncentiveForSDKTx

	// wasm
	case sdk.MsgTypeURL(&msgregistry.MsgExecuteContract{}):
		msgResponse = &msgregistry.MsgExecuteContractResponse{}
		rewardType = types.KeyAutoTxIncentiveForWasmTx

	case sdk.MsgTypeURL(&msgregistry.MsgInstantiateContract{}):
		msgResponse = &msgregistry.MsgInstantiateContractResponse{}
		rewardType = types.KeyAutoTxIncentiveForWasmTx

	// osmo
	case sdk.MsgTypeURL(&msgregistry.MsgSwapExactAmountIn{}):
		msgResponse = &msgregistry.MsgSwapExactAmountInResponse{}
		rewardType = types.KeyAutoTxIncentiveForOsmoTx

	case sdk.MsgTypeURL(&msgregistry.MsgSwapExactAmountOut{}):
		msgResponse = &msgregistry.MsgSwapExactAmountOutResponse{}
		rewardType = types.KeyAutoTxIncentiveForOsmoTx

	case sdk.MsgTypeURL(&msgregistry.MsgJoinPool{}):
		msgResponse = &msgregistry.MsgJoinPoolResponse{}
		rewardType = types.KeyAutoTxIncentiveForOsmoTx

	case sdk.MsgTypeURL(&msgregistry.MsgExitPool{}):
		msgResponse = &msgregistry.MsgExitPoolResponse{}
		rewardType = types.KeyAutoTxIncentiveForOsmoTx

	default:
		return "", -1, nil
	}

	if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
		return "", -1, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal response message: %s", err.Error())
	}

	return msgResponse.String(), rewardType, nil
}
