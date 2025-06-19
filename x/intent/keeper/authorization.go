package keeper

import (
	"fmt"

	msgv1 "cosmossdk.io/api/cosmos/msg/v1"
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-proto/anyutil"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	cosmosproto "github.com/cosmos/gogoproto/proto"
	"github.com/trstlabs/intento/x/intent/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
)

func (k Keeper) SignerOk(ctx sdk.Context, codec codec.Codec, flowInfo types.FlowInfo) error {
	for _, message := range flowInfo.Msgs {
		if err := k.validateMessage(ctx, codec, flowInfo, message); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) validateMessage(ctx sdk.Context, codec codec.Codec, flowInfo types.FlowInfo, message *codectypes.Any) error {
	var sdkMsg sdk.Msg
	if err := codec.UnpackAny(message, &sdkMsg); err != nil {
		return errorsmod.Wrap(err, "failed to unpack message")
	}

	switch {
	case isAuthzMsgExec(message):
		// Validate Authz MsgExec messages.
		return k.validateAuthzMsg(ctx, codec, flowInfo, message)

	case isLocalMessage(flowInfo):
		// Validate local messages to ensure the signer matches the owner.
		return k.validateSigners(ctx, codec, flowInfo, message)

	case isHostedICAMessage(flowInfo):
		// Restrict Hosted ICA messages to MsgExec for security.
		if message.TypeUrl != sdk.MsgTypeURL(&authztypes.MsgExec{}) {
			return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "only MsgExec is allowed for Hosted ICA messages")
		}
		return nil

	case isSelfHostedICAMessage(flowInfo):
		// Allow Self-hosted ICA messages without additional validation.
		return nil

	default:
		// Unsupported message type.
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "unsupported message type")
	}
}

func isAuthzMsgExec(message *codectypes.Any) bool {
	return message.TypeUrl == sdk.MsgTypeURL(&authztypes.MsgExec{})
}

func isLocalMessage(flowInfo types.FlowInfo) bool {
	return (flowInfo.ICAConfig == nil || flowInfo.ICAConfig.ConnectionID == "") && (flowInfo.HostedICAConfig == nil || flowInfo.HostedICAConfig.HostedAddress == "")
}

func isHostedICAMessage(flowInfo types.FlowInfo) bool {
	return flowInfo.HostedICAConfig != nil && flowInfo.HostedICAConfig.HostedAddress != ""
}

func isSelfHostedICAMessage(flowInfo types.FlowInfo) bool {
	return flowInfo.ICAConfig != nil && flowInfo.ICAConfig.ConnectionID != ""
}

// validateAuthzMsg validates an authz MsgExec message.
func (k Keeper) validateAuthzMsg(ctx sdk.Context, codec codec.Codec, flowInfo types.FlowInfo, message *codectypes.Any) error {
	var authzMsg authztypes.MsgExec
	if err := cosmosproto.Unmarshal(message.Value, &authzMsg); err != nil {
		return errorsmod.Wrap(err, "failed to unmarshal MsgExec")
	}

	for _, innerMessage := range authzMsg.Msgs {
		if err := k.validateSigners(ctx, codec, flowInfo, innerMessage); err != nil {
			return err
		}
	}
	return nil
}

// validateSigners checks the signers of a message against the owner, ICA, and hosted accounts.
func (k Keeper) validateSigners(ctx sdk.Context, codec codec.Codec, flowInfo types.FlowInfo, message *codectypes.Any) error {

	protoReflectMsg, err := unpackV2Any(codec, message)
	if err != nil {
		return errorsmod.Wrap(err, "failed to unpack message")
	}

	signers, err := extractSigners(protoReflectMsg)
	if err != nil {
		return errorsmod.Wrap(err, "failed to get message signers")
	}
	if len(signers) < 1 {
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "no valid signers found")
	}
	ownerAddr, err := sdk.AccAddressFromBech32(flowInfo.Owner)
	if err != nil {
		return errorsmod.Wrap(err, "failed to parse owner address")
	}
	signer, err := parseAccAddressFromAnyPrefix(signers[0])
	if err != nil {
		return errorsmod.Wrap(err, "failed to parse owner address")
	}

	k.Logger(ctx).Debug("Signer validation", "owner", flowInfo.Owner, "signer", signers[0])
	if !signer.Equals(ownerAddr) {
		errorString := fmt.Sprintf("signer address %s does not match expected owner address %s", signer.String(), ownerAddr.String())
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, errorString)

	}

	return nil
}

// extractSigners takes a proto.Message and returns a slice of signer addresses as strings.
func extractSigners(protoReflectMsg protoreflect.Message) ([]string, error) {

	descriptor := protoReflectMsg.Descriptor()
	signerFields, err := getSignerFieldNames(descriptor)
	if err != nil {
		return nil, err
	}

	var addresses []string
	for _, fieldName := range signerFields {
		field := descriptor.Fields().ByName(protoreflect.Name(fieldName))
		if field == nil {
			return nil, fmt.Errorf("field %s not found in message %s", fieldName, descriptor.FullName())
		}

		if field.Kind() != protoreflect.StringKind {
			return nil, fmt.Errorf("unexpected field type %s for field %s in message %s; only string fields are supported", field.Kind(), fieldName, descriptor.FullName())
		}

		fieldValue := protoReflectMsg.Get(field)
		if field.IsList() {
			list := fieldValue.List()
			for i := 0; i < list.Len(); i++ {
				addresses = append(addresses, list.Get(i).String())
			}
		} else {
			addresses = append(addresses, fieldValue.String())
		}
	}

	return addresses, nil
}

func (k Keeper) FlowIsToSourceChain(ctx sdk.Context, destinationChannelID, portID, flowConnectionID string, txMsgsAnys []*codectypes.Any, hostedAddress string) bool {
	ics20ConnectionID, _ := k.GetConnectionID(ctx, portID, destinationChannelID)

	// Case 1: All messages are MsgTransfer and match port/channel
	allTransfers := true
	for _, msg := range txMsgsAnys {
		if msg.TypeUrl != "/ibc.applications.transfer.v1.MsgTransfer" {
			allTransfers = false
			break
		}
		transferMsg, err := types.GetTransferMsg(k.cdc, msg)
		if err != nil {
			return false
		}

		if transferMsg.SourcePort != portID || transferMsg.SourceChannel != destinationChannelID {
			return false
		}
	}

	if allTransfers {
		return true
	}

	// Case 2: Hosted account, validate connection ID
	if hostedAddress != "" {
		hosted, err := k.TryGetHostedAccount(ctx, hostedAddress)
		if err != nil {
			return false
		}
		if hosted.ICAConfig.ConnectionID == ics20ConnectionID {
			return true
		}
	}

	// Case 3: ICA flow, validate connection ID
	if flowConnectionID == "" {
		return false
	}

	return ics20ConnectionID == flowConnectionID
}

func getSignerFieldNames(descriptor protoreflect.MessageDescriptor) ([]string, error) {
	// Retrieve the signer fields directly from the extension
	signersFields, ok := proto.GetExtension(descriptor.Options(), msgv1.E_Signer).([]string)
	if !ok || len(signersFields) == 0 {
		return nil, fmt.Errorf("no cosmos.msg.v1.signer option found for message %s; use DefineCustomGetSigners to specify a custom getter", descriptor.FullName())
	}

	return signersFields, nil
}

func unpackV2Any(cdc codec.Codec, msg *codectypes.Any) (protoreflect.Message, error) {
	msgv2, err := anyutil.Unpack(&anypb.Any{
		TypeUrl: msg.TypeUrl,
		Value:   msg.Value,
	}, cdc.InterfaceRegistry(), nil)
	if err != nil {
		return nil, err
	}

	return msgv2.ProtoReflect(), nil
}

func parseAccAddressFromAnyPrefix(bech32str string) (sdk.AccAddress, error) {
	if len(bech32str) == 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "address is empty")
	}

	_, bz, err := bech32.DecodeAndConvert(bech32str)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to decode Bech32 address")
	}

	return sdk.AccAddress(bz), nil
}
