package keeper

import (
	"fmt"
	"testing"
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"
)

func TestDistributeCoinsNotRecurring(t *testing.T) {

	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))

	autoTxAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))

	keeper.SetParams(ctx, types.Params{
		AutoTxFundsCommission:      10,
		AutoTxConstantFee:          1_000_000,                 // 1trst
		AutoTxFlexFeeMul:           100,                       // 100/100 = 1 = gasUsed
		RecurringAutoTxConstantFee: 1_000_000,                 // 1trst
		MaxAutoTxDuration:          time.Hour * 24 * 366 * 10, // a little over 10 years
		MinAutoTxDuration:          time.Second * 60,
		MinAutoTxInterval:          time.Second * 20,
	})
	// fmt.Printf("%v", ctx.WithBlockGasMeter(sdk.NewInfiniteGasMeter()))

	pub2 := secp256k1.GenPrivKey().PubKey()
	addr2 := sdk.AccAddress(pub2.Address())
	types.Denom = "stake"

	autoTxInfo := types.AutoTxInfo{
		TxID: 0, Owner: addr2.String(), FeeAddress: autoTxAddr.String(), Msgs: NewMsg(), Duration: time.Minute, Interval: time.Second * 20, StartTime: time.Now().Add(time.Hour * -1), EndTime: time.Now().Add(time.Second * 20), PortID: "ibccontoller-test", ConnectionID: "connection-0",
	}

	val := keeper.stakingKeeper.ValidatorByConsAddr(ctx, sdk.ConsAddress(ctx.BlockHeader().ProposerAddress))
	require.Equal(t, sdk.ZeroDec(), keeper.distrKeeper.GetValidatorCurrentRewards(ctx, val.GetOperator()).Rewards.AmountOf(sdk.DefaultBondDenom))

	_, err := keeper.DistributeCoins(ctx, autoTxInfo, sdk.NewInt(time.Minute.Milliseconds()), false, ctx.BlockHeader().ProposerAddress)
	require.Empty(t, err)

	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0)), keeper.bankKeeper.GetBalance(ctx, autoTxAddr, sdk.DefaultBondDenom))
	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt((3_000_000*0.9-1_000_000-time.Minute.Milliseconds()))), keeper.bankKeeper.GetBalance(ctx, addr2, sdk.DefaultBondDenom))
	// check validator current rewards
	require.Equal(t, sdk.NewDec(time.Minute.Milliseconds()), keeper.distrKeeper.GetValidatorCurrentRewards(ctx, val.GetOperator()).Rewards.AmountOf(sdk.DefaultBondDenom))
}

func TestDistributeCoinsFeesFromUser(t *testing.T) {

	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins())

	keeper.SetParams(ctx, types.Params{
		AutoTxFundsCommission:      2,
		AutoTxConstantFee:          1_000_000,                 // 1trst
		AutoTxFlexFeeMul:           100,                       // 100/100 = 1 = gasUsed
		RecurringAutoTxConstantFee: 1_000_000,                 // 1trst
		MaxAutoTxDuration:          time.Hour * 24 * 366 * 10, // a little over 10 years
		MinAutoTxDuration:          time.Second * 60,
		MinAutoTxInterval:          time.Second * 20,
	})

	pub1 := secp256k1.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pub1.Address())
	addr2, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	types.Denom = "stake"

	autoTxInfo := types.AutoTxInfo{
		TxID: 0, Owner: addr2.String(), FeeAddress: addr1.String(), Msgs: NewMsg(), Duration: time.Minute, Interval: time.Second * 20, StartTime: time.Now().Add(time.Hour * -1), EndTime: time.Now().Add(time.Second * 20), PortID: "ibccontoller-test", ConnectionID: "connection-0",
	}

	_, err := keeper.DistributeCoins(ctx, autoTxInfo, sdk.NewInt(time.Minute.Milliseconds()), true, ctx.BlockHeader().ProposerAddress)
	require.Empty(t, err)

	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0)), keeper.bankKeeper.GetBalance(ctx, addr1, sdk.DefaultBondDenom))
	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_940_000)), keeper.bankKeeper.GetBalance(ctx, addr2, sdk.DefaultBondDenom))

	// check validator current rewards
	val := keeper.stakingKeeper.ValidatorByConsAddr(ctx, sdk.ConsAddress(ctx.BlockHeader().ProposerAddress))
	require.Equal(t, sdk.NewDec(time.Minute.Milliseconds()), keeper.distrKeeper.GetValidatorCurrentRewards(ctx, val.GetOperator()).Rewards.AmountOf(sdk.DefaultBondDenom))

}

func TestDistributeCoinsEmptyAutoTxBalanceNotLastExec(t *testing.T) {

	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins())
	autoTxAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))

	keeper.SetParams(ctx, types.Params{
		AutoTxFundsCommission:      2,
		AutoTxConstantFee:          1_000_000,                 // 1trst
		AutoTxFlexFeeMul:           100,                       // 100/100 = 1 = gasUsed
		RecurringAutoTxConstantFee: 1_000_000,                 // 1trst
		MaxAutoTxDuration:          time.Hour * 24 * 366 * 10, // a little over 10 years
		MinAutoTxDuration:          time.Second * 60,
		MinAutoTxInterval:          time.Second * 60,
	})

	//pub1 := secp256k1.GenPrivKey().PubKey()
	pub2 := secp256k1.GenPrivKey().PubKey()
	//addr1 := sdk.AccAddress(pub1.Address())
	addr2 := sdk.AccAddress(pub2.Address())
	types.Denom = "stake"

	autoTxInfo := types.AutoTxInfo{
		TxID: 0, Owner: addr2.String(), FeeAddress: autoTxAddr.String(), Msgs: NewMsg(), Duration: time.Minute, Interval: time.Second * 20, StartTime: time.Now().Add(time.Hour * -1), EndTime: time.Now().Add(time.Second * 20), PortID: "ibccontoller-test", ConnectionID: "connection-0",
	}

	_, err := keeper.DistributeCoins(ctx, autoTxInfo, sdk.NewInt(time.Minute.Milliseconds()), true, ctx.BlockHeader().ProposerAddress)
	require.Empty(t, err)
	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_940_000)), keeper.bankKeeper.GetBalance(ctx, autoTxAddr, sdk.DefaultBondDenom))
	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0)), keeper.bankKeeper.GetBalance(ctx, addr2, sdk.DefaultBondDenom))

	// check validator current rewards
	val := keeper.stakingKeeper.ValidatorByConsAddr(ctx, sdk.ConsAddress(ctx.BlockHeader().ProposerAddress))
	require.Equal(t, sdk.NewDec(time.Minute.Milliseconds()), keeper.distrKeeper.GetValidatorCurrentRewards(ctx, val.GetOperator()).Rewards.AmountOf(sdk.DefaultBondDenom))

}

func TestDistributeCoinsEmptyAutoTxBalanceAndMultipliedFlexFee(t *testing.T) {

	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins())
	autoTxAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	keeper.SetParams(ctx, types.Params{
		AutoTxFundsCommission:      2,
		AutoTxConstantFee:          1_000_000,                 // 1trst
		AutoTxFlexFeeMul:           250,                       // 250/100 = 2.5x
		RecurringAutoTxConstantFee: 1_000_000,                 // 1trst
		MaxAutoTxDuration:          time.Hour * 24 * 366 * 10, // a little over 10 years
		MinAutoTxDuration:          time.Second * 60,
		MinAutoTxInterval:          time.Second * 20,
	})

	//pub1 := secp256k1.GenPrivKey().PubKey()
	pub2 := secp256k1.GenPrivKey().PubKey()
	//addr1 := sdk.AccAddress(pub1.Address())
	addr2 := sdk.AccAddress(pub2.Address())
	types.Denom = "stake"

	autoTxInfo := types.AutoTxInfo{
		TxID: 0, Owner: addr2.String(), FeeAddress: autoTxAddr.String(), Msgs: NewMsg(), Duration: time.Minute, Interval: time.Second * 20, StartTime: time.Now().Add(time.Hour * -1), EndTime: time.Now().Add(time.Second * 20), PortID: "ibccontoller-test", ConnectionID: "connection-0",
	}

	_, err := keeper.DistributeCoins(ctx, autoTxInfo, sdk.NewInt(time.Minute.Milliseconds()), false, ctx.BlockHeader().ProposerAddress)
	require.Empty(t, err)
	fmt.Printf("%v\n", keeper.bankKeeper.GetBalance(ctx, autoTxAddr, sdk.DefaultBondDenom))
	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0)), keeper.bankKeeper.GetBalance(ctx, autoTxAddr, sdk.DefaultBondDenom))
	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(3_000_000*0.98-1_000_000-time.Minute.Milliseconds()*250/100)), keeper.bankKeeper.GetBalance(ctx, addr2, sdk.DefaultBondDenom))

	// check validator current rewards
	val := keeper.stakingKeeper.ValidatorByConsAddr(ctx, sdk.ConsAddress(ctx.BlockHeader().ProposerAddress))
	require.Equal(t, sdk.NewDec(time.Minute.Milliseconds()*250/100), keeper.distrKeeper.GetValidatorCurrentRewards(ctx, val.GetOperator()).Rewards.AmountOf(sdk.DefaultBondDenom))

}

func TestDistributeCoinsEmptyAutoTxBalanceAndDiscountedFlexFee(t *testing.T) {

	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins())
	autoTxAddr, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))
	keeper.SetParams(ctx, types.Params{
		AutoTxFundsCommission:      2,
		AutoTxConstantFee:          1_000_000,                 // 1trst
		AutoTxFlexFeeMul:           25,                        // 250/100 = 0.25x
		RecurringAutoTxConstantFee: 1_000_000,                 // 1trst
		MaxAutoTxDuration:          time.Hour * 24 * 366 * 10, // a little over 10 years
		MinAutoTxDuration:          time.Second * 60,
		MinAutoTxInterval:          time.Second * 20,
	})

	//pub1 := secp256k1.GenPrivKey().PubKey()
	pub2 := secp256k1.GenPrivKey().PubKey()
	//addr1 := sdk.AccAddress(pub1.Address())
	addr2 := sdk.AccAddress(pub2.Address())
	types.Denom = "stake"

	autoTxInfo := types.AutoTxInfo{
		TxID: 0, Owner: addr2.String(), FeeAddress: autoTxAddr.String(), Msgs: NewMsg(), Duration: time.Minute, Interval: time.Second * 20, StartTime: time.Now().Add(time.Hour * -1), EndTime: time.Now().Add(time.Second * 20), PortID: "ibccontoller-test", ConnectionID: "connection-0",
	}

	_, err := keeper.DistributeCoins(ctx, autoTxInfo, sdk.NewInt(time.Minute.Milliseconds()), false, ctx.BlockHeader().ProposerAddress)
	require.Empty(t, err)
	fmt.Printf("%v\n", keeper.bankKeeper.GetBalance(ctx, autoTxAddr, sdk.DefaultBondDenom))
	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0)), keeper.bankKeeper.GetBalance(ctx, autoTxAddr, sdk.DefaultBondDenom))
	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(3_000_000*0.98-1_000_000-time.Minute.Milliseconds()*25/100)), keeper.bankKeeper.GetBalance(ctx, addr2, sdk.DefaultBondDenom))

	// check validator current rewards
	val := keeper.stakingKeeper.ValidatorByConsAddr(ctx, sdk.ConsAddress(ctx.BlockHeader().ProposerAddress))
	require.Equal(t, sdk.NewDec(time.Minute.Milliseconds()*25/100), keeper.distrKeeper.GetValidatorCurrentRewards(ctx, val.GetOperator()).Rewards.AmountOf(sdk.DefaultBondDenom))

}

func TestDistributeCoinsLargeFee(t *testing.T) {

	ctx, keeper, _, _, _, _ := setupTest(t, sdk.NewCoins())

	keeper.SetParams(ctx, types.Params{
		AutoTxFundsCommission:      2,
		AutoTxConstantFee:          1_000_000,                 // 1trst
		AutoTxFlexFeeMul:           100,                       // 100/100 = 1 = gasUsed
		RecurringAutoTxConstantFee: 1_000_000,                 // 1trst
		MaxAutoTxDuration:          time.Hour * 24 * 366 * 10, // a little over 10 years
		MinAutoTxDuration:          time.Second * 60,
		MinAutoTxInterval:          time.Second * 20,
	})

	pub1 := secp256k1.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pub1.Address())
	addr2, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000_000)))
	types.Denom = "stake"

	autoTxInfo := types.AutoTxInfo{
		TxID: 0, Owner: addr2.String(), FeeAddress: addr1.String(), Msgs: NewMsg(), Duration: time.Minute, Interval: time.Second * 20, StartTime: time.Now().Add(time.Hour * -1), EndTime: time.Now().Add(time.Second * 20), PortID: "ibccontoller-test", ConnectionID: "connection-0",
	}

	_, err := keeper.DistributeCoins(ctx, autoTxInfo, sdk.NewInt(time.Hour.Milliseconds()*24*30), true, ctx.BlockHeader().ProposerAddress)
	require.Empty(t, err)

	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0)), keeper.bankKeeper.GetBalance(ctx, addr1, sdk.DefaultBondDenom))
	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(407000000)), keeper.bankKeeper.GetBalance(ctx, addr2, sdk.DefaultBondDenom))

	// check validator current rewards
	val := keeper.stakingKeeper.ValidatorByConsAddr(ctx, sdk.ConsAddress(ctx.BlockHeader().ProposerAddress))
	require.Equal(t, sdk.NewDec(time.Hour.Milliseconds()*24*30), keeper.distrKeeper.GetValidatorCurrentRewards(ctx, val.GetOperator()).Rewards.AmountOf(sdk.DefaultBondDenom))

}

func NewMsg() []*sdktypes.Any {
	msg, _ := sdktypes.NewAnyWithValue(&types.MsgSubmitTx{})
	return []*sdktypes.Any{msg}
}
