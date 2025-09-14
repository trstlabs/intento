package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
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

func (s *KeeperTestSuite) SetupTest() {
	s.app, s.ctx = createTestApp(true)

	s.app.AllocKeeper.SetParams(s.ctx, alloctypes.DefaultParams())

	s.app.StakingKeeper.SetParams(s.ctx, stakingtypes.DefaultParams())

	s.app.DistrKeeper.FeePool.Set(s.ctx, distrtypes.InitialFeePool())
	s.app.ClaimKeeper.CreateModuleAccount(s.ctx, sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(10000000)))
	startTime := time.Now()

	s.msgSrvr = keeper.NewMsgServerImpl(s.app.ClaimKeeper)
	s.app.ClaimKeeper.SetParams(s.ctx, types.Params{
		AirdropStartTime:       startTime,
		DurationUntilDecay:     types.DefaultDurationUntilDecay,
		DurationOfDecay:        types.DefaultDurationOfDecay,
		ClaimDenom:             sdk.DefaultBondDenom,
		DurationVestingPeriods: types.DefaultDurationVestingPeriods,
	})
	s.ctx = s.ctx.WithBlockTime(startTime)
}

func (s *KeeperTestSuite) TestModuleAccountCreated() {
	app, ctx := s.app, s.ctx
	moduleAddress := app.AccountKeeper.GetModuleAddress(types.ModuleName)
	balance := app.BankKeeper.GetBalance(ctx, moduleAddress, sdk.DefaultBondDenom)
	s.Require().Equal(fmt.Sprintf("10000000%s", sdk.DefaultBondDenom), balance.String())
}

func (s *KeeperTestSuite) TestClaimClaimable() {

	pub1 := ed25519.GenPrivKey().PubKey()

	valAddr := sdk.ValAddress(pub1.Address())

	addr1 := s.createAccount()
	//addr3 := s.createAccount()
	claimRecords := s.createClaimRecords(addr1, 2000)

	s.app.ClaimKeeper.SetParams(s.ctx, types.Params{
		AirdropStartTime:       time.Now().Add(time.Hour * -1),
		ClaimDenom:             sdk.DefaultBondDenom,
		DurationUntilDecay:     time.Hour,
		DurationOfDecay:        time.Hour * 4,
		DurationVestingPeriods: types.DefaultDurationVestingPeriods,
	})
	err := s.app.ClaimKeeper.SetClaimRecords(s.ctx, claimRecords)
	s.Require().NoError(err)

	// Attempt claim - unauthorized address
	msgClaimClaimable := types.NewMsgClaimClaimable(addr1)

	_, err = s.msgSrvr.ClaimClaimable(s.ctx, msgClaimClaimable)
	s.Require().Error(err)
	s.Contains(err.Error(), "address does not have claimable tokens right now")

	// Setup and process claims
	s.app.ClaimKeeper.AfterActionLocal(s.ctx, addr1)
	record, err := s.app.ClaimKeeper.GetClaimRecord(s.ctx, addr1)
	s.Require().NoError(err)
	s.Require().True(record.Status[0].ActionCompleted)

	// Validator setup for delegation
	validator, err := stakingtypes.NewValidator(valAddr.String(), pub1, stakingtypes.Description{})
	s.Require().NoError(err)
	validator = stakingkeeper.TestingUpdateValidator(s.app.StakingKeeper, s.ctx, validator, true)
	s.app.StakingKeeper.Hooks().AfterValidatorCreated(s.ctx, valAddr)
	validator, _ = validator.AddTokensFromDel(sdk.TokensFromConsensusPower(1, sdk.DefaultPowerReduction))

	balanceBeforeAction := s.app.BankKeeper.GetAllBalances(s.ctx, addr1)
	// Delegate tokens and process claim
	_, err = s.app.StakingKeeper.Delegate(s.ctx, addr1, math.NewInt(67), stakingtypes.Unbonded, validator, true)
	s.NoError(err)

	balance := s.app.BankKeeper.GetAllBalances(s.ctx, addr1)

	//Initial Claimable Action 1
	s.Require().Equal(balance.AmountOf(sdk.DefaultBondDenom).Sub(balanceBeforeAction[0].Amount).Add(math.NewInt(67)), math.NewInt(2000/4/5))

	_, err = s.msgSrvr.ClaimClaimable(s.ctx, msgClaimClaimable)
	s.Contains(err.Error(), "address does not have claimable tokens right now")

	//happends in endblocker
	record.Status[0].VestingPeriodsCompleted = []bool{true, false, false, false}
	err = s.app.ClaimKeeper.SetClaimRecord(s.ctx, record)
	s.Require().NoError(err)

	_, err = s.msgSrvr.ClaimClaimable(s.ctx, msgClaimClaimable)
	s.Require().NoError(err)

	//Claimable #1
	balanceAfterClaim := s.app.BankKeeper.GetAllBalances(s.ctx, addr1)
	s.Require().Equal(balanceAfterClaim.AmountOf(sdk.DefaultBondDenom).Sub(balance[0].Amount), math.NewInt(2000/4/5))

	// Process second claim and validate
	s.app.ClaimKeeper.AfterActionICA(s.ctx, addr1)
	record, err = s.app.ClaimKeeper.GetClaimRecord(s.ctx, addr1)
	s.Require().NoError(err)
	s.Require().True(record.Status[1].ActionCompleted)

	//Initial Claimable Action 2
	balance2 := s.app.BankKeeper.GetAllBalances(s.ctx, addr1)
	s.Require().Equal(balance2.Sub(balanceAfterClaim[0])[0].Amount, math.NewInt(2000/4/5))

	//happends in endblocker
	record, err = s.app.ClaimKeeper.GetClaimRecord(s.ctx, addr1)
	s.Require().NoError(err)
	record.Status[1].VestingPeriodsCompleted = []bool{true, false, false, false}
	err = s.app.ClaimKeeper.SetClaimRecord(s.ctx, record)
	s.Require().NoError(err)

	//claim #3
	_, err = s.msgSrvr.ClaimClaimable(s.ctx, msgClaimClaimable)
	s.Require().NoError(err)

	//Claimable #2
	balanceAfterClaim2 := s.app.BankKeeper.GetAllBalances(s.ctx, addr1)
	s.Require().Equal(balanceAfterClaim2.AmountOf(sdk.DefaultBondDenom).Sub(balance2[0].Amount), math.NewInt(2000/4/5))

	_, err = s.app.StakingKeeper.Delegate(s.ctx, addr1, math.NewInt(67), stakingtypes.Unbonded, validator, true)
	s.NoError(err)
	balanceAfterDelegate2 := s.app.BankKeeper.GetAllBalances(s.ctx, addr1)

	//claim #4
	record, err = s.app.ClaimKeeper.GetClaimRecord(s.ctx, addr1)
	s.Require().NoError(err)
	record.Status[1].VestingPeriodsCompleted = []bool{true, true, false, false}
	err = s.app.ClaimKeeper.SetClaimRecord(s.ctx, record)
	s.Require().NoError(err)

	_, err = s.msgSrvr.ClaimClaimable(s.ctx, msgClaimClaimable)
	s.Require().NoError(err)

	//Claimable #3
	balanceAfterClaim3 := s.app.BankKeeper.GetAllBalances(s.ctx, addr1)

	s.Require().Equal(balanceAfterClaim3.AmountOf(sdk.DefaultBondDenom).Sub(balanceAfterDelegate2[0].Amount), math.NewInt(2000/4/5))
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

	addr1 := s.createAccount()

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
		ClaimDenom:             sdk.DefaultBondDenom,
		AirdropStartTime:       airdropStartTime,
		DurationUntilDecay:     time.Hour,
		DurationOfDecay:        time.Hour * 4,
		DurationVestingPeriods: types.DefaultDurationVestingPeriods,
	})

	addr1 := s.createAccount()
	claimRecords := s.createClaimRecords(addr1, 2000)

	err := s.app.ClaimKeeper.SetClaimRecords(s.ctx, claimRecords)
	s.Require().NoError(err)

	coins, err := s.app.ClaimKeeper.GetTotalClaimableForAddr(s.ctx, addr1)
	s.NoError(err)
	// Now, it is before starting air drop, so this value should return the empty coins
	s.True(coins.IsZero())

	totalAction, err := s.app.ClaimKeeper.GetTotalClaimableAmountPerAction(s.ctx, addr1)
	s.NoError(err)
	// Now, it is before starting air drop, so this value should return the empty coins
	s.True(totalAction.IsZero())

	s.app.ClaimKeeper.AfterDelegationModified(s.ctx, addr1, sdk.ValAddress(addr1))
	balances := s.app.BankKeeper.GetAllBalances(s.ctx, addr1)
	// Now, it is before starting air drop, so claim module should not send the balances to the user after delegate.
	s.True(balances.Empty())

	s.app.ClaimKeeper.AfterDelegationModified(s.ctx.WithBlockTime(airdropStartTime), addr1, sdk.ValAddress(addr1))
	balances = s.app.BankKeeper.GetAllBalances(s.ctx, addr1)
	//fmt.Printf("%v \n", balances)
	// Now, it is the time for air drop, so claim module should send the balances to the user after delegate.
	s.Equal(claimRecords[0].MaximumClaimableAmount.Amount.Quo(math.NewInt(int64(len(types.Action_value)))).Quo(math.NewInt(types.ClaimsPortions)), balances.AmountOf(sdk.DefaultBondDenom))
}

func (s *KeeperTestSuite) TestAirdropDisabled() {
	s.SetupTest()

	airdropStartTime := time.Now().Add(time.Hour)

	s.app.ClaimKeeper.SetParams(s.ctx, types.Params{
		ClaimDenom:             sdk.DefaultBondDenom,
		DurationUntilDecay:     time.Hour,
		DurationOfDecay:        time.Hour * 4,
		DurationVestingPeriods: types.DefaultDurationVestingPeriods,
	})
	addr1 := s.createAccount()
	claimRecords := s.createClaimRecords(addr1, 2000)

	err := s.app.ClaimKeeper.SetClaimRecords(s.ctx, claimRecords)
	s.Require().NoError(err)

	coins, err := s.app.ClaimKeeper.GetTotalClaimableForAddr(s.ctx, addr1)
	s.NoError(err)
	// Now, it is before starting air drop, so this value should return the empty coins
	s.True(coins.IsZero())

	total, err := s.app.ClaimKeeper.GetTotalClaimableAmountPerAction(s.ctx, addr1)
	s.NoError(err)
	// Now, it is before starting air drop, so this value should return the empty coins
	s.True(total.IsZero())

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
	s.Equal(claimRecords[0].MaximumClaimableAmount.Amount.Quo(math.NewInt(int64(len(types.Action_value)))).Quo(math.NewInt(types.ClaimsPortions)), balances.AmountOf(sdk.DefaultBondDenom))
}
func (s *KeeperTestSuite) TestDuplicatedActionNotWithdrawRepeatedly() {
	s.SetupTest()

	addr1 := s.createAccount()
	claimRecords := s.createClaimRecords(addr1, 2000)

	err := s.app.ClaimKeeper.SetClaimRecords(s.ctx, claimRecords)
	s.Require().NoError(err)

	// Initial claimable amount
	initialCoins, err := s.app.ClaimKeeper.GetTotalClaimableForAddr(s.ctx, addr1)
	s.Require().NoError(err)
	s.Require().Equal(initialCoins, claimRecords[0].MaximumClaimableAmount)

	// First action triggers claim
	s.triggerDelegationAction(addr1)
	claim := s.getClaimRecord(addr1)
	s.True(claim.Status[3].ActionCompleted)

	balance := s.app.BankKeeper.GetAllBalances(s.ctx, addr1)
	expectedClaim := claimRecords[0].MaximumClaimableAmount.Amount.Quo(math.NewInt(20))
	s.Require().Equal(expectedClaim, balance.AmountOf(sdk.DefaultBondDenom))

	// Repeat action should not double the claim
	s.triggerDelegationAction(addr1)
	claim = s.getClaimRecord(addr1)
	s.True(claim.Status[3].ActionCompleted)

	balance = s.app.BankKeeper.GetAllBalances(s.ctx, addr1)
	s.Require().Equal(expectedClaim, balance.AmountOf(sdk.DefaultBondDenom))
}

func (s *KeeperTestSuite) TestNotRunningGenesisBlock() {
	s.ctx = s.ctx.WithBlockHeight(1)

	s.app.ClaimKeeper.SetParams(s.ctx, types.Params{
		AirdropStartTime:       time.Now().Add(-time.Hour),
		ClaimDenom:             sdk.DefaultBondDenom,
		DurationUntilDecay:     time.Hour,
		DurationOfDecay:        time.Hour * 4,
		DurationVestingPeriods: types.DefaultDurationVestingPeriods,
	})

	addr1 := s.createAccount()
	claimRecords := s.createClaimRecords(addr1, 2000)

	err := s.app.ClaimKeeper.SetClaimRecords(s.ctx, claimRecords)
	s.Require().NoError(err)

	// Initial claimable amount
	initialCoins, err := s.app.ClaimKeeper.GetTotalClaimableForAddr(s.ctx, addr1)
	s.Require().NoError(err)
	s.Require().Equal(initialCoins, claimRecords[0].MaximumClaimableAmount)

	// Action should mark claim as completed
	s.triggerDelegationAction(addr1)
	claim := s.getClaimRecord(addr1)
	s.True(claim.Status[3].ActionCompleted)

	// Claimable amount remains consistent
	finalCoins, err := s.app.ClaimKeeper.GetTotalClaimableForAddr(s.ctx, addr1)
	s.Require().NoError(err)
	s.Require().Equal(finalCoins, claimRecords[0].MaximumClaimableAmount)
}

func (s *KeeperTestSuite) TestEndAirdrop() {
	s.app.ClaimKeeper.SetParams(s.ctx, types.Params{
		AirdropStartTime:       time.Now().Add(-time.Hour),
		ClaimDenom:             sdk.DefaultBondDenom,
		DurationUntilDecay:     time.Hour,
		DurationOfDecay:        time.Hour * 4,
		DurationVestingPeriods: types.DefaultDurationVestingPeriods,
	})

	addr1 := s.createAccount()
	addr2 := s.createAccount()
	claimRecords := []types.ClaimRecord{
		s.createClaimRecord(addr1, 1000),
		s.createClaimRecord(addr2, 1000),
	}

	err := s.app.ClaimKeeper.SetClaimRecords(s.ctx, claimRecords)
	s.Require().NoError(err)

	// End the airdrop
	err = s.app.ClaimKeeper.EndAirdrop(s.ctx)
	s.Require().NoError(err)

	// Module account should have no remaining balance
	moduleAccAddr := s.app.AccountKeeper.GetModuleAddress(types.ModuleName)
	moduleBalance := s.app.BankKeeper.GetBalance(s.ctx, moduleAccAddr, sdk.DefaultBondDenom)
	s.Require().Equal(sdk.NewInt64Coin(sdk.DefaultBondDenom, 0).String(), moduleBalance.String())
}

func (s *KeeperTestSuite) createAccount() sdk.AccAddress {
	pubKey := ed25519.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubKey.Address())
	s.app.AccountKeeper.NewAccount(s.ctx, authtypes.NewBaseAccount(addr, nil, 0, 0))
	return addr
}

func (s *KeeperTestSuite) createClaimRecords(addr sdk.AccAddress, amount int64) []types.ClaimRecord {
	status := types.Status{
		ActionCompleted:         false,
		VestingPeriodsCompleted: []bool{false, false, false, false},
		VestingPeriodsClaimed:   []bool{false, false, false, false},
	}
	return []types.ClaimRecord{
		{
			Address:                addr.String(),
			MaximumClaimableAmount: sdk.NewInt64Coin(sdk.DefaultBondDenom, amount),
			Status:                 []types.Status{status, status, status, status},
		},
	}
}

func (s *KeeperTestSuite) createClaimRecord(addr sdk.AccAddress, amount int64) types.ClaimRecord {
	return s.createClaimRecords(addr, amount)[0]
}

func (s *KeeperTestSuite) triggerDelegationAction(addr sdk.AccAddress) {
	s.app.ClaimKeeper.AfterDelegationModified(s.ctx, addr, sdk.ValAddress(addr))
}

func (s *KeeperTestSuite) getClaimRecord(addr sdk.AccAddress) types.ClaimRecord {
	claim, err := s.app.ClaimKeeper.GetClaimRecord(s.ctx, addr)
	s.Require().NoError(err)
	return claim
}
func (s *KeeperTestSuite) TestVestingQueue() {
	addr1 := s.createAccount()
	action := types.ACTION_ACTION_LOCAL

	// Create initial claim record
	status := types.Status{
		ActionCompleted:         true,
		VestingPeriodsCompleted: []bool{false, false, false, false},
		VestingPeriodsClaimed:   []bool{false, false, false, false},
	}
	record := types.ClaimRecord{
		Address: addr1.String(),
		Status:  []types.Status{status, status, status, status},
	}
	s.Require().NoError(s.app.ClaimKeeper.SetClaimRecord(s.ctx, record))

	// Insert 4 vesting periods
	s.Require().NoError(s.app.ClaimKeeper.InsertEntriesIntoVestingQueue(s.ctx, addr1.String(), byte(action), s.ctx.BlockTime()))

	params, _ := s.app.ClaimKeeper.GetParams(s.ctx)
	vestDuration := params.DurationVestingPeriods[byte(action)]

	// Process all periods
	for period := 0; period < 4; period++ {
		// advance block time to after vesting period
		s.ctx = s.ctx.WithBlockTime(s.ctx.BlockTime().Add(vestDuration))

		s.app.ClaimKeeper.IterateVestingQueue(s.ctx, s.ctx.BlockHeader().Time, func(recipientAddr sdk.AccAddress, a int32, p int32, endTime time.Time) bool {
			claimRecord, err := s.app.ClaimKeeper.GetClaimRecord(s.ctx, recipientAddr)
			if err != nil {
				panic("Failed to get claim record")
			}

			// mark vesting period completed
			if int(action) < len(claimRecord.Status) &&
				int(period) < len(claimRecord.Status[action].VestingPeriodsCompleted) {
				claimRecord.Status[action].VestingPeriodsCompleted[period] = true
			}

			// persist record
			if err := s.app.ClaimKeeper.SetClaimRecord(s.ctx, claimRecord); err != nil {
				panic("Failed to set claim record")
			}

			// remove from queue
			s.app.ClaimKeeper.RemoveEntryFromVestingQueue(s.ctx, recipientAddr.String(), endTime, byte(action), byte(period))

			return false // keep iterating
		})

		// Verify completion
		updated, err := s.app.ClaimKeeper.GetClaimRecord(s.ctx, addr1)
		s.Require().NoError(err)
		s.Require().True(updated.Status[action].VestingPeriodsCompleted[period], "Period %d not completed", period)

		// Verify queue removal
		found := false
		s.app.ClaimKeeper.IterateVestingQueue(s.ctx, s.ctx.BlockTime(), func(addr sdk.AccAddress, a int32, p int32, endTime time.Time) bool {
			if a == int32(action) && p == int32(period) {
				found = true
				return true
			}
			return false
		})
		s.Require().False(found, "Queue entry for period %d not removed", period)
	}
}
