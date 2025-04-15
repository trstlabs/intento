package keeper

import (
	errorsmod "cosmossdk.io/errors"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	proto "github.com/cosmos/gogoproto/proto"
	osmosisgammv1beta1 "github.com/trstlabs/intento/x/intent/msg_registry/osmosis/gamm/v1beta1"
	"github.com/trstlabs/intento/x/intent/types"
)

func handleMsgData(msgData *sdk.MsgData) (proto.Message, int, error) {
	// fmt.Printf("handling data for typeurl: %v and data: %v\n", msgData.MsgType, msgData.Data)

	var msgResponse proto.Message
	var rewardType int
	switch msgData.MsgType {

	// authz
	case sdk.MsgTypeURL(&authztypes.MsgExec{}):
		msgResponse = &authztypes.MsgExecResponse{}
		rewardType = types.KeyFlowIncentiveForAuthzTx

	// sdk
	case sdk.MsgTypeURL(&banktypes.MsgSend{}):
		msgResponse = &banktypes.MsgSendResponse{}
		rewardType = types.KeyFlowIncentiveForSDKTx

	case sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}):
		msgResponse = &stakingtypes.MsgDelegateResponse{}
		rewardType = types.KeyFlowIncentiveForSDKTx

	case sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}):
		msgResponse = &stakingtypes.MsgUndelegateResponse{}
		rewardType = types.KeyFlowIncentiveForSDKTx

	case sdk.MsgTypeURL(&distrtypes.MsgWithdrawDelegatorReward{}):
		msgResponse = &distrtypes.MsgWithdrawDelegatorRewardResponse{}
		rewardType = types.KeyFlowIncentiveForSDKTx

	// wasm
	case sdk.MsgTypeURL(&wasmtypes.MsgExecuteContract{}):
		msgResponse = &wasmtypes.MsgExecuteContractResponse{}
		rewardType = types.KeyFlowIncentiveForWasmTx

	case sdk.MsgTypeURL(&wasmtypes.MsgInstantiateContract{}):
		msgResponse = &wasmtypes.MsgInstantiateContractResponse{}
		rewardType = types.KeyFlowIncentiveForWasmTx

	// osmo
	case sdk.MsgTypeURL(&osmosisgammv1beta1.MsgSwapExactAmountIn{}):
		msgResponse = &osmosisgammv1beta1.MsgSwapExactAmountInResponse{}
		rewardType = types.KeyFlowIncentiveForOsmoTx

	case sdk.MsgTypeURL(&osmosisgammv1beta1.MsgSwapExactAmountOut{}):
		msgResponse = &osmosisgammv1beta1.MsgSwapExactAmountOutResponse{}
		rewardType = types.KeyFlowIncentiveForOsmoTx

	case sdk.MsgTypeURL(&osmosisgammv1beta1.MsgJoinPool{}):
		msgResponse = &osmosisgammv1beta1.MsgJoinPoolResponse{}
		rewardType = types.KeyFlowIncentiveForOsmoTx

	case sdk.MsgTypeURL(&osmosisgammv1beta1.MsgExitPool{}):
		msgResponse = &osmosisgammv1beta1.MsgExitPoolResponse{}
		rewardType = types.KeyFlowIncentiveForOsmoTx

	default:
		return nil, -1, nil
	}

	if err := proto.Unmarshal(msgData.Data, msgResponse); err != nil {
		return nil, -1, errorsmod.Wrapf(types.ErrJSONUnmarshal, "cannot unmarshal response message: %s", err.Error())
	}

	return msgResponse, rewardType, nil
}

func getMsgRewardType(typeUrl string) int {
	var rewardType int

	switch typeUrl {

	// authz
	case sdk.MsgTypeURL(&authztypes.MsgExec{}):
		rewardType = types.KeyFlowIncentiveForAuthzTx

	// sdk
	case sdk.MsgTypeURL(&banktypes.MsgSend{}):
		rewardType = types.KeyFlowIncentiveForSDKTx

	case sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}):
		rewardType = types.KeyFlowIncentiveForSDKTx

	case sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}):
		rewardType = types.KeyFlowIncentiveForSDKTx

	case sdk.MsgTypeURL(&distrtypes.MsgWithdrawDelegatorReward{}):
		rewardType = types.KeyFlowIncentiveForSDKTx

	// wasm
	case sdk.MsgTypeURL(&wasmtypes.MsgExecuteContract{}):
		rewardType = types.KeyFlowIncentiveForWasmTx

	case sdk.MsgTypeURL(&wasmtypes.MsgInstantiateContract{}):
		rewardType = types.KeyFlowIncentiveForWasmTx

	// osmo
	case sdk.MsgTypeURL(&osmosisgammv1beta1.MsgSwapExactAmountIn{}):
		rewardType = types.KeyFlowIncentiveForOsmoTx

	case sdk.MsgTypeURL(&osmosisgammv1beta1.MsgSwapExactAmountOut{}):
		rewardType = types.KeyFlowIncentiveForOsmoTx

	case sdk.MsgTypeURL(&osmosisgammv1beta1.MsgJoinPool{}):
		rewardType = types.KeyFlowIncentiveForOsmoTx

	case sdk.MsgTypeURL(&osmosisgammv1beta1.MsgExitPool{}):
		rewardType = types.KeyFlowIncentiveForOsmoTx

	default:
		return -1
	}

	return rewardType
}
