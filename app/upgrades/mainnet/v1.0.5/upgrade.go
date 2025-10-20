package v105

import (
	"context"
	"time"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	upgrades "github.com/trstlabs/intento/app/upgrades"
)

// UpgradeName defines the on-chain upgrade name for the Intento v1.0.5 upgrade
const UpgradeName = "v1.0.5"

// Upgrade defines the v1.0.5 upgrade
var Upgrade = upgrades.Upgrade{
	UpgradeName: UpgradeName,
	CreateUpgradeHandler: func(mm *module.Manager, cfg module.Configurator, _ upgrades.IntentoKeepers) upgradetypes.UpgradeHandler {
		return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			startTime := time.Now()
			wctx := sdk.UnwrapSDKContext(ctx)
			wctx.Logger().Info("starting upgrade", "upgrade_name", UpgradeName)

			// Run migrations before any other logic
			migrations, err := mm.RunMigrations(ctx, cfg, fromVM)
			if err != nil {
				return nil, err
			}

			wctx.Logger().Info("upgrade completed", "duration_ms", time.Since(startTime).Milliseconds())
			return migrations, nil
		}
	},
	// Add any new modules or store upgrades here
	StoreUpgrades: storetypes.StoreUpgrades{
		Added: []string{},
	},
}
