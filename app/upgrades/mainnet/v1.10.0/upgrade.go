package v1100

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/trstlabs/intento/app/upgrades"
)

const UpgradeName = "v1.10.0"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        storetypes.StoreUpgrades{},
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers upgrades.IntentoKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// Run DeICS logic
		// We pass nil for readyValopers so that logic relies on default behavior (filtering only jails etc, checking allowlist if it was non-nil)
		// effectively trying to migrate all governors that meet criteria.
		// NOTE: DeICS implementation checks: if readyValopers != nil && !readyValopers[valoperStr].
		// So if readyValopers is nil, it skips that check, allowing all non-jailed governors to migrate.
		err := DeICS(
			sdk.UnwrapSDKContext(ctx),
			*keepers.StakingKeeper,
			keepers.ConsumerKeeper,
			*keepers.StakingKeeper,
			keepers.BankKeeper,
			nil,
		)
		if err != nil {
			return nil, err
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
