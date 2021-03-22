package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	//"cosmos/base/v1beta1/coin.proto"
)

var _ sdk.Msg = &MsgCreateEstimator{}

func NewMsgCreateEstimator(estimator string, estimation int64, itemid string, deposit int64, interested bool, comment string) *MsgCreateEstimator {
	return &MsgCreateEstimator{
		Estimator:  estimator,
		Estimation: estimation,
		//Estimatorestimationhash: estimatorestimationhash,
		Itemid:     itemid,
		Deposit:    deposit,
		Interested: interested,
		Comment:    comment,
	}
}

func (msg *MsgCreateEstimator) Route() string {
	return RouterKey
}

func (msg *MsgCreateEstimator) Type() string {
	return "CreateEstimator"
}

func (msg *MsgCreateEstimator) GetSigners() []sdk.AccAddress {
	estimator, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{estimator}
}

func (msg *MsgCreateEstimator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateEstimator) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if len(msg.Comment) > 50 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "comment too large")
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateEstimator{}

func NewMsgUpdateEstimator(estimator string, itemid string, interested bool) *MsgUpdateEstimator {
	return &MsgUpdateEstimator{
		Itemid:    itemid,
		Estimator: estimator,

		Interested: interested,
	}
}

func (msg *MsgUpdateEstimator) Route() string {
	return RouterKey
}

func (msg *MsgUpdateEstimator) Type() string {
	return "UpdateEstimator"
}

func (msg *MsgUpdateEstimator) GetSigners() []sdk.AccAddress {
	estimator, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{estimator}
}

func (msg *MsgUpdateEstimator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateEstimator) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateEstimator{}

func NewMsgDeleteEstimator(estimator string, itemid string) *MsgDeleteEstimator {
	return &MsgDeleteEstimator{
		Itemid:    itemid,
		Estimator: estimator,
	}
}
func (msg *MsgDeleteEstimator) Route() string {
	return RouterKey
}

func (msg *MsgDeleteEstimator) Type() string {
	return "DeleteEstimator"
}

func (msg *MsgDeleteEstimator) GetSigners() []sdk.AccAddress {
	estimator, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{estimator}
}

func (msg *MsgDeleteEstimator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteEstimator) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid estimator address (%s)", err)
	}
	return nil
}

func NewMsgCreateFlag(estimator string, flag bool, itemid string) *MsgCreateFlag {
	return &MsgCreateFlag{
		Itemid:    itemid,
		Flag:      flag,
		Estimator: estimator,
	}
}

func (msg *MsgCreateFlag) Route() string {
	return RouterKey
}

func (msg *MsgCreateFlag) Type() string {
	return "CreateFlag"
}

func (msg *MsgCreateFlag) GetSigners() []sdk.AccAddress {
	estimator, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{estimator}
}

func (msg *MsgCreateFlag) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateFlag) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCreateEstimator{}
