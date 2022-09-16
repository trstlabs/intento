package keeper

import (
	"bytes"
	"strconv"

	//"encoding/json"
	"fmt"
	"time"

	codedctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	sdktxsigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	api "github.com/trstlabs/trst/go-cosmwasm/api"
	wasmTypes "github.com/trstlabs/trst/go-cosmwasm/types"
	"github.com/trstlabs/trst/x/compute/internal/types"
)

// Create uploads and compiles a WASM contract, returning a short identifier for the contract
func (k Keeper) Create(ctx sdk.Context, creator sdk.AccAddress, wasmCode []byte, source string, builder string, duration time.Duration, interval time.Duration, title string, description string) (codeID uint64, err error) {

	wasmCode, err = uncompress(wasmCode)
	if err != nil {
		return 0, sdkerrors.Wrap(types.ErrCreateFailed, err.Error())
	}
	ctx.GasMeter().ConsumeGas(types.CompileCost*uint64(len(wasmCode)), "Compiling WASM Bytecode")

	codeHash, err := k.wasmer.Create(wasmCode)
	if err != nil {
		return 0, sdkerrors.Wrap(types.ErrCreateFailed, err.Error())
	}

	store := ctx.KVStore(k.storeKey)
	codeID = k.autoIncrementID(ctx, types.KeyLastCodeID)
	/*
		if instantiateAccess == nil {
			defaultAccessConfig := k.getInstantiateAccessConfig(ctx).With(creator)
			instantiateAccess = &defaultAccessConfig
		}
	*/
	codeInfo := types.NewCodeInfo(codeHash, creator, source, builder, duration, interval, title, description)
	// 0x01 | codeID (uint64) -> ContractInfo
	store.Set(types.GetCodeKey(codeID), k.cdc.MustMarshal(&codeInfo))

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
	store.Set(key, k.cdc.MustMarshal(&codeInfo))
	return nil
}

// Instantiate creates an instance of a WASM contract
func (k Keeper) Instantiate(ctx sdk.Context, codeID uint64, creator /* , admin */ sdk.AccAddress, msg []byte, autoMsg []byte, id string, funds sdk.Coins, callbackSig []byte, duration time.Duration, interval time.Duration, startTime time.Time) (sdk.AccAddress, []byte, error) {

	ctx.GasMeter().ConsumeGas(types.InstanceCost, "Loading CosmWasm module: init")

	signBytes := []byte{}
	signMode := sdktxsigning.SignMode_SIGN_MODE_UNSPECIFIED
	modeInfoBytes := []byte{}
	pkBytes := []byte{}
	signerSig := []byte{}
	var err error

	// If no callback signature - we should send the actual msg sender sign bytes and signature
	if callbackSig == nil {
		signBytes, signMode, modeInfoBytes, pkBytes, signerSig, err = k.GetSignerInfo(ctx, creator)
		if err != nil {
			return nil, nil, err
		}
		fmt.Printf("Init by account \n")
	}

	verificationInfo := types.NewVerificationInfo(signBytes, signMode, modeInfoBytes, pkBytes, signerSig, callbackSig)

	store := ctx.KVStore(k.storeKey)
	existingAddress := store.Get(types.GetContractIdPrefix(id))

	if existingAddress != nil {
		return nil, nil, sdkerrors.Wrap(types.ErrAccountExists, id)
	}

	contractAddress := k.generateContractAddress(ctx, codeID)
	existingAcct := k.accountKeeper.GetAccount(ctx, contractAddress)
	if existingAcct != nil {
		return nil, nil, sdkerrors.Wrap(types.ErrAccountExists, existingAcct.GetAddress().String())
	}

	// deposit initial contract funds
	if !funds.IsZero() {
		if k.bankKeeper.BlockedAddr(creator) {
			return nil, nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "blocked address can not be used")
		}
		sdkerr := k.bankKeeper.SendCoins(ctx, creator, contractAddress, funds)
		if sdkerr != nil {
			return nil, nil, sdkerr
		}
	} else {
		// create an empty account (so we don't have issues later)
		contractAccount := k.accountKeeper.NewAccountWithAddress(ctx, contractAddress)
		k.accountKeeper.SetAccount(ctx, contractAccount)
	}

	// get contact info
	bz := store.Get(types.GetCodeKey(codeID))
	if bz == nil {
		return nil, nil, sdkerrors.Wrap(types.ErrNotFound, "code")
	}
	var codeInfo types.CodeInfo
	k.cdc.MustUnmarshal(bz, &codeInfo)

	//fmt.Printf("funds %v \n", funds)
	// prepare env for contract instantiate call
	env := types.NewEnv(ctx, creator, funds, contractAddress, nil)

	// create prefixed data store
	// 0x03 | contractAddress (sdk.AccAddress)
	prefixStoreKey := types.GetContractStorePrefixKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)

	// prepare querier
	querier := QueryHandler{
		Ctx:     ctx,
		Plugins: k.queryPlugins,
	}

	autoMsgToSend := make([]byte, 256)
	if autoMsg != nil {
		autoMsgToSend = autoMsg
	}
	// instantiate wasm contract
	gas := gasForContract(ctx)
	res, key, callbackSig, gasUsed, errData, err := k.wasmer.Instantiate(codeInfo.CodeHash, env, msg, autoMsgToSend, prefixStore, cosmwasmAPI, querier, ctx.GasMeter(), gas, verificationInfo, contractAddress)
	consumeGas(ctx, gasUsed)
	if err != nil {
		return nil, errData, sdkerrors.Wrap(types.ErrInstantiateFailed, err.Error())
	}

	if err != nil {
		return contractAddress, res.Data, sdkerrors.Wrap(types.ErrInstantiateFailed, err.Error())
	}

	//fmt.Printf("Attributes: %v \n", res.Attributes)
	// emit all events from this contract itself
	//events := types.ParseEvents(res.Attributes, contractAddress)
	//ctx.EventManager().EmitEvents(events)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeInstantiate,
		sdk.NewAttribute(types.AttributeKeyContractAddr, contractAddress.String()),
		sdk.NewAttribute(types.AttributeKeyCodeID, strconv.FormatUint(codeID, 10)),
	))

	// persist instance
	createdAt := types.NewAbsoluteTxPosition(ctx)

	var endTime time.Time
	var execTime time.Time
	//if default duration, this duration is executed. Instantiators can have a custom duration and interval (validated in handler).
	if codeInfo.DefaultDuration != 0 {
		endTime = startTime.Add(codeInfo.DefaultDuration)
		if duration != 0 {
			endTime = startTime.Add(duration)
		}
		if codeInfo.DefaultInterval == 0 {
			k.InsertContractQueue(ctx, contractAddress.String(), endTime)
			execTime = endTime
		} else {
			endTime = startTime.Add(codeInfo.DefaultDuration)
			execTime = startTime.Add(codeInfo.DefaultInterval)
			if interval != 0 {
				execTime = startTime.Add(interval)
			} else {
				interval = codeInfo.DefaultInterval
			}
			k.InsertContractQueue(ctx, contractAddress.String(), execTime)
		}
	}

	contractInfo := types.NewContractInfo(codeID, creator /* admin, */, id, createdAt, startTime, execTime, endTime, interval, autoMsg, callbackSig)
	// check for IBC flag
	report, err := k.wasmer.AnalyzeCode(codeInfo.CodeHash)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(types.ErrInstantiateFailed, err.Error())
	}
	if report.HasIBCEntryPoints {
		// register IBC port
		ibcPort, err := k.ensureIbcPort(ctx, contractAddress)
		if err != nil {
			return nil, nil, err
		}
		contractInfo.IBCPortID = ibcPort
	}

	store.Set(types.GetContractAddressKey(contractAddress), k.cdc.MustMarshal(&contractInfo))

	codeInfo.Instances = codeInfo.Instances + 1
	store.Set(types.GetCodeKey(codeID), k.cdc.MustMarshal(&codeInfo))
	store.Set(types.GetContractEnclaveKey(contractAddress), key)
	store.Set(types.GetContractIdPrefix(id), contractAddress)

	err = k.SetContractPublicState(ctx, contractAddress, res.Attributes)
	if err != nil {
		return nil, nil, err
	}

	data, err := k.handleContractResponse(ctx, contractAddress, contractInfo.IBCPortID, res, res.Messages, res.Events, res.Data, msg, verificationInfo)
	if err != nil {
		return nil, nil, err
	}

	//both airdrop actions are performed through callbacksig
	if callbackSig != nil {
		k.SetAirdropAction(ctx, res.Attributes)
	}
	return contractAddress, data, nil
}

// Execute executes the contract instance
func (k Keeper) Execute(ctx sdk.Context, contractAddress sdk.AccAddress, caller sdk.AccAddress, msg []byte, funds sdk.Coins, callbackSig []byte) (*sdk.Result, error) {

	ctx.GasMeter().ConsumeGas(types.InstanceCost, "Loading compute module: execute")

	signBytes := []byte{}
	signMode := sdktxsigning.SignMode_SIGN_MODE_UNSPECIFIED
	modeInfoBytes := []byte{}
	pkBytes := []byte{}
	signerSig := []byte{}
	var err error

	// If no callback signature - we should send the actual msg sender sign bytes and signature
	if callbackSig == nil {
		signBytes, signMode, modeInfoBytes, pkBytes, signerSig, err = k.GetSignerInfo(ctx, caller)
		if err != nil {

			return nil, err
		}
		fmt.Printf("Execute by account \n")
	}

	verificationInfo := types.NewVerificationInfo(signBytes, signMode, modeInfoBytes, pkBytes, signerSig, callbackSig)

	contractInfo, codeInfo, prefixStore, err := k.contractInstance(ctx, contractAddress)
	if err != nil {
		return nil, err
	}

	store := ctx.KVStore(k.storeKey)

	// add funds
	if !funds.IsZero() {
		if k.bankKeeper.BlockedAddr(caller) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "blocked address can not be used")
		}

		sdkerr := k.bankKeeper.SendCoins(ctx, caller, contractAddress, funds)
		if sdkerr != nil {
			return nil, sdkerr
		}
	}

	contractKey := store.Get(types.GetContractEnclaveKey(contractAddress))
	if contractKey == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "contract key not found")
	}
	params := types.NewEnv(ctx, caller, funds, contractAddress, contractKey)

	// prepare querier
	querier := QueryHandler{
		Ctx:     ctx,
		Plugins: k.queryPlugins,
	}

	gas := gasForContract(ctx)
	res, gasUsed, errData, err := k.wasmer.Execute(codeInfo.CodeHash, params, msg, prefixStore, cosmwasmAPI, querier, gasMeter(ctx), gas, verificationInfo, wasmTypes.HandleTypeExecute)
	consumeGas(ctx, gasUsed)

	if err != nil {
		return &sdk.Result{Data: errData}, sdkerrors.Wrap(types.ErrExecuteFailed, err.Error())
	}

	err = k.SetContractPublicState(ctx, contractAddress, res.Attributes)
	if err != nil {
		return nil, err
	}

	/*ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeExecute,
		sdk.NewAttribute(types.AttributeKeyContractAddr, contractAddress.String()),
	))*/

	data, err := k.handleContractResponse(ctx, contractAddress, contractInfo.IBCPortID, res, res.Messages, res.Events, res.Data, msg, verificationInfo)
	if err != nil {
		return nil, err
	}

	//return &sdk.Result{Data: res.Data,Log: res.Log[1].Value}, nil //used for item module compatibilily

	return &sdk.Result{Data: data}, nil

}

func (k Keeper) GetSignerInfo(ctx sdk.Context, signer sdk.AccAddress) ([]byte, sdktxsigning.SignMode, []byte, []byte, []byte, error) {
	tx := sdktx.Tx{}
	err := k.cdc.Unmarshal(ctx.TxBytes(), &tx)
	if err != nil {
		return nil, 0, nil, nil, nil, sdkerrors.Wrap(types.ErrSigFailed, fmt.Sprintf("Unable to decode transaction from bytes: %s", err.Error()))
	}

	// for MsgInstantiateContract, there is only one signer which is msg.Sender
	// (https://github.com/trstlabs/trst/blob/d7813792fa07b93a10f0885eaa4c5e0a0a698854/x/compute/internal/types/msg.go#L192-L194)
	signerAcc, err := ante.GetSignerAcc(ctx, k.accountKeeper, signer)
	if err != nil {

		return nil, 0, nil, nil, nil, sdkerrors.Wrap(types.ErrSigFailed, fmt.Sprintf("Unable to retrieve account by address: %s", err.Error()))
	}

	txConfig := authtx.NewTxConfig(k.cdc.(*codec.ProtoCodec), authtx.DefaultSignModes)
	modeHandler := txConfig.SignModeHandler()
	signingData := authsigning.SignerData{
		ChainID:       ctx.ChainID(),
		AccountNumber: signerAcc.GetAccountNumber(),
		Sequence:      signerAcc.GetSequence() - 1,
	}

	protobufTx := authtx.WrapTx(&tx).GetTx()

	pubKeys, err := protobufTx.GetPubKeys()
	if err != nil {
		return nil, 0, nil, nil, nil, sdkerrors.Wrap(types.ErrSigFailed, fmt.Sprintf("Unable to get public keys for instantiate: %s", err.Error()))
	}

	pkIndex := -1
	var _signers [][]byte // This is just used for the error message below
	for index, pubKey := range pubKeys {
		thisSigner := pubKey.Address().Bytes()
		_signers = append(_signers, thisSigner)
		if bytes.Equal(thisSigner, signer.Bytes()) {
			pkIndex = index
		}
	}
	if pkIndex == -1 {
		return nil, 0, nil, nil, nil, sdkerrors.Wrap(types.ErrSigFailed, fmt.Sprintf("Message sender: %s is not found in the tx signer set: %v, callback signature not provided", signer.String(), _signers))
	}

	signatures, _ := protobufTx.GetSignaturesV2()
	var signMode sdktxsigning.SignMode
	switch signData := signatures[pkIndex].Data.(type) {
	case *sdktxsigning.SingleSignatureData:
		signMode = signData.SignMode
	case *sdktxsigning.MultiSignatureData:
		signMode = sdktxsigning.SignMode_SIGN_MODE_LEGACY_AMINO_JSON
	}

	signBytes, err := modeHandler.GetSignBytes(signMode, signingData, protobufTx)
	if err != nil {
		return nil, 0, nil, nil, nil, sdkerrors.Wrap(types.ErrSigFailed, fmt.Sprintf("Unable to recreate sign bytes for the tx: %s", err.Error()))
	}

	modeInfoBytes, err := sdktxsigning.SignatureDataToProto(signatures[pkIndex].Data).Marshal()
	if err != nil {
		return nil, 0, nil, nil, nil, sdkerrors.Wrap(types.ErrSigFailed, "couldn't marshal mode info")
	}

	var pkBytes []byte
	pubKey := pubKeys[pkIndex]
	anyPubKey, err := codedctypes.NewAnyWithValue(pubKey)
	if err != nil {
		return nil, 0, nil, nil, nil, sdkerrors.Wrap(types.ErrSigFailed, "couldn't turn public key into Any")
	}

	pkBytes, err = k.cdc.Marshal(anyPubKey)
	if err != nil {
		return nil, 0, nil, nil, nil, sdkerrors.Wrap(types.ErrSigFailed, "couldn't marshal public key")
	}

	return signBytes, signMode, modeInfoBytes, pkBytes, tx.Signatures[pkIndex], nil
}

// Delete deletes data related to the contract instance after the contract ends
func (k Keeper) Delete(ctx sdk.Context, contractAddress sdk.AccAddress) error {

	_, _, prefixStore, err := k.contractInstance(ctx, contractAddress)
	if err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)

	contractBz := store.Get(types.GetContractAddressKey(contractAddress))
	if contractBz == nil {
		return sdkerrors.Wrap(types.ErrNotFound, "contract")
	}
	var contract types.ContractInfo
	k.cdc.MustUnmarshal(contractBz, &contract)

	//delete wasm contract instance data
	prefixStoreKey := types.GetContractStorePrefixKey(contractAddress)
	prefixStore.Delete(prefixStoreKey)

	//store.Delete(types.GetContractEnclaveKey(contractAddress))
	store.Delete(types.GetContractIdPrefix(contract.ContractId))
	// contract info
	//store.Delete(types.GetContractAddressKey(contractAddress))

	return nil
}

// QueryPrivate queries the smart contract itself.
func (k Keeper) QueryPrivate(ctx sdk.Context, contractAddr sdk.AccAddress, req []byte, useDefaultGasLimit bool) ([]byte, error) {
	return k.queryPrivateContractImpl(ctx, contractAddr, req, useDefaultGasLimit, false)
}

// queryPrivateRecursive queries the smart contract itself. This should only be called when running inside another query recursively.
func (k Keeper) queryPrivateRecursive(ctx sdk.Context, contractAddr sdk.AccAddress, req []byte, useDefaultGasLimit bool) ([]byte, error) {
	return k.queryPrivateContractImpl(ctx, contractAddr, req, useDefaultGasLimit, true)
}

// queryPrivateContractImpl queries the contract itself. This should only be called when running inside another query recursively.
func (k Keeper) queryPrivateContractImpl(ctx sdk.Context, contractAddr sdk.AccAddress, req []byte, useDefaultGasLimit bool, recursive bool) ([]byte, error) {
	defer telemetry.MeasureSince(time.Now(), "compute", "keeper", "query")

	if useDefaultGasLimit {
		ctx = ctx.WithGasMeter(sdk.NewGasMeter(k.queryGasLimit))
	}
	ctx.GasMeter().ConsumeGas(types.InstanceCost, "Loading CosmWasm module: query")

	_, codeInfo, prefixStore, err := k.contractInstance(ctx, contractAddr)
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

	params := types.NewEnv(
		ctx,
		sdk.AccAddress{}, /* empty because it's unused in queries */
		[]sdk.Coin{},     /* empty because it's unused in queries */
		contractAddr,
		contractKey,
	)
	params.Recursive = recursive
	queryResult, gasUsed, qErr := k.wasmer.Query(codeInfo.CodeHash, params, req, prefixStore, cosmwasmAPI, querier, gasMeter(ctx), gasForContract(ctx))
	consumeGas(ctx, gasUsed)
	//fmt.Printf("Query queryResult %+v \n", queryResult)
	if qErr != nil {

		return nil, sdkerrors.Wrap(types.ErrQueryFailed, qErr.Error())
	}
	return queryResult, nil
}

//QueryPublic queries the public contract state
func (k Keeper) QueryPublic(ctx sdk.Context, contractAddress sdk.AccAddress, key string) []byte {
	value := k.GetContractPublicStateValue(ctx, contractAddress, key)
	return value
}

//QueryPublicForAddr queries the public contract state for a given address
func (k Keeper) QueryPublicForAddr(ctx sdk.Context, contractAddress sdk.AccAddress, accountAddress sdk.AccAddress, key string) []byte {
	value := k.GetContractPublicStateValueForAddr(ctx, contractAddress, accountAddress, key)
	return value
}

func (k Keeper) contractInstance(ctx sdk.Context, contractAddress sdk.AccAddress) (types.ContractInfo, types.CodeInfo, prefix.Store, error) {
	store := ctx.KVStore(k.storeKey)
	contractBz := store.Get(types.GetContractAddressKey(contractAddress))
	if contractBz == nil {
		return types.ContractInfo{}, types.CodeInfo{}, prefix.Store{}, sdkerrors.Wrap(types.ErrNotFound, "contract")
	}
	var contractInfo types.ContractInfo
	k.cdc.MustUnmarshal(contractBz, &contractInfo)

	codeInfoBz := store.Get(types.GetCodeKey(contractInfo.CodeID))
	if codeInfoBz == nil {
		return types.ContractInfo{}, types.CodeInfo{}, prefix.Store{}, sdkerrors.Wrap(types.ErrNotFound, "contract info")
	}
	var codeInfo types.CodeInfo
	k.cdc.MustUnmarshal(codeInfoBz, &codeInfo)
	prefixStoreKey := types.GetContractStorePrefixKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	return contractInfo, codeInfo, prefixStore, nil
}

// CreateCommunityPoolCallbackSig creates a callback sig which can be used to execute a specific message for a specific code for the community pool.
// When callback signature is made, any node can 'run' the message at any time on the community pool's behalf, therefore, anyone can create outputs for the distribution module account.
// By hardcoding the distribution module address in the enclave, we can use this for contract instantiation and execution over governance.
func (k Keeper) CreateCommunityPoolCallbackSig(ctx sdk.Context, msg []byte, codeID uint64, funds sdk.Coins) (callbackSig []byte, encryptedMessage []byte, err error) {
	// get contact info
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCodeKey(codeID))
	if bz == nil {
		return nil, nil, sdkerrors.Wrap(types.ErrNotFound, "code")
	}
	var codeInfo types.CodeInfo
	k.cdc.MustUnmarshal(bz, &codeInfo)
	msgInfo := types.NewMsgInfo(codeInfo.CodeHash, funds)
	//fmt.Printf("code hash: \t %v \n", msgInfo.CodeHash)
	callbackSig, encryptedMessage, err = api.GetCallbackSig(msg, msgInfo)
	if err != nil {
		return nil, nil, err
	}

	//fmt.Printf("callbackSig: \t %v \n", callbackSig)

	return callbackSig, encryptedMessage, nil
}

// DiscardAutoMsg cancels the automessage for a given contract on request of the instantiator
func (k Keeper) DiscardAutoMsg(ctx sdk.Context, info types.ContractInfo, contractAddress sdk.AccAddress, sender sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)

	// have a sufficient runway before discarding the contract (can be adjusted later on)
	min, err := time.ParseDuration("1h")
	if err != nil {
		return err
	}
	if info.EndTime.Before(ctx.BlockHeader().Time.Add(min)) {
		return sdkerrors.Wrap(types.ErrNotFound, "contract info")
	}
	k.RemoveFromContractQueue(ctx, contractAddress.String(), info.EndTime)
	info.AutoMsg = nil
	info.EndTime = ctx.BlockHeader().Time
	store.Set(types.GetContractAddressKey(contractAddress), k.cdc.MustMarshal(&info))
	//info: \t %v \n", info)

	return nil
}
