package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/trstlabs/intento/x/claim/types"
)

func (k Keeper) AfterActionAuthz(ctx sdk.Context, recipient sdk.AccAddress) {
	k.ClaimInitialCoinsForAction(ctx, recipient, types.ActionActionAuthz)
}

func (k Keeper) AfterActionWasm(ctx sdk.Context, recipient sdk.AccAddress) {
	k.ClaimInitialCoinsForAction(ctx, recipient, types.ActionActionWasm)
}

func (k Keeper) AfterGovernanceVoted(ctx sdk.Context, recipient sdk.AccAddress) {
	k.ClaimInitialCoinsForAction(ctx, recipient, types.ActionGovernanceVote)
}

func (k Keeper) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	k.ClaimInitialCoinsForAction(ctx, delAddr, types.ActionDelegateStake)
}

/* func (k Keeper) AfterAutoSwap(ctx sdk.Context, recipient sdk.AccAddress) {
	k.ClaimInitialCoinsForAction(ctx, recipient, types.ActionAutoSwap)

}

func (k Keeper) AfterRecurringSend(ctx sdk.Context, recipient sdk.AccAddress) {
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
func (h Hooks) AfterProposalSubmission(ctx sdk.Context, proposalID uint64) {}
func (h Hooks) AfterProposalDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress) {

}
func (h Hooks) AfterProposalVotingPeriodEnded(ctx sdk.Context, proposalID uint64) {}
func (h Hooks) AfterProposalFailedMinDeposit(ctx sdk.Context, proposalID uint64)  {}

func (h Hooks) AfterProposalVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {
	h.k.AfterGovernanceVoted(ctx, voterAddr)

}

func (h Hooks) AfterProposalInactive(ctx sdk.Context, proposalID uint64) {}
func (h Hooks) AfterProposalActive(ctx sdk.Context, proposalID uint64)   {}

// staking hooks
func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterUnbondingInitiated(ctx sdk.Context, id uint64) error {
	return nil
}
func (h Hooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	h.k.AfterDelegationModified(ctx, delAddr, valAddr)
	return nil
}
func (h Hooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) error {
	return nil
}

// AutoIBCTX hooks
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


func (k Keeper) AfterItemTokenized(ctx sdk.Context, creator sdk.AccAddress) {
    _, err := k.ClaimInitialCoinsForAction(ctx, creator, types.ActionItemTokenized)
    if err != nil {
        panic(err.Error())
    }
}*/

//var _ itemtypes.ItemHooks = Hooks{}

/*
// item hooks

func (h Hooks) AfterItemTokenized(ctx sdk.Context, recipientAddr sdk.AccAddress) {
	//h.k.AfterItemTokenized(ctx, recipientAddr)
}
func (h Hooks) AfterItemBought(ctx sdk.Context, recipientAddr sdk.AccAddress) {
	//h.k.AfterItemBought(ctx, recipientAddr)
}
//func (h Hooks) AfterItemEstimated(ctx sdk.Context, proposalID uint64) {}
*/
