package keeper

import (

	//"encoding/json"
	"fmt"

	//"log"
	"path/filepath"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	mintkeeper "github.com/trstlabs/trst/x/mint/keeper"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
	wasm "github.com/trstlabs/trst/go-cosmwasm"

	"github.com/trstlabs/trst/x/compute/internal/types"
)

// Keeper will have a reference to Wasmer with it's own data directory.
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      codec.BinaryCodec
	//legacyAmino   codec.LegacyAmino
	accountKeeper authkeeper.AccountKeeper
	bankKeeper    bankkeeper.Keeper
	distrKeeper   distrkeeper.Keeper
	wasmer        wasm.Wasmer
	queryPlugins  QueryPlugins
	messenger     MessageHandler
	// queryGasLimit is the max wasm gas that can be spent on executing a query with a contract
	queryGasLimit uint64
	paramSpace    paramtypes.Subspace
	hooks         ComputeHooks

	// authZPolicy   AuthorizationPolicy

}

// NewKeeper creates a new contract Keeper instance
// If customEncoders is non-nil, we can use this to override some of the message handler, especially custom
func NewKeeper(cdc codec.BinaryCodec /*legacyAmino codec.LegacyAmino,*/, storeKey sdk.StoreKey, accountKeeper authkeeper.AccountKeeper, bankKeeper bankkeeper.Keeper /*(govKeeper govkeeper.Keeper,*/, distKeeper distrkeeper.Keeper, mintKeeper mintkeeper.Keeper, stakingKeeper stakingkeeper.Keeper,
	router sdk.Router, homeDir string, wasmConfig *types.WasmConfig, supportedFeatures string, customEncoders *MessageEncoders, customPlugins *QueryPlugins, paramSpace paramtypes.Subspace, gh ComputeHooks) Keeper {
	wasmer, err := wasm.NewWasmer(filepath.Join(homeDir, "wasm"), supportedFeatures, wasmConfig.CacheSize, wasmConfig.EnclaveCacheSize)
	if err != nil {
		panic(err)
	}

	addr := accountKeeper.GetModuleAddress(types.ModuleName)
	// ensure module account is set
	if addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(ParamKeyTable())
	}

	keeper := Keeper{
		storeKey: storeKey,
		cdc:      cdc,
		//legacyAmino:   legacyAmino,
		wasmer:        *wasmer,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		distrKeeper:   distKeeper,
		messenger:     NewMessageHandler(router, customEncoders),
		queryGasLimit: wasmConfig.SmartQueryGasLimit,
		paramSpace:    paramSpace,
		hooks:         gh,

		// authZPolicy:   DefaultAuthorizationPolicy{},

	}
	keeper.queryPlugins = DefaultQueryPlugins( /*govKeeper,*/ distKeeper, mintKeeper, bankKeeper, stakingKeeper, &keeper).Merge(customPlugins)
	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
