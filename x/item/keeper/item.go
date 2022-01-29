package keeper

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trstlabs/trst/x/item/types"
)

// GetItemCount get the total number of items
func (k Keeper) GetItemCount(ctx sdk.Context) int64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ItemCountKey))
	byteKey := types.KeyPrefix(types.ItemCountKey)
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

// SetItemCount set the total number of items
func (k Keeper) SetItemCount(ctx sdk.Context, count int64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ItemCountKey))
	byteKey := types.KeyPrefix(types.ItemCountKey)
	bz := []byte(strconv.FormatInt(count, 10))
	store.Set(byteKey, bz)
}

// CreateItem creates a item with a new id and update the count
func (k Keeper) CreateItem(ctx sdk.Context, msg types.MsgCreateItem) error {

	// Create the item
	count := k.GetItemCount(ctx)

	submitTime := ctx.BlockHeader().Time

	activePeriod := k.GetParams(ctx).MaxActivePeriod
	endTime := submitTime.Add(activePeriod)

	//k.computeKeeper.Instantiate()

	userAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return err
	}

	contractAddr, err := k.computeKeeper.Instantiate(ctx, uint64(1), userAddress, msg.InitMsg, msg.AutoMsg, "Item "+fmt.Sprint(count), sdk.NewCoins(sdk.NewCoin("utrst", sdk.ZeroInt())), nil, activePeriod)
	if err != nil {
		return err
	}

	var item = types.Item{
		Creator:         msg.Creator,
		Seller:          msg.Creator,
		Id:              uint64(count) + 1,
		Title:           msg.Title,
		Description:     msg.Description,
		ShippingCost:    msg.ShippingCost,
		LocalPickup:     msg.LocalPickup,
		EstimationCount: msg.EstimationCount,
		Contract:        contractAddr.String(),

		Tags: msg.Tags,

		Condition:      msg.Condition,
		ShippingRegion: msg.ShippingRegion,
		DepositAmount:  msg.DepositAmount,
		SubmitTime:     submitTime,
		EndTime:        endTime,
	}

	k.BindItemSeller(ctx, item.Id, msg.Creator)
	//works 100% with endtime tx.BlockHeader().Time
	k.InsertListedItemQueue(ctx, item.Id, item, endTime)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeItemCreated, sdk.NewAttribute(types.AttributeKeyCreator, item.Creator), sdk.NewAttribute(types.AttributeKeyItemID, strconv.FormatUint(item.Id, 10))),
	)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ItemKey))
	key := append(types.KeyPrefix(types.ItemKey), types.Uint64ToByte(item.Id)...)
	value := k.cdc.MustMarshal(&item)
	store.Set(key, value)

	// Update item count
	k.SetItemCount(ctx, count+1)

	return nil
}

// SetItem set a specific item in the store
func (k Keeper) SetItem(ctx sdk.Context, item types.Item) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ItemKey))
	b := k.cdc.MustMarshal(&item)
	store.Set(append(types.KeyPrefix(types.ItemKey), types.Uint64ToByte(item.Id)...), b)
}

// GetItem returns a item from its id
func (k Keeper) GetItem(ctx sdk.Context, id uint64) types.Item {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ItemKey))
	var item types.Item
	k.cdc.MustUnmarshal(store.Get(append(types.KeyPrefix(types.ItemKey), types.Uint64ToByte(id)...)), &item)
	return item
}

// HasItem checks if the item exists
func (k Keeper) HasItem(ctx sdk.Context, id uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ItemKey))
	return store.Has(append(types.KeyPrefix(types.ItemKey), types.Uint64ToByte(id)...))

}

// GetItemBuyer returns the seller of the item
func (k Keeper) GetItemBuyer(ctx sdk.Context, id uint64) string {
	return k.GetItem(ctx, id).Seller
}

// DeleteItem deletes a item
func (k Keeper) DeleteItemContract(ctx sdk.Context, contract string) error {
	contractAddress, err := sdk.AccAddressFromBech32(contract)
	if err != nil {
		return err
	}
	err = k.computeKeeper.Delete(ctx, contractAddress)

	if err != nil {
		return err
	}
	return nil
}

// DeleteItem deletes a item
func (k Keeper) DeleteItem(ctx sdk.Context, key uint64) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ItemKey))
	store.Delete(append(types.KeyPrefix(types.ItemKey), types.Uint64ToByte(key)...))
}

// GetAllItem returns all item
func (k Keeper) GetAllItem(ctx sdk.Context) (msgs []types.Item) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ItemKey))
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.ItemKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var msg types.Item
		k.cdc.MustUnmarshal(iterator.Value(), &msg)
		msgs = append(msgs, msg)
	}

	return
}

// GetAllListedItems returns all inactive item
func (k Keeper) GetAllListedItems(ctx sdk.Context) (msgs []*types.Item) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ListedItemQueuePrefix)
	iterator := sdk.KVStorePrefixIterator(store, types.ListedItemQueuePrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var msg types.Item
		k.cdc.MustUnmarshal(iterator.Value(), &msg)
		msgs = append(msgs, &msg)
	}

	return
}

// SendPaymentToAccount handles payment
func (k Keeper) SendPaymentToAccount(ctx sdk.Context, address string, coinToSend sdk.Coin) {

	userAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		panic(err)
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, userAddress, sdk.NewCoins(coinToSend))
	if err != nil {
		panic(err)
	}
	k.hooks.AfterItemBought(ctx, userAddress)
}

// MintReward mints coins to a module account
func (k Keeper) MintReward(ctx sdk.Context, coinToSend sdk.Coin) {
	k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(coinToSend))

}

// HandleReward handles reward
func (k Keeper) HandleBuyerReward(ctx sdk.Context, coinToSend sdk.Coin, buyer sdk.AccAddress) {
	itemIncentivesAddr := k.accountKeeper.GetModuleAddress(types.ItemIncentivesModuleAcctName)
	balance := k.bankKeeper.GetBalance(ctx, itemIncentivesAddr, "utrst")
	//distribute to the buyer
	params := k.GetParams(ctx)
	if coinToSend.Amount.SubRaw(params.MaxBuyerReward).IsNegative() {
		coinToSend.Amount = sdk.NewInt(params.MaxBuyerReward)
	}
	if coinToSend.Sub(balance).IsNegative() {
		coinToSend = balance
	}
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ItemIncentivesModuleAcctName, buyer, sdk.NewCoins(coinToSend))
}

// TokenizeItem mints coins to a module account
func (k Keeper) TokenizeItem(ctx sdk.Context, itemId uint64, addr string) error {
	userAddress, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return err
	}
	coin := sdk.NewCoins(sdk.NewCoin(("TRSTITEM" + strconv.FormatUint(itemId, 10)), sdk.OneInt()))
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coin)
	if err != nil {
		return err ///panic(err)

	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, userAddress, coin)
	if err != nil {
		return err ///panic(err)
	}
	//k.hooks.AfterItemTokenized(ctx, userAddress)
	return nil
}

// TokenizeItem burns coins from a module account
func (k Keeper) UnTokenizeItem(ctx sdk.Context, itemId uint64, addr string) error {
	userAddress, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return err
	}
	coin := sdk.NewCoins(sdk.NewCoin(("TRSTITEM" + strconv.FormatUint(itemId, 10)), sdk.OneInt()))
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, userAddress, types.ModuleName, coin)
	if err != nil {
		return err ///panic(err)
	}

	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coin)
	if err != nil {
		return err ///panic(err)

	}
	return nil
}

// BurnCoins burns coins from a module account
func (k Keeper) BurnCoins(ctx sdk.Context, coinToBurn sdk.Coin) error {
	err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coinToBurn))
	if err != nil {
		return err ///panic(err)

	}
	return nil
}

// RevealEstimation reveals an item
func (k Keeper) RevealEstimation(ctx sdk.Context, item types.Item, msg types.MsgRevealEstimation) error {

	creatorAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return err ///panic(err)

	}
	contractAddr, err := sdk.AccAddressFromBech32(item.Contract)
	if err != nil {
		return err ///panic(err)
	}
	//fmt.Printf("executing contract: %s", item.Contract)
	res, err := k.computeKeeper.Execute(ctx, contractAddr, creatorAddress, msg.RevealMsg, sdk.NewCoins(sdk.NewCoin("utrst", sdk.ZeroInt())), nil)
	if err != nil {
		fmt.Printf("err executing: ")
		return err ///panic(err)
	}
	//	fmt.Printf("res for item %s: %s\n", res.Log, contractAddr)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeItemCreated, sdk.NewAttribute(types.AttributeKeyCreator, item.Creator), sdk.NewAttribute(types.AttributeKeyItemID, strconv.FormatUint(item.Id, 10))),
	)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ItemKey))
	key := append(types.KeyPrefix(types.ItemKey), types.Uint64ToByte(item.Id)...)
	value := k.cdc.MustMarshal(&item)
	store.Set(key, value)

	var result types.RevealResult
	json.Unmarshal([]byte(res.Log), &result)

	b := make([]int64, len(result.RevealEstimation.EstimationList))
	for i, v := range result.RevealEstimation.EstimationList {
		b[i] = int64(v)
	}

	if result.RevealEstimation.Status == "Success" {
		item.BestEstimator = result.RevealEstimation.BestEstimator
		item.EstimationPrice = int64(result.RevealEstimation.BestEstimation)
		item.Comments = result.RevealEstimation.Comments
		item.EstimationList = b
		k.SetItem(ctx, item)
	} else {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "reveal not possible")
	}
	return nil
}

// RevealEstimation reveals an item
func (k Keeper) Transferable(ctx sdk.Context, item types.Item, msg types.MsgItemTransferable) error {

	creatorAddress, err := sdk.AccAddressFromBech32(msg.Seller)
	if err != nil {
		return err ///panic(err)

	}
	contractAddr, err := sdk.AccAddressFromBech32(item.Contract)
	if err != nil {
		return err ///panic(err)
	}
	fmt.Printf("executing contract: %s", item.Contract)
	res, err := k.computeKeeper.Execute(ctx, contractAddr, creatorAddress, msg.TransferableMsg, sdk.NewCoins(sdk.NewCoin("utrst", sdk.ZeroInt())), nil)
	if err != nil {
		fmt.Printf("err executing: %s", err)
		return err ///panic(err)
	}
	fmt.Printf("res for item %s: %s\n", res.Log, contractAddr)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeItemCreated, sdk.NewAttribute(types.AttributeKeyCreator, item.Creator), sdk.NewAttribute(types.AttributeKeyItemID, strconv.FormatUint(item.Id, 10))),
	)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ItemKey))
	key := append(types.KeyPrefix(types.ItemKey), types.Uint64ToByte(item.Id)...)
	value := k.cdc.MustMarshal(&item)
	store.Set(key, value)

	var result types.TransferableResult
	json.Unmarshal([]byte(res.Log), &result)
	fmt.Println(json.Unmarshal([]byte(res.Log), &result))
	//fmt.Printf("log: Got Unmarshal msg for item %s: %s\n", strconv.Itoa(result.), contractAddr)
	//fmt.Printf("log: Got Unmarshal msg for item %s: %s\n", result.RevealEstimation.Comments[1], contractAddr)
	//	fmt.Printf("log: Got Unmarshal msg for item %s: %s\n", res.Log, contractAddr)

	if result.Transferable.Status == "Success" {
		/*item.BestEstimator = result.RevealEstimation.BestEstimator
		item.EstimationPrice = int64(result.RevealEstimation.Bestestimation)
		item.Comments = result.RevealEstimation.Comments*/
		item.Transferable = true
		k.SetItem(ctx, item)
	} else {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "setting item not possible")
	}
	return nil
}
