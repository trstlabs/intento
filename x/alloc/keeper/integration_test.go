package keeper_test

import (
	"encoding/json"

	"cosmossdk.io/log"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	app "github.com/trstlabs/intento/app"
	"github.com/trstlabs/intento/x/mint/types"
)

// returns context and an app with updated mint keeper
func createTestApp(isCheckTx bool, dir string) (*app.IntoApp, sdk.Context) {
	app := setup(isCheckTx, dir)

	ctx := app.BaseApp.NewContext(isCheckTx)
	app.MintKeeper.SetParams(ctx, types.DefaultParams())
	app.MintKeeper.SetMinter(ctx, types.DefaultInitialMinter())

	return app, ctx
}

func setup(isCheckTx bool, dir string) *app.IntoApp {
	app, genesisState := genApp(!isCheckTx, dir)

	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			&abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: &tmproto.ConsensusParams{},
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

func genApp(withGenesis bool, dir string) (*app.IntoApp, app.GenesisState) {
	db := dbm.NewMemDB()

	app.SetupConfig()

	appOptions := make(sims.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = dir
	IntoApp := app.NewIntoApp(
		log.NewNopLogger(),
		db,
		nil,
		true,

		appOptions,
		[]wasmkeeper.Option{},
	)

	if withGenesis {
		return IntoApp, app.NewDefaultGenesisState(IntoApp.AppCodec())
	}

	return IntoApp, app.GenesisState{}
}
