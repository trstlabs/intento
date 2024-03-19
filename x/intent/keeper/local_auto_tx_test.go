package keeper

import (
	"testing"
	"time"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	//tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/trstlabs/intento/x/intent/types"
)

func newFakeMsgWithdrawDelegatorReward(delegator sdk.AccAddress, validator stakingtypes.Validator) *distrtypes.MsgWithdrawDelegatorReward {
	msgWithdrawDelegatorReward := &distrtypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: delegator.String(),
		ValidatorAddress: validator.GetOperator().String(),
	}
	return msgWithdrawDelegatorReward
}

func newFakeMsgSend(fromAddr sdk.AccAddress, toAddr sdk.AccAddress) *banktypes.MsgSend {
	msgSend := &banktypes.MsgSend{
		FromAddress: fromAddr.String(),
		ToAddress:   toAddr.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))),
	}
	return msgSend
}

func TestSendLocalTx(t *testing.T) {
	ctx, keepers, addr1, _, addr2, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))

	autoTxAddr, _ := CreateFakeFundedAccount(ctx, keepers.accountKeeper, keepers.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))

	types.Denom = "stake"

	localMsg := newFakeMsgSend(addr1, addr2)
	anys, err := types.PackTxMsgAnys([]sdk.Msg{localMsg})
	require.NoError(t, err)

	autoTxInfo := createBaseActionInfo(addr1, autoTxAddr)
	autoTxInfo.Msgs = anys

	err, executedLocally, msgResponses := keepers.SendAction(ctx, &autoTxInfo)
	require.NoError(t, err)
	require.NotNil(t, msgResponses)
	require.True(t, executedLocally)
}

func TestSendLocalTxAutoCompound(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))

	autoTxAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))

	types.Denom = "stake"

	// Set baseline
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)

	autoTxInfo := createBaseActionInfo(delAddr, autoTxAddr)
	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	autoTxInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward})

	err, executedLocally, _ := keeper.SendAction(ctx, &autoTxInfo)
	require.NoError(t, err)
	require.True(t, executedLocally)

	delegations := keeper.stakingKeeper.GetAllDelegatorDelegations(ctx, delAddr)
	require.Greater(t, delegations[0].Shares.TruncateInt64(), sdk.NewDec(77).TruncateInt64())
}

func delegateTokens(t *testing.T, ctx sdk.Context, keepers Keeper, delAddr sdk.AccAddress) (stakingtypes.Validator, sdk.Context) {
	val := keepers.stakingKeeper.GetAllValidators(ctx)[0]
	require.NotEmpty(t, val)

	val.Tokens = sdk.NewInt(5000)
	val.DelegatorShares = sdk.NewDecFromInt(val.Tokens)
	val.Commission = stakingtypes.NewCommission(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	keepers.stakingKeeper.SetValidator(ctx, val)

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
	keepers.distrKeeper.IncrementValidatorPeriod(ctx, val)
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
	return val, ctx
}

func createBaseActionInfo(ownerAddr sdk.AccAddress, autoTxAddr sdk.AccAddress) types.ActionInfo {
	autoTxInfo := types.ActionInfo{
		ID:            1,
		Owner:         ownerAddr.String(),
		FeeAddress:    autoTxAddr.String(),
		Msgs:          []*cdctypes.Any{},
		Interval:      time.Second * 20,
		StartTime:     time.Now().Add(time.Hour * -1),
		EndTime:       time.Now().Add(time.Second * 20),
		ICAConfig:     &types.ICAConfig{},
		Configuration: &types.ExecutionConfiguration{SaveMsgResponses: true},
	}
	return autoTxInfo
}

// This will commit the current set, update the block height, and set historic info
// Basically, it lets blocks pass
func nextStakingBlocks(ctx sdk.Context, stakingKeeper stakingkeeper.Keeper, count int) sdk.Context {
	// for i := 0; i < count; i++ {
	staking.EndBlocker(ctx, &stakingKeeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	staking.BeginBlocker(ctx, &stakingKeeper)
	return ctx
	// }
}
