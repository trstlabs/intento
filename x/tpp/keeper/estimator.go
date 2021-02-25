package keeper

import (
	"github.com/tendermint/tendermint/crypto"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/danieljdd/tpp/x/tpp/types"
	"strconv"
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

// CreateEstimator creates a estimator with a new id and update the count
func (k Keeper) CreateEstimator(ctx sdk.Context, msg types.MsgCreateEstimator) {
	// Create the estimator
	count := k.GetEstimatorCount(ctx)
	var estimator = types.Estimator{
		Estimator:                 msg.Estimator,
		Estimation:              msg.Estimation,
		Estimatorestimationhash: msg.Estimatorestimationhash,
		Itemid:                  msg.Itemid,
		Deposit:                 msg.Deposit,
		Interested:              msg.Interested,
		Comment:                 msg.Comment,
	}
	

	estimatoraddress, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		panic(err)
	}

	moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))

	/*sdkError := k.bankKeeper.SendCoinsFromAccountToModule(ctx, estimatoraddress, moduleAcct.String(), sdk.NewCoins(msg.Deposit))
	if sdkError != nil {
		return
	}*/

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, estimatoraddress, moduleAcct.String(), sdk.NewCoins(msg.Deposit)); err != nil {
		panic(err)
	}


	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorKey))
	key := types.KeyPrefix(types.EstimatorKey + estimator.Itemid + "-" + estimator.Estimator)
	value := k.cdc.MustMarshalBinaryBare(&estimator)
	store.Set(key, value)

	// Update estimator count
	k.SetEstimatorCount(ctx, count+1)
}

// SetEstimator set a specific estimator in the store
func (k Keeper) SetEstimator(ctx sdk.Context, estimator types.Estimator) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorKey))
	b := k.cdc.MustMarshalBinaryBare(&estimator)
	store.Set(types.KeyPrefix(types.EstimatorKey+estimator.Itemid + "-" + estimator.Estimator), b)
}

// GetEstimator returns a estimator from its key
func (k Keeper) GetEstimator(ctx sdk.Context, key string) types.Estimator {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorKey))
	var estimator types.Estimator
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.KeyPrefix(types.EstimatorKey+key)), &estimator)
	return estimator
}

// HasEstimator checks if the estimator exists
func (k Keeper) HasEstimator(ctx sdk.Context, id string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorKey))
	return store.Has(types.KeyPrefix(types.EstimatorKey + id))
}

// GetEstimatorOwner returns the creator of the estimator
func (k Keeper) GetEstimatorOwner(ctx sdk.Context, key string) string {
	return k.GetEstimator(ctx, key).Estimator
}

// DeleteEstimator deletes a estimator
func (k Keeper) DeleteEstimator(ctx sdk.Context, key string){
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EstimatorKey))
	var estimator types.Estimator
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.KeyPrefix(types.EstimatorKey+key)), &estimator)
	estimatoraddress, err := sdk.AccAddressFromBech32(estimator.Estimator)
	if err != nil {
		panic(err)
	}
	moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
	sdkErrorEstimator := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, moduleAcct.String(), estimatoraddress, sdk.NewCoins(estimator.Deposit))
		if sdkErrorEstimator != nil {
			panic(err)
		}
	
	store.Delete(types.KeyPrefix(types.EstimatorKey + key))
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
