package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	types "github.com/trstlabs/trst/go-cosmwasm/types"
)

// IBCContractKeeper IBC lifecycle event handler
type IBCContractKeeper interface {
	OnOpenChannel(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		msg types.IBCChannelOpenMsg,
	) (string, error)
	OnConnectChannel(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		msg types.IBCChannelConnectMsg,
	) error
	OnCloseChannel(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		msg types.IBCChannelCloseMsg,
	) error
	OnRecvPacket(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		msg types.IBCPacketReceiveMsg,
	) ([]byte, error)
	OnAckPacket(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		acknowledgement types.IBCPacketAckMsg,
	) error
	OnTimeoutPacket(
		ctx sdk.Context,
		contractAddr sdk.AccAddress,
		msg types.IBCPacketTimeoutMsg,
	) error
	// ClaimCapability allows the transfer module to claim a capability
	// that IBC module passes to it
	ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error
	// AuthenticateCapability wraps the scopedKeeper's AuthenticateCapability function
	AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool
}
