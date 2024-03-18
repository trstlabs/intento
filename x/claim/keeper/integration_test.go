package keeper_test

import (
	"encoding/json"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/app"
	"github.com/trstlabs/intento/x/mint/types"
)

// returns context and an app with updated mint keeper
func createTestApp(isCheckTx bool) (*app.IntoApp, sdk.Context) {
	app := setup(isCheckTx)

	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.MintKeeper.SetParams(ctx, types.DefaultParams())
	app.MintKeeper.SetMinter(ctx, types.DefaultInitialMinter())

	return app, ctx
}

func setup(isCheckTx bool) *app.IntoApp {
	app, genesisState := genApp(!isCheckTx, 5)

	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: &tmproto.ConsensusParams{},
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

func genApp(withGenesis bool, invCheckPeriod uint) (*app.IntoApp, app.GenesisState) {
	db := dbm.NewMemDB()
	encCdc := app.MakeEncodingConfig()
	IntoApp := app.NewIntoApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		encCdc,
		app.EmptyAppOptions{},
		// compute.GetConfig(simapp.EmptyAppOptions{}),
		// app.GetEnabledProposals(),
	)

	// if withGenesis {
	// 	encCdc := app.MakeEncodingConfig()
	// 	return IntoApp, app.NewDefaultGenesisState(encCdc.Codec)
	// }

	return IntoApp, app.GenesisState{}
}
