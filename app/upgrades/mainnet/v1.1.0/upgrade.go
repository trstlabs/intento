package v1100

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/trstlabs/intento/app/upgrades"
)

const UpgradeName = "v1.1.0"

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
		// Load ready validators from file
		readyValopers, err := GetReadyValidators()
		if err != nil {
			return nil, err
		}

		// Run DeICS logic
		err = DeICS(
			sdk.UnwrapSDKContext(ctx),
			*keepers.StakingKeeper,
			keepers.BankKeeper,
			readyValopers,
		)
		if err != nil {
			return nil, err
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
