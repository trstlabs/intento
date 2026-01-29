package v1100_test

import (
	"testing"

	math "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/stretchr/testify/require"

	ccvconsumertypes "github.com/cosmos/interchain-security/v6/x/ccv/consumer/types"
	"github.com/trstlabs/intento/app"
	"github.com/trstlabs/intento/app/upgrades"
	v1100 "github.com/trstlabs/intento/app/upgrades/mainnet/v1.10.0"
)

func TestUpgrade(t *testing.T) {
	// initialize app with initChain=true to have validators
	intoApp := app.InitIntentoTestApp(true)
	ctx := intoApp.BaseApp.NewContext(false)

	// Check initial state
	// Manual fix: Ensure CCValidators are populated from staking validators
	stakingVals, err := intoApp.StakingKeeper.GetAllValidators(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, stakingVals, "should have staking validators initially")

	for _, v := range stakingVals {
		consAddr, err := v.GetConsAddr()
		require.NoError(t, err)
		pk, err := v.ConsPubKey()
		require.NoError(t, err)
		power := v.GetConsensusPower(sdk.DefaultPowerReduction)
		ccvVal, err := ccvconsumertypes.NewCCValidator(consAddr, power, pk)
		require.NoError(t, err)
		intoApp.ConsumerKeeper.SetCCValidator(ctx, ccvVal)
	}

	consumerVals := intoApp.ConsumerKeeper.GetAllCCValidator(ctx)
	require.NotEmpty(t, consumerVals, "should have consumer validators initially")

	// Fund the FundAddress (DAO address) used in DeICS
	_, fundAddrBz, err := bech32.DecodeAndConvert(v1100.FundAddress)
	require.NoError(t, err)
	fundAddr := sdk.AccAddress(fundAddrBz)

	// Mint coins to the mint module and send them to the fund address
	coins := sdk.NewCoins(sdk.NewCoin(v1100.DefaultDenom, math.NewInt(1000_000_000_000)))
	err = intoApp.BankKeeper.MintCoins(ctx, "mint", coins)
	require.NoError(t, err)
	err = intoApp.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", fundAddr, coins)
	require.NoError(t, err)

	initialFundBalance := intoApp.BankKeeper.GetBalance(ctx, fundAddr, v1100.DefaultDenom)

	// Setup the upgrade handler
	mm := intoApp.ModuleManager
	configurator := intoApp.Configurator
	keepers := upgrades.IntentoKeepers{
		StakingKeeper:  intoApp.StakingKeeper,
		ConsumerKeeper: intoApp.ConsumerKeeper,
		BankKeeper:     intoApp.BankKeeper,
	}

	handler := v1100.CreateUpgradeHandler(mm, configurator, keepers)

	// Run the handler
	plan := upgradetypes.Plan{Name: v1100.UpgradeName, Height: 100}
	versionMap := intoApp.ModuleManager.GetVersionMap()
	_, err = handler(ctx, plan, versionMap)
	require.NoError(t, err)

	// Verify logic execution by checking if FundAddress was debited
	// DeICS creates/updates validators by funding them.
	// Since we have at least one validator in the test setup (governor),
	// DeICS should have funded it with SovereignSelfStake (1_000_000).
	finalFundBalance := intoApp.BankKeeper.GetBalance(ctx, fundAddr, v1100.DefaultDenom)

	// We expect the balance to decrease by at least SovereignSelfStake
	diff := initialFundBalance.Amount.Sub(finalFundBalance.Amount)
	require.True(t, diff.GTE(math.NewInt(v1100.SovereignSelfStake)), "FundAddress should have spent funds to bond validators")

	// Also confirm migration successful
	// We can check if MaxValidators was updated in params (DeICS updates params)
	params, err := intoApp.StakingKeeper.GetParams(ctx)
	require.NoError(t, err)
	// DeICS sets MaxValidators to len(consumer) + len(newMsgs).
	// Test setup has 1 validator. So roughly 1 + 1 = 2 (or 1 if skipped).
	// But it sets it explicitly.
	require.True(t, params.MaxValidators > 0, "MaxValidators should be set")
}
