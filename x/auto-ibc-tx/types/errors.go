package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrIBCAccountAlreadyExist = sdkerrors.Register(ModuleName, 2, "interchain account already registered")
	ErrIBCAccountNotExist     = sdkerrors.Register(ModuleName, 3, "interchain account not exist")
	ErrAccountExists          = sdkerrors.Register(ModuleName, 6, "fee account already exists")
	ErrDuplicate              = sdkerrors.Register(ModuleName, 14, "duplicate")
	ErrInvalid                = sdkerrors.Register(ModuleName, 1, "custom error message")
	ErrEmpty                  = sdkerrors.Register(ModuleName, 11, "empty")
	ErrAutoTxContinue         = sdkerrors.Register(ModuleName, 7, "max retries reached or not tx not acknowledged after execution on host chain (yet)")

	//ics20 hooks
	ErrMsgValidation        = sdkerrors.Register("ics20-hooks", 2, "error in ics20 hook message validation")
	ErrMarshaling           = sdkerrors.Register("ics20-hooks", 3, "cannot marshal the ICS20 packet")
	ErrInvalidPacket        = sdkerrors.Register("ics20-hooks", 4, "invalid packet data")
	ErrBadResponse          = sdkerrors.Register("ics20-hooks", 5, "cannot create response")
	ErrIcs20Error           = sdkerrors.Register("ics20-hooks", 6, "ics20 hook error")
	ErrBadSender            = sdkerrors.Register("ics20-hooks", 7, "bad sender")
	ErrBadMetadataFormatMsg = "auto_tx metadata not properly formatted for: '%v'. %s"
	ErrBadAutoTxMsg         = "cannot create autoTx: %v"
)
