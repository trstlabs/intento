package app

import (
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgrades "github.com/trstlabs/intento/app/upgrades"
	mainnetupgradesv103 "github.com/trstlabs/intento/app/upgrades/mainnet/v1.0.3"
	mainnetupgradesv105 "github.com/trstlabs/intento/app/upgrades/mainnet/v1.0.5"
	mainnetupgradesv110 "github.com/trstlabs/intento/app/upgrades/mainnet/v1.10.0"
)

var Upgrades = []upgrades.Upgrade{
	// mainnet upgrades
	mainnetupgradesv103.Upgrade,
	mainnetupgradesv105.Upgrade,
	mainnetupgradesv110.Upgrade,
}

func (app IntoApp) RegisterUpgradeHandlers(configurator module.Configurator) {

	for _, u := range Upgrades {
		app.UpgradeKeeper.SetUpgradeHandler(
			u.UpgradeName,
			u.CreateUpgradeHandler(app.ModuleManager, configurator, upgrades.IntentoKeepers{}),
		)
	}

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	for _, u := range Upgrades {
		u := u
		if upgradeInfo.Name == u.UpgradeName {
			app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &u.StoreUpgrades))
		}
	}
}
