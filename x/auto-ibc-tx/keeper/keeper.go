package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/keeper"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/trstlabs/trst/x/auto-ibc-tx/types"

	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

type Keeper struct {
	cdc codec.Codec

	storeKey sdk.StoreKey

	scopedKeeper        capabilitykeeper.ScopedKeeper
	icaControllerKeeper icacontrollerkeeper.Keeper
	bankKeeper          bankkeeper.Keeper
	distrKeeper         distrkeeper.Keeper
	stakingKeeper       stakingkeeper.Keeper
	accountKeeper       authkeeper.AccountKeeper
	paramSpace          paramtypes.Subspace
}

func NewKeeper(cdc codec.Codec, storeKey sdk.StoreKey, iaKeeper icacontrollerkeeper.Keeper, scopedKeeper capabilitykeeper.ScopedKeeper, bankKeeper bankkeeper.Keeper, distrKeeper distrkeeper.Keeper, stakingKeeper stakingkeeper.Keeper, accountKeeper authkeeper.AccountKeeper, paramSpace paramtypes.Subspace) Keeper {
	moduleAccAddr := accountKeeper.GetModuleAddress(types.ModuleName)
	// ensure module account is set
	if moduleAccAddr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(ParamKeyTable())
	}

	return Keeper{
		cdc:                 cdc,
		storeKey:            storeKey,
		scopedKeeper:        scopedKeeper,
		icaControllerKeeper: iaKeeper,
		paramSpace:          paramSpace,
		bankKeeper:          bankKeeper,
		distrKeeper:         distrKeeper,
		stakingKeeper:       stakingKeeper,
		accountKeeper:       accountKeeper,
	}
}

// ClaimCapability claims the channel capability passed via the OnOpenChanInit callback
func (k *Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

// RegisterInterchainAccount registers account
func (k Keeper) RegisterInterchainAccount(ctx sdk.Context, connectionId, owner string) error {
	if err := k.icaControllerKeeper.RegisterInterchainAccount(ctx, connectionId, owner); err != nil {
		return err
	}
	return nil
}

// Logger returns the application logger, scoped to the associated module
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s-%s", host.ModuleName, types.ModuleName))
}
