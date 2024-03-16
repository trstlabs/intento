package intentoibctesting

import (
	"encoding/json"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/trstlabs/intento/app"
	autoibctxtypes "github.com/trstlabs/intento/x/auto-ibc-tx/types"
	// "github.com/trstlabs/intento/x/compute"
)

type TestChain struct {
	*ibctesting.TestChain
}

func SetupTestingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	encCdc := app.MakeEncodingConfig()
	TrstApp := *app.NewTrstApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, encCdc, app.EmptyAppOptions{} /* , compute.DefaultWasmConfig(), app.GetEnabledProposals() */)
	TrstApp.AutoIbcTxKeeper.SetParams(TrstApp.GetBaseApp().NewContext(true, tmproto.Header{Height: TrstApp.LastBlockHeight()}), autoibctxtypes.DefaultParams())
	//encCdc := app.MakeEncodingConfig()
	return &TrstApp, app.NewDefaultGenesisState(encCdc.Codec)
}

// GetTrstApp returns the current chain's app as an TrstApp
func (chain *TestChain) GetTrstApp() *app.TrstApp {
	v, _ := chain.App.(*app.TrstApp)

	return v
}
