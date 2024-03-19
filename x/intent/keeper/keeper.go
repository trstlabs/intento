package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	"github.com/trstlabs/intento/x/intent/types"

	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
)

type Keeper struct {
	cdc                 codec.Codec
	storeKey            storetypes.StoreKey
	scopedKeeper        capabilitykeeper.ScopedKeeper
	icaControllerKeeper icacontrollerkeeper.Keeper
	bankKeeper          bankkeeper.Keeper
	distrKeeper         distrkeeper.Keeper
	stakingKeeper       stakingkeeper.Keeper
	transferKeeper      ibctransferkeeper.Keeper
	accountKeeper       authkeeper.AccountKeeper
	paramSpace          paramtypes.Subspace
	hooks               IntentHooks
	msgRouter           MessageRouter
}

func NewKeeper(cdc codec.Codec, storeKey storetypes.StoreKey, icaKeeper icacontrollerkeeper.Keeper, scopedKeeper capabilitykeeper.ScopedKeeper, bankKeeper bankkeeper.Keeper, distrKeeper distrkeeper.Keeper, stakingKeeper stakingkeeper.Keeper, transferKeeper ibctransferkeeper.Keeper, accountKeeper authkeeper.AccountKeeper, paramSpace paramtypes.Subspace, ah IntentHooks, msgRouter MessageRouter) Keeper {
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
		icaControllerKeeper: icaKeeper,
		paramSpace:          paramSpace,
		bankKeeper:          bankKeeper,
		distrKeeper:         distrKeeper,
		stakingKeeper:       stakingKeeper,
		transferKeeper:      transferKeeper,
		accountKeeper:       accountKeeper,
		hooks:               ah,
		msgRouter:           msgRouter,
	}
}

// RegisterInterchainAccount registers account
func (k Keeper) RegisterInterchainAccount(ctx sdk.Context, connectionId, owner, version string) error {
	if err := k.icaControllerKeeper.RegisterInterchainAccount(ctx, connectionId, owner, version); err != nil {
		return err
	}
	return nil
}

// Logger returns the application logger, scoped to the associated module
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// MessageRouter ADR 031 request type routing
type MessageRouter interface {
	Handler(msg sdk.Msg) baseapp.MsgServiceHandler
}
