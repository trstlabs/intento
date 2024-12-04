package apptesting

import (
	"encoding/json"
	"os"

	"cosmossdk.io/log"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/testutil/sims"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"
	"github.com/trstlabs/intento/app"
)

type TestChain struct {
	*ibctesting.TestChain
}

func SetupTestingApp() (ibctesting.TestingApp, map[string]json.RawMessage) {
	dir, _ := os.MkdirTemp("", "ibctest")
	appOptions := sims.NewAppOptionsWithFlagHome(dir)
	IntoApp := app.NewIntoApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, appOptions, []wasmkeeper.Option{})
	return IntoApp, app.NewDefaultGenesisState(IntoApp.AppCodec())
}

// GetIntoApp returns the current chain's app as an IntoApp
func (chain *TestChain) GetIntoApp() *app.IntoApp {
	v, _ := chain.App.(*app.IntoApp)

	return v
}
