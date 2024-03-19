package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"
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
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	"github.com/spf13/cobra"
	alloctypes "github.com/trstlabs/intento/x/alloc/types"
	claimtypes "github.com/trstlabs/intento/x/claim/types"
	intenttypes "github.com/trstlabs/intento/x/intent/types"

	// compute "github.com/trstlabs/intento/x/compute"
	minttypes "github.com/trstlabs/intento/x/mint/types"
	//	itemtypes "github.com/trstlabs/intento/x/item/types"
)

const (
	HumanCoinUnit = "into"
	BaseCoinUnit  = "uinto"
	IntoExponent  = 6

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = "into"
)

type GenesisParams struct {
	AirdropSupply            sdkmath.Int
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
	// ComputeParams   compute.Params
	IcaParams    icatypes.Params
	IntentParams intenttypes.Params
	//	ItemParams     itemtypes.Params
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
			genDoc.ConsensusParams = genesisParams.ConsensusParams

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
	slashingGenStateBz, err := cdc.MarshalJSON(slashingGenState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal slashing genesis state: %w", err)
	}
	appState[slashingtypes.ModuleName] = slashingGenStateBz

	// claim module genesis
	claimGenState := claimtypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
	claimGenState.Params = genesisParams.ClaimParams
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
	autoTxGenState := intenttypes.DefaultGenesis()
	autoTxGenState.Params = genesisParams.IntentParams
	autoTxGenStateBz, err := cdc.MarshalJSON(autoTxGenState)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal autoIbcTx genesis state: %w", err)
	}
	appState[intenttypes.ModuleName] = autoTxGenStateBz

	// compute module genesis
	// computeGenState := compute.DefaultGenesis()
	// computeGenState.Params = genesisParams.ComputeParams
	// computeGenStateBz, err := cdc.MarshalJSON(computeGenState)
	// if err != nil {
	// 	return nil, nil, fmt.Errorf("failed to marshal autoIbcTx genesis state: %w", err)
	// }
	// appState[compute.ModuleName] = computeGenStateBz

	// ica module genesis
	// icaGenState := icagenesistypes.DefaultGenesis()
	// icaGenState.HostGenesisState.Params = genesisParams.IcaParams
	// icaGenStateBz, err := cdc.MarshalJSON(icaGenState)
	// if err != nil {
	// 	return nil, nil, fmt.Errorf("failed to marshal ica genesis state: %w", err)
	// }
	// appState[icagenesistypes.ModuleName] = icaGenStateBz

	// item module genesis
	/*
		itemGenState := itemtypes.GetGenesisStateFromAppState(cdc, appState)
		itemGenState.Params = genesisParams.ItemParams
		itemGenStateBz, err := cdc.MarshalJSON(itemGenState)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal slashing genesis state: %w", err)
		}
		appState[itemtypes.ModuleName] = itemGenStateBz
	*/
	// return appState and genDoc
	return appState, nil
}

func MainnetGenesisParams() GenesisParams {
	genParams := GenesisParams{}

	genParams.AirdropSupply = sdk.NewInt(100_000_000_000_000)             // (100M INTO)
	genParams.GenesisTime = time.Date(2023, 02, 1, 17, 0, 0, 0, time.UTC) // 2023 2 Feb - 17:00 UTC
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
		Staking:                     sdk.MustNewDecFromStr("0.45"),
		CommunityPool:               sdk.MustNewDecFromStr("0.45"),
		TrustlessContractIncentives: sdk.MustNewDecFromStr("0.00"),
		RelayerIncentives:           sdk.MustNewDecFromStr("0.10"),
		//ItemIncentives:              sdk.MustNewDecFromStr("0.05"),
		ContributorRewards: sdk.MustNewDecFromStr("0.00"),
	}
	//genParams.AllocParams.WeightedContributorRewardsReceivers = []alloctypes.WeightedAddress{}

	// mint
	genParams.MintParams = minttypes.DefaultParams()
	genParams.MintParams.MintDenom = BaseCoinUnit
	genParams.MintParams.StartTime = genParams.GenesisTime.AddDate(0, 6, 0)
	genParams.MintParams.InitialAnnualProvisions = sdk.NewDec(150_000_000_000_000)

	genParams.MintParams.ReductionFactor = sdk.NewDec(3).QuoInt64(4)
	//31,536,000 seconds a year
	genParams.MintParams.BlocksPerYear = uint64(31540000 / 2) //assuming 2s average block times, param to be updated periodically
	// staking
	genParams.StakingParams = stakingtypes.DefaultParams()
	genParams.StakingParams.UnbondingTime = time.Hour * 24 * 21 //3 weeks
	genParams.StakingParams.MaxValidators = 50
	genParams.StakingParams.BondDenom = genParams.NativeCoinMetadatas[0].Base

	// genParams.StakingParams.MinCommissionRate = sdk.MustNewDecFromStr("0.05")
	// distr
	genParams.DistributionParams = distributiontypes.DefaultParams()

	genParams.DistributionParams.CommunityTax = sdk.MustNewDecFromStr("0.05")
	genParams.DistributionParams.WithdrawAddrEnabled = true
	// gov
	genParams.GovParams = govtypesv1.DefaultParams()
	maxDepositPeriod := time.Hour * 24 * 14 // 2 weeks
	genParams.GovParams.MaxDepositPeriod = &maxDepositPeriod
	genParams.GovParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(
		genParams.NativeCoinMetadatas[0].Base,
		sdk.NewInt(1_000_000_000),
	))
	genParams.GovParams.Quorum = "0.200000000000000000" // 20%
	votingPeriod := time.Hour * 24 * 3                  // 3 days
	genParams.GovParams.VotingPeriod = &votingPeriod
	// crisis
	genParams.CrisisConstantFee = sdk.NewCoin(
		genParams.NativeCoinMetadatas[0].Base,
		sdk.NewInt(100_000_000_000),
	)
	// slash
	genParams.SlashingParams = slashingtypes.DefaultParams()
	genParams.SlashingParams.SignedBlocksWindow = int64(25000)                       // ~41 hr at 6 second blocks
	genParams.SlashingParams.MinSignedPerWindow = sdk.MustNewDecFromStr("0.05")      // 5% minimum liveness
	genParams.SlashingParams.DowntimeJailDuration = time.Minute                      // 1 minute jail period
	genParams.SlashingParams.SlashFractionDoubleSign = sdk.MustNewDecFromStr("0.05") // 5% double sign slashing
	genParams.SlashingParams.SlashFractionDowntime = sdk.MustNewDecFromStr("0.0001") // 0.01% liveness slashing               // 0% liveness slashing
	//item
	/*
		genParams.ItemParams.MaxActivePeriod = time.Hour * 24 * 30
		genParams.ItemParams.MaxEstimatorCreatorRatio = 50
		genParams.ItemParams.MaxBuyerReward = 5000000000
		genParams.ItemParams.EstimationRatioForNewItem = 0
		genParams.ItemParams.CreateItemFee = 0
	*/

	//compute
	// genParams.ComputeParams = compute.DefaultParams()
	// genParams.ComputeParams.MaxContractDuration = time.Hour * 24 * 366 * 2
	// genParams.ComputeParams.MinContractDuration = time.Second * 30
	// genParams.ComputeParams.MinContractInterval = time.Second * 60
	// genParams.ComputeParams.AutoMsgFundsCommission = 2
	// genParams.ComputeParams.AutoMsgConstantFee = 5_000
	// genParams.ComputeParams.AutoMsgFlexFeeMul = 100
	// genParams.ComputeParams.RecurringAutoMsgConstantFee = 10_000
	// genParams.ComputeParams.MinContractDurationForIncentive = time.Hour * 24 * 4
	// genParams.ComputeParams.MinContractBalanceForIncentive = 50_000_000
	// genParams.ComputeParams.MaxContractIncentive = 500_000_000

	//AutoIBCTx
	genParams.IntentParams = intenttypes.DefaultParams()
	genParams.IntentParams.MaxActionDuration = time.Hour * 24 * 366 * 2
	genParams.IntentParams.MinActionDuration = time.Second * 60
	genParams.IntentParams.MinActionInterval = time.Second * 60
	genParams.IntentParams.ActionFundsCommission = 2
	genParams.IntentParams.ActionConstantFee = 7_000
	genParams.IntentParams.ActionFlexFeeMul = 3
	genParams.IntentParams.RecurringActionConstantFee = 5_000
	genParams.IntentParams.RelayerRewards = []int64{10_000, 15_000, 18_000, 22_000}

	//claim
	genParams.ClaimParams = claimtypes.DefaultGenesis().Params
	genParams.ClaimParams.AirdropStartTime = genParams.GenesisTime.Add(time.Hour * 24 * 365) // 1 year (will be changed through gov)
	genParams.ClaimParams.DurationUntilDecay = time.Hour * 24 * 60                           // 60 days = ~2 months
	genParams.ClaimParams.DurationOfDecay = time.Hour * 24 * 120                             // 120 days = ~4 months
	genParams.ClaimParams.ClaimDenom = genParams.NativeCoinMetadatas[0].Base
	genParams.ClaimParams.DurationVestingPeriods = []time.Duration{time.Hour * 24 * 7, time.Hour * 24 * 7 * 6, time.Hour * 24 * 3, time.Hour * 24 * 7 * 4}

	//consensus
	genParams.ConsensusParams = tmtypes.DefaultConsensusParams()
	genParams.ConsensusParams.Block.MaxBytes = 5 * 1024 * 1024
	genParams.ConsensusParams.Block.MaxGas = 6_000_000
	genParams.ConsensusParams.Evidence.MaxAgeDuration = genParams.StakingParams.UnbondingTime
	genParams.ConsensusParams.Evidence.MaxAgeNumBlocks = int64(genParams.StakingParams.UnbondingTime.Seconds()) / 3
	genParams.ConsensusParams.Version.App = 1
	genParams.DistributedAccounts = []banktypes.Balance{}

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
		sdk.NewInt(1_000_000), // 1 INTO
	))
	genParams.GovParams.Quorum = "0.100000000000000000" // 10%
	votingPeriod := time.Minute
	genParams.GovParams.VotingPeriod = &votingPeriod

	//action
	genParams.IntentParams.MinActionDuration = time.Second * 40
	genParams.IntentParams.MinActionInterval = time.Second * 40

	//claim
	genParams.ClaimParams.AirdropStartTime = genParams.GenesisTime
	genParams.ClaimParams.DurationUntilDecay = time.Hour * 24 * 5 // 5 days
	genParams.ClaimParams.DurationOfDecay = time.Hour * 24 * 5    // 5 days
	genParams.ClaimParams.DurationVestingPeriods = []time.Duration{time.Minute, time.Minute * 2, time.Minute * 5, time.Minute}

	//31,536,000 seconds a year and estimated 4s block times
	genParams.MintParams.BlocksPerYear = uint64(31540000 / 4)
	//compute
	// genParams.ComputeParams.MaxContractDuration = time.Hour * 24 * 60
	// genParams.ComputeParams.MinContractDuration = time.Second * 10
	// genParams.ComputeParams.MinContractInterval = time.Second * 20
	// genParams.ComputeParams.MinContractDurationForIncentive = time.Second
	// genParams.ComputeParams.MinContractBalanceForIncentive = 50000
	// genParams.ComputeParams.MaxContractIncentive = 500000

	//item
	/*
		genParams.ItemParams.MaxActivePeriod = time.Hour * 24 * 5 // 5 days
		genParams.ItemParams.MaxEstimatorCreatorRatio = 100
		genParams.ItemParams.MaxBuyerReward = 500000000000
	*/
	return genParams
}
