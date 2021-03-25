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
	cdc.RegisterConcrete(&MsgCreateFlag{}, "tpp/CreateFlag", nil)

	cdc.RegisterConcrete(&MsgCreateBuyer{}, "tpp/CreateBuyer", nil)
	cdc.RegisterConcrete(&MsgUpdateBuyer{}, "tpp/UpdateBuyer", nil)
	cdc.RegisterConcrete(&MsgDeleteBuyer{}, "tpp/DeleteBuyer", nil)
	cdc.RegisterConcrete(&MsgItemThank{}, "tpp/ItemThank", nil)
	cdc.RegisterConcrete(&MsgItemTransfer{}, "tpp/ItemTransfer", nil)


	cdc.RegisterConcrete(&MsgCreateItem{}, "tpp/CreateItem", nil)
	cdc.RegisterConcrete(&MsgUpdateItem{}, "tpp/UpdateItem", nil)
	cdc.RegisterConcrete(&MsgDeleteItem{}, "tpp/DeleteItem", nil)
	cdc.RegisterConcrete(&MsgRevealEstimation{}, "tpp/RevealEstimation", nil)
	cdc.RegisterConcrete(&MsgItemTransferable{}, "tpp/ItemTransferable", nil)
	cdc.RegisterConcrete(&MsgItemShipping{}, "tpp/ItemShipping", nil)
	cdc.RegisterConcrete(&MsgItemResell{}, "tpp/ItemResell", nil)

}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateEstimator{},
		&MsgUpdateEstimator{},
		&MsgDeleteEstimator{},
		&MsgCreateFlag{},

		&MsgCreateBuyer{},
		&MsgUpdateBuyer{},
		&MsgDeleteBuyer{},
		&MsgItemTransfer{},

		&MsgCreateItem{},
		&MsgUpdateItem{},
		&MsgDeleteItem{},
		&MsgRevealEstimation{},
		&MsgItemTransferable{},
		&MsgItemShipping{},
		&MsgItemThank{},
		&MsgItemResell{},
	)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
