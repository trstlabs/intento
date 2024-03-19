package msg_registry

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
// If the validator address is not same as delegator's, then the validator must
// sign the msg as well.

func (msg MsgExecuteContract) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg MsgExecuteContract) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrap(err, "sender")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Contract); err != nil {
		return errorsmod.Wrap(err, "contract")
	}

	if !msg.Funds.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "sentFunds")
	}

	return nil
}

func (msg MsgInstantiateContract) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrap(err, "sender")
	}

	if msg.CodeID == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "code id is required")
	}

	// if err := ValidateLabel(msg.Label); err != nil {
	// 	return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "label is required")
	// }

	if !msg.Funds.IsValid() {
		return sdkerrors.ErrInvalidCoins
	}

	if len(msg.Admin) != 0 {
		if _, err := sdk.AccAddressFromBech32(msg.Admin); err != nil {
			return errorsmod.Wrap(err, "admin")
		}
	}
	// if err := msg.Msg.ValidateBasic(); err != nil {
	// 	return errorsmod.Wrap(err, "payload msg")
	// }
	return nil
}

func (msg MsgInstantiateContract) GetSigners() []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil { // should never happen as valid basic rejects invalid addresses
		panic(err.Error())
	}
	return []sdk.AccAddress{senderAddr}
}

func (msg MsgInstantiateContract2) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrap(err, "sender")
	}

	if msg.CodeID == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "code id is required")
	}

	// if err := ValidateLabel(msg.Label); err != nil {
	// 	return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "label is required")
	// }

	if !msg.Funds.IsValid() {
		return sdkerrors.ErrInvalidCoins
	}

	if len(msg.Admin) != 0 {
		if _, err := sdk.AccAddressFromBech32(msg.Admin); err != nil {
			return errorsmod.Wrap(err, "admin")
		}
	}
	// if err := msg.Msg.ValidateBasic(); err != nil {
	// 	return errorsmod.Wrap(err, "payload msg")
	// }
	return nil
}

func (msg MsgInstantiateContract2) GetSigners() []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil { // should never happen as valid basic rejects invalid addresses
		panic(err.Error())
	}
	return []sdk.AccAddress{senderAddr}
}
