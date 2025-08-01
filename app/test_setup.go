package app

import (
	"encoding/json"
	"os"
	"time"

	"cosmossdk.io/log"
	math "cosmossdk.io/math"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	tmtypes "github.com/cometbft/cometbft/types"
	cometbftdb "github.com/cosmos/cosmos-db"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	ibccommitmenttypes "github.com/cosmos/ibc-go/v8/modules/core/23-commitment/types"
	ibctmtypes "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"
	ccvconsumertypes "github.com/cosmos/interchain-security/v6/x/ccv/consumer/types"
	ccvprovidertypes "github.com/cosmos/interchain-security/v6/x/ccv/provider/types"

	ccvtypes "github.com/cosmos/interchain-security/v6/x/ccv/types"
)

func init() {
	SetupConfig()
}

func SetupConfig() {
	const Bech32Prefix = "cosmos"

	config := sdk.GetConfig()
	valoper := sdk.PrefixValidator + sdk.PrefixOperator
	valoperpub := sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	config.SetBech32PrefixForAccount(Bech32Prefix, Bech32Prefix+sdk.PrefixPublic)
	config.SetBech32PrefixForValidator(Bech32Prefix+valoper, Bech32Prefix+valoperpub)
	//BaseCoinUnit := sdk.DefaultBondDenom
	// if err := sdk.RegisterDenom(BaseCoinUnit, math.LegacyNewDecWithPrec(1, 6)); err != nil {
	// 	panic(err)
	// }
	// deno, _ := sdk.GetDenomUnit(BaseCoinUnit)
	// fmt.Printf("DEM %s", deno)
	// den, err := sdk.GetBaseDenom()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("DEM %s", den)
	// if err := sdk.SetBaseDenom(BaseCoinUnit); err != nil {
	// 	panic(err)
	// }

}

// Initializes a new IntoApp without IBC functionality
func InitIntentoTestApp(initChain bool) *IntoApp {
	db := cometbftdb.NewMemDB()
	dir, _ := os.MkdirTemp("", "ibctest")
	appOptions := sims.NewAppOptionsWithFlagHome(dir)

	app := NewIntoApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		appOptions,
		[]wasmkeeper.Option{},
	)
	if initChain {

		genesisState := GenesisStateWithValSet(app)
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		app.InitChain(
			&abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: sims.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)

	}

	return app
}

func GenesisStateWithValSet(app *IntoApp) GenesisState {
	privVal := mock.NewPV()
	pubKey, _ := privVal.GetPubKey()
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	senderPrivKey.PubKey().Address()
	acc := authtypes.NewBaseAccountWithAddress(senderPrivKey.PubKey().Address().Bytes())
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100000000000000))),
	}

	//////////////////////
	balances := []banktypes.Balance{balance}
	genesisState := NewDefaultGenesisState(app.appCodec)
	genAccs := []authtypes.GenesisAccount{acc}
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.DefaultPowerReduction
	initValPowers := []abci.ValidatorUpdate{}

	for _, val := range valSet.Validators {
		pk, _ := cryptocodec.FromCmtPubKeyInterface(val.PubKey)
		pkAny, _ := codectypes.NewAnyWithValue(pk)
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   math.LegacyOneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
			MinSelfDelegation: math.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress().String(), sdk.ValAddress(val.Address).String(), math.LegacyOneDec()))

		// add initial validator powers so consumer InitGenesis runs correctly
		pub, _ := val.ToProto()
		initValPowers = append(initValPowers, abci.ValidatorUpdate{
			Power:  val.VotingPower,
			PubKey: pub.PubKey,
		})
	}
	// set validators and delegations
	stakingGenesis := stakingtypes.NewGenesisState(stakingtypes.DefaultParams(), validators, delegations)
	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	for range delegations {
		// add delegated tokens to total supply
		totalSupply = totalSupply.Add(sdk.NewCoin(sdk.DefaultBondDenom, bondAmt))
	}

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, bondAmt)},
	})

	// update total supply
	bankGenesis := banktypes.NewGenesisState(
		banktypes.DefaultGenesisState().Params,
		balances,
		totalSupply,
		[]banktypes.Metadata{},
		[]banktypes.SendEnabled{},
	)
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	vals, err := tmtypes.PB2TM.ValidatorUpdates(initValPowers)
	if err != nil {
		panic("failed to get vals")
	}

	consumerGenesisState := CreateMinimalConsumerTestGenesis()
	consumerGenesisState.Provider.InitialValSet = initValPowers
	consumerGenesisState.Provider.ConsensusState.NextValidatorsHash = tmtypes.NewValidatorSet(vals).Hash()
	consumerGenesisState.Params.Enabled = true
	if err := consumerGenesisState.Validate(); err != nil {
		panic(err)
	}
	genesisState[ccvconsumertypes.ModuleName] = app.AppCodec().MustMarshalJSON(consumerGenesisState)

	return genesisState
}

// Initializes a new Intento App casted as a TestingApp for IBC support
func InitIntentoIBCTestingApp(initValPowers []abci.ValidatorUpdate) func() (ibctesting.TestingApp, map[string]json.RawMessage) {
	return func() (ibctesting.TestingApp, map[string]json.RawMessage) {

		app := InitIntentoTestApp(false)
		genesisState := NewDefaultGenesisState(app.appCodec)
		// Feed consumer genesis with provider validators
		var consumerGenesis ccvconsumertypes.GenesisState
		app.appCodec.MustUnmarshalJSON(genesisState[ccvconsumertypes.ModuleName], &consumerGenesis)
		consumerGenesis.Provider.InitialValSet = initValPowers
		consumerGenesis.Params.Enabled = true
		genesisState[ccvconsumertypes.ModuleName] = app.appCodec.MustMarshalJSON(&consumerGenesis)

		return app, genesisState
	}
}

// This function creates consumer module genesis state that is used as starting point for modifications
// that allow chain to be started locally without having to start the provider chain and the relayer.
// It is also used in tests that are starting the chain node.
func CreateMinimalConsumerTestGenesis() *ccvconsumertypes.GenesisState {
	genesisState := ccvconsumertypes.DefaultGenesisState()
	genesisState.Params.Enabled = true
	genesisState.NewChain = true
	genesisState.Provider.ClientState = ccvprovidertypes.DefaultParams().TemplateClient
	genesisState.Provider.ClientState.ChainId = "intento"
	genesisState.Provider.ClientState.LatestHeight = ibcclienttypes.Height{RevisionNumber: 0, RevisionHeight: 1}
	genesisState.Provider.ClientState.TrustingPeriod, _ = ccvtypes.CalculateTrustPeriod(genesisState.Params.UnbondingPeriod, ccvprovidertypes.DefaultTrustingPeriodFraction)

	genesisState.Provider.ClientState.UnbondingPeriod = genesisState.Params.UnbondingPeriod
	genesisState.Provider.ClientState.MaxClockDrift = ccvprovidertypes.DefaultMaxClockDrift
	genesisState.Provider.ConsensusState = &ibctmtypes.ConsensusState{
		Timestamp: time.Now().UTC(),
		Root:      ibccommitmenttypes.MerkleRoot{Hash: []byte("dummy")},
	}

	return genesisState
}
