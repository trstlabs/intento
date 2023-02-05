package keeper

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	"github.com/tendermint/tendermint/crypto"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

// GetAutoTxInfo
func (k Keeper) GetAutoTxInfo(ctx sdk.Context, autoTxID uint64) types.AutoTxInfo {
	store := ctx.KVStore(k.storeKey)
	var autoTx types.AutoTxInfo
	autoTxBz := store.Get(types.GetAutoTxKey(autoTxID))

	k.cdc.MustUnmarshal(autoTxBz, &autoTx)
	return autoTx
}

// TryGetAutoTxInfo
func (k Keeper) TryGetAutoTxInfo(ctx sdk.Context, autoTxID uint64) (types.AutoTxInfo, error) {
	store := ctx.KVStore(k.storeKey)
	var autoTx types.AutoTxInfo
	autoTxBz := store.Get(types.GetAutoTxKey(autoTxID))

	err := k.cdc.Unmarshal(autoTxBz, &autoTx)
	if err != nil {
		return types.AutoTxInfo{}, err
	}
	return autoTx, nil
}

func (k Keeper) SetAutoTxInfo(ctx sdk.Context, autoTx *types.AutoTxInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetAutoTxKey(autoTx.TxID), k.cdc.MustMarshal(autoTx))
}

func (k Keeper) SendAutoTx(ctx sdk.Context, autoTxInfo types.AutoTxInfo) error {

	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: autoTxInfo.Data,
	}

	channelID, found := k.icaControllerKeeper.GetActiveChannelID(ctx, autoTxInfo.ConnectionID, autoTxInfo.PortID)
	if !found {
		return sdkerrors.Wrapf(icatypes.ErrActiveChannelNotFound, "failed to retrieve active channel for port %s", autoTxInfo.PortID)
	}

	chanCap, found := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(autoTxInfo.PortID, channelID))
	if !found {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	timeoutTimestamp := time.Now().Add(time.Hour).UnixNano()
	//to make sure timeout does not happen soon
	if autoTxInfo.Interval > time.Minute*2 {
		timeoutTimestamp = timeoutTimestamp + autoTxInfo.Interval.Nanoseconds() - time.Minute.Nanoseconds()
	}
	sequence, err := k.icaControllerKeeper.SendTx(ctx, chanCap, autoTxInfo.ConnectionID, autoTxInfo.PortID, packetData, uint64(timeoutTimestamp))
	if err != nil {
		return err
	}

	k.setTmpAutoTxID(ctx, autoTxInfo.TxID, autoTxInfo.PortID, sequence)
	return nil
}

func (k Keeper) CreateAutoTx(ctx sdk.Context, owner sdk.AccAddress, portID string, data []byte, connectionId string, duration time.Duration, interval time.Duration, startAt time.Time, feeFunds sdk.Coins, retries uint64, dependsOn []uint64) error {

	txID := k.autoIncrementID(ctx, types.KeyLastTxID)
	autoTxAddress, err := k.createFeeAccount(ctx, txID, owner, feeFunds)
	if err != nil {
		return err
	}

	endTime, execTime, interval := k.calculateAndInsertQueue(ctx, startAt, duration, txID, interval)

	autoTx := types.AutoTxInfo{
		TxID:           txID,
		FeeAddress:     autoTxAddress.String(),
		Owner:          owner.String(),
		Data:           data,
		Interval:       interval,
		Duration:       duration,
		StartTime:      startAt,
		ExecTime:       execTime,
		EndTime:        endTime,
		PortID:         portID,
		ConnectionID:   connectionId,
		MaxRetries:     retries,
		DependsOnTxIds: dependsOn,
	}

	k.SetAutoTxInfo(ctx, &autoTx)
	k.addToAutoTxOwnerIndex(ctx, owner, txID)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAutoTx,
			sdk.NewAttribute(types.AttributeKeyAutoTxID, strconv.FormatUint(txID, 10)),
		))
	return nil
}

func (k Keeper) createFeeAccount(ctx sdk.Context, txID uint64, owner sdk.AccAddress, feeFunds sdk.Coins) (sdk.AccAddress, error) {
	autoTxAddress := k.generateAutoTxInfoAddress(ctx, txID)
	existingAcct := k.accountKeeper.GetAccount(ctx, autoTxAddress)
	if existingAcct != nil {
		return nil, sdkerrors.Wrap(types.ErrAccountExists, existingAcct.GetAddress().String())
	}

	// deposit initial autoTx funds
	if !feeFunds.IsZero() {
		if k.bankKeeper.BlockedAddr(owner) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "blocked address can not be used")
		}
		sdkerr := k.bankKeeper.SendCoins(ctx, owner, autoTxAddress, feeFunds)
		if sdkerr != nil {
			return nil, sdkerr
		}
	} else {
		// create an empty account (so we don't have issues later)
		autoTxAccount := k.accountKeeper.NewAccountWithAddress(ctx, autoTxAddress)
		k.accountKeeper.SetAccount(ctx, autoTxAccount)
	}
	return autoTxAddress, nil
}

// generates a autoTx address from txID + instanceID
func (k Keeper) generateAutoTxInfoAddress(ctx sdk.Context, txID uint64) sdk.AccAddress {
	instanceID := k.autoIncrementID(ctx, types.KeyLastTxAddrID)
	return autoTxAddress(txID, instanceID)
}

func autoTxAddress(txID, instanceID uint64) sdk.AccAddress {
	// NOTE: It is possible to get a duplicate address if either txID or instanceID
	// overflow 32 bits. This is highly improbable, but something that could be refactored.
	autoTxID := txID<<32 + instanceID
	return addrFromUint64(autoTxID)

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

func addrFromUint64(id uint64) sdk.AccAddress {
	addr := make([]byte, 20)
	addr[0] = 'C'
	binary.PutUvarint(addr[1:], id)
	return sdk.AccAddress(crypto.AddressHash(addr))
}

func (k Keeper) calculateAndInsertQueue(ctx sdk.Context, startTime time.Time, duration time.Duration, autoTxID uint64, interval time.Duration) (time.Time, time.Time, time.Duration) {
	endTime, execTime := calculateEndAndExecTimes(startTime, duration, interval)
	k.InsertAutoTxQueue(ctx, autoTxID, execTime)

	return endTime, execTime, interval
}

func calculateEndAndExecTimes(startTime time.Time, duration time.Duration, interval time.Duration) (time.Time, time.Time) {
	endTime := startTime.Add(duration)

	execTime := calculateExecTime(duration, interval, startTime)

	return endTime, execTime
}

func calculateExecTime(duration, interval time.Duration, startTime time.Time) time.Time {
	if startTime.After(time.Now().Add(time.Minute)) {
		return startTime
	}
	if interval != 0 {
		return startTime.Add(interval)
	}
	return startTime.Add(duration)

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

func (k Keeper) importAutoTxInfo(ctx sdk.Context, autoTxId uint64, autoTxInfo types.AutoTxInfo) error {

	store := ctx.KVStore(k.storeKey)
	key := types.GetAutoTxKey(autoTxId)
	if store.Has(key) {
		return sdkerrors.Wrapf(types.ErrDuplicate, "duplicate code: %d", autoTxId)
	}
	// 0x01 | autoTxId (uint64) -> autoTxInfo
	store.Set(key, k.cdc.MustMarshal(&autoTxInfo))
	return nil
}

func (k Keeper) IterateAutoTxInfos(ctx sdk.Context, cb func(uint64, types.AutoTxInfo) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.AutoTxKeyPrefix)
	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var c types.AutoTxInfo
		k.cdc.MustUnmarshal(iter.Value(), &c)
		// cb returns true to stop early
		if cb(binary.BigEndian.Uint64(iter.Key()), c) {
			return
		}
	}
}

// addToAutoTxOwnerIndex adds element to the index for autoTxs-by-creator queries
func (k Keeper) addToAutoTxOwnerIndex(ctx sdk.Context, ownerAddress sdk.AccAddress, autoTxID uint64) error {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetAutoTxByOwnerIndexKey(ownerAddress, autoTxID), []byte{})
	return nil
}

// IterateAutoTxsByOwner iterates over all autoTxs with given creator address in order of creation time asc.
func (k Keeper) IterateAutoTxsByOwner(ctx sdk.Context, owner sdk.AccAddress, cb func(address sdk.AccAddress) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetAutoTxsByOwnerPrefix(owner))
	for iter := prefixStore.Iterator(nil, nil); iter.Valid(); iter.Next() {
		key := iter.Key()
		if cb(key /* [types.TimeTimeLen:] */) {
			return
		}
	}
}

/*

// GetLatestAutoTxByICAPort returns the id for a port
func (k Keeper) GetLatestAutoTxByICAPort(ctx sdk.Context, port string) (uint64, error) {
	owner := port[14:]
	ownerAddress, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		return 0, err
	}
	var id []byte
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetAutoTxsByOwnerPrefix(ownerAddress))
	for iter := prefixStore.ReverseIterator(nil, nil); iter.Valid(); {
		id = iter.Key()
		fmt.Printf("GetLatestAutoTxByICAPort: %v", id)
		//iter.Key(), nil
		break

	}
	fmt.Printf("GetLatestAutoTxByICAPort uint: %v", binary.BigEndian.Uint64(id))
	return binary.BigEndian.Uint64(id), nil
} */

// SetAutoTxResult sets the result of the last executed TxID set at SendAutoTx. As Interchain accounts IBC channels are ORDERED, this should work perfectly..
func (k Keeper) SetAutoTxResult(ctx sdk.Context, port string, rewardType int, seq uint64) error {
	id := k.getTmpAutoTxID(ctx, port, seq)
	if id <= 0 {
		return nil
	}

	k.Logger(ctx).Debug("auto_tx_result", "id", id)

	txInfo := k.GetAutoTxInfo(ctx, id)
	fmt.Printf("Reward Type: %v\n", rewardType)

	owner, err := sdk.AccAddressFromBech32(txInfo.Owner)
	if err != nil {
		return err
	}
	//airdrop reward hooks
	if rewardType == 3 {
		k.hooks.AfterAutoTxAuthz(ctx, owner)
	} else if rewardType == 1 {
		k.hooks.AfterAutoTxWasm(ctx, owner)
	}

	txInfo.AutoTxHistory[len(txInfo.AutoTxHistory)-1].ExecutedOnHost = true
	k.SetAutoTxInfo(ctx, &txInfo)

	return nil
}

// checks if dependent transactions have executed on the host chain
func (k Keeper) AllowedToExecute(ctx sdk.Context, autoTx *types.AutoTxInfo) bool {
	//check if dependent tx executions succeeded
	for _, autoTxId := range autoTx.DependsOnTxIds {
		autoTxInfo := k.GetAutoTxInfo(ctx, autoTxId)
		if len(autoTxInfo.AutoTxHistory) == 0 {
			return true
		}
		if !autoTx.AutoTxHistory[len(autoTx.AutoTxHistory)-1].ExecutedOnHost {
			// we could reinsert the entry into the queue if desired
			// if autoTx.AutoTxHistory[len(autoTx.AutoTxHistory)-1].Retries <= autoTx.MaxRetries {
			// 	k.InsertAutoTxQueue(ctx, autoTx.TxID, autoTx.ExecTime)
			// }
			return false
		}
	}
	return true
}

// getTmpAutoTxID for a certain port and sequence
func (k Keeper) getTmpAutoTxID(ctx sdk.Context, portID string, seq uint64) uint64 {
	store := ctx.KVStore(k.storeKey)
	autoTxIDBz := store.Get(append((append(types.TmpAutoTxIDLatestTX, []byte(portID)...)), types.GetBytesForUint(seq)...))

	return types.GetIDFromBytes(autoTxIDBz)
}
func (k Keeper) setTmpAutoTxID(ctx sdk.Context, autoTxID uint64, portID string, seq uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(append((append(types.TmpAutoTxIDLatestTX, []byte(portID)...)), types.GetBytesForUint(seq)...), types.GetBytesForUint(autoTxID))
}
