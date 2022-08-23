package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	appParams "github.com/trstlabs/trst/app/params"
	alloctypes "github.com/trstlabs/trst/x/alloc/types"
	claimtypes "github.com/trstlabs/trst/x/claim/types"
	compute "github.com/trstlabs/trst/x/compute"
	minttypes "github.com/trstlabs/trst/x/mint/types"
	//	itemtypes "github.com/trstlabs/trst/x/item/types"
)

type GenesisParams struct {
	AirdropSupply            sdk.Int
	StrategicReserveAccounts []banktypes.Balance
	DistributedAccounts      []banktypes.Balance
	ConsensusParams          *tmproto.ConsensusParams

	GenesisTime         time.Time
	NativeCoinMetadatas []banktypes.Metadata

	StakingParams      stakingtypes.Params
	DistributionParams distributiontypes.Params
	GovParams          govtypes.Params

	CrisisConstantFee sdk.Coin

	SlashingParams slashingtypes.Params
	AllocParams    alloctypes.Params
	ClaimParams    claimtypes.Params
	MintParams     minttypes.Params
	ComputeParams  compute.Params
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
	starsd prepare-genesis mainnet stargaze-1
	- Check input genesis:
		file is at ~/.starsd/config/genesis.json
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
			appState, genDoc, err = PrepareGenesis(clientCtx, appState, genDoc, genesisParams, chainID)
			if err != nil {
				return fmt.Errorf("failed to prepare genesis: %w", err)
			}

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
	genDoc *tmtypes.GenesisDoc,
	genesisParams GenesisParams,
	chainID string,
) (map[string]json.RawMessage, *tmtypes.GenesisDoc, error) {
	cdc := clientCtx.Codec
	// chain params genesis
	genDoc.GenesisTime = genesisParams.GenesisTime
	genDoc.ChainID = chainID
	genDoc.ConsensusParams = genesisParams.ConsensusParams

	// ---
	// bank module genesis
	bankGenState := banktypes.DefaultGenesisState()
	bankGenState.Params.DefaultSendEnabled = true
	bankGenStateBz, err := cdc.MarshalJSON(bankGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal bank genesis state: %w", err)
	}
	appState[banktypes.ModuleName] = bankGenStateBz

	// IBC transfer module genesis
	ibcGenState := ibctransfertypes.DefaultGenesisState()
	ibcGenState.Params.SendEnabled = true
	ibcGenState.Params.ReceiveEnabled = true
	ibcGenStateBz, err := cdc.MarshalJSON(ibcGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal IBC transfer genesis state: %w", err)
	}
	appState[ibctransfertypes.ModuleName] = ibcGenStateBz

	// mint module genesis

	mintGenState := minttypes.DefaultGenesisState()
	mintGenState.Params = genesisParams.MintParams

	mintGenStateBz, err := cdc.MarshalJSON(mintGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal mint genesis state: %w", err)
	}
	appState[minttypes.ModuleName] = mintGenStateBz

	// staking module genesis
	stakingGenState := stakingtypes.GetGenesisStateFromAppState(cdc, appState)
	stakingGenState.Params = genesisParams.StakingParams
	stakingGenStateBz, err := cdc.MarshalJSON(stakingGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal staking genesis state: %w", err)
	}
	appState[stakingtypes.ModuleName] = stakingGenStateBz

	// distribution module genesis
	distributionGenState := distributiontypes.DefaultGenesisState()
	distributionGenState.Params = genesisParams.DistributionParams
	distributionGenStateBz, err := cdc.MarshalJSON(distributionGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal distribution genesis state: %w", err)
	}
	appState[distributiontypes.ModuleName] = distributionGenStateBz

	// gov module genesis
	govGenState := govtypes.DefaultGenesisState()
	govGenState.DepositParams = genesisParams.GovParams.DepositParams
	govGenState.TallyParams = genesisParams.GovParams.TallyParams
	govGenState.VotingParams = genesisParams.GovParams.VotingParams
	govGenStateBz, err := cdc.MarshalJSON(govGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal gov genesis state: %w", err)
	}
	appState[govtypes.ModuleName] = govGenStateBz

	// crisis module genesis
	crisisGenState := crisistypes.DefaultGenesisState()
	crisisGenState.ConstantFee = genesisParams.CrisisConstantFee
	crisisGenStateBz, err := cdc.MarshalJSON(crisisGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal crisis genesis state: %w", err)
	}
	appState[crisistypes.ModuleName] = crisisGenStateBz

	// slashing module genesis
	slashingGenState := slashingtypes.DefaultGenesisState()
	slashingGenState.Params = genesisParams.SlashingParams
	slashingGenStateBz, err := cdc.MarshalJSON(slashingGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal slashing genesis state: %w", err)
	}
	appState[slashingtypes.ModuleName] = slashingGenStateBz

	// claim module genesis
	claimGenState := claimtypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
	claimGenState.Params = genesisParams.ClaimParams
	claimGenStateBz, err := cdc.MarshalJSON(claimGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal claim genesis state: %w", err)
	}
	appState[claimtypes.ModuleName] = claimGenStateBz

	// alloc module genesis
	allocGenState := alloctypes.GetGenesisStateFromAppState(cdc, appState)
	allocGenState.Params = genesisParams.AllocParams
	allocGenStateBz, err := cdc.MarshalJSON(allocGenState)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal alloc genesis state: %w", err)
	}
	appState[alloctypes.ModuleName] = allocGenStateBz

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
	return appState, genDoc, nil
}

func MainnetGenesisParams() GenesisParams {
	genParams := GenesisParams{}

	genParams.AirdropSupply = sdk.NewInt(100_000_000_000_000)             // (100M TRST)
	genParams.GenesisTime = time.Date(2023, 02, 1, 17, 0, 0, 0, time.UTC) // 2023 2 Feb - 17:00 UTC
	genParams.NativeCoinMetadatas = []banktypes.Metadata{
		{
			Description: "The native token of TRST",
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    appParams.BaseCoinUnit,
					Exponent: 0,
					Aliases:  nil,
				},
				{
					Denom:    appParams.HumanCoinUnit,
					Exponent: appParams.TrstExponent,
					Aliases:  nil,
				},
			},
			Base:    appParams.BaseCoinUnit,
			Display: appParams.HumanCoinUnit,
		},
	}
	// alloc
	genParams.AllocParams = alloctypes.DefaultParams()
	genParams.AllocParams.DistributionProportions = alloctypes.DistributionProportions{
		Staking:                     sdk.MustNewDecFromStr("0.60"), // 25%
		CommunityPool:               sdk.MustNewDecFromStr("0.25"), // 5%
		TrustlessContractIncentives: sdk.MustNewDecFromStr("0.10"), // 45%
		//ItemIncentives:              sdk.MustNewDecFromStr("0.05"), // 45%
		ContributorRewards: sdk.MustNewDecFromStr("0.05"), // 25%

	}
	genParams.AllocParams.WeightedContributorRewardsReceivers = []alloctypes.WeightedAddress{
		{
			Address: "trust1sns5l9cvkgf4fy770nmg98e7uzet5xhhmv8njv",
			Weight:  sdk.NewDecWithPrec(100, 2),
		},
	}

	// mint
	genParams.MintParams = minttypes.DefaultParams()
	genParams.MintParams.MintDenom = appParams.BaseCoinUnit
	genParams.MintParams.StartTime = genParams.GenesisTime.AddDate(0, 6, 0)
	genParams.MintParams.InitialAnnualProvisions = sdk.NewDec(250_000_000_000_000)
	genParams.MintParams.ReductionFactor = sdk.NewDec(2).QuoInt64(3)
	genParams.MintParams.BlocksPerYear = uint64(5737588)
	// staking
	genParams.StakingParams = stakingtypes.DefaultParams()
	genParams.StakingParams.UnbondingTime = time.Hour * 24 * 21 //3 weeks
	genParams.StakingParams.MaxValidators = 50
	genParams.StakingParams.BondDenom = genParams.NativeCoinMetadatas[0].Base

	// genParams.StakingParams.MinCommissionRate = sdk.MustNewDecFromStr("0.05")
	// distr
	genParams.DistributionParams = distributiontypes.DefaultParams()
	genParams.DistributionParams.BaseProposerReward = sdk.MustNewDecFromStr("0.01")
	genParams.DistributionParams.BonusProposerReward = sdk.MustNewDecFromStr("0.04")
	genParams.DistributionParams.CommunityTax = sdk.MustNewDecFromStr("0.05")
	genParams.DistributionParams.WithdrawAddrEnabled = true
	// gov
	genParams.GovParams = govtypes.DefaultParams()
	genParams.GovParams.DepositParams.MaxDepositPeriod = time.Hour * 24 * 14 // 2 weeks
	genParams.GovParams.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(
		genParams.NativeCoinMetadatas[0].Base,
		sdk.NewInt(1_000_000_000),
	))
	genParams.GovParams.TallyParams.Quorum = sdk.MustNewDecFromStr("0.2") // 20%
	genParams.GovParams.VotingParams.VotingPeriod = time.Hour * 24 * 3    // 3 days
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
	genParams.ComputeParams.MaxContractDuration = time.Hour * 24 * 366
	genParams.ComputeParams.MinContractDuration = time.Second * 30
	genParams.ComputeParams.MinContractInterval = time.Second * 60
	genParams.ComputeParams.AutoMsgFundsCommission = 2
	genParams.ComputeParams.AutoMsgConstantFee = 1000000
	genParams.ComputeParams.RecurringAutoMsgConstantFee = 1000000
	genParams.ComputeParams.MinContractDurationForIncentive = time.Hour * 24 * 4
	genParams.ComputeParams.MinContractBalanceForIncentive = 50000000
	genParams.ComputeParams.MaxContractIncentive = 500000000

	//claim
	genParams.ClaimParams.AirdropStartTime = genParams.GenesisTime.Add(time.Hour * 24 * 365) // 1 year (will be changed through gov)
	genParams.ClaimParams.DurationUntilDecay = time.Hour * 24 * 60                           // 60 days = ~2 months
	genParams.ClaimParams.DurationOfDecay = time.Hour * 24 * 120                             // 120 days = ~4 months
	genParams.ClaimParams.ClaimDenom = genParams.NativeCoinMetadatas[0].Base
	genParams.ClaimParams.DurationVestingPeriods = []time.Duration{time.Hour * 24 * 7, time.Hour * 24 * 7 * 6, time.Hour * 24 * 3, time.Hour * 24 * 7 * 4}
	genParams.ConsensusParams = tmtypes.DefaultConsensusParams()
	genParams.ConsensusParams.Block.MaxBytes = 5 * 1024 * 1024
	genParams.ConsensusParams.Block.MaxGas = 6_000_000
	genParams.ConsensusParams.Evidence.MaxAgeDuration = genParams.StakingParams.UnbondingTime
	genParams.ConsensusParams.Evidence.MaxAgeNumBlocks = int64(genParams.StakingParams.UnbondingTime.Seconds()) / 3
	genParams.ConsensusParams.Version.AppVersion = 1
	genParams.DistributedAccounts = []banktypes.Balance{}

	return genParams
}

func TestnetGenesisParams() GenesisParams {

	genParams := MainnetGenesisParams()

	// genParams.GenesisTime = time.Now()
	genParams.GenesisTime = time.Now()
	genParams.MintParams.StartTime = time.Now()
	genParams.StakingParams.UnbondingTime = time.Hour * 24 * 3 // 3 days

	//gov
	genParams.GovParams.DepositParams.MinDeposit = sdk.NewCoins(sdk.NewCoin(
		genParams.NativeCoinMetadatas[0].Base,
		sdk.NewInt(1_000_000), // 1 TRST
	))
	genParams.GovParams.TallyParams.Quorum = sdk.MustNewDecFromStr("0.1") // 10%
	genParams.GovParams.VotingParams.VotingPeriod = time.Minute           //time.Hour * 24 * 1    // 1 day

	//claim
	genParams.ClaimParams.AirdropStartTime = genParams.GenesisTime
	genParams.ClaimParams.DurationUntilDecay = time.Hour * 24 * 5 // 5 days
	genParams.ClaimParams.DurationOfDecay = time.Hour * 24 * 5    // 5 days
	genParams.ClaimParams.DurationVestingPeriods = []time.Duration{time.Minute, time.Minute * 2, time.Minute * 5, time.Minute}

	//compute
	genParams.ComputeParams.MaxContractDuration = time.Hour * 24 * 60
	genParams.ComputeParams.MinContractDuration = time.Second * 10
	genParams.ComputeParams.MinContractInterval = time.Second * 20
	genParams.ComputeParams.MinContractDurationForIncentive = time.Second
	genParams.ComputeParams.MinContractBalanceForIncentive = 50000
	genParams.ComputeParams.MaxContractIncentive = 500000

	//item
	/*
		genParams.ItemParams.MaxActivePeriod = time.Hour * 24 * 5 // 5 days
		genParams.ItemParams.MaxEstimatorCreatorRatio = 100
		genParams.ItemParams.MaxBuyerReward = 500000000000
	*/
	return genParams
}
