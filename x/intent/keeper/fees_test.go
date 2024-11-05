package keeper

import (
	"testing"
	"time"

	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdktypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/trstlabs/intento/x/intent/types"
)

func TestDistributeCoinsNotRecurring(t *testing.T) {

	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))

	feeAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 30_000_000)))

	keeper.SetParams(ctx, types.Params{
		ActionFundsCommission: 10,
		ActionConstantFee:     1_000_000, // 1trst
		ActionFlexFeeMul:      10,
		GasFeeCoins:           sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1))),
		MaxActionDuration:     time.Hour * 24 * 366 * 10, // a little over 10 years
		MinActionDuration:     time.Second * 60,
		MinActionInterval:     time.Second * 20,
	})

	pub2 := secp256k1.GenPrivKey().PubKey()
	ownerAddr := sdk.AccAddress(pub2.Address())
	types.Denom = "stake"
	lastTime := time.Now().Add(time.Second * 20)
	actionInfo := types.ActionInfo{
		ID: 0, Owner: ownerAddr.String(), FeeAddress: feeAddr.String(), Msgs: NewMsg(), StartTime: time.Now().Add(time.Hour * -1), EndTime: lastTime, ExecTime: lastTime,
	}

	val := keeper.stakingKeeper.ValidatorByConsAddr(ctx, sdk.ConsAddress(ctx.BlockHeader().ProposerAddress))
	require.Equal(t, sdk.ZeroDec(), keeper.distrKeeper.GetValidatorCurrentRewards(ctx, val.GetOperator()).Rewards.AmountOf(sdk.DefaultBondDenom))
	acc, denom, err := keeper.GetFeeAccountForMinFees(ctx, actionInfo, 1_0000_000)
	require.Nil(t, err)
	fee, err := keeper.DistributeCoins(ctx, actionInfo, acc, types.Denom, ctx.BlockHeader().ProposerAddress)
	require.Nil(t, err)

	feePool := keeper.distrKeeper.GetFeePool(ctx).CommunityPool
	require.Equal(t, feePool.Sort().AmountOf(types.Denom).TruncateInt().String(), fee.Amount.String())
	require.Equal(t, sdk.NewInt64Coin(denom, 30_000_000), keeper.bankKeeper.GetBalance(ctx, ownerAddr, sdk.DefaultBondDenom).Add(fee))

}

func TestDistributeCoinsOwnerFeeFallbackLastExec(t *testing.T) {

	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins())

	keeper.SetParams(ctx, types.Params{
		ActionFundsCommission: 2,
		ActionConstantFee:     1_000_000, // 1trst
		ActionFlexFeeMul:      10,
		GasFeeCoins:           sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1))),
		MaxActionDuration:     time.Hour * 24 * 366 * 10, // a little over 10 years
		MinActionDuration:     time.Second * 60,
		MinActionInterval:     time.Second * 20,
	})

	pub1 := secp256k1.GenPrivKey().PubKey()
	feeAddr := sdk.AccAddress(pub1.Address())
	ownerAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 30_000_000)))
	types.Denom = "stake"
	lastTime := time.Now().Add(time.Second * 20)
	actionInfo := types.ActionInfo{
		ID: 0, Owner: ownerAddr.String(), FeeAddress: feeAddr.String(), Msgs: NewMsg(), StartTime: time.Now().Add(time.Hour * -1), EndTime: lastTime, ExecTime: lastTime, Configuration: &types.ExecutionConfiguration{FallbackToOwnerBalance: true},
	}
	ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	acc, denom, err := keeper.GetFeeAccountForMinFees(ctx, actionInfo, 1_0000_000)
	require.Nil(t, err)
	require.NotEmpty(t, denom)
	require.Equal(t, acc, ownerAddr)
	fee, err := keeper.DistributeCoins(ctx, actionInfo, acc, denom, ctx.BlockHeader().ProposerAddress)
	require.Nil(t, err)

	require.Equal(t, sdk.NewCoin(types.Denom, sdk.NewInt(0)), keeper.bankKeeper.GetBalance(ctx, feeAddr, types.Denom))

	feePool := keeper.distrKeeper.GetFeePool(ctx).CommunityPool
	require.Equal(t, feePool.Sort().AmountOf(types.Denom).TruncateInt(), fee.Amount)
	require.Equal(t, sdk.NewInt64Coin(denom, 30_000_000), keeper.bankKeeper.GetBalance(ctx, ownerAddr, types.Denom).Add(fee))

}

func TestDistributeCoinsEmptyActionBalance(t *testing.T) {

	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins())
	feeAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 0)))

	keeper.SetParams(ctx, types.Params{
		ActionFundsCommission: 2,
		ActionConstantFee:     1_000_000, // 1trst
		ActionFlexFeeMul:      100,
		GasFeeCoins:           sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1))),
		MaxActionDuration:     time.Hour * 24 * 366 * 10, // a little over 10 years
		MinActionDuration:     time.Second * 60,
		MinActionInterval:     time.Second * 60,
	})

	pub2 := secp256k1.GenPrivKey().PubKey()

	ownerAddr := sdk.AccAddress(pub2.Address())
	types.Denom = "stake"

	actionInfo := types.ActionInfo{
		ID: 0, Owner: ownerAddr.String(), FeeAddress: feeAddr.String(), Msgs: NewMsg(), Interval: time.Second * 20, StartTime: time.Now().Add(time.Hour * -1), EndTime: time.Now().Add(time.Second * 20), ICAConfig: &types.ICAConfig{PortID: "ibccontoller-test", ConnectionID: "connection-0"},
	}

	ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	_, denom, err := keeper.GetFeeAccountForMinFees(ctx, actionInfo, 1_0000_000)
	require.Nil(t, err)
	require.Empty(t, denom)
}

func TestDistributeCoinsEmptyOwnerBalanceAndMultipliedFlexFee(t *testing.T) {

	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins())
	feeAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 300_000_000)))
	keeper.SetParams(ctx, types.Params{
		ActionFundsCommission: 2,
		ActionConstantFee:     1_000_000, // 1trst
		ActionFlexFeeMul:      250,       // 250/100 = 2.5x
		GasFeeCoins:           sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1))),
		MaxActionDuration:     time.Hour * 24 * 366 * 10, // a little over 10 years
		MinActionDuration:     time.Second * 60,
		MinActionInterval:     time.Second * 20,
	})

	pub2 := secp256k1.GenPrivKey().PubKey()

	ownerAddr := sdk.AccAddress(pub2.Address())
	types.Denom = "stake"

	actionInfo := types.ActionInfo{
		ID: 0, Owner: ownerAddr.String(), FeeAddress: feeAddr.String(), Msgs: NewMsg(), Interval: time.Second * 20, StartTime: time.Now().Add(time.Hour * -1), EndTime: time.Now().Add(time.Second * 20), ICAConfig: &types.ICAConfig{PortID: "ibccontoller-test", ConnectionID: "connection-0"},
	}

	ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	acc, denom, err := keeper.GetFeeAccountForMinFees(ctx, actionInfo, 1_0000_000)
	require.Nil(t, err)
	require.NotEmpty(t, denom)
	require.NotEmpty(t, acc)
	fee, err := keeper.DistributeCoins(ctx, actionInfo, acc, denom, ctx.BlockHeader().ProposerAddress)
	require.Nil(t, err)

	feePool := keeper.distrKeeper.GetFeePool(ctx).CommunityPool
	require.Equal(t, feePool.Sort().AmountOf(types.Denom).TruncateInt(), fee.Amount)
	require.Equal(t, sdk.NewInt64Coin(denom, 300_000_000), keeper.bankKeeper.GetBalance(ctx, feeAddr, sdk.DefaultBondDenom).Add(fee))
	require.Equal(t, sdk.NewInt64Coin(denom, 0), keeper.bankKeeper.GetBalance(ctx, ownerAddr, sdk.DefaultBondDenom))

}

func NewMsg() []*sdktypes.Any {
	msg, _ := sdktypes.NewAnyWithValue(&types.MsgSubmitTx{})
	return []*sdktypes.Any{msg}
}
