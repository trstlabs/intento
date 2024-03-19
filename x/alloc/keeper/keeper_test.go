package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
	"github.com/trstlabs/intento/app"
	"github.com/trstlabs/intento/x/alloc/types"
)

type KeeperTestSuite struct {
	suite.Suite
	ctx sdk.Context
	app *app.IntoApp
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(true)

	suite.app.AllocKeeper.SetParams(suite.ctx, types.DefaultParams())

	suite.app.StakingKeeper.SetParams(suite.ctx, stakingtypes.DefaultParams())

	suite.app.DistrKeeper.SetFeePool(suite.ctx, distrtypes.InitialFeePool())
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func fundModuleAccount(bankKeeper bankkeeper.Keeper, ctx sdk.Context, recipientMod string, amounts sdk.Coins) error {
	if err := bankKeeper.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}
	return bankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, recipientMod, amounts)
}

func (suite *KeeperTestSuite) TestDistribution() {
	suite.SetupTest()

	denom := suite.app.StakingKeeper.BondDenom(suite.ctx)

	allocKeeper := suite.app.AllocKeeper
	params := suite.app.AllocKeeper.GetParams(suite.ctx)
	contributorRewardsReceiver := sdk.AccAddress([]byte("addr1---------------"))
	params.DistributionProportions.CommunityPool = sdk.NewDecWithPrec(25, 2)
	params.DistributionProportions.ContributorRewards = sdk.NewDecWithPrec(5, 2)
	params.DistributionProportions.Staking = sdk.NewDecWithPrec(50, 2)
	// params.DistributionProportions.TrustlessContractIncentives = sdk.NewDecWithPrec(10, 2)
	params.DistributionProportions.RelayerIncentives = sdk.NewDecWithPrec(20, 2)
	params.WeightedContributorRewardsReceivers = []types.WeightedAddress{
		{
			Address: contributorRewardsReceiver.String(),
			Weight:  sdk.NewDec(1),
		},
	}
	suite.app.AllocKeeper.SetParams(suite.ctx, params)

	feePool := suite.app.DistrKeeper.GetFeePool(suite.ctx)
	feeCollector := suite.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	suite.Equal(
		"0",
		suite.app.BankKeeper.GetAllBalances(suite.ctx, feeCollector).AmountOf(denom).String())
	suite.Equal(
		sdk.NewDec(0),
		feePool.CommunityPool.AmountOf(denom))

	mintCoin := sdk.NewCoin(denom, sdk.NewInt(100_000))
	mintCoins := sdk.Coins{mintCoin}
	feeCollectorAccount := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, authtypes.FeeCollectorName)
	suite.Require().NotNil(feeCollectorAccount)

	suite.Require().NoError(fundModuleAccount(suite.app.BankKeeper, suite.ctx, feeCollectorAccount.GetName(), mintCoins))

	feeCollector = suite.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	suite.Equal(
		mintCoin.Amount.String(),
		suite.app.BankKeeper.GetAllBalances(suite.ctx, feeCollector).AmountOf(denom).String())

	suite.Equal(
		sdk.NewDec(0),
		feePool.CommunityPool.AmountOf(denom))

	allocKeeper.DistributeInflation(suite.ctx)

	feeCollector = suite.app.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	/* 	totalDistr := params.DistributionProportions.TrustlessContractIncentives.
	Add(params.DistributionProportions.ContributorRewards.Add(params.DistributionProportions.CommunityPool)).Add(params.DistributionProportions.RelayerIncentives) // 40% */
	totalDistr := params.DistributionProportions.ContributorRewards.Add(params.DistributionProportions.CommunityPool).Add(params.DistributionProportions.RelayerIncentives) // 40%
	autoIbcTxModule := suite.app.AccountKeeper.GetModuleAddress("intent")
	suite.Equal(
		sdk.NewDecFromInt(mintCoin.Amount).Mul(params.DistributionProportions.RelayerIncentives).RoundInt().String(),
		suite.app.BankKeeper.GetAllBalances(suite.ctx, autoIbcTxModule).AmountOf(denom).String())

	// remaining going to fee collector for staking rewards should be 100% - 40% = 60%
	suite.Equal(
		sdk.NewDecFromInt(mintCoin.Amount).Mul(sdk.NewDecWithPrec(100, 2).Sub(totalDistr)).RoundInt().String(),
		suite.app.BankKeeper.GetAllBalances(suite.ctx, feeCollector).AmountOf(denom).String())

	suite.Equal(
		sdk.NewDecFromInt(mintCoin.Amount).Mul(params.DistributionProportions.ContributorRewards).TruncateInt(),
		suite.app.BankKeeper.GetBalance(suite.ctx, contributorRewardsReceiver, denom).Amount)

	feePool = suite.app.DistrKeeper.GetFeePool(suite.ctx)
	suite.Equal(
		sdk.NewDecFromInt(mintCoin.Amount).Mul(params.DistributionProportions.CommunityPool),
		feePool.CommunityPool.AmountOf(denom))
}
