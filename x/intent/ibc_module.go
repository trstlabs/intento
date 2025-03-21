package intent

import (
	errorsmod "cosmossdk.io/errors"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/cosmos/gogoproto/proto"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/trstlabs/intento/x/intent/keeper"
	"github.com/trstlabs/intento/x/intent/types"
)

var _ porttypes.IBCModule = IBCModule{}

// IBCModule implements the ICS26 interface for interchain accounts controller chains
type IBCModule struct {
	keeper keeper.Keeper
}

// NewIBCModule creates a new IBCModule given the keeper
func NewIBCModule(k keeper.Keeper) IBCModule {
	return IBCModule{
		keeper: k,
	}
}

// OnChanOpenInit implements the IBCModule interface
func (im IBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) (string, error) {
	return version, nil
}

// OnChanOpenTry implements the IBCModule interface
func (im IBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (version string, err error) {
	return "", nil
}

// OnChanOpenAck implements the IBCModule interface
func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartychannelID string,
	counterpartyVersion string,
) error {
	return nil
}

// OnChanOpenConfirm implements the IBCModule interface
func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnChanCloseInit implements the IBCModule interface
func (im IBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnChanCloseConfirm implements the IBCModule interface
func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnRecvPacket implements the IBCModule interface. A successful acknowledgement
// is returned if the packet data is successfully decoded and the receive application
// logic returns without error.
func (im IBCModule) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	return channeltypes.NewErrorAcknowledgement(errorsmod.Wrapf(types.ErrInvalidRequest, "cannot receive packet via interchain accounts authentication module"))
}

// OnAcknowledgementPacket implements the IBCModule interface
func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {

	var ack channeltypes.Acknowledgement
	if err := channeltypes.SubModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return errorsmod.Wrapf(types.ErrUnknownRequest, "cannot unmarshal ICS-27 packet acknowledgement: %v", err)
	}
	if !ack.Success() {
		errorString := "error handling packet on host chain: see host chain events for details, error: " + ack.GetError()
		im.keeper.SetFlowError(ctx, packet.SourcePort, packet.SourceChannel, packet.Sequence, errorString)
		return nil

	}
	var txMsgData sdk.TxMsgData
	if err := proto.Unmarshal(ack.GetResult(), &txMsgData); err != nil {
		return errorsmod.Wrapf(types.ErrUnknownRequest, "cannot unmarshal ICS-27 tx message data: %v", err)
	}

	//handle message data
	switch len(txMsgData.Data) {
	case 0:
		// for SDK 0.46 and above
		im.handleMsgResponses(ctx, txMsgData.GetMsgResponses(), relayer, packet)

		//we process errors internally and return nil so acknoledgement is succesfull and ordered channel stays active.
		return nil
	default:
		return nil // im.handleDeprecatedMsgResponses(ctx, txMsgData, relayer, packet)
	}
}

func (im IBCModule) handleMsgResponses(ctx sdk.Context, msgResponses []*cdctypes.Any, relayer sdk.AccAddress, packet channeltypes.Packet) {
	if len(msgResponses) == 0 {
		err := errorsmod.Wrapf(types.ErrInvalidType, "no messages in ICS-27 message response: %v", msgResponses)
		im.keeper.SetFlowError(ctx, packet.SourcePort, packet.SourceChannel, packet.Sequence, err.Error())
		return
	}

	// handle response (trigger next messages if response parsing) set result in Flow history
	err := im.keeper.HandleResponseAndSetFlowResult(ctx, packet.SourcePort, packet.SourceChannel, relayer, packet.Sequence, msgResponses)
	if err != nil {
		im.keeper.SetFlowError(ctx, packet.SourcePort, packet.SourceChannel, packet.Sequence, err.Error())
	}
}

// OnTimeoutPacket implements the IBCModule interface.
func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	// fmt.Println("TIMED OUT, FAILED ATTEMPT")
	//set result in flow history
	err := im.keeper.SetFlowOnTimeout(ctx, packet.SourcePort, packet.SourceChannel, packet.Sequence)
	if err != nil {
		im.keeper.SetFlowError(ctx, packet.SourcePort, packet.SourceChannel, packet.Sequence, err.Error())
		return err
	}
	return nil
}

// NegotiateAppVersion implements the IBCModule interface
func (im IBCModule) NegotiateAppVersion(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionID string,
	portID string,
	counterparty channeltypes.Counterparty,
	proposedVersion string,
) (string, error) {
	return "", nil
}
