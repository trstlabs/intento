package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/keeper"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	"github.com/trstlabs/intento/internal/collcompat"
	"github.com/trstlabs/intento/x/intent/types"
	interchainquerykeeper "github.com/trstlabs/intento/x/interchainquery/keeper"
)

type Keeper struct {
	cdc                   codec.Codec
	storeService          corestoretypes.KVStoreService
	Schema                collections.Schema
	scopedKeeper          capabilitykeeper.ScopedKeeper
	icaControllerKeeper   icacontrollerkeeper.Keeper
	bankKeeper            bankkeeper.Keeper
	distrKeeper           distrkeeper.Keeper
	stakingKeeper         stakingkeeper.Keeper
	transferKeeper        ibctransferkeeper.Keeper
	accountKeeper         authkeeper.AccountKeeper
	interchainQueryKeeper interchainquerykeeper.Keeper
	hooks                 IntentHooks
	msgRouter             MessageRouter
	interfaceRegistry     cdctypes.InterfaceRegistry
	Params                collections.Item[types.Params]
	authority             string
}

func NewKeeper(cdc codec.Codec, storeService corestoretypes.KVStoreService, icaKeeper icacontrollerkeeper.Keeper, scopedKeeper capabilitykeeper.ScopedKeeper, bankKeeper bankkeeper.Keeper, distrKeeper distrkeeper.Keeper, stakingKeeper stakingkeeper.Keeper, transferKeeper ibctransferkeeper.Keeper, accountKeeper authkeeper.AccountKeeper, interchainQueryKeeper interchainquerykeeper.Keeper, ah IntentHooks, msgRouter MessageRouter, interfaceRegistry cdctypes.InterfaceRegistry, authority string,
) Keeper {
	moduleAccAddr := accountKeeper.GetModuleAddress(types.ModuleName)
	// ensure module account is set
	if moduleAccAddr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	sb := collections.NewSchemaBuilder(storeService)

	keeper := Keeper{
		cdc:                   cdc,
		storeService:          storeService,
		scopedKeeper:          scopedKeeper,
		icaControllerKeeper:   icaKeeper,
		bankKeeper:            bankKeeper,
		distrKeeper:           distrKeeper,
		stakingKeeper:         stakingKeeper,
		transferKeeper:        transferKeeper,
		accountKeeper:         accountKeeper,
		interchainQueryKeeper: interchainQueryKeeper,
		hooks:                 ah,
		msgRouter:             msgRouter,
		interfaceRegistry:     interfaceRegistry,
		Params: collections.NewItem(
			sb,
			types.ParamsKey,
			"params",
			collcompat.ProtoValue[types.Params](cdc),
		),
		authority: authority,
	}
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	keeper.Schema = schema
	return keeper

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

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}
