package keeper

import (

	//"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/trstlabs/trst/x/compute/internal/types"
)

// ContractPayoutCreator pays the creator of the contract
func (k Keeper) ContractPayoutCreator(ctx sdk.Context, contractAddress sdk.AccAddress) error {

	store := ctx.KVStore(k.storeKey)
	contractBz := store.Get(types.GetContractAddressKey(contractAddress))
	if contractBz == nil {
		return sdkerrors.Wrap(types.ErrNotFound, "contract")
	}
	var contract types.ContractInfo
	k.cdc.MustUnmarshal(contractBz, &contract)
	//fmt.Print("balances..")
	//payout contract coins to the creator
	balance := k.bankKeeper.GetAllBalances(ctx, contractAddress)
	if !balance.Empty() {
		//returning the trst tokens
		commission := k.GetParams(ctx).Commission
		percentageCreator := sdk.NewDecWithPrec(100-commission, 2)
		percentageCommission := sdk.NewDecWithPrec(commission, 2)

		toCommission := percentageCommission.MulInt(balance.AmountOf("utrst")).Ceil().TruncateInt()
		toCreator := percentageCreator.MulInt(balance.AmountOf("utrst")).TruncateInt()
		balance = balance.Sub(sdk.NewCoins(sdk.NewCoin("utrst", toCreator)))

		err := k.distrKeeper.FundCommunityPool(ctx, sdk.NewCoins(sdk.NewCoin("utrst", toCommission)), contractAddress)
		if err != nil {
			return err
		}

		err = k.bankKeeper.SendCoins(ctx, contractAddress, contract.Creator, balance)
		if err != nil {
			return err
		}
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
