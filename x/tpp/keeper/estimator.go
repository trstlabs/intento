package keeper

import (
	//"github.com/tendermint/tendermint/crypto"

	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"strconv"

	"github.com/danieljdd/tpp/x/tpp/types"
)

// GetEstimatorCount get the total number of estimator
func (k Keeper) GetEstimatorCount(ctx sdk.Context) int64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorCountKey))
	byteKey := types.KeyPrefix(types.EstimatorCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseInt(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to int64
		panic("cannot decode count")
	}

	return count
}

// SetEstimatorCount set the total number of estimator
func (k Keeper) SetEstimatorCount(ctx sdk.Context, count int64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorCountKey))
	byteKey := types.KeyPrefix(types.EstimatorCountKey)
	bz := []byte(strconv.FormatInt(count, 10))
	store.Set(byteKey, bz)
}

// CreateEstimation creates a estimator with a new id and update the count
func (k Keeper) CreateEstimation(ctx sdk.Context, msg types.MsgCreateEstimation) error {

	item := k.GetItem(ctx, msg.Itemid)

	fmt.Printf("Keeper  item: %X\n", item.Contract)
	// Create the estimator
	count := k.GetEstimatorCount(ctx)
	deposit := sdk.NewInt64Coin("tpp", msg.Deposit)
	var estimator = types.Estimator{
		Estimator: msg.Estimator,
		//	Estimation: msg.Estimation,
		//Estimatemsg: msg.Estimatemsg,
		//Estimatorestimationhash: msg.Estimatorestimationhash,
		Itemid:     msg.Itemid,
		Deposit:    deposit,
		Interested: msg.Interested,
		//Msg:        msg.Msg,
		//	Comment: msg.Comment,
	}

	estimatorAddress, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "address invalid")
	}

	coins := sdk.NewCoins(deposit)

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, estimatorAddress, types.ModuleName, coins)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "sending coins failed")
	}

	contractAddr, err := sdk.AccAddressFromBech32(item.Contract)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "contract address invalid")
	}
	fmt.Printf("executing contract: %X\n", item.Contract)
	fmt.Printf("executing contract addr: %s", item.Contract)
	res, err := k.computeKeeper.Execute(ctx, contractAddr, estimatorAddress, msg.Estimatemsg, sdk.NewCoins(sdk.NewCoin("tpp", sdk.ZeroInt())), nil)
	if err != nil {
		return sdkerrors.Wrap(err, "Execution failed")
	}
	//fmt.Printf("result: Got result for item %s: %s\n", contractAddr, res)
	//fmt.Printf("result: Got result for item %s: %s\n", contractAddr, string(res.Data))

	//fmt.Printf("result: Got log for item %s: %s\n", contractAddr, res.Log)
	var raw map[string]json.RawMessage
	_ = json.Unmarshal([]byte(res.Log), &raw)
	//fmt.Printf("log: Got Unmarshal raw for item %s: %s\n", raw, contractAddr)
	//var values map[string]interface{}
	//fmt.Println(json.Unmarshal([]byte(res.Log), &msg))
	//fmt.Printf("log: Got Unmarshal msg for item %s: %s\n", values, contractAddr)

	var result types.EstimateResult
	json.Unmarshal([]byte(res.Log), &result)
	//fmt.Printf("log: Got Unmarshal msg for item %s: %s\n", strconv.Itoa(result.Estimation.TotalCount), contractAddr)

	if result.Estimation.Status != "" {
		item.Estimationtotal = int64(result.Estimation.TotalCount)
		store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorKey))
		key := append(append([]byte(types.EstimatorKey), types.Uint64ToByte(estimator.Itemid)...), []byte(estimator.Estimator)...)
		value := k.cdc.MustMarshal(&estimator)
		store.Set(key, value)
		k.SetItem(ctx, item)

		// Update estimator count
		k.SetEstimatorCount(ctx, count+1)
		//if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, estimatorAddress, moduleAcct.String(), ); err != nil {
		//	panic(err)
		//}

	} else {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "error after executing estimation")
	}
	return nil
}

// SetEstimator set a specific estimator in the store
func (k Keeper) SetEstimator(ctx sdk.Context, estimator types.Estimator) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorKey))
	b := k.cdc.MustMarshal(&estimator)
	appended := append([]byte(types.EstimatorKey), types.Uint64ToByte(estimator.Itemid)...)
	store.Set(append(appended, estimator.Estimator...), b)
}

// GetEstimator returns a estimator from its key
func (k Keeper) GetEstimator(ctx sdk.Context, key []byte) types.Estimator {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorKey))
	var estimator types.Estimator
	k.cdc.MustUnmarshal(store.Get(append(types.KeyPrefix(types.EstimatorKey), key...)), &estimator)
	return estimator
}

// HasEstimator checks if the estimator exists
func (k Keeper) HasEstimator(ctx sdk.Context, key []byte) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorKey))
	return store.Has(append(types.KeyPrefix(types.EstimatorKey), key...))
}

// GetEstimatorOwner returns the creator of the estimator
func (k Keeper) GetEstimatorOwner(ctx sdk.Context, key []byte) string {
	return k.GetEstimator(ctx, key).Estimator
}

// DeleteEstimation deletes a estimator
func (k Keeper) DeleteEstimation(ctx sdk.Context, key []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorKey))
	var estimator types.Estimator
	k.cdc.MustUnmarshal(store.Get(append(types.KeyPrefix(types.EstimatorKey), key...)), &estimator)
	estimatorAddress, err := sdk.AccAddressFromBech32(estimator.Estimator)
	if err == nil {
		//panic(err)
		//moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, estimatorAddress, sdk.NewCoins(estimator.Deposit))
		if err != nil {
			panic(err)
		}

		store.Delete(append(types.KeyPrefix(types.EstimatorKey), key...))
	}

}

// DeleteEstimationWithReward deletes a estimator
func (k Keeper) DeleteEstimationWithReward(ctx sdk.Context, key []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorKey))
	var estimator types.Estimator
	k.cdc.MustUnmarshal(store.Get(append(types.KeyPrefix(types.EstimatorKey), key...)), &estimator)
	estimatorAddress, err := sdk.AccAddressFromBech32(estimator.Estimator)
	if err == nil {
		//panic(err)
		//moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, estimatorAddress, sdk.NewCoins(estimator.Deposit.Add(estimator.Deposit)))
		if err != nil {
			panic(err)
		}
		store.Delete(append(types.KeyPrefix(types.EstimatorKey), key...))
	}

}

// DeleteEstimationWithoutDeposit deletes a estimator without returing a deposit
func (k Keeper) DeleteEstimationWithoutDeposit(ctx sdk.Context, key []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorKey))
	var estimator types.Estimator
	k.cdc.MustUnmarshal(store.Get(append(types.KeyPrefix(types.EstimatorKey), key...)), &estimator)

	store.Delete(append(types.KeyPrefix(types.EstimatorKey), key...))
}

// GetAllEstimator returns all estimator
func (k Keeper) GetAllEstimator(ctx sdk.Context) (msgs []types.Estimator) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorKey))
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.EstimatorKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var msg types.Estimator
		k.cdc.MustUnmarshal(iterator.Value(), &msg)
		msgs = append(msgs, msg)
	}

	return
}

// FlagMessage flags an item
func (k Keeper) Flag(ctx sdk.Context, item types.Item, msg types.MsgFlagItem) error {

	estimatorAddress, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		return err ///panic(err)

	}
	contractAddr, err := sdk.AccAddressFromBech32(item.Contract)
	if err != nil {
		return err ///panic(err)
	}
	fmt.Printf("executing contract: %s", item.Contract)
	res, err := k.computeKeeper.Execute(ctx, contractAddr, estimatorAddress, msg.Flagmsg, sdk.NewCoins(sdk.NewCoin("tpp", sdk.ZeroInt())), nil)
	if err != nil {
		fmt.Printf("err executing: ")
		//return sdkerrors.Wrapf(types.ErrInvalid, "err %s must be greater %d ",err, msg.Flagmsg)
		return err ///panic(err)
	}
	fmt.Printf("res for item %s: %s\n", res.Log, contractAddr)

	var result types.StatusResult

	_ = json.Unmarshal([]byte(res.Log), &result)
	if result.StatusOnly.Status == "Success" {
		for _, element := range item.Estimatorlist {

			key := append(types.Uint64ToByte(item.Id), []byte(element)...)

			k.DeleteEstimation(ctx, key)
		}
		k.RemoveFromListedItemQueue(ctx, item.Id, item.Endtime)
		_ = k.DeleteItemContract(ctx, item.Contract)
		k.DeleteItem(ctx, item.Id)
		k.RemoveFromItemSeller(ctx, item.Id, item.Seller)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(types.EventTypeItemRemoved, sdk.NewAttribute(types.AttributeKeyCreator, item.Title), sdk.NewAttribute(types.AttributeKeyItemID, strconv.FormatUint(item.Id, 10))),
		)
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeItemFlagged, sdk.NewAttribute(types.AttributeKeyCreator, item.Creator), sdk.NewAttribute(types.AttributeKeyItemID, strconv.FormatUint(item.Id, 10))),
	)

	return nil
}

// DeleteEstimation deletes an estimation
func (k Keeper) DeleteEncryptedEstimation(ctx sdk.Context, item types.Item, msg types.MsgDeleteEstimation) error {

	estimatorAddress, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		return err ///panic(err)

	}
	contractAddr, err := sdk.AccAddressFromBech32(item.Contract)
	if err != nil {
		return err ///panic(err)
	}
	fmt.Printf("executing contract: %s", item.Contract)
	res, err := k.computeKeeper.Execute(ctx, contractAddr, estimatorAddress, msg.Deletemsg, sdk.NewCoins(sdk.NewCoin("tpp", sdk.ZeroInt())), nil)
	if err != nil {
		fmt.Printf("err executing: ")
		return err ///panic(err)
	}
	fmt.Printf("res for item %s: %s\n", res.Log, contractAddr)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeItemCreated, sdk.NewAttribute(types.AttributeKeyCreator, item.Creator), sdk.NewAttribute(types.AttributeKeyItemID, strconv.FormatUint(item.Id, 10))),
	)

	return nil
}
