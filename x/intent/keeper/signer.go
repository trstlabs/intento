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

func (k Keeper) SignerOk(ctx sdk.Context, codec codec.Codec, actionInfo types.ActionInfo) error {
	for _, message := range actionInfo.Msgs {
		if err := k.validateMessage(ctx, codec, actionInfo, message); err != nil {
			return err
		}
	}
	return nil
}

// validateMessage handles validation for a single message.
//   - Local messages: Validates that the signer matches the action owner.
//   - Authz MsgExec messages: Validates the signers for all contained messages.
//   - ICA messages: Skips validation because signer authentication happens
//     via IBC for the packet sender. The user controls which ICA to use,
//     and the expected configuration is already established within the module.
//     Additional verification is unnecessary.
func (k Keeper) validateMessage(ctx sdk.Context, codec codec.Codec, actionInfo types.ActionInfo, message *codectypes.Any) error {
	var sdkMsg sdk.Msg
	if err := codec.UnpackAny(message, &sdkMsg); err != nil {
		return errorsmod.Wrap(err, "failed to unpack message")
	}

	if isAuthzMsgExec(message) {
		// Validate MsgExec messages for nested signer verification.
		return k.validateAuthzMsg(ctx, codec, actionInfo, message)
	}

	if isLocalMessage(actionInfo) {
		// Validate local messages by ensuring the signer matches the action owner.
		return k.validateSigners(ctx, codec, actionInfo, message)
	}

	// ICA messages are trusted due to IBC authentication and controlled configuration.
	// No further validation is needed for these cases.
	return nil
}

func isAuthzMsgExec(message *codectypes.Any) bool {
	return message.TypeUrl == sdk.MsgTypeURL(&authztypes.MsgExec{})
}

func isLocalMessage(actionInfo types.ActionInfo) bool {
	return (actionInfo.ICAConfig == nil || actionInfo.ICAConfig.ConnectionID == "") && (actionInfo.HostedConfig == nil || actionInfo.HostedConfig.HostedAddress == "")
}

// validateAuthzMsg validates an authz MsgExec message.
func (k Keeper) validateAuthzMsg(ctx sdk.Context, codec codec.Codec, actionInfo types.ActionInfo, message *codectypes.Any) error {
	var authzMsg authztypes.MsgExec
	if err := cosmosproto.Unmarshal(message.Value, &authzMsg); err != nil {
		return errorsmod.Wrap(err, "failed to unmarshal MsgExec")
	}

	for _, innerMessage := range authzMsg.Msgs {
		if err := k.validateSigners(ctx, codec, actionInfo, innerMessage); err != nil {
			return err
		}
	}
	return nil
}

// validateSigners checks the signers of a message against the owner, ICA, and hosted accounts.
func (k Keeper) validateSigners(ctx sdk.Context, codec codec.Codec, actionInfo types.ActionInfo, message *codectypes.Any) error {

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
	ownerAddr, err := sdk.AccAddressFromBech32(actionInfo.Owner)
	if err != nil {
		return errorsmod.Wrap(err, "failed to parse owner address")
	}
	signer, err := parseAccAddressFromAnyPrefix(signers[0])
	if err != nil {
		return errorsmod.Wrap(err, "failed to parse owner address")
	}
	// fmt.Printf("Owner %s \n", actionInfo.Owner)
	// fmt.Printf("Signer %s \n", signers[0])
	k.Logger(ctx).Debug("Signer validation", "owner", actionInfo.Owner, "signer", signers[0])
	if !signer.Equals(ownerAddr) {
		//if !checkICA {
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "signer address does not match expected owner address")
		//}
		//return k.validateHostedOrICAAccount(ctx, actionInfo, signer)
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

// validateHostedOrICAAccount checks if the signer matches a hosted or ICA account.
// func (k Keeper) validateHostedOrICAAccount(ctx sdk.Context, actionInfo types.ActionInfo, signer sdk.AccAddress) error {
// 	// Check Hosted Config
// 	if actionInfo.HostedConfig != nil && actionInfo.HostedConfig.HostedAddress != "" {
// 		ica, err := k.TryGetHostedAccount(ctx, actionInfo.HostedConfig.HostedAddress)
// 		if err != nil {
// 			return errorsmod.Wrap(err, "failed to get hosted account")
// 		}

// 		hostedAccAddr, err := parseAccAddressFromAnyPrefix(ica.HostedAddress)
// 		if err != nil {
// 			return errorsmod.Wrap(err, "failed to parse hosted address")
// 		}
// 		if signer.Equals(hostedAccAddr) {
// 			return nil
// 		}
// 	}

// 	// Check ICA Account
// 	icaAddrString, found := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, actionInfo.ICAConfig.ConnectionID, actionInfo.ICAConfig.PortID)
// 	if !found {
// 		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "ICA account not found")
// 	}
// 	icaAddr, err := parseAccAddressFromAnyPrefix(icaAddrString)
// 	if err != nil {
// 		return errorsmod.Wrap(err, "failed to parse ica address")
// 	}
// 	if !signer.Equals(icaAddr) {
// 		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "signer does not match any authorized account")

// 	}

// 	return nil
// }
