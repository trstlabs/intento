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
// Fix:
//
//   - Sovereign/ready validators: fund them so they can self-bond. They stay
//     Bonded and staking emits proper power>0 updates to CometBFT.
//
//   - Non-ready governance validators: remove from power index and zero
//     LastValidatorPower. The endblocker only processes validators where
//     last power != current power. Zeroing last power means no state
//     transition is attempted, avoiding the "bad state transition
//     bondedToUnbonding" panic that occurs when setting status=Unbonded
//     directly on a previously-bonded validator.
//
//   - ICS validators: registered into staking with 1uinto stake (VP=0) and
//     force-bonded so staking emits a power=0 update to CometBFT, cleanly
//     evicting them from CometBFT's active set.
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
	// These were never part of CometBFT consensus on the ICS chain.
	// -------------------------------------------------------------------------
	validators, err := stakingKeeper.GetAllValidators(ctx)
	if err != nil {
		return err
	}

	for _, val := range validators {
		valoper := val.GetOperator()
		valAddr, err := sdk.ValAddressFromBech32(valoper)
		if err != nil {
			return err
		}

		if readyValopers[valoper] {
			// Fund if below MinStake so the validator can self-bond as sovereign.
			accAddr := sdk.AccAddress(valAddr)
			balance := bankKeeper.GetBalance(ctx, accAddr, Denom)
			if balance.Amount.LT(math.NewInt(MinStake)) {
				coins := sdk.NewCoins(sdk.NewCoin(Denom, math.NewInt(MinStake)))
				if err := bankKeeper.SendCoins(ctx, DAOaddr, accAddr, coins); err != nil {
					return err
				}
			}
			// Validator stays Bonded — staking emits a proper power>0 update.
		} else {
			// NOT ready. We must NOT:
			//   - Call Jail(): emits power=0 comet update for a key comet has
			//     never seen → "validator does not exist" panic.
			//   - Set status=Unbonded directly: endblocker sees last power>0,
			//     tries bondedToUnbonding, finds already Unbonded →
			//     "bad state transition bondedToUnbonding" panic (seen in logs:
			//     val_state_change.go:276).
			//
			// Correct approach: remove from the power index and zero
			// LastValidatorPower. The endblocker (ApplyAndReturnValidatorSetUpdates)
			// only processes validators where last power != current computed power.
			// With last power=0 and validator absent from the power index,
			// it is ignored entirely this block.
			if err := stakingKeeper.DeleteValidatorByPowerIndex(ctx, val); err != nil {
				return fmt.Errorf("failed to remove power index for %s: %w", valoper, err)
			}
			if err := stakingKeeper.SetLastValidatorPower(ctx, valAddr, 0); err != nil {
				return fmt.Errorf("failed to zero last power for %s: %w", valoper, err)
			}
		}
	}

	// -------------------------------------------------------------------------
	// Step 2: Register ICS validators into x/staking and force-bond them so
	// staking emits power=0 updates to CometBFT, cleanly evicting them.
	//
	// ICS validators have 1uinto stake → VP = 1/1_000_000 = 0 in staking math.
	// Staking endblocker removes them from the active set this same block.
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

		// Fund so the ICS validator can self-bond 1uinto.
		if err := bk.SendCoins(ctx, DAOaddr, accAddr, sdk.NewCoins(sdk.NewCoin(Denom, math.NewInt(ICSSelfStake)))); err != nil {
			return fmt.Errorf("failed to fund ICS validator %s: %w", accAddr, err)
		}

		valoperAddr, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), accAddr)
		if err != nil {
			return err
		}

		consPubKey, err := v.ConsPubKey()
		if err != nil {
			return err
		}

		// Skip if consensus key already registered — an ICS validator that
		// reused their key when joining the sovereign set.
		_, lookupErr := sk.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(consPubKey.(cryptotypes.PubKey)))
		if lookupErr == nil {
			skippedValidators++
			continue
		} else if !errors.Is(lookupErr, stakingtypes.ErrNoValidatorFound) {
			return lookupErr
		}

		_, err = srv.CreateValidator(ctx, &stakingtypes.MsgCreateValidator{
			Description: stakingtypes.Description{Moniker: fmt.Sprintf("ics %d", i)},
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

		valAddr := sdk.ValAddress(accAddr)

		// Only emit a removal update to CometBFT if CometBFT actually had this
		// validator in its set. We detect this by checking LastValidatorPower:
		// the ICS consumer module sets this for validators it tracks. If it is
		// already 0 (or absent), CometBFT has no record of this key — emitting
		// a power=0 update would cause "failed to find validator to remove".
		existingPower, err := sk.GetLastValidatorPower(ctx, valAddr)
		if err != nil || existingPower == 0 {
			// CometBFT does not know this key. The validator is registered in
			// staking with VP=0 and will fall out naturally. No comet update needed.
			skippedValidators++
			continue
		}

		// CometBFT had this validator. Set last power=1 so the endblocker sees
		// a delta (1 → 0) and emits a clean power=0 removal update to CometBFT.
		if err := sk.SetLastValidatorPower(ctx, valAddr, 1); err != nil {
			return fmt.Errorf("failed to set last validator power for %s: %w", valoperAddr, err)
		}

		savedVal, err := sk.GetValidator(ctx, valAddr)
		if err != nil {
			return fmt.Errorf("failed to get saved ICS validator %s: %w", valoperAddr, err)
		}

		// Force-bond so the validator enters the active set this block.
		// With VP=0 (1uinto / 1_000_000), the endblocker immediately moves
		// it to unbonding and emits the removal update to CometBFT.
		if _, err := bondValidator(ctx, sk, savedVal); err != nil {
			return fmt.Errorf("failed to bond ICS validator %s: %w", valoperAddr, err)
		}
	}

	// Reconcile pool balances: force-bonded validators had stake in
	// NotBondedPool; move it to BondedPool to keep invariants intact.
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
// normal unbonding queue. Intentional during the ICS→sovereign upgrade so that
// staking emits a validator update to CometBFT in the same block.
//
// Copied from cosmos-sdk x/staking/keeper/val_state_change.go.
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
