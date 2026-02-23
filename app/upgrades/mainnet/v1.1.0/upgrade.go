package v1100

import (
	"context"

	"cosmossdk.io/math"
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
			keepers.ConsumerKeeper,
			readyValopers,
		)
		if err != nil {
			return nil, err
		}

		// Update Slashing Params
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		slashingParams, err := keepers.SlashingKeeper.GetParams(sdkCtx)
		if err != nil {
			return nil, err
		}
		slashingParams.MinSignedPerWindow = math.LegacyNewDecWithPrec(5, 1) // 50%
		if err := keepers.SlashingKeeper.SetParams(sdkCtx, slashingParams); err != nil {
			return nil, err
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
