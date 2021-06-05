package keeper

import (
	//"github.com/tendermint/tendermint/crypto"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
func (k Keeper) CreateEstimation(ctx sdk.Context, msg types.MsgCreateEstimation) {
	// Create the estimator
	count := k.GetEstimatorCount(ctx)
	deposit := sdk.NewInt64Coin("tpp", msg.Deposit)
	var estimator = types.Estimator{
		Estimator:  msg.Estimator,
		Estimation: msg.Estimation,
		//Estimatorestimationhash: msg.Estimatorestimationhash,
		Itemid:     msg.Itemid,
		Deposit:    deposit,
		Interested: msg.Interested,
		Comment:    msg.Comment,
	}

	estimatorAddress, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		panic(err)
	}

	Coins := sdk.NewCoins(deposit)

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, estimatorAddress, types.ModuleName, Coins)
	if err != nil {
		panic(err)
	}

	//if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, estimatorAddress, moduleAcct.String(), ); err != nil {
	//	panic(err)
	//}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorKey))

	key := append(append([]byte(types.EstimatorKey), types.Uint64ToByte(estimator.Itemid)...), []byte(estimator.Estimator)...)
	value := k.cdc.MustMarshalBinaryBare(&estimator)
	store.Set(key, value)

	// Update estimator count
	k.SetEstimatorCount(ctx, count+1)
}

// SetEstimator set a specific estimator in the store
func (k Keeper) SetEstimator(ctx sdk.Context, estimator types.Estimator) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorKey))
	b := k.cdc.MustMarshalBinaryBare(&estimator)
	appended := append([]byte(types.EstimatorKey), types.Uint64ToByte(estimator.Itemid)...)
	store.Set(append(appended, estimator.Estimator...), b)
}

// GetEstimator returns a estimator from its key
func (k Keeper) GetEstimator(ctx sdk.Context, key []byte) types.Estimator {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorKey))
	var estimator types.Estimator
	k.cdc.MustUnmarshalBinaryBare(store.Get(append(types.KeyPrefix(types.EstimatorKey), key...)), &estimator)
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
	k.cdc.MustUnmarshalBinaryBare(store.Get(append(types.KeyPrefix(types.EstimatorKey), key...)), &estimator)
	estimatorAddress, err := sdk.AccAddressFromBech32(estimator.Estimator)
	if err != nil {
		panic(err)
	}
	//moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, estimatorAddress, sdk.NewCoins(estimator.Deposit))
	if err != nil {
		panic(err)
	}

	store.Delete(append(types.KeyPrefix(types.EstimatorKey), key...))
}

// DeleteEstimationWithReward deletes a estimator
func (k Keeper) DeleteEstimationWithReward(ctx sdk.Context, key []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorKey))
	var estimator types.Estimator
	k.cdc.MustUnmarshalBinaryBare(store.Get(append(types.KeyPrefix(types.EstimatorKey), key...)), &estimator)
	estimatorAddress, err := sdk.AccAddressFromBech32(estimator.Estimator)
	if err != nil {
		panic(err)
	}
	//moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, estimatorAddress, sdk.NewCoins(estimator.Deposit.Add(estimator.Deposit)))
	if err != nil {
		panic(err)
	}

	store.Delete(append(types.KeyPrefix(types.EstimatorKey), key...))
}

// DeleteEstimationWithoutDeposit deletes a estimator without returing a deposit
func (k Keeper) DeleteEstimationWithoutDeposit(ctx sdk.Context, key []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorKey))
	var estimator types.Estimator
	k.cdc.MustUnmarshalBinaryBare(store.Get(append(types.KeyPrefix(types.EstimatorKey), key...)), &estimator)

	store.Delete(append(types.KeyPrefix(types.EstimatorKey), key...))
}

// GetAllEstimator returns all estimator
func (k Keeper) GetAllEstimator(ctx sdk.Context) (msgs []types.Estimator) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorKey))
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.EstimatorKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var msg types.Estimator
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &msg)
		msgs = append(msgs, msg)
	}

	return
}
