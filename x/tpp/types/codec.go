package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2
	cdc.RegisterConcrete(&MsgCreateEstimation{}, "tpp/CreateEstimation", nil)
	cdc.RegisterConcrete(&MsgUpdateLike{}, "tpp/UpdateLike", nil)
	cdc.RegisterConcrete(&MsgDeleteEstimation{}, "tpp/DeleteEstimation", nil)
	cdc.RegisterConcrete(&MsgFlagItem{}, "tpp/FlagItem", nil)

	cdc.RegisterConcrete(&MsgPrepayment{}, "tpp/Prepayment", nil)
	cdc.RegisterConcrete(&MsgUpdateBuyer{}, "tpp/UpdateBuyer", nil)
	cdc.RegisterConcrete(&MsgWithdrawal{}, "tpp/Withdrawal", nil)
	cdc.RegisterConcrete(&MsgItemRating{}, "tpp/ItemRating", nil)
	cdc.RegisterConcrete(&MsgItemTransfer{}, "tpp/ItemTransfer", nil)

	cdc.RegisterConcrete(&MsgCreateItem{}, "tpp/CreateItem", nil)
	cdc.RegisterConcrete(&MsgUpdateItem{}, "tpp/UpdateItem", nil)
	cdc.RegisterConcrete(&MsgDeleteItem{}, "tpp/DeleteItem", nil)
	cdc.RegisterConcrete(&MsgRevealEstimation{}, "tpp/RevealEstimation", nil)
	cdc.RegisterConcrete(&MsgItemTransferable{}, "tpp/ItemTransferable", nil)
	cdc.RegisterConcrete(&MsgItemShipping{}, "tpp/ItemShipping", nil)
	cdc.RegisterConcrete(&MsgItemResell{}, "tpp/ItemResell", nil)
	cdc.RegisterConcrete(&MsgTokenizeItem{}, "tpp/TokenizeItem", nil)
	cdc.RegisterConcrete(&MsgUnTokenizeItem{}, "tpp/UnTokenizeItem", nil)

}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateEstimation{},
		&MsgUpdateLike{},
		&MsgDeleteEstimation{},
		&MsgFlagItem{},

		&MsgPrepayment{},

		&MsgWithdrawal{},
		&MsgItemTransfer{},

		&MsgCreateItem{},

		&MsgDeleteItem{},
		&MsgRevealEstimation{},
		&MsgItemTransferable{},
		&MsgItemShipping{},
		&MsgItemRating{},
		&MsgItemResell{},
		&MsgUnTokenizeItem{},
		&MsgTokenizeItem{},
	)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
