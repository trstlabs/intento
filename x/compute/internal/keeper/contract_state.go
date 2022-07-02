package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"

	//"encoding/json"
	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	wasmTypes "github.com/trstlabs/trst/go-cosmwasm/types"
	"github.com/trstlabs/trst/x/compute/internal/types"
)

func (k Keeper) GetContractKey(ctx sdk.Context, contractAddress sdk.AccAddress) []byte {
	store := ctx.KVStore(k.storeKey)

	contractKey := store.Get(types.GetContractEnclaveKey(contractAddress))

	return contractKey
}

func (k Keeper) GetContractAddress(ctx sdk.Context, id string) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)

	contractAddress := store.Get(types.GetContractIdPrefix(id))

	return contractAddress
}

func (k Keeper) GetContractHash(ctx sdk.Context, contractAddress sdk.AccAddress) []byte {

	info := k.GetContractInfo(ctx, contractAddress)

	hash := k.GetCodeInfo(ctx, info.CodeID).CodeHash

	return hash
}

//GetContractInfo
func (k Keeper) GetContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress) *types.ContractInfo {
	store := ctx.KVStore(k.storeKey)
	var contract types.ContractInfo
	contractBz := store.Get(types.GetContractAddressKey(contractAddress))
	if contractBz == nil {
		return nil
	}
	k.cdc.MustUnmarshal(contractBz, &contract)
	return &contract
}

//GetContractInfoWithAddress
func (k Keeper) GetContractInfoWithAddress(ctx sdk.Context, contractAddress sdk.AccAddress) types.ContractInfoWithAddress {
	store := ctx.KVStore(k.storeKey)

	var contract types.ContractInfoWithAddress

	contractBz := store.Get(types.GetContractAddressKey(contractAddress))
	if contractBz == nil {
		return types.ContractInfoWithAddress{} //sdkerrors.Wrap(types.ErrNotFound, "contract")
	}

	var info types.ContractInfo
	k.cdc.MustUnmarshal(contractBz, &info)

	contract.ContractInfo = &info

	contract.Address = contractAddress

	return contract
}

func (k Keeper) containsContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetContractAddressKey(contractAddress))
}

func (k Keeper) setContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress, contract *types.ContractInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetContractAddressKey(contractAddress), k.cdc.MustMarshal(contract))
}

// SetContractPublicState sets the result of the contract from wasm attributes, it overrides existing keys
func (k Keeper) SetContractPublicState(ctx sdk.Context, contrAddr sdk.AccAddress, caller sdk.AccAddress, result []wasmTypes.LogAttribute) {
	prefixStoreKey := types.GetContractPubDbKey(contrAddr)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	prefixAccStoreKey := types.GetContractAccPubDbKey(contrAddr, caller)
	prefixAccStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixAccStoreKey)
	for _, attr := range result {
		if attr.PubDb {
			if attr.AccPubDb {
				fmt.Printf(" caller: %s \n", caller.String())
				fmt.Printf("set acc state key %s \n,", string(attr.Key))
				prefixAccStore.Set([]byte(attr.Key), attr.Value)
				break
			}
			fmt.Printf(" contrAddr: %s \n", contrAddr.String())
			fmt.Printf("set state key %s \n,", string(attr.Key))
			prefixStore.Set([]byte(attr.Key), attr.Value)
		}
	}
}

/*
// SetContractPublicState sets the result of the contract from wasm attributes, it overrides existing keys
func (k Keeper) SetContractPublicState(ctx sdk.Context, contrAddr sdk.AccAddress, accAddr sdk.AccAddress, result []wasmTypes.LogAttribute) {
	for _, attr := range result {
		if attr.PubDb {
			if attr.AccPubDb {
				k.SetContractPublicStateAccVal(ctx, attr, contrAddr, accAddr)
				break
			}
			k.SetContractPublicStateVal(ctx, attr, contrAddr)
		}
	}
}
*/
/*
// TEST SetContractPublicState sets the result of the contract from wasm attributes, it overrides existing keys
func (k Keeper) SetContractPublicState(ctx sdk.Context, contrAddr sdk.AccAddress, result []wasmTypes.LogAttribute) {
	prefixStoreKey := types.GetContractPubDbKey(contrAddr)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	for i, attr := range result {
		fmt.Printf(" contrAddr: %s \n", contrAddr.String())
		fmt.Printf(" key: %s \n", attr.Key)
		fmt.Printf(" val: %s \n", attr.Value)
		fmt.Printf(" pub: %v \n", attr.PubDb)
		fmt.Printf(" acc: %v \n", attr.AccPubDb)
		fmt.Printf(" enc: %v \n", attr.Encrypted)
		if attr.PubDb && !attr.AccPubDb {
			fmt.Printf("set key: %s \n", attr.Key)
			fmt.Printf("set val: %s \n", attr.Value)
			//up until here works for this..
			prefixStore.Set([]byte(attr.Key), attr.Value)
			prefixStore.Set([]byte("pub key test"), []byte("pub key test"))
			prefixStore.Set([]byte("pub key test2"), []byte("pub key test2"))
		} else if attr.AccPubDb {
			fmt.Printf("set acc key: %s \n", attr.Key)
			fmt.Printf("set acc val: %s \n", attr.Value)
			//up until here works for this..
			prefixStore.Set([]byte("acc key test"), []byte("acc key test"))
			prefixStore.Set([]byte("acc key test"), []byte("acc key test"))
		} else if !attr.PubDb {
			fmt.Printf("set encrypted2 key: %s \n", attr.Key)
			fmt.Printf("set encrypted2 val: %s \n", attr.Value)
			prefixStore.Set([]byte("encrypted key2 test"), []byte("encrypted key2 test"))
			prefixStore.Set([]byte("encrypted key2 test2"), []byte("encrypted key2 test"))
			prefixStore.Set([]byte("encrypted key2 test3"), []byte("encrypted key2 test"))
		} else {
			fmt.Printf("else: %s \n", attr.Key)
		}
		fmt.Printf("done for val: %s \n", attr.Key)
		fmt.Printf("done for: %d \n", i)
		prefixStore.Set([]byte("done for key"), attr.Value)
	}
}

//SetContractPublicStateVal sets value for public contract state
func (k Keeper) SetContractPublicStateVal(ctx sdk.Context, kv wasmTypes.LogAttribute, contrAddr sdk.AccAddress) {
	prefixStoreKey := types.GetContractPubDbKey(contrAddr)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	prefixStore.Set([]byte(kv.Key), kv.Value)
	fmt.Printf("set state key %s \n", kv.Key)
}

//SetContractPublicStateVal sets value for an account specific public contract state
func (k Keeper) SetContractPublicStateAccVal(ctx sdk.Context, kv wasmTypes.LogAttribute, contrAddr sdk.AccAddress, accAddr sdk.AccAddress) {
	prefixAccStoreKey := types.GetContractAccPubDbKey(contrAddr, accAddr)
	prefixAccStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixAccStoreKey)
	prefixAccStore.Set([]byte(kv.Key), kv.Value)
}


// SetContractPublicState sets the result of the contract from wasm attributes, it overrides existing keys
func (k Keeper) SetContractPublicState(ctx sdk.Context, contrAddr sdk.AccAddress, accAddr sdk.AccAddress, result []wasmTypes.LogAttribute) {
	var toSet []*types.KeyPair
	var toSetAcc []*types.KeyPair

	store := ctx.KVStore(k.storeKey)
	prefixStoreKey := types.GetContractPubDbKey(contrAddr)
	pS := store.Get(prefixStoreKey)
	err := proto.Unmarshal(&pS, toSet)
	if err != nil {
		panic(err)
	}

	prefixStoreKeyAcc := types.GetContractAccPubDbKey(contrAddr, accAddr)
	pSAcc := store.Get(prefixStoreKeyAcc)
	err = proto.Unmarshal(&pS, toSetAcc)
	if err != nil {
		panic(err)
	}
	for _, attr := range result {
		if attr.PubDb {
			if attr.AccPubDb {
				toSetAcc = append(toSetAcc, &types.KeyPair{Key: attr.Key, Value: string(attr.Value)})
				break
			}
			toSet = append(toSet, &types.KeyPair{Key: attr.Key, Value: string(attr.Value)})
		}
	}

	bz, err := proto.Marshal(&toSet)
	if err != nil {
		panic(err)
	}
	//prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	store.Set(prefixStoreKey, bz)

}
*/
/*
// WORKS OLD SetContractPublicState sets the result of the contract from wasm attributes, it overrides existing keys
func (k Keeper) SetContractPublicState(ctx sdk.Context, contractAddress sdk.AccAddress, result []wasmTypes.LogAttribute) {
	prefixStoreKey := types.GetContractPubDbKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	for _, attr := range result {
		if !attr.Encrypted {
			prefixStore.Set([]byte(attr.Key), []byte(attr.Value))
		}
	}
}*/

//GetContractPublicState gets the public contract state
func (k Keeper) GetContractPublicState(ctx sdk.Context, contractAddress sdk.AccAddress) []*types.KeyPair {
	prefixStoreKey := types.GetContractPubDbKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	iter := prefixStore.Iterator(nil, nil)

	defer iter.Close()

	var pKeyPair []*types.KeyPair
	for ; iter.Valid(); iter.Next() {
		fmt.Printf("get key %s \n", []byte(iter.Key()))
		pKeyPair = append(pKeyPair, &types.KeyPair{Key: string(iter.Key()), Value: string(iter.Value())})
	}
	return pKeyPair
}

//GetContractPublicStateForAccount
func (k Keeper) GetContractPublicStateForAccount(ctx sdk.Context, contractAddress sdk.AccAddress, account sdk.AccAddress) []*types.KeyPair {
	prefixStoreKey := types.GetContractAccPubDbKey(contractAddress, account)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	iter := prefixStore.Iterator(nil, nil)
	var pKeyPair []*types.KeyPair
	for ; iter.Valid(); iter.Next() {
		pKeyPair = append(pKeyPair, &types.KeyPair{Key: string(iter.Key()), Value: string(iter.Value())})
	}
	return pKeyPair
}

//GetContractPublicStateByKey
func (k Keeper) GetContractPublicStateByKey(ctx sdk.Context, contractAddress sdk.AccAddress, key []byte) types.KeyPair {
	prefixStoreKey := types.GetContractPubDbKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		if bytes.Equal(iter.Key(), key) {
			return types.KeyPair{Key: string(iter.Key()), Value: string(iter.Value())}
		}
	}
	return types.KeyPair{}
}

//GetContractPublicStateValue gets the value from the key-value store of the public state
func (k Keeper) GetContractPublicStateValue(ctx sdk.Context, contractAddress sdk.AccAddress, key []byte) []byte {
	prefixStoreKey := types.GetContractPubDbKey(contractAddress)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		if bytes.Equal(iter.Key(), key) {
			return iter.Value()
		}
	}
	return nil
}

//GetContractPublicStateValueForAddr gets the value from the key-value store of the public state for a given address
func (k Keeper) GetContractPublicStateValueForAddr(ctx sdk.Context, contractAddress sdk.AccAddress, accAddr sdk.AccAddress, key []byte) []byte {
	prefixStoreKey := types.GetContractAccPubDbKey(contractAddress, accAddr)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)
	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		if bytes.Equal(iter.Key(), key) {
			return iter.Value()
		}
	}
	return nil
}

func (k Keeper) IterateContractInfo(ctx sdk.Context, cb func(sdk.AccAddress, types.ContractInfo) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ContractKeyPrefix)
	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var contract types.ContractInfo
		k.cdc.MustUnmarshal(iter.Value(), &contract)
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
	k.cdc.MustUnmarshal(codeInfoBz, &codeInfo)
	return &codeInfo
}

func (k Keeper) GetCodeHash(ctx sdk.Context, codeID uint64) (codeHash []byte) {
	store := ctx.KVStore(k.storeKey)
	var codeInfo types.CodeInfo
	codeInfoBz := store.Get(types.GetCodeKey(codeID))
	if codeInfoBz == nil {
		return nil
	}

	k.cdc.MustUnmarshal(codeInfoBz, &codeInfo)
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
		k.cdc.MustUnmarshal(iter.Value(), &c)
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
	k.cdc.MustUnmarshal(codeInfoBz, &codeInfo)
	return k.wasmer.GetCode(codeInfo.CodeHash)
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
