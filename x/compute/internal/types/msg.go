package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (msg MsgStoreCode) Route() string {
	return RouterKey
}

func (msg MsgStoreCode) Type() string {
	return "store-code"
}

func (msg MsgStoreCode) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if err := validateWasmCode(msg.WASMByteCode); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "code bytes %s", err.Error())
	}

	if err := validateSourceURL(msg.Source); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "source %s", err.Error())
	}

	if err := validateBuilder(msg.Builder); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "builder %s", err.Error())
	}

	if len(msg.Title) > 100 {
		return sdkerrors.Wrap(sdkerrors.ErrMemoTooLarge, "Title length too long")
	}

	if len(msg.Description) > 1000 {
		return sdkerrors.Wrap(sdkerrors.ErrMemoTooLarge, "Description length too long")
	}
	/*
		if msg.InstantiatePermission != nil {
			if err := msg.InstantiatePermission.ValidateBasic(); err != nil {
				return sdkerrors.Wrap(err, "instantiate permission")
			}
		}
	*/
	return nil
}

func (msg MsgStoreCode) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgStoreCode) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg MsgInstantiateContract) Route() string {
	return RouterKey
}

func (msg MsgInstantiateContract) Type() string {
	return "instantiate"
}

func (msg MsgInstantiateContract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if msg.CodeID == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "code ID is required")
	}

	if err := validateContractId(msg.ContractId); err != nil {
		return err
	}

	if !msg.InitFunds.IsValid() {
		return sdkerrors.ErrInvalidCoins
	}

	/*
		if len(msg.Admin) != 0 {
			if err := sdk.VerifyAddressFormat(msg.Admin); err != nil {
				return err
			}
		}
		if !json.Valid(msg.InitMsg) {
			return sdkerrors.Wrap(ErrInvalid, "init msg json")
		}
	*/
	return nil
}

func (msg MsgInstantiateContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgInstantiateContract) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func validateContractId(contractId string) error {
	if contractId == "" {
		return sdkerrors.Wrap(ErrEmpty, "is required")
	}
	if len(contractId) > MaxContractIdSize {
		return sdkerrors.Wrap(ErrLimit, "cannot be longer than 128 characters")
	}
	return nil
}

func (msg MsgExecuteContract) Route() string {
	return RouterKey
}

func (msg MsgExecuteContract) Type() string {
	return "execute"
}

func (msg MsgExecuteContract) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if !msg.SentFunds.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "sentFunds")
	}

	/*
		if !json.Valid(msg.Msg) {
			return sdkerrors.Wrap(ErrInvalid, "msg json")
		}
	*/
	return nil
}

func (msg MsgExecuteContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgExecuteContract) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

// TypeMsgDiscardAutoMsg defines the type value for a MsgDiscardAutoMsg.
const TypeMsgDiscardAutoMsg = "msg_cancel_auto_msg"

var _ sdk.Msg = &MsgDiscardAutoMsg{}

// NewMsgDiscardAutoMsg returns a reference to a new MsgDiscardAutoMsg.
//nolint:interfacer
func NewMsgDiscardAutoMsg(sender sdk.AccAddress, contractAddr sdk.AccAddress) *MsgDiscardAutoMsg {
	return &MsgDiscardAutoMsg{
		Sender:          sender,
		ContractAddress: contractAddr,
	}
}

// Route returns the message route for a MsgDiscardAutoMsg.
func (msg MsgDiscardAutoMsg) Route() string { return RouterKey }

// Type returns the message type for a MsgDiscardAutoMsg.
func (msg MsgDiscardAutoMsg) Type() string { return TypeMsgDiscardAutoMsg }

// ValidateBasic Implements Msg.
func (msg MsgDiscardAutoMsg) ValidateBasic() error {
	/*if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid 'from' address: %s", err)
	}*/

	return nil
}

// GetSignBytes returns the bytes all expected signers must sign over for a
// MsgDiscardAutoMsg.
func (msg MsgDiscardAutoMsg) GetSignBytes() []byte {
	return sdk.MustSortJSON(amino.MustMarshalJSON(&msg))
}

// GetSigners returns the expected signers for a MsgDiscardAutoMsg.
func (msg MsgDiscardAutoMsg) GetSigners() []sdk.AccAddress {
	/*(addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}*/
	return []sdk.AccAddress{msg.Sender}
}

/*
type MsgMigrateContract struct {
	Sender     sdk.AccAddress  `json:"sender" yaml:"sender"`
	Contract   sdk.AccAddress  `json:"contract" yaml:"contract"`
	CodeID     uint64          `json:"code_id" yaml:"code_id"`
	MigrateMsg json.RawMessage `json:"msg" yaml:"msg"`
}

func (msg MsgMigrateContract) Route() string {
	return RouterKey
}

func (msg MsgMigrateContract) Type() string {
	return "migrate"
}

func (msg MsgMigrateContract) ValidateBasic() error {
	if msg.CodeID == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "code_id is required")
	}
	if err := sdk.VerifyAddressFormat(msg.Sender); err != nil {
		return sdkerrors.Wrap(err, "sender")
	}
	if err := sdk.VerifyAddressFormat(msg.Contract); err != nil {
		return sdkerrors.Wrap(err, "contract")
	}
	if !json.Valid(msg.MigrateMsg) {
		return sdkerrors.Wrap(ErrInvalid, "migrate msg json")
	}

	return nil
}

func (msg MsgMigrateContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgMigrateContract) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender} }

type MsgUpdateAdmin struct {
	Sender   sdk.AccAddress `json:"sender" yaml:"sender"`
	NewAdmin sdk.AccAddress `json:"new_admin" yaml:"new_admin"`
	Contract sdk.AccAddress `json:"contract" yaml:"contract"`
}

func (msg MsgUpdateAdmin) Route() string {
	return RouterKey
}

func (msg MsgUpdateAdmin) Type() string {
	return "update-contract-admin"
}

func (msg MsgUpdateAdmin) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(msg.Sender); err != nil {
		return sdkerrors.Wrap(err, "sender")
	}
	if err := sdk.VerifyAddressFormat(msg.Contract); err != nil {
		return sdkerrors.Wrap(err, "contract")
	}
	if err := sdk.VerifyAddressFormat(msg.NewAdmin); err != nil {
		return sdkerrors.Wrap(err, "new admin")
	}
	if msg.Sender.Equals(msg.NewAdmin) {
		return sdkerrors.Wrap(ErrInvalidMsg, "new admin is the same as the old")
	}
	return nil
}

func (msg MsgUpdateAdmin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgUpdateAdmin) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender} }

type MsgClearAdmin struct {
	Sender   sdk.AccAddress `json:"sender" yaml:"sender"`
	Contract sdk.AccAddress `json:"contract" yaml:"contract"`
}

func (msg MsgClearAdmin) Route() string {
	return RouterKey
}

func (msg MsgClearAdmin) Type() string {
	return "clear-contract-admin"
}

func (msg MsgClearAdmin) ValidateBasic() error {
	if err := sdk.VerifyAddressFormat(msg.Sender); err != nil {
		return sdkerrors.Wrap(err, "sender")
	}
	if err := sdk.VerifyAddressFormat(msg.Contract); err != nil {
		return sdkerrors.Wrap(err, "contract")
	}
	return nil
}

func (msg MsgClearAdmin) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgClearAdmin) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender} }
*/
