package keeper

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strings"
	"testing"
	"time"

	stypes "github.com/cosmos/cosmos-sdk/store/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/libs/log"

	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cosmwasm "github.com/trstlabs/trst/go-cosmwasm/types"

	"github.com/trstlabs/trst/x/compute/internal/types"
)

type ContractEvent []cosmwasm.Attribute

type TestContract struct {
	CosmWasmVersion string
	IsCosmWasmV1    bool
	WasmFilePath    string
}

var testContract = TestContract{
	WasmFilePath: "./testdata/test-contract/contract.wasm",
}

// if codeID isn't 0, it will try to use that. Otherwise will take the contractAddress
func testEncrypt(t *testing.T, keeper Keeper, ctx sdk.Context, contractAddress sdk.AccAddress, codeId uint64, msg []byte) ([]byte, error) {
	var hash []byte
	if codeId != 0 {
		inf, _ := keeper.GetCodeInfo(ctx, codeId)
		hash = inf.CodeHash
	} else {
		hash, _ = keeper.GetContractHash(ctx, contractAddress)
	}

	if hash == nil {
		return nil, cosmwasm.StdError{}
	}

	intMsg := types.ContractMsg{
		CodeHash: []byte(hex.EncodeToString(hash)),
		Msg:      msg,
	}

	queryBz, err := wasmCtx.Encrypt(intMsg.Serialize())
	require.NoError(t, err)

	return queryBz, nil
}

func setupTest(t *testing.T, wasmPath string, additionalCoinsInWallets sdk.Coins) (sdk.Context, Keeper, uint64, string, sdk.AccAddress, crypto.PrivKey, sdk.AccAddress, crypto.PrivKey) {
	encodingConfig := MakeEncodingConfig()
	var transferPortSource types.ICS20TransferPortSource
	transferPortSource = MockIBCTransferKeeper{GetPortFn: func(ctx sdk.Context) string {
		return "myTransferPort"
	}}
	encoders := DefaultEncoders(transferPortSource, encodingConfig.Marshaler)
	ctx, keepers := CreateTestInput(t, false, SupportedFeatures, &encoders, nil)
	accKeeper, keeper := keepers.AccountKeeper, keepers.WasmKeeper

	walletA, privKeyA := CreateFakeFundedAccount(ctx, accKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("denom", 200000)).Add(additionalCoinsInWallets...))
	walletB, privKeyB := CreateFakeFundedAccount(ctx, accKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("denom", 5000)).Add(additionalCoinsInWallets...))

	wasmCode, err := ioutil.ReadFile(wasmPath)
	require.NoError(t, err)

	codeID, err := keeper.Create(ctx, walletA, wasmCode, "", "", 0, 0, "title", "descr")
	require.NoError(t, err)

	info, _ := keeper.GetCodeInfo(ctx, codeID)
	codeHash := hex.EncodeToString(info.CodeHash)

	keeper.SetParams(ctx, types.Params{
		AutoMsgFundsCommission:          2,
		AutoMsgConstantFee:              1_000_000,                 // 1trst
		AutoMsgFlexFeeMul:               100,                       // 100/100 = 1 = gasUsed
		RecurringAutoMsgConstantFee:     1_000_000,                 // 1trst
		MaxContractDuration:             time.Hour * 24 * 366 * 10, // a little over 10 years
		MinContractDuration:             time.Second * 40,
		MinContractInterval:             time.Second * 20,
		MinContractDurationForIncentive: time.Second * 20, // time.Hour * 24 // 1 day
		MaxContractIncentive:            5_000_000,        // 5trst
		ContractIncentiveMul:            100,              //  100/100 = 1 = full incentive
		MinContractBalanceForIncentive:  50_000_000,       // 50trst
	})
	return ctx, keeper, codeID, codeHash, walletA, privKeyA, walletB, privKeyB
}

// getDecryptedWasmEvents gets all "wasm" events and decrypt what's necessary
// Returns all "wasm" events, including from contract callbacks
func getDecryptedWasmEvents(t *testing.T, ctx sdk.Context, nonce []byte) []ContractEvent {
	events := ctx.EventManager().Events()
	var res []ContractEvent
	for _, e := range events {
		if e.Type == "wasm" {
			newEvent := []cosmwasm.Attribute{}
			for _, oldLog := range e.Attributes {
				newLog := cosmwasm.Attribute{
					Key:   string(oldLog.Key),
					Value: oldLog.Value,
				}

				if newLog.Key != "contract_address" {
					// key
					keyCipherBz, err := base64.StdEncoding.DecodeString(newLog.Key)
					require.NoError(t, err)
					keyPlainBz, err := wasmCtx.Decrypt(keyCipherBz, nonce)
					require.NoError(t, err)
					newLog.Key = string(keyPlainBz)

				}

				newEvent = append(newEvent, newLog)
			}
			res = append(res, newEvent)
		}
	}
	return res
}

func decryptAttribute(attr cosmwasm.Attribute, nonce []byte) (cosmwasm.Attribute, error) {
	if len(nonce) == 0 {
		return attr, nil
	}
	var newAttr cosmwasm.Attribute

	keyCipherBz, err := base64.StdEncoding.DecodeString(attr.Key)
	if err != nil {
		fmt.Printf("keyCipherBz err %s", err.Error())
		return attr, fmt.Errorf("Failed DecodeString for key %+v\n", attr.Key)
	}
	keyPlainBz, err := wasmCtx.Decrypt(keyCipherBz, nonce)
	if err != nil {
		return attr, fmt.Errorf("Failed Decrypt for key %+v\n", keyCipherBz)
	}

	newAttr.Key = string(keyPlainBz)
	//fmt.Printf("newAttr.Key %s \n", newAttr.Key)

	valuePlainBz, err := wasmCtx.Decrypt(attr.Value, nonce)
	if err != nil {
		return attr, fmt.Errorf("Failed Decrypt for value %+v\n", attr.Value)
	}
	newAttr.Value = valuePlainBz

	return newAttr, nil
}

func parseAndDecryptAttributes(attrs []abci.EventAttribute, nonce []byte) ([]cosmwasm.Attribute, error) {
	var newAttrs []cosmwasm.Attribute
	var err error
	for _, a := range attrs {

		var attr cosmwasm.Attribute
		attr.Key = string(a.Key)
		attr.Value = a.Value
		if attr.Key == "contract_address" {
			newAttrs = append(newAttrs, attr)
			continue
		}

		newAttr, error := decryptAttribute(attr, nonce)

		fmt.Printf("attr.Key err %s \n", err)
		fmt.Printf("attr.Key %s \n", newAttr.Key)

		newAttrs = append(newAttrs, newAttr)
		err = error

	}

	return newAttrs, err
}

// tryDecryptWasmEvents gets all "wasm" events and try to decrypt what it can.
// Returns all "wasm" events, including from contract callbacks.
// The difference between this and getDecryptedWasmEvents is that it is aware of plaintext logs.
func tryDecryptWasmEvents(ctx sdk.Context, nonce []byte, shouldSkipAttributes ...bool) []ContractEvent {
	events := ctx.EventManager().Events()
	var res []ContractEvent
	for _, e := range events {
		if strings.HasPrefix(e.Type, "wasm") {
			newEvent := []cosmwasm.Attribute{}
			for _, oldLog := range e.Attributes {
				newLog := cosmwasm.Attribute{
					Key:   string(oldLog.Key),
					Value: oldLog.Value,
				}
				newEvent = append(newEvent, newLog)

				if newLog.Key != "contract_address" {
					// key
					newAttr, err := decryptAttribute(newLog, nonce)
					fmt.Printf("newAttr %+v \n", newAttr)
					fmt.Printf("err %+v \n", err)
					if err != nil {
						continue
					}

					newEvent[len(newEvent)-1] = newAttr
				}
			}
			res = append(res, newEvent)

		}
	}
	fmt.Printf("decrypted attributes: %+v \n", res)
	return res
}

// getDecryptedData decrypts the output of the first function to be called
// Only returns the data, logs and messages from the first function call
func getDecryptedData(t *testing.T, data []byte, nonce []byte) []byte {
	if len(data) == 0 {
		return data
	}

	dataPlaintextBase64, err := wasmCtx.Decrypt(data, nonce)
	require.NoError(t, err)

	dataPlaintext, err := base64.StdEncoding.DecodeString(string(dataPlaintextBase64))
	require.NoError(t, err)

	return dataPlaintext
}

var contractErrorRegex = regexp.MustCompile(`.*encrypted: (.+): (?:instantiate|execute|query|reply to) contract failed`)

func extractInnerError(t *testing.T, err error, nonce []byte, isEncrypted bool) cosmwasm.StdError {
	match := contractErrorRegex.FindAllStringSubmatch(err.Error(), -1)
	if match == nil {
		require.True(t, !isEncrypted, fmt.Sprintf("Error message should be plaintext but was: %v", err))
		return cosmwasm.StdError{GenericErr: &cosmwasm.GenericErr{Msg: err.Error()}}
	}

	require.True(t, isEncrypted, "Error message should be encrypted")
	require.NotEmpty(t, match)
	require.Equal(t, 1, len(match))
	require.Equal(t, 2, len(match[0]))
	errorCipherB64 := match[0][1]

	errorCipherBz, err := base64.StdEncoding.DecodeString(errorCipherB64)
	require.NoError(t, err)
	errorPlainBz, err := wasmCtx.Decrypt(errorCipherBz, nonce)
	require.NoError(t, err)

	//if !isV1Contract {
	//err = json.Unmarshal(errorPlainBz, &innerErr)
	//require.NoError(t, err)
	//} else {
	innerErr := cosmwasm.StdError{GenericErr: &cosmwasm.GenericErr{Msg: string(errorPlainBz)}}
	//}

	return innerErr
}

const defaultGasForTests uint64 = 120_000

// wrap the default gas meter with a counter of wasm calls
// in order to verify that every wasm call consumes gas
type WasmCounterGasMeter struct {
	wasmCounter uint64
	gasMeter    sdk.GasMeter
}

func (wasmGasMeter *WasmCounterGasMeter) RefundGas(amount stypes.Gas, descriptor string) {}

func (wasmGasMeter *WasmCounterGasMeter) GasConsumed() sdk.Gas {
	return wasmGasMeter.gasMeter.GasConsumed()
}

func (wasmGasMeter *WasmCounterGasMeter) GasConsumedToLimit() sdk.Gas {
	return wasmGasMeter.gasMeter.GasConsumedToLimit()
}

func (wasmGasMeter *WasmCounterGasMeter) Limit() sdk.Gas {
	return wasmGasMeter.gasMeter.Limit()
}

func (wasmGasMeter *WasmCounterGasMeter) ConsumeGas(amount sdk.Gas, descriptor string) {
	//fmt.Printf("condsume gas desc %s", descriptor)
	if (descriptor == "wasm contract" || descriptor == "contract sub-query") && amount > 0 {
		wasmGasMeter.wasmCounter++
	}
	wasmGasMeter.gasMeter.ConsumeGas(amount, descriptor)
}

func (wasmGasMeter *WasmCounterGasMeter) IsPastLimit() bool {
	return wasmGasMeter.gasMeter.IsPastLimit()
}

func (wasmGasMeter *WasmCounterGasMeter) IsOutOfGas() bool {
	return wasmGasMeter.gasMeter.IsOutOfGas()
}

func (wasmGasMeter *WasmCounterGasMeter) String() string {
	return fmt.Sprintf("WasmCounterGasMeter: %+v %+v\n", wasmGasMeter.wasmCounter, wasmGasMeter.gasMeter)
}

func (wasmGasMeter *WasmCounterGasMeter) GetWasmCounter() uint64 {
	return wasmGasMeter.wasmCounter
}

var _ sdk.GasMeter = (*WasmCounterGasMeter)(nil) // check interface

func queryHelper(
	t *testing.T, keeper Keeper, ctx sdk.Context,
	contractAddr sdk.AccAddress, input string,
	isErrorEncrypted bool, gas uint64,
) (string, cosmwasm.StdError) {
	return queryHelperImpl(t, keeper, ctx, contractAddr, input, isErrorEncrypted, gas, -1)
}

func queryHelperImpl(
	t *testing.T, keeper Keeper, ctx sdk.Context,
	contractAddr sdk.AccAddress, input string,
	isErrorEncrypted bool, gas uint64, wasmCallCount int64,
) (string, cosmwasm.StdError) {

	codeHash, _ := keeper.GetContractHash(ctx, contractAddr)

	hashStr := hex.EncodeToString(codeHash)

	msg := types.ContractMsg{
		CodeHash: []byte(hashStr),
		Msg:      []byte(input),
	}

	queryBz, err := wasmCtx.Encrypt(msg.Serialize())
	require.NoError(t, err)
	nonce := queryBz[0:32]

	// create new ctx with the same storage and set our gas meter
	// this is to reset the event manager, so we won't get
	// events from past calls
	gasMeter := &WasmCounterGasMeter{0, sdk.NewGasMeter(gas)}
	ctx = sdk.NewContext(
		ctx.MultiStore(),
		ctx.BlockHeader(),
		ctx.IsCheckTx(),
		log.NewNopLogger(),
	).WithGasMeter(gasMeter)

	resultCipherBz, err := keeper.QueryPrivate(ctx, contractAddr, queryBz, false)

	if wasmCallCount < 0 {
		// default, just check that at least 1 call happened
		require.NotZero(t, gasMeter.GetWasmCounter(), err)
	} else {
		require.Equal(t, uint64(wasmCallCount), gasMeter.GetWasmCounter(), err)
	}

	if err != nil {
		return "", extractInnerError(t, err, nonce, isErrorEncrypted)
	}

	resultPlainBz, err := wasmCtx.Decrypt(resultCipherBz, nonce)
	require.NoError(t, err)

	resultBz, err := base64.StdEncoding.DecodeString(string(resultPlainBz))
	require.NoError(t, err)

	return string(resultBz), cosmwasm.StdError{}
}

func execHelper(
	t *testing.T, keeper Keeper, ctx sdk.Context,
	contractAddress sdk.AccAddress, txSender sdk.AccAddress, senderPrivKey crypto.PrivKey, execMsg string,
	isErrorEncrypted bool, gas uint64, coin int64, shouldSkipAttributes ...bool,
) ([]byte, sdk.Context, []byte, []ContractEvent, uint64, cosmwasm.StdError) {
	return execHelperImpl(t, keeper, ctx, contractAddress, txSender, senderPrivKey, execMsg, isErrorEncrypted, gas, coin, -1, shouldSkipAttributes...)
}

func execHelperImpl(
	t *testing.T, keeper Keeper, ctx sdk.Context,
	contractAddress sdk.AccAddress, txSender sdk.AccAddress, senderPrivKey crypto.PrivKey, execMsg string,
	isErrorEncrypted bool, gas uint64, coin int64, wasmCallCount int64, shouldSkipAttributes ...bool,
) ([]byte, sdk.Context, []byte, []ContractEvent, uint64, cosmwasm.StdError) {
	hash, _ := keeper.GetContractHash(ctx, contractAddress)
	hashStr := hex.EncodeToString(hash)

	msg := types.ContractMsg{
		CodeHash: []byte(hashStr),
		Msg:      []byte(execMsg),
	}

	execMsgBz, err := wasmCtx.Encrypt(msg.Serialize())
	require.NoError(t, err)
	nonce := execMsgBz[0:32]

	// create new ctx with the same storage and a gas limit
	// this is to reset the event manager, so we won't get
	// events from past calls
	gasMeter := &WasmCounterGasMeter{0, sdk.NewGasMeter(gas)}
	ctx = sdk.NewContext(
		ctx.MultiStore(),
		ctx.BlockHeader(),
		ctx.IsCheckTx(),
		log.NewNopLogger(),
	).WithGasMeter(gasMeter)

	ctx = PrepareExecSignedTx(t, keeper, ctx, txSender, senderPrivKey, execMsgBz, contractAddress, sdk.NewCoins(sdk.NewInt64Coin("denom", coin)))

	gasBefore := ctx.GasMeter().GasConsumed()
	execResult, err := keeper.Execute(ctx, contractAddress, txSender, execMsgBz, sdk.NewCoins(sdk.NewInt64Coin("denom", coin)), nil)
	//fmt.Printf("gas %+v\n", gasBefore)
	fmt.Printf("wasmCallCount %+v \n", wasmCallCount)
	gasAfter := ctx.GasMeter().GasConsumed()
	//fmt.Printf("gas %+v\n", gasAfter)
	gasUsed := gasAfter - gasBefore

	if wasmCallCount < 0 {
		// default, just check that at least 1 call happened
		//require.NotZero(t, gasMeter.GetWasmCounter(), err)
	} else {
		require.Equal(t, uint64(wasmCallCount), gasMeter.GetWasmCounter(), err)
	}

	if err != nil {
		return nil, ctx, nil, nil, 0, extractInnerError(t, err, nonce, isErrorEncrypted)
	}

	// wasmEvents comes from all the callbacks as well
	wasmEvents := tryDecryptWasmEvents(ctx, nonce, shouldSkipAttributes...)

	// TODO check if we can extract the messages from ctx

	// Data is the output of only the first call
	data := getDecryptedData(t, execResult.Data, nonce)

	return nonce, ctx, data, wasmEvents, gasUsed, cosmwasm.StdError{}
}

func initHelper(
	t *testing.T, keeper Keeper, ctx sdk.Context,
	codeID uint64, creator sdk.AccAddress, creatorPrivKey crypto.PrivKey, initMsg string,
	isErrorEncrypted bool, gas uint64, shouldSkipAttributes ...bool,
) ([]byte, sdk.Context, sdk.AccAddress, []ContractEvent, cosmwasm.StdError) {
	return initHelperImpl(t, keeper, ctx, codeID, creator, creatorPrivKey, initMsg, isErrorEncrypted, gas, -1, sdk.NewCoins(), shouldSkipAttributes...)
}

func initHelperImpl(
	t *testing.T, keeper Keeper, ctx sdk.Context,
	codeID uint64, creator sdk.AccAddress, creatorPrivKey crypto.PrivKey, initMsg string,
	isErrorEncrypted bool, gas uint64, wasmCallCount int64, sentFunds sdk.Coins, shouldSkipAttributes ...bool,
) ([]byte, sdk.Context, sdk.AccAddress, []ContractEvent, cosmwasm.StdError) {

	info, _ := keeper.GetCodeInfo(ctx, codeID)
	hashStr := hex.EncodeToString(info.CodeHash)

	msg := types.ContractMsg{
		CodeHash: []byte(hashStr),
		Msg:      []byte(initMsg),
	}

	initMsgBz, err := wasmCtx.Encrypt(msg.Serialize())
	require.NoError(t, err)
	nonce := initMsgBz[0:32]

	// create new ctx with the same storage and a gas limit
	// this is to reset the event manager, so we won't get
	// events from past calls
	gasMeter := &WasmCounterGasMeter{0, sdk.NewGasMeter(gas)}
	ctx = sdk.NewContext(
		ctx.MultiStore(),
		ctx.BlockHeader(),
		ctx.IsCheckTx(),
		log.NewNopLogger(),
	).WithGasMeter(gasMeter)

	ctx = PrepareInitSignedTx(t, keeper, ctx, creator, creatorPrivKey, initMsgBz, codeID, sentFunds)
	// make the label a random base64 string, because why not?
	contractAddress, _, err := keeper.Instantiate(ctx, codeID, creator /* nil,*/, initMsgBz, nil, base64.RawURLEncoding.EncodeToString(nonce), sentFunds, nil, 0, 0, time.Now(), nil)

	if wasmCallCount < 0 {
		// default, just check that at least 1 call happened
		require.NotZero(t, gasMeter.GetWasmCounter(), err)
	} else {
		require.Equal(t, uint64(wasmCallCount), gasMeter.GetWasmCounter(), err)
	}

	if err != nil {
		return nil, ctx, nil, nil, extractInnerError(t, err, nonce, isErrorEncrypted)
	}

	// wasmEvents comes from all the callbacks as well
	wasmEvents := tryDecryptWasmEvents(ctx, nonce, shouldSkipAttributes...)

	// TODO check if we can extract the messages from ctx

	return nonce, ctx, contractAddress, wasmEvents, cosmwasm.StdError{}
}

func TestCallbackSanity(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	// init
	_, _, contractAddress, initEvents, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	require.Equal(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "init", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		initEvents,
	)

	_, _, _, execEvents, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, fmt.Sprintf(`{"a":{"contract_addr":"%s","code_hash":"%s","x":2,"y":3}}`, contractAddress.String(), codeHash), true, defaultGasForTests, 0)
	fmt.Printf("TestCallbackSanity ev  %+v \n", execEvents)
	require.Empty(t, err)
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "banana", Value: []byte("🍌"), AccAddr: "", Encrypted: false, PubDb: false},
	}, execEvents[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "kiwi", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
	}, execEvents[1])
	require.ElementsMatch(t, ContractEvent{
		{Key: "watermelon", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
	}, execEvents[2])

}

/*
	func TestSanity(t *testing.T) {
		ctx, keeper, codeID, _, walletA, privKeyA, walletB, _ := setupTest(t, "./testdata/erc20.wasm", sdk.NewCoins())

		// init
		initMsg := fmt.Sprintf(`{"decimals":10,"initial_balances":[{"address":"%s","amount":"108"},{"address":"%s","amount":"53"}],"name":"ReuvenPersonalRustCoin","symbol":"RPRC"}`, walletA.String(), walletB.String())

		_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, initMsg, true, false, defaultGasForTests)
		require.Empty(t, err)
		// require.Empty(t, initEvents)

		// check state after init
		qRes, qErr := queryHelper(t, keeper, ctx, contractAddress, fmt.Sprintf(`{"balance":{"address":"%s"}}`, walletA.String()), true, false, defaultGasForTests)
		require.Empty(t, qErr)
		require.JSONEq(t, `{"balance":"108"}`, qRes)

		qRes, qErr = queryHelper(t, keeper, ctx, contractAddress, fmt.Sprintf(`{"balance":{"address":"%s"}}`, walletB.String()), true, false, defaultGasForTests)
		require.Empty(t, qErr)
		require.JSONEq(t, `{"balance":"53"}`, qRes)

		// transfer 10 from A to B
		_, _, data, wasmEvents, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA,
			fmt.Sprintf(`{"transfer":{"amount":"10","recipient":"%s"}}`, walletB.String()), true, false, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Empty(t, data)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "action", Value: []byte("transfer"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "sender", Value: walletA.Bytes(), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "recipient", Value: walletB.Bytes(), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			wasmEvents,
		)

		// check state after transfer
		qRes, qErr = queryHelper(t, keeper, ctx, contractAddress, fmt.Sprintf(`{"balance":{"address":"%s"}}`, walletA.String()), true, false, defaultGasForTests)
		require.Empty(t, qErr)
		require.JSONEq(t, `{"balance":"98"}`, qRes)

		qRes, qErr = queryHelper(t, keeper, ctx, contractAddress, fmt.Sprintf(`{"balance":{"address":"%s"}}`, walletB.String()), true, false, defaultGasForTests)
		require.Empty(t, qErr)
		require.JSONEq(t, `{"balance":"63"}`, qRes)
	}
*/
func TestInitLogs(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, initEvents, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)
	require.Equal(t, 1, len(initEvents))
	require.Equal(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "init", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		initEvents,
	)

}

func TestEmptyLogKeyValue(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, _, _, execEvents, _, execErr := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"empty_log_key_value":{}}`, true, defaultGasForTests, 0)

	require.Empty(t, execErr)
	require.Equal(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "my value is empty", Value: nil, AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "", Value: []byte("my key is empty"), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		execEvents,
	)

}

/*
func TestEmptyData(t *testing.T) {

			ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

			_, _, contractAddress, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
			require.Empty(t, initErr)

			_, _, data, _, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"empty_data":{}}`, true, defaultGasForTests, 0)

			require.Empty(t, err)
			require.Empty(t, data)

}*/

func TestNoData(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, _, data, _, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"no_data":{}}`, true, defaultGasForTests, 0)

	require.Empty(t, err)
	require.Empty(t, data)

}

func TestExecuteIllegalInputError(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, _, _, _, _, execErr := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `bad input`, true, defaultGasForTests, 0)

	require.NotNil(t, execErr.GenericErr)
	require.Contains(t, execErr.GenericErr.Msg, "Error parsing")

}

func TestInitIllegalInputError(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, _, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `bad input`, true, defaultGasForTests)

	require.NotNil(t, initErr.GenericErr)
	require.Contains(t, initErr.GenericErr.Msg, "Error parsing")

}

func TestCallbackFromInitAndCallbackEvents(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	// init first contract so we'd have someone to callback
	_, _, firstContractAddress, initEvents, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	require.Equal(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: []byte(firstContractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "init", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		initEvents,
	)

	// init second contract and callback to the first contract
	_, _, contractAddress, initEvents, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"callback":{"contract_addr":"%s", "code_hash": "%s"}}`, firstContractAddress.String(), codeHash), true, defaultGasForTests)
	require.Empty(t, initErr)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "init with a callback", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, initEvents[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(firstContractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "watermelon", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
	}, initEvents[1])

	/*
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "init with a callback", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				{
					{Key: "contract_address", Value: []byte(firstContractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "watermelon", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			initEvents,
		)*/

}

/*
	func TestQueryInputParamError(t *testing.T) {
		ctx, keeper, codeID, _, walletA, privKeyA, walletB, _ := setupTest(t, "./testdata/erc20.wasm", sdk.NewCoins())

		// init
		initMsg := fmt.Sprintf(`{"decimals":10,"initial_balances":[{"address":"%s","amount":"108"},{"address":"%s","amount":"53"}],"name":"ReuvenPersonalRustCoin","symbol":"RPRC"}`, walletA.String(), walletB.String())

		_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, initMsg, true, false, defaultGasForTests)
		require.Empty(t, err)
		// require.Empty(t, initEvents)

		_, qErr := queryHelper(t, keeper, ctx, contractAddress, `{"balance":{"address":"blabla"}}`, true, false, defaultGasForTests)

		require.NotNil(t, qErr.GenericErr)
		require.Equal(t, "canonicalize_address errored: invalid length", qErr.GenericErr.Msg)
	}
*/
func TestUnicodeData(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, _, data, _, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"unicode_data":{}}`, true, defaultGasForTests, 0)

	require.Empty(t, err)
	require.Equal(t, "🍆🥑🍄", string(data))

}

func TestInitContractError(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	t.Run("generic_err", func(t *testing.T) {
		_, _, _, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"contract_error":{"error_type":"generic_err"}}`, true, defaultGasForTests)

		require.NotNil(t, err.GenericErr)
		require.Contains(t, err.GenericErr.Msg, "la la 🤯")
	})
	t.Run("invalid_base64", func(t *testing.T) {
		_, _, _, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"contract_error":{"error_type":"invalid_base64"}}`, true, defaultGasForTests)

		require.NotNil(t, err.GenericErr)
		require.Contains(t, err.GenericErr.Msg, "ra ra 🤯")

	})
	t.Run("invalid_utf8", func(t *testing.T) {
		_, _, _, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"contract_error":{"error_type":"invalid_utf8"}}`, true, defaultGasForTests)

		require.NotNil(t, err.GenericErr)
		require.Contains(t, err.GenericErr.Msg, "ka ka 🤯")

	})
	t.Run("not_found", func(t *testing.T) {
		_, _, _, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"contract_error":{"error_type":"not_found"}}`, true, defaultGasForTests)

		require.NotNil(t, err.GenericErr)
		require.Contains(t, err.GenericErr.Msg, "za za 🤯")

	})
	t.Run("parse_err", func(t *testing.T) {
		_, _, _, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"contract_error":{"error_type":"parse_err"}}`, true, defaultGasForTests)

		require.NotNil(t, err.GenericErr)
		require.Contains(t, err.GenericErr.Msg, "na na 🤯")
		require.Contains(t, err.GenericErr.Msg, "pa pa 🤯")

	})
	t.Run("serialize_err", func(t *testing.T) {
		_, _, _, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"contract_error":{"error_type":"serialize_err"}}`, true, defaultGasForTests)

		require.NotNil(t, err.GenericErr)
		require.Contains(t, err.GenericErr.Msg, "ba ba 🤯")
		require.Contains(t, err.GenericErr.Msg, "ga ga 🤯")

	})
	t.Run("unauthorized", func(t *testing.T) {
		_, _, _, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"contract_error":{"error_type":"unauthorized"}}`, true, defaultGasForTests)

		// Not supported in V1
		require.NotNil(t, err.GenericErr)
		require.Contains(t, err.GenericErr.Msg, "catch-all 🤯")

	})
	t.Run("underflow", func(t *testing.T) {
		_, _, _, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"contract_error":{"error_type":"underflow"}}`, true, defaultGasForTests)

		// Not supported in V1
		require.NotNil(t, err.GenericErr)
		require.Contains(t, err.GenericErr.Msg, "catch-all 🤯")

	})

}

func TestExecContractError(t *testing.T) {
	t.Run("TestExecContractError", func(t *testing.T) {
		ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

		_, _, contractAddr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
		require.Empty(t, initErr)

		t.Run("generic_err", func(t *testing.T) {
			_, _, _, _, _, err := execHelper(t, keeper, ctx, contractAddr, walletA, privKeyA, `{"contract_error":{"error_type":"generic_err"}}`, true, defaultGasForTests, 0)

			require.NotNil(t, err.GenericErr)
			require.Contains(t, err.GenericErr.Msg, "la la 🤯")
		})
		t.Run("invalid_base64", func(t *testing.T) {
			_, _, _, _, _, err := execHelper(t, keeper, ctx, contractAddr, walletA, privKeyA, `{"contract_error":{"error_type":"invalid_base64"}}`, true, defaultGasForTests, 0)

			require.NotNil(t, err.GenericErr)
			require.Contains(t, err.GenericErr.Msg, "ra ra 🤯")

		})
		t.Run("invalid_utf8", func(t *testing.T) {
			_, _, _, _, _, err := execHelper(t, keeper, ctx, contractAddr, walletA, privKeyA, `{"contract_error":{"error_type":"invalid_utf8"}}`, true, defaultGasForTests, 0)

			require.NotNil(t, err.GenericErr)
			require.Contains(t, err.GenericErr.Msg, "ka ka 🤯")

		})
		t.Run("not_found", func(t *testing.T) {
			_, _, _, _, _, err := execHelper(t, keeper, ctx, contractAddr, walletA, privKeyA, `{"contract_error":{"error_type":"not_found"}}`, true, defaultGasForTests, 0)

			require.NotNil(t, err.GenericErr)
			require.Contains(t, err.GenericErr.Msg, "za za 🤯")

		})
		t.Run("parse_err", func(t *testing.T) {
			_, _, _, _, _, err := execHelper(t, keeper, ctx, contractAddr, walletA, privKeyA, `{"contract_error":{"error_type":"parse_err"}}`, true, defaultGasForTests, 0)

			require.NotNil(t, err.GenericErr)
			require.Contains(t, err.GenericErr.Msg, "na na 🤯")
			require.Contains(t, err.GenericErr.Msg, "pa pa 🤯")

		})
		t.Run("serialize_err", func(t *testing.T) {
			_, _, _, _, _, err := execHelper(t, keeper, ctx, contractAddr, walletA, privKeyA, `{"contract_error":{"error_type":"serialize_err"}}`, true, defaultGasForTests, 0)

			require.NotNil(t, err.GenericErr)
			require.Contains(t, err.GenericErr.Msg, "ba ba 🤯")
			require.Contains(t, err.GenericErr.Msg, "ga ga 🤯")

		})
		t.Run("unauthorized", func(t *testing.T) {
			_, _, _, _, _, err := execHelper(t, keeper, ctx, contractAddr, walletA, privKeyA, `{"contract_error":{"error_type":"unauthorized"}}`, true, defaultGasForTests, 0)

			require.NotNil(t, err.GenericErr)
			require.Contains(t, err.GenericErr.Msg, "catch-all 🤯")

		})
		t.Run("underflow", func(t *testing.T) {
			_, _, _, _, _, err := execHelper(t, keeper, ctx, contractAddr, walletA, privKeyA, `{"contract_error":{"error_type":"underflow"}}`, true, defaultGasForTests, 0)

			require.NotNil(t, err.GenericErr)
			require.Contains(t, err.GenericErr.Msg, "catch-all 🤯")

		})
	})
}

func TestQueryContractError(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	t.Run("generic_err", func(t *testing.T) {
		_, err := queryHelper(t, keeper, ctx, contractAddr, `{"contract_error":{"error_type":"generic_err"}}`, true, defaultGasForTests)

		require.NotNil(t, err.GenericErr)
		require.Contains(t, err.GenericErr.Msg, "la la 🤯")
	})
	t.Run("invalid_base64", func(t *testing.T) {
		_, err := queryHelper(t, keeper, ctx, contractAddr, `{"contract_error":{"error_type":"invalid_base64"}}`, true, defaultGasForTests)

		require.NotNil(t, err.GenericErr)
		require.Contains(t, err.GenericErr.Msg, "ra ra 🤯")

	})
	t.Run("invalid_utf8", func(t *testing.T) {
		_, err := queryHelper(t, keeper, ctx, contractAddr, `{"contract_error":{"error_type":"invalid_utf8"}}`, true, defaultGasForTests)

		require.NotNil(t, err.GenericErr)
		require.Contains(t, err.GenericErr.Msg, "ka ka 🤯")

	})
	t.Run("not_found", func(t *testing.T) {
		_, err := queryHelper(t, keeper, ctx, contractAddr, `{"contract_error":{"error_type":"not_found"}}`, true, defaultGasForTests)

		require.NotNil(t, err.GenericErr)
		require.Contains(t, err.GenericErr.Msg, "za za 🤯")

	})
	t.Run("parse_err", func(t *testing.T) {
		_, err := queryHelper(t, keeper, ctx, contractAddr, `{"contract_error":{"error_type":"parse_err"}}`, true, defaultGasForTests)

		require.NotNil(t, err.GenericErr)
		require.Contains(t, err.GenericErr.Msg, "na na 🤯")
		require.Contains(t, err.GenericErr.Msg, "pa pa 🤯")

	})
	t.Run("serialize_err", func(t *testing.T) {
		_, err := queryHelper(t, keeper, ctx, contractAddr, `{"contract_error":{"error_type":"serialize_err"}}`, true, defaultGasForTests)

		require.NotNil(t, err.GenericErr)
		require.Contains(t, err.GenericErr.Msg, "ba ba 🤯")
		require.Contains(t, err.GenericErr.Msg, "ga ga 🤯")

	})
	t.Run("unauthorized", func(t *testing.T) {
		_, err := queryHelper(t, keeper, ctx, contractAddr, `{"contract_error":{"error_type":"unauthorized"}}`, true, defaultGasForTests)

		require.NotNil(t, err.GenericErr)
		require.Contains(t, err.GenericErr.Msg, "catch-all 🤯")

	})
	t.Run("underflow", func(t *testing.T) {
		_, err := queryHelper(t, keeper, ctx, contractAddr, `{"contract_error":{"error_type":"underflow"}}`, true, defaultGasForTests)

		require.NotNil(t, err.GenericErr)
		require.Contains(t, err.GenericErr.Msg, "catch-all 🤯")

	})

}

func TestInitParamError(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	codeHash := "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	msg := fmt.Sprintf(`{"callback":{"contract_addr":"notanaddress", "code_hash":"%s"}}`, codeHash)

	_, _, _, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, msg, false, defaultGasForTests)

	require.Contains(t, initErr.Error(), "invalid address")

}

func TestCallbackExecuteParamError(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	msg := fmt.Sprintf(`{"a":{"code_hash":"%s","contract_addr":"notanaddress","x":2,"y":3}}`, codeHash)

	_, _, _, _, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, msg, false, defaultGasForTests, 0)

	require.Contains(t, err.Error(), "invalid address")

}

/*
	func TestQueryInputStructureError(t *testing.T) {
		ctx, keeper, codeID, _, walletA, privKeyA, walletB, _ := setupTest(t, "./testdata/erc20.wasm", sdk.NewCoins())

		// init
		initMsg := fmt.Sprintf(`{"decimals":10,"initial_balances":[{"address":"%s","amount":"108"},{"address":"%s","amount":"53"}],"name":"ReuvenPersonalRustCoin","symbol":"RPRC"}`, walletA.String(), walletB.String())

		_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, initMsg, true, false, defaultGasForTests)
		require.Empty(t, err)
		// require.Empty(t, initEvents)

		_, qErr := queryHelper(t, keeper, ctx, contractAddress, `{"balance":{"invalidkey":"invalidval"}}`, true, false, defaultGasForTests)

		require.NotNil(t, qErr.ParseErr)
		require.Contains(t, qErr.ParseErr.Msg, "missing field `address`")
	}
*/
func TestInitNotEncryptedInputError(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKey, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	//ctx = sdk.NewContext(
	//	ctx.MultiStore(),
	//	ctx.BlockHeader(),
	//	ctx.IsCheckTx(),
	//	log.NewNopLogger(),
	//).WithGasMeter(sdk.NewGasMeter(defaultGas))

	initMsg := []byte(`{"nop":{}`)

	ctx = PrepareInitSignedTx(t, keeper, ctx, walletA, privKey, initMsg, codeID, nil)

	// init
	_, _, err := keeper.Instantiate(ctx, codeID, walletA /* nil, */, initMsg, nil, "some label", sdk.NewCoins(sdk.NewInt64Coin("denom", 0)), nil, 0, 0, time.Now(), nil)
	require.Error(t, err)

	require.Contains(t, err.Error(), "failed to decrypt data")

}

func TestExecuteNotEncryptedInputError(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	//ctx = sdk.NewContext(
	//	ctx.MultiStore(),
	//	ctx.BlockHeader(),
	//	ctx.IsCheckTx(),
	//	log.NewNopLogger(),
	//).WithGasMeter(sdk.NewGasMeter(defaultGas))

	execMsg := []byte(`{"empty_log_key_value":{}}`)

	ctx = PrepareExecSignedTx(t, keeper, ctx, walletA, privKeyA, execMsg, contractAddress, nil)

	_, err := keeper.Execute(ctx, contractAddress, walletA, execMsg, sdk.NewCoins(sdk.NewInt64Coin("denom", 0)), nil)
	require.Error(t, err)

	require.Contains(t, err.Error(), "failed to decrypt data")

}

func TestQueryNotEncryptedInputError(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, err := keeper.QueryPrivate(ctx, contractAddress, []byte(`{"owner":{}}`), false)
	require.Error(t, err)

	require.Contains(t, err.Error(), "failed to decrypt data")

}

func TestInitNoLogs(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	// init
	_, _, _, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"no_logs":{}}`, true, defaultGasForTests)

	require.Empty(t, initErr)
	////require.Empty(t, initEvents)

}

func TestExecNoLogs(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	// init
	_, _, contractAddress, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, _, _, _, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"no_logs":{}}`, true, defaultGasForTests, 0)

	require.Empty(t, err)
	// require.Empty(t, execEvents)

}

func TestExecCallbackToInit(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	// init first contract
	_, _, contractAddress, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	// init second contract and callback to the first contract
	_, _, execData, execEvents, _, execErr := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, fmt.Sprintf(`{"callback_to_init":{"code_id":%d, "code_hash":"%s"}}`, codeID, codeHash), true, defaultGasForTests, 0)
	require.Empty(t, execErr)
	require.Empty(t, execData)

	require.Equal(t, 2, len(execEvents))
	require.Equal(t,
		ContractEvent{
			{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
			{Key: "instantiating a new contract", Value: []byte("🪂"), AccAddr: "", Encrypted: false, PubDb: false},
		},
		execEvents[0],
	)
	require.Contains(t,
		execEvents[1],
		cosmwasm.Attribute{Key: "init", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
	)
	//require.Contains(t, execEvents[1], "contract_address")
	var secondContractAddressBech32 string
	for _, ev := range execEvents[1] {
		if ev.Key == "contract_address" {
			secondContractAddressBech32 = string(ev.Value)
		}
	}
	secondContractAddress, _ := sdk.AccAddressFromBech32(secondContractAddressBech32)
	fmt.Printf("execEvents %+v \n", execEvents)
	fmt.Printf("contractAddress %v \n", secondContractAddress.String())

	_, _, data, _, _, err := execHelper(t, keeper, ctx, secondContractAddress, walletA, privKeyA, `{"unicode_data":{}}`, true, defaultGasForTests, 0)

	require.Empty(t, err)

	require.Equal(t, "🍆🥑🍄", string(data))

}

func TestInitCallbackToInit(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, initEvents, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"callback_to_init":{"code_id":%d, "code_hash":"%s"}}`, codeID, codeHash), true, defaultGasForTests)
	require.Empty(t, initErr)

	require.Equal(t, 2, len(initEvents))
	require.Equal(t,
		ContractEvent{
			{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
			{Key: "instantiating a new contract from init!", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
		},
		initEvents[0],
	)
	require.Contains(t,
		initEvents[1],
		cosmwasm.Attribute{Key: "init", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
	)
	//require.Contains(t, initEvents[1], "contract_address")
	var secondContractAddressBech32 string
	for _, ev := range initEvents[1] {
		if ev.Key == "contract_address" {
			secondContractAddressBech32 = string(ev.Value)
		}
	}
	//secondContractAddressBech32 := string(initEvents[1][0].Value)
	secondContractAddress, _ := sdk.AccAddressFromBech32(secondContractAddressBech32)

	_, _, data, _, _, err := execHelper(t, keeper, ctx, secondContractAddress, walletA, privKeyA, `{"unicode_data":{}}`, true, defaultGasForTests, 0)

	require.Empty(t, err)
	// require.Empty(t, execEvents)
	require.Equal(t, "🍆🥑🍄", string(data))

}

func TestInitCallbackContractError(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, initEvents, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)
	require.Equal(t, 1, len(initEvents))

	_, _, secondContractAddress, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"callback_contract_error":{"contract_addr":"%s", "code_hash":"%s"}}`, contractAddress, codeHash), true, defaultGasForTests)

	require.NotNil(t, initErr.GenericErr)
	require.Contains(t, initErr.GenericErr.Msg, "la la 🤯")
	require.Empty(t, secondContractAddress)
	// require.Empty(t, initEvents)

}

func TestExecCallbackContractError(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	// init
	_, _, contractAddress, initEvents, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)
	require.Equal(t, 1, len(initEvents))

	_, _, data, _, _, execErr := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, fmt.Sprintf(`{"callback_contract_error":{"contract_addr":"%s","code_hash":"%s"}}`, contractAddress, codeHash), true, defaultGasForTests, 0)

	require.NotNil(t, execErr.GenericErr)
	require.Contains(t, execErr.GenericErr.Msg, "la la 🤯")
	// require.Empty(t, execEvents)
	require.Empty(t, data)

}

func TestExecCallbackBadParam(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	// init
	_, _, contractAddress, initEvents, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)
	require.Equal(t, 1, len(initEvents))

	_, _, data, _, _, execErr := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, fmt.Sprintf(`{"callback_contract_bad_param":{"contract_addr":"%s"}}`, contractAddress), true, defaultGasForTests, 0)

	require.NotNil(t, execErr.GenericErr)
	require.Contains(t, execErr.GenericErr.Msg, "v1_sanity_contract::msg::ExecuteMsg")
	require.Contains(t, execErr.GenericErr.Msg, "unknown variant `callback_contract_bad_param`")

	// require.Empty(t, execEvents)
	require.Empty(t, data)

}

/*
func TestInitCallbackBadParam(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	// init first
	_, _, contractAddress, initEvents, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)
	require.Equal(t, 1, len(initEvents))

	_, _, secondContractAddress, initEvents, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"callback_contract_bad_param":{"contract_addr":"%s"}}`, contractAddress), true, defaultGasForTests)
	require.Empty(t, secondContractAddress)
	// require.Empty(t, initEvents)

	require.NotNil(t, initErr.GenericErr)
	require.Contains(t, initErr.GenericErr.Msg, "v1_sanity_contract::msg::InstantiateMsg")
	require.Contains(t, initErr.GenericErr.Msg, "unknown variant `callback_contract_bad_param`")

}
*/
func TestState(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	// init
	_, _, contractAddress, initEvents, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)
	require.Equal(t, 1, len(initEvents))

	_, _, data, _, _, execErr := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"get_state":{"key":"banana"}}`, true, defaultGasForTests, 0)
	require.Empty(t, execErr)
	require.Empty(t, data)

	_, _, _, _, _, execErr = execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"set_state":{"key":"banana","value":"🍌"}}`, true, defaultGasForTests, 0)
	require.Empty(t, execErr)

	_, _, data, _, _, execErr = execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"get_state":{"key":"banana"}}`, true, defaultGasForTests, 0)
	require.Empty(t, execErr)
	require.NotEmpty(t, data)
	//require.Equal(t, "🍌", string(data))

	_, _, _, _, _, execErr = execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"remove_state":{"key":"banana"}}`, true, defaultGasForTests, 0)
	require.Empty(t, execErr)

	_, _, data, _, _, execErr = execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"get_state":{"key":"banana"}}`, true, defaultGasForTests, 0)
	require.Empty(t, execErr)
	require.Empty(t, data)

}

func TestCanonicalizeAddressErrors(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, initEvents, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)
	require.Equal(t, 1, len(initEvents))

	// this function should handle errors internally and return gracefully
	_, _, data, _, _, execErr := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"test_canonicalize_address_errors":{}}`, true, defaultGasForTests, 0)
	require.Empty(t, execErr)
	require.NotEmpty(t, data)
	//require.Equal(t, "🤟", string(data))

}

func TestInitPanic(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, _, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"panic":{}}`, false, defaultGasForTests)

	require.NotNil(t, initErr.GenericErr)
	require.Contains(t, initErr.GenericErr.Msg, "the contract panicked")

}

func TestExecPanic(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, _, _, _, _, execErr := execHelper(t, keeper, ctx, addr, walletA, privKeyA, `{"panic":{}}`, false, defaultGasForTests, 0)

	require.NotNil(t, execErr.GenericErr)
	require.Contains(t, execErr.GenericErr.Msg, "the contract panicked")

}

func TestQueryPanic(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, queryErr := queryHelper(t, keeper, ctx, addr, `{"panic":{}}`, false, defaultGasForTests)
	require.NotNil(t, queryErr.GenericErr)
	require.Contains(t, queryErr.GenericErr.Msg, "the contract panicked")

}

func TestAllocateOnHeapFailBecauseMemoryLimit(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, _, data, _, _, execErr := execHelper(t, keeper, ctx, addr, walletA, privKeyA, `{"allocate_on_heap":{"bytes":13631488}}`, false, defaultGasForTests, 0)

	// this should fail with memory error because 13MiB is more than the allowed 12MiB

	require.Empty(t, data)

	require.NotNil(t, execErr.GenericErr)
	require.Contains(t, execErr.GenericErr.Msg, "the contract panicked")

}

func TestAllocateOnHeapFailBecauseGasLimit(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	// ensure we get an out of gas panic
	defer func() {
		r := recover()
		require.NotNil(t, r)
		_, ok := r.(sdk.ErrorOutOfGas)
		require.True(t, ok, "%+v\n", r)
	}()

	_, _, _, _, _, _ = execHelper(t, keeper, ctx, addr, walletA, privKeyA, `{"allocate_on_heap":{"bytes":1073741824}}`, false, defaultGasForTests, 0)

	// this should fail with out of gas because 1GiB will ask for
	// 134,217,728 gas units (8192 per page times 16,384 pages)
	// the default gas limit in ctx is 200,000 which translates into
	// 20,000,000 WASM gas units, so before the memory_grow opcode is reached
	// the gas metering sees a request that'll cost 134mn and the limit
	// is 20mn, so it throws an out of gas exception

	require.True(t, false)

}

func TestAllocateOnHeapMoreThanSGXHasFailBecauseMemoryLimit(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, _, data, _, _, execErr := execHelper(t, keeper, ctx, addr, walletA, privKeyA, `{"allocate_on_heap":{"bytes":1073741824}}`, false, 9_000_000, 0)

	// this should fail with memory error because 1GiB is more
	// than the allowed 12MiB, gas is 9mn so WASM gas is 900mn
	// which is bigger than the 134mn from the previous test

	require.Empty(t, data)

	require.NotNil(t, execErr.GenericErr)
	require.Contains(t, execErr.GenericErr.Msg, "the contract panicked")

}

func TestPassNullPointerToImports(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	tests := []string{
		"read_db_key",
		"write_db_key",
		"write_db_value",
		"remove_db_key",
		"canonicalize_address_input",
		"humanize_address_input",
	}

	for _, passType := range tests {
		t.Run(passType, func(t *testing.T) {
			_, _, _, _, _, execErr := execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"pass_null_pointer_to_imports_should_throw":{"pass_type":"%s"}}`, passType), false, defaultGasForTests, 0)

			require.NotNil(t, execErr.GenericErr)

			require.Contains(t, execErr.GenericErr.Msg, "execute contract failed")

		})
	}

}

func TestExternalQueryWorks(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, _, data, _, _, execErr := execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"send_external_query":{"to":"%s","code_hash":"%s"}}`, addr.String(), codeHash), true, defaultGasForTests, 0)

	require.Empty(t, execErr)
	require.NotEmpty(t, data)
	//require.Equal(t, []byte{3}, data)

}

func TestExternalQueryCalleePanic(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	_, _, _, _, _, err = execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"send_external_query_panic":{"to":"%s","code_hash":"%s"}}`, addr.String(), codeHash), true, defaultGasForTests, 0)

	require.NotNil(t, err.GenericErr)
	require.Contains(t, err.GenericErr.Msg, "the contract panicked")

}

func TestExternalQueryCalleeStdError(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	_, _, _, _, _, err = execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"send_external_query_error":{"to":"%s","code_hash":"%s"}}`, addr.String(), codeHash), true, defaultGasForTests, 0)

	require.NotNil(t, err.GenericErr)
	require.Contains(t, err.GenericErr.Msg, "la la 🤯")

}

func TestExternalQueryCalleeDoesntExist(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	_, _, _, _, _, err = execHelper(t, keeper, ctx, addr, walletA, privKeyA, `{"send_external_query_error":{"to":"trust13l72vhjngmg55ykajxdnlalktwglyqjqv9pkq4","code_hash":"bla bla"}}`, true, defaultGasForTests, 0)

	require.NotNil(t, err.GenericErr)
	require.Contains(t, err.GenericErr.Msg, "invalid address") //"not found")

}

func TestExternalQueryBadSenderABI(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	_, _, _, _, _, err = execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"send_external_query_bad_abi":{"to":"%s","code_hash":"%s"}}`, addr.String(), codeHash), true, defaultGasForTests, 0)

	require.NotNil(t, err.GenericErr)
	require.Contains(t, err.GenericErr.Msg, "v1_sanity_contract::msg::QueryMsg")
	require.Contains(t, err.GenericErr.Msg, "Invalid type")

}

func TestExternalQueryBadReceiverABI(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	_, _, _, _, _, err = execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"send_external_query_bad_abi_receiver":{"to":"%s","code_hash":"%s"}}`, addr.String(), codeHash), true, defaultGasForTests, 0)

	require.NotNil(t, err.GenericErr)
	require.Contains(t, err.GenericErr.Msg, "alloc::string::String")
	require.Contains(t, err.GenericErr.Msg, "Invalid type")

}

func TestMsgSenderInCallback(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	_, _, _, events, _, err := execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"callback_to_log_msg_sender":{"to":"%s","code_hash":"%s"}}`, addr.String(), codeHash), true, defaultGasForTests, 0)

	require.Empty(t, err)
	/*require.Equal(t, []ContractEvent{
		{
			{Key: "contract_address", Value: []byte(addr.String()), AccAddr: "", Encrypted: false, PubDb: false},
			{Key: "hi", Value: []byte("hey"), AccAddr: "", Encrypted: false, PubDb: false},
		},
		{
			{Key: "contract_address", Value: []byte(addr.String()), AccAddr: "", Encrypted: false, PubDb: false},
			{Key: "msg.sender", Value: []byte(addr.String()), AccAddr: "", Encrypted: false, PubDb: false},
		},
	}, events)*/
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(addr.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "hi", Value: []byte("hey"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(addr.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "msg.sender", Value: []byte(addr.String()), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[1])

}

func TestInfiniteQueryLoopKilledGracefullyByOOM(t *testing.T) {
	t.SkipNow() // We no longer expect to hit OOM trivially

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, 10*defaultGasForTests)
	require.Empty(t, err)

	data, err := queryHelper(t, keeper, ctx, addr, fmt.Sprintf(`{"send_external_query_infinite_loop":{"to":"%s","code_hash":"%s"}}`, addr.String(), codeHash), true, defaultGasForTests)

	require.Empty(t, data)
	require.NotNil(t, err.GenericErr)
	require.Equal(t, err.GenericErr.Msg, "query contract failed: Execution error: Enclave: enclave ran out of heap memory")

}

func TestQueryRecursionLimitEnforcedInQueries(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	data, err := queryHelper(t, keeper, ctx, addr, fmt.Sprintf(`{"send_external_query_recursion_limit":{"to":"%s","code_hash":"%s", "depth":1}}`, addr.String(), codeHash), true, 10*defaultGasForTests)
	fmt.Printf("data %+v\n", data)
	fmt.Printf("err %+v \n", err)
	//require.NotEmpty(t, data)
	require.Equal(t, data, "\"Recursion limit was correctly enforced\"")

	require.Nil(t, err.GenericErr)

}

func TestQueryRecursionLimitEnforcedInHandles(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	_, _, data, _, _, err := execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"send_external_query_recursion_limit":{"to":"%s","code_hash":"%s", "depth":1}}`, addr.String(), codeHash), true, 10*defaultGasForTests, 0)

	require.NotEmpty(t, data)
	//require.Equal(t, string(data), "\"Recursion limit was correctly enforced\"")
	require.Nil(t, err.GenericErr)

}

func TestQueryRecursionLimitEnforcedInInits(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	// Initialize a contract that we will be querying
	_, _, addr, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	// Initialize the contract that will be running the test
	_, _, addr, events, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"send_external_query_recursion_limit":{"to":"%s","code_hash":"%s", "depth":1}}`, addr.String(), codeHash), true, 10*defaultGasForTests)
	require.Empty(t, err)

	require.Nil(t, err.GenericErr)

	require.Equal(t, []ContractEvent{
		{
			{Key: "contract_address", Value: []byte(addr.String()), AccAddr: "", Encrypted: false, PubDb: false},
			{Key: "message", Value: []byte("Recursion limit was correctly enforced"), AccAddr: "", Encrypted: false, PubDb: false},
		},
	}, events)

}

func TestWriteToStorageDuringQuery(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, queryErr := queryHelper(t, keeper, ctx, addr, `{"write_to_storage": {}}`, false, defaultGasForTests)
	require.NotNil(t, queryErr.GenericErr)
	require.Contains(t, queryErr.GenericErr.Msg, "contract tried to write to storage during a query")

}

func TestRemoveFromStorageDuringQuery(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, queryErr := queryHelper(t, keeper, ctx, addr, `{"remove_from_storage": {}}`, false, defaultGasForTests)
	require.NotNil(t, queryErr.GenericErr)
	require.Contains(t, queryErr.GenericErr.Msg, "contract tried to write to storage during a query")

}

func TestDepositToContract(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	contractCoinsBefore := keeper.bankKeeper.GetAllBalances(ctx, addr)
	walletCoinsBefore := keeper.bankKeeper.GetAllBalances(ctx, walletA)

	require.Equal(t, "", contractCoinsBefore.String())
	require.Equal(t, "200000denom", walletCoinsBefore.String())

	_, _, data, _, _, execErr := execHelper(t, keeper, ctx, addr, walletA, privKeyA, `{"deposit_to_contract":{}}`, false, defaultGasForTests, 17)

	require.Empty(t, execErr)

	contractCoinsAfter := keeper.bankKeeper.GetAllBalances(ctx, addr)
	walletCoinsAfter := keeper.bankKeeper.GetAllBalances(ctx, walletA)

	require.Equal(t, "17denom", contractCoinsAfter.String())
	require.Equal(t, "199983denom", walletCoinsAfter.String())

	//require.Equal(t, `[{"denom":"denom","amount":"17"}]`, string(data))
	require.NotEmpty(t, data)

}

func TestContractSendFunds(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, _, _, _, _, execErr := execHelper(t, keeper, ctx, addr, walletA, privKeyA, `{"deposit_to_contract":{}}`, false, defaultGasForTests, 17)

	require.Empty(t, execErr)

	contractCoinsBefore := keeper.bankKeeper.GetAllBalances(ctx, addr)
	walletCoinsBefore := keeper.bankKeeper.GetAllBalances(ctx, walletA)

	require.Equal(t, "17denom", contractCoinsBefore.String())
	require.Equal(t, "199983denom", walletCoinsBefore.String())

	_, _, _, _, _, execErr = execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"send_funds":{"from":"%s","to":"%s","denom":"%s","amount":%d}}`, addr.String(), walletA.String(), "denom", 17), false, defaultGasForTests, 0)

	contractCoinsAfter := keeper.bankKeeper.GetAllBalances(ctx, addr)
	walletCoinsAfter := keeper.bankKeeper.GetAllBalances(ctx, walletA)

	require.Equal(t, "", contractCoinsAfter.String())
	require.Equal(t, "200000denom", walletCoinsAfter.String())

	require.Empty(t, execErr)

}

/*
// In V1 there is no "from" field in Bank message functionality which means it shouldn't be tested

	func TestContractTryToSendFundsFromSomeoneElse(t *testing.T) {
		ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

		_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, false, defaultGasForTests)
		require.Empty(t, initErr)

		_, _, _, _, _, execErr := execHelper(t, keeper, ctx, addr, walletA, privKeyA, `{"deposit_to_contract":{}}`, false, false, defaultGasForTests, 17)

		require.Empty(t, execErr)

		contractCoinsBefore := keeper.bankKeeper.GetAllBalances(ctx, addr)
		walletCoinsBefore := keeper.bankKeeper.GetAllBalances(ctx, walletA)

		require.Equal(t, "17denom", contractCoinsBefore.String())
		require.Equal(t, "199983denom", walletCoinsBefore.String())

		_, _, _, _, _, execErr = execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"send_funds":{"from":"%s","to":"%s","denom":"%s","amount":%d}}`, walletA.String(), addr.Bytes(), "denom", 17), false, false, defaultGasForTests, 0)

		require.NotNil(t, execErr.GenericErr)
		require.Contains(t, execErr.GenericErr.Msg, "contract doesn't have permission")
	}
*/
func TestContractSendFundsToInitCallback(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	contractCoinsBefore := keeper.bankKeeper.GetAllBalances(ctx, addr)
	walletCoinsBefore := keeper.bankKeeper.GetAllBalances(ctx, walletA)

	require.Equal(t, "", contractCoinsBefore.String())
	require.Equal(t, "200000denom", walletCoinsBefore.String())

	_, _, _, execEvents, _, execErr := execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"send_funds_to_init_callback":{"code_id":%d,"denom":"%s","amount":%d,"code_hash":"%s"}}`, codeID, "denom", 17, codeHash), true, defaultGasForTests, 17)

	require.Empty(t, execErr)
	require.NotEmpty(t, execEvents)

	contractCoinsAfter := keeper.bankKeeper.GetAllBalances(ctx, addr)
	walletCoinsAfter := keeper.bankKeeper.GetAllBalances(ctx, walletA)

	var newContract string
	for _, ev := range execEvents[1] {
		if ev.Key == "contract_address" {
			newContract = string(ev.Value)
		}
	}
	newAddr, _ := sdk.AccAddressFromBech32(newContract)
	newContractCoins := keeper.bankKeeper.GetAllBalances(ctx, newAddr)

	require.Equal(t, "", contractCoinsAfter.String())
	require.Equal(t, "199983denom", walletCoinsAfter.String())
	require.Equal(t, "17denom", newContractCoins.String())

}

func TestContractSendFundsToInitCallbackNotEnough(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	contractCoinsBefore := keeper.bankKeeper.GetAllBalances(ctx, addr)
	walletCoinsBefore := keeper.bankKeeper.GetAllBalances(ctx, walletA)

	require.Equal(t, "", contractCoinsBefore.String())
	require.Equal(t, "200000denom", walletCoinsBefore.String())

	_, _, _, _, _, execErr := execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"send_funds_to_init_callback":{"code_id":%d,"denom":"%s","amount":%d,"code_hash":"%s"}}`, codeID, "denom", 18, codeHash), false, defaultGasForTests, 17)

	// require.Empty(t, execEvents)

	require.NotNil(t, execErr.GenericErr)
	require.Contains(t, execErr.GenericErr.Msg, "insufficient funds")

	contractCoinsAfter := keeper.bankKeeper.GetAllBalances(ctx, addr)
	walletCoinsAfter := keeper.bankKeeper.GetAllBalances(ctx, walletA)

	require.Equal(t, "17denom", contractCoinsAfter.String())
	require.Equal(t, "199983denom", walletCoinsAfter.String())

}

func TestContractSendFundsToExecCallback(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, _, addr2, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	contractCoinsBefore := keeper.bankKeeper.GetAllBalances(ctx, addr)
	contract2CoinsBefore := keeper.bankKeeper.GetAllBalances(ctx, addr2)
	walletCoinsBefore := keeper.bankKeeper.GetAllBalances(ctx, walletA)

	require.Equal(t, "", contractCoinsBefore.String())
	require.Equal(t, "", contract2CoinsBefore.String())
	require.Equal(t, "200000denom", walletCoinsBefore.String())

	_, _, _, _, _, execErr := execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"send_funds_to_exec_callback":{"to":"%s","denom":"%s","amount":%d,"code_hash":"%s"}}`, addr2.String(), "denom", 17, codeHash), true, defaultGasForTests, 17)

	require.Empty(t, execErr)

	contractCoinsAfter := keeper.bankKeeper.GetAllBalances(ctx, addr)
	contract2CoinsAfter := keeper.bankKeeper.GetAllBalances(ctx, addr2)
	walletCoinsAfter := keeper.bankKeeper.GetAllBalances(ctx, walletA)

	require.Equal(t, "", contractCoinsAfter.String())
	require.Equal(t, "17denom", contract2CoinsAfter.String())
	require.Equal(t, "199983denom", walletCoinsAfter.String())

}

func TestContractSendFundsToExecCallbackNotEnough(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, _, addr2, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	contractCoinsBefore := keeper.bankKeeper.GetAllBalances(ctx, addr)
	contract2CoinsBefore := keeper.bankKeeper.GetAllBalances(ctx, addr2)
	walletCoinsBefore := keeper.bankKeeper.GetAllBalances(ctx, walletA)

	require.Equal(t, "", contractCoinsBefore.String())
	require.Equal(t, "", contract2CoinsBefore.String())
	require.Equal(t, "200000denom", walletCoinsBefore.String())

	_, _, _, _, _, execErr := execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"send_funds_to_exec_callback":{"to":"%s","denom":"%s","amount":%d,"code_hash":"%s"}}`, addr2.String(), "denom", 19, codeHash), false, defaultGasForTests, 17)

	require.NotNil(t, execErr.GenericErr)
	require.Contains(t, execErr.GenericErr.Msg, "insufficient funds")

	contractCoinsAfter := keeper.bankKeeper.GetAllBalances(ctx, addr)
	contract2CoinsAfter := keeper.bankKeeper.GetAllBalances(ctx, addr2)
	walletCoinsAfter := keeper.bankKeeper.GetAllBalances(ctx, walletA)

	require.Equal(t, "17denom", contractCoinsAfter.String())
	require.Equal(t, "", contract2CoinsAfter.String())
	require.Equal(t, "199983denom", walletCoinsAfter.String())

}

func TestSleep(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, _, _, _, _, execErr := execHelper(t, keeper, ctx, addr, walletA, privKeyA, `{"sleep":{"ms":3000}}`, false, defaultGasForTests, 0)

	require.Error(t, execErr)
	require.Error(t, execErr.GenericErr)
	require.Contains(t, execErr.GenericErr.Msg, "the contract panicked")

}

func TestGasIsChargedForInitCallbackToInit(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, _, _, err := initHelperImpl(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"callback_to_init":{"code_id":%d,"code_hash":"%s"}}`, codeID, codeHash), true, defaultGasForTests, 2, sdk.NewCoins())
	require.Empty(t, err)

}

func TestGasIsChargedForInitCallbackToExec(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, _, _, _, err := initHelperImpl(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"callback":{"contract_addr":"%s","code_hash":"%s"}}`, addr, codeHash), true, defaultGasForTests, 2, sdk.NewCoins())
	require.Empty(t, err)

}

func TestGasIsChargedForExecCallbackToInit(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	// exec callback to init
	_, _, _, _, _, err := execHelperImpl(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"callback_to_init":{"code_id":%d,"code_hash":"%s"}}`, codeID, codeHash), true, defaultGasForTests, 0, 2)
	require.Empty(t, err)

}

func TestGasIsChargedForExecCallbackToExec(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	// exec callback to exec
	_, _, _, _, _, err := execHelperImpl(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"a":{"contract_addr":"%s","code_hash":"%s","x":1,"y":2}}`, addr, codeHash), true, defaultGasForTests, 0, 3)
	require.Empty(t, err)

}

func TestGasIsChargedForExecExternalQuery(t *testing.T) {
	t.SkipNow() // as of v0.10 CosmWasm are overriding the default gas meter

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, _, _, _, _, err := execHelperImpl(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"send_external_query_depth_counter":{"to":"%s","depth":2,"code_hash":"%s"}}`, addr.String(), codeHash), true, defaultGasForTests, 0, 3)
	require.Empty(t, err)

}

func TestGasIsChargedForInitExternalQuery(t *testing.T) {
	t.SkipNow() // as of v0.10 CosmWasm are overriding the default gas meter

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, _, _, _, err := initHelperImpl(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"send_external_query_depth_counter":{"to":"%s","depth":2,"code_hash":"%s"}}`, addr.String(), codeHash), true, defaultGasForTests, 3, sdk.NewCoins())
	require.Empty(t, err)

}

func TestGasIsChargedForQueryExternalQuery(t *testing.T) {
	t.SkipNow() // as of v0.10 CosmWasm are overriding the default gas meter

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	_, err := queryHelperImpl(t, keeper, ctx, addr, fmt.Sprintf(`{"send_external_query_depth_counter":{"to":"%s","depth":2,"code_hash":"%s"}}`, addr.String(), codeHash), true, defaultGasForTests, 3)
	require.Empty(t, err)

}

/*
	func TestWasmTooHighInitialMemoryRuntimeFail(t *testing.T) {
		ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/too-high-initial-memory.wasm", sdk.NewCoins())

		_, _, _, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, false, false, defaultGasForTests)
		require.NotNil(t, err.GenericErr)
		require.Contains(t, err.GenericErr.Msg, "failed to initialize wasm memory")
	}

	func TestWasmTooHighInitialMemoryStaticFail(t *testing.T) {
		encodingConfig := MakeEncodingConfig()
		var transferPortSource types.ICS20TransferPortSource
		transferPortSource = MockIBCTransferKeeper{GetPortFn: func(ctx sdk.Context) string {
			return "myTransferPort"
		}}
		encoders := DefaultEncoders(transferPortSource, encodingConfig.Marshaler)
		ctx, keepers := CreateTestInput(t, false, SupportedFeatures, &encoders, nil)
		accKeeper, keeper := keepers.AccountKeeper, keepers.WasmKeeper

		walletA, _ := CreateFakeFundedAccount(ctx, accKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("denom", 1)))

		wasmCode, err := ioutil.ReadFile("./testdata/test-contract/static-too-high-initial-memory.wasm")
		require.NoError(t, err)

		_, err = keeper.Create(ctx, walletA, wasmCode, "", "", 0, 0, "title", "descr")
		require.Error(t, err)
		require.Contains(t, err.Error(), "Error during static Wasm validation: Wasm contract memory's minimum must not exceed 512 pages")
	}

func TestWasmWithFloatingPoints(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract_with_floats.wasm", sdk.NewCoins())

	_, _, _, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, false, defaultGasForTests)
	require.NotNil(t, err.GenericErr)
	require.Contains(t, err.GenericErr.Msg, "found floating point operation in module code")

}
*/
func TestCodeHashInvalid(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privWalletA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())
	initMsg := []byte(`AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA{"nop":{}`)

	enc, _ := wasmCtx.Encrypt(initMsg)

	ctx = PrepareInitSignedTx(t, keeper, ctx, walletA, privWalletA, enc, codeID, sdk.NewCoins(sdk.NewInt64Coin("denom", 0)))
	_, _, err := keeper.Instantiate(ctx, codeID, walletA /* nil, */, enc, nil, "some label", sdk.NewCoins(sdk.NewInt64Coin("denom", 0)), nil, 0, 0, time.Now(), nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to validate transaction")

}

func TestCodeHashEmpty(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privWalletA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())
	initMsg := []byte(`{"nop":{}`)

	enc, _ := wasmCtx.Encrypt(initMsg)

	ctx = PrepareInitSignedTx(t, keeper, ctx, walletA, privWalletA, enc, codeID, sdk.NewCoins(sdk.NewInt64Coin("denom", 0)))
	_, _, err := keeper.Instantiate(ctx, codeID, walletA /* nil, */, enc, nil, "some label", sdk.NewCoins(sdk.NewInt64Coin("denom", 0)), nil, 0, 0, time.Now(), nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to validate transaction")

}

func TestCodeHashNotHex(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privWalletA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())
	initMsg := []byte(`🍉🍉🍉🍉🍉🍉🍉🍉🍉🍉🍉🍉🍉🍉🍉🍉{"nop":{}}`)

	enc, _ := wasmCtx.Encrypt(initMsg)

	ctx = PrepareInitSignedTx(t, keeper, ctx, walletA, privWalletA, enc, codeID, sdk.NewCoins(sdk.NewInt64Coin("denom", 0)))
	_, _, err := keeper.Instantiate(ctx, codeID, walletA /* nil, */, enc, nil, "some label", sdk.NewCoins(sdk.NewInt64Coin("denom", 0)), nil, 0, 0, time.Now(), nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to validate transaction")

}

func TestCodeHashTooSmall(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privWalletA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	initMsg := []byte(codeHash[0:63] + `{"nop":{}`)

	enc, _ := wasmCtx.Encrypt(initMsg)

	ctx = PrepareInitSignedTx(t, keeper, ctx, walletA, privWalletA, enc, codeID, sdk.NewCoins(sdk.NewInt64Coin("denom", 0)))
	_, _, err := keeper.Instantiate(ctx, codeID, walletA /* nil, */, enc, nil, "some label", sdk.NewCoins(sdk.NewInt64Coin("denom", 0)), nil, 0, 0, time.Now(), nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to validate transaction")

}

func TestCodeHashTooBig(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privWalletA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	initMsg := []byte(codeHash + "a" + `{"nop":{}`)

	enc, _ := wasmCtx.Encrypt(initMsg)

	ctx = PrepareInitSignedTx(t, keeper, ctx, walletA, privWalletA, enc, codeID, sdk.NewCoins(sdk.NewInt64Coin("denom", 0)))
	_, _, err := keeper.Instantiate(ctx, codeID, walletA /* nil, */, enc, nil, "some label", sdk.NewCoins(sdk.NewInt64Coin("denom", 0)), nil, 0, 0, time.Now(), nil)
	require.Error(t, err)

	initErr := extractInnerError(t, err, enc[0:32], true)
	require.NotEmpty(t, initErr)
	require.Contains(t, initErr.Error(), "Expected to parse either a `true`, `false`, or a `null`.")

}

func TestCodeHashWrong(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privWalletA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	initMsg := []byte(`e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855{"nop":{}`)

	enc, _ := wasmCtx.Encrypt(initMsg)

	ctx = PrepareInitSignedTx(t, keeper, ctx, walletA, privWalletA, enc, codeID, sdk.NewCoins(sdk.NewInt64Coin("denom", 0)))
	_, _, err := keeper.Instantiate(ctx, codeID, walletA /* nil, */, enc, nil, "some label", sdk.NewCoins(sdk.NewInt64Coin("denom", 0)), nil, 0, 0, time.Now(), nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to validate transaction")

}

func TestCodeHashInitCallInit(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	t.Run("GoodCodeHash", func(t *testing.T) {
		_, _, addr, events, err := initHelperImpl(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"call_to_init":{"code_id":%d,"code_hash":"%s","msg":"%s","contract_id":"1"}}`, codeID, codeHash, `{\"nop\":{}}`), true, defaultGasForTests, 2, sdk.NewCoins())

		require.Empty(t, err)
		/*require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(addr.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "a", Value: []byte("a"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				{
					{Key: "contract_address", Value: events[1][0].Value},
					{Key: "init", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)*/
		require.ElementsMatch(t, ContractEvent{
			{Key: "contract_address", Value: []byte(addr.String()), AccAddr: "", Encrypted: false, PubDb: false},
			{Key: "a", Value: []byte("a"), AccAddr: "", Encrypted: false, PubDb: false},
		}, events[0])
		require.Contains(t,
			events[1],
			cosmwasm.Attribute{Key: "init", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
		)
	})
	/*
		t.Run("EmptyCodeHash", func(t *testing.T) {
			_, _, _, _, err := initHelperImpl(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"call_to_init":{"code_id":%d,"code_hash":"","msg":"%s","contract_id":"2"}}`, codeID, `{\"nop\":{}}`), false, defaultGasForTests, 2, sdk.NewCoins())

			require.NotEmpty(t, err)
			require.Contains(t,
				err.Error(),
				"failed to validate transaction",
			)
		})
		t.Run("TooBigCodeHash", func(t *testing.T) {
				_, _, _, _, err := initHelperImpl(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"call_to_init":{"code_id":%d,"code_hash":"%sa","msg":"%s","contract_id":"3"}}`, codeID, codeHash, `{\"nop\":{}}`), true, defaultGasForTests, 2, sdk.NewCoins())

				require.NotEmpty(t, err)
				require.Contains(t,
					err.Error(),
					"Expected to parse either a `true`, `false`, or a `null`.",
				)
			})
			t.Run("TooSmallCodeHash", func(t *testing.T) {
				_, _, _, _, err := initHelperImpl(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"call_to_init":{"code_id":%d,"code_hash":"%s","msg":"%s","contract_id":"4"}}`, codeID, codeHash[0:63], `{\"nop\":{}}`), false, defaultGasForTests, 2, sdk.NewCoins())

				require.NotEmpty(t, err)
				require.Contains(t,
					err.Error(),
					"failed to validate transaction",
				)
			})
			t.Run("IncorrectCodeHash", func(t *testing.T) {
				_, _, _, _, err := initHelperImpl(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"call_to_init":{"code_id":%d,"code_hash":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","msg":"%s","contract_id":"5"}}`, codeID, `{\"nop\":{}}`), false, defaultGasForTests, 2, sdk.NewCoins())

				require.NotEmpty(t, err)
				require.Contains(t,
					err.Error(),
					"failed to validate transaction",
				)
			})*/

}

func TestCodeHashInitCallExec(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, err := initHelperImpl(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests, 1, sdk.NewCoins())
	require.Empty(t, err)

	t.Run("GoodCodeHash", func(t *testing.T) {
		_, _, addr2, events, err := initHelperImpl(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"call_to_exec":{"addr":"%s","code_hash":"%s","msg":"%s"}}`, addr.String(), codeHash, `{\"c\":{\"x\":1,\"y\":1}}`), true, defaultGasForTests, 2, sdk.NewCoins())

		require.Empty(t, err)
		require.ElementsMatch(t, ContractEvent{
			{Key: "contract_address", Value: []byte(addr2.String()), AccAddr: "", Encrypted: false, PubDb: false},
			{Key: "b", Value: []byte("b"), AccAddr: "", Encrypted: false, PubDb: false},
		}, events[0])
		require.ElementsMatch(t, ContractEvent{
			{Key: "contract_address", Value: []byte(addr.String()), AccAddr: "", Encrypted: false, PubDb: false},
			{Key: "watermelon", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
		}, events[1])
		/*	require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(addr2.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "b", Value: []byte("b"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				{
					{Key: "contract_address", Value: []byte(addr.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "watermelon", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)*/
	})
	t.Run("EmptyCodeHash", func(t *testing.T) {
		_, _, _, _, err = initHelperImpl(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"call_to_exec":{"addr":"%s","code_hash":"","msg":"%s"}}`, addr.String(), `{\"c\":{\"x\":1,\"y\":1}}`), false, defaultGasForTests, 2, sdk.NewCoins())

		require.NotEmpty(t, err)
		require.Contains(t,
			err.Error(),
			"failed to validate transaction",
		)
	})
	t.Run("TooBigCodeHash", func(t *testing.T) {
		_, _, _, _, err = initHelperImpl(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"call_to_exec":{"addr":"%s","code_hash":"%sa","msg":"%s"}}`, addr.String(), codeHash, `{\"c\":{\"x\":1,\"y\":1}}`), true, defaultGasForTests, 2, sdk.NewCoins())

		require.NotEmpty(t, err)
		require.Contains(t,
			err.Error(),
			"Expected to parse either a `true`, `false`, or a `null`.",
		)
	})
	t.Run("TooSmallCodeHash", func(t *testing.T) {
		_, _, _, _, err = initHelperImpl(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"call_to_exec":{"addr":"%s","code_hash":"%s","msg":"%s"}}`, addr.String(), codeHash[0:63], `{\"c\":{\"x\":1,\"y\":1}}`), false, defaultGasForTests, 2, sdk.NewCoins())

		require.NotEmpty(t, err)
		require.Contains(t,
			err.Error(),
			"failed to validate transaction",
		)
	})
	t.Run("IncorrectCodeHash", func(t *testing.T) {
		_, _, _, _, err = initHelperImpl(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"call_to_exec":{"addr":"%s","code_hash":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","msg":"%s"}}`, addr.String(), `{\"c\":{\"x\":1,\"y\":1}}`), false, defaultGasForTests, 2, sdk.NewCoins())

		require.NotEmpty(t, err)
		require.Contains(t,
			err.Error(),
			"failed to validate transaction",
		)
	})

}

func TestCodeHashInitCallQuery(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	t.Run("GoodCodeHash", func(t *testing.T) {
		_, _, addr2, events, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"call_to_query":{"addr":"%s","code_hash":"%s","msg":"%s"}}`, addr.String(), codeHash, `{\"receive_external_query\":{\"num\":1}}`), true, defaultGasForTests)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(addr2.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "c", Value: []byte("2"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("EmptyCodeHash", func(t *testing.T) {
		_, _, _, _, err = initHelper(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"call_to_query":{"addr":"%s","code_hash":"","msg":"%s"}}`, addr.String(), `{\"receive_external_query\":{\"num\":1}}`), true, defaultGasForTests)

		require.NotEmpty(t, err)
		require.Contains(t,
			err.Error(),
			"failed to validate transaction",
		)
	})
	t.Run("TooBigCodeHash", func(t *testing.T) {
		_, _, _, _, err = initHelper(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"call_to_query":{"addr":"%s","code_hash":"%sa","msg":"%s"}}`, addr.String(), codeHash, `{\"receive_external_query\":{\"num\":1}}`), true, defaultGasForTests)

		require.NotEmpty(t, err)
		require.Contains(t,
			err.Error(),
			"Expected to parse either a `true`, `false`, or a `null`.",
		)
	})
	t.Run("TooSmallCodeHash", func(t *testing.T) {
		_, _, _, _, err = initHelper(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"call_to_query":{"addr":"%s","code_hash":"%s","msg":"%s"}}`, addr.String(), codeHash[0:63], `{\"receive_external_query\":{\"num\":1}}`), true, defaultGasForTests)

		require.NotEmpty(t, err)
		require.Contains(t,
			err.Error(),
			"failed to validate transaction",
		)
	})
	t.Run("IncorrectCodeHash", func(t *testing.T) {
		_, _, _, _, err = initHelper(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"call_to_query":{"addr":"%s","code_hash":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","msg":"%s"}}`, addr.String(), `{\"receive_external_query\":{\"num\":1}}`), true, defaultGasForTests)

		require.NotEmpty(t, err)
		require.Contains(t,
			err.Error(),
			"failed to validate transaction",
		)
	})

}

func TestCodeHashExecCallInit(t *testing.T) {
	t.Run("TestCodeHashExecCallInit", func(t *testing.T) {
		ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

		_, _, addr, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
		require.Empty(t, err)

		t.Run("GoodCodeHash", func(t *testing.T) {
			_, _, _, events, _, err := execHelperImpl(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"call_to_init":{"code_id":%d,"code_hash":"%s","msg":"%s","contract_id":"1"}}`, codeID, codeHash, `{\"nop\":{}}`), true, defaultGasForTests, 0, 2)

			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(addr.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "a", Value: []byte("a"), AccAddr: "", Encrypted: false, PubDb: false},
			}, events[0])
			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: events[1][0].Value},
				{Key: "init", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			}, events[1])
		})
		t.Run("EmptyCodeHash", func(t *testing.T) {
			_, _, _, _, _, err := execHelperImpl(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"call_to_init":{"code_id":%d,"code_hash":"","msg":"%s","contract_id":"2"}}`, codeID, `{\"nop\":{}}`), false, defaultGasForTests, 0, 2)

			require.NotEmpty(t, err)
			require.Contains(t,
				err.Error(),
				"failed to validate transaction",
			)
		})
		t.Run("TooBigCodeHash", func(t *testing.T) {
			_, _, _, _, _, err := execHelperImpl(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"call_to_init":{"code_id":%d,"code_hash":"%sa","msg":"%s","contract_id":"3"}}`, codeID, codeHash, `{\"nop\":{}}`), true, defaultGasForTests, 0, 2)

			require.NotEmpty(t, err)

			require.Contains(t,
				err.Error(),
				"v1_sanity_contract::msg::InstantiateMsg: Expected to parse either a `true`, `false`, or a `null`.",
			)

		})
		t.Run("TooSmallCodeHash", func(t *testing.T) {
			_, _, _, _, _, err := execHelperImpl(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"call_to_init":{"code_id":%d,"code_hash":"%s","msg":"%s","contract_id":"4"}}`, codeID, codeHash[0:63], `{\"nop\":{}}`), false, defaultGasForTests, 0, 2)

			require.NotEmpty(t, err)
			require.Contains(t,
				err.Error(),
				"failed to validate transaction",
			)
		})
		t.Run("IncorrectCodeHash", func(t *testing.T) {
			_, _, _, _, _, err := execHelperImpl(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"call_to_init":{"code_id":%d,"code_hash":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","msg":"%s","contract_id":"5"}}`, codeID, `{\"nop\":{}}`), false, defaultGasForTests, 0, 2)

			require.NotEmpty(t, err)
			require.Contains(t,
				err.Error(),
				"failed to validate transaction",
			)
		})
	})
}

func TestContractIdCollisionWhenMultipleCallbacksToInitFromSameContract(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	_, _, _, _, _, err = execHelperImpl(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"call_to_init":{"code_id":%d,"code_hash":"%s","msg":"%s","contract_id":"1"}}`, codeID, codeHash, `{\"nop\":{}}`), true, defaultGasForTests, 0, 2)
	require.Empty(t, err)

	_, _, _, _, _, err = execHelperImpl(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"call_to_init":{"code_id":%d,"code_hash":"%s","msg":"%s","contract_id":"1"}}`, codeID, codeHash, `{\"nop\":{}}`), false, defaultGasForTests, 0, 1)
	require.NotEmpty(t, err)
	require.NotNil(t, err.GenericErr)
	require.Contains(t, err.GenericErr.Msg, "contract account already exists")

}

func TestCodeHashExecCallExec(t *testing.T) {
	t.Run("TestCodeHashExecCallExec", func(t *testing.T) {
		ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

		_, _, addr, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
		require.Empty(t, err)

		t.Run("GoodCodeHash", func(t *testing.T) {
			_, _, _, events, _, err := execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"call_to_exec":{"addr":"%s","code_hash":"%s","msg":"%s"}}`, addr, codeHash, `{\"c\":{\"x\":1,\"y\":1}}`), true, defaultGasForTests, 0)

			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(addr.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "b", Value: []byte("b"), AccAddr: "", Encrypted: false, PubDb: false},
			}, events[0])
			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: events[1][0].Value},
				{Key: "watermelon", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
			}, events[1])
		})
		t.Run("EmptyCodeHash", func(t *testing.T) {
			_, _, _, _, _, err := execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"call_to_exec":{"addr":"%s","code_hash":"","msg":"%s"}}`, addr, `{\"c\":{\"x\":1,\"y\":1}}`), false, defaultGasForTests, 0)

			require.NotEmpty(t, err)
			require.Contains(t,
				err.Error(),
				"failed to validate transaction",
			)
		})
		t.Run("TooBigCodeHash", func(t *testing.T) {
			_, _, _, _, _, err := execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"call_to_exec":{"addr":"%s","code_hash":"%sa","msg":"%s"}}`, addr, codeHash, `{\"c\":{\"x\":1,\"y\":1}}`), true, defaultGasForTests, 0)

			require.NotEmpty(t, err)

			require.Contains(t,
				err.Error(),
				"v1_sanity_contract::msg::ExecuteMsg: Expected to parse either a `true`, `false`, or a `null`.",
			)

		})
		t.Run("TooSmallCodeHash", func(t *testing.T) {
			_, _, _, _, _, err := execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"call_to_exec":{"addr":"%s","code_hash":"%s","msg":"%s"}}`, addr, codeHash[0:63], `{\"c\":{\"x\":1,\"y\":1}}`), false, defaultGasForTests, 0)

			require.NotEmpty(t, err)
			require.Contains(t,
				err.Error(),
				"failed to validate transaction",
			)
		})
		t.Run("IncorrectCodeHash", func(t *testing.T) {
			_, _, _, _, _, err := execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"call_to_exec":{"addr":"%s","code_hash":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","msg":"%s"}}`, addr, `{\"c\":{\"x\":1,\"y\":1}}`), false, defaultGasForTests, 0)

			require.NotEmpty(t, err)
			require.Contains(t,
				err.Error(),
				"failed to validate transaction",
			)
		})
	})
}

func TestQueryGasPrice(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	t.Run("Query to Self Gas Price", func(t *testing.T) {
		_, _, _, _, gasUsed, err := execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"call_to_query":{"addr":"%s","code_hash":"%s","msg":"%s"}}`, addr.String(), codeHash, `{\"receive_external_query\":{\"num\":1}}`), true, defaultGasForTests, 0)
		require.Empty(t, err)
		// require that more gas was used than the 15K
		require.Greater(t, gasUsed, uint64(15_000))
	})

}

func TestCodeHashExecCallQuery(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	t.Run("GoodCodeHash", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"call_to_query":{"addr":"%s","code_hash":"%s","msg":"%s"}}`, addr.String(), codeHash, `{\"receive_external_query\":{\"num\":1}}`), true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(addr.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "c", Value: []byte("2"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("EmptyCodeHash", func(t *testing.T) {
		_, _, _, _, _, err = execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"call_to_query":{"addr":"%s","code_hash":"","msg":"%s"}}`, addr.String(), `{\"receive_external_query\":{\"num\":1}}`), true, defaultGasForTests, 0)

		require.NotEmpty(t, err)
		require.Contains(t,
			err.Error(),
			"failed to validate transaction",
		)
	})
	t.Run("TooBigCodeHash", func(t *testing.T) {
		_, _, _, _, _, err = execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"call_to_query":{"addr":"%s","code_hash":"%sa","msg":"%s"}}`, addr.String(), codeHash, `{\"receive_external_query\":{\"num\":1}}`), true, defaultGasForTests, 0)

		require.NotEmpty(t, err)

		require.Contains(t,
			err.Error(),
			"Expected to parse either a `true`, `false`, or a `null`",
		)

	})
	t.Run("TooSmallCodeHash", func(t *testing.T) {
		_, _, _, _, _, err = execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"call_to_query":{"addr":"%s","code_hash":"%s","msg":"%s"}}`, addr.String(), codeHash[0:63], `{\"receive_external_query\":{\"num\":1}}`), true, defaultGasForTests, 0)

		require.NotEmpty(t, err)
		require.Contains(t,
			err.Error(),
			"failed to validate transaction",
		)
	})
	t.Run("IncorrectCodeHash", func(t *testing.T) {
		_, _, _, _, _, err = execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"call_to_query":{"addr":"%s","code_hash":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","msg":"%s"}}`, addr.String(), `{\"receive_external_query\":{\"num\":1}}`), true, defaultGasForTests, 0)

		require.NotEmpty(t, err)
		require.Contains(t,
			err.Error(),
			"failed to validate transaction",
		)
	})

}

func TestExecCallQueryPublic(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	fakePubState := []cosmwasm.Attribute{{
		Key:       "hi",
		Value:     []byte("hi back"),
		PubDb:     true,
		AccAddr:   "",
		Encrypted: false,
	}}

	keeper.SetContractPublicState(ctx, addr, fakePubState)
	require.Empty(t, err)

	_, _, _, events, _, err := execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"call_to_query_public":{"addr":"%s","key":"%s"}}`, addr.String(), "hi"), true, defaultGasForTests, 0)
	require.Empty(t, err)
	require.Contains(t,
		events[0],
		cosmwasm.Attribute{Key: "public", Value: []byte("hi back"), AccAddr: "", Encrypted: false, PubDb: false},
	)

}

func TestExecCallQueryPublicAddr(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	fakePubState := []cosmwasm.Attribute{{
		Key:       "hi",
		Value:     []byte("hey"),
		PubDb:     true,
		AccAddr:   addr.String(),
		Encrypted: false,
	}}

	keeper.SetContractPublicState(ctx, addr, fakePubState)
	require.Empty(t, err)

	_, _, _, events, _, err := execHelper(t, keeper, ctx, addr, walletA, privKeyA, fmt.Sprintf(`{"call_to_query_public_addr":{"addr":"%s","key":"%s"}}`, addr.String(), "hi"), true, defaultGasForTests, 0)
	require.Empty(t, err)
	require.Contains(t,
		events[0],
		cosmwasm.Attribute{Key: "public", Value: []byte("hey"), AccAddr: "", Encrypted: false, PubDb: false},
	)

}

func TestCodeHashQueryCallQuery(t *testing.T) {

	ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, addr, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	t.Run("GoodCodeHash", func(t *testing.T) {
		output, err := queryHelper(t, keeper, ctx, addr, fmt.Sprintf(`{"call_to_query":{"addr":"%s","code_hash":"%s","msg":"%s"}}`, addr.String(), codeHash, `{\"receive_external_query\":{\"num\":1}}`), true, defaultGasForTests)

		require.Empty(t, err)
		require.Equal(t, "2", output)
	})
	t.Run("EmptyCodeHash", func(t *testing.T) {
		_, err := queryHelper(t, keeper, ctx, addr, fmt.Sprintf(`{"call_to_query":{"addr":"%s","code_hash":"","msg":"%s"}}`, addr.String(), `{\"receive_external_query\":{\"num\":1}}`), true, defaultGasForTests)

		require.NotEmpty(t, err)
		require.Contains(t,
			err.Error(),
			"failed to validate transaction",
		)
	})
	t.Run("TooBigCodeHash", func(t *testing.T) {
		_, err := queryHelper(t, keeper, ctx, addr, fmt.Sprintf(`{"call_to_query":{"addr":"%s","code_hash":"%sa","msg":"%s"}}`, addr.String(), codeHash, `{\"receive_external_query\":{\"num\":1}}`), true, defaultGasForTests)

		require.NotEmpty(t, err)

		require.Contains(t,
			err.Error(),
			"Expected to parse either a `true`, `false`, or a `null`",
		)

	})
	t.Run("TooSmallCodeHash", func(t *testing.T) {
		_, err := queryHelper(t, keeper, ctx, addr, fmt.Sprintf(`{"call_to_query":{"addr":"%s","code_hash":"%s","msg":"%s"}}`, addr.String(), codeHash[0:63], `{\"receive_external_query\":{\"num\":1}}`), true, defaultGasForTests)

		require.NotEmpty(t, err)
		require.Contains(t,
			err.Error(),
			"failed to validate transaction",
		)
	})
	t.Run("IncorrectCodeHash", func(t *testing.T) {
		_, err := queryHelper(t, keeper, ctx, addr, fmt.Sprintf(`{"call_to_query":{"addr":"%s","code_hash":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","msg":"%s"}}`, addr.String(), `{\"receive_external_query\":{\"num\":1}}`), true, defaultGasForTests)

		require.NotEmpty(t, err)
		require.Contains(t,
			err.Error(),
			"failed to validate transaction",
		)
	})

}

func TestSecp256k1Verify(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	// https://paulmillr.com/noble/

	t.Run("CorrectCompactPubkey", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"secp256k1_verify":{"iterations":1,"pubkey":"A0ZGrlBHMWtCMNAIbIrOxofwCxzZ0dxjT2yzWKwKmo//","sig":"/hZeEYHs9trj+Akeb+7p3UAtXjcDNYP9/D/hj/ALIUAG9bfrJltxkfpMz/9Jn5K3c5QjLuvaNT2jgr7P/AEW8A==","msg_hash":"ARp3VEHssUlDEwoW8AzdQYGKg90ENy8yWePKcjfjzao="}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("true"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("CorrectLongPubkey", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"secp256k1_verify":{"iterations":1,"pubkey":"BEZGrlBHMWtCMNAIbIrOxofwCxzZ0dxjT2yzWKwKmo///ne03QpL+5WFHztzVceB3WD4QY/Ipl0UkHr/R8kDpVk=","sig":"/hZeEYHs9trj+Akeb+7p3UAtXjcDNYP9/D/hj/ALIUAG9bfrJltxkfpMz/9Jn5K3c5QjLuvaNT2jgr7P/AEW8A==","msg_hash":"ARp3VEHssUlDEwoW8AzdQYGKg90ENy8yWePKcjfjzao="}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("true"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("IncorrectMsgHashCompactPubkey", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"secp256k1_verify":{"iterations":1,"pubkey":"A0ZGrlBHMWtCMNAIbIrOxofwCxzZ0dxjT2yzWKwKmo//","sig":"/hZeEYHs9trj+Akeb+7p3UAtXjcDNYP9/D/hj/ALIUAG9bfrJltxkfpMz/9Jn5K3c5QjLuvaNT2jgr7P/AEW8A==","msg_hash":"ARp3VEHssUlDEwoW8AzdQYGKg90ENy8yWePKcjfjzas="}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("false"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("IncorrectMsgHashLongPubkey", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"secp256k1_verify":{"iterations":1,"pubkey":"BEZGrlBHMWtCMNAIbIrOxofwCxzZ0dxjT2yzWKwKmo///ne03QpL+5WFHztzVceB3WD4QY/Ipl0UkHr/R8kDpVk=","sig":"/hZeEYHs9trj+Akeb+7p3UAtXjcDNYP9/D/hj/ALIUAG9bfrJltxkfpMz/9Jn5K3c5QjLuvaNT2jgr7P/AEW8A==","msg_hash":"ARp3VEHssUlDEwoW8AzdQYGKg90ENy8yWePKcjfjzas="}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("false"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("IncorrectSigCompactPubkey", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"secp256k1_verify":{"iterations":1,"pubkey":"A0ZGrlBHMWtCMNAIbIrOxofwCxzZ0dxjT2yzWKwKmo//","sig":"rhZeEYHs9trj+Akeb+7p3UAtXjcDNYP9/D/hj/ALIUAG9bfrJltxkfpMz/9Jn5K3c5QjLuvaNT2jgr7P/AEW8A==","msg_hash":"ARp3VEHssUlDEwoW8AzdQYGKg90ENy8yWePKcjfjzao="}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("false"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("IncorrectSigLongPubkey", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"secp256k1_verify":{"iterations":1,"pubkey":"BEZGrlBHMWtCMNAIbIrOxofwCxzZ0dxjT2yzWKwKmo///ne03QpL+5WFHztzVceB3WD4QY/Ipl0UkHr/R8kDpVk=","sig":"rhZeEYHs9trj+Akeb+7p3UAtXjcDNYP9/D/hj/ALIUAG9bfrJltxkfpMz/9Jn5K3c5QjLuvaNT2jgr7P/AEW8A==","msg_hash":"ARp3VEHssUlDEwoW8AzdQYGKg90ENy8yWePKcjfjzao="}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("false"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("IncorrectCompactPubkey", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"secp256k1_verify":{"iterations":1,"pubkey":"AoSdDHH9J0Bfb9pT8GFn+bW9cEVkgIh4bFsepMWmczXc","sig":"/hZeEYHs9trj+Akeb+7p3UAtXjcDNYP9/D/hj/ALIUAG9bfrJltxkfpMz/9Jn5K3c5QjLuvaNT2jgr7P/AEW8A==","msg_hash":"ARp3VEHssUlDEwoW8AzdQYGKg90ENy8yWePKcjfjzao="}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("false"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("IncorrectLongPubkey", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"secp256k1_verify":{"iterations":1,"pubkey":"BISdDHH9J0Bfb9pT8GFn+bW9cEVkgIh4bFsepMWmczXcFWl11YCgu65hzvNDQE2Qo1hwTMQ/42Xif8O/MrxzvxI=","sig":"/hZeEYHs9trj+Akeb+7p3UAtXjcDNYP9/D/hj/ALIUAG9bfrJltxkfpMz/9Jn5K3c5QjLuvaNT2jgr7P/AEW8A==","msg_hash":"ARp3VEHssUlDEwoW8AzdQYGKg90ENy8yWePKcjfjzao="}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("false"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})

}

func TestBenchmarkSecp256k1VerifyAPI(t *testing.T) {
	t.SkipNow()
	// Assaf: I wrote the benchmark like this because the init functions take testing.T
	// and not testing.B and I just wanted to quickly get a feel for the perf improvements

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, _ := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)

	start := time.Now()
	// https://paulmillr.com/noble/
	execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"secp256k1_verify":{"iterations":10,"pubkey":"A0ZGrlBHMWtCMNAIbIrOxofwCxzZ0dxjT2yzWKwKmo//","sig":"/hZeEYHs9trj+Akeb+7p3UAtXjcDNYP9/D/hj/ALIUAG9bfrJltxkfpMz/9Jn5K3c5QjLuvaNT2jgr7P/AEW8A==","msg_hash":"ARp3VEHssUlDEwoW8AzdQYGKg90ENy8yWePKcjfjzao="}}`, true, defaultGasForTests, 0)
	elapsed := time.Since(start)
	fmt.Printf("TestBenchmarkSecp256k1VerifyAPI took %s\n", elapsed)

}

func TestBenchmarkSecp256k1VerifyCrate(t *testing.T) {
	t.SkipNow()
	// Assaf: I wrote the benchmark like this because the init functions take testing.T
	// and not testing.B and I just wanted to quickly get a feel for the perf improvements

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, _ := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)

	start := time.Now()
	// https://paulmillr.com/noble/
	execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"secp256k1_verify_from_crate":{"iterations":10,"pubkey":"A0ZGrlBHMWtCMNAIbIrOxofwCxzZ0dxjT2yzWKwKmo//","sig":"/hZeEYHs9trj+Akeb+7p3UAtXjcDNYP9/D/hj/ALIUAG9bfrJltxkfpMz/9Jn5K3c5QjLuvaNT2jgr7P/AEW8A==","msg_hash":"ARp3VEHssUlDEwoW8AzdQYGKg90ENy8yWePKcjfjzao="}}`, true, 100_000_000, 0)
	elapsed := time.Since(start)
	fmt.Printf("TestBenchmarkSecp256k1VerifyCrate took %s\n", elapsed)

}

func TestEd25519Verify(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	// https://paulmillr.com/noble/
	t.Run("Correct", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"ed25519_verify":{"iterations":1,"pubkey":"LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","sig":"8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","msg":"YXNzYWYgd2FzIGhlcmU="}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("true"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("IncorrectMsg", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"ed25519_verify":{"iterations":1,"pubkey":"LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","sig":"8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","msg":"YXNzYWYgd2FzIGhlcmUK"}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("false"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("IncorrectSig", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"ed25519_verify":{"iterations":1,"pubkey":"LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","sig":"8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDw==","msg":"YXNzYWYgd2FzIGhlcmU="}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("false"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("IncorrectPubkey", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"ed25519_verify":{"iterations":1,"pubkey":"DV1lgRdKw7nt4hvl8XkGZXMzU9S3uM9NLTK0h0qMbUs=","sig":"8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","msg":"YXNzYWYgd2FzIGhlcmU="}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("false"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})

}

func TestEd25519BatchVerify(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	// https://paulmillr.com/noble/
	t.Run("Correct", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"ed25519_batch_verify":{"iterations":1,"pubkeys":["LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA="],"sigs":["8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg=="],"msgs":["YXNzYWYgd2FzIGhlcmU="]}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("true"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("100Correct", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"ed25519_batch_verify":{"iterations":1,"pubkeys":["LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA="],"sigs":["8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg=="],"msgs":["YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU="]}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("true"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("IncorrectPubkey", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"ed25519_batch_verify":{"iterations":1,"pubkeys":["DV1lgRdKw7nt4hvl8XkGZXMzU9S3uM9NLTK0h0qMbUs="],"sigs":["8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg=="],"msgs":["YXNzYWYgd2FzIGhlcmU="]}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("false"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("IncorrectMsg", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"ed25519_batch_verify":{"iterations":1,"pubkeys":["LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA="],"sigs":["8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg=="],"msgs":["YXNzYWYgd2FzIGhlcmUK"]}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("false"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("IncorrectSig", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"ed25519_batch_verify":{"iterations":1,"pubkeys":["LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA="],"sigs":["8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDw=="],"msgs":["YXNzYWYgd2FzIGhlcmU="]}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("false"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("CorrectEmptySigsEmptyMsgsOnePubkey", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"ed25519_batch_verify":{"iterations":1,"pubkeys":["LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA="],"sigs":[],"msgs":[]}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("true"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("CorrectEmpty", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"ed25519_batch_verify":{"iterations":1,"pubkeys":[],"sigs":[],"msgs":[]}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("true"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("CorrectEmptyPubkeysEmptySigsOneMsg", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"ed25519_batch_verify":{"iterations":1,"pubkeys":[],"sigs":[],"msgs":["YXNzYWYgd2FzIGhlcmUK"]}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("true"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("CorrectMultisig", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"ed25519_batch_verify":{"iterations":1,"pubkeys":["LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","2ukhmWRNmcgCrB9fpLP9/HZVuJn6AhpITf455F4GsbM="],"sigs":["8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","bp/N4Ub2WFk9SE9poZVEanU1l46WMrFkTd5wQIXi6QJKjvZUi7+GTzmTe8y2yzgpBI+GWQmt0/QwYbnSVxq/Cg=="],"msgs":["YXNzYWYgd2FzIGhlcmU="]}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("true"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})
	t.Run("CorrectMultiMsgOneSigner", func(t *testing.T) {
		_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"ed25519_batch_verify":{"iterations":1,"pubkeys":["2ukhmWRNmcgCrB9fpLP9/HZVuJn6AhpITf455F4GsbM="],"sigs":["bp/N4Ub2WFk9SE9poZVEanU1l46WMrFkTd5wQIXi6QJKjvZUi7+GTzmTe8y2yzgpBI+GWQmt0/QwYbnSVxq/Cg==","uuNxLEzAYDbuJ+BiYN94pTqhD7UhvCJNbxAbnWz0B9DivkPXmqIULko0DddP2/tVXPtjJ90J20faiWCEC3QkDg=="],"msgs":["YXNzYWYgd2FzIGhlcmU=","cGVhY2Ugb3V0"]}}`, true, defaultGasForTests, 0)

		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "result", Value: []byte("true"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	})

}

func TestSecp256k1RecoverPubkey(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	// https://paulmillr.com/noble/
	_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"secp256k1_recover_pubkey":{"iterations":1,"recovery_param":0,"sig":"/hZeEYHs9trj+Akeb+7p3UAtXjcDNYP9/D/hj/ALIUAG9bfrJltxkfpMz/9Jn5K3c5QjLuvaNT2jgr7P/AEW8A==","msg_hash":"ARp3VEHssUlDEwoW8AzdQYGKg90ENy8yWePKcjfjzao="}}`, true, defaultGasForTests, 0)

	require.Empty(t, err)
	require.Equal(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "result", Value: []byte("A0ZGrlBHMWtCMNAIbIrOxofwCxzZ0dxjT2yzWKwKmo//"), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		events,
	)

	_, _, _, events, _, err = execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"secp256k1_recover_pubkey":{"iterations":1,"recovery_param":1,"sig":"/hZeEYHs9trj+Akeb+7p3UAtXjcDNYP9/D/hj/ALIUAG9bfrJltxkfpMz/9Jn5K3c5QjLuvaNT2jgr7P/AEW8A==","msg_hash":"ARp3VEHssUlDEwoW8AzdQYGKg90ENy8yWePKcjfjzao="}}`, true, defaultGasForTests, 0)

	require.Empty(t, err)
	require.Equal(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "result", Value: []byte("Ams198xOCEVnc/ESvxF2nxnE3AVFO8ahB22S1ZgX2vSR"), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		events,
	)

}

func TestSecp256k1Sign(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	// priv iadRiuRKNZvAXwolxqzJvr60uiMDJTxOEzEwV8OK2ao=
	// pub ArQojoh5TVlSSNA1HFlH5HcQsv0jnrpeE7hgwR/N46nS
	// msg d2VuIG1vb24=
	// msg_hash K9vGEuzCYCUcIXlhMZu20ke2K4mJhreguYct5MqAzhA=

	// https://paulmillr.com/noble/
	_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"secp256k1_sign":{"iterations":1,"msg":"d2VuIG1vb24=","privkey":"iadRiuRKNZvAXwolxqzJvr60uiMDJTxOEzEwV8OK2ao="}}`, true, defaultGasForTests, 0)
	require.Empty(t, err)

	signature := events[0][1].Value

	_, _, _, events, _, err = execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, fmt.Sprintf(`{"secp256k1_verify":{"iterations":1,"pubkey":"ArQojoh5TVlSSNA1HFlH5HcQsv0jnrpeE7hgwR/N46nS","sig":"%s","msg_hash":"K9vGEuzCYCUcIXlhMZu20ke2K4mJhreguYct5MqAzhA="}}`, signature), true, defaultGasForTests, 0)

	require.Empty(t, err)
	require.Equal(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "result", Value: []byte("true"), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		events,
	)

}

func TestEd25519Sign(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, initErr := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, initErr)

	// priv z01UNefH2yjRslwZMmcHssdHmdEjzVvbxjr+MloUEYo=
	// pub jh58UkC0FDsiupZBLdaqKUqYubJbk3LDaruZiJiy0Po=
	// msg d2VuIG1vb24=
	// msg_hash K9vGEuzCYCUcIXlhMZu20ke2K4mJhreguYct5MqAzhA=

	// https://paulmillr.com/noble/
	_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"ed25519_sign":{"iterations":1,"msg":"d2VuIG1vb24=","privkey":"z01UNefH2yjRslwZMmcHssdHmdEjzVvbxjr+MloUEYo="}}`, true, defaultGasForTests, 0)
	require.Empty(t, err)

	signature := events[0][1].Value

	_, _, _, events, _, err = execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, fmt.Sprintf(`{"ed25519_verify":{"iterations":1,"pubkey":"jh58UkC0FDsiupZBLdaqKUqYubJbk3LDaruZiJiy0Po=","sig":"%s","msg":"d2VuIG1vb24="}}`, signature), true, defaultGasForTests, 0)

	require.Empty(t, err)
	require.Equal(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "result", Value: []byte("true"), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		events,
	)

}

func TestBenchmarkEd25519BatchVerifyAPI(t *testing.T) {
	t.SkipNow()
	// Assaf: I wrote the benchmark like this because the init functions take testing.T
	// and not testing.B and I just wanted to quickly get a feel for the performance improvements

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, _ := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)

	start := time.Now()
	_, _, _, _, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"ed25519_batch_verify":{"iterations":1000,"pubkeys":["LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA=","LO2+Bt+/FIjomSaPB+I++LXkxgxwfnrKHLyvCic72rA="],"sigs":["8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg==","8O7nwhM71/B9srKwe8Ps39z5lAsLMMs6LxdvoPk0HXjEM97TNhKbdU6gEePT2MaaIUSiMEmoG28HIZMgMRTCDg=="],"msgs":["YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU=","YXNzYWYgd2FzIGhlcmU="]}}`, true, math.MaxUint64, 0)

	require.Empty(t, err)

	elapsed := time.Since(start)
	fmt.Printf("TestBenchmarkEd25519BatchVerifyAPI took %s\n", elapsed)

}

type GetResponse struct {
	Count uint32 `json:"count"`
}
type v1QueryResponse struct {
	Get GetResponse `json:"get"`
}

func TestV1EndpointsSanity(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, _ := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"counter":{"counter":10, "expires":100}}`, true, defaultGasForTests)

	_, _, _, _, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"increment":{"addition": 13}}`, true, math.MaxUint64, 0)

	require.Empty(t, err)
	//	require.Equal(t, uint32(23), binary.BigEndian.Uint32(data))

	queryRes, qErr := queryHelper(t, keeper, ctx, contractAddress, `{"get":{}}`, true, math.MaxUint64)
	require.Empty(t, qErr)
	// assert result is 32 byte sha256 hash (if hashed), or contractAddr if not
	var resp v1QueryResponse
	e := json.Unmarshal([]byte(queryRes), &resp)
	require.NoError(t, e)

	require.Equal(t, uint32(23), resp.Get.Count)
}

func TestV1QueryWorksWithEnv(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, _ := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"counter":{"counter":10, "expires":0}}`, true, defaultGasForTests)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 10)

	queryRes, qErr := queryHelper(t, keeper, ctx, contractAddress, `{"get":{}}`, true, math.MaxUint64)
	require.Empty(t, qErr)

	// assert result is 32 byte sha256 hash (if hashed), or contractAddr if not
	var resp v1QueryResponse
	e := json.Unmarshal([]byte(queryRes), &resp)
	require.NoError(t, e)
	require.Equal(t, uint32(0), resp.Get.Count)
}

func TestV1ReplySanity(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"counter":{"counter":10, "expires":100}}`, true, defaultGasForTests)
	require.Empty(t, err)

	_, _, _, ev, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"increment":{"addition": 13}}`, true, math.MaxUint64, 0)

	require.Empty(t, err)
	require.Contains(t, ev[0], cosmwasm.Attribute{Key: "resp", Value: []byte("23"), Encrypted: false, PubDb: false, AccAddr: ""})

	_, _, _, ev, _, err = execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"transfer_money":{"amount": 10213}}`, true, math.MaxUint64, 0)

	require.Empty(t, err)
	require.Contains(t, ev[1], cosmwasm.Attribute{Key: "resp", Value: []byte("23"), Encrypted: false, PubDb: false, AccAddr: ""})

	_, _, _, ev, _, err = execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"recursive_reply":{}}`, true, math.MaxUint64, 0)

	require.Empty(t, err)
	require.Contains(t, ev[1], cosmwasm.Attribute{Key: "resp", Value: []byte("25"), Encrypted: false, PubDb: false, AccAddr: ""})

	_, _, _, ev, _, err = execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"recursive_reply_fail":{}}`, true, math.MaxUint64, 0)

	require.Empty(t, err)
	require.Contains(t, ev[1], cosmwasm.Attribute{Key: "resp", Value: []byte("10"), Encrypted: false, PubDb: false, AccAddr: ""})

	_, _, _, ev, _, err = execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"init_new_contract":{}}`, true, math.MaxUint64, 0)

	require.Empty(t, err)
	require.Contains(t, ev[2], cosmwasm.Attribute{Key: "resp", Value: []byte("150"), Encrypted: false, PubDb: false, AccAddr: ""})

	_, _, _, ev, _, err = execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"init_new_contract_with_error":{}}`, true, math.MaxUint64, 0)

	require.Empty(t, err)
	require.Contains(t, ev[1], cosmwasm.Attribute{Key: "resp", Value: []byte("1337"), Encrypted: false, PubDb: false, AccAddr: ""})

	queryRes, qErr := queryHelper(t, keeper, ctx, contractAddress, `{"get":{}}`, true, math.MaxUint64)
	require.Empty(t, qErr)

	// assert result is 32 byte sha256 hash (if hashed), or contractAddr if not
	var resp v1QueryResponse
	e := json.Unmarshal([]byte(queryRes), &resp)
	require.NoError(t, e)
	require.Equal(t, uint32(1337), resp.Get.Count)
}

func TestV1ReplyOnMultipleSubmessages(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, _ := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"counter":{"counter":10, "expires":100}}`, true, defaultGasForTests)

	_, _, data, ev, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"multiple_sub_messages":{}}`, true, math.MaxUint64, 0)

	require.Empty(t, err)
	fmt.Printf("data %v\n", data)
	require.Equal(t, uint32(102), binary.BigEndian.Uint32(data))

	require.Contains(t, ev[4], cosmwasm.Attribute{Key: "resp", Value: []byte("102"), Encrypted: false, PubDb: false, AccAddr: ""})

}

func TestV1MultipleSubmessagesNoReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, _ := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"counter":{"counter":10, "expires":100}}`, true, defaultGasForTests)

	_, _, _, ev, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"multiple_sub_messages_no_reply":{}}`, true, math.MaxUint64, 0)

	require.Empty(t, err)
	//require.Equal(t, uint32(10), binary.BigEndian.Uint32(data))
	fmt.Printf("ev %+v\n", ev)
	require.Contains(t, ev[0], cosmwasm.Attribute{Key: "resp", Value: []byte("10"), Encrypted: false, PubDb: false, AccAddr: ""})
}

func TestV1ReplyLoop(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"counter":{"counter":10, "expires":100}}`, true, defaultGasForTests)
	require.Empty(t, err)
	_, _, _, ev, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"sub_msg_loop":{"iter": 10}}`, true, math.MaxUint64, 0)
	require.Empty(t, err)
	fmt.Printf("ev %+v \n", ev)
	require.Contains(t, ev[21], cosmwasm.Attribute{Key: "resp", Value: []byte("20"), Encrypted: false, PubDb: false, AccAddr: ""})

}

func TestBankMsgSend(t *testing.T) {

	for _, callType := range []string{"init", "exec"} {
		t.Run(callType, func(t *testing.T) {
			for _, test := range []struct {
				description    string
				input          string
				isSuccuss      bool
				errorMsg       string
				balancesBefore string
				balancesAfter  string
			}{
				{
					description:    "regular",
					input:          `[{"amount":"2","denom":"denom"}]`,
					isSuccuss:      true,
					balancesBefore: "5000assaf,200000denom 5000assaf,5000denom",
					balancesAfter:  "4998assaf,199998denom 5000assaf,5002denom",
				},
				{
					description:    "multi-coin",
					input:          `[{"amount":"1","denom":"assaf"},{"amount":"1","denom":"denom"}]`,
					isSuccuss:      true,
					balancesBefore: "5000assaf,200000denom 5000assaf,5000denom",
					balancesAfter:  "4998assaf,199998denom 5001assaf,5001denom",
				},
				/*	{
							description:    "zero",
							input:          `[{"amount":"0","denom":"denom"}]`,
							isSuccuss:      false,
							errorMsg:       "encrypted: submessages: 0denom: invalid coins",
							balancesBefore: "5000assaf,200000denom 5000assaf,5000denom",
							balancesAfter:  "4998assaf,199998denom 5000assaf,5000denom",
						},
						{
							description:    "insufficient funds",
							input:          `[{"amount":"3","denom":"denom"}]`,
							isSuccuss:      false,
							balancesBefore: "5000assaf,200000denom 5000assaf,5000denom",
							balancesAfter:  "4998assaf,199998denom 5000assaf,5000denom",
							errorMsg:       "encrypted: submessages: 2denom is smaller than 3denom: insufficient funds",
						},
					{
							description:    "non-existing denom",
							input:          `[{"amount":"1","denom":"blabla"}]`,
							isSuccuss:      false,
							balancesBefore: "5000assaf,200000denom 5000assaf,5000denom",
							balancesAfter:  "4998assaf,199998denom 5000assaf,5000denom",
							errorMsg:       "encrypted: submessages: 0blabla is smaller than 1blabla: insufficient funds",
						},*/
				{
					description:    "none",
					input:          `[]`,
					isSuccuss:      true,
					balancesBefore: "5000assaf,200000denom 5000assaf,5000denom",
					balancesAfter:  "4998assaf,199998denom 5000assaf,5000denom",
				},
			} {
				t.Run(test.description, func(t *testing.T) {
					ctx, keeper, codeID, _, walletA, privKeyA, walletB, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins(sdk.NewInt64Coin("assaf", 5000)))

					walletACoinsBefore := keeper.bankKeeper.GetAllBalances(ctx, walletA)
					walletBCoinsBefore := keeper.bankKeeper.GetAllBalances(ctx, walletB)

					require.Equal(t, test.balancesBefore, walletACoinsBefore.String()+" "+walletBCoinsBefore.String())

					var err cosmwasm.StdError
					var contractAddress sdk.AccAddress

					if callType == "init" {
						_, _, _, _, _ = initHelperImpl(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"bank_msg_send":{"to":"%s","amount":%s}}`, walletB.String(), test.input), false, defaultGasForTests, -1, sdk.NewCoins(sdk.NewInt64Coin("denom", 2), sdk.NewInt64Coin("assaf", 2)))
					} else {
						_, _, contractAddress, _, _ = initHelperImpl(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, false, defaultGasForTests, -1, sdk.NewCoins(sdk.NewInt64Coin("denom", 2), sdk.NewInt64Coin("assaf", 2)))

						_, _, _, _, _, err = execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, fmt.Sprintf(`{"bank_msg_send":{"to":"%s","amount":%s}}`, walletB.String(), test.input), false, math.MaxUint64, 0)
					}

					if test.isSuccuss {
						require.Empty(t, err)
					} else {
						require.NotEmpty(t, err)
						require.Equal(t, err.Error(), test.errorMsg)
					}

					walletACoinsAfter := keeper.bankKeeper.GetAllBalances(ctx, walletA)
					walletBCoinsAfter := keeper.bankKeeper.GetAllBalances(ctx, walletB)

					require.Equal(t, test.balancesAfter, walletACoinsAfter.String()+" "+walletBCoinsAfter.String())
				})
			}
		})
	}

}

func TestBankMsgBurn(t *testing.T) {
	t.Run("v1", func(t *testing.T) {
		for _, callType := range []string{"init", "exec"} {
			t.Run(callType, func(t *testing.T) {
				for _, test := range []struct {
					description string
					sentFunds   sdk.Coins
				}{
					{
						description: "try to burn coins it has",
						sentFunds:   sdk.NewCoins(sdk.NewInt64Coin("denom", 1)),
					},
					{
						description: "try to burn coins it doesnt have",
						sentFunds:   sdk.NewCoins(),
					},
				} {
					t.Run(test.description, func(t *testing.T) {
						ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

						var err cosmwasm.StdError
						var contractAddress sdk.AccAddress

						if callType == "init" {
							_, _, _, _, err = initHelperImpl(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"bank_msg_burn":{"amount":[{"amount":"1","denom":"denom"}]}}`), false, defaultGasForTests, -1, test.sentFunds)
						} else {
							_, _, contractAddress, _, _ = initHelperImpl(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, false, defaultGasForTests, -1, test.sentFunds)

							_, _, _, _, _, err = execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, fmt.Sprintf(`{"bank_msg_burn":{"amount":[{"amount":"1","denom":"denom"}]}}`), false, math.MaxUint64, 0)
						}

						require.NotEmpty(t, err)
						require.Contains(t, err.Error(), "Unknown variant of Bank: invalid CosmosMsg from the contract")
					})
				}
			})
		}
	})
}

func TestCosmosMsgCustom(t *testing.T) {

	for _, callType := range []string{"init", "exec"} {
		t.Run(callType, func(t *testing.T) {
			ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

			var err cosmwasm.StdError
			var contractAddress sdk.AccAddress

			if callType == "init" {
				_, _, contractAddress, _, err = initHelperImpl(t, keeper, ctx, codeID, walletA, privKeyA, fmt.Sprintf(`{"cosmos_msg_custom":{}}`), false, defaultGasForTests, -1, sdk.NewCoins())
			} else {
				_, _, contractAddress, _, err = initHelperImpl(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, false, defaultGasForTests, -1, sdk.NewCoins())

				_, _, _, _, _, err = execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, fmt.Sprintf(`{"cosmos_msg_custom":{}}`), false, math.MaxUint64, 0)
			}

			require.NotEmpty(t, err)

			require.Contains(t, err.Error(), "Custom variant not supported: invalid CosmosMsg from the contract")

		})
	}

}

/*
	func TestV1InitV010ContractNoReplyWithError(t *testing.T) {
		ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

		wasmCode, err := ioutil.ReadFile("./testdata/test-contract/contract.wasm")
		require.NoError(t, err)

		v010CodeID, err := keeper.Create(ctx, walletA, wasmCode, "", "", 0, 0, "title", "descr")
		require.NoError(t, err)

		v010CodeHash := hex.EncodeToString(keeper.GetCodeInfo(ctx, v010CodeID).CodeHash)

		_, _, contractAddress, _, _ := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"counter":{"counter":10, "expires":100}}`, true. defaultGasForTests)
		msg := fmt.Sprintf(`{"init_v10_no_reply_with_error":{"code_id":%d, "code_hash":"%s"}}`, v010CodeID, v010CodeHash)

		_, _, data, _, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, msg, true. math.MaxUint64, 0)

		require.NotEmpty(t, err)
		require.Nil(t, data)
	}

	func TestV1ExecuteV010ContractNoReplyWithError(t *testing.T) {
		ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

		wasmCode, err := ioutil.ReadFile("./testdata/test-contract/contract.wasm")
		require.NoError(t, err)

		v010CodeID, err := keeper.Create(ctx, walletA, wasmCode, "", "", 0, 0, "title", "descr")
		require.NoError(t, err)

		v010CodeHash := hex.EncodeToString(keeper.GetCodeInfo(ctx, v010CodeID).CodeHash)

		_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true. defaultGasForTests)
		require.Empty(t, err)
		_, _, v010ContractAddress, _, err := initHelper(t, keeper, ctx, v010CodeID, walletA, privKeyA, `{"nop":{}}`, true, false, defaultGasForTests)
		require.Empty(t, err)

		msg := fmt.Sprintf(`{"exec_v10_no_reply_with_error":{"address":"%s", "code_hash":"%s"}}`, v010ContractAddress, v010CodeHash)

		_, _, data, _, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, msg, true. math.MaxUint64, 0)

		require.NotEmpty(t, err)
		require.Nil(t, data)
	}

	func TestV1QueryV010ContractWithError(t *testing.T) {
		ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

		wasmCode, err := ioutil.ReadFile("./testdata/test-contract/contract.wasm")
		require.NoError(t, err)

		v010CodeID, err := keeper.Create(ctx, walletA, wasmCode, "", "", 0, 0, "title", "descr")
		require.NoError(t, err)

		v010CodeHash := hex.EncodeToString(keeper.GetCodeInfo(ctx, v010CodeID).CodeHash)

		_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true. defaultGasForTests)
		require.Empty(t, err)
		_, _, v010ContractAddress, _, err := initHelper(t, keeper, ctx, v010CodeID, walletA, privKeyA, `{"nop":{}}`, true, false, defaultGasForTests)
		require.Empty(t, err)

		msg := fmt.Sprintf(`{"query_v10_with_error":{"address":"%s", "code_hash":"%s"}}`, v010ContractAddress, v010CodeHash)

		_, _, data, _, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, msg, true. math.MaxUint64, 0)

		require.NotEmpty(t, err)
		require.Nil(t, data)
	}

	func TestV010InitV1ContractFromInitWithOkResponse(t *testing.T) {
		ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

		wasmCode, err := ioutil.ReadFile("./testdata/test-contract/contract.wasm")
		require.NoError(t, err)

		v010CodeID, err := keeper.Create(ctx, walletA, wasmCode, "", "", 0, 0, "title", "descr")
		require.NoError(t, err)

		_, _, contractAddress, initEvents, err := initHelper(t, keeper, ctx, v010CodeID, walletA, privKeyA, fmt.Sprintf(`{"callback_to_init":{"code_id":%d, "code_hash":"%s"}}`, codeID, codeHash), true. defaultGasForTests)
		queryRes, qErr := queryHelper(t, keeper, ctx, contractAddress, `{"get_contract_version":{}}`, true, false, math.MaxUint64)
		require.Empty(t, qErr)

		require.Equal(t, queryRes, "10")

		require.Empty(t, err)
		accAddress := sdk.AccAddress(initEvents[1][0].Value)
		require.Empty(t, err)

		queryRes, qErr = queryHelper(t, keeper, ctx, accAddress, `{"get_contract_version":{}}`, true, false, math.MaxUint64)
		require.Empty(t, qErr)

		require.Equal(t, queryRes, "1")
	}

	func TestV010InitV1ContractFromExecuteWithOkResponse(t *testing.T) {
		ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

		wasmCode, err := ioutil.ReadFile("./testdata/test-contract/contract.wasm")
		require.NoError(t, err)

		v010CodeID, err := keeper.Create(ctx, walletA, wasmCode, "", "", 0, 0, "title", "descr")
		require.NoError(t, err)

		_, _, contractAddress, _, err := initHelper(t, keeper, ctx, v010CodeID, walletA, privKeyA, `{"nop":{}}`, true, false, defaultGasForTests)
		require.Empty(t, err)

		queryRes, qErr := queryHelper(t, keeper, ctx, contractAddress, `{"get_contract_version":{}}`, true, false, math.MaxUint64)
		require.Empty(t, qErr)

		require.Equal(t, queryRes, "10")

		_, _, execData, execEvents, _, execErr := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, fmt.Sprintf(`{"callback_to_init":{"code_id":%d, "code_hash":"%s"}}`, codeID, codeHash), true. defaultGasForTests, 0)
		require.Empty(t, execErr)
		require.Empty(t, execData)

		accAddress := sdk.AccAddress(execEvents[1][0].Value)
		require.Empty(t, err)

		queryRes, qErr = queryHelper(t, keeper, ctx, accAddress, `{"get_contract_version":{}}`, true, false, math.MaxUint64)
		require.Empty(t, qErr)

		require.Equal(t, queryRes, "1")
	}

	func TestV010ExecuteV1ContractFromInitWithOkResponse(t *testing.T) {
		ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

		wasmCode, err := ioutil.ReadFile("./testdata/test-contract/contract.wasm")
		require.NoError(t, err)

		v010CodeID, err := keeper.Create(ctx, walletA, wasmCode, "", "", 0, 0, "title", "descr")
		require.NoError(t, err)

		_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"counter":{"counter":199, "expires":100}}`, true. defaultGasForTests)
		require.Empty(t, err)
		_, _, _, _, err = initHelper(t, keeper, ctx, v010CodeID, walletA, privKeyA, fmt.Sprintf(`{"call_to_exec":{"addr":"%s","code_hash":"%s","msg":"%s"}}`, contractAddress.String(), codeHash, `{\"increment\":{\"addition\": 1}}`), true. defaultGasForTests)
		require.Empty(t, err)

		queryRes, qErr := queryHelper(t, keeper, ctx, contractAddress, `{"get":{}}`, true. math.MaxUint64)
		require.Empty(t, qErr)

		// assert result is 32 byte sha256 hash (if hashed), or contractAddr if not
		var resp v1QueryResponse
		e := json.Unmarshal([]byte(queryRes), &resp)
		require.NoError(t, e)
		require.Equal(t, uint32(200), resp.Get.Count)
	}

	func TestV010ExecuteV1ContractFromExecuteWithOkResponse(t *testing.T) {
		ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

		wasmCode, err := ioutil.ReadFile("./testdata/test-contract/contract.wasm")
		require.NoError(t, err)

		v010CodeID, err := keeper.Create(ctx, walletA, wasmCode, "", "", 0, 0, "title", "descr")
		require.NoError(t, err)

		_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"counter":{"counter":299, "expires":100}}`, true. defaultGasForTests)
		_, _, v010ContractAddress, _, err := initHelper(t, keeper, ctx, v010CodeID, walletA, privKeyA, `{"nop":{}}`, true, false, defaultGasForTests)

		_, _, _, _, _, err = execHelper(t, keeper, ctx, v010ContractAddress, walletA, privKeyA, fmt.Sprintf(`{"call_to_exec":{"addr":"%s","code_hash":"%s","msg":"%s"}}`, contractAddress.String(), codeHash, `{\"increment\":{\"addition\": 1}}`), true. defaultGasForTests, 0)
		require.Empty(t, err)

		queryRes, qErr := queryHelper(t, keeper, ctx, contractAddress, `{"get":{}}`, true. math.MaxUint64)
		require.Empty(t, qErr)

		// assert result is 32 byte sha256 hash (if hashed), or contractAddr if not
		var resp v1QueryResponse
		e := json.Unmarshal([]byte(queryRes), &resp)
		require.NoError(t, e)
		require.Equal(t, uint32(300), resp.Get.Count)
	}

	func TestV010QueryV1ContractFromInitWithOkResponse(t *testing.T) {
		ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

		wasmCode, err := ioutil.ReadFile("./testdata/test-contract/contract.wasm")
		require.NoError(t, err)

		v010CodeID, err := keeper.Create(ctx, walletA, wasmCode, "", "", 0, 0, "title", "descr")
		require.NoError(t, err)

		_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true. defaultGasForTests)
		_, _, v010ContractAddress, events, err := initHelper(t, keeper, ctx, v010CodeID, walletA, privKeyA, fmt.Sprintf(`{"call_to_query":{"addr":"%s","code_hash":"%s","msg":"%s"}}`, contractAddress.String(), codeHash, `{\"receive_external_query_v1\":{\"num\":1}}`), true. defaultGasForTests)
		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: v010ContractAddress.Bytes(), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "c", Value: []byte("2"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	}

	func TestV010QueryV1ContractFromExecuteWithOkResponse(t *testing.T) {
		ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

		wasmCode, err := ioutil.ReadFile("./testdata/test-contract/contract.wasm")
		require.NoError(t, err)

		v010CodeID, err := keeper.Create(ctx, walletA, wasmCode, "", "", 0, 0, "title", "descr")
		require.NoError(t, err)

		_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true. defaultGasForTests)
		_, _, v010ContractAddress, _, err := initHelper(t, keeper, ctx, v010CodeID, walletA, privKeyA, `{"nop":{}}`, true, false, defaultGasForTests)

		_, _, _, events, _, err := execHelper(t, keeper, ctx, v010ContractAddress, walletA, privKeyA, fmt.Sprintf(`{"call_to_query":{"addr":"%s","code_hash":"%s","msg":"%s"}}`, contractAddress.String(), codeHash, `{\"receive_external_query_v1\":{\"num\":1}}`), true. defaultGasForTests, 0)
		require.Empty(t, err)
		require.Equal(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: v010ContractAddress.Bytes(), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "c", Value: []byte("2"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)
	}

	func TestV010InitV1ContractFromInitWithErrResponse(t *testing.T) {
		ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

		wasmCode, err := ioutil.ReadFile("./testdata/test-contract/contract.wasm")
		require.NoError(t, err)

		v010CodeID, err := keeper.Create(ctx, walletA, wasmCode, "", "", 0, 0, "title", "descr")
		require.NoError(t, err)

		_, _, _, _, err = initHelper(t, keeper, ctx, v010CodeID, walletA, privKeyA, fmt.Sprintf(`{"call_to_init":{"code_id":%d, "code_hash":"%s","contract_id":"blabla", "msg":"%s"}}`, codeID, codeHash, `{\"counter\":{\"counter\":0, \"expires\":100}}`), true. defaultGasForTests)
		require.Contains(t, fmt.Sprintf("%+v\n", err), "got wrong counter on init")
	}

	func TestV010InitV1ContractFromExecuteWithErrResponse(t *testing.T) {
		ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

		wasmCode, err := ioutil.ReadFile("./testdata/test-contract/contract.wasm")
		require.NoError(t, err)

		v010CodeID, err := keeper.Create(ctx, walletA, wasmCode, "", "", 0, 0, "title", "descr")
		require.NoError(t, err)

		_, _, contractAddress, _, err := initHelper(t, keeper, ctx, v010CodeID, walletA, privKeyA, `{"nop":{}}`, true, false, defaultGasForTests)

		queryRes, qErr := queryHelper(t, keeper, ctx, contractAddress, `{"get_contract_version":{}}`, true, false, math.MaxUint64)
		require.Empty(t, qErr)

		require.Equal(t, queryRes, "10")

		_, _, _, _, _, err = execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, fmt.Sprintf(`{"call_to_init":{"code_id":%d, "code_hash":"%s","contract_id":"blabla", "msg":"%s"}}`, codeID, codeHash, `{\"counter\":{\"counter\":0, \"expires\":100}}`), true. defaultGasForTests, 0)
		require.Contains(t, fmt.Sprintf("%+v\n", err), "got wrong counter on init")
	}

	func TestV010ExecuteV1ContractFromInitWithErrResponse(t *testing.T) {
		ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

		wasmCode, err := ioutil.ReadFile("./testdata/test-contract/contract.wasm")
		require.NoError(t, err)

		v010CodeID, err := keeper.Create(ctx, walletA, wasmCode, "", "", 0, 0, "title", "descr")
		require.NoError(t, err)

		_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"counter":{"counter":199, "expires":100}}`, true. defaultGasForTests)
		_, _, _, _, err = initHelper(t, keeper, ctx, v010CodeID, walletA, privKeyA, fmt.Sprintf(`{"call_to_exec":{"addr":"%s","code_hash":"%s","msg":"%s"}}`, contractAddress.String(), codeHash, `{\"increment\":{\"addition\": 0}}`), true. defaultGasForTests)

		require.Contains(t, fmt.Sprintf("%+v\n", err), "got wrong counter on increment")
	}

	func TestV010ExecuteV1ContractFromExecuteWithErrResponse(t *testing.T) {
		ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

		wasmCode, err := ioutil.ReadFile("./testdata/test-contract/contract.wasm")
		require.NoError(t, err)

		v010CodeID, err := keeper.Create(ctx, walletA, wasmCode, "", "", 0, 0, "title", "descr")
		require.NoError(t, err)

		_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"counter":{"counter":299, "expires":100}}`, true. defaultGasForTests)
		_, _, v010ContractAddress, _, err := initHelper(t, keeper, ctx, v010CodeID, walletA, privKeyA, `{"nop":{}}`, true, false, defaultGasForTests)

		_, _, _, _, _, err = execHelper(t, keeper, ctx, v010ContractAddress, walletA, privKeyA, fmt.Sprintf(`{"call_to_exec":{"addr":"%s","code_hash":"%s","msg":"%s"}}`, contractAddress.String(), codeHash, `{\"increment\":{\"addition\": 0}}`), true. defaultGasForTests, 0)
		require.Contains(t, fmt.Sprintf("%+v\n", err), "got wrong counter on increment")
	}

	func TestV010QueryV1ContractFromInitWithErrResponse(t *testing.T) {
		ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

		wasmCode, err := ioutil.ReadFile("./testdata/test-contract/contract.wasm")
		require.NoError(t, err)

		v010CodeID, err := keeper.Create(ctx, walletA, wasmCode, "", "", 0, 0, "title", "descr")
		require.NoError(t, err)

		_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true. defaultGasForTests)
		_, _, _, _, err = initHelper(t, keeper, ctx, v010CodeID, walletA, privKeyA, fmt.Sprintf(`{"call_to_query":{"addr":"%s","code_hash":"%s","msg":"%s"}}`, contractAddress.String(), codeHash, `{\"contract_error\":{\"error_type\":\"generic_err\"}}`), true. defaultGasForTests)
		require.Contains(t, fmt.Sprintf("%+v\n", err), "la la 🤯")
	}

	func TestV010QueryV1ContractFromExecuteWithErrResponse(t *testing.T) {
		ctx, keeper, codeID, codeHash, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

		wasmCode, err := ioutil.ReadFile("./testdata/test-contract/contract.wasm")
		require.NoError(t, err)

		v010CodeID, err := keeper.Create(ctx, walletA, wasmCode, "", "", 0, 0, "title", "descr")
		require.NoError(t, err)

		_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true. defaultGasForTests)
		_, _, v010ContractAddress, _, err := initHelper(t, keeper, ctx, v010CodeID, walletA, privKeyA, `{"nop":{}}`, true, false, defaultGasForTests)

		_, _, _, _, _, err = execHelper(t, keeper, ctx, v010ContractAddress, walletA, privKeyA, fmt.Sprintf(`{"call_to_query":{"addr":"%s","code_hash":"%s","msg":"%s"}}`, contractAddress.String(), codeHash, `{\"contract_error\":{\"error_type\":\"generic_err\"}}`), true. defaultGasForTests, 0)
		require.Contains(t, fmt.Sprintf("%+v\n", err), "la la 🤯")
	}
*/
func TestSendEncryptedAttributesFromInitWithoutSubmessageWithoutReply(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, events, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_attributes":{}}`, true, defaultGasForTests)
	require.Empty(t, err)

	require.Equal(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		events,
	)

}

func TestSendEncryptedAttributesFromInitWithSubmessageWithoutReply(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, events, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_attributes_with_submessage":{"id":0}}`, true, defaultGasForTests)
	require.Empty(t, err)
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[1])

	/*require.ElementsMatch(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			},
			{

				{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		events,
	)*/

}

func TestV1SendsEncryptedAttributesFromInitWithSubmessageWithReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, events, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_attributes_with_submessage":{"id":2200}}`, true, defaultGasForTests)
	require.Empty(t, err)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[1])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr5", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr6", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[2])
}

func TestSendEncryptedAttributesFromExecuteWithoutSubmessageWithoutReply(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_attributes":{}}`, true, defaultGasForTests, 0)
	require.Empty(t, err)

	require.Equal(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		events,
	)

}

func TestSendEncryptedAttributesFromExecuteWithSubmessageWithoutReply(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_attributes_with_submessage":{"id":0}}`, true, defaultGasForTests, 0)
	require.Empty(t, err)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[1])

	/*
		require.ElementsMatch(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				{

					{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)*/

}

func TestV1SendsEncryptedAttributesFromExecuteWithSubmessageWithReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_attributes_with_submessage":{"id":2200}}`, true, defaultGasForTests, 0)
	require.Empty(t, err)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[1])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr5", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr6", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[2])

}

func TestSendPlaintextFromInitWithoutSubmessageWithoutReply(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, events, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_plaintext_attributes":{}}`, true, defaultGasForTests, true)
	require.Empty(t, err)

	require.Equal(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		events,
	)

}

func TestSendPlaintextAttributesFromInitWithSubmessageWithoutReply(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, events, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_plaintext_attributes_with_submessage":{"id":0}}`, true, defaultGasForTests, true)
	require.Empty(t, err)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[1])
	/*
		require.ElementsMatch(t,
			[]ContractEvent{
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
				},
			},
			events,
		)*/

}

func TestV1SendsPlaintextAttributesFromInitWithSubmessageWithReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, events, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_plaintext_attributes_with_submessage":{"id":2300}}`, true, defaultGasForTests, true)
	require.Empty(t, err)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[1])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr5", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr6", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[2])
}

func TestSendPlaintextAttributesFromExecuteWithoutSubmessageWithoutReply(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_plaintext_attributes":{}}`, true, defaultGasForTests, 0, true)
	require.Empty(t, err)

	require.Equal(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		events,
	)

}

func TestSendPlaintextAttributesFromExecuteWithSubmessageWithoutReply(t *testing.T) {

	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, testContract.WasmFilePath, sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_plaintext_attributes_with_submessage":{"id":0}}`, true, defaultGasForTests, 0, true)
	require.Empty(t, err)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[1])

}

func TestV1SendsPlaintextAttributesFromExecuteWithSubmessageWithReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	_, _, _, events, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_plaintext_attributes_with_submessage":{"id":2300}}`, true, defaultGasForTests, 0, true)
	require.Empty(t, err)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[1])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr5", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr6", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, events[2])

}

func TestV1SendsEncryptedEventsFromInitWithoutSubmessageWithoutReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	nonce, ctx, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_events":{}}`, true, defaultGasForTests)

	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber2 := false
	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.Equal(t,
				[]cosmwasm.Attribute{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				attrs,
			)

			hadCyber1 = true
		}

		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.Equal(t,
				[]cosmwasm.Attribute{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				attrs,
			)

			hadCyber2 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber2)
}

func TestV1SendsEncryptedEventsFromInitWithSubmessageWithoutReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	nonce, ctx, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_events_with_submessage":{"id":0}}`, true, defaultGasForTests)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber2 := false
	hadCyber3 := false
	hadCyber4 := false

	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, _ := parseAndDecryptAttributes(e.Attributes, nonce)
			//require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber1 = true
		}

		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber2 = true
		}

		if e.Type == "wasm-cyber3" {
			require.False(t, hadCyber3)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "attr1", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🤟"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber3 = true
		}

		if e.Type == "wasm-cyber4" {
			require.False(t, hadCyber4)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("😅"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber4 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber2)
	require.True(t, hadCyber3)
	require.True(t, hadCyber4)
}

func TestV1SendsEncryptedEventsFromInitWithSubmessageWithReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	nonce, ctx, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_events_with_submessage":{"id":2400}}`, true, defaultGasForTests)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber2 := false
	hadCyber3 := false
	hadCyber4 := false
	hadCyber5 := false
	hadCyber6 := false

	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber1 = true
		}

		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber2 = true
		}

		if e.Type == "wasm-cyber3" {
			require.False(t, hadCyber3)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🤟"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber3 = true
		}

		if e.Type == "wasm-cyber4" {
			require.False(t, hadCyber4)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("😅"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber4 = true
		}

		if e.Type == "wasm-cyber5" {
			require.False(t, hadCyber5)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("😗"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("😋"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber5 = true
		}

		if e.Type == "wasm-cyber6" {
			require.False(t, hadCyber6)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("😉"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("😊"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber6 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber2)
	require.True(t, hadCyber3)
	require.True(t, hadCyber4)
	require.True(t, hadCyber5)
	require.True(t, hadCyber6)
}

func TestV1SendsEncryptedEventsFromExecuteWithoutSubmessageWithoutReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	nonce, ctx, _, _, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_events":{}}`, true, defaultGasForTests, 0)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber2 := false

	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.Equal(t,
				[]cosmwasm.Attribute{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				attrs,
			)

			hadCyber1 = true
		}

		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber2 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber2)
}

func TestV1SendsEncryptedEventsFromExecuteWithSubmessageWithoutReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	nonce, ctx, _, _, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_events_with_submessage":{"id":0}}`, true, defaultGasForTests, 0)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber2 := false
	hadCyber3 := false
	hadCyber4 := false

	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber1 = true
		}

		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber2 = true
		}

		if e.Type == "wasm-cyber3" {
			require.False(t, hadCyber3)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🤟"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber3 = true
		}

		if e.Type == "wasm-cyber4" {
			require.False(t, hadCyber4)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("😅"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber4 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber2)
	require.True(t, hadCyber3)
	require.True(t, hadCyber4)
}

func TestV1SendsEncryptedEventsFromExecuteWithSubmessageWithReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	nonce, ctx, _, _, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_events_with_submessage":{"id":2400}}`, true, defaultGasForTests, 0)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber2 := false
	hadCyber3 := false
	hadCyber4 := false
	hadCyber5 := false
	hadCyber6 := false

	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.Equal(t,
				[]cosmwasm.Attribute{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				attrs,
			)

			hadCyber1 = true
		}

		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber2 = true
		}

		if e.Type == "wasm-cyber3" {
			require.False(t, hadCyber3)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🤟"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber3 = true
		}

		if e.Type == "wasm-cyber4" {
			require.False(t, hadCyber4)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("😅"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber4 = true
		}

		if e.Type == "wasm-cyber5" {
			require.False(t, hadCyber5)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("😗"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("😋"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber5 = true
		}

		if e.Type == "wasm-cyber6" {
			require.False(t, hadCyber6)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("😉"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("😊"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber6 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber2)
	require.True(t, hadCyber3)
	require.True(t, hadCyber4)
	require.True(t, hadCyber5)
	require.True(t, hadCyber6)
}

func TestV1SendsMixedLogsFromInitWithoutSubmessageWithoutReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	nonce, ctx, contractAddress, logs, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_mixed_attributes_and_events":{}}`, true, defaultGasForTests, true)

	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, hasPlaintext := parseAndDecryptAttributes(e.Attributes, nonce)
			require.NotEmpty(t, hasPlaintext)
			require.Equal(t,
				[]cosmwasm.Attribute{
					{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				attrs,
			)

			hadCyber1 = true
		}
	}

	require.True(t, hadCyber1)

	require.ElementsMatch(t,
		ContractEvent{

			{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
			{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
			{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		},
		logs[0],
	)
}

func TestV1SendsMixedAttributesAndEventsFromInitWithSubmessageWithoutReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	nonce, ctx, contractAddress, logs, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_mixed_attributes_and_events_with_submessage":{"id":0}}`, true, defaultGasForTests)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber2 := false

	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber1 = true
		}

		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr5", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr6", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber2 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber2)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, logs[0])

}

func TestV1SendsMixedAttributesAndEventsFromInitWithSubmessageWithReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	nonce, ctx, contractAddress, logs, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"add_mixed_attributes_and_events_with_submessage":{"id":2500}}`, true, defaultGasForTests)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber2 := false
	hadCyber3 := false

	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber1 = true
		}

		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr5", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr6", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber2 = true
		}

		if e.Type == "wasm-cyber3" {
			require.False(t, hadCyber3)
			attrs, _ := parseAndDecryptAttributes(e.Attributes, nonce)
			//require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr9", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr10", Value: []byte("🤟"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber3 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber2)
	require.True(t, hadCyber3)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, logs[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
	}, logs[1])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr11", Value: []byte("😉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr12", Value: []byte("😊"), AccAddr: "", Encrypted: false, PubDb: false},
	}, logs[2])

}

func TestV1SendsMixedAttributesAndEventsFromExecuteWithoutSubmessageWithoutReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	nonce, ctx, _, logs, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_mixed_attributes_and_events":{}}`, true, defaultGasForTests, 0, true)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false

	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, hasPlaintext := parseAndDecryptAttributes(e.Attributes, nonce)
			require.NotEmpty(t, hasPlaintext)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber1 = true
		}
	}

	require.True(t, hadCyber1)

	require.ElementsMatch(t,
		ContractEvent{

			{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
			{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
			{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
		},
		logs[0],
	)
}

func TestV1SendsMixedAttributesAndEventsFromExecuteWithSubmessageWithoutReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	nonce, ctx, _, logs, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_mixed_attributes_and_events_with_submessage":{"id":0}}`, true, defaultGasForTests, 0)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber2 := false

	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber1 = true
		}

		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr5", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr6", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber2 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber2)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, logs[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr7", Value: []byte("😅"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr8", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, logs[2])

}

func TestV1SendsMixedAttributesAndEventsFromExecuteWithSubmessageWithReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	_, _, contractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true, defaultGasForTests)
	require.Empty(t, err)
	nonce, ctx, _, logs, _, err := execHelper(t, keeper, ctx, contractAddress, walletA, privKeyA, `{"add_mixed_attributes_and_events_with_submessage":{"id":2500}}`, true, defaultGasForTests, 0)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber2 := false
	hadCyber3 := false

	for _, e := range events {
		fmt.Printf("ev %+v \n", e)
		fmt.Printf("Attributes %+v \n", e.Attributes)
		fmt.Printf("Attribute2 %s \n", e.Attributes[1].Key)
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)
			hadCyber1 = true
		}

		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr5", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr6", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber2 = true
		}

		if e.Type == "wasm-cyber3" {
			require.False(t, hadCyber3)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.ElementsMatch(t, ContractEvent{
				{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr9", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr10", Value: []byte("🤟"), AccAddr: "", Encrypted: false, PubDb: false},
			}, attrs)

			hadCyber3 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber2)
	require.True(t, hadCyber3)

	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, logs[0])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr7", Value: []byte("😅"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr8", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
	}, logs[2])
	require.ElementsMatch(t, ContractEvent{
		{Key: "contract_address", Value: []byte(contractAddress.String()), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr11", Value: []byte("😉"), AccAddr: "", Encrypted: false, PubDb: false},
		{Key: "attr12", Value: []byte("😊"), AccAddr: "", Encrypted: false, PubDb: false},
	}, logs[4])
}

/*
func TestV1SendsLogsMixedWithV010WithoutReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	wasmCode, err := ioutil.ReadFile("./testdata/test-contract/contract.wasm")
	require.NoError(t, err)

	v010CodeID, err := keeper.Create(ctx, walletA, wasmCode, "", "", 0, 0, "title", "descr")
	require.NoError(t, err)

	v010CodeHash := hex.EncodeToString(keeper.GetCodeInfo(ctx, v010CodeID).CodeHash)

	_, _, v010ContractAddress, _, err := initHelper(t, keeper, ctx, v010CodeID, walletA, privKeyA, `{"nop":{}}`, true, false, defaultGasForTests)
	require.Empty(t, err)
	_, _, v1ContractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true. defaultGasForTests)
	require.Empty(t, err)
	nonce, ctx, _, logs, _, err := execHelper(t, keeper, ctx, v1ContractAddress, walletA, privKeyA, fmt.Sprintf(`{"add_attributes_from_v010":{"addr":"%s","code_hash":"%s", "id":0}}`, v010ContractAddress, v010CodeHash), true. defaultGasForTests, 0)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false

	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.Equal(t,
				[]cosmwasm.Attribute{
					{Key: "contract_address", Value: v1ContractAddress.Bytes(), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				attrs,
			)

			hadCyber1 = true
		}
	}

	require.True(t, hadCyber1)

	require.Equal(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: v1ContractAddress.Bytes(), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			},
			{
				{Key: "contract_address", Value: v010ContractAddress.Bytes(), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		logs,
	)
}

func TestV1SendsLogsMixedWithV010WithReply(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	wasmCode, err := ioutil.ReadFile("./testdata/test-contract/contract.wasm")
	require.NoError(t, err)

	v010CodeID, err := keeper.Create(ctx, walletA, wasmCode, "", "", 0, 0, "title", "descr")
	require.NoError(t, err)

	v010CodeHash := hex.EncodeToString(keeper.GetCodeInfo(ctx, v010CodeID).CodeHash)

	_, _, v010ContractAddress, _, err := initHelper(t, keeper, ctx, v010CodeID, walletA, privKeyA, `{"nop":{}}`, true, false, defaultGasForTests)
	require.Empty(t, err)
	_, _, v1ContractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true. defaultGasForTests)
	require.Empty(t, err)
	nonce, ctx, _, logs, _, err := execHelper(t, keeper, ctx, v1ContractAddress, walletA, privKeyA, fmt.Sprintf(`{"add_attributes_from_v010":{"addr":"%s","code_hash":"%s", "id":2500}}`, v010ContractAddress, v010CodeHash), true. defaultGasForTests, 0)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber1 := false
	hadCyber3 := false

	for _, e := range events {
		if e.Type == "wasm-cyber1" {
			require.False(t, hadCyber1)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.Equal(t,
				[]cosmwasm.Attribute{
					{Key: "contract_address", Value: v1ContractAddress.Bytes(), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				attrs,
			)

			hadCyber1 = true
		}

		if e.Type == "wasm-cyber3" {
			require.False(t, hadCyber3)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.Equal(t,
				[]cosmwasm.Attribute{
					{Key: "contract_address", Value: v1ContractAddress.Bytes(), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr9", Value: []byte("🤯"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr10", Value: []byte("🤟"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				attrs,
			)

			hadCyber3 = true
		}
	}

	require.True(t, hadCyber1)
	require.True(t, hadCyber3)

	require.Equal(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: v1ContractAddress.Bytes(), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			},
			{
				{Key: "contract_address", Value: v010ContractAddress.Bytes(), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr3", Value: []byte("🍉"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr4", Value: []byte("🥝"), AccAddr: "", Encrypted: false, PubDb: false},
			},
			{
				{Key: "contract_address", Value: v1ContractAddress.Bytes(), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr11", Value: []byte("😉"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr12", Value: []byte("😊"), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		logs,
	)
}

func TestV010SendsLogsMixedWithV1(t *testing.T) {
	ctx, keeper, codeID, v1CodeHash, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	wasmCode, err := ioutil.ReadFile("./testdata/test-contract/contract.wasm")
	require.NoError(t, err)

	v010CodeID, err := keeper.Create(ctx, walletA, wasmCode, "", "", 0, 0, "title", "descr")
	require.NoError(t, err)

	_, _, v010ContractAddress, _, err := initHelper(t, keeper, ctx, v010CodeID, walletA, privKeyA, `{"nop":{}}`, true, false, defaultGasForTests)
	require.Empty(t, err)
	_, _, v1ContractAddress, _, err := initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"nop":{}}`, true. defaultGasForTests)
	require.Empty(t, err)
	nonce, ctx, _, logs, _, err := execHelper(t, keeper, ctx, v010ContractAddress, walletA, privKeyA, fmt.Sprintf(`{"add_mixed_events_and_attributes_from_v1":{"addr":"%s","code_hash":"%s"}}`, v1ContractAddress, v1CodeHash), true. defaultGasForTests, 0)
	require.Empty(t, err)

	events := ctx.EventManager().Events()

	hadCyber2 := false

	for _, e := range events {
		if e.Type == "wasm-cyber2" {
			require.False(t, hadCyber2)
			attrs, err := parseAndDecryptAttributes(e.Attributes, nonce)
			require.Empty(t, err)

			require.Equal(t,
				[]cosmwasm.Attribute{
					{Key: "contract_address", Value: v1ContractAddress.Bytes(), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr5", Value: []byte("🐙"), AccAddr: "", Encrypted: false, PubDb: false},
					{Key: "attr6", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				},
				attrs,
			)

			hadCyber2 = true
		}
	}

	require.True(t, hadCyber2)

	require.Equal(t,
		[]ContractEvent{
			{
				{Key: "contract_address", Value: v010ContractAddress.Bytes(), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr1", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr2", Value: []byte("🌈"), AccAddr: "", Encrypted: false, PubDb: false},
			},
			{
				{Key: "contract_address", Value: v1ContractAddress.Bytes(), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr7", Value: []byte("😅"), AccAddr: "", Encrypted: false, PubDb: false},
				{Key: "attr8", Value: []byte("🦄"), AccAddr: "", Encrypted: false, PubDb: false},
			},
		},
		logs,
	)
}
*/

/*todo find out why it does not retreive panic
func TestSubmessageGasExceedingMessageGas(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	defer func() {
		r := recover()
		require.NotNil(t, r)
		_, ok := r.(sdk.ErrorOutOfGas)
		require.True(t, ok, "%+v\n", r)
	}()
	_, _, _, _, _ = initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"measure_gas_for_submessage":{"id":0}}`, false, defaultGasForTests)
}

func TestReplyGasExceedingMessageGas(t *testing.T) {
	ctx, keeper, codeID, _, walletA, privKeyA, _, _ := setupTest(t, "./testdata/test-contract/contract.wasm", sdk.NewCoins())

	defer func() {
		r := recover()
		require.NotNil(t, r)
		_, ok := r.(sdk.ErrorOutOfGas)
		require.True(t, ok, "%+v\n", r)
	}()
	_, _, _, _, _ = initHelper(t, keeper, ctx, codeID, walletA, privKeyA, `{"measure_gas_for_submessage":{"id":2600}}`, false, defaultGasForTests)
}
*/
