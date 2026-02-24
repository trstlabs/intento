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
	stakingValidators, err := sk.GetAllValidators(ctx)
	if err != nil {
		return err
	}

	// Build a map of ICS validators by consensus address for quick lookup
	icsValsByConsAddr := make(map[string]ccvconsumertypes.CrossChainValidator)
	for _, icsVal := range consumerValidators {
		consPubKey, _ := icsVal.ConsPubKey()
		consAddr := sdk.ConsAddress(consPubKey.Address())
		icsValsByConsAddr[consAddr.String()] = icsVal
	}

	// Process existing staking validators - only FUND the ready ones
	for _, val := range stakingValidators {
		valoper := val.GetOperator()
		valAddr, _ := sdk.ValAddressFromBech32(valoper)
		consAddr, _ := val.GetConsAddr()

		_, isICSVal := icsValsByConsAddr[sdk.ConsAddress(consAddr).String()]
		isGovReady := readyValopers[valoper]

		// Only fund validators that are governance-ready (whether ICS or not)
		if isGovReady {
			accAddr := sdk.AccAddress(valAddr)
			balance := bk.GetBalance(ctx, accAddr, Denom)
			if balance.Amount.LT(math.NewInt(MinStake)) {
				diff := math.NewInt(MinStake).Sub(balance.Amount)
				if err := bk.SendCoins(ctx, DAOaddr, accAddr, sdk.NewCoins(sdk.NewCoin(Denom, diff))); err != nil {
					return err
				}
				if isICSVal {
					fmt.Printf("✓ Funded ICS→native ready validator %s with %s\n", valoper, diff.String())
				} else {
					fmt.Printf("✓ Funded ready validator %s with %s\n", valoper, diff.String())
				}
			}
			// Remove from ICS map so we don't recreate them
			if isICSVal {
				delete(icsValsByConsAddr, sdk.ConsAddress(consAddr).String())
			}
		} else {
			// Not ready - just log, don't touch them
			// They'll have low power and won't affect consensus much
			if isICSVal {
				fmt.Printf("⚠ ICS validator %s exists but not governance-ready (leaving as-is)\n", valoper)
				delete(icsValsByConsAddr, sdk.ConsAddress(consAddr).String())
			} else {
				fmt.Printf("⚠ Non-ICS validator %s not governance-ready (leaving as-is)\n", valoper)
			}
		}
	}

	// Now handle ICS validators that DON'T exist in staking yet
	return createNewICSValidators(ctx, sk, bk, DAOaddr, icsValsByConsAddr, readyValopers)
}

func createNewICSValidators(
	ctx sdk.Context,
	sk stakingkeeper.Keeper,
	bk bankkeeper.Keeper,
	DAOaddr sdk.AccAddress,
	icsValsByConsAddr map[string]ccvconsumertypes.CrossChainValidator,
	readyValopers map[string]bool,
) error {
	srv := stakingkeeper.NewMsgServerImpl(&sk)

	i := 0
	for _, icsVal := range icsValsByConsAddr {
		consPubKey, _ := icsVal.ConsPubKey()
		consAddr := sdk.ConsAddress(consPubKey.Address())
		valAddr := sdk.ValAddress(consAddr)
		valoperAddr, _ := bech32.ConvertAndEncode(sdk.GetConfig().GetBech32ValidatorAddrPrefix(), valAddr)

		// Check if this ICS validator is governance-ready
		isReady := readyValopers[valoperAddr]

		var fundAmount math.Int
		var minSelfDel math.Int
		var moniker string

		if isReady {
			// Ready validators get full funding
			fundAmount = math.NewInt(MinStake)
			minSelfDel = math.NewInt(MinStake)
			moniker = fmt.Sprintf("ics-%d", i)
			fmt.Printf("✓ Creating ready ICS validator %s with %s stake\n", moniker, fundAmount.String())
		} else {
			// Non-ready validators get minimal funding (1 token)
			// They'll exist but have negligible voting power
			fundAmount = math.NewInt(ICSSelfStake)
			minSelfDel = math.NewInt(ICSSelfStake)
			moniker = fmt.Sprintf("ics-unready-%d", i)
			fmt.Printf("⚠ Creating non-ready ICS validator %s with minimal stake\n", moniker)
		}

		// Fund the account
		accAddr := sdk.AccAddress(consAddr)
		if err := bk.SendCoins(ctx, DAOaddr, accAddr, sdk.NewCoins(sdk.NewCoin(Denom, fundAmount))); err != nil {
			return err
		}

		// Create the validator - DON'T try to manipulate status
		// Just create and let staking module handle everything naturally
		_, err := srv.CreateValidator(ctx, &stakingtypes.MsgCreateValidator{
			Description: stakingtypes.Description{Moniker: moniker},
			Commission: stakingtypes.CommissionRates{
				Rate:          math.LegacyMustNewDecFromStr("0.1"),
				MaxRate:       math.LegacyMustNewDecFromStr("0.1"),
				MaxChangeRate: math.LegacyMustNewDecFromStr("0.1"),
			},
			MinSelfDelegation: minSelfDel,
			ValidatorAddress:  valoperAddr,
			Pubkey:            icsVal.GetPubkey(),
			Value:             sdk.NewCoin(Denom, fundAmount),
		})
		if err != nil {
			return fmt.Errorf("failed to create ICS validator %d: %w", i, err)
		}

		i++
	}

	return nil
}
