package v1100

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"strings"

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

// cometValidatorSet is the exact set of validator consensus addresses (hex)
// that CometBFT held at height 10632699 — the block immediately before the
// upgrade. Source: GET /validators?height=10632699&per_page=100 on the live chain.
//
// This is the ground truth for which validators need explicit power=0 removal
// updates. Any validator not in this set must NOT receive a removal update —
// doing so causes "failed to find validator to remove" at commit.
//
// Includes both dual-role validators (in x/staking + ICS consumer set) and
// pure ICS validators (only in consumer set). Does NOT include the slashed
// validator 9DA88F50 which was in the consumer store but already removed from
// CometBFT's set before the upgrade.
var cometValidatorSet = map[string]bool{
	"37A59AE2151D67B83C1AAF49A78070EDF8336060": true, // NodeStake
	"85A85F24E9B212F0B7CBBD1F7F032B8775B69258": true, // POSTHUMAN
	"857E6467573CCBE514B631A8DD6D20286CC28D21": true, // pure ICS
	"0C1551FDE29A3EA2F5FA7EFDD88C225C621B8402": true, // Interstellar
	"1345682C0374388FD4885F7C2AA2933A78E4096B": true, // ECO Stake
	"C082B0D6060B6DC1B752BB01E47961A2AC3B9AB4": true, // 01node
	"28BC56792C823F1E46FEC0C12264C7EE20CD7923": true, // Stake&Relax
	"DF866F356499AE8C7F222169FEBBC0BE94F1FCDC": true, // Atlas Staking
	"D373EF87CA0935D7AA1BFAF8770332A85B4B0252": true, // Apeiron Nodes (slashed, but still in comet set)
	"009037C2C75632F3BF9E39A11C0E81EACB262D9E": true, // pure ICS
}

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

// isInCometSet returns true if the given consensus address was in CometBFT's
// validator set at the upgrade height. Uses the hardcoded ground-truth map.
func isInCometSet(consAddr sdk.ConsAddress) bool {
	return cometValidatorSet[strings.ToUpper(fmt.Sprintf("%X", consAddr))]
}

// DeICS migrates the chain from ICS consumer to sovereign.
//
// Validator set topology at upgrade height (confirmed from live chain state):
//
//   - 8 dual-role validators: present in both x/staking (governance) and the
//     ICS consumer set. They were signing blocks — CometBFT knows their keys.
//
//   - 2 pure ICS validators: only in the consumer store, not in x/staking.
//     Were signing blocks — CometBFT knows their keys.
//
//   - 7 governance-only validators: in x/staking but NOT in CometBFT's set.
//     Never participated in local consensus.
//
//   - 1 ghost validator: in consumer store with LastValidatorPower>0 but
//     already removed from CometBFT's set (slashed/tombstoned before upgrade).
//     Must NOT receive a power=0 update.
//
// The cometValidatorSet map is the authoritative source for which validators
// CometBFT will accept removal updates for.
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
	// Sets MaxValidators high enough to hold all validators alive during this
	// block, and copies UnbondingTime from consumer params.
	// -------------------------------------------------------------------------
	cp := consumerKeeper.GetConsumerParams(ctx)
	_, err = srv.UpdateParams(ctx, &stakingtypes.MsgUpdateParams{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Params: stakingtypes.Params{
			UnbondingTime:     cp.UnbondingPeriod,
			MaxValidators:     uint32(len(consumerValidators) + len(readyValopers) + 10), //nolint:gosec
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
	// Step 2: Handle x/staking governance validators.
	//
	// Dual-role validators (in CometBFT set): delete from power index only.
	// LastValidatorPower is left intact so the endblocker finds them in the
	// `last` map, sees power=0 (not in index), and emits a clean power=0
	// removal update to CometBFT via the valid Bonded→Unbonding transition.
	//
	// Governance-only validators (NOT in CometBFT set): delete from power
	// index AND zero LastValidatorPower so the endblocker ignores them
	// entirely — emitting any update for these keys would crash CometBFT.
	//
	// Ready validators: fund their account so they can self-bond as sovereign.
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

		consAddr, err := val.GetConsAddr()
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
			// Stays Bonded — endblocker emits power>0 update to CometBFT.
		} else if isInCometSet(consAddr) {
			// Dual-role: in CometBFT's set. Remove from power index only.
			// Endblocker will handle the Bonded→Unbonding transition and
			// emit the power=0 removal update to CometBFT correctly.
			if err := stakingKeeper.DeleteValidatorByPowerIndex(ctx, val); err != nil {
				return fmt.Errorf("failed to remove power index for dual-role %s: %w", valoper, err)
			}
			// LastValidatorPower intentionally left intact.
		} else {
			// Governance-only: NOT in CometBFT's set. Remove from power index
			// AND zero LastValidatorPower so endblocker ignores entirely.
			if err := stakingKeeper.DeleteValidatorByPowerIndex(ctx, val); err != nil {
				return fmt.Errorf("failed to remove power index for governance %s: %w", valoper, err)
			}
			// IMPORTANT: LastValidatorPower is keyed by sdk.ValAddress(consAddr),
			// NOT by the operator valAddr from bech32. The ICS consumer module
			// sets last power using the consensus-address-derived valAddr.
			// Using the operator address here would zero the wrong key and
			// the endblocker would still find the old entry and emit an update
			// for a key CometBFT never had → "failed to find validator to remove".
			consValAddr := sdk.ValAddress(consAddr)
			if err := stakingKeeper.SetLastValidatorPower(ctx, consValAddr, 0); err != nil {
				return fmt.Errorf("failed to zero last power for governance %s: %w", valoper, err)
			}
		}
	}

	// -------------------------------------------------------------------------
	// Step 3: Register pure ICS validators (not in x/staking) into staking
	// and force-bond them so the endblocker emits power=0 removal updates.
	//
	// Uses cometValidatorSet to gate which validators actually need the
	// bond+removal treatment — validators in the consumer store but absent
	// from CometBFT's set (e.g. previously slashed/tombstoned) are registered
	// inert and skipped.
	// -------------------------------------------------------------------------
	if err := moveICSToStaking(ctx, stakingKeeper, bankKeeper, DAOaddr, consumerValidators); err != nil {
		return err
	}

	return nil
}

// moveICSToStaking registers ICS consumer validators into x/staking.
//
// For validators in cometValidatorSet: registered with 1uinto stake (VP=0),
// force-bonded so staking emits a power=0 removal update to CometBFT.
//
// For validators NOT in cometValidatorSet (e.g. previously slashed and already
// removed from CometBFT's internal set): registered inert with VP=0, no bond.
// Emitting a removal update for these would crash with "failed to find
// validator to remove".
//
// Skips validators whose consensus key is already in x/staking (dual-role
// validators handled in Step 2).
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
		// Derive consensus address from pubkey — this is how CometBFT indexes
		// validators internally.
		consPubKey, err := v.ConsPubKey()
		if err != nil {
			return fmt.Errorf("ics validator %d: failed to get cons pubkey: %w", i, err)
		}
		pk := consPubKey.(cryptotypes.PubKey)
		consAddr := sdk.ConsAddress(pk.Address())

		// Skip dual-role validators — already handled in Step 2 via x/staking.
		if _, lookupErr := sk.GetValidatorByConsAddr(ctx, consAddr); lookupErr == nil {
			continue
		} else if !errors.Is(lookupErr, stakingtypes.ErrNoValidatorFound) {
			return lookupErr
		}

		// Use consAddr-derived valAddr for staking operations.
		valAddr := sdk.ValAddress(consAddr)
		valoperAddr, err := bech32.ConvertAndEncode(
			sdk.GetConfig().GetBech32ValidatorAddrPrefix(), valAddr,
		)
		if err != nil {
			return fmt.Errorf("ics validator %d: failed to encode valoper: %w", i, err)
		}

		// Fund the account that will self-bond.
		accAddr := sdk.AccAddress(consAddr)
		if err := bk.SendCoins(ctx, DAOaddr, accAddr, sdk.NewCoins(
			sdk.NewCoin(Denom, math.NewInt(ICSSelfStake)),
		)); err != nil {
			return fmt.Errorf("failed to fund ICS validator %s: %w", valoperAddr, err)
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

		if !isInCometSet(consAddr) {
			// Not in CometBFT's set — register inert, no removal update needed.
			// (e.g. previously slashed validator still in consumer store)
			continue
		}

		// In CometBFT's set — force-bond so endblocker emits power=0 removal.
		if err := sk.SetLastValidatorPower(ctx, valAddr, 1); err != nil {
			return fmt.Errorf("failed to set last power for %s: %w", valoperAddr, err)
		}

		savedVal, err := sk.GetValidator(ctx, valAddr)
		if err != nil {
			return fmt.Errorf("failed to get ICS validator %s after creation: %w", valoperAddr, err)
		}

		if _, err := bondValidator(ctx, sk, savedVal); err != nil {
			return fmt.Errorf("failed to bond ICS validator %s: %w", valoperAddr, err)
		}

		bondedCount++
	}

	// Reconcile pool: force-bonded validators had stake in NotBondedPool.
	coins := sdk.NewCoins(sdk.NewCoin(Denom, math.NewInt(int64(bondedCount)*ICSSelfStake))) //nolint:gosec
	return bk.SendCoinsFromModuleToModule(ctx, stakingtypes.NotBondedPoolName, stakingtypes.BondedPoolName, coins)
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
