package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	//"cosmos/base/v1beta1/coin.proto"
)

var _ sdk.Msg = &MsgCreateBuyer{}

func NewMsgCreateBuyer(buyer string, itemid string, deposit int64) *MsgCreateBuyer {


	return &MsgCreateBuyer{
		Buyer:  buyer,
		Itemid: itemid,
		Deposit: deposit,
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

func NewMsgUpdateBuyer(buyer string, itemid string, deposit int64) *MsgUpdateBuyer {
	return &MsgUpdateBuyer{

		Buyer:        buyer,
		Itemid:       itemid,
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
		Itemid: itemid,
		Buyer:  buyer,
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

var _ sdk.Msg = &MsgCreateBuyer{}

func NewMsgItemTransfer(buyer string, itemid string) *MsgItemTransfer {
	return &MsgItemTransfer{
		Buyer:        buyer,
		Itemid:       itemid,

	}
}

func (msg *MsgItemTransfer) Route() string {
	return RouterKey
}

func (msg *MsgItemTransfer) Type() string {
	return "ItemTransfer"
}

func (msg *MsgItemTransfer) GetSigners() []sdk.AccAddress {
	buyer, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{buyer}
}

func (msg *MsgItemTransfer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgItemTransfer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid buyer address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateBuyer{}

func NewMsgItemRating(buyer string, itemid string, rating int64, note string) *MsgItemRating {
	return &MsgItemRating{
		Buyer:  buyer,
		Itemid: itemid,
		Rating:  rating,
		Note: note,
	}
}

func (msg *MsgItemRating) Route() string {
	return RouterKey
}

func (msg *MsgItemRating) Type() string {
	return "ItemRating"
}

func (msg *MsgItemRating) GetSigners() []sdk.AccAddress {
	buyer, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{buyer}
}

func (msg *MsgItemRating) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgItemRating) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid buyer address (%s)", err)
	}
	if msg.Rating > 5 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid item condition")
	}

	if len(msg.Note) > 240 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "note too long")
	}


	return nil
}

var _ sdk.Msg = &MsgCreateBuyer{}
