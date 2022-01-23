package keeper

import (
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
		//fmt.Printf("addr is:  %s ", addr)

		contract := k.GetContractInfoWithAddress(ctx, contractAddr)

		//fmt.Printf("info creator is:  %s \n", contract.ContractInfo.Creator)

		if cb(contract) {
			break
		}
	}
}

/*
// IterateContractsQueue iterates over the items in the inactive item queue
// and performs a callback function
func (k Keeper) IterateContractsInQueueForReward(ctx sdk.Context, endTime time.Time, cb func(contract types.ContractInfoWithAddress) (stop bool)) {
	iterator := k.ContractQueueIterator(ctx, endTime)
	params := k.GetParams(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {

		addr, _ := types.SplitContractQueueKey(iterator.Key())
		contractAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return
		}

		contract := k.GetContractInfoWithAddress(ctx, contractAddr)
		info := k.GetCodeInfo(ctx, contract.ContractInfo.CodeID)
		if info.Duration >= params.MinContractDurationForIncentive {

			if cb(contract) {
				fmt.Printf("found addr:  %s ", addr)
				break
			}
		} else {
			return
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
}*/

// GetContractAddressesForBlock returns all expiring contracts for a block
func (k Keeper) GetContractAddressesForBlock(ctx sdk.Context) (incentiveList []string, contracts []types.ContractInfoWithAddress) {
	params := k.GetParams(ctx)
	k.IterateContractQueue(ctx, ctx.BlockHeader().Time, func(contract types.ContractInfoWithAddress) bool {

		if contract.ContractInfo.AutoMsg != nil {
			info := k.GetCodeInfo(ctx, contract.ContractInfo.CodeID)
			if info.Duration >= params.MinContractDurationForIncentive {
				if k.bankKeeper.GetBalance(ctx, contract.Address, "utrst").Amount.SubRaw(params.MinContractBalanceForIncentive).IsPositive() {
					incentiveList = append(incentiveList, contract.Address.String())
				}
			}
		}

		contracts = append(contracts, contract)
		return false
	})
	return
}

/*
// GetContractAddressesForBlock returns all expiring contracts for a block
func (k Keeper) GetContractAddresses(ctx sdk.Context) (incentiveList []string) {
	params := k.GetParams(ctx)
	k.IterateContractQueue(ctx, ctx.BlockHeader().Time, func(contract types.ContractInfoWithAddress) bool {

		if contract.ContractInfo.AutoMsg != nil {
			info := k.GetCodeInfo(ctx, contract.ContractInfo.CodeID)
			if info.Duration >= params.MinContractDurationForIncentive {
				//if k.bankKeeper.GetBalance(contract.Address, "utrst") > params.MinContractBalanceForIncentive {
				incentiveList = append(incentiveList, contract.Address.String())
				//}
			}
		}

		return false
	})
	return
}
*/
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
