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
func (k Keeper) SetContractPublicState(ctx sdk.Context, contrAddr sdk.AccAddress, result []wasmTypes.Attribute) error {
	prefixStoreKey := types.GetContractPubDbKey(contrAddr)
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixStoreKey)

	for _, attr := range result {

		if attr.Encrypted {
			continue
		} else if len(attr.AccAddr) == 44 {
			accAddr, err := sdk.AccAddressFromBech32(attr.AccAddr)
			if err != nil {
				return err
			}
			prefixAccStoreKey := types.GetContractAccPubDbKey(contrAddr, accAddr)
			prefixAccStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixAccStoreKey)

			prefixAccStore.Set([]byte(attr.Key), attr.Value)
		} else if attr.PubDb {

			prefixStore.Set([]byte(attr.Key), attr.Value)
		}
	}
	return nil
}

// SetAirdropAction sets the airdrop from the contract attributes
func (k Keeper) SetAirdropAction(ctx sdk.Context, result []wasmTypes.Attribute) error {

	for _, attr := range result {
		if attr.Key == "init_auto_swap" {
			acc, err := sdk.AccAddressFromBech32(string(attr.Value))
			if err != nil {
				return err
			}
			k.hooks.AfterAutoSwap(ctx, acc)
		} else if attr.Key == "init_recurring_send" {
			//fmt.Printf("initiated key: %s \n,", (attr.Key))
			//fmt.Printf("initiated val: %s \n,", string(attr.Value))
			acc, err := sdk.AccAddressFromBech32(string(attr.Value))
			if err != nil {
				return err
			}
			k.hooks.AfterRecurringSend(ctx, acc)
		}
	}

	return nil
}

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
