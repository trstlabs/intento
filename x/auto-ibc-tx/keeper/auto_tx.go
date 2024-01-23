package keeper

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"

	errorsmod "cosmossdk.io/errors"
	"github.com/cometbft/cometbft/crypto"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
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

func (k Keeper) SendAutoTx(ctx sdk.Context, autoTxInfo *types.AutoTxInfo) (error, bool, []*cdctypes.Any) {

	//check if autoTx is local
	if autoTxInfo.ConnectionID == "" {
		txMsgs := autoTxInfo.GetTxMsgs(k.cdc)
		err, msgResponses := handleLocalAutoTx(k, ctx, txMsgs, *autoTxInfo)
		return err, err == nil, msgResponses
	}

	//if message contains ICA_ADDR, the ICA address is retrieved and parsed
	txMsgs, err := k.parseAndSetMsgs(ctx, autoTxInfo)
	if err != nil {
		return err, false, nil
	}
	data, err := icatypes.SerializeCosmosTx(k.cdc, txMsgs)
	if err != nil {
		return err, false, nil
	}
	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}

	relativeTimeoutTimestamp := uint64(time.Minute.Nanoseconds())
	//idea: we can add logic to ensure timeout does not result in channel closing before next execution
	// if autoTxInfo.Interval > time.Minute*2 {
	// 	relativeTimeoutTimestamp = uint64((autoTxInfo.Interval - time.Minute).Nanoseconds())
	// }

	msgServer := icacontrollerkeeper.NewMsgServerImpl(&k.icaControllerKeeper)
	icaMsg := icacontrollertypes.NewMsgSendTx(autoTxInfo.Owner, autoTxInfo.ConnectionID, relativeTimeoutTimestamp, packetData)

	res, err := msgServer.SendTx(ctx, icaMsg)
	if err != nil {
		return err, false, nil
	}

	k.Logger(ctx).Debug("auto_tx", "ibc_sequence", res.Sequence)
	k.setTmpAutoTxID(ctx, autoTxInfo.TxID, autoTxInfo.PortID, res.Sequence)
	return nil, false, nil
}

func handleLocalAutoTx(k Keeper, ctx sdk.Context, txMsgs []sdk.Msg, autoTxInfo types.AutoTxInfo) (error, []*cdctypes.Any) {
	// CacheContext returns a new context with the multi-store branched into a cached storage object
	// writeCache is called only if all msgs succeed, performing state transitions atomically
	var msgResponses []*cdctypes.Any

	cacheCtx, writeCache := ctx.CacheContext()
	for _, msg := range txMsgs {
		handler := k.msgRouter.Handler(msg)
		for _, acct := range msg.GetSigners() {
			if acct.String() != autoTxInfo.Owner {
				return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "owner doesn't have permission to send this message"), nil
			}
		}

		res, err := handler(cacheCtx, msg)
		if err != nil {
			return err, nil
		}

		msgResponses = append(msgResponses, res.MsgResponses...)
		//autocompound
		if sdk.MsgTypeURL(msg) == "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward" {
			validator := ""
			amount := sdk.NewCoin(types.Denom, sdk.ZeroInt())
			for _, ev := range res.Events {
				if ev.Type == distrtypes.EventTypeWithdrawRewards {
					for _, attr := range ev.Attributes {
						if string(attr.Key) == distrtypes.AttributeKeyValidator {
							validator = string(attr.Value)
						}
						if string(attr.Key) == sdk.AttributeKeyAmount {
							amount, err = sdk.ParseCoinNormalized(string(attr.Value))
							if err != nil {
								return err, nil
							}
						}
					}

					msgDelegate := stakingtypes.MsgDelegate{DelegatorAddress: autoTxInfo.Owner, ValidatorAddress: validator, Amount: amount}
					handler := k.msgRouter.Handler(&msgDelegate)
					_, err = handler(cacheCtx, &msgDelegate)
					if err != nil {
						return err, nil
					}
				}
			}

		}

	}
	writeCache()
	if !autoTxInfo.Configuration.SaveMsgResponses {
		msgResponses = nil
	}
	return nil, msgResponses
}

func (k Keeper) CreateAutoTx(ctx sdk.Context, owner sdk.AccAddress, label string, portID string, msgs []*cdctypes.Any, connectionId string, duration time.Duration, interval time.Duration, startAt time.Time, feeFunds sdk.Coins /*  retries uint64,  */, configuration types.ExecutionConfiguration) error {

	txID := k.autoIncrementID(ctx, types.KeyLastTxID)
	autoTxAddress, err := k.createFeeAccount(ctx, txID, owner, feeFunds)
	if err != nil {
		return err
	}

	endTime, execTime := k.calculateTimeAndInsertQueue(ctx, startAt, duration, txID, interval)

	autoTx := types.AutoTxInfo{
		TxID:          txID,
		Owner:         owner.String(),
		Label:         label,
		FeeAddress:    autoTxAddress.String(),
		Msgs:          msgs,
		Interval:      interval,
		StartTime:     startAt,
		ExecTime:      execTime,
		EndTime:       endTime,
		PortID:        portID,
		ConnectionID:  connectionId,
		Configuration: &configuration,
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
		return nil, errorsmod.Wrap(types.ErrAccountExists, existingAcct.GetAddress().String())
	}

	// deposit initial autoTx funds
	if !feeFunds.IsZero() && !feeFunds[0].Amount.IsZero() {
		if k.bankKeeper.BlockedAddr(owner) {
			return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "blocked address can not be used")
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

func (k Keeper) calculateTimeAndInsertQueue(ctx sdk.Context, startTime time.Time, duration time.Duration, autoTxID uint64, interval time.Duration) (time.Time, time.Time) {
	endTime, execTime := calculateEndAndExecTimes(ctx, startTime, duration, interval)
	k.InsertAutoTxQueue(ctx, autoTxID, execTime)

	return endTime, execTime
}

func calculateEndAndExecTimes(ctx sdk.Context, startTime time.Time, duration time.Duration, interval time.Duration) (time.Time, time.Time) {
	endTime := startTime.Add(duration)
	execTime := calculateExecTime(ctx, duration, interval, startTime)

	return endTime, execTime
}

func calculateExecTime(ctx sdk.Context, duration, interval time.Duration, startTime time.Time) time.Time {
	if startTime.After(ctx.BlockTime()) {
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
		return errorsmod.Wrapf(types.ErrDuplicate, "autoincrement id: %s", string(lastIDKey))
	}
	bz := sdk.Uint64ToBigEndian(val)
	store.Set(lastIDKey, bz)
	return nil
}

func (k Keeper) importAutoTxInfo(ctx sdk.Context, autoTxId uint64, autoTxInfo types.AutoTxInfo) error {

	store := ctx.KVStore(k.storeKey)
	key := types.GetAutoTxKey(autoTxId)
	if store.Has(key) {
		return errorsmod.Wrapf(types.ErrDuplicate, "duplicate code: %d", autoTxId)
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
		if cb(key) {
			return
		}
	}
}

// SetAutoTxResult sets the result of the last executed TxID set at SendAutoTx.
func (k Keeper) SetAutoTxResult(ctx sdk.Context, port string, rewardType int, seq uint64, msgResponses []*cdctypes.Any) error {
	id := k.getTmpAutoTxID(ctx, port, seq)
	if id <= 0 {
		return nil
	}

	k.Logger(ctx).Debug("auto_tx", "executed", "on host")

	autoTxInfo := k.GetAutoTxInfo(ctx, id)
	fmt.Println(autoTxInfo.TxID)
	k.UpdateAutoTxIbcUsage(ctx, autoTxInfo)
	owner, err := sdk.AccAddressFromBech32(autoTxInfo.Owner)
	if err != nil {
		return err
	}
	//airdrop reward hooks
	if rewardType == 3 {
		k.hooks.AfterAutoTxAuthz(ctx, owner)
	} else if rewardType == 1 {
		k.hooks.AfterAutoTxWasm(ctx, owner)
	}

	autoTxInfo.AutoTxHistory[len(autoTxInfo.AutoTxHistory)-1].Executed = true

	if autoTxInfo.Configuration.SaveMsgResponses {
		autoTxInfo.AutoTxHistory[len(autoTxInfo.AutoTxHistory)-1].MsgResponses = msgResponses
	}

	k.SetAutoTxInfo(ctx, &autoTxInfo)

	return nil
}

// SetAutoTxOnTimeout sets the AutoTx timeout result to the AutoTx
func (k Keeper) SetAutoTxOnTimeout(ctx sdk.Context, sourcePort string, seq uint64) error {
	id := k.getTmpAutoTxID(ctx, sourcePort, seq)
	if id <= 0 {
		return nil
	}

	k.Logger(ctx).Debug("auto_tx packet timed out", "auto_tx_id", id)

	txInfo := k.GetAutoTxInfo(ctx, id)

	txInfo.AutoTxHistory[len(txInfo.AutoTxHistory)-1].TimedOut = true
	k.SetAutoTxInfo(ctx, &txInfo)

	return nil
}

// SetAutoTxOnTimeout sets the AutoTx timeout result to the AutoTx
func (k Keeper) SetAutoTxError(ctx sdk.Context, sourcePort string, seq uint64, err string) {
	id := k.getTmpAutoTxID(ctx, sourcePort, seq)
	if id <= 0 {
		return
	}

	k.Logger(ctx).Debug("auto_tx", "id", id, "error", err)

	txInfo := k.GetAutoTxInfo(ctx, id)

	txInfo.AutoTxHistory[len(txInfo.AutoTxHistory)-1].Errors = append(txInfo.AutoTxHistory[len(txInfo.AutoTxHistory)-1].Errors, err)
	k.SetAutoTxInfo(ctx, &txInfo)

	return
}

// AllowedToExecute checks if execution conditons are met, e.g. if dependent transactions have executed on the host chain
// insert the next entry when execution has not happend yet
func (k Keeper) AllowedToExecute(ctx sdk.Context, autoTx types.AutoTxInfo) bool {
	allowedToExecute := true
	// shouldRecur := autoTx.ExecTime.Before(autoTx.EndTime) && autoTx.ExecTime.Add(autoTx.Interval).Before(autoTx.EndTime)
	// conditions := autoTx.Conditions

	// //check if dependent tx executions succeeded
	// for _, autoTxId := range conditions.StopOnSuccessOf {
	// 	dependentTx := k.GetAutoTxInfo(ctx, uint64(autoTxId))
	// 	if len(dependentTx.AutoTxHistory) != 0 {
	// 		success := dependentTx.AutoTxHistory[len(dependentTx.AutoTxHistory)-1].Executed && dependentTx.AutoTxHistory[len(dependentTx.AutoTxHistory)-1].Errors != nil
	// 		if !success {
	// 			allowedToExecute = false
	// 			shouldRecur = false
	// 		}
	// 	}
	// }

	// //check if dependent tx executions failed
	// for _, autoTxId := range conditions.StopOnFailureOf {
	// 	dependentTx := k.GetAutoTxInfo(ctx, uint64(autoTxId))
	// 	if len(dependentTx.AutoTxHistory) != 0 {
	// 		success := dependentTx.AutoTxHistory[len(dependentTx.AutoTxHistory)-1].Executed && dependentTx.AutoTxHistory[len(dependentTx.AutoTxHistory)-1].Errors != nil
	// 		if success {
	// 			allowedToExecute = false
	// 			shouldRecur = false
	// 		}
	// 	}
	// }

	// //check if dependent tx executions succeeded
	// for _, autoTxId := range conditions.skipOnFailureOf {
	// 	dependentTx := k.GetAutoTxInfo(ctx, uint64(autoTxId))
	// 	if len(dependentTx.AutoTxHistory) != 0 {
	// 		success := dependentTx.AutoTxHistory[len(dependentTx.AutoTxHistory)-1].Executed && dependentTx.AutoTxHistory[len(dependentTx.AutoTxHistory)-1].Errors != nil
	// 		if !success {
	// 			allowedToExecute = false
	// 		}
	// 	}
	// }

	// //check if dependent tx executions failed
	// for _, autoTxId := range conditions.skipOnSuccessOf {
	// 	dependentTx := k.GetAutoTxInfo(ctx, uint64(autoTxId))
	// 	if len(dependentTx.AutoTxHistory) != 0 {
	// 		success := dependentTx.AutoTxHistory[len(dependentTx.AutoTxHistory)-1].Executed && dependentTx.AutoTxHistory[len(dependentTx.AutoTxHistory)-1].Errors != nil
	// 		if success {
	// 			allowedToExecute = false
	// 		}
	// 	}
	// }

	// //if not allowed to execute, remove entry
	// if !allowedToExecute {
	// 	k.RemoveFromAutoTxQueue(ctx, autoTx)
	// 	//insert the next entry given a recurring tx
	// 	if shouldRecur {
	// 		// adding next execTime and a new entry into the queue based on interval
	// 		k.InsertAutoTxQueue(ctx, autoTx.TxID, autoTx.ExecTime.Add(autoTx.Interval))
	// 	}
	// }

	return allowedToExecute
}

// getTmpAutoTxID getds tmp AutoTxId for a certain port and sequence. This is used to set results and timeouts.
func (k Keeper) getTmpAutoTxID(ctx sdk.Context, portID string, seq uint64) uint64 {
	store := ctx.KVStore(k.storeKey)
	autoTxIDBz := store.Get(append((append(types.TmpAutoTxIDLatestTX, []byte(portID)...)), types.GetBytesForUint(seq)...))

	return types.GetIDFromBytes(autoTxIDBz)
}

func (k Keeper) setTmpAutoTxID(ctx sdk.Context, autoTxID uint64, portID string, seq uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(append((append(types.TmpAutoTxIDLatestTX, []byte(portID)...)), types.GetBytesForUint(seq)...), types.GetBytesForUint(autoTxID))
}

func (k Keeper) parseAndSetMsgs(ctx sdk.Context, autoTxInfo *types.AutoTxInfo) (protoMsgs []proto.Message, err error) {

	if len(autoTxInfo.AutoTxHistory) != 0 {
		txMsgs := autoTxInfo.GetTxMsgs(k.cdc)
		for _, msg := range txMsgs {
			protoMsgs = append(protoMsgs, msg)
		}
		return protoMsgs, nil
	}

	var txMsgs []sdk.Msg
	var parsedIcaAddr bool

	for _, msg := range autoTxInfo.Msgs {
		var txMsg sdk.Msg
		err := k.cdc.UnpackAny(msg, &txMsg)
		if err != nil {
			return nil, err
		}
		// Marshal the message into a JSON string
		msgJSON, err := k.cdc.MarshalInterfaceJSON(txMsg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal %s message", msg)
		}

		msgJSONString := string(msgJSON)

		index := strings.Index(msgJSONString, types.ParseICAValue)
		if index == -1 {
			protoMsgs = append(protoMsgs, txMsg)
			txMsgs = append(txMsgs, txMsg)
			continue
		}

		ica, found := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, autoTxInfo.ConnectionID, autoTxInfo.PortID)
		if !found {
			return nil, errorsmod.Wrapf(sdkerrors.ErrKeyNotFound, "ICA address not found")
		}

		// Replace the text "ICA_ADDR" in the JSON string
		msgJSONString = strings.ReplaceAll(msgJSONString, types.ParseICAValue, ica)
		// Unmarshal the modified JSON string back into a proto message
		var updatedMsg sdk.Msg
		err = k.cdc.UnmarshalInterfaceJSON([]byte(msgJSONString), &updatedMsg)
		if err != nil {
			return nil, err
		}
		protoMsgs = append(protoMsgs, updatedMsg)

		txMsgs = append(txMsgs, updatedMsg)
		parsedIcaAddr = true

	}

	if parsedIcaAddr {
		anys, err := types.PackTxMsgAnys(txMsgs)
		if err != nil {
			return nil, err
		}
		autoTxInfo.Msgs = anys
	}

	return protoMsgs, nil
}
