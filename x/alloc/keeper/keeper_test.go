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
	"github.com/trstlabs/trst/app"
	"github.com/trstlabs/trst/x/alloc/types"
)

type KeeperTestSuite struct {
	suite.Suite
	ctx sdk.Context
	app *app.TrstApp
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(true)

	suite.app.AppKeepers.AllocKeeper.SetParams(suite.ctx, types.DefaultParams())

	suite.app.AppKeepers.StakingKeeper.SetParams(suite.ctx, stakingtypes.DefaultParams())

	suite.app.AppKeepers.DistrKeeper.SetFeePool(suite.ctx, distrtypes.InitialFeePool())
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

	denom := suite.app.AppKeepers.StakingKeeper.BondDenom(suite.ctx)

	allocKeeper := suite.app.AppKeepers.AllocKeeper
	params := suite.app.AppKeepers.AllocKeeper.GetParams(suite.ctx)
	contributorRewardsReceiver := sdk.AccAddress([]byte("addr1---------------"))
	params.DistributionProportions.CommunityPool = sdk.NewDecWithPrec(25, 2)
	params.DistributionProportions.ContributorRewards = sdk.NewDecWithPrec(5, 2)
	params.DistributionProportions.Staking = sdk.NewDecWithPrec(60, 2)
	params.DistributionProportions.TrustlessContractIncentives = sdk.NewDecWithPrec(10, 2)
	params.WeightedContributorRewardsReceivers = []types.WeightedAddress{
		{
			Address: contributorRewardsReceiver.String(),
			Weight:  sdk.NewDec(1),
		},
	}
	suite.app.AppKeepers.AllocKeeper.SetParams(suite.ctx, params)

	feePool := suite.app.AppKeepers.DistrKeeper.GetFeePool(suite.ctx)
	feeCollector := suite.app.AppKeepers.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	suite.Equal(
		"0",
		suite.app.AppKeepers.BankKeeper.GetAllBalances(suite.ctx, feeCollector).AmountOf(denom).String())
	suite.Equal(
		sdk.NewDec(0),
		feePool.CommunityPool.AmountOf(denom))

	mintCoin := sdk.NewCoin(denom, sdk.NewInt(100_000))
	mintCoins := sdk.Coins{mintCoin}
	feeCollectorAccount := suite.app.AppKeepers.AccountKeeper.GetModuleAccount(suite.ctx, authtypes.FeeCollectorName)
	suite.Require().NotNil(feeCollectorAccount)

	suite.Require().NoError(fundModuleAccount(suite.app.AppKeepers.BankKeeper, suite.ctx, feeCollectorAccount.GetName(), mintCoins))

	feeCollector = suite.app.AppKeepers.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	suite.Equal(
		mintCoin.Amount.String(),
		suite.app.AppKeepers.BankKeeper.GetAllBalances(suite.ctx, feeCollector).AmountOf(denom).String())

	suite.Equal(
		sdk.NewDec(0),
		feePool.CommunityPool.AmountOf(denom))

	allocKeeper.DistributeInflation(suite.ctx)

	feeCollector = suite.app.AppKeepers.AccountKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	totalDistr := params.DistributionProportions.TrustlessContractIncentives.
		Add(params.DistributionProportions.ContributorRewards.Add(params.DistributionProportions.CommunityPool)) // 15%

	// remaining going to fee collector should be 100% - 15% = 85%
	suite.Equal(
		mintCoin.Amount.ToDec().Mul(sdk.NewDecWithPrec(100, 2).Sub(totalDistr)).RoundInt().String(),
		suite.app.AppKeepers.BankKeeper.GetAllBalances(suite.ctx, feeCollector).AmountOf(denom).String())

	suite.Equal(
		mintCoin.Amount.ToDec().Mul(params.DistributionProportions.ContributorRewards).TruncateInt(),
		suite.app.AppKeepers.BankKeeper.GetBalance(suite.ctx, contributorRewardsReceiver, denom).Amount)

	feePool = suite.app.AppKeepers.DistrKeeper.GetFeePool(suite.ctx)
	suite.Equal(
		mintCoin.Amount.ToDec().Mul(params.DistributionProportions.CommunityPool),
		feePool.CommunityPool.AmountOf(denom))
}
