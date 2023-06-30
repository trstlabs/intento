package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/trstlabs/trst/x/compute/internal/types"
)

// IterateContractsQueue iterates over the items in the inactive contract queue
// and performs a callback function
func (k Keeper) IterateContractQueue(ctx sdk.Context, endTime time.Time, cb func(contract types.ContractInfoWithAddress) (stop bool)) {
	iterator := k.ContractQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {

		addr, _ := types.SplitContractQueueKey(iterator.Key())
		contractAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return
		}
		contract := k.GetContractInfoWithAddress(ctx, contractAddr)

		//fmt.Printf("info creator is:  %s \n", contract.ContractInfo.Creator)

		if cb(contract) {
			break
		}
	}
}

// GetContractAddressesForBlock returns all expiring contracts for a block
func (k Keeper) GetContractAddressesForBlock(ctx sdk.Context) (contracts []types.ContractInfoWithAddress) {
	k.IterateContractQueue(ctx, ctx.BlockHeader().Time, func(contract types.ContractInfoWithAddress) bool {

		contracts = append(contracts, contract)
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
