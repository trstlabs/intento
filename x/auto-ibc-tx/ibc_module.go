package autoibctx

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	proto "github.com/cosmos/gogoproto/proto"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/trstlabs/trst/x/auto-ibc-tx/keeper"
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

	txMsgData := &sdk.TxMsgData{}
	if err := proto.Unmarshal(ack.GetResult(), txMsgData); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 tx message data: %v", err)
	}

	//rewardType is used for Relayer Reward and Airdrop Reward
	//handle message data
	//set result in auto-tx history
	switch len(txMsgData.Data) {
	case 0:
		return im.handleMsgResponses(ctx, txMsgData.GetMsgResponses(), relayer, packet)

	default:
		return im.handleDeprecatedMsgResponses(ctx, txMsgData, relayer, packet)
	}
}

func (im IBCModule) handleMsgResponses(ctx sdk.Context, msgResp []*cdctypes.Any, relayer sdk.AccAddress, packet channeltypes.Packet) error {

	rewardType := -1

	for index, anyResp := range msgResp {
		im.keeper.Logger(ctx).Info("msg response in ICS-27 packet", "response", anyResp.GoString(), "typeURL", anyResp.GetTypeUrl())
		rewardClass := getMsgRewardType(ctx, anyResp.GetTypeUrl())
		if index == 0 {
			rewardType = rewardClass
		}
	}

	if rewardType >= 0 {
		im.keeper.HandleRelayerReward(ctx, relayer, rewardType)

	}

	err := im.keeper.SetAutoTxResult(ctx, packet.SourcePort, rewardType, packet.Sequence)
	if err != nil {
		im.keeper.SetAutoTxError(ctx, packet.SourcePort, packet.Sequence, err.Error())
		return err
	}

	return nil
}

func (im IBCModule) handleDeprecatedMsgResponses(ctx sdk.Context, txMsgData *sdk.TxMsgData, relayer sdk.AccAddress, packet channeltypes.Packet) error {
	rewardType := -1

	for index, msgData := range txMsgData.Data {
		response, rewardClass, err := handleMsgData(ctx, msgData)
		if err != nil {
			fmt.Printf("handleMsgData err: %v\n", err)
			return err
		}

		im.keeper.Logger(ctx).Debug("message response in ICS-27 packet response", "response", response)
		if index == 0 {
			rewardType = rewardClass
		}
	}

	if rewardType >= 0 {
		im.keeper.HandleRelayerReward(ctx, relayer, rewardType)

	}

	err := im.keeper.SetAutoTxResult(ctx, packet.SourcePort, rewardType, packet.Sequence)
	if err != nil {
		im.keeper.SetAutoTxError(ctx, packet.SourcePort, packet.Sequence, err.Error())
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
	//set result in auto-tx history

	err := im.keeper.SetAutoTxOnTimeout(ctx, packet.SourcePort, packet.Sequence)
	if err != nil {
		im.keeper.SetAutoTxError(ctx, packet.SourcePort, packet.Sequence, err.Error())
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
