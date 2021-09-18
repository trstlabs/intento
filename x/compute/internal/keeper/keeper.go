package keeper

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	//"log"
	"path/filepath"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	sdktxsigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	wasm "github.com/danieljdd/tpp/go-cosmwasm"
	wasmTypes "github.com/danieljdd/tpp/go-cosmwasm/types"

	"github.com/danieljdd/tpp/x/compute/internal/types"
	"github.com/tendermint/tendermint/libs/log"
)

// GasMultiplier is how many cosmwasm gas points = 1 sdk gas point
// SDK reference costs can be found here: https://github.com/cosmos/cosmos-sdk/blob/02c6c9fafd58da88550ab4d7d494724a477c8a68/store/types/gas.go#L153-L164
// A write at ~3000 gas and ~200us = 10 gas per us (microsecond) cpu/io
// Rough timing have 88k gas at 90us, which is equal to 1k sdk gas... (one read)
//
// Please not that all gas prices returned to the wasmer engine should have this multiplied
const GasMultiplier uint64 = 100

// MaxGas for a contract is 10 billion wasmer gas (enforced in rust to prevent overflow)
// The limit for v0.9.3 is defined here: https://github.com/CosmWasm/cosmwasm/blob/v0.9.3/packages/vm/src/backends/singlepass.rs#L15-L23
// (this will be increased in future releases)
const MaxGas = 10_000_000_000

// InstanceCost is how much SDK gas we charge each time we load a WASM instance.
// Creating a new instance is costly, and this helps put a recursion limit to contracts calling contracts.
const InstanceCost uint64 = 40_000

// CompileCost is how much SDK gas we charge *per byte* for compiling WASM code.
const CompileCost uint64 = 2

// Keeper will have a reference to Wasmer with it's own data directory.
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.Marshaler
	legacyAmino   codec.LegacyAmino
	accountKeeper authkeeper.AccountKeeper
	bankKeeper    bankkeeper.Keeper
	wasmer        wasm.Wasmer
	queryPlugins  QueryPlugins
	messenger     MessageHandler
	// queryGasLimit is the max wasm gas that can be spent on executing a query with a contract
	queryGasLimit uint64
	paramSpace    paramtypes.Subspace

	// authZPolicy   AuthorizationPolicy
	//paramSpace    subspace.Subspace
}

// NewKeeper creates a new contract Keeper instance
// If customEncoders is non-nil, we can use this to override some of the message handler, especially custom
func NewKeeper(cdc codec.Marshaler, legacyAmino codec.LegacyAmino, storeKey sdk.StoreKey, accountKeeper authkeeper.AccountKeeper,
	bankKeeper bankkeeper.Keeper, govKeeper govkeeper.Keeper, distKeeper distrkeeper.Keeper, mintKeeper mintkeeper.Keeper, stakingKeeper stakingkeeper.Keeper,
	router sdk.Router, homeDir string, wasmConfig types.WasmConfig, supportedFeatures string, customEncoders *MessageEncoders, customPlugins *QueryPlugins, paramSpace paramtypes.Subspace) Keeper {
	wasmer, err := wasm.NewWasmer(filepath.Join(homeDir, "wasm"), supportedFeatures, wasmConfig.CacheSize)
	if err != nil {
		panic(err)
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(ParamKeyTable())
	}

	keeper := Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		legacyAmino:   legacyAmino,
		wasmer:        *wasmer,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		messenger:     NewMessageHandler(router, customEncoders),
		queryGasLimit: wasmConfig.SmartQueryGasLimit,
		paramSpace:    paramSpace,
		// authZPolicy:   DefaultAuthorizationPolicy{},
		//paramSpace:    paramSpace,
	}
	keeper.queryPlugins = DefaultQueryPlugins(govKeeper, distKeeper, mintKeeper, bankKeeper, stakingKeeper, &keeper).Merge(customPlugins)
	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Create uploads and compiles a WASM contract, returning a short identifier for the contract
func (k Keeper) Create(ctx sdk.Context, creator sdk.AccAddress, wasmCode []byte, source string, builder string, endTime time.Duration) (codeID uint64, err error) {

	wasmCode, err = uncompress(wasmCode)
	if err != nil {
		return 0, sdkerrors.Wrap(types.ErrCreateFailed, err.Error())
	}
	ctx.GasMeter().ConsumeGas(CompileCost*uint64(len(wasmCode)), "Compiling WASM Bytecode")

	codeHash, err := k.wasmer.Create(wasmCode)
	if err != nil {
		// return 0, sdkerrors.Wrap(err, "cosmwasm create")
		return 0, sdkerrors.Wrap(types.ErrCreateFailed, err.Error())
	}

	//hash = string(codeHash)
	store := ctx.KVStore(k.storeKey)
	codeID = k.autoIncrementID(ctx, types.KeyLastCodeID)
	/*
		if instantiateAccess == nil {
			defaultAccessConfig := k.getInstantiateAccessConfig(ctx).With(creator)
			instantiateAccess = &defaultAccessConfig
		}
	*/
	codeInfo := types.NewCodeInfo(codeHash, creator, source, builder, endTime /* , *instantiateAccess */)
	// 0x01 | codeID (uint64) -> ContractInfo
	store.Set(types.GetCodeKey(codeID), k.cdc.MustMarshalBinaryBare(&codeInfo))

	return codeID, nil
}

func (k Keeper) importCode(ctx sdk.Context, codeID uint64, codeInfo types.CodeInfo, wasmCode []byte) error {
	wasmCode, err := uncompress(wasmCode)
	if err != nil {
		return sdkerrors.Wrap(types.ErrCreateFailed, err.Error())
	}
	newCodeHash, err := k.wasmer.Create(wasmCode)
	if err != nil {
		return sdkerrors.Wrap(types.ErrCreateFailed, err.Error())
	}
	if !bytes.Equal(codeInfo.CodeHash, newCodeHash) {
		return sdkerrors.Wrap(types.ErrInvalid, "code hashes not same")
	}

	store := ctx.KVStore(k.storeKey)
	key := types.GetCodeKey(codeID)
	if store.Has(key) {
		return sdkerrors.Wrapf(types.ErrDuplicate, "duplicate code: %d", codeID)
	}
	// 0x01 | codeID (uint64) -> ContractInfo
	store.Set(key, k.cdc.MustMarshalBinaryBare(&codeInfo))
	return nil
}

func (k Keeper) GetSignerInfo(ctx sdk.Context, signer sdk.AccAddress) ([]byte, []byte, error) {
	tx := sdktx.Tx{}
	err := k.cdc.UnmarshalBinaryBare(ctx.TxBytes(), &tx)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(types.ErrInstantiateFailed, fmt.Sprintf("Unable to decode transaction from bytes: %s", err.Error()))
	}

	// for MsgInstantiateContract, there is only one signer which is msg.Sender
	// (https://github.com/danieljdd/tpp/blob/d7813792fa07b93a10f0885eaa4c5e0a0a698854/x/compute/internal/types/msg.go#L192-L194)
	signerAcc, err := ante.GetSignerAcc(ctx, k.accountKeeper, signer)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(types.ErrInstantiateFailed, fmt.Sprintf("Unable to retrieve account by address: %s", err.Error()))
	}

	txConfig := authtx.NewTxConfig(k.cdc.(*codec.ProtoCodec), []sdktxsigning.SignMode{sdktxsigning.SignMode_SIGN_MODE_DIRECT})
	modeHandler := txConfig.SignModeHandler()
	signingData := authsigning.SignerData{
		ChainID:       ctx.ChainID(),
		AccountNumber: signerAcc.GetAccountNumber(),
	}

	protobufTx := authtx.WrapTx(&tx).GetTx()
	signBytes, err := modeHandler.GetSignBytes(sdktxsigning.SignMode_SIGN_MODE_DIRECT, signingData, protobufTx)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(types.ErrInstantiateFailed, fmt.Sprintf("Unable to recreate sign bytes for the tx: %s", err.Error()))
	}

	// The first signature is the signature of the message sender,
	// according to the docstring of `tx.AuthInfo.SignerInfos`
	return tx.Signatures[0], signBytes, nil
}

// Instantiate creates an instance of a WASM contract
func (k Keeper) Instantiate(ctx sdk.Context, codeID uint64, creator /* , admin */ sdk.AccAddress, initMsg []byte, label string, deposit sdk.Coins, callbackSig []byte) (sdk.AccAddress, error) {

	ctx.GasMeter().ConsumeGas(InstanceCost, "Loading CosmWasm module: init")

	signerSig := []byte{}
	signBytes := []byte{}
	var err error

	// If no callback signature - we should send the actual msg sender sign bytes and signature
	if callbackSig == nil {
		signerSig, signBytes, err = k.GetSignerInfo(ctx, creator)
		if err != nil {
			return nil, err
		}
	}

	verificationInfo := types.NewVerificationInfo(signBytes, signerSig, callbackSig)

	// create contract address

	store := ctx.KVStore(k.storeKey)
	existingAddress := store.Get(types.GetContractLabelPrefix(label))

	if existingAddress != nil {
		return nil, sdkerrors.Wrap(types.ErrAccountExists, label)
	}

	contractAddress := k.generateContractAddress(ctx, codeID)
	existingAcct := k.accountKeeper.GetAccount(ctx, contractAddress)
	if existingAcct != nil {
		return nil, sdkerrors.Wrap(types.ErrAccountExists, existingAcct.GetAddress().String())
	}

	// deposit initial contract funds
	if !deposit.IsZero() {
		if k.bankKeeper.BlockedAddr(creator) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "blocked address can not be used")
		}
		sdkerr := k.bankKeeper.SendCoins(ctx, creator, contractAddress, deposit)
		if sdkerr != nil {
			return nil, sdkerr
		}
	} else {
		// create an empty account (so we don't have issues later)
		// TODO: can we remove this?
		contractAccount := k.accountKeeper.NewAccountWithAddress(ctx, contractAddress)
		k.accountKeeper.SetAccount(ctx, contractAccount)
	}

	// get contact info

	bz := store.Get(types.GetCodeKey(codeID))
	if bz == nil {
		return nil, sdkerrors.Wrap(types.ErrNotFound, "code")
	}

	var codeInfo types.CodeInfo
	k.cdc.MustUnmarshalBinaryBare(bz, &codeInfo)

	// if !authZ.CanInstantiateContract(codeInfo.InstantiateConfig, creator) {
	// 	return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "can not instantiate")
	// }

	// prepare params for contract instantiate call
	params := types.NewEnv(ctx, creator, deposit, contractAddress, nil)

	// create prefixed data store
	// 0x03 | contractAddress (sdk.AccAddress)
	prefixStoreKey := types.GetContractStorePrefixKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)

	// prepare querier
	querier := QueryHandler{
		Ctx:     ctx,
		Plugins: k.queryPlugins,
	}
	//fmt.Printf("Init message before wasm is %s \n", base64.StdEncoding.EncodeToString(initMsg))
	// instantiate wasm contract
	gas := gasForContract(ctx)
	res, key, gasUsed, err := k.wasmer.Instantiate(codeInfo.CodeHash, params, initMsg, prefixStore, cosmwasmAPI, querier, ctx.GasMeter(), gas, verificationInfo)
	consumeGas(ctx, gasUsed)
	//fmt.Print("Init message after wasm is  \n", res.Messages[0])
	//fmt.Print("Init message after wasm is  \n", res.Messages[1])
	//fmt.Printf("Init message after wasm is  \n", res.Log)
	if err != nil {
		return contractAddress, sdkerrors.Wrap(types.ErrInstantiateFailed, err.Error())
	}

	// emit all events from this contract itself
	events := types.ParseEvents(res.Log, contractAddress)
	ctx.EventManager().EmitEvents(events)

	// persist instance
	createdAt := types.NewAbsoluteTxPosition(ctx)
	var instance = types.ContractInfo{
		CodeID:     codeID,
		Creator:    creator,
		ContractId: label,
		Created:    createdAt,
	}
	//	instance.CodeID = codeID
	//	instance.Creator = creator
	//	instance.Label = label
	//	instance.Created = createdAt
	//instance.LastMsg = nil
	//instance = types.NewContractInfo(codeID, creator /* admin, */, label, createdAt, nil)
	store.Set(types.GetContractAddressKey(contractAddress), k.cdc.MustMarshalBinaryBare(&instance))

	//	fmt.Printf("Storing key: %s for account %s\n", base64.StdEncoding.EncodeToString(key), contractAddress)

	store.Set(types.GetContractEnclaveKey(contractAddress), key)

	store.Set(types.GetContractLabelPrefix(label), contractAddress)

	err = k.dispatchMessages(ctx, contractAddress, res.Messages)
	if err != nil {
		return nil, err
	}

	// k.appendToContractHistory(ctx, contractAddress, instance.InitialHistory(initMsg))
	return contractAddress, nil
}

// Execute executes the contract instance
func (k Keeper) Execute(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, msg []byte, coins sdk.Coins, callbackSig []byte) (*sdk.Result, error) {
	ctx.GasMeter().ConsumeGas(InstanceCost, "Loading Compute module: execute")

	signerSig := []byte{}
	signBytes := []byte{}
	var err error

	if callbackSig == nil {
		signerSig, signBytes, err = k.GetSignerInfo(ctx, caller)
		if err != nil {
			return nil, err
		}
	}

	verificationInfo := types.NewVerificationInfo(signBytes, signerSig, callbackSig)

	codeInfo, prefixStore, err := k.contractInstance(ctx, contractAddress)
	if err != nil {
		return nil, err
	}

	store := ctx.KVStore(k.storeKey)

	// add funds
	if !coins.IsZero() {
		if k.bankKeeper.BlockedAddr(caller) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "blocked address can not be used")
		}

		sdkerr := k.bankKeeper.SendCoins(ctx, caller, contractAddress, coins)
		if sdkerr != nil {
			return nil, sdkerr
		}
	}

	contractKey := store.Get(types.GetContractEnclaveKey(contractAddress))
	if contractKey == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "contract key not found")
	}
	fmt.Printf("Contract Execute: Got contract Key for contract %s: %s\n", contractAddress, base64.StdEncoding.EncodeToString(contractKey))
	params := types.NewEnv(ctx, caller, coins, contractAddress, contractKey)
	fmt.Printf("Contract Execute: key from params %s \n", params.Key)

	// prepare querier
	querier := QueryHandler{
		Ctx:     ctx,
		Plugins: k.queryPlugins,
	}

	gas := gasForContract(ctx)
	//	fmt.Printf("Execute message before wasm is %s \n", base64.StdEncoding.EncodeToString(msg))
	res, gasUsed, execErr := k.wasmer.Execute(codeInfo.CodeHash, params, msg, prefixStore, cosmwasmAPI, querier, gasMeter(ctx), gas, verificationInfo)
	consumeGas(ctx, gasUsed)

	if execErr != nil {
		return nil, sdkerrors.Wrap(types.ErrExecuteFailed, execErr.Error())
	}

	//var res wasmTypes.CosmosResponse
	//err = json.Unmarshal(res, &res)
	//if err != nil {
	//	return sdk.Result{}, err
	//}
	//fmt.Printf("result: Got result for item data: %s | log: %s\n", res.Data, res.Log)
	var raw map[string]json.RawMessage
	_ = json.Unmarshal([]byte(res.Log[0].Value), &raw)
	//fmt.Printf("log: Got res raw for item %s: %s\n", raw, res)

	// emit all events from this contract itself
	events := types.ParseEvents(res.Log, contractAddress)
	ctx.EventManager().EmitEvents(events)

	// TODO: capture events here as well
	err = k.dispatchMessages(ctx, contractAddress, res.Messages)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{
		//	Data: []byte(res.Log[0].Value),
		Data: res.Data,
		Log:  res.Log[0].Value,
	}, nil
}

/*
// Delete deletes the contract instance
func (k Keeper) ContractPayout(ctx sdk.Context, contractAddress sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)
	contractBz := store.Get(types.GetContractAddressKey(contractAddress))
	if contractBz == nil {
		return sdkerrors.Wrap(types.ErrNotFound, "contract")
	}
	var contract types.ContractInfo
	k.cdc.MustUnmarshalBinaryBare(contractBz, &contract)

	//payout contract coins to the creator
	balance := k.bankKeeper.GetAllBalances(ctx, contractAddress)
	if !balance.Empty() {
		k.bankKeeper.SendCoins(ctx, contractAddress, contract.Creator, balance)
	}
	return nil
}

// CallLastMsg executes a final message before end-blocker deletion
func (k Keeper) CallLastMsg(ctx sdk.Context, contractAddress sdk.AccAddress) (err error) {


	//msgLast =
	//var err
	lastMsg := types.SecretMsg{}
	last := types.ParseLast{}

	//initMsg.Msg = []byte("{\"estimationcount\": \"3\"}")
	lastMsg.Msg, err = json.Marshal(last)
	if err != nil {
		return err
	}

	//queryClient := types.NewQueryClient(ctx.Context())

	//get codeid first
	hash := k.GetCodeHash(ctx, codeid)

	var encryptedMsg []byte
	lastMsg.CodeHash = []byte(hex.EncodeToString(hash))
	encryptedMsg, err = wasmCtx.Encrypt(deletegMsg.Serialize())
	if err != nil {
		return err
	}

	res, err := k.Execute(ctx, contractAddress, contractAddress, encryptedMsg, sdk.NewCoins(sdk.NewCoin("tpp", sdk.ZeroInt())), nil)
	if err != nil {
		panic(err)
	}

	return nil
}
*/
// Delete deletes the contract instance
func (k Keeper) Delete(ctx sdk.Context, contractAddress sdk.AccAddress) error {

	_, prefixStore, err := k.contractInstance(ctx, contractAddress)
	if err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)

	contractBz := store.Get(types.GetContractAddressKey(contractAddress))
	if contractBz == nil {
		return sdkerrors.Wrap(types.ErrNotFound, "contract")
	}
	var contract types.ContractInfo
	k.cdc.MustUnmarshalBinaryBare(contractBz, &contract)

	//	var codeInfo types.CodeInfo
	//k.cdc.MustUnmarshalBinaryBare(contractInfoBz, &codeInfo)
	prefixStoreKey := types.GetContractStorePrefixKey(contractAddress)
	//prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	prefixStore.Delete(prefixStoreKey)

	//prefixStore.Delete()
	store.Delete(types.GetContractEnclaveKey(contractAddress))
	store.Delete(types.GetContractLabelPrefix(contract.ContractId))
	store.Delete(types.GetContractAddressKey(contractAddress))

	//store.Delete(types.GetCodeKey(contract.CodeID))

	return nil
}

/*
func (k Keeper) DeleteContract(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, msg []byte) (*sdk.Result, error) {
	ctx.GasMeter().ConsumeGas(InstanceCost, "Loading CosmWasm module: Delete Contract")

	//lets see if the module acc is able to pass signer info
	signerSig := []byte{}
	signBytes := []byte{}
	var err error
	signerSig, signBytes, err = k.GetSignerInfo(ctx, caller)
	if err != nil {
		return nil, err
	}

	verificationInfo := types.NewVerificationInfo(signBytes, signerSig, nil)
	/*
		_ = authtypes.StdSignature{
			PubKey:    secp256k1.PubKeySecp256k1{},
			Signature: []byte{},
		}
		_ = []byte{}

		tx := authtypes.StdTx{}
		txBytes := ctx.TxBytes()
		err := k.cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
		if err != nil {
			return &sdk.Result{}, sdkerrors.Wrap(types.ErrInstantiateFailed, fmt.Sprintf("Unable to decode transaction from bytes: %s", err.Error()))
		}

		contractInfo := k.GetContractInfo(ctx, contractAddress)
		if contractInfo == nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "unknown contract")
		}	*/
/*if !k.authZPolicy.CanModifyContract(contractInfo.Creator, caller) {
	return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "can not Delete Contract")
}*/

/*newCodeInfo := k.GetCodeInfo(ctx, newCodeID)
	if newCodeInfo == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "unknown code")
	}


	codeInfo, prefixStore, err := k.contractInstance(ctx, contractAddress)
	if err != nil {
		return nil, err
	}

	store := ctx.KVStore(k.storeKey)
	contractKey := store.Get(types.GetContractEnclaveKey(contractAddress))

	var noDeposit sdk.Coins
	params := types.NewEnv(ctx, caller, noDeposit, contractAddress, contractKey)

	fmt.Printf("Contract Delete: key from params %s \n", params.Key)

	// prepare querier
	querier := QueryHandler{
		Ctx:     ctx,
		Plugins: k.queryPlugins,
	}

	//prefixStoreKey := types.GetContractStorePrefixKey(contractAddress)
	//prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	gas := gasForContract(ctx)
	//res, gasUsed, err := k.wasmer.Migrate(codeInfo.CodeHash, params, msg, &prefixStore, cosmwasmAPI, &querier, gasMeter(ctx), gas)
	//consumeGas(ctx, gasUsed)

	res, gasUsed, err := k.wasmer.Execute(codeInfo.CodeHash, params, msg, &prefixStore, cosmwasmAPI, &querier, gasMeter(ctx), gas, verificationInfo)
	consumeGas(ctx, gasUsed)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrMigrationFailed, err.Error())
	}

	// emit all events from this contract itself
	events := types.ParseEvents(res.Log, contractAddress)
	ctx.EventManager().EmitEvents(events)

	//historyEntry := contractInfo.AddMigration(ctx, newCodeID, msg)
	//k.appendToContractHistory(ctx, contractAddress, historyEntry)
	//k.setContractInfo(ctx, contractAddress, contractInfo)

	////k.DeleteContractInfo
	if err := k.dispatchMessages(ctx, contractAddress, res.Messages); err != nil {
		return nil, sdkerrors.Wrap(err, "dispatch")
	}

	return &sdk.Result{
		Data: res.Data,
	}, nil
}


// We don't use this function currently. It's here for upstream compatibility
// Migrate allows to upgrade a contract to a new code with data migration.
func (k Keeper) Migrate(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, newCodeID uint64, msg []byte) (*sdk.Result, error) {
	return k.migrate(ctx, contractAddress, caller, newCodeID, msg, k.authZPolicy)
}

func (k Keeper) migrate(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, newCodeID uint64, msg []byte, authZ AuthorizationPolicy) (*sdk.Result, error) {
	ctx.GasMeter().ConsumeGas(InstanceCost, "Loading CosmWasm module: migrate")
	_ = authtypes.StdSignature{
		PubKey:    secp256k1.PubKeySecp256k1{},
		Signature: []byte{},
	}
	_ = []byte{}

	tx := authtypes.StdTx{}
	txBytes := ctx.TxBytes()
	err := k.cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
	if err != nil {
		return &sdk.Result{}, sdkerrors.Wrap(types.ErrInstantiateFailed, fmt.Sprintf("Unable to decode transaction from bytes: %s", err.Error()))
	}

	contractInfo := k.GetContractInfo(ctx, contractAddress)
	if contractInfo == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "unknown contract")
	}
	if !authZ.CanModifyContract(contractInfo.Admin, caller) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "can not migrate")
	}

	newCodeInfo := k.GetCodeInfo(ctx, newCodeID)
	if newCodeInfo == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "unknown code")
	}

	store := ctx.KVStore(k.storeKey)
	contractKey := store.Get(types.GetContractEnclaveKey(contractAddress))

	var noDeposit sdk.Coins
	params := types.NewEnv(ctx, caller, noDeposit, contractAddress, contractKey)

	// prepare querier
	querier := QueryHandler{
		Ctx:     ctx,
		Plugins: k.queryPlugins,
	}

	prefixStoreKey := types.GetContractStorePrefixKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	gas := gasForContract(ctx)
	res, gasUsed, err := k.wasmer.Migrate(newCodeInfo.CodeHash, params, msg, &prefixStore, cosmwasmAPI, &querier, gasMeter(ctx), gas)
	consumeGas(ctx, gasUsed)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrMigrationFailed, err.Error())
	}

	// emit all events from this contract itself
	events := types.ParseEvents(res.Log, contractAddress)
	ctx.EventManager().EmitEvents(events)

	historyEntry := contractInfo.AddMigration(ctx, newCodeID, msg)
	k.appendToContractHistory(ctx, contractAddress, historyEntry)
	k.setContractInfo(ctx, contractAddress, contractInfo)

	if err := k.dispatchMessages(ctx, contractAddress, res.Messages); err != nil {
		return nil, sdkerrors.Wrap(err, "dispatch")
	}

	return &sdk.Result{
		Data: res.Data,
	}, nil
}

// UpdateContractAdmin sets the admin value on the ContractInfo. It must be a valid address (use ClearContractAdmin to remove it)
func (k Keeper) UpdateContractAdmin(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, newAdmin sdk.AccAddress) error {
	return k.setContractAdmin(ctx, contractAddress, caller, newAdmin, k.authZPolicy)
}

// ClearContractAdmin sets the admin value on the ContractInfo to nil, to disable further migrations/ updates.
func (k Keeper) ClearContractAdmin(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress) error {
	return k.setContractAdmin(ctx, contractAddress, caller, nil, k.authZPolicy)
}

func (k Keeper) setContractAdmin(ctx sdk.Context, contractAddress, caller, newAdmin sdk.AccAddress, authZ AuthorizationPolicy) error {
	contractInfo := k.GetContractInfo(ctx, contractAddress)
	if contractInfo == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "unknown contract")
	}
	if !authZ.CanModifyContract(contractInfo.Admin, caller) {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "can not modify contract")
	}
	contractInfo.Admin = newAdmin
	k.setContractInfo(ctx, contractAddress, contractInfo)
	return nil
}

func (k Keeper) appendToContractHistory(ctx sdk.Context, contractAddr sdk.AccAddress, newEntries ...types.ContractCodeHistoryEntry) {
	entries := append(k.GetContractHistory(ctx, contractAddr), newEntries...)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ContractHistoryStorePrefix)
	prefixStore.Set(contractAddr, k.cdc.MustMarshalBinaryBare(&entries))
}

func (k Keeper) GetContractHistory(ctx sdk.Context, contractAddr sdk.AccAddress) []types.ContractCodeHistoryEntry {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ContractHistoryStorePrefix)
	var entries []types.ContractCodeHistoryEntry
	bz := prefixStore.Get(contractAddr)
	if bz != nil {
		k.cdc.MustUnmarshalBinaryBare(bz, &entries)
	}
	return entries
}
*/

// QuerySmart queries the smart contract itself.
func (k Keeper) QuerySmart(ctx sdk.Context, contractAddr sdk.AccAddress, req []byte, useDefaultGasLimit bool) ([]byte, error) {
	if useDefaultGasLimit {
		ctx = ctx.WithGasMeter(sdk.NewGasMeter(k.queryGasLimit))
	}
	ctx.GasMeter().ConsumeGas(InstanceCost, "Loading CosmWasm module: query")

	codeInfo, prefixStore, err := k.contractInstance(ctx, contractAddr)
	if err != nil {
		return nil, err
	}

	// prepare querier
	querier := QueryHandler{
		Ctx:     ctx,
		Plugins: k.queryPlugins,
	}

	store := ctx.KVStore(k.storeKey)
	// 0x01 | codeID (uint64) -> ContractInfo
	contractKey := store.Get(types.GetContractEnclaveKey(contractAddr))

	queryResult, gasUsed, qErr := k.wasmer.Query(codeInfo.CodeHash, append(contractKey[:], req[:]...), prefixStore, cosmwasmAPI, querier, gasMeter(ctx), gasForContract(ctx))
	consumeGas(ctx, gasUsed)

	if qErr != nil {
		return nil, sdkerrors.Wrap(types.ErrQueryFailed, qErr.Error())
	}
	return queryResult, nil
}

// We don't use this function since we have an encrypted state. It's here for upstream compatibility
// QueryRaw returns the contract's state for give key. For a `nil` key a empty slice result is returned.
func (k Keeper) QueryRaw(ctx sdk.Context, contractAddress sdk.AccAddress, key []byte) []types.Model {
	result := make([]types.Model, 0)
	if key == nil {
		return result
	}
	prefixStoreKey := types.GetContractStorePrefixKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)

	if val := prefixStore.Get(key); val != nil {
		return append(result, types.Model{
			Key:   key,
			Value: val,
		})
	}
	return result
}

func (k Keeper) contractInstance(ctx sdk.Context, contractAddress sdk.AccAddress) (types.CodeInfo, prefix.Store, error) {
	store := ctx.KVStore(k.storeKey)

	contractBz := store.Get(types.GetContractAddressKey(contractAddress))
	if contractBz == nil {
		return types.CodeInfo{}, prefix.Store{}, sdkerrors.Wrap(types.ErrNotFound, "contract")
	}
	var contract types.ContractInfo
	k.cdc.MustUnmarshalBinaryBare(contractBz, &contract)

	contractInfoBz := store.Get(types.GetCodeKey(contract.CodeID))
	if contractInfoBz == nil {
		return types.CodeInfo{}, prefix.Store{}, sdkerrors.Wrap(types.ErrNotFound, "contract info")
	}
	var codeInfo types.CodeInfo
	k.cdc.MustUnmarshalBinaryBare(contractInfoBz, &codeInfo)
	prefixStoreKey := types.GetContractStorePrefixKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	return codeInfo, prefixStore, nil
}

func (k Keeper) GetContractKey(ctx sdk.Context, contractAddress sdk.AccAddress) []byte {
	store := ctx.KVStore(k.storeKey)

	contractKey := store.Get(types.GetContractEnclaveKey(contractAddress))

	return contractKey
}

func (k Keeper) GetContractAddress(ctx sdk.Context, label string) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)

	contractAddress := store.Get(types.GetContractLabelPrefix(label))

	return contractAddress
}

func (k Keeper) GetContractHash(ctx sdk.Context, contractAddress sdk.AccAddress) []byte {

	info, _ := k.GetContractInfo(ctx, contractAddress)

	hash := k.GetCodeInfo(ctx, info.CodeID).CodeHash

	return hash
}

//GetContractInfo (if you see panic error, try commenting out this)
func (k Keeper) GetContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress) (types.ContractInfo, error) {
	store := ctx.KVStore(k.storeKey)
	var contract types.ContractInfo
	contractBz := store.Get(types.GetContractAddressKey(contractAddress))
	if contractBz == nil {
		return types.ContractInfo{}, sdkerrors.Wrap(types.ErrNotFound, "contract info")
	}
	k.cdc.MustUnmarshalBinaryBare(contractBz, &contract)
	return contract, nil
}

//GetContractResult  (if you see panic error, try commenting out this)
func (k Keeper) GetContractResult(ctx sdk.Context, contractAddress sdk.AccAddress) ([]byte, error) {
	store := ctx.KVStore(k.storeKey)
	//	var contract types.ContractInfo
	contractBz := store.Get(types.GetContractResultKey(contractAddress))
	/*if contractBz == nil {
		return byte, sdkerrors.Wrap(types.ErrNotFound, "contract info")
	}
	k.cdc.MustUnmarshalBinaryBare(contractBz, &contract)*/
	return contractBz, nil
}

//GetContractInfoWithAddress  (if you see panic error, try commenting out this)
func (k Keeper) GetContractInfoWithAddress(ctx sdk.Context, contractAddress sdk.AccAddress) types.ContractInfoWithAddress {
	store := ctx.KVStore(k.storeKey)
	fmt.Printf("Getting info")
	var contract types.ContractInfoWithAddress

	contractBz := store.Get(types.GetContractAddressKey(contractAddress))
	if contractBz == nil {
		return types.ContractInfoWithAddress{} //sdkerrors.Wrap(types.ErrNotFound, "contract")
	}

	fmt.Printf("Unmarshalling..")
	var Info types.ContractInfo
	k.cdc.MustUnmarshalBinaryBare(contractBz, &Info)
	fmt.Printf("Setting")
	contract.ContractInfo = &Info

	/*contractBz := store.Get(types.GetContractAddressKey(contractAddress))
	if contractBz == nil {
		return nil
	}
		err := k.cdc.UnmarshalBinaryBare(contractBz, contract.ContractInfo)
	if err != nil {
		return nil
	}
	*/
	//contract.ContractInfo = k.GetContractInfo(ctx, contractAddress)

	fmt.Printf("info creator is:  %s ", contract.ContractInfo.Creator.String())

	contract.Address = contractAddress
	fmt.Printf("info Address is:  %s ", contract.Address.String())
	return contract
}

func (k Keeper) containsContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetContractAddressKey(contractAddress))
}

func (k Keeper) setContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress, contract *types.ContractInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetContractAddressKey(contractAddress), k.cdc.MustMarshalBinaryBare(contract))
}

// SetContractResult sets the result of the contract
func (k Keeper) SetContractResult(ctx sdk.Context, contractAddress sdk.AccAddress, res *sdk.Result) error {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetContractResultKey(contractAddress), k.cdc.MustMarshalBinaryBare(res))
	return nil
}
func (k Keeper) IterateContractInfo(ctx sdk.Context, cb func(sdk.AccAddress, types.ContractInfo) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ContractKeyPrefix)
	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var contract types.ContractInfo
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &contract)
		// cb returns true to stop early
		if cb(iter.Key(), contract) {
			break
		}
	}
}

func (k Keeper) GetContractState(ctx sdk.Context, contractAddress sdk.AccAddress) sdk.Iterator {
	prefixStoreKey := types.GetContractStorePrefixKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	return prefixStore.Iterator(nil, nil)
}

func (k Keeper) importContractState(ctx sdk.Context, contractAddress sdk.AccAddress, models []types.Model) error {
	prefixStoreKey := types.GetContractStorePrefixKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	for _, model := range models {
		if model.Value == nil {
			model.Value = []byte{}
		}
		if prefixStore.Has(model.Key) {
			return sdkerrors.Wrapf(types.ErrDuplicate, "duplicate key: %x", model.Key)
		}
		prefixStore.Set(model.Key, model.Value)
	}
	return nil
}

func (k Keeper) GetCodeInfo(ctx sdk.Context, codeID uint64) *types.CodeInfo {
	store := ctx.KVStore(k.storeKey)
	var codeInfo types.CodeInfo
	codeInfoBz := store.Get(types.GetCodeKey(codeID))
	if codeInfoBz == nil {
		return nil
	}
	k.cdc.MustUnmarshalBinaryBare(codeInfoBz, &codeInfo)
	return &codeInfo
}

func (k Keeper) GetCodeHash(ctx sdk.Context, codeID uint64) (codeHash []byte) {
	store := ctx.KVStore(k.storeKey)
	var codeInfo types.CodeInfo
	codeInfoBz := store.Get(types.GetCodeKey(codeID))
	if codeInfoBz == nil {
		return nil
	}

	k.cdc.MustUnmarshalBinaryBare(codeInfoBz, &codeInfo)
	fmt.Printf("GetCodeHash codeInfo: %X\n", codeInfo)
	fmt.Printf("GetCodeHash codeInfo: %X\n", codeInfo.CodeHash)
	return codeInfo.CodeHash
}

func (k Keeper) containsCodeInfo(ctx sdk.Context, codeID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetCodeKey(codeID))
}

func (k Keeper) IterateCodeInfos(ctx sdk.Context, cb func(uint64, types.CodeInfo) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.CodeKeyPrefix)
	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var c types.CodeInfo
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &c)
		// cb returns true to stop early
		if cb(binary.BigEndian.Uint64(iter.Key()), c) {
			return
		}
	}
}

func (k Keeper) GetByteCode(ctx sdk.Context, codeID uint64) ([]byte, error) {
	store := ctx.KVStore(k.storeKey)
	var codeInfo types.CodeInfo
	codeInfoBz := store.Get(types.GetCodeKey(codeID))
	if codeInfoBz == nil {
		return nil, nil
	}
	k.cdc.MustUnmarshalBinaryBare(codeInfoBz, &codeInfo)
	return k.wasmer.GetCode(codeInfo.CodeHash)
}

func (k Keeper) dispatchMessages(ctx sdk.Context, contractAddr sdk.AccAddress, msgs []wasmTypes.CosmosMsg) error {
	for _, msg := range msgs {
		if err := k.messenger.Dispatch(ctx, contractAddr, msg); err != nil {
			return err
		}
	}
	return nil
}

func gasForContract(ctx sdk.Context) uint64 {
	meter := ctx.GasMeter()
	remaining := (meter.Limit() - meter.GasConsumed()) * GasMultiplier
	if remaining > MaxGas {
		return MaxGas
	}
	return remaining
}

func consumeGas(ctx sdk.Context, gas uint64) {
	consumed := (gas / GasMultiplier) + 1
	ctx.GasMeter().ConsumeGas(consumed, "wasm contract")
	// throw OutOfGas error if we ran out (got exactly to zero due to better limit enforcing)
	if ctx.GasMeter().IsOutOfGas() {
		panic(sdk.ErrorOutOfGas{})
	}
}

// generates a contract address from codeID + instanceID
func (k Keeper) generateContractAddress(ctx sdk.Context, codeID uint64) sdk.AccAddress {
	instanceID := k.autoIncrementID(ctx, types.KeyLastInstanceID)
	return contractAddress(codeID, instanceID)
}

func contractAddress(codeID, instanceID uint64) sdk.AccAddress {
	// NOTE: It is possible to get a duplicate address if either codeID or instanceID
	// overflow 32 bits. This is highly improbable, but something that could be refactored.
	contractID := codeID<<32 + instanceID
	return addrFromUint64(contractID)

}

func (k Keeper) GetNextCodeID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyLastCodeID)
	id := uint64(1)
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}
	return id
}

func (k Keeper) autoIncrementID(ctx sdk.Context, lastIDKey []byte) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(lastIDKey)
	id := uint64(1)
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}
	bz = sdk.Uint64ToBigEndian(id + 1)
	store.Set(lastIDKey, bz)
	return id
}

// peekAutoIncrementID reads the current value without incrementing it.
func (k Keeper) peekAutoIncrementID(ctx sdk.Context, lastIDKey []byte) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(lastIDKey)
	id := uint64(1)
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}
	return id
}

func (k Keeper) importAutoIncrementID(ctx sdk.Context, lastIDKey []byte, val uint64) error {
	store := ctx.KVStore(k.storeKey)
	if store.Has(lastIDKey) {
		return sdkerrors.Wrapf(types.ErrDuplicate, "autoincrement id: %s", string(lastIDKey))
	}
	bz := sdk.Uint64ToBigEndian(val)
	store.Set(lastIDKey, bz)
	return nil
}

func (k Keeper) importContract(ctx sdk.Context, contractAddr sdk.AccAddress, c *types.ContractInfo, state []types.Model) error {
	if !k.containsCodeInfo(ctx, c.CodeID) {
		return sdkerrors.Wrapf(types.ErrNotFound, "code id: %d", c.CodeID)
	}
	if k.containsContractInfo(ctx, contractAddr) {
		return sdkerrors.Wrapf(types.ErrDuplicate, "contract: %s", contractAddr)
	}

	// historyEntry := c.ResetFromGenesis(ctx)
	// k.appendToContractHistory(ctx, contractAddr, historyEntry)
	k.setContractInfo(ctx, contractAddr, c)
	return k.importContractState(ctx, contractAddr, state)
}

func addrFromUint64(id uint64) sdk.AccAddress {
	addr := make([]byte, 20)
	addr[0] = 'C'
	binary.PutUvarint(addr[1:], id)
	return sdk.AccAddress(crypto.AddressHash(addr))
}

// MultipliedGasMeter wraps the GasMeter from context and multiplies all reads by out defined multiplier
type MultipiedGasMeter struct {
	originalMeter sdk.GasMeter
}

var _ wasm.GasMeter = MultipiedGasMeter{}

func (m MultipiedGasMeter) GasConsumed() sdk.Gas {
	return m.originalMeter.GasConsumed() * GasMultiplier
}

func gasMeter(ctx sdk.Context) MultipiedGasMeter {
	return MultipiedGasMeter{
		originalMeter: ctx.GasMeter(),
	}
}
