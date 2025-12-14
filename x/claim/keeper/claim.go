package keeper

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trstlabs/intento/x/claim/types"
)

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

	// Local copy for atomic update
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

	// Track the amount claimed after the first 1/5 portion
	claimedAfterFirstPeriod := math.ZeroInt()

	for action, status := range updatedRecord.Status {
		if !status.ActionCompleted {
			continue
		}

		var toClaimPeriodsForAction int64 = 0
		var claimedPeriodsForAction int64 = 0

		for period, completed := range status.VestingPeriodsCompleted {
			if completed && !status.VestingPeriodsClaimed[period] {
				toClaimPeriodsForAction++
				updatedRecord.Status[action].VestingPeriodsClaimed[period] = true
			} else if completed && status.VestingPeriodsClaimed[period] {
				claimedPeriodsForAction++

				// Only count periods after the first for staking requirement
				if period != 0 {
					portion := math.LegacyNewDecFromInt(totalClaimableAmountForAction).Quo(math.LegacyNewDec(5))
					claimedAfterFirstPeriod = claimedAfterFirstPeriod.Add(portion.TruncateInt())
				}
			}
		}

		// Calculate claimable for this action
		if toClaimPeriodsForAction != 0 {
			toClaimPercent := math.LegacyNewDec(toClaimPeriodsForAction).Quo(math.LegacyNewDec(5))
			claimableDec := math.LegacyNewDecFromInt(totalClaimableAmountForAction).Mul(toClaimPercent)
			claimableCoin = claimableCoin.AddAmount(claimableDec.TruncateInt())
			toClaimPeriods += toClaimPeriodsForAction
		}

		// Track total claimed (for reporting, not staking check)
		claimedPart := math.LegacyNewDec(claimedPeriodsForAction).Quo(math.LegacyNewDec(5))
		claimedCoin = claimedCoin.AddAmount(math.LegacyNewDecFromInt(totalClaimableAmountForAction).Mul(claimedPart).TruncateInt())
		claimedPeriods += claimedPeriodsForAction
	}

	if toClaimPeriods == 0 || claimableCoin.Amount.IsZero() {
		return errorsmod.Wrap(sdkerrors.ErrNotFound, "address does not have claimable tokens right now")
	}

	// Compute total staked tokens
	totalStaked := math.ZeroInt()
	delegations, _ := k.stakingKeeper.GetAllDelegatorDelegations(ctx, addr)
	for _, del := range delegations {
		valAddr, err := sdk.ValAddressFromBech32(del.ValidatorAddress)
		if err != nil {
			continue
		}
		val, err := k.stakingKeeper.GetValidator(ctx, valAddr)
		if err != nil {
			continue
		}
		tokens := val.TokensFromShares(del.Shares)
		totalStaked = totalStaked.Add(tokens.TruncateInt())
	}

	// Apply staking minimum only for claims after the first 1/5
	if !claimedAfterFirstPeriod.IsZero() {
		minRequired := claimedAfterFirstPeriod.MulRaw(67).QuoRaw(100)
		if totalStaked.LT(minRequired) {
			return errorsmod.Wrapf(
				sdkerrors.ErrInsufficientFunds,
				"address does not have enough tokens staked to claim: staked: %s, required: %s",
				totalStaked.String(), minRequired.String(),
			)
		}
	}

	// Transfer claimable tokens to the user
	if err := k.TransferToUser(ctx, addr, claimableCoin, moduleAccountBalance.Amount, p.ClaimDenom); err != nil {
		return err
	}

	// Atomically update claim record
	k.SetClaimRecord(ctx, updatedRecord)

	// Emit claim event
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaim,
			sdk.NewAttribute(sdk.AttributeKeySender, addr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, claimableCoin.String()),
		),
	})

	return nil
}

func (k Keeper) GetTotalClaimableAmountPerAction(ctx sdk.Context, addr sdk.AccAddress) (math.Int, error) {
	record, err := k.GetClaimRecord(ctx, addr)
	if err != nil {
		return math.ZeroInt(), errors.New("claim record not found")
	}

	// Get module params (contains airdrop start, decay durations, denom, etc.)
	p, err := k.GetParams(ctx)
	if err != nil {
		return math.ZeroInt(), err
	}

	// Before airdrop start → nothing is claimable
	if ctx.BlockTime().Before(p.AirdropStartTime) {
		return math.ZeroInt(), nil
	}

	// Calculate elapsed time since airdrop started
	timeElapsed := ctx.BlockTime().Sub(p.AirdropStartTime)

	// Base amount per action (divide total allocation by number of actions)
	baseAmount := record.MaximumClaimableAmount.Amount.Quo(math.NewInt(int64(len(types.Action_name))))

	// If we are still before the decay period starts → full base amount is claimable
	if timeElapsed <= p.DurationUntilDecay {
		return baseAmount, nil
	}

	// Time since decay started
	decayElapsed := timeElapsed - p.DurationUntilDecay
	if decayElapsed < 0 {
		decayElapsed = 0 // safety check
	}

	// Linear decay fraction
	// Fraction of claimable remaining: 1 → just started decay, 0 → fully decayed
	remainingFrac := math.LegacyNewDecFromInt(
		math.NewInt(int64(p.DurationOfDecay - decayElapsed)),
	).Quo(math.LegacyNewDecFromInt(math.NewInt(int64(p.DurationOfDecay))))

	// Clamp fraction to [0,1] to avoid negative or >100%
	if remainingFrac.LT(math.LegacyZeroDec()) {
		remainingFrac = math.LegacyZeroDec()
	} else if remainingFrac.GT(math.LegacyOneDec()) {
		remainingFrac = math.LegacyOneDec()
	}

	// Apply decay fraction to base amount
	decayedAmount := math.LegacyNewDecFromInt(baseAmount).Mul(remainingFrac).TruncateInt()

	// Final adjusted amount — no delegation multiplier applied
	adjustedAmount := decayedAmount

	// Safety check: Ensure we never exceed base allocation per action
	if adjustedAmount.GT(baseAmount) {
		adjustedAmount = baseAmount
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
	k.clearVestingQueue(ctx)
	return nil
}

// FundRemainingsToCommunity fund remainings to the community when airdrop period ends
func (k Keeper) fundRemainingsToCommunity(ctx sdk.Context) error {
	moduleAccAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	amt := k.GetModuleAccountBalance(ctx)
	return k.distrKeeper.FundCommunityPool(ctx, sdk.NewCoins(amt), moduleAccAddr)
}
