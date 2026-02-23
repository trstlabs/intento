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
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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
// Two separate validator sets exist on an ICS democracy chain:
//
//  1. x/staking — local governance/democracy validators. Never part of CometBFT
//     consensus. Their consensus keys are unknown to CometBFT.
//
//  2. x/ccv/consumer — ICS validators assigned by the provider (Cosmos Hub).
//     These drove CometBFT consensus. CometBFT only knows these keys.
//
// Migration strategy (validated against Neutron's v6 reference implementation):
//
//   - UpdateParams first: set MaxValidators >= ICS + sovereign count and copy
//     UnbondingTime from consumer params. Without this the begin blocker panics.
//
//   - Ready governance validators: fund their account so they can self-bond.
//     Leave them Bonded — staking emits proper power>0 updates to CometBFT.
//
//   - Non-ready governance validators: zero LastValidatorPower and remove from
//     power index. DO NOT call Jail() (emits update for key CometBFT never saw →
//     "validator does not exist" panic) and DO NOT set status=Unbonded directly
//     (endblocker sees last power>0, attempts bondedToUnbonding on already-Unbonded
//     validator → "bad state transition" panic).
//
//   - ICS validators: register into staking using v.GetAddress() throughout
//     (matching Neutron's reference — GetAddress() returns the same address the
//     ICS consumer module used when calling SetLastValidatorPower). Force-bond
//     them so staking emits power=0 removal updates to CometBFT this block.
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

	consumerValidators := consumerKeeper.GetAllCCValidator(ctx)
	srv := stakingkeeper.NewMsgServerImpl(&stakingKeeper)

	// -------------------------------------------------------------------------
	// Step 1: UpdateParams — must happen before any validator set manipulation.
	//
	// MaxValidators must be >= ICS validators + ready sovereign validators to
	// avoid a begin blocker panic. UnbondingTime is copied from consumer params
	// so the chain has a valid unbonding period from day one as sovereign.
	// -------------------------------------------------------------------------
	cp := consumerKeeper.GetConsumerParams(ctx)
	readyCount := len(readyValopers)
	_, err = srv.UpdateParams(ctx, &stakingtypes.MsgUpdateParams{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Params: stakingtypes.Params{
			UnbondingTime: cp.UnbondingPeriod,
			// Must cover all validators alive during this block: ICS (being
			// removed) + sovereign (being added). Safe to reduce next block.
			MaxValidators:     uint32(len(consumerValidators) + readyCount + 10), //nolint:gosec
			MaxEntries:        7,
			HistoricalEntries: 10_000,
			BondDenom:         Denom,
			MinCommissionRate: math.LegacyMustNewDecFromStr("0.0"),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update staking params: %w", err)
	}

	// -------------------------------------------------------------------------
	// Step 2: Handle governance/democracy validators in x/staking.
	//
	// These were never part of CometBFT consensus on the ICS consumer chain.
	// Their keys are invisible to CometBFT — do not emit any consensus updates
	// for them.
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
			// Fund account so the validator can self-bond as sovereign.
			accAddr := sdk.AccAddress(valAddr)
			balance := bankKeeper.GetBalance(ctx, accAddr, Denom)
			if balance.Amount.LT(math.NewInt(MinStake)) {
				coins := sdk.NewCoins(sdk.NewCoin(Denom, math.NewInt(MinStake)))
				if err := bankKeeper.SendCoins(ctx, DAOaddr, accAddr, coins); err != nil {
					return fmt.Errorf("failed to fund ready validator %s: %w", valoper, err)
				}
			}
			// Stays Bonded — staking endblocker emits power>0 update to CometBFT.
		} else {
			// Non-ready governance validator. These validators ARE in CometBFT's
			// active set (they participated in ICS consensus alongside the pure
			// ICS validators). CometBFT must receive an explicit power=0 update
			// for each of them or it will crash at commit with "failed to find
			// validator to remove".
			//
			// Correct approach: only delete from the power index. Leave
			// LastValidatorPower untouched. The staking endblocker will:
			//   1. Find this validator in the `last` map (last power > 0)
			//   2. See current power = 0 (absent from power index)
			//   3. Execute the valid Bonded → Unbonding transition
			//   4. Emit a power=0 update to CometBFT ✓
			//
			// Do NOT zero LastValidatorPower — that would make the endblocker
			// skip the validator entirely, leaving CometBFT with a stale entry.
			// Do NOT set status=Unbonded directly — that causes the endblocker
			// to find it already Unbonded mid-transition → "bad state transition
			// bondedToUnbonding" panic.
			if err := stakingKeeper.DeleteValidatorByPowerIndex(ctx, val); err != nil {
				return fmt.Errorf("failed to remove power index for %s: %w", valoper, err)
			}
			// LastValidatorPower intentionally left as-is.
		}
	}

	// -------------------------------------------------------------------------
	// Step 3: Register ICS validators into x/staking and force-bond them so
	// staking emits power=0 removal updates to CometBFT this block.
	//
	// Following Neutron's reference implementation exactly:
	//   - Use v.GetAddress() for both funding and as the valoper address.
	//     This is the address the ICS consumer module used when it called
	//     SetLastValidatorPower, so GetLastValidatorPower lookups will match.
	//   - ICS validators get 1uinto stake → VP = 1/1_000_000 = 0.
	//   - Force-bond so they enter the active set; endblocker kicks them out
	//     the same block and emits the power=0 update to CometBFT.
	// -------------------------------------------------------------------------
	if err := moveICSToStaking(ctx, stakingKeeper, bankKeeper, DAOaddr, consumerValidators); err != nil {
		return err
	}

	return nil
}

// moveICSToStaking registers ICS validators into x/staking and force-bonds
// them so the staking endblocker emits power=0 removal updates to CometBFT.
//
// Address handling follows Neutron's reference implementation (v5/v6):
// v.GetAddress() is used for both funding and as the staking valoper address.
// This matches the address the ICS consumer module used internally when it
// called SetLastValidatorPower — ensuring LastValidatorPower lookups hit the
// correct store key and CometBFT's validator set is cleanly updated.
//
// The flow per ICS validator:
//  1. Fund v.GetAddress() with 1uinto for the self-bond
//  2. CreateValidator with valoper = v.GetAddress(), stake = 1uinto (VP=0)
//  3. SetLastValidatorPower = 1 so endblocker sees a delta (1→0)
//  4. bondValidator: force to Bonded so it enters the active set this block
//  5. Endblocker: VP=0 → moves to unbonding, emits power=0 update to CometBFT
//  6. SendCoinsFromModuleToModule: reconcile NotBondedPool → BondedPool
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
		// v.GetAddress() is the address the ICS consumer module used when
		// syncing this validator's power into x/staking. Use it everywhere.
		valAddr := sdk.ValAddress(v.GetAddress())
		accAddr := sdk.AccAddress(v.GetAddress())

		valoperAddr, err := bech32.ConvertAndEncode(
			sdk.GetConfig().GetBech32ValidatorAddrPrefix(), valAddr,
		)
		if err != nil {
			return fmt.Errorf("ics validator %d: failed to encode valoper: %w", i, err)
		}

		// Fund the account that will self-bond.
		if err := bk.SendCoins(ctx, DAOaddr, accAddr, sdk.NewCoins(
			sdk.NewCoin(Denom, math.NewInt(ICSSelfStake)),
		)); err != nil {
			return fmt.Errorf("failed to fund ICS validator %s: %w", valoperAddr, err)
		}

		// Derive consAddr from the consensus pubkey for the duplicate-key check.
		// This is the only place we need the pubkey-derived address.
		consPubKey, err := v.ConsPubKey()
		if err != nil {
			return fmt.Errorf("ics validator %d: failed to get cons pubkey: %w", i, err)
		}
		consAddr := sdk.GetConsAddress(consPubKey.(cryptotypes.PubKey))

		// Skip if this consensus key is already registered — an ICS validator
		// that reused their key when joining the sovereign set.
		if _, lookupErr := sk.GetValidatorByConsAddr(ctx, consAddr); lookupErr == nil {
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
			MinSelfDelegation: math.NewInt(ICSSelfStake),
			ValidatorAddress:  valoperAddr,
			Pubkey:            v.GetPubkey(),
			Value:             sdk.NewCoin(Denom, math.NewInt(ICSSelfStake)),
		})
		if err != nil {
			return fmt.Errorf("failed to create ICS validator %s: %w", valoperAddr, err)
		}

		// Set last power=1 so the endblocker sees a delta (1→0) and emits a
		// clean power=0 removal update to CometBFT for this key.
		if err := sk.SetLastValidatorPower(ctx, valAddr, 1); err != nil {
			return fmt.Errorf("failed to set last power for %s: %w", valoperAddr, err)
		}

		savedVal, err := sk.GetValidator(ctx, valAddr)
		if err != nil {
			return fmt.Errorf("failed to get ICS validator %s after creation: %w", valoperAddr, err)
		}

		// Force-bond: validator enters active set this block with VP=0.
		// Endblocker moves it to unbonding and emits the power=0 CometBFT update.
		if _, err := bondValidator(ctx, sk, savedVal); err != nil {
			return fmt.Errorf("failed to bond ICS validator %s: %w", valoperAddr, err)
		}
	}

	// Reconcile pool balances: force-bonded validators had their stake sitting
	// in NotBondedPool. Move exactly that amount to BondedPool to keep the
	// module account invariants intact.
	bondedCount := len(consumerValidators) - skippedValidators
	coins := sdk.NewCoins(sdk.NewCoin(Denom, math.NewInt(int64(bondedCount)*ICSSelfStake))) //nolint:gosec
	return bk.SendCoinsFromModuleToModule(ctx, stakingtypes.NotBondedPoolName, stakingtypes.BondedPoolName, coins)
}

// bondValidator force-transitions a validator to Bonded status, bypassing the
// normal unbonding queue. This is intentional during the ICS→sovereign upgrade:
// we need the validator in the active set so staking emits a validator update
// to CometBFT in the same block.
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
