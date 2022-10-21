package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/trstlabs/trst/x/alloc/types"
)

type (
	Keeper struct {
		cdc           codec.BinaryCodec
		storeKey      sdk.StoreKey
		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper
		stakingKeeper types.StakingKeeper
		distrKeeper   types.DistrKeeper
		paramstore    paramtypes.Subspace
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper, stakingKeeper types.StakingKeeper, distrKeeper types.DistrKeeper,
	ps paramtypes.Subspace,
) *Keeper {

	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		accountKeeper: accountKeeper, bankKeeper: bankKeeper, stakingKeeper: stakingKeeper, distrKeeper: distrKeeper, //computeKeeper: ck,
		paramstore: ps,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetBalance gets balance
func (k Keeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	return k.bankKeeper.GetBalance(ctx, addr, denom)
}

// DistributeInflation distributes module-specific inflation
func (k Keeper) DistributeInflation(ctx sdk.Context) error {
	blockInflationAddr := k.accountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName).GetAddress()
	blockInflation := k.bankKeeper.GetBalance(ctx, blockInflationAddr, k.stakingKeeper.BondDenom(ctx))
	//blockInflationDec := sdk.NewDecFromInt(blockInflation.Amount)

	params := k.GetParams(ctx)
	proportions := params.DistributionProportions

	contrIncentiveCoins := sdk.NewCoins(k.GetProportions(ctx, blockInflation, proportions.TrustlessContractIncentives))
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, authtypes.FeeCollectorName, "compute", contrIncentiveCoins)
	if err != nil {
		return err
	}
	/*itemIncentiveCoins := sdk.NewCoins(k.GetProportions(ctx, blockInflation, proportions.ItemIncentives))
	err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, authtypes.FeeCollectorName, "item_incentives", itemIncentiveCoins)
	if err != nil {
		return err
	}*/

	k.Logger(ctx).Debug("funded trustless contracts", "amount", contrIncentiveCoins.String(), "from", blockInflationAddr)

	//staking incentives stay in the fee collector account and are to be moved to on next begin blocker
	stakingIncentivesCoins := sdk.NewCoins(k.GetProportions(ctx, blockInflation, proportions.Staking))

	rewardCoin := k.GetProportions(ctx, blockInflation, proportions.ContributorRewards)
	rewardCoins := sdk.NewCoins(rewardCoin)

	for _, w := range params.WeightedContributorRewardsReceivers {
		rewardPortionCoins := sdk.NewCoins(k.GetProportions(ctx, rewardCoin, w.Weight))
		if w.Address == "" {
			err := k.distrKeeper.FundCommunityPool(ctx, rewardPortionCoins, blockInflationAddr)
			if err != nil {
				return err
			}
		} else {
			contributorRewardsAddr, err := sdk.AccAddressFromBech32(w.Address)
			if err != nil {
				return err
			}
			err = k.bankKeeper.SendCoins(ctx, blockInflationAddr, contributorRewardsAddr, rewardPortionCoins)
			if err != nil {
				return err
			}
			k.Logger(ctx).Debug("sent coins to contributor", "amount", rewardPortionCoins.String(), "from", blockInflationAddr)
		}
	}

	// subtract from original provision to ensure no coins left over after the allocations
	communityPoolCoins := sdk.NewCoins(blockInflation).Sub(stakingIncentivesCoins). /*.Sub(itemIncentiveCoins)*/ Sub(contrIncentiveCoins).Sub(rewardCoins)

	err = k.distrKeeper.FundCommunityPool(ctx, communityPoolCoins, blockInflationAddr)
	if err != nil {
		return err
	}

	return nil
}

// GetProportions gets the balance of the `MintedDenom` from minted coins
// and returns coins according to the `AllocationRatio`
func (k Keeper) GetProportions(ctx sdk.Context, mintedCoin sdk.Coin, ratio sdk.Dec) sdk.Coin {
	return sdk.NewCoin(mintedCoin.Denom, mintedCoin.Amount.ToDec().Mul(ratio).TruncateInt())
}
