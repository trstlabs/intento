package keeper

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdktypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/trstlabs/intento/x/intent/types"
)

func TestDistributeCoinsNotRecurring(t *testing.T) {

	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000))))

	feeAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 30_000_000)))

	keeper.SetParams(ctx, types.Params{
		FlowFundsCommission: 10,
		FlowConstantFee:     1_000_000, // 1trst
		FlowFlexFeeMul:      10,
		GasFeeCoins:         sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1))),
		MaxFlowDuration:     time.Hour * 24 * 366 * 10, // a little over 10 years
		MinFlowDuration:     time.Second * 60,
		MinFlowInterval:     time.Second * 20,
	})

	pub2 := secp256k1.GenPrivKey().PubKey()
	ownerAddr := sdk.AccAddress(pub2.Address())
	types.Denom = "stake"
	lastTime := time.Now().Add(time.Second * 20)
	flowInfo := types.FlowInfo{
		ID: 0, Owner: ownerAddr.String(), FeeAddress: feeAddr.String(), Msgs: NewMsg(), StartTime: time.Now().Add(time.Hour * -1), EndTime: lastTime, ExecTime: lastTime,
	}

	val, _ := keeper.stakingKeeper.ValidatorByConsAddr(ctx, sdk.ConsAddress(ctx.BlockHeader().ProposerAddress))
	rewards, err := keeper.distrKeeper.GetValidatorCurrentRewards(ctx, sdk.ValAddress(val.GetOperator()))
	require.Nil(t, err)

	require.Equal(t, math.LegacyZeroDec(), rewards.Rewards.AmountOf(sdk.DefaultBondDenom))
	acc, denom, err := keeper.GetFeeAccountForMinFees(ctx, flowInfo, 1_0000_000)
	require.Nil(t, err)
	fee, err := keeper.DistributeCoins(ctx, flowInfo, acc, types.Denom, ctx.BlockHeader().ProposerAddress)
	require.Nil(t, err)

	feePool, err := keeper.distrKeeper.FeePool.Get(ctx)
	require.Nil(t, err)
	require.Equal(t, feePool.CommunityPool.AmountOf(types.Denom).TruncateInt().String(), fee.Amount.String())
	require.Equal(t, sdk.NewInt64Coin(denom, 30_000_000), keeper.bankKeeper.GetBalance(ctx, ownerAddr, sdk.DefaultBondDenom).Add(fee))

}

func TestDistributeCoinsOwnerFeeFallbackLastExec(t *testing.T) {

	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins())

	keeper.SetParams(ctx, types.Params{
		FlowFundsCommission: 2,
		FlowConstantFee:     1_000_000, // 1trst
		FlowFlexFeeMul:      10,
		GasFeeCoins:         sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1))),
		MaxFlowDuration:     time.Hour * 24 * 366 * 10, // a little over 10 years
		MinFlowDuration:     time.Second * 60,
		MinFlowInterval:     time.Second * 20,
	})

	pub1 := secp256k1.GenPrivKey().PubKey()
	feeAddr := sdk.AccAddress(pub1.Address())
	ownerAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 30_000_000)))
	types.Denom = "stake"
	lastTime := time.Now().Add(time.Second * 20)
	flowInfo := types.FlowInfo{
		ID: 0, Owner: ownerAddr.String(), FeeAddress: feeAddr.String(), Msgs: NewMsg(), StartTime: time.Now().Add(time.Hour * -1), EndTime: lastTime, ExecTime: lastTime, Configuration: &types.ExecutionConfiguration{FallbackToOwnerBalance: true},
	}
	ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
	acc, denom, err := keeper.GetFeeAccountForMinFees(ctx, flowInfo, 1_0000_000)
	require.Nil(t, err)
	require.NotEmpty(t, denom)
	require.Equal(t, acc, ownerAddr)
	fee, err := keeper.DistributeCoins(ctx, flowInfo, acc, denom, ctx.BlockHeader().ProposerAddress)
	require.Nil(t, err)

	require.Equal(t, sdk.NewCoin(types.Denom, math.NewInt(0)), keeper.bankKeeper.GetBalance(ctx, feeAddr, types.Denom))

	feePool, err := keeper.distrKeeper.FeePool.Get(ctx)
	require.Nil(t, err)
	require.Equal(t, feePool.CommunityPool.AmountOf(types.Denom).TruncateInt().String(), fee.Amount.String())
	require.Equal(t, sdk.NewInt64Coin(denom, 30_000_000), keeper.bankKeeper.GetBalance(ctx, ownerAddr, types.Denom).Add(fee))

}

func TestDistributeCoinsEmptyFlowBalance(t *testing.T) {

	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins())
	feeAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))

	keeper.SetParams(ctx, types.Params{
		FlowFundsCommission: 2,
		FlowConstantFee:     1_000_000, // 1trst
		FlowFlexFeeMul:      100,
		GasFeeCoins:         sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1))),
		MaxFlowDuration:     time.Hour * 24 * 366 * 10, // a little over 10 years
		MinFlowDuration:     time.Second * 60,
		MinFlowInterval:     time.Second * 60,
	})

	pub2 := secp256k1.GenPrivKey().PubKey()

	ownerAddr := sdk.AccAddress(pub2.Address())
	types.Denom = "stake"

	flowInfo := types.FlowInfo{
		ID: 0, Owner: ownerAddr.String(), FeeAddress: feeAddr.String(), Msgs: NewMsg(), Interval: time.Second * 20, StartTime: time.Now().Add(time.Hour * -1), EndTime: time.Now().Add(time.Second * 20), ICAConfig: &types.ICAConfig{PortID: "ibccontoller-test", ConnectionID: "connection-0"},
	}

	ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
	_, denom, err := keeper.GetFeeAccountForMinFees(ctx, flowInfo, 1_0000_000)
	require.Nil(t, err)
	require.Empty(t, denom)
}

func TestDistributeCoinsEmptyOwnerBalanceAndMultipliedFlexFee(t *testing.T) {

	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins())
	feeAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 300_000_000)))
	keeper.SetParams(ctx, types.Params{
		FlowFundsCommission: 2,
		FlowConstantFee:     1_000_000, // 1trst
		FlowFlexFeeMul:      250,       // 250/100 = 2.5x
		GasFeeCoins:         sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1))),
		MaxFlowDuration:     time.Hour * 24 * 366 * 10, // a little over 10 years
		MinFlowDuration:     time.Second * 60,
		MinFlowInterval:     time.Second * 20,
	})

	pub2 := secp256k1.GenPrivKey().PubKey()

	ownerAddr := sdk.AccAddress(pub2.Address())
	types.Denom = "stake"

	flowInfo := types.FlowInfo{
		ID: 0, Owner: ownerAddr.String(), FeeAddress: feeAddr.String(), Msgs: NewMsg(), Interval: time.Second * 20, StartTime: time.Now().Add(time.Hour * -1), EndTime: time.Now().Add(time.Second * 20), ICAConfig: &types.ICAConfig{PortID: "ibccontoller-test", ConnectionID: "connection-0"},
	}

	ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
	acc, denom, err := keeper.GetFeeAccountForMinFees(ctx, flowInfo, 1_0000_000)
	require.Nil(t, err)
	require.NotEmpty(t, denom)
	require.NotEmpty(t, acc)
	fee, err := keeper.DistributeCoins(ctx, flowInfo, acc, denom, ctx.BlockHeader().ProposerAddress)
	require.Nil(t, err)

	feePool, err := keeper.distrKeeper.FeePool.Get(ctx)
	require.Nil(t, err)
	require.Equal(t, feePool.CommunityPool.AmountOf(types.Denom).TruncateInt().String(), fee.Amount.String())
	require.Equal(t, sdk.NewInt64Coin(denom, 300_000_000), keeper.bankKeeper.GetBalance(ctx, feeAddr, sdk.DefaultBondDenom).Add(fee))
	require.Equal(t, sdk.NewInt64Coin(denom, 0), keeper.bankKeeper.GetBalance(ctx, ownerAddr, sdk.DefaultBondDenom))

}

func NewMsg() []*sdktypes.Any {
	msg, _ := sdktypes.NewAnyWithValue(&types.MsgSubmitTx{})
	return []*sdktypes.Any{msg}
}
