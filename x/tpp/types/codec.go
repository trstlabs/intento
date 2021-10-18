package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2
	cdc.RegisterConcrete(&MsgCreateEstimation{}, "trst/CreateEstimation", nil)
	cdc.RegisterConcrete(&MsgUpdateLike{}, "trst/UpdateLike", nil)
	cdc.RegisterConcrete(&MsgDeleteEstimation{}, "trst/DeleteEstimation", nil)
	cdc.RegisterConcrete(&MsgFlagItem{}, "trst/FlagItem", nil)

	cdc.RegisterConcrete(&MsgPrepayment{}, "trst/Prepayment", nil)
	cdc.RegisterConcrete(&MsgUpdateBuyer{}, "trst/UpdateBuyer", nil)
	cdc.RegisterConcrete(&MsgWithdrawal{}, "trst/Withdrawal", nil)
	cdc.RegisterConcrete(&MsgItemRating{}, "trst/ItemRating", nil)
	cdc.RegisterConcrete(&MsgItemTransfer{}, "trst/ItemTransfer", nil)

	cdc.RegisterConcrete(&MsgCreateItem{}, "trst/CreateItem", nil)
	cdc.RegisterConcrete(&MsgUpdateItem{}, "trst/UpdateItem", nil)
	cdc.RegisterConcrete(&MsgDeleteItem{}, "trst/DeleteItem", nil)
	cdc.RegisterConcrete(&MsgRevealEstimation{}, "trst/RevealEstimation", nil)
	cdc.RegisterConcrete(&MsgItemTransferable{}, "trst/ItemTransferable", nil)
	cdc.RegisterConcrete(&MsgItemShipping{}, "trst/ItemShipping", nil)
	cdc.RegisterConcrete(&MsgItemResell{}, "trst/ItemResell", nil)
	cdc.RegisterConcrete(&MsgTokenizeItem{}, "trst/TokenizeItem", nil)
	cdc.RegisterConcrete(&MsgUnTokenizeItem{}, "trst/UnTokenizeItem", nil)

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
