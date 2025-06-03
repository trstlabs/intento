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
	cdc.RegisterConcrete(MsgSubmitFlow{}, "intent/MsgSendFlow", nil)
	cdc.RegisterConcrete(MsgRegisterAccountAndSubmitFlow{}, "intent/MsgRegisterAccountAndSubmitFlow", nil)
	cdc.RegisterConcrete(MsgUpdateFlow{}, "intent/MsgUpdateFlow", nil)
	cdc.RegisterConcrete(MsgUpdateParams{}, "intent/MsgUpdateParams", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgRegisterAccount{},
		&MsgSubmitTx{},
		&MsgSubmitFlow{},
		&MsgRegisterAccountAndSubmitFlow{},
		&MsgUpdateFlow{},
		&MsgUpdateParams{},
	)

}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
