package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgRegisterAccount{}, "intent/MsgRegisterAccount", nil)
	cdc.RegisterConcrete(MsgSubmitTx{}, "intent/MsgSendTx", nil)
	cdc.RegisterConcrete(MsgSubmitAction{}, "intent/MsgSendAction", nil)
	cdc.RegisterConcrete(MsgRegisterAccountAndSubmitAction{}, "intent/MsgRegisterAccountAndSubmitAction", nil)
	cdc.RegisterConcrete(MsgUpdateAction{}, "intent/MsgUpdateAction", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgRegisterAccount{},
		&MsgSubmitTx{},
		&MsgSubmitAction{},
		&MsgRegisterAccountAndSubmitAction{},
		&MsgUpdateAction{},
	)

}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
