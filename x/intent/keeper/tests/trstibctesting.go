package intentoibctesting

import (
	"encoding/json"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/trstlabs/intento/app"
	intenttypes "github.com/trstlabs/intento/x/intent/types"
	// "github.com/trstlabs/intento/x/compute"
)

type TestChain struct {
	*ibctesting.TestChain
}

func SetupTestingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	encCdc := app.MakeEncodingConfig()
	IntoApp := *app.NewIntoApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, encCdc, app.EmptyAppOptions{} /* , compute.DefaultWasmConfig(), app.GetEnabledProposals() */)
	IntoApp.IntentKeeper.SetParams(IntoApp.GetBaseApp().NewContext(true, tmproto.Header{Height: IntoApp.LastBlockHeight()}), intenttypes.DefaultParams())
	//encCdc := app.MakeEncodingConfig()
	return &IntoApp, app.NewDefaultGenesisState(encCdc.Codec)
}

// GetIntoApp returns the current chain's app as an IntoApp
func (chain *TestChain) GetIntoApp() *app.IntoApp {
	v, _ := chain.App.(*app.IntoApp)

	return v
}
