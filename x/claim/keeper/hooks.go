package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/trstlabs/trst/x/claim/types"
	itemtypes "github.com/trstlabs/trst/x/item/types"
)

func (k Keeper) AfterComputeExecuted(ctx sdk.Context, sender sdk.AccAddress) {
	_, err := k.ClaimCoinsForAction(ctx, sender, types.ActionComputeExecute, "60s") //"1440h") //vest for 240 days /8 months
	if err != nil {
		panic(err.Error())
	}
}

func (k Keeper) AfterComputeInstantiated(ctx sdk.Context, sender sdk.AccAddress) {
	_, err := k.ClaimCoinsForAction(ctx, sender, types.ActionComputeInstantiate, "60s") // "168h") //vest for 4 weeks
	if err != nil {
		panic(err.Error())
	}
}

func (k Keeper) AfterItemBought(ctx sdk.Context, sender sdk.AccAddress) {
	_, err := k.ClaimCoinsForAction(ctx, sender, types.ActionItemBought, "336h") //vest for 8 weeks
	if err != nil {
		panic(err.Error())
	}
}

/*
func (k Keeper) AfterItemTokenized(ctx sdk.Context, creator sdk.AccAddress) {
	_, err := k.ClaimCoinsForAction(ctx, creator, types.ActionItemTokenized)
	if err != nil {
		panic(err.Error())
	}
}*/

func (k Keeper) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	_, err := k.ClaimCoinsForAction(ctx, delAddr, types.ActionDelegateStake, "60s") //"720h") //vest for 120 days/ 4 months
	if err != nil {
		panic(err.Error())
	}
}

// ________________________________________________________________________________________

// Hooks wrapper struct for claims keeper
type Hooks struct {
	k Keeper
}

var _ stakingtypes.StakingHooks = Hooks{}

var _ itemtypes.ItemHooks = Hooks{}

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

}

func (h Hooks) AfterProposalInactive(ctx sdk.Context, proposalID uint64) {}
func (h Hooks) AfterProposalActive(ctx sdk.Context, proposalID uint64)   {}

// staking hooks
func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress)   {}
func (h Hooks) BeforeValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) {}
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.k.AfterDelegationModified(ctx, delAddr, valAddr)
}
func (h Hooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {}

// item hooks

func (h Hooks) AfterItemTokenized(ctx sdk.Context, senderAddr sdk.AccAddress) {
	//h.k.AfterItemTokenized(ctx, senderAddr)
}
func (h Hooks) AfterItemBought(ctx sdk.Context, senderAddr sdk.AccAddress) {
	h.k.AfterItemBought(ctx, senderAddr)
}

// Compute hooks
func (h Hooks) AfterComputeExecuted(ctx sdk.Context, senderAddr sdk.AccAddress) {
	h.k.AfterComputeExecuted(ctx, senderAddr)
}
func (h Hooks) AfterComputeInstantiated(ctx sdk.Context, senderAddr sdk.AccAddress) {
	h.k.AfterComputeInstantiated(ctx, senderAddr)
}

//func (h Hooks) AfterItemEstimated(ctx sdk.Context, proposalID uint64) {}
