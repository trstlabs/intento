package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"cosmossdk.io/math"
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
	MaxCap = 5000000000 //5000 tokens from one address on one chain

	TotalIntoAirdropAmountCosmos = 100000000000000 //Cosmos Hub
	TotalIntoAirdropAmountCosm   = 50000000000000  //Cosm is a placeholder for some second Cosmos ecoystem chain to airdrop tokens to
	TotalIntoAirdropAmountSecret = 30000000000000  //Secret may or may not be possible to airdrop tokens to. TODO: check if code works and otherwise remove
)

// ExportAirdropSnapshotCmd generates a snapshot.json from a provided cosmos-sdk v0.36 genesis export.
func ExportAirdropSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-airdrop-snapshot [first-input-snapshot-file] [second-input-snapshot-file] [third-input-snapshot-file] [output-file]",
		Short: "Export a quadratic fairdrop snapshot from provided genesis exports",
		Long: `Export a quadratic fairdrop snapshot from provided  genesis exports
Example:
intentod export-airdrop-snapshot ~/genesisfiles/genesis.cosmoshub-4.json ~/genesisfiles/genesis_secret_3.json  ~/genesisfiles/columbus-5-genesis.json ./snapshot.json
	- Check input genesis:
		file is at ~/.tsrtd/config/genesis.json
	- Snapshot
		file is at "../snapshot.json"
`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			firstGenesisFile := args[0]
			secondGenesisFile := args[1]
			thirdGenesisFile := args[2]
			snapshotOutput := args[3]

			var snapshot Snapshot
			snapshot.Accounts = make(map[string]SnapshotAccount)
			snapshot.TotalIntoAirdropAmount = sdk.ZeroInt()

			snapshot = exportSnapShotFromGenesisFile(clientCtx, firstGenesisFile, "uatom", snapshot)

			var snapshotSecret Snapshot
			snapshotSecret.Accounts = make(map[string]SnapshotAccount)
			snapshotSecret.TotalIntoAirdropAmount = sdk.ZeroInt()

			snapshotSecret = exportSecretSnapShotFromGenesisFile(clientCtx, secondGenesisFile, "uscrt", snapshotSecret)

			var snapshotCosm Snapshot
			snapshotCosm.Accounts = make(map[string]SnapshotAccount)
			snapshotCosm.TotalIntoAirdropAmount = sdk.ZeroInt()

			snapshotCosm = exportSnapShotFromGenesisFile(clientCtx, thirdGenesisFile, "ucosm", snapshotCosm)

			fmt.Printf("Cosmos atom amount %s \n", snapshot.TotalTokenAmount)
			fmt.Printf("Cosmos Into amount %s \n", snapshot.TotalIntoAirdropAmount)
			fmt.Printf("# Cosmos atom accounts %v \n", len(snapshot.Accounts))
			fmt.Printf("Secret scrt amount %s \n", snapshotSecret.TotalTokenAmount)
			fmt.Printf("Secret Into amount %s \n", snapshotSecret.TotalIntoAirdropAmount)
			fmt.Printf("# Secret scrt accounts %v \n", len(snapshotSecret.Accounts))
			fmt.Printf("Cosm luna amount %s \n", snapshotCosm.TotalTokenAmount)
			fmt.Printf("Cosm Into amount %s \n", snapshotCosm.TotalIntoAirdropAmount)
			fmt.Printf("# Cosm luna accounts %v \n", len(snapshotCosm.Accounts))
			snapshot.TotalIntoAirdropAmount = snapshot.TotalIntoAirdropAmount.Add(snapshotSecret.TotalIntoAirdropAmount.Add(snapshotCosm.TotalIntoAirdropAmount))
			snapshot.TotalTokenAmount = snapshot.TotalTokenAmount.Add(snapshotSecret.TotalTokenAmount.Add(snapshotCosm.TotalTokenAmount))

			snapshot.Accounts = removeDuplicatesFromSnapshot(snapshot.Accounts, snapshotSecret.Accounts, snapshotCosm.Accounts)

			fmt.Printf("# total INTO accounts %v \n", len(snapshot.Accounts))
			fmt.Printf("# total INTO coins %s \n", snapshot.TotalIntoAirdropAmount)
			//remove any other duplicates from the lists
			//var snapshotClean = snapshot
			//	snapshotClean.Accounts = removeDuplicates(snapshot.Accounts)

			// export snapshot json
			snapshotJSON, err := json.MarshalIndent(snapshot, "", "    ")
			if err != nil {
				return fmt.Errorf("failed to marshal snapshot: %w", err)
			}

			err = os.WriteFile(snapshotOutput, snapshotJSON, 0644)
			return err
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// compare balance with max cap
func getMin(balance math.LegacyDec) math.LegacyDec {
	if balance.GTE(sdk.NewDec(MaxCap)) {
		atomSqrt, err := sdk.NewDecFromInt(math.NewInt(MaxCap)).ApproxSqrt()
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

func getDenominator(snapshotAccs map[string]SnapshotAccount) math.Int {
	denominator := sdk.ZeroInt()
	for _, acc := range snapshotAccs {
		//add so we ensure suffiient balance
		allTokens := sdk.NewDecFromInt(acc.TokenBalance)
		denominator = denominator.Add(getMin(allTokens).RoundInt())
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

					continue
				}
				var account Account
				if err := codec.NewLegacyAmino().UnmarshalJSON(byteAccounts, &account); err != nil {
					continue
				}

				snapshotAccs[account.Address] = SnapshotAccount{
					TokenAddress:         account.Address,
					TokenBalance:         sdk.ZeroInt(),
					TokenUnstakedBalance: sdk.ZeroInt(),
					TokenStakedBalance:   sdk.ZeroInt(),
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

		acc.TokenBalance = acc.TokenBalance.Add(balance)
		acc.TokenUnstakedBalance = acc.TokenUnstakedBalance.Add(balance)

		snapshotAccs[account.Address] = acc

	}

	for _, unbonding := range stakingGenState.UnbondingDelegations {
		address := unbonding.DelegatorAddress
		acc, ok := snapshotAccs[address]
		if !ok {
			fmt.Printf("No account found for unbonding %s \n", address)
			continue
		}

		unbondingTokens := sdk.NewInt(0)
		for _, entry := range unbonding.Entries {
			unbondingTokens = unbondingTokens.Add(entry.Balance)
		}

		acc.TokenBalance = acc.TokenBalance.Add(unbondingTokens)
		acc.TokenUnstakedBalance = acc.TokenUnstakedBalance.Add(unbondingTokens)

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
		stakedTokens := delegation.Shares.MulInt(val.Tokens).Quo(val.DelegatorShares).RoundInt()

		acc.TokenBalance = acc.TokenBalance.Add(stakedTokens)
		acc.TokenStakedBalance = acc.TokenStakedBalance.Add(stakedTokens)

		snapshotAccs[address] = acc
	}

	denominator := getDenominator(snapshotAccs)
	totalBalance := sdk.ZeroInt()
	totalTokenBalance := sdk.NewInt(0)
	for address, acc := range snapshotAccs {
		allTokens := sdk.NewDecFromInt(acc.TokenBalance)

		allTokenSqrt := getMin(allTokens).RoundInt()

		if denominator.IsZero() {
			acc.TokenOwnershipPercent = sdk.NewDecFromInt(sdk.NewInt(0))
		} else {
			acc.TokenOwnershipPercent = sdk.NewDecFromInt(allTokenSqrt).QuoInt(denominator)
		}

		if allTokens.IsZero() {
			acc.TokenStakedPercent = sdk.ZeroDec()
			acc.IntoBalance = sdk.ZeroInt()
			snapshotAccs[address] = acc
			continue
		}

		stakedTokens := sdk.NewDecFromInt(acc.TokenStakedBalance)
		stakedPercent := stakedTokens.Quo(allTokens)

		acc.TokenStakedPercent = stakedPercent
		acc.IntoBalance = acc.TokenOwnershipPercent.MulInt(sdk.NewInt(TotalIntoAirdropAmountCosmos)).RoundInt()

		totalBalance = totalBalance.Add(acc.IntoBalance)
		snapshotAccount, ok := snapshot.Accounts[address]
		if !ok {
			snapshot.Accounts[address] = acc
			totalTokenBalance = totalTokenBalance.Add(acc.TokenBalance)
		} else {
			if snapshotAccount.IntoBalance.IsNil() {
				snapshotAccount.IntoBalance = sdk.ZeroInt()
			}
			snapshotAccount.IntoBalance = snapshotAccount.IntoBalance.Add(acc.IntoBalance)
			snapshotAccount.TokenBalance = snapshotAccount.TokenBalance.Add(acc.TokenBalance)
			snapshotAccount.TokenUnstakedBalance = snapshotAccount.TokenUnstakedBalance.Add(acc.TokenUnstakedBalance)
			snapshot.Accounts[address] = snapshotAccount

			totalTokenBalance = totalTokenBalance.Add(acc.TokenBalance)
		}
	}
	snapshot.TotalTokenAmount = totalTokenBalance
	snapshot.TotalIntoAirdropAmount = snapshot.TotalIntoAirdropAmount.Add(totalBalance)
	snapshot.NumberAccounts = snapshot.NumberAccounts + uint64(len(snapshot.Accounts))

	fmt.Printf("Complete read genesis file %s \n", genesisFile)
	fmt.Printf("# atom accounts: %d\n", len(snapshotAccs))
	fmt.Printf("# atom accounts in snapshot: %d\n", len(snapshot.Accounts))
	fmt.Printf("Token TotalSupply: %s\n", totalTokenBalance.String())
	fmt.Printf("Into TotalSupply: %s\n", totalBalance.String())
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
			if err != nil {
				fmt.Printf("coins failed to marshal: %s\n", err)
			}
			err = json.Unmarshal(coins, &secretAcc.Value.Coins)
			if err != nil {
				fmt.Printf("Secret auth module failed to unmarshal: %s\n", err)
			}

			balance := secretAcc.Value.Coins.AmountOf(denom)

			snapshotAccs[account.Value.Address] = SnapshotAccount{
				TokenAddress:         account.Value.Address,
				TokenBalance:         balance,
				TokenUnstakedBalance: balance,
				TokenStakedBalance:   sdk.ZeroInt(),
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

		unbondingTokens := sdk.NewInt(0)
		for _, entry := range unbonding.Entries {
			intb, _ := sdk.NewIntFromString(entry.Balance)
			unbondingTokens = unbondingTokens.Add(intb)
		}

		acc.TokenBalance = acc.TokenBalance.Add(unbondingTokens)
		acc.TokenUnstakedBalance = acc.TokenUnstakedBalance.Add(unbondingTokens)

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
		shares, _ := sdk.NewDecFromStr(delegation.Shares)
		tokens, _ := sdk.NewIntFromString(val.Tokens)
		valShares, _ := sdk.NewDecFromStr(val.DelegatorShares)
		stakedTokens := shares.MulInt(tokens).Quo(valShares).RoundInt()

		acc.TokenBalance = acc.TokenBalance.Add(stakedTokens)
		acc.TokenStakedBalance = acc.TokenStakedBalance.Add(stakedTokens)
		snapshotAccs[address] = acc
	}

	denominator := getDenominator(snapshotAccs)
	totalBalance := sdk.ZeroInt()
	totalTokenBalance := sdk.NewInt(0)
	for address, acc := range snapshotAccs {
		allTokens := sdk.NewDecFromInt(acc.TokenBalance)

		allTokenSqrt := getMin(allTokens).RoundInt()

		if denominator.IsZero() {
			acc.TokenOwnershipPercent = sdk.NewDecFromInt(sdk.NewInt(0))
		} else {
			acc.TokenOwnershipPercent = sdk.NewDecFromInt(allTokenSqrt).QuoInt(denominator)
		}

		if allTokens.IsZero() {
			acc.TokenStakedPercent = sdk.ZeroDec()
			acc.IntoBalance = sdk.ZeroInt()
			snapshotAccs[address] = acc
			continue
		}

		stakedTokens := sdk.NewDecFromInt(acc.TokenStakedBalance)
		stakedPercent := stakedTokens.Quo(allTokens)

		acc.TokenStakedPercent = stakedPercent
		acc.IntoBalance = acc.TokenOwnershipPercent.MulInt(sdk.NewInt(TotalIntoAirdropAmountSecret)).RoundInt()

		totalBalance = totalBalance.Add(acc.IntoBalance)
		snapshotAccount, ok := snapshot.Accounts[address]
		if !ok {
			snapshot.Accounts[address] = acc
			totalTokenBalance = totalTokenBalance.Add(acc.TokenBalance)
		} else {
			if snapshotAccount.IntoBalance.IsNil() {
				snapshotAccount.IntoBalance = sdk.ZeroInt()
			}
			snapshotAccount.IntoBalance = snapshotAccount.IntoBalance.Add(acc.IntoBalance)
			snapshotAccount.TokenBalance = snapshotAccount.TokenBalance.Add(acc.TokenBalance)
			snapshotAccount.TokenUnstakedBalance = snapshotAccount.TokenUnstakedBalance.Add(acc.TokenUnstakedBalance)
			snapshot.Accounts[address] = snapshotAccount

			totalTokenBalance = totalTokenBalance.Add(acc.TokenBalance)
		}
	}
	snapshot.TotalTokenAmount = totalTokenBalance
	snapshot.TotalIntoAirdropAmount = snapshot.TotalIntoAirdropAmount.Add(totalBalance)
	snapshot.NumberAccounts = snapshot.NumberAccounts + uint64(len(snapshot.Accounts))

	fmt.Printf("Complete read genesis file %s \n", genesisFile)
	fmt.Printf("# scrt accounts: %d\n", len(snapshotAccs))
	fmt.Printf("# scrt accounts in snapshot: %d\n", len(snapshot.Accounts))
	fmt.Printf("Scrt TotalSupply: %s\n", totalTokenBalance.String())
	fmt.Printf("Into TotalSupply: %s\n", totalBalance.String())

	return snapshot
}

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
	TotalTokenAmount       math.Int `json:"total_atom_amount"`
	TotalIntoAirdropAmount math.Int `json:"total_into_amount"`
	NumberAccounts         uint64   `json:"num_accounts"`

	Accounts map[string]SnapshotAccount `json:"accounts"`
}

// SnapshotAccount provide fields of snapshot per account
type SnapshotAccount struct {
	TokenAddress string `json:"atom_address"` // Token Balance = TokenStakedBalance + TokenUnstakedBalance

	TokenBalance          math.Int       `json:"atom_balance"`
	TokenOwnershipPercent math.LegacyDec `json:"atom_ownership_percent"`

	TokenStakedBalance   math.Int       `json:"atom_staked_balance"`
	TokenUnstakedBalance math.Int       `json:"atom_unstaked_balance"` // TokenStakedPercent = TokenStakedBalance / TokenBalance
	TokenStakedPercent   math.LegacyDec `json:"atom_staked_percent"`

	IntoBalance math.Int `json:"into_balance"`
	Denominator math.Int `json:"denominator"`
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
	Delegations         []Delegation `json:"delegations"`
	Exported            bool         `json:"exported"`
	LastTotalPower      string       `json:"last_total_power"`
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
	UnbondingDelegations []UnbondingDelegations `json:"unbonding_delegations"`
	Validators           []Validator            `json:"validators"`
}

type UnbondingDelegations struct {
	DelegatorAddress string `json:"delegator_address"`
	Entries          []struct {
		Balance        string    `json:"balance"`
		CompletionTime time.Time `json:"completion_time"`
		CreationHeight string    `json:"creation_height"`
		InitialBalance string    `json:"initial_balance"`
	} `json:"entries"`
	ValidatorAddress string `json:"validator_address"`
}
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

type Delegation struct {
	DelegatorAddress string `json:"delegator_address"`
	Shares           string `json:"shares"`
	ValidatorAddress string `json:"validator_address"`
}

func removeDuplicatesFromSnapshot(first map[string]SnapshotAccount, second map[string]SnapshotAccount, third map[string]SnapshotAccount) map[string]SnapshotAccount {
	list := make([]string, 0, len(first))

	var duplicatesSecret uint64
	var duplicatesCosm uint64
	fmt.Printf("length cosmos accs %d\n", len(first))
	fmt.Printf("length secret accs %d\n", len(second))
	fmt.Printf("length terra accs %d\n", len(third))
	//append accounts to list so we have list of all acc
	for acc, _ := range first {

		list = append(list, acc)

	}
	fmt.Println("appended cosmos accs")
	//append accounts to list so we have list of all secret acc
	for acc, value := range second {
		newAcc, _ := cosmosConvertBech32(acc)
		//if it does not contain the acc we append to it
		if !contains(list, newAcc) {

			list = append(list, newAcc)
			first[newAcc] = value

		} else {
			//	fmt.Println("delete secret acc  %s\n", acc)
			var acct = first[newAcc]
			acct.IntoBalance = acct.IntoBalance.Add(value.IntoBalance)
			delete(second, acc)
			first[newAcc] = acct
			fmt.Println("merged scrt acc", acc)
			duplicatesSecret = duplicatesSecret + 1
		}
	}
	fmt.Println("removed duplicates secret")
	//append accounts to list so we have list of all terra acc
	for acc, value := range third {
		newAcc, _ := cosmosConvertBech32(acc)
		//if it does not contain the acc we append to it
		if !contains(list, newAcc) {

			list = append(list, newAcc)
			first[newAcc] = value
		} else {
			//	fmt.Println("delete secret acc  %s\n", acc)
			var acct = first[newAcc]
			acct.IntoBalance = acct.IntoBalance.Add(value.IntoBalance)
			first[newAcc] = acct

			fmt.Println("merged terra acc", acc)
			delete(third, acc)
			duplicatesCosm = duplicatesCosm + 1
		}
	}

	fmt.Printf("final length secret accs: %d\n", len(second))
	fmt.Printf("final length terra accs: %d\n", len(third))
	fmt.Printf("length list accs: %d\n", len(list))

	fmt.Printf("duplicates secret accs: %d\n", duplicatesSecret)
	fmt.Printf("duplicates terra accs: %d\n", duplicatesCosm)
	//all := append(first, second)

	fmt.Println("removed duplicates")
	return first
}
