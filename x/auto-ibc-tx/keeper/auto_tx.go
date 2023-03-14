package keeper

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	icatypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v4/modules/core/24-host"
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

	//check if autoTx is local
	if autoTxInfo.ConnectionID == "" {
		txMsgs := autoTxInfo.GetTxMsgs()
		return handleLocalAutoTx(k, ctx, txMsgs, autoTxInfo)
	}
	//if message contains ICA_ADDR, the ICA address is retreived and parsed
	txMsgs, err := k.replaceTextInMsg(ctx, autoTxInfo)
	if err != nil {
		return err
	}
	data, err := icatypes.SerializeCosmosTx(k.cdc, txMsgs)
	if err != nil {
		return err
	}
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}

	channelID, found := k.icaControllerKeeper.GetActiveChannelID(ctx, autoTxInfo.ConnectionID, autoTxInfo.PortID)
	if !found {
		return sdkerrors.Wrapf(icatypes.ErrActiveChannelNotFound, "failed to retrieve active channel for port %s", autoTxInfo.PortID)
	}

	chanCap, found := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(autoTxInfo.PortID, channelID))
	if !found {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	timeoutTimestamp := ctx.BlockTime().Add(time.Minute).UnixNano()
	//to ensure timeout does not result in channel closing
	if autoTxInfo.Interval > time.Minute*2 {
		timeoutTimestamp = ctx.BlockTime().Add(autoTxInfo.Interval).UnixNano()
	}

	sequence, err := k.icaControllerKeeper.SendTx(ctx, chanCap, autoTxInfo.ConnectionID, autoTxInfo.PortID, packetData, uint64(timeoutTimestamp))
	if err != nil {
		return err
	}

	k.setTmpAutoTxID(ctx, autoTxInfo.TxID, autoTxInfo.PortID, sequence)
	return nil
}

func handleLocalAutoTx(k Keeper, ctx sdk.Context, txMsgs []sdk.Msg, autoTxInfo types.AutoTxInfo) error {
	for _, msg := range txMsgs {
		handler := k.msgRouter.Handler(msg)
		for _, acct := range msg.GetSigners() {
			if acct.String() != autoTxInfo.Owner {
				return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "owner doesn't have permission to send this message")
			}
		}
		res, err := handler(ctx, msg)
		if err != nil {
			return err
		}

		//autocompound
		if sdk.MsgTypeURL(msg) == "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward" {
			validator := ""
			amount := sdk.NewCoin(types.Denom, sdk.ZeroInt())
			for _, ev := range res.Events {
				if ev.Type == distrtypes.EventTypeWithdrawRewards {
					for _, attr := range ev.Attributes {
						fmt.Printf("event %v\n", string(attr.Key))
						if string(attr.Key) == distrtypes.AttributeKeyValidator {
							validator = string(attr.Value)
						}
						if string(attr.Key) == sdk.AttributeKeyAmount {
							amount, err = sdk.ParseCoinNormalized(string(attr.Value))
							if err != nil {
								return err
							}
						}
					}

					msgDelegate := stakingtypes.MsgDelegate{DelegatorAddress: autoTxInfo.Owner, ValidatorAddress: validator, Amount: amount}
					handler := k.msgRouter.Handler(&msgDelegate)
					_, err = handler(ctx, &msgDelegate)
					if err != nil {
						return err
					}
				}
			}

		}

	}
	return nil
}

func (k Keeper) CreateAutoTx(ctx sdk.Context, owner sdk.AccAddress, label string, portID string, msgs []*cdctypes.Any, connectionId string, duration time.Duration, interval time.Duration, startAt time.Time, feeFunds sdk.Coins /*  retries uint64,  */, dependsOn []uint64) error {

	txID := k.autoIncrementID(ctx, types.KeyLastTxID)
	autoTxAddress, err := k.createFeeAccount(ctx, txID, owner, feeFunds)
	if err != nil {
		return err
	}

	endTime, execTime, interval := k.calculateAndInsertQueue(ctx, startAt, duration, txID, interval)

	autoTx := types.AutoTxInfo{
		TxID:         txID,
		Owner:        owner.String(),
		Label:        label,
		FeeAddress:   autoTxAddress.String(),
		Msgs:         msgs,
		Interval:     interval,
		StartTime:    startAt,
		ExecTime:     execTime,
		EndTime:      endTime,
		PortID:       portID,
		ConnectionID: connectionId,
		//MaxRetries:     retries,
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
	autoTxAddress := k.generateAutoTxFeeAddress(ctx, txID)
	existingAcct := k.accountKeeper.GetAccount(ctx, autoTxAddress)
	if existingAcct != nil {
		return nil, sdkerrors.Wrap(types.ErrAccountExists, existingAcct.GetAddress().String())
	}

	// deposit initial autoTx funds
	if !feeFunds.IsZero() && !feeFunds[0].Amount.IsZero() {
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
func (k Keeper) generateAutoTxFeeAddress(ctx sdk.Context, txID uint64) sdk.AccAddress {
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
	if startTime.After(time.Now().Add(time.Second * 30)) {
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

// SetAutoTxResult sets the result of the last executed TxID set at SendAutoTx.
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

	txInfo.AutoTxHistory[len(txInfo.AutoTxHistory)-1].Executed = true
	k.SetAutoTxInfo(ctx, &txInfo)

	return nil
}

// SetAutoTxOnTimeout sets the AutoTx timeout result to the AutoTx
func (k Keeper) SetAutoTxOnTimeout(ctx sdk.Context, sourcePort string, seq uint64) error {
	id := k.getTmpAutoTxID(ctx, sourcePort, seq)
	if id <= 0 {
		return nil
	}

	k.Logger(ctx).Debug("auto_tx_timeout", "id", id)

	txInfo := k.GetAutoTxInfo(ctx, id)

	txInfo.AutoTxHistory[len(txInfo.AutoTxHistory)-1].TimedOut = true
	k.SetAutoTxInfo(ctx, &txInfo)

	return nil
}

// SetAutoTxOnTimeout sets the AutoTx timeout result to the AutoTx
func (k Keeper) SetAutoTxError(ctx sdk.Context, sourcePort string, seq uint64, err string) error {
	id := k.getTmpAutoTxID(ctx, sourcePort, seq)
	if id <= 0 {
		return nil
	}

	k.Logger(ctx).Debug("auto_tx_error", "id", id)

	txInfo := k.GetAutoTxInfo(ctx, id)

	txInfo.AutoTxHistory[len(txInfo.AutoTxHistory)-1].Error = err
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
		if !autoTx.AutoTxHistory[len(autoTx.AutoTxHistory)-1].Executed {
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

func (k Keeper) replaceTextInMsg(ctx sdk.Context, autoTxInfo types.AutoTxInfo) (sdkMsgs []sdk.Msg, err error) {
	var txMsgs []sdk.Msg
	for _, message := range autoTxInfo.Msgs {
		var txMsg sdk.Msg
		err := k.cdc.UnpackAny(message, &txMsg)
		if err != nil {
			return nil, err
		}
		txMsgs = append(txMsgs, txMsg)
	}
	for _, msg := range txMsgs {

		// Marshal the message into a JSON string
		msgJSON, err := k.cdc.MarshalInterfaceJSON(msg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal %s message containing ICA_ADDR placeholder", msg)
		}
		msgJSONString := string(msgJSON)
		icaAddrToParse := "ICA_ADDR"
		index := strings.Index(msgJSONString, icaAddrToParse)
		if index == -1 {
			return txMsgs, nil
		}

		ica, found := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, autoTxInfo.ConnectionID, autoTxInfo.PortID)
		if !found {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrKeyNotFound, "ICA address not found")
		}

		// Replace the text "ICA_ADDR" in the JSON string
		msgJSONString = strings.ReplaceAll(msgJSONString, icaAddrToParse, ica)
		// Unmarshal the modified JSON string back into a message
		var updatedMsg sdk.Msg
		err = k.cdc.UnmarshalInterfaceJSON([]byte(msgJSONString), &updatedMsg)
		if err != nil {
			return nil, err
		}
		sdkMsgs = append(sdkMsgs, updatedMsg)

	}

	anys, err := types.PackTxMsgAnys(sdkMsgs)
	if err != nil {
		return nil, err
	}
	autoTxInfo.Msgs = anys
	k.SetAutoTxInfo(ctx, &autoTxInfo)

	return sdkMsgs, nil
}
