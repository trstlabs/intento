package types

import (
	"encoding/json"
)

//-------- Queries --------

type QueryResponse struct {
	Ok  []byte    `json:"Ok,omitempty"`
	Err *StdError `json:"Err,omitempty"`
}

//-------- Querier -----------

type Querier interface {
	Query(request QueryRequest, gasLimit uint64) ([]byte, error)
	GasConsumed() uint64
}

// this is a thin wrapper around the desired Go API to give us types closer to Rust FFI
func RustQuery(querier Querier, binRequest []byte, gasLimit uint64) QuerierResult {
	var request QueryRequest
	err := json.Unmarshal(binRequest, &request)
	if err != nil {
		return ToQuerierResult(nil, UnsupportedRequest{err.Error()})
	}
	bz, err := querier.Query(request, gasLimit)
	return ToQuerierResult(bz, err)
}

// This is a 2-level result
type QuerierResult struct {
	Ok  *QueryResponse `json:"Ok,omitempty"`
	Err *SystemError   `json:"Err,omitempty"`
}

func ToQuerierResult(response []byte, err error) QuerierResult {
	if err == nil {
		return QuerierResult{
			Ok: &QueryResponse{
				Ok: response,
			},
		}
	}
	syserr := ToSystemError(err)
	if syserr != nil {

		return QuerierResult{
			Err: syserr,
		}
	}
	stderr := ToStdError(err)

	return QuerierResult{
		Ok: &QueryResponse{
			Err: stderr,
		},
	}
}

// QueryRequest is an rust enum and only (exactly) one of the fields should be set
// Should we do a cleaner approach in Go? (type/data?)
type QueryRequest struct {
	Bank     *BankQuery      `json:"bank,omitempty"`
	Custom   json.RawMessage `json:"custom,omitempty"`
	Staking  *StakingQuery   `json:"staking,omitempty"`
	Wasm     *WasmQuery      `json:"wasm,omitempty"`
	Dist     *DistQuery      `json:"dist,omitempty"`
	Mint     *MintQuery      `json:"mint,omitempty"`
	Gov      *GovQuery       `json:"gov,omitempty"`
	IBC      *IBCQuery       `json:"ibc,omitempty"`
	Stargate *StargateQuery  `json:"stargate,omitempty"`
}

type BankQuery struct {
	Balance     *BalanceQuery     `json:"balance,omitempty"`
	AllBalances *AllBalancesQuery `json:"all_balances,omitempty"`
}

type BalanceQuery struct {
	Address string `json:"address"`
	Denom   string `json:"denom"`
}

// BalanceResponse is the expected response to BalanceQuery
type BalanceResponse struct {
	Amount Coin `json:"amount"`
}

type AllBalancesQuery struct {
	Address string `json:"address"`
}

// AllBalancesResponse is the expected response to AllBalancesQuery
type AllBalancesResponse struct {
	Amount Coins `json:"amount"`
}

type StakingQuery struct {
	Validators           *ValidatorsQuery         `json:"validators,omitempty"`
	Validator            *ValidatorQuery          `json:"validator,omitempty"`
	AllDelegations       *AllDelegationsQuery     `json:"all_delegations,omitempty"`
	Delegation           *DelegationQuery         `json:"delegation,omitempty"`
	UnBondingDelegations *UnbondingDeletionsQuery `json:"unbonding_delegations, omitempty"`
	BondedDenom          *struct{}                `json:"bonded_denom,omitempty"`
}

type ValidatorQuery struct {
	/// Address is the validator's address (e.g. cosmosvaloper1...)
	Address string `json:"address"`
}

// ValidatorResponse is the expected response to ValidatorQuery
type ValidatorResponse struct {
	Validator *Validator `json:"validator"` // serializes to `null` when unset which matches Rust's Option::None serialization
}

type UnbondingDeletionsQuery struct {
	Delegator string `json:"delegator"`
}

type ValidatorsQuery struct{}

// ValidatorsResponse is the expected response to ValidatorsQuery
type ValidatorsResponse struct {
	Validators Validators `json:"validators"`
}

// TODO: Validators must JSON encode empty array as []
type Validators []Validator

// MarshalJSON ensures that we get [] for empty arrays
func (v Validators) MarshalJSON() ([]byte, error) {
	if len(v) == 0 {
		return []byte("[]"), nil
	}
	var raw []Validator = v
	return json.Marshal(raw)
}

// UnmarshalJSON ensures that we get [] for empty arrays
func (v *Validators) UnmarshalJSON(data []byte) error {
	// make sure we deserialize [] back to null
	if string(data) == "[]" || string(data) == "null" {
		return nil
	}
	var raw []Validator
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*v = raw
	return nil
}

type Validator struct {
	Address string `json:"address"`
	// decimal string, eg "0.02"
	Commission string `json:"commission"`
	// decimal string, eg "0.02"
	MaxCommission string `json:"max_commission"`
	// decimal string, eg "0.02"
	MaxChangeRate string `json:"max_change_rate"`
}

type AllDelegationsQuery struct {
	Delegator string `json:"delegator"`
}

type DelegationQuery struct {
	Delegator string `json:"delegator"`
	Validator string `json:"validator"`
}

// AllDelegationsResponse is the expected response to AllDelegationsQuery
type AllDelegationsResponse struct {
	Delegations Delegations `json:"delegations"`
}

type Delegations []Delegation

// MarshalJSON ensures that we get [] for empty arrays
func (d Delegations) MarshalJSON() ([]byte, error) {
	if len(d) == 0 {
		return []byte("[]"), nil
	}
	var raw []Delegation = d
	return json.Marshal(raw)
}

// UnmarshalJSON ensures that we get [] for empty arrays
func (d *Delegations) UnmarshalJSON(data []byte) error {
	// make sure we deserialize [] back to null
	if string(data) == "[]" || string(data) == "null" {
		return nil
	}
	var raw []Delegation
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*d = raw
	return nil
}

type Delegation struct {
	Delegator string `json:"delegator"`
	Validator string `json:"validator"`
	Amount    Coin   `json:"amount"`
}

// DelegationResponse is the expected response to DelegationsQuery
type DelegationResponse struct {
	Delegation *FullDelegation `json:"delegation,omitempty"`
}

type FullDelegation struct {
	Delegator          string `json:"delegator"`
	Validator          string `json:"validator"`
	Amount             Coin   `json:"amount"`
	AccumulatedRewards Coins  `json:"accumulated_rewards"`
	CanRedelegate      Coin   `json:"can_redelegate"`
}

type UnbondingDelegationsResponse struct {
	Delegations Delegations `json:"delegations"`
}

type BondedDenomResponse struct {
	Denom string `json:"denom"`
}

type WasmQuery struct {
	Private       *PrivateQuery    `json:"private,omitempty"`
	Public        *PublicQuery     `json:"public,omitempty"`
	PublicForAddr *PublicAddrQuery `json:"public_for_addr,omitempty"`
}

// PrivateQuery queries the private contract state
type PrivateQuery struct {
	ContractAddr string `json:"contract_addr"`
	Msg          []byte `json:"msg"`
}

// PublicQuery queries the public contract state
type PublicQuery struct {
	ContractAddr string `json:"contract_addr"`
	Key          string `json:"key"`
}

// PublicQuery queries the public contract state
type PublicAddrQuery struct {
	ContractAddr string `json:"contract_addr"`
	AccountAddr  string `json:"account_addr"`
	Key          string `json:"key"`
}

type DistQuery struct {
	Rewards *RewardsQuery `json:"rewards,omitempty"`
}

type GovQuery struct {
	Proposals *ProposalsQuery `json:"proposals,omitempty"`
}
type MintQuery struct {
	Inflation   *MintingInflationQuery   `json:"inflation,omitempty"`
	BondedRatio *MintingBondedRatioQuery `json:"bonded_ratio,omitempty"`
}

type MintingBondedRatioQuery struct{}
type MintingInflationQuery struct{}

type MintingInflationResponse struct {
	InflationRate string `json:"inflation_rate"`
}

type MintingBondedRatioResponse struct {
	BondedRatio string `json:"bonded_ratio"`
}

type ProposalsQuery struct{}

// DelegationResponse is the expected response to DelegationsQuery
type ProposalsResponse struct {
	Proposals []Proposal `json:"proposals,omitempty"`
}

type Proposal struct {
	ProposalID      uint64 `json:"id" yaml:"id"`                               //  ID of the proposal
	VotingStartTime uint64 `json:"voting_start_time" yaml:"voting_start_time"` // Time of the block where MinDeposit was reached. -1 if MinDeposit is not reached
	VotingEndTime   uint64 `json:"voting_end_time" yaml:"voting_end_time"`     // Time that the VotingPeriod for this proposal will end and votes will be tallied
}

type RewardsQuery struct {
	Delegator string `json:"delegator"`
}

// DelegationResponse is the expected response to DelegationsQuery
type RewardsResponse struct {
	Rewards []Rewards   `json:"rewards,omitempty"`
	Total   RewardCoins `json:"total,omitempty"`
}

type Rewards struct {
	Validator string      `json:"validator_address"`
	Reward    RewardCoins `json:"reward"`
}

type RewardCoins []Coin

// MarshalJSON ensures that we get [] for empty arrays
func (d RewardCoins) MarshalJSON() ([]byte, error) {
	if len(d) == 0 {
		return []byte("[]"), nil
	}
	var raw []Coin = d
	return json.Marshal(raw)
}

// UnmarshalJSON ensures that we get [] for empty arrays
func (d *RewardCoins) UnmarshalJSON(data []byte) error {
	// make sure we deserialize [] back to null
	if string(data) == "[]" || string(data) == "null" {
		return nil
	}
	var raw []Coin
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*d = raw
	return nil
}

// MarshalJSON ensures that we get [] for empty arrays
func (d ProposalsResponse) MarshalJSON() ([]byte, error) {
	if len(d.Proposals) == 0 {
		return []byte("{\"proposals\": []}"), nil
	}
	var raw = d.Proposals
	asBytes, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	return append(append([]byte("{\"proposals\": "), asBytes...), []byte("}")...), nil
}

// UnmarshalJSON ensures that we get [] for empty arrays
func (d *ProposalsResponse) UnmarshalJSON(data []byte) error {
	// make sure we deserialize [] back to null
	if string(data) == "{\"proposals\": []}" || string(data) == "null" || string(data) == "{\"proposals\":[]}" {
		return nil
	}
	var raw []Proposal
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	d.Proposals = raw
	return nil
}

// StargateQuery is encoded the same way as abci_query, with path and protobuf encoded request data.
// The format is defined in [ADR-21](https://github.com/cosmos/cosmos-sdk/blob/master/docs/architecture/adr-021-protobuf-query-encoding.md).
// The response is protobuf encoded data directly without a JSON response wrapper.
// The caller is responsible for compiling the proper protobuf definitions for both requests and responses.
type StargateQuery struct {
	// this is the fully qualified service path used for routing,
	// eg. custom/cosmos_sdk.x.bank.v1.Query/QueryBalance
	Path string `json:"path"`
	// this is the expected protobuf message type (not any), binary encoded
	Data []byte `json:"data"`
}

// IBCQuery defines a query request from the contract into the chain.
// This is the counterpart of [IbcQuery](https://github.com/CosmWasm/cosmwasm/blob/v0.14.0-beta1/packages/std/src/ibc.rs#L61-L83).
type IBCQuery struct {
	PortID       *PortIDQuery       `json:"port_id,omitempty"`
	ListChannels *ListChannelsQuery `json:"list_channels,omitempty"`
	Channel      *ChannelQuery      `json:"channel,omitempty"`
}

type PortIDQuery struct{}

type PortIDResponse struct {
	PortID string `json:"port_id"`
}

// ListChannelsQuery is an IBCQuery that lists all channels that are bound to a given port.
// If `PortID` is unset, this list all channels bound to the contract's port.
// Returns a `ListChannelsResponse`.
// This is the counterpart of [IbcQuery::ListChannels](https://github.com/CosmWasm/cosmwasm/blob/v0.14.0-beta1/packages/std/src/ibc.rs#L70-L73).
type ListChannelsQuery struct {
	// optional argument
	PortID string `json:"port_id,omitempty"`
}

type ListChannelsResponse struct {
	Channels IBCChannels `json:"channels"`
}

// IBCChannels must JSON encode empty array as [] (not null) for consistency with Rust parser
type IBCChannels []IBCChannel

// MarshalJSON ensures that we get [] for empty arrays
func (e IBCChannels) MarshalJSON() ([]byte, error) {
	if len(e) == 0 {
		return []byte("[]"), nil
	}
	var raw []IBCChannel = e
	return json.Marshal(raw)
}

// UnmarshalJSON ensures that we get [] for empty arrays
func (e *IBCChannels) UnmarshalJSON(data []byte) error {
	// make sure we deserialize [] back to null
	if string(data) == "[]" || string(data) == "null" {
		return nil
	}
	var raw []IBCChannel
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*e = raw
	return nil
}

// IBCEndpoints must JSON encode empty array as [] (not null) for consistency with Rust parser
type IBCEndpoints []IBCEndpoint

// MarshalJSON ensures that we get [] for empty arrays
func (e IBCEndpoints) MarshalJSON() ([]byte, error) {
	if len(e) == 0 {
		return []byte("[]"), nil
	}
	var raw []IBCEndpoint = e
	return json.Marshal(raw)
}

// UnmarshalJSON ensures that we get [] for empty arrays
func (e *IBCEndpoints) UnmarshalJSON(data []byte) error {
	// make sure we deserialize [] back to null
	if string(data) == "[]" || string(data) == "null" {
		return nil
	}
	var raw []IBCEndpoint
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*e = raw
	return nil
}

type ChannelQuery struct {
	// optional argument
	PortID    string `json:"port_id,omitempty"`
	ChannelID string `json:"channel_id"`
}

type ChannelResponse struct {
	// may be empty if there is no matching channel
	Channel *IBCChannel `json:"channel,omitempty"`
}
