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

func GetReadyValidators() (map[string]bool, error) {
	ready := make(map[string]bool)
	err := fs.WalkDir(Vals, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
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
		err = json.Unmarshal(data, &skval)
		if err != nil {
			return err
		}
		ready[skval.OperatorAddress] = true
		return nil
	})
	return ready, err
}

// DeICS migrates the chain from ICS consumer to sovereign.
//
// The original approach called stakingKeeper.Jail() on non-ready governance
// validators. This crashed with "validator does not exist" in FinalizeBlock
// because:
//
//  1. On an ICS consumer chain, CometBFT's active validator set is the
//     provider-assigned ICS set (stored in x/ccv/consumer), NOT the local
//     governance/democracy validators in x/staking.
//
//  2. Calling Jail() emits a power=0 validator update to CometBFT for that
//     consensus key. If CometBFT has never seen the key (governance validators
//     never participated in local consensus), it panics.
//
//  3. The ICS validators that CometBFT DOES know about were never in
//     x/staking at all, so nothing was cleaning them out of CometBFT's set.
//
// Fix (borrowed from the advanced v1100 implementation):
//
//   - Sovereign/ready validators: created normally via MsgCreateValidator so
//     staking emits a proper power>0 update to CometBFT.
//
//   - ICS validators: registered into staking with 1uinto stake (VP rounds to
//     0) and force-bonded so staking emits a power=0 update. CometBFT sees a
//     clean removal.
//
//   - Non-ready governance validators: set directly to Unbonded in the staking
//     store WITHOUT calling Jail, so no consensus update is emitted.
func DeICS(
	ctx sdk.Context,
	stakingKeeper stakingkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	consumerKeeper ccvconsumerkeeper.Keeper,
	readyValopers map[string]bool,
) error {
	_, DAOaddrBz, err := bech32.DecodeAndConvert(FundAddress)
	if err != nil {
		return err
	}
	DAOaddr := sdk.AccAddress(DAOaddrBz)

	// -------------------------------------------------------------------------
	// Step 1: Handle existing governance/democracy validators in x/staking.
	//
	// These were never part of CometBFT consensus on the ICS chain. Do NOT call
	// Jail() on them — that would emit a consensus update for unknown keys.
	// -------------------------------------------------------------------------
	validators, err := stakingKeeper.GetAllValidators(ctx)
	if err != nil {
		return err
	}

	for _, val := range validators {
		valoper := val.GetOperator()

		if readyValopers[valoper] {
			// Fund the validator account if below MinStake so it can self-bond.
			valAddr, err := sdk.ValAddressFromBech32(valoper)
			if err != nil {
				return err
			}
			accAddr := sdk.AccAddress(valAddr)

			balance := bankKeeper.GetBalance(ctx, accAddr, Denom)
			if balance.Amount.LT(math.NewInt(MinStake)) {
				coins := sdk.NewCoins(sdk.NewCoin(Denom, math.NewInt(MinStake)))
				if err := bankKeeper.SendCoins(ctx, DAOaddr, accAddr, coins); err != nil {
					return err
				}
			}
			// Validator stays bonded — staking will emit a proper power>0 update.
		} else {
			// NOT ready: set unbonded directly in the store.
			// Calling Jail() would emit a power=0 CometBFT update for a key
			// CometBFT has never seen → "validator does not exist" panic.
			val.Status = stakingtypes.Unbonded
			if err := stakingKeeper.SetValidator(ctx, val); err != nil {
				return fmt.Errorf("failed to unbond non-ready governor %s: %w", valoper, err)
			}
			if err := stakingKeeper.DeleteValidatorByPowerIndex(ctx, val); err != nil {
				return fmt.Errorf("failed to remove power index for %s: %w", valoper, err)
			}
		}
	}

	// -------------------------------------------------------------------------
	// Step 2: Register ICS validators into x/staking and force-bond them so
	// staking emits power=0 updates to CometBFT, cleanly evicting them.
	//
	// ICS validators have 1uinto stake → VP = 1/1_000_000 = 0 in staking math.
	// Staking endblocker will remove them from the active set this same block.
	// -------------------------------------------------------------------------
	consumerValidators := consumerKeeper.GetAllCCValidator(ctx)
	if err := moveICSToStaking(ctx, stakingKeeper, bankKeeper, DAOaddr, consumerValidators); err != nil {
		return err
	}

	return nil
}

// moveICSToStaking registers each ICS (CCV) validator into x/staking with
// minimal stake and immediately bonds them, causing staking to emit a power=0
// validator update to CometBFT. This is the only safe way to evict ICS
// validators from CometBFT's active set during the sovereignty migration.
func moveICSToStaking(
	ctx sdk.Context,
	sk stakingkeeper.Keeper,
	bk bankkeeper.Keeper,
	DAOaddr sdk.AccAddress,
	consumerValidators []ccvconsumertypes.CrossChainValidator,
) error {
	srv := stakingkeeper.NewMsgServerImpl(&sk)

	skippedValidators := 0

	for i, v := range consumerValidators {
		accAddr := v.GetAddress()

		// Fund the ICS validator address so it can self-bond 1uinto.
		if err := bk.SendCoins(ctx, DAOaddr, accAddr, sdk.NewCoins(sdk.NewCoin(Denom, math.NewInt(ICSSelfStake)))); err != nil {
			return fmt.Errorf("failed to fund ICS validator %s: %w", accAddr, err)
		}

		valoperAddr, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), accAddr)
		if err != nil {
			return err
		}

		// ConsPubKey() returns (cryptotypes.PubKey, error) on the concrete type.
		consPubKey, err := v.ConsPubKey()
		if err != nil {
			return err
		}

		// Skip if this consensus key is already registered (e.g. an ICS validator
		// that chose to reuse their key as a sovereign validator).
		_, lookupErr := sk.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(consPubKey.(cryptotypes.PubKey)))
		if lookupErr == nil {
			skippedValidators++
			continue
		} else if !errors.Is(lookupErr, stakingtypes.ErrNoValidatorFound) {
			return lookupErr
		}

		// GetPubkey() returns *codectypes.Any which is exactly what MsgCreateValidator.Pubkey expects.
		_, err = srv.CreateValidator(ctx, &stakingtypes.MsgCreateValidator{
			Description: stakingtypes.Description{
				Moniker: fmt.Sprintf("ics %d", i),
			},
			Commission: stakingtypes.CommissionRates{
				Rate:          math.LegacyMustNewDecFromStr("0.1"),
				MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
				MaxChangeRate: math.LegacyMustNewDecFromStr("0.1"),
			},
			MinSelfDelegation: math.NewInt(1),
			ValidatorAddress:  valoperAddr,
			Pubkey:            v.GetPubkey(),
			Value:             sdk.NewCoin(Denom, math.NewInt(ICSSelfStake)),
		})
		if err != nil {
			return fmt.Errorf("failed to create ICS validator %s in staking: %w", valoperAddr, err)
		}

		// SetLastValidatorPower so staking knows this key previously had power,
		// ensuring it emits a power=0 update to CometBFT when the endblocker
		// sees VP = 1uinto / 1_000_000 = 0.
		valAddr := sdk.ValAddress(accAddr)
		if err := sk.SetLastValidatorPower(ctx, valAddr, 1); err != nil {
			return fmt.Errorf("failed to set last validator power for %s: %w", valoperAddr, err)
		}

		savedVal, err := sk.GetValidator(ctx, valAddr)
		if err != nil {
			return fmt.Errorf("failed to get saved ICS validator %s: %w", valoperAddr, err)
		}

		// Force-bond so the validator enters the active set this block.
		// Staking will then compute VP=0 and emit the removal update to CometBFT.
		if _, err := bondValidator(ctx, sk, savedVal); err != nil {
			return fmt.Errorf("failed to bond ICS validator %s: %w", valoperAddr, err)
		}
	}

	// Reconcile pool balances: we force-bonded ICS validators whose stake was
	// sitting in NotBondedPool, so move it to BondedPool.
	bondedCount := len(consumerValidators) - skippedValidators
	if bondedCount > 0 {
		coins := sdk.NewCoins(sdk.NewCoin(Denom, math.NewInt(int64(bondedCount)*ICSSelfStake))) //nolint:gosec
		if err := bk.SendCoinsFromModuleToModule(ctx, stakingtypes.NotBondedPoolName, stakingtypes.BondedPoolName, coins); err != nil {
			return fmt.Errorf("failed to reconcile bonded pool: %w", err)
		}
	}

	return nil
}

// bondValidator force-transitions a validator to Bonded status, bypassing the
// normal unbonding queue. This is intentional during the ICS→sovereign upgrade
// so that staking emits a validator update to CometBFT in the same block.
//
// Copied from cosmos-sdk staking keeper val_state_change.go.
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
