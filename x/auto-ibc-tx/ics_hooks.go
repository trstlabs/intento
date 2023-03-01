package autoibctx

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	//codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v3/modules/core/exported"
	"github.com/trstlabs/trst/x/auto-ibc-tx/keeper"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

var _ porttypes.Middleware = &IBCMiddleware{}

type IBCMiddleware struct {
	app         porttypes.IBCModule
	ics4Wrapper porttypes.ICS4Wrapper
	keeper      keeper.Keeper //add a keeper for stateful middleware
	registry    codectypes.InterfaceRegistry
	//wrapper  ICS4Wrapper
}

// IBCMiddleware creates a new IBCMiddleware given the associated keeper and underlying application
func NewIBCMiddleware(app porttypes.IBCModule, k keeper.Keeper, registry codectypes.InterfaceRegistry, wrapper porttypes.ICS4Wrapper) IBCMiddleware {
	return IBCMiddleware{
		app:         app,
		keeper:      k,
		registry:    registry,
		ics4Wrapper: wrapper,
	}
}

/*
type ICS4Wrapper interface {
	SendPacket(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet ibcexported.PacketI) error
	WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet ibcexported.PacketI, ack ibcexported.Acknowledgement) error
	GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool)
} */

// OnChanOpenInit implements the IBCMiddleware interface
func (im IBCMiddleware) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	channelCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) /* string,  */ error {

	/* finalVersion, */
	err := im.app.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version)

	return /*  version, */ err
}

// OnChanOpenTry implements the IBCMiddleware interface
func (im IBCMiddleware) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	channelCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {

	version, err := im.app.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, counterpartyVersion)

	return version, err
}

// OnChanOpenAck implements the IBCMiddleware interface
func (im IBCMiddleware) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {

	err := im.app.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)

	return err
}

// OnChanOpenConfirm implements the IBCMiddleware interface
func (im IBCMiddleware) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {

	err := im.app.OnChanOpenConfirm(ctx, portID, channelID)

	return err
}

// OnChanCloseInit implements the IBCMiddleware interface
func (im IBCMiddleware) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {

	err := im.app.OnChanCloseInit(ctx, portID, channelID)

	return err
}

// OnChanCloseConfirm implements the IBCMiddleware interface
func (im IBCMiddleware) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {

	err := im.app.OnChanCloseConfirm(ctx, portID, channelID)

	return err
}

// OnRecvPacket implements the IBCMiddleware interface
func (im IBCMiddleware) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	ack := OnRecvPacketOverride(im, ctx, packet, relayer)

	//ack := im.app.OnRecvPacket(ctx, packet, relayer)

	return ack
}

// OnAcknowledgementPacket implements the IBCMiddleware interface
func (im IBCMiddleware) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {

	err := im.app.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)

	return err
}

// OnTimeoutPacket implements the IBCMiddleware interface
func (im IBCMiddleware) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {

	err := im.app.OnTimeoutPacket(ctx, packet, relayer)

	return err
}

// SendPacket implements the ICS4 Wrapper interface
func (im IBCMiddleware) SendPacket(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	packet ibcexported.PacketI,
) error {
	return im.ics4Wrapper.SendPacket(ctx, chanCap, packet)
}

// WriteAcknowledgement implements the ICS4 Wrapper interface
func (im IBCMiddleware) WriteAcknowledgement(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	packet ibcexported.PacketI,
	ack ibcexported.Acknowledgement,
) error {
	return im.ics4Wrapper.WriteAcknowledgement(ctx, chanCap, packet, ack)
}

/*
func (im IBCMiddleware) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	return im.ics4Wrapper.GetAppVersion(ctx, portID, channelID)
} */

type ContractAck struct {
	ContractResult []byte `json:"contract_result"`
	IbcAck         []byte `json:"ibc_ack"`
}

type Ics20Hooks struct {
	autoTxKeeper        *keeper.Keeper
	bech32PrefixAccAddr string
}

func NewIcs20Hooks(autoTxKeeper *keeper.Keeper, bech32PrefixAccAddr string) Ics20Hooks {
	return Ics20Hooks{
		autoTxKeeper:        autoTxKeeper,
		bech32PrefixAccAddr: bech32PrefixAccAddr,
	}
}

func /* (h Ics20Hooks) */ OnRecvPacketOverride(im IBCMiddleware, ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) ibcexported.Acknowledgement {
	/* if im.keeper != nil {
		// Not configured
		return im.app.OnRecvPacket(ctx, packet, relayer)
	} */
	isIcs20, data := isIcs20Packet(packet)
	if !isIcs20 {
		return im.app.OnRecvPacket(ctx, packet, relayer)
	}

	// Validate the memo
	isAutoTxRouted, ownerAddr, msgBytes, label, connectionID, duration, interval, startAt, registeICA, err := ValidateAndParseMemo(data.GetMemo(), data.Receiver)
	if !isAutoTxRouted {
		return im.app.OnRecvPacket(ctx, packet, relayer)
	}
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(err.Error())
	}
	if msgBytes == nil || ownerAddr == nil { // This should never happen
		return channeltypes.NewErrorAcknowledgement(err.Error())
	}
	var txMsgsAny codectypes.Any
	cdc := codec.NewProtoCodec(im.registry)

	if err := cdc.UnmarshalJSON(msgBytes, &txMsgsAny); err != nil {
		return channeltypes.NewErrorAcknowledgement(fmt.Sprintf(types.ErrBadMetadataFormatMsg, "error unmarshalling sdk msg file"))

	}

	// Calculate the receiver / contract caller based on the packet's channel and sender
	channel := packet.GetDestChannel()
	sender := data.GetSender()
	senderBech32, err := deriveIntermediateSender(channel, sender, "trust")
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(fmt.Sprintf("cannot convert sender address %s/%s to bech32: %s", channel, sender, err.Error()))
	}

	// The funds sent on this packet need to be transferred to the intermediary account for the sender.
	// For this, we override the ICS20 packet's Receiver (essentially hijacking the funds to this new address)
	// and execute the underlying OnRecvPacket() call (which should eventually land on the transfer app's
	// relay.go and send the sunds to the intermediary account.
	//
	// If that succeeds, we make the contract call
	data.Receiver = senderBech32
	bz, err := json.Marshal(data)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(sdkerrors.Wrapf(types.ErrMarshaling, err.Error()).Error())
	}
	packet.Data = bz

	// Execute the receive
	ack := im.app.OnRecvPacket(ctx, packet, relayer)
	if !ack.Success() {
		return ack
	}

	amount, ok := sdk.NewIntFromString(data.GetAmount())
	if !ok {
		// This should never happen, as it should've been caught in the underlaying call to OnRecvPacket,
		// but returning here for completeness
		return channeltypes.NewErrorAcknowledgement(sdkerrors.Wrapf(types.ErrInvalidPacket, "Amount is not an int").Error())
	}

	// The packet's denom is the denom in the sender chain. This needs to be converted to the local denom.
	denom := MustExtractDenomFromPacketOnRecv(packet)
	funds := sdk.NewCoins(sdk.NewCoin(denom, amount))

	// Build the message
	if registeICA {
		msg := types.MsgRegisterAccountAndSubmitAutoTx{
			Owner:        ownerAddr.String(),
			Msgs:         []*codectypes.Any{&txMsgsAny},
			FeeFunds:     funds,
			Label:        label,
			ConnectionId: connectionID,
			Duration:     duration,
			Interval:     interval,
			StartAt:      startAt,
		}
		response, err := /* h. */ registerAndSubmitTx(im.keeper, ctx, &msg)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(sdkerrors.Wrapf(types.ErrIcs20Error, err.Error()).Error())
		}
		bz, err = json.Marshal(response)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(sdkerrors.Wrapf(types.ErrBadResponse, err.Error()).Error())
		}

		return channeltypes.NewResultAcknowledgement(bz)
	} else {
		msg := types.MsgSubmitAutoTx{
			//Owner:    senderBech32,
			Owner:        ownerAddr.String(),
			Msgs:         []*codectypes.Any{&txMsgsAny},
			FeeFunds:     funds,
			Label:        label,
			ConnectionId: connectionID,
			Duration:     duration,
			Interval:     interval,
			StartAt:      startAt,
		}
		response, err := submitTx(im.keeper, ctx, &msg)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(sdkerrors.Wrapf(types.ErrIcs20Error, err.Error()).Error())
		}
		bz, err = json.Marshal(response)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(sdkerrors.Wrapf(types.ErrBadResponse, err.Error()).Error())
		}

		return channeltypes.NewResultAcknowledgement(bz)
	}

}

func /* (h Ics20Hooks) */ registerAndSubmitTx(k keeper.Keeper, ctx sdk.Context, autoTxMsg *types.MsgRegisterAccountAndSubmitAutoTx) (*types.MsgRegisterAccountAndSubmitAutoTxResponse, error) {
	if err := autoTxMsg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf(types.ErrBadAutoTxMsg, err.Error())
	}
	ics20MsgServer := keeper.NewMsgServerImpl(k)
	return ics20MsgServer.RegisterAccountAndSubmitAutoTx(sdk.WrapSDKContext(ctx), autoTxMsg)
}

func /* (h Ics20Hooks) */ submitTx(k keeper.Keeper, ctx sdk.Context, autoTxMsg *types.MsgSubmitAutoTx) (*types.MsgSubmitAutoTxResponse, error) {
	if err := autoTxMsg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf(types.ErrBadAutoTxMsg, err.Error())
	}
	ics20MsgServer := keeper.NewMsgServerImpl(k)
	return ics20MsgServer.SubmitAutoTx(sdk.WrapSDKContext(ctx), autoTxMsg)
}

func isIcs20Packet(packet channeltypes.Packet) (isIcs20 bool, ics20data transfertypes.FungibleTokenPacketData) {
	var data transfertypes.FungibleTokenPacketData
	if err := json.Unmarshal(packet.GetData(), &data); err != nil {
		return false, data
	}
	return true, data
}

// jsonStringHasKey parses the memo as a json object and checks if it contains the key.
func jsonStringHasKey(memo, key string) (found bool, jsonObject map[string]interface{}) {
	jsonObject = make(map[string]interface{})

	// If there is no memo, the packet was either sent with an earlier version of IBC, or the memo was
	// intentionally left blank. Nothing to do here. Ignore the packet and pass it down the stack.
	if len(memo) == 0 {
		return false, jsonObject
	}

	// the jsonObject must be a valid JSON object
	err := json.Unmarshal([]byte(memo), &jsonObject)
	if err != nil {
		return false, jsonObject
	}

	// If the key doesn't exist, there's nothing to do on this hook. Continue by passing the packet
	// down the stack
	_, ok := jsonObject[key]
	if !ok {
		return false, jsonObject
	}

	return true, jsonObject
}

func ValidateAndParseMemo(memo string, receiver string) (isAutoTxRouted bool, ownerAddr sdk.AccAddress, msgBytes []byte, label, connectionID, duration, interval string, startAt uint64, registerICA bool, err error) {
	isAutoTxRouted, metadata := jsonStringHasKey(memo, "auto_tx")
	if !isAutoTxRouted {
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false, nil
	}

	ics20Raw := metadata["auto_tx"]

	// Make sure the ics20 key is a map. If it isn't, ignore this packet
	ics20, ok := ics20Raw.(map[string]interface{})
	if !ok {
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, "auto_tx metadata is not a valid JSON map object")
	}

	// Get the owner
	owner, ok := ics20["owner"].(string)
	if !ok {
		// The tokens will be returned
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `auto_tx["owner"]`)
	}

	ownerAddr, err = sdk.AccAddressFromBech32(owner)
	if err != nil {
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `ics20["owner"] is not a valid bech32 address`)
	}

	// The owner and the receiver should be the same for the packet to be valid
	if owner != receiver {
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `ics20["owner"] should be the same as the receiver of the packet`)
	}

	// Ensure the message key is provided
	if ics20["msg"] == nil {
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `ics20["msg"]`)
	}

	// Make sure the msg key is a map. If it isn't, return an error
	_, ok = ics20["msg"].(map[string]interface{})
	if !ok {
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `ics20["msg"] is not a map object`)
	}

	// Get the label
	label, ok = ics20["label"].(string)
	if !ok {
		// The tokens will be returned
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `auto_tx["label"]`)
	}
	// Get the portID
	connectionID, ok = ics20["portID"].(string)
	if !ok {
		// The tokens will be returned
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `auto_tx["portID"]`)
	}
	// Get the duration
	duration, ok = ics20["duration"].(string)
	if !ok {
		// The tokens will be returned
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `auto_tx["duration"]`)
	}
	// Get the interval
	interval, ok = ics20["interval"].(string)
	if !ok {
		// The tokens will be returned
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `auto_tx["interval"]`)
	}
	// Get the label
	startAtString, ok := ics20["startAt"].(string)
	if !ok {
		// The tokens will be returned
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `auto_tx["startAt"]`)
	}
	startAt, err = strconv.ParseUint(startAtString, 10, 64)
	if err != nil {
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `auto_tx["startAt"]`)
	}

	// see if register_ica
	registerICAString, ok := ics20["register_ica"].(string)
	if ok && registerICAString == "true" {
		registerICA = true
	}

	// Get the message string by serializing the map
	msgBytes, err = json.Marshal(ics20["msg"])
	if err != nil {
		// The tokens will be returned
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, err.Error())
	}

	return isAutoTxRouted, ownerAddr, msgBytes, label, connectionID, duration, interval, startAt, registerICA, nil
}

/*
func (h Ics20Hooks) SendPacketOverride(i ICS4Middleware, ctx sdk.Context, chanCap *capabilitytypes.Capability, packet ibcexported.PacketI) error {
	concretePacket, ok := packet.(channeltypes.Packet)
	if !ok {
		return i.channel.SendPacket(ctx, chanCap, packet) // continue
	}

	isIcs20, data := isIcs20Packet(concretePacket)
	if !isIcs20 {
		return i.channel.SendPacket(ctx, chanCap, packet) // continue
	}

	isCallbackRouted, metadata := jsonStringHasKey(data.GetMemo(), types.IBCCallbackKey)
	if !isCallbackRouted {
		return i.channel.SendPacket(ctx, chanCap, packet) // continue
	}

	// We remove the callback metadata from the memo as it has already been processed.

	// If the only available key in the memo is the callback, we should remove the memo
	// from the data completely so the packet is sent without it.
	// This way receiver chains that are on old versions of IBC will be able to process the packet

	callbackRaw := metadata[types.IBCCallbackKey] // This will be used later.
	delete(metadata, types.IBCCallbackKey)
	bzMetadata, err := json.Marshal(metadata)
	if err != nil {
		return sdkerrors.Wrap(err, "Send packet with callback error")
	}
	stringMetadata := string(bzMetadata)
	if stringMetadata == "{}" {
		data.Memo = ""
	} else {
		data.Memo = stringMetadata
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return sdkerrors.Wrap(err, "Send packet with callback error")
	}

	packetWithoutCallbackMemo := channeltypes.Packet{
		Sequence:           concretePacket.Sequence,
		SourcePort:         concretePacket.SourcePort,
		SourceChannel:      concretePacket.SourceChannel,
		DestinationPort:    concretePacket.DestinationPort,
		DestinationChannel: concretePacket.DestinationChannel,
		Data:               dataBytes,
		TimeoutTimestamp:   concretePacket.TimeoutTimestamp,
		TimeoutHeight:      concretePacket.TimeoutHeight,
	}

	err = i.channel.SendPacket(ctx, chanCap, packetWithoutCallbackMemo)
	if err != nil {
		return err
	}

	// Make sure the callback contract is a string and a valid bech32 addr. If it isn't, ignore this packet
	contract, ok := callbackRaw.(string)
	if !ok {
		return nil
	}
	_, err = sdk.AccAddressFromBech32(contract)
	if err != nil {
		return nil
	}

	h.autoTxKeeper.StorePacketCallback(ctx, packet.GetSourceChannel(), packet.GetSequence(), contract)
	return nil
}

func (h Ics20Hooks) OnAcknowledgementPacketOverride(im IBCMiddleware, ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
	err := im.app.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
	if err != nil {
		return err
	}

	if !h.ProperlyConfigured() {
		// Not configured. Return from the underlying implementation
		return nil
	}

	contract := h.autoTxKeeper.GetPacketCallback(ctx, packet.GetSourceChannel(), packet.GetSequence())
	if contract == "" {
		// No callback configured
		return nil
	}

	contractAddr, err := sdk.AccAddressFromBech32(contract)
	if err != nil {
		return sdkerrors.Wrap(err, "Ack callback error") // The callback configured is not a bech32. Error out
	}

	success := "false"
	if !osmoutils.IsAckError(acknowledgement) {
		success = "true"
	}

	// Notify the sender that the ack has been received
	ackAsJson, err := json.Marshal(acknowledgement)
	if err != nil {
		// If the ack is not a json object, error
		return err
	}

	sudoMsg := []byte(fmt.Sprintf(
		`{"ibc_lifecycle_complete": {"ibc_ack": {"channel": "%s", "sequence": %d, "ack": %s, "success": %s}}}`,
		packet.SourceChannel, packet.Sequence, ackAsJson, success))
	_, err = h.ContractKeeper.Sudo(ctx, contractAddr, sudoMsg)
	if err != nil {
		// error processing the callback
		// ToDo: Open Question: Should we also delete the callback here?
		return sdkerrors.Wrap(err, "Ack callback error")
	}
	h.autoTxKeeper.DeletePacketCallback(ctx, packet.GetSourceChannel(), packet.GetSequence())
	return nil
}

func (h Ics20Hooks) OnTimeoutPacketOverride(im IBCMiddleware, ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) error {
	err := im.app.OnTimeoutPacket(ctx, packet, relayer)
	if err != nil {
		return err
	}

	if !h.ProperlyConfigured() {
		// Not configured. Return from the underlying implementation
		return nil
	}

	contract := h.autoTxKeeper.GetPacketCallback(ctx, packet.GetSourceChannel(), packet.GetSequence())
	if contract == "" {
		// No callback configured
		return nil
	}

	contractAddr, err := sdk.AccAddressFromBech32(contract)
	if err != nil {
		return sdkerrors.Wrap(err, "Timeout callback error") // The callback configured is not a bech32. Error out
	}

	sudoMsg := []byte(fmt.Sprintf(
		`{"ibc_lifecycle_complete": {"ibc_timeout": {"channel": "%s", "sequence": %d}}}`,
		packet.SourceChannel, packet.Sequence))
	_, err = h.ContractKeeper.Sudo(ctx, contractAddr, sudoMsg)
	if err != nil {
		// error processing the callback. This could be because the contract doesn't implement the message type to
		// process the callback. Retrying this will not help, so we can delete the callback from storage.
		// Since the packet has timed out, we don't expect any other responses that may trigger the callback.
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				"ibc-timeout-callback-error",
				sdk.NewAttribute("contract", contractAddr.String()),
				sdk.NewAttribute("message", string(sudoMsg)),
				sdk.NewAttribute("error", err.Error()),
			),
		})
	}
	h.autoTxKeeper.DeletePacketCallback(ctx, packet.GetSourceChannel(), packet.GetSequence())
	return nil
} */

func deriveIntermediateSender(channel, originalSender, bech32Prefix string) (string, error) {
	senderStr := fmt.Sprintf("%s/%s", channel, originalSender)
	senderHash32 := address.Hash(types.SenderPrefix, []byte(senderStr))
	sender := sdk.AccAddress(senderHash32[:])
	return sdk.Bech32ifyAddressBytes(bech32Prefix, sender)
}

// MustExtractDenomFromPacketOnRecv takes a packet with a valid ICS20 token data in the Data field and returns the
// denom as represented in the local chain.
// If the data cannot be unmarshalled this function will panic
func MustExtractDenomFromPacketOnRecv(packet ibcexported.PacketI) string {
	var data transfertypes.FungibleTokenPacketData
	if err := json.Unmarshal(packet.GetData(), &data); err != nil {
		panic("unable to unmarshal ICS20 packet data")
	}

	var denom string
	if transfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), data.Denom) {
		// remove prefix added by sender chain
		voucherPrefix := transfertypes.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())

		unprefixedDenom := data.Denom[len(voucherPrefix):]

		// coin denomination used in sending from the escrow address
		denom = unprefixedDenom

		// The denomination used to send the coins is either the native denom or the hash of the path
		// if the denomination is not native.
		denomTrace := transfertypes.ParseDenomTrace(unprefixedDenom)
		if denomTrace.Path != "" {
			denom = denomTrace.IBCDenom()
		}
	} else {
		prefixedDenom := transfertypes.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel()) + data.Denom
		denom = transfertypes.ParseDenomTrace(prefixedDenom).IBCDenom()
	}
	return denom
}
