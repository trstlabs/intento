package keeper

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	//tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

func TestSendLocalTx(t *testing.T) {

	ctx, keepers, addr1, _, addr2, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))

	autoTxAddr, _ := CreateFakeFundedAccount(ctx, keepers.accountKeeper, keepers.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))

	types.Denom = "stake"

	localMsg := &banktypes.MsgSend{
		FromAddress: addr2.String(),
		ToAddress:   addr1.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
	}
	anys, err := types.PackTxMsgAnys([]sdk.Msg{localMsg})
	require.NoError(t, err)

	autoTxInfo := types.AutoTxInfo{
		TxID: 0, Owner: addr2.String(), FeeAddress: autoTxAddr.String(), Msgs: anys, Interval: time.Second * 20, StartTime: time.Now().Add(time.Hour * -1), EndTime: time.Now().Add(time.Second * 20), PortID: "", ConnectionID: "",
	}

	err, executedLocally := keepers.SendAutoTx(ctx, autoTxInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)

}

func TestSendLocalTxAutoCompound(t *testing.T) {

	ctx, keepers, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))

	autoTxAddr, _ := CreateFakeFundedAccount(ctx, keepers.accountKeeper, keepers.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))

	types.Denom = "stake"

	val := keepers.stakingKeeper.GetAllValidators(ctx)[0]
	require.NotEmpty(t, val)
	val.Tokens = sdk.NewInt(5000)
	val.DelegatorShares = sdk.NewDecFromInt(val.Tokens)
	val.Commission = stakingtypes.NewCommission(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	keepers.stakingKeeper.SetValidator(ctx, val)

	//setting baseline
	keepers.distrKeeper.SetValidatorHistoricalRewards(ctx, val.GetOperator(), 2, distrtypes.ValidatorHistoricalRewards{
		CumulativeRewardRatio: sdk.DecCoins{},
		ReferenceCount:        2,
	})
	keepers.distrKeeper.SetValidatorCurrentRewards(ctx, val.GetOperator(), distrtypes.ValidatorCurrentRewards{
		Rewards: sdk.DecCoins{},
		Period:  3,
	})
	count := keepers.distrKeeper.GetValidatorHistoricalReferenceCount(ctx)
	require.Equal(t, uint64(2), count)
	rewards := keepers.distrKeeper.GetValidatorCurrentRewards(ctx, val.GetOperator())
	require.Equal(t, uint64(3), rewards.Period)

	newShares, err := keepers.stakingKeeper.Delegate(ctx, delAddr, sdk.NewInt(77), stakingtypes.Unbonded, val, true)
	require.NoError(t, err)
	require.Equal(t, newShares, sdk.NewDec(77))

	decCoins := sdk.NewDecCoins(sdk.NewDecCoin("stake", sdk.NewInt(6666)))
	keepers.distrKeeper.AllocateTokensToValidator(ctx, val, decCoins)
	keepers.distrKeeper.SetValidatorCurrentRewards(ctx, val.GetOperator(), distrtypes.NewValidatorCurrentRewards(decCoins, 3))
	/* endingPeriod := */ keepers.distrKeeper.IncrementValidatorPeriod(ctx, val)
	ctx = nextStakingBlocks(ctx, keepers.stakingKeeper, 1)

	keepers.distrKeeper.SetValidatorHistoricalRewards(ctx, val.GetOperator(), 3, distrtypes.ValidatorHistoricalRewards{
		CumulativeRewardRatio: decCoins,
		ReferenceCount:        2,
	})

	rewards = keepers.distrKeeper.GetValidatorCurrentRewards(ctx, val.GetOperator())
	require.Equal(t, uint64(4), rewards.Period)

	count = keepers.distrKeeper.GetValidatorHistoricalReferenceCount(ctx)
	require.Equal(t, uint64(2), count)

	keepers.distrKeeper.SetValidatorCurrentRewards(ctx, val.GetOperator(), distrtypes.ValidatorCurrentRewards{
		Rewards: decCoins,
		Period:  4,
	})

	autoTxInfo := createLocalAutoTxInfo(delAddr, val, autoTxAddr)
	err, executedLocally := keepers.SendAutoTx(ctx, autoTxInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)

	delegations := keepers.stakingKeeper.GetAllDelegatorDelegations(ctx, delAddr)
	require.Greater(t, delegations[0].Shares.TruncateInt64(), sdk.NewDec(77).TruncateInt64())

}

func createLocalAutoTxInfo(addr2 sdk.AccAddress, val stakingtypes.Validator, autoTxAddr sdk.AccAddress) types.AutoTxInfo {
	localMsg := &distrtypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: addr2.String(),
		ValidatorAddress: val.GetOperator().String(),
	}
	anys, _ := types.PackTxMsgAnys([]sdk.Msg{localMsg})
	autoTxInfo := types.AutoTxInfo{
		TxID: 0, Owner: addr2.String(), FeeAddress: autoTxAddr.String(), Msgs: anys, Interval: time.Second * 20, StartTime: time.Now().Add(time.Hour * -1), EndTime: time.Now().Add(time.Second * 20), PortID: "", ConnectionID: "",
	}
	return autoTxInfo
}

// this will commit the current set, update the block height and set historic info
// basically, letting blocks pass
func nextStakingBlocks(ctx sdk.Context, stakingKeeper stakingkeeper.Keeper, count int) sdk.Context {
	// for i := 0; i < count; i++ {
	staking.EndBlocker(ctx, &stakingKeeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	staking.BeginBlocker(ctx, &stakingKeeper)
	return ctx
	// }

}
