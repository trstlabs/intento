package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"

	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"
	claimtypes "github.com/trstlabs/intento/x/claim/types"
)

const (
	flagVestingStart = "vesting-start-time"
	flagVestingEnd   = "vesting-end-time"
	flagVestingAmt   = "vesting-amount"
)

// SnapshotEntry represents an entry in the snapshot.
type SnapshotEntry struct {
	Address string   `json:"address"`
	Weight  math.Int `json:"weight"`
}

// Snapshot represents the overall snapshot as a slice of SnapshotEntry.
type Snapshot []SnapshotEntry

// AddGenesisAccountCmd returns add-genesis-account cobra Command.
func AddGenesisAccountCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-account [address_or_key_name] [coin][,[coin]]",
		Short: "Add a genesis account to genesis.json",
		Long: `Add a genesis account to genesis.json. The provided account must specify
the account address or key name and a list of initial coins. If a key name is given,
the address will be looked up in the local Keybase. The list of initial tokens must
contain valid denominations. Accounts may optionally be supplied with vesting parameters.
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			depCdc := clientCtx.Codec
			cdc := depCdc

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			coins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return fmt.Errorf("failed to parse coins: %w", err)
			}

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				keyringBackend, err := cmd.Flags().GetString(flags.FlagKeyringBackend)
				if err != nil {
					return err
				}

				// attempt to lookup address from Keybase if no address was provided
				kb, err := keyring.New(sdk.KeyringServiceName(), keyringBackend, clientCtx.HomeDir, inBuf, clientCtx.Codec)
				if err != nil {
					return err
				}

				info, err := kb.Key(args[0])
				if err != nil {
					return fmt.Errorf("failed to get address from Keybase: %w", err)
				}

				addr, err = info.GetAddress()
				if err != nil {
					return err
				}
			}

			vestingStart, err := cmd.Flags().GetInt64(flagVestingStart)
			if err != nil {
				return err
			}
			vestingEnd, err := cmd.Flags().GetInt64(flagVestingEnd)
			if err != nil {
				return err
			}
			vestingAmtStr, err := cmd.Flags().GetString(flagVestingAmt)
			if err != nil {
				return err
			}

			vestingAmt, err := sdk.ParseCoinsNormalized(vestingAmtStr)
			if err != nil {
				return fmt.Errorf("failed to parse vesting amount: %w", err)
			}

			// create concrete account type based on input parameters
			var genAccount authtypes.GenesisAccount

			balances := banktypes.Balance{Address: addr.String(), Coins: coins.Sort()}
			baseAccount := authtypes.NewBaseAccount(addr, nil, 0, 0)

			if !vestingAmt.IsZero() {
				baseVestingAccount, err := authvesting.NewBaseVestingAccount(baseAccount, vestingAmt.Sort(), vestingEnd)
				if err != nil {
					return err
				}
				if (balances.Coins.IsZero() && !baseVestingAccount.OriginalVesting.IsZero()) ||
					baseVestingAccount.OriginalVesting.IsAnyGT(balances.Coins) {
					return errors.New("vesting amount cannot be greater than total amount")
				}

				switch {
				case vestingStart != 0 && vestingEnd != 0:
					genAccount = authvesting.NewContinuousVestingAccountRaw(baseVestingAccount, vestingStart)

				case vestingEnd != 0:
					genAccount = authvesting.NewDelayedVestingAccountRaw(baseVestingAccount)

				default:
					return errors.New("invalid vesting parameters; must supply start and end time or end time")
				}
			} else {
				genAccount = baseAccount
			}

			if err := genAccount.Validate(); err != nil {
				return fmt.Errorf("failed to validate new genesis account: %w", err)
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			authGenState := authtypes.GetGenesisStateFromAppState(cdc, appState)

			accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
			if err != nil {
				return fmt.Errorf("failed to get accounts from any: %w", err)
			}

			if accs.Contains(addr) {
				return fmt.Errorf("cannot add account at existing address %s", addr)
			}

			// Add the new account to the set of genesis accounts and sanitize the
			// accounts afterwards.
			accs = append(accs, genAccount)
			accs = authtypes.SanitizeGenesisAccounts(accs)

			genAccs, err := authtypes.PackAccounts(accs)
			if err != nil {
				return fmt.Errorf("failed to convert accounts into any's: %w", err)
			}
			authGenState.Accounts = genAccs

			authGenStateBz, err := cdc.MarshalJSON(&authGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}

			appState[authtypes.ModuleName] = authGenStateBz

			bankGenState := banktypes.GetGenesisStateFromAppState(depCdc, appState)
			bankGenState.Balances = append(bankGenState.Balances, balances)
			bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)

			bankGenStateBz, err := cdc.MarshalJSON(bankGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal bank genesis state: %w", err)
			}

			appState[banktypes.ModuleName] = bankGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test)")
	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.Flags().String(flagVestingAmt, "", "amount of coins for vesting accounts")
	cmd.Flags().Int64(flagVestingStart, 0, "schedule start time (unix epoch) for vesting accounts")
	cmd.Flags().Int64(flagVestingEnd, 0, "schedule end time (unix epoch) for vesting accounts")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
func ImportGenesisAccountsFromSnapshotCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import-genesis-accounts-from-snapshot [input-snapshot-file] [input-non-airdrop-accounts-file]",
		Short: "Import genesis accounts from a snapshot and distribute an airdrop amount.",
		Long: `Import genesis accounts from a snapshot file and a non-airdrop accounts file.
Distribute the specified total airdrop amount among the snapshot accounts proportionally based on weights.
Also include non-airdrop accounts with specified destinations.

Example:
intentod import-genesis-accounts-from-snapshot ../snapshot.json ../non-airdrop-accounts.json --airdrop-amount=1000000000
- Input genesis file: ~/.intentod/config/genesis.json`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			config.SetRoot(clientCtx.HomeDir)

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			authGenState := authtypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
			accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
			if err != nil {
				return fmt.Errorf("failed to get accounts from genesis state: %w", err)
			}

			airdropAmountFlag, err := cmd.Flags().GetInt64("airdrop-amount")
			if err != nil || airdropAmountFlag <= 0 {
				return fmt.Errorf("invalid or missing --airdrop-amount flag")
			}
			airdropAmount := math.NewInt(airdropAmountFlag)

			// Load snapshot
			snapshotFile := args[0]
			snapshotJSON, err := os.Open(snapshotFile)
			if err != nil {
				return fmt.Errorf("failed to open snapshot file: %w", err)
			}
			defer snapshotJSON.Close()

			var snapshot Snapshot
			if err := json.NewDecoder(snapshotJSON).Decode(&snapshot); err != nil {
				return fmt.Errorf("failed to parse snapshot file: %w", err)
			}

			nonAirdropFile := args[1]
			nonAirdropJSON, err := os.Open(nonAirdropFile)
			if err != nil {
				return fmt.Errorf("failed to open non-airdrop accounts file: %w", err)
			}
			defer nonAirdropJSON.Close()

			var nonAirdropAccounts []struct {
				Address string `json:"address"`
				Amount  int64  `json:"amount"`
				Name    string `json:"name"`
			}
			if err := json.NewDecoder(nonAirdropJSON).Decode(&nonAirdropAccounts); err != nil {
				return fmt.Errorf("failed to parse non-airdrop accounts file: %w", err)
			}

			// Deduplicate snapshot first
			addressToWeight := make(map[string]math.Int)
			addressToEntry := make(map[string]SnapshotEntry)
			for _, account := range snapshot {
				if prevWeight, exists := addressToWeight[account.Address]; exists {
					if account.Weight.GT(prevWeight) {
						addressToWeight[account.Address] = account.Weight
						addressToEntry[account.Address] = account
					}
				} else {
					addressToWeight[account.Address] = account.Weight
					addressToEntry[account.Address] = account
				}
			}
			dedupedSnapshot := make([]SnapshotEntry, 0, len(addressToEntry))
			for _, entry := range addressToEntry {
				dedupedSnapshot = append(dedupedSnapshot, entry)
			}
			// Now calculate totalWeight and normalizationFactor using dedupedSnapshot
			var totalWeight math.LegacyDec = math.LegacyNewDec(0)
			for _, account := range dedupedSnapshot {
				totalWeight = totalWeight.Add(math.LegacyNewDecFromInt(account.Weight))
			}
			normalizationFactor := math.LegacyNewDecFromInt(airdropAmount).Quo(totalWeight)

			bankGenState := banktypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
			liquidBalances := bankGenState.Balances
			claimRecords := []claimtypes.ClaimRecord{}
			claimModuleBalance := math.NewInt(0)
			// ...existing code...
			accountMap := make(map[string]bool)
			totalDistributed := math.NewInt(0)

			// Track largest claimable

			largestClaimable := math.NewInt(0)
			claimableAmounts := make([]math.Int, len(dedupedSnapshot))
			addresses := make([]string, len(dedupedSnapshot))
			for i, account := range dedupedSnapshot {
				address, err := ConvertBech32(account.Address)
				if err != nil {
					return fmt.Errorf("invalid address in snapshot: %w", err)
				}

				if accountMap[address] {
					// Should not happen after deduplication
					continue
				}
				accountMap[address] = true

				airdropShare := math.LegacyNewDecFromInt(account.Weight).Mul(normalizationFactor).TruncateInt()

				liquidAmount := airdropShare.MulRaw(2).QuoRaw(10) // 20%
				maxLiquid := math.NewInt(2_000_000)
				if liquidAmount.GT(maxLiquid) {
					liquidAmount = maxLiquid
				}

				claimableAmount := airdropShare.Sub(liquidAmount)
				claimableAmounts[i] = claimableAmount
				addresses[i] = address
				if claimableAmount.GT(largestClaimable) {
					largestClaimable = claimableAmount

				}
				if !airdropShare.Equal(liquidAmount.Add(claimableAmount)) {
					return fmt.Errorf("math error: liquid + claimable ≠ airdropShare for %s", address)
				}

				liquidBalances = append(liquidBalances, banktypes.Balance{
					Address: address,
					Coins:   sdk.NewCoins(sdk.NewCoin("uinto", liquidAmount)),
				})
				claimModuleBalance = claimModuleBalance.Add(claimableAmount)
				totalDistributed = totalDistributed.Add(airdropShare)

				accAddr, err := sdk.AccAddressFromBech32(address)
				if err != nil {
					return fmt.Errorf("bad address: %w", err)
				}
				baseAccount := authtypes.NewBaseAccount(accAddr, nil, 0, 0)
				if err := baseAccount.Validate(); err != nil {
					return fmt.Errorf("invalid base account: %w", err)
				}
				accs = append(accs, baseAccount)

				status := claimtypes.Status{ActionCompleted: false, VestingPeriodsCompleted: []bool{false, false, false, false}, VestingPeriodsClaimed: []bool{false, false, false, false}}
				claimRecords = append(claimRecords, claimtypes.ClaimRecord{
					Address:                address,
					MaximumClaimableAmount: sdk.NewCoin("uinto", claimableAmount),
					Status:                 []claimtypes.Status{status, status, status, status},
				})
			}

			// Handle airdrop remainder
			remainder := airdropAmount.Sub(totalDistributed)
			if remainder.GT(math.NewInt(0)) {
				// Option 2: Send to burner address
				burnerAddress, err := sdk.AccAddressFromHexUnsafe("7C4954EAE77FF15A4C67C5F821C5241008ED966D") // Cosmos random address with no keypair
				if err != nil {
					panic(err)
				}
				fmt.Printf("Burner address: %s\n", burnerAddress.String())
				status := claimtypes.Status{ActionCompleted: false, VestingPeriodsCompleted: []bool{false, false, false, false}, VestingPeriodsClaimed: []bool{false, false, false, false}}
				claimRecords = append(claimRecords, claimtypes.ClaimRecord{
					Address:                burnerAddress.String(),
					MaximumClaimableAmount: sdk.NewCoin("uinto", remainder),
					Status:                 []claimtypes.Status{status, status, status, status},
				})
				claimModuleBalance = claimModuleBalance.Add(remainder)
				fmt.Printf("Airdrop remainder %s uinto sent to burner address %s\n", remainder.String(), burnerAddress)
			}

			// Add non-airdrop accounts
			for _, info := range nonAirdropAccounts {
				if info.Address == "" {
					return fmt.Errorf("missing address for non-airdrop account: %s", info.Name)
				}

				if accountMap[info.Address] {
					fmt.Printf("Skipping duplicate non-airdrop address: %s\n", info.Address)
					continue
				}
				accountMap[info.Address] = true

				address, err := sdk.AccAddressFromBech32(info.Address)
				if err != nil {
					return fmt.Errorf("invalid non-airdrop address: %w", err)
				}

				baseAccount := authtypes.NewBaseAccount(address, nil, 0, 0)
				if err := baseAccount.Validate(); err != nil {
					return fmt.Errorf("failed to validate non-airdrop account: %w", err)
				}
				accs = append(accs, baseAccount)
				liquidBalances = append(liquidBalances, banktypes.Balance{
					Address: address.String(),
					Coins:   sdk.NewCoins(sdk.NewCoin("uinto", math.NewInt(info.Amount))),
				})
				fmt.Printf("Non-airdrop account added: %s (%s) with %d uinto\n", info.Address, info.Name, info.Amount)
			}

			// Totals check
			totalLiquid := math.NewInt(0)
			for _, balance := range liquidBalances {
				totalLiquid = totalLiquid.Add(balance.Coins.AmountOf("uinto"))
			}
			fmt.Printf("Total Airdrop: %s uinto\n", airdropAmount)
			fmt.Printf("Total Liquid: %s uinto\n", totalLiquid)
			fmt.Printf("Total Claimable: %s uinto\n", claimModuleBalance)
			fmt.Printf("Sum of Distributed Airdrop Shares: %s uinto\n", totalDistributed)

			// After handling remainder, totalDistributed + remainder should equal airdropAmount
			if !totalDistributed.Add(remainder).Equal(airdropAmount) {
				return fmt.Errorf("mismatch: totalDistributed + remainder (%s) ≠ airdropAmount (%s)", totalDistributed.Add(remainder), airdropAmount)
			}

			// Update genesis state
			authGenState.Accounts, err = authtypes.PackAccounts(accs)
			if err != nil {
				return fmt.Errorf("failed to PackAccounts: %w", err)
			}
			authGenStateBz, err := clientCtx.Codec.MarshalJSON(&authGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}
			appState[authtypes.ModuleName] = authGenStateBz

			bankGenState.Balances = banktypes.SanitizeGenesisBalances(liquidBalances)
			bankGenStateBz, err := clientCtx.Codec.MarshalJSON(bankGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal bank genesis state: %w", err)
			}
			appState[banktypes.ModuleName] = bankGenStateBz

			claimGenState := claimtypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
			claimGenState.ModuleAccountBalance = sdk.NewCoin("uinto", claimModuleBalance)
			claimGenState.ClaimRecords = claimRecords
			claimGenStateBz, err := clientCtx.Codec.MarshalJSON(claimGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal claim genesis state: %w", err)
			}
			appState[claimtypes.ModuleName] = claimGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}
			genDoc.AppState = appStateJSON

			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().Int64("airdrop-amount", 0, "Total amount to distribute in the airdrop")
	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func ConvertBech32(address string) (string, error) {
	_, bz, err := bech32.DecodeAndConvert(address)
	if err != nil {
		panic(err)
		//return "", err
	}

	bech32Addr, err := bech32.ConvertAndEncode("into", bz)
	if err != nil {
		panic(err)
	}
	return bech32Addr, err
}
