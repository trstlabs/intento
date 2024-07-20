package intent

import (
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"

	//codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/trstlabs/intento/x/intent/keeper"
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
	sourcePort string,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (uint64, error) {
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
