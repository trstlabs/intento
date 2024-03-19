package msg_registry

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {

	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		//cosmwasm v1
		&MsgInstantiateContract{},
		&MsgInstantiateContract2{},
		&MsgExecuteContract{},
		//osmosis gamm v1beta1
		&MsgExitPool{},
		&MsgExitSwapExternAmountOut{},
		&MsgExitSwapShareAmountIn{},
		&MsgJoinPool{},
		&MsgJoinSwapExternAmountIn{},
		&MsgJoinSwapShareAmountOut{},
		&MsgSwapExactAmountIn{},
		&MsgSwapExactAmountOut{},
	)
}

type RawContractMessage []byte
