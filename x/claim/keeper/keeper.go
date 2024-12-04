package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trstlabs/intento/internal/collcompat"
	"github.com/trstlabs/intento/x/claim/types"
)

// Keeper struct
type Keeper struct {
	cdc           codec.Codec
	storeService  corestoretypes.KVStoreService
	Schema        collections.Schema
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	stakingKeeper types.StakingKeeper
	distrKeeper   types.DistrKeeper
	params        collections.Item[types.Params]
	authority     string
}

// NewKeeper returns keeper
func NewKeeper(
	cdc codec.Codec,
	storeService corestoretypes.KVStoreService,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	sk types.StakingKeeper,
	dk types.DistrKeeper,
	authority string) Keeper {

	sb := collections.NewSchemaBuilder(storeService)
	keeper := Keeper{
		cdc:           cdc,
		storeService:  storeService,
		accountKeeper: ak,
		bankKeeper:    bk,
		stakingKeeper: sk,
		distrKeeper:   dk,
		authority:     authority,
		params: collections.NewItem(
			sb,
			types.ParamsKey,
			"params",
			collcompat.ProtoValue[types.Params](cdc),
		),
	}
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	keeper.Schema = schema
	return keeper

}

// Logger returns logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
