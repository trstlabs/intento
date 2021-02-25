package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	//"cosmos/base/v1beta1/coin.proto"
)

var _ sdk.Msg = &MsgCreateBuyer{}

func NewMsgCreateBuyer(buyer string, itemid string, transferable bool, deposit sdk.Coin) *MsgCreateBuyer {
	return &MsgCreateBuyer{
		Buyer:      buyer,
		Itemid:       itemid,
		Transferable: transferable,
		Deposit:      deposit,
	}
}

func (msg *MsgCreateBuyer) Route() string {
	return RouterKey
}

func (msg *MsgCreateBuyer) Type() string {
	return "CreateBuyer"
}

func (msg *MsgCreateBuyer) GetSigners() []sdk.AccAddress {
	buyer, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{buyer}
}

func (msg *MsgCreateBuyer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateBuyer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateBuyer{}

func NewMsgUpdateBuyer(buyer string, itemid string, transferable bool, deposit sdk.Coin) *MsgUpdateBuyer {
	return &MsgUpdateBuyer{

		Buyer:      buyer,
		Itemid:       itemid,
		Transferable: transferable,
		Deposit:      deposit,
	}
}

func (msg *MsgUpdateBuyer) Route() string {
	return RouterKey
}

func (msg *MsgUpdateBuyer) Type() string {
	return "UpdateBuyer"
}

func (msg *MsgUpdateBuyer) GetSigners() []sdk.AccAddress {
	buyer, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{buyer}
}

func (msg *MsgUpdateBuyer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateBuyer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateBuyer{}

func NewMsgDeleteBuyer(buyer string, itemid string) *MsgDeleteBuyer {
	return &MsgDeleteBuyer{
		Itemid:      itemid,
		Buyer: buyer,
	}
}
func (msg *MsgDeleteBuyer) Route() string {
	return RouterKey
}

func (msg *MsgDeleteBuyer) Type() string {
	return "DeleteBuyer"
}

func (msg *MsgDeleteBuyer) GetSigners() []sdk.AccAddress {
	buyer, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{buyer}
}

func (msg *MsgDeleteBuyer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteBuyer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
