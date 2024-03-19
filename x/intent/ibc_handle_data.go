package intent

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	proto "github.com/cosmos/gogoproto/proto"
	msgregistry "github.com/trstlabs/intento/x/intent/msg_registry"
	"github.com/trstlabs/intento/x/intent/types"
)

func handleMsgData(ctx sdk.Context, msgData *sdk.MsgData) (proto.Message, int, error) {
	// fmt.Printf("handling data for typeurl: %v and data: %v\n", msgData.MsgType, msgData.Data)

	var msgResponse proto.Message
	var rewardType int
	switch msgData.MsgType {

	// authz
	case sdk.MsgTypeURL(&authztypes.MsgExec{}):
		msgResponse = &authztypes.MsgExecResponse{}
		rewardType = types.KeyActionIncentiveForAuthzTx

	// sdk
	case sdk.MsgTypeURL(&banktypes.MsgSend{}):
		msgResponse = &banktypes.MsgSendResponse{}
		rewardType = types.KeyActionIncentiveForSDKTx

	case sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}):
		msgResponse = &stakingtypes.MsgDelegateResponse{}
		rewardType = types.KeyActionIncentiveForSDKTx

	case sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}):
		msgResponse = &stakingtypes.MsgUndelegateResponse{}
		rewardType = types.KeyActionIncentiveForSDKTx

	case sdk.MsgTypeURL(&distrtypes.MsgWithdrawDelegatorReward{}):
		msgResponse = &distrtypes.MsgWithdrawDelegatorRewardResponse{}
		rewardType = types.KeyActionIncentiveForSDKTx

	// wasm
	case sdk.MsgTypeURL(&msgregistry.MsgExecuteContract{}):
		msgResponse = &msgregistry.MsgExecuteContractResponse{}
		rewardType = types.KeyActionIncentiveForWasmTx

	case sdk.MsgTypeURL(&msgregistry.MsgInstantiateContract{}):
		msgResponse = &msgregistry.MsgInstantiateContractResponse{}
		rewardType = types.KeyActionIncentiveForWasmTx

	// osmo
	case sdk.MsgTypeURL(&msgregistry.MsgSwapExactAmountIn{}):
		msgResponse = &msgregistry.MsgSwapExactAmountInResponse{}
		rewardType = types.KeyActionIncentiveForOsmoTx

	case sdk.MsgTypeURL(&msgregistry.MsgSwapExactAmountOut{}):
		msgResponse = &msgregistry.MsgSwapExactAmountOutResponse{}
		rewardType = types.KeyActionIncentiveForOsmoTx

	case sdk.MsgTypeURL(&msgregistry.MsgJoinPool{}):
		msgResponse = &msgregistry.MsgJoinPoolResponse{}
		rewardType = types.KeyActionIncentiveForOsmoTx

	case sdk.MsgTypeURL(&msgregistry.MsgExitPool{}):
		msgResponse = &msgregistry.MsgExitPoolResponse{}
		rewardType = types.KeyActionIncentiveForOsmoTx

	default:
		return nil, -1, nil
	}

	if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
		return nil, -1, errorsmod.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal response message: %s", err.Error())
	}

	return msgResponse, rewardType, nil
}

func getMsgRewardType(ctx sdk.Context, typeUrl string) int {
	var rewardType int

	switch typeUrl {

	// authz
	case sdk.MsgTypeURL(&authztypes.MsgExec{}):
		rewardType = types.KeyActionIncentiveForAuthzTx

	// sdk
	case sdk.MsgTypeURL(&banktypes.MsgSend{}):
		rewardType = types.KeyActionIncentiveForSDKTx

	case sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}):
		rewardType = types.KeyActionIncentiveForSDKTx

	case sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}):
		rewardType = types.KeyActionIncentiveForSDKTx

	case sdk.MsgTypeURL(&distrtypes.MsgWithdrawDelegatorReward{}):
		rewardType = types.KeyActionIncentiveForSDKTx

	// wasm
	case sdk.MsgTypeURL(&msgregistry.MsgExecuteContract{}):
		rewardType = types.KeyActionIncentiveForWasmTx

	case sdk.MsgTypeURL(&msgregistry.MsgInstantiateContract{}):
		rewardType = types.KeyActionIncentiveForWasmTx

	// osmo
	case sdk.MsgTypeURL(&msgregistry.MsgSwapExactAmountIn{}):
		rewardType = types.KeyActionIncentiveForOsmoTx

	case sdk.MsgTypeURL(&msgregistry.MsgSwapExactAmountOut{}):
		rewardType = types.KeyActionIncentiveForOsmoTx

	case sdk.MsgTypeURL(&msgregistry.MsgJoinPool{}):
		rewardType = types.KeyActionIncentiveForOsmoTx

	case sdk.MsgTypeURL(&msgregistry.MsgExitPool{}):
		rewardType = types.KeyActionIncentiveForOsmoTx

	default:
		return -1
	}

	return rewardType
}
