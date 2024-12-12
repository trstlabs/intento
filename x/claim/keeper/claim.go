package keeper

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trstlabs/intento/x/claim/types"
)

// ClaimClaimableForAddr refactored for improved clarity, edge case handling, and atomicity.
func (k Keeper) ClaimClaimableForAddr(ctx sdk.Context, addr sdk.AccAddress) error {
	record, err := k.GetClaimRecord(ctx, addr)
	if err != nil {
		return err
	}

	p, err := k.GetParams(ctx)
	if err != nil {
		return err
	}

	moduleAccountBalance := k.GetModuleAccountBalance(ctx)

	// Local copy for mutation, updated atomically
	updatedRecord := record

	totalClaimableAmountForAction, err := k.GetTotalClaimableAmountPerAction(ctx, addr)
	if err != nil {
		return err
	}
	if totalClaimableAmountForAction.IsZero() {
		return errorsmod.Wrap(sdkerrors.ErrNotFound, "address does not have claimable tokens")
	}

	claimableCoin := sdk.NewCoin(p.ClaimDenom, math.ZeroInt())
	claimedCoin := sdk.NewCoin(p.ClaimDenom, math.ZeroInt())
	var toClaimPeriods int64 = 0
	var claimedPeriods int64 = 0
	for action, status := range updatedRecord.Status {
		if !status.ActionCompleted {
			continue
		}

		var toClaimPeriodsForAction int64 = 0
		var claimedPeriodsForAction int64 = 1
		for period, completed := range status.VestingPeriodsCompleted {
			//fmt.Printf("period %v ; completed: %v\n", period, completed)
			if completed && !status.VestingPeriodsClaimed[period] {
				toClaimPeriodsForAction++
				updatedRecord.Status[action].VestingPeriodsClaimed[period] = true
			} else if completed && status.VestingPeriodsClaimed[period] {
				claimedPeriodsForAction++
			}

		}

		if toClaimPeriodsForAction != 0 {
			toClaimPercent := math.LegacyNewDec(toClaimPeriodsForAction).Quo(math.LegacyNewDec(5))
			claimableTotalDec := math.LegacyNewDecFromInt(totalClaimableAmountForAction)
			claimableDec := claimableTotalDec.Mul(toClaimPercent)
			claimableCoin = claimableCoin.AddAmount(claimableDec.TruncateInt())
			toClaimPeriods = toClaimPeriods + toClaimPeriodsForAction
		}
		claimedPart := math.LegacyNewDec(claimedPeriodsForAction).Quo(math.LegacyNewDec(5))
		claimedCoin = claimedCoin.AddAmount(math.LegacyNewDecFromInt(totalClaimableAmountForAction).Mul(claimedPart).TruncateInt())
		claimedPeriods = claimedPeriods + claimedPeriodsForAction

	}

	if toClaimPeriods == 0 || claimedCoin.Amount == math.ZeroInt() {
		return errorsmod.Wrap(sdkerrors.ErrNotFound, "address does not have claimable tokens right now")
	}
	// Perform staking check before transferring
	delegationInfo, err := k.stakingKeeper.GetAllDelegatorDelegations(ctx, addr)
	if err != nil {
		return err
	}

	totalDelegations := math.LegacyZeroDec()
	for _, delegation := range delegationInfo {
		totalDelegations = totalDelegations.Add(delegation.Shares)
	}

	minBonded := math.LegacyNewDecWithPrec(67, 2).MulInt(claimableCoin.Amount)
	if totalDelegations.Sub(minBonded).IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInsufficientFunds, "address does not have enough tokens staked to claim: staked "+totalDelegations.BigInt().String()+"required: "+minBonded.BigInt().String()+"claimable :"+claimableCoin.Amount.String())
	}

	// Transfer claimable amount to the user
	err = k.TransferToUser(ctx, addr, claimableCoin, moduleAccountBalance.Amount, p.ClaimDenom)
	if err != nil {
		return err
	}

	// Update the record atomically
	k.SetClaimRecord(ctx, updatedRecord)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaim,
			sdk.NewAttribute(sdk.AttributeKeySender, addr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, claimableCoin.String()),
		),
	})

	return nil
}

// GetTotalClaimableAmountPerAction
func (k Keeper) GetTotalClaimableAmountPerAction(ctx sdk.Context, addr sdk.AccAddress) (math.Int, error) {
	record, err := k.GetClaimRecord(ctx, addr)
	if err != nil {
		return math.ZeroInt(), errors.New("claim record not found")
	}

	p, err := k.GetParams(ctx)
	if err != nil {
		return math.ZeroInt(), err
	}

	// If we are before the start time, do nothing.
	// This case _shouldn't_ occur on chain, since the
	// start time ought to be chain start time.
	if ctx.BlockTime().Before(p.AirdropStartTime) {
		return math.ZeroInt(), nil
	}

	totalDelegations := k.GetTotalDelegations(ctx, addr)

	timeElapsed := ctx.BlockTime().Sub(p.AirdropStartTime)
	timeLeft := timeElapsed - p.DurationUntilDecay + p.DurationOfDecay
	if timeElapsed < p.DurationUntilDecay {
		return record.MaximumClaimableAmount.Amount.Quo(math.NewInt(int64(len(types.Action_name)))), nil
	}

	if timeLeft < 0 {
		timeLeft = 0
	}
	decayPercent := math.LegacyNewDecFromInt(math.NewInt((int64(timeLeft)))).QuoInt64(int64(p.DurationOfDecay)).Ceil().TruncateInt()

	baseAmount := record.MaximumClaimableAmount.Amount.QuoRaw(int64(len(types.Action_name)))
	adjustedAmount := baseAmount.Mul(decayPercent)

	// Delegation adjustment
	if totalDelegations.IsPositive() {
		adjustedAmount = adjustedAmount.Mul(totalDelegations)
	}

	return adjustedAmount, nil
}

// TransferToUser safely handles fund transfers and vesting setup.
func (k Keeper) TransferToUser(ctx sdk.Context, addr sdk.AccAddress, claimableCoin sdk.Coin, moduleAccountBalance math.Int, denom string) error {
	if claimableCoin.IsZero() || claimableCoin.Amount.GT(moduleAccountBalance) {
		return errors.New("insufficient funds in module account")
	}

	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.NewCoins(claimableCoin))
	if err != nil {
		return err
	}

	// Additional vesting logic if necessary (placeholder)
	// k.SetupVestingSchedule(ctx, addr, amount, vestingPeriod)

	return nil
}

// GetTotalDelegations retrieves total delegations for an address.
func (k Keeper) GetTotalDelegations(ctx sdk.Context, addr sdk.AccAddress) math.Int {
	delegations, err := k.stakingKeeper.GetDelegatorDelegations(ctx, addr, 10)
	if err != nil {
		return math.ZeroInt()
	}
	total := math.ZeroInt()
	for _, delegation := range delegations {
		total = total.Add(delegation.Shares.RoundInt())
	}

	return total
}

// GetModuleAccountBalance gets the airdrop coin balance of module account
func (k Keeper) GetModuleAccountBalance(ctx sdk.Context) sdk.Coin {
	moduleAccAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	p, err := k.GetParams(ctx)
	if err != nil {
		return sdk.Coin{}
	}

	return k.bankKeeper.GetBalance(ctx, moduleAccAddr, p.ClaimDenom)
}

// ClaimInitialCoinsForAction remove claimable amount entry and transfer it to recipient's account
func (k Keeper) ClaimInitialCoinsForAction(ctx sdk.Context, addr sdk.AccAddress, action types.Action) error {
	claimRecord, err := k.GetClaimRecord(ctx, addr)
	if err != nil {
		return err
	}

	p, err := k.GetParams(ctx)
	if err != nil {
		return err
	}

	if claimRecord.Status[action].ActionCompleted {
		return nil //errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "address already claimed tokens for this action")
	}
	claimable, err := k.GetTotalClaimableAmountPerAction(ctx, addr)
	if err != nil {
		return err
	}

	if claimable.IsZero() {
		return nil
	}

	//we distribute 20% after action completion and the remaining become claimable after each vesting period
	claimsPortionCoins := sdk.NewCoins(sdk.NewCoin(p.ClaimDenom, claimable.QuoRaw(types.ClaimsPortions)))

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, claimsPortionCoins)
	if err != nil {
		return err
	}

	//creates entries into the endblocker ( 4 )
	err = k.InsertEntriesIntoVestingQueue(ctx, addr.String(), byte(action), ctx.BlockHeader().Time)
	if err != nil {
		return err
	}

	//set claim record
	claimRecord.Status[action].ActionCompleted = true
	err = k.SetClaimRecord(ctx, claimRecord)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaim,
			sdk.NewAttribute("claimer", addr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, claimable.String()),
		),
	})

	return nil
}

// GetTotalClaimableForAddr returns claimable amount for a specific action done by an address
func (k Keeper) GetTotalClaimableForAddr(ctx sdk.Context, addr sdk.AccAddress) (sdk.Coin, error) {
	claimRecord, err := k.GetClaimRecord(ctx, addr)
	if err != nil {
		return sdk.Coin{}, err
	}
	if claimRecord.Address == "" {
		return sdk.Coin{}, nil
	}
	p, err := k.GetParams(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}

	claimableForAction, err := k.GetTotalClaimableAmountPerAction(ctx, addr)
	if err != nil {
		return sdk.Coin{}, err
	}
	totalClaimable := claimableForAction.MulRaw(int64(len(types.Action_name)))

	return sdk.NewCoin(p.ClaimDenom, totalClaimable), nil
}

func (k Keeper) EndAirdrop(ctx sdk.Context) error {
	err := k.fundRemainingsToCommunity(ctx)
	if err != nil {
		return err
	}
	k.clearInitialClaimables(ctx)
	return nil
}

// FundRemainingsToCommunity fund remainings to the community when airdrop period ends
func (k Keeper) fundRemainingsToCommunity(ctx sdk.Context) error {
	moduleAccAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	amt := k.GetModuleAccountBalance(ctx)
	return k.distrKeeper.FundCommunityPool(ctx, sdk.NewCoins(amt), moduleAccAddr)
}
