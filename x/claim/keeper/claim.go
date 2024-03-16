package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/trstlabs/intento/x/claim/types"
)

// GetModuleAccountAddress gets the airdrop coin balance of module account
func (k Keeper) GetModuleAccountAddress(ctx sdk.Context) sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleName)
}

// GetModuleAccountBalance gets the airdrop coin balance of module account
func (k Keeper) GetModuleAccountBalance(ctx sdk.Context) sdk.Coin {
	moduleAccAddr := k.GetModuleAccountAddress(ctx)
	params, err := k.GetParams(ctx)
	if err != nil {
		return sdk.Coin{}
	}
	return k.bankKeeper.GetBalance(ctx, moduleAccAddr, params.ClaimDenom)
}

// CreateModuleAccount sets balance of airdrop module
func (k Keeper) CreateModuleAccount(ctx sdk.Context, amount sdk.Coin) {
	moduleAcc := authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Minter)
	k.accountKeeper.SetModuleAccount(ctx, moduleAcc)

	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(amount))
	if err != nil {
		panic(err)
	}
}

func (k Keeper) EndAirdrop(ctx sdk.Context) error {
	err := k.fundRemainingsToCommunity(ctx)
	if err != nil {
		return err
	}
	k.clearInitialClaimables(ctx)
	return nil
}

// clearInitialClaimables clear claimable amounts
func (k Keeper) clearInitialClaimables(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.ClaimRecordsStorePrefix))
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		store.Delete(key)
	}
}

// SetClaimRecords sets claim records
func (k Keeper) SetClaimRecords(ctx sdk.Context, claimRecords []types.ClaimRecord) error {
	for _, claimRecord := range claimRecords {
		err := k.SetClaimRecord(ctx, claimRecord)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetClaimRecords get claim records for genesis export
func (k Keeper) GetClaimRecords(ctx sdk.Context) []types.ClaimRecord {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, []byte(types.ClaimRecordsStorePrefix))

	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	claimRecords := []types.ClaimRecord{}
	for ; iterator.Valid(); iterator.Next() {

		claimRecord := types.ClaimRecord{}

		err := proto.Unmarshal(iterator.Value(), &claimRecord)
		if err != nil {
			panic(err)
		}

		claimRecords = append(claimRecords, claimRecord)
	}
	return claimRecords
}

// GetClaimRecord returns the claim record for a specific address
func (k Keeper) GetClaimRecord(ctx sdk.Context, addr sdk.AccAddress) (types.ClaimRecord, error) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, []byte(types.ClaimRecordsStorePrefix))
	if !prefixStore.Has(addr) {
		return types.ClaimRecord{}, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "address does not have claim record")
	}
	bz := prefixStore.Get(addr)

	claimRecord := types.ClaimRecord{}
	err := proto.Unmarshal(bz, &claimRecord)
	if err != nil {
		return types.ClaimRecord{}, err
	}
	if len(claimRecord.Status) == 0 {
		return types.ClaimRecord{}, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "address does not have claim record")
	}
	return claimRecord, nil
}

// SetClaimRecord sets a claim record for an address in store
func (k Keeper) SetClaimRecord(ctx sdk.Context, claimRecord types.ClaimRecord) error {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, []byte(types.ClaimRecordsStorePrefix))

	bz, err := proto.Marshal(&claimRecord)
	if err != nil {
		return err
	}

	addr, err := sdk.AccAddressFromBech32(claimRecord.Address)
	if err != nil {
		return err
	}

	prefixStore.Set(addr, bz)
	return nil
}

// GetTotalClaimableAmountForAction returns total claimable amount for a specific action for the current time
// this includes claimed tokens
func (k Keeper) GetTotalClaimableAmountForAction(
	ctx sdk.Context, addr sdk.AccAddress, action types.Action) (sdk.Coins, error) {

	claimRecord, err := k.GetClaimRecord(ctx, addr)
	if err != nil {
		return nil, err
	}

	if claimRecord.Address == "" {
		return sdk.Coins{}, nil
	}

	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	// If we are before the start time, do nothing.
	// This case _shouldn't_ occur on chain, since the
	// start time ought to be chain start time.
	if ctx.BlockTime().Before(params.AirdropStartTime) {
		return sdk.Coins{}, nil
	}

	InitialClaimablePerAction := sdk.Coins{}
	for _, coin := range claimRecord.InitialClaimableAmount {
		InitialClaimablePerAction = InitialClaimablePerAction.Add(
			sdk.NewCoin(coin.Denom,
				coin.Amount.QuoRaw(int64(len(types.Action_name))),
			),
		)
	}
	//fmt.Printf("InitialClaimablePerAction %v \n", InitialClaimablePerAction)
	timeElapsed := ctx.BlockTime().Sub(params.AirdropStartTime)
	timeLeft := timeElapsed - params.DurationUntilDecay + params.DurationOfDecay
	//timeLeftUntilDecay := timeElapsed - params.DurationUntilDecay

	//vestingPeriod := params.DurationVestingPeriods[action]
	// The entire airdrop has completed
	if timeLeft <= 0 {
		return sdk.Coins{}, nil
	}
	// Early enough in the airdrop there is nov decay
	if timeElapsed <= params.DurationUntilDecay {
		return InitialClaimablePerAction, nil
	}

	// Positive, since goneTime > params.DurationUntilDecay
	decayTime := timeElapsed - params.DurationUntilDecay
	decayPercent := sdk.NewDec(decayTime.Nanoseconds()).QuoInt64(params.DurationOfDecay.Nanoseconds())
	claimablePercent := sdk.OneDec().Sub(decayPercent)
	claimableCoins := sdk.Coins{}
	for _, coin := range InitialClaimablePerAction {
		claimableCoins = claimableCoins.Add(sdk.NewCoin(coin.Denom, sdk.NewDecFromInt(coin.Amount).Mul(claimablePercent).RoundInt()))
	}

	return claimableCoins, nil
}

// GetTotalClaimableForAddr returns claimable amount for a specific action done by an address
func (k Keeper) GetTotalClaimableForAddr(ctx sdk.Context, addr sdk.AccAddress) (sdk.Coins, error) {
	claimRecord, err := k.GetClaimRecord(ctx, addr)
	if err != nil {
		return sdk.Coins{}, err
	}
	if claimRecord.Address == "" {
		return sdk.Coins{}, nil
	}

	totalClaimable := sdk.Coins{}

	for action := range types.Action_name {
		claimableForAction, err := k.GetTotalClaimableAmountForAction(ctx, addr, types.Action(action))
		if err != nil {
			return sdk.Coins{}, err
		}
		totalClaimable = totalClaimable.Add(claimableForAction...)
	}
	return totalClaimable, nil
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
	claimableCoin, err := k.GetTotalClaimableAmountForAction(ctx, addr, action)
	if err != nil {
		return err
	}

	if claimableCoin.Empty() {
		return nil
	}

	//we distribute 20% after action completion and the remaining become claimable after each vesting period (4)
	claimsPortion := sdk.NewCoins(sdk.NewCoin(p.ClaimDenom, claimableCoin.AmountOf(p.ClaimDenom).QuoRaw(5)))

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, claimsPortion)
	if err != nil {
		return err
	}

	//creates entries into the endblocker ( 4 )
	err = k.InsertEntriesIntoVestingQueue(ctx, addr.String(), byte(action), ctx.BlockHeader().Time)
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return err
	}

	//set claim record
	//fmt.Printf("claimRecord %v \n", claimRecord)
	claimRecord.Status[action].ActionCompleted = true
	err = k.SetClaimRecord(ctx, claimRecord)
	if err != nil {
		fmt.Printf("err: %v \n", err)
		return err
	}
	//fmt.Printf("test %v \n", claimableCoin)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaim,
			sdk.NewAttribute(sdk.AttributeKeySender, addr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, claimableCoin.String()),
		),
	})

	return nil
}

// ClaimClaimableForAddr remove claimable amount entries and transfer it to the sender account
func (k Keeper) ClaimClaimableForAddr(ctx sdk.Context, addr sdk.AccAddress) error {
	claimRecord, err := k.GetClaimRecord(ctx, addr)
	if err != nil {

		return err
	}
	p, err := k.GetParams(ctx)
	if err != nil {
		return err
	}

	claimableCoin := sdk.NewCoin(p.ClaimDenom, sdk.ZeroInt())
	claimedCoin := sdk.NewCoin(p.ClaimDenom, sdk.ZeroInt())
	var toClaimPeriods int64 = 0
	var claimedPeriods int64 = 0
	for action, status := range claimRecord.Status {
		if !status.ActionCompleted {
			continue
		}
		totalClaimableCoinsForAction, err := k.GetTotalClaimableAmountForAction(ctx, addr, types.Action(action))
		if err != nil {
			return err
		}

		var toClaimPeriodsForAction int64 = 0
		var claimedPeriodsForAction int64 = 1
		for period, completed := range status.VestingPeriodCompleted {
			//	fmt.Printf("period %v ; completed: %v\n", period, completed)
			if completed && !status.VestingPeriodClaimed[period] {
				toClaimPeriodsForAction = toClaimPeriodsForAction + 1
				claimRecord.Status[action].VestingPeriodClaimed[period] = true
			} else if completed && status.VestingPeriodClaimed[period] {
				claimedPeriodsForAction = claimedPeriodsForAction + 1
			}

		}

		if toClaimPeriodsForAction != 0 {

			toClaimPercent := sdk.NewDec(toClaimPeriodsForAction).Quo(sdk.NewDec(5))
			claimableTotalDec := sdk.NewDecFromInt(totalClaimableCoinsForAction.AmountOf(p.ClaimDenom))
			claimableDec := claimableTotalDec.Mul(toClaimPercent)
			claimableCoin = claimableCoin.AddAmount(claimableDec.TruncateInt())
		}
		claimedPercent := sdk.NewDec(claimedPeriodsForAction).Quo(sdk.NewDec(5))
		claimedCoin = claimedCoin.AddAmount(sdk.NewDecFromInt(totalClaimableCoinsForAction.AmountOf(p.ClaimDenom)).Mul(claimedPercent).TruncateInt())

		claimedPeriods = claimedPeriods + claimedPeriodsForAction
		toClaimPeriods = toClaimPeriods + toClaimPeriodsForAction
	}
	if toClaimPeriods == 0 || claimableCoin.Amount == sdk.ZeroInt() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "address does not have claimable tokens right now")
	}

	//get delegations and calculate min bonded ratio for claim
	delegationInfo := k.stakingKeeper.GetAllDelegatorDelegations(ctx, addr)
	//fmt.Printf("delegationInfo %v\n", delegationInfo)
	totalDelegations := sdk.ZeroDec()
	for _, delegation := range delegationInfo {
		totalDelegations = totalDelegations.Add(delegation.Shares)
	}
	//fmt.Printf("totalDelegations %v\n", totalDelegations)
	minBonded := sdk.NewDecWithPrec(67, 2).MulInt(claimedCoin.Amount)
	//fmt.Printf("minBonded amount: %v\n", minBonded)
	if totalDelegations.Sub(minBonded).IsNegative() {
		return errorsmod.Wrap(sdkerrors.ErrInsufficientFunds, "address does not have enough tokens staked")
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.NewCoins(claimableCoin))
	if err != nil {
		return err
	}
	fmt.Printf("record %v \n", claimRecord)
	err = k.SetClaimRecord(ctx, claimRecord)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaim,
			sdk.NewAttribute(sdk.AttributeKeySender, addr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, claimableCoin.String()),
		),
	})

	return nil
}

// FundRemainingsToCommunity fund remainings to the community when airdrop period ends
func (k Keeper) fundRemainingsToCommunity(ctx sdk.Context) error {
	moduleAccAddr := k.GetModuleAccountAddress(ctx)
	amt := k.GetModuleAccountBalance(ctx)
	return k.distrKeeper.FundCommunityPool(ctx, sdk.NewCoins(amt), moduleAccAddr)
}
