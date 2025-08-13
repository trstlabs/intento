package keeper

import (
	"context"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"

	"github.com/trstlabs/intento/x/intent/types"
)

var _ types.QueryServer = &QueryServer{}

// QueryServer implements the module gRPC query service.
type QueryServer struct {
	keeper Keeper
}

// NewQueryServer creates a new gRPC query server.
func NewQueryServer(keeper Keeper) *QueryServer {
	return &QueryServer{
		keeper: keeper,
	}
}

// Params returns params of the alloc module.
func (q QueryServer) Params(goCtx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params, err := q.keeper.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, err
}

// InterchainAccountFromAddress implements the Query/InterchainAccountFromAddress gRPC method
func (q QueryServer) InterchainAccountFromAddress(goCtx context.Context, req *types.QueryInterchainAccountFromAddressRequest) (*types.QueryInterchainAccountFromAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	portID, err := icatypes.NewControllerPortID(req.Owner)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "could not find account: %s", err)
	}

	ica, found := q.keeper.icaControllerKeeper.GetInterchainAccountAddress(ctx, req.ConnectionId, portID)
	if !found {
		return nil, status.Errorf(codes.NotFound, "no account found for connectionID: %s and portID: %s",
			req.ConnectionId, portID)
	}

	return types.NewQueryInterchainAccountResponse(ica), nil
}

// Flow implements the Query/FlowgRPC method
func (q QueryServer) Flow(c context.Context, req *types.QueryFlowRequest) (*types.QueryFlowResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	id, err := strconv.ParseUint(req.Id, 10, 64)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	flow, err := q.keeper.TryGetflow(ctx, id)
	if err != nil {
		return nil, err
	}
	// for msg := range flow.Msgs{
	// 	makeReadableMsgData(&flow, msg)
	// }

	return &types.QueryFlowResponse{
		Flow: flow,
	}, nil
}

func (q QueryServer) FlowHistory(ctx context.Context, req *types.QueryFlowHistoryRequest) (*types.QueryFlowHistoryResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	// Convert id from the request
	id, err := strconv.ParseUint(req.Id, 10, 64)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request")
	}

	// Assuming FlowHistoryEntry items are stored with keys prefixed by id

	store := runtime.KVStoreAdapter(q.keeper.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.GetFlowHistoryKey(id))
	// Paginate over the prefixed store
	var historyEntries []types.FlowHistoryEntry
	pageRes, err := query.Paginate(prefixStore, req.Pagination, func(key []byte, value []byte) error {
		var historyEntry types.FlowHistoryEntry
		if err := q.keeper.cdc.Unmarshal(value, &historyEntry); err != nil {
			return err
		}
		historyEntries = append(historyEntries, historyEntry)
		return nil
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "paginate: %v", err)
	}

	// Return paginated results
	return &types.QueryFlowHistoryResponse{
		History:    historyEntries,
		Pagination: pageRes,
	}, nil
}

// Flows implements the Query/Flows gRPC method
func (q QueryServer) Flows(c context.Context, req *types.QueryFlowsRequest) (*types.QueryFlowsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	flows := make([]types.Flow, 0)
	store := runtime.KVStoreAdapter(q.keeper.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.FlowKeyPrefix)

	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(_ []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			var c types.Flow
			q.keeper.cdc.MustUnmarshal(value, &c)
			flows = append(flows, c)

		}
		return true, nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryFlowsResponse{
		Flows:      flows,
		Pagination: pageRes,
	}, nil
}

// FlowsForOwner implements the Query/FlowsForOwner gRPC method
func (q QueryServer) FlowsForOwner(c context.Context, req *types.QueryFlowsForOwnerRequest) (*types.QueryFlowsForOwnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	flows := make([]types.Flow, 0)

	ownerAddress, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, err
	}
	store := runtime.KVStoreAdapter(q.keeper.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.GetFlowsByOwnerPrefix(ownerAddress))

	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(key []byte, _ []byte, accumulate bool) (bool, error) {
		if accumulate {
			flowID := types.GetIDFromBytes(key)
			flow := q.keeper.Getflow(ctx, flowID)

			flows = append(flows, flow)

		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryFlowsForOwnerResponse{
		Flows:      flows,
		Pagination: pageRes,
	}, nil
}

// TrustlessAgent implements the Query/TrustlessAgent gRPC method
func (q QueryServer) TrustlessAgent(c context.Context, req *types.QueryTrustlessAgentRequest) (*types.QueryTrustlessAgentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	hosted, err := q.keeper.TryGetTrustlessAgent(ctx, req.AgentAddress)
	if err != nil {
		return nil, err
	}

	return &types.QueryTrustlessAgentResponse{
		TrustlessAgent: hosted,
	}, nil
}

// TrustlessAgents implements the Query/TrustlessAgents gRPC method
func (q QueryServer) TrustlessAgents(c context.Context, req *types.QueryTrustlessAgentsRequest) (*types.QueryTrustlessAgentsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	trustlessAgents := make([]types.TrustlessAgent, 0)
	store := runtime.KVStoreAdapter(q.keeper.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.TrustlessAgentKeyPrefix)

	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(_ []byte, value []byte, accumulate bool) (bool, error) {
		if accumulate {
			var c types.TrustlessAgent
			q.keeper.cdc.MustUnmarshal(value, &c)
			trustlessAgents = append(trustlessAgents, c)

		}
		return true, nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryTrustlessAgentsResponse{
		TrustlessAgents: trustlessAgents,
		Pagination:      pageRes,
	}, nil
}

// TrustlessAgentsByFeeAdmin implements the Query/TrustlessExTrustlessAgentsByFeeAdminecutionAgentsByAdmin gRPC method
func (q QueryServer) TrustlessAgentsByFeeAdmin(c context.Context, req *types.QueryTrustlessAgentsByFeeAdminRequest) (*types.QueryTrustlessAgentsByFeeAdminResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	trustlessAgents := make([]types.TrustlessAgent, 0)

	admin, err := sdk.AccAddressFromBech32(req.FeeAdmin)
	if err != nil {
		return nil, err
	}
	store := runtime.KVStoreAdapter(q.keeper.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.GetTrustlessAgentsByAdminPrefix(admin))
	pageRes, err := query.FilteredPaginate(prefixStore, req.Pagination, func(key []byte, _ []byte, accumulate bool) (bool, error) {
		if accumulate {
			trustlessAgentAddress := string(key)
			flow := q.keeper.GetTrustlessAgent(ctx, trustlessAgentAddress)

			trustlessAgents = append(trustlessAgents, flow)

		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryTrustlessAgentsByFeeAdminResponse{
		TrustlessAgents: trustlessAgents,
		Pagination:      pageRes,
	}, nil
}
