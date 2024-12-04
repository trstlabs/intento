package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/trstlabs/intento/x/claim/types"
)

func (k Keeper) AfterActionAuthz(ctx sdk.Context, recipient sdk.AccAddress) {
	err := k.ClaimInitialCoinsForAction(ctx, recipient, types.ActionActionAuthz)
	// fmt.Printf("ClaimInitialCoinsForAction err %v \n", err)
	k.Logger(ctx).Debug("ClaimInitialCoinsForAction", "error", err, "for", recipient)
}

func (k Keeper) AfterActionWasm(ctx sdk.Context, recipient sdk.AccAddress) {
	err := k.ClaimInitialCoinsForAction(ctx, recipient, types.ActionActionWasm)
	k.Logger(ctx).Debug("ClaimInitialCoinsForAction", "error", err, "for", recipient)
}

func (k Keeper) AfterGovernanceVoted(ctx sdk.Context, recipient sdk.AccAddress) {
	err := k.ClaimInitialCoinsForAction(ctx, recipient, types.ActionGovernanceVote)
	k.Logger(ctx).Debug("ClaimInitialCoinsForAction", "error", err, "for", recipient)
}

func (k Keeper) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	err := k.ClaimInitialCoinsForAction(ctx, delAddr, types.ActionDelegateStake)
	k.Logger(ctx).Debug("ClaimInitialCoinsForAction", "error", err, "for", delAddr)
}

/* func (k Keeper) AfterAutoSwap(ctx context.Context, recipient sdk.AccAddress) {
	k.ClaimInitialCoinsForAction(ctx, recipient, types.ActionAutoSwap)

}

func (k Keeper) AfterRecurringSend(ctx context.Context, recipient sdk.AccAddress) {
	k.ClaimInitialCoinsForAction(ctx, recipient, types.ActionRecurringSend)
} */

// ________________________________________________________________________________________

// Hooks wrapper struct for claims keeper
type Hooks struct {
	k Keeper
}

var _ stakingtypes.StakingHooks = Hooks{}

var _ govtypes.GovHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// governance hooks
func (h Hooks) AfterProposalSubmission(ctx context.Context, proposalID uint64) error {
	return nil
}
func (h Hooks) AfterProposalDeposit(ctx context.Context, proposalID uint64, depositorAddr sdk.AccAddress) error {
	return nil
}
func (h Hooks) AfterProposalVotingPeriodEnded(ctx context.Context, proposalID uint64) error {
	return nil
}
func (h Hooks) AfterProposalFailedMinDeposit(ctx context.Context, proposalID uint64) error {
	return nil
}

func (h Hooks) AfterProposalVote(ctx context.Context, proposalID uint64, voterAddr sdk.AccAddress) error {
	h.k.AfterGovernanceVoted(sdk.UnwrapSDKContext(ctx), voterAddr)
	return nil
}

func (h Hooks) AfterProposalInactive(ctx context.Context, proposalID uint64) {}
func (h Hooks) AfterProposalActive(ctx context.Context, proposalID uint64)   {}

// staking hooks
func (h Hooks) AfterValidatorCreated(ctx context.Context, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeValidatorModified(ctx context.Context, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterValidatorRemoved(ctx context.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterValidatorBonded(ctx context.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterValidatorBeginUnbonding(ctx context.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterUnbondingInitiated(ctx context.Context, id uint64) error {
	return nil
}
func (h Hooks) BeforeDelegationCreated(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeDelegationSharesModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeDelegationRemoved(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterDelegationModified(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	h.k.AfterDelegationModified(sdk.UnwrapSDKContext(ctx), delAddr, valAddr)
	return nil
}
func (h Hooks) BeforeValidatorSlashed(ctx context.Context, valAddr sdk.ValAddress, fraction math.LegacyDec) error {
	return nil
}

// intent hooks
func (h Hooks) AfterActionAuthz(ctx sdk.Context, recipientAddr sdk.AccAddress) {
	h.k.AfterActionAuthz(ctx, recipientAddr)
}
func (h Hooks) AfterActionWasm(ctx sdk.Context, recipientAddr sdk.AccAddress) {
	h.k.AfterActionWasm(ctx, recipientAddr)
}

func (h Hooks) AfterAutoSwap(ctx sdk.Context, recipientAddr sdk.AccAddress) {
	//h.k.AfterAutoSwap(ctx, recipientAddr)
}
func (h Hooks) AfterRecurringSend(ctx sdk.Context, recipientAddr sdk.AccAddress) {
	// h.k.AfterRecurringSend(ctx, recipientAddr)
}

// ________________________________________________________________________________________

// for future reference
/*

// Compute hooks


func (k Keeper) AfterItemTokenized(ctx context.Context, creator sdk.AccAddress) {
    _, err := k.ClaimInitialCoinsForAction(ctx, creator, types.ActionItemTokenized)
    if err != nil {
        panic(err.Error())
    }
}*/

//var _ itemtypes.ItemHooks = Hooks{}

/*
// item hooks

func (h Hooks) AfterItemTokenized(ctx context.Context, recipientAddr sdk.AccAddress) {
	//h.k.AfterItemTokenized(ctx, recipientAddr)
}
func (h Hooks) AfterItemBought(ctx context.Context, recipientAddr sdk.AccAddress) {
	//h.k.AfterItemBought(ctx, recipientAddr)
}
//func (h Hooks) AfterItemEstimated(ctx context.Context, proposalID uint64) {}
*/
