package keeper_test

import (
	"encoding/json"

	"github.com/tendermint/spm/cosmoscmd"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
	app "github.com/trstlabs/trst/app"
	"github.com/trstlabs/trst/x/mint/types"
	// "github.com/trstlabs/trst/x/compute"
)

// returns context and an app with updated mint keeper
func createTestApp(isCheckTx bool) (*app.TrstApp, sdk.Context) {
	app := setup(isCheckTx)

	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.AppKeepers.MintKeeper.SetParams(ctx, types.DefaultParams())
	app.AppKeepers.MintKeeper.SetMinter(ctx, types.DefaultInitialMinter())

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
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

func genApp(withGenesis bool, invCheckPeriod uint) (*app.TrstApp, app.GenesisState) {
	db := dbm.NewMemDB()
	encCdc := cosmoscmd.MakeEncodingConfig(app.ModuleBasics())
	TrstApp := app.NewTrstApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		simapp.DefaultNodeHome,
		invCheckPeriod,
		true,
		simapp.EmptyAppOptions{},
		// compute.GetConfig(simapp.EmptyAppOptions{}),
		// app.GetEnabledProposals(),
	)

	if withGenesis {
		return TrstApp, app.NewDefaultGenesisState(encCdc.Marshaler)
	}

	return TrstApp, app.GenesisState{}
}
