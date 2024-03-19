package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	"github.com/spf13/cobra"
	claimtypes "github.com/trstlabs/intento/x/claim/types"
)

// change these before using fn
const (
	MinNumCosmosTxs = 100 //min num cosmos txs for given address to counter bots
	//additonal rules may be added
	MaxIntoPerSupportedConnection = 100_000_000_000_000 //uinto
)

type TestnetSnapshot struct {
	TotalIntoAirdropAmount sdkmath.Int                       `json:"total_into_amount"`
	NumAccounts            uint64                            `json:"num_accounts"`
	Accounts               map[string]TestnetSnapshotAccount `json:"accounts"`
}
type TestnetSnapshotAccount struct {
	IntoAddress                string            `json:"into_address"`
	TestnetParticpationPercent sdkmath.LegacyDec `json:"testnet_participation_percent"`
	IntoTotalClaimable         sdkmath.Int       `json:"into_total_claimable"`
}

type Transaction struct {
	Time          time.Time `json:"time"`
	Chain         string    `json:"chain"`
	Address       string    `json:"address"`
	DiscordID     string    `json:"discord_id"`
	Amount        string    `json:"amount"`
	TxHash        string    `json:"tx_hash"`
	FaucetBalance string    `json:"faucet_balance"`
	NumCosmosTxs  int64     `json:"num_cosmos_txs"`
}

type ActionIbcUsage struct {
	Address string      `json:"address"`
	Txs     []ActionAck `json:"actions"`
}

type ActionAck struct {
	Coin struct {
		Amount string `json:"amount"`
		Denom  string `json:"denom"`
	} `json:"coin"`
	ConnectionID string `json:"connection_id"`
}

type SupportedConnectionToken struct {
	ConnectionID            string
	IbcDenom                string
	Allocation              int64
	MaxTokenAmountForAction int64
}

var supportedList = []SupportedConnectionToken{{ConnectionID: "connection-0", IbcDenom: "ibc/BE4E6F930E604F5E410DCA660E8F4DB6F2BDE1F8E4730CAE2C4B15982EFB42B8", Allocation: 40_000_000_000_000, MaxTokenAmountForAction: 5_000_000_000}, {ConnectionID: "connection-1", IbcDenom: "ibc/8E2FEFCBD754FA3C97411F0126B9EC76191BAA1B3959CB73CECF396A4037BBF0", Allocation: 40_000_000_000_000, MaxTokenAmountForAction: 5_000_000_000}}

func ExportTestnetSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-testnet-snapshot [usage-file] [faucet-file] [output-file]",
		Short: "Export a testnet snapshot from provided exports",
		Long: `Export a  testnet snapshot from provided usage file and faucet details
Example:
intentod export-testnet-snapshot ~/into/analytics.json ~/cosmos-discord-faucet/transactions.json ./testnet-snapshot.json
`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			usageFile := args[0]
			faucetFile := args[1]
			snapshotOutput := args[2]

			faucetFileBytes, err := os.ReadFile(faucetFile)
			if err != nil {
				return err
			}

			var faucetTransactions []Transaction
			err = json.Unmarshal(faucetFileBytes, &faucetTransactions)
			if err != nil {
				return err
			}

			autoTxUsage := make([]ActionIbcUsage, 0)
			err = json.Unmarshal([]byte(usageFile), &autoTxUsage)
			if err != nil {
				return err
			}

			// make a new []ActionIbcUsage{} called autoTxUsageByUser and put in only
			// the address that matches a unique discordID from faucetTransactions based
			// on the highest amount of NumCosmosTxs

			var autoTxUsageByUser []ActionIbcUsage
			uniqueDiscordIDs := make(map[string]bool)

			// iterate through faucetTransactions and get the unique Discord IDs
			// and the faucet transaction with the highest NumCosmosTxs for each ID
			for _, tx := range faucetTransactions {
				if _, ok := uniqueDiscordIDs[tx.DiscordID]; !ok {
					// if this DiscordID is not in uniqueDiscordIDs, add it
					uniqueDiscordIDs[tx.DiscordID] = true

					// find the faucet transaction with the highest NumCosmosTxs for this DiscordID
					highestNumCosmosTxs := int64(0)
					var highestNumCosmosTxsTx Transaction
					for _, tx2 := range faucetTransactions {
						if tx2.DiscordID == tx.DiscordID && tx2.NumCosmosTxs > highestNumCosmosTxs {
							highestNumCosmosTxs = tx2.NumCosmosTxs
							highestNumCosmosTxsTx = tx2
						}
					}

					// create a new ActionIbcUsage with the address and Txs from the highest NumCosmosTxs faucet transaction
					var autoTxIbcUsage ActionIbcUsage
					autoTxIbcUsage.Address = highestNumCosmosTxsTx.Address
					autoTxIbcUsage.Txs = []ActionAck{}
					autoTxUsageByUser = append(autoTxUsageByUser, autoTxIbcUsage)
				}
			}

			//2. get total amount of IbcDenom for each Tx in Txs

			totalIbcDenomAmount := make(map[string]int64)
			for _, usage := range autoTxUsageByUser {
				for _, tx := range usage.Txs {
					for _, supported := range supportedList {
						if tx.Coin.Denom == supported.IbcDenom && tx.ConnectionID == supported.ConnectionID {
							totalIbcDenomAmount[tx.Coin.Denom] += sdk.MustNewDecFromStr(tx.Coin.Amount).RoundInt64()
							break
						}
					}
				}
			}

			// 3. remove entries from autoTxUsageByUser where the coin amount is higher than MaxIntoPerSupportedConnection
			for i := 0; i < len(autoTxUsageByUser); i++ {
				usage := &autoTxUsageByUser[i]
				for j := 0; j < len(usage.Txs); j++ {
					tx := &usage.Txs[j]
					coinAmount, err := strconv.ParseInt(tx.Coin.Amount, 10, 64)
					if err != nil {
						return err
					}
					if coinAmount > MaxIntoPerSupportedConnection {
						tx.Coin.Amount = fmt.Sprint(MaxIntoPerSupportedConnection) //append(tx.Coin[:k], tx.Coin[k+1:]...)
					}

					// If the usage has no more entries, remove it from autoTxUsageByUser
					if len(usage.Txs) == 0 {
						autoTxUsageByUser = append(autoTxUsageByUser[:i], autoTxUsageByUser[i+1:]...)
						i--
						break
					}
				}
			}

			// 4. For each entry in `supportedList`, calculate `ParticipationPercentage` and `IntoTotalClaimable` for each `autoTxUsageByUser` address based on the total amount of supported tokens. Then, create a `TestnetSnapshotAccount` for each combination of `supported.ConnectionID` and `usage.Address`, and add it to the `snapshot.Accounts` map.
			// Finally, update the `snapshot.TotalIntoAirdropAmount` variable for each supported token.
			//
			// To do this, we first initialize `snapshot` and `trstTotalClaimable`, and loop over the `supportedList` to calculate the total IBC amount for all supported tokens. Then, we loop over the `autoTxUsageByUser` slice and calculate the `participationPercentage` and `trstTotalClaimable` for each `autoTxUsageByUser` address based on the total IBC amount.
			// We also update the `TestnetSnapshotAccount` for each address with the calculated values.
			// Finally, we add the updated `TestnetSnapshotAccount` to the `snapshot.Accounts` map, and update the `snapshot.TotalIntoAirdropAmount` variable for each supported token.
			var snapshot TestnetSnapshot
			snapshot.TotalIntoAirdropAmount = sdk.ZeroInt()
			snapshot.Accounts = make(map[string]TestnetSnapshotAccount)

			var trstTotalClaimable sdkmath.Int
			totalAirdropAmount := int64(0)
			for _, supported := range supportedList {
				totalAirdropAmount += supported.Allocation //
			}

			for _, usage := range autoTxUsageByUser {
				var account TestnetSnapshotAccount
				account.IntoAddress = usage.Address
				for _, tx := range usage.Txs {
					for _, supported := range supportedList {
						if tx.Coin.Denom == supported.IbcDenom && tx.ConnectionID == supported.ConnectionID {
							coinAmount, _ := strconv.ParseInt(tx.Coin.Amount, 10, 64)
							if coinAmount > supported.MaxTokenAmountForAction {
								coinAmount = supported.MaxTokenAmountForAction
							}
							//for example if tx.Coin.Amount is 5, out of 50 (totalIbcDenomAmount) is 10% and supported.Allocation is 100, trstAmount should be 10
							trstAmount := sdk.NewInt(coinAmount).Quo(sdk.NewInt(totalIbcDenomAmount[supported.IbcDenom])).Mul(sdk.NewInt(supported.Allocation))
							trstTotalClaimable = trstTotalClaimable.Add(trstAmount.Quo(sdk.NewInt(10)))
							participationPercentage := sdk.NewDecFromInt(trstAmount).QuoInt(sdk.NewInt(totalAirdropAmount))

							account.TestnetParticpationPercent = account.TestnetParticpationPercent.Add(participationPercentage)
							account.IntoTotalClaimable = account.IntoTotalClaimable.Add(trstAmount)
							snapshot.TotalIntoAirdropAmount = snapshot.TotalIntoAirdropAmount.Add(trstAmount)
						}
					}
				}
				snapshot.Accounts[usage.Address] = account
			}

			if trstTotalClaimable != sdk.NewInt(totalAirdropAmount) {
				fmt.Printf("total airdrop amount: %v\n", totalAirdropAmount)
				fmt.Printf("total airdrop tokens: %v\n", trstTotalClaimable)
				return fmt.Errorf("total claimable is different from set total %v", trstTotalClaimable.Int64())
			}

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

func ImportTestnetSnapshotCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import-testnet-snapshot [input-snapshot-file]",
		Short: "Import testnet accounts from the testnet snapshot.json",
		Long: `Import testnet accounts from the testnet claimable snapshot.json
		INTO is illiquid, a small amount INTO of of the airdropped coins is liquid in accounts.
		The illiqud INTO is placed in the claims module, to be claimed after performing several actions.
		Example:
		intentod import-testnet-snapshot ~/into/snapshot.json
		- Check input genesis:
			file is at ~/.intentod/config/genesis.json
`,
		Args: cobra.ExactArgs(1),
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
				return fmt.Errorf("failed to get accounts from any: %w", err)
			}

			// Read snapshot file
			snapshotInput := args[0]
			snapshotJSON, err := os.Open(snapshotInput)
			if err != nil {
				return err
			}
			defer snapshotJSON.Close()

			var snapshot TestnetSnapshot
			decoder := json.NewDecoder(snapshotJSON)
			if err := decoder.Decode(&snapshot); err != nil {
				return err
			}

			// get genesis params
			genesisParams := MainnetGenesisParams()

			fmt.Printf("total snapshot accounts: %v\n", len(snapshot.Accounts))

			fmt.Printf("snapshot total supply: %s\n", snapshot.TotalIntoAirdropAmount)

			// donating the remainder to the rainforrest foundation https://rainforestfoundation.org/
			if !genesisParams.AirdropSupply.Sub(snapshot.TotalIntoAirdropAmount).IsNil() {
				return fmt.Errorf("snapshot supply does not match genesis")
			}

			bankGenState := banktypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
			liquidBalances := bankGenState.Balances

			claimRecords := []claimtypes.ClaimRecord{}
			claimModuleAccountBalance := sdk.NewInt(0)

			// for each account in the snapshot
			for _, acc := range snapshot.Accounts {
				// read address from snapshot
				//we cannot have duplicate accounts

				address, err := sdk.AccAddressFromBech32(acc.IntoAddress)
				if err != nil {
					return fmt.Errorf("into acc from snapshot invalid: %s", err)
				}

				// initial liquid amounts
				liquidCoins := sdk.NewCoins(sdk.NewCoin(genesisParams.NativeCoinMetadatas[0].Base, sdk.NewInt(696969)))

				// Add the new account to the set of genesis accounts
				baseAccount := authtypes.NewBaseAccount(address, nil, 0, 0)
				if err := baseAccount.Validate(); err != nil {
					return fmt.Errorf("failed to validate new genesis account: %w", err)
				}
				accs = append(accs, baseAccount)

				// claimable balances
				claimableAmount := acc.IntoTotalClaimable.Sub(sdk.NewInt(696969))
				if acc.IntoTotalClaimable.Sub(sdk.NewInt(696969)).IsNegative() {
					claimableAmount = sdk.ZeroInt()
					liquidCoins = sdk.NewCoins(sdk.NewCoin("uinto", acc.IntoTotalClaimable))
				}

				liquidBalances = append(liquidBalances, banktypes.Balance{
					Address: address.String(),
					Coins:   liquidCoins,
				})

				status := claimtypes.Status{ActionCompleted: false, VestingPeriodCompleted: []bool{false, false, false, false}, VestingPeriodClaimed: []bool{false, false, false, false}}
				claimRecords = append(claimRecords, claimtypes.ClaimRecord{
					Address:                address.String(),
					InitialClaimableAmount: sdk.NewCoins(sdk.NewCoin(genesisParams.NativeCoinMetadatas[0].Base, claimableAmount)),
					Status:                 []claimtypes.Status{status, status, status, status},
				})

				claimModuleAccountBalance = claimModuleAccountBalance.Add(claimableAmount)

			}
			var totalAirdrop sdk.Coins
			for _, balance := range liquidBalances {
				totalAirdrop = totalAirdrop.Add(balance.Coins...)
			}

			fmt.Printf("Total liquid coins %s \n", totalAirdrop.AmountOf("uinto"))
			fmt.Printf("Total airdrop coins %s \n", claimModuleAccountBalance.Add(totalAirdrop.AmountOf("uinto")))

			var total sdk.Coins
			for _, balance := range liquidBalances {
				total = total.Add(balance.Coins...)
			}
			fmt.Printf("Total non-airdrop %s \n", total.Sub(totalAirdrop[0]))

			fmt.Printf("Total balances %s \n", claimModuleAccountBalance.Add(totalAirdrop.AmountOf("uinto")).Add(total.Sub(totalAirdrop[0]).AmountOf("uinto")))

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
				return fmt.Errorf("error marshal ibc transfer: %w", err)
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
				return fmt.Errorf("error unmarshal ibc transfer: %w", err)
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
