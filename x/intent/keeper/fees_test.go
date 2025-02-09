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
	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(0))))

	feeAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 30_000_000)))

	keeper.SetParams(ctx, types.Params{
		FlowFundsCommission: 10,
		BurnFeePerMsg:       2_000_000, // 2INTO
		FlowFlexFeeMul:      10,
		GasFeeCoins:         sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1))),
		MaxFlowDuration:     time.Hour * 24 * 366 * 10, // a little over 10 years
		MinFlowDuration:     time.Second * 60,
		MinFlowInterval:     time.Second * 20,
	})

	pub2 := secp256k1.GenPrivKey().PubKey()
	ownerAddr := sdk.AccAddress(pub2.Address())
	types.Denom = sdk.DefaultBondDenom

	lastTime := time.Now().Add(time.Second * 20)
	flowInfo := types.FlowInfo{
		ID: 0, Owner: ownerAddr.String(), FeeAddress: feeAddr.String(), Msgs: NewMsg(), StartTime: time.Now().Add(time.Hour * -1), EndTime: lastTime, ExecTime: lastTime,
	}

	acc, denom, err := keeper.GetFeeAccountForMinFees(ctx, flowInfo, 1_000_000)
	require.Nil(t, err)
	fee, err := keeper.DistributeCoins(ctx, flowInfo, acc, types.Denom)
	require.Nil(t, err)

	// Check if the BurnFeePerMsg is being burned
	burnedAmount := sdk.NewCoin(types.Denom, math.NewInt(2_000_000))
	expectedCommunityFeeAmount := fee.Amount.Sub(burnedAmount.Amount)

	// Check fee pool
	feePool, err := keeper.distrKeeper.FeePool.Get(ctx)
	require.Nil(t, err)
	require.Equal(t, feePool.CommunityPool.AmountOf(types.Denom).TruncateInt().String(), expectedCommunityFeeAmount.String())

	// When not recurring the owner is returned the feeAddr tokens
	require.Equal(t, sdk.NewInt64Coin(denom, 30_000_000), keeper.bankKeeper.GetBalance(ctx, ownerAddr, sdk.DefaultBondDenom).Add(fee))
}

func TestDistributeCoinsOwnerFeeFallbackNotRecurring(t *testing.T) {
	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins())

	keeper.SetParams(ctx, types.Params{
		FlowFundsCommission: 2,
		BurnFeePerMsg:       1_000_000, // 1into
		FlowFlexFeeMul:      10,
		GasFeeCoins:         sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1))),
		MaxFlowDuration:     time.Hour * 24 * 366 * 10,
		MinFlowDuration:     time.Second * 60,
		MinFlowInterval:     time.Second * 20,
	})

	pub1 := secp256k1.GenPrivKey().PubKey()
	feeAddr := sdk.AccAddress(pub1.Address())
	ownerAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 30_000_000)))
	types.Denom = sdk.DefaultBondDenom
	lastTime := time.Now().Add(time.Second * 20)
	flowInfo := types.FlowInfo{
		ID: 0, Owner: ownerAddr.String(), FeeAddress: feeAddr.String(), Msgs: NewMsg(), StartTime: time.Now().Add(time.Hour * -1), EndTime: lastTime, ExecTime: lastTime, Configuration: &types.ExecutionConfiguration{FallbackToOwnerBalance: true},
	}

	//tokens from the owner will be used
	require.Equal(t, sdk.NewCoin(types.Denom, math.NewInt(0)), keeper.bankKeeper.GetBalance(ctx, feeAddr, types.Denom))

	ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
	acc, denom, err := keeper.GetFeeAccountForMinFees(ctx, flowInfo, 1_000_000)
	require.Nil(t, err)
	require.NotEmpty(t, denom)
	require.Equal(t, acc, ownerAddr)
	fee, err := keeper.DistributeCoins(ctx, flowInfo, acc, denom)
	require.Nil(t, err)

	burnedAmount := sdk.NewCoin(types.Denom, math.NewInt(1_000_000)) // assuming BurnFeePerMsg is burned
	expectedCommunityFeeAmount := fee.Amount.Sub(burnedAmount.Amount)

	// Check fee pool
	feePool, err := keeper.distrKeeper.FeePool.Get(ctx)
	require.Nil(t, err)

	// Check if the BurnFeePerMsg is being burned
	require.Equal(t, feePool.CommunityPool.AmountOf(types.Denom).TruncateInt().String(), expectedCommunityFeeAmount.String())

	// When not recurring the owner is returned the feeAddr tokens
	require.Equal(t, sdk.NewInt64Coin(denom, 30_000_000), keeper.bankKeeper.GetBalance(ctx, ownerAddr, sdk.DefaultBondDenom).Add(fee))

}

func TestDistributeCoinsEmptyFlowBalance(t *testing.T) {
	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins())
	feeAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)))

	keeper.SetParams(ctx, types.Params{
		FlowFundsCommission: 2,
		BurnFeePerMsg:       1_000_000, // 1into
		FlowFlexFeeMul:      100,
		GasFeeCoins:         sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1))),
		MaxFlowDuration:     time.Hour * 24 * 366 * 10,
		MinFlowDuration:     time.Second * 60,
		MinFlowInterval:     time.Second * 60,
	})

	pub2 := secp256k1.GenPrivKey().PubKey()

	ownerAddr := sdk.AccAddress(pub2.Address())
	types.Denom = sdk.DefaultBondDenom

	flowInfo := types.FlowInfo{
		ID: 0, Owner: ownerAddr.String(), FeeAddress: feeAddr.String(), Msgs: NewMsg(), Interval: time.Second * 20, StartTime: time.Now().Add(time.Hour * -1), EndTime: time.Now().Add(time.Second * 20), ICAConfig: &types.ICAConfig{PortID: "ibccontoller-test", ConnectionID: "connection-0"},
	}

	ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
	acc, denom, err := keeper.GetFeeAccountForMinFees(ctx, flowInfo, 1_000_000)
	require.Nil(t, err)
	require.Empty(t, denom)
	fee, err := keeper.DistributeCoins(ctx, flowInfo, acc, denom)
	require.Error(t, err)
	require.Empty(t, fee.Denom, fee.Amount)
}

func TestDistributeCoinsEmptyOwnerBalanceAndMultipliedFlexFee(t *testing.T) {
	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins())
	feeAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 300_000_000)))

	keeper.SetParams(ctx, types.Params{
		FlowFundsCommission: 2,
		BurnFeePerMsg:       1_000_000, // fixed burn fee
		FlowFlexFeeMul:      250,       // flex fee multiplier (2.5x)
		GasFeeCoins:         sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1))),
		MaxFlowDuration:     time.Hour * 24 * 366 * 10,
		MinFlowDuration:     time.Second * 60,
		MinFlowInterval:     time.Second * 20,
	})

	pub2 := secp256k1.GenPrivKey().PubKey()
	ownerAddr := sdk.AccAddress(pub2.Address())

	flowInfo := types.FlowInfo{
		ID: 0, Owner: ownerAddr.String(), FeeAddress: feeAddr.String(), Msgs: NewMsg(), Interval: time.Second * 20, StartTime: time.Now().Add(time.Hour * -1), EndTime: time.Now().Add(time.Second * 20), ICAConfig: &types.ICAConfig{PortID: "ibccontoller-test", ConnectionID: "connection-0"},
	}

	ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
	acc, denom, err := keeper.GetFeeAccountForMinFees(ctx, flowInfo, 1_000_000)
	require.Nil(t, err)
	require.NotEmpty(t, denom)
	require.NotEmpty(t, acc)

	types.Denom = sdk.DefaultBondDenom

	fee, err := keeper.DistributeCoins(ctx, flowInfo, acc, denom)
	require.Nil(t, err)

	burnedAmount := sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1_000_000)) // assuming BurnFeePerMsg is burned
	expectedCommunityFeeAmount := fee.Amount.Sub(burnedAmount.Amount)

	// Check fee pool
	feePool, err := keeper.distrKeeper.FeePool.Get(ctx)
	require.Nil(t, err)
	require.Equal(t, feePool.CommunityPool.AmountOf(sdk.DefaultBondDenom).TruncateInt().String(), expectedCommunityFeeAmount.String())

	// feeAddr tokens +fee = start balance
	require.Equal(t, sdk.NewInt64Coin(denom, 300_000_000), keeper.bankKeeper.GetBalance(ctx, feeAddr, sdk.DefaultBondDenom).Add(fee))
}

func TestDistributeCoinsDifferentDenom(t *testing.T) {
	// Setup context and keeper with an initial coin balance of a different denom
	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin("uabcd", math.NewInt(1_000_000))))

	feeAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("uabcd", 30_000_000)))

	// Set params
	keeper.SetParams(ctx, types.Params{
		FlowFundsCommission: 10,
		BurnFeePerMsg:       1_000_000, // 1INTO
		FlowFlexFeeMul:      10,
		GasFeeCoins:         sdk.NewCoins(sdk.NewCoin("uabcd", math.NewInt(1))),
		MaxFlowDuration:     time.Hour * 24 * 366 * 10, // a little over 10 years
		MinFlowDuration:     time.Second * 60,
		MinFlowInterval:     time.Second * 20,
	})

	// Generate a new public key for the owner and feeAddr
	pub2 := secp256k1.GenPrivKey().PubKey()
	ownerAddr := sdk.AccAddress(pub2.Address())

	// Simulate the flow info for testing
	lastTime := time.Now().Add(time.Second * 20)
	flowInfo := types.FlowInfo{
		ID: 0, Owner: ownerAddr.String(), FeeAddress: feeAddr.String(), Msgs: NewMsg(), StartTime: time.Now().Add(time.Hour * -1), EndTime: lastTime, ExecTime: lastTime,
	}

	// Get the validator rewards and check for no rewards yet
	val, _ := keeper.stakingKeeper.ValidatorByConsAddr(ctx, sdk.ConsAddress(ctx.BlockHeader().ProposerAddress))
	rewards, err := keeper.distrKeeper.GetValidatorCurrentRewards(ctx, sdk.ValAddress(val.GetOperator()))
	require.Nil(t, err)
	require.Equal(t, math.LegacyZeroDec(), rewards.Rewards.AmountOf("uabcd"))

	// Call the DistributeCoins method
	acc, denom, err := keeper.GetFeeAccountForMinFees(ctx, flowInfo, 1_000_000)
	require.Nil(t, err)
	require.Equal(t, denom, "uabcd")
	fee, err := keeper.DistributeCoins(ctx, flowInfo, acc, denom)
	require.Nil(t, err)

	// Check fee pool and ensure the community pool has the correct amount of coins
	feePool, err := keeper.distrKeeper.FeePool.Get(ctx)
	require.Nil(t, err)
	require.Equal(t, feePool.CommunityPool.AmountOf("uabcd").TruncateInt().String(), fee.Amount.String())

	// Check the owner's balance after fee distribution and burning
	require.Equal(t, sdk.NewInt64Coin(denom, 30_000_000), keeper.bankKeeper.GetBalance(ctx, ownerAddr, "uabcd").Add(fee))
}

func NewMsg() []*sdktypes.Any {
	msg, _ := sdktypes.NewAnyWithValue(&types.MsgSubmitTx{})
	return []*sdktypes.Any{msg}
}
