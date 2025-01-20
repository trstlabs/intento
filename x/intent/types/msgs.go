package types

import (
	fmt "fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/cosmos/gogoproto/proto"
)

var (
	_ sdk.Msg = &MsgRegisterAccount{}
	_ sdk.Msg = &MsgSubmitTx{}
	_ sdk.Msg = &MsgSubmitAction{}
	_ sdk.Msg = &MsgRegisterAccountAndSubmitAction{}
	_ sdk.Msg = &MsgUpdateAction{}
	_ sdk.Msg = &MsgCreateHostedAccount{}
	_ sdk.Msg = &MsgUpdateHostedAccount{}

	_ codectypes.UnpackInterfacesMessage = MsgSubmitTx{}
	_ codectypes.UnpackInterfacesMessage = MsgSubmitAction{}
	_ codectypes.UnpackInterfacesMessage = MsgRegisterAccountAndSubmitAction{}
	_ codectypes.UnpackInterfacesMessage = MsgUpdateAction{}
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
		return errorsmod.Wrap(ErrInvalidAddress, "missing sender address")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "failed to parse address: %s", msg.Owner)
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
			return nil, fmt.Errorf("cannot proto marshal %T", message)
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
		return fmt.Errorf("cannot execute an empty msg")
	}

	if msg.ConnectionId == "" {
		return fmt.Errorf("cannot execute an empty ConnectionId")
	}

	return nil
}

// NewMsgSubmitAction creates a new NewMsgSubmitAction instance
func NewMsgSubmitAction(owner, label string, sdkMsgs []sdk.Msg, connectionID string, hostConnectionID string, duration string, interval string, startAt uint64, feeFunds sdk.Coins, hostedAddress string, hostedFeeLimit sdk.Coin, configuration *ExecutionConfiguration, conditions *ExecutionConditions) (*MsgSubmitAction, error) {
	anys, err := PackTxMsgAnys(sdkMsgs)
	if err != nil {
		return nil, err
	}

	return &MsgSubmitAction{
		Owner:            owner,
		Label:            label,
		Msgs:             anys,
		Duration:         duration,
		Interval:         interval,
		StartAt:          startAt,
		FeeFunds:         feeFunds,
		Configuration:    configuration,
		ConnectionId:     connectionID,
		HostConnectionId: hostConnectionID,
		HostedConfig: &HostedConfig{HostedAddress: hostedAddress,
			FeeCoinLimit: hostedFeeLimit},
		Conditions: conditions,
	}, nil
}

// UnpackInterfaces implements codectypes.UnpackInterfacesMessage
func (msg MsgSubmitAction) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var sdkMsgs []sdk.Msg
	for _, message := range msg.Msgs {
		unpacker.UnpackAny(message, &sdkMsgs)
	}
	return nil
}

// GetTxMsgs fetches cached any messages
func (msg *MsgSubmitAction) GetTxMsgs() []sdk.Msg {
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
func (msg MsgSubmitAction) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// ValidateBasic implements sdk.Msg
func (msg MsgSubmitAction) ValidateBasic() error {
	if len(msg.Msgs) == 0 {
		return fmt.Errorf("msg.Msgs is empty, at least one message is required")
	}

	if len(msg.Msgs) >= 10 {
		return fmt.Errorf("cannot execute more than 9 messages")
	}
	if msg.Conditions != nil {
		err := checkConditions(*msg.Conditions, len(msg.Msgs))
		if err != nil {
			return err
		}
	}
	if len(msg.Label) > 50 {
		return fmt.Errorf("label must be shorter than 50 characters")
	}

	for _, message := range msg.GetTxMsgs() {
		// check if the msgs contain valid inputs
		m, ok := message.(sdk.HasValidateBasic)
		if !ok {
			continue
		}

		if err := m.ValidateBasic(); err != nil {
			return errorsmod.Wrapf(ErrUnknownRequest, "cannot validate action message: %s", err.Error())
		}

	}

	return nil
}

// NewMsgSend creates a new MsgSend instance
func NewMsgRegisterAccountAndSubmitAction(owner, label string, sdkMsgs []sdk.Msg, connectionID string, duration string, interval string, startAt uint64, feeFunds sdk.Coins, configuration *ExecutionConfiguration, version string) (*MsgRegisterAccountAndSubmitAction, error) {
	anys, err := PackTxMsgAnys(sdkMsgs)
	if err != nil {
		return nil, err
	}

	return &MsgRegisterAccountAndSubmitAction{
		Owner:         owner,
		Label:         label,
		ConnectionId:  connectionID,
		Msgs:          anys,
		Duration:      duration,
		Interval:      interval,
		StartAt:       startAt,
		FeeFunds:      feeFunds,
		Configuration: configuration,
		Version:       version,
	}, nil
}

// UnpackInterfaces implements codectypes.UnpackInterfacesMessage
func (msg MsgRegisterAccountAndSubmitAction) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var sdkMsgs []sdk.Msg
	for _, message := range msg.Msgs {
		unpacker.UnpackAny(message, &sdkMsgs)
	}
	return nil
}

// GetTxMsgs fetches cached any messages
func (msg *MsgRegisterAccountAndSubmitAction) GetTxMsgs() []sdk.Msg {
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
func (msg MsgRegisterAccountAndSubmitAction) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// ValidateBasic implements sdk.Msg
func (msg MsgRegisterAccountAndSubmitAction) ValidateBasic() error {
	if len(msg.Msgs) == 0 {
		return fmt.Errorf("msg.Msgs is empty, at least one message is required")
	}

	if msg.Conditions != nil {
		err := checkConditions(*msg.Conditions, len(msg.Msgs))
		if err != nil {
			return err
		}
	}

	for _, message := range msg.GetTxMsgs() {
		// check if the msgs contain valid inputs
		m, ok := message.(sdk.HasValidateBasic)
		if !ok {
			continue
		}

		if err := m.ValidateBasic(); err != nil {
			return errorsmod.Wrapf(ErrUnknownRequest, "cannot validate action message: %s", err.Error())
		}

	}
	if len(msg.Label) > 50 {
		return fmt.Errorf("label must be shorter than 50 characters")
	}
	return nil
}

func checkConditions(conditions ExecutionConditions, lenMsgMsgs int) error {

	if conditions.FeedbackLoops != nil {
		if len(conditions.FeedbackLoops) > 5 {
			return fmt.Errorf("cannot create more than 5 feedbackloops")
		}
		for _, feedbackLoop := range conditions.FeedbackLoops {
			if feedbackLoop.MsgKey == "" || feedbackLoop.ValueType == "" {
				return errorsmod.Wrapf(ErrUnknownRequest, "condition FeedbackLoops fields are not complete: %+v", feedbackLoop)
			}
			if int(feedbackLoop.ResponseIndex) >= lenMsgMsgs {
				return errorsmod.Wrapf(ErrInvalidRequest, "response index: %v must be shorter than length msgs array", feedbackLoop.ResponseIndex)
			}
			if feedbackLoop.ICQConfig != nil {
				if feedbackLoop.ICQConfig.TimeoutDuration == 0 || feedbackLoop.ICQConfig.ChainId == "" || feedbackLoop.ICQConfig.QueryKey == "" {
					return errorsmod.Wrapf(ErrUnknownRequest, "Loop ICQ Config fields are not complete: %+v", feedbackLoop.ICQConfig)
				}
			}
		}

	}
	if conditions.Comparisons != nil {
		for _, comparison := range conditions.Comparisons {
			if len(conditions.Comparisons) > 5 {
				return fmt.Errorf("cannot create more than 5 Comparisons")
			}
			if comparison.Operator <= 0 || comparison.ValueType == "" {
				return errorsmod.Wrapf(ErrUnknownRequest, "condition Comparision fields are not complete: %+v", conditions)
			}
			if comparison.ICQConfig != nil {
				if comparison.ICQConfig.TimeoutDuration == 0 || comparison.ICQConfig.ChainId == "" || comparison.ICQConfig.QueryKey == "" {
					return errorsmod.Wrapf(ErrUnknownRequest, "Comparison ICQ Config fields are not complete: %+v", comparison.ICQConfig)
				}
			}
		}

	}
	return nil
}

// NewMsgUpdateAction creates a new NewMsgUpdateAction instance
func NewMsgUpdateAction(owner string, id uint64, label string, sdkMsgs []sdk.Msg, connectionID string, endTime uint64, interval string, startAt uint64, feeFunds sdk.Coins, hostedAddress string, hostedFeeLimit sdk.Coin, configuration *ExecutionConfiguration, conditions *ExecutionConditions) (*MsgUpdateAction, error) {
	anys, err := PackTxMsgAnys(sdkMsgs)
	if err != nil {
		return nil, err
	}

	return &MsgUpdateAction{
		Owner:         owner,
		ID:            id,
		Label:         label,
		ConnectionId:  connectionID,
		Msgs:          anys,
		EndTime:       endTime,
		StartAt:       startAt,
		Interval:      interval,
		Configuration: configuration,
		FeeFunds:      feeFunds,
		HostedConfig: &HostedConfig{HostedAddress: hostedAddress,
			FeeCoinLimit: hostedFeeLimit},
		Conditions: conditions,
	}, nil
}

// UnpackInterfaces implements codectypes.UnpackInterfacesMessage
func (msg MsgUpdateAction) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var sdkMsgs []sdk.Msg
	for _, message := range msg.Msgs {
		unpacker.UnpackAny(message, &sdkMsgs)
	}
	return nil
}

// GetTxMsgs fetches cached any messages
func (msg *MsgUpdateAction) GetTxMsgs() []sdk.Msg {
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
func (msg MsgUpdateAction) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{owner}
}

// ValidateBasic implements sdk.Msg
func (msg MsgUpdateAction) ValidateBasic() error {
	if strings.TrimSpace(msg.Owner) == "" {
		return errorsmod.Wrap(ErrInvalidAddress, "missing creator address")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "failed to parse address: %s", msg.Owner)
	}

	if len(msg.Msgs) >= 10 {
		return fmt.Errorf("cannot execute more than 9 messages")
	}

	for _, message := range msg.GetTxMsgs() {
		// check if the msgs contain valid inputs
		m, ok := message.(sdk.HasValidateBasic)
		if !ok {
			continue
		}

		if err := m.ValidateBasic(); err != nil {
			return errorsmod.Wrapf(ErrUnknownRequest, "cannot validate action message: %s", err.Error())
		}

	}
	if len(msg.Label) > 50 {
		return fmt.Errorf("label must be shorter than 50 characters")
	}

	return nil
}

// NewMsgCreateHostedAccount creates a new MsgCreateHostedAccount instance
func NewMsgCreateHostedAccount(creator, connectionID, version string, feeFundsSupported sdk.Coins) *MsgCreateHostedAccount {
	return &MsgCreateHostedAccount{
		Creator:          creator,
		ConnectionId:     connectionID,
		Version:          version,
		FeeCoinsSuported: feeFundsSupported,
	}
}

// ValidateBasic implements sdk.Msg
func (msg MsgCreateHostedAccount) ValidateBasic() error {
	if strings.TrimSpace(msg.Creator) == "" {
		return errorsmod.Wrap(ErrInvalidAddress, "missing creator address")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "failed to parse address: %s", msg.Creator)
	}
	return nil
}

// GetSigners implements sdk.Msg
func (msg MsgCreateHostedAccount) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{accAddr}
}

// NewMsgUpdateHostedAccount creates a new NewMsgUpdateHostedAccount instance
func NewMsgUpdateHostedAccount(admin, hostedAddress, connectionID, hostConnectionID, newAdmin string, feeFundsSupported sdk.Coins) *MsgUpdateHostedAccount {

	return &MsgUpdateHostedAccount{
		Admin:            admin,
		HostedAddress:    hostedAddress,
		ConnectionId:     connectionID,
		HostConnectionId: hostConnectionID,
		HostFeeConfig:    &HostFeeConfig{FeeCoinsSuported: feeFundsSupported, Admin: newAdmin},
	}
}

// ValidateBasic implements sdk.Msg
func (msg MsgUpdateHostedAccount) ValidateBasic() error {
	if strings.TrimSpace(msg.Admin) == "" {
		return errorsmod.Wrap(ErrInvalidAddress, "missing creator address")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Admin); err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "failed to parse address: %s", msg.Admin)
	}
	return nil
}

// GetSigners implements sdk.Msg
func (msg MsgUpdateHostedAccount) GetSigners() []sdk.AccAddress {
	admin, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{admin}
}
