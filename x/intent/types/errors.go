package types

import (
	"cosmossdk.io/errors"
)

// Module error codes (1xxx range)
const (
	// 1xx - General validation errors
	codeErrInvalidGenesis = 100
	codeErrInvalidRequest = 101
	codeErrUnknownRequest = 102
	codeErrInvalidType    = 103
	codeErrInvalidAddress = 104
	codeErrInvalidTime    = 105
	codeErrEmpty          = 106

	// 2xx - Authorization and permission errors
	codeErrUnauthorized          = 200
	codeErrSignerNotOk           = 201
	codeErrTrustlessAgentLimit   = 202
	codeErrInvalidTrustlessAgent = 203

	// 3xx - State and storage errors
	codeErrAccountExists = 300
	codeErrDuplicate     = 301
	codeErrNotFound      = 302
	codeErrUpdateFlow    = 303

	// 4xx - Processing errors
	codeErrUnexpectedFeeCalculation = 400
	codeErrAckErr                   = 401
	codeErrAnyUnmarshal             = 402
	codeErrMsgResponsesHandling     = 403

	// 5xx - Flow and scheduling errors
	codeErrInvalidFeedbackLoop = 500
	codeErrInvalidICQTimeout   = 501
	codeErrInvalidScheduling   = 502
)

// ICS20 Hooks error codes (2xxx range)
const (
	codeIcs20ErrMsgValidation     = 2000
	codeIcs20ErrMsgMsgsValidation = 2001
	codeIcs20ErrMarshaling        = 2002
	codeIcs20ErrInvalidPacket     = 2003
	codeIcs20ErrBadResponse       = 2004
	codeIcs20ErrBadSender         = 2006
)

// Module errors
var (
	// General validation errors (1xx)
	ErrInvalidGenesis = errors.Register(ModuleName, codeErrInvalidGenesis, "invalid ids upon genesis")
	ErrInvalidRequest = errors.Register(ModuleName, codeErrInvalidRequest, "invalid request")
	ErrUnknownRequest = errors.Register(ModuleName, codeErrUnknownRequest, "unknown request")
	ErrInvalidType    = errors.Register(ModuleName, codeErrInvalidType, "invalid type")
	ErrInvalidAddress = errors.Register(ModuleName, codeErrInvalidAddress, "invalid address")
	ErrInvalidTime    = errors.Register(ModuleName, codeErrInvalidTime, "time must be more than 1 minute from now")
	ErrEmpty          = errors.Register(ModuleName, codeErrEmpty, "empty")

	// Authorization and permission errors (2xx)
	ErrUnauthorized           = errors.Register(ModuleName, codeErrUnauthorized, "unauthorized")
	ErrSignerNotOk            = errors.Register(ModuleName, codeErrSignerNotOk, "message signer is not message sender")
	ErrTrustlessAgentFeeLimit = errors.Register(ModuleName, codeErrTrustlessAgentLimit, "fee limit reached for trustless agent")
	ErrInvalidTrustlessAgent  = errors.Register(ModuleName, codeErrInvalidTrustlessAgent, "invalid msg response")

	// State and storage errors (3xx)
	ErrAccountExists = errors.Register(ModuleName, codeErrAccountExists, "fee account already exists")
	ErrDuplicate     = errors.Register(ModuleName, codeErrDuplicate, "duplicate")
	ErrNotFound      = errors.Register(ModuleName, codeErrNotFound, "not found")
	ErrUpdateFlow    = errors.Register(ModuleName, codeErrUpdateFlow, "cannot update flow parameter")

	// Processing errors (4xx)
	ErrUnexpectedFeeCalculation = errors.Register(ModuleName, codeErrUnexpectedFeeCalculation, "unexpected error during fee calculation")
	ErrAckErr                   = errors.Register(ModuleName, codeErrAckErr, "acknowledgement error")
	ErrAnyUnmarshal             = errors.Register(ModuleName, codeErrAnyUnmarshal, "failed unmarshal proto any")
	ErrMsgResponsesHandling     = errors.Register(ModuleName, codeErrMsgResponsesHandling, "error handling msg responses")

	// Flow and scheduling errors (5xx)
	ErrInvalidFeedbackLoop = errors.Register(ModuleName, codeErrInvalidFeedbackLoop, "invalid feedback loop")
	ErrInvalidICQTimeout   = errors.Register(ModuleName, codeErrInvalidICQTimeout, "invalid icq timeout")
	ErrInvalidScheduling   = errors.Register(ModuleName, codeErrInvalidScheduling, "invalid scheduling")
)

// ICS20 Hooks errors
var (
	ErrMsgValidation     = errors.Register("ics20-hooks", codeIcs20ErrMsgValidation, "error in ics20 hook message validation")
	ErrMsgMsgsValidation = errors.Register("ics20-hooks", codeIcs20ErrMsgMsgsValidation, "error in ics20 hook message validation for a message inside the flow")
	ErrMarshaling        = errors.Register("ics20-hooks", codeIcs20ErrMarshaling, "cannot marshal the ICS20 packet")
	ErrInvalidPacket     = errors.Register("ics20-hooks", codeIcs20ErrInvalidPacket, "invalid packet data")
	ErrBadResponse       = errors.Register("ics20-hooks", codeIcs20ErrBadResponse, "cannot create response")
	ErrBadSender         = errors.Register("ics20-hooks", codeIcs20ErrBadSender, "bad sender")
)

// Flow execution error messages (used for FlowHistory entries)
// These are not registered errors but format strings used in flow execution
const (
	ErrBadMetadataFormatMsg = "metadata not properly formatted for: '%v'. %s"
	ErrFlowConditions       = "conditions to execute not met: %v"
	ErrFlowFeeDistribution  = "distribution error: %s"
	ErrFlowMsgHandling      = "msg handling error: %s"
	ErrFlowResponseUseValue = "msg handling error using response value: %s"
	ErrSettingFlowResult    = "setting flow result: %s"
	ErrBalanceTooLow        = "balance too low to deduct expected fee"
)
