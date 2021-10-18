package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	//"cosmos/base/v1beta1/coin.proto"
)

var _ sdk.Msg = &MsgPrepayment{}

func NewMsgPrepayment(buyer string, itemid uint64, deposit int64) *MsgPrepayment {

	return &MsgPrepayment{
		Buyer:   buyer,
		Itemid:  itemid,
		Deposit: deposit,
	}
}

func (msg *MsgPrepayment) Route() string {
	return RouterKey
}

func (msg *MsgPrepayment) Type() string {
	return "ItemPrepayment"
}

func (msg *MsgPrepayment) GetSigners() []sdk.AccAddress {
	buyer, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{buyer}
}

func (msg *MsgPrepayment) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgPrepayment) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid buyer address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgPrepayment{}

func NewMsgWithdrawal(buyer string, itemid uint64) *MsgWithdrawal {
	return &MsgWithdrawal{
		Itemid: itemid,
		Buyer:  buyer,
	}
}
func (msg *MsgWithdrawal) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawal) Type() string {
	return "Withdrawal"
}

func (msg *MsgWithdrawal) GetSigners() []sdk.AccAddress {
	buyer, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{buyer}
}

func (msg *MsgWithdrawal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawal) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid buyer address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgPrepayment{}

func NewMsgItemTransfer(buyer string, itemid uint64) *MsgItemTransfer {
	return &MsgItemTransfer{
		Buyer:  buyer,
		Itemid: itemid,
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

var _ sdk.Msg = &MsgPrepayment{}

func NewMsgItemRating(buyer string, itemid uint64, rating int64, note string) *MsgItemRating {
	return &MsgItemRating{
		Buyer:  buyer,
		Itemid: itemid,
		Rating: rating,
		Note:   note,
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

var _ sdk.Msg = &MsgPrepayment{}
