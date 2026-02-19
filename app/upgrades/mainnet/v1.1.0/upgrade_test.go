package v1100_test

import (
	"testing"

	math "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/trstlabs/intento/app"
	"github.com/trstlabs/intento/app/upgrades"
	v1100 "github.com/trstlabs/intento/app/upgrades/mainnet/v1.1.0"
)

func TestUpgrade(t *testing.T) {
	// Set config for prefixes
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("into", "intopub")
	config.SetBech32PrefixForValidator("intovaloper", "intovaloperpub")
	config.SetBech32PrefixForConsensusNode("intovalcons", "intovalconspub")
	// config.Seal() // Don't seal to avoid panic if other tests run

	// initialize app with initChain=true to have validators
	intoApp := app.InitIntentoTestApp(true)
	ctx := intoApp.BaseApp.NewContext(false)

	// Fund the DAO address
	_, daoAddrBz, err := bech32.DecodeAndConvert(v1100.FundAddress)
	require.NoError(t, err)
	daoAddr := sdk.AccAddress(daoAddrBz)

	fundCoins := sdk.NewCoins(sdk.NewCoin(v1100.Denom, math.NewInt(100_000_000_000)))
	err = intoApp.BankKeeper.MintCoins(ctx, "mint", fundCoins)
	require.NoError(t, err)
	err = intoApp.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", daoAddr, fundCoins)
	require.NoError(t, err)

	initialDaoBalance := intoApp.BankKeeper.GetBalance(ctx, daoAddr, v1100.Denom)

	// Create a "ready" validator that matches one in the embedded JSON files
	// Address from validators/staking/sk8.json
	readyValoperStr := "intovaloper1mn86qd3w8050e065jq9xgsegc9yr6k34zamc9r"
	valAddr, err := sdk.ValAddressFromBech32(readyValoperStr)
	require.NoError(t, err)
	accAddr := sdk.AccAddress(valAddr)

	// Ensure this ready validator has 0 balance initially
	require.True(t, intoApp.BankKeeper.GetBalance(ctx, accAddr, v1100.Denom).IsZero())

	// Create the validator in staking keeper
	pk := ed25519.GenPrivKey().PubKey()
	pkAny, err := codectypes.NewAnyWithValue(pk)
	require.NoError(t, err)

	readyVal := stakingtypes.Validator{
		OperatorAddress: readyValoperStr,
		ConsensusPubkey: pkAny,
		Status:          stakingtypes.Bonded,
		Tokens:          math.NewInt(100),
		DelegatorShares: math.LegacyNewDec(100),
	}
	intoApp.StakingKeeper.SetValidator(ctx, readyVal)
	err = intoApp.StakingKeeper.SetValidatorByConsAddr(ctx, readyVal)
	require.NoError(t, err)
	// We don't need to set power index unless we care about voting power updates,
	// but DeICS iterates GetAllValidators which iterates the store by operator address, so SetValidator is key.

	// Also have the existing random validators from InitChain (ValA, ValB etc)
	// They won't be in the ready set.

	// Setup the upgrade handler
	mm := intoApp.ModuleManager
	configurator := intoApp.Configurator
	keepers := upgrades.IntentoKeepers{
		StakingKeeper:  intoApp.StakingKeeper,
		ConsumerKeeper: intoApp.ConsumerKeeper,
		BankKeeper:     intoApp.BankKeeper,
		SlashingKeeper: intoApp.SlashingKeeper,
	}

	handler := v1100.CreateUpgradeHandler(mm, configurator, keepers)

	// Run the handler
	plan := upgradetypes.Plan{Name: v1100.UpgradeName, Height: 100}
	versionMap := intoApp.ModuleManager.GetVersionMap()
	_, err = handler(ctx, plan, versionMap)
	require.NoError(t, err)

	// Verify ready validator was funded
	finalBalance := intoApp.BankKeeper.GetBalance(ctx, accAddr, v1100.Denom)
	require.Equal(t, int64(v1100.MinStake), finalBalance.Amount.Int64(), "Ready validator should be funded to MinStake")

	// Verify DAO balance decreased exactly by MinStake
	finalDaoBalance := intoApp.BankKeeper.GetBalance(ctx, daoAddr, v1100.Denom)
	diff := initialDaoBalance.Amount.Sub(finalDaoBalance.Amount)
	require.Equal(t, int64(v1100.MinStake), diff.Int64(), "DAO should spend MinStake")

	// Verify random validators are jailed
	allVals, err := intoApp.StakingKeeper.GetAllValidators(ctx)
	require.NoError(t, err)

	foundReady := false
	for _, v := range allVals {
		if v.OperatorAddress == readyValoperStr {
			foundReady = true
			require.False(t, v.IsJailed(), "Ready validator should NOT be jailed")
		} else {
			require.True(t, v.IsJailed(), "Random validator %s should be jailed", v.OperatorAddress)
		}
	}
	require.True(t, foundReady, "Should find the ready validator in the store")

	// Verify Slashing Params
	params, err := intoApp.SlashingKeeper.GetParams(ctx)
	require.NoError(t, err)
	require.True(t, math.LegacyNewDecWithPrec(5, 1).Equal(params.MinSignedPerWindow), "MinSignedPerWindow should be 0.5")
}
