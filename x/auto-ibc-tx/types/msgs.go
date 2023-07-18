package types

import (
	fmt "fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	proto "github.com/cosmos/gogoproto/proto"
)

var (
	_ sdk.Msg = &MsgRegisterAccount{}
	_ sdk.Msg = &MsgSubmitTx{}
	_ sdk.Msg = &MsgSubmitAutoTx{}
	_ sdk.Msg = &MsgRegisterAccountAndSubmitAutoTx{}
	_ sdk.Msg = &MsgUpdateAutoTx{}

	_ codectypes.UnpackInterfacesMessage = MsgSubmitTx{}
	_ codectypes.UnpackInterfacesMessage = MsgSubmitAutoTx{}
	_ codectypes.UnpackInterfacesMessage = MsgRegisterAccountAndSubmitAutoTx{}
	_ codectypes.UnpackInterfacesMessage = MsgUpdateAutoTx{}
)

// NewMsgRegisterAccount creates a new MsgRegisterAccount instance
func NewMsgRegisterAccount(owner, connectionID string, version string) *MsgRegisterAccount {
	return &MsgRegisterAccount{
		Owner:        owner,
		ConnectionId: connectionID,
		Version:      version,
	}
}

// ValidateBasic implements sdk.Msg
func (msg MsgRegisterAccount) ValidateBasic() error {
	if strings.TrimSpace(msg.Owner) == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "failed to parse address: %s", msg.Owner)
	}
	return nil
}

// GetSigners implements sdk.Msg
func (msg MsgRegisterAccount) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{accAddr}
}

// NewMsgSend creates a new MsgSend instance
func NewMsgSubmitTx(owner string, sdkMsg sdk.Msg, connectionID string) (*MsgSubmitTx, error) {
	anys, err := PackTxMsgAnys([]sdk.Msg{sdkMsg})
	if err != nil {
		return nil, err
	}

	return &MsgSubmitTx{
		Owner:        owner,
		ConnectionId: connectionID,
		Msg:          anys[0],
	}, nil
}

// PackTxMsgAnys marshals the sdk.Msg payload to a protobuf Any type
func PackTxMsgAnys(sdkMsgs []sdk.Msg) ([]*codectypes.Any, error) {
	var anys []*codectypes.Any
	for _, message := range sdkMsgs {
		msg, ok := message.(proto.Message)
		if !ok {
			return nil, fmt.Errorf("can't proto marshal %T", message)
		}

		any, err := codectypes.NewAnyWithValue(msg)
		if err != nil {
			return nil, err
		}
		anys = append(anys, any)
	}
	return anys, nil
}

// UnpackInterfaces implements codectypes.UnpackInterfacesMessage
func (msg MsgSubmitTx) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var sdkMsg sdk.Msg

	return unpacker.UnpackAny(msg.Msg, &sdkMsg)
}

// GetTxMsg fetches the cached any message
func (msg *MsgSubmitTx) GetTxMsg() proto.Message {
	sdkMsg, ok := msg.Msg.GetCachedValue().(sdk.Msg)
	if !ok {
		return nil
	}

	return sdkMsg
}

// GetSigners implements sdk.Msg
func (msg MsgSubmitTx) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// ValidateBasic implements sdk.Msg
func (msg MsgSubmitTx) ValidateBasic() error {
	if len(msg.Msg.GetValue()) == 0 {
		return fmt.Errorf("can't execute an empty msg")
	}

	if msg.ConnectionId == "" {
		return fmt.Errorf("can't execute an empty ConnectionId")
	}

	return nil
}

// NewMsgSubmitAutoTx creates a new NewMsgSubmitAutoTx instance
func NewMsgSubmitAutoTx(owner, label string, sdkMsgs []sdk.Msg, connectionID string, duration string, interval string, startAt uint64, feeFunds sdk.Coins, dependsOn []uint64 /*  retries uint64 */) (*MsgSubmitAutoTx, error) {
	anys, err := PackTxMsgAnys(sdkMsgs)
	if err != nil {
		return nil, err
	}

	return &MsgSubmitAutoTx{
		Owner:          owner,
		Label:          label,
		ConnectionId:   connectionID,
		Msgs:           anys,
		Duration:       duration,
		Interval:       interval,
		StartAt:        startAt,
		FeeFunds:       feeFunds,
		DependsOnTxIds: dependsOn,
		//Retries:        retries,
	}, nil
}

// UnpackInterfaces implements codectypes.UnpackInterfacesMessage
func (msg MsgSubmitAutoTx) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var sdkMsgs []sdk.Msg
	for _, message := range msg.Msgs {
		unpacker.UnpackAny(message, &sdkMsgs)
	}
	return nil
}

// GetTxMsgs fetches cached any messages
func (msg *MsgSubmitAutoTx) GetTxMsgs() []sdk.Msg {
	var sdkMsgs []sdk.Msg
	for _, message := range msg.Msgs {
		sdkMsg, ok := message.GetCachedValue().(sdk.Msg)
		if !ok {
			return nil
		}
		sdkMsgs = append(sdkMsgs, sdkMsg)
	}

	return sdkMsgs
}

// GetSigners implements sdk.Msg
func (msg MsgSubmitAutoTx) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// ValidateBasic implements sdk.Msg
func (msg MsgSubmitAutoTx) ValidateBasic() error {
	if len(msg.Msgs) == 0 {
		return fmt.Errorf("msg.Msgs is empty, at least one message is required")
	}
	// if len(msg.Msgs[0].GetValue()) == 0 {
	// 	return fmt.Errorf("can't execute an empty msg")
	// }
	if len(msg.Msgs) >= 10 {
		return fmt.Errorf("can't execute more than 9 messages")
	}

	for _, message := range msg.GetTxMsgs() {
		// check if the msgs contain valid inputs
		err := message.ValidateBasic()
		if err != nil && !strings.Contains(err.Error(), "bech32") && !strings.Contains(err.Error(), "Bech32") {
			return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot validate autoTx message: %s", err.Error())
		}
	}

	return nil
}

// NewMsgSend creates a new MsgSend instance
func NewMsgRegisterAccountAndSubmitAutoTx(owner, label string, sdkMsgs []sdk.Msg, connectionID string, duration string, interval string, startAt uint64, feeFunds sdk.Coins, dependsOn []uint64 /*  retries uint64 */, version string) (*MsgRegisterAccountAndSubmitAutoTx, error) {
	anys, err := PackTxMsgAnys(sdkMsgs)
	if err != nil {
		return nil, err
	}

	return &MsgRegisterAccountAndSubmitAutoTx{
		Owner:          owner,
		Label:          label,
		ConnectionId:   connectionID,
		Msgs:           anys,
		Duration:       duration,
		Interval:       interval,
		StartAt:        startAt,
		FeeFunds:       feeFunds,
		DependsOnTxIds: dependsOn,
		Version:        version,
		//Retries:        retries,
	}, nil
}

// UnpackInterfaces implements codectypes.UnpackInterfacesMessage
func (msg MsgRegisterAccountAndSubmitAutoTx) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var sdkMsgs []sdk.Msg
	for _, message := range msg.Msgs {
		unpacker.UnpackAny(message, &sdkMsgs)
	}
	return nil
}

// GetTxMsgs fetches cached any messages
func (msg *MsgRegisterAccountAndSubmitAutoTx) GetTxMsgs() []sdk.Msg {
	var sdkMsgs []sdk.Msg
	for _, message := range msg.Msgs {
		sdkMsg, ok := message.GetCachedValue().(sdk.Msg)
		if !ok {
			return nil
		}
		sdkMsgs = append(sdkMsgs, sdkMsg)
	}

	return sdkMsgs
}

// GetSigners implements sdk.Msg
func (msg MsgRegisterAccountAndSubmitAutoTx) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// ValidateBasic implements sdk.Msg
func (msg MsgRegisterAccountAndSubmitAutoTx) ValidateBasic() error {
	if len(msg.Msgs) == 0 {
		return fmt.Errorf("msg.Msgs is empty, at least one message is required")
	}
	// if len(msg.Msgs[0].GetValue()) == 0 {
	// 	return fmt.Errorf("can't execute an empty msg")
	// }
	if len(msg.Msgs) >= 10 {
		return fmt.Errorf("can't execute more than 9 messages")
	}

	for _, message := range msg.GetTxMsgs() {
		// check if the msgs contain valid inputs
		err := message.ValidateBasic()
		if err != nil && !strings.Contains(err.Error(), "bech32") && !strings.Contains(err.Error(), "Bech32") {
			return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot validate autoTx message: %s", err.Error())
		}
	}

	return nil
}

// NewMsgUpdateAutoTx creates a new NewMsgUpdateAutoTx instance
func NewMsgUpdateAutoTx(owner string, txID uint64, label string, sdkMsgs []sdk.Msg, connectionID string, endTime uint64, interval string, startAt uint64, feeFunds sdk.Coins, dependsOn []uint64 /*  retries uint64 */) (*MsgUpdateAutoTx, error) {
	anys, err := PackTxMsgAnys(sdkMsgs)
	if err != nil {
		return nil, err
	}

	return &MsgUpdateAutoTx{
		Owner:          owner,
		TxId:           txID,
		Label:          label,
		ConnectionId:   connectionID,
		Msgs:           anys,
		EndTime:        endTime,
		StartAt:        startAt,
		Interval:       interval,
		DependsOnTxIds: dependsOn,
		FeeFunds:       feeFunds,
		//Retries:        retries,
	}, nil
}

// UnpackInterfaces implements codectypes.UnpackInterfacesMessage
func (msg MsgUpdateAutoTx) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var sdkMsgs []sdk.Msg
	for _, message := range msg.Msgs {
		unpacker.UnpackAny(message, &sdkMsgs)
	}
	return nil
}

// GetTxMsgs fetches cached any messages
func (msg *MsgUpdateAutoTx) GetTxMsgs() []sdk.Msg {
	var sdkMsgs []sdk.Msg
	for _, message := range msg.Msgs {
		sdkMsg, ok := message.GetCachedValue().(sdk.Msg)
		if !ok {
			return nil
		}
		sdkMsgs = append(sdkMsgs, sdkMsg)
	}

	return sdkMsgs
}

// GetSigners implements sdk.Msg
func (msg MsgUpdateAutoTx) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// ValidateBasic implements sdk.Msg
func (msg MsgUpdateAutoTx) ValidateBasic() error {
	if len(msg.Msgs) >= 10 {
		return fmt.Errorf("can't execute more than 9 messages")
	}

	for _, message := range msg.GetTxMsgs() {
		// check if the msgs contain valid inputs
		err := message.ValidateBasic()
		if err != nil && !strings.Contains(err.Error(), "bech32") && !strings.Contains(err.Error(), "Bech32") {
			return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot validate autoTx message: %s", err.Error())
		}
	}

	return nil
}
