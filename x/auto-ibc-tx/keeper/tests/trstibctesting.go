package trstibctesting

import (
	"encoding/json"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/trstlabs/trst/app"
	autoibctxtypes "github.com/trstlabs/trst/x/auto-ibc-tx/types"
	// "github.com/trstlabs/trst/x/compute"
)

type TestChain struct {
	*ibctesting.TestChain
}

func SetupTestingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	TrstApp := *app.NewTrstApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, app.EmptyAppOptions{} /* , compute.DefaultWasmConfig(), app.GetEnabledProposals() */)
	TrstApp.AutoIbcTxKeeper.SetParams(TrstApp.GetBaseApp().NewContext(true, tmproto.Header{Height: TrstApp.LastBlockHeight()}), autoibctxtypes.DefaultParams())
	//encCdc := app.MakeEncodingConfig()
	return &TrstApp, app.NewDefaultGenesisState(TrstApp.AppCodec())
}

// GetTrstApp returns the current chain's app as an TrstApp
func (chain *TestChain) GetTrstApp() *app.TrstApp {
	v, _ := chain.App.(*app.TrstApp)

	return v
}
