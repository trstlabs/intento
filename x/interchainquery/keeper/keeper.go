package keeper

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	"github.com/trstlabs/intento/x/interchainquery/types"
)

// Keeper of this module maintains collections of registered zones.
type Keeper struct {
	cdc          codec.Codec
	storeService corestoretypes.KVStoreService
	callbacks    map[string]types.QueryCallbacks
	IBCKeeper    *ibckeeper.Keeper
	authority    string
}

// NewKeeper returns a new instance of zones Keeper
func NewKeeper(cdc codec.Codec, storeService corestoretypes.KVStoreService, ibckeeper *ibckeeper.Keeper) Keeper {
	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		callbacks:    make(map[string]types.QueryCallbacks),
		IBCKeeper:    ibckeeper,
	}
}

func (k *Keeper) SetCallbackHandler(module string, handler types.QueryCallbacks) error {
	_, found := k.callbacks[module]
	if found {
		return fmt.Errorf("callback handler already set for %s", module)
	}
	k.callbacks[module] = handler.RegisterICQCallbacks()
	return nil
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) SubmitICQRequest(ctx sdk.Context, query types.Query) error {
	k.Logger(ctx).Info("Submitting ICQ Request", "callbackId", query.CallbackId, "connectionId", query.ConnectionId, "queryType", query.QueryType, "timeout_duration", query.TimeoutDuration)

	if err := k.ValidateQuery(ctx, query); err != nil {
		return err
	}

	// Set the timeout using the block time and timeout duration
	timeoutTimestamp := uint64(ctx.BlockTime().UnixNano() + query.TimeoutDuration.Nanoseconds())
	query.TimeoutTimestamp = timeoutTimestamp

	query.RequestSent = false

	// Set the submission height on the Query to the latest light client height
	// In the query response, this will be used to verify that the query wasn't historical
	connection, found := k.IBCKeeper.ConnectionKeeper.GetConnection(ctx, query.ConnectionId)
	if !found {
		return errorsmod.Wrap(connectiontypes.ErrConnectionNotFound, query.ConnectionId)
	}
	clientState, found := k.IBCKeeper.ClientKeeper.GetClientState(ctx, connection.ClientId)
	if !found {
		return errorsmod.Wrap(clienttypes.ErrClientNotFound, connection.ClientId)
	}
	query.SubmissionHeight = clientState.GetLatestHeight().GetRevisionHeight()

	// Save the query to the store
	// If the same query is re-requested, it will get replace in the store with an updated TTL
	//  and the RequestSent bool reset to false
	k.SetQuery(ctx, query)

	return nil
}

// Re-submit an ICQ, generally used after a timeout
func (k *Keeper) RetryICQRequest(ctx sdk.Context, query types.Query) error {
	k.Logger(ctx).Info("Queuing ICQ Retry", "Query Type", query.CallbackId, "Query ID", query.Id)

	// Delete old query
	k.DeleteQuery(ctx, query.Id)

	//custom logic: we retry once, and then it should be rejected to prevent looping.
	query.TimeoutPolicy = types.TimeoutPolicy_REJECT_QUERY_RESPONSE
	// Submit a new query
	if err := k.SubmitICQRequest(ctx, query); err != nil {
		return errorsmod.Wrap(err, types.ErrFailedToRetryQuery.Error())
	}

	return nil
}

// Returns a log string with a chain Id and callback as a prefix
// callbackType is either ICACALLBACK or ICQCALLBACK
// Format:
//
//	|   CHAIN-ID    |  {CALLBACK_ID} {CALLBACK_TYPE}  |  string
func logCallbackWithHostChain(chainId string, callbackId string, callbackType string, s string, a ...any) string {
	msg := fmt.Sprintf(s, a...)
	return fmt.Sprintf("|   %-13s |  %s %s  |  %s", strings.ToUpper(chainId), strings.ToUpper(callbackId), callbackType, msg)
}

// Returns a log string with a chain Id and icacallback as a prefix
// Ex:
//
//	| COSMOSHUB-4   |  DELEGATE ICACALLBACK  |  string
func LogICQCallbackWithHostChain(chainId string, callbackId string, s string, a ...any) string {
	return logCallbackWithHostChain(chainId, callbackId, "ICQCALLBACK", s, a...)
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}
