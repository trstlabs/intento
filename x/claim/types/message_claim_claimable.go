package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// TypeMsgClaimClaimable defines the type value for a MsgClaimClaimable.
const TypeMsgClaimClaimable = "msg_claim_claimable"

var _ sdk.Msg = &MsgClaimClaimable{}

// NewMsgClaimClaimable returns a reference to a new MsgClaimClaimable.
//nolint:interfacer
func NewMsgClaimClaimable(sender sdk.AccAddress) *MsgClaimClaimable {
	return &MsgClaimClaimable{
		Sender: sender.String(),
	}
}

// Route returns the message route for a MsgClaimClaimable.
func (msg MsgClaimClaimable) Route() string { return RouterKey }

// Type returns the message type for a MsgClaimClaimable.
func (msg MsgClaimClaimable) Type() string { return TypeMsgClaimClaimable }

// ValidateBasic Implements Msg.
func (msg MsgClaimClaimable) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid 'from' address: %s", err)
	}

	return nil
}

// GetSignBytes returns the bytes all expected signers must sign over for a
// MsgClaimClaimable.
func (msg MsgClaimClaimable) GetSignBytes() []byte {
	return sdk.MustSortJSON(amino.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for a MsgClaimClaimable.
func (msg MsgClaimClaimable) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}
