package types

import (
	"encoding/base64"
	fmt "fmt"
	"strings"
	"time"

	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktxsigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/spf13/cast"
	wasmTypes "github.com/trstlabs/trst/go-cosmwasm/types"
)

var Denom = "utrst"

const defaultLRUCacheSize = uint64(0)
const defaultEnclaveLRUCacheSize = uint8(5) // can safely go up to 15
const defaultQueryGasLimit = uint64(3000000)

// base64 of a 64 byte key
type ContractKey string

func (m Model) ValidateBasic() error {
	if len(m.Key) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "key")
	}
	return nil
}

func (c CodeInfo) ValidateBasic() error {
	if len(c.CodeHash) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "code hash")
	}
	if err := sdk.VerifyAddressFormat(c.Creator); err != nil {
		return sdkerrors.Wrap(err, "creator")
	}
	if err := validateSourceURL(c.Source); err != nil {
		return sdkerrors.Wrap(err, "source")
	}
	if err := validateBuilder(c.Builder); err != nil {
		return sdkerrors.Wrap(err, "builder")
	}
	/*
		if err := c.InstantiateConfig.ValidateBasic(); err != nil {
			return sdkerrors.Wrap(err, "instantiate config")
		}
	*/
	return nil
}

// ParseEvents converts wasm LogAttributes into an sdk.Events (with 0 or 1 elements)
func ContractLogsToSdkEvents(logs []wasmTypes.Attribute, contractAddr sdk.AccAddress) sdk.Events {
	// we always tag with the contract address issuing this event
	attrs := []sdk.Attribute{sdk.NewAttribute(AttributeKeyContractAddr, contractAddr.String())}
	// append attributes from wasm to the sdk.Event
	for _, l := range logs {
		// and reserve the contract_address key for our use (not contract)
		if l.Key != AttributeKeyContractAddr {
			attr := sdk.NewAttribute(l.Key, string(l.Value))
			attrs = append(attrs, attr)
		}
	}
	// each wasm invokation always returns one sdk.Event
	return sdk.Events{sdk.NewEvent(CustomEventType, attrs...)}
}

const eventTypeMinLength = 2

// NewCustomEvents converts wasm events from a contract response to sdk type events
func NewCustomEvents(evts wasmTypes.Events, contractAddr sdk.AccAddress) (sdk.Events, error) {
	events := make(sdk.Events, 0, len(evts))
	for _, e := range evts {
		typ := strings.TrimSpace(e.Type)
		if len(typ) <= eventTypeMinLength {
			return nil, sdkerrors.Wrap(ErrInvalidEvent, fmt.Sprintf("Event type too short: '%s'", typ))
		}
		attributes, err := contractSDKEventAttributes(e.Attributes, contractAddr)
		if err != nil {
			return nil, err
		}
		events = append(events, sdk.NewEvent(fmt.Sprintf("%s%s", CustomContractEventPrefix, typ), attributes...))
	}
	return events, nil
}

// convert and add contract address issuing this event
func contractSDKEventAttributes(customAttributes []wasmTypes.Attribute, contractAddr sdk.AccAddress) ([]sdk.Attribute, error) {
	attrs := []sdk.Attribute{sdk.NewAttribute(AttributeKeyContractAddr, contractAddr.String())}
	// append attributes from wasm to the sdk.Event
	for _, l := range customAttributes {
		if l.PubDb {
			continue
		}
		// ensure key and value are non-empty (and trim what is there)
		key := strings.TrimSpace(l.Key)
		if len(key) == 0 {
			return nil, sdkerrors.Wrap(ErrInvalidEvent, fmt.Sprintf("Empty attribute key. Value: %s", l.Value))
		}
		value := strings.TrimSpace(string(l.Value))
		// TODO: check if this is legal in the SDK - if it is, we can remove this check
		if len(value) == 0 {
			return nil, sdkerrors.Wrap(ErrInvalidEvent, fmt.Sprintf("Empty attribute value. Key: %s", key))
		}
		// and reserve all _* keys for our use (not contract)
		if strings.HasPrefix(key, AttributeReservedPrefix) {
			return nil, sdkerrors.Wrap(ErrInvalidEvent, fmt.Sprintf("Attribute key starts with reserved prefix %s: '%s'", AttributeReservedPrefix, key))
		}
		attrs = append(attrs, sdk.NewAttribute(key, value))
	}
	return attrs, nil
}

// NewCodeInfo fills a new Contract struct
func NewCodeInfo(codeHash []byte, creator sdk.AccAddress, source string, builder string, default_duration time.Duration, default_interval time.Duration /* , instantiatePermission AccessConfig */, title string, description string) CodeInfo {
	return CodeInfo{
		CodeHash:        codeHash,
		Creator:         creator,
		Source:          source,
		Builder:         builder,
		DefaultDuration: default_duration,
		DefaultInterval: default_interval,
		Title:           title,
		Description:     description,
		Instances:       0,
		// InstantiateConfig: instantiatePermission,
	}
}

/*
type ContractCodeHistoryOperationType string

const (
	InitContractCodeHistoryType    ContractCodeHistoryOperationType = "Init"
	MigrateContractCodeHistoryType ContractCodeHistoryOperationType = "Migrate"
	GenesisContractCodeHistoryType ContractCodeHistoryOperationType = "Genesis"
)

var AllCodeHistoryTypes = []ContractCodeHistoryOperationType{InitContractCodeHistoryType, MigrateContractCodeHistoryType}

// ContractCodeHistoryEntry stores code updates to a contract.
type ContractCodeHistoryEntry struct {
	Operation ContractCodeHistoryOperationType `json:"operation"`
	CodeID    uint64                           `json:"code_id"`
	Updated   *AbsoluteTxPosition              `json:"updated,omitempty"`
	Msg       json.RawMessage                  `json:"msg,omitempty"`
}
*/

// NewContractInfo creates a new instance of a given WASM contract info
func NewContractInfo(codeID uint64, creator /* , admin */ sdk.AccAddress, label string, createdAt *AbsoluteTxPosition, startTime time.Time, execTime time.Time, endTime time.Time, duration time.Duration, interval time.Duration, autoMsg []byte, callbackSig []byte, owner sdk.AccAddress) ContractInfo {
	return ContractInfo{
		CodeID:  codeID,
		Creator: creator,
		Owner:   owner,
		// Admin:   admin,
		ContractId:  label,
		Created:     createdAt,
		StartTime:   startTime,
		ExecTime:    execTime,
		EndTime:     endTime,
		Duration:    duration,
		Interval:    interval,
		AutoMsg:     autoMsg,
		CallbackSig: callbackSig,
	}
}
func (c *ContractInfo) ValidateBasic() error {
	if c.CodeID == 0 {
		return sdkerrors.Wrap(ErrEmpty, "code id")
	}
	if err := sdk.VerifyAddressFormat(c.Creator); err != nil {
		return sdkerrors.Wrap(err, "creator")
	}
	/*
		if c.Admin != nil {
			if err := sdk.VerifyAddressFormat(c.Admin); err != nil {
				return sdkerrors.Wrap(err, "admin")
			}
		}
	*/
	if err := validateContractId(c.ContractId); err != nil {
		return sdkerrors.Wrap(err, "label")
	}
	return nil
}

/*
func (c ContractInfo) InitialHistory(msg []byte) ContractCodeHistoryEntry {
	return ContractCodeHistoryEntry{
		Operation: InitContractCodeHistoryType,
		CodeID:    c.CodeID,
		Updated:   c.Created,
		Msg:       msg,
	}
}

func (c *ContractInfo) AddMigration(ctx sdk.Context, codeID uint64, msg []byte) ContractCodeHistoryEntry {
	h := ContractCodeHistoryEntry{
		Operation: MigrateContractCodeHistoryType,
		CodeID:    codeID,
		Updated:   NewAbsoluteTxPosition(ctx),
		Msg:       msg,
	}
	c.CodeID = codeID
	return h
}

// ResetFromGenesis resets contracts timestamp and history.
func (c *ContractInfo) ResetFromGenesis(ctx sdk.Context) ContractCodeHistoryEntry {
	c.Created = NewAbsoluteTxPosition(ctx)
	return ContractCodeHistoryEntry{
		Operation: GenesisContractCodeHistoryType,
		CodeID:    c.CodeID,
		Updated:   c.Created,
	}
}
*/

// LessThan can be used to sort
func (a *AbsoluteTxPosition) LessThan(b *AbsoluteTxPosition) bool {
	if a == nil {
		return true
	}
	if b == nil {
		return false
	}
	return a.BlockHeight < b.BlockHeight || (a.BlockHeight == b.BlockHeight && a.TxIndex < b.TxIndex)
}

// NewAbsoluteTxPosition gets a timestamp from the context
func NewAbsoluteTxPosition(ctx sdk.Context) *AbsoluteTxPosition {
	// we must safely handle nil gas meters
	var index uint64
	meter := ctx.BlockGasMeter()
	if meter != nil {
		index = meter.GasConsumed()
	}
	return &AbsoluteTxPosition{
		BlockHeight: ctx.BlockHeight(),
		TxIndex:     index,
	}
}

// NewEnv initializes the environment for a contract instance
func NewEnv(ctx sdk.Context, creator sdk.AccAddress, deposit sdk.Coins, contractAddr sdk.AccAddress, contractKey []byte) wasmTypes.Env {
	// safety checks before casting below
	if ctx.BlockHeight() < 0 {
		panic("Block height must never be negative")
	}
	if ctx.BlockTime().Unix() < 0 {
		panic("Block (unix) time must never be negative ")
	}
	env := wasmTypes.Env{
		Block: wasmTypes.BlockInfo{
			Height:  uint64(ctx.BlockHeight()),
			Time:    uint64(ctx.BlockTime().Unix()),
			ChainID: ctx.ChainID(),
		},
		Message: wasmTypes.MessageInfo{
			Sender: creator.String(),
			Funds:  NewWasmCoins(deposit),
		},
		Contract: wasmTypes.ContractInfo{
			Address: contractAddr.String(),
		},
		Key: wasmTypes.ContractKey(base64.StdEncoding.EncodeToString(contractKey)),
	}
	return env
}

// NewWasmCoins translates between Cosmos SDK coins and Wasm coins
func NewWasmCoins(cosmosCoins sdk.Coins) (wasmCoins []wasmTypes.Coin) {
	for _, coin := range cosmosCoins {
		wasmCoin := wasmTypes.Coin{
			Denom:  coin.Denom,
			Amount: coin.Amount.String(),
		}
		wasmCoins = append(wasmCoins, wasmCoin)
	}
	return wasmCoins
}

const CustomEventType = "wasm"
const EventTypeContractExpired = "contract_expired"
const EventTypeAutoMsgContract = "eontract_executed"
const AttributeKeyContractAddr = "contract_address"

/*
// ParseEvents converts wasm Attributes into an sdk.Events (with 0 or 1 elements)
func ParseEvents(logs []wasmTypes.Attribute, contractAddr sdk.AccAddress) sdk.Events {
	if len(logs) == 0 {
		return nil
	}
	// we always tag with the contract address issuing this event
	attrs := []sdk.Attribute{sdk.NewAttribute(AttributeKeyContractAddr, contractAddr.String())}
	for _, l := range logs {
		fmt.Printf("Log Key: %v \n", l.Key)
		if l.PubDb {
			continue
		}
		// and reserve the contract_address key for our use (not contract)
		if string(l.Key) != AttributeKeyContractAddr {
			attr := sdk.NewAttribute(l.Key, base64.StdEncoding.EncodeToString(l.Value))
			attrs = append(attrs, attr)
		}
	}
	return sdk.Events{sdk.NewEvent(CustomEventType, attrs...)}
}*/

// WasmConfig is the extra config required for wasm
type WasmConfig struct {
	SmartQueryGasLimit uint64
	CacheSize          uint64
	EnclaveCacheSize   uint8
}

// DefaultWasmConfig returns the default settings for WasmConfig
func DefaultWasmConfig() *WasmConfig {
	return &WasmConfig{
		SmartQueryGasLimit: defaultQueryGasLimit,
		CacheSize:          defaultLRUCacheSize,
		EnclaveCacheSize:   defaultEnclaveLRUCacheSize,
	}
}

type ContractMsg struct {
	CodeHash []byte
	Msg      []byte
}

func (m ContractMsg) Serialize() []byte {
	return append(m.CodeHash, m.Msg...)
}

func NewVerificationInfo(
	signBytes []byte, signMode sdktxsigning.SignMode, modeInfo []byte, publicKey []byte, signature []byte, callbackSig []byte,
) wasmTypes.VerificationInfo {
	return wasmTypes.VerificationInfo{
		Bytes:             signBytes,
		SignMode:          signMode.String(),
		ModeInfo:          modeInfo,
		Signature:         signature,
		PublicKey:         publicKey,
		CallbackSignature: callbackSig,
	}
}

func NewMsgInfo(
	code_hash []byte, funds sdk.Coins,
) wasmTypes.MsgInfo {
	wasmFunds := NewWasmCoins(funds)
	return wasmTypes.MsgInfo{
		CodeHash: code_hash,
		Funds:    wasmFunds,
	}
}

type ParseAuto struct {
	AutoMsg struct {
	} `json:"auto_msg"`
}

// GetConfig load config values from the app options
func GetConfig(appOpts servertypes.AppOptions) *WasmConfig {
	return &WasmConfig{
		SmartQueryGasLimit: cast.ToUint64(appOpts.Get("wasm.contract-query-gas-limit")),
		CacheSize:          cast.ToUint64(appOpts.Get("wasm.contract-memory-cache-size")),
		EnclaveCacheSize:   cast.ToUint8(appOpts.Get("wasm.contract-memory-enclave-cache-size")),
	}
}

// DefaultConfigTemplate default config template for wasm module
const DefaultConfigTemplate = `
[wasm]
# The maximum gas amount can be spent for contract query.
# The contract query will invoke contract execution vm,
# so we need to restrict the max usage to prevent DoS attack
contract-query-gas-limit = "{{ .WASMConfig.SmartQueryGasLimit }}"

# The WASM VM memory cache size in MiB not bytes
contract-memory-cache-size = "{{ .WASMConfig.CacheSize }}"

# The WASM VM memory cache size in number of cached modules. Can safely go up to 15, but not recommended for validators
contract-memory-enclave-cache-size = "{{ .WASMConfig.EnclaveCacheSize }}"
`
