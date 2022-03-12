package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v2/modules/apps/transfer/types"
	"github.com/spf13/cobra"
	claimtypes "github.com/trstlabs/trst/x/claim/types"
)

const (
	flagVestingStart = "vesting-start-time"
	flagVestingEnd   = "vesting-end-time"
	flagVestingAmt   = "vesting-amount"
)

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

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			var kr keyring.Keyring
			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				keyringBackend, _ := cmd.Flags().GetString(flags.FlagKeyringBackend)

				if keyringBackend != "" && clientCtx.Keyring == nil {
					var err error
					kr, err = keyring.New(sdk.KeyringServiceName(), keyringBackend, clientCtx.HomeDir, inBuf)
					if err != nil {
						return err
					}
				} else {
					kr = clientCtx.Keyring
				}

				k, err := kr.Key(args[0])
				if err != nil {
					return fmt.Errorf("failed to get address from Keyring: %w", err)
				}

				addr = k.GetAddress()
				if err != nil {
					return err
				}
			}

			coins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return fmt.Errorf("failed to parse coins: %w", err)
			}

			vestingStart, _ := cmd.Flags().GetInt64(flagVestingStart)
			vestingEnd, _ := cmd.Flags().GetInt64(flagVestingEnd)
			vestingAmtStr, _ := cmd.Flags().GetString(flagVestingAmt)

			vestingAmt, err := sdk.ParseCoinsNormalized(vestingAmtStr)
			if err != nil {
				return fmt.Errorf("failed to parse vesting amount: %w", err)
			}

			// create concrete account type based on input parameters
			var genAccount authtypes.GenesisAccount

			balances := banktypes.Balance{Address: addr.String(), Coins: coins.Sort()}
			baseAccount := authtypes.NewBaseAccount(addr, nil, 0, 0)

			if !vestingAmt.IsZero() {
				baseVestingAccount := authvesting.NewBaseVestingAccount(baseAccount, vestingAmt.Sort(), vestingEnd)

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

			authGenState := authtypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)

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

			authGenStateBz, err := clientCtx.Codec.MarshalJSON(&authGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}

			appState[authtypes.ModuleName] = authGenStateBz

			bankGenState := banktypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
			bankGenState.Balances = append(bankGenState.Balances, balances)
			bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)
			//bankGenState.Supply = bankGenState.Supply.Add(balances.Coins...)

			bankGenStateBz, err := clientCtx.Codec.MarshalJSON(bankGenState)
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

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test)")
	cmd.Flags().String(flagVestingAmt, "", "amount of coins for vesting accounts")
	cmd.Flags().Int64(flagVestingStart, 0, "schedule start time (unix epoch) for vesting accounts")
	cmd.Flags().Int64(flagVestingEnd, 0, "schedule end time (unix epoch) for vesting accounts")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
func ImportGenesisAccountsFromSnapshotCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import-genesis-accounts-from-snapshot [input-snapshot-file] [input-non-airdrop-accounts-file]",
		Short: "Import genesis accounts from fairdrop snapshot.json and an non-airdrop-accounts.json",
		Long: `Import genesis accounts from fairdrop snapshot.json
		1TRST of the airdrop coins to be received is liquid in accounts.
		The remaining is placed in the claims module, to be claimed.
		Must also pass in an snapshot.json file to airdrop genesis TRST coins
		Example:
		trstd import-genesis-accounts-from-snapshot ../snapshot.json ../non-airdrop-accounts.json
		- Check input genesis:
			file is at ~/.trstd/config/genesis.json
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			// aminoCodec := clientCtx.LegacyAmino.Amino

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
				return fmt.Errorf("failed to get accounts from any: %w", err)
			}

			// Read snapshot file
			snapshotInput := args[0]
			snapshotJSON, err := os.Open(snapshotInput)
			if err != nil {
				return err
			}
			defer snapshotJSON.Close()
			byteValue, _ := ioutil.ReadAll(snapshotJSON)
			snapshot := Snapshot{}

			json.Unmarshal(byteValue, &snapshot)
			if err != nil {
				return err
			}

			// Read ions file
			trstInput := args[1]
			trstJSON, err := os.Open(trstInput)
			if err != nil {
				return err
			}
			defer trstJSON.Close()
			byteValue2, _ := ioutil.ReadAll(trstJSON)
			var trstAmts map[string]int64
			json.Unmarshal(byteValue2, &trstAmts)
			if err != nil {
				return err
			}

			// get genesis params
			genesisParams := MainnetGenesisParams()
			nonAirdropAccs := make(map[string]sdk.Coins)

			for _, acc := range genesisParams.DistributedAccounts {
				nonAirdropAccs[acc.Address] = acc.GetCoins()
			}

			for addr, amt := range trstAmts {
				// set atom bech32 prefixes
				bech32Addr, err := ConvertBech32(addr)
				if err != nil {
					return err
				}

				address, err := sdk.AccAddressFromBech32(bech32Addr)
				if err != nil {
					return err
				}

				if val, ok := nonAirdropAccs[address.String()]; ok {
					nonAirdropAccs[address.String()] = val.Add(sdk.NewCoin("utrst", sdk.NewInt(amt).MulRaw(1_000_000)))
				} else {
					nonAirdropAccs[address.String()] = sdk.NewCoins(sdk.NewCoin("utrst", sdk.NewInt(amt).MulRaw(1_000_000)))
				}
			}

			// figure out normalizationFactor to normalize snapshot balances to desired airdrop supply
			normalizationFactor := genesisParams.AirdropSupply.ToDec().QuoInt(snapshot.TotalTrstAirdropAmount)

			fmt.Printf("total snapshot accounts: %s\n", len(snapshot.Accounts))
			fmt.Printf("normalization factor: %s\n", normalizationFactor)
			fmt.Printf("snapshot total supply: %s\n", snapshot.TotalTrstAirdropAmount)

			// donating the remainder to the rainforrest foundation https://rainforestfoundation.org/
			toDonate := genesisParams.AirdropSupply.Sub(snapshot.TotalTrstAirdropAmount)
			rainForrestFoundationAddr, _ := sdk.AccAddressFromHex("17F318875145240F05259C65FCAB0E9C3DB92C0B")
			nonAirdropAccs[rainForrestFoundationAddr.String()] = sdk.NewCoins(sdk.NewCoin("utrst", toDonate))
			fmt.Printf("donated: %s \n", toDonate)

			bankGenState := banktypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
			liquidBalances := bankGenState.Balances

			claimRecords := []claimtypes.ClaimRecord{}
			claimModuleAccountBalance := sdk.NewInt(0)

			// for each account in the snapshot
			for _, acc := range snapshot.Accounts {
				// read address from snapshot
				//we cannot have duplicate accounts

				bech32Addr, err := ConvertBech32(acc.TokenAddress)
				if err != nil {
					return err
				}

				address, err := sdk.AccAddressFromBech32(bech32Addr)
				if err != nil {
					return err
				}
				//if !accs.Contains(address) {
				// initial liquid amounts
				// We consistently round down to the nearest utrst
				// get normalized trst balance for account
				normalizedTrstBalance := acc.TrstBalance.ToDec().Mul(normalizationFactor)

				//liquidCoins := sdk.NewCoins(sdk.NewCoin(genesisParams.NativeCoinMetadatas[0].Base, acc.TrstBalance))
				//liquidAmount := normalizedTrstBalance.Mul(sdk.MustNewDecFromStr("0.2")).TruncateInt() // 20% of airdrop amount
				liquidCoins := sdk.NewCoins(sdk.NewCoin(genesisParams.NativeCoinMetadatas[0].Base, sdk.NewInt(696969)))

				if coins, ok := nonAirdropAccs[address.String()]; ok {
					liquidCoins = liquidCoins.Add(coins...)
					delete(nonAirdropAccs, address.String())
				}

				liquidBalances = append(liquidBalances, banktypes.Balance{
					Address: address.String(),
					Coins:   liquidCoins,
				})
				//supply = supply.Add(liquidCoins...)

				// Add the new account to the set of genesis accounts
				baseAccount := authtypes.NewBaseAccount(address, nil, 0, 0)
				if err := baseAccount.Validate(); err != nil {
					return fmt.Errorf("failed to validate new genesis account: %w", err)
				}
				accs = append(accs, baseAccount)

				// claimable balances
				claimableAmount := normalizedTrstBalance.TruncateInt().Sub(sdk.NewInt(696969)) //.Mul(sdk.MustNewDecFromStr("0.8")).TruncateInt()
				if normalizedTrstBalance.Sub(sdk.NewDec(696969)).IsNegative() {
					claimableAmount = sdk.ZeroInt()
					liquidCoins = sdk.NewCoins(sdk.NewCoin("utrst", normalizedTrstBalance.TruncateInt()))
				}
				claimRecords = append(claimRecords, claimtypes.ClaimRecord{
					Address:                address.String(),
					InitialClaimableAmount: sdk.NewCoins(sdk.NewCoin(genesisParams.NativeCoinMetadatas[0].Base, claimableAmount)),
					ActionCompleted:        []bool{false, false, false, false},
				})

				claimModuleAccountBalance = claimModuleAccountBalance.Add(claimableAmount)

				//} else {
				//	fmt.Printf(" new genesis account contains double account: %w", address)
				//	}
			}
			var totalAirdrop sdk.Coins
			for _, balance := range liquidBalances {
				totalAirdrop = totalAirdrop.Add(balance.Coins...)
			}
			fmt.Printf("Total liquid coins %s \n", totalAirdrop.AmountOf("utrst"))
			//fmt.Printf("Total airdrop coins approx %s", totalAirdrop.AmountOf("utrst").MulRaw(5))
			fmt.Printf("Total airdrop coins %s \n", claimModuleAccountBalance.Add(totalAirdrop.AmountOf("utrst")))
			//fmt.Println(supply.AmountOf("utrst"))

			// distribute remaining trst to accounts not in fairdrop
			for addr, coin := range nonAirdropAccs {
				// read address from snapshot
				address, err := sdk.AccAddressFromBech32(addr)
				if err != nil {
					return err
				}

				liquidBalances = append(liquidBalances, banktypes.Balance{
					Address: address.String(),
					Coins:   coin,
				})
				//supply = supply.Add(coin...)

				// Add the new account to the set of genesis accounts
				baseAccount := authtypes.NewBaseAccount(address, nil, 0, 0)
				if err := baseAccount.Validate(); err != nil {
					return fmt.Errorf("failed to validate new genesis account: %w", err)
				}
				accs = append(accs, baseAccount)
			}
			var total sdk.Coins
			for _, balance := range liquidBalances {
				total = total.Add(balance.Coins...)
			}
			fmt.Printf("Total non-airdrop %s \n", total.Sub(totalAirdrop).AmountOf("utrst"))

			fmt.Printf("Total balances %s \n", claimModuleAccountBalance.Add(totalAirdrop.AmountOf("utrst")).Add(total.Sub(totalAirdrop).AmountOf("utrst")))

			// auth module genesis
			accs = authtypes.SanitizeGenesisAccounts(accs)
			genAccs, err := authtypes.PackAccounts(accs)
			if err != nil {
				return fmt.Errorf("failed to convert accounts into any's: %w", err)
			}
			authGenState.Accounts = genAccs
			authGenStateBz, err := clientCtx.Codec.MarshalJSON(&authGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}
			appState[authtypes.ModuleName] = authGenStateBz

			// bank module genesis
			bankGenState.Balances = banktypes.SanitizeGenesisBalances(liquidBalances)
			//bankGenState.Supply = supply
			bankGenStateBz, err := clientCtx.Codec.MarshalJSON(bankGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal bank genesis state: %w", err)
			}
			appState[banktypes.ModuleName] = bankGenStateBz

			byteIBCTransfer, err := appState[ibctransfertypes.ModuleName].MarshalJSON()
			if err != nil {
				return fmt.Errorf("Error marshal ibc transfer: %w", err)
			}

			// claim module genesis
			claimGenState := claimtypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
			claimGenState.ModuleAccountBalance = sdk.NewCoin(genesisParams.NativeCoinMetadatas[0].Base, claimModuleAccountBalance)

			claimGenState.ClaimRecords = claimRecords
			claimGenStateBz, err := clientCtx.Codec.MarshalJSON(claimGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal claim genesis state: %w", err)
			}
			appState[claimtypes.ModuleName] = claimGenStateBz

			var ibcGenState ibctransfertypes.GenesisState
			err = ibctransfertypes.ModuleCdc.UnmarshalJSON(byteIBCTransfer, &ibcGenState)
			if err != nil {
				return fmt.Errorf("Error unmarshal ibc transfer: %w", err)
			}
			ibcGenState.Params = ibctransfertypes.NewParams(false, false)
			ibcGenStateBz, err := clientCtx.Codec.MarshalJSON(&ibcGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal ibc genesis state: %w", err)
			}
			appState[ibctransfertypes.ModuleName] = ibcGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}
			genDoc.AppState = appStateJSON

			err = genutil.ExportGenesisFile(genDoc, genFile)
			return err

		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func ConvertBech32(address string) (string, error) {
	_, bz, err := bech32.DecodeAndConvert(address)
	if err != nil {
		panic(err)
	}

	bech32Addr, err := bech32.ConvertAndEncode("trust", bz)
	if err != nil {
		panic(err)
	}
	return bech32Addr, err
}

func cosmosConvertBech32(address string) (string, error) {
	_, bz, err := bech32.DecodeAndConvert(address)
	if err != nil {
		panic(err)
	}

	bech32Addr, err := bech32.ConvertAndEncode("cosmos", bz)
	if err != nil {
		panic(err)
	}
	return bech32Addr, err
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
