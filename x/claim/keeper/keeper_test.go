package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
	"github.com/trstlabs/intento/app"
	alloctypes "github.com/trstlabs/intento/x/alloc/types"
	"github.com/trstlabs/intento/x/claim/keeper"
	"github.com/trstlabs/intento/x/claim/types"
)

type KeeperTestSuite struct {
	suite.Suite
	ctx     sdk.Context
	app     *app.IntoApp
	msgSrvr types.MsgServer
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(true)

	suite.app.AllocKeeper.SetParams(suite.ctx, alloctypes.DefaultParams())

	suite.app.StakingKeeper.SetParams(suite.ctx, stakingtypes.DefaultParams())

	suite.app.DistrKeeper.SetFeePool(suite.ctx, distrtypes.InitialFeePool())
	suite.app.ClaimKeeper.CreateModuleAccount(suite.ctx, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000)))
	startTime := time.Now()

	suite.msgSrvr = keeper.NewMsgServerImpl(suite.app.ClaimKeeper)
	suite.app.ClaimKeeper.SetParams(suite.ctx, types.Params{
		//AirdropEnabled:     true,
		AirdropStartTime:       startTime,
		DurationUntilDecay:     types.DefaultDurationUntilDecay,
		DurationOfDecay:        types.DefaultDurationOfDecay,
		ClaimDenom:             sdk.DefaultBondDenom,
		DurationVestingPeriods: types.DefaultDurationVestingPeriods,
	})
	suite.ctx = suite.ctx.WithBlockTime(startTime)
}

func (s *KeeperTestSuite) TestModuleAccountCreated() {
	app, ctx := s.app, s.ctx
	moduleAddress := app.AccountKeeper.GetModuleAddress(types.ModuleName)
	balance := app.BankKeeper.GetBalance(ctx, moduleAddress, sdk.DefaultBondDenom)
	s.Require().Equal(fmt.Sprintf("10000000%s", sdk.DefaultBondDenom), balance.String())
}

func (s *KeeperTestSuite) TestClaimClaimable() {
	pub1 := ed25519.GenPrivKey().PubKey()
	pub2 := ed25519.GenPrivKey().PubKey()
	pub3 := ed25519.GenPrivKey().PubKey()
	pub4 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pub1.Address())
	addr2 := sdk.AccAddress(pub2.Address())
	addr3 := sdk.AccAddress(pub3.Address())
	valAddr := sdk.ValAddress(pub4.Address())

	status := types.Status{ActionCompleted: false, VestingPeriodCompleted: []bool{false, false, false, false}, VestingPeriodClaimed: []bool{false, false, false, false}}
	claimRecords := []types.ClaimRecord{
		{
			Address:                addr1.String(),
			InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 2000)),
			Status:                 []types.Status{status, status, status, status},
		},
		{
			Address:                addr2.String(),
			InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 2000)),
			Status:                 []types.Status{status, status, status, status},
		},
	}

	s.app.AccountKeeper.SetAccount(s.ctx, authtypes.NewBaseAccount(addr1, nil, 0, 0))
	s.app.AccountKeeper.SetAccount(s.ctx, authtypes.NewBaseAccount(addr2, nil, 0, 0))

	s.app.ClaimKeeper.SetParams(s.ctx, types.Params{
		//AirdropEnabled:     false,
		AirdropStartTime:       time.Now().Add(time.Hour * -1),
		ClaimDenom:             sdk.DefaultBondDenom,
		DurationUntilDecay:     time.Hour,
		DurationOfDecay:        time.Hour * 4,
		DurationVestingPeriods: types.DefaultDurationVestingPeriods,
		//AllowedClaimers:    make([]types.ClaimAuthorization, 0),
	})
	err := s.app.ClaimKeeper.SetClaimRecords(s.ctx, claimRecords)
	s.Require().NoError(err)
	msgClaimClaimable := types.NewMsgClaimClaimable(addr1)
	ctx := sdk.WrapSDKContext(s.ctx)
	_, err = s.msgSrvr.ClaimClaimable(ctx, msgClaimClaimable)
	s.Error(err)
	s.Contains(err.Error(), "address does not have claimable tokens right now")

	// unauthorized
	msgClaimClaimable = types.NewMsgClaimClaimable(addr3)
	_, err = s.msgSrvr.ClaimClaimable(ctx, msgClaimClaimable)
	s.Require().Error(err)
	s.Contains(err.Error(), "address does not have claim record")

	// claim
	s.app.ClaimKeeper.AfterActionAuthz(s.ctx, addr1)

	record, err := s.app.ClaimKeeper.GetClaimRecord(s.ctx, addr1)
	s.Require().NoError(err)
	s.Require().True(record.Status[0].ActionCompleted)

	record.Status[0].VestingPeriodCompleted = []bool{true, false, false, false}

	err = s.app.ClaimKeeper.SetClaimRecord(s.ctx, record)
	s.Require().NoError(err)

	// addr1 has to delegate
	validator, err := stakingtypes.NewValidator(valAddr, pub1, stakingtypes.Description{})
	s.Require().NoError(err)
	validator = stakingkeeper.TestingUpdateValidator(s.app.StakingKeeper, s.ctx, validator, true)
	s.app.StakingKeeper.Hooks().AfterValidatorCreated(s.ctx, validator.GetOperator())

	validator, _ = validator.AddTokensFromDel(sdk.TokensFromConsensusPower(1, sdk.DefaultPowerReduction))

	_, err = s.app.StakingKeeper.Delegate(s.ctx, addr1, record.InitialClaimableAmount.AmountOf(sdk.DefaultBondDenom).Quo(sdk.NewInt(20)), stakingtypes.Unbonded, validator, true)
	s.NoError(err)
	_, err = s.app.StakingKeeper.Delegate(s.ctx, addr1, record.InitialClaimableAmount.AmountOf(sdk.DefaultBondDenom).Quo(sdk.NewInt(20)), stakingtypes.Unbonded, validator, true)
	s.NoError(err)

	msgClaimClaimable = types.NewMsgClaimClaimable(addr1)
	_, err = s.msgSrvr.ClaimClaimable(ctx, msgClaimClaimable)
	s.Require().NoError(err)
	claimedCoins := s.app.BankKeeper.GetAllBalances(s.ctx, addr1)

	s.Require().Equal(claimedCoins.AmountOf(sdk.DefaultBondDenom), sdk.NewInt(100))

	s.Require().Equal(true, record.Status[0].ActionCompleted)

	s.Require().Equal([]bool{false, false, false, false}, record.Status[0].VestingPeriodClaimed)
	s.Require().Equal([]bool{true, false, false, false}, record.Status[0].VestingPeriodCompleted)

	// claim 2
	s.app.ClaimKeeper.AfterActionWasm(s.ctx, addr1)

	record, err = s.app.ClaimKeeper.GetClaimRecord(s.ctx, addr1)
	s.Require().NoError(err)
	s.Require().True(record.Status[1].ActionCompleted)

	record.Status[1].VestingPeriodCompleted = []bool{true, false, false, false}

	err = s.app.ClaimKeeper.SetClaimRecord(s.ctx, record)
	s.Require().NoError(err)

	s.app.StakingKeeper.Delegate(s.ctx, addr1, sdk.NewInt(68), stakingtypes.Unbonded, validator, true)

	msgClaimClaimable = types.NewMsgClaimClaimable(addr1)
	_, err = s.msgSrvr.ClaimClaimable(ctx, msgClaimClaimable)
	s.Require().NoError(err)

	//get delegations and calculate min bonded ratio for claim
	delegationInfo := s.app.StakingKeeper.GetAllDelegatorDelegations(s.ctx, addr1)
	totalDelegations := sdk.ZeroDec()
	for _, delegation := range delegationInfo {
		totalDelegations = totalDelegations.Add(delegation.Shares)
	}
	balanceCoins := s.app.BankKeeper.GetAllBalances(s.ctx, addr1)

	s.Require().Equal(
		balanceCoins.AmountOf(sdk.DefaultBondDenom), sdk.NewInt(232))
	s.Require().Equal(
		totalDelegations, sdk.NewDec(268))

	record, err = s.app.ClaimKeeper.GetClaimRecord(s.ctx, addr1)
	s.Require().NoError(err)
	s.Require().Equal(true, record.Status[0].ActionCompleted)
	s.Require().Equal(true, record.Status[1].ActionCompleted)
	s.Require().Equal([]bool{true, false, false, false}, record.Status[1].VestingPeriodClaimed)
	s.Require().Equal([]bool{true, false, false, false}, record.Status[1].VestingPeriodCompleted)

	// claim second address
	s.app.ClaimKeeper.AfterActionAuthz(s.ctx, addr2)
	s.app.ClaimKeeper.AfterActionWasm(s.ctx, addr2)
	record, err = s.app.ClaimKeeper.GetClaimRecord(s.ctx, addr2)
	s.Require().NoError(err)
	s.Require().True(record.Status[0].ActionCompleted)
	s.Require().True(record.Status[1].ActionCompleted)

	record.Status[0].VestingPeriodCompleted = []bool{true, false, false, false}
	record.Status[1].VestingPeriodCompleted = []bool{true, false, false, false}

	err = s.app.ClaimKeeper.SetClaimRecord(s.ctx, record)
	s.Require().NoError(err)
	s.app.StakingKeeper.Delegate(s.ctx, addr2, sdk.NewInt(200), stakingtypes.Unbonded, validator, true)
	s.app.StakingKeeper.Delegate(s.ctx, addr2, sdk.NewInt(68), stakingtypes.Unbonded, validator, true)
	msgClaimClaimable = types.NewMsgClaimClaimable(addr2)
	_, err = s.msgSrvr.ClaimClaimable(ctx, msgClaimClaimable)
	s.Require().NoError(err)

	//get delegations and calculate min bonded ratio for claim
	delegationInfo = s.app.StakingKeeper.GetAllDelegatorDelegations(s.ctx, addr2)
	totalDelegations = sdk.ZeroDec()
	for _, delegation := range delegationInfo {
		totalDelegations = totalDelegations.Add(delegation.Shares)
	}
	balanceCoins = s.app.BankKeeper.GetAllBalances(s.ctx, addr2)

	s.Require().Equal(
		balanceCoins.AmountOf(sdk.DefaultBondDenom), sdk.NewInt(232))
	s.Require().Equal(
		totalDelegations, sdk.NewDec(268))

	record, err = s.app.ClaimKeeper.GetClaimRecord(s.ctx, addr2)
	s.Require().NoError(err)
	s.Require().Equal(true, record.Status[0].ActionCompleted)
	s.Require().Equal(true, record.Status[1].ActionCompleted)
	s.Require().Equal([]bool{true, false, false, false}, record.Status[0].VestingPeriodClaimed)
	s.Require().Equal([]bool{true, false, false, false}, record.Status[0].VestingPeriodCompleted)
	s.Require().Equal([]bool{true, false, false, false}, record.Status[1].VestingPeriodClaimed)
	s.Require().Equal([]bool{true, false, false, false}, record.Status[1].VestingPeriodCompleted)

}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func FundAccount(bankKeeper bankkeeper.Keeper, ctx sdk.Context, addr sdk.AccAddress, amounts sdk.Coins) error {
	if err := bankKeeper.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}
	return bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, amounts)
}

func (s *KeeperTestSuite) TestHookOfUnclaimableAccount() {
	pub1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pub1.Address())
	s.app.AccountKeeper.SetAccount(s.ctx, authtypes.NewBaseAccount(addr1, nil, 0, 0))

	claim, err := s.app.ClaimKeeper.GetClaimRecord(s.ctx, addr1)
	s.Contains(err.Error(), "address does not have claim record")
	s.Equal(types.ClaimRecord{}, claim)

	s.app.ClaimKeeper.AfterDelegationModified(s.ctx, addr1, sdk.ValAddress(addr1))

	balances := s.app.BankKeeper.GetAllBalances(s.ctx, addr1)
	s.Equal(sdk.Coins{}, balances)
}

func (s *KeeperTestSuite) TestHookBeforeAirdropStart() {
	s.SetupTest()

	airdropStartTime := time.Now().Add(time.Hour)

	s.app.ClaimKeeper.SetParams(s.ctx, types.Params{
		//AirdropEnabled:     true,
		ClaimDenom:             sdk.DefaultBondDenom,
		AirdropStartTime:       airdropStartTime,
		DurationUntilDecay:     time.Hour,
		DurationOfDecay:        time.Hour * 4,
		DurationVestingPeriods: types.DefaultDurationVestingPeriods,
	})

	pub1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pub1.Address())
	status := types.Status{ActionCompleted: false, VestingPeriodCompleted: []bool{false, false, false, false}, VestingPeriodClaimed: []bool{false, false, false, false}}
	claimRecords := []types.ClaimRecord{
		{
			Address:                addr1.String(),
			InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)),
			Status:                 []types.Status{status, status, status, status},
		},
	}
	s.app.AccountKeeper.SetAccount(s.ctx, authtypes.NewBaseAccount(addr1, nil, 0, 0))

	err := s.app.ClaimKeeper.SetClaimRecords(s.ctx, claimRecords)
	s.Require().NoError(err)

	coins, err := s.app.ClaimKeeper.GetTotalClaimableForAddr(s.ctx, addr1)
	s.NoError(err)
	// Now, it is before starting air drop, so this value should return the empty coins
	s.True(coins.Empty())

	coins, err = s.app.ClaimKeeper.GetTotalClaimableAmountForAction(s.ctx, addr1, types.ActionDelegateStake)
	s.NoError(err)
	// Now, it is before starting air drop, so this value should return the empty coins
	s.True(coins.Empty())

	s.app.ClaimKeeper.AfterDelegationModified(s.ctx, addr1, sdk.ValAddress(addr1))
	balances := s.app.BankKeeper.GetAllBalances(s.ctx, addr1)
	// Now, it is before starting air drop, so claim module should not send the balances to the user after delegate.
	s.True(balances.Empty())

	s.app.ClaimKeeper.AfterDelegationModified(s.ctx.WithBlockTime(airdropStartTime), addr1, sdk.ValAddress(addr1))
	balances = s.app.BankKeeper.GetAllBalances(s.ctx, addr1)
	//fmt.Printf("%v \n", balances)
	// Now, it is the time for air drop, so claim module should send the balances to the user after delegate.
	s.Equal(claimRecords[0].InitialClaimableAmount.AmountOf(sdk.DefaultBondDenom).Quo(sdk.NewInt(int64(len(types.Action_value)))).Quo(sdk.NewInt(5)), balances.AmountOf(sdk.DefaultBondDenom))
}

func (s *KeeperTestSuite) TestAirdropDisabled() {
	s.SetupTest()

	airdropStartTime := time.Now().Add(time.Hour)

	s.app.ClaimKeeper.SetParams(s.ctx, types.Params{
		//AirdropEnabled:     false,
		ClaimDenom:             sdk.DefaultBondDenom,
		DurationUntilDecay:     time.Hour,
		DurationOfDecay:        time.Hour * 4,
		DurationVestingPeriods: types.DefaultDurationVestingPeriods,
	})

	pub1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pub1.Address())

	status := types.Status{ActionCompleted: false, VestingPeriodCompleted: []bool{false, false, false, false}, VestingPeriodClaimed: []bool{false, false, false, false}}
	claimRecords := []types.ClaimRecord{
		{
			Address:                addr1.String(),
			InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)),
			Status:                 []types.Status{status, status, status, status},
		},
	}
	s.app.AccountKeeper.SetAccount(s.ctx, authtypes.NewBaseAccount(addr1, nil, 0, 0))

	err := s.app.ClaimKeeper.SetClaimRecords(s.ctx, claimRecords)
	s.Require().NoError(err)

	coins, err := s.app.ClaimKeeper.GetTotalClaimableForAddr(s.ctx, addr1)
	s.NoError(err)
	// Now, it is before starting air drop, so this value should return the empty coins
	s.True(coins.Empty())

	coins, err = s.app.ClaimKeeper.GetTotalClaimableAmountForAction(s.ctx, addr1, types.ActionDelegateStake)
	s.NoError(err)
	// Now, it is before starting air drop, so this value should return the empty coins
	s.True(coins.Empty())

	s.app.ClaimKeeper.AfterDelegationModified(s.ctx, addr1, sdk.ValAddress(addr1))
	balances := s.app.BankKeeper.GetAllBalances(s.ctx, addr1)
	// Now, it is before starting air drop, so claim module should not send the balances to the user after delegate.
	s.True(balances.Empty())

	s.app.ClaimKeeper.AfterGovernanceVoted(s.ctx, addr1)
	balances = s.app.BankKeeper.GetAllBalances(s.ctx, addr1)
	// Now, it is before starting air drop, so claim module should not send the balances to the user after vote.
	s.True(balances.Empty())

	// set airdrop enabled but with invalid date
	s.app.ClaimKeeper.SetParams(s.ctx, types.Params{
		//AirdropEnabled:     true,
		ClaimDenom:             sdk.DefaultBondDenom,
		DurationUntilDecay:     time.Hour,
		DurationOfDecay:        time.Hour * 4,
		DurationVestingPeriods: types.DefaultDurationVestingPeriods,
	})

	s.app.ClaimKeeper.AfterDelegationModified(s.ctx, addr1, sdk.ValAddress(addr1))
	balances = s.app.BankKeeper.GetAllBalances(s.ctx, addr1)
	// Now airdrop is enabled but a potential misconfiguraion on start time
	s.True(balances.Empty())

	s.app.ClaimKeeper.AfterGovernanceVoted(s.ctx, addr1)
	balances = s.app.BankKeeper.GetAllBalances(s.ctx, addr1)
	// Now airdrop is enabled but a potential misconfiguraion on start time, so claim module should not send the balances to the user after vote.
	s.True(balances.Empty())

	// set airdrop enabled but with date in the future
	s.app.ClaimKeeper.SetParams(s.ctx, types.Params{
		//AirdropEnabled:     true,
		AirdropStartTime:       airdropStartTime.Add(time.Hour),
		ClaimDenom:             sdk.DefaultBondDenom,
		DurationUntilDecay:     time.Hour,
		DurationOfDecay:        time.Hour * 4,
		DurationVestingPeriods: types.DefaultDurationVestingPeriods,
	})

	s.app.ClaimKeeper.AfterDelegationModified(s.ctx, addr1, sdk.ValAddress(addr1))
	balances = s.app.BankKeeper.GetAllBalances(s.ctx, addr1)
	// Now airdrop is enabled  and date is not empty but block time still behid
	s.True(balances.Empty())

	s.app.ClaimKeeper.AfterGovernanceVoted(s.ctx, addr1)
	balances = s.app.BankKeeper.GetAllBalances(s.ctx, addr1)
	// Now airdrop is enabled  and date is not empty but block time still behid
	s.True(balances.Empty())

	// add extra 2 hours
	s.app.ClaimKeeper.AfterDelegationModified(s.ctx.WithBlockTime(airdropStartTime.Add(time.Hour*2)), addr1, sdk.ValAddress(addr1))
	balances = s.app.BankKeeper.GetAllBalances(s.ctx, addr1)
	// Now, it is the time for air drop, so claim module should send the balances to the user after delegate.
	s.Equal(claimRecords[0].InitialClaimableAmount.AmountOf(sdk.DefaultBondDenom).Quo(sdk.NewInt(int64(len(types.Action_value)))).Quo(sdk.NewInt(5)), balances.AmountOf(sdk.DefaultBondDenom))
}

func (s *KeeperTestSuite) TestDuplicatedActionNotWithdrawRepeatedly() {
	s.SetupTest()

	pub1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pub1.Address())
	status := types.Status{ActionCompleted: false, VestingPeriodCompleted: []bool{false, false, false, false}, VestingPeriodClaimed: []bool{false, false, false, false}}
	claimRecords := []types.ClaimRecord{
		{
			Address:                addr1.String(),
			InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 2000)),
			Status:                 []types.Status{status, status, status, status},
		},
	}
	s.app.AccountKeeper.SetAccount(s.ctx, authtypes.NewBaseAccount(addr1, nil, 0, 0))

	err := s.app.ClaimKeeper.SetClaimRecords(s.ctx, claimRecords)
	s.Require().NoError(err)

	coins1, err := s.app.ClaimKeeper.GetTotalClaimableForAddr(s.ctx, addr1)
	s.Require().NoError(err)
	s.Require().Equal(coins1, claimRecords[0].InitialClaimableAmount)

	s.app.ClaimKeeper.AfterDelegationModified(s.ctx, addr1, sdk.ValAddress(addr1))
	claim, err := s.app.ClaimKeeper.GetClaimRecord(s.ctx, addr1)
	s.NoError(err)
	s.True(claim.Status[3].ActionCompleted)

	claimedCoins := s.app.BankKeeper.GetAllBalances(s.ctx, addr1)
	s.Require().Equal(claimedCoins.AmountOf(sdk.DefaultBondDenom), claimRecords[0].InitialClaimableAmount.AmountOf(sdk.DefaultBondDenom).Quo(sdk.NewInt(20)))

	s.app.ClaimKeeper.AfterDelegationModified(s.ctx, addr1, sdk.ValAddress(addr1))
	claim, err = s.app.ClaimKeeper.GetClaimRecord(s.ctx, addr1)
	s.NoError(err)
	s.True(claim.Status[3].ActionCompleted)

	claimedCoins = s.app.BankKeeper.GetAllBalances(s.ctx, addr1)
	s.Require().Equal(claimedCoins.AmountOf(sdk.DefaultBondDenom), claimRecords[0].InitialClaimableAmount.AmountOf(sdk.DefaultBondDenom).Quo(sdk.NewInt(20)))
}

func (s *KeeperTestSuite) TestNotRunningGenesisBlock() {
	s.ctx = s.ctx.WithBlockHeight(1)
	s.app.ClaimKeeper.CreateModuleAccount(s.ctx, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000)))
	// set airdrop enabled but with date in the future
	s.app.ClaimKeeper.SetParams(s.ctx, types.Params{
		//AirdropEnabled:     true,
		AirdropStartTime:       time.Now().Add(time.Hour * -1),
		ClaimDenom:             sdk.DefaultBondDenom,
		DurationUntilDecay:     time.Hour,
		DurationOfDecay:        time.Hour * 4,
		DurationVestingPeriods: types.DefaultDurationVestingPeriods,
	})

	pub1 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pub1.Address())
	status := types.Status{ActionCompleted: false, VestingPeriodCompleted: []bool{false, false, false, false}, VestingPeriodClaimed: []bool{false, false, false, false}}
	claimRecords := []types.ClaimRecord{
		{
			Address:                addr1.String(),
			InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 2000)),
			Status:                 []types.Status{status, status, status, status},
		},
	}
	s.app.AccountKeeper.SetAccount(s.ctx, authtypes.NewBaseAccount(addr1, nil, 0, 0))

	err := s.app.ClaimKeeper.SetClaimRecords(s.ctx, claimRecords)
	s.Require().NoError(err)

	coins1, err := s.app.ClaimKeeper.GetTotalClaimableForAddr(s.ctx, addr1)
	s.Require().NoError(err)
	s.Require().Equal(coins1, claimRecords[0].InitialClaimableAmount)

	s.app.ClaimKeeper.AfterDelegationModified(s.ctx, addr1, sdk.ValAddress(addr1))
	claim, err := s.app.ClaimKeeper.GetClaimRecord(s.ctx, addr1)
	s.NoError(err)
	s.True(claim.Status[3].ActionCompleted)

	coins1, err = s.app.ClaimKeeper.GetTotalClaimableForAddr(s.ctx, addr1)
	s.Require().NoError(err)
	s.Require().Equal(coins1, claimRecords[0].InitialClaimableAmount)
}

func (s *KeeperTestSuite) TestEndAirdrop() {
	// set airdrop enabled but with date in the future
	s.app.ClaimKeeper.SetParams(s.ctx, types.Params{
		//AirdropEnabled:     true,
		AirdropStartTime:       time.Now().Add(time.Hour * -1),
		ClaimDenom:             sdk.DefaultBondDenom,
		DurationUntilDecay:     time.Hour,
		DurationOfDecay:        time.Hour * 4,
		DurationVestingPeriods: types.DefaultDurationVestingPeriods,
	})

	pub1 := ed25519.GenPrivKey().PubKey()
	pub2 := ed25519.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pub1.Address())
	addr2 := sdk.AccAddress(pub2.Address())
	status := types.Status{ActionCompleted: false, VestingPeriodCompleted: []bool{false, false, false, false}, VestingPeriodClaimed: []bool{false, false, false, false}}
	claimRecords := []types.ClaimRecord{
		{
			Address:                addr1.String(),
			InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)),
			Status:                 []types.Status{status, status, status, status},
		},
		{
			Address:                addr2.String(),
			InitialClaimableAmount: sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)),
			Status:                 []types.Status{status, status, status, status},
		},
	}

	s.app.AccountKeeper.SetAccount(s.ctx, authtypes.NewBaseAccount(addr1, nil, 0, 0))
	s.app.AccountKeeper.SetAccount(s.ctx, authtypes.NewBaseAccount(addr2, nil, 0, 0))

	err := s.app.ClaimKeeper.SetClaimRecords(s.ctx, claimRecords)
	s.Require().NoError(err)

	err = s.app.ClaimKeeper.EndAirdrop(s.ctx)
	s.Require().NoError(err)

	moduleAccAddr := s.app.AccountKeeper.GetModuleAddress(types.ModuleName)
	coins := s.app.BankKeeper.GetBalance(s.ctx, moduleAccAddr, sdk.DefaultBondDenom)
	s.Require().Equal(sdk.NewInt64Coin(sdk.DefaultBondDenom, 0).String(), coins.String())
}
