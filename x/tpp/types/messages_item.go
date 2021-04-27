package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	//"cosmos/base/v1beta1/coin.proto"
)

var _ sdk.Msg = &MsgCreateItem{}

func NewMsgCreateItem(creator string, title string, description string, shippingcost int64, localpickup bool, estimationcount int64, tags []string, condition int64, shippingregion []string, depositamount int64) *MsgCreateItem {
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
		Depositamount:   depositamount,
	}
}

func (msg *MsgCreateItem) Route() string {
	return RouterKey
}

func (msg *MsgCreateItem) Type() string {
	return "CreateItem"
}

func (msg *MsgCreateItem) GetSigners() []sdk.AccAddress {
	seller, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{seller}
}

func (msg *MsgCreateItem) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateItem) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid seller address (%s)", err)
	}

	if len(msg.Tags) > 5 || len(msg.Tags) < 1 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "tags invalid")
	}
	for _, tags := range msg.Tags {
		if len(tags) > 16 {
			return sdkerrors.Wrap(sdkerrors.ErrMemoTooLarge, "tag too long")
		}
	}

	if len(msg.Shippingregion) > 6 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Region invalid")
	}

	for _, region := range msg.Shippingregion {
		if len(region) > 2{
			return sdkerrors.Wrap(sdkerrors.ErrMemoTooLarge, "Region too long")
		}
	}
	
	if len(msg.Description) > 800 {
		return sdkerrors.Wrap(sdkerrors.ErrMemoTooLarge, "description too long")
	}

	if msg.Shippingcost == 0 && msg.Localpickup != true {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Provide either shipping or localpickup")
	}
	if msg.Condition > 6 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid item condition")
	}
	if msg.Estimationcount > 24 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid estimation count")
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateItem{}

func NewMsgUpdateItem(seller string, id string, shippingcost int64, localpickup bool, shippingregion []string) *MsgUpdateItem {
	return &MsgUpdateItem{
		Id:             id,
		Seller:        seller,
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
	seller, err := sdk.AccAddressFromBech32(msg.Seller)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{seller}
}

func (msg *MsgUpdateItem) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateItem) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Seller)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid seller address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateItem{}

func NewMsgDeleteItem(seller string, id string) *MsgDeleteItem {
	return &MsgDeleteItem{
		Id:      id,
		Seller: seller,
	}
}
func (msg *MsgDeleteItem) Route() string {
	return RouterKey
}

func (msg *MsgDeleteItem) Type() string {
	return "DeleteItem"
}

func (msg *MsgDeleteItem) GetSigners() []sdk.AccAddress {
	seller, err := sdk.AccAddressFromBech32(msg.Seller)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{seller}
}

func (msg *MsgDeleteItem) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteItem) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Seller)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid seller address (%s)", err)
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid seller address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateItem{}

func NewMsgItemTransferable(seller string, transferable bool, itemid string) *MsgItemTransferable {
	return &MsgItemTransferable{

		Seller:      seller,
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
	seller, err := sdk.AccAddressFromBech32(msg.Seller)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{seller}
}

func (msg *MsgItemTransferable) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgItemTransferable) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Seller)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid seller address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateItem{}

func NewMsgItemShipping(seller string, tracking bool, itemid string) *MsgItemShipping {
	return &MsgItemShipping{

		Seller:  seller,
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
	seller, err := sdk.AccAddressFromBech32(msg.Seller)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{seller}
}

func (msg *MsgItemShipping) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgItemShipping) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Seller)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid seller address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateItem{}


func NewMsgItemResell(seller string, itemid string, shippingcost int64, discount int64, localpickup bool, shippingregion []string, note string) *MsgItemResell {
	return &MsgItemResell {
		Seller:  seller,
		Itemid:   itemid,
		Shippingcost:   shippingcost,
		Discount: discount,
		Localpickup:    localpickup,
		Shippingregion: shippingregion,
		Note: note,
		
	}
}

func (msg *MsgItemResell) Route() string {
	return RouterKey
}

func (msg *MsgItemResell) Type() string {
	return "ItemResell"
}

func (msg *MsgItemResell) GetSigners() []sdk.AccAddress {
	seller, err := sdk.AccAddressFromBech32(msg.Seller)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{seller}
}

func (msg *MsgItemResell) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgItemResell) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Seller)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid seller address (%s)", err)
	}
	if len(msg.Note) > 240 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "note too long")
	}
	if len(msg.Shippingregion) > 6  {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Region invalid")
	}

	for _, region := range msg.Shippingregion {
		if len(region) > 2{
			return sdkerrors.Wrap(sdkerrors.ErrMemoTooLarge, "Region too long")
		}
	}
	if msg.Shippingcost == 0 && msg.Localpickup != true {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Provide either shipping or localpickup")
	}

	return nil
}

var _ sdk.Msg = &MsgCreateItem{}