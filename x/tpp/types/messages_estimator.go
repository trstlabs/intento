package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	//"cosmos/base/v1beta1/coin.proto"
)

var _ sdk.Msg = &MsgCreateEstimation{}

func NewMsgCreateEstimation(estimator string, estimatemsg []byte, itemid uint64, deposit int64, interested bool) *MsgCreateEstimation {
	fmt.Printf("new item msg: %X\n", estimator)
	return &MsgCreateEstimation{
		Estimator:   estimator,
		Estimatemsg: estimatemsg,
		//Estimatorestimationhash: estimatorestimationhash,
		Itemid:     itemid,
		Deposit:    deposit,
		Interested: interested,
		//Comment:    comment,
	}
}

func (msg *MsgCreateEstimation) Route() string {
	return RouterKey
}

func (msg *MsgCreateEstimation) Type() string {
	return "CreateEstimation"
}

func (msg *MsgCreateEstimation) GetSigners() []sdk.AccAddress {
	estimator, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{estimator}
}

func (msg *MsgCreateEstimation) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateEstimation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid estimator address (%s)", err)
	}
	/*if len(msg.Comment) > 100 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "comment too large")
	}*/
	return nil
}

var _ sdk.Msg = &MsgUpdateLike{}

func NewMsgUpdateLike(estimator string, itemid uint64, interested bool) *MsgUpdateLike {
	return &MsgUpdateLike{
		Itemid:    itemid,
		Estimator: estimator,

		Interested: interested,
	}
}

func (msg *MsgUpdateLike) Route() string {
	return RouterKey
}

func (msg *MsgUpdateLike) Type() string {
	return "UpdateLike"
}

func (msg *MsgUpdateLike) GetSigners() []sdk.AccAddress {
	estimator, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{estimator}
}

func (msg *MsgUpdateLike) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateLike) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid estimator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateEstimation{}

func NewMsgDeleteEstimation(estimator string, itemid uint64, deletemsg []byte) *MsgDeleteEstimation {
	return &MsgDeleteEstimation{
		Itemid:    itemid,
		Estimator: estimator,
		Deletemsg: deletemsg,
	}
}
func (msg *MsgDeleteEstimation) Route() string {
	return RouterKey
}

func (msg *MsgDeleteEstimation) Type() string {
	return "DeleteEstimation"
}

func (msg *MsgDeleteEstimation) GetSigners() []sdk.AccAddress {
	estimator, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{estimator}
}

func (msg *MsgDeleteEstimation) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteEstimation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid estimator address (%s)", err)
	}
	return nil
}

func NewMsgFlagItem(estimator string, itemid uint64, flagmsg []byte) *MsgFlagItem {
	return &MsgFlagItem{
		Itemid: itemid,

		Estimator: estimator,
		Flagmsg:   flagmsg,
	}
}

func (msg *MsgFlagItem) Route() string {
	return RouterKey
}

func (msg *MsgFlagItem) Type() string {
	return "FlagItem"
}

func (msg *MsgFlagItem) GetSigners() []sdk.AccAddress {
	estimator, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{estimator}
}

func (msg *MsgFlagItem) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgFlagItem) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateEstimation{}
