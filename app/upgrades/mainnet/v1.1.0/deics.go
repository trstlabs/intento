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
)

const (
	FundAddress = "into1mhd977xqvd8pl7efsrtyltucw0dhf7h4mpv2ve"
	MinStake    = 2_000_000
	Denom       = "uinto"
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

func DeICS(
	ctx sdk.Context,
	stakingKeeper stakingkeeper.Keeper,
	bankKeeper bankkeeper.Keeper,
	readyValopers map[string]bool,
) error {
	_, DAOaddrBz, err := bech32.DecodeAndConvert(FundAddress)
	if err != nil {
		return err
	}
	DAOaddr := sdk.AccAddress(DAOaddrBz)

	validators, err := stakingKeeper.GetAllValidators(ctx)
	if err != nil {
		return err
	}

	for _, val := range validators {
		valoper := val.GetOperator()

		if readyValopers[valoper] {
			// Check balance and fund if needed
			valAddr, err := sdk.ValAddressFromBech32(valoper)
			if err != nil {
				return err
			}
			accAddr := sdk.AccAddress(valAddr)

			balance := bankKeeper.GetBalance(ctx, accAddr, Denom)
			if balance.Amount.LT(math.NewInt(MinStake)) {
				// Fund 2M uinto
				coins := sdk.NewCoins(sdk.NewCoin(Denom, math.NewInt(MinStake)))
				if err := bankKeeper.SendCoins(ctx, DAOaddr, accAddr, coins); err != nil {
					return err
				}
			}
		} else {
			consAddr, err := val.GetConsAddr()
			if err != nil {
				return err
			}

			if err := stakingKeeper.Jail(ctx, consAddr); err != nil {
				return err
			}
		}
	}

	return nil
}
