package keeper

import (
	"testing"
	"time"

	math "cosmossdk.io/math"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	//tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/trstlabs/intento/x/intent/types"
)

func newFakeMsgWithdrawDelegatorReward(delegator sdk.AccAddress, validator stakingtypes.Validator) *distrtypes.MsgWithdrawDelegatorReward {
	msgWithdrawDelegatorReward := &distrtypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: delegator.String(),
		//ValidatorAddress: validator.GetOperator().String(),
		ValidatorAddress: validator.GetOperator(),
	}
	return msgWithdrawDelegatorReward
}

func newFakeMsgDelegate(delegator sdk.AccAddress, validator stakingtypes.Validator) *stakingtypes.MsgDelegate {
	MsgDelegate := &stakingtypes.MsgDelegate{
		DelegatorAddress: delegator.String(),
		ValidatorAddress: validator.GetOperator(),
		Amount:           sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000)),
	}
	return MsgDelegate
}

func newFakeMsgSend(fromAddr sdk.AccAddress, toAddr sdk.AccAddress) *banktypes.MsgSend {
	msgSend := &banktypes.MsgSend{
		FromAddress: fromAddr.String(),
		ToAddress:   toAddr.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))),
	}
	return msgSend
}

func TestSendLocalTx(t *testing.T) {
	ctx, keepers, addr1, _, addr2, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))

	flowAddr, _ := CreateFakeFundedAccount(ctx, keepers.accountKeeper, keepers.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))

	types.Denom = "stake"

	localMsg := newFakeMsgSend(addr1, addr2)
	anys, err := types.PackTxMsgAnys([]sdk.Msg{localMsg})
	require.NoError(t, err)

	flowInfo := createBaseFlowInfo(addr1, flowAddr)
	flowInfo.Msgs = anys

	executedLocally, msgResponses, err := keepers.TriggerFlow(ctx, &flowInfo)
	require.NoError(t, err)
	require.NotNil(t, msgResponses)
	require.True(t, executedLocally)
}

func TestSendLocalTxAutocompound(t *testing.T) {
	ctx, keeper, _, _, delAddr, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))

	flowAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))

	types.Denom = "stake"

	// Set baseline
	val, ctx := delegateTokens(t, ctx, keeper, delAddr)

	flowInfo := createBaseFlowInfo(delAddr, flowAddr)
	msgWithdrawDelegatorReward := newFakeMsgWithdrawDelegatorReward(delAddr, val)
	msgDelegate := newFakeMsgDelegate(delAddr, val)
	flowInfo.Msgs, _ = types.PackTxMsgAnys([]sdk.Msg{msgWithdrawDelegatorReward, msgDelegate})
	flowInfo.Conditions = &types.ExecutionConditions{FeedbackLoops: []*types.FeedbackLoop{{ResponseIndex: 0, ResponseKey: "Amount.[0].Amount", MsgsIndex: 1, MsgKey: "Amount", ValueType: "sdk.Int"}}}
	delegations, _ := keeper.stakingKeeper.GetAllDelegatorDelegations(ctx, delAddr)
	require.Equal(t, delegations[0].Shares.TruncateInt64(), math.LegacyNewDec(77).TruncateInt64())
	keeper.HandleFlow(ctx, ctx.Logger(), flowInfo, time.Now(), nil)

	history, _ := keeper.GetFlowHistory(ctx, flowInfo.ID)

	delegations, _ = keeper.stakingKeeper.GetAllDelegatorDelegations(ctx, delAddr)
	require.Greater(t, delegations[0].Shares.TruncateInt64(), math.LegacyNewDec(77).TruncateInt64())

	require.Equal(t, len(history[0].MsgResponses), 2)

	///also test feedbackloop response via IBC handling
	keeper.setTmpFlowID(ctx, flowInfo.ID, "port-1", "channel-1", 0)
	err := keeper.HandleResponseAndSetFlowResult(ctx, "port-1", "channel-1", delAddr, 0, history[0].MsgResponses)
	require.NoError(t, err)
	history, _ = keeper.GetFlowHistory(ctx, flowInfo.ID)
	require.Equal(t, len(history[0].MsgResponses), 4)
}

func delegateTokens(t *testing.T, ctx sdk.Context, keepers Keeper, delAddr sdk.AccAddress) (stakingtypes.Validator, sdk.Context) {
	vals, _ := keepers.stakingKeeper.GetAllValidators(ctx)
	require.NotEmpty(t, vals)
	val := vals[0]
	val.Tokens = math.NewInt(5000)
	val.DelegatorShares = math.LegacyNewDecFromInt(val.Tokens)
	valAddr, _ := sdk.ValAddressFromBech32(val.OperatorAddress)
	val.Commission = stakingtypes.NewCommission(math.LegacyNewDecWithPrec(5, 1), math.LegacyNewDecWithPrec(5, 1), math.LegacyNewDec(0))

	keepers.stakingKeeper.SetValidator(ctx, val)

	keepers.distrKeeper.SetValidatorHistoricalRewards(ctx, valAddr, 2, distrtypes.ValidatorHistoricalRewards{
		CumulativeRewardRatio: sdk.DecCoins{},
		ReferenceCount:        2,
	})
	keepers.distrKeeper.SetValidatorCurrentRewards(ctx, valAddr, distrtypes.ValidatorCurrentRewards{
		Rewards: sdk.DecCoins{},
		Period:  3,
	})
	count := keepers.distrKeeper.GetValidatorHistoricalReferenceCount(ctx)
	require.Equal(t, uint64(2), count)
	rewards, _ := keepers.distrKeeper.GetValidatorCurrentRewards(ctx, valAddr)
	require.Equal(t, uint64(3), rewards.Period)

	newShares, err := keepers.stakingKeeper.Delegate(ctx, delAddr, math.NewInt(77), stakingtypes.Unbonded, val, true)
	require.NoError(t, err)
	require.Equal(t, newShares, math.LegacyNewDec(77))

	decCoins := sdk.NewDecCoins(sdk.NewDecCoin("stake", math.NewInt(6666)))
	keepers.distrKeeper.AllocateTokensToValidator(ctx, val, decCoins)
	keepers.distrKeeper.SetValidatorCurrentRewards(ctx, valAddr, distrtypes.NewValidatorCurrentRewards(decCoins, 3))
	keepers.distrKeeper.IncrementValidatorPeriod(ctx, val)
	ctx = nextStakingBlock(ctx, keepers.stakingKeeper)

	keepers.distrKeeper.SetValidatorHistoricalRewards(ctx, valAddr, 3, distrtypes.ValidatorHistoricalRewards{
		CumulativeRewardRatio: decCoins,
		ReferenceCount:        2,
	})

	rewards, _ = keepers.distrKeeper.GetValidatorCurrentRewards(ctx, valAddr)
	require.Equal(t, uint64(4), rewards.Period)

	count = keepers.distrKeeper.GetValidatorHistoricalReferenceCount(ctx)
	require.Equal(t, uint64(2), count)

	keepers.distrKeeper.SetValidatorCurrentRewards(ctx, valAddr, distrtypes.ValidatorCurrentRewards{
		Rewards: decCoins,
		Period:  4,
	})
	return val, ctx
}

func createBaseFlowInfo(ownerAddr sdk.AccAddress, flowAddr sdk.AccAddress) types.FlowInfo {
	flowInfo := types.FlowInfo{
		ID:            1,
		Owner:         ownerAddr.String(),
		FeeAddress:    flowAddr.String(),
		Msgs:          []*cdctypes.Any{},
		Interval:      time.Second * 20,
		StartTime:     time.Now().Add(time.Hour * -1),
		EndTime:       time.Now().Add(time.Second * 20),
		ICAConfig:     &types.ICAConfig{},
		Configuration: &types.ExecutionConfiguration{SaveResponses: true},
	}
	return flowInfo
}

// This will commit the current set, update the block height, and set historic info
// Basically, it lets blocks pass
func nextStakingBlock(ctx sdk.Context, stakingKeeper stakingkeeper.Keeper) sdk.Context {
	// for i := 0; i < count; i++ {
	stakingKeeper.EndBlocker(ctx)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)
	stakingKeeper.BeginBlocker(ctx)
	//staking.BeginBlocker(ctx, &stakingKeeper)
	return ctx
	// }
}
