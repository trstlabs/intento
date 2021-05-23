package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	//"cosmos/base/v1beta1/coin.proto"
)

var _ sdk.Msg = &MsgCreateEstimation{}

func NewMsgCreateEstimation(estimator string, estimation int64, itemid uint64, deposit int64, interested bool, comment string) *MsgCreateEstimation {
	return &MsgCreateEstimation{
		Estimator:  estimator,
		Estimation: estimation,
		//Estimatorestimationhash: estimatorestimationhash,
		Itemid:     itemid,
		Deposit:    deposit,
		Interested: interested,
		Comment:    comment,
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
	if len(msg.Comment) > 100 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "comment too large")
	}
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

func NewMsgDeleteEstimation(estimator string, itemid uint64) *MsgDeleteEstimation {
	return &MsgDeleteEstimation{
		Itemid:    itemid,
		Estimator: estimator,
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

func NewMsgFlagItem(estimator string, itemid uint64) *MsgFlagItem {
	return &MsgFlagItem{
		Itemid: itemid,

		Estimator: estimator,
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
