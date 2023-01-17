package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	// this line is used by starport scaffolding # 1
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgRegisterAccount{}, "auto-ibc-tx/MsgRegisterAccount", nil)
	cdc.RegisterConcrete(MsgSubmitTx{}, "auto-ibc-tx/MsgSendTx", nil)
	cdc.RegisterConcrete(MsgSubmitAutoTx{}, "auto-ibc-tx/MsgSendAutoTx", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgRegisterAccount{},
		&MsgSubmitTx{},
		&MsgSubmitAutoTx{},
	)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
