package intent

import (
	errorsmod "cosmossdk.io/errors"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	proto "github.com/cosmos/gogoproto/proto"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/trstlabs/intento/x/intent/keeper"
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
	return channeltypes.NewErrorAcknowledgement(errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "cannot receive packet via interchain accounts authentication module"))
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
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 packet acknowledgement: %v", err)
	}
	if !ack.Success() {
		im.keeper.SetActionError(ctx, packet.SourcePort, packet.SourceChannel, packet.Sequence, "error handling packet on host chain: see host chain events for details")
	}
	var txMsgData sdk.TxMsgData
	if err := proto.Unmarshal(ack.GetResult(), &txMsgData); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 tx message data: %v", err)
	}

	//handle message data
	switch len(txMsgData.Data) {
	case 0:
		// for SDK 0.46 and above
		im.handleMsgResponses(ctx, txMsgData.GetMsgResponses(), relayer, packet)

		//we process errors internally and return nil so acknoledgement is succesfull and ordered channel stays active.
		return nil
	default:
		return im.handleDeprecatedMsgResponses(ctx, txMsgData, relayer, packet)
	}
}

func (im IBCModule) handleMsgResponses(ctx sdk.Context, msgResponses []*cdctypes.Any, relayer sdk.AccAddress, packet channeltypes.Packet) {
	if len(msgResponses) == 0 {
		err := errorsmod.Wrapf(sdkerrors.ErrInvalidType, "no messages in ICS-27 message response: %v", msgResponses)
		im.keeper.SetActionError(ctx, packet.SourcePort, packet.SourceChannel, packet.Sequence, err.Error())
		return
	}
	//msgClass is used for Relayer Reward and Airdrop Reward
	var msgClass int // Initialize msgClass outside the loop
	for index, anyResp := range msgResponses {
		im.keeper.Logger(ctx).Debug("msg response in ICS-27 packet", "response", anyResp.GoString(), "typeURL", anyResp.GetTypeUrl())

		rewardClass := getMsgRewardType(ctx, anyResp.GetTypeUrl())
		if index == 0 && rewardClass > 0 {
			msgClass = rewardClass
			im.keeper.HandleRelayerReward(ctx, relayer, msgClass)
		}
	}

	// set result in Action history
	err := im.keeper.SetActionResult(ctx, packet.SourcePort, packet.SourceChannel, msgClass, packet.Sequence, msgResponses)
	if err != nil {
		im.keeper.SetActionError(ctx, packet.SourcePort, packet.SourceChannel, packet.Sequence, err.Error())
		return
	}

	return
}

func (im IBCModule) handleDeprecatedMsgResponses(ctx sdk.Context, txMsgData sdk.TxMsgData, relayer sdk.AccAddress, packet channeltypes.Packet) error {
	//msgClass is used for Relayer Reward and Airdrop Reward
	var msgClass int // Initialize msgClass outside the loop

	for index, msgData := range txMsgData.Data {
		response, rewardClass, err := handleMsgData(ctx, msgData)
		if err != nil {
			im.keeper.Logger(ctx).Debug("message response error in ICS-27 packet response", "error", err)
			return err
		}

		im.keeper.Logger(ctx).Debug("message response in ICS-27 packet response", "response", response.String())
		if index == 0 && rewardClass > 0 {
			msgClass = rewardClass
			im.keeper.HandleRelayerReward(ctx, relayer, rewardClass)
		}

	}

	// set result in Action history
	err := im.keeper.SetActionResult(ctx, packet.SourcePort, packet.SourceChannel, msgClass, packet.Sequence, nil)
	if err != nil {
		im.keeper.SetActionError(ctx, packet.SourcePort, packet.SourceChannel, packet.Sequence, err.Error())
		return err
	}

	return nil
}

// OnTimeoutPacket implements the IBCModule interface.
func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	// fmt.Println("TIMED OUT, FAILED ATTEMPT")
	//set result in action history
	err := im.keeper.SetActionOnTimeout(ctx, packet.SourcePort, packet.SourceChannel, packet.Sequence)
	if err != nil {
		im.keeper.SetActionError(ctx, packet.SourcePort, packet.SourceChannel, packet.Sequence, err.Error())
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
