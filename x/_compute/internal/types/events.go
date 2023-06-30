package types

// Governance module event types
const (
	CustomContractEventPrefix      = "wasm-"
	EventTypeInstantiate           = "instantiate"
	EventTypeExecute               = "execute"
	EventTypeReply                 = "reply"
	EventTypeDistributedToContract = "DistributedContractIncentive"
	EventTypeGovContractResult     = "gov_contract_result"
	AttributeKeyAddress            = "address"
	AttributeKeyResultDataHex      = "result"
	AttributeReservedPrefix        = "_"
)
