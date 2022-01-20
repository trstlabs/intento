package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/trstlabs/trst/x/compute/internal/types"
)

// IterateContractsQueue iterates over the items in the inactive item queue
// and performs a callback function
func (k Keeper) IterateContractQueue(ctx sdk.Context, endTime time.Time, cb func(contract types.ContractInfoWithAddress) (stop bool)) {
	iterator := k.ContractQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {

		//fmt.Printf("iterator key is:  %s ", string(iterator.Key()))
		//get the contract from endTime (key)
		addr, _ := types.SplitContractQueueKey(iterator.Key())
		contractAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return
		}
		fmt.Printf("addr is:  %s ", addr)

		contract := k.GetContractInfoWithAddress(ctx, contractAddr)

		fmt.Printf("info creator is:  %s \n", contract.ContractInfo.Creator)

		if cb(contract) {
			break
		}
	}
}

// IterateContractsQueue iterates over the items in the inactive item queue
// and performs a callback function
func (k Keeper) IterateContractQueueAddressOnly(ctx sdk.Context, endTime time.Time, cb func(addr string) (stop bool)) {
	iterator := k.ContractQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		addr, _ := types.SplitContractQueueKey(iterator.Key())
		if cb(addr) {
			break
		}
	}
}

// GetAllContractAddresses returns all seller items on chain based on the seller
func (k Keeper) GetAllContractAddresses(ctx sdk.Context) (addressList []*string) {
	k.IterateContractQueueAddressOnly(ctx, ctx.BlockHeader().Time, func(address string) bool {
		addressList = append(addressList, &address)
		return false
	})
	return
}

// ContractQueueIterator returns an sdk.Iterator for all the items in the Inactive Queue that expire by endTime
func (k Keeper) ContractQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.ContractQueuePrefix, sdk.PrefixEndBytes(types.ContractByTimeKey(endTime))) //we check the end of the bites array for the end time
}

// InsertContractQueue Inserts a contract into the inactive item queue at endTime
func (k Keeper) InsertContractQueue(ctx sdk.Context, contractAddr string, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := []byte(contractAddr)

	//here the key is time+contract appended (as bytes) and value is contract in bytes
	store.Set(types.ContractQueueKey(contractAddr, endTime), bz)
}

// RemoveFromContractQueue removes a contract from the Inactive Item Queue
func (k Keeper) RemoveFromContractQueue(ctx sdk.Context, contractAddr string, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ContractQueueKey(contractAddr, endTime))
}
