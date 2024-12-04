package keeper_test

import (
	"encoding/json"

	"cosmossdk.io/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	abci "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	app "github.com/trstlabs/intento/app"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
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
	IntoApp := app.NewIntoApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		app.EmptyAppOptions{},
		[]wasmkeeper.Option{},
	)

	if withGenesis {
		return IntoApp, app.NewDefaultGenesisState(IntoApp.AppCodec())
	}

	return IntoApp, app.GenesisState{}
}
