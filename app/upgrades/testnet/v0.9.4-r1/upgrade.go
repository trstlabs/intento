package v094r1

import (
	"context"
	"time"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	upgrades "github.com/trstlabs/intento/app/upgrades"
)

// next upgrade name
const UpgradeName = "v094r1"

var Upgrade = upgrades.Upgrade{
	UpgradeName: UpgradeName,
	CreateUpgradeHandler: func(mm *module.Manager, cfg module.Configurator, _ upgrades.IntentoKeepers) upgradetypes.UpgradeHandler {
		return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			startTime := time.Now()
			wctx := sdk.UnwrapSDKContext(ctx)
			wctx.Logger().Info("upgrade started", "upgrade_name", UpgradeName)
			migrations, err := mm.RunMigrations(ctx, cfg, fromVM)
			if err != nil {
				return nil, err
			}
			wctx.Logger().Info("upgrade completed", "duration_ms", time.Since(startTime).Milliseconds())
			return migrations, nil
		}
	},
	StoreUpgrades: storetypes.StoreUpgrades{
		Added: []string{},
	},
}
