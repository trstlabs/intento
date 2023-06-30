package keeper_test

import (
	"encoding/json"

	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	app "github.com/trstlabs/trst/app"
	"github.com/trstlabs/trst/x/mint/types"
	// "github.com/trstlabs/trst/x/compute"
)

// returns context and an app with updated mint keeper
func createTestApp(isCheckTx bool) (*app.TrstApp, sdk.Context) {
	app := setup(isCheckTx)

	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.MintKeeper.SetParams(ctx, types.DefaultParams())
	app.MintKeeper.SetMinter(ctx, types.DefaultInitialMinter())

	return app, ctx
}

func setup(isCheckTx bool) *app.TrstApp {
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

func genApp(withGenesis bool, invCheckPeriod uint) (*app.TrstApp, app.GenesisState) {
	db := dbm.NewMemDB()
	encCdc := app.MakeEncodingConfig()
	TrstApp := app.NewTrstApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		app.EmptyAppOptions{},
	)

	if withGenesis {
		return TrstApp, app.NewDefaultGenesisState(encCdc.Codec)
	}

	return TrstApp, app.GenesisState{}
}
