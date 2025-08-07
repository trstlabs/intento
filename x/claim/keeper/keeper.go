package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	corestoretypes "cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gogoproto/proto"
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
	Params        collections.Item[types.Params]
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
		Params: collections.NewItem(
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

func (k Keeper) CreateModuleAccount(ctx sdk.Context, amount sdk.Coin) {
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(amount))
	if err != nil {
		panic(err)
	}
}

// SetClaimRecords sets claim records
func (k Keeper) SetClaimRecords(ctx sdk.Context, claimRecords []types.ClaimRecord) error {
	for _, claimRecord := range claimRecords {
		err := k.SetClaimRecord(ctx, claimRecord)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetClaimRecords get claim records for genesis export
func (k Keeper) GetClaimRecords(ctx sdk.Context) []types.ClaimRecord {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, []byte(types.ClaimRecordsStorePrefix))

	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	claimRecords := []types.ClaimRecord{}
	for ; iterator.Valid(); iterator.Next() {

		claimRecord := types.ClaimRecord{}

		err := proto.Unmarshal(iterator.Value(), &claimRecord)
		if err != nil {
			panic(err)
		}

		claimRecords = append(claimRecords, claimRecord)
	}
	return claimRecords
}

// GetClaimRecord returns the claim record for a specific address
func (k Keeper) GetClaimRecord(ctx sdk.Context, addr sdk.AccAddress) (types.ClaimRecord, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, []byte(types.ClaimRecordsStorePrefix))
	if !prefixStore.Has(addr) {
		return types.ClaimRecord{}, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "address does not have claim record")
	}
	bz := prefixStore.Get(addr)

	claimRecord := types.ClaimRecord{}
	err := proto.Unmarshal(bz, &claimRecord)
	if err != nil {
		return types.ClaimRecord{}, err
	}
	if len(claimRecord.Status) == 0 {
		return types.ClaimRecord{}, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "address does not have claim record")
	}

	return claimRecord, nil
}

// SetClaimRecord sets a claim record for an address in store
func (k Keeper) SetClaimRecord(ctx sdk.Context, claimRecord types.ClaimRecord) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, []byte(types.ClaimRecordsStorePrefix))

	bz, err := proto.Marshal(&claimRecord)
	if err != nil {
		return err
	}

	addr, err := sdk.AccAddressFromBech32(claimRecord.Address)
	if err != nil {
		return err
	}

	prefixStore.Set(addr, bz)
	return nil
}

// clearInitialClaimables clear claimable amounts
func (k Keeper) clearInitialClaimables(ctx sdk.Context) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, []byte(types.ClaimRecordsStorePrefix))
	iterator := prefixStore.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		store.Delete(key)
	}
}

// GetAuthority returns the x/alloc module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}
