package keeper

import (

	//"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/trstlabs/trst/x/compute/internal/types"
)

// ContractPayoutCreator pays the creator of the contract
func (k Keeper) ContractPayoutCreator(ctx sdk.Context, contractAddress sdk.AccAddress) error {
	balance := k.bankKeeper.GetAllBalances(ctx, contractAddress)

	if !balance.Empty() {
		store := ctx.KVStore(k.storeKey)
		contractBz := store.Get(types.GetContractAddressKey(contractAddress))
		if contractBz == nil {
			return sdkerrors.Wrap(types.ErrNotFound, "contract")
		}
		var contract types.ContractInfo
		k.cdc.MustUnmarshal(contractBz, &contract)

		//returning the trst tokens
		commission := k.GetParams(ctx).Commission
		//percentageCreator := sdk.NewDecWithPrec(100-commission, 2)
		percentageCommission := sdk.NewDecWithPrec(commission, 2)

		toCommission := percentageCommission.MulInt(balance.AmountOf("utrst")).Ceil().TruncateInt()
		toCommissionCoins := sdk.NewCoins(sdk.NewCoin("utrst", toCommission))
		//toCreator := percentageCreator.MulInt(balance.AmountOf("utrst")).TruncateInt()
		balance = balance.Sub(toCommissionCoins)

		err := k.distrKeeper.FundCommunityPool(ctx, toCommissionCoins, contractAddress)
		if err != nil {
			return err
		}

		err = k.bankKeeper.SendCoins(ctx, contractAddress, contract.Creator, balance)
		if err != nil {
			return err
		}
	} else {
		k.Logger(ctx).Info("compute", "contract", "has no balance")
	}
	return nil
}

/*
// CallAutoMsg executes a final message before end-blocker deletion
func (k Keeper) CallAutoMsg(ctx sdk.Context, contractAddress sdk.AccAddress) (err error) {

	//get codeid first
	info, err := k.GetContractInfo(ctx, contractAddress)
	if err != nil {
		return err
	}

	if info.AutoMsg != nil {
		res, err := k.Execute(ctx, contractAddress, contractAddress, info.AutoMsg, sdk.NewCoins(sdk.NewCoin("utrst", sdk.ZeroInt())), nil)
		if err != nil {
			return err
		}
		k.SetContractResult(ctx, contractAddress, res)
	}
	return nil
}
*/

// SetIncentiveCoins distributes coins to the contracts in the compute module
func (k Keeper) SetIncentiveCoins(ctx sdk.Context, addressList []string) {
	params := k.GetParams(ctx)
	//if len(addressList) > 0 {
	total := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress("compute"), "utrst")
	k.Logger(ctx).Info("sent", "total", total)

	amount := total.Amount.QuoRaw(int64(len(addressList)))
	if amount.Int64() > params.MaxContractIncentive {
		amount = sdk.NewInt(params.MaxContractIncentive)
	}
	k.Logger(ctx).Info("sent", "amount", amount)
	//coins := sdk.NewCoins(sdk.NewCoin("utrst", amount))

	for _, addr := range addressList {
		sdkAddr, _ := sdk.AccAddressFromBech32(addr)
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdkAddr, sdk.NewCoins(sdk.NewCoin("utrst", amount)))
		if err != nil {
			k.Logger(ctx).Info("sent", "err", err)
		}

		k.Logger(ctx).Info("allocated", "contract", addr, "coins", amount)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeDistributedToContract,
				sdk.NewAttribute(types.AttributeKeyAddress, addr),
			),
		)

	}

}
