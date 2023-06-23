package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	proto "github.com/gogo/protobuf/proto"
	msgregistry "github.com/trstlabs/trst/x/auto-ibc-tx/msg_registry"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
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

// GetAutoTxIbcUsage
func (k Keeper) TryGetAutoTxIbcUsage(ctx sdk.Context, owner string) (types.AutoTxIbcUsage, error) {
	store := ctx.KVStore(k.storeKey)
	var autoTxUsage types.AutoTxIbcUsage
	autoTxBz := store.Get(append(types.AutoTxIbcUsageKeyPrefix, []byte(owner)...))

	err := k.cdc.Unmarshal(autoTxBz, &autoTxUsage)
	if err != nil {
		return types.AutoTxIbcUsage{}, err
	}
	return autoTxUsage, nil
}

func (k Keeper) SetAutoTxIbcUsage(ctx sdk.Context, autoTxUsage *types.AutoTxIbcUsage) {
	store := ctx.KVStore(k.storeKey)
	store.Set(append(types.AutoTxIbcUsageKeyPrefix, []byte(autoTxUsage.Address)...), k.cdc.MustMarshal(autoTxUsage))
}

func (k Keeper) UpdateAutoTxIbcUsage(ctx sdk.Context, autoTx types.AutoTxInfo) {
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

			k.appendToAutoTxIbcUsage(ctx, autoTx.Owner, &types.AutoIbcTxAck{
				Coin:         coin,
				ConnectionId: autoTx.ConnectionID,
			})
		}
	}
}

func (k Keeper) appendToAutoTxIbcUsage(ctx sdk.Context, owner string, autoTxAck *types.AutoIbcTxAck) {
	autoIbcUsage, err := k.TryGetAutoTxIbcUsage(ctx, owner)
	if err != nil {
		autoIbcUsage.Txs = append(autoIbcUsage.Txs, autoTxAck)
	} else {
		autoIbcUsage.Address = owner
		autoIbcUsage.Txs = []*types.AutoIbcTxAck{autoTxAck}
	}
	k.SetAutoTxIbcUsage(ctx, &autoIbcUsage)
}

func (k Keeper) IterateAutoTxUsage(ctx sdk.Context) []types.AutoTxIbcUsage {
	// Get an instance of the KVStore for the given storeKey
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.AutoTxIbcUsageKeyPrefix)
	iter := prefixStore.Iterator(nil, nil)

	// Defer closing the iterator until the function returns
	defer iter.Close()

	// Create a slice to hold the values
	var values []types.AutoTxIbcUsage

	// Loop over the iterator and append each value to the slice
	for ; iter.Valid(); iter.Next() {
		var c types.AutoTxIbcUsage
		k.cdc.MustUnmarshal(iter.Value(), &c)
		values = append(values, c)
	}

	// Return the slice of values
	return values
}
