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
	msg, _ := sdktypes.NewAnyWithValue(&types.MsgSubmitTx{})
	flow := types.Flow{
		ID: 0, Owner: ownerAddr.String(), FeeAddress: feeAddr.String(), Msgs: []*sdktypes.Any{msg}, StartTime: time.Now().Add(time.Hour * -1), EndTime: lastTime, ExecTime: lastTime,
	}

	acc, denom, err := keeper.GetFeeAccountForMinFees(ctx, flow, 1_000_000)
	require.Nil(t, err)
	fee, err := keeper.DistributeCoins(ctx, flow, acc, types.Denom)
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
	msg, _ := sdktypes.NewAnyWithValue(&types.MsgSubmitTx{})
	flow := types.Flow{
		ID: 0, Owner: ownerAddr.String(), FeeAddress: feeAddr.String(), Msgs: []*sdktypes.Any{msg}, StartTime: time.Now().Add(time.Hour * -1), EndTime: lastTime, ExecTime: lastTime, Configuration: &types.ExecutionConfiguration{WalletFallback: true},
	}

	//tokens from the owner will be used
	require.Equal(t, sdk.NewCoin(types.Denom, math.NewInt(0)), keeper.bankKeeper.GetBalance(ctx, feeAddr, types.Denom))

	ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
	acc, denom, err := keeper.GetFeeAccountForMinFees(ctx, flow, 1_000_000)
	require.Nil(t, err)
	require.NotEmpty(t, denom)
	require.Equal(t, acc, ownerAddr)
	fee, err := keeper.DistributeCoins(ctx, flow, acc, denom)
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

	msg, _ := sdktypes.NewAnyWithValue(&types.MsgSubmitTx{})
	flow := types.Flow{
		ID: 0, Owner: ownerAddr.String(), FeeAddress: feeAddr.String(), Msgs: []*sdktypes.Any{msg}, Interval: time.Second * 20, StartTime: time.Now().Add(time.Hour * -1), EndTime: time.Now().Add(time.Second * 20), SelfHostedICA: &types.ICAConfig{PortID: "ibccontoller-test", ConnectionID: "connection-0"},
	}

	ctx = ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
	acc, denom, err := keeper.GetFeeAccountForMinFees(ctx, flow, 1_000_000)
	require.Nil(t, err)
	require.NotEmpty(t, denom)
	require.NotEmpty(t, acc)

	types.Denom = sdk.DefaultBondDenom

	fee, err := keeper.DistributeCoins(ctx, flow, acc, denom)
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
	msg, _ := sdktypes.NewAnyWithValue(&types.MsgSubmitTx{})
	flow := types.Flow{
		ID: 0, Owner: ownerAddr.String(), FeeAddress: feeAddr.String(), Msgs: []*sdktypes.Any{msg}, StartTime: time.Now().Add(time.Hour * -1), EndTime: lastTime, ExecTime: lastTime,
	}

	// Get the validator rewards and check for no rewards yet
	val, _ := keeper.stakingKeeper.ValidatorByConsAddr(ctx, sdk.ConsAddress(ctx.BlockHeader().ProposerAddress))
	rewards, err := keeper.distrKeeper.GetValidatorCurrentRewards(ctx, sdk.ValAddress(val.GetOperator()))
	require.Nil(t, err)
	require.Equal(t, math.LegacyZeroDec(), rewards.Rewards.AmountOf("uabcd"))

	// Call the DistributeCoins method
	acc, denom, err := keeper.GetFeeAccountForMinFees(ctx, flow, 1_000_000)
	require.Nil(t, err)
	require.Equal(t, denom, "uabcd")
	fee, err := keeper.DistributeCoins(ctx, flow, acc, denom)
	require.Nil(t, err)

	// Check fee pool and ensure the community pool has the correct amount of coins
	feePool, err := keeper.distrKeeper.FeePool.Get(ctx)
	require.Nil(t, err)
	require.Equal(t, feePool.CommunityPool.AmountOf("uabcd").TruncateInt().String(), fee.Amount.String())

	// Check the owner's balance after fee distribution and burning
	require.Equal(t, sdk.NewInt64Coin(denom, 30_000_000), keeper.bankKeeper.GetBalance(ctx, ownerAddr, "uabcd").Add(fee))
}

func TestSendFeesToTrustlessAgentFeeAdmin(t *testing.T) {
	denom := "ibc/17409F270CB2FE874D5E3F339E958752DEC39319E5A44AD0399D2D1284AD499C"
	feeAmount := sdk.NewInt64Coin(denom, 5000)
	fullBalance := sdk.NewCoins(sdk.NewInt64Coin(denom, 10000))

	tests := []struct {
		name            string
		setup           func(*testing.T) (sdk.Context, Keeper, sdk.AccAddress, sdk.AccAddress, types.Flow, types.TrustlessAgent)
		expectErr       bool
		expectedBalance sdk.Coins
		expectFallback  bool
	}{
		{
			name: "fee address pays successfully",
			setup: func(t *testing.T) (sdk.Context, Keeper, sdk.AccAddress, sdk.AccAddress, types.Flow, types.TrustlessAgent) {
				ctx, keeper, _, _, _, _ := setupTest(t, nil)

				feeAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, fullBalance)
				ownerAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
				adminAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

				flow := types.Flow{
					Owner:      ownerAddr.String(),
					FeeAddress: feeAddr.String(),
					TrustlessAgent: &types.TrustlessAgentConfig{
						FeeLimit: sdk.NewCoins(feeAmount),
					},
					Configuration: &types.ExecutionConfiguration{WalletFallback: false},
				}

				agent := types.TrustlessAgent{
					FeeConfig: &types.TrustlessAgentFeeConfig{
						FeeCoinsSupported: sdk.NewCoins(feeAmount),
						FeeAdmin:          adminAddr.String(),
					},
				}

				return ctx, keeper, feeAddr, adminAddr, flow, agent
			},
			expectErr:       false,
			expectedBalance: sdk.NewCoins(feeAmount),
			expectFallback:  false,
		},
		{
			name: "fallback to owner",
			setup: func(t *testing.T) (sdk.Context, Keeper, sdk.AccAddress, sdk.AccAddress, types.Flow, types.TrustlessAgent) {
				ctx, keeper, _, _, _, _ := setupTest(t, nil)

				feeAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins()) // empty
				ownerAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
				adminAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

				// fund owner
				require.NoError(t, keeper.bankKeeper.MintCoins(ctx, types.ModuleName, fullBalance))
				require.NoError(t, keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ownerAddr, fullBalance))

				flow := types.Flow{
					Owner:      ownerAddr.String(),
					FeeAddress: feeAddr.String(),
					TrustlessAgent: &types.TrustlessAgentConfig{
						FeeLimit: sdk.NewCoins(feeAmount),
					},
					Configuration: &types.ExecutionConfiguration{WalletFallback: true},
				}

				agent := types.TrustlessAgent{
					FeeConfig: &types.TrustlessAgentFeeConfig{
						FeeCoinsSupported: sdk.NewCoins(feeAmount),
						FeeAdmin:          adminAddr.String(),
					},
				}

				return ctx, keeper, feeAddr, adminAddr, flow, agent
			},
			expectErr:       false,
			expectedBalance: sdk.NewCoins(feeAmount),
			expectFallback:  true,
		},
		{
			name: "insufficient funds and no fallback",
			setup: func(t *testing.T) (sdk.Context, Keeper, sdk.AccAddress, sdk.AccAddress, types.Flow, types.TrustlessAgent) {
				ctx, keeper, _, _, _, _ := setupTest(t, nil)

				// fee addr has less than required
				feeAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper,
					sdk.NewCoins(sdk.NewInt64Coin(denom, 1000)),
				)
				ownerAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
				adminAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

				flow := types.Flow{
					Owner:      ownerAddr.String(),
					FeeAddress: feeAddr.String(),
					TrustlessAgent: &types.TrustlessAgentConfig{
						FeeLimit: sdk.NewCoins(feeAmount),
					},
					Configuration: &types.ExecutionConfiguration{WalletFallback: false},
				}

				agent := types.TrustlessAgent{
					FeeConfig: &types.TrustlessAgentFeeConfig{
						FeeCoinsSupported: sdk.NewCoins(feeAmount),
						FeeAdmin:          adminAddr.String(),
					},
				}

				return ctx, keeper, feeAddr, adminAddr, flow, agent
			},
			expectErr:       true,
			expectedBalance: sdk.NewCoins(),
			expectFallback:  false,
		},
		{
			name: "multiple fee denoms paid from fee address",
			setup: func(t *testing.T) (sdk.Context, Keeper, sdk.AccAddress, sdk.AccAddress, types.Flow, types.TrustlessAgent) {
				ctx, keeper, _, _, _, _ := setupTest(t, nil)

				denomA := "ibc/AAA..."
				denomB := "ibc/BBB..."
				feeCoins := sdk.NewCoins(
					sdk.NewInt64Coin(denomA, 5000),
					sdk.NewInt64Coin(denomB, 2000),
				)

				// Fund fee address with both
				feeAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, feeCoins)
				ownerAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
				adminAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

				flow := types.Flow{
					Owner:      ownerAddr.String(),
					FeeAddress: feeAddr.String(),
					TrustlessAgent: &types.TrustlessAgentConfig{
						FeeLimit: feeCoins,
					},
					Configuration: &types.ExecutionConfiguration{WalletFallback: false},
				}

				agent := types.TrustlessAgent{
					FeeConfig: &types.TrustlessAgentFeeConfig{
						FeeCoinsSupported: feeCoins,
						FeeAdmin:          adminAddr.String(),
					},
				}

				return ctx, keeper, feeAddr, adminAddr, flow, agent
			},
			expectErr: false,
			expectedBalance: sdk.NewCoins(
				sdk.NewInt64Coin("ibc/AAA...", 5000),
			),
			expectFallback: false,
		},
		{
			name: "uinto set first and used first",
			setup: func(t *testing.T) (sdk.Context, Keeper, sdk.AccAddress, sdk.AccAddress, types.Flow, types.TrustlessAgent) {
				ctx, keeper, _, _, _, _ := setupTest(t, nil)

				denomA := "uinto"
				denomB := "ibc/BBB..."
				denomC := "ibc/CCC..."
				feeCoinsSupported := sdk.NewCoins(
					sdk.NewInt64Coin(denomA, 3456),
					sdk.NewInt64Coin(denomB, 1455),
					sdk.NewInt64Coin(denomC, 234),
				)

				// Fund fee address with both
				feeAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, feeCoinsSupported)
				ownerAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
				adminAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

				flow := types.Flow{
					Owner:      ownerAddr.String(),
					FeeAddress: feeAddr.String(),
					TrustlessAgent: &types.TrustlessAgentConfig{
						FeeLimit: sdk.NewCoins(
							sdk.NewInt64Coin(denomA, 5000),
							sdk.NewInt64Coin(denomB, 2000),
							sdk.NewInt64Coin(denomC, 1000),
						),
					},
					Configuration: &types.ExecutionConfiguration{WalletFallback: false},
				}

				agent := types.TrustlessAgent{
					FeeConfig: &types.TrustlessAgentFeeConfig{
						FeeCoinsSupported: feeCoinsSupported,
						FeeAdmin:          adminAddr.String(),
					},
				}

				return ctx, keeper, feeAddr, adminAddr, flow, agent
			},
			expectErr: false,
			expectedBalance: sdk.NewCoins(
				sdk.NewInt64Coin("uinto", 3456),
			),
			expectFallback: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, keeper, feeAddr, adminAddr, flow, agent := tc.setup(t)

			err := keeper.SendFeesToTrustlessAgentFeeAdmin(ctx, flow, agent)
			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			adminBal := keeper.bankKeeper.GetAllBalances(ctx, adminAddr)
			require.Equal(t, adminBal, tc.expectedBalance)

			if tc.expectFallback {
				feeBal := keeper.bankKeeper.GetAllBalances(ctx, feeAddr)
				require.True(t, feeBal.IsZero(), "fee address should not have been used in fallback")
			}
		})
	}
}

func TestGetFeeAccountForMinFees_WithMultipleBalanceDenoms(t *testing.T) {
	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins())

	feeAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(
		sdk.NewInt64Coin("ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9", 9970),
		sdk.NewInt64Coin("ibc/17409F270CB2FE874D5E3F339E958752DEC39319E5A44AD0399D2D1284AD499C", 38180),
		sdk.NewInt64Coin("ibc/92E0120F15D037353CFB73C14651FC8930ADC05B93100FD7754D3A689E53B333", 39290),
		sdk.NewInt64Coin("uinto", 983262761587),
	))
	ownerAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins())

	params := types.DefaultParams()
	params.BurnFeePerMsg = 10000
	params.GasFeeCoins = sdk.NewCoins(
		sdk.NewInt64Coin("ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9", 100),
		sdk.NewInt64Coin("ibc/17409F270CB2FE874D5E3F339E958752DEC39319E5A44AD0399D2D1284AD499C", 1),
		sdk.NewInt64Coin("ibc/92E0120F15D037353CFB73C14651FC8930ADC05B93100FD7754D3A689E53B333", 10),
		sdk.NewInt64Coin("uinto", 10),
	)

	params.FlowFlexFeeMul = 10
	_ = keeper.SetParams(ctx, params)
	msg, _ := sdktypes.NewAnyWithValue(&types.MsgSubmitTx{})
	flow := types.Flow{
		ID:         1,
		FeeAddress: feeAddr.String(),
		Owner:      ownerAddr.String(),
		Msgs:       []*sdktypes.Any{msg, msg, msg, msg, msg, msg, msg, msg, msg, msg},
	}

	expectedGas := uint64(1000000)
	acc, denom, err := keeper.GetFeeAccountForMinFees(ctx, flow, expectedGas)
	require.NoError(t, err)
	require.NotNil(t, acc)
	require.NotEmpty(t, denom)
	require.Equal(t, feeAddr.String(), acc.String())
	require.Equal(t, "ibc/17409F270CB2FE874D5E3F339E958752DEC39319E5A44AD0399D2D1284AD499C", denom)
}

func TestTotalBurnt(t *testing.T) {
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
	msg, _ := sdktypes.NewAnyWithValue(&types.MsgSubmitTx{})
	flow := types.Flow{
		ID: 0, Owner: ownerAddr.String(), FeeAddress: feeAddr.String(), Msgs: []*sdktypes.Any{msg}, StartTime: time.Now().Add(time.Hour * -1), EndTime: lastTime, ExecTime: lastTime,
	}

	// Initially, total burnt coins should be empty
	totalBurnt := keeper.GetTotalBurnt(ctx)
	require.True(t, totalBurnt == sdk.Coin{})

	acc, _, err := keeper.GetFeeAccountForMinFees(ctx, flow, 1_000_000)
	require.Nil(t, err)
	_, err = keeper.DistributeCoins(ctx, flow, acc, types.Denom)
	require.Nil(t, err)

	// After burning, total burnt coins should contain the burnt amount
	totalBurnt = keeper.GetTotalBurnt(ctx)
	expectedBurnt := sdk.NewCoin(types.Denom, math.NewInt(2_000_000))
	require.Equal(t, expectedBurnt, totalBurnt)

	// Run another flow to accumulate more burnt coins
	lastTime = time.Now().Add(time.Second * 40)
	flow2 := types.Flow{
		ID: 1, Owner: ownerAddr.String(), FeeAddress: feeAddr.String(), Msgs: []*sdktypes.Any{msg}, StartTime: time.Now().Add(time.Hour * -1), EndTime: lastTime, ExecTime: lastTime,
	}

	_, err = keeper.DistributeCoins(ctx, flow2, acc, types.Denom)
	require.Nil(t, err)

	// Total burnt coins should now be 4_000_000 (2_000_000 * 2)
	totalBurnt = keeper.GetTotalBurnt(ctx)
	expectedBurnt = sdk.NewCoin(types.Denom, math.NewInt(4_000_000))
	require.Equal(t, expectedBurnt, totalBurnt)
}
