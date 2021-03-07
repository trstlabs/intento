package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	//"cosmos/base/v1beta1/coin.proto"
)

var _ sdk.Msg = &MsgCreateItem{}

func NewMsgCreateItem(creator string, title string, description string, shippingcost int64, localpickup bool, estimationcount int64, tags []string, condition int64, shippingregion []string) *MsgCreateItem {
	return &MsgCreateItem{

		Creator:         creator,
		Title:           title,
		Description:     description,
		Shippingcost:    shippingcost,
		Localpickup:     localpickup,
		Estimationcount: estimationcount,
		Tags:            tags,
		Condition:       condition,
		Shippingregion:  shippingregion,
	}
}

func (msg *MsgCreateItem) Route() string {
	return RouterKey
}

func (msg *MsgCreateItem) Type() string {
	return "CreateItem"
}

func (msg *MsgCreateItem) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateItem) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateItem) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.Tags) > 5 || len(msg.Tags) < 1 {
		return sdkerrors.Wrap(sdkerrors.ErrMemoTooLarge, "tags invalid")
	}
	if len(msg.Shippingregion) > 5 || len(msg.Shippingregion) < 1 {
		return sdkerrors.Wrap(sdkerrors.ErrMemoTooLarge, "Region invalid")
	}
	if len(msg.Description) > 500 {
		return sdkerrors.Wrap(sdkerrors.ErrMemoTooLarge, "description too long")
	}
	if msg.Condition > 6 {
		return sdkerrors.Wrap(sdkerrors.ErrMemoTooLarge, "invalid item condition")
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateItem{}

func NewMsgUpdateItem(creator string, id string, shippingcost int64, localpickup bool, shippingregion []string) *MsgUpdateItem {
	return &MsgUpdateItem{
		Id:             id,
		Creator:        creator,
		Shippingcost:   shippingcost,
		Localpickup:    localpickup,
		Shippingregion: shippingregion,
	}
}

func (msg *MsgUpdateItem) Route() string {
	return RouterKey
}

func (msg *MsgUpdateItem) Type() string {
	return "UpdateItem"
}

func (msg *MsgUpdateItem) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateItem) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateItem) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateItem{}

func NewMsgDeleteItem(creator string, id string) *MsgDeleteItem {
	return &MsgDeleteItem{
		Id:      id,
		Creator: creator,
	}
}
func (msg *MsgDeleteItem) Route() string {
	return RouterKey
}

func (msg *MsgDeleteItem) Type() string {
	return "DeleteItem"
}

func (msg *MsgDeleteItem) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteItem) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteItem) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func NewMsgRevealEstimation(creator string, itemid string) *MsgRevealEstimation {
	return &MsgRevealEstimation{

		Creator: creator,
		Itemid:  itemid,
	}
}

func (msg *MsgRevealEstimation) Route() string {
	return RouterKey
}

func (msg *MsgRevealEstimation) Type() string {
	return "RevealEstimation"
}

func (msg *MsgRevealEstimation) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRevealEstimation) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRevealEstimation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateItem{}

func NewMsgItemTransferable(creator string, transferable bool, itemid string) *MsgItemTransferable {
	return &MsgItemTransferable{

		Creator:      creator,
		Transferable: transferable,
		Itemid:       itemid,
	}
}

func (msg *MsgItemTransferable) Route() string {
	return RouterKey
}

func (msg *MsgItemTransferable) Type() string {
	return "ItemTransferable"
}

func (msg *MsgItemTransferable) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgItemTransferable) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgItemTransferable) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateItem{}

func NewMsgItemShipping(creator string, tracking bool, itemid string) *MsgItemShipping {
	return &MsgItemShipping{

		Creator:  creator,
		Tracking: tracking,
		Itemid:   itemid,
	}
}

func (msg *MsgItemShipping) Route() string {
	return RouterKey
}

func (msg *MsgItemShipping) Type() string {
	return "ItemShipping"
}

func (msg *MsgItemShipping) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgItemShipping) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgItemShipping) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgItemShipping{}
