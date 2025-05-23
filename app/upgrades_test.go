package app

// import (

// 	// "github.com/cosmos/cosmos-sdk/testutil/testdata_pulsar"

// 	"testing"

// 	upgradetypes "cosmossdk.io/x/upgrade/types"

// 	"github.com/cosmos/cosmos-sdk/types/module"
// 	"github.com/stretchr/testify/require"

// 	upgrades "github.com/trstlabs/intento/app/upgrades"
// 	testnetupgradesv094r1 "github.com/trstlabs/intento/app/upgrades/testnet/v0.9.4-r1"
// )

// // func TestUpgradeHandler(t *testing.T) {
// // 	// Setup a mock context and keepers (can be in-memory stores)

// // 	app := InitIntentoTestApp(true)
// // 	ctx := app.BaseApp.NewContext(true)
// // 	// Setup module manager and configurator mocks or real instances depending on your app structure
// // 	mm := app.ModuleManager
// // 	configurator := app.Configurator // or mock if possible

// // 	// Use the actual upgrade handler function you defined
// // 	handler := testnetupgradesv094r1.Upgrade.CreateUpgradeHandler(mm, configurator, upgrades.IntentoKeepers{})

// // 	// Prepare a mock VersionMap (simulate the from-version of modules)
// // 	fromVM := module.VersionMap{}

// // 	// Call the handler
// // 	newVM, err := handler(ctx, upgradetypes.Plan{Name: testnetupgradesv094r1.Upgrade.UpgradeName}, fromVM)
// // 	require.NoError(t, err)

// // 	// Assert the new version map has been updated (modules migrated)
// // 	for modName, version := range newVM {
// // 		t.Logf("Module %s migrated to version %d", modName, version)
// // 	}

// // 	// Optionally: check that specific migrations ran, or store upgrades happened
// // 	// by querying your keeper or store state here
// // }
