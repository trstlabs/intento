package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"cosmossdk.io/math"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ccvconsumertypes "github.com/cosmos/interchain-security/v6/x/ccv/consumer/types"
	"github.com/spf13/cobra"
	"github.com/trstlabs/intento/app"
	alloctypes "github.com/trstlabs/intento/x/alloc/types"
	claimtypes "github.com/trstlabs/intento/x/claim/types"
	intenttypes "github.com/trstlabs/intento/x/intent/types"
	minttypes "github.com/trstlabs/intento/x/mint/types"
)

const (
	HumanCoinUnit = "into"
	BaseCoinUnit  = "uinto"
	IntoExponent  = 6

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = "into"
)

type GenesisParams struct {
	AirdropSupply            math.Int
	StrategicReserveAccounts []banktypes.Balance
	DistributedAccounts      []banktypes.Balance
	ConsensusParams          *tmtypes.ConsensusParams

	GenesisTime         time.Time
	NativeCoinMetadatas []banktypes.Metadata

	StakingParams      stakingtypes.Params
	DistributionParams distributiontypes.Params
	GovParams          govtypesv1.Params

	CrisisConstantFee sdk.Coin

	SlashingParams slashingtypes.Params
	AllocParams    alloctypes.Params
	ClaimParams    claimtypes.Params
	MintParams     minttypes.Params

	IcaParams            icatypes.Params
	IntentParams         intenttypes.Params
	WasmParams           wasmtypes.Params
	ConsumerGenesisState ccvconsumertypes.GenesisState
}

func PrepareGenesisCmd(defaultNodeHome string, mbm module.BasicManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prepare-genesis [network] [chainID]",
		Short: "Prepare a genesis file with initial setup",
		Long: `Prepare a genesis file with initial setup.
Examples include:
	- Setting module initial params
	- Setting denom metadata
Example:
	intentod prepare-genesis mainnet intento-1
	- Check input genesis:
		file is at ~/.intentod/config/genesis.json
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			cdc := clientCtx.Codec
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			// read genesis file
			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			// get genesis params
			var genesisParams GenesisParams
			network := args[0]
			switch network {
			case "testnet":
				genesisParams = TestnetGenesisParams()
			case "mainnet":
				genesisParams = MainnetGenesisParams()
			default:
				return fmt.Errorf("please choose 'mainnet' or 'testnet'")
			}

			// get genesis params
			chainID := args[1]

			// run Prepare Genesis
			appState, err = PrepareGenesis(clientCtx, appState, genesisParams)
			if err != nil {
				return fmt.Errorf("failed to prepare genesis: %w", err)
			}
			genDoc.GenesisTime = genesisParams.GenesisTime
			genDoc.ChainID = chainID
			// genDoc.ConsensusParams = genesisParams.ConsensusParams

			// validate genesis state
			if err = mbm.ValidateGenesis(cdc, clientCtx.TxConfig, appState); err != nil {
				return fmt.Errorf("error validating genesis file: %s", err.Error())
			}

			// save genesis
			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON

			//fmt.Printf("%v \n", string(appStateJSON))
			err = genutil.ExportGenesisFile(genDoc, genFile)
			return err
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func PrepareGenesis(
	clientCtx client.Context,
	appState map[string]json.RawMessage,
	genesisParams GenesisParams,
) (map[string]json.RawMessage, error) {
	cdc := clientCtx.Codec
	// ---
	// bank module genesis

	bankGenState := banktypes.DefaultGenesisState()
	bankGenState.Params.DefaultSendEnabled = true
	bankGenStateBz, err := cdc.MarshalJSON(bankGenState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal bank genesis state: %w", err)
	}
	appState[banktypes.ModuleName] = bankGenStateBz

	// IBC transfer module genesis
	ibcGenState := ibctransfertypes.DefaultGenesisState()
	ibcGenState.Params.SendEnabled = true
	ibcGenState.Params.ReceiveEnabled = true
	ibcGenStateBz, err := cdc.MarshalJSON(ibcGenState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal IBC transfer genesis state: %w", err)
	}
	appState[ibctransfertypes.ModuleName] = ibcGenStateBz

	// mint module genesis

	mintGenState := minttypes.DefaultGenesisState()
	mintGenState.Params = genesisParams.MintParams

	mintGenStateBz, err := cdc.MarshalJSON(mintGenState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal mint genesis state: %w", err)
	}
	appState[minttypes.ModuleName] = mintGenStateBz

	// staking module genesis
	stakingGenState := stakingtypes.GetGenesisStateFromAppState(cdc, appState)
	stakingGenState.Params = genesisParams.StakingParams
	stakingGenStateBz, err := cdc.MarshalJSON(stakingGenState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal staking genesis state: %w", err)
	}
	appState[stakingtypes.ModuleName] = stakingGenStateBz

	// distribution module genesis
	distributionGenState := distributiontypes.DefaultGenesisState()
	distributionGenState.Params = genesisParams.DistributionParams
	distributionGenStateBz, err := cdc.MarshalJSON(distributionGenState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal distribution genesis state: %w", err)
	}
	appState[distributiontypes.ModuleName] = distributionGenStateBz

	// // gov module genesis
	govGenState := govtypesv1.DefaultGenesisState()
	govGenStateBz, err := cdc.MarshalJSON(govGenState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal gov genesis state: %w", err)
	}
	appState[govtypes.ModuleName] = govGenStateBz

	// crisis module genesis
	crisisGenState := crisistypes.DefaultGenesisState()
	crisisGenState.ConstantFee = genesisParams.CrisisConstantFee
	crisisGenStateBz, err := cdc.MarshalJSON(crisisGenState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal crisis genesis state: %w", err)
	}
	appState[crisistypes.ModuleName] = crisisGenStateBz

	// slashing module genesis
	slashingGenState := slashingtypes.DefaultGenesisState()
	slashingGenState.Params = genesisParams.SlashingParams
	slashingGenState.Params.SignedBlocksWindow = 48000                              // equal to 16 hours @ 1.3s block time
	slashingGenState.Params.SlashFractionDowntime = math.LegacyNewDecWithPrec(0, 4) // 0% liveness slashing, as missing out rewards is enough penalty
	slashingGenStateBz, err := cdc.MarshalJSON(slashingGenState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal slashing genesis state: %w", err)
	}
	appState[slashingtypes.ModuleName] = slashingGenStateBz

	// claim module genesis
	claimGenState := claimtypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
	claimGenState.Params = genesisParams.ClaimParams
	claimGenState.Params.AirdropStartTime = genesisParams.GenesisTime.Add(time.Hour * 24 * 365) // will be adjusted by governance
	claimGenState.Params.DurationUntilDecay = time.Hour * 24 * 28                               // 28 days or 4 weeks
	claimGenState.Params.DurationOfDecay = time.Hour * 24 * 56                                  // 56 days or 8 weeks
	claimGenState.Params.DurationVestingPeriods = []time.Duration{time.Minute, time.Minute * 2, time.Minute * 5, time.Minute}
	claimGenStateBz, err := cdc.MarshalJSON(claimGenState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal claim genesis state: %w", err)
	}
	appState[claimtypes.ModuleName] = claimGenStateBz

	// alloc module genesis
	allocGenState := alloctypes.GetGenesisStateFromAppState(cdc, appState)
	allocGenState.Params = genesisParams.AllocParams
	allocGenStateBz, err := cdc.MarshalJSON(allocGenState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal alloc genesis state: %w", err)
	}
	appState[alloctypes.ModuleName] = allocGenStateBz

	// Intent module genesis
	intentGenState := intenttypes.DefaultGenesis()
	intentGenState.Params = genesisParams.IntentParams
	intentGenStateBz, err := cdc.MarshalJSON(intentGenState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal autoIbcTx genesis state: %w", err)
	}
	appState[intenttypes.ModuleName] = intentGenStateBz

	// return appState and genDoc
	return appState, nil
}

func MainnetGenesisParams() GenesisParams {
	genParams := GenesisParams{}

	genParams.AirdropSupply = math.NewInt(90_000_000_000_000)                       // (90M INTO)
	genParams.GenesisTime = time.Date(2025, time.August, 27, 14, 0, 0, 0, time.UTC) // 2025 27 Aug - 14:00 UTC
	genParams.NativeCoinMetadatas = []banktypes.Metadata{
		{
			Description: "The native token of INTO",
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    BaseCoinUnit,
					Exponent: 0,
					Aliases:  nil,
				},
				{
					Denom:    HumanCoinUnit,
					Exponent: IntoExponent,
					Aliases:  nil,
				},
			},
			Base:    BaseCoinUnit,
			Display: HumanCoinUnit,
		},
	}

	// alloc
	genParams.AllocParams = alloctypes.DefaultParams()
	genParams.AllocParams.DistributionProportions = alloctypes.DistributionProportions{
		CommunityPool:     math.LegacyNewDecWithPrec(70, 2),
		DeveloperRewards:  math.LegacyNewDecWithPrec(0, 2),
		RelayerIncentives: math.LegacyNewDecWithPrec(5, 2),
	}
	genParams.AllocParams.WeightedDeveloperRewardsReceivers = []alloctypes.WeightedAddress{}

	// mint
	genParams.MintParams = minttypes.DefaultParams()
	genParams.MintParams.MintDenom = BaseCoinUnit
	genParams.MintParams.StartTime = genParams.GenesisTime.AddDate(0, 6, 0)
	genParams.MintParams.InitialAnnualProvisions = math.LegacyNewDec(150_000_000_000_000)

	genParams.MintParams.ReductionFactor = math.LegacyNewDec(3).QuoInt64(4)
	genParams.MintParams.BlocksPerYear = uint64(24261538) //assuming 1.3s average block times, param to be updated periodically

	// staking
	genParams.StakingParams = stakingtypes.DefaultParams()
	genParams.StakingParams.UnbondingTime = time.Hour * 24 * 14 //2 weeks
	genParams.StakingParams.MaxValidators = 50
	genParams.StakingParams.BondDenom = genParams.NativeCoinMetadatas[0].Base

	// distr
	genParams.DistributionParams = distributiontypes.DefaultParams()
	genParams.DistributionParams.CommunityTax = math.LegacyNewDecWithPrec(5, 2)
	genParams.DistributionParams.WithdrawAddrEnabled = true

	// gov
	genParams.GovParams = govtypesv1.DefaultParams()
	maxDepositPeriod := time.Hour * 24 * 14 // 2 weeks
	genParams.GovParams.MaxDepositPeriod = &maxDepositPeriod
	genParams.GovParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(
		genParams.NativeCoinMetadatas[0].Base,
		math.NewInt(1_000_000_000),
	))
	genParams.GovParams.Quorum = "0.200000000000000000" // 20%
	votingPeriod := time.Hour * 24 * 3                  // 3 days
	genParams.GovParams.VotingPeriod = &votingPeriod

	// crisis
	genParams.CrisisConstantFee = sdk.NewCoin(
		genParams.NativeCoinMetadatas[0].Base,
		math.NewInt(100_000_000_000),
	)
	// slash
	genParams.SlashingParams = slashingtypes.DefaultParams()
	genParams.SlashingParams.SignedBlocksWindow = int64(25000)                         // ~41 hr at 6 second blocks
	genParams.SlashingParams.MinSignedPerWindow = math.LegacyNewDecWithPrec(5, 2)      // 5% minimum liveness
	genParams.SlashingParams.DowntimeJailDuration = time.Minute                        // 1 minute jail period
	genParams.SlashingParams.SlashFractionDoubleSign = math.LegacyNewDecWithPrec(5, 2) // 5% double sign slashing
	genParams.SlashingParams.SlashFractionDowntime = math.LegacyNewDecWithPrec(1, 4)   // 0.01% liveness slashing               // 0% liveness slashing

	//wasm included for future integrations
	genParams.WasmParams = wasmtypes.DefaultParams()
	genParams.WasmParams.CodeUploadAccess = wasmtypes.AllowNobody
	genParams.WasmParams.InstantiateDefaultPermission = wasmtypes.AccessTypeNobody

	//intent
	genParams.IntentParams = intenttypes.DefaultParams()
	genParams.IntentParams.MaxFlowDuration = time.Hour * 24 * 366 * 3
	genParams.IntentParams.MinFlowDuration = time.Second * 600 // 10 minutes
	genParams.IntentParams.MinFlowInterval = time.Second * 300 // 5 minutes
	genParams.IntentParams.FlowFundsCommission = 2
	genParams.IntentParams.BurnFeePerMsg = 10_000
	genParams.IntentParams.FlowFlexFeeMul = 5
	genParams.IntentParams.GasFeeCoins = sdk.Coins(sdk.NewCoins(sdk.NewCoin(BaseCoinUnit, math.NewInt(80)), sdk.NewCoin("ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9", math.NewInt(20))))
	genParams.IntentParams.ConnectionRewards = []*intenttypes.ConnectionRelayerReward{
		{
			ConnectionID:   "connection-0",
			RelayerRewards: []int64{20_000, 30_000, 40_000, 50_000},
		},
	}

	//claim
	genParams.ClaimParams = claimtypes.DefaultGenesis().Params
	genParams.ClaimParams.AirdropStartTime = genParams.GenesisTime.Add(time.Hour * 24 * 365) // 1 year (will be changed through gov)
	genParams.ClaimParams.DurationUntilDecay = time.Hour * 24 * 60                           // 60 days = ~2 months
	genParams.ClaimParams.DurationOfDecay = time.Hour * 24 * 120                             // 120 days = ~4 months
	genParams.ClaimParams.ClaimDenom = genParams.NativeCoinMetadatas[0].Base
	genParams.ClaimParams.DurationVestingPeriods = []time.Duration{time.Hour * 24 * 7, time.Hour * 24 * 7 * 6, time.Hour * 24 * 3, time.Hour * 24 * 7 * 4}

	//consensus
	genParams.ConsensusParams = tmtypes.DefaultConsensusParams()
	genParams.ConsensusParams.Block.MaxBytes = 22020096
	genParams.ConsensusParams.Block.MaxGas = -1
	genParams.ConsensusParams.Evidence.MaxAgeDuration = time.Second * 120
	genParams.ConsensusParams.Evidence.MaxAgeNumBlocks = int64(genParams.StakingParams.UnbondingTime.Seconds()) / 3
	genParams.ConsensusParams.Version.App = 1
	genParams.DistributedAccounts = []banktypes.Balance{}

	consumerGenesisState := app.CreateMinimalConsumerTestGenesis()
	genParams.ConsumerGenesisState = *consumerGenesisState

	//interchain accounts host
	genParams.IcaParams.AllowMessages = []string{"*"} // allow all msgs
	genParams.IcaParams.HostEnabled = true
	return genParams
}

func TestnetGenesisParams() GenesisParams {

	genParams := MainnetGenesisParams()

	genParams.GenesisTime = time.Now()
	genParams.MintParams.StartTime = time.Now()
	genParams.StakingParams.UnbondingTime = time.Hour * 24 * 3 // 3 days

	//gov
	genParams.GovParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(
		genParams.NativeCoinMetadatas[0].Base,
		math.NewInt(1_000_000), // 1 INTO
	))
	genParams.GovParams.Quorum = "0.100000000000000000" // 10%
	votingPeriod := time.Minute
	genParams.GovParams.VotingPeriod = &votingPeriod

	//flow
	genParams.IntentParams.MinFlowDuration = time.Second * 40
	genParams.IntentParams.MinFlowInterval = time.Second * 40
	//genParams.IntentParams.MaxFlowDuration = time.Hour * 8

	//slasing window
	genParams.SlashingParams.SignedBlocksWindow = 10000 //shorter for testnet

	genParams.ClaimParams.AirdropStartTime = genParams.GenesisTime
	genParams.ClaimParams.DurationUntilDecay = time.Hour * 24 * 5 // 5 days
	genParams.ClaimParams.DurationOfDecay = time.Hour * 24 * 5    // 5 days
	genParams.ClaimParams.DurationVestingPeriods = []time.Duration{time.Minute, time.Minute * 2, time.Minute * 5, time.Minute}

	genParams.WasmParams.CodeUploadAccess = wasmtypes.AllowEverybody
	genParams.WasmParams.InstantiateDefaultPermission = wasmtypes.AccessTypeEverybody
	return genParams
}
