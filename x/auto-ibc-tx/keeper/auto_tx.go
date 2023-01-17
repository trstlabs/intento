package keeper

import (
	"encoding/binary"
	"time"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
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
func (k Keeper) SetAutoTxInfo(ctx sdk.Context, autoTxID uint64, autoTx *types.AutoTxInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetAutoTxKey(autoTxID), k.cdc.MustMarshal(autoTx))
}

func (k Keeper) SendAutoTx(ctx sdk.Context, autoTxInfo types.AutoTxInfo) error {
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: autoTxInfo.Data,
	}

	_, err := k.icaControllerKeeper.SendTx(ctx, &capabilitytypes.Capability{Index: autoTxInfo.ChannelCapabilityIndex}, autoTxInfo.ConnectionID, autoTxInfo.PortID, packetData, ^uint64(0))
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) CreateAutoTx(ctx sdk.Context, owner sdk.AccAddress, portID string, data []byte, connectionId string, duration time.Duration, interval time.Duration, startAt time.Time, feeFunds sdk.Coins) error {

	channelID, found := k.icaControllerKeeper.GetActiveChannelID(ctx, connectionId, portID)
	if !found {
		return sdkerrors.Wrapf(icatypes.ErrActiveChannelNotFound, "failed to retrieve active channel for port %s", portID)
	}

	chanCap, found := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	if !found {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	txID := k.autoIncrementID(ctx, types.KeyLastTxID)
	autoTxAddress, err := k.createFeeAccount(ctx, txID, owner, feeFunds)
	if err != nil {
		return err
	}
	endTime, execTime, interval := k.calculateAndInsertQueue(ctx, startAt, duration, txID, interval)
	autoTx := types.AutoTxInfo{
		TxID:                   txID,
		Address:                autoTxAddress,
		Owner:                  owner,
		Data:                   data,
		Interval:               interval,
		Duration:               duration,
		StartTime:              startAt,
		ExecTime:               execTime,
		EndTime:                endTime,
		ChannelCapabilityIndex: chanCap.Index,
		PortID:                 portID,
	}

	k.SetAutoTxInfo(ctx, txID, &autoTx)
	k.addToAutoTxOwnerSecondaryIndex(ctx, owner /*  startAt, */, txID)
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

// addToAutoTxOwnerSecondaryIndex adds element to the index for autoTxs-by-creator queries
func (k Keeper) addToAutoTxOwnerSecondaryIndex(ctx sdk.Context, ownerAddress sdk.AccAddress, autoTxID uint64) error {
	store := ctx.KVStore(k.storeKey)
	/* timeBytes, err := startTime.MarshalBinary()
	if err != nil {
		return sdkerrors.Wrapf(types.ErrInvalid, "time bytes %s", string(timeBytes))
	} */
	store.Set(types.GetAutoTxByOwnerSecondaryIndexKey(ownerAddress, autoTxID), []byte{})
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
