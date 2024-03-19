package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	proto "github.com/cosmos/gogoproto/proto"
	msgregistry "github.com/trstlabs/intento/x/intent/msg_registry"
	"github.com/trstlabs/intento/x/intent/types"
)

func (k Keeper) HandleRelayerReward(ctx sdk.Context, relayer sdk.AccAddress, rewardType int) {
	if !k.GetRelayerRewardsAvailability(ctx) {
		return
	}
	p := k.GetParams(ctx)

	// fmt.Printf("p.RelayerRewards %v\n", p.RelayerRewards)
	incentiveCoin := sdk.NewCoin(types.Denom, sdk.NewInt(p.RelayerRewards[rewardType]))

	// fmt.Printf("relayer %v\n", relayer.String())
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, relayer, sdk.NewCoins(incentiveCoin))
	if err != nil {
		//set incentives unavailable
		k.SetRelayerRewardsAvailability(ctx, false)
	}

}

func (k Keeper) SetRelayerRewardsAvailability(ctx sdk.Context, rewardsAvailable bool) {
	store := ctx.KVStore(k.storeKey)
	value := []byte("false")
	if rewardsAvailable {
		value = []byte("true")
	}
	store.Set([]byte(types.KeyRelayerRewardsAvailability), value)
}

// GetRelayerRewardsAvailability returns the rewards availability
func (k Keeper) GetRelayerRewardsAvailability(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	value := store.Get([]byte(types.KeyRelayerRewardsAvailability))
	return string(value) == "true"
}

// GetActionIbcUsage
func (k Keeper) TryGetActionIbcUsage(ctx sdk.Context, owner string) (types.ActionIbcUsage, error) {
	store := ctx.KVStore(k.storeKey)
	var autoTxUsage types.ActionIbcUsage
	autoTxBz := store.Get(append(types.ActionIbcUsageKeyPrefix, []byte(owner)...))

	err := k.cdc.Unmarshal(autoTxBz, &autoTxUsage)
	if err != nil {
		return types.ActionIbcUsage{}, err
	}
	return autoTxUsage, nil
}

func (k Keeper) SetActionIbcUsage(ctx sdk.Context, autoTxUsage *types.ActionIbcUsage) {
	store := ctx.KVStore(k.storeKey)
	store.Set(append(types.ActionIbcUsageKeyPrefix, []byte(autoTxUsage.Address)...), k.cdc.MustMarshal(autoTxUsage))
}

func (k Keeper) UpdateActionIbcUsage(ctx sdk.Context, autoTx types.ActionInfo) {
	for _, msg := range autoTx.Msgs {
		if msg.TypeUrl != sdk.MsgTypeURL(&authztypes.MsgExec{}) {
			return
		}

		msgExec := &authztypes.MsgExec{}
		if err := proto.Unmarshal(msg.Value, msgExec); err != nil {
			return
		}

		for _, msgInMsgExec := range msgExec.Msgs {
			var coin sdk.Coin

			switch msgInMsgExec.TypeUrl {
			case sdk.MsgTypeURL(&banktypes.MsgSend{}):
				{
					msgValue := &banktypes.MsgSend{}
					if err := proto.Unmarshal(msgInMsgExec.Value, msgValue); err != nil {
						return
					}
					coin = msgValue.Amount[0]
				}

			case sdk.MsgTypeURL(&msgregistry.MsgExecuteContract{}):
				{
					msgValue := &msgregistry.MsgExecuteContract{}
					if err := proto.Unmarshal(msgInMsgExec.Value, msgValue); err != nil {
						return
					}
					coin = msgValue.Funds[0]
				}

			case sdk.MsgTypeURL(&msgregistry.MsgSwapExactAmountIn{}):
				{
					msgValue := &msgregistry.MsgSwapExactAmountIn{}
					if err := proto.Unmarshal(msgInMsgExec.Value, msgValue); err != nil {
						return
					}
					coin = msgValue.TokenIn
				}

			case sdk.MsgTypeURL(&msgregistry.MsgSwapExactAmountOut{}):
				{
					msgValue := &msgregistry.MsgSwapExactAmountOut{}
					if err := proto.Unmarshal(msgInMsgExec.Value, msgValue); err != nil {
						return
					}
					coin = msgValue.TokenOut
				}

			default:
				return
			}

			k.appendToActionIbcUsage(ctx, autoTx.Owner, &types.ActionAck{
				Coin:         coin,
				ConnectionId: autoTx.ICAConfig.ConnectionID,
			})
		}
	}
}

func (k Keeper) appendToActionIbcUsage(ctx sdk.Context, owner string, autoTxAck *types.ActionAck) {
	autoIbcUsage, err := k.TryGetActionIbcUsage(ctx, owner)
	if err != nil {
		autoIbcUsage.Txs = append(autoIbcUsage.Txs, autoTxAck)
	} else {
		autoIbcUsage.Address = owner
		autoIbcUsage.Txs = []*types.ActionAck{autoTxAck}
	}
	k.SetActionIbcUsage(ctx, &autoIbcUsage)
}

func (k Keeper) IterateActionUsage(ctx sdk.Context) []types.ActionIbcUsage {
	// Get an instance of the KVStore for the given storeKey
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ActionIbcUsageKeyPrefix)
	iter := prefixStore.Iterator(nil, nil)

	// Defer closing the iterator until the function returns
	defer iter.Close()

	// Create a slice to hold the values
	var values []types.ActionIbcUsage

	// Loop over the iterator and append each value to the slice
	for ; iter.Valid(); iter.Next() {
		var c types.ActionIbcUsage
		k.cdc.MustUnmarshal(iter.Value(), &c)
		values = append(values, c)
	}

	// Return the slice of values
	return values
}
