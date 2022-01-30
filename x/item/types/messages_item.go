package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	//"cosmos/base/v1beta1/coin.proto"
)

var _ sdk.Msg = &MsgCreateItem{}

func NewMsgCreateItem(creator string, title string, description string, shippingCost int64, localPickup string, estimationCount int64, tags []string, condition int64, shippingRegion []string, depositAmount int64, initMsg []byte, autoMsg []byte, photos []string, tokenURI string) *MsgCreateItem {
	return &MsgCreateItem{

		Creator:         creator,
		Title:           title,
		Description:     description,
		ShippingCost:    shippingCost,
		LocalPickup:     localPickup,
		EstimationCount: estimationCount,
		Tags:            tags,
		Condition:       condition,
		ShippingRegion:  shippingRegion,
		DepositAmount:   depositAmount,
		InitMsg:         initMsg,
		AutoMsg:         autoMsg,
		Photos:          photos,
		TokenUri:        tokenURI,
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
		if len(tags) > 24 {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "a tag was too long")
		}
	}

	if len(msg.ShippingRegion) > 9 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Region list too long")
	}

	for _, region := range msg.ShippingRegion {
		if len(region) > 2 {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "A Region cannot be longer than 2")
		}
	}
	if len(msg.Title) > 100 {
		return sdkerrors.Wrap(sdkerrors.ErrMemoTooLarge, "Title length too long")
	}

	if len(msg.Description) > 1000 {
		return sdkerrors.Wrap(sdkerrors.ErrMemoTooLarge, "Description length too long")
	}

	if len(msg.LocalPickup) > 48 {
		return sdkerrors.Wrap(sdkerrors.ErrMemoTooLarge, "Local pickup location too long")
	}

	if msg.Condition > 5 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid item condition")
	}
	if msg.EstimationCount > 24 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid estimation count")
	}

	if len(msg.Photos) > 9 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "too many photos")
	}
	for _, photo := range msg.Photos {
		if len(photo) > 300 {
			return sdkerrors.Wrap(sdkerrors.ErrMemoTooLarge, "photo url too long")
		}
	}
	if msg.TokenUri == "" && msg.Photos == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "provide either photos or a token URI")
	}
	return nil
}

var _ sdk.Msg = &MsgCreateItem{}

func NewMsgDeleteItem(seller string, id uint64) *MsgDeleteItem {
	return &MsgDeleteItem{
		Id:     id,
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

func NewMsgRevealEstimation(creator string, itemid uint64, revealMsg []byte) *MsgRevealEstimation {
	return &MsgRevealEstimation{

		Creator:   creator,
		Itemid:    itemid,
		RevealMsg: revealMsg,
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

func NewMsgItemTransferable(seller string, transferable []byte, itemid uint64) *MsgItemTransferable {
	return &MsgItemTransferable{

		Seller:          seller,
		TransferableMsg: transferable,
		Itemid:          itemid,
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

func NewMsgItemShipping(seller string, tracking bool, itemid uint64) *MsgItemShipping {
	return &MsgItemShipping{

		Seller:   seller,
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

func NewMsgItemResell(seller string, itemid uint64, shippingCost int64, discount int64, localPickup string, shippingRegion []string, note string) *MsgItemResell {
	return &MsgItemResell{
		Seller:         seller,
		Itemid:         itemid,
		ShippingCost:   shippingCost,
		Discount:       discount,
		LocalPickup:    localPickup,
		ShippingRegion: shippingRegion,
		Note:           note,
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
	if len(msg.ShippingRegion) > 6 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Regions too long")
	}

	for _, region := range msg.ShippingRegion {
		if len(region) > 2 {
			return sdkerrors.Wrap(sdkerrors.ErrMemoTooLarge, "Region too long")
		}
	}
	if msg.ShippingCost == 0 && msg.LocalPickup == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Provide either shipping or localPickup")
	}

	if len(msg.LocalPickup) > 25 {
		return sdkerrors.Wrap(sdkerrors.ErrMemoTooLarge, "Local pickup too long")
	}

	return nil
}

func NewMsgTokenizeItem(sender string, id uint64) *MsgTokenizeItem {
	return &MsgTokenizeItem{
		Id:     id,
		Sender: sender,
	}
}
func (msg *MsgTokenizeItem) Route() string {
	return RouterKey
}

func (msg *MsgTokenizeItem) Type() string {
	return "DeleteItem"
}

func (msg *MsgTokenizeItem) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgTokenizeItem) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTokenizeItem) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

func NewMsgUnTokenizeItem(sender string, id uint64) *MsgUnTokenizeItem {
	return &MsgUnTokenizeItem{
		Id:     id,
		Sender: sender,
	}
}
func (msg *MsgUnTokenizeItem) Route() string {
	return RouterKey
}

func (msg *MsgUnTokenizeItem) Type() string {
	return "DeleteItem"
}

func (msg *MsgUnTokenizeItem) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgUnTokenizeItem) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUnTokenizeItem) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}
