package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrUnauthorized             = errorsmod.Register(ModuleName, 2, "unauthorized")
	ErrAccountExists            = errorsmod.Register(ModuleName, 6, "fee account already exists")
	ErrDuplicate                = errorsmod.Register(ModuleName, 14, "duplicate")
	ErrSignerNotOk              = errorsmod.Register(ModuleName, 19, "message signer is not message sender")
	ErrInvalidGenesis           = errorsmod.Register(ModuleName, 1, "invalid ids upon genesis")
	ErrEmpty                    = errorsmod.Register(ModuleName, 11, "empty")
	ErrInvalidRequest           = errorsmod.Register(ModuleName, 8, "invalid request")
	ErrUnknownRequest           = errorsmod.Register(ModuleName, 4, "unknown request")
	ErrUnexpectedFeeCalculation = errorsmod.Register(ModuleName, 3, "unexpected error during fee calculation")
	ErrAckErr                   = errorsmod.Register(ModuleName, 33, "acknowledgement error")
	ErrNotFound                 = errorsmod.Register(ModuleName, 16, "not found")
	ErrHostedFeeLimit           = errorsmod.Register(ModuleName, 17, "fee limit reached for hosted account")
	ErrInvalidType              = errorsmod.Register(ModuleName, 12, "invalid type")
	ErrInvalidAddress           = errorsmod.Register(ModuleName, 5, "invalid address")
	ErrJSONUnmarshal            = errorsmod.Register(ModuleName, 15, "failed unmarshal json")
	ErrMsgResponsesHandling     = errorsmod.Register(ModuleName, 18, "error handling msg responses")
	//ics20 hooks
	ErrMsgValidation = errorsmod.Register("ics20-hooks", 20, "error in ics20 hook message validation")
	ErrMarshaling    = errorsmod.Register("ics20-hooks", 21, "cannot marshal the ICS20 packet")
	ErrInvalidPacket = errorsmod.Register("ics20-hooks", 22, "invalid packet data")
	ErrBadResponse   = errorsmod.Register("ics20-hooks", 23, "cannot create response")
	ErrIcs20Error    = errorsmod.Register("ics20-hooks", 24, "ics20 hook error")
	ErrBadSender     = errorsmod.Register("ics20-hooks", 25, "bad sender")

	ErrInvalidTime            = errorsmod.Register(ModuleName, 30, "time must be longer than 2 minutes from now")
	ErrUpdateFlow             = errorsmod.Register(ModuleName, 31, "cannot update flow parameter")
	ErrValidateMsgRegistryMsg = errorsmod.Register(ModuleName, 32, "could not validate Flow message")

	//errors specific to Flow execution that are to be appended to FlowHistory entries
	ErrBadMetadataFormatMsg = "metadata not properly formatted for: '%v'. %s"
	ErrBadFlowMsg           = "cannot create flow: %v"
	ErrFlowConditions       = "conditions to execute not met: %v"
	ErrFlowFeeDistribution  = "distribution error: %s"
	ErrFlowMsgHandling      = "msg handling error: %s"
	ErrFlowResponseUseValue = "msg handling error using response value: %s"
	ErrSettingFlowResult    = "setting flow result:  %s"
	ErrBalanceTooLow        = "balance too low to deduct expected fee"
)
