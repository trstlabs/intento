package keeper_test

import (
	"encoding/json"
	"os"

	"cosmossdk.io/log"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/app"
	"github.com/trstlabs/intento/x/mint/types"
)

// returns context and an app with updated mint keeper
func createTestApp(isCheckTx bool) (*app.IntoApp, sdk.Context) {
	app := setup(isCheckTx)

	ctx := app.BaseApp.NewContext(isCheckTx)
	app.MintKeeper.SetParams(ctx, types.DefaultParams())
	app.MintKeeper.SetMinter(ctx, types.DefaultInitialMinter())

	return app, ctx
}

func setup(isCheckTx bool) *app.IntoApp {
	app, genesisState := genApp(!isCheckTx)

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

func genApp(withGenesis bool) (*app.IntoApp, app.GenesisState) {
	db := dbm.NewMemDB()

	dir, _ := os.MkdirTemp("", "ibctest")
	appOptions := sims.NewAppOptionsWithFlagHome(dir)
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
