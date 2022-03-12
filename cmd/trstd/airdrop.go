package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	MaxCap = 5000000000 //5000 ATOM/LUNA/SCRT

	TotalTrstAirdropAmountCosmos = 100000000000000 //We want to support cosmonauts and we love Cosmos <3
	TotalTrstAirdropAmountTerra  = 50000000000000  //Terra is a huge ecosystem and growing, we would love to see developers use Trustless Contracts
	TotalTrstAirdropAmountSecret = 30000000000000  //Without Secret's development this project would not be where it is, this total amount is lower as there are significantly fewer accounts in the snapshot, so percentage-wise these accounts should be well off nontheless
)

// ExportAirdropSnapshotCmd generates a snapshot.json from a provided cosmos-sdk v0.36 genesis export.
func ExportAirdropSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-airdrop-snapshot [first-input-snapshot-file] [second-input-snapshot-file] [third-input-snapshot-file] [output-file]",
		Short: "Export a quadratic fairdrop snapshot from provided genesis exports",
		Long: `Export a quadratic fairdrop snapshot from provided  genesis exports
Example:
trstd export-airdrop-snapshot ~/genesisfiles/genesis.cosmoshub-4.json ~/genesisfiles/genesis_secret_3.json  ~/genesisfiles/columbus-5-genesis.json ./snapshot.json
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
			snapshot.TotalTrstAirdropAmount = sdk.ZeroInt()

			snapshot = exportSnapShotFromGenesisFile(clientCtx, firstGenesisFile, "uatom", snapshot)

			var snapshotSecret Snapshot
			snapshotSecret.Accounts = make(map[string]SnapshotAccount)
			snapshotSecret.TotalTrstAirdropAmount = sdk.ZeroInt()

			snapshotSecret = exportSecretSnapShotFromGenesisFile(clientCtx, secondGenesisFile, "uscrt", snapshotSecret)

			var snapshotTerra Snapshot
			snapshotTerra.Accounts = make(map[string]SnapshotAccount)
			snapshotTerra.TotalTrstAirdropAmount = sdk.ZeroInt()

			snapshotTerra = exportTerraSnapShotFromGenesisFile(clientCtx, thirdGenesisFile, "uluna", snapshotTerra)

			fmt.Printf("Cosmos atom amount %s \n", snapshot.TotalTokenAmount)
			fmt.Printf("Cosmos Trst amount %s \n", snapshot.TotalTrstAirdropAmount)
			fmt.Printf("# Cosmos atom accounts %s \n", len(snapshot.Accounts))
			fmt.Printf("Secret scrt amount %s \n", snapshotSecret.TotalTokenAmount)
			fmt.Printf("Secret Trst amount %s \n", snapshotSecret.TotalTrstAirdropAmount)
			fmt.Printf("# Secret scrt accounts %s \n", len(snapshotSecret.Accounts))
			fmt.Printf("Terra luna amount %s \n", snapshotTerra.TotalTokenAmount)
			fmt.Printf("Terra Trst amount %s \n", snapshotTerra.TotalTrstAirdropAmount)
			fmt.Printf("# Terra luna accounts %s \n", len(snapshotTerra.Accounts))
			snapshot.TotalTrstAirdropAmount = snapshot.TotalTrstAirdropAmount.Add(snapshotSecret.TotalTrstAirdropAmount.Add(snapshotTerra.TotalTrstAirdropAmount))
			snapshot.TotalTokenAmount = snapshot.TotalTokenAmount.Add(snapshotSecret.TotalTokenAmount.Add(snapshotTerra.TotalTokenAmount))

			snapshot.Accounts = removeDuplicatesFromSnapshot(snapshot.Accounts, snapshotSecret.Accounts, snapshotTerra.Accounts)

			fmt.Printf("# total TRST accounts %s \n", len(snapshot.Accounts))
			fmt.Printf("# total TRST coins %s \n", snapshot.TotalTrstAirdropAmount)
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
		//add so we ensure suffiient balance
		allTokens := acc.TokenBalance.ToDec()
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
					fmt.Printf(err.Error())
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
		allTokens := acc.TokenBalance.ToDec()

		allTokenSqrt := getMin(allTokens).RoundInt()

		if denominator.IsZero() {
			acc.TokenOwnershipPercent = sdk.NewInt(0).ToDec()
		} else {
			acc.TokenOwnershipPercent = allTokenSqrt.ToDec().QuoInt(denominator)
		}

		if allTokens.IsZero() {
			acc.TokenStakedPercent = sdk.ZeroDec()
			acc.TrstBalance = sdk.ZeroInt()
			snapshotAccs[address] = acc
			continue
		}

		stakedTokens := acc.TokenStakedBalance.ToDec()
		stakedPercent := stakedTokens.Quo(allTokens)

		acc.TokenStakedPercent = stakedPercent
		acc.TrstBalance = acc.TokenOwnershipPercent.MulInt(sdk.NewInt(TotalTrstAirdropAmountCosmos)).RoundInt()

		totalBalance = totalBalance.Add(acc.TrstBalance)
		snapshotAccount, ok := snapshot.Accounts[address]
		if !ok {
			snapshot.Accounts[address] = acc
			totalTokenBalance = totalTokenBalance.Add(acc.TokenBalance)
		} else {
			if snapshotAccount.TrstBalance.IsNil() {
				snapshotAccount.TrstBalance = sdk.ZeroInt()
			}
			snapshotAccount.TrstBalance = snapshotAccount.TrstBalance.Add(acc.TrstBalance)
			snapshotAccount.TokenBalance = snapshotAccount.TokenBalance.Add(acc.TokenBalance)
			snapshotAccount.TokenUnstakedBalance = snapshotAccount.TokenUnstakedBalance.Add(acc.TokenUnstakedBalance)
			snapshot.Accounts[address] = snapshotAccount

			totalTokenBalance = totalTokenBalance.Add(acc.TokenBalance)
		}
	}
	snapshot.TotalTokenAmount = totalTokenBalance
	snapshot.TotalTrstAirdropAmount = snapshot.TotalTrstAirdropAmount.Add(totalBalance)
	snapshot.NumberAccounts = snapshot.NumberAccounts + uint64(len(snapshot.Accounts))

	fmt.Printf("Complete read genesis file %s \n", genesisFile)
	fmt.Printf("# atom accounts: %d\n", len(snapshotAccs))
	fmt.Printf("# atom accounts in snapshot: %d\n", len(snapshot.Accounts))
	fmt.Printf("Token TotalSupply: %s\n", totalTokenBalance.String())
	fmt.Printf("Trst TotalSupply: %s\n", totalBalance.String())
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
		allTokens := acc.TokenBalance.ToDec()

		allTokenSqrt := getMin(allTokens).RoundInt()

		if denominator.IsZero() {
			acc.TokenOwnershipPercent = sdk.NewInt(0).ToDec()
		} else {
			acc.TokenOwnershipPercent = allTokenSqrt.ToDec().QuoInt(denominator)
		}

		if allTokens.IsZero() {
			acc.TokenStakedPercent = sdk.ZeroDec()
			acc.TrstBalance = sdk.ZeroInt()
			snapshotAccs[address] = acc
			continue
		}

		stakedTokens := acc.TokenStakedBalance.ToDec()
		stakedPercent := stakedTokens.Quo(allTokens)

		acc.TokenStakedPercent = stakedPercent
		acc.TrstBalance = acc.TokenOwnershipPercent.MulInt(sdk.NewInt(TotalTrstAirdropAmountSecret)).RoundInt()

		totalBalance = totalBalance.Add(acc.TrstBalance)
		snapshotAccount, ok := snapshot.Accounts[address]
		if !ok {
			snapshot.Accounts[address] = acc
			totalTokenBalance = totalTokenBalance.Add(acc.TokenBalance)
		} else {
			if snapshotAccount.TrstBalance.IsNil() {
				snapshotAccount.TrstBalance = sdk.ZeroInt()
			}
			snapshotAccount.TrstBalance = snapshotAccount.TrstBalance.Add(acc.TrstBalance)
			snapshotAccount.TokenBalance = snapshotAccount.TokenBalance.Add(acc.TokenBalance)
			snapshotAccount.TokenUnstakedBalance = snapshotAccount.TokenUnstakedBalance.Add(acc.TokenUnstakedBalance)
			snapshot.Accounts[address] = snapshotAccount

			totalTokenBalance = totalTokenBalance.Add(acc.TokenBalance)
		}
	}
	snapshot.TotalTokenAmount = totalTokenBalance
	snapshot.TotalTrstAirdropAmount = snapshot.TotalTrstAirdropAmount.Add(totalBalance)
	snapshot.NumberAccounts = snapshot.NumberAccounts + uint64(len(snapshot.Accounts))

	fmt.Printf("Complete read genesis file %s \n", genesisFile)
	fmt.Printf("# scrt accounts: %d\n", len(snapshotAccs))
	fmt.Printf("# scrt accounts in snapshot: %d\n", len(snapshot.Accounts))
	fmt.Printf("Scrt TotalSupply: %s\n", totalTokenBalance.String())
	fmt.Printf("Trst TotalSupply: %s\n", totalBalance.String())

	return snapshot
}

func exportTerraSnapShotFromGenesisFile(clientCtx client.Context, genesisFile string, denom string, snapshot Snapshot) Snapshot {
	appState, _, _ := genutiltypes.GenesisStateFromGenFile(genesisFile)

	var authGenState = appState["auth"]
	var auth AuthTerra
	err := json.Unmarshal(authGenState, &auth)
	if err != nil {
		fmt.Printf("auth module failed to unmarshal: %s\n", err)
	}

	var stakingGenState = appState["staking"]
	var staking StakingTerra
	err = json.Unmarshal(stakingGenState, &staking)

	if err != nil {
		fmt.Printf("staking module failed to unmarshal: %s\n", err)
	}

	snapshotAccs := make(map[string]SnapshotAccount)
	for _, account := range auth.Accounts {

		if account.Type == "/cosmos.auth.v1beta1.BaseAccount" || account.Type == "core/Account" {

			snapshotAccs[account.Address] = SnapshotAccount{
				TokenAddress:         account.Address,
				TokenBalance:         sdk.ZeroInt(),
				TokenUnstakedBalance: sdk.ZeroInt(),
				TokenStakedBalance:   sdk.ZeroInt(),
			}

		}
	}
	var bankGenState = appState["bank"]
	var bank BankTerra
	err = json.Unmarshal(bankGenState, &bank)
	// Produce the map of address to total atom balance, both staked and unstaked
	fmt.Printf("accounts %d \n", len(snapshotAccs))
	for _, account := range bank.Balances {

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
	fmt.Printf("Balances %s\n", len(bank.Balances))
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
	fmt.Printf("UnbondingDelegations %d \n", len(staking.UnbondingDelegations))
	// Make a map from validator operator address to the v42 validator type

	validators := make(map[string]ValidatorTerra)
	for _, validator := range staking.Validators {
		//	add, _ := sdk.ValAddressFromBech32(validator.Address)

		//	val, _ := stakingtypes.NewValidator(add, validator.PubKey, validator.Name)
		validators[validator.OperatorAddress] = validator
	}
	fmt.Printf("validator %d \n", len(staking.Validators))

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
		/*val := validators[delegation.ValidatorAddress]
		share, _ := strconv.ParseInt(delegation.Shares, 10, 64)
		valShare, _ := strconv.ParseInt(val.DelegatorShares, 10, 64)

		token, _ := sdk.NewIntFromString(val.Tokens)

		if valShare != 0 && val.Tokens != "0" {
			stakedTokens := sdk.NewDec(share).MulInt(token).Quo(sdk.NewDec(valShare)).RoundInt()

			acc.TokenBalance = acc.TokenBalance.Add(stakedTokens)
			acc.TokenStakedBalance = acc.TokenStakedBalance.Add(stakedTokens)
			snapshotAccs[address] = acc
		} else {
			//fmt.Printf("Zero val token balance for %s \n", address)
		}
		*/

	}
	fmt.Printf("delegations %d \n", len(staking.Delegations))

	denominator := getDenominator(snapshotAccs)
	totalBalance := sdk.ZeroInt()
	totalTokenBalance := sdk.NewInt(0)
	for address, acc := range snapshotAccs {
		allTokens := acc.TokenBalance.ToDec()

		allTokenSqrt := getMin(allTokens).RoundInt()

		if denominator.IsZero() {
			acc.TokenOwnershipPercent = sdk.NewInt(0).ToDec()
		} else {
			acc.TokenOwnershipPercent = allTokenSqrt.ToDec().QuoInt(denominator)
		}

		if allTokens.IsZero() {
			acc.TokenStakedPercent = sdk.ZeroDec()
			acc.TrstBalance = sdk.ZeroInt()
			snapshotAccs[address] = acc
			continue
		}

		stakedTokens := acc.TokenStakedBalance.ToDec()
		stakedPercent := stakedTokens.Quo(allTokens)

		acc.TokenStakedPercent = stakedPercent
		acc.TrstBalance = acc.TokenOwnershipPercent.MulInt(sdk.NewInt(TotalTrstAirdropAmountTerra)).RoundInt()

		totalBalance = totalBalance.Add(acc.TrstBalance)
		snapshotAccount, ok := snapshot.Accounts[address]
		if !ok {
			snapshot.Accounts[address] = acc
			totalTokenBalance = totalTokenBalance.Add(acc.TokenBalance)
		} else {
			if snapshotAccount.TrstBalance.IsNil() {
				snapshotAccount.TrstBalance = sdk.ZeroInt()
			}
			snapshotAccount.TrstBalance = snapshotAccount.TrstBalance.Add(acc.TrstBalance)
			snapshotAccount.TokenBalance = snapshotAccount.TokenBalance.Add(acc.TokenBalance)
			snapshotAccount.TokenUnstakedBalance = snapshotAccount.TokenUnstakedBalance.Add(acc.TokenUnstakedBalance)
			snapshot.Accounts[address] = snapshotAccount

			totalTokenBalance = totalTokenBalance.Add(acc.TokenBalance)
		}
	}
	snapshot.TotalTokenAmount = totalTokenBalance
	snapshot.TotalTrstAirdropAmount = snapshot.TotalTrstAirdropAmount.Add(totalBalance)
	snapshot.NumberAccounts = snapshot.NumberAccounts + uint64(len(snapshot.Accounts))

	fmt.Printf("Complete read genesis file %s \n", genesisFile)
	fmt.Printf("# luna accounts: %d\n", len(snapshotAccs))
	fmt.Printf("# luna accounts in snapshot: %d\n", len(snapshot.Accounts))
	fmt.Printf("Luna TotalSupply: %s\n", totalTokenBalance.String())
	fmt.Printf("Trst TotalSupply: %s\n", totalBalance.String())

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
	TotalTokenAmount       sdk.Int `json:"total_atom_amount"`
	TotalTrstAirdropAmount sdk.Int `json:"total_trst_amount"`
	NumberAccounts         uint64  `json:"num_accounts"`

	Accounts map[string]SnapshotAccount `json:"accounts"`
}

// SnapshotAccount provide fields of snapshot per account
type SnapshotAccount struct {
	TokenAddress string `json:"atom_address"` // Token Balance = TokenStakedBalance + TokenUnstakedBalance

	TokenBalance          sdk.Int `json:"atom_balance"`
	TokenOwnershipPercent sdk.Dec `json:"atom_ownership_percent"`

	TokenStakedBalance   sdk.Int `json:"atom_staked_balance"`
	TokenUnstakedBalance sdk.Int `json:"atom_unstaked_balance"` // TokenStakedPercent = TokenStakedBalance / TokenBalance
	TokenStakedPercent   sdk.Dec `json:"atom_staked_percent"`

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

type AuthTerra struct {
	Accounts []struct {
		Type          string `json:"@type"`
		AccountNumber string `json:"account_number"`
		Address       string `json:"address"`
		PubKey        struct {
			Type string `json:"@type"`
			Key  string `json:"key"`
		} `json:"pub_key"`
		Sequence string `json:"sequence"`
	} `json:"accounts"`
	Params struct {
		MaxMemoCharacters      string `json:"max_memo_characters"`
		SigVerifyCostEd25519   string `json:"sig_verify_cost_ed25519"`
		SigVerifyCostSecp256K1 string `json:"sig_verify_cost_secp256k1"`
		TxSigLimit             string `json:"tx_sig_limit"`
		TxSizeCostPerByte      string `json:"tx_size_cost_per_byte"`
	} `json:"params"`
}
type StakingTerra struct {
	Delegations []struct {
		DelegatorAddress string `json:"delegator_address"`
		Shares           string `json:"shares"`
		ValidatorAddress string `json:"validator_address"`
	} `json:"delegations"`
	Exported            bool   `json:"exported"`
	LastTotalPower      string `json:"last_total_power"`
	LastValidatorPowers []struct {
		Address string `json:"address"`
		Power   string `json:"power"`
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
	Validators           []ValidatorTerra       `json:"validators"`
}

type BankTerra struct {
	Params struct {
		SendEnabled        []interface{} `json:"send_enabled"`
		DefaultSendEnabled bool          `json:"default_send_enabled"`
	} `json:"params"`
	Balances []struct {
		Address string    `json:"address"`
		Coins   sdk.Coins `json:"coins"`
	} `json:"balances"`
	Supply []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"supply"`
	DenomMetadata []interface{} `json:"denom_metadata"`
}

type Delegation struct {
	DelegatorAddress string `json:"delegator_address"`
	Shares           string `json:"shares"`
	ValidatorAddress string `json:"validator_address"`
}

type ValidatorTerra struct {
	Commission struct {
		CommissionRates struct {
			MaxChangeRate string `json:"max_change_rate"`
			MaxRate       string `json:"max_rate"`
			Rate          string `json:"rate"`
		} `json:"commission_rates"`
		UpdateTime time.Time `json:"update_time"`
	} `json:"commission"`
	ConsensusPubkey struct {
		Type string `json:"@type"`
		Key  string `json:"key"`
	} `json:"consensus_pubkey"`
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
	Status            string    `json:"status"`
	Tokens            string    `json:"tokens"`
	UnbondingHeight   string    `json:"unbonding_height"`
	UnbondingTime     time.Time `json:"unbonding_time"`
}

func removeDuplicates(m map[string]SnapshotAccount) map[string]SnapshotAccount {
	list := make([]string, 0, len(m))

	//append accounts to list so we have list of all acc
	for acc, _ := range m {

		list = append(list, acc)
	}

	for _, item := range list {

		if contains(list, item) == false {
			_, ok := m[item]
			if ok {
				delete(m, item)
			}

			//list = append(list, item)
		}
	}

	return m
}

func removeDuplicatesFromSnapshot(first map[string]SnapshotAccount, second map[string]SnapshotAccount, third map[string]SnapshotAccount) map[string]SnapshotAccount {
	list := make([]string, 0, len(first))

	var duplicatesSecret uint64
	var duplicatesTerra uint64
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
		if contains(list, newAcc) == false {

			list = append(list, newAcc)
			first[newAcc] = value

		} else {
			//	fmt.Println("delete secret acc  %s\n", acc)
			var acct = first[newAcc]
			acct.TrstBalance = acct.TrstBalance.Add(value.TrstBalance)
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
		if contains(list, newAcc) == false {

			list = append(list, newAcc)
			first[newAcc] = value
		} else {
			//	fmt.Println("delete secret acc  %s\n", acc)
			var acct = first[newAcc]
			acct.TrstBalance = acct.TrstBalance.Add(value.TrstBalance)
			first[newAcc] = acct

			fmt.Println("merged terra acc", acc)
			delete(third, acc)
			duplicatesTerra = duplicatesTerra + 1
		}
	}

	fmt.Printf("final length secret accs: c\n", len(second))
	fmt.Printf("final length terra accs: %d\n", len(third))
	fmt.Printf("length list accs: %d\n", len(list))

	fmt.Printf("duplicates secret accs: %d\n", duplicatesSecret)
	fmt.Printf("duplicates terra accs: %d\n", duplicatesTerra)
	//all := append(first, second)

	fmt.Println("removed duplicates")
	return first
}

/*
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
*/

/*
type TerraGenesis struct {
	AppHash  string `json:"app_hash"`
	AppState struct {
		Auth struct {
			Accounts []struct {
				Type          string `json:"@type"`
				AccountNumber string `json:"account_number"`
				Address       string `json:"address"`
				PubKey        struct {
					Type string `json:"@type"`
					Key  string `json:"key"`
				} `json:"pub_key"`
				Sequence string `json:"sequence"`
			} `json:"accounts"`
			Params struct {
				MaxMemoCharacters      string `json:"max_memo_characters"`
				SigVerifyCostEd25519   string `json:"sig_verify_cost_ed25519"`
				SigVerifyCostSecp256K1 string `json:"sig_verify_cost_secp256k1"`
				TxSigLimit             string `json:"tx_sig_limit"`
				TxSizeCostPerByte      string `json:"tx_size_cost_per_byte"`
			} `json:"params"`
		} `json:"auth"`
		Authz struct {
			Authorization []struct {
				Authorization struct {
					Type       string `json:"@type"`
					SpendLimit []struct {
						Amount string `json:"amount"`
						Denom  string `json:"denom"`
					} `json:"spend_limit"`
				} `json:"authorization"`
				Expiration time.Time `json:"expiration"`
				Grantee    string    `json:"grantee"`
				Granter    string    `json:"granter"`
			} `json:"authorization"`
		} `json:"authz"`
		Bank struct {
			Balances []struct {
				Address string `json:"address"`
				Coins   []struct {
					Amount string `json:"amount"`
					Denom  string `json:"denom"`
				} `json:"coins"`
			} `json:"balances"`
			DenomMetadata []struct {
				Base       string `json:"base"`
				DenomUnits []struct {
					Aliases  []string `json:"aliases"`
					Denom    string   `json:"denom"`
					Exponent int      `json:"exponent"`
				} `json:"denom_units"`
				Description string `json:"description"`
				Display     string `json:"display"`
				Name        string `json:"name"`
				Symbol      string `json:"symbol"`
			} `json:"denom_metadata"`
			Params struct {
				DefaultSendEnabled bool          `json:"default_send_enabled"`
				SendEnabled        []interface{} `json:"send_enabled"`
			} `json:"params"`
			Supply []struct {
				Amount string `json:"amount"`
				Denom  string `json:"denom"`
			} `json:"supply"`
		} `json:"bank"`
		Capability struct {
			Index  string        `json:"index"`
			Owners []interface{} `json:"owners"`
		} `json:"capability"`
		Crisis struct {
			ConstantFee struct {
				Amount string `json:"amount"`
				Denom  string `json:"denom"`
			} `json:"constant_fee"`
		} `json:"crisis"`
		Distribution struct {
			DelegatorStartingInfos []struct {
				DelegatorAddress string `json:"delegator_address"`
				StartingInfo     struct {
					Height         string `json:"height"`
					PreviousPeriod string `json:"previous_period"`
					Stake          string `json:"stake"`
				} `json:"starting_info"`
				ValidatorAddress string `json:"validator_address"`
			} `json:"delegator_starting_infos"`
			DelegatorWithdrawInfos []struct {
				DelegatorAddress string `json:"delegator_address"`
				WithdrawAddress  string `json:"withdraw_address"`
			} `json:"delegator_withdraw_infos"`
			FeePool struct {
				CommunityPool []struct {
					Amount string `json:"amount"`
					Denom  string `json:"denom"`
				} `json:"community_pool"`
			} `json:"fee_pool"`
			OutstandingRewards []struct {
				OutstandingRewards []struct {
					Amount string `json:"amount"`
					Denom  string `json:"denom"`
				} `json:"outstanding_rewards"`
				ValidatorAddress string `json:"validator_address"`
			} `json:"outstanding_rewards"`
			Params struct {
				BaseProposerReward  string `json:"base_proposer_reward"`
				BonusProposerReward string `json:"bonus_proposer_reward"`
				CommunityTax        string `json:"community_tax"`
				WithdrawAddrEnabled bool   `json:"withdraw_addr_enabled"`
			} `json:"params"`
			PreviousProposer                string `json:"previous_proposer"`
			ValidatorAccumulatedCommissions []struct {
				Accumulated struct {
					Commission []struct {
						Amount string `json:"amount"`
						Denom  string `json:"denom"`
					} `json:"commission"`
				} `json:"accumulated"`
				ValidatorAddress string `json:"validator_address"`
			} `json:"validator_accumulated_commissions"`
			ValidatorCurrentRewards []struct {
				Rewards struct {
					Period  string `json:"period"`
					Rewards []struct {
						Amount string `json:"amount"`
						Denom  string `json:"denom"`
					} `json:"rewards"`
				} `json:"rewards"`
				ValidatorAddress string `json:"validator_address"`
			} `json:"validator_current_rewards"`
			ValidatorHistoricalRewards []struct {
				Period  string `json:"period"`
				Rewards struct {
					CumulativeRewardRatio []interface{} `json:"cumulative_reward_ratio"`
					ReferenceCount        int           `json:"reference_count"`
				} `json:"rewards"`
				ValidatorAddress string `json:"validator_address"`
			} `json:"validator_historical_rewards"`
			ValidatorSlashEvents []struct {
				Height              string `json:"height"`
				Period              string `json:"period"`
				ValidatorAddress    string `json:"validator_address"`
				ValidatorSlashEvent struct {
					Fraction        string `json:"fraction"`
					ValidatorPeriod string `json:"validator_period"`
				} `json:"validator_slash_event"`
			} `json:"validator_slash_events"`
		} `json:"distribution"`
		Evidence struct {
			Evidence []interface{} `json:"evidence"`
		} `json:"evidence"`
		Genutil struct {
			GenTxs []interface{} `json:"gen_txs"`
		} `json:"genutil"`
		Gov struct {
			DepositParams struct {
				MaxDepositPeriod string `json:"max_deposit_period"`
				MinDeposit       []struct {
					Amount string `json:"amount"`
					Denom  string `json:"denom"`
				} `json:"min_deposit"`
			} `json:"deposit_params"`
			Deposits  []interface{} `json:"deposits"`
			Proposals []struct {
				Content struct {
					Type        string `json:"@type"`
					Description string `json:"description"`
					Title       string `json:"title"`
				} `json:"content"`
				DepositEndTime   time.Time `json:"deposit_end_time"`
				FinalTallyResult struct {
					Abstain    string `json:"abstain"`
					No         string `json:"no"`
					NoWithVeto string `json:"no_with_veto"`
					Yes        string `json:"yes"`
				} `json:"final_tally_result"`
				ProposalID   string    `json:"proposal_id"`
				Status       string    `json:"status"`
				SubmitTime   time.Time `json:"submit_time"`
				TotalDeposit []struct {
					Amount string `json:"amount"`
					Denom  string `json:"denom"`
				} `json:"total_deposit"`
				VotingEndTime   time.Time `json:"voting_end_time"`
				VotingStartTime time.Time `json:"voting_start_time"`
			} `json:"proposals"`
			StartingProposalID string `json:"starting_proposal_id"`
			TallyParams        struct {
				Quorum        string `json:"quorum"`
				Threshold     string `json:"threshold"`
				VetoThreshold string `json:"veto_threshold"`
			} `json:"tally_params"`
			Votes        []interface{} `json:"votes"`
			VotingParams struct {
				VotingPeriod string `json:"voting_period"`
			} `json:"voting_params"`
		} `json:"gov"`
		Ibc struct {
			ChannelGenesis struct {
				AckSequences        []interface{} `json:"ack_sequences"`
				Acknowledgements    []interface{} `json:"acknowledgements"`
				Channels            []interface{} `json:"channels"`
				Commitments         []interface{} `json:"commitments"`
				NextChannelSequence string        `json:"next_channel_sequence"`
				Receipts            []interface{} `json:"receipts"`
				RecvSequences       []interface{} `json:"recv_sequences"`
				SendSequences       []interface{} `json:"send_sequences"`
			} `json:"channel_genesis"`
			ClientGenesis struct {
				Clients            []interface{} `json:"clients"`
				ClientsConsensus   []interface{} `json:"clients_consensus"`
				ClientsMetadata    []interface{} `json:"clients_metadata"`
				CreateLocalhost    bool          `json:"create_localhost"`
				NextClientSequence string        `json:"next_client_sequence"`
				Params             struct {
					AllowedClients []string `json:"allowed_clients"`
				} `json:"params"`
			} `json:"client_genesis"`
			ConnectionGenesis struct {
				ClientConnectionPaths  []interface{} `json:"client_connection_paths"`
				Connections            []interface{} `json:"connections"`
				NextConnectionSequence string        `json:"next_connection_sequence"`
				Params                 struct {
					MaxExpectedTimePerBlock string `json:"max_expected_time_per_block"`
				} `json:"params"`
			} `json:"connection_genesis"`
		} `json:"ibc"`
		Market struct {
			Params struct {
				BasePool           string `json:"base_pool"`
				MinStabilitySpread string `json:"min_stability_spread"`
				PoolRecoveryPeriod string `json:"pool_recovery_period"`
			} `json:"params"`
			TerraPoolDelta string `json:"terra_pool_delta"`
		} `json:"market"`
		Mint struct {
			Minter struct {
				AnnualProvisions string `json:"annual_provisions"`
				Inflation        string `json:"inflation"`
			} `json:"minter"`
			Params struct {
				BlocksPerYear       string `json:"blocks_per_year"`
				GoalBonded          string `json:"goal_bonded"`
				InflationMax        string `json:"inflation_max"`
				InflationMin        string `json:"inflation_min"`
				InflationRateChange string `json:"inflation_rate_change"`
				MintDenom           string `json:"mint_denom"`
			} `json:"params"`
		} `json:"mint"`
		Oracle struct {
			AggregateExchangeRatePrevotes []struct {
				Hash        string `json:"hash"`
				SubmitBlock string `json:"submit_block"`
				Voter       string `json:"voter"`
			} `json:"aggregate_exchange_rate_prevotes"`
			AggregateExchangeRateVotes []struct {
				ExchangeRateTuples []struct {
					Denom        string `json:"denom"`
					ExchangeRate string `json:"exchange_rate"`
				} `json:"exchange_rate_tuples"`
				Voter string `json:"voter"`
			} `json:"aggregate_exchange_rate_votes"`
			ExchangeRates []struct {
				Denom        string `json:"denom"`
				ExchangeRate string `json:"exchange_rate"`
			} `json:"exchange_rates"`
			FeederDelegations []struct {
				FeederAddress    string `json:"feeder_address"`
				ValidatorAddress string `json:"validator_address"`
			} `json:"feeder_delegations"`
			MissCounters []struct {
				MissCounter      string `json:"miss_counter"`
				ValidatorAddress string `json:"validator_address"`
			} `json:"miss_counters"`
			Params struct {
				MinValidPerWindow        string `json:"min_valid_per_window"`
				RewardBand               string `json:"reward_band"`
				RewardDistributionWindow string `json:"reward_distribution_window"`
				SlashFraction            string `json:"slash_fraction"`
				SlashWindow              string `json:"slash_window"`
				VotePeriod               string `json:"vote_period"`
				VoteThreshold            string `json:"vote_threshold"`
				Whitelist                []struct {
					Name     string `json:"name"`
					TobinTax string `json:"tobin_tax"`
				} `json:"whitelist"`
			} `json:"params"`
			TobinTaxes []struct {
				Denom    string `json:"denom"`
				TobinTax string `json:"tobin_tax"`
			} `json:"tobin_taxes"`
		} `json:"oracle"`
		Slashing struct {
			MissedBlocks []struct {
				Address      string `json:"address"`
				MissedBlocks []struct {
					Index  string `json:"index"`
					Missed bool   `json:"missed"`
				} `json:"missed_blocks"`
			} `json:"missed_blocks"`
			Params struct {
				DowntimeJailDuration    string `json:"downtime_jail_duration"`
				MinSignedPerWindow      string `json:"min_signed_per_window"`
				SignedBlocksWindow      string `json:"signed_blocks_window"`
				SlashFractionDoubleSign string `json:"slash_fraction_double_sign"`
				SlashFractionDowntime   string `json:"slash_fraction_downtime"`
			} `json:"params"`
			SigningInfos []struct {
				Address              string `json:"address"`
				ValidatorSigningInfo struct {
					Address             string    `json:"address"`
					IndexOffset         string    `json:"index_offset"`
					JailedUntil         time.Time `json:"jailed_until"`
					MissedBlocksCounter string    `json:"missed_blocks_counter"`
					StartHeight         string    `json:"start_height"`
					Tombstoned          bool      `json:"tombstoned"`
				} `json:"validator_signing_info"`
			} `json:"signing_infos"`
		} `json:"slashing"`
		Staking struct {
			Delegations         []Delegation `json:"delegations"`
			Exported            bool         `json:"exported"`
			LastTotalPower      string       `json:"last_total_power"`
			LastValidatorPowers []struct {
				Address string `json:"address"`
				Power   string `json:"power"`
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
			Validators           []struct {
				Commission struct {
					CommissionRates struct {
						MaxChangeRate string `json:"max_change_rate"`
						MaxRate       string `json:"max_rate"`
						Rate          string `json:"rate"`
					} `json:"commission_rates"`
					UpdateTime time.Time `json:"update_time"`
				} `json:"commission"`
				ConsensusPubkey struct {
					Type string `json:"@type"`
					Key  string `json:"key"`
				} `json:"consensus_pubkey"`
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
				Status            string    `json:"status"`
				Tokens            string    `json:"tokens"`
				UnbondingHeight   string    `json:"unbonding_height"`
				UnbondingTime     time.Time `json:"unbonding_time"`
			} `json:"validators"`
		} `json:"staking"`
		Transfer struct {
			DenomTraces []interface{} `json:"denom_traces"`
			Params      struct {
				ReceiveEnabled bool `json:"receive_enabled"`
				SendEnabled    bool `json:"send_enabled"`
			} `json:"params"`
			PortID string `json:"port_id"`
		} `json:"transfer"`
		Treasury struct {
			EpochInitialIssuance []struct {
				Amount string `json:"amount"`
				Denom  string `json:"denom"`
			} `json:"epoch_initial_issuance"`
			EpochStates []struct {
				Epoch             string `json:"epoch"`
				SeigniorageReward string `json:"seigniorage_reward"`
				TaxReward         string `json:"tax_reward"`
				TotalStakedLuna   string `json:"total_staked_luna"`
			} `json:"epoch_states"`
			Params struct {
				MiningIncrement string `json:"mining_increment"`
				RewardPolicy    struct {
					Cap struct {
						Amount string `json:"amount"`
						Denom  string `json:"denom"`
					} `json:"cap"`
					ChangeRateMax string `json:"change_rate_max"`
					RateMax       string `json:"rate_max"`
					RateMin       string `json:"rate_min"`
				} `json:"reward_policy"`
				SeigniorageBurdenTarget string `json:"seigniorage_burden_target"`
				TaxPolicy               struct {
					Cap struct {
						Amount string `json:"amount"`
						Denom  string `json:"denom"`
					} `json:"cap"`
					ChangeRateMax string `json:"change_rate_max"`
					RateMax       string `json:"rate_max"`
					RateMin       string `json:"rate_min"`
				} `json:"tax_policy"`
				WindowLong      string `json:"window_long"`
				WindowProbation string `json:"window_probation"`
				WindowShort     string `json:"window_short"`
			} `json:"params"`
			RewardWeight string `json:"reward_weight"`
			TaxCaps      []struct {
				Denom  string `json:"denom"`
				TaxCap string `json:"tax_cap"`
			} `json:"tax_caps"`
			TaxProceeds []struct {
				Amount string `json:"amount"`
				Denom  string `json:"denom"`
			} `json:"tax_proceeds"`
			TaxRate string `json:"tax_rate"`
		} `json:"treasury"`
		Upgrade struct {
		} `json:"upgrade"`
		Wasm struct {
			Codes []struct {
				CodeBytes string `json:"code_bytes"`
				CodeInfo  struct {
					CodeHash string `json:"code_hash"`
					CodeID   string `json:"code_id"`
					Creator  string `json:"creator"`
				} `json:"code_info"`
			} `json:"codes"`
			Contracts []struct {
				ContractInfo struct {
					Address string `json:"address"`
					Admin   string `json:"admin"`
					CodeID  string `json:"code_id"`
					Creator string `json:"creator"`
					InitMsg struct {
						IsPublic            bool   `json:"is_public"`
						PoolManager         string `json:"pool_manager"`
						SparFactoryContract string `json:"spar_factory_contract"`
					} `json:"init_msg"`
				} `json:"contract_info"`
				ContractStore []struct {
					Key   string `json:"key"`
					Value string `json:"value"`
				} `json:"contract_store"`
			} `json:"contracts"`
			LastCodeID     string `json:"last_code_id"`
			LastInstanceID string `json:"last_instance_id"`
			Params         struct {
				MaxContractGas     string `json:"max_contract_gas"`
				MaxContractMsgSize string `json:"max_contract_msg_size"`
				MaxContractSize    string `json:"max_contract_size"`
			} `json:"params"`
		} `json:"wasm"`
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
			MaxBytes        string `json:"max_bytes"`
		} `json:"evidence"`
		Validator struct {
			PubKeyTypes []string `json:"pub_key_types"`
		} `json:"validator"`
		Version struct {
		} `json:"version"`
	} `json:"consensus_params"`
	GenesisTime   time.Time `json:"genesis_time"`
	InitialHeight string    `json:"initial_height"`
	Validators    []struct {
		Address string `json:"address"`
		Name    string `json:"name"`
		Power   string `json:"power"`
		PubKey  struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"pub_key"`
	} `json:"validators"`
}
*/
