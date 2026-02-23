package v1100

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec/address"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
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
	stakingKeeper stakingkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	consumerKeeper ccvconsumerkeeper.Keeper,
	readyValopers map[string]bool,
) error {
	stakingKeeper.IterateLastValidatorPowers(ctx, func(addr sdk.ValAddress, power int64) bool {
		fmt.Println("LastValPower:", addr.String(), power)
		return false
	})
	_, DAOaddrBz, err := bech32.DecodeAndConvert(FundAddress)
	if err != nil {
		return err
	}
	DAOaddr := sdk.AccAddress(DAOaddrBz)

	consumerValidators := consumerKeeper.GetAllCCValidator(ctx)
	fmt.Printf("Consumer validators: %d\n", len(consumerValidators))

	if err != nil {
		return fmt.Errorf("failed to update staking params: %w", err)
	}

	validators, err := stakingKeeper.GetAllValidators(ctx)
	if err != nil {
		return err
	}
	fmt.Println("Validators staking:", len(validators))
	for _, val := range validators {
		valoper := val.GetOperator()
		valAddr, err := sdk.ValAddressFromBech32(valoper)
		if err != nil {
			return err
		}

		consAddr, err := val.GetConsAddr()
		if err != nil {
			return err
		}
		fmt.Printf("Processing  consAddr %s  val Addr %s\n", consAddr, valAddr)
		if readyValopers[valoper] {
			accAddr := sdk.AccAddress(valAddr)
			balance := bankKeeper.GetBalance(ctx, accAddr, Denom)
			if balance.Amount.LT(math.NewInt(MinStake)) {
				coins := sdk.NewCoins(sdk.NewCoin(Denom, math.NewInt(MinStake)))
				if err := bankKeeper.SendCoins(ctx, DAOaddr, accAddr, coins); err != nil {
					return err
				}
			}
			continue
		}

		inComet, err := wasInCometSet(ctx, stakingKeeper, consAddr)
		if err != nil {
			return err
		}

		if err := stakingKeeper.DeleteValidatorByPowerIndex(ctx, val); err != nil {
			return err
		}
		fmt.Printf("Removed validator %s (in comet: %t)\n", valoper, inComet)
		if !inComet {
			// governance-only — prevent removal update
			consValAddr := sdk.ValAddress(consAddr)
			if err := stakingKeeper.SetLastValidatorPower(ctx, consValAddr, 0); err != nil {
				return err
			}
		}
	}

	return moveICSToStaking(ctx, stakingKeeper, bankKeeper, DAOaddr, consumerValidators)
}

func moveICSToStaking(
	ctx sdk.Context,
	sk stakingkeeper.Keeper,
	bk bankkeeper.Keeper,
	DAOaddr sdk.AccAddress,
	consumerValidators []ccvconsumertypes.CrossChainValidator,
) error {

	srv := stakingkeeper.NewMsgServerImpl(&sk)
	bondedCount := 0

	for i, v := range consumerValidators {
		consPubKey, err := v.ConsPubKey()
		if err != nil {
			return err
		}
		fmt.Printf("Processing consumer validator %d\n", i)

		pk, ok := consPubKey.(cryptotypes.PubKey)
		if !ok {
			return fmt.Errorf("unexpected pubkey type %T", consPubKey)
		}

		consAddr := sdk.ConsAddress(pk.Address())

		if _, err := sk.GetValidatorByConsAddr(ctx, consAddr); err == nil {
			continue
		} else if !errors.Is(err, stakingtypes.ErrNoValidatorFound) {
			return err
		}

		inComet, err := wasInCometSet(ctx, sk, consAddr)
		if err != nil {
			return err
		}

		valAddr := sdk.ValAddress(consAddr)

		valoperAddr, err := bech32.ConvertAndEncode(
			sdk.GetConfig().GetBech32ValidatorAddrPrefix(), valAddr,
		)
		fmt.Printf("Processing  valoperAddr %s  val Addr %s\n", valoperAddr, valAddr)

		if err != nil {
			return err
		}

		accAddr := sdk.AccAddress(consAddr)
		if err := bk.SendCoins(ctx, DAOaddr, accAddr, sdk.NewCoins(
			sdk.NewCoin(Denom, math.NewInt(ICSSelfStake)),
		)); err != nil {
			return err
		}

		_, err = srv.CreateValidator(ctx, &stakingtypes.MsgCreateValidator{
			Description: stakingtypes.Description{Moniker: fmt.Sprintf("ics %d", i)},
			Commission: stakingtypes.CommissionRates{
				Rate:          math.LegacyMustNewDecFromStr("0.1"),
				MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
				MaxChangeRate: math.LegacyMustNewDecFromStr("0.1"),
			},
			MinSelfDelegation: math.NewInt(ICSSelfStake),
			ValidatorAddress:  valoperAddr,
			Pubkey:            v.GetPubkey(),
			Value:             sdk.NewCoin(Denom, math.NewInt(ICSSelfStake)),
		})
		fmt.Printf("Created validator %s (in comet: %t)\n", valoperAddr, inComet)
		if err != nil {
			return err
		}

		if !inComet {
			continue
		}

		if err := sk.SetLastValidatorPower(ctx, valAddr, 1); err != nil {
			return err
		}

		savedVal, err := sk.GetValidator(ctx, valAddr)
		if err != nil {
			return err
		}

		if _, err := bondValidator(ctx, sk, savedVal); err != nil {
			return err
		}
		fmt.Printf("Bonded validator %s (in comet: %t)\n", valoperAddr, inComet)

		bondedCount++
	}

	if bondedCount == 0 {
		return nil
	}

	coins := sdk.NewCoins(
		sdk.NewCoin(Denom, math.NewInt(int64(bondedCount)*ICSSelfStake)),
	)

	return bk.SendCoinsFromModuleToModule(
		ctx,
		stakingtypes.NotBondedPoolName,
		stakingtypes.BondedPoolName,
		coins,
	)
}

func bondValidator(ctx context.Context, k stakingkeeper.Keeper, validator stakingtypes.Validator) (stakingtypes.Validator, error) {
	if err := k.DeleteValidatorByPowerIndex(ctx, validator); err != nil {
		return validator, err
	}

	validator = validator.UpdateStatus(stakingtypes.Bonded)

	if err := k.SetValidator(ctx, validator); err != nil {
		return validator, err
	}
	if err := k.SetValidatorByPowerIndex(ctx, validator); err != nil {
		return validator, err
	}
	if err := k.DeleteValidatorQueue(ctx, validator); err != nil {
		return validator, err
	}

	consAddr, err := validator.GetConsAddr()
	if err != nil {
		return validator, err
	}

	codec := address.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())
	valAddrBz, err := codec.StringToBytes(validator.GetOperator())
	if err != nil {
		return validator, err
	}

	if err := k.Hooks().AfterValidatorBonded(ctx, consAddr, valAddrBz); err != nil {
		return validator, err
	}

	return validator, nil
}
