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
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	ccvconsumerkeeper "github.com/cosmos/interchain-security/v6/x/ccv/consumer/keeper"
	types2 "github.com/cosmos/interchain-security/v6/x/ccv/consumer/types"
)

const (
	// set of constants defines self delegation amount for newly created validator
	// ICS and Sovereign ones
	SovereignMinSelfDelegation = 1_000_000
	SovereignSelfStake         = 1_000_000
	ICSMinSelfDelegation       = 1
	ICSSelfStake               = 1
	DefaultDenom               = "uinto"
	FundAddress                = "into1mhd977xqvd8pl7efsrtyltucw0dhf7h4mpv2ve"
)

//go:embed validators/staking
var Vals embed.FS

type StakingValidator struct {
	Moniker         string         `json:"moniker"`
	Valoper         string         `json:"valoper"`
	PK              ed25519.PubKey `json:"pk"`
	Identity        string         `json:"identity,omitempty"`
	Website         string         `json:"website,omitempty"`
	SecurityContact string         `json:"security_contact,omitempty"`
	Details         string         `json:"details,omitempty"`
}

func GatherStakingMsgs() ([]types.MsgCreateValidator, error) {
	msgs := make([]types.MsgCreateValidator, 0)
	errWalk := fs.WalkDir(Vals, ".", func(path string, d fs.DirEntry, err error) error {
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
		skval := StakingValidator{}
		err = json.Unmarshal(data, &skval)
		if err != nil {
			return err
		}
		msg := StakingValMsg(skval.Moniker, SovereignSelfStake, skval.Valoper, skval.PK, skval.Identity, skval.Website, skval.SecurityContact, skval.Details)
		msgs = append(msgs, msg)

		return nil
	})
	return msgs, errWalk
}

func StakingValMsg(moniker string, stake int64, valoper string, pk ed25519.PubKey, identity, website, securityContact, details string) types.MsgCreateValidator {
	pubkey, err := codectypes.NewAnyWithValue(&pk)
	if err != nil {
		panic(err)
	}
	return types.MsgCreateValidator{
		Description: types.Description{
			Moniker:         moniker,
			Identity:        identity,
			Website:         website,
			SecurityContact: securityContact,
			Details:         details,
		},
		Commission: types.CommissionRates{
			Rate:          math.LegacyMustNewDecFromStr("0.1"),
			MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
			MaxChangeRate: math.LegacyMustNewDecFromStr("0.1"),
		},
		MinSelfDelegation: math.NewInt(SovereignMinSelfDelegation),
		DelegatorAddress:  "",
		// WARN: Operator must have enough funds
		ValidatorAddress: valoper,
		Pubkey:           pubkey,
		Value: sdk.Coin{
			Denom:  DefaultDenom,
			Amount: math.NewInt(stake),
		},
	}
}

func GetReadyValidators() (map[string]bool, error) {
	msgs, err := GatherStakingMsgs()
	if err != nil {
		return nil, err
	}
	ready := make(map[string]bool)
	for _, msg := range msgs {
		ready[msg.ValidatorAddress] = true
	}
	return ready, nil
}

// GatherGovernorStakingMsgs builds MsgCreateValidator messages for all
// "governor" validators that are explicitly marked as ready for sovereign
// validation.
//
// readyValopers is a set of valoper strings (e.g. "intovaloper1...") that
// have confirmed they will run a sovereign node. Only those will be migrated.
//
// It returns:
//   - msgs: MsgCreateValidator messages for ready governors
//   - notMigrated: valoper addresses for governors that were *not* migrated
//     (either not in readyValopers or skipped for other reasons)
func GatherGovernorStakingMsgs(
	ctx sdk.Context,
	stakingKeeper stakingkeeper.Keeper,
	readyValopers map[string]bool,
) (msgs []types.MsgCreateValidator, notMigrated [][]byte, err error) {
	msgs = make([]types.MsgCreateValidator, 0)
	notMigrated = make([][]byte, 0)

	govVals, err := stakingKeeper.GetAllValidators(ctx)
	if err != nil {
		return nil, nil, err
	}

	for _, govVal := range govVals {
		valoperAddr := govVal.GetOperator() // sdk.ValAddress
		consAddr, err := govVal.GetConsAddr()
		if err != nil {
			return nil, nil, err
		}
		// 1) Skip jailed governors outright
		if govVal.IsJailed() {
			notMigrated = append(notMigrated, consAddr)
			continue
		}

		// 2) Only migrate if explicitly marked as ready
		valoperStr := valoperAddr // "intovaloper1..."
		if readyValopers != nil && !readyValopers[valoperStr] {
			// They did not signal readiness / not on the allowlist
			notMigrated = append(notMigrated, consAddr)
			continue
		}

		// 3) Use existing description for the sovereign validator
		desc := govVal.Description

		// 4) Pack consensus pubkey
		consPk, err := govVal.ConsPubKey()
		if err != nil {
			return nil, nil, fmt.Errorf(
				"gather governor msgs: failed cons pubkey for %s: %w",
				valoperStr, err,
			)
		}

		pkAny, err := codectypes.NewAnyWithValue(consPk)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"gather governor msgs: failed to pack pubkey for %s: %w",
				valoperStr, err,
			)
		}

		// 5) Build MsgCreateValidator with same params as sovereign vals
		msg := types.MsgCreateValidator{
			Description: desc,
			Commission: types.CommissionRates{
				Rate:          math.LegacyMustNewDecFromStr("0.1"),
				MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
				MaxChangeRate: math.LegacyMustNewDecFromStr("0.1"),
			},
			MinSelfDelegation: math.NewInt(SovereignMinSelfDelegation),
			DelegatorAddress:  "",
			ValidatorAddress:  valoperStr,
			Pubkey:            pkAny,
			Value: sdk.Coin{
				Denom:  DefaultDenom,
				Amount: math.NewInt(SovereignSelfStake),
			},
		}

		msgs = append(msgs, msg)
	}

	return msgs, notMigrated, nil
}

// MoveICSToStaking creates CCV (ICS) validators in staking module, forces them to
// change status to bonded to generate valsetupdate with 0 vp in the same block.
func MoveICSToStaking(ctx sdk.Context, sk stakingkeeper.Keeper, bk bankkeeper.Keeper, consumerValidators []types2.CrossChainValidator) error {

	skippedValidators := 0
	srv := stakingkeeper.NewMsgServerImpl(&sk)
	_, DAOaddrBz, err := bech32.DecodeAndConvert(FundAddress)
	if err != nil {
		return err
	}
	DAOaddr := sdk.AccAddress(DAOaddrBz)

	// Add all ICS validators to staking module
	for i, v := range consumerValidators {
		// funding ICS valopers from DAO to stake a coin
		err := bk.SendCoins(ctx, DAOaddr, v.GetAddress(), sdk.NewCoins(sdk.Coin{
			Denom:  DefaultDenom,
			Amount: math.NewInt(ICSSelfStake),
		}))
		if err != nil {
			return err
		}

		valoperAddr, err := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), v.GetAddress())
		if err != nil {
			return err
		}

		consPubKey, err := v.ConsPubKey()
		if err != nil {
			return err
		}

		_, err = sk.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(consPubKey))
		if err == nil {
			// The validator is already created during the configuring of the sovereign/staking validator group.
			// This is possible if an ICS validator moves to a sovereign group and decides to use the same pubkey after the transition.
			skippedValidators++
			continue
		} else if !errors.Is(err, types.ErrNoValidatorFound) {
			return err
		}

		_, err = srv.CreateValidator(ctx, &types.MsgCreateValidator{
			Description: types.Description{
				Moniker:         fmt.Sprintf("ics %d", i),
				Identity:        "",
				Website:         "",
				SecurityContact: "",
				Details:         "",
			},
			Commission: types.CommissionRates{
				Rate:          math.LegacyMustNewDecFromStr("0.1"),
				MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
				MaxChangeRate: math.LegacyMustNewDecFromStr("0.1"),
			},
			MinSelfDelegation: math.NewInt(ICSMinSelfDelegation),
			// WARN: valoper must have enough funds to selfbond
			ValidatorAddress: valoperAddr,
			Pubkey:           v.GetPubkey(),
			Value: sdk.Coin{
				Denom:  DefaultDenom,
				Amount: math.NewInt(ICSSelfStake),
			},
		})
		if err != nil {
			return err
		}

		err = sk.SetLastValidatorPower(ctx, v.GetAddress(), 1)
		if err != nil {
			return err
		}

		savedVal, err := sk.GetValidator(ctx, v.GetAddress())
		if err != nil {
			return err
		}
		// add validator to active set to remove him from endblocker the same block
		// validator will be kicked out from active set due to the fact, voting power is calculated as `staked_amount/1_000_000`
		// staked amount for ICS validator is 1uinto => vp = 0
		_, err = bondValidator(ctx, sk, savedVal)
		if err != nil {
			return err
		}
	}

	coins := sdk.NewCoins(sdk.NewCoin(DefaultDenom, math.NewInt(int64((len(consumerValidators)-skippedValidators)*ICSSelfStake))))
	// since we forced to set bond status for ics validators during the upgrade, we have to move ICS staked funds from NotBondedPoolName to BondedPoolName
	return bk.SendCoinsFromModuleToModule(ctx, types.NotBondedPoolName, types.BondedPoolName, coins)
}

// DeICS - does the deics. The whole point of the method is to force the staking module
// to remove ICS validators and add STAKING (sovereign) ones, by generating valset updates.
// We add STAKING and ICS to staking module in a special way.
// STAKING added in natural staking way, just submit `MsgCreateValidator` msg with vp >=1.
// ICS added with the message (with vp = 0). And to force staking to remove the ICS
// validators the same block, we force validators to join "active" set by bonding them
// with `bondValidator` and move stake from nonbonded pool to bonded.
//
// With the updated logic, sovereign validators are built from:
//   - ready stakingKeeper ccvstaking "governor" validators (allowlisted),
//   - plus any extra JSON-based validators from validators/staking (optional).
//
// Signature extended to include:
//   - stakingKeeper: stakingKeeper keeper for governor set
//   - readyValopers: allowlist of governors to migrate to sovereign set
func DeICS(
	ctx sdk.Context,
	sk stakingkeeper.Keeper,
	consumerKeeper ccvconsumerkeeper.Keeper,
	stakingKeeper stakingkeeper.Keeper,
	bk bankkeeper.Keeper,
	readyValopers map[string]bool,
) error {
	srv := stakingkeeper.NewMsgServerImpl(&sk)
	consumerValidators := consumerKeeper.GetAllCCValidator(ctx)

	// msgs to create new sovereign validators from governors (allowlisted)
	govValMsgs, notMigratedGovs, err := GatherGovernorStakingMsgs(ctx, stakingKeeper, readyValopers)
	if err != nil {
		return err
	}

	// optional extra validators from embedded JSON
	extraValMsgs, err := GatherStakingMsgs()
	if err != nil {
		return err
	}

	// final list of new sovereign validators
	newValMsgs := append(govValMsgs, extraValMsgs...)

	cp := consumerKeeper.GetConsumerParams(ctx)

	p := types.Params{
		UnbondingTime: cp.UnbondingPeriod,
		// During migration MaxValidators MUST be >= all the validators number, old and new ones.
		// i.e. chain managed by 150 ICS validators, and we are switching to 70 STAKING, MaxValidators MUST be at least 220,
		// otherwise panic during staking begin blocker happens
		// It's allowed to change the value at the very next block
		MaxValidators:     uint32(len(consumerValidators) + len(newValMsgs)), //nolint:gosec
		MaxEntries:        7,
		HistoricalEntries: 10_000,
		BondDenom:         DefaultDenom,
		MinCommissionRate: math.LegacyMustNewDecFromStr("0.0"),
	}

	_, err = srv.UpdateParams(ctx, &types.MsgUpdateParams{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Params:    p,
	})
	if err != nil {
		return err
	}

	_, DAOaddrBz, err := bech32.DecodeAndConvert(FundAddress)
	if err != nil {
		return err
	}
	DAOaddr := sdk.AccAddress(DAOaddrBz)

	// Create sovereign validators (governors that opted in + any extras)
	// Deduplicate validators by address to prevent double funding/processing
	processedValidators := make(map[string]bool)

	for _, msg := range newValMsgs {
		_, valAddrBz, err := bech32.DecodeAndConvert(msg.ValidatorAddress)
		if err != nil {
			return err
		}
		valAddr := sdk.ValAddress(valAddrBz)
		accAddr := sdk.AccAddress(valAddrBz)

		// Check if we already processed this validator in this loop
		if processedValidators[valAddr.String()] {
			continue
		}
		processedValidators[valAddr.String()] = true

		// Re-encode address to use the current global prefix to avoid 'hrp does not match' errors
		msg.ValidatorAddress = valAddr.String()

		// prefund validator to make selfbond
		err = bk.SendCoins(ctx, DAOaddr, accAddr, sdk.NewCoins(sdk.Coin{
			Denom:  DefaultDenom,
			Amount: math.NewInt(SovereignSelfStake),
		}))
		if err != nil {
			return err
		}

		_, err = sk.GetValidator(ctx, valAddr)
		if err == nil {
			// Validator already exists, skip creation
			continue
		}

		_, err = srv.CreateValidator(ctx, &msg)
		if err != nil {
			return err
		}
	}

	// Run ICS -> staking bridge logic to remove ICS validators from the set
	err = MoveICSToStaking(ctx, sk, bk, consumerValidators)
	if err != nil {
		return err
	}

	// Optionally jail non-migrated governors so they don't pretend to be active validators.
	// Adjust to your actual ccvstaking keeper API if the method name/signature differs.
	for _, valAddrCons := range notMigratedGovs {
		if err := stakingKeeper.Jail(ctx, valAddrCons); err != nil {
			// Don't fail the whole migration if we can't jail someone; just log.
			fmt.Printf("failed to jail non-migrated governor %s: %v\n", valAddrCons, err)
		}
	}

	return nil
}

// copied from staking module https://github.com/cosmos/cosmos-sdk/blob/v0.50.6/x/staking/keeper/val_state_change.go#L336
func bondValidator(ctx context.Context, k stakingkeeper.Keeper, validator types.Validator) (types.Validator, error) {
	// delete the validator by power index, as the key will change
	if err := k.DeleteValidatorByPowerIndex(ctx, validator); err != nil {
		return validator, err
	}

	validator = validator.UpdateStatus(types.Bonded)

	// save the now bonded validator record to the two referenced stores
	if err := k.SetValidator(ctx, validator); err != nil {
		return validator, err
	}

	if err := k.SetValidatorByPowerIndex(ctx, validator); err != nil {
		return validator, err
	}

	// delete from queue if present
	if err := k.DeleteValidatorQueue(ctx, validator); err != nil {
		return validator, err
	}

	// trigger hook
	consAddr, err := validator.GetConsAddr()
	if err != nil {
		return validator, err
	}
	codec := address.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())
	str, err := codec.StringToBytes(validator.GetOperator())
	if err != nil {
		return validator, err
	}

	if err := k.Hooks().AfterValidatorBonded(ctx, consAddr, str); err != nil {
		return validator, err
	}

	return validator, err
}
