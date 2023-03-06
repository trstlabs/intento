package autoibctx

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	//codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	//"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/types/address"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v4/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v4/modules/core/exported"
	"github.com/trstlabs/trst/x/auto-ibc-tx/keeper"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

var _ porttypes.Middleware = &IBCMiddleware{}

type IBCMiddleware struct {
	app      porttypes.IBCModule
	keeper   keeper.Keeper //add a keeper for stateful middleware
	registry codectypes.InterfaceRegistry
}

// IBCMiddleware creates a new IBCMiddleware given the associated keeper and underlying application
func NewIBCMiddleware(app porttypes.IBCModule, k keeper.Keeper, registry codectypes.InterfaceRegistry) IBCMiddleware {
	return IBCMiddleware{
		app:      app,
		keeper:   k,
		registry: registry,
	}
}

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
) (string, error) {

	finalVersion, err := im.app.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version)

	return finalVersion, err
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

	ack := onRecvPacketOverride(im, ctx, packet, relayer)
	fmt.Printf("ack %v\n", ack.Acknowledgement())
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
	panic("SendPacket not supported forics middleware")
}

// WriteAcknowledgement implements the ICS4 Wrapper interface
func (im IBCMiddleware) WriteAcknowledgement(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	packet ibcexported.PacketI,
	ack ibcexported.Acknowledgement,
) error {
	panic("WriteAcknowledgement not supported forics middleware")
}

// GetAppVersion returns the interchain accounts metadata.
func (im IBCMiddleware) GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool) {
	//return im.app.GetAppVersion(ctx, portID, channelID)
	panic("GetAppVersion not supported forics middleware")
}

func onRecvPacketOverride(im IBCMiddleware, ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) ibcexported.Acknowledgement {

	isIcs20, data := isIcs20Packet(packet)
	if !isIcs20 {
		return im.app.OnRecvPacket(ctx, packet, relayer)
	}

	// Validate the memo
	isAutoTxRouted, ownerAddr, msgsBytes, label, connectionID, duration, interval, startAt, registerICA, err := ValidateAndParseMemo(data.GetMemo(), data.Receiver)

	if !isAutoTxRouted {
		return im.app.OnRecvPacket(ctx, packet, relayer)
	}
	if err != nil {

		return channeltypes.NewErrorAcknowledgement(err)
	}
	if msgsBytes == nil /* || ownerAddr == nil  */ { // This should never happen
		return channeltypes.NewErrorAcknowledgement(sdkerrors.Wrapf(types.ErrMsgValidation, err.Error()))
	}
	var txMsgsAny []*codectypes.Any
	for _, msgBytes := range msgsBytes {
		var txMsgAny codectypes.Any
		cdc := codec.NewProtoCodec(im.registry)
		if err := cdc.UnmarshalJSON(msgBytes, &txMsgAny); err != nil {
			return channeltypes.NewErrorAcknowledgement(types.ErrMsgValidation)
		}
		txMsgsAny = append(txMsgsAny, &txMsgAny)
	}
	if ownerAddr.Empty() {

		// Calculate the receiver / contract caller based on the packet's channel and sender
		channel := packet.GetDestChannel()
		sender := data.GetSender()
		senderLocalAddr := derivePlaceholderSender(channel, sender)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(sdkerrors.Wrapf(types.ErrBadSender, fmt.Sprintf("cannot convert sender address %s/%s to bech32: %s", channel, sender, err)))
		}

		// The funds sent on this packet need to be transferred to an intermediary account for the sender.
		// For this, we override the ICS20 packet's Receiver (essentially hijacking the funds to this new address)
		// and execute the underlying OnRecvPacket() call. Hereafter we send the funds from the intermediary account to the AutoTx FeeFunds address

		data.Receiver = senderLocalAddr.String()
		bz, err := json.Marshal(data)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(sdkerrors.Wrapf(types.ErrMarshaling, err.Error()))
		}
		packet.Data = bz
		ownerAddr = senderLocalAddr
	}
	// Execute the receive of funds
	ack := im.app.OnRecvPacket(ctx, packet, relayer)
	if !ack.Success() {
		return ack
	}

	amount, ok := sdk.NewIntFromString(data.GetAmount())
	if !ok {
		// This should never happen, as it should've been caught in the underlaying call to OnRecvPacket,
		// but returning here for completeness
		return channeltypes.NewErrorAcknowledgement(sdkerrors.Wrapf(types.ErrInvalidPacket, "Amount is not an int"))
	}

	// The packet's denom is the denom in the sender chain. This needs to be converted to the local denom.
	denom := MustExtractDenomFromPacketOnRecv(packet)
	funds := sdk.NewCoins(sdk.NewCoin(denom, amount))

	// Build the message to handle
	if registerICA {
		msg := types.MsgRegisterAccountAndSubmitAutoTx{
			Owner:        ownerAddr.String(),
			Msgs:         txMsgsAny,
			FeeFunds:     funds,
			Label:        label,
			ConnectionId: connectionID,
			Duration:     duration,
			Interval:     interval,
			StartAt:      startAt,
			//Version:      version,
		}
		response, err := registerAndSubmitTx(im.keeper, ctx, &msg)
		if err != nil {
			// fmt.Println(err.Error())
			return channeltypes.NewErrorAcknowledgement(sdkerrors.Wrapf(types.ErrIcs20Error, err.Error()))
		}
		bz, err := json.Marshal(response)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(sdkerrors.Wrapf(types.ErrBadResponse, err.Error()))
		}

		return channeltypes.NewResultAcknowledgement(bz)
	} else {
		msg := types.MsgSubmitAutoTx{
			Owner:        ownerAddr.String(),
			Msgs:         txMsgsAny,
			FeeFunds:     funds,
			Label:        label,
			ConnectionId: connectionID,
			Duration:     duration,
			Interval:     interval,
			StartAt:      startAt,
		}
		response, err := submitTx(im.keeper, ctx, &msg)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(sdkerrors.Wrapf(types.ErrIcs20Error, err.Error()))
		}
		bz, err := json.Marshal(response)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(sdkerrors.Wrapf(types.ErrBadResponse, err.Error()))
		}

		return channeltypes.NewResultAcknowledgement(bz)
	}

}

func registerAndSubmitTx(k keeper.Keeper, ctx sdk.Context, autoTxMsg *types.MsgRegisterAccountAndSubmitAutoTx) (*types.MsgRegisterAccountAndSubmitAutoTxResponse, error) {
	if err := autoTxMsg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf(types.ErrBadAutoTxMsg, err.Error())
	}
	ics20MsgServer := keeper.NewMsgServerImpl(k)
	return ics20MsgServer.RegisterAccountAndSubmitAutoTx(sdk.WrapSDKContext(ctx), autoTxMsg)
}

func submitTx(k keeper.Keeper, ctx sdk.Context, autoTxMsg *types.MsgSubmitAutoTx) (*types.MsgSubmitAutoTxResponse, error) {
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

func ValidateAndParseMemo(memo string, receiver string) (isAutoTxRouted bool, ownerAddr sdk.AccAddress, msgsBytes [][]byte, label, connectionID, duration, interval string, startAt uint64 /* version string, */, registerICA bool, err error) {
	isAutoTxRouted, metadata := jsonStringHasKey(memo, "auto_tx")
	if !isAutoTxRouted {
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false, nil
	}

	ics20Raw := metadata["auto_tx"]

	// Make sure the ics20 key is a map. If it isn't, ignore this packet
	autoTx, ok := ics20Raw.(map[string]interface{})
	if !ok {
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, "auto_tx metadata is not a valid JSON map object")
	}

	// Get the owner
	owner, ok := autoTx["owner"].(string)
	if !ok {
		// The tokens will follow regular MsgTransfer pattern
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `auto_tx["owner"]`)
	}

	// The owner and the receiver should be the same for the packet to be valid
	if ok && owner != "" {
		if owner != receiver {
			return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
				fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `auto_tx["owner"] should be the same as the receiver of the packet`)
		}
		ownerAddr, err = sdk.AccAddressFromBech32(owner)
		if err != nil {
			return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
				fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `auto_tx["owner"] is not a valid bech32 address`)
		}

	}

	// Ensure the message key is provided
	if autoTx["msgs"] == nil {
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `auto_tx["msgs"]`)
	}

	// Make sure the msg key is an array of maps. If it isn't, return an error
	msgs, ok := autoTx["msgs"].([]interface{})
	if !ok {
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `auto_tx["msgs"] is not an array`)
	}

	// Get the label
	label, ok = autoTx["label"].(string)
	if !ok {
		// The tokens will be returned
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `auto_tx["label"]`)
	}

	// Get the connectionID
	connectionID, ok = autoTx["connection_id"].(string)
	if !ok {
		connectionID = ""
	}
	// Get the duration
	duration, ok = autoTx["duration"].(string)
	if !ok {
		// The tokens will be returned
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `auto_tx["duration"]`)
	}
	// Get the interval
	interval, ok = autoTx["interval"].(string)
	if !ok {
		// The tokens will be returned
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `auto_tx["interval"]`)
	}
	// Get the label
	startAtString, ok := autoTx["start_at"].(string)
	if !ok {
		// The tokens will be returned
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `auto_tx["start_at"]`)
	}
	startAt, err = strconv.ParseUint(startAtString, 10, 64)
	if err != nil {
		return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `auto_tx["start_at"]`)
	}
	/* // Get the version
	version, ok = autoTx["version"].(string)
	if !ok {
		version = ""
	} */
	registerICAString, ok := autoTx["register_ica"].(string)
	if ok && registerICAString == "true" {
		registerICA = true
	}

	//var msgsBytes [][]byte
	// Get the message string by serializing the map
	for _, msg := range msgs {
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			// The tokens will be returned

			return isAutoTxRouted, sdk.AccAddress{}, nil, "", "", "", "", 0, false,
				fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, err.Error())
		}
		msgsBytes = append(msgsBytes, msgBytes)
	}

	return isAutoTxRouted, ownerAddr, msgsBytes, label, connectionID, duration, interval, startAt, registerICA /* version */, nil
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

func derivePlaceholderSender(channel, originalSender string) sdk.AccAddress {
	senderStr := fmt.Sprintf("%s/%s", channel, originalSender)
	senderHash32 := address.Hash(types.SenderPrefix, []byte(senderStr))
	sender := sdk.AccAddress(senderHash32[:])
	return sender
}
