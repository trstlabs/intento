package keeper

import (

	//"log"

	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/trstlabs/trst/x/compute/internal/types"
)

func TestContractIncentive(t *testing.T) {

	ctx, keeper, _, _, walletA, _, _, _ := setupTest(t, TestContractPaths[ibcContract], sdk.NewCoins())

	keeper.SetParams(ctx, types.Params{
		AutoMsgFundsCommission:          2,
		AutoMsgConstantFee:              1_000_000,                 // 1trst
		AutoMsgFlexFeeMul:               100,                       // 100/100 = 1 = gasUsed
		RecurringAutoMsgConstantFee:     1_000_000,                 // 1trst
		MaxContractDuration:             time.Hour * 24 * 366 * 10, // a little over 10 years
		MinContractDuration:             time.Second * 40,
		MinContractInterval:             time.Second * 20,
		MinContractDurationForIncentive: time.Second * 20, // time.Hour * 24 // 1 day
		MaxContractIncentive:            5_000_000,        // 5trst
		ContractIncentiveMul:            100,              //  100/100 = 1 = full incentive
		MinContractBalanceForIncentive:  50_000_000,       // 50trst
	})
	gasCoin := sdk.NewInt64Coin(sdk.DefaultBondDenom, 731_241)
	/*if err := keeper.bankKeeper.MintCoins(ctx, "compute", gasCoins); err != nil {
		panic(err)
	}*/
	//keeper.ContractIncentive(ctx, contract, 798_012, true)
	_, err := keeper.ContractIncentive(ctx, gasCoin, walletA)
	require.Empty(t, err)
	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(731_241)), keeper.bankKeeper.GetBalance(ctx, walletA, sdk.DefaultBondDenom))

	//require.Equal(t, 10, len(contracts))

}

func TestContractIncentiveWithLowDenom(t *testing.T) {

	ctx, keeper, _, _, walletA, _, _, _ := setupTest(t, TestContractPaths[ibcContract], sdk.NewCoins())

	keeper.SetParams(ctx, types.Params{
		AutoMsgFundsCommission:          2,
		AutoMsgConstantFee:              1_000_000,                 // 1trst
		AutoMsgFlexFeeMul:               100,                       // 100/100 = 1 = gasUsed
		RecurringAutoMsgConstantFee:     1_000_000,                 // 1trst
		MaxContractDuration:             time.Hour * 24 * 366 * 10, // a little over 10 years
		MinContractDuration:             time.Second * 40,
		MinContractInterval:             time.Second * 20,
		MinContractDurationForIncentive: time.Second * 20, // time.Hour * 24 // 1 day
		MaxContractIncentive:            5_000_000,        // 5trst
		ContractIncentiveMul:            23,               //  100/100 = 1 = full incentive
		MinContractBalanceForIncentive:  50_000_000,       // 50trst
	})
	gasCoin := sdk.NewInt64Coin(sdk.DefaultBondDenom, 731_241)

	_, err := keeper.ContractIncentive(ctx, gasCoin, walletA)
	require.Empty(t, err)
	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(168_185)), keeper.bankKeeper.GetBalance(ctx, walletA, sdk.DefaultBondDenom))

}

func TestDistributeCoinsWithoutIncentive(t *testing.T) {

	ctx, keeper, _, _, _, _, _, _ := setupTest(t, TestContractPaths[ibcContract], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1_000_000))))
	addr1, _ := CreateFakeFundedAccount(ctx, keeper.accountKeeper, keeper.bankKeeper, sdk.NewCoins(sdk.NewInt64Coin("stake", 3_000_000)))

	keeper.SetParams(ctx, types.Params{
		AutoMsgFundsCommission:          10,
		AutoMsgConstantFee:              1_000_000,                 // 1trst
		AutoMsgFlexFeeMul:               100,                       // 100/100 = 1 = gasUsed
		RecurringAutoMsgConstantFee:     1_000_000,                 // 1trst
		MaxContractDuration:             time.Hour * 24 * 366 * 10, // a little over 10 years
		MinContractDuration:             time.Second * 40,
		MinContractInterval:             time.Second * 20,
		MinContractDurationForIncentive: time.Second * 20, // time.Hour * 24 // 1 day
		MaxContractIncentive:            5_000_000,        // 5trst
		ContractIncentiveMul:            0,                //  0/100 = 0 = no incentive
		MinContractBalanceForIncentive:  50_000_000,       // 50trst
	})
	pub2 := secp256k1.GenPrivKey().PubKey()
	addr2 := sdk.AccAddress(pub2.Address())
	types.Denom = "stake"

	contrInfo := types.ContractInfo{
		CodeID: 0, Creator: addr2, Owner: addr2, ContractId: "test", Created: types.NewAbsoluteTxPosition(ctx), AutoMsg: []byte("test"), Duration: time.Minute, Interval: time.Second * 20, StartTime: time.Now().Add(time.Hour * -1), EndTime: time.Now().Add(time.Second * 20), IBCPortID: "", CallbackSig: nil,
	}
	contract := types.ContractInfoWithAddress{
		Address:      addr1,
		ContractInfo: &contrInfo,
	}

	err := keeper.DistributeCoins(ctx, contract, 800_000, true)
	require.Empty(t, err)

	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0)), keeper.bankKeeper.GetBalance(ctx, contract.Address, sdk.DefaultBondDenom))
	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt((3_000_000*0.9-1_800_000))), keeper.bankKeeper.GetBalance(ctx, addr2, sdk.DefaultBondDenom))
}

func TestDistributeCoinsEmptyContractBalanceWithMaxIncentive(t *testing.T) {

	ctx, keeper, _, _, _, _, _, _ := setupTest(t, TestContractPaths[ibcContract], sdk.NewCoins())

	keeper.SetParams(ctx, types.Params{
		AutoMsgFundsCommission:          2,
		AutoMsgConstantFee:              1_000_000,                 // 1trst
		AutoMsgFlexFeeMul:               100,                       // 100/100 = 1 = gasUsed
		RecurringAutoMsgConstantFee:     1_000_000,                 // 1trst
		MaxContractDuration:             time.Hour * 24 * 366 * 10, // a little over 10 years
		MinContractDuration:             time.Second * 40,
		MinContractInterval:             time.Second * 20,
		MinContractDurationForIncentive: time.Second * 20, // time.Hour * 24 // 1 day
		MaxContractIncentive:            5_000_000,        // 5trst
		ContractIncentiveMul:            100,              //  100/100 = 1 = full incentive
		MinContractBalanceForIncentive:  50_000_000,       // 50trst
	})

	pub1 := secp256k1.GenPrivKey().PubKey()
	pub2 := secp256k1.GenPrivKey().PubKey()
	addr1 := sdk.AccAddress(pub1.Address())
	addr2 := sdk.AccAddress(pub2.Address())
	types.Denom = "stake"

	contrInfo := types.ContractInfo{
		CodeID: 0, Creator: addr2, Owner: addr2, ContractId: "test", Created: types.NewAbsoluteTxPosition(ctx), AutoMsg: []byte("test"), Duration: time.Minute, Interval: time.Second * 20, StartTime: time.Now().Add(time.Hour * -1), EndTime: time.Now().Add(time.Second * 20), IBCPortID: "", CallbackSig: nil,
	}
	contract := types.ContractInfoWithAddress{
		Address:      addr1,
		ContractInfo: &contrInfo,
	}

	err := keeper.DistributeCoins(ctx, contract, 731_241, true)
	require.Empty(t, err)
	/*info := keeper.GetContractInfo(ctx, contract.Address)
	fmt.Print("info", info)
	require.Equal(t, types.ContractInfo{
		CodeID: 0, Creator: addr2, Owner: addr2, ContractId: "test", Created: types.NewAbsoluteTxPosition(ctx), AutoMsg: []byte("test"), Duration: time.Minute, Interval: time.Second * 20, StartTime: time.Now().Add(time.Hour * -1), EndTime: time.Now().Add(time.Second * 20), IBCPortID: "", CallbackSig: nil,
	}, *info)*/
	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0)), keeper.bankKeeper.GetBalance(ctx, contract.Address, sdk.DefaultBondDenom))
	require.Equal(t, sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0)), keeper.bankKeeper.GetBalance(ctx, addr2, sdk.DefaultBondDenom))

}

func TestSelfExecute(t *testing.T) {
	//todo
}
