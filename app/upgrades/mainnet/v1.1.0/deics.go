package v1100

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ccvconsumerkeeper "github.com/cosmos/interchain-security/v6/x/ccv/consumer/keeper"
	ccvconsumertypes "github.com/cosmos/interchain-security/v6/x/ccv/consumer/types"
)

const (
	FundAddress  = "into1mhd977xqvd8pl7efsrtyltucw0dhf7h4mpv2ve"
	MinStake     = 2_000_000
	ICSSelfStake = 1
	Denom        = "uinto"
)

//go:embed validators/staking
var Vals embed.FS

type StakingValidator struct {
	OperatorAddress string `json:"operator_address"`
}

func wasInCometSet(
	ctx sdk.Context,
	sk stakingkeeper.Keeper,
	consAddr sdk.ConsAddress,
) (bool, error) {

	cometKnownVals := make(map[string]bool)
	sk.IterateLastValidatorPowers(ctx, func(addr sdk.ValAddress, power int64) bool {
		cometKnownVals[addr.String()] = true
		fmt.Println("LastValPower:", addr.String(), power)
		return false
	})
	return cometKnownVals[consAddr.String()], nil
}

func GetReadyValidators() (map[string]bool, error) {
	ready := make(map[string]bool)
	err := fs.WalkDir(Vals, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := Vals.ReadFile(path)
		if err != nil {
			return err
		}
		var skval StakingValidator
		if err := json.Unmarshal(data, &skval); err != nil {
			return err
		}
		ready[skval.OperatorAddress] = true
		return nil
	})
	return ready, err
}

func DeICS(
	ctx sdk.Context,
	sk stakingkeeper.Keeper,
	bk bankkeeper.Keeper,
	ck ccvconsumerkeeper.Keeper,
	readyValopers map[string]bool,
) error {
	_, DAOaddrBz, err := bech32.DecodeAndConvert(FundAddress)
	if err != nil {
		return err
	}
	DAOaddr := sdk.AccAddress(DAOaddrBz)

	consumerValidators := ck.GetAllCCValidator(ctx)
	validators, err := sk.GetAllValidators(ctx)
	if err != nil {
		return err
	}

	for _, val := range validators {
		valoper := val.GetOperator()
		valAddr, _ := sdk.ValAddressFromBech32(valoper)

		if readyValopers[valoper] {
			// Ensure ready validators meet the MinStake
			accAddr := sdk.AccAddress(valAddr)
			balance := bk.GetBalance(ctx, accAddr, Denom)
			if balance.Amount.LT(math.NewInt(MinStake)) {
				diff := math.NewInt(MinStake).Sub(balance.Amount)
				if err := bk.SendCoins(ctx, DAOaddr, accAddr, sdk.NewCoins(sdk.NewCoin(Denom, diff))); err != nil {
					return err
				}
			}
			continue
		}

		// SOFT REMOVAL: Instead of deleting, we set power to 0.
		// This prevents the "validator does not exist" panic in EndBlocker.
		val.Status = stakingtypes.Unbonded
		val.Tokens = math.ZeroInt()
		if err := sk.SetValidator(ctx, val); err != nil {
			return err
		}

		// Signal to CometBFT that this validator now has 0 weight.
		if err := sk.SetLastValidatorPower(ctx, valAddr, 0); err != nil {
			return err
		}
		fmt.Printf("Deactivated validator %s (set power 0)\n", valoper)
	}

	return moveICSToStaking(ctx, sk, bk, DAOaddr, consumerValidators)
}

func moveICSToStaking(
	ctx sdk.Context,
	sk stakingkeeper.Keeper,
	bk bankkeeper.Keeper,
	DAOaddr sdk.AccAddress,
	consumerValidators []ccvconsumertypes.CrossChainValidator,
) error {
	srv := stakingkeeper.NewMsgServerImpl(&sk)

	for i, v := range consumerValidators {
		consPubKey, _ := v.ConsPubKey()
		consAddr := sdk.ConsAddress(consPubKey.Address())

		// If they already exist in staking, skip creation
		if _, err := sk.GetValidatorByConsAddr(ctx, consAddr); err == nil {
			continue
		}

		// Use the consensus address bytes to form the validator address
		valAddr := sdk.ValAddress(consAddr)
		valoperAddr, _ := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), valAddr)

		// 1. Fund the account for self-delegation
		accAddr := sdk.AccAddress(consAddr)
		if err := bk.SendCoins(ctx, DAOaddr, accAddr, sdk.NewCoins(sdk.NewCoin(Denom, math.NewInt(ICSSelfStake)))); err != nil {
			return err
		}

		// 2. Create the validator object
		_, err := srv.CreateValidator(ctx, &stakingtypes.MsgCreateValidator{
			Description: stakingtypes.Description{Moniker: fmt.Sprintf("ics-%d", i)},
			Commission: stakingtypes.CommissionRates{
				Rate:          math.LegacyMustNewDecFromStr("0.1"),
				MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
				MaxChangeRate: math.LegacyMustNewDecFromStr("0.1"),
			}, MinSelfDelegation: math.NewInt(ICSSelfStake),
			ValidatorAddress: valoperAddr,
			Pubkey:           v.GetPubkey(),
			Value:            sdk.NewCoin(Denom, math.NewInt(ICSSelfStake)),
		})
		if err != nil {
			return err
		}

		// 3. Force them into the bonded set if they were active in ICS
		// This ensures the transition is seamless.
		newVal, _ := sk.GetValidator(ctx, valAddr)
		newVal.Status = stakingtypes.Bonded
		if err := sk.SetValidator(ctx, newVal); err != nil {
			return err
		}
		if err := sk.SetLastValidatorPower(ctx, valAddr, ICSSelfStake); err != nil {
			return err
		}
	}

	return nil
}
