package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2
	cdc.RegisterConcrete(&MsgCreateEstimator{}, "tpp/CreateEstimator", nil)
	cdc.RegisterConcrete(&MsgUpdateEstimator{}, "tpp/UpdateEstimator", nil)
	cdc.RegisterConcrete(&MsgDeleteEstimator{}, "tpp/DeleteEstimator", nil)

	cdc.RegisterConcrete(&MsgCreateBuyer{}, "tpp/CreateBuyer", nil)
	cdc.RegisterConcrete(&MsgUpdateBuyer{}, "tpp/UpdateBuyer", nil)
	cdc.RegisterConcrete(&MsgDeleteBuyer{}, "tpp/DeleteBuyer", nil)

	cdc.RegisterConcrete(&MsgCreateItem{}, "tpp/CreateItem", nil)
	cdc.RegisterConcrete(&MsgUpdateItem{}, "tpp/UpdateItem", nil)
	cdc.RegisterConcrete(&MsgDeleteItem{}, "tpp/DeleteItem", nil)

}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateEstimator{},
		&MsgUpdateEstimator{},
		&MsgDeleteEstimator{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateBuyer{},
		&MsgUpdateBuyer{},
		&MsgDeleteBuyer{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateItem{},
		&MsgUpdateItem{},
		&MsgDeleteItem{},
	)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
