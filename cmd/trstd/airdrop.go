package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
)

const (
	MaxCap = 50000000000

	TotalTrstAirdropAmount = 20000000000000 // 0.5% * 200000000
)

type Validator struct {
	Commission struct {
		CommissionRates struct {
			MaxChangeRate string `json:"max_change_rate"`
			MaxRate       string `json:"max_rate"`
			Rate          string `json:"rate"`
		} `json:"commission_rates"`
		UpdateTime time.Time `json:"update_time"`
	} `json:"commission"`
	ConsensusPubkey string `json:"consensus_pubkey"`
	DelegatorShares string `json:"delegator_shares"`
	Description     struct {
		Details         string `json:"details"`
		Identity        string `json:"identity"`
		Moniker         string `json:"moniker"`
		SecurityContact string `json:"security_contact"`
		Website         string `json:"website"`
	} `json:"description"`
	Jailed            bool      `json:"jailed"`
	MinSelfDelegation string    `json:"min_self_delegation"`
	OperatorAddress   string    `json:"operator_address"`
	Status            int       `json:"status"`
	Tokens            string    `json:"tokens"`
	UnbondingHeight   string    `json:"unbonding_height"`
	UnbondingTime     time.Time `json:"unbonding_time"`
}

type GenesisFile struct {
	AppHash  string `json:"app_hash"`
	AppState struct {
		Auth struct {
			Accounts []struct {
				Type  string `json:"type"`
				Value struct {
					AccountNumber int    `json:"account_number"`
					Address       string `json:"address"`
					Coins         []byte `json:"coins"`
					PublicKey     string `json:"public_key"`
					Sequence      int    `json:"sequence"`
				} `json:"value"`
			} `json:"accounts"`
			Params struct {
				MaxMemoCharacters      string `json:"max_memo_characters"`
				SigVerifyCostEd25519   string `json:"sig_verify_cost_ed25519"`
				SigVerifyCostSecp256K1 string `json:"sig_verify_cost_secp256k1"`
				TxSigLimit             string `json:"tx_sig_limit"`
				TxSizeCostPerByte      string `json:"tx_size_cost_per_byte"`
			} `json:"params"`
		} `json:"auth"`
		Staking struct {
			Delegations []struct {
				DelegatorAddress string `json:"delegator_address"`
				Shares           string `json:"shares"`
				ValidatorAddress string `json:"validator_address"`
			} `json:"delegations"`
			Exported            bool   `json:"exported"`
			LastTotalPower      string `json:"last_total_power"`
			LastValidatorPowers []struct {
				Address string `json:"Address"`
				Power   string `json:"Power"`
			} `json:"last_validator_powers"`
			Params struct {
				BondDenom         string `json:"bond_denom"`
				HistoricalEntries int    `json:"historical_entries"`
				MaxEntries        int    `json:"max_entries"`
				MaxValidators     int    `json:"max_validators"`
				UnbondingTime     string `json:"unbonding_time"`
			} `json:"params"`
			Redelegations []struct {
				DelegatorAddress string `json:"delegator_address"`
				Entries          []struct {
					CompletionTime time.Time `json:"completion_time"`
					CreationHeight string    `json:"creation_height"`
					InitialBalance string    `json:"initial_balance"`
					SharesDst      string    `json:"shares_dst"`
				} `json:"entries"`
				ValidatorDstAddress string `json:"validator_dst_address"`
				ValidatorSrcAddress string `json:"validator_src_address"`
			} `json:"redelegations"`
			UnbondingDelegations []struct {
				DelegatorAddress string `json:"delegator_address"`
				Entries          []struct {
					Balance        string    `json:"balance"`
					CompletionTime time.Time `json:"completion_time"`
					CreationHeight string    `json:"creation_height"`
					InitialBalance string    `json:"initial_balance"`
				} `json:"entries"`
				ValidatorAddress string `json:"validator_address"`
			} `json:"unbonding_delegations"`
			Validators []Validator `json:"validators"`
		} `json:"staking"`
	} `json:"app_state"`
	ChainID         string `json:"chain_id"`
	ConsensusParams struct {
		Block struct {
			MaxBytes   string `json:"max_bytes"`
			MaxGas     string `json:"max_gas"`
			TimeIotaMs string `json:"time_iota_ms"`
		} `json:"block"`
		Evidence struct {
			MaxAgeDuration  string `json:"max_age_duration"`
			MaxAgeNumBlocks string `json:"max_age_num_blocks"`
		} `json:"evidence"`
		Validator struct {
			PubKeyTypes []string `json:"pub_key_types"`
		} `json:"validator"`
	} `json:"consensus_params"`
	GenesisTime time.Time `json:"genesis_time"`
	Validators  []struct {
		Address string `json:"address"`
		Name    string `json:"name"`
		Power   string `json:"power"`
		PubKey  struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"pub_key"`
	} `json:"validators"`
}

type Snapshot struct {
	TotalAtomAmount        sdk.Int `json:"total_atom_amount"`
	TotalTrstAirdropAmount sdk.Int `json:"total_trst_amount"`
	NumberAccounts         uint64  `json:"num_accounts"`

	Accounts map[string]SnapshotAccount `json:"accounts"`
}

// SnapshotAccount provide fields of snapshot per account
type SnapshotAccount struct {
	AtomAddress string `json:"atom_address"` // Atom Balance = AtomStakedBalance + AtomUnstakedBalance

	AtomBalance          sdk.Int `json:"atom_balance"`
	AtomOwnershipPercent sdk.Dec `json:"atom_ownership_percent"`

	AtomStakedBalance   sdk.Int `json:"atom_staked_balance"`
	AtomUnstakedBalance sdk.Int `json:"atom_unstaked_balance"` // AtomStakedPercent = AtomStakedBalance / AtomBalance
	AtomStakedPercent   sdk.Dec `json:"atom_staked_percent"`

	TrstBalance sdk.Int `json:"trst_balance"`
	Denominator sdk.Int `json:"denominator"`
}

type Account struct {
	Address       string `json:"address,omitempty"`
	AccountNumber uint64 `json:"account_number,omitempty"`
	Sequence      uint64 `json:"sequence,omitempty"`
}

type SecretAccount struct {
	Type  string `json:"type"`
	Value struct {
		AccountNumber int       `json:"account_number"`
		Address       string    `json:"address"`
		Coins         sdk.Coins `json:"coins"`

		PublicKey string `json:"public_key"`
		Sequence  int    `json:"sequence"`
	} `json:"value"`
}

type Staking struct {
	Delegations []struct {
		DelegatorAddress string `json:"delegator_address"`
		Shares           string `json:"shares"`
		ValidatorAddress string `json:"validator_address"`
	} `json:"delegations"`
	Exported            bool   `json:"exported"`
	LastTotalPower      string `json:"last_total_power"`
	LastValidatorPowers []struct {
		Address string `json:"Address"`
		Power   string `json:"Power"`
	} `json:"last_validator_powers"`
	Params struct {
		BondDenom         string `json:"bond_denom"`
		HistoricalEntries int    `json:"historical_entries"`
		MaxEntries        int    `json:"max_entries"`
		MaxValidators     int    `json:"max_validators"`
		UnbondingTime     string `json:"unbonding_time"`
	} `json:"params"`
	Redelegations []struct {
		DelegatorAddress string `json:"delegator_address"`
		Entries          []struct {
			CompletionTime time.Time `json:"completion_time"`
			CreationHeight string    `json:"creation_height"`
			InitialBalance string    `json:"initial_balance"`
			SharesDst      string    `json:"shares_dst"`
		} `json:"entries"`
		ValidatorDstAddress string `json:"validator_dst_address"`
		ValidatorSrcAddress string `json:"validator_src_address"`
	} `json:"redelegations"`
	UnbondingDelegations []struct {
		DelegatorAddress string `json:"delegator_address"`
		Entries          []struct {
			Balance        string    `json:"balance"`
			CompletionTime time.Time `json:"completion_time"`
			CreationHeight string    `json:"creation_height"`
			InitialBalance string    `json:"initial_balance"`
		} `json:"entries"`
		ValidatorAddress string `json:"validator_address"`
	} `json:"unbonding_delegations"`
	Validators []Validator `json:"validators"`
}

/*
type Auth struct {
	Accounts []SecretAccount `json:"accounts"`
	Params   struct {
		MaxMemoCharacters      string `json:"max_memo_characters"`
		SigVerifyCostEd25519   string `json:"sig_verify_cost_ed25519"`
		SigVerifyCostSecp256K1 string `json:"sig_verify_cost_secp256k1"`
		TxSigLimit             string `json:"tx_sig_limit"`
		TxSizeCostPerByte      string `json:"tx_size_cost_per_byte"`
	} `json:"params"`
}*/
type Auth struct {
	Accounts []struct {
		Type  string `json:"type"`
		Value struct {
			AccountNumber int    `json:"account_number"`
			Address       string `json:"address"`
			Coins         []struct {
				Amount string `json:"amount"`
				Denom  string `json:"denom"`
			} `json:"coins"`
			PublicKey string `json:"public_key"`
			Sequence  int    `json:"sequence"`
		} `json:"value"`
	} `json:"accounts"`
	Params struct {
		MaxMemoCharacters      string `json:"max_memo_characters"`
		SigVerifyCostEd25519   string `json:"sig_verify_cost_ed25519"`
		SigVerifyCostSecp256K1 string `json:"sig_verify_cost_secp256k1"`
		TxSigLimit             string `json:"tx_sig_limit"`
		TxSizeCostPerByte      string `json:"tx_size_cost_per_byte"`
	} `json:"params"`
}

// ExportAirdropSnapshotCmd generates a snapshot.json from a provided cosmos-sdk v0.36 genesis export.
func ExportAirdropSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-airdrop-snapshot [first-input-snapshot-file] [second-input-snapshot-file] [output-file]",
		Short: "Export a quadratic fairdrop snapshot from a provided cosmos-sdk  genesis export",
		Long: `Export a quadratic fairdrop snapshot from a provided cosmos-sdk genesis export
Example:
trstd export-airdrop-snapshot ~/genesisfiles/genesis.cosmoshub-4.json ~/genesisfiles/genesis_secret_3.json  ./snapshot.json
	- Check input genesis:
		file is at ~/.tsrtd/config/genesis.json
	- Snapshot
		file is at "../snapshot.json"
`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			firstGenesisFile := args[0]
			secondGenesisFile := args[1]
			snapshotOutput := args[2]

			var snapshot Snapshot
			snapshot.Accounts = make(map[string]SnapshotAccount)
			snapshot.TotalTrstAirdropAmount = sdk.ZeroInt()

			snapshot = exportSnapShotFromGenesisFile(clientCtx, firstGenesisFile, "uatom", snapshot)

			var snapshotSecret Snapshot
			snapshotSecret.Accounts = make(map[string]SnapshotAccount)
			snapshotSecret.TotalTrstAirdropAmount = sdk.ZeroInt()

			snapshotSecret = exportSecretSnapShotFromGenesisFile(clientCtx, secondGenesisFile, "uscrt", snapshotSecret)

			fmt.Printf("atom amount %s \n", snapshot.TotalAtomAmount)
			fmt.Printf("scrt amount %s \n", snapshotSecret.TotalAtomAmount)
			snapshot.TotalTrstAirdropAmount = snapshot.TotalTrstAirdropAmount.Add(snapshotSecret.TotalTrstAirdropAmount)
			snapshot.TotalAtomAmount = snapshot.TotalAtomAmount.Add(snapshotSecret.TotalAtomAmount)

			snapshot.Accounts = removeDuplicatesFromSnapshot(snapshot.Accounts, snapshotSecret.Accounts)

			fmt.Println("cleaning list")
			//remove any other duplicates from the lists
			//var snapshotClean = snapshot
			//	snapshotClean.Accounts = removeDuplicates(snapshot.Accounts)

			// export snapshot json
			snapshotJSON, err := json.MarshalIndent(snapshot, "", "    ")
			if err != nil {
				return fmt.Errorf("failed to marshal snapshot: %w", err)
			}

			err = ioutil.WriteFile(snapshotOutput, snapshotJSON, 0644)
			return err
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// compare balance with max cap
func getMin(balance sdk.Dec) sdk.Dec {
	if balance.GTE(sdk.NewDec(MaxCap)) {
		atomSqrt, err := sdk.NewInt(MaxCap).ToDec().ApproxSqrt()
		if err != nil {
			panic(fmt.Sprintf("failed to root atom balance: %s", err))
		}
		return atomSqrt
	} else {
		atomSqrt, err := balance.ApproxSqrt()
		if err != nil {
			panic(fmt.Sprintf("failed to root atom balance: %s", err))
		}
		return atomSqrt
	}
}

func getDenominator(snapshotAccs map[string]SnapshotAccount) sdk.Int {
	denominator := sdk.ZeroInt()
	for _, acc := range snapshotAccs {
		allAtoms := acc.AtomBalance.ToDec()
		denominator = denominator.Add(getMin(allAtoms).RoundInt())
	}
	return denominator
}
func exportSnapShotFromGenesisFile(clientCtx client.Context, genesisFile string, denom string, snapshot Snapshot) Snapshot {
	appState, _, _ := genutiltypes.GenesisStateFromGenFile(genesisFile)
	bankGenState := banktypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
	stakingGenState := stakingtypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
	authGenState := authtypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)

	snapshotAccs := make(map[string]SnapshotAccount)
	for _, account := range authGenState.GetAccounts() {

		if account.TypeUrl == "/cosmos.auth.v1beta1.BaseAccount" {
			_, ok := account.GetCachedValue().(authtypes.GenesisAccount)
			if ok {
				var byteAccounts []byte
				// Reason is prefix is tsrt --> getAddress will be empty
				// Marshal construct and convert to new struct to get address
				byteAccounts, err := codec.NewLegacyAmino().MarshalJSON(account.GetCachedValue())
				if err != nil {
					fmt.Printf("No account found for bank balance %s \n", string(byteAccounts))
					fmt.Println(err.Error())
					continue
				}
				var account Account
				if err := codec.NewLegacyAmino().UnmarshalJSON(byteAccounts, &account); err != nil {
					continue
				}

				snapshotAccs[account.Address] = SnapshotAccount{
					AtomAddress:         account.Address,
					AtomBalance:         sdk.ZeroInt(),
					AtomUnstakedBalance: sdk.ZeroInt(),
					AtomStakedBalance:   sdk.ZeroInt(),
				}
			}
		}
	}

	// Produce the map of address to total atom balance, both staked and unstaked

	for _, account := range bankGenState.Balances {

		acc, ok := snapshotAccs[account.Address]
		if !ok {
			fmt.Printf("No account found for bank balance %s \n", account.Address)
			continue
		}
		balance := account.Coins.AmountOf(denom)

		acc.AtomBalance = acc.AtomBalance.Add(balance)
		acc.AtomUnstakedBalance = acc.AtomUnstakedBalance.Add(balance)

		snapshotAccs[account.Address] = acc

	}

	for _, unbonding := range stakingGenState.UnbondingDelegations {
		address := unbonding.DelegatorAddress
		acc, ok := snapshotAccs[address]
		if !ok {
			fmt.Printf("No account found for unbonding %s \n", address)
			continue
		}

		unbondingAtoms := sdk.NewInt(0)
		for _, entry := range unbonding.Entries {
			unbondingAtoms = unbondingAtoms.Add(entry.Balance)
		}

		acc.AtomBalance = acc.AtomBalance.Add(unbondingAtoms)
		acc.AtomUnstakedBalance = acc.AtomUnstakedBalance.Add(unbondingAtoms)

		snapshotAccs[address] = acc
	}

	// Make a map from validator operator address to the v42 validator type
	validators := make(map[string]stakingtypes.Validator)
	for _, validator := range stakingGenState.Validators {
		validators[validator.OperatorAddress] = validator
	}

	for _, delegation := range stakingGenState.Delegations {
		address := delegation.DelegatorAddress

		acc, ok := snapshotAccs[address]
		if !ok {
			fmt.Printf("No account found for delegation address %s \n", address)
			continue
		}

		val := validators[delegation.ValidatorAddress]
		stakedAtoms := delegation.Shares.MulInt(val.Tokens).Quo(val.DelegatorShares).RoundInt()

		acc.AtomBalance = acc.AtomBalance.Add(stakedAtoms)
		acc.AtomStakedBalance = acc.AtomStakedBalance.Add(stakedAtoms)

		snapshotAccs[address] = acc
	}

	denominator := getDenominator(snapshotAccs)
	totalBalance := sdk.ZeroInt()
	totalAtomBalance := sdk.NewInt(0)
	for address, acc := range snapshotAccs {
		allAtoms := acc.AtomBalance.ToDec()

		allAtomSqrt := getMin(allAtoms).RoundInt()

		if denominator.IsZero() {
			acc.AtomOwnershipPercent = sdk.NewInt(0).ToDec()
		} else {
			acc.AtomOwnershipPercent = allAtomSqrt.ToDec().QuoInt(denominator)
		}

		if allAtoms.IsZero() {
			acc.AtomStakedPercent = sdk.ZeroDec()
			acc.TrstBalance = sdk.ZeroInt()
			snapshotAccs[address] = acc
			continue
		}

		stakedAtoms := acc.AtomStakedBalance.ToDec()
		stakedPercent := stakedAtoms.Quo(allAtoms)

		acc.AtomStakedPercent = stakedPercent
		acc.TrstBalance = acc.AtomOwnershipPercent.MulInt(sdk.NewInt(TotalTrstAirdropAmount)).RoundInt()

		totalBalance = totalBalance.Add(acc.TrstBalance)
		snapshotAccount, ok := snapshot.Accounts[address]
		if !ok {
			snapshot.Accounts[address] = acc
			totalAtomBalance = totalAtomBalance.Add(acc.AtomBalance)
		} else {
			if snapshotAccount.TrstBalance.IsNil() {
				snapshotAccount.TrstBalance = sdk.ZeroInt()
			}
			snapshotAccount.TrstBalance = snapshotAccount.TrstBalance.Add(acc.TrstBalance)
			snapshotAccount.AtomBalance = snapshotAccount.AtomBalance.Add(acc.AtomBalance)
			snapshotAccount.AtomUnstakedBalance = snapshotAccount.AtomUnstakedBalance.Add(acc.AtomUnstakedBalance)
			snapshot.Accounts[address] = snapshotAccount

			totalAtomBalance = totalAtomBalance.Add(acc.AtomBalance)
		}
	}
	snapshot.TotalAtomAmount = totalAtomBalance
	snapshot.TotalTrstAirdropAmount = snapshot.TotalTrstAirdropAmount.Add(totalBalance)
	snapshot.NumberAccounts = snapshot.NumberAccounts + uint64(len(snapshot.Accounts))

	fmt.Printf("Complete read genesis file %s \n", genesisFile)
	fmt.Printf("# accounts: %d\n", len(snapshotAccs))
	fmt.Printf("atomTotalSupply: %s\n", totalAtomBalance.String())
	fmt.Printf("trstTotalSupply: %s\n", totalBalance.String())
	return snapshot
}

func exportSecretSnapShotFromGenesisFile(clientCtx client.Context, genesisFile string, denom string, snapshot Snapshot) Snapshot {
	appState, _, _ := genutiltypes.GenesisStateFromGenFile(genesisFile)

	var authGenState = appState["auth"]
	var auth Auth
	err := json.Unmarshal(authGenState, &auth)
	if err != nil {
		fmt.Printf("auth module failed to unmarshal: %s\n", err)
	}

	var stakingGenState = appState["staking"]
	var staking Staking
	err = json.Unmarshal(stakingGenState, &staking)
	if err != nil {
		fmt.Printf("staking module failed to unmarshal: %s\n", err)
	}

	snapshotAccs := make(map[string]SnapshotAccount)
	for _, account := range auth.Accounts {

		if account.Type == "cosmos-sdk/Account" {

			var secretAcc SecretAccount
			coins, err := json.Marshal(account.Value.Coins)
			err = json.Unmarshal(coins, &secretAcc.Value.Coins)
			if err != nil {
				fmt.Printf("Secret auth module failed to unmarshal: %s\n", err)
			}

			balance := secretAcc.Value.Coins.AmountOf(denom)

			snapshotAccs[account.Value.Address] = SnapshotAccount{
				AtomAddress:         account.Value.Address,
				AtomBalance:         balance,
				AtomUnstakedBalance: balance,
				AtomStakedBalance:   sdk.ZeroInt(),
			}

		}
	}

	for _, unbonding := range staking.UnbondingDelegations {
		address := unbonding.DelegatorAddress
		acc, ok := snapshotAccs[address]
		if !ok {
			fmt.Printf("No account found for unbonding %s \n", address)
			continue
		}

		unbondingAtoms := sdk.NewInt(0)
		for _, entry := range unbonding.Entries {
			intb, _ := sdk.NewIntFromString(entry.Balance)
			unbondingAtoms = unbondingAtoms.Add(intb)
		}

		acc.AtomBalance = acc.AtomBalance.Add(unbondingAtoms)
		acc.AtomUnstakedBalance = acc.AtomUnstakedBalance.Add(unbondingAtoms)

		snapshotAccs[address] = acc
	}

	// Make a map from validator operator address to the v42 validator type

	validators := make(map[string]Validator)
	for _, validator := range staking.Validators {
		//	add, _ := sdk.ValAddressFromBech32(validator.Address)

		//	val, _ := stakingtypes.NewValidator(add, validator.PubKey, validator.Name)
		validators[validator.OperatorAddress] = validator
	}

	for _, delegation := range staking.Delegations {
		address := delegation.DelegatorAddress

		acc, ok := snapshotAccs[address]
		if !ok {
			fmt.Printf("No account found for delegation address %s \n", address)
			continue
		}

		val := validators[delegation.ValidatorAddress]
		share, _ := strconv.ParseInt(delegation.Shares, 10, 64)
		valShare, _ := strconv.ParseInt(val.DelegatorShares, 10, 64)

		token, _ := sdk.NewIntFromString(val.Tokens)

		if valShare != 0 && val.Tokens != "0" {
			stakedAtoms := sdk.NewDec(share).MulInt(token).Quo(sdk.NewDec(valShare)).RoundInt()

			acc.AtomBalance = acc.AtomBalance.Add(stakedAtoms)
			acc.AtomStakedBalance = acc.AtomStakedBalance.Add(stakedAtoms)
		} else {
			//fmt.Printf("Zero val token balance for %s \n", address)
		}

		snapshotAccs[address] = acc
	}

	denominator := getDenominator(snapshotAccs)
	totalBalance := sdk.ZeroInt()
	totalAtomBalance := sdk.NewInt(0)
	for address, acc := range snapshotAccs {
		allAtoms := acc.AtomBalance.ToDec()

		allAtomSqrt := getMin(allAtoms).RoundInt()

		if denominator.IsZero() {
			acc.AtomOwnershipPercent = sdk.NewInt(0).ToDec()
		} else {
			acc.AtomOwnershipPercent = allAtomSqrt.ToDec().QuoInt(denominator)
		}

		if allAtoms.IsZero() {
			acc.AtomStakedPercent = sdk.ZeroDec()
			acc.TrstBalance = sdk.ZeroInt()
			snapshotAccs[address] = acc
			continue
		}

		stakedAtoms := acc.AtomStakedBalance.ToDec()
		stakedPercent := stakedAtoms.Quo(allAtoms)

		acc.AtomStakedPercent = stakedPercent
		acc.TrstBalance = acc.AtomOwnershipPercent.MulInt(sdk.NewInt(TotalTrstAirdropAmount)).RoundInt()

		totalBalance = totalBalance.Add(acc.TrstBalance)
		snapshotAccount, ok := snapshot.Accounts[address]
		if !ok {
			snapshot.Accounts[address] = acc
			totalAtomBalance = totalAtomBalance.Add(acc.AtomBalance)
		} else {
			if snapshotAccount.TrstBalance.IsNil() {
				snapshotAccount.TrstBalance = sdk.ZeroInt()
			}
			snapshotAccount.TrstBalance = snapshotAccount.TrstBalance.Add(acc.TrstBalance)
			snapshotAccount.AtomBalance = snapshotAccount.AtomBalance.Add(acc.AtomBalance)
			snapshotAccount.AtomUnstakedBalance = snapshotAccount.AtomUnstakedBalance.Add(acc.AtomUnstakedBalance)
			snapshot.Accounts[address] = snapshotAccount

			totalAtomBalance = totalAtomBalance.Add(acc.AtomBalance)
		}
	}
	snapshot.TotalAtomAmount = totalAtomBalance
	snapshot.TotalTrstAirdropAmount = snapshot.TotalTrstAirdropAmount.Add(totalBalance)
	snapshot.NumberAccounts = snapshot.NumberAccounts + uint64(len(snapshot.Accounts))

	fmt.Printf("Complete read genesis file %s \n", genesisFile)
	fmt.Printf("# accounts: %d\n", len(snapshotAccs))
	fmt.Printf("atomTotalSupply: %s\n", totalAtomBalance.String())
	fmt.Printf("trstTotalSupply: %s\n", totalBalance.String())

	return snapshot
}
