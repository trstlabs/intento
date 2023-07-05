package autoibctx

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	proto "github.com/cosmos/gogoproto/proto"
	msgregistry "github.com/trstlabs/trst/x/auto-ibc-tx/msg_registry"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

func handleMsgData(ctx sdk.Context, msgData *sdk.MsgData) (string, int, error) {
	//fmt.Printf("handling data for typeurl: %v and data: %v\n", msgData.MsgType, msgData.Data)

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
		return "", -1, errorsmod.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal response message: %s", err.Error())
	}

	return msgResponse.String(), rewardType, nil
}

func getMsgRewardType(ctx sdk.Context, typeUrl string) int {
	var rewardType int

	switch typeUrl {

	// authz
	case sdk.MsgTypeURL(&authztypes.MsgExec{}):
		rewardType = types.KeyAutoTxIncentiveForAuthzTx

	// sdk
	case sdk.MsgTypeURL(&banktypes.MsgSend{}):
		rewardType = types.KeyAutoTxIncentiveForSDKTx

	case sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}):
		rewardType = types.KeyAutoTxIncentiveForSDKTx

	case sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}):
		rewardType = types.KeyAutoTxIncentiveForSDKTx

	case sdk.MsgTypeURL(&distrtypes.MsgWithdrawDelegatorReward{}):
		rewardType = types.KeyAutoTxIncentiveForSDKTx

	// wasm
	case sdk.MsgTypeURL(&msgregistry.MsgExecuteContract{}):
		rewardType = types.KeyAutoTxIncentiveForWasmTx

	case sdk.MsgTypeURL(&msgregistry.MsgInstantiateContract{}):
		rewardType = types.KeyAutoTxIncentiveForWasmTx

	// osmo
	case sdk.MsgTypeURL(&msgregistry.MsgSwapExactAmountIn{}):
		rewardType = types.KeyAutoTxIncentiveForOsmoTx

	case sdk.MsgTypeURL(&msgregistry.MsgSwapExactAmountOut{}):
		rewardType = types.KeyAutoTxIncentiveForOsmoTx

	case sdk.MsgTypeURL(&msgregistry.MsgJoinPool{}):
		rewardType = types.KeyAutoTxIncentiveForOsmoTx

	case sdk.MsgTypeURL(&msgregistry.MsgExitPool{}):
		rewardType = types.KeyAutoTxIncentiveForOsmoTx

	default:
		return -1
	}

	return rewardType
}
