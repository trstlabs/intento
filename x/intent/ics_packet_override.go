package intent

import (
	"encoding/json"
	"fmt"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"

	//codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/address"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/trstlabs/intento/x/intent/keeper"
	"github.com/trstlabs/intento/x/intent/types"
)

func onRecvPacketOverride(im IBCMiddleware, ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) ibcexported.Acknowledgement {

	isIcs20, data := isIcs20Packet(packet)
	if !isIcs20 {
		return im.app.OnRecvPacket(ctx, packet, relayer)
	}

	// Validate the memo
	isFlowRouted, ownerAddr, msgsBytes, label, connectionID, hostConnectionID, duration, interval, startAt, endTime, registerICA, hostedAddress, hostedFeeLimit, configuration, conditions, version, err := ValidateAndParseMemo(data.GetMemo(), data.Receiver)
	if !isFlowRouted {
		im.keeper.Logger(ctx).Debug("ics20 packet not routed")
		return im.app.OnRecvPacket(ctx, packet, relayer)
	}
	if err != nil {
		im.keeper.Logger(ctx).Debug("handling ICS20 packet memo content", "error", err.Error())
		return channeltypes.NewErrorAcknowledgement(err)
	}
	if msgsBytes == nil /* || ownerAddr == nil  */ { // This should never happen
		return channeltypes.NewErrorAcknowledgement(types.ErrMsgValidation)
	}
	var txMsgsAny []*codectypes.Any
	for _, msgBytes := range msgsBytes {
		var txMsgAny codectypes.Any
		cdc := codec.NewProtoCodec(im.registry)
		if err := cdc.UnmarshalJSON(msgBytes, &txMsgAny); err != nil {
			im.keeper.Logger(ctx).Debug("ICS20 packet unmarshalling flow message in msg array", "error", err.Error())
			return channeltypes.NewErrorAcknowledgement(types.ErrMsgValidation)
		}
		txMsgsAny = append(txMsgsAny, &txMsgAny)
	}
	//im.keeper.Logger(ctx).Info("ics20 got messages in array", "first", txMsgsAny[0].TypeUrl)
	// Calculate the receiver / contract caller based on the packet's channel and sender
	// The funds sent on this packet need to be transferred to an intermediary account for the sender.
	// For this, we override the ICS20 packet's Receiver (essentially hijacking the funds to this new address)
	// and execute the underlying OnRecvPacket() call. Hereafter we send the funds from the intermediary account to the flow FeeFunds address
	ownerAddr, errAck := makeOwnerForChannelSender(ownerAddr, &packet, data)
	if errAck != nil {
		return errAck
	}
	// Execute the receive of funds
	ack := im.app.OnRecvPacket(ctx, packet, relayer)
	if !ack.Success() {
		return ack
	}

	amount, ok := math.NewIntFromString(data.GetAmount())
	if !ok {
		// This should never happen, as it should've been caught in the underlaying call to OnRecvPacket,
		// but returning here for completeness
		return channeltypes.NewErrorAcknowledgement(errorsmod.Wrapf(types.ErrInvalidPacket, "Amount is not an int"))
	}

	// The packet's denom is the denom in the sender chain. This needs to be converted to the local denom.
	denom := MustExtractDenomFromPacketOnRecv(packet)
	funds := sdk.NewCoins(sdk.NewCoin(denom, amount))

	// Build the message to handle
	if registerICA {
		msg := types.MsgRegisterAccountAndSubmitFlow{
			Owner:            ownerAddr.String(),
			Msgs:             txMsgsAny,
			FeeFunds:         funds,
			Label:            label,
			ConnectionID:     connectionID,
			Duration:         duration,
			Interval:         interval,
			StartAt:          startAt,
			Configuration:    &configuration,
			HostConnectionID: hostConnectionID,
			Conditions:       &conditions,
			Version:          version,
		}
		response, err := registerAndSubmitTx(im.keeper, ctx, &msg)
		if err != nil {
			im.keeper.Logger(ctx).Debug("error handling ICS20 packet flow", err.Error())
			return channeltypes.NewErrorAcknowledgement(errorsmod.Wrapf(types.ErrIcs20Error, err.Error()))
		}
		bz, err := json.Marshal(response)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(errorsmod.Wrapf(types.ErrBadResponse, err.Error()))
		}

		return channeltypes.NewResultAcknowledgement(bz)
	} else if endTime != 0 {
		parsedOwnerAddr, errAck := makeOwnerForChannelSender(ownerAddr, &packet, data)
		if errAck != nil {
			return errAck
		}
		msg := types.MsgUpdateFlow{
			Owner:         parsedOwnerAddr.String(),
			Msgs:          txMsgsAny,
			FeeFunds:      funds,
			Label:         label,
			ConnectionID:  connectionID,
			EndTime:       endTime,
			Interval:      interval,
			StartAt:       startAt,
			Configuration: &configuration,
			Conditions:    &conditions,
			HostedICAConfig: &types.HostedICAConfig{HostedAddress: hostedAddress,
				FeeCoinLimit: hostedFeeLimit},
		}
		response, err := updateFlow(im.keeper, ctx, &msg)
		if err != nil {
			im.keeper.Logger(ctx).Debug("error handling ICS20 packet flow update", err.Error())
			return channeltypes.NewErrorAcknowledgement(errorsmod.Wrapf(types.ErrIcs20Error, err.Error()))
		}
		bz, err := json.Marshal(response)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(errorsmod.Wrapf(types.ErrBadResponse, err.Error()))
		}

		return channeltypes.NewResultAcknowledgement(bz)
	} else {
		msg := types.MsgSubmitFlow{
			Owner:         ownerAddr.String(),
			Msgs:          txMsgsAny,
			FeeFunds:      funds,
			Label:         label,
			Duration:      duration,
			Interval:      interval,
			StartAt:       startAt,
			Configuration: &configuration,
			Conditions:    &conditions,
			ConnectionID:  connectionID,
			HostedICAConfig: &types.HostedICAConfig{HostedAddress: hostedAddress,
				FeeCoinLimit: hostedFeeLimit},
		}
		response, err := submitFlow(im.keeper, ctx, &msg)
		if err != nil {
			im.keeper.Logger(ctx).Debug("error handling ICS20 packet flow submission", err.Error())
			return channeltypes.NewErrorAcknowledgement(errorsmod.Wrapf(types.ErrIcs20Error, err.Error()))
		}
		bz, err := json.Marshal(response)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(errorsmod.Wrapf(types.ErrBadResponse, err.Error()))
		}
		im.keeper.Logger(ctx).Debug("flow via ics20 submitted sucesssfully")
		return channeltypes.NewResultAcknowledgement(bz)
	}

}

func makeOwnerForChannelSender(ownerAddr sdk.AccAddress, packet *channeltypes.Packet, data transfertypes.FungibleTokenPacketData) (sdk.AccAddress, ibcexported.Acknowledgement) {
	if ownerAddr.Empty() {
		channel := packet.GetDestChannel()
		sender := data.GetSender()
		senderLocalAddr := derivePlaceholderSender(channel, sender)
		// if err != nil {
		// 	return nil, channeltypes.NewErrorAcknowledgement(errorsmod.Wrapf(types.ErrBadSender, fmt.Sprintf("cannot convert sender address %s/%s to bech32: %s", channel, sender, err)))
		// }

		data.Receiver = senderLocalAddr.String()
		bz, err := json.Marshal(data)
		if err != nil {
			return nil, channeltypes.NewErrorAcknowledgement(errorsmod.Wrapf(types.ErrMarshaling, err.Error()))
		}
		packet.Data = bz
		ownerAddr = senderLocalAddr
	}
	return ownerAddr, nil
}

func registerAndSubmitTx(k keeper.Keeper, ctx sdk.Context, ics20ParsedMsg *types.MsgRegisterAccountAndSubmitFlow) (*types.MsgRegisterAccountAndSubmitFlowResponse, error) {
	if err := ics20ParsedMsg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf(types.ErrBadFlowMsg, err.Error())
	}
	ics20MsgServer := keeper.NewMsgServerImpl(k)
	return ics20MsgServer.RegisterAccountAndSubmitFlow(ctx, ics20ParsedMsg)
}

func submitFlow(k keeper.Keeper, ctx sdk.Context, ics20ParsedMsg *types.MsgSubmitFlow) (*types.MsgSubmitFlowResponse, error) {
	if err := ics20ParsedMsg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf(types.ErrBadFlowMsg, err.Error())
	}
	ics20MsgServer := keeper.NewMsgServerImpl(k)
	return ics20MsgServer.SubmitFlow(ctx, ics20ParsedMsg)
}

func updateFlow(k keeper.Keeper, ctx sdk.Context, ics20ParsedMsg *types.MsgUpdateFlow) (*types.MsgUpdateFlowResponse, error) {
	if err := ics20ParsedMsg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf(types.ErrBadFlowMsg, err.Error())
	}
	ics20MsgServer := keeper.NewMsgServerImpl(k)
	return ics20MsgServer.UpdateFlow(ctx, ics20ParsedMsg)
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

func ValidateAndParseMemo(memo string, receiver string) (isFlowRouted bool, ownerAddr sdk.AccAddress, msgsBytes [][]byte, label, connectionID, hostConnectionID, duration, interval string, startAt uint64, endTime uint64, registerICA bool, hostedAddress string, hostedFeeLimit sdk.Coin, configuration types.ExecutionConfiguration, conditions types.ExecutionConditions, version string, err error) {
	isFlowRouted, metadata := jsonStringHasKey(memo, "flow")
	if !isFlowRouted {
		return isFlowRouted, sdk.AccAddress{}, nil, "", "", "", "", "", 0, 0, false, "", sdk.Coin{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, "", nil
	}
	ics20Raw := metadata["flow"]

	// Make sure the ics20 key is a map. If it isn't, ignore this packet
	flow, ok := ics20Raw.(map[string]interface{})
	if !ok {
		return isFlowRouted, sdk.AccAddress{}, nil, "", "", "", "", "", 0, 0, false, "", sdk.Coin{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, "",
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, "flow metadata is not a valid JSON map object")
	}

	// Get the owner
	owner, ok := flow["owner"].(string)
	if !ok {
		owner = ""
	}

	// Owner is optional and the owner and the receiver should be the same for the packet to be valid
	if ok && owner != "" {
		if owner != receiver {
			return isFlowRouted, sdk.AccAddress{}, nil, "", "", "", "", "", 0, 0, false, "", sdk.Coin{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, "",
				fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `flow["owner"] should be the same as the receiver of the packet`)
		}
		ownerAddr, err = sdk.AccAddressFromBech32(owner)
		if err != nil {
			return isFlowRouted, sdk.AccAddress{}, nil, "", "", "", "", "", 0, 0, false, "", sdk.Coin{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, "",
				fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `flow["owner"] is not a valid bech32 address`)
		}

	}

	// Ensure the message key is provided
	if flow["msgs"] == nil {
		return isFlowRouted, sdk.AccAddress{}, nil, "", "", "", "", "", 0, 0, false, "", sdk.Coin{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, "",
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `flow["msgs"]`)
	}

	// Make sure the msg key is an array of maps. If it isn't, return an error
	msgs, ok := flow["msgs"].([]interface{})
	if !ok {
		return isFlowRouted, sdk.AccAddress{}, nil, "", "", "", "", "", 0, 0, false, "", sdk.Coin{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, "",
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `flow["msgs"] is not an array of interfaces`)
	}

	// Get the label
	label, ok = flow["label"].(string)
	if !ok {
		// The tokens will be returned
		return isFlowRouted, sdk.AccAddress{}, nil, "", "", "", "", "", 0, 0, false, "", sdk.Coin{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, "",
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `flow["label"]`)
	}

	// Get the connectionID. To save space we write cid instead of connection_id
	connectionID, _ = flow["cid"].(string)

	// Get the version
	hostConnectionID, _ = flow["host_cid"].(string)

	// optional for updating trigger end time
	endTimeString, ok := flow["end_time"].(string)
	if ok {
		endTime, err = strconv.ParseUint(endTimeString, 10, 64)
		if err != nil {
			return isFlowRouted, sdk.AccAddress{}, nil, "", "", "", "", "", 0, 0, false, "", sdk.Coin{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, "",
				fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `flow["end_time"]`)
		}
	}

	// Get the duration
	duration, ok = flow["duration"].(string)
	// A sumbitFlow should have a duration key, an updateFlow should have an endTime key
	if !ok && endTime == 0 {
		// The tokens will be returned
		return isFlowRouted, sdk.AccAddress{}, nil, "", "", "", "", "", 0, 0, false, "", sdk.Coin{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, "",
			fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `flow["duration"]`)
	}
	// Get the interval,optional
	interval, _ = flow["interval"].(string)

	// Get the label
	startAtString, ok := flow["start_at"].(string)
	if ok {
		startAt, err = strconv.ParseUint(startAtString, 10, 64)

		if err != nil {
			return isFlowRouted, sdk.AccAddress{}, nil, "", "", "", "", "", 0, 0, false, "", sdk.Coin{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, "",
				fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `flow["start_at"]`)
		}
	}

	conditionsString := flow["conditions"].(string)
	if conditionsString != "" {
		if err := json.Unmarshal([]byte(conditionsString), &conditions); err != nil {
			return isFlowRouted, sdk.AccAddress{}, nil, "", "", "", "", "", 0, 0, false, "", sdk.Coin{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, "", fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `flow["conditions"]`)
		}
	}

	//optional hosted account
	hostedAddress, ok = flow["hosted_account"].(string)
	if !ok {
		hostedAddress = ""
	}

	hostedFeeLimitString, ok := flow["hosted_fee_limit"].(string)
	if ok {
		// return isFlowRouted, sdk.AccAddress{}, nil, "", "", "", "", "", 0, 0, false, "", sdk.Coin{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, "", fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `flow["hosted_fee_limit"]`)
		hostedFeeLimit, err = sdk.ParseCoinNormalized(hostedFeeLimitString)
		if err != nil {
			return isFlowRouted, sdk.AccAddress{}, nil, "", "", "", "", "", 0, 0, false, "", sdk.Coin{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, "", fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, `flow["hosted_fee_limit"]`)
		}
	}

	registerICAString, ok := flow["register_ica"].(string)
	if ok && registerICAString == "true" {
		registerICA = true
	}

	updateDisabled := false
	updateDisabledString, ok := flow["update_disabled"].(string)
	if ok && updateDisabledString == "true" {
		updateDisabled = true
	}

	SaveResponses := false
	SaveResponsesString, ok := flow["save_responses"].(string)
	if ok && SaveResponsesString == "true" {
		SaveResponses = true
	}

	stopOnSuccess := false
	stopOnSuccessString, ok := flow["stop_on_success"].(string)
	if ok && stopOnSuccessString == "true" {
		stopOnSuccess = true
	}

	stopOnFailure := false
	stopOnFailureString, ok := flow["stop_on_fail"].(string)
	if ok && stopOnFailureString == "true" {
		stopOnFailure = true
	}

	stopOnTimeout := false
	stopOnTimeoutString, ok := flow["stop_on_timeout"].(string)
	if ok && stopOnTimeoutString == "true" {
		stopOnTimeout = true
	}

	fallbackOwner := false
	fallbackOwnerString, ok := flow["fallback"].(string)
	if ok && fallbackOwnerString == "true" {
		fallbackOwner = true
	}

	configuration = types.ExecutionConfiguration{
		SaveResponses:          SaveResponses,
		UpdatingDisabled:       updateDisabled,
		StopOnSuccess:          stopOnSuccess,
		StopOnFailure:          stopOnFailure,
		StopOnTimeout:          stopOnTimeout,
		FallbackToOwnerBalance: fallbackOwner,
	}

	version = string(icatypes.ModuleCdc.MustMarshalJSON(&icatypes.Metadata{
		Version:                icatypes.Version,
		ControllerConnectionId: connectionID,
		HostConnectionId:       hostConnectionID,
		Encoding:               icatypes.EncodingProtobuf,
		TxType:                 icatypes.TxTypeSDKMultiMsg,
	}))
	//var msgsBytes [][]byte
	// Get the message string by serializing the map
	for _, msg := range msgs {
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			// The tokens will be returned
			return isFlowRouted, sdk.AccAddress{}, nil, "", "", "", "", "", 0, 0, false, "", sdk.Coin{}, types.ExecutionConfiguration{}, types.ExecutionConditions{}, "",
				fmt.Errorf(types.ErrBadMetadataFormatMsg, memo, err.Error())
		}
		msgsBytes = append(msgsBytes, msgBytes)
	}

	return isFlowRouted, ownerAddr, msgsBytes, label, connectionID, hostConnectionID, duration, interval, startAt, endTime, registerICA, hostedAddress, hostedFeeLimit, configuration, conditions, version, nil
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
