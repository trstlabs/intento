package keeper

import (
	//"github.com/tendermint/tendermint/crypto"

	"encoding/json"
	"fmt"

	//"github.com/coreos/etcd/store"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"strconv"

	"github.com/trstlabs/trst/x/item/types"
)

// GetProfileCount get the total number of profiles
func (k Keeper) GetProfileCount(ctx sdk.Context) int64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProfileCountKey))
	byteKey := types.KeyPrefix(types.ProfileCountKey)
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

// SetProfileCount sets the total number of profiles
func (k Keeper) SetProfileCount(ctx sdk.Context, count int64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProfileCountKey))
	byteKey := types.KeyPrefix(types.ProfileCountKey)
	bz := []byte(strconv.FormatInt(count, 10))
	store.Set(byteKey, bz)
}

// CreateEstimation creates a estimator with a new id and update the count
func (k Keeper) CreateEstimation(ctx sdk.Context, msg types.MsgCreateEstimation) error {

	item := k.GetItem(ctx, msg.Itemid)

	//var profile types.Profile
	var estimationInfo = types.EstimationInfo{
		Itemid:      msg.Itemid,
		Interested:  msg.Interested,
		ItemCreator: item.Creator,
	}
	//	fmt.Printf("executiddng consstract: %X\n", item.Contract)
	var key = append([]byte(types.ProfileKey), []byte(msg.Estimator)...)
	var profile types.Profile
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProfileKey))
	if !store.Has(key) {
		var count = k.GetProfileCount(ctx)
		profile.Owner = msg.Estimator
		k.SetProfileCount(ctx, count+1)

	} else {
		profile = k.GetProfile(ctx, msg.Estimator)
	}
	//	fmt.Printf("executaing contract: %X\n", item.Contract)
	amountEstimations := len(profile.Estimations)
	var amountSameCreator int
	for _, estimation := range profile.Estimations {
		if estimation.ItemCreator == estimationInfo.ItemCreator {
			amountSameCreator = amountSameCreator + 1
		}
	}
	if amountSameCreator != 0 {
		params := k.GetParams(ctx)
		if (amountSameCreator/amountEstimations)*100 > int(params.MaxEstimatorCreatorRatio) {
			return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "cannot estimate again for this creator")
		}
	}
	profile.Estimations = append(profile.Estimations, &estimationInfo)

	//fmt.Printf("execusadfting contract: %X\n", item.Contract)
	estimatorAddress, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "address invalid")
	}

	contractAddr, err := sdk.AccAddressFromBech32(item.Estimation.Contract)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "contract address invalid")
	}

	res, err := k.computeKeeper.Execute(ctx, contractAddr, estimatorAddress, msg.EstimateMsg, sdk.NewCoins(sdk.NewInt64Coin("utrst", msg.Deposit)), nil)
	if err != nil {
		return sdkerrors.Wrap(err, "Execution failed")
	}

	var raw map[string]json.RawMessage
	_ = json.Unmarshal([]byte(res.Log), &raw)

	var result types.EstimateResult
	json.Unmarshal([]byte(res.Log), &result)

	if result.Estimation.Status != "" {

		item.Estimation.EstimationTotal = int64(result.Estimation.TotalCount)

		k.SetItem(ctx, item)

		//b := k.cdc.MustMarshal(&profile)
		k.SetProfile(ctx, profile, msg.Estimator)

	} else {
		//fmt.Printf("result: got result for estimation %s: %s\n", contractAddr, result.Estimation.Status)
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "error during execution of estimation")
	}
	return nil
}

// SetProfile set a specific profile in the store
func (k Keeper) SetProfile(ctx sdk.Context, profile types.Profile, owner string) {
	var count = k.GetProfileCount(ctx)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProfileKey))
	b := k.cdc.MustMarshal(&profile)
	appended := append(types.KeyPrefix(types.ProfileKey), []byte(owner)...)

	k.SetProfileCount(ctx, count+1)
	store.Set(appended, b)
}

// GetProfile returns profile info from its owner
func (k Keeper) GetProfile(ctx sdk.Context, owner string) types.Profile {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProfileKey))
	var profile types.Profile
	k.cdc.MustUnmarshal(store.Get(append(types.KeyPrefix(types.ProfileKey), types.KeyPrefix(owner)...)), &profile)
	return profile
}

// UpdateEstimationInfo gets a info from a specific profile in the store
func (k Keeper) UpdateEstimationInfo(ctx sdk.Context, estimationInfo types.EstimationInfo, estimator string) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProfileKey))
	var profile types.Profile
	var key = append(types.KeyPrefix(types.ProfileKey), []byte(estimator)...)
	k.cdc.MustUnmarshal(store.Get(append(types.KeyPrefix(types.ProfileKey), key...)), &profile)

	for index, estimation := range profile.Estimations {
		if estimation.Itemid == estimationInfo.Itemid {
			profile.Estimations[index] = &estimationInfo

		} else {
			return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "estimation info not found")
		}
	}

	b := k.cdc.MustMarshal(&profile)
	store.Set(key, b)
	return nil
}

// HasProfile checks if the estimator exists
func (k Keeper) HasProfile(ctx sdk.Context, key []byte) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProfileKey))
	return store.Has(append(types.KeyPrefix(types.ProfileKey), key...))
}

// DeleteProfile deletes a profile
func (k Keeper) DeleteProfile(ctx sdk.Context, key []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProfileKey))
	var profile types.Profile
	k.cdc.MustUnmarshal(store.Get(append(types.KeyPrefix(types.ProfileKey), key...)), &profile)

	store.Delete(append(types.KeyPrefix(types.ProfileKey), key...))
}

// GetAllProfiles returns all estimator profiles
func (k Keeper) GetAllProfiles(ctx sdk.Context) (msgs []types.Profile) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ProfileKey))
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.ProfileKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var msg types.Profile
		k.cdc.MustUnmarshal(iterator.Value(), &msg)
		msgs = append(msgs, msg)
	}

	return
}

// Flag flags an item in the contract
func (k Keeper) Flag(ctx sdk.Context, item types.Item, msg types.MsgFlagItem) error {

	estimatorAddress, err := sdk.AccAddressFromBech32(msg.Estimator)
	if err != nil {
		return err ///panic(err)

	}
	contractAddr, err := sdk.AccAddressFromBech32(item.Estimation.Contract)
	if err != nil {
		return err ///panic(err)
	}
	//fmt.Printf("executing contract: %s", item.Contract)
	res, err := k.computeKeeper.Execute(ctx, contractAddr, estimatorAddress, msg.FlagMsg, sdk.NewCoins(sdk.NewCoin("utrst", sdk.ZeroInt())), nil)
	if err != nil {
		fmt.Printf("err executing: ")
		//return sdkerrors.Wrapf(types.ErrInvalid, "err %s must be greater %d ",err, msg.Flagmsg)
		return err ///panic(err)
	}
	//fmt.Printf("res for item %s: %s\n", res.Log, contractAddr)

	var result types.StatusResult

	_ = json.Unmarshal([]byte(res.Log), &result)
	if result.StatusOnly.Status == "Success" {
		/*for _, element := range item.Estimation.EstimatorList {

			key := append(types.Uint64ToByte(item.Id), []byte(element)...)

			k.DeleteEstimation(ctx, key)
		}*/
		k.RemoveFromListedItemQueue(ctx, item.Id, item.ListingDuration.EndTime)
		_ = k.DeleteItemContract(ctx, item.Estimation.Contract)
		k.DeleteItem(ctx, item.Id)
		k.RemoveFromSellerItems(ctx, item.Id, item.Transfer.Seller)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(types.EventTypeItemRemoved, sdk.NewAttribute(types.AttributeKeyCreator, item.Title), sdk.NewAttribute(types.AttributeKeyItemID, strconv.FormatUint(item.Id, 10))),
		)
	} else {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "flagging item not possible")
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
	contractAddr, err := sdk.AccAddressFromBech32(item.Estimation.Contract)
	if err != nil {
		return err ///panic(err)
	}
	fmt.Printf("executing contract: %s", item.Estimation.Contract)
	res, err := k.computeKeeper.Execute(ctx, contractAddr, estimatorAddress, msg.DeleteMsg, sdk.NewCoins(sdk.NewCoin("utrst", sdk.ZeroInt())), nil)
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
