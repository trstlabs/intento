package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrAccountExists  = sdkerrors.Register(ModuleName, 6, "fee account already exists")
	ErrDuplicate      = sdkerrors.Register(ModuleName, 14, "duplicate")
	ErrInvalidGenesis = sdkerrors.Register(ModuleName, 1, "invalid ids upon genesis")
	ErrEmpty          = sdkerrors.Register(ModuleName, 11, "empty")

	//ics20 hooks
	ErrMsgValidation = sdkerrors.Register("ics20-hooks", 20, "error in ics20 hook message validation")
	ErrMarshaling    = sdkerrors.Register("ics20-hooks", 21, "cannot marshal the ICS20 packet")
	ErrInvalidPacket = sdkerrors.Register("ics20-hooks", 22, "invalid packet data")
	ErrBadResponse   = sdkerrors.Register("ics20-hooks", 23, "cannot create response")
	ErrIcs20Error    = sdkerrors.Register("ics20-hooks", 24, "ics20 hook error")
	ErrBadSender     = sdkerrors.Register("ics20-hooks", 25, "bad sender")

	ErrInvalidTime  = sdkerrors.Register(ModuleName, 30, "time must be longer than 2 minutes from now")
	ErrUpdateAutoTx = sdkerrors.Register(ModuleName, 31, "cannot update auto_tx parameter")

	ErrBadMetadataFormatMsg = "metadata not properly formatted for: '%v'. %s"
	ErrBadAutoTxMsg         = "cannot create autoTx: %v"
	ErrAutoTxConditions     = "conditions to execute not met"
	ErrAutoTxDistribution   = "distribution error: %s"
	ErrAutoTxMsgHandling    = "msg handling error: %s"
)
