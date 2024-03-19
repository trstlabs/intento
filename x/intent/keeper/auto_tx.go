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
	"github.com/trstlabs/intento/x/intent/types"
)

// GetActionInfo
func (k Keeper) GetActionInfo(ctx sdk.Context, autoID uint64) types.ActionInfo {
	store := ctx.KVStore(k.storeKey)
	var autoTx types.ActionInfo
	autoTxBz := store.Get(types.GetActionKey(autoID))

	k.cdc.MustUnmarshal(autoTxBz, &autoTx)
	return autoTx
}

// TryGetActionInfo
func (k Keeper) TryGetActionInfo(ctx sdk.Context, autoID uint64) (types.ActionInfo, error) {
	store := ctx.KVStore(k.storeKey)
	var autoTx types.ActionInfo
	autoTxBz := store.Get(types.GetActionKey(autoID))

	err := k.cdc.Unmarshal(autoTxBz, &autoTx)
	if err != nil {
		return types.ActionInfo{}, err
	}
	return autoTx, nil
}

func (k Keeper) SetActionInfo(ctx sdk.Context, autoTx *types.ActionInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetActionKey(autoTx.ID), k.cdc.MustMarshal(autoTx))
}

func (k Keeper) SendAction(ctx sdk.Context, autoTx *types.ActionInfo) (error, bool, []*cdctypes.Any) {
	//check if autoTx is local
	if autoTx.ICAConfig == nil || autoTx.ICAConfig.ConnectionID == "" {
		txMsgs := autoTx.GetTxMsgs(k.cdc)
		err, msgResponses := handleLocalAction(k, ctx, txMsgs, *autoTx)
		return err, err == nil, msgResponses
	}

	channelID, found := k.icaControllerKeeper.GetActiveChannelID(ctx, autoTx.ICAConfig.ConnectionID, autoTx.ICAConfig.PortID)
	if !found {
		return icatypes.ErrActiveChannelNotFound, false, nil
	}

	//if message contains ICA_ADDR, the ICA address is retrieved and parsed
	txMsgs, err := k.parseAndSetMsgs(ctx, autoTx)
	if err != nil {
		fmt.Printf("ERrrr")
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

	msgServer := icacontrollerkeeper.NewMsgServerImpl(&k.icaControllerKeeper)
	icaMsg := icacontrollertypes.NewMsgSendTx(autoTx.Owner, autoTx.ICAConfig.ConnectionID, relativeTimeoutTimestamp, packetData)

	res, err := msgServer.SendTx(ctx, icaMsg)
	if err != nil {
		return err, false, nil
	}

	k.Logger(ctx).Debug("action", "ibc_sequence", res.Sequence)
	k.setTmpActionID(ctx, autoTx.ID, autoTx.ICAConfig.PortID, channelID, res.Sequence)
	return nil, false, nil
}

func handleLocalAction(k Keeper, ctx sdk.Context, txMsgs []sdk.Msg, autoTx types.ActionInfo) (error, []*cdctypes.Any) {
	// CacheContext returns a new context with the multi-store branched into a cached storage object
	// writeCache is called only if all msgs succeed, performing state transitions atomically
	var msgResponses []*cdctypes.Any

	cacheCtx, writeCache := ctx.CacheContext()
	for _, msg := range txMsgs {
		// if sdk.MsgTypeURL(msg) == "/ibc.applications.transfer.v1.MsgTransfer" {
		// 	transferMsg, err := types.GetTransferMsg(k.cdc, autoTx.Msgs[index])
		// 	if err != nil {
		// 		return err, nil
		// 	}
		// 	_, err = k.transferKeeper.Transfer(ctx, &transferMsg)
		// 	if err != nil {
		// 		return err, nil
		// 	}
		// 	continue
		// }

		handler := k.msgRouter.Handler(msg)
		for _, acct := range msg.GetSigners() {
			if acct.String() != autoTx.Owner {
				return errorsmod.Wrap(sdkerrors.ErrUnauthorized, "owner doesn't have permission to send this message"), nil
			}
		}

		res, err := handler(cacheCtx, msg)
		if err != nil {
			return err, nil
		}

		msgResponses = append(msgResponses, res.MsgResponses...)
		//autocompound example
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

					msgDelegate := stakingtypes.MsgDelegate{DelegatorAddress: autoTx.Owner, ValidatorAddress: validator, Amount: amount}
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
	if !autoTx.Configuration.SaveMsgResponses {
		msgResponses = nil
	}
	return nil, msgResponses
}

func (k Keeper) CreateAction(ctx sdk.Context, owner sdk.AccAddress, label string, msgs []*cdctypes.Any, duration time.Duration, interval time.Duration, startAt time.Time, feeFunds sdk.Coins, configuration types.ExecutionConfiguration, portID string, connectionId string, hostConnectionId string) error {

	id := k.autoIncrementID(ctx, types.KeyLastID)
	autoTxAddress, err := k.createFeeAccount(ctx, id, owner, feeFunds)
	if err != nil {
		return err
	}

	endTime, execTime := k.calculateTimeAndInsertQueue(ctx, startAt, duration, id, interval)

	icaConfig := types.ICAConfig{}
	if portID != "" && connectionId != "" && hostConnectionId != "" {
		icaConfig = types.ICAConfig{
			PortID:           portID,
			ConnectionID:     connectionId,
			HostConnectionID: hostConnectionId,
		}
	}

	autoTx := types.ActionInfo{
		ID:            id,
		Owner:         owner.String(),
		Label:         label,
		FeeAddress:    autoTxAddress.String(),
		Msgs:          msgs,
		Interval:      interval,
		StartTime:     startAt,
		ExecTime:      execTime,
		EndTime:       endTime,
		ICAConfig:     &icaConfig,
		Configuration: &configuration,
	}

	k.SetActionInfo(ctx, &autoTx)
	k.addToActionOwnerIndex(ctx, owner, id)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAction,
			sdk.NewAttribute(types.AttributeKeyActionID, strconv.FormatUint(id, 10)),
		))
	return nil
}

func (k Keeper) createFeeAccount(ctx sdk.Context, id uint64, owner sdk.AccAddress, feeFunds sdk.Coins) (sdk.AccAddress, error) {
	autoTxAddress := k.generateActionFeeAddress(ctx, id)
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

// generates a autoTx address from id + instanceID
func (k Keeper) generateActionFeeAddress(ctx sdk.Context, id uint64) sdk.AccAddress {
	instanceID := k.autoIncrementID(ctx, types.KeyLastTxAddrID)
	return autoTxAddress(id, instanceID)
}

func autoTxAddress(id, instanceID uint64) sdk.AccAddress {
	// NOTE: It is possible to get a duplicate address if either id or instanceID
	// overflow 32 bits. This is highly improbable, but something that could be refactored.
	autoID := id<<32 + instanceID
	return addrFromUint64(autoID)

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

func (k Keeper) calculateTimeAndInsertQueue(ctx sdk.Context, startTime time.Time, duration time.Duration, autoID uint64, interval time.Duration) (time.Time, time.Time) {
	endTime, execTime := calculateEndAndExecTimes(ctx, startTime, duration, interval)
	k.InsertActionQueue(ctx, autoID, execTime)

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

func (k Keeper) importActionInfo(ctx sdk.Context, actionId uint64, autoTx types.ActionInfo) error {

	store := ctx.KVStore(k.storeKey)
	key := types.GetActionKey(actionId)
	if store.Has(key) {
		return errorsmod.Wrapf(types.ErrDuplicate, "duplicate code: %d", actionId)
	}
	// 0x01 | actionId (uint64) -> autoTx
	store.Set(key, k.cdc.MustMarshal(&autoTx))
	return nil
}

func (k Keeper) IterateActionInfos(ctx sdk.Context, cb func(uint64, types.ActionInfo) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ActionKeyPrefix)
	iter := prefixStore.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		var c types.ActionInfo
		k.cdc.MustUnmarshal(iter.Value(), &c)
		// cb returns true to stop early
		if cb(binary.BigEndian.Uint64(iter.Key()), c) {
			return
		}
	}
}

// addToActionOwnerIndex adds element to the index for autoTxs-by-creator queries
func (k Keeper) addToActionOwnerIndex(ctx sdk.Context, ownerAddress sdk.AccAddress, autoID uint64) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetActionByOwnerIndexKey(ownerAddress, autoID), []byte{})
}

// IterateActionsByOwner iterates over all autoTxs with given creator address in order of creation time asc.
func (k Keeper) IterateActionsByOwner(ctx sdk.Context, owner sdk.AccAddress, cb func(address sdk.AccAddress) bool) {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetActionsByOwnerPrefix(owner))
	for iter := prefixStore.Iterator(nil, nil); iter.Valid(); iter.Next() {
		key := iter.Key()
		if cb(key) {
			return
		}
	}
}

// SetActionResult sets the result of the last executed ID set at SendAction.
func (k Keeper) SetActionResult(ctx sdk.Context, portID string, channelID string, rewardType int, seq uint64, msgResponses []*cdctypes.Any) error {
	id := k.getTmpActionID(ctx, portID, channelID, seq)
	if id <= 0 {
		return nil
	}

	k.Logger(ctx).Debug("action", "executed", "on host")

	autoTx := k.GetActionInfo(ctx, id)

	k.UpdateActionIbcUsage(ctx, autoTx)
	owner, err := sdk.AccAddressFromBech32(autoTx.Owner)
	if err != nil {
		return err
	}
	//airdrop reward hooks
	if rewardType == 3 {
		k.hooks.AfterActionAuthz(ctx, owner)
	} else if rewardType == 1 {
		k.hooks.AfterActionWasm(ctx, owner)
	}

	autoTxHistory := k.GetActionHistory(ctx, id)

	autoTxHistory.History[len(autoTxHistory.History)-1].Executed = true

	if autoTx.Configuration.SaveMsgResponses {
		autoTxHistory.History[len(autoTxHistory.History)-1].MsgResponses = msgResponses
	}

	k.SetActionInfo(ctx, &autoTx)

	k.SetActionHistory(ctx, autoTx.ID, &autoTxHistory)

	return nil
}

// SetActionOnTimeout sets the Action timeout result to the Action

func (k Keeper) SetActionOnTimeout(ctx sdk.Context, sourcePort string, channelID string, seq uint64) error {
	id := k.getTmpActionID(ctx, sourcePort, channelID, seq)
	if id <= 0 {
		return nil
	}
	autoTx := k.GetActionInfo(ctx, id)
	if autoTx.Configuration.ReregisterICAAfterTimeout {
		autoTx := k.GetActionInfo(ctx, id)
		metadataString := icatypes.NewDefaultMetadataString(autoTx.ICAConfig.ConnectionID, autoTx.ICAConfig.HostConnectionID)
		err := k.RegisterInterchainAccount(ctx, autoTx.ICAConfig.ConnectionID, autoTx.Owner, metadataString)
		if err != nil {
			return err
		}
	} else {
		k.RemoveFromActionQueue(ctx, autoTx)
	}
	k.Logger(ctx).Debug("action packet timed out", "action_id", id)

	autoTxHistory := k.GetActionHistory(ctx, id)

	autoTxHistory.History[len(autoTxHistory.History)-1].TimedOut = true
	k.SetActionHistory(ctx, id, &autoTxHistory)

	return nil
}

// SetActionOnTimeout sets the Action timeout result to the Action
func (k Keeper) SetActionError(ctx sdk.Context, sourcePort string, channelID string, seq uint64, err string) {
	id := k.getTmpActionID(ctx, sourcePort, channelID, seq)
	if id <= 0 {
		return
	}

	k.Logger(ctx).Debug("action", "id", id, "error", err)

	autoTxHistory := k.GetActionHistory(ctx, id)

	autoTxHistory.History[len(autoTxHistory.History)-1].Errors = append(autoTxHistory.History[len(autoTxHistory.History)-1].Errors, err)
	k.SetActionHistory(ctx, id, &autoTxHistory)
}

// AllowedToExecute checks if execution conditons are met, e.g. if dependent transactions have executed on the host chain
// insert the next entry when execution has not happend yet
func (k Keeper) AllowedToExecute(ctx sdk.Context, autoTx types.ActionInfo) bool {
	allowedToExecute := true
	// shouldRecur := autoTx.ExecTime.Before(autoTx.EndTime) && autoTx.ExecTime.Add(autoTx.Interval).Before(autoTx.EndTime)
	// conditions := autoTx.Conditions

	// //check if dependent tx executions succeeded
	// for _, actionId := range conditions.StopOnSuccessOf {
	// 	dependentTx := k.GetActionInfo(ctx, uint64(actionId))
	// 	if len(dependentTx.ActionHistory) != 0 {
	// 		success := dependentTx.ActionHistory[len(dependentTx.ActionHistory)-1].Executed && dependentTx.ActionHistory[len(dependentTx.ActionHistory)-1].Errors != nil
	// 		if !success {
	// 			allowedToExecute = false
	// 			shouldRecur = false
	// 		}
	// 	}
	// }

	// //check if dependent tx executions failed
	// for _, actionId := range conditions.StopOnFailureOf {
	// 	dependentTx := k.GetActionInfo(ctx, uint64(actionId))
	// 	if len(dependentTx.ActionHistory) != 0 {
	// 		success := dependentTx.ActionHistory[len(dependentTx.ActionHistory)-1].Executed && dependentTx.ActionHistory[len(dependentTx.ActionHistory)-1].Errors != nil
	// 		if success {
	// 			allowedToExecute = false
	// 			shouldRecur = false
	// 		}
	// 	}
	// }

	// //check if dependent tx executions succeeded
	// for _, actionId := range conditions.skipOnFailureOf {
	// 	dependentTx := k.GetActionInfo(ctx, uint64(actionId))
	// 	if len(dependentTx.ActionHistory) != 0 {
	// 		success := dependentTx.ActionHistory[len(dependentTx.ActionHistory)-1].Executed && dependentTx.ActionHistory[len(dependentTx.ActionHistory)-1].Errors != nil
	// 		if !success {
	// 			allowedToExecute = false
	// 		}
	// 	}
	// }

	// //check if dependent tx executions failed
	// for _, actionId := range conditions.skipOnSuccessOf {
	// 	dependentTx := k.GetActionInfo(ctx, uint64(actionId))
	// 	if len(dependentTx.ActionHistory) != 0 {
	// 		success := dependentTx.ActionHistory[len(dependentTx.ActionHistory)-1].Executed && dependentTx.ActionHistory[len(dependentTx.ActionHistory)-1].Errors != nil
	// 		if success {
	// 			allowedToExecute = false
	// 		}
	// 	}
	// }

	// //if not allowed to execute, remove entry
	// if !allowedToExecute {
	// 	k.RemoveFromActionQueue(ctx, autoTx)
	// 	//insert the next entry given a recurring tx
	// 	if shouldRecur {
	// 		// adding next execTime and a new entry into the queue based on interval
	// 		k.InsertActionQueue(ctx, autoTx.ID, autoTx.ExecTime.Add(autoTx.Interval))
	// 	}
	// }

	return allowedToExecute
}

// getTmpActionID getds tmp actionId for a certain port and sequence. This is used to set results and timeouts.
func (k Keeper) getTmpActionID(ctx sdk.Context, portID string, channelID string, seq uint64) uint64 {
	store := ctx.KVStore(k.storeKey)
	// Append both portID and channelID to the key
	key := append(types.TmpActionIDLatestTX, []byte(portID)...)
	key = append(key, []byte(channelID)...)          // Append channelID after portID
	key = append(key, types.GetBytesForUint(seq)...) // Append sequence number

	autoIDBz := store.Get(key)

	return types.GetIDFromBytes(autoIDBz)
}

func (k Keeper) setTmpActionID(ctx sdk.Context, autoID uint64, portID string, channelID string, seq uint64) {
	store := ctx.KVStore(k.storeKey)
	// Append both portID and channelID to the key
	key := append(types.TmpActionIDLatestTX, []byte(portID)...)
	key = append(key, []byte(channelID)...)          // Append channelID after portID
	key = append(key, types.GetBytesForUint(seq)...) // Append sequence number

	store.Set(key, types.GetBytesForUint(autoID))
}

func (k Keeper) parseAndSetMsgs(ctx sdk.Context, autoTx *types.ActionInfo) (protoMsgs []proto.Message, err error) {
	fmt.Printf("denom %s\n", types.Denom)
	store := ctx.KVStore(k.storeKey)
	if store.Has(types.GetActionHistoryKey(autoTx.ID)) {
		txMsgs := autoTx.GetTxMsgs(k.cdc)
		for _, msg := range txMsgs {
			protoMsgs = append(protoMsgs, msg)
		}
		return protoMsgs, nil
	}

	var txMsgs []sdk.Msg
	var parsedIcaAddr bool

	for _, msg := range autoTx.Msgs {
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

		ica, found := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, autoTx.ICAConfig.ConnectionID, autoTx.ICAConfig.PortID)
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
		autoTx.Msgs = anys
	}
	fmt.Printf("parsed %v\n", autoTx.Msgs)
	return protoMsgs, nil
}
